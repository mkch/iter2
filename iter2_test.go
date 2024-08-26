package iter2

import (
	"maps"
	"slices"
	"testing"
	"time"
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

func TestMerge(t *testing.T) {
	var seq1 = func(yield func(int) bool) {
		ticker := time.NewTicker(time.Millisecond * 20)
		defer ticker.Stop()
		for i := range 3 {
			<-ticker.C
			if !yield(i) {
				return
			}
		}
	}

	var seq2 = func(yield func(int) bool) {
		time.Sleep(time.Millisecond * 30)
		yield(-1)
	}

	seq := Merge(seq1, seq2)
	if s := slices.Collect(seq); !slices.Equal(s, []int{0, -1, 1, 2}) {
		t.Fatal(s)
	}

	seq = Merge(seq1, seq2)
	var s []int
	for v := range seq {
		s = append(s, v)
		if len(s) == 3 {
			break
		}
	}
	if !slices.Equal(s, []int{0, -1, 1}) {
		t.Fatal(s)
	}
}

func TestMergeSlice(t *testing.T) {
	seq1 := slices.Values([]int{1, 2, 3})
	seq2 := slices.Values([]int{4, 5})
	seq := Merge(seq1, seq2)
	var m = make(map[int]int)
	for v := range seq {
		m[v] = 0
	}
	if !maps.Equal(m, map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}) {
		t.Fatal(m)
	}
}
