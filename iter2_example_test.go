package iter2_test

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mkch/iter2"
)

func ExampleZip() {
	ks := []int{1, 2, 3}
	vs := []string{"one", "two", "three"}

	type pair struct {
		N   int
		Str string
	}

	s := slices.Collect(iter2.Zip(slices.Values(ks), slices.Values(vs), func(i int, s string) pair { return pair{i, s} }))
	fmt.Println(s)
	// Output: [{1 one} {2 two} {3 three}]
}

func ExampleZip2() {
	ks := []int{1, 2, 3}
	vs := []string{"one", "two", "three"}
	zipped := iter2.Zip2(slices.Values(ks), slices.Values(vs))
	m := maps.Collect(zipped)
	for k, v := range m {
		fmt.Println(k, v)
	}
	// Unordered output:
	// 1 one
	// 2 two
	// 3 three
}

func ExampleConcat() {
	seq1 := slices.Values([]int{1, 2, 3})
	seq2 := slices.Values([]int{4, 5})
	seq := iter2.Concat(seq1, seq2)
	fmt.Println(slices.Collect(seq))
	// Output: [1 2 3 4 5]
}

func ExampleMerge() {
	seq1 := slices.Values([]int{1, 2, 3})
	seq2 := slices.Values([]int{4, 5})
	seq := iter2.Merge(seq1, seq2)
	for v := range seq {
		fmt.Println(v)
	}
	// Unordered output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleMap() {
	iter := slices.Values([]int{1, 2, 3})
	seq := iter2.Map(iter, func(v int) string { return strconv.Itoa(v + 1) })
	fmt.Printf("%q", slices.Collect(seq))
	// Output: ["2" "3" "4"]
}

func ExampleMap2() {
	iter := slices.All([]int{1, 2, 3})
	seq := iter2.Map2(iter, func(i int, v int) (int, string) { return i + 1, strings.Repeat("a", v) })
	for k, v := range seq {
		fmt.Printf("%v %q\n", k, v)
	}
	// Output:
	// 1 "a"
	// 2 "aa"
	// 3 "aaa"
}

func ExampleMap1To2() {
	seq1 := slices.Values([]int{1, 2, 3})
	seq2 := iter2.Map1To2(seq1, func(v int) (byte, string) { return byte(v), strconv.Itoa(v) })
	for k, v := range seq2 {
		fmt.Printf("%v %q\n", k, v)
	}
	// Output:
	// 1 "1"
	// 2 "2"
	// 3 "3"
}

func ExampleKeys() {
	seq2 := func(yield func(int, string) bool) {
		if !yield(0, "zero") {
			return
		}
		if !yield(1, "one") {
			return
		}
	}
	keys := iter2.Keys(seq2)
	fmt.Println(slices.Collect(keys))
	// Output: [0 1]
}

func ExampleValues() {
	seq2 := func(yield func(int, string) bool) {
		if !yield(0, "zero") {
			return
		}
		if !yield(1, "one") {
			return
		}
	}
	values := iter2.Values(seq2)
	fmt.Println(slices.Collect(values))
	// Output: [zero one]
}

func ExampleTake() {
	seq := slices.Values([]int{1, 2, 3, 4, 5})
	seq = iter2.Take(seq, 2)
	fmt.Println(slices.Collect(seq))
	// Output: [1 2]
}

func ExampleWalkDir() {
	dirs := iter2.WalkDir(os.DirFS("testdata"), ".")
	for d, err := range dirs {
		if err != nil {
			continue // continue to ignore the error
		}
		if d.Path == "should_skip" {
			d.SkipDir()
			continue
		}
		fmt.Printf("Walk: %v\n", d.Path)
	}
}

func ExamplePush() {
	seq, yield, stop := iter2.Push[int]()
	defer stop()
	time.AfterFunc(time.Millisecond*1, func() {
		yield(1)
		time.AfterFunc(time.Millisecond*1, func() {
			yield(2)
			stop()
		})
	})

	for t := range seq {
		fmt.Println(t)
	}
	// Output:
	// 1
	// 2
}

func ExampleFilter() {
	s := []int{1, 2, 3, 4}
	even := iter2.Filter(slices.Values(s),
		func(n int) bool { return n%2 == 0 })
	fmt.Println(slices.Collect(even))
	// Output: [2 4]
}

func ExampleFilter2() {
	s := []int{1, 2, 3, 4}
	even := iter2.Filter2(slices.All(s),
		func(i, n int) bool { return n%2 == 0 })
	for i, n := range even {
		fmt.Println(i, n)
	}
	// Unordered output:
	// 1 2
	// 3 4
}

func ExampleJust() {
	seq := iter2.Just(1, 2, 3)
	fmt.Println(slices.Collect(seq))
	// Output: [1 2 3]
}

func ExampleJust2() {
	type KV = struct {
		K int
		V string
	}
	seq := iter2.Just2(KV{1, "one"}, KV{2, "two"})
	for k, v := range seq {
		fmt.Println(k, v)
	}
	// Output:
	// 1 one
	// 2 two
}
