package parser

import "testing"

func TestRedshiftTemplateAndShowDocsForms(t *testing.T) {
	// CREATE / ALTER / DROP TEMPLATE
	for _, tc := range []struct {
		sql     string
		command string
	}{
		{"CREATE TEMPLATE t FOR COPY", "create"},
		{"CREATE TEMPLATE İn FOR COPY", "create"},
		{"CREATE OR REPLACE TEMPLATE db.s.t FOR COPY AS FORMAT AS parquet emptyasnull, maxerror 10", "create or replace"},
		{"CREATE TEMPLATE s.t FOR COPY FORMAT csv delimiter ',', ignoreheader 1", "create"},
		{"CREATE TEMPLATE t FOR COPY parquet emptyasnull", "create"},
		{"ALTER TEMPLATE db.s.t RENAME TO t2", "alter"},
		{"ALTER TEMPLATE s.t OWNER TO alice", "alter"},
		{"ALTER TEMPLATE t ADD delimiter AS ','", "alter"},
		{"ALTER TEMPLATE t ADD emptyasnull", "alter"},
		{"ALTER TEMPLATE t DROP emptyasnull", "alter"},
		{"ALTER TEMPLATE t SET search_path TO public, timezone TO DEFAULT", "alter"},
		{"DROP TEMPLATE db.s.t", "drop"},
	} {
		stmt := firstRedshiftObjectStmt(t, tc.sql)
		if stmt.Command != tc.command || stmt.ObjectType != "template" {
			t.Errorf("%q: got (%q, %q), want (%q, \"template\")", tc.sql, stmt.Command, stmt.ObjectType, tc.command)
		}
	}

	// Documented SHOW commands
	for _, tc := range []struct {
		sql     string
		command string
	}{
		{"SHOW USER LOCKOUT", "user lockout"},
		{"SHOW COLUMN GRANTS ON TABLE s.t FOR ROLE r LIMIT 5", "column grants"},
		{"SHOW COLUMN GRANTS ON TABLE db.s.t FOR PUBLIC", "column grants"},
		{"SHOW COLUMN GRANTS ON TABLE s.t FOR ROLE", "column grants"},
		{"SHOW CONSTRAINTS PRIMARY KEYS FROM TABLE s.t LIMIT 10", "constraints"},
		{"SHOW CONSTRAINTS FOREIGN KEYS EXPORTED FROM TABLE db.s.t", "constraints"},
		{"SHOW FUNCTIONS FROM SCHEMA db.s LIKE 'f%' LIMIT 10", "functions"},
		{"SHOW PARAMETERS OF FUNCTION s.f(INT, VARCHAR(20)) LIKE 'p%'", "parameters"},
		{"SHOW PARAMETERS OF PROCEDURE db.s.p(SUPER)", "parameters"},
		{"SHOW RLS POLICIES ON db.s.t FOR ROLE r LIMIT 5", "rls policies"},
		{"SHOW MASKING POLICIES FROM DATABASE dev LIMIT 5", "masking policies"},
		{"SHOW PROCEDURES FROM SCHEMA s LIKE 'p%' LIMIT 5", "procedures"},
		{"SHOW TEMPLATE db.s.t", "template"},
		{"SHOW TEMPLATE İn", "template"},
		{"SHOW TEMPLATES FROM SCHEMA s LIKE 't%' LIMIT 5", "templates"},
	} {
		stmt := firstRedshiftShowStmt(t, tc.sql)
		if stmt.Command != tc.command {
			t.Errorf("%q: Command = %q, want %q", tc.sql, stmt.Command, tc.command)
		}
	}

	// Negative checks: undocumented shapes, wrong part counts, reserved or quoted keywords
	for _, bad := range []string{
		"CREATE TEMPLATE t",
		`CREATE "TEMPLATE" t FOR COPY`,
		`CREATE TEMPLATE t FOR "COPY"`,
		"CREATE TEMPLATE a.b.c.d FOR COPY",
		"CREATE TEMPLATE select FOR COPY",
		"CREATE TEMPLATE t FOR COPY ASC",
		"CREATE TEMPLATE t FOR COPY a,",
		"CREATE TEMPLATE t FOR COPY a b c d",
		"CREATE TEMPLATE t FOR COPY FORMAT parquet, emptyasnull",
		"ALTER TEMPLATE t",
		"ALTER TEMPLATE a.b.c.d RENAME TO u",
		"ALTER TEMPLATE t RENAME TO s.u",
		"ALTER TEMPLATE t OWNER TO s.u",
		"ALTER TEMPLATE t ADD a ASC",
		"ALTER TEMPLATE t ADD a b c",
		"ALTER TEMPLATE t DROP a b",
		"ALTER TEMPLATE t SET a TO b,",
		"DROP TEMPLATE a.b.c.d",
		"DROP TEMPLATE IF EXISTS t",
		"DROP TEMPLATE t CASCADE",
		"SHOW USER LOCKOUT extra",
		`SHOW USER "LOCKOUT"`,
		"SHOW COLUMN GRANTS ON TABLE t",
		"SHOW COLUMN GRANTS ON TABLE a.b.c.d",
		"SHOW COLUMN GRANTS ON TABLE s.t LIMIT 'x'",
		"SHOW CONSTRAINTS FROM TABLE s.t",
		"SHOW CONSTRAINTS PRIMARY KEYS FROM TABLE t",
		`SHOW CONSTRAINTS PRIMARY "KEYS" FROM TABLE s.t`,
		"SHOW CONSTRAINTS PRIMARY KEYS FROM TABLE s.t",
		"SHOW FUNCTIONS FROM SCHEMA a.b.c",
		"SHOW PARAMETERS OF FUNCTION f(INT)",
		"SHOW PARAMETERS OF FUNCTION s.f()",
		"SHOW PARAMETERS OF FUNCTION s.f(SETOF INT)",
		`SHOW "RLS" POLICIES`,
		"SHOW RLS POLICIES ON t",
		"SHOW RLS POLICIES FROM DATABASE s.d",
		"SHOW TEMPLATE a.b.c.d",
		"SHOW TEMPLATE t extra",
		"SHOW TEMPLATES FROM SCHEMA a.b.c",
	} {
		if _, err := Parse(bad); err == nil {
			t.Errorf("Parse(%q) succeeded, want syntax error", bad)
		}
	}
}
