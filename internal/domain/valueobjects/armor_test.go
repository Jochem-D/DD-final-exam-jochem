package valueobjects

import "testing"

func TestGetArmorInfoLightArmor(t *testing.T) {
	tests := []struct {
		name         string
		armorName    string
		expectFound  bool
		expectedBase int
		expectedType ArmorType
	}{
		{"Leather armor exact", "leather", true, 11, ArmorTypeLight},
		{"Leather with suffix", "leather armor", true, 11, ArmorTypeLight},
		{"Studded leather", "studded leather", true, 12, ArmorTypeLight},
		{"Padded", "padded", true, 11, ArmorTypeLight},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, found := GetArmorInfo(tt.armorName)
			assertArmorInfo(t, tt.armorName, info, found, tt.expectFound, tt.expectedBase, tt.expectedType)
		})
	}
}

func TestGetArmorInfoMediumArmor(t *testing.T) {
	tests := []struct {
		name         string
		armorName    string
		expectedBase int
	}{
		{"Hide", "hide", 12},
		{"Chain shirt", "chain shirt", 13},
		{"Breastplate", "breastplate", 14},
		{"Half plate", "half plate", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, found := GetArmorInfo(tt.armorName)
			assertArmorInfo(t, tt.armorName, info, found, true, tt.expectedBase, ArmorTypeMedium)
		})
	}
}

func TestGetArmorInfoHeavyArmor(t *testing.T) {
	tests := []struct {
		name         string
		armorName    string
		expectedBase int
	}{
		{"Chain mail", "chain mail", 16},
		{"Plate", "plate", 18},
		{"Splint", "splint", 17},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, found := GetArmorInfo(tt.armorName)
			assertArmorInfo(t, tt.armorName, info, found, true, tt.expectedBase, ArmorTypeHeavy)
		})
	}
}

func TestGetArmorInfoCaseInsensitive(t *testing.T) {
	tests := []struct {
		name      string
		armorName string
	}{
		{"Upper case", "LEATHER"},
		{"Mixed case", "ChAiN MaIl"},
		{"With whitespace", "  leather  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, found := GetArmorInfo(tt.armorName)
			if !found {
				t.Errorf("GetArmorInfo(%q) should find armor", tt.armorName)
			}
		})
	}
}

func TestGetArmorInfoNotFound(t *testing.T) {
	tests := []struct {
		name      string
		armorName string
	}{
		{"Invalid armor", "chainmail bikini"},
		{"Empty string", ""},
		{"Random text", "something else"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, found := GetArmorInfo(tt.armorName)
			if found {
				t.Errorf("GetArmorInfo(%q) should not find armor", tt.armorName)
			}
		})
	}
}

func assertArmorInfo(t *testing.T, armorName string, info ArmorInfo, found bool, expectFound bool, expectedBase int, expectedType ArmorType) {
	t.Helper()
	if found != expectFound {
		t.Errorf("GetArmorInfo(%q) found = %v; want %v", armorName, found, expectFound)
	}
	if found {
		if info.Base != expectedBase {
			t.Errorf("GetArmorInfo(%q).Base = %d; want %d", armorName, info.Base, expectedBase)
		}
		if info.Type != expectedType {
			t.Errorf("GetArmorInfo(%q).Type = %s; want %s", armorName, info.Type, expectedType)
		}
	}
}
