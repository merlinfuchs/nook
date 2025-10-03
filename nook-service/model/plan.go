package model

type Plan struct {
	ID          string
	Title       string
	Description string
	Default     bool
	Popular     bool
	Hidden      bool

	PaddleMonthlyPriceID string
	PaddleYearlyPriceID  string

	DiscordRoleID string

	Features Features
}

type Features struct {
	BasicAccess bool
}

func (f Features) Merge(other Features) Features {
	return Features{
		BasicAccess: f.BasicAccess || other.BasicAccess,
	}
}

func maxOrMinusOne(a, b int) int {
	if a == -1 || b == -1 {
		return -1
	}
	return max(a, b)
}
