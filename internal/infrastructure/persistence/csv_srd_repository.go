package persistence

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strings"

	"ddsheetfinal/internal/domain/repositories"
)

const (
	errFailedReadSpellsCSV = "failed to read spells CSV: %w"
	errFailedOpenSpellsCSV = "failed to open spells CSV: %w"
)

// CSVSRDRepository implements SRDRepository using CSV files
type CSVSRDRepository struct {
	equipmentCSVPath string
	spellsCSVPath    string
}

// NewCSVSRDRepository creates a new CSV-based SRD repository
func NewCSVSRDRepository(equipmentCSVPath, spellsCSVPath string) repositories.SRDRepository {
	return &CSVSRDRepository{
		equipmentCSVPath: equipmentCSVPath,
		spellsCSVPath:    spellsCSVPath,
	}
}

// IsValidSpell checks if a spell exists in the SRD
func (r *CSVSRDRepository) IsValidSpell(spellName string) bool {
	key := r.canonKey(spellName)
	
	f, err := os.Open(r.spellsCSVPath)
	if err != nil {
		return false
	}
	defer f.Close()
	
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}
	
	for _, rec := range records {
		if len(rec) == 0 {
			continue
		}
		name := r.canonKey(rec[0])
		if name == key {
			return true
		}
	}
	
	return false
}

// IsValidEquipment checks if equipment exists in the SRD
func (r *CSVSRDRepository) IsValidEquipment(equipmentName string) bool {
	key := r.canonKey(equipmentName)
	cands := r.equipAlternates(key)
	
	f, err := os.Open(r.equipmentCSVPath)
	if err != nil {
		return false
	}
	defer f.Close()
	
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}
	
	for _, rec := range records {
		if len(rec) == 0 {
			continue
		}
		name := r.canonKey(rec[0])
		// direct match
		if name == key {
			return true
		}
		// match any candidate
		if slices.Contains(cands, name) {
			return true
		}
	}
	
	return false
}

// IsSpellForClass checks if a spell is available to a class
func (r *CSVSRDRepository) IsSpellForClass(spellName, className string) (bool, error) {
	key := r.canonKey(spellName)
	classLower := strings.ToLower(strings.TrimSpace(className))
	
	f, err := os.Open(r.spellsCSVPath)
	if err != nil {
		return false, fmt.Errorf(errFailedOpenSpellsCSV, err)
	}
	defer f.Close()
	
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return false, fmt.Errorf(errFailedReadSpellsCSV, err)
	}
	
	// CSV format: name,level,class
	// Find the spell and check if the class is in its class list
	for i, rec := range records {
		if i == 0 || len(rec) < 3 {
			continue // skip header or invalid rows
		}
		name := r.canonKey(rec[0])
		if name == key {
			// Classes are in column 2 (index 2)
			classList := strings.ToLower(rec[2])
			return strings.Contains(classList, classLower), nil
		}
	}
	
	return false, nil
}

// GetLearnableSpells returns spells a character can learn
func (r *CSVSRDRepository) GetLearnableSpells(className string, knownSpells []string) ([]string, error) {
	classLower := strings.ToLower(strings.TrimSpace(className))
	knownMap := make(map[string]bool)
	for _, s := range knownSpells {
		knownMap[r.canonKey(s)] = true
	}
	
	f, err := os.Open(r.spellsCSVPath)
	if err != nil {
		return nil, fmt.Errorf(errFailedOpenSpellsCSV, err)
	}
	defer f.Close()
	
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf(errFailedReadSpellsCSV, err)
	}
	
	// CSV format: name,level,class
	var learnable []string
	for i, rec := range records {
		if i == 0 || len(rec) < 3 {
			continue // skip header or invalid rows
		}
		name := rec[0]
		key := r.canonKey(name)
		
		// Skip if already known
		if knownMap[key] {
			continue
		}
		
		// Check if available to class (column 2)
		classList := strings.ToLower(rec[2])
		if strings.Contains(classList, classLower) {
			learnable = append(learnable, name)
		}
	}
	
	return learnable, nil
}

// canonKey normalizes a name to a consistent lookup key
func (r *CSVSRDRepository) canonKey(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.NewReplacer("'", "'", "–", "-", "—", "-", "/", " ").Replace(s)
	s = strings.Join(strings.Fields(s), " ") // collapse repeated spaces
	return s
}

// equipAlternates returns a list of alternate lookup keys for an equipment name
func (r *CSVSRDRepository) equipAlternates(key string) []string {
	out := []string{key}
	armorSuffix := " armor"
	
	// toggle " armor" suffix
	if strings.HasSuffix(key, armorSuffix) {
		out = append(out, strings.TrimSuffix(key, armorSuffix))
	} else {
		out = append(out, key+armorSuffix)
		out = append(out, key+" armor")
	}
	
	// common armor aliases
	switch key {
	case "leather":
		out = append(out, "leather armor")
	case "studded leather":
		out = append(out, "studded leather armor")
	case "scale mail":
		out = append(out, "scale mail armor")
	case "chain mail":
		out = append(out, "chain mail armor")
	case "ring mail":
		out = append(out, "ring mail armor")
	case "splint":
		out = append(out, "splint armor")
	case "plate":
		out = append(out, "plate armor")
	case "half plate":
		out = append(out, "half plate armor")
	}
	
	// de-duplicate
	seen := map[string]struct{}{}
	uniq := make([]string, 0, len(out))
	for _, k := range out {
		k = strings.TrimSpace(k)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		uniq = append(uniq, k)
	}
	return uniq
}

// GetSpellLevel returns the level of a spell (0-9)
// CSV format: name,level,class
func (r *CSVSRDRepository) GetSpellLevel(spellName string) (int, error) {
	key := r.canonKey(spellName)
	
	f, err := os.Open(r.spellsCSVPath)
	if err != nil {
		return 0, fmt.Errorf(errFailedOpenSpellsCSV, err)
	}
	defer f.Close()
	
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf(errFailedReadSpellsCSV, err)
	}
	
	for _, rec := range records {
		if len(rec) < 2 {
			continue
		}
		name := r.canonKey(rec[0])
		if name == key {
			// Parse level from column 1
			level := 0
			fmt.Sscanf(rec[1], "%d", &level)
			return level, nil
		}
	}
	
	return 0, fmt.Errorf("spell '%s' not found", spellName)
}
