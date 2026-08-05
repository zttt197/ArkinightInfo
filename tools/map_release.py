import json, re, os

os.chdir(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

# Load character table and build name→id mapping
with open("data/gamedata/character_table.json", "r", encoding="utf-8") as f:
    ct = json.load(f)

name2id = {}
for cid, data in ct.items():
    name = data.get("name", "")
    appellation = data.get("appellation", "")
    if name:
        name2id[name] = cid
    if appellation and appellation != name:
        name2id[appellation] = cid

# PRTS data
prts_raw = open("tools/prts_dates.txt", "r", encoding="utf-8").read()

result = {}
matched = 0
unmatched = []

for line in prts_raw.strip().split("\n"):
    line = line.strip()
    if not line or "|" not in line:
        continue
    parts = line.split("|")
    name = parts[0].strip()
    date_raw = parts[1].strip() if len(parts) > 1 else ""

    m = re.search(r'(\d{4})年(\d{1,2})月(\d{1,2})日', date_raw)
    if m:
        date = f"{m.group(1)}-{int(m.group(2)):02d}-{int(m.group(3)):02d}"
    else:
        date = date_raw

    # Remove parenthetical suffixes for matching
    clean_name = re.sub(r'\([^)]*\)', '', name).strip()

    # Try exact first, then cleaned
    for n in [name, clean_name]:
        if n in name2id:
            cid = name2id[n]
            if not cid.startswith("token_") and not cid.startswith("trap_"):
                result[cid] = date
                matched += 1
                break
    else:
        unmatched.append(name)

print(f"Matched: {matched}")
print(f"Unmatched: {len(unmatched)}")

with open("data/gamedata/release_dates.json", "w", encoding="utf-8") as f:
    json.dump(result, f, ensure_ascii=False, indent=2)

print(f"Saved {len(result)} entries")
