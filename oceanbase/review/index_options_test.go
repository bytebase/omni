package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkThroughOceanBaseIndexOptions: the options change nothing the walk follows, and a
// qualified index name is held to its table's database, each read from 4.3.5.6.
func TestWalkThroughOceanBaseIndexOptions(t *testing.T) {
	shop := review.Target{Schema: oceanBaseSchema(), LowerCaseTableNames: 1}
	oneFindingOB(t, "CREATE INDEX n ON t (a) PARALLEL 2 STORING (a);\nCREATE INDEX n ON t (a) NOPARALLEL;", 1, "Index `n` already exists in `t`.", shop)
	noFindingOB(t, "CREATE INDEX n ON t (a) WITH ROWID;\nDROP INDEX n ON t;", shop)
	oneFindingOB(t, "CREATE INDEX shop.n ON t (a);\nCREATE INDEX n ON t (a);", 1, "Index `n` already exists in `t`.", shop)
	oneFindingOB(t, "CREATE INDEX shop.i ON t (a);", 0, "Index `i` already exists in `t`.", shop)
	noFindingOB(t, "CREATE INDEX shop.n ON t (a);\nDROP INDEX n ON t;", shop)
	// The sync lists one database, so the walk cannot tell 1049 from 1235 here. Either way
	// the statement fails, and a database the change dropped is the 1049 it can name.
	oneFindingOB(t, "CREATE INDEX nosuch.n ON t (a);", 0, "The index is named in database `nosuch`, which `t` is not in.", shop)
	oneFindingOB(t, "CREATE DATABASE q;\nDROP DATABASE q;\nCREATE INDEX q.n ON t (a);", 2, "Database `q` doesn't exist.", shop)
	oneFindingOB(t, "CREATE DATABASE q;\nCREATE INDEX q.n ON t (a);", 1, "The index is named in database `q`, which `t` is not in.", shop)
}
