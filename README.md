# iter2

 Go iter utilities.

1. Zip:

    ```go
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
    ```

2. Concat:

    ```go
    func ExampleConcat() {
        seq1 := slices.Values([]int{1, 2, 3})
        seq2 := slices.Values([]int{4, 5})
        seq := iter2.Concat(seq1, seq2)
        fmt.Println(slices.Collect(seq))
        // Output:
        // [1 2 3 4 5]
    }
    ```
