// ----------------------- helpers -----------------------
function signed(n) {
  const x = Number(n) || 0;
  return (x >= 0 ? "+" : "") + x;
}
function abilityMod(score) {
  const n = Number(score);
  if (!Number.isFinite(n)) return 0;
  return Math.floor((n - 10) / 2);
}
function profBonusFromLevel(level) {
  const L = Number(level) || 1;
  if (L >= 17) return 6;
  if (L >= 13) return 5;
  if (L >= 9)  return 4;
  if (L >= 5)  return 3;
  return 2;
}
function readProfBonus(form) {
  const el = form.querySelector('[name="proficiencybonus"]');
  if (!el) return 0;
  const n = Number(el.value);
  return Number.isFinite(n) ? n : 0;
}
// name setter

// ---------------- Saving throw profs by class ----------------
function setSaveProficienciesByClass(cls, form) {
  const map = {
    barbarian: ["Strength","Constitution"],
    bard:      ["Dexterity","Charisma"],
    cleric:    ["Wisdom","Charisma"],
    druid:     ["Intelligence","Wisdom"],
    fighter:   ["Strength","Constitution"],
    monk:      ["Strength","Dexterity"],
    paladin:   ["Wisdom","Charisma"],
    ranger:    ["Strength","Dexterity"],
    rogue:     ["Dexterity","Intelligence"],
    sorcerer:  ["Constitution","Charisma"],
    warlock:   ["Wisdom","Charisma"],
    wizard:    ["Intelligence","Wisdom"],
  };
  const profs = map[(cls || "").toLowerCase()] || [];
  for (const a of ["Strength","Dexterity","Constitution","Intelligence","Wisdom","Charisma"]) {
    const cb = form.querySelector(`[name="${a}-save-prof"]`);
    if (cb) cb.checked = profs.includes(a);
  }
}

// ---------------- Skills setup ----------------

const SKILL_TO_ABILITY = {
  "Acrobatics":      "Dexterity",
  "Animal Handling": "Wisdom",
  "Arcana":          "Intelligence",
  "Athletics":       "Strength",
  "Deception":       "Charisma",
  "History":         "Intelligence",
  "Insight":         "Wisdom",
  "Intimidation":    "Charisma",
  "Investigation":   "Intelligence",
  "Medicine":        "Wisdom",
  "Nature":          "Intelligence",
  "Perception":      "Wisdom",
  "Performance":     "Charisma",
  "Persuasion":      "Charisma",
  "Religion":        "Intelligence",
  "Sleight of Hand": "Dexterity",
  "Stealth":         "Dexterity",
  "Survival":        "Wisdom",
};

// Backgrounds and their fixed skills
const BACKGROUND_FIXED_SKILLS = {
  "acolyte":          ["Insight","Religion"],
  "charlatan":        ["Deception","Sleight of Hand"],
  "criminal":         ["Deception","Stealth"],
  "entertainer":      ["Acrobatics","Performance"],
  "folk hero":        ["Animal Handling","Survival"],
  "guild artisan":    ["Insight","Persuasion"],
  "hermit":           ["Medicine","Religion"],
  "noble":            ["History","Persuasion"],
  "outlander":        ["Athletics","Survival"],
  "sage":             ["Arcana","History"],
  "sailor":           ["Athletics","Perception"],
  "soldier":          ["Athletics","Intimidation"],
  "urchin":           ["Sleight of Hand","Stealth"],
};

// Tick skill prof checkboxes from JSON and (if empty) background
function setSkillProficienciesFromData(data, form) {
  // normalize incoming possible skill lists into a single array, then into a Set
  let profsArray = [];
  if (Array.isArray(data.skill_proficiencies)) {
    profsArray = data.skill_proficiencies.map(x => String(x).toLowerCase());
  } else if (Array.isArray(data.SkillProficiencies)) {
    profsArray = data.SkillProficiencies.map(x => String(x).toLowerCase());
  }
  const have = new Set(profsArray);

  // If none provided, auto-apply background fixed skills
  if (have.size === 0) {
    const bg = (data.background || data.Background || "").toString().toLowerCase();
    const fixed = BACKGROUND_FIXED_SKILLS[bg] || [];
    for (const s of fixed) have.add(s.toLowerCase());
  }

  // skip saving throw boxes (they end with -save-prof)
  for (const skill of Object.keys(SKILL_TO_ABILITY)) {
    const nameTitle = `${skill}-prof`;        // e.g., "Perception-prof"
    const nameLower = `${skill.toLowerCase()}-prof`; // e.g., "perception-prof"
    const cb = form.querySelector(`[name="${nameTitle}"]`) ||
               form.querySelector(`[name="${nameLower}"]`);
    if (cb) cb.checked = have.has(skill.toLowerCase());
  }
}

/* ---------------- Client-side derived (always runs) ---------------- */

// helpers to read ability scores from "*score" fields
function readAbilityScore(form, ability) {
  const map = {
    Strength: "Strengthscore",
    Dexterity: "Dexterityscore",
    Constitution: "Constitutionscore",
    Wisdom: "Wisdomscore",
    Intelligence: "Intelligencescore",
    Charisma: "Charismascore",
  };
  const el = form.querySelector(`[name="${map[ability]}"]`);
  const n = el ? Number(el.value) : Number.NaN;
  return Number.isFinite(n) ? n : 10;
}

function readLevelFromClassLevel(form) {
  const el = form.querySelector('[name="classlevel"]');
  if (!el) return 1;
  const parts = (el.value || "").trim().split(/\s+/);
  const maybe = Number.parseInt(parts.at(-1), 10);
  return Number.isFinite(maybe) ? maybe : 1;
}

function readProficiencyBonusOrInfer(form) {
  const pb = readProfBonus(form);
  if (pb) return pb;
  return profBonusFromLevel(readLevelFromClassLevel(form));
}

function setIfNotManual(form, name, value) {
  const el = form.querySelector(`[name="${name}"]`);
  if (el && !el.dataset.manual) el.value = value;
}

// Compute ability modifiers and fill ability mod fields when not manual.
function computeAbilityMods(form) {
  const abilities = ["Strength","Dexterity","Constitution","Intelligence","Wisdom","Charisma"];
  const mods = {};
  for (const ab of abilities) {
    const score = readAbilityScore(form, ab);
    mods[ab] = abilityMod(score);
    // fill ability mod fields if available
    setIfNotManual(form, `${ab}mod`, signed(mods[ab]));
  }
  return mods;
}

function applySavingThrows(form, mods, pb) {
  for (const ab of Object.keys(mods)) {
    const prof = !!form.querySelector(`[name="${ab}-save-prof"]`)?.checked;
    const total = mods[ab] + (prof ? pb : 0);
    setIfNotManual(form, `${ab}-save`, signed(total));
  }
}

function applySkillsLocal(form, mods, pb) {
  for (const [skill, ability] of Object.entries(SKILL_TO_ABILITY)) {
    const out = form.querySelector(`[name="${skill}"]`);
    if (!out) continue;
    const prof = !!form.querySelector(`[name="${skill}-prof"]`)?.checked;
    const expertise = !!form.querySelector(`[name="${skill}-expertise"]`)?.checked;
    let pbMult = 0;
    if (expertise) {
      pbMult = 2;
    } else if (prof) {
      pbMult = 1;
    }
    const total = mods[ability] + pb * pbMult;
    // debug: show calculation details for troubleshooting (outcommented)
    // console.log(`skill=${skill} prof=${prof} expertise=${expertise} pb=${pb} pbMult=${pbMult} mod=${mods[ability]} total=${total}`);
    if (!out.dataset.manual) out.value = signed(total);
  }
}

function applyInitiative(form, mods) {
  setIfNotManual(form, "initiative", signed(mods.Dexterity));
}

function applyPassivePerception(form) {
  const ppEl = form.querySelector('[name="passiveperception"]');
  const perceptionEl = form.querySelector('[name="Perception"]');
  if (ppEl && perceptionEl && !ppEl.dataset.manual) {
    const n = Number(String(perceptionEl.value).replaceAll(/[+]/g, "")) || 0;
    ppEl.value = 10 + n;
  }
}

function applyDerivedLocal(form) {
  if (!form) return;
  const mods = computeAbilityMods(form);
  const pb = readProficiencyBonusOrInfer(form);

  applySavingThrows(form, mods, pb);
  applySkillsLocal(form, mods, pb);
  applyInitiative(form, mods);
  applyPassivePerception(form);
}

// Clone helper
function cloneCharacterData(src) {
  if (!src || typeof src !== 'object') return {};
  if (typeof structuredClone === 'function') {
    try {
      return structuredClone(src);
    } catch (e) {
      if (typeof console !== 'undefined' && console.debug) console.debug('structuredClone failed, falling back to JSON clone', e);
      // fall through to JSON fallback
    }
  }
  // structuredClone isn't available or failed;
  try {
    const deepCloneFallback = (obj, seen = new WeakMap()) => {
      if (obj === null) return null;
      if (typeof obj !== 'object') return obj;
      if (seen.has(obj)) return seen.get(obj);
  if (obj instanceof Date) return new Date(obj);
      if (Array.isArray(obj)) {
        const out = [];
        seen.set(obj, out);
        for (let i = 0; i < obj.length; i++) out[i] = deepCloneFallback(obj[i], seen);
        return out;
      }
      // plain object
      const out = {};
      seen.set(obj, out);
      for (const k of Object.keys(obj)) {
        out[k] = deepCloneFallback(obj[k], seen);
      }
      return out;
    };
    return deepCloneFallback(src);
  } catch (e) {
    if (typeof console !== 'undefined' && console.debug) console.debug('deep clone fallback failed', e);
    return {};
  }
}

// ---------------- Derived numbers renderer ----------------
function renderDerived(form) {
  // debug: count invocations to detect duplicate handlers (outcommented)
  // try { console.count('renderDerived called'); console.log('renderDerived called'); } catch (e) {}
  // Always compute locally so the sheet works without a backend
  applyDerivedLocal(form);

  // Clone the global character data into a payload we can mutate safely.
  const payload = cloneCharacterData(globalThis.__characterData);

  // helper to read numeric inputs
  const valNum = (name) => {
    const el = form.querySelector(`[name="${name}"]`);
    if (!el) return undefined;
    const n = Number(el.value);
    return Number.isFinite(n) ? n : undefined;
  };

  // map form field names to payload fields
  if (valNum("Strengthscore") !== undefined) payload.Str = valNum("Strengthscore");
  if (valNum("Dexterityscore") !== undefined) payload.Dex = valNum("Dexterityscore");
  if (valNum("Constitutionscore") !== undefined) payload.Con = valNum("Constitutionscore");
  if (valNum("Intelligencescore") !== undefined) payload.Int = valNum("Intelligencescore");
  if (valNum("Wisdomscore") !== undefined) payload.Wis = valNum("Wisdomscore");
  if (valNum("Charismascore") !== undefined) payload.Cha = valNum("Charismascore");

  // level 
  const lvlEl = form.querySelector('[name="level"],[name="classlevel"],[name="Level"]');
  if (lvlEl) {
    const lv = Number(lvlEl.value);
    if (Number.isFinite(lv)) payload.Level = lv;
  }

  // equipment: send simple array of equipment names so server can detect shields
  const equipmentEl = form.querySelector('textarea[name="equipment"]');
  if (equipmentEl) payload.Equipment = equipmentEl.value.split(/\r?\n/).filter(Boolean);

  if (globalThis.__characterData?.enriched) {
    payload.enriched = globalThis.__characterData.enriched;
  }

  // Optional backend derive: use async/await with AbortController and a short timeout to
  // avoid leaving non-numeric placeholders in the AC field. If the server returns fields
  // we set them; otherwise we silently keep local values.
  (async function doServerDerive() {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 1500);
    try {
      const res = await fetch("/api/derive", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
        signal: controller.signal,
      });
      clearTimeout(timeout);
      if (!res.ok) return; // leave local values intact
      const obj = await res.json();
      if (!obj) return;

      // helper to set a field only when present and not manually edited
      const setIf = (selector, v, fmt) => {
        if (v === undefined) return;
        const el = form.querySelector(selector);
        if (!el || el.dataset.manual) return;
        el.value = (fmt ? fmt(v) : v);
      };

      // ability mods
      if (obj.ability_mods) {
        const set = (name, v) => setIf(`[name="${name}"]`, v, (x) => (x >= 0 ? "+" : "") + x);
        set("Strengthmod", obj.ability_mods.Strength);
        set("Dexteritymod", obj.ability_mods.Dexterity);
        set("Constitutionmod", obj.ability_mods.Constitution);
        set("Intelligencemod", obj.ability_mods.Intelligence);
        set("Wisdommod", obj.ability_mods.Wisdom);
        set("Charismamod", obj.ability_mods.Charisma);
      }
      // prof bonus
      setIf('[name="proficiencybonus"]', obj.proficiency_bonus);
      // initiative & passive perception
      setIf('[name="initiative"]', obj.initiative);
      setIf('[name="passiveperception"]', obj.passive_perception);
      // AC
      setIf('[name="ac"],[name="ArmorClass"]', obj.armor_class);

      // Re-run local derivations in case server values changed inputs
      applyDerivedLocal(form);
    } catch (e) {
      // Abort or network error — keep local values; log for visibility
      if (typeof console !== 'undefined' && console.debug) console.debug('doServerDerive failed or aborted:', e);
    } finally {
      clearTimeout(timeout);
    }
  })();
}

// ---------------- load & wire-up ----------------
document.addEventListener("DOMContentLoaded", function () {
  function qs(param) { return new URLSearchParams(location.search).get(param); }
  const name = qs("name");
  if (!name) return;

  async function tryFetchJson(url) {
    try {
      const res = await fetch(url, { cache: "no-store" });
      if (!res.ok) return null;
      return await res.json();
    } catch (e) { return null; }
  }

  

  (async function findAndLoad() {
    let charPath = "/characters/" + encodeURIComponent(name) + ".json";
    const manifest = await tryFetchJson("/frontend/manifest.json");
    if (Array.isArray(manifest)) {
      for (const p of manifest) {
        if (typeof p === "string" && p.toLowerCase().endsWith(("/" + name + ".json").toLowerCase())) {
          charPath = p;
          break;
        }
      }
    }
    const data = await tryFetchJson(charPath);
    if (!data) {
      console.warn("character not found at", charPath);
      return;
    }
  // Prefer enriched character if available
  globalThis.__characterData = data;
    // try to load enriched version (characters/<name>-enriched/<name>.json)
  const enrichedPath = "/enrichments/" + encodeURIComponent(name) + ".json";
    try {
      const enriched = await tryFetchJson(enrichedPath);
      if (enriched) {
        // use enriched
        globalThis.__characterData = enriched;
        populateForm(enriched);
        return;
      }
      // if enriched not found, try generating it via the server-side enrich endpoint
        const resp = await tryFetchJson(`/api/enrich?name=${encodeURIComponent(name)}`);
        if (resp?.output_path) {
          // try to load the newly created enriched file
          const enriched2 = await tryFetchJson(enrichedPath);
          if (enriched2) {
            globalThis.__characterData = enriched2;
            populateForm(enriched2);
            return;
          }
        }
    } catch (e) {
      // ignore and fall back to the non-enriched data
    }

    // fallback: use the plain character data
    populateForm(data);
  })();

  function setIf(name, value) {
    const el = document.querySelector(`[name="${name}"]`);
    if (!el) return;
    if (el.type === "checkbox") el.checked = !!value;
    else el.value = value === undefined || value === null ? "" : value;
  }

  function populateForm(data) {
    setProfileFields(data);
    setAbilityScores(data);
    setVitalsAndHp(data);
    setMoneyEquipmentSpells(data);

    const form = document.querySelector("form.charsheet");
    if (form) {
      // auto saving-throw profs from class
      const inferredClass = (data.Class || data.class || "").toString();
      setSaveProficienciesByClass(inferredClass, form);
      // auto skill profs from JSON or (fallback) background
      setSkillProficienciesFromData(data, form);
      renderDerived(form);
      // populate attacks table after derived values (ability mods / prof) are available
      try { populateAttacks(data); } catch (e) { if (typeof console !== 'undefined' && console.debug) console.debug('populateAttacks failed:', e); }
    }
  }

  function setProfileFields(data) {
    setIf("charname", data.Name || data.name || "");
    const cls = data.Class || data.class || "";
    const lvl = data.Level || data.level || "";
    setIf("classlevel", cls ? (cls + (lvl ? " " + lvl : "")) : "");
    setIf("background", data.Background || data.background || "");
    setIf("playername", data.Player || data.player || "");
    setIf("race", data.Race || data.race || "");
    setIf("alignment", data.Alignment || data.alignment || "");
    setIf("experiencepoints", data.XP || data.xp || data.ExperiencePoints || "");
  }

  function setAbilityScores(data) {
    setIf("Strengthscore", data.Str || data.Strength || data.str);
    setIf("Dexterityscore", data.Dex || data.Dexterity || data.dex);
    setIf("Constitutionscore", data.Con || data.Constitution || data.con);
    setIf("Wisdomscore", data.Wis || data.Wisdom || data.wis);
    setIf("Intelligencescore", data.Int || data.Intelligence || data.int);
    setIf("Charismascore", data.Cha || data.Charisma || data.cha);
  }

  function setVitalsAndHp(data) {
    setIf("passiveperception", data.PassivePerception || data.passivePerception || "");
    setIf("ac", data.ArmorClass || data.ac || "");
    setIf("initiative", data.Initiative || data.initiative || "");
    setIf("speed", data.Speed || data.speed || "");
  }

  function setMoneyEquipmentSpells(data) {
    setMoneyFields(data);
    setEquipmentTextArea(data);
    injectEquippedItems(data);
    setSpellsTextArea(data);
  }

  // money helpers
  function setMoneyFields(data) {
    setIf("cp", (data.Money?.cp) || data.cp || "");
    setIf("sp", (data.Money?.sp) || data.sp || "");
    setIf("gp", (data.Money?.gp) || data.gp || "");
  }

  // equipment textarea helpers
  function setEquipmentTextArea(data) {
    const equipmentText = Array.isArray(data.Equipment) ? data.Equipment.join("\n") : (data.Equipment || data.equipment || "");
    const equipEl = document.querySelector('textarea[name="equipment"]');
    if (!equipEl) return;
    equipEl.value = equipmentText;
    // show placeholder hint if empty
    if (!equipEl.value || !String(equipEl.value).trim()) {
      equipEl.placeholder = "(no equipment equipped)";
    } else {
      equipEl.placeholder = "";
    }
  }

  // detect and inject compact equipped items like Weapon / Armor / Shield / Off-hand
  // small helpers to keep cognitive complexity low
  function getEquipElement() {
    return document.querySelector('textarea[name="equipment"]');
  }

  function pickEquipmentNames(data) {
    const getVal = (candidates) => {
      for (const k of candidates) {
        if (data[k] !== undefined && data[k] !== null && String(data[k]).trim() !== "") return String(data[k]);
      }
      return "";
    };
    return {
      weaponName: getVal(["weapon", "Weapon", "WeaponName", "mainhand", "main_hand"]),
      armorName: getVal(["armor", "Armor", "ArmorName"]),
      shieldName: getVal(["shield", "Shield", "offhand", "off_hand", "OffHand"]),
      offHandName: getVal(["off_hand", "offhand", "OffHand"]),
    };
  }

  function buildEquipLines({weaponName, armorName, shieldName, offHandName}) {
    const lines = [];
    if (weaponName) lines.push(`Weapon: ${weaponName}`);
    if (armorName)  lines.push(`Armor: ${armorName}`);
    if (offHandName) lines.push(`Off-hand: ${offHandName}`);
    else if (shieldName) lines.push(`Shield: ${shieldName}`);
    return lines;
  }

  function mergeEquipLinesIntoTextarea(equipEl, equipLines) {
    // Inject into equipment textarea; avoid duplicates and preserve user-entered lines
    const existing = equipEl.value ? String(equipEl.value).split(/\r?\n/).map(s => s.trim()).filter(Boolean) : [];
    const existingLower = new Set(existing.map(e => e.toLowerCase()));
    const toPrepend = [];
    const toPrependLower = new Set();
    for (const line of equipLines) {
      const key = line.toLowerCase();
      if (!existingLower.has(key) && !toPrependLower.has(key)) {
        toPrepend.push(line);
        toPrependLower.add(key);
      }
    }
    if (toPrepend.length) {
      equipEl.value = toPrepend.join("\n") + (existing.length ? "\n" + existing.join("\n") : "");
    }
  }

  function injectEquippedItems(data) {
    const equipEl = getEquipElement();
    if (!equipEl) return;
    try {
      const names = pickEquipmentNames(data);
      const equipLines = buildEquipLines(names);
      if (!equipLines.length) return;
      try {
        mergeEquipLinesIntoTextarea(equipEl, equipLines);
      } catch (e) {
        if (typeof console !== 'undefined' && console.debug) console.debug('equipment injection failed:', e);
      }
    } catch (e) {
      // non-fatal: don't block the rest of the form; log for visibility
      if (typeof console !== 'undefined' && console.debug) console.debug('injectEquippedItems outer error:', e);
    }
  }

  function setSpellsTextArea(data) {
    const spellsText = Array.isArray(data.Spells) ? data.Spells.join("\n") : (data.Spells || data.spells || "");
    const spellsEl = document.querySelector('textarea[name="spellsarea"]');
    if (spellsEl) spellsEl.value = spellsText;
  }

  // Populate the attacks table (atknameN, atkbonusN, atkdamageN)
  function populateAttacks(data) {
    const form = document.querySelector("form.charsheet");
    if (!form) return;

    const setIfNotManual = (selector, v) => {
      const el = form.querySelector(selector);
      if (!el) return;
      if (el.dataset.manual) return;
      el.value = v === undefined || v === null ? "" : v;
    };

    // small helpers to read ability mods and prof
    const readAbility = (ability) => {
      const n = readAbilityScore(form, ability);
      return abilityMod(n);
    };
    const prof = readProficiencyBonusOrInfer(form) || 0;

  // get candidate weapon names: prefer explicit weapon/offhand fields
    const names = [];
    const namesPicked = pickEquipmentNames(data);
  if (typeof console !== 'undefined' && console.debug) console.debug('populateAttacks namesPicked=', namesPicked);
    if (namesPicked.weaponName) names.push(namesPicked.weaponName);
    if (namesPicked.offHandName && namesPicked.offHandName !== namesPicked.weaponName) names.push(namesPicked.offHandName);

    // include any equipment entries that look like weapons or have enriched damage
    const seen = new Set(names.map(s => (s || "").toLowerCase()));
    const enrichedMap = (data.enriched && data.enriched.equipment) ? data.enriched.equipment : {};
  if (typeof console !== 'undefined' && console.debug) console.debug('populateAttacks enriched keys=', Object.keys(enrichedMap));

    if (Array.isArray(data.Equipment)) {
      for (const entry of data.Equipment) {
        if (!entry) continue;
        const eName = String(entry).trim();
        if (!eName) continue;
        const key = eName.toLowerCase();
        if (seen.has(key)) continue;
        // prefer items that have enriched damage info
        const info = enrichedMap && (enrichedMap[eName] || enrichedMap[eName.toLowerCase()] || enrichedMap[key]);
        if ((info && info.damage_text) || /sword|dagger|axe|mace|scimitar|shortsword|longsword|bow|crossbow|spear|trident|club|halberd|rapier|whip/i.test(eName)) {
          names.push(eName);
          seen.add(key);
        }
        if (names.length >= 3) break;
      }
    }

    // limit to 3 attacks
    names.length = Math.min(3, names.length);

    for (let i = 0; i < 3; i++) {
      const n = names[i];
      const nameField = `[name="atkname${i+1}"]`;
      const bonusField = `[name="atkbonus${i+1}"]`;
      const damageField = `[name="atkdamage${i+1}"]`;

      if (!n) {
        // clear fields when not present and not manual
        setIfNotManual(nameField, "");
        setIfNotManual(bonusField, "");
        setIfNotManual(damageField, "");
        continue;
      }

  // detect equipped status: mark if it came from the compact weapon fields
  const equipped = (namesPicked.weaponName && namesPicked.weaponName === n) || (namesPicked.offHandName && namesPicked.offHandName === n);
  // displayName is identical to the item name; keep equipped boolean for future UI markers
  setIfNotManual(nameField, n);

      // determine damage text and average from enriched data when available
      let damageText = "";
      let damageAvg = undefined;
      const info = enrichedMap && (enrichedMap[n] || enrichedMap[n.toLowerCase()]);
      if (info) {
        if (info.damage_text) damageText = info.damage_text;
        if (info.damage_avg !== undefined && info.damage_avg !== null) damageAvg = info.damage_avg;
      }
        // heuristic: choose ability based on name / range (must compute before atk and damage avg)
        let ability = "Strength";
        const lname = n.toLowerCase();
        if (/bow|crossbow|sling|dart|shortbow|longbow|hand crossbow/.test(lname)) ability = "Dexterity";
        if (/dagger|rapier|scimitar|shortsword|finesse|whip|rapier/.test(lname)) ability = "Dexterity";

        const abilityModVal = readAbility(ability);
        const atk = abilityModVal + prof;
        setIfNotManual(bonusField, signed(atk));

        // prefer showing dice text followed by ability modifier and average in parentheses when available
        // include the ability modifier in the displayed average damage (DnD rules: add ability mod to damage)
        const totalAvg = (damageAvg !== undefined && damageAvg !== null) ? (Number(damageAvg) + abilityModVal) : (abilityModVal || undefined);
        let displayDamage = "";
        if (damageText) {
          const modStr = abilityModVal ? ` ${abilityModVal >= 0 ? '+' : ''}${abilityModVal}` : '';
          displayDamage = damageText + (modStr ? (` ${modStr}`) : '');
          if (totalAvg !== undefined) displayDamage += ` (avg ${totalAvg})`;
        } else {
          if (totalAvg !== undefined) displayDamage = `(avg ${totalAvg})`;
        }
        setIfNotManual(damageField, displayDamage);
        // store avg as a data attribute for debugging or future UI tweaks
        const dmgEl = form.querySelector(damageField);
        if (dmgEl) {
          if (totalAvg !== undefined) dmgEl.dataset.avg = String(totalAvg);
          else delete dmgEl.dataset.avg;
        }

        if (typeof console !== 'undefined' && console.debug) console.debug('populateAttacks set', nameField, displayName, bonusField, signed(atk), damageField, displayDamage);
    }
  }
  // live updates
  const formForLive = document.querySelector("form.charsheet");
  if (formForLive) {
    for (const n of ["Strengthscore","Dexterityscore","Constitutionscore","Wisdomscore","Intelligencescore","Charismascore"]) {
      const el = formForLive.querySelector(`[name="${n}"]`);
      if (!el) continue;
      el.addEventListener("input", () => renderDerived(formForLive));
    }
    const initEl = formForLive.querySelector('[name="initiative"]');
    if (initEl) initEl.addEventListener("input", () => initEl.dataset.manual = "1");
    const ppEl = formForLive.querySelector('[name="passiveperception"]');
    if (ppEl) ppEl.addEventListener("input", () => ppEl.dataset.manual = "1");

    // re-apply saving throw profs with changes
    const classEl = formForLive.querySelector('[name="classlevel"]');
    if (classEl) {
      classEl.addEventListener("input", () => {
        const text = (classEl.value || "").trim();
        const parts = text.split(/\s+/);
        let cls = text;
        const maybeLevel = Number.parseInt(parts.at(-1), 10);
        if (!Number.isNaN(maybeLevel)) { parts.pop(); cls = parts.join(" "); }
        setSaveProficienciesByClass(cls, formForLive);
        renderDerived(formForLive);
      });
    }

    // recalc when something changes
    formForLive.addEventListener("change", (e) => {
      const t = e.target;
      if (!(t instanceof HTMLInputElement)) return;
      if (/-(prof|expertise)$/i.test(t.name)) renderDerived(formForLive);
      if (/^proficiencybonus$/i.test(t.name)) renderDerived(formForLive);
    });
  } 
});
