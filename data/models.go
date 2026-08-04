package data

type rangeGridDto struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type rangeDto struct {
	ID        string         `json:"id"`
	Direction int            `json:"direction"`
	Grids     []rangeGridDto `json:"grids"`
}

// ---------- JSON DTO (matching game data format) ----------

type operatorDto struct {
	Name              string           `json:"name"`
	Appellation       string           `json:"appellation"`
	Rarity            string           `json:"rarity"`
	Profession        string           `json:"profession"`
	SubProfessionID   string           `json:"subProfessionId"`
	Position          string           `json:"position"`
	NationID          string           `json:"nationId"`
	DisplayNumber     string           `json:"displayNumber"`
	IsNotObtainable   bool             `json:"isNotObtainable"`
	TagList           []string         `json:"tagList"`
	Phases            []phaseDto       `json:"phases"`
	Skills            []skillRefDto    `json:"skills"`
	Talents           []talentDto      `json:"talents"`
	FavorKeyFrames    []favorFrameDto  `json:"favorKeyFrames"`
}

type phaseDto struct {
	MaxLevel              int                    `json:"maxLevel"`
	RangeID               string                 `json:"rangeId"`
	AttributesKeyFrames   []attributeKeyFrameDto `json:"attributesKeyFrames"`
}

type attributeKeyFrameDto struct {
	Level int          `json:"level"`
	Data  *attributeDto `json:"data"`
}

type attributeDto struct {
	MaxHp            float64 `json:"maxHp"`
	Atk              float64 `json:"atk"`
	Def              float64 `json:"def"`
	MagicResistance  float64 `json:"magicResistance"`
	Cost             float64 `json:"cost"`
	BlockCnt         float64 `json:"blockCnt"`
	RespawnTime      float64 `json:"respawnTime"`
}

type skillRefDto struct {
	SkillID    string        `json:"skillId"`
	UnlockCond *unlockCondDto `json:"unlockCond"`
}

type unlockCondDto struct {
	Phase string `json:"phase"`
	Level int    `json:"level"`
}

type talentDto struct {
	Candidates []talentCandidateDto `json:"candidates"`
}

type talentCandidateDto struct {
	UnlockCondition      *unlockCondDto `json:"unlockCondition"`
	RequiredPotentialRank int           `json:"requiredPotentialRank"`
	Name                 string         `json:"name"`
	Description          string         `json:"description"`
}

type favorFrameDto struct {
	Level int            `json:"level"`
	Data  *favorDataDto   `json:"data"`
}

type favorDataDto struct {
	MaxHp float64 `json:"maxHp"`
	Atk   float64 `json:"atk"`
}

type skillDto struct {
	SkillID string          `json:"skillId"`
	IconID  string          `json:"iconId"`
	Levels  []skillLevelDto `json:"levels"`
}

type skillLevelDto struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	SkillType   interface{}     `json:"skillType"`
	SpData      *spDataDto      `json:"spData"`
	Duration    *float64        `json:"duration"`
	RangeID     string          `json:"rangeId"`
	Blackboard  []blackboardDto `json:"blackboard"`
}

type spDataDto struct {
	SpType  interface{} `json:"spType"`
	SpCost  *float64    `json:"spCost"`
	InitSp  *float64    `json:"initSp"`
}

type blackboardDto struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type buildingDataDto struct {
	Chars map[string]buildingCharDto `json:"chars"`
	Buffs map[string]buffDto         `json:"buffs"`
}

type buildingCharDto struct {
	BuffChar []buffCharGroupDto `json:"buffChar"`
}

type buffCharGroupDto struct {
	BuffData []buffRefDto `json:"buffData"`
}

type buffRefDto struct {
	BuffID string `json:"buffId"`
}

type buffDto struct {
	BuffID      string `json:"buffId"`
	BuffName    string `json:"buffName"`
	Description string `json:"description"`
	RoomType    string `json:"roomType"`
}

// ---------- Domain models (returned to frontend) ----------

type Operator struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Appellation   string        `json:"appellation"`
	Rarity        int           `json:"rarity"`
	Stars         string        `json:"stars"`
	RarityText    string        `json:"rarityText"`
	ClassLabel    string        `json:"classLabel"`
	PositionLabel string        `json:"positionLabel"`
	NationLabel   string        `json:"nationLabel"`
	TagsText      string        `json:"tagsText"`
	Initial       string        `json:"initial"`
	DeployCost    string        `json:"deployCost"`
	RedeployText  string        `json:"redeployText"`
	TrustText     string        `json:"trustText"`
	Phases        []PhaseRow    `json:"phases"`
	Talents       []TalentItem  `json:"talents"`
	Skills        []SkillItem   `json:"skills"`
	BaseSkills    []BaseSkillItem `json:"baseSkills"`
	AvatarPath    string        `json:"avatarPath"`
}

type PhaseRow struct {
	Label     string    `json:"label"`
	Hp        string    `json:"hp"`
	Atk       string    `json:"atk"`
	Def       string    `json:"def"`
	Res       string    `json:"res"`
	Cost      string    `json:"cost"`
	Redeploy  string    `json:"redeploy"`
	Block     string    `json:"block"`
	RangeGrid RangeGrid `json:"rangeGrid"`
}

type RangeGrid struct {
	Rows int      `json:"rows"`
	Cols int      `json:"cols"`
	Grid []string `json:"grid"` // "self" | "range" | "empty"
}

type TalentItem struct {
	Title string `json:"title"`
	Meta  string `json:"meta"`
	Desc  string `json:"desc"`
}

type SkillItem struct {
	Title     string    `json:"title"`
	Meta      string    `json:"meta"`
	Desc      string    `json:"desc"`
	RangeGrid RangeGrid `json:"rangeGrid"`
}

type BaseSkillItem struct {
	Title string `json:"title"`
	Meta  string `json:"meta"`
	Desc  string `json:"desc"`
}
