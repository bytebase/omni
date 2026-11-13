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
	final := v1review.Schema{ObjectCaseSensitive: true, Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
		{
			Name: "public",
			Tables: []*metadata.TableMetadata{
				{Name: "t", Indexes: indexes(3), Partitions: []*metadata.TablePartitionMetadata{{Name: "p0"}}},
				{Name: "u", Indexes: indexes(1)},
				{Name: "Big", Indexes: indexes(4)},
			},
		},
		{Name: "s", Tables: []*metadata.TableMetadata{{Name: "t", Indexes: indexes(5)}}},
	}}}
	sqls := []string{
		"CREATE INDEX i ON t (a)",
		"ALTER TABLE u ADD CONSTRAINT uk UNIQUE (a)",
		"CREATE TABLE \"Big\" (a int)",
		"ALTER TABLE \"Big\" ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (a)",
		"ALTER TABLE \"Big\" ADD COLUMN c int UNIQUE",
		"CREATE TABLE missing (a int PRIMARY KEY)",
		"CREATE UNIQUE INDEX i ON p0 (a)",
		"ALTER TABLE s.t ADD PRIMARY KEY (a)",
		"CREATE TABLE t (a int, UNIQUE (a))",
		"ALTER TABLE big ADD PRIMARY KEY (a)",
	}
	got, err := IndexTotalNumberLimit(parseEach(t, sqls), 2, final)
	if err != nil {
		t.Fatal(err)
	}
	type total struct {
		Schema, Table string
		Count         int
		Statement     int
	}
	var gotTotals []total
	for _, f := range got {
		gotTotals = append(gotTotals, total{f.Schema, f.Table, countOf(f), f.Statement})
	}
	// t is marked last by its CREATE TABLE; p0 finds t's indexes; big is not
	// Big in a case sensitive snapshot.
	want := []total{{"public", "Big", 4, 4}, {"public", "p0", 3, 6}, {"s", "t", 5, 7}, {"public", "t", 3, 8}}
	if diff := cmp.Diff(want, gotTotals); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}

	if got, err := IndexTotalNumberLimit(parseEach(t, []string{"CREATE TABLE t (a int)", "DROP INDEX i"}), 0, v1review.Schema{}); err != nil || got != nil {
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
