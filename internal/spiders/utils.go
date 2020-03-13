package spiders

import (
	"regexp"
	"strconv"
)

func parseHeat(text string) int32 {
	re := regexp.MustCompile(`^(\d+) ([万亿])`)
	m := re.FindStringSubmatch(text)
	if len(m) < 3 {
		return 0
	}
	heat, _ := strconv.Atoi(m[1])
	if m[2] == "万" {
		heat *= 10000
	} else if m[2] == "亿" {
		heat *= 100000000
	}
	return int32(heat)
}
