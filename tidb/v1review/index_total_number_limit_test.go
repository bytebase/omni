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
			{Name: "v", Indexes: indexes(5)},
		},
	}}}}
	sql := "CREATE INDEX i ON t (a);\n" +
		"ALTER TABLE u ADD INDEX i (a);\n" +
		"ALTER TABLE big ADD COLUMN c INT;\n" +
		"ALTER TABLE big ADD FOREIGN KEY (a) REFERENCES p (a);\n" +
		"ALTER TABLE big ADD VECTOR INDEX vi ((VEC_COSINE_DISTANCE(a)));\n" +
		"ALTER TABLE BIG MODIFY COLUMN c INT UNIQUE;\n" +
		"CREATE TABLE missing (a INT);\n" +
		"ALTER TABLE p0 ADD FULLTEXT INDEX f (a);\n" +
		"ALTER TABLE v ADD (c INT, KEY (c));\n" +
		"CREATE TABLE t (a INT)"
	got, err := IndexTotalNumberLimit(parseScript(t, sql), 2, final)
	if err != nil {
		t.Fatal(err)
	}
	type total struct {
		Table     string
		Count     int
		Statement int
		Ranged    bool
	}
	var gotTotals []total
	for _, f := range got {
		gotTotals = append(gotTotals, total{f.Table, countOf(f), f.Statement, f.Start != 0 || f.End != 0})
	}
	// t is marked last by its CREATE TABLE; p0 finds t's indexes too.
	want := []total{{"BIG", 4, 5, false}, {"p0", 3, 7, false}, {"t", 3, 9, false}}
	if diff := cmp.Diff(want, gotTotals); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}

	if got, err := IndexTotalNumberLimit(parseScript(t, "ALTER TABLE t DROP INDEX i"), 0, v1review.Schema{}); err != nil || got != nil {
		t.Errorf("a script that marks no table: got %v, %v", got, err)
	}
	if _, err := IndexTotalNumberLimit(parseScript(t, "CREATE INDEX i ON t (a)"), 0, v1review.Schema{}); !errors.Is(err, v1review.ErrSchemaMissing) {
		t.Errorf("err = %v, want ErrSchemaMissing", err)
	}
}

// countOf reads the count a total-count finding holds in Value.
func countOf(f v1review.Finding) int {
	n, _ := strconv.Atoi(f.Value)
	return n
}
