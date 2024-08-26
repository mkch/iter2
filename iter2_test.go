package iter2

import (
	"maps"
	"slices"
	"testing"
)

func TestZip(t *testing.T) {
	ks := []int{1, 2, 3}
	vs := []string{"one", "two", "three"}

	m := maps.Collect(Zip(slices.Values(ks), slices.Values(vs)))
	if !maps.Equal(m, map[int]string{1: "one", 2: "two", 3: "three"}) {
		t.Fatal(m)
	}
}

func TestConcat(t *testing.T) {
	seq1 := slices.Values([]int{1, 2, 3})
	seq2 := slices.Values([]int{4, 5})

	var s []int
	var i = 0
	for v := range Concat(seq1, seq2) {
		s = append(s, v)
		i++
		if i == 3 {
			break
		}
	}
	if slices.Equal(s, []int{1, 2, 3, 4}) {
		t.Fatal(s)
	}
}
