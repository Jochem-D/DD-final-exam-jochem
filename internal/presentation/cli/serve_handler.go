package cli

import (
	"context"
	"ddsheetfinal/internal/application/dtos"
	"ddsheetfinal/internal/application/usecases"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const charactersPath = "/characters/"

const (
	jsonExtension   = ".json"
	jsonContentType = "application/json"
	headerContentType = "Content-Type"
	errMethodNotAllowed = "method not allowed"
)

// ServeHandler handles the HTTP server command
type ServeHandler struct {
	getCharacterUseCase    *usecases.GetCharacterUseCase
	saveCharacterUseCase   *usecases.SaveCharacterUseCase
	deleteCharacterUseCase *usecases.DeleteCharacterUseCase
	listCharactersUseCase  *usecases.ListCharactersUseCase
	deriveStatsUseCase     *usecases.DeriveCharacterStatsUseCase
	enrichDataUseCase      *usecases.EnrichDataUseCase
}

// NewServeHandler creates a new serve handler
func NewServeHandler(
	getCharacterUseCase *usecases.GetCharacterUseCase,
	saveCharacterUseCase *usecases.SaveCharacterUseCase,
	deleteCharacterUseCase *usecases.DeleteCharacterUseCase,
	listCharactersUseCase *usecases.ListCharactersUseCase,
	deriveStatsUseCase *usecases.DeriveCharacterStatsUseCase,
	enrichDataUseCase *usecases.EnrichDataUseCase,
) *ServeHandler {
	return &ServeHandler{
		getCharacterUseCase:    getCharacterUseCase,
		saveCharacterUseCase:   saveCharacterUseCase,
		deleteCharacterUseCase: deleteCharacterUseCase,
		listCharactersUseCase:  listCharactersUseCase,
		deriveStatsUseCase:     deriveStatsUseCase,
		enrichDataUseCase:      enrichDataUseCase,
	}
}

// Handle starts the HTTP server
func (h *ServeHandler) Handle(args []string) {
	port := "8080"
	if len(args) > 0 && (args[0] == "-port" || args[0] == "--port") && len(args) > 1 {
		port = args[1]
	}

	// Get the current working directory (should be project root)
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	dataPath := filepath.Join(projectRoot, "data")
	charDir := filepath.Join(dataPath, "characters")
	webPath := filepath.Join(projectRoot, "web", "static")

	mux := http.NewServeMux()

	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir(webPath))
	mux.Handle("/", fs)

	// Dynamic manifest listing all character JSON files (at both paths for compatibility)
	mux.HandleFunc(charactersPath+"manifest.json", h.handleManifest(charDir))
	mux.HandleFunc("/frontend/manifest.json", h.handleManifest(charDir))

	// Character CRUD endpoints
	mux.HandleFunc(charactersPath, h.handleCharacters(charDir))

	// API: Derive endpoint (compute AC, saves, etc.)
	mux.HandleFunc("/api/derive", h.handleDerive())

	// API: Enrich endpoint (fetch from dnd5eapi)
	mux.HandleFunc("/api/enrich", h.handleEnrich())

	addr := ":" + port
	log.Printf("Server running at http://localhost%s", addr)
	log.Printf("Open http://localhost%s/index.html to view characters", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// handleManifest returns a handler that lists all character JSON files
func (h *ServeHandler) handleManifest(charDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := os.ReadDir(charDir)
		if err != nil {
			http.Error(w, "cannot read characters directory", http.StatusInternalServerError)
			return
		}
		
		var out []string
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasSuffix(strings.ToLower(name), jsonExtension) {
				out = append(out, charactersPath+name)
			}
		}
		
		w.Header().Set(headerContentType, jsonContentType)
		json.NewEncoder(w).Encode(out)
	}
}
// handleCharacters returns a handler for GET/PUT/DELETE on character files
func (h *ServeHandler) handleCharacters(charDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rel := strings.TrimPrefix(r.URL.Path, charactersPath)
		rel = strings.TrimPrefix(rel, "/")
		if rel == "" || rel == "manifest.json" {
			http.Error(w, "missing filename", http.StatusBadRequest)
			return
		}
		// Prevent path traversal
		if strings.Contains(rel, "..") || strings.HasPrefix(rel, "/") {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		
		// Extract character name (remove .json extension)
		charName := strings.TrimSuffix(rel, jsonExtension)

		switch r.Method {
		case http.MethodGet:
			// Use GetCharacterUseCase
			charDTO, err := h.getCharacterUseCase.Execute(charName)
			if err != nil {
				http.Error(w, "character not found: "+err.Error(), http.StatusNotFound)
				return
			}
			
			w.Header().Set(headerContentType, jsonContentType)
			json.NewEncoder(w).Encode(charDTO)
			return

		case http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "read error", http.StatusBadRequest)
				return
			}
			log.Printf("PUT %s (%d bytes)", charName, len(body))
			
			// Validate JSON and parse into DTO
			var charDTO dtos.CharacterDTO
			if err := json.Unmarshal(body, &charDTO); err != nil {
				http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
				return
			}
			
			// Use SaveCharacterUseCase
			if err := h.saveCharacterUseCase.Execute(&charDTO); err != nil {
				http.Error(w, "save failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
			
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return

		case http.MethodDelete:
			// Use DeleteCharacterUseCase
			if err := h.deleteCharacterUseCase.Execute(charName); err != nil {
				if strings.Contains(err.Error(), "not found") {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				http.Error(w, "delete failed", http.StatusInternalServerError)
				return
			}
			log.Printf("DELETE %s OK", charName)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return

		default:
			http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	}
}

// handleDerive returns a handler that computes derived stats (AC, saves, etc.)
func (h *ServeHandler) handleDerive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			CharacterName string `json:"character_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Use DeriveCharacterStatsUseCase
		derivedStats, err := h.deriveStatsUseCase.Execute(req.CharacterName)
		if err != nil {
			http.Error(w, "error deriving stats: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return snake_case keys to match frontend expectations
		result := map[string]interface{}{
			"armor_class":        derivedStats.AC,
			"armor_class_desc":   derivedStats.ACCalculation,
			"initiative":         derivedStats.DexMod,
			"passive_perception": 10 + derivedStats.WisMod, // Simplified, should come from use case
			"proficiency_bonus":  derivedStats.ProficiencyBonus,
			"ability_mods": map[string]int{
				"Strength":     derivedStats.StrMod,
				"Dexterity":    derivedStats.DexMod,
				"Constitution": derivedStats.ConMod,
				"Intelligence": derivedStats.IntMod,
				"Wisdom":       derivedStats.WisMod,
				"Charisma":     derivedStats.ChaMod,
			},
			"saving_throws": map[string]int{
				"Strength":     derivedStats.StrSave,
				"Dexterity":    derivedStats.DexSave,
				"Constitution": derivedStats.ConSave,
				"Intelligence": derivedStats.IntSave,
				"Wisdom":       derivedStats.WisSave,
				"Charisma":     derivedStats.ChaSave,
			},
		}

		w.Header().Set(headerContentType, jsonContentType)
		json.NewEncoder(w).Encode(result)
	}
}

// handleEnrich returns a handler that fetches spell/equipment data from dnd5eapi
func (h *ServeHandler) handleEnrich() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimSpace(r.URL.Query().Get("name"))
		if name == "" {
			http.Error(w, "missing ?name= parameter", http.StatusBadRequest)
			return
		}

		start := time.Now()
		ctx := context.Background()

		// Load character using GetCharacterUseCase
		charDTO, err := h.getCharacterUseCase.Execute(name)
		if err != nil {
			http.Error(w, "character not found: "+err.Error(), http.StatusNotFound)
			return
		}

		// Enrich spells
		enrichedSpells := make(map[string]interface{})
		for _, spell := range charDTO.Spells {
			spellIndex := h.enrichDataUseCase.SanitizeIndex(spell)
			spellInfo, err := h.enrichDataUseCase.ExecuteSpell(ctx, spellIndex)
			if err == nil && spellInfo != nil {
				enrichedSpells[spell] = spellInfo
			}
		}

		// Enrich equipment
		enrichedEquipment := make(map[string]interface{})
		equipmentItems := []string{}
		if charDTO.Weapon != "" {
			equipmentItems = append(equipmentItems, charDTO.Weapon)
		}
		if charDTO.OffHand != "" {
			equipmentItems = append(equipmentItems, charDTO.OffHand)
		}
		if charDTO.Armor != "" {
			equipmentItems = append(equipmentItems, charDTO.Armor)
		}
		if charDTO.Shield != "" {
			equipmentItems = append(equipmentItems, charDTO.Shield)
		}
		equipmentItems = append(equipmentItems, charDTO.Inventory...)
		
		for _, equip := range equipmentItems {
			if equip == "" {
				continue
			}
			equipIndex := h.enrichDataUseCase.SanitizeIndex(equip)
			equipInfo, err := h.enrichDataUseCase.ExecuteEquipment(ctx, equipIndex)
			if err == nil && equipInfo != nil {
				enrichedEquipment[equip] = equipInfo
			}
		}

		// Build enriched character with embedded enrichment data
		enrichedChar := map[string]interface{}{
			"name":                charDTO.Name,
			"race":                charDTO.Race,
			"class":               charDTO.Class,
			"background":          charDTO.Background,
			"level":               charDTO.Level,
			"str":                 charDTO.Str,
			"dex":                 charDTO.Dex,
			"con":                 charDTO.Con,
			"int":                 charDTO.Int,
			"wis":                 charDTO.Wis,
			"cha":                 charDTO.Cha,
			"skill_proficiencies": charDTO.SkillProficiencies,
			"spells":              charDTO.Spells,
			"equipment":           charDTO.Inventory,
			"weapon":              charDTO.Weapon,
			"armor":               charDTO.Armor,
			"shield":              charDTO.Shield,
			"enriched": map[string]interface{}{
				"spells":    enrichedSpells,
				"equipment": enrichedEquipment,
			},
		}

		// Save enriched character to file
		projectRoot, _ := os.Getwd()
		enrichDir := filepath.Join(projectRoot, "data", "enrichments")
		os.MkdirAll(enrichDir, 0755)
		enrichPath := filepath.Join(enrichDir, name+jsonExtension)
		
		enrichedJSON, err := json.MarshalIndent(enrichedChar, "", "  ")
		if err == nil {
			os.WriteFile(enrichPath, enrichedJSON, 0644)
		}

		result := map[string]interface{}{
			"OK":             true,
			"CharacterName":  name,
			"Spells":         enrichedSpells,
			"Equipment":      enrichedEquipment,
			"FetchedSpells":  len(enrichedSpells),
			"FetchedEquip":   len(enrichedEquipment),
			"ElapsedMs":      time.Since(start).Milliseconds(),
			"output_path":    enrichPath,
		}

		w.Header().Set(headerContentType, jsonContentType)
		json.NewEncoder(w).Encode(result)
	}
}
