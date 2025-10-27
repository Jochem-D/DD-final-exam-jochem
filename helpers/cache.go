// helpers/cache.go
package helpers

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func cacheDir() string {
	d := ".cache"
	_ = os.MkdirAll(d, 0o755)
	return d
}

func loadJSON[T any](p string, zero T) (T, error) {
	var out T = zero
	b, err := os.ReadFile(p)
	if err != nil {
		return out, err
	}
	_ = json.Unmarshal(b, &out)
	return out, nil
}
func saveJSON[T any](p string, v T) error {
	b, _ := json.MarshalIndent(v, "", "  ")
	return os.WriteFile(p, b, 0o644)
}

// --- concrete caches ---

type SpellInfo struct {
	Name   string `json:"name"`
	School string `json:"school"`
	Range  string `json:"range"`
	// Description is a short text/paragraph joined from the API's desc array
	Description string `json:"description,omitempty"`
}
type EquipmentInfo struct {
	Name                string `json:"name"`
	Category            string `json:"category"`
	RangeNormal         int    `json:"range_normal,omitempty"`
	TwoHanded           bool   `json:"two_handed,omitempty"`
	ArmorBase           int    `json:"armor_base,omitempty"`
	DexAllowed          bool   `json:"dex_allowed,omitempty"`
	MaxDex              *int   `json:"max_dex,omitempty"`
	StrMinimum          int    `json:"str_minimum,omitempty"`
	StealthDisadvantage bool   `json:"stealth_disadvantage,omitempty"`
	PriceGP    float64 `json:"price_gp,omitempty"`    // canonical price in gp
	DamageText string  `json:"damage_text,omitempty"` // exact dice text from API (e.g. "1d8")
	DamageAvg  float64 `json:"damage_avg,omitempty"`  // computed average damage
}

func LoadSpellCache() (map[string]SpellInfo, error) {
	return loadJSON(filepath.Join(cacheDir(), "spells.json"), map[string]SpellInfo{})
}
func SaveSpellCache(m map[string]SpellInfo) error {
	return saveJSON(filepath.Join(cacheDir(), "spells.json"), m)
}
func LoadEquipCache() (map[string]EquipmentInfo, error) {
	return loadJSON(filepath.Join(cacheDir(), "equipment.json"), map[string]EquipmentInfo{})
}
func SaveEquipCache(m map[string]EquipmentInfo) error {
	return saveJSON(filepath.Join(cacheDir(), "equipment.json"), m)
}
