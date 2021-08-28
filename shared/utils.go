package shared

import (
	"regexp"
	"strings"
)

var (
	matchFirst = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAll = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnake(s string) string {
	snake := matchFirst.ReplaceAllString(s, "${1}_${2}")
	snake = matchAll.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}