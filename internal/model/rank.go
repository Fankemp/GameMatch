package model

type ValorantRank string

const (
	RankIron        ValorantRank = "Iron"
	RankBronze      ValorantRank = "Bronze"
	RankSilver      ValorantRank = "Silver"
	RankGold        ValorantRank = "Gold"
	RankPlatinum    ValorantRank = "Platinum"
	RankDiamond     ValorantRank = "Diamond"
	RankAscendant   ValorantRank = "Ascendant"
	RankImmortal    ValorantRank = "Immortal"
	RankRadiant     ValorantRank = "Radiant"
)

func (r ValorantRank) IsValid() bool {
	switch r {
	case RankIron, RankBronze, RankSilver, RankGold, RankPlatinum,
		RankDiamond, RankAscendant, RankImmortal, RankRadiant:
		return true
	}
	return false
}
