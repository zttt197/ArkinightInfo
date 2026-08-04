package data

import (
	"fmt"
	"regexp"
	"strings"
)

var Professions = map[string]string{
	"PIONEER": "先锋",
	"WARRIOR": "近卫",
	"TANK":    "重装",
	"SNIPER":  "狙击",
	"CASTER":  "术师",
	"MEDIC":   "医疗",
	"SUPPORT": "辅助",
	"SPECIAL": "特种",
}

var Positions = map[string]string{
	"MELEE":  "近战位",
	"RANGED": "远程位",
	"ALL":    "近战位/远程位",
	"NONE":   "—",
}

var Nations = map[string]string{
	"rhodes":    "罗德岛",
	"lungmen":   "龙门",
	"ursus":     "乌萨斯",
	"laterano":  "拉特兰",
	"columbia":  "哥伦比亚",
	"victoria":  "维多利亚",
	"yan":       "炎国",
	"sargon":    "萨尔贡",
	"kjerag":    "谢拉格",
	"higashi":   "东国",
	"bolivar":   "玻利瓦尔",
	"iberia":    "伊比利亚",
	"leithanien": "莱塔尼亚",
	"sami":      "萨米",
	"siracusa":  "叙拉古",
	"minos":     "米诺斯",
	"egir":      "阿戈尔",
	"kazimierz": "卡西米尔",
	"rim":       "雷姆必拓",
}

var SpTypes = map[string]string{
	"INCREASE_WITH_TIME":             "自动回复",
	"INCREASE_WITH_ATTACK":           "攻击回复",
	"INCREASE_WITH_ATTACK_AND_TIME":  "攻击/受击回复",
	"INCREASE_WITH_KILL":             "击杀回复",
	"INCREASE_WITH_TREASURY":         "弹药回复",
}

var SkillTypes = map[string]string{
	"AUTO":    "自动触发",
	"MANUAL":  "手动触发",
	"PASSIVE": "被动",
}

var Rooms = map[string]string{
	"CONTROL":    "控制中枢",
	"DORMITORY":  "宿舍",
	"MANUFACTURE": "制造站",
	"TRADING":    "贸易站",
	"POWER":      "发电站",
	"TRAINING":   "训练室",
	"WORKSHOP":   "加工站",
	"MEETING":    "会客室",
	"HIRE":       "办公室",
}

func mapLookup(m map[string]string, key string, fallback string) string {
	if key == "" {
		return fallback
	}
	if v, ok := m[key]; ok {
		return v
	}
	return key
}

func phaseLabel(phase string) string {
	switch phase {
	case "PHASE_0":
		return "初始"
	case "PHASE_1":
		return "精英1"
	case "PHASE_2":
		return "精英2"
	default:
		if phase == "" {
			return "—"
		}
		return phase
	}
}

func jsonText(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

var (
	placeholderRe = regexp.MustCompile(`\{([A-Za-z0-9_@\.\[\]]+)(?::([^}]*))?\}`)
	tagOpenRe     = regexp.MustCompile(`<@[^>]*>`)
	tagCloseRe    = regexp.MustCompile(`<[^>]*>`)
)

// CleanText removes game rich text tags and substitutes {variables} with values from blackboard.
func CleanText(text string, blackboard map[string]float64) string {
	if text == "" {
		return ""
	}
	s := strings.ReplaceAll(text, "\\n", "\n")

	s = placeholderRe.ReplaceAllStringFunc(s, func(m string) string {
		parts := placeholderRe.FindStringSubmatch(m)
		key := parts[1]
		f := parts[2]
		if v, ok := blackboard[key]; ok {
			if strings.Contains(f, "%") {
				return strings.TrimRight(
					strings.TrimRight(fmt.Sprintf("%.2f", v*100), "0"), ".") + "%"
			}
			if v == float64(int64(v)) {
				return fmt.Sprintf("%d", int64(v))
			}
			return strings.TrimRight(
				strings.TrimRight(fmt.Sprintf("%.2f", v), "0"), ".")
		}
		return key
	})

	s = tagOpenRe.ReplaceAllString(s, "")
	s = tagCloseRe.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}
