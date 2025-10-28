package dtos

import "ddsheetfinal/internal/domain/entities"

// Mapper functions to convert between domain entities and DTOs
// This layer prevents the presentation layer from knowing about domain entities

// ToCharacterDTO converts a domain Character entity to a DTO
func ToCharacterDTO(char *entities.Character) *CharacterDTO {
	if char == nil {
		return nil
	}
	
	return &CharacterDTO{
		Name:               char.Name,
		Race:               char.Race,
		Class:              char.Class,
		Background:         char.Background,
		Level:              char.Level,
		Str:                char.Str,
		Dex:                char.Dex,
		Con:                char.Con,
		Int:                char.Int,
		Wis:                char.Wis,
		Cha:                char.Cha,
		ProficiencyBonus:   char.ProficiencyBonus,
		SkillProficiencies: copyStringSlice(char.SkillProficiencies),
		Weapon:             char.Weapon,
		OffHand:            char.OffHand,
		Armor:              char.Armor,
		Shield:             char.Shield,
		Inventory:          copyStringSlice(char.Inventory),
		Spells:             copyStringSlice(char.Spells),
		PreparedSpells:     copyStringSlice(char.PreparedSpells),
	}
}

// ToCharacterEntity converts a DTO to a domain Character entity
func ToCharacterEntity(dto *CharacterDTO) *entities.Character {
	if dto == nil {
		return nil
	}
	
	return &entities.Character{
		Name:               dto.Name,
		Race:               dto.Race,
		Class:              dto.Class,
		Background:         dto.Background,
		Level:              dto.Level,
		Str:                dto.Str,
		Dex:                dto.Dex,
		Con:                dto.Con,
		Int:                dto.Int,
		Wis:                dto.Wis,
		Cha:                dto.Cha,
		ProficiencyBonus:   dto.ProficiencyBonus,
		SkillProficiencies: copyStringSlice(dto.SkillProficiencies),
		Weapon:             dto.Weapon,
		OffHand:            dto.OffHand,
		Armor:              dto.Armor,
		Shield:             dto.Shield,
		Inventory:          copyStringSlice(dto.Inventory),
		Spells:             copyStringSlice(dto.Spells),
		PreparedSpells:     copyStringSlice(dto.PreparedSpells),
	}
}

// ToCharacterListDTO converts a list of names to a DTO
func ToCharacterListDTO(names []string) *CharacterListDTO {
	return &CharacterListDTO{
		Names: copyStringSlice(names),
	}
}

// copyStringSlice creates a copy of a string slice to prevent mutations
func copyStringSlice(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
