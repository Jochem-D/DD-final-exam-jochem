package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ddsheetfinal/internal/domain/entities"
	"ddsheetfinal/internal/domain/repositories"
)

// JSONCharacterRepository implements CharacterRepository using JSON file storage
type JSONCharacterRepository struct {
	baseDir string
}

// NewJSONCharacterRepository creates a new JSON-based character repository
func NewJSONCharacterRepository(baseDir string) repositories.CharacterRepository {
	// Ensure the directory exists
	_ = os.MkdirAll(baseDir, 0o755)
	
	return &JSONCharacterRepository{
		baseDir: baseDir,
	}
}

// Save persists a character to a JSON file
func (r *JSONCharacterRepository) Save(character *entities.Character) error {
	data, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal character: %w", err)
	}

	filename := r.getFilePath(character.Name)
	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return fmt.Errorf("failed to write character file: %w", err)
	}

	return nil
}

// FindByName retrieves a character by name
func (r *JSONCharacterRepository) FindByName(name string) (*entities.Character, error) {
	filename := r.getFilePath(name)
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("character \"%s\" not found", name)
		}
		return nil, fmt.Errorf("failed to read character file: %w", err)
	}

	var character entities.Character
	if err := json.Unmarshal(data, &character); err != nil {
		return nil, fmt.Errorf("failed to parse character data: %w", err)
	}

	return &character, nil
}

// Delete removes a character by name
func (r *JSONCharacterRepository) Delete(name string) error {
	filename := r.getFilePath(name)
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete character file: %w", err)
	}
	return nil
}

// List returns all character names
func (r *JSONCharacterRepository) List() ([]string, error) {
	files, err := os.ReadDir(r.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read characters directory: %w", err)
	}

	var names []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			name := strings.TrimSuffix(file.Name(), ".json")
			names = append(names, name)
		}
	}

	return names, nil
}

// Exists checks if a character exists
func (r *JSONCharacterRepository) Exists(name string) bool {
	filename := r.getFilePath(name)
	_, err := os.Stat(filename)
	return err == nil
}

func (r *JSONCharacterRepository) getFilePath(name string) string {
	return filepath.Join(r.baseDir, fmt.Sprintf("%s.json", name))
}
