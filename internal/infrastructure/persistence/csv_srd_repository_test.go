package persistence

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	testDataDir              = "../../../assets/srd"
	spellsCSV                = "5e-SRD-Spells.csv"
	equipmentCSV             = "5e-SRD-Equipment.csv"
	nonExistentEquipmentPath = "/nonexistent/equipment.csv"
	nonExistentSpellsPath    = "/nonexistent/spells.csv"
	spellFireball            = "Fireball"
	spellMagicMissile        = "Magic Missile"
	errUnexpected            = "Unexpected error: %v"
	errExpectedFileNotExist  = "Expected error when CSV file doesn't exist"
)

func getTestCSVPaths(t *testing.T) (string, string) {
	t.Helper()
	
	// Get absolute paths to CSV files
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	
	equipPath := filepath.Join(cwd, testDataDir, equipmentCSV)
	spellPath := filepath.Join(cwd, testDataDir, spellsCSV)
	
	// Verify files exist
	if _, err := os.Stat(equipPath); err != nil {
		t.Fatalf("Equipment CSV not found at %s: %v", equipPath, err)
	}
	if _, err := os.Stat(spellPath); err != nil {
		t.Fatalf("Spells CSV not found at %s: %v", spellPath, err)
	}
	
	return equipPath, spellPath
}

// ========== IsValidSpell Tests ==========

func TestIsValidSpellSuccess(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	tests := []struct {
		name      string
		spellName string
	}{
		{"exact match", spellFireball},
		{"lowercase", "fireball"},
		{"uppercase", "FIREBALL"},
		{"with spaces", " Fireball "},
		{"cantrip", "Acid Splash"},
		{"level 1", spellMagicMissile},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := repo.IsValidSpell(tt.spellName)
			if !valid {
				t.Errorf("Expected '%s' to be valid spell", tt.spellName)
			}
		})
	}
}

func TestIsValidSpellNotFound(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	tests := []string{
		"NonExistentSpell",
		"Super Mega Fireball",
		"",
		"   ",
	}
	
	for _, spellName := range tests {
		t.Run(spellName, func(t *testing.T) {
			valid := repo.IsValidSpell(spellName)
			if valid {
				t.Errorf("Expected '%s' to be invalid spell", spellName)
			}
		})
	}
}

func TestIsValidSpellFileNotFound(t *testing.T) {
	repo := NewCSVSRDRepository(nonExistentEquipmentPath, nonExistentSpellsPath)
	
	valid := repo.IsValidSpell(spellFireball)
	if valid {
		t.Error(errExpectedFileNotExist)
	}
}

// ========== IsValidEquipment Tests ==========

func TestIsValidEquipmentSuccess(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	tests := []struct {
		name          string
		equipmentName string
	}{
		{"weapon", "Longsword"},
		{"lowercase weapon", "longsword"},
		{"armor", "Chain Mail"},
		{"armor without suffix", "Leather"},
		{"armor with suffix", "Leather Armor"},
		{"shield", "Shield"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := repo.IsValidEquipment(tt.equipmentName)
			if !valid {
				t.Errorf("Expected '%s' to be valid equipment", tt.equipmentName)
			}
		})
	}
}

func TestIsValidEquipmentNotFound(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	tests := []string{
		"NonExistentWeapon",
		"Super Sword",
		"",
		"   ",
	}
	
	for _, equipName := range tests {
		t.Run(equipName, func(t *testing.T) {
			valid := repo.IsValidEquipment(equipName)
			if valid {
				t.Errorf("Expected '%s' to be invalid equipment", equipName)
			}
		})
	}
}

func TestIsValidEquipmentFileNotFound(t *testing.T) {
	repo := NewCSVSRDRepository(nonExistentEquipmentPath, nonExistentSpellsPath)
	
	valid := repo.IsValidEquipment("Longsword")
	if valid {
		t.Error(errExpectedFileNotExist)
	}
}

// ========== IsSpellForClass Tests ==========

func TestIsSpellForClassSuccess(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	tests := []struct {
		spell     string
		class     string
		shouldBe  bool
	}{
		{spellFireball, "Wizard", true},
		{spellFireball, "Sorcerer", true},
		{spellFireball, "Cleric", false},
		{"Cure Wounds", "Cleric", true},
		{"Cure Wounds", "Wizard", false},
		{spellMagicMissile, "Wizard", true},
		{"magic missile", "wizard", true}, // lowercase
		{"MAGIC MISSILE", "WIZARD", true}, // uppercase
	}
	
	for _, tt := range tests {
		t.Run(tt.spell+"_for_"+tt.class, func(t *testing.T) {
			available, err := repo.IsSpellForClass(tt.spell, tt.class)
			if err != nil {
				t.Fatalf(errUnexpected, err)
			}
			if available != tt.shouldBe {
				t.Errorf("Expected '%s' for %s to be %v, got %v", 
					tt.spell, tt.class, tt.shouldBe, available)
			}
		})
	}
}

func TestIsSpellForClassNotFound(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	available, err := repo.IsSpellForClass("NonExistentSpell", "Wizard")
	if err != nil {
		t.Errorf("Expected no error for non-existent spell, got %v", err)
	}
	if available {
		t.Error("Expected false for non-existent spell")
	}
}

func TestIsSpellForClassFileError(t *testing.T) {
	repo := NewCSVSRDRepository(nonExistentEquipmentPath, nonExistentSpellsPath)
	
	_, err := repo.IsSpellForClass(spellFireball, "Wizard")
	if err == nil {
		t.Error(errExpectedFileNotExist)
	}
}

// ========== GetLearnableSpells Tests ==========

func TestGetLearnableSpellsSuccess(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	// Wizard with no known spells should get many learnable spells
	spells, err := repo.GetLearnableSpells("Wizard", []string{})
	if err != nil {
		t.Fatalf(errUnexpected, err)
	}
	
	if len(spells) == 0 {
		t.Error("Expected Wizard to have learnable spells")
	}
	
	// Should include Fireball
	hasFireball := false
	for _, s := range spells {
		if s == spellFireball {
			hasFireball = true
			break
		}
	}
	if !hasFireball {
		t.Error("Expected Wizard learnable spells to include Fireball")
	}
}

func TestGetLearnableSpellsExcludesKnown(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	knownSpells := []string{spellFireball, spellMagicMissile}
	spells, err := repo.GetLearnableSpells("Wizard", knownSpells)
	if err != nil {
		t.Fatalf(errUnexpected, err)
	}
	
	// Should not include Fireball or Magic Missile
	for _, s := range spells {
		if s == spellFireball || s == spellMagicMissile {
			t.Errorf("Learnable spells should not include already known spell: %s", s)
		}
	}
}

func TestGetLearnableSpellsClassSpecific(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	// Cleric should not get Wizard-only spells
	spells, err := repo.GetLearnableSpells("Cleric", []string{})
	if err != nil {
		t.Fatalf(errUnexpected, err)
	}
	
	hasFireball := false
	for _, s := range spells {
		if s == spellFireball {
			hasFireball = true
			break
		}
	}
	if hasFireball {
		t.Error("Cleric should not be able to learn Fireball (Wizard/Sorcerer spell)")
	}
}

func TestGetLearnableSpellsFileError(t *testing.T) {
	repo := NewCSVSRDRepository(nonExistentEquipmentPath, nonExistentSpellsPath)
	
	_, err := repo.GetLearnableSpells("Wizard", []string{})
	if err == nil {
		t.Error(errExpectedFileNotExist)
	}
}

// ========== GetSpellLevel Tests ==========

func TestGetSpellLevelSuccess(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	tests := []struct {
		spell string
		level int
	}{
		{"Acid Splash", 0},      // Cantrip
		{spellMagicMissile, 1},  // Level 1
		{spellFireball, 3},      // Level 3
		{"fireball", 3},         // Case insensitive
		{"FIREBALL", 3},         // Uppercase
	}
	
	for _, tt := range tests {
		t.Run(tt.spell, func(t *testing.T) {
			level, err := repo.GetSpellLevel(tt.spell)
			if err != nil {
				t.Fatalf(errUnexpected, err)
			}
			if level != tt.level {
				t.Errorf("Expected %s to be level %d, got %d", tt.spell, tt.level, level)
			}
		})
	}
}

func TestGetSpellLevelNotFound(t *testing.T) {
	equipPath, spellPath := getTestCSVPaths(t)
	repo := NewCSVSRDRepository(equipPath, spellPath)
	
	_, err := repo.GetSpellLevel("NonExistentSpell")
	if err == nil {
		t.Error("Expected error for non-existent spell")
	}
}

func TestGetSpellLevelFileError(t *testing.T) {
	repo := NewCSVSRDRepository(nonExistentEquipmentPath, nonExistentSpellsPath)
	
	_, err := repo.GetSpellLevel(spellFireball)
	if err == nil {
		t.Error(errExpectedFileNotExist)
	}
}

// ========== Helper Function Tests ==========

func TestCanonKey(t *testing.T) {
	repo := &CSVSRDRepository{}
	
	tests := []struct {
		input    string
		expected string
	}{
		{"Fireball", "fireball"},
		{"FIREBALL", "fireball"},
		{" Fireball ", "fireball"},
		{"Fire  Ball", "fire ball"},
		{"Magic Missile", "magic missile"},
		{"Acid'Splash", "acid'splash"},
		{"Test–Dash", "test-dash"},
		{"Test/Slash", "test slash"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := repo.canonKey(tt.input)
			if result != tt.expected {
				t.Errorf("canonKey(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Helper function to check if a slice contains all required strings
func containsAll(alternates []string, required []string) (bool, string) {
	for _, req := range required {
		found := false
		for _, alt := range alternates {
			if alt == req {
				found = true
				break
			}
		}
		if !found {
			return false, req
		}
	}
	return true, ""
}

func TestEquipAlternates(t *testing.T) {
	repo := &CSVSRDRepository{}
	
	tests := []struct {
		input    string
		mustHave []string
	}{
		{"leather", []string{"leather", "leather armor"}},
		{"chain mail", []string{"chain mail", "chain mail armor"}},
		{"plate", []string{"plate", "plate armor"}},
		{"studded leather", []string{"studded leather", "studded leather armor"}},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			alternates := repo.equipAlternates(tt.input)
			
			hasAll, missing := containsAll(alternates, tt.mustHave)
			if !hasAll {
				t.Errorf("equipAlternates(%q) missing %q, got %v", 
					tt.input, missing, alternates)
			}
		})
	}
}
