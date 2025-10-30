package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const (
	testCharName          = "TestCharacter"
	errExpectedNameFormat = "Expected name '%s', got '%s'"
	errInplaceFalse       = "Expected inplace to be false"
	errInplaceTrue        = "Expected inplace to be true"
	errForceFalse         = "Expected force to be false"
	errForceTrue          = "Expected force to be true"
	errFetchAllFalse      = "Expected fetchAll to be false"
	errFetchAllTrue       = "Expected fetchAll to be true"
	errUnexpected         = "Unexpected error: %v"
	flagInplace           = "--inplace"
	flagForce             = "--force"
	flagFetchAll          = "--fetch-all"
)

// ========== parseArgs Tests ==========
// Note: These tests use nil for the use case since parseArgs doesn't use it

func TestParseArgsNameOnly(t *testing.T) {
	handler := NewEnrichHandler(nil)
	
	inplace, force, fetchAll, name := handler.parseArgs([]string{testCharName})
	
	if name != testCharName {
		t.Errorf(errExpectedNameFormat, testCharName, name)
	}
	if inplace {
		t.Error(errInplaceFalse)
	}
	if force {
		t.Error(errForceFalse)
	}
	if fetchAll {
		t.Error(errFetchAllFalse)
	}
}

func TestParseArgsWithInplace(t *testing.T) {
	handler := NewEnrichHandler(nil)
	
	// Flags must come before positional arguments in Go's flag package
	inplace, force, fetchAll, name := handler.parseArgs([]string{flagInplace, testCharName})
	
	if name != testCharName {
		t.Errorf(errExpectedNameFormat, testCharName, name)
	}
	if !inplace {
		t.Error(errInplaceTrue)
	}
	if force {
		t.Error(errForceFalse)
	}
	if fetchAll {
		t.Error(errFetchAllFalse)
	}
}

func TestParseArgsWithAllFlags(t *testing.T) {
	handler := NewEnrichHandler(nil)
	
	inplace, force, fetchAll, name := handler.parseArgs([]string{flagInplace, flagForce, flagFetchAll, testCharName})
	
	if name != testCharName {
		t.Errorf(errExpectedNameFormat, testCharName, name)
	}
	if !inplace {
		t.Error(errInplaceTrue)
	}
	if !force {
		t.Error(errForceTrue)
	}
	if !fetchAll {
		t.Error(errFetchAllTrue)
	}
}

func TestParseArgsFlagsBeforeName(t *testing.T) {
	handler := NewEnrichHandler(nil)
	
	inplace, force, fetchAll, name := handler.parseArgs([]string{flagInplace, flagForce, testCharName})
	
	t.Logf("Parsed: inplace=%v, force=%v, fetchAll=%v, name=%s", inplace, force, fetchAll, name)
	
	if name != testCharName {
		t.Errorf(errExpectedNameFormat, testCharName, name)
	}
	if !inplace {
		t.Error(errInplaceTrue)
	}
	if !force {
		t.Error(errForceTrue)
	}
	if fetchAll {
		t.Error(errFetchAllFalse)
	}
}

func TestParseArgsNoName(t *testing.T) {
	handler := NewEnrichHandler(nil)
	
	inplace, force, fetchAll, name := handler.parseArgs([]string{flagFetchAll})
	
	if name != "" {
		t.Errorf("Expected empty name, got '%s'", name)
	}
	if inplace {
		t.Error(errInplaceFalse)
	}
	if force {
		t.Error(errForceFalse)
	}
	if !fetchAll {
		t.Error(errFetchAllTrue)
	}
}

func TestParseArgsNameWithSpaces(t *testing.T) {
	handler := NewEnrichHandler(nil)
	
	inplace, force, fetchAll, name := handler.parseArgs([]string{" TestCharacter "})
	
	if name != testCharName {
		t.Errorf("Expected name 'TestCharacter' (trimmed), got '%s'", name)
	}
	if inplace {
		t.Error(errInplaceFalse)
	}
	if force {
		t.Error(errForceFalse)
	}
	if fetchAll {
		t.Error(errFetchAllFalse)
	}
}

// ========== saveEnrichedCharacter Tests ==========

func TestSaveEnrichedCharacterInplace(t *testing.T) {
	handler := NewEnrichHandler(nil)
	
	// Create temporary test directory
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Change to temp directory
	os.Chdir(tmpDir)
	
	// Create data/characters directory
	os.MkdirAll(filepath.Join("data", "characters"), 0755)
	
	testData := map[string]interface{}{
		"name":  "TestCharacter",
		"class": "Wizard",
		"level": 5,
	}
	
	outPath, err := handler.saveEnrichedCharacter("TestCharacter", testData, true)
	if err != nil {
		t.Fatalf(errUnexpected, err)
	}
	
	expectedPath := filepath.Join("data", "characters", "TestCharacter.json")
	if outPath != expectedPath {
		t.Errorf("Expected path '%s', got '%s'", expectedPath, outPath)
	}
	
	// Verify file was created
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}
	
	// Verify content
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	if result["name"] != "TestCharacter" {
		t.Errorf("Expected name 'TestCharacter', got '%v'", result["name"])
	}
}

func TestSaveEnrichedCharacterEnrichmentDir(t *testing.T) {
	handler := NewEnrichHandler(nil)
	
	// Create temporary test directory
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Change to temp directory
	os.Chdir(tmpDir)
	
	testData := map[string]interface{}{
		"name":  "TestCharacter",
		"class": "Wizard",
		"level": 5,
	}
	
	outPath, err := handler.saveEnrichedCharacter("TestCharacter", testData, false)
	if err != nil {
		t.Fatalf(errUnexpected, err)
	}
	
	expectedPath := filepath.Join("data", "enrichments", "TestCharacter.json")
	if outPath != expectedPath {
		t.Errorf("Expected path '%s', got '%s'", expectedPath, outPath)
	}
	
	// Verify file was created
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}
	
	// Verify content
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	if result["name"] != "TestCharacter" {
		t.Errorf("Expected name 'TestCharacter', got '%v'", result["name"])
	}
}

func TestSaveEnrichedCharacterCreatesDirectory(t *testing.T) {
	handler := NewEnrichHandler(nil)
	
	// Create temporary test directory
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Change to temp directory
	os.Chdir(tmpDir)
	
	testData := map[string]interface{}{
		"name": "TestCharacter",
	}
	
	// enrichments directory doesn't exist yet
	enrichDir := filepath.Join("data", "enrichments")
	if _, err := os.Stat(enrichDir); !os.IsNotExist(err) {
		t.Fatal("enrichments directory should not exist yet")
	}
	
	_, err := handler.saveEnrichedCharacter("TestCharacter", testData, false)
	if err != nil {
		t.Fatalf(errUnexpected, err)
	}
	
	// Verify directory was created
	if _, err := os.Stat(enrichDir); os.IsNotExist(err) {
		t.Error("Expected enrichments directory to be created")
	}
}
