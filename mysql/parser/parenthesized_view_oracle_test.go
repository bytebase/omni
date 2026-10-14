package parser

import (
	"strings"
	"testing"
)

// TestParenthesizedViewDefinitionOracle replays the BYT-10316 report against a
// live MySQL 8.0: the report's statements and the definitions the server
// stores for its views must all parse. MySQL stores a compound join ON
// condition with redundant parens — `on(((a = b) and (c = d)))` — and Bytebase
// parses a stored view definition to analyze any query that reads the view, so
// a definition omni rejects blocks every query on that view and on the views
// built over it (r5 over r4 in the report).
func TestParenthesizedViewDefinitionOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, ddl := range []string{
		"CREATE TABLE a (id INT, k VARCHAR(20))",
		"CREATE TABLE b (id INT, k VARCHAR(20))",
		"CREATE VIEW vb AS SELECT id, k FROM b",
		"CREATE VIEW r2 AS SELECT t1.id, t2.k FROM a t1 JOIN vb t2 ON t1.id = t2.id",
		"CREATE VIEW r4 AS SELECT t1.id, t2.k FROM a t1 JOIN vb t2 ON t1.id = t2.id AND t1.k = t2.k AND t2.k <> 'x' AND t2.k <> 'y'",
		"CREATE VIEW r5 AS SELECT id FROM r2 UNION ALL SELECT id FROM r4",
		"CREATE VIEW r8 AS SELECT t1.id FROM a t1 JOIN b t2 ON t1.id = t2.id AND t1.k = t2.k",
		"CREATE VIEW r9 AS SELECT t1.id FROM a t1 JOIN b t2 ON t1.id = t2.id AND t2.k <> 'x'",
		"CREATE VIEW r10 AS SELECT t1.id FROM a t1 JOIN b t2 ON t1.id = t2.id WHERE t2.k <> 'x'",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}

	for _, sql := range []string{
		"SELECT count(*) FROM a WHERE ((a.id = 1) AND (a.k = 'x'))",
		"SELECT count(*) FROM a WHERE (((a.id = 1) AND (a.k = 'x')))",
		"SELECT count(*) FROM a t1 JOIN b t2 ON (((t1.id = t2.id) AND (t1.k = t2.k)))",
	} {
		if _, err := o.db.ExecContext(o.ctx, sql); err != nil {
			t.Errorf("mysql rejects %q: %v", sql, err)
		}
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}

	definitions := map[string]string{}
	for _, view := range []string{"vb", "r2", "r4", "r5", "r8", "r9", "r10"} {
		var definition string
		if err := o.db.QueryRowContext(o.ctx,
			"SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = ?", view,
		).Scan(&definition); err != nil {
			t.Fatalf("read VIEW_DEFINITION of %s: %v", view, err)
		}
		var name, createView, charset, collation string
		if err := o.db.QueryRowContext(o.ctx, "SHOW CREATE VIEW "+view).Scan(&name, &createView, &charset, &collation); err != nil {
			t.Fatalf("SHOW CREATE VIEW %s: %v", view, err)
		}
		definitions[view] = definition
		for _, sql := range []string{definition, createView} {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		}
	}

	// The case is covered only while the server still stores the redundant
	// parens the report hit.
	if !strings.Contains(definitions["r4"], "on(((") {
		t.Errorf("r4 is stored without the redundant ON parens: %s", definitions["r4"])
	}
}
