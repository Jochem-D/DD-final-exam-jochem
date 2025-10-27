package repositories

import "ddsheetfinal/internal/domain/entities"

// CharacterRepository defines the contract for character persistence
type CharacterRepository interface {
	// Save persists a character
	Save(character *entities.Character) error
	
	// FindByName retrieves a character by name
	FindByName(name string) (*entities.Character, error)
	
	// Delete removes a character by name
	Delete(name string) error
	
	// List returns all character names
	List() ([]string, error)
	
	// Exists checks if a character exists
	Exists(name string) bool
}
