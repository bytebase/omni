package catalog

import "testing"

// TestOracle_OceanBaseCheckExpr holds checkTableCheckExpr to the engine: a CHECK's
// columns exist and its expression is a condition, both before its name is checked.
func TestOracle_OceanBaseCheckExpr(t *testing.T) {
	runExecSteps(t, "oceanbase_catalog_check_expr", []struct {
		name  string
		steps []engineStep
	}{
		{"ALTER TABLE", []engineStep{
			{"CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT k CHECK (c > 0))", taken},
			{"ALTER TABLE a ADD CHECK (c)", 3812},
			{"ALTER TABLE a ADD CHECK ('s')", 3812},
			{"ALTER TABLE a ADD CHECK (c + 1)", 3812},
			{"ALTER TABLE a ADD CHECK (c > zz)", 3820},
			{"ALTER TABLE a ADD CONSTRAINT k CHECK (zz)", 3820},
			{"ALTER TABLE a ADD CONSTRAINT k CHECK (c)", 3812},
			{"ALTER TABLE a ADD CONSTRAINT k CHECK (c > 1)", 3822},
			{"ALTER TABLE a ADD CHECK (1)", taken},
		}},
		{"CREATE TABLE", []engineStep{
			{"CREATE TABLE t2 (b INT PRIMARY KEY, CHECK (b))", 3812},
			{"CREATE TABLE t2 (b INT PRIMARY KEY, CHECK (zz > 0))", 3820},
		}},
	})
}
