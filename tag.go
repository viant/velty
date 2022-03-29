package velty

import (
	"strings"
)

const (
	nameSeparator = "|"
	velty         = "velty"
)

//Tag represent field tag
type Tag struct {
	Names  []string
	Prefix string
	Omit   bool
}

//Parse parses tag
func Parse(tagString string) *Tag {
	tag := &Tag{}
	tag.Names = make([]string, 0)

	elements := strings.Split(tagString, ",")
	if len(elements) == 0 {
		return tag
	}
	for i, element := range elements {
		nv := strings.Split(element, "=")
		if len(nv) == 2 {
			switch strings.ToLower(strings.TrimSpace(nv[0])) {
			case "names":
				tag.Names = strings.Split(strings.TrimSpace(nv[1]), nameSeparator)
			case "name":
				tag.Names = []string{strings.TrimSpace(nv[1])}
			case "prefix":
				tag.Prefix = strings.TrimSpace(nv[1])
			}

			continue
		}

		if len(nv) == 1 && strings.TrimSpace(nv[0]) == "-" {
			tag.Omit = true
			continue
		}

		if i == 0 {
			columnName := strings.TrimSpace(element)
			if len(columnName) > 0 {
				tag.Names = strings.Split(columnName, nameSeparator)
			}
		}
	}
	return tag
}

func (t *Tag) nameEqual(value string) bool {
	for _, name := range t.Names {
		if name == value {
			return true
		}
	}
	return false
}
