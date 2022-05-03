package utils

import "strings"

func UpperCaseFirstLetter(id string) string {
	switch len(id) {
	case 0:
		return ""
	case 1:
		return strings.Title(id)
	default:
		return strings.Title(id[:1]) + id[1:]

	}
}
