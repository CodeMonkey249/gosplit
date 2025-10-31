package gosplit

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

	ret := ""

	switch len(parts) {
	case 1:
		ret = fmt.Sprintf("%s.%s", s, ms)
	case 2:
		ret = fmt.Sprintf("%s:%s.%s", m, s, ms)
	case 3:
		ret = fmt.Sprintf("%s:%s:%s.%s", h, m, s, ms)
	default:
		ret = fmt.Sprintf("0.%s", ms)
	}

	mainPart, msPart = splitMilliseconds(ret)
	parts = strings.Split(mainPart, ":")
	modParts := parts
	for i, part := range parts {
		if i == len(parts) - 1 {
			modParts = []string{"0"}
			break
		}
		if part == "00" {
			modParts = parts[i+1:]
		} else if parts[i][0] == '0' && parts[i][1] != '0' {
			parts[i] = parts[i][1:]
			break
		} else {
			break
		}
	}
	parts = modParts
	ret = strings.Join([]string{strings.Join(parts, ":"), msPart}, ".")

	return ret
}

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

func timeAsString(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	str := fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
	return normalizeTimeString(str)
}

func parseSeconds(parts []string, i int) (int, int) {
	secParts := strings.Split(parts[i], ".")
	seconds, err := strconv.Atoi(secParts[0])
	if err != nil {
		panic(err)
	}

	millis, err := strconv.Atoi(secParts[1])
	if err != nil {
		panic(err)
	}
	return seconds, millis
}

func parseMinutes(parts []string, i int) int {
	minutes, err := strconv.Atoi(parts[i])
	if err != nil {
		panic(fmt.Errorf("invalid minutes: %w", err))
	}
	return minutes
}

func parseHours(parts []string, i int) int {
	hours, err := strconv.Atoi(parts[i])
	if err != nil {
		panic(fmt.Errorf("invalid hours: %w", err))
	}
	return hours
}

func timeAsDuration(s string) (time.Duration) {
	var hours, minutes, seconds, millis int
	parts := strings.Split(s, ":")

	switch len(parts) {
	case 1:
		seconds, millis = parseSeconds(parts, 0)
		minutes = 0
		hours = 0
	case 2:
		seconds, millis = parseSeconds(parts, 1)
		minutes = parseMinutes(parts, 0)
		hours = 0
	case 3:
		seconds, millis = parseSeconds(parts, 2)
		minutes = parseMinutes(parts, 1)
		hours = parseHours(parts, 0)
	default:
		seconds = 0
		minutes = 0
		hours = 0
	}

	total := time.Duration(hours)*time.Hour +
	time.Duration(minutes)*time.Minute +
	time.Duration(seconds)*time.Second +
	time.Duration(millis)*time.Millisecond

	return total
}
