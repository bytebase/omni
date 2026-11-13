package catalog

import (
	"errors"
	"testing"
)

// TestOracle_OceanBaseKeyParts holds checkKeyParts to the engine: a key's block size
// and parts are refused before its name is looked up.
func TestOracle_OceanBaseKeyParts(t *testing.T) {
	const table = "CREATE TABLE a (b INT PRIMARY KEY, c INT, v VARCHAR(10), KEY idx_c (c))"
	runExecSteps(t, "oceanbase_catalog_key_parts", []struct {
		name  string
		steps []engineStep
	}{
		{"ALTER TABLE and CREATE INDEX", []engineStep{
			{table, taken},
			{"ALTER TABLE a ADD KEY idx_c (c, c)", 1060},
			{"ALTER TABLE a ADD KEY (zz, zz)", 1060},
			{"ALTER TABLE a ADD KEY ((c))", 3762},
			{"ALTER TABLE a ADD KEY (c(1))", 1089},
			{"ALTER TABLE a ADD KEY (c) BLOCK_SIZE 1", 5258},
			{"ALTER TABLE a ADD SPATIAL (c)", 1687},
			{"CREATE SPATIAL INDEX idx_c ON a (c)", 1687},
			{"CREATE VECTOR INDEX idx_c ON a (c)", 7601},
			{"CREATE FULLTEXT INDEX idx_c ON a (c)", 1283},
			{"CREATE INDEX i ON a (c, c)", 1060},
			{"ALTER TABLE a ADD KEY ((c + 1)), ADD KEY (v(1)) BLOCK_SIZE 16384", taken},
		}},
		{"CREATE TABLE", []engineStep{
			{"CREATE TABLE t2 (b INT PRIMARY KEY, KEY (b, b))", 1060},
			{"CREATE TABLE t2 (b INT PRIMARY KEY, KEY ((b)))", 3762},
			{"CREATE TABLE t2 (b INT PRIMARY KEY, KEY (b(1)))", 1089},
			{"CREATE TABLE t2 (b INT PRIMARY KEY, KEY (b) BLOCK_SIZE 1)", 5258},
		}},
	})
}

// TestKeyPartsSQLState holds each refusal's SQLSTATE to what
// TestOracle_OceanBaseKeyParts's tape recorded, which runExecSteps does not compare.
func TestKeyPartsSQLState(t *testing.T) {
	for _, tc := range []struct {
		sql   string
		code  int
		state string
	}{
		{"ALTER TABLE a ADD KEY (c, c)", 1060, "42S21"},
		{"ALTER TABLE a ADD KEY ((c))", 3762, "HY000"},
		{"ALTER TABLE a ADD KEY (c(1))", 1089, "HY000"},
		{"ALTER TABLE a ADD KEY (c) BLOCK_SIZE 1", 5258, "HY000"},
		{"ALTER TABLE a ADD SPATIAL (c)", 1687, "42000"},
		{"CREATE VECTOR INDEX i ON a (c)", 7601, "HY000"},
		{"CREATE FULLTEXT INDEX i ON a (c)", 1283, "42000"},
	} {
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE TABLE a (b INT PRIMARY KEY, c INT)")
		if err != nil {
			t.Fatal(err)
		}
		results, err := c.Exec(tc.sql, nil)
		if err != nil {
			t.Fatalf("%s: %v", tc.sql, err)
		}
		var refusal *Error
		if !errors.As(results[0].Error, &refusal) || refusal.Code != tc.code || refusal.SQLState != tc.state {
			t.Errorf("%s: got %v, want %d (%s)", tc.sql, results[0].Error, tc.code, tc.state)
		}
	}
}
