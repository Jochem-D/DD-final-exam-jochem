package cli

import (
	"context"
	"ddsheetfinal/internal/application/usecases"
	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/valueobjects"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ServeHandler handles the HTTP server command
type ServeHandler struct {
	serveUseCase *usecases.ServeHTTPUseCase
}

// NewServeHandler creates a new serve handler
func NewServeHandler(serveUseCase *usecases.ServeHTTPUseCase) *ServeHandler {
	return &ServeHandler{
		serveUseCase: serveUseCase,
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

	webDir := filepath.Join(projectRoot, "web", "static")
	dataDir := filepath.Join(projectRoot, "data")
	charDir := filepath.Join(dataDir, "characters")

	mux := http.NewServeMux()

	// Serve static frontend files
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	// Dynamic manifest listing all character JSON files (at both paths for compatibility)
	mux.HandleFunc("/characters/manifest.json", h.handleManifest(charDir))
	mux.HandleFunc("/frontend/manifest.json", h.handleManifest(charDir))

	// Character CRUD endpoints
	mux.HandleFunc("/characters/", h.handleCharacters(charDir))

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
			if strings.HasSuffix(strings.ToLower(name), ".json") {
				out = append(out, "/characters/"+name)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	}
}

// handleCharacters returns a handler for GET/PUT/DELETE on character files
func (h *ServeHandler) handleCharacters(charDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rel := strings.TrimPrefix(r.URL.Path, "/characters/")
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
		target := filepath.Join(charDir, filepath.Clean(rel))

		switch r.Method {
		case http.MethodGet:
			http.ServeFile(w, r, target)
			return

		case http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "read error", http.StatusBadRequest)
				return
			}
			log.Printf("PUT %s (%d bytes)", target, len(body))
			
			// Validate JSON
			var char entities.Character
			if err := json.Unmarshal(body, &char); err != nil {
				http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
				return
			}
			
			// Ensure characters dir exists
			if err := os.MkdirAll(charDir, 0755); err != nil {
				http.Error(w, "cannot create characters dir", http.StatusInternalServerError)
				return
			}
			
			if err := os.WriteFile(target, body, 0644); err != nil {
				http.Error(w, "save failed", http.StatusInternalServerError)
				return
			}
			
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return

		case http.MethodDelete:
			if err := os.Remove(target); err != nil {
				if os.IsNotExist(err) {
					http.Error(w, "not found", http.StatusNotFound)
					return
				}
				http.Error(w, "delete failed", http.StatusInternalServerError)
				return
			}
			log.Printf("DELETE %s OK", target)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

// handleDerive returns a handler that computes derived stats (AC, saves, etc.)
func (h *ServeHandler) handleDerive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var char entities.Character
		if err := json.NewDecoder(r.Body).Decode(&char); err != nil {
			http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
			return
		}

	// Calculate derived stats using the character service
	service := h.serveUseCase.GetCharacterService()
	
	ac, acDesc := service.ComputeArmorClass(&char)
	pb := valueobjects.ProficiencyBonus(char.Level)
	
	// Build ability modifiers map
	abilityMods := map[string]int{
		"Strength":     valueobjects.AbilityModifier(char.Str),
		"Dexterity":    valueobjects.AbilityModifier(char.Dex),
		"Constitution": valueobjects.AbilityModifier(char.Con),
		"Intelligence": valueobjects.AbilityModifier(char.Int),
		"Wisdom":       valueobjects.AbilityModifier(char.Wis),
		"Charisma":     valueobjects.AbilityModifier(char.Cha),
	}
	
	// Build saving throws map
	saves := make(map[string]int)
	for ability, value := range map[string]int{
		"Strength":     char.Str,
		"Dexterity":    char.Dex,
		"Constitution": char.Con,
		"Intelligence": char.Int,
		"Wisdom":       char.Wis,
		"Charisma":     char.Cha,
	} {
		modifier := valueobjects.AbilityModifier(value)
		if service.IsSaveProficient(&char, ability) {
			saves[ability] = modifier + pb
		} else {
			saves[ability] = modifier
		}
	}
	
	// Return snake_case keys to match frontend expectations
	result := map[string]interface{}{
		"armor_class":        ac,
		"armor_class_desc":   acDesc,
		"initiative":         service.ComputeInitiativeBonus(&char),
		"passive_perception": service.ComputePassivePerception(&char),
		"proficiency_bonus":  pb,
		"ability_mods":       abilityMods,
		"saving_throws":      saves,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	}
}

// handleEnrich returns a handler that fetches spell/equipment data from dnd5eapi
func (h *ServeHandler) handleEnrich() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimSpace(r.URL.Query().Get("name"))
		if name == "" {
			http.Error(w, "missing ?name= parameter", http.StatusBadRequest)
			return
		}

		start := time.Now()
		ctx := context.Background()

		// Load character
		charRepo := h.serveUseCase.GetCharacterRepository()
		char, err := charRepo.FindByName(name)
		if err != nil {
			http.Error(w, "character not found: "+err.Error(), http.StatusNotFound)
			return
		}

		// Enrich spells and equipment
		enrichRepo := h.serveUseCase.GetEnrichmentRepository()
		
		enrichedSpells := make(map[string]interface{})
		for _, spell := range char.Spells {
			spellInfo, err := enrichRepo.FetchSpellInfo(ctx, spell)
			if err == nil && spellInfo != nil {
				enrichedSpells[spell] = spellInfo
			}
		}

		enrichedEquipment := make(map[string]interface{})
		// Equipment is stored in Weapon, OffHand, Armor, Shield, and Inventory
		equipmentItems := []string{}
		if char.Weapon != "" {
			equipmentItems = append(equipmentItems, char.Weapon)
		}
		if char.OffHand != "" {
			equipmentItems = append(equipmentItems, char.OffHand)
		}
		if char.Armor != "" {
			equipmentItems = append(equipmentItems, char.Armor)
		}
		if char.Shield != "" {
			equipmentItems = append(equipmentItems, char.Shield)
		}
		equipmentItems = append(equipmentItems, char.Inventory...)
		
		for _, equip := range equipmentItems {
			if equip == "" {
				continue
			}
			equipInfo, err := enrichRepo.FetchEquipmentInfo(ctx, equip)
			if err == nil && equipInfo != nil {
				enrichedEquipment[equip] = equipInfo
			}
		}

		result := map[string]interface{}{
			"OK":             true,
			"CharacterName":  name,
			"Spells":         enrichedSpells,
			"Equipment":      enrichedEquipment,
			"FetchedSpells":  len(enrichedSpells),
			"FetchedEquip":   len(enrichedEquipment),
			"ElapsedMs":      time.Since(start).Milliseconds(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
