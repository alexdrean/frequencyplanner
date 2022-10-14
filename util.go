package main

import (
	"golang.org/x/exp/constraints"
	"sort"
)

func sortSlice[T constraints.Ordered](s []T) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}
