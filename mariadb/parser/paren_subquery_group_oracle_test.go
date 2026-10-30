package parser

import (
	"strings"
	"testing"
)

// parenSubqueryViews create views over derived tables in joins and over
// subquery operands; MariaDB writes their definitions back with the shapes
// TestParenSubqueryGroup pins.
var parenSubqueryViews = []string{
	"CREATE VIEW v1 AS SELECT * FROM (SELECT k1 FROM j1) d JOIN j2 USING (k1)",
	"CREATE VIEW v2 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN j2 ON d.k1 = j2.k1 JOIN j3 ON j3.k1 = d.k1",
	"CREATE VIEW v3 AS SELECT d.k1 FROM (SELECT k1 FROM j1 UNION SELECT k1 FROM j2) d LEFT JOIN j3 ON j3.k1 = d.k1",
	"CREATE VIEW v4 AS SELECT d.x FROM (SELECT k1 FROM j1) d (x) JOIN j2 ON d.x = j2.k1",
	"CREATE VIEW v6 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d, j2",
	"CREATE VIEW v7 AS SELECT ((SELECT a FROM n1 LIMIT 1) + 0) AS x FROM n1",
	"CREATE VIEW v8 AS SELECT d.k1 FROM j2 JOIN ((SELECT k1 FROM j1) d JOIN j3 USING (k1)) ON j2.k1 = d.k1",
	"CREATE VIEW v9 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN (SELECT k1 FROM j2) e ON d.k1 = e.k1",
	"CREATE VIEW v10 AS SELECT a FROM n1 WHERE ((SELECT a FROM n1 LIMIT 1) = 1) AND a IN ((SELECT a FROM n1 LIMIT 1) + 0, 2)",
	"CREATE VIEW v11 AS SELECT ((SELECT a FROM n1 LIMIT 1) IS NULL) AS x, if(((SELECT count(0) FROM n1) = 0), 'NO', 'YES') AS y",
	"CREATE VIEW v12 AS SELECT d.k1 FROM ((SELECT k1 FROM j1) UNION (SELECT k1 FROM j2)) d JOIN j3 ON d.k1 = j3.k1",
}

// TestParenSubqueryGroupOracle asks MariaDB 11.8 about every statement of
// TestParenSubqueryGroup, and parses the definitions MariaDB stores for views
// over derived tables in joins and over subquery operands. MariaDB writes
// `from ((select ...) d join j2 on(...))` for a derived table joined to a
// table, and `((select ...) + 0)` for a subquery operand kept in
// parentheses; Bytebase parses a synced view definition to analyze queries
// over the view.
func TestParenSubqueryGroupOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, ddl := range []string{
		"CREATE TABLE j1 (k1 INT, a INT)",
		"CREATE TABLE j2 (k1 INT, b INT)",
		"CREATE TABLE j3 (k1 INT, c INT)",
		"CREATE TABLE n1 (a INT)",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	check := func(sql string, wantAccept bool) {
		t.Helper()
		v, code, msg := o.classify(sql)
		if v == vErrored {
			t.Fatalf("MariaDB errored on %q: %s", sql, msg)
		}
		if (v == vAccepted) != wantAccept {
			t.Errorf("MariaDB on %q: accepted %v (%d %s), want %v", sql, v == vAccepted, code, msg, wantAccept)
		}
		if _, err := Parse(sql); (err == nil) != wantAccept {
			t.Errorf("Parse(%q) = %v, want accepted %v", sql, err, wantAccept)
		}
	}
	for _, sql := range parenSubqueryAccepted {
		check(sql, true)
	}
	for _, sql := range parenSubqueryRejected {
		check(sql, false)
	}

	groups, operands := 0, 0
	for i, ddl := range parenSubqueryViews {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("%s: %v", ddl, err)
		}
		name := strings.Fields(ddl)[2]
		var definition, view, createView, charset, collation string
		if err := o.db.QueryRowContext(o.ctx,
			"SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?", name,
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
		if strings.Contains(definition, "from ((select ") || strings.Contains(definition, "from (((select ") {
			groups++
		}
		if strings.Contains(definition, "((select ") && strings.Contains(definition, ") + 0") {
			operands++
		}
	}
	// The views cover the reported shapes only while MariaDB still writes
	// them.
	if groups < 5 || operands < 1 {
		t.Errorf("MariaDB wrote %d derived-first join groups and %d subquery operands in parentheses, want at least 5 and 1", groups, operands)
	}
}
