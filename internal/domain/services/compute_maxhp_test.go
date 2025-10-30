package services

import (
    "testing"

    "ddsheetfinal/internal/domain/entities"
)

// TestComputeMaxHPExamples verifies the two exam cases explicitly
func TestComputeMaxHPExamples(t *testing.T) {
    svc := NewCharacterService()

    rogue := &entities.Character{
        Name:  "Rogue-Test",
        Race:  "Halfling",
        Class: "Rogue",
        Level: 2,
        Str:   10,
        Dex:   16,
        Con:   12,
        Int:   13,
        Wis:   11,
        Cha:   14,
    }

    got := svc.ComputeMaxHP(rogue)
    want := 15
    if got != want {
        t.Fatalf("Rogue max HP = %d; want %d", got, want)
    }

    barb := &entities.Character{
        Name:  "Barbarian-Test",
        Race:  "Hill Dwarf",
        Class: "Barbarian",
        Level: 6,
        Str:   10,
        Dex:   10,
        Con:   14,
        Int:   10,
        Wis:   10,
        Cha:   10,
    }

    got2 := svc.ComputeMaxHP(barb)
    want2 := 59
    if got2 != want2 {
        t.Fatalf("Barbarian max HP = %d; want %d", got2, want2)
    }
}
