// Simplified character sheet JS - all calculations done server-side
// This only handles: fetching data, populating fields, and saving

// ============================================================================
// Helper Functions
// ============================================================================

function getQueryParam(param) {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get(param);
}

function setFieldValue(fieldName, value) {
  const el = document.querySelector(`[name="${fieldName}"]`);
  if (!el) return;
  
  if (el.type === "checkbox") {
    el.checked = !!value;
  } else {
    el.value = value === undefined || value === null ? "" : value;
  }
}

function signed(n) {
  const x = Number(n) || 0;
  return (x >= 0 ? "+" : "") + x;
}

// ============================================================================
// Main Character Loading
// ============================================================================

document.addEventListener("DOMContentLoaded", async function () {
  const characterName = getQueryParam("name");
  if (!characterName) {
    console.error("No character name in URL");
    return;
  }

  try {
    // Fetch enriched character data (includes ALL calculated stats)
    const response = await fetch(`/api/character/enriched?name=${encodeURIComponent(characterName)}`);
    if (!response.ok) {
      throw new Error(`Failed to load character: ${response.statusText}`);
    }
    
    const data = await response.json();
    console.log("Loaded enriched character:", data);
    
    populateCharacterSheet(data);
  } catch (error) {
    console.error("Error loading character:", error);
    alert("Failed to load character: " + error.message);
  }
});

// ============================================================================
// Populate Character Sheet
// ============================================================================

function populateCharacterSheet(data) {
  // Basic Info
  setFieldValue("charname", data.name || "");
  
  const classLevel = data.class && data.level ? `${data.class} ${data.level}` : "";
  setFieldValue("classlevel", classLevel);
  
  setFieldValue("background", data.background || "");
  setFieldValue("race", data.race || "");
  setFieldValue("alignment", data.alignment || "");
  setFieldValue("experiencepoints", data.xp || "");
  setFieldValue("playername", data.player || "");

  // Ability Scores
  setFieldValue("Strengthscore", data.str);
  setFieldValue("Dexterityscore", data.dex);
  setFieldValue("Constitutionscore", data.con);
  setFieldValue("Intelligencescore", data.int);
  setFieldValue("Wisdomscore", data.wis);
  setFieldValue("Charismascore", data.cha);

  // Ability Modifiers (calculated by server)
  setFieldValue("Strengthmod", signed(data.str_mod));
  setFieldValue("Dexteritymod", signed(data.dex_mod));
  setFieldValue("Constitutionmod", signed(data.con_mod));
  setFieldValue("Intelligencemod", signed(data.int_mod));
  setFieldValue("Wisdommod", signed(data.wis_mod));
  setFieldValue("Charismamod", signed(data.cha_mod));

  // Derived Stats
  setFieldValue("proficiencybonus", signed(data.proficiency_bonus));
  setFieldValue("ac", data.armor_class);
  setFieldValue("initiative", signed(data.initiative));
  setFieldValue("speed", data.speed);
  setFieldValue("passiveperception", data.passive_perception);
  
  // HP
  setFieldValue("maxhp", data.max_hp);
  setFieldValue("currenthp", data.current_hp);

  // Saving Throws
  setFieldValue("Strength-save", signed(data.str_save));
  setFieldValue("Dexterity-save", signed(data.dex_save));
  setFieldValue("Constitution-save", signed(data.con_save));
  setFieldValue("Intelligence-save", signed(data.int_save));
  setFieldValue("Wisdom-save", signed(data.wis_save));
  setFieldValue("Charisma-save", signed(data.cha_save));

  // Saving Throw Proficiencies
  setFieldValue("Strength-save-prof", data.str_save_prof);
  setFieldValue("Dexterity-save-prof", data.dex_save_prof);
  setFieldValue("Constitution-save-prof", data.con_save_prof);
  setFieldValue("Intelligence-save-prof", data.int_save_prof);
  setFieldValue("Wisdom-save-prof", data.wis_save_prof);
  setFieldValue("Charisma-save-prof", data.cha_save_prof);

  // Skills (calculated by server with proficiency applied)
  if (data.skills && data.skill_profs) {
    populateSkills(data.skills, data.skill_profs);
  }

  // Equipment & Money
  setFieldValue("cp", data.money?.cp || "");
  setFieldValue("sp", data.money?.sp || "");
  setFieldValue("gp", data.money?.gp || "");
  
  // Equipment textarea - combine inventory items
  populateEquipment(data);
  
  // Attacks table - populate from weapon/off_hand
  populateAttacks(data);
  
  if (Array.isArray(data.spells)) {
    const spellsText = data.spells.join("\n");
    const spellsEl = document.querySelector('textarea[name="spells"]');
    if (spellsEl) spellsEl.value = spellsText;
  }
}

function populateSkills(skills, skillProfs) {
  const skillMapping = {
    "acrobatics": "Acrobatics",
    "animal handling": "Animal Handling",
    "arcana": "Arcana",
    "athletics": "Athletics",
    "deception": "Deception",
    "history": "History",
    "insight": "Insight",
    "intimidation": "Intimidation",
    "investigation": "Investigation",
    "medicine": "Medicine",
    "nature": "Nature",
    "perception": "Perception",
    "performance": "Performance",
    "persuasion": "Persuasion",
    "religion": "Religion",
    "sleight of hand": "Sleight of Hand",
    "stealth": "Stealth",
    "survival": "Survival"
  };

  for (const [skillKey, displayName] of Object.entries(skillMapping)) {
    const modifier = skills[skillKey];
    const isProficient = skillProfs[skillKey];
    
    // Skill modifier field (just the display name)
    if (modifier !== undefined) {
      setFieldValue(displayName, signed(modifier));
    }
    // Proficiency checkbox
    if (isProficient !== undefined) {
      setFieldValue(`${displayName}-prof`, isProficient);
    }
  }
}

function populateEquipment(data) {
  const equipmentLines = [];
  
  // Add equipped items first
  if (data.weapon) equipmentLines.push(`Weapon: ${data.weapon}`);
  if (data.off_hand) equipmentLines.push(`Off-hand: ${data.off_hand}`);
  if (data.armor) equipmentLines.push(`Armor: ${data.armor}`);
  if (data.shield) equipmentLines.push(`Shield: ${data.shield}`);
  
  // Add inventory items
  if (Array.isArray(data.inventory)) {
    equipmentLines.push(...data.inventory);
  }
  
  const equipmentText = equipmentLines.join("\n");
  const equipEl = document.querySelector('textarea[name="equipment"]');
  if (equipEl) {
    equipEl.value = equipmentText;
    if (!equipmentText.trim()) {
      equipEl.placeholder = "(no equipment)";
    }
  }
}

function populateAttacks(data) {
  const attacks = [];
  
  // Add weapon as first attack if present
  if (data.weapon) {
    const weaponAttack = calculateWeaponAttack(data.weapon, data);
    attacks.push(weaponAttack);
  }
  
  // Add off-hand as second attack if present
  if (data.off_hand) {
    const offHandAttack = calculateWeaponAttack(data.off_hand, data);
    attacks.push(offHandAttack);
  }
  
  // Populate attack fields
  for (let i = 0; i < attacks.length && i < 3; i++) {
    const attack = attacks[i];
    const rowNum = i + 1;
    
    setFieldValue(`atkname${rowNum}`, attack.name);
    setFieldValue(`atkbonus${rowNum}`, attack.bonus);
    setFieldValue(`atkdamage${rowNum}`, attack.damage);
  }
}

function calculateWeaponAttack(weaponName, data) {
  // Simple weapon damage lookup (can be expanded)
  const weaponData = {
    "dagger": { damage: "1d4", type: "piercing", finesse: true },
    "shortsword": { damage: "1d6", type: "piercing", finesse: true },
    "longsword": { damage: "1d8", type: "slashing", finesse: false },
    "greatsword": { damage: "2d6", type: "slashing", finesse: false },
    "handaxe": { damage: "1d6", type: "slashing", finesse: false },
    "battleaxe": { damage: "1d8", type: "slashing", finesse: false },
    "greataxe": { damage: "1d12", type: "slashing", finesse: false },
    "quarterstaff": { damage: "1d6", type: "bludgeoning", finesse: false },
    "mace": { damage: "1d6", type: "bludgeoning", finesse: false },
    "rapier": { damage: "1d8", type: "piercing", finesse: true },
    "scimitar": { damage: "1d6", type: "slashing", finesse: true },
  };
  
  const weaponLower = weaponName.toLowerCase().trim();
  const weapon = weaponData[weaponLower] || { damage: "1d6", type: "bludgeoning", finesse: false };
  
  // Determine which ability mod to use (finesse weapons can use DEX or STR, whichever is higher)
  let abilityMod = data.str_mod || 0;
  if (weapon.finesse && (data.dex_mod || 0) > abilityMod) {
    abilityMod = data.dex_mod || 0;
  }
  
  // Attack bonus = ability mod + proficiency bonus
  const attackBonus = abilityMod + (data.proficiency_bonus || 0);
  
  return {
    name: weaponName,
    bonus: signed(attackBonus),
    damage: `${weapon.damage} ${signed(abilityMod)} ${weapon.type}`
  };
}
