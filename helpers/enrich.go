// helpers/enrich.go
package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// getRateFromEnv reads DND5EAPI_RPS and returns a sane rate (min 1).
func getRateFromEnv(def int) int {
	if s := os.Getenv("DND5EAPI_RPS"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			return v
		}
	}
	return def
}

type EnrichOutput struct {
	OK         bool           `json:"ok"`
	OutputPath string         `json:"output_path,omitempty"`
	Fetched    map[string]int `json:"fetched,omitempty"`
	Cached     map[string]int `json:"cached,omitempty"`
	Logs       []string       `json:"logs,omitempty"`
	Error      string         `json:"error,omitempty"`
	ElapsedMs  int64          `json:"elapsed_ms"`
}

// package-level job and result types used by multiple helpers
type refJob struct{ kind, idx string }
type fetchResult struct {
	kind  string
	idx   string
	spell SpellInfo
	equip EquipmentInfo
}

// FetchSpellInfoPublic is an exported wrapper for fetchSpellInfo
func FetchSpellInfoPublic(ctx context.Context, idx string, logs LogSink) (SpellInfo, error) {
	return fetchSpellInfo(ctx, idx, logs)
}

// FetchEquipmentInfoPublic is an exported wrapper for fetchEquipmentInfo
func FetchEquipmentInfoPublic(ctx context.Context, idx string, logs LogSink) (EquipmentInfo, error) {
	return fetchEquipmentInfo(ctx, idx, logs)
}

// --- Deduped fetch/parse helpers ---

// fetchSpellInfo fetches one spell and converts it to SpellInfo.
func fetchSpellInfo(ctx context.Context, idx string, logs LogSink) (SpellInfo, error) {
	var body struct {
		Name   string `json:"name"`
		Range  string `json:"range"`
		School struct {
			Name string `json:"name"`
		} `json:"school"`
		Desc []string `json:"desc"`
	}
	if err := getJSONWithFallback(ctx, SpellPath(idx), &body, logs); err != nil {
		return SpellInfo{}, fmt.Errorf("spell %s: %w", idx, err)
	}
	desc := strings.Join(body.Desc, "\n\n")
	return SpellInfo{
		Name:        body.Name,
		School:      body.School.Name,
		Range:       body.Range,
		Description: desc,
	}, nil
}

// fetchEquipmentInfo fetches one equipment item and converts it to EquipmentInfo.
func fetchEquipmentInfo(ctx context.Context, idx string, logs LogSink) (EquipmentInfo, error) {
	var raw map[string]any
	if err := getJSONWithFallback(ctx, EquipmentPath(idx), &raw, logs); err != nil {
		return EquipmentInfo{}, fmt.Errorf("equipment %s: %w", idx, err)
	}
	info := EquipmentInfo{}
	// Basic fields
	if v, ok := raw["name"].(string); ok {
		info.Name = v
	}
	if v, ok := raw["equipment_category"].(map[string]any); ok {
		if n, ok2 := v["name"].(string); ok2 {
			info.Category = n
		}
	}

	// parse weapon-specific fields
	parseWeaponFields(raw, &info)

	// parse armor-specific fields
	parseArmorFields(raw, &info)

	// common fields
	if sd, ok := raw["stealth_disadvantage"].(bool); ok {
		info.StealthDisadvantage = sd
	}
	if sm, ok := raw["str_minimum"].(float64); ok {
		info.StrMinimum = int(sm)
	}
	// parse cost if available: {"quantity":..., "unit":"gp"}
	if cost, ok := raw["cost"].(map[string]any); ok {
		if gp := parseCostToGP(cost); gp > 0 {
			info.PriceGP = gp
		}
	}
	// parse damage if available (weapon entries)
	if dmg, ok := raw["damage"].(map[string]any); ok {
		if dmgStr, ok2 := dmg["damage_dice"].(string); ok2 {
			info.DamageText = strings.TrimSpace(dmgStr)
			info.DamageAvg = parseDiceAvg(info.DamageText)
		}
	}
	return info, nil
}

// parseCostToGP converts API cost object to gp floats (assumes unit is cp/sp/gp/pp)
func parseCostToGP(cost map[string]any) float64 {
	qf := 0.0
	if q, ok := cost["quantity"].(float64); ok {
		qf = q
	}
	unit := ""
	if u, ok := cost["unit"].(string); ok {
		unit = strings.ToLower(strings.TrimSpace(u))
	}
	switch unit {
	case "gp":
		return qf
	case "sp":
		return qf / 10.0
	case "cp":
		return qf / 100.0
	case "pp":
		return qf * 10.0
	default:
		return 0
	}
}

// parseDiceAvg computes the average expected value of a dice expression like "1d8" or "2d6".
// This is conservative: it only parses simple NdM forms and returns 0.0 for unknown formats.
func parseDiceAvg(s string) float64 {
	s = strings.TrimSpace(s)
	// expect forms like "1d8" or "2d6"
	parts := strings.Split(s, "d")
	if len(parts) != 2 {
		return 0
	}
	n, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || n <= 0 || m <= 0 {
		return 0
	}
	// average of 1..m is (m+1)/2, so NdM average is N*(m+1)/2
	return float64(n) * float64(m+1) / 2.0
}

// parseWeaponFields extracts weapon-specific data into info if present.
func parseWeaponFields(raw map[string]any, info *EquipmentInfo) {
	if _, ok := raw["weapon_range"].(string); !ok {
		return
	}
	if rng, ok := raw["range"].(map[string]any); ok {
		if n, ok2 := rng["normal"].(float64); ok2 {
			info.RangeNormal = int(n)
		}
	}
	if props, ok := raw["properties"].([]any); ok {
		for _, p := range props {
			if pm, ok := p.(map[string]any); ok {
				if n, ok := pm["name"].(string); ok && strings.EqualFold(n, "Two-Handed") {
					info.TwoHanded = true
					break
				}
			}
		}
	}
}

// parseArmorFields extracts armor-specific data into info if present.
func parseArmorFields(raw map[string]any, info *EquipmentInfo) {
	ac, ok := raw["armor_class"].(map[string]any)
	if !ok {
		return
	}
	if b, ok := ac["base"].(float64); ok {
		info.ArmorBase = int(b)
	}
	if db, ok := ac["dex_bonus"].(bool); ok {
		info.DexAllowed = db
		if db {
			if mb, ok := ac["max_bonus"].(float64); ok {
				m := int(mb)
				info.MaxDex = &m
			} else {
				info.MaxDex = nil
			}
		} else {
			z := 0
			info.MaxDex = &z
		}
	}
}

// marshalIndentToFile writes pretty JSON to path with the given permissions.
func marshalIndentToFile(path string, v any, perm os.FileMode) error {
	pretty, _ := json.MarshalIndent(v, "", "  ")
	return os.WriteFile(path, pretty, perm)
}

// fetchReferenceLists fetches the spell and equipment lists from the API.
func fetchReferenceLists(ctx context.Context, logs LogSink) (spellIndices []string, equipIndices []string, err error) {
	var spellsList struct {
		Results []struct {
			Index string `json:"index"`
		} `json:"results"`
	}
	if err := getJSONWithFallback(ctx, "/spells", &spellsList, logs); err != nil {
		return nil, nil, fmt.Errorf("fetch spells list: %w", err)
	}
	var equipList struct {
		Results []struct {
			Index string `json:"index"`
		} `json:"results"`
	}
	if err := getJSONWithFallback(ctx, "/equipment", &equipList, logs); err != nil {
		return nil, nil, fmt.Errorf("fetch equipment list: %w", err)
	}

	for _, r := range spellsList.Results {
		spellIndices = append(spellIndices, r.Index)
	}
	for _, r := range equipList.Results {
		equipIndices = append(equipIndices, r.Index)
	}
	return spellIndices, equipIndices, nil
}

// FetchAllReferences fetches the full lists of spells and equipment from the
// dnd5eapi and stores each item into the local .cache (spells.json, equipment.json).
// It uses the same rate-limited parallel runner as RunEnrich.
func FetchAllReferences(ctx context.Context, logs LogSink) error {
	spellIndices, equipIndices, err := fetchReferenceLists(ctx, logs)
	if err != nil {
		return err
	}

	spellCache, _ := LoadSpellCache()
	equipCache, _ := LoadEquipCache()

	rate := getRateFromEnv(8)
	if logs != nil {
		logs.Add(fmt.Sprintf("NOTE: fetch-all will use %d requests/sec. Please test with small batches (8-10) and be kind to the public API.", rate))
	}

	items := createRefJobs(spellIndices, equipIndices)

	fn := func(it refJob) error {
		switch it.kind {
		case "spell":
			sp, err := fetchSpellInfo(ctx, it.idx, logs)
			if err != nil {
				return err
			}
			spellCache[it.idx] = sp
		case "equip":
			ei, err := fetchEquipmentInfo(ctx, it.idx, logs)
			if err != nil {
				return err
			}
			equipCache[it.idx] = ei
		}
		return nil
	}

	genItems := make([]refJob, len(items))
	copy(genItems, items)
	if err := RunParallel(ctx, genItems, rate, 12, fn); err != nil {
		if logs != nil {
			logs.Add("fetch-all partial failures: " + err.Error())
		}
		return fmt.Errorf("fetch-all partial failures: %w", err)
	}

	if err := SaveSpellCache(spellCache); err != nil {
		return err
	}
	if err := SaveEquipCache(equipCache); err != nil {
		return err
	}
	return nil
}

// createRefJobs creates a slice of refJob for spells and equipment.
func createRefJobs(spellIndices, equipIndices []string) []refJob {
	var items []refJob
	for _, idx := range spellIndices {
		items = append(items, refJob{"spell", idx})
	}
	for _, idx := range equipIndices {
		items = append(items, refJob{"equip", idx})
	}
	return items
}

// RunEnrich loads /characters/<name>.json, fetches missing data from the API with concurrency,
// writes to enrichments/<name>.json (or in-place), and returns an EnrichOutput with logs.
func RunEnrich(ctx context.Context, name string, inplace bool, force bool, logs LogSink) (EnrichOutput, error) {
	buf := NewLogBuf()
	if logs == nil {
		logs = buf
	}
	start := time.Now()

	charPath := filepath.Join("characters", name+".json")
	b, err := os.ReadFile(charPath)
	if err != nil {
		return EnrichOutput{}, fmt.Errorf("read %s: %w", charPath, err)
	}

	var ch map[string]any
	if err := json.Unmarshal(b, &ch); err != nil {
		return EnrichOutput{}, fmt.Errorf("parse %s: %w", charPath, err)
	}

	spellNames := collectSpellNames(ch)
	equipNames := collectEquipmentNames(ch)
	buf.Add(fmt.Sprintf("character=%q spells=%d equipment=%d", name, len(spellNames), len(equipNames)))

	fetched, cached, skip, outPath := checkEnrichedExists(name, inplace, force, buf)
	if skip {
		return EnrichOutput{OK: true, OutputPath: outPath, Cached: cached, Logs: buf.Lines(), ElapsedMs: time.Since(start).Milliseconds()}, nil
	}

	spellCache, _ := LoadSpellCache()
	equipCache, _ := LoadEquipCache()

	jobs, idxSpells, idxEquip := createEnrichJobs(spellNames, equipNames, spellCache, equipCache, force, cached)

	fetchResults := fetchEnrichJobs(ctx, jobs, logs)

	mergeEnrichResults(fetchResults, spellCache, equipCache, fetched)

	_ = SaveSpellCache(spellCache)
	_ = SaveEquipCache(equipCache)

	spellsOut := fanOutSpellNames(idxSpells, spellCache)
	equipOut := fanOutEquipNames(idxEquip, equipCache)

	ch["enriched"] = map[string]any{
		"spells":    spellsOut,
		"equipment": equipOut,
	}

	outPath, werr := writeCharacter(name, ch, inplace)
	if werr != nil {
		return EnrichOutput{}, werr
	}

	buf.Add(fmt.Sprintf("wrote %s (inplace=%v)", outPath, inplace))

	out := buildEnrichOutput(outPath, fetched, cached, buf, start)
	return out, nil
}

// Helper to check if enriched file exists and should be skipped
func checkEnrichedExists(name string, inplace, force bool, buf LogSink) (map[string]int, map[string]int, bool, string) {
	fetched := map[string]int{"spells": 0, "equipment": 0}
	cached := map[string]int{"spells": 0, "equipment": 0}
	if !inplace {
		outDir := filepath.Join("enrichments")
		outPath := filepath.Join(outDir, name+".json")
		if _, statErr := os.Stat(outPath); statErr == nil {
			if !force {
				buf.Add("enriched file already exists; skipping fetch: " + outPath)
				return fetched, cached, true, outPath
			}
			buf.Add("enriched file exists but force=true, continuing to fetch missing caches: " + outPath)
		}
		return fetched, cached, false, outPath
	}
	return fetched, cached, false, ""
}

// Helper to create jobs for enrichment
func createEnrichJobs(spellNames, equipNames []string, spellCache map[string]SpellInfo, equipCache map[string]EquipmentInfo, force bool, cached map[string]int) ([]refJob, map[string]string, map[string]string) {
	var jobs []refJob
	idxSpells := UniqueIndices(spellNames, true)
	for _, idx := range idxSpells {
		if _, ok := spellCache[idx]; ok && !force {
			cached["spells"]++
			continue
		}
		jobs = append(jobs, refJob{"spell", idx})
	}
	idxEquip := UniqueIndices(equipNames, false)
	for _, idx := range idxEquip {
		if _, ok := equipCache[idx]; ok && !force {
			cached["equipment"]++
			continue
		}
		jobs = append(jobs, refJob{"equip", idx})
	}
	return jobs, idxSpells, idxEquip
}

// Helper to fetch jobs in parallel
func fetchEnrichJobs(ctx context.Context, jobs []refJob, logs LogSink) []fetchResult {
	results := make(chan fetchResult, len(jobs))
	fn := func(j refJob) error {
		switch j.kind {
		case "spell":
			sp, err := fetchSpellInfo(ctx, j.idx, logs)
			if err != nil {
				return err
			}
			results <- fetchResult{kind: "spell", idx: j.idx, spell: sp}
		case "equip":
			ei, err := fetchEquipmentInfo(ctx, j.idx, logs)
			if err != nil {
				return err
			}
			results <- fetchResult{kind: "equip", idx: j.idx, equip: ei}
		}
		return nil
	}
	rate := getRateFromEnv(8)
	if err := RunParallel(ctx, jobs, rate, 6, fn); err != nil {
		if logs != nil {
			logs.Add("fetch warnings: " + err.Error())
		}
	}
	close(results)
	var fetchResults []fetchResult
	for r := range results {
		fetchResults = append(fetchResults, r)
	}
	return fetchResults
}

// Helper to merge fetch results into caches and update fetched counters
func mergeEnrichResults(fetchResults []fetchResult, spellCache map[string]SpellInfo, equipCache map[string]EquipmentInfo, fetched map[string]int) {
	for _, r := range fetchResults {
		switch r.kind {
		case "spell":
			spellCache[r.idx] = r.spell
			fetched["spells"]++
		case "equip":
			equipCache[r.idx] = r.equip
			fetched["equipment"]++
		}
	}
}

// Helper to fan-out spell names
func fanOutSpellNames(idxSpells map[string]string, spellCache map[string]SpellInfo) map[string]SpellInfo {
	spellsOut := map[string]SpellInfo{}
	for name, idx := range idxSpells {
		if v, ok := spellCache[idx]; ok {
			spellsOut[name] = v
		}
	}
	return spellsOut
}

// Helper to fan-out equipment names
func fanOutEquipNames(idxEquip map[string]string, equipCache map[string]EquipmentInfo) map[string]EquipmentInfo {
	equipOut := map[string]EquipmentInfo{}
	for name, idx := range idxEquip {
		if v, ok := equipCache[idx]; ok {
			equipOut[name] = v
		}
	}
	return equipOut
}

// Helper to build EnrichOutput
func buildEnrichOutput(outPath string, fetched, cached map[string]int, buf *logBuf, start time.Time) EnrichOutput {
	return EnrichOutput{
		OK:         true,
		OutputPath: outPath,
		Fetched:    fetched,
		Cached:     cached,
		Logs:       buf.Lines(),
		ElapsedMs:  time.Since(start).Milliseconds(),
	}
}

// --- HTTP handler (hooked into serve mux) ---

func EnrichHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if name == "" {
		http.Error(w, "missing ?name=", http.StatusBadRequest)
		return
	}
	inplace := strings.EqualFold(r.URL.Query().Get("inplace"), "true")
	force := strings.EqualFold(r.URL.Query().Get("force"), "true")

	ctx := r.Context()
	start := time.Now()
	// use a buffer sink for HTTP handlers so we can return structured logs
	bufSink := NewLogBuf()
	out, err := RunEnrich(ctx, name, inplace, force, bufSink)

	// prefer the detailed output from RunEnrich; ensure ElapsedMs reflects handler time
	resp := out
	resp.OK = err == nil && out.OK
	if resp.ElapsedMs == 0 {
		resp.ElapsedMs = time.Since(start).Milliseconds()
	}
	if err != nil {
		resp.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// --- tiny character IO + name collectors (no duplicates with helpers.go) ---

func writeCharacter(name string, ch map[string]any, inplace bool) (string, error) {
	if inplace {
		outPath := filepath.Join("characters", name+".json")
		if err := marshalIndentToFile(outPath, ch, 0o644); err != nil {
			return "", fmt.Errorf("write %s: %w", outPath, err)
		}
		return outPath, nil
	}
	outDir := filepath.Join("enrichments")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", outDir, err)
	}
	outPath := filepath.Join(outDir, name+".json")
	if err := marshalIndentToFile(outPath, ch, 0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", outPath, err)
	}
	return outPath, nil
}

func collectSpellNames(ch map[string]any) []string {
	set := map[string]struct{}{}
	addList(set, toStringSlice(ch["spells"]))
	addList(set, toStringSlice(ch["prepared_spells"]))
	addList(set, toStringSlice(ch["Spells"]))
	addList(set, toStringSlice(ch["PreparedSpells"]))
	return setToSortedSlice(set)
}
func collectEquipmentNames(ch map[string]any) []string {
	set := map[string]struct{}{}
	addList(set, toStringSlice(ch["equipment"]))
	addOne(set, toString(ch["weapon"]))
	addOne(set, toString(ch["offhand"]))
	addOne(set, toString(ch["armor"]))
	addOne(set, toString(ch["shield"]))
	addList(set, toStringSlice(ch["Equipment"]))
	addOne(set, toString(ch["Weapon"]))
	addOne(set, toString(ch["OffHand"]))
	addOne(set, toString(ch["Armor"]))
	addOne(set, toString(ch["Shield"]))
	return setToSortedSlice(set)
}

// small helpers scoped to this file
func addList(set map[string]struct{}, items []string) {
	for _, s := range items {
		s = strings.TrimSpace(s)
		if s != "" {
			set[s] = struct{}{}
		}
	}
}
func addOne(set map[string]struct{}, s string) {
	s = strings.TrimSpace(s)
	if s != "" {
		set[s] = struct{}{}
	}
}
func setToSortedSlice(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for s := range set {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
func toStringSlice(v any) []string {
	switch vv := v.(type) {
	case []any:
		out := make([]string, 0, len(vv))
		for _, e := range vv {
			if s, ok := e.(string); ok && strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	case []string:
		out := make([]string, 0, len(vv))
		for _, s := range vv {
			if strings.TrimSpace(s) != "" {
				out = append(out, strings.TrimSpace(s))
			}
		}
		return out
	default:
		return nil
	}
}
func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
