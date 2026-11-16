package driver

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/cosmosdb/internal/testcontainer"
)

// The multi-page fixture: the queries whose results a range can only deliver
// over several pages, and the operators that compose with each other in one
// plan. Both need a container whose rows pass the gateway's per-page cap.

// pagedSeqs runs a query against the multi-page fixture and returns the seq of
// every row it emitted, in the order the query engine emitted them.
func pagedSeqs(t *testing.T, d *Driver, statement string) []int {
	t.Helper()
	rows := query(t, d, testcontainer.PagedContainer, statement)
	seqs := make([]int, 0, len(rows))
	for _, row := range rows {
		var value struct {
			Seq *int `json:"seq"`
		}
		require.NoError(t, json.Unmarshal(row, &value), "row: %s", row)
		require.NotNil(t, value.Seq, "row missing seq: %s", row)
		seqs = append(seqs, *value.Seq)
	}
	return seqs
}

// pagedQuery runs a query against the multi-page fixture and returns its rows
// as raw JSON, in the order the query engine emitted them.
func pagedQuery(t *testing.T, d *Driver, statement string) []string {
	t.Helper()
	rows := query(t, d, testcontainer.PagedContainer, statement)
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, string(row))
	}
	return out
}

// ascendingSeqs is the fixture's seq values in order: every query below either
// returns this sequence or a slice of it.
func ascendingSeqs(from, to int) []int {
	seqs := make([]int, 0, to-from)
	for seq := from; seq < to; seq++ {
		seqs = append(seqs, seq)
	}
	return seqs
}

// TestCrossPartitionMultiPage covers the query classes whose results a range
// can only deliver over several pages. Each one reads past the gateway's
// per-page cap, which an engine that treats a continuation token as the end of
// a range cannot do: before the engine followed them, an unbounded ORDER BY
// failed outright on any container holding more than PageSize rows in one
// range, and so did DISTINCT and GROUP BY over that many values and a TOP
// above the cap.
func TestCrossPartitionMultiPage(t *testing.T) {
	d := newDriver(t)
	rows := testcontainer.PagedRows()

	t.Run("order_by_unbounded", func(t *testing.T) {
		assert.Equal(t, ascendingSeqs(0, rows),
			pagedSeqs(t, d, `SELECT c.seq FROM c ORDER BY c.seq ASC`))
	})

	t.Run("order_by_descending", func(t *testing.T) {
		want := ascendingSeqs(0, rows)
		slices.Reverse(want)
		assert.Equal(t, want, pagedSeqs(t, d, `SELECT c.seq FROM c ORDER BY c.seq DESC`))
	})

	// A page of a range ends mid-result, so the rows either side of the cap
	// are the ones a resumed range can misplace.
	t.Run("offset_limit_across_a_page_boundary", func(t *testing.T) {
		assert.Equal(t, ascendingSeqs(testcontainer.PageSize-5, testcontainer.PageSize+5),
			pagedSeqs(t, d, fmt.Sprintf(
				`SELECT c.seq FROM c ORDER BY c.seq ASC OFFSET %d LIMIT 10`, testcontainer.PageSize-5)))
	})

	t.Run("distinct_high_cardinality", func(t *testing.T) {
		seqs := pagedSeqs(t, d, `SELECT DISTINCT c.seq FROM c`)
		slices.Sort(seqs)
		assert.Equal(t, ascendingSeqs(0, rows), seqs)
	})

	t.Run("group_by_many_groups", func(t *testing.T) {
		rows := pagedQuery(t, d, `SELECT c.seq, COUNT(1) AS n FROM c GROUP BY c.seq`)
		seqs := make([]int, 0, len(rows))
		for _, row := range rows {
			var group struct {
				Seq int `json:"seq"`
				N   int `json:"n"`
			}
			require.NoError(t, json.Unmarshal([]byte(row), &group), "row: %s", row)
			assert.Equal(t, 1, group.N, "seq %d", group.Seq)
			seqs = append(seqs, group.Seq)
		}
		slices.Sort(seqs)
		assert.Equal(t, ascendingSeqs(0, testcontainer.PagedRows()), seqs)
	})

	// An aggregate answers with one partial document per range however many
	// rows that range holds, which is why aggregates were the one class this
	// fixture's size never broke. Kept as the control case.
	t.Run("aggregate_over_the_whole_fixture", func(t *testing.T) {
		rows := pagedQuery(t, d, `SELECT COUNT(1) AS n, SUM(c.seq) AS total FROM c`)
		require.Len(t, rows, 1)
		want := 0
		for _, seq := range ascendingSeqs(0, testcontainer.PagedRows()) {
			want += seq
		}
		assert.JSONEq(t, fmt.Sprintf(`{"n":%d,"total":%d}`, testcontainer.PagedRows(), want), rows[0])
	})
}

// pagedBucketValues reads the bucket out of every row, for both spellings of a
// projection: `{"bucket":"big-1"}` and the VALUE form's bare `"big-1"`.
func pagedBucketValues(t *testing.T, rows []string) []string {
	t.Helper()
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		var object struct {
			Bucket *string `json:"bucket"`
		}
		if err := json.Unmarshal([]byte(row), &object); err == nil && object.Bucket != nil {
			out = append(out, *object.Bucket)
			continue
		}
		var scalar string
		require.NoError(t, json.Unmarshal([]byte(row), &scalar), "row: %s", row)
		out = append(out, scalar)
	}
	return out
}

// TestCrossPartitionDistinctComposed covers DISTINCT alongside an ORDER BY or
// a page. The service reports distinctType "Unordered" for the object form
// even under ORDER BY, so the engine dropped the DISTINCT and returned every
// row in the container; a paged DISTINCT ran the user's own query per range,
// so each range applied the LIMIT to itself before anything was deduplicated.
func TestCrossPartitionDistinctComposed(t *testing.T) {
	d := newDriver(t)
	buckets := testcontainer.SortedPagedBuckets()

	t.Run("order_by_object_form", func(t *testing.T) {
		assert.Equal(t, buckets, pagedBucketValues(t,
			pagedQuery(t, d, `SELECT DISTINCT c.bucket FROM c ORDER BY c.bucket`)))
	})

	t.Run("order_by_value_form", func(t *testing.T) {
		assert.Equal(t, buckets, pagedBucketValues(t,
			pagedQuery(t, d, `SELECT DISTINCT VALUE c.bucket FROM c ORDER BY c.bucket`)))
	})

	t.Run("order_by_descending", func(t *testing.T) {
		want := testcontainer.SortedPagedBuckets()
		slices.Reverse(want)
		assert.Equal(t, want, pagedBucketValues(t,
			pagedQuery(t, d, `SELECT DISTINCT c.bucket FROM c ORDER BY c.bucket DESC`)))
	})

	// An ordered DISTINCT pins the page against exact values; the unordered
	// forms below can only be held to the count, since Cosmos defines no order
	// for them.
	t.Run("order_by_with_a_page", func(t *testing.T) {
		assert.Equal(t, buckets[2:5], pagedBucketValues(t,
			pagedQuery(t, d, `SELECT DISTINCT c.bucket FROM c ORDER BY c.bucket OFFSET 2 LIMIT 3`)))
	})

	t.Run("limit_value_form", func(t *testing.T) {
		got := pagedBucketValues(t,
			pagedQuery(t, d, `SELECT DISTINCT VALUE c.bucket FROM c OFFSET 0 LIMIT 3`))
		assert.Len(t, got, 3)
		assert.Subset(t, buckets, got)
		assert.Len(t, slices.Compact(slices.Sorted(slices.Values(got))), 3, "values must be distinct")
	})

	t.Run("offset_limit_object_form", func(t *testing.T) {
		got := pagedBucketValues(t,
			pagedQuery(t, d, `SELECT DISTINCT c.bucket FROM c OFFSET 2 LIMIT 3`))
		assert.Len(t, got, 3)
		assert.Subset(t, buckets, got)
		assert.Len(t, slices.Compact(slices.Sorted(slices.Values(got))), 3, "values must be distinct")
	})
}

// TestCrossPartitionTopWithOrderBy covers TOP alongside an ORDER BY. The
// rewritten per-range query carries the TOP, so the merge saw n rows from
// every range and emitted all of them: n x partitions rows, in order.
func TestCrossPartitionTopWithOrderBy(t *testing.T) {
	d := newDriver(t)
	rows := testcontainer.PagedRows()

	t.Run("ascending", func(t *testing.T) {
		assert.Equal(t, ascendingSeqs(0, 5),
			pagedSeqs(t, d, `SELECT TOP 5 c.seq FROM c ORDER BY c.seq ASC`))
	})

	t.Run("descending", func(t *testing.T) {
		want := ascendingSeqs(rows-5, rows)
		slices.Reverse(want)
		assert.Equal(t, want, pagedSeqs(t, d, `SELECT TOP 5 c.seq FROM c ORDER BY c.seq DESC`))
	})

	// A TOP past the gateway's page cap has to keep reading the ranges after
	// the cap, and still stop at n.
	t.Run("above_the_page_cap", func(t *testing.T) {
		assert.Equal(t, ascendingSeqs(0, testcontainer.PageSize+25),
			pagedSeqs(t, d, fmt.Sprintf(
				`SELECT TOP %d c.seq FROM c ORDER BY c.seq ASC`, testcontainer.PageSize+25)))
	})
}

// TestCrossPartitionGroupByOrder covers the order GROUP BY rows come back in.
// The pipeline emitted them straight out of a Go map, so the same query
// answered in a different order every run, and a GROUP BY the user ordered
// reaches the engine as a plan with no order fields at all — the ORDER BY
// survives only inside the rewritten query.
func TestCrossPartitionGroupByOrder(t *testing.T) {
	d := newDriver(t)

	type group struct {
		Bucket string `json:"bucket"`
		N      int    `json:"n"`
	}
	groupsOf := func(t *testing.T, rows []string) []group {
		t.Helper()
		out := make([]group, 0, len(rows))
		for _, row := range rows {
			var g group
			require.NoError(t, json.Unmarshal([]byte(row), &g), "row: %s", row)
			out = append(out, g)
		}
		return out
	}
	want := make([]group, 0, len(testcontainer.SortedPagedBuckets()))
	for _, bucket := range testcontainer.SortedPagedBuckets() {
		want = append(want, group{Bucket: bucket, N: testcontainer.PagedBucketRows(bucket)})
	}

	t.Run("order_by_the_group_key", func(t *testing.T) {
		assert.Equal(t, want, groupsOf(t, pagedQuery(t, d,
			`SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket ORDER BY c.bucket`)))
	})

	t.Run("order_by_descending", func(t *testing.T) {
		reversed := slices.Clone(want)
		slices.Reverse(reversed)
		assert.Equal(t, reversed, groupsOf(t, pagedQuery(t, d,
			`SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket ORDER BY c.bucket DESC`)))
	})

	// Without an ORDER BY the order is Cosmos's to choose, but it must be the
	// same on every run.
	t.Run("unordered_is_stable_across_runs", func(t *testing.T) {
		const sql = `SELECT c.bucket, COUNT(1) AS n FROM c GROUP BY c.bucket`
		first := groupsOf(t, pagedQuery(t, d, sql))
		assert.ElementsMatch(t, want, first)
		for i := 0; i < 2; i++ {
			assert.Equal(t, first, groupsOf(t, pagedQuery(t, d, sql)))
		}
	})
}

// TestCrossPartitionEmptyAggregates covers an aggregate whose filter matches
// no row. MIN and MAX failed the query outright, because an empty partition's
// partial is an absent one; the aliased form answered `{"m":null}` where the
// gateway omits the alias. COUNT and SUM really are 0 over an empty set and
// are here as controls.
func TestCrossPartitionEmptyAggregates(t *testing.T) {
	d := newDriver(t)
	const noRows = `WHERE c.seq > 99999999`

	for _, tc := range []struct {
		name string
		sql  string
		want []string
	}{
		{"min_value", `SELECT VALUE MIN(c.seq) FROM c ` + noRows, nil},
		{"max_value", `SELECT VALUE MAX(c.seq) FROM c ` + noRows, nil},
		{"avg_value", `SELECT VALUE AVG(c.seq) FROM c ` + noRows, nil},
		{"min_aliased", `SELECT MIN(c.seq) AS m FROM c ` + noRows, []string{`{}`}},
		{"min_max_aliased", `SELECT MIN(c.seq) AS lo, MAX(c.seq) AS hi FROM c ` + noRows, []string{`{}`}},
		{"count_value", `SELECT VALUE COUNT(1) FROM c ` + noRows, []string{`0`}},
		{"count_aliased", `SELECT COUNT(1) AS n FROM c ` + noRows, []string{`{"n":0}`}},
		{"sum_value", `SELECT VALUE SUM(c.seq) FROM c ` + noRows, []string{`0`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := pagedQuery(t, d, tc.sql)
			require.Len(t, got, len(tc.want))
			for i := range tc.want {
				assert.JSONEq(t, tc.want[i], got[i])
			}
		})
	}

	// The same MIN and MAX over rows that do match: the VALUE form reads its
	// partial from item2, which is where the count lives.
	t.Run("min_max_over_rows_that_match", func(t *testing.T) {
		assert.Equal(t, []string{`0`}, pagedQuery(t, d, `SELECT VALUE MIN(c.seq) FROM c`))
		assert.Equal(t, []string{fmt.Sprint(testcontainer.PagedRows() - 1)},
			pagedQuery(t, d, `SELECT VALUE MAX(c.seq) FROM c`))
	})
}

// TestCrossPartitionOffsetLimitAcrossPages covers ordinary result paging,
// which the engine rejected although it advertises OffsetAndLimit, and a window
// wider than the gateway's page cap, which reads several pages of a range
// before it is full. Both queries order their rows: a page over an unordered
// result keeps whichever rows the merge reached first, so the answer would be
// the server's choice and not the fixture's.
func TestCrossPartitionOffsetLimitAcrossPages(t *testing.T) {
	d := newDriver(t)

	assert.Equal(t, ascendingSeqs(2, 5),
		pagedSeqs(t, d, `SELECT c.seq FROM c ORDER BY c.seq OFFSET 2 LIMIT 3`))
	assert.Equal(t, ascendingSeqs(0, testcontainer.PageSize+25),
		pagedSeqs(t, d, fmt.Sprintf(`SELECT c.seq FROM c ORDER BY c.seq OFFSET 0 LIMIT %d`, testcontainer.PageSize+25)))
}

// TestCrossPartitionAggregateWithOffsetLimit covers a page applied to an
// aggregate. The query produces exactly one row, so OFFSET 1 skips past it and
// the gateway answers with nothing; the engine returned the row regardless.
func TestCrossPartitionAggregateWithOffsetLimit(t *testing.T) {
	d := newDriver(t)

	assert.Empty(t, pagedQuery(t, d, `SELECT VALUE COUNT(1) FROM c OFFSET 1 LIMIT 1`))
	assert.Equal(t, []string{fmt.Sprint(testcontainer.PagedRows())},
		pagedQuery(t, d, `SELECT VALUE COUNT(1) FROM c OFFSET 0 LIMIT 1`))
	assert.Empty(t, pagedQuery(t, d, `SELECT COUNT(1) AS n FROM c OFFSET 1 LIMIT 1`))
}

// TestCrossPartitionDistinctWithOrderByThroughTheDriver runs the worst of these
// through Query, the entry point the SQL editor uses, so the statement
// validation is in the path too.
func TestCrossPartitionDistinctWithOrderByThroughTheDriver(t *testing.T) {
	d := newDriver(t)
	rows := query(t, d, testcontainer.CitiesContainer, `SELECT DISTINCT c.country FROM c ORDER BY c.country`)
	got := make([]string, 0, len(rows))
	for _, row := range rows {
		var object struct {
			Country string `json:"country"`
		}
		require.NoError(t, json.Unmarshal(row, &object))
		got = append(got, object.Country)
	}
	assert.Equal(t, slices.Sorted(slices.Values(testcontainer.Cities)), got)
}

// TestCrossPartitionSparseGroupBy groups on a column the documents do not all
// define, which is ordinary in a schemaless store. The service answers the
// undefined rows as a group of their own, counted apart from an explicit null
// and spelled with the alias left out of the row.
func TestCrossPartitionSparseGroupBy(t *testing.T) {
	d := newDriver(t)

	want := []string{fmt.Sprintf(`{"n":%d}`, testcontainer.PagedRows()-len(testcontainer.PagedShapes))}
	for _, shape := range testcontainer.PagedShapes {
		want = append(want, fmt.Sprintf(`{"k":%s,"n":1}`, shape.Value))
	}
	rows := pagedQuery(t, d, `SELECT c.v AS k, COUNT(1) AS n FROM c GROUP BY c.v`)
	assert.ElementsMatch(t, want, rows)

	// The groups also come back in key order, the undefined one first, which is
	// what makes two runs of one query agree.
	require.NotEmpty(t, rows)
	assert.Equal(t, want[0], rows[0])
	assert.ElementsMatch(t, []string{`{"k":{"a":1},"n":1}`, `{"k":{"b":1},"n":1}`}, rows[len(rows)-2:],
		"objects sort last")
}

// TestCrossPartitionMixedTypeOrder checks ORDER BY against the order the
// service itself returns: undefined < null < bool < number < string < array <
// object.
//
// Two different arrays, or two different objects, it cannot order when they
// reach the merge from different ranges, because the service orders those by
// its index's own encoding — [1,2] < [2] < [1,3] on a real container. Which
// range holds which value is the emulator's to decide, so the engine's refusal
// is covered in merge_test.go, where the ranges are given.
func TestCrossPartitionMixedTypeOrder(t *testing.T) {
	d := newDriver(t)

	scalars := make([]int, 0, len(testcontainer.PagedShapes))
	for _, shape := range testcontainer.PagedShapes {
		if shape.Value[0] != '[' && shape.Value[0] != '{' {
			scalars = append(scalars, shape.Seq)
		}
	}
	got := pagedSeqs(t, d,
		`SELECT c.seq FROM c WHERE IS_DEFINED(c.v) AND NOT IS_ARRAY(c.v) AND NOT IS_OBJECT(c.v) ORDER BY c.v`)
	assert.Equal(t, scalars, got, "the scalar classes are ordered as the service orders them")

	// One composite among them still orders, since the class ranking is the
	// service's own: an array sorts after every string.
	withOneArray := append(slices.Clone(scalars), testcontainer.PagedShapes[7].Seq)
	got = pagedSeqs(t, d, fmt.Sprintf(
		`SELECT c.seq FROM c WHERE IS_DEFINED(c.v) AND NOT IS_OBJECT(c.v) AND (NOT IS_ARRAY(c.v) OR c.seq = %d) ORDER BY c.v`,
		testcontainer.PagedShapes[7].Seq))
	assert.Equal(t, withOneArray, got, "an array sorts after the scalars")
}

// ------------------------------- assertion helpers carried from Bytebase

// assertOrderedByPopulation walks rows and asserts they are sorted by
// c.population under Cosmos item ordering (undefined < null < bool < number <
// string). ascending=true checks nondecreasing; ascending=false checks
// nonincreasing.
func assertOrderedByPopulation(t *testing.T, rows [][]byte, ascending bool) {
	t.Helper()
	var prevRank int
	var prevVal any
	for i, row := range rows {
		var object map[string]any
		require.NoError(t, json.Unmarshal(row, &object), "row %d: %s", i, row)
		rank, value := cosmosTypeRank(object["population"])
		if i > 0 {
			if ascending {
				assert.LessOrEqual(t, prevRank, rank, "rank broken at row %d", i)
			} else {
				assert.GreaterOrEqual(t, prevRank, rank, "rank broken at row %d", i)
			}
			if prevRank == rank {
				if ascending {
					assert.LessOrEqual(t, prevVal, value, "value broken at row %d", i)
				} else {
					assert.GreaterOrEqual(t, prevVal, value, "value broken at row %d", i)
				}
			}
		}
		prevRank, prevVal = rank, value
	}
}

// assertAscendingByName walks rows and asserts their name field is ascending,
// returning the names.
func assertAscendingByName(t *testing.T, rows [][]byte) []string {
	t.Helper()
	names := make([]string, 0, len(rows))
	var prev string
	for i, row := range rows {
		var object map[string]any
		require.NoError(t, json.Unmarshal(row, &object), "row %d: %s", i, row)
		name, ok := object["name"].(string)
		require.True(t, ok, "row %d missing name: %s", i, row)
		names = append(names, name)
		if i > 0 {
			assert.LessOrEqual(t, prev, name, "ascending broken at row %d: prev=%q cur=%q", i, prev, name)
		}
		prev = name
	}
	return names
}

// cosmosTypeRank tags a JSON value with its Cosmos-ordering type class. Used to
// compare rows that straddle the number→string boundary: Cosmos orders
// undefined < null < bool < number < string.
func cosmosTypeRank(v any) (int, any) {
	switch x := v.(type) {
	case float64:
		return 3, x
	case string:
		return 4, x
	case bool:
		return 2, x
	case nil:
		return 1, x
	default:
		return 0, x
	}
}
