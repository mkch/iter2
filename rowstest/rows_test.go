package rowstest

import (
	"database/sql"
	"maps"
	"slices"
	"testing"

	"github.com/mkch/iter2"
	_ "modernc.org/sqlite"
)

func TestAll(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	q := `
create temp table uid (id bigint); -- Create temp table for queries.
insert into uid values (1); -- Populate temp table.
insert into uid values (2);
insert into uid values (3);

-- First result set.
select * from uid;
`

	r, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	seq := iter2.AllRows(r)
	ids := slices.Collect(iter2.Map(seq, func(row iter2.Row) int {
		var id int
		row.Scan(&id)
		return id
	}))
	if !slices.Equal(ids, []int{1, 2, 3}) {
		t.Fatal(ids)
	}
	if err := r.Err(); err != nil {
		panic(err)
	}
}

func TestTableMap(t *testing.T) {
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

	r, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	seq2 := iter2.Map1To2(iter2.AllRows(r), func(row iter2.Row) (int, string) {
		var id int
		var name string
		row.Scan(&id, &name)
		return id, name
	})
	m := maps.Collect(seq2)

	if !maps.Equal(m, map[int]string{1: "User1", 2: "User2", 3: "User3"}) {
		t.Fatal(m)
	}
}

func TestTableMapStruct(t *testing.T) {
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

	r, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	type User struct {
		ID   int
		Name string
	}

	seq := iter2.Map(iter2.AllRows(r), func(row iter2.Row) User {
		var id int
		var name string
		row.Scan(&id, &name)
		return User{id, name}
	})
	users := slices.Collect(seq)

	if !slices.Equal(users, []User{{1, "User1"}, {2, "User2"}, {3, "User3"}}) {
		t.Fatal(users)
	}
}
