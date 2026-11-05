package review

import (
	"context"
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkThroughTiDBVector: a vector index in a change is an index the walk goes on
// from. A later statement that depends on it, or on its being new, reads it, and a
// finding after it still reports: the walk does not end at the statement, and the index
// is not lost.
func TestWalkThroughTiDBVector(t *testing.T) {
	shop := review.Target{Schema: walkSchema()}
	noFinding(t, "CREATE TABLE x (id int, b VECTOR(3), VECTOR INDEX vi ((VEC_COSINE_DISTANCE(b))));\nCREATE VECTOR INDEX vj ON x ((VEC_L2_DISTANCE(b)));\nDROP INDEX vj ON x;\nALTER TABLE x DROP INDEX vi;\nALTER TABLE x DROP COLUMN b;", shop)
	oneFinding(t, "CREATE VECTOR INDEX vi ON t ((a));\nCREATE VECTOR INDEX VI ON t ((a));", 1, "Index `VI` already exists in `t`.", shop)
	oneFinding(t, "CREATE VECTOR INDEX vi ON t ((a));\nDROP INDEX vi ON t;\nDROP INDEX vi ON t;", 2, "Index `vi` doesn't exist in `t`.", shop)
	oneFinding(t, "CREATE VECTOR INDEX I ON t ((a));", 0, "Index `I` already exists in `t`.", shop)
	noFinding(t, "CREATE VECTOR INDEX IF NOT EXISTS I ON t ((a));", shop)
	oneFinding(t, "CREATE VECTOR INDEX vi ON nosuch ((a));", 0, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "CREATE VECTOR INDEX vi ON t ((a));\nALTER TABLE nosuch ADD COLUMN c int;", 1, "Table `nosuch` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE t ADD VECTOR INDEX vi ((a));\nALTER TABLE t ADD VECTOR INDEX vi ((a));", 1, "Index `vi` already exists in `t`.", shop)
	noFinding(t, "ALTER TABLE t ADD VECTOR INDEX vi ((a));\nALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS vi ((a));\nALTER TABLE t DROP INDEX vi;", shop)
	oneFinding(t, "CREATE TABLE x (b VECTOR(3), VECTOR INDEX vi ((b)));\nCREATE INDEX vi ON x ((b));", 1, "Index `vi` already exists in `x`.", shop)
	// The column type changes nothing about what the walk follows.
	oneFinding(t, "ALTER TABLE t ADD COLUMN c VECTOR(3);\nALTER TABLE t ADD COLUMN C VECTOR;", 1, "Column `C` already exists in `t`.", shop)
	noFinding(t, "ALTER TABLE t ADD COLUMN c VECTOR(3);\nALTER TABLE t MODIFY c VECTOR(4);\nALTER TABLE t CHANGE c d VECTOR<FLOAT>(4);\nALTER TABLE t DROP COLUMN d;", shop)
}

// TestVectorShapeRefusedFinding: what TiDB's grammar actions refuse is a syntax error to
// the review, at the position the parser gives, and ends it as any other does.
func TestVectorShapeRefusedFinding(t *testing.T) {
	for _, tc := range []struct {
		sql       string
		statement int
		position  int
	}{
		{"CREATE VECTOR INDEX i ON t (a);", 0, 0},
		{"SELECT 1;\nCREATE VECTOR INDEX i ON t ((a)) USING BTREE;", 1, 10},
		{"CREATE INDEX i ON t (a) USING HNSW;", 0, 0},
		{"ALTER TABLE t ADD VECTOR INDEX i (a);", 0, 18},
		{"ALTER TABLE t ADD COLUMN c int, ADD (VECTOR INDEX i (a));", 0, 37},
		{"CREATE TABLE x (b VECTOR(3), VECTOR INDEX (b));", 0, 29},
	} {
		result, err := Review(context.Background(), tc.sql, review.Options{}, nil)
		if err != nil {
			t.Fatalf("%s: %v", tc.sql, err)
		}
		if len(result.Findings) != 1 {
			t.Errorf("%s: findings = %+v, want one", tc.sql, result.Findings)
			continue
		}
		f := result.Findings[0]
		if f.Rule != review.Syntax || len(f.Codes) != 1 || f.Codes[0] != review.SyntaxError {
			t.Errorf("%s: finding %+v, want the syntax finding", tc.sql, f)
		}
		if f.Statement != tc.statement || f.Range.Start != tc.position || f.Range.End != tc.position {
			t.Errorf("%s: at statement %d [%d,%d), want %d at %d", tc.sql, f.Statement, f.Range.Start, f.Range.End, tc.statement, tc.position)
		}
	}
	result, err := Review(context.Background(), "CREATE VECTOR INDEX i ON t ((a));", review.Options{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range result.Findings {
		if f.Rule == review.Syntax {
			t.Errorf("a good shape got %+v", f)
		}
	}
}
