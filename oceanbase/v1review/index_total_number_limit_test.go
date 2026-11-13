package v1review

import (
	"errors"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

func TestIndexTotalNumberLimit(t *testing.T) {
	indexes := func(n int) []*metadata.IndexMetadata {
		var list []*metadata.IndexMetadata
		for range n {
			list = append(list, &metadata.IndexMetadata{})
		}
		return list
	}
	final := v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{
			{Name: "t", Indexes: indexes(3), Partitions: []*metadata.TablePartitionMetadata{{Name: "p0"}}},
			{Name: "u", Indexes: indexes(1)},
			{Name: "Big", Indexes: indexes(4)},
		},
	}}}}
	sqls := []string{
		"CREATE INDEX i ON t (a)",
		"ALTER TABLE u ADD INDEX i (a)",
		"ALTER TABLE big ADD COLUMN c INT",
		"ALTER TABLE big ADD SPATIAL INDEX s (g)",
		"ALTER TABLE BIG MODIFY COLUMN c INT UNIQUE",
		"CREATE TABLE missing (a INT)",
		"ALTER TABLE p0 ADD FULLTEXT INDEX f (a)",
		"CREATE TABLE t (a INT)",
	}
	got, err := IndexTotalNumberLimit(parseEach(t, sqls), 2, final)
	if err != nil {
		t.Fatal(err)
	}
	type total struct {
		Table     string
		Count     int
		Statement int
	}
	var gotTotals []total
	for _, f := range got {
		gotTotals = append(gotTotals, total{f.Table, countOf(f), f.Statement})
	}
	// t is marked last by its CREATE TABLE; p0 finds t's indexes too.
	want := []total{{"BIG", 4, 4}, {"p0", 3, 6}, {"t", 3, 7}}
	if diff := cmp.Diff(want, gotTotals); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}

	if got, err := IndexTotalNumberLimit(parseEach(t, []string{"ALTER TABLE t DROP INDEX i"}), 0, v1review.Schema{}); err != nil || got != nil {
		t.Errorf("a script that marks no table: got %v, %v", got, err)
	}
	if _, err := IndexTotalNumberLimit(parseEach(t, sqls[:1]), 0, v1review.Schema{}); !errors.Is(err, v1review.ErrSchemaMissing) {
		t.Errorf("err = %v, want ErrSchemaMissing", err)
	}
}

// countOf reads the count a total-count finding holds in Value.
func countOf(f v1review.Finding) int {
	n, _ := strconv.Atoi(f.Value)
	return n
}
