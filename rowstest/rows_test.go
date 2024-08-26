package rowstest

import (
	"database/sql"
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

	seq := iter2.AllRows[int](r)
	ids := slices.Collect(iter2.Map(seq, func(p *int) int { return *p }))
	if !slices.Equal(ids, []int{1, 2, 3}) {
		t.Fatal(ids)
	}
	if err := r.Err(); err != nil {
		panic(err)
	}
}
