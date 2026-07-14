package extractor

import (
	"regexp"
	"strings"
	"unicode"
)

var tokenRegex = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9\-_]*`)

func LegacyExtractTokens(rawText string) (allTokesn []string) {
	fields := strings.Fields(rawText)
	tokens := make([]string, 0, len(fields))

	for idx := range fields {
		trimmed := strings.TrimFunc(fields[idx], func(r rune) bool {
			return unicode.IsPunct(r) && r != '-' && r != '_'
		})

		lower := strings.ToLower(trimmed)

		if lower == "" || isNumeric(lower) || len(lower) < 2 {
			continue
		}

		tokens = append(tokens, lower)
	}

	return tokens
}

func ExtractTokens(rawText string) (allTokens []string) {
	matchs := tokenRegex.FindAllString(rawText, -1)

	seen := make(map[string]bool)

	tokens := make([]string, 0, len(matchs))

	for idx := range matchs {
		lower := strings.ToLower(matchs[idx])
		lower = strings.TrimSpace(lower)

		if len(lower) < 2 || isNumeric(lower) {
			continue
		}

		if !seen[lower] {
			seen[lower] = true
			tokens = append(tokens, lower)
		}

	}

	return tokens
}

func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) && r != '.' {
			return false
		}
	}
	return true
}
