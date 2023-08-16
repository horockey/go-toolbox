package comparer

var _ Comparer[string] = &StringComparer{}

type StringComparer struct{}

func (comp *StringComparer) Compare(a, b string) int {
	if a > b {
		return -1
	}
	if a < b {
		return 1
	}
	return 0
}
