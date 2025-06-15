package types

import "github.com/groundcover-com/groundcover-sdk-go/pkg/models"

const MatchTypeEqual = models.MatchType(MatchEqual)
const MatchTypeNotEqual = models.MatchType(MatchNotEqual)
const MatchTypeRegexp = models.MatchType(MatchRegexp)
const MatchTypeNotRegexp = models.MatchType(MatchNotRegexp)

// MatchType is an enum for label matching types.
type MatchType int

// Possible MatchTypes.
const (
	MatchEqual MatchType = iota
	MatchNotEqual
	MatchRegexp
	MatchNotRegexp
)

func (m MatchType) String() string {
	typeToStr := map[MatchType]string{
		MatchEqual:     "=",
		MatchNotEqual:  "!=",
		MatchRegexp:    "=~",
		MatchNotRegexp: "!~",
	}
	if str, ok := typeToStr[m]; ok {
		return str
	}
	panic("unknown match type")
}
