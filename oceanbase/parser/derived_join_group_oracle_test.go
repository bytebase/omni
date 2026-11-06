package parser

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

// derivedJoinGroupViews create views over derived tables in joins and over
// parenthesized subquery operands; OceanBase writes their definitions back with
// the shapes TestDerivedJoinGroup pins.
var derivedJoinGroupViews = []string{
	"CREATE VIEW v1 AS SELECT * FROM (SELECT k1 FROM j1) d JOIN j2 USING (k1)",
	"CREATE VIEW v2 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN j2 ON d.k1 = j2.k1 JOIN j3 ON j3.k1 = d.k1",
	"CREATE VIEW v3 AS SELECT d.k1 FROM (SELECT k1 FROM j1 UNION SELECT k1 FROM j2) d LEFT JOIN j3 ON j3.k1 = d.k1",
	"CREATE VIEW v4 AS SELECT d.x FROM (SELECT k1 FROM j1) d (x) JOIN j2 ON d.x = j2.k1",
	"CREATE VIEW v5 AS SELECT j1.k1, d.y FROM j1, LATERAL (SELECT j1.k1 + 1 AS y) d",
	"CREATE VIEW v6 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d, j2",
	"CREATE VIEW v7 AS SELECT (SELECT a FROM n1 LIMIT 1) + 0 AS x",
	"CREATE VIEW v8 AS SELECT d.k1 FROM j2 JOIN ((SELECT k1 FROM j1) d JOIN j3 USING (k1)) ON j2.k1 = d.k1",
	"CREATE VIEW v9 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN (SELECT k1 FROM j2) e ON d.k1 = e.k1",
	"CREATE VIEW v10 AS SELECT a FROM n1 WHERE (SELECT a FROM n1 LIMIT 1) = 1 AND a IN ((SELECT a FROM n1 LIMIT 1) + 0, 2)",
	"CREATE VIEW v11 AS SELECT ((SELECT a FROM n1 LIMIT 1) IS NULL) AS x, if((SELECT count(0) FROM n1) = 0, 'NO', 'YES') AS y",
	"CREATE VIEW v12 AS SELECT d.k1 FROM ((SELECT k1 FROM j1) UNION (SELECT k1 FROM j2)) d JOIN j3 ON d.k1 = j3.k1",
}

// derivedFirstGroup finds a join group whose first operand is a derived table,
// at any depth of parentheses: OceanBase writes `from ((select ...) d join ...`
// and, for a group nested in another, `from (((select ...`.
var derivedFirstGroup = regexp.MustCompile(`from \(\(+select `)

// TestDerivedJoinGroupOracle parses the definitions OceanBase stores for views
// over derived tables in joins. OceanBase writes `SELECT * FROM (SELECT k1 FROM
// j1) d JOIN j2 USING (k1)` back as `from ((select ...) d join j2 on (...))`, a
// join group whose first operand is a derived table, and Bytebase parses a
// synced view definition to analyze queries over the view.
// TestDerivedJoinGroupOceanBase asks OceanBase about the statements themselves.
func TestDerivedJoinGroupOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, ddl := range []string{
		"CREATE TABLE j1 (k1 INT, a INT, KEY i (k1))",
		"CREATE TABLE j2 (k1 INT, b INT)",
		"CREATE TABLE j3 (k1 INT, c INT)",
		"CREATE TABLE n1 (a INT)",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	groups := 0
	for i, ddl := range derivedJoinGroupViews {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("%s: %v", ddl, err)
		}
		name := strings.Fields(ddl)[2]
		var definition, view, createView, charset, collation string
		if err := o.db.QueryRowContext(o.ctx,
			"SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", name,
		).Scan(&definition); err != nil {
			t.Fatalf("read VIEW_DEFINITION of %s: %v", name, err)
		}
		if err := o.db.QueryRowContext(o.ctx, "SHOW CREATE VIEW "+name).Scan(&view, &createView, &charset, &collation); err != nil {
			t.Fatalf("SHOW CREATE VIEW %s: %v", name, err)
		}
		for _, sql := range []string{definition, createView} {
			if _, err := Parse(sql); err != nil {
				t.Errorf("view %d: Parse(%q): %v", i+1, sql, err)
			}
		}
		if derivedFirstGroup.MatchString(definition) {
			groups++
		}
	}
	// The views cover the reported shape only while OceanBase still writes a
	// join group that begins with a derived table: six of them on 4.3.5.6
	// (2026-11-05).
	if groups < 5 {
		t.Errorf("OceanBase wrote %d view definitions with a derived-first join group, want at least 5", groups)
	}
}

// TestDerivedJoinGroupOceanBase asks OceanBase about every statement of
// TestDerivedJoinGroup. It accepts what MySQL accepts and refuses what MySQL
// refuses, apart from loneParenthesizedFactors, which it accepts, and
// derivedJoinGroupMySQLOnly, which it refuses (read on 2026-11-05).
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
func TestDerivedJoinGroupOceanBase(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obderived",
		"USE obderived",
		"DROP TABLE IF EXISTS j1, j2, j3",
		"CREATE TABLE j1 (k1 INT, a INT, KEY i (k1))",
		"CREATE TABLE j2 (k1 INT, b INT)",
		"CREATE TABLE j3 (k1 INT, c INT)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obderived")
	})
	check := func(sql string, want bool) {
		t.Helper()
		engineTakes, err := engineAccepts(ctx, db, sql)
		if err != nil {
			t.Fatalf("the oracle gave no verdict on %q: %v", sql, err)
		}
		if engineTakes != want {
			t.Errorf("OceanBase accepts=%v, want %v for %q", engineTakes, want, sql)
		}
		_, perr := Parse(sql)
		if (perr == nil) != engineTakes {
			t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v", perr == nil, engineTakes, sql, perr)
		}
	}
	for _, sql := range derivedJoinGroupAccepted {
		check(sql, true)
	}
	for _, sql := range loneParenthesizedFactors {
		check(sql, true)
	}
	for _, sql := range derivedJoinGroupRejected {
		check(sql, false)
	}
	for _, sql := range derivedJoinGroupMySQLOnly {
		check(sql, false)
	}
}
