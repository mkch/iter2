package iter2_test

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/mkch/iter2"
)

func ExampleZip() {
	ks := []int{1, 2, 3}
	vs := []string{"one", "two", "three"}
	zipped := iter2.Zip(slices.Values(ks), slices.Values(vs))
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
	// Output:
	// [1 2 3 4 5]
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
	// Output:
	// ["2" "3" "4"]
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
