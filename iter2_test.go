package iter2

import (
	"maps"
	"os"
	"slices"
	"strconv"
	"testing"
	"time"
)

func TestZip(t *testing.T) {
	ks := []int{1, 2, 3}
	vs := []string{"one", "two", "three"}

	type pair struct {
		N   int
		Str string
	}

	s := slices.Collect(Zip(slices.Values(ks), slices.Values(vs), func(i int, s string) pair { return pair{i, s} }))
	if !slices.Equal(s, []pair{{1, "one"}, {2, "two"}, {3, "three"}}) {
		t.Fatal(s)
	}
}

func TestZip2(t *testing.T) {
	ks := []int{1, 2, 3}
	vs := []string{"one", "two", "three"}

	m := maps.Collect(Zip2(slices.Values(ks), slices.Values(vs)))
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

func TestMap(t *testing.T) {
	seq := Map(slices.Values([]int{1, 2, 3}), func(v int) string { return strconv.Itoa(v + 1) })
	if s := slices.Collect(seq); !slices.Equal(s, []string{"2", "3", "4"}) {
		t.Fatal(s)
	}
}

func TestMap2(t *testing.T) {
	seq := Map2(slices.All([]string{"a", "b", "c"}), func(k int, v string) (string, string) { return strconv.Itoa(k + 1), "V" + v })
	if m := maps.Collect(seq); !maps.Equal(m, map[string]string{"1": "Va", "2": "Vb", "3": "Vc"}) {
		t.Fatal(m)
	}
}

func TestMap1To2(t *testing.T) {
	seq1 := slices.Values([]int{1, 2, 3})
	seq2 := Map1To2(seq1, func(v int) (byte, string) { return byte(v), strconv.Itoa(v) })
	if m := maps.Collect(seq2); !maps.Equal(m, map[byte]string{1: "1", 2: "2", 3: "3"}) {
		t.Fatal(m)
	}
}

func TestKeys(t *testing.T) {
	seq2 := func(yield func(int, string) bool) {
		if !yield(1, "one") {
			return
		}
		if !yield(2, "two") {
			return
		}
	}
	var keys []int
	if keys = slices.Collect(Keys(seq2)); !slices.Equal(keys, []int{1, 2}) {
		t.Fatal(keys)
	}

	// early stop
	keys = slices.Collect(Take(Keys(seq2), 1))
	if !slices.Equal(keys, []int{1}) {
		t.Fatal(keys)
	}
}

func TestValues(t *testing.T) {
	seq2 := func(yield func(int, string) bool) {
		if !yield(1, "one") {
			return
		}
		if !yield(2, "two") {
			return
		}
	}
	var values []string
	if values = slices.Collect(Values(seq2)); !slices.Equal(values, []string{"one", "two"}) {
		t.Fatal(values)
	}

	// early stop
	values = slices.Collect(Take(Values(seq2), 1))
	if !slices.Equal(values, []string{"one"}) {
		t.Fatal(values)
	}
}

func TestTake(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3})
	seq = Take(seq, 2)
	if s := slices.Collect(seq); !slices.Equal(s, []int{1, 2}) {
		t.Fatal(s)
	}

	seq = Take(seq, 0)
	if s := slices.Collect(seq); !slices.Equal(s, []int{}) {
		t.Fatal(s)
	}

	var panicked any
	func() {
		defer func() {
			panicked = recover()
		}()
		seq = Take(seq, -2)
	}()

	if panicked == nil {
		t.Fatal("should panic")
	}
}

func TestTake2(t *testing.T) {
	seq2 := slices.All([]string{"one", "two", "three"})
	seq2 = Take2(seq2, 2)
	if m := maps.Collect(seq2); !maps.Equal(m, map[int]string{0: "one", 1: "two"}) {
		t.Fatal(m)
	}

	seq2 = Take2(seq2, 0)
	if m := maps.Collect(seq2); !maps.Equal(m, map[int]string{}) {
		t.Fatal(m)
	}

	var panicked any
	func() {
		defer func() {
			panicked = recover()
		}()
		seq2 = Take2(seq2, -2)
	}()
	if panicked == nil {
		t.Fatal("should panic")
	}
}

func TestTakeCount(t *testing.T) {
	var yieldCount = 0
	seq := func(y func(int) bool) {
		yield := func(i int) bool {
			yieldCount++
			return y(i)
		}
		if !yield(0) {
			return
		}
		if !yield(1) {
			return
		}
	}
	if s := slices.Collect(Take(seq, 1)); !slices.Equal(s, []int{0}) {
		t.Fatal(s)
	}
	if yieldCount != 1 {
		t.Fatal(yieldCount)
	}

	yieldCount = 0
	seq2 := func(y func(int, int) bool) {
		yield := func(k, v int) bool {
			yieldCount++
			return y(k, v)
		}
		if !yield(0, 0) {
			return
		}
		if !yield(1, 1) {
			return
		}
		if !yield(2, 2) {
			return
		}
	}
	if s := slices.Collect(Keys(Take2(seq2, 2))); !slices.Equal(s, []int{0, 1}) {
		t.Fatal(s)
	}
	if yieldCount != 2 {
		t.Fatal(yieldCount)
	}
}

func TestWalkDir(t *testing.T) {
	seq := WalkDir(os.DirFS("testdata"), ".")

	files := Keys(Map2(seq, func(d *DirEntry, err error) (string, error) {
		if err != nil {
			panic(err)
		}
		return d.Path, nil
	}))
	s := slices.Collect(files)
	if !slices.Equal(s, []string{".", "a", "b", "dir1", "dir1/a", "e"}) {
		t.Fatal(s)
	}

	// test early stop
	s = slices.Collect(Take(files, 2))
	if !slices.Equal(s, []string{".", "a"}) {
		t.Fatal(s)
	}

	// test skip
	s = nil
	for d, err := range seq {
		if err != nil {
			panic(err)
		}
		if d.Entry.IsDir() && d.Path == "dir1" {
			d.SkipAll()
			continue
		}
		s = append(s, d.Path)
	}
	if !slices.Equal(s, []string{".", "a", "b"}) {
		t.Fatal(s)
	}

	s = nil
	for d, err := range seq {
		if err != nil {
			panic(err)
		}
		if d.Entry.IsDir() && d.Path == "dir1" {
			d.SkipDir()
			continue
		}
		s = append(s, d.Path)
	}
	if !slices.Equal(s, []string{".", "a", "b", "e"}) {
		t.Fatal(s)
	}
}

func TestWalkDirErr(t *testing.T) {
	seq := WalkDir(os.DirFS("testdata"), "NO THIS FILE")

	var s []string
	var err error
	var d *DirEntry
	for d, err = range seq {
		if err != nil {
			continue
		}
		s = append(s, d.Path)
	}
	if !slices.Equal(s, []string{}) {
		t.Fatal(s)
	}
	if err == nil || !os.IsNotExist(err) {
		t.Fatal("should be not exist")
	}
}

func TestPush(t *testing.T) {
	seq, yield, stop := Push[int]()

	go func() {
		defer stop()
		if !yield(1) {
			return
		}
		if !yield(2) {
			return
		}
		if !yield(3) {
			return
		}
	}()

	if s := slices.Collect(seq); !slices.Equal(s, []int{1, 2, 3}) {
		t.Fatal(s)
	}

	// early stop
	seq, yield, stop = Push[int]()

	go func() {
		defer stop()
		if !yield(1) {
			return
		}
		if !yield(2) {
			return
		}
		if !yield(3) {
			return
		}
	}()

	if s := slices.Collect(Take(seq, 1)); !slices.Equal(s, []int{1}) {
		t.Fatal(s)
	}
}

func TestPush2(t *testing.T) {
	seq, yield, stop := Push2[int, string]()

	go func() {
		defer stop()
		if !yield(1, "one") {
			return
		}
		if !yield(2, "two") {
			return
		}
		if !yield(3, "three") {
			return
		}
	}()

	if m := maps.Collect(seq); !maps.Equal(m, map[int]string{1: "one", 2: "two", 3: "three"}) {
		t.Fatal(m)
	}

	// early stop
	seq, yield, stop = Push2[int, string]()

	go func() {
		defer stop()
		if !yield(1, "one") {
			return
		}
		if !yield(2, "two") {
			return
		}
		if !yield(3, "three") {
			return
		}
	}()

	if m := maps.Collect(Take2(seq, 1)); !maps.Equal(m, map[int]string{1: "one"}) {
		t.Fatal(m)
	}
}
