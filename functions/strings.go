package functions

import "strings"

type Strings struct {
}

func (s Strings) ToLower(val string) string {
	return strings.ToLower(val)
}

func (s Strings) ToUpper(val string) string {
	return strings.ToUpper(val)
}

func (s Strings) ReplaceAll(val, old, new string) string {
	return strings.ReplaceAll(val, old, new)
}

func (s Strings) HasSuffix(val, suffix string) bool {
	return strings.HasSuffix(val, suffix)
}

func (s Strings) HasPrefix(val, suffix string) bool {
	return strings.HasSuffix(val, suffix)
}

func (s Strings) Index(val, substr string) int {
	return strings.Index(val, substr)
}

func (s Strings) Fields(val string) []string {
	return strings.Fields(val)
}

func (s Strings) Split(val, sep string) []string {
	return strings.Split(val, sep)
}

func (s Strings) TrimSpace(val string) string {
	return strings.TrimSpace(val)
}

func (s Strings) Trim(val, cutset string) string {
	return strings.Trim(val, cutset)
}

func (s Strings) LastIndex(val, substr string) int {
	return strings.LastIndex(val, substr)
}
