package usecases

import (
    "errors"
)

// mockSRDRepository is a small test mock implementing repositories.SRDRepository.
// It exposes function fields so tests can override behavior per-case.
type mockSRDRepository struct {
    IsValidSpellFunc    func(spellName string) bool
    IsValidEquipmentFunc func(equipmentName string) bool
    IsSpellForClassFunc func(spellName, className string) (bool, error)
    GetLearnableSpellsFunc func(className string, knownSpells []string) ([]string, error)
    GetSpellLevelFunc    func(spellName string) (int, error)
}

func (m *mockSRDRepository) IsValidSpell(spellName string) bool {
    if m == nil || m.IsValidSpellFunc == nil {
        return false
    }
    return m.IsValidSpellFunc(spellName)
}

func (m *mockSRDRepository) IsValidEquipment(equipmentName string) bool {
    if m == nil || m.IsValidEquipmentFunc == nil {
        return false
    }
    return m.IsValidEquipmentFunc(equipmentName)
}

func (m *mockSRDRepository) IsSpellForClass(spellName, className string) (bool, error) {
    if m == nil || m.IsSpellForClassFunc == nil {
        return false, errors.New("IsSpellForClass not implemented in mock")
    }
    return m.IsSpellForClassFunc(spellName, className)
}

func (m *mockSRDRepository) GetLearnableSpells(className string, knownSpells []string) ([]string, error) {
    if m == nil || m.GetLearnableSpellsFunc == nil {
        return nil, errors.New("GetLearnableSpells not implemented in mock")
    }
    return m.GetLearnableSpellsFunc(className, knownSpells)
}

func (m *mockSRDRepository) GetSpellLevel(spellName string) (int, error) {
    if m == nil || m.GetSpellLevelFunc == nil {
        return 0, errors.New("GetSpellLevel not implemented in mock")
    }
    return m.GetSpellLevelFunc(spellName)
}

// Optional: implement adapter methods to return repositories.EquipmentInfo / SpellInfo
// if future tests need richer metadata. For now the interface methods above suffice.
