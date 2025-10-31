package services

import (
    "testing"
    "ddsheetfinal/internal/domain/entities"
)

func TestComputeMaxHPExamCases(t *testing.T) {
    svc := NewCharacterService()

    cases := []struct{
        name string
        ch   *entities.Character
        want int
    }{
        {
            name: "Rogue L2, CON 12 -> 15",
            ch: &entities.Character{
                Race: "Halfling", Class: "Rogue", Level: 2,
                Str:10, Dex:16, Con:12, Int:13, Wis:11, Cha:14,
            },
            want: 15,
        },
        {
            name: "Barbarian L6, CON 14 -> 59",
            ch: &entities.Character{
                Race: "Hill Dwarf", Class: "Barbarian", Level: 6,
                Str:10, Dex:10, Con:14, Int:10, Wis:10, Cha:10,
            },
            want: 59,
        },
        {
            name: "Level 0 returns 0 HP",
            ch: &entities.Character{
                Race: "Human", Class: "Fighter", Level: 0,
                Str:10, Dex:10, Con:10, Int:10, Wis:10, Cha:10,
            },
            want: 0,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got := svc.ComputeMaxHP(tc.ch)
            if got != tc.want {
                t.Fatalf("got %d, want %d", got, tc.want)
            }
        })
    }
}
