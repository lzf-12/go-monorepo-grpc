package regexp

import (
	"regexp"
	"strings"
)

var (
	compactWS = regexp.MustCompile(`\s+`)
)

func ReplaceWhitesWithSingleSpace(q string) string {
	q = compactWS.ReplaceAllString(q, " ")
	return strings.TrimSpace(q)
}
