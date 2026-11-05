package review

import (
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkThroughOceanBaseVector: a vector index in a change is an index the walk goes on
// from. OceanBase adds one with CREATE VECTOR INDEX and in CREATE TABLE only.
func TestWalkThroughOceanBaseVector(t *testing.T) {
	shop := review.Target{Schema: oceanBaseSchema(), LowerCaseTableNames: 1}
	const with = "WITH (distance=l2, type=hnsw, lib=vsag)"
	noFindingOB(t, "CREATE TABLE x (id int, v VECTOR(3) NOT NULL, VECTOR INDEX vi (v) "+with+");\nDROP INDEX vi ON x;\nCREATE VECTOR INDEX vj ON x (v) "+with+";\nALTER TABLE x DROP INDEX vj;\nALTER TABLE x DROP COLUMN v;", shop)
	oneFindingOB(t, "CREATE VECTOR INDEX vi ON t (a) "+with+";\nCREATE VECTOR INDEX VI ON t (a) "+with+";", 1, "Index `VI` already exists in `t`.", shop)
	oneFindingOB(t, "CREATE VECTOR INDEX vi ON t (a) "+with+";\nDROP INDEX vi ON t;\nDROP INDEX vi ON t;", 2, "Index `vi` doesn't exist in `t`.", shop)
	oneFindingOB(t, "CREATE VECTOR INDEX I ON t (a) "+with+";", 0, "Index `I` already exists in `t`.", shop)
	noFindingOB(t, "CREATE VECTOR INDEX IF NOT EXISTS I ON t (a) "+with+";", shop)
	oneFindingOB(t, "CREATE VECTOR INDEX vi ON nosuch (a) "+with+";", 0, "Table `nosuch` doesn't exist.", shop)
	oneFindingOB(t, "CREATE VECTOR INDEX vi ON t (a) "+with+";\nALTER TABLE nosuch ADD COLUMN c int;", 1, "Table `nosuch` doesn't exist.", shop)
	oneFindingOB(t, "CREATE TABLE x (v VECTOR(3) NOT NULL, VECTOR KEY vi (v) "+with+");\nCREATE INDEX vi ON x (id);", 1, "Index `vi` already exists in `x`.", shop)
	oneFindingOB(t, "ALTER TABLE t ADD COLUMN c VECTOR(3);\nALTER TABLE t ADD COLUMN C VECTOR(3);", 1, "Column `C` already exists in `t`.", shop)
}
