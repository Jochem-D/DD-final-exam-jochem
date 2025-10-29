package persistence

import (
	"os"
	"path/filepath"
	"testing"

	"ddsheetfinal/internal/domain/entities"
)

const errFailedToSave = "Failed to save character: %v"

func TestJSONRepositorySaveAndFind(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewJSONCharacterRepository(tmpDir)

	char := &entities.Character{
		Name:  "TestChar",
		Race:  "Human",
		Class: "Wizard",
		Level: 5,
		Str:   10,
		Dex:   14,
		Con:   12,
		Int:   18,
		Wis:   14,
		Cha:   10,
	}

	err := repo.Save(char)
	if err != nil {
		t.Fatalf(errFailedToSave, err)
	}

	filePath := filepath.Join(tmpDir, "TestChar.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Character file was not created")
	}

	found, err := repo.FindByName("TestChar")
	if err != nil {
		t.Fatalf("Failed to find character: %v", err)
	}

	assertCharacterMatch(t, found, char)
}

func TestJSONRepositoryFindNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewJSONCharacterRepository(tmpDir)

	_, err := repo.FindByName("NonExistent")
	if err == nil {
		t.Error("Expected error when finding non-existent character")
	}
}

func TestJSONRepositoryDelete(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewJSONCharacterRepository(tmpDir)

	char := &entities.Character{
		Name:  "ToDelete",
		Race:  "Elf",
		Class: "Ranger",
		Level: 3,
	}

	err := repo.Save(char)
	if err != nil {
		t.Fatalf(errFailedToSave, err)
	}

	if !repo.Exists("ToDelete") {
		t.Error("Character should exist before deletion")
	}

	err = repo.Delete("ToDelete")
	if err != nil {
		t.Fatalf("Failed to delete character: %v", err)
	}

	if repo.Exists("ToDelete") {
		t.Error("Character should not exist after deletion")
	}
}

func TestJSONRepositoryList(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewJSONCharacterRepository(tmpDir)

	chars := []*entities.Character{
		{Name: "Char1", Race: "Human", Class: "Fighter", Level: 1},
		{Name: "Char2", Race: "Elf", Class: "Wizard", Level: 2},
		{Name: "Char3", Race: "Dwarf", Class: "Cleric", Level: 3},
	}

	for _, char := range chars {
		err := repo.Save(char)
		if err != nil {
			t.Fatalf("Failed to save character %s: %v", char.Name, err)
		}
	}

	names, err := repo.List()
	if err != nil {
		t.Fatalf("Failed to list characters: %v", err)
	}

	if len(names) != 3 {
		t.Errorf("Expected 3 characters, got %d", len(names))
	}
}

func TestJSONRepositoryExists(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewJSONCharacterRepository(tmpDir)

	char := &entities.Character{
		Name:  "ExistsTest",
		Race:  "Human",
		Class: "Rogue",
		Level: 4,
	}

	if repo.Exists("ExistsTest") {
		t.Error("Character should not exist before creation")
	}

	err := repo.Save(char)
	if err != nil {
		t.Fatalf(errFailedToSave, err)
	}

	if !repo.Exists("ExistsTest") {
		t.Error("Character should exist after creation")
	}
}

func TestJSONRepositorySaveOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	repo := NewJSONCharacterRepository(tmpDir)

	char := &entities.Character{
		Name:  "UpdateTest",
		Race:  "Human",
		Class: "Fighter",
		Level: 1,
	}

	err := repo.Save(char)
	if err != nil {
		t.Fatalf("Failed to save initial character: %v", err)
	}

	char.Level = 10
	char.Class = "Paladin"

	err = repo.Save(char)
	if err != nil {
		t.Fatalf("Failed to save updated character: %v", err)
	}

	found, err := repo.FindByName("UpdateTest")
	if err != nil {
		t.Fatalf("Failed to find updated character: %v", err)
	}

	if found.Level != 10 {
		t.Errorf("Expected level 10, got %d", found.Level)
	}
	if found.Class != "Paladin" {
		t.Errorf("Expected class Paladin, got %s", found.Class)
	}
}

func assertCharacterMatch(t *testing.T, found *entities.Character, expected *entities.Character) {
	t.Helper()
	if found.Name != expected.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, found.Name)
	}
	if found.Class != expected.Class {
		t.Errorf("Expected class %s, got %s", expected.Class, found.Class)
	}
	if found.Level != expected.Level {
		t.Errorf("Expected level %d, got %d", expected.Level, found.Level)
	}
}
