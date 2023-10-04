package string_comparer

import "github.com/horockey/go-toolbox/datastructs/pkg/comparer"

var _ comparer.Comparer[string] = &stringComparer{}

type stringComparer struct{}

func New() *stringComparer {
	return &stringComparer{}
}

func (comp *stringComparer) Compare(a, b string) int {
	if a > b {
		return -1
	}
	if a < b {
		return 1
	}
	return 0
}
