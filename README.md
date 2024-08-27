# iter2

 Go iter utilities.

1. Operations of iterator:

    Zipping two Seqs to one Seq of pairs:

    ```go
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
    ```

    Zipping tow Seqs to one Seq2:

    ```go
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
    ```

    Concating 2 Seqs:

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

2. Map DB rows:

    ```go
    import (
        "database/sql"
        _ "modernc.org/sqlite"
        // and more
    )

    func ExampleAllRows_tableStruct() {
        db, err := sql.Open("sqlite", ":memory:")
        if err != nil {
            panic(err)
        }
        q := `
    create temp table users (id integer, name text); -- Create temp table for queries.
    insert into users values (1, "User1"); -- Populate temp table.
    insert into users values (2, "User2");
    insert into users values (3, "User3");

    -- First result set.
    select * from users;
    `
        type User struct {
            ID   int
            Name string
        }

        users := slices.Collect(
            iter2.Map(
                iter2.MustAllRows(db.Query(q)), func(row iter2.Row) (user User) {
                    row.Scan(&user.ID, &user.Name)
                    return
                }))
        fmt.Println(users)
        // Should output:
        // [{1 User1} {2 User2} {3 User3}]
    }
    ```

More [examples](https://pkg.go.dev/github.com/mkch/iter2#pkg-examples) here.
