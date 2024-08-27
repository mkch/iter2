package iter2_test

import (
	"database/sql"
	"fmt"
	"slices"

	"github.com/mkch/iter2"
)

func ExampleAllRows() {
	var db *sql.DB // Open db.
	q := `
-- Create temp table for queries.
create temp table uid (id bigint);
-- Populate temp table. 
insert into uid values (1); 
insert into uid values (2);
insert into uid values (3);
-- Do query.
select * from uid order by id;
`
	r, err := db.Query(q)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	for row := range iter2.AllRows(r) {
		var id int
		row.Scan(&id)
		fmt.Println(id)
	}
	// Should output:
	// 1
	// 2
	// 3
}

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
			iter2.MustAllRows(db.Query(q)), func(row iter2.Row) User {
				var id int
				var name string
				row.Scan(&id, &name)
				return User{id, name}
			}))
	fmt.Println(users)
	// Should output:
	// [{1 User1} {2 User2} {3 User3}]
}
