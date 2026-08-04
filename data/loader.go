package data

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const cacheFile = "operators.cache"

// GameData holds all parsed and filtered operator data.
type GameData struct {
	Version string
	Ops     []Operator
}

// ResolveDataRoot walks up from the executable to find the "data" directory.
func ResolveDataRoot() string {
	exe, err := os.Executable()
	if err != nil {
		exe = "."
	}
	exeDir := filepath.Dir(exe)

	// Search paths in priority order: exe-relative, then walk up
	dirs := []string{exeDir}
	dir := exeDir
	for i := 0; i < 10; i++ {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dirs = append(dirs, parent)
		dir = parent
	}

	// Also check working directory
	if wd, err := os.Getwd(); err == nil {
		dirs = append(dirs, wd)
	}

	for _, d := range dirs {
		candidate := filepath.Join(d, "data")
		// Verify it's the game data directory (contains gamedata/), not our Go package
		gamedata := filepath.Join(candidate, "gamedata")
		if info, err := os.Stat(gamedata); err == nil && info.IsDir() {
			return candidate
		}
	}
	return filepath.Join(exeDir, "data")
}

// ReadVersion extracts version info from data_version.txt.
func ReadVersion(dataRoot string) string {
	path := filepath.Join(dataRoot, "gamedata", "data_version.txt")
	b, err := os.ReadFile(path)
	if err != nil {
		return "未知"
	}
	raw := strings.TrimSpace(string(b))
	lines := strings.Split(raw, "\n")
	var version, date string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "VersionControl:") {
			version = strings.TrimPrefix(line, "VersionControl:")
		}
		if strings.HasPrefix(line, "Change:") {
			parts := strings.SplitN(strings.TrimPrefix(line, "Change:"), " on ", 2)
			if len(parts) == 2 {
				date = strings.TrimSpace(parts[1])
			}
		}
	}
	version = strings.TrimSpace(version)
	if version == "" {
		version = "未知"
	}
	if date != "" {
		return fmt.Sprintf("%s（%s）", version, date)
	}
	return version
}

// Load loads operator data with caching.
// Returns data from gob cache if available, otherwise parses JSON and creates cache.
func Load(dataRoot string) (*GameData, error) {
	version := ReadVersion(dataRoot)
	cachePath := filepath.Join(dataRoot, cacheFile)

	// Try cache first
	if info, err := os.Stat(cachePath); err == nil {
		// Check if cache is newer than data files
		charPath := filepath.Join(dataRoot, "gamedata", "character_table.json")
		if charInfo, err := os.Stat(charPath); err == nil {
			if info.ModTime().After(charInfo.ModTime()) {
				if ops, err := loadCache(cachePath); err == nil {
					return &GameData{Version: version, Ops: ops}, nil
				}
			}
		}
	}

	// Parse JSON
	ops, err := parseAll(dataRoot)
	if err != nil {
		return nil, err
	}

	// Save cache
	saveCache(cachePath, ops)

	return &GameData{Version: version, Ops: ops}, nil
}

func parseAll(dataRoot string) ([]Operator, error) {
	excel := filepath.Join(dataRoot, "gamedata")
	charPath := filepath.Join(excel, "character_table.json")
	skillPath := filepath.Join(excel, "skill_table.json")
	buildingPath := filepath.Join(excel, "building_data.json")

	if _, err := os.Stat(charPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("未找到数据文件 character_table.json")
	}
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("未找到数据文件 skill_table.json")
	}

	charTable, err := readJSON[map[string]operatorDto](charPath)
	if err != nil {
		return nil, fmt.Errorf("解析 character_table.json 失败: %w", err)
	}

	skillTable, err := readJSON[map[string]skillDto](skillPath)
	if err != nil {
		return nil, fmt.Errorf("解析 skill_table.json 失败: %w", err)
	}

	var building *buildingDataDto
	if _, err := os.Stat(buildingPath); err == nil {
		b, err := readJSON[buildingDataDto](buildingPath)
		if err == nil {
			building = &b
		}
	}

	var ops []Operator
	for id, dto := range charTable {
		// Filter: must have favor data (2+ frames) to be a playable operator
		if len(dto.FavorKeyFrames) < 2 {
			continue
		}
		// Filter: exclude summons / non-obtainable
		if strings.HasPrefix(id, "token_") || strings.HasPrefix(id, "trap_") {
			continue
		}
		if dto.IsNotObtainable {
			continue
		}

		rarity := parseRarity(dto.Rarity)
		phases := buildPhases(dto.Phases)
		if len(phases) == 0 {
			continue
		}

		talents := buildTalents(dto.Talents)
		skills := buildSkills(dto.Skills, skillTable)
		baseSkills := buildBaseSkills(id, building)
		trust := buildTrust(dto.FavorKeyFrames)

		var deployCost, redeployText string
		last := phases[len(phases)-1]
		deployCost = last.Cost
		redeployText = last.Redeploy

		stars := strings.Repeat("★", rarity)
		name := dto.Name
		if name == "" {
			name = id
		}
		initial := "?"
		if len(name) > 0 {
			initial = string([]rune(name)[0])
		}

		// Check for local avatar file, serve via /avatars/ route
		avatarPath := ""
		avatarFile := filepath.Join(dataRoot, "avatars", id+".png")
		if _, err := os.Stat(avatarFile); err == nil {
			avatarPath = "/avatars/" + id + ".png"
		}

		ops = append(ops, Operator{
			ID:            id,
			Name:          name,
			Appellation:   dto.Appellation,
			Rarity:        rarity,
			Stars:         stars,
			RarityText:    fmt.Sprintf("%d星", rarity),
			ClassLabel:    mapLookup(Professions, dto.Profession, "—"),
			PositionLabel: mapLookup(Positions, dto.Position, "—"),
			NationLabel:   mapLookup(Nations, dto.NationID, "—"),
			TagsText:      formatTags(dto.TagList),
			Initial:       initial,
			DeployCost:    deployCost,
			RedeployText:  redeployText,
			TrustText:     trust,
			Phases:        phases,
			Talents:       talents,
			Skills:        skills,
			BaseSkills:    baseSkills,
			AvatarPath:    avatarPath,
		})
	}

	// Sort by rarity desc, then name asc
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Rarity != ops[j].Rarity {
			return ops[i].Rarity > ops[j].Rarity
		}
		return ops[i].Name < ops[j].Name
	})

	return ops, nil
}

func readJSON[T any](path string) (T, error) {
	var zero T
	b, err := os.ReadFile(path)
	if err != nil {
		return zero, err
	}
	var val T
	if err := json.Unmarshal(b, &val); err != nil {
		return zero, err
	}
	return val, nil
}

func loadCache(path string) ([]Operator, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var ops []Operator
	if err := gob.NewDecoder(f).Decode(&ops); err != nil {
		return nil, err
	}
	return ops, nil
}

func saveCache(path string, ops []Operator) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(ops)
}

func parseRarity(r string) int {
	if !strings.HasPrefix(r, "TIER_") {
		return 0
	}
	var v int
	fmt.Sscanf(r, "TIER_%d", &v)
	return v
}

func buildPhases(phases []phaseDto) []PhaseRow {
	var result []PhaseRow
	for i, p := range phases {
		frames := p.AttributesKeyFrames
		if len(frames) == 0 {
			continue
		}
		last := frames[len(frames)-1]
		if last.Data == nil {
			continue
		}
		d := last.Data
		result = append(result, PhaseRow{
			Label:    fmt.Sprintf("精英%d 满级 Lv.%d", i, last.Level),
			Hp:       formatNum(d.MaxHp),
			Atk:      formatNum(d.Atk),
			Def:      formatNum(d.Def),
			Res:      formatNum(d.MagicResistance),
			Cost:     formatNum(d.Cost),
			Redeploy: formatNum(d.RespawnTime) + "s",
			Block:    formatNum(d.BlockCnt),
		})
	}
	return result
}

func buildTalents(talents []talentDto) []TalentItem {
	var result []TalentItem
	for i, t := range talents {
		if len(t.Candidates) == 0 {
			continue
		}
		// Take the last candidate with requiredPotentialRank == 0
		var cand talentCandidateDto
		for _, c := range t.Candidates {
			if c.RequiredPotentialRank == 0 {
				cand = c
			}
		}
		if cand.Name == "" {
			cand = t.Candidates[len(t.Candidates)-1]
		}
		var unlock string
		if cand.UnlockCondition != nil {
			unlock = phaseLabel(cand.UnlockCondition.Phase)
		}
		result = append(result, TalentItem{
			Title: fmt.Sprintf("天赋%d · %s", i+1, cand.Name),
			Meta:  unlock,
			Desc:  CleanText(cand.Description, nil),
		})
	}
	return result
}

func buildSkills(refs []skillRefDto, skillTable map[string]skillDto) []SkillItem {
	var result []SkillItem
	for i, r := range refs {
		if r.SkillID == "" {
			continue
		}
		skill, ok := skillTable[r.SkillID]
		if !ok || len(skill.Levels) == 0 {
			continue
		}
		maxLevel := skill.Levels[len(skill.Levels)-1]

		blackboard := make(map[string]float64)
		for _, b := range maxLevel.Blackboard {
			if b.Key != "" {
				blackboard[b.Key] = b.Value
			}
		}

		name := maxLevel.Name
		if name == "" {
			name = r.SkillID
		}

		spType := mapLookup(SpTypes, jsonText(maxLevel.SpData.SpType), "")
		skillType := mapLookup(SkillTypes, jsonText(maxLevel.SkillType), "")

		var metaParts []string
		if r.UnlockCond != nil {
			metaParts = append(metaParts, phaseLabel(r.UnlockCond.Phase))
		}
		if spType != "" {
			metaParts = append(metaParts, spType)
		}
		if skillType != "" {
			metaParts = append(metaParts, skillType)
		}
		if maxLevel.SpData != nil && maxLevel.SpData.SpCost != nil {
			metaParts = append(metaParts, fmt.Sprintf("技力消耗 %s", formatNum(*maxLevel.SpData.SpCost)))
		}
		if maxLevel.SpData != nil && maxLevel.SpData.InitSp != nil && *maxLevel.SpData.InitSp > 0 {
			metaParts = append(metaParts, fmt.Sprintf("初始技力 %s", formatNum(*maxLevel.SpData.InitSp)))
		}
		if maxLevel.Duration != nil && *maxLevel.Duration > 0 {
			metaParts = append(metaParts, fmt.Sprintf("持续 %ss", formatNum(*maxLevel.Duration)))
		}

		result = append(result, SkillItem{
			Title: fmt.Sprintf("技能%d · %s", i+1, name),
			Meta:  strings.Join(metaParts, " · "),
			Desc:  CleanText(maxLevel.Description, blackboard),
		})
	}
	return result
}

func buildBaseSkills(charID string, building *buildingDataDto) []BaseSkillItem {
	var result []BaseSkillItem
	if building == nil || building.Chars == nil {
		return result
	}
	bc, ok := building.Chars[charID]
	if !ok || len(bc.BuffChar) == 0 {
		return result
	}

	// Show base skills from the last elite phase group
	lastGroup := bc.BuffChar[len(bc.BuffChar)-1]
	for i, ref := range lastGroup.BuffData {
		if ref.BuffID == "" || building.Buffs == nil {
			continue
		}
		buff, ok := building.Buffs[ref.BuffID]
		if !ok {
			continue
		}
		name := buff.BuffName
		if name == "" {
			name = ref.BuffID
		}
		result = append(result, BaseSkillItem{
			Title: fmt.Sprintf("基建技能%d · %s", i+1, name),
			Meta:  mapLookup(Rooms, buff.RoomType, "—"),
			Desc:  CleanText(buff.Description, nil),
		})
	}
	return result
}

func buildTrust(frames []favorFrameDto) string {
	if len(frames) == 0 {
		return "—"
	}
	last := frames[len(frames)-1]
	if last.Data == nil {
		return "—"
	}
	hp, atk := int(last.Data.MaxHp), int(last.Data.Atk)
	if hp == 0 && atk == 0 {
		return "—"
	}
	return fmt.Sprintf("生命 +%d / 攻击 +%d", hp, atk)
}

func formatTags(tags []string) string {
	if len(tags) == 0 {
		return "—"
	}
	return strings.Join(tags, " / ")
}

func formatNum(v float64) string {
	if v == float64(int64(v)) {
		return fmt.Sprintf("%d", int64(v))
	}
	return strings.TrimRight(
		strings.TrimRight(fmt.Sprintf("%.2f", v), "0"), ".")
}
