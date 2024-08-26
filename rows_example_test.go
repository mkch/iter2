package iter2_test

import (
	"database/sql"
	"fmt"

	"github.com/mkch/iter2"
)

func ExampleAll() {
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

	for id := range iter2.AllRows[int](r) {
		fmt.Println(*id)
	}
	// Should output:
	// 1
	// 2
	// 3
}
