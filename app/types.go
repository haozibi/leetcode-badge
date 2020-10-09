package app

type BadgeType int

const (
	BadgeTypeProfile BadgeType = iota + 1
	BadgeTypeUserNotFound
	BadgeTypeRanking
	BadgeTypeSolved
	BadgeTypeSolvedRate
	BadgeTypeAccepted
	BadgeTypeAcceptedRate
	BadgeTypeChartRanking
	BadgeTypeChartSolved
)

func (b BadgeType) String() string {
	switch b {
	case BadgeTypeProfile:
		return "profile"
	case BadgeTypeUserNotFound:
		return "user not found"
	case BadgeTypeRanking:
		return "ranking"
	case BadgeTypeSolved:
		return "solved"
	case BadgeTypeSolvedRate:
		return "solved-rate"
	case BadgeTypeAccepted:
		return "accepted"
	case BadgeTypeAcceptedRate:
		return "accepted-rate"
	case BadgeTypeChartRanking:
		return "chat-ranking"
	case BadgeTypeChartSolved:
		return "chat-solved"
	}

	return ""
}
