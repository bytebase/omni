package catalog

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

// TestSimulationCausesOracle holds the refusals and stored forms
// testdata/simulate.tsv records to the pinned OceanBase cause by cause: a
// repartition applied instead of ignored, its repeated partition name and its
// primary-key prefix rule, a character set or collation the server does not
// have, and the forms OceanBase stores differently from the catalog (SERIAL,
// and a key whose STORING repeats a column it keys on). Each case is answered
// by the pinned server and by the catalog, and both must agree with the code
// the case states; a case the server takes also compares what the two store.
func TestSimulationCausesOracle(t *testing.T) {
	conn := testcontainer.OpenEmpty(t)
	ctx := context.Background()
	database := "omni_t_simulation"
	reset := func() {
		if _, err := conn.ExecContext(ctx, "DROP DATABASE IF EXISTS "+database); err != nil {
			t.Fatalf("drop %s: %v", database, err)
		}
	}
	t.Cleanup(reset)

	oneTable := "CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))"
	// A table whose primary key starts with b, so a partitioning named a is not
	// a prefix of it.
	aColumnTable := "CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))"
	cases := []struct {
		name  string
		setup []string
		stmt  string
		code  int
		state string
	}{
		{
			name:  "repartition by an implicit key",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY KEY()",
		},
		{
			name:  "repartition by an implicit key with a named partition",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY KEY()(PARTITION a)",
		},
		{
			name:  "repartition by automatic range",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY RANGE()",
		},
		{
			name:  "repartition by automatic range columns",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY RANGE COLUMNS()",
		},
		{
			name:  "repartition by an explicit range key with no list",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY RANGE(b)",
		},
		{
			name:  "repartition by string bounds one of which spells MAXVALUE",
			setup: []string{"CREATE TABLE a (s VARCHAR(10) PRIMARY KEY, c INT)"},
			stmt:  "ALTER TABLE a PARTITION BY RANGE COLUMNS(s)(PARTITION p0 VALUES LESS THAN ('MAXVALUE'), PARTITION p1 VALUES LESS THAN (MAXVALUE))",
		},
		{
			name:  "a SERIAL column with an explicit NULL",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a SERIAL NULL",
			code:  5252,
			state: "HY000",
		},
		{
			name:  "a collation written in another case",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4 COLLATE UTF8MB4_GENERAL_CI",
		},
		{
			name:  "repartition with a repeated partition name",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a PARTITION BY KEY()(PARTITION a, PARTITION a)",
			code:  1517,
			state: "HY000",
		},
		{
			name:  "repartition by a range key that is not the primary key prefix",
			setup: []string{aColumnTable},
			stmt:  "ALTER TABLE a PARTITION BY RANGE(a)",
			code:  1235,
			state: "0A000",
		},
		{
			name:  "a character set the server does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a CHARSET a",
			code:  1115,
			state: "42000",
		},
		{
			name:  "a converted character set the server does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a CONVERT TO CHARACTER SET a",
			code:  1115,
			state: "42000",
		},
		{
			name:  "a collation the server does not have",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COLLATE a",
			code:  1273,
			state: "HY000",
		},
		{
			name:  "a SERIAL column added to a table",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a ADD a SERIAL",
		},
		{
			name:  "a key storing a column it already keys on",
			setup: []string{aColumnTable},
			stmt:  "ALTER TABLE a ADD KEY(a) STORING(a)",
		},
		{
			name:  "a collation the server does not have under a charset it does",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a COLLATE utf8mb4_not_a_collation",
			code:  1273,
			state: "HY000",
		},
		{
			name:  "a collation of another character set",
			setup: []string{oneTable},
			stmt:  "ALTER TABLE a CONVERT TO CHARACTER SET latin1 COLLATE utf8mb4_general_ci",
			code:  1253,
			state: "42000",
		},
		{
			name:  "a SERIAL replacement of an unindexed column",
			setup: []string{aColumnTable},
			stmt:  "ALTER TABLE a MODIFY a SERIAL",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			build := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=utf8mb4; USE %s;", database, database)
			if _, err := conn.ExecContext(ctx, build); err != nil {
				t.Fatalf("%s: %v", build, err)
			}
			for _, stmt := range tc.setup {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("%s: %v", stmt, err)
				}
			}
			_, err := conn.ExecContext(ctx, tc.stmt)
			engineCode, engineState, engineMessage := simulationCauseEngineRefusal(t, tc.stmt, err)
			if engineCode != tc.code {
				t.Fatalf("oceanbase: %s answered %d (%s), want %d", tc.stmt, engineCode, engineMessage, tc.code)
			}
			if tc.state != "" && engineState != tc.state {
				t.Fatalf("oceanbase: %s answered SQLSTATE %s, want %s", tc.stmt, engineState, tc.state)
			}

			c, err := LoadSQL(build + "\n" + joinStatements(tc.setup))
			if err != nil {
				t.Fatalf("the catalog does not load the fixture: %v", err)
			}
			catalogCode, catalogState, catalogMessage := simulationCauseCatalogRefusal(t, c, tc.stmt)
			if catalogCode != engineCode {
				t.Fatalf("catalog: %s answered %d (%s), OceanBase answered %d (%s)",
					tc.stmt, catalogCode, catalogMessage, engineCode, engineMessage)
			}
			if catalogState != engineState {
				t.Fatalf("catalog: %s answered SQLSTATE %s, OceanBase answered %s", tc.stmt, catalogState, engineState)
			}
			if tc.code != 0 {
				return
			}

			var name, engineCreate string
			if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE a").Scan(&name, &engineCreate); err != nil {
				t.Fatalf("SHOW CREATE TABLE a: %v", err)
			}
			want, err := LoadSQL(build + "\n" + engineCreate)
			if err != nil {
				t.Fatalf("the catalog does not load what OceanBase stores: %v\n%s", err, engineCreate)
			}
			n := NormalizerFor(MySQL80)
			if d := DiffWithNormalizer(c, want, n); !d.IsEmpty() {
				t.Fatalf("%s: the catalog does not store what OceanBase stores: %s\nOceanBase:\n%s", tc.stmt, describeDiff(d), engineCreate)
			}
			if d := DiffWithNormalizer(want, c, n); !d.IsEmpty() {
				t.Fatalf("%s: the catalog stores what OceanBase does not: %s\nOceanBase:\n%s", tc.stmt, describeDiff(d), engineCreate)
			}
		})
	}
}

func joinStatements(stmts []string) string {
	out := ""
	for _, stmt := range stmts {
		out += stmt + ";\n"
	}
	return out
}

// simulationCauseEngineRefusal reads the server's answer as the error number
// and message, 0 for a statement it took.
func simulationCauseEngineRefusal(t *testing.T, stmt string, err error) (int, string, string) {
	t.Helper()
	if err == nil {
		return 0, "", ""
	}
	var refusal *gomysql.MySQLError
	if !errors.As(err, &refusal) {
		t.Fatalf("%s: not an OceanBase error: %v", stmt, err)
	}
	return int(refusal.Number), string(refusal.SQLState[:]), refusal.Message
}

// simulationCauseCatalogRefusal reads the catalog's answer as the error number
// and message, 0 for a statement it took.
func simulationCauseCatalogRefusal(t *testing.T, c *Catalog, stmt string) (int, string, string) {
	t.Helper()
	results, err := c.Exec(stmt, nil)
	if err != nil {
		t.Fatalf("Exec(%q): %v", stmt, err)
	}
	if len(results) != 1 {
		t.Fatalf("Exec(%q): %d results, want 1", stmt, len(results))
	}
	if results[0].Error == nil {
		return 0, "", ""
	}
	var refusal *Error
	if !errors.As(results[0].Error, &refusal) {
		t.Fatalf("%s: not a catalog error: %v", stmt, results[0].Error)
	}
	return refusal.Code, refusal.SQLState, refusal.Message
}
