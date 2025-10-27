// frontend/list.js
// List all characters from the manifest, let user open or create.

const listEl = document.getElementById("charList");
const emptyMsg = document.getElementById("emptyMsg");

// Helpers
const manifestURL = "/frontend/manifest.json";
const toName = (path) => {
  const base = path.split("/").pop() || "";
  return base.endsWith(".json") ? base.slice(0, -5) : base;
};
const charPath = (name) => `/characters/${encodeURIComponent(name)}.json`;

async function loadManifest() {
  const res = await fetch(manifestURL, { cache: "no-store" });
  if (!res.ok) throw new Error("Cannot load manifest");
  return await res.json(); // array of /characters/*.json
}

async function fetchMeta(path) {
  try {
    const res = await fetch(path, { cache: "no-store" });
    if (!res.ok) return null;
    const j = await res.json();
    // Nice-to-have summary (all optional)
    const cls = j.Class || j.class || "";
    const lvl = j.Level || j.level || "";
    const race = j.Race || j.race || "";
    return {
      display: (j.Name || j.name || toName(path)),
      subtitle: [race, cls && (lvl ? `${cls} ${lvl}` : cls)]
        .filter(Boolean)
        .join(" — ")
    };
  } catch {
    return null;
  }
}

function makeRow(name, path, meta) {
  const li = document.createElement("li");
  li.style.border = "1px solid #ddd";
  li.style.borderRadius = "10px";
  li.style.padding = "10px 12px";
  li.style.marginBottom = "10px";
  li.style.display = "flex";
  li.style.alignItems = "center";
  li.style.justifyContent = "space-between";
  li.style.gap = "12px";

  const left = document.createElement("div");
  const title = document.createElement("div");
  title.textContent = meta?.display || name;
  title.style.fontWeight = "600";
  const sub = document.createElement("div");
  sub.textContent = meta?.subtitle || "";
  sub.className = "meta";
  sub.style.color = "#666";
  sub.style.fontSize = "14px";
  left.appendChild(title);
  if (sub.textContent) left.appendChild(sub);

  const actions = document.createElement("div");
  actions.style.display = "flex";
  actions.style.gap = "8px";

  const open = document.createElement("a");
  open.href = `charactersheet.html?name=${encodeURIComponent(name)}`;
  open.textContent = "Open";
  open.className = "btn";
  open.style.textDecoration = "none";
  open.style.border = "1px solid #000";
  open.style.borderRadius = "10px";
  open.style.padding = "6px 10px";

  // Deletion is CLI-only
  actions.append(open);
  li.append(left, actions);
  return li;
}

// CREATE NEW
function addCreateUI() {
  const wrap = document.createElement("div");
  wrap.style.display = "flex";
  wrap.style.gap = "8px";
  wrap.style.margin = "16px 0";

  const input = document.createElement("input");
  input.placeholder = "New character name (file-safe)";
  input.style.flex = "1";
  input.maxLength = 60;

  const btn = document.createElement("button");
  btn.textContent = "Create";
  btn.onclick = async () => {
    const raw = (input.value || "").trim();
    if (!raw) return;
    // Keep filename simple
    const name = raw.replaceAll(/[^a-zA-Z0-9._ -]/g, "_");
    const payload = {
      Name: raw,
      Class: "",
      Level: 1,
      Race: "",
      Alignment: "",
      XP: 0,
      Str: 10, Dex: 10, Con: 10, Int: 10, Wis: 10, Cha: 10,
      ArmorClass: 10, Initiative: 0, Speed: 30,
      MaxHP: 10, CurrentHP: 10, TempHP: 0,
      Equipment: [], Spells: [], Money: { cp: 0, sp: 0, gp: 0 },
      SkillProficiencies: []
    };
    const res = await fetch(charPath(name), {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload, null, 2)
    });
    if (!res.ok) {
      alert("Create failed: " + (await res.text()));
      return;
    }
    location.href = `charactersheet.html?name=${encodeURIComponent(name)}`;
  };

  wrap.append(input, btn);
  listEl.parentNode.insertBefore(wrap, listEl);
}

// Render list
try {
  addCreateUI();
  const files = await loadManifest();
  if (!files || files.length === 0) {
    emptyMsg.style.display = "block";
  } else {
    emptyMsg.style.display = "none";
    files.sort((a, b) => a.localeCompare(b, undefined, { sensitivity: "base" }));
    for (const p of files) {
      const name = toName(p);
      // Normalize to the canonical characters path so fetchMeta requests
      // the correct location regardless of how the manifest listed it.
      const path = charPath(name);
      const meta = await fetchMeta(path);
      listEl.appendChild(makeRow(name, path, meta));
    }
  }
} catch (e) {
  emptyMsg.textContent = "Failed to load characters.";
  emptyMsg.style.display = "block";
  console.error(e);
}
