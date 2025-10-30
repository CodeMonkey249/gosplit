package gosplit

import (
	"fmt"
	"strings"
)

// splitMilliseconds separates the fractional part if present
func splitMilliseconds(s string) (string, string) {
	if strings.Contains(s, ".") {
		parts := strings.SplitN(s, ".", 2)
		return parts[0], parts[1]
	}
	return s, ""
}

func normalizeTimeString(input string) string {
	var h, m, s string
	mainPart, msPart := splitMilliseconds(input)
	parts := strings.Split(mainPart, ":")

	switch len(parts) {
	case 1:
		s = parts[0]
	case 2:
		m = parts[0]
		s = parts[1]
	case 3:
		h = parts[0]
		m = parts[1]
		s = parts[2]
	default:
		return "0.000"
	}

	ms := normalizeMilliseconds(msPart)

	s, m, h = padIfNeeded(s, m, h)

	switch len(parts) {
	case 1:
		return fmt.Sprintf("%s.%s", s, ms)
	case 2:
		return fmt.Sprintf("%s:%s.%s", m, s, ms)
	case 3:
		return fmt.Sprintf("%s:%s:%s.%s", h, m, s, ms)
	default:
		return fmt.Sprintf("0.%s", ms)
	}
}

// normalizeMilliseconds ensures ms has exactly 3 digits
func normalizeMilliseconds(ms string) string {
	if ms == "" {
		return "000"
	}
	if len(ms) > 3 {
		return ms[:3]
	}
	for len(ms) < 3 {
		ms += "0"
	}
	return ms
}

// padIfNeeded adds a leading zero if single-digit but non-empty
func padIfNeeded(s string, m string, h string) (string, string, string) {
	if s == "" {
		s = "0"
	}
	if len(s) == 1 && len(m) > 0 {
		s = "0" + s
	}
	if len(m) == 1 && len(h) > 0 {
		m = "0" + m
	}
	return s, m, h
}
