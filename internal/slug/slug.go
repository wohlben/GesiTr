package slug

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var (
	nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)
	leadTrailDash   = regexp.MustCompile(`^-+|-+$`)
)

// Generate creates a URL-friendly slug from a string by lowercasing,
// stripping diacritics, replacing non-alphanumeric runs with hyphens,
// and trimming leading/trailing hyphens.
func Generate(s string) string {
	s = strings.ToLower(s)
	// Decompose into base + combining marks, then strip combining marks
	s = strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Mn, r) { // Mn = Mark, Nonspacing (diacritics)
			return -1
		}
		return r
	}, norm.NFD.String(s))
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = leadTrailDash.ReplaceAllString(s, "")
	return s
}
