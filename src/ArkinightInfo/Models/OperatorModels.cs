using System.Collections.Generic;

namespace ArkinightInfo.Models;

// ---------- 游戏 JSON 原始结构（解析用 DTO） ----------

public sealed class OperatorDto
{
    public string? name { get; set; }
    public string? appellation { get; set; }
    public string? rarity { get; set; }
    public string? profession { get; set; }
    public string? subProfessionId { get; set; }
    public string? position { get; set; }
    public string? nationId { get; set; }
    public string? displayNumber { get; set; }
    public bool isNotObtainable { get; set; }
    public List<string>? tagList { get; set; }
    public List<PhaseDto>? phases { get; set; }
    public List<SkillRefDto>? skills { get; set; }
    public List<TalentDto>? talents { get; set; }
    public List<FavorFrameDto>? favorKeyFrames { get; set; }
}

public sealed class PhaseDto
{
    public int maxLevel { get; set; }
    public List<AttributeKeyFrameDto>? attributesKeyFrames { get; set; }
}

public sealed class AttributeKeyFrameDto
{
    public int level { get; set; }
    public AttributeDto? data { get; set; }
}

public sealed class AttributeDto
{
    public double maxHp { get; set; }
    public double atk { get; set; }
    public double def { get; set; }
    public double magicResistance { get; set; }
    public double cost { get; set; }
    public double blockCnt { get; set; }
    public double respawnTime { get; set; }
}

public sealed class SkillRefDto
{
    public string? skillId { get; set; }
    public UnlockCondDto? unlockCond { get; set; }
}

public sealed class UnlockCondDto
{
    public string? phase { get; set; }
    public int level { get; set; }
}

public sealed class TalentDto
{
    public List<TalentCandidateDto>? candidates { get; set; }
}

public sealed class TalentCandidateDto
{
    public UnlockCondDto? unlockCondition { get; set; }
    public int requiredPotentialRank { get; set; }
    public string? name { get; set; }
    public string? description { get; set; }
}

public sealed class FavorFrameDto
{
    public int level { get; set; }
    public FavorDataDto? data { get; set; }
}

public sealed class FavorDataDto
{
    public double maxHp { get; set; }
    public double atk { get; set; }
}

public sealed class SkillDto
{
    public string? skillId { get; set; }
    public string? iconId { get; set; }
    public List<SkillLevelDto>? levels { get; set; }
}

public sealed class SkillLevelDto
{
    public string? name { get; set; }
    public string? description { get; set; }
    public string? skillType { get; set; }
    public SpDataDto? spData { get; set; }
    public double duration { get; set; }
    public List<BlackboardDto>? blackboard { get; set; }
}

public sealed class SpDataDto
{
    public string? spType { get; set; }
    public double spCost { get; set; }
    public double initSp { get; set; }
}

public sealed class BlackboardDto
{
    public string? key { get; set; }
    public double value { get; set; }
}

public sealed class BuildingDataDto
{
    public Dictionary<string, BuildingCharDto>? chars { get; set; }
    public Dictionary<string, BuffDto>? buffs { get; set; }
}

public sealed class BuildingCharDto
{
    public List<BuffCharGroupDto>? buffChar { get; set; }
}

public sealed class BuffCharGroupDto
{
    public List<BuffRefDto>? buffData { get; set; }
}

public sealed class BuffRefDto
{
    public string? buffId { get; set; }
}

public sealed class BuffDto
{
    public string? buffId { get; set; }
    public string? buffName { get; set; }
    public string? description { get; set; }
    public string? roomType { get; set; }
}

// ---------- 展示用领域模型 ----------

public sealed class Operator
{
    public required string Id { get; init; }
    public required string Name { get; init; }
    public string Appellation { get; init; } = "";
    public int Rarity { get; init; }                 // 0..6，星级 = Rarity + 1
    public string Profession { get; init; } = "";
    public string SubProfessionId { get; init; } = "";
    public string Position { get; init; } = "";
    public string NationId { get; init; } = "";
    public string DisplayNumber { get; init; } = "";
    public IReadOnlyList<string> Tags { get; init; } = [];
    public IReadOnlyList<PhaseInfo> Phases { get; init; } = [];
    public IReadOnlyList<TalentInfo> Talents { get; init; } = [];
    public IReadOnlyList<SkillInfo> Skills { get; init; } = [];
    public IReadOnlyList<BaseSkillInfo> BaseSkills { get; init; } = [];
    public (int Hp, int Atk) TrustBonus { get; init; }
    public string? AvatarPath { get; init; }
}

public sealed class PhaseInfo
{
    public int Index { get; init; }                  // 0/1/2 = 精英0/1/2
    public int MaxLevel { get; init; }
    public AttributeSet MaxAttributes { get; init; } = new();
}

public sealed class AttributeSet
{
    public double MaxHp { get; init; }
    public double Atk { get; init; }
    public double Def { get; init; }
    public double MagicResistance { get; init; }
    public double Cost { get; init; }
    public double BlockCnt { get; init; }
    public double RespawnTime { get; init; }
}

public sealed class TalentInfo
{
    public string Name { get; init; } = "";
    public string Description { get; init; } = "";
    public string UnlockText { get; init; } = "";
}

public sealed class SkillInfo
{
    public string Name { get; init; } = "";
    public string Description { get; init; } = "";
    public string UnlockText { get; init; } = "";
    public string SpTypeLabel { get; init; } = "";
    public string SkillTypeLabel { get; init; } = "";
    public double SpCost { get; init; }
    public double InitSp { get; init; }
    public double Duration { get; init; }
}

public sealed class BaseSkillInfo
{
    public string Name { get; init; } = "";
    public string Description { get; init; } = "";
    public string RoomLabel { get; init; } = "";
}
