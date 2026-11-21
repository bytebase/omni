package catalog

import (
	"errors"
	"testing"
)

func execCode(t *testing.T, setupSQL, stmt string) (string, *Catalog) {
	t.Helper()
	c, err := LoadSQL(setupSQL)
	if err != nil {
		t.Fatalf("LoadSQL(%q): %v", setupSQL, err)
	}
	results, err := c.Exec(stmt, nil)
	if err != nil {
		var ce *Error
		if errors.As(err, &ce) {
			return ce.Code, c
		}
		t.Fatalf("Exec(%q): unexpected non-catalog error %v", stmt, err)
	}
	for _, r := range results {
		if r.Error != nil {
			var ce *Error
			if errors.As(r.Error, &ce) {
				return ce.Code, c
			}
			t.Fatalf("Exec(%q): unexpected non-catalog error %v", stmt, r.Error)
		}
	}
	return "", c
}

func TestCommentOnTableAndViewRelKind(t *testing.T) {
	setup := `
CREATE TABLE t (b int PRIMARY KEY, a int, d int);
CREATE VIEW v AS SELECT b, a, d FROM t;
CREATE MATERIALIZED VIEW mv AS SELECT b, a, d FROM t;
`
	cases := []struct {
		stmt string
		want string
	}{
		{"COMMENT ON TABLE t IS 'ok'", ""},
		{"COMMENT ON VIEW v IS 'ok'", ""},
		{"COMMENT ON TABLE v IS 'nope'", CodeWrongObjectType},
		{"COMMENT ON TABLE mv IS 'nope'", CodeWrongObjectType},
		{"COMMENT ON VIEW t IS 'nope'", CodeWrongObjectType},
		{"COMMENT ON VIEW mv IS 'nope'", CodeWrongObjectType},
	}
	for _, tc := range cases {
		if got, _ := execCode(t, setup, tc.stmt); got != tc.want {
			t.Errorf("%s: got code %q, want %q", tc.stmt, got, tc.want)
		}
	}
}

func TestAlterTablePermissionsAcrossObjectKinds(t *testing.T) {
	viewSetup := `CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE VIEW a AS SELECT b, a, d FROM t;`
	matviewSetup := `CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE MATERIALIZED VIEW a AS SELECT b, a, d FROM t;`
	seqSetup := `CREATE TABLE t (b int PRIMARY KEY, c int, d int); CREATE SEQUENCE a;`

	cases := []struct {
		name  string
		setup string
		stmt  string
		want  string
	}{
		{"alter view drop col", viewSetup, "ALTER VIEW a DROP a", CodeWrongObjectType},
		{"alter view if exists drop col", viewSetup, "ALTER VIEW IF EXISTS a DROP a", CodeWrongObjectType},
		{"alter table view drop default", viewSetup, "ALTER TABLE a ALTER a DROP DEFAULT", ""},
		{"alter table view reset relopt", viewSetup, "ALTER TABLE a RESET (security_barrier)", ""},
		{"alter table view bad relopt", viewSetup, "ALTER TABLE a SET (unknown_opt = 1)", CodeInvalidParameterValue},
		{"alter table view owner missing role", viewSetup, "ALTER TABLE a OWNER TO no_such_role", CodeUndefinedObject},
		{"alter table view set schema missing", viewSetup, "ALTER TABLE a SET SCHEMA no_such_schema", CodeUndefinedSchema},
		{"alter table view add identity nullable", viewSetup, "ALTER TABLE a ALTER a ADD GENERATED ALWAYS AS IDENTITY", CodeObjectNotInPrerequisiteState},
		{"alter table view drop identity if exists", viewSetup, "ALTER TABLE a ALTER a DROP IDENTITY IF EXISTS", ""},
		{"create expr index on view analyzes expr first", viewSetup, "CREATE INDEX ON a(nonexistent_fn(a))", CodeUndefinedFunction},

		{"alter matview drop col", matviewSetup, "ALTER MATERIALIZED VIEW a DROP a", CodeWrongObjectType},
		{"alter table matview set statistics", matviewSetup, "ALTER TABLE a ALTER a SET STATISTICS 10", ""},
		{"alter table matview set statistics by num", matviewSetup, "ALTER TABLE a ALTER 1 SET STATISTICS 10", CodeFeatureNotSupported},
		{"alter table matview set without cluster", matviewSetup, "ALTER TABLE a SET WITHOUT CLUSTER", ""},
		{"alter table matview replica identity full", matviewSetup, "ALTER TABLE a REPLICA IDENTITY FULL", ""},
		{"alter table matview set access method default", matviewSetup, "ALTER TABLE a SET ACCESS METHOD DEFAULT", ""},
		{"alter table matview set tablespace missing", matviewSetup, "ALTER TABLE a SET TABLESPACE no_such_ts", CodeUndefinedObject},
		{"alter table matview set schema missing", matviewSetup, "ALTER TABLE a SET SCHEMA no_such_schema", CodeUndefinedSchema},

		{"alter seq drop col", seqSetup, "ALTER SEQUENCE a DROP a", CodeWrongObjectType},
		{"alter table seq set unlogged", seqSetup, "ALTER TABLE a SET UNLOGGED", ""},
		{"alter table seq set logged", seqSetup, "ALTER TABLE a SET LOGGED", ""},
		{"alter table seq set tablespace", seqSetup, "ALTER TABLE a SET TABLESPACE no_such_ts", CodeWrongObjectType},
		{"alter table seq set schema missing", seqSetup, "ALTER TABLE a SET SCHEMA no_such_schema", CodeUndefinedSchema},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, _ := execCode(t, tc.setup, tc.stmt); got != tc.want {
				t.Errorf("%s: got code %q, want %q", tc.stmt, got, tc.want)
			}
		})
	}
}

func TestAlterSequenceOptionsValidation(t *testing.T) {
	setup := `CREATE SEQUENCE a;`
	cases := []struct {
		stmt string
		want string
	}{
		{"ALTER SEQUENCE a AS no_such_type", CodeUndefinedObject},
		{"ALTER SEQUENCE a AS text", CodeFeatureNotSupported},
		{"ALTER SEQUENCE a AS smallint", ""},
		{"ALTER SEQUENCE a CACHE 1.5", CodeInvalidTextRepresentation},
		{"ALTER SEQUENCE a INCREMENT 1.5", CodeInvalidTextRepresentation},
		{"ALTER SEQUENCE a MINVALUE 1.5", CodeInvalidTextRepresentation},
		{"ALTER SEQUENCE a MAXVALUE 1.5", CodeInvalidTextRepresentation},
		{"ALTER SEQUENCE a START 1.5", CodeInvalidTextRepresentation},
		{"ALTER SEQUENCE a RESTART 1.5", CodeInvalidTextRepresentation},
		{"ALTER SEQUENCE a OWNED BY a", CodeSyntaxError},
		{"ALTER SEQUENCE a SEQUENCE NAME a", CodeSyntaxError},
		{"ALTER SEQUENCE a LOGGED", CodeInternalError},
		{"ALTER SEQUENCE a UNLOGGED", CodeInternalError},
		{"ALTER SEQUENCE a SET UNLOGGED", ""},
		{"ALTER SEQUENCE a SET LOGGED", ""},
	}
	for _, tc := range cases {
		if got, _ := execCode(t, setup, tc.stmt); got != tc.want {
			t.Errorf("%s: got code %q, want %q", tc.stmt, got, tc.want)
		}
	}
}

func TestPartitionedTableBehavior(t *testing.T) {
	setup := `
CREATE TABLE a (b int, a int, d int, PRIMARY KEY (b, a)) PARTITION BY RANGE (a);
CREATE TABLE a_p1 PARTITION OF a FOR VALUES FROM (0) TO (100);
`
	cases := []struct {
		stmt string
		want string
	}{
		{"ALTER TABLE a DETACH PARTITION a", CodeUndefinedTable},
		{"ALTER TABLE a ATTACH PARTITION a DEFAULT", CodeDuplicateTable},
		{"ALTER TABLE a ATTACH PARTITION a FOR VALUES IN (1)", CodeInvalidTableDefinition},
		{"ALTER TABLE a ATTACH PARTITION a FOR VALUES FROM (a) TO (a)", CodeFeatureNotSupported},
		{"ALTER TABLE a ALTER a DROP NOT NULL", CodeInvalidTableDefinition},
		{"ALTER TABLE a SET WITHOUT CLUSTER", CodeFeatureNotSupported},
		{"ALTER TABLE a INHERIT a", CodeWrongObjectType},
		{"ALTER TABLE a ALTER a TYPE no_such_type", CodeInvalidTableDefinition},
		{"ALTER TABLE a SET (badns.foo = 1)", CodeInvalidParameterValue},
		{"ALTER TABLE a SET (fillfactor = 50)", CodeWrongObjectType},
	}
	for _, tc := range cases {
		if got, _ := execCode(t, setup, tc.stmt); got != tc.want {
			t.Errorf("%s: got code %q, want %q", tc.stmt, got, tc.want)
		}
	}

	// DROP TABLE on a partitioned table auto-drops child partitions without CASCADE.
	if got, c := execCode(t, setup, "DROP TABLE a"); got != "" {
		t.Fatalf("DROP TABLE a: got code %q, want ok", got)
	} else if c.GetRelation("public", "a_p1") != nil {
		t.Fatalf("DROP TABLE a left partition a_p1 behind")
	}

	// ALTER TABLE ... SET UNLOGGED on a partitioned table is a no-op on persistence.
	if got, c := execCode(t, setup, "ALTER TABLE a SET UNLOGGED"); got != "" {
		t.Fatalf("ALTER TABLE a SET UNLOGGED: got code %q, want ok", got)
	} else if rel := c.GetRelation("public", "a"); rel == nil || rel.Persistence != 'p' {
		t.Fatalf("expected partitioned table persistence 'p', got %+v", rel)
	}

	// CREATE INDEX on a partitioned table clones to existing partitions (attached).
	if got, c := execCode(t, setup, "CREATE INDEX ON a (a)"); got != "" {
		t.Fatalf("CREATE INDEX ON a (a): got code %q, want ok", got)
	} else {
		p1 := c.GetRelation("public", "a_p1")
		found := false
		for _, idx := range c.IndexesOf(p1.OID) {
			if idx.ParentIndexOID != 0 && idx.ConstraintOID == 0 {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected CREATE INDEX ON a (a) to clone an attached index onto partition a_p1")
		}
	}

	// ALTER TABLE ... ADD UNIQUE on a partitioned table links child constraint via ConParentID.
	if got, c := execCode(t, setup, "ALTER TABLE a ADD UNIQUE (a)"); got != "" {
		t.Fatalf("ALTER TABLE a ADD UNIQUE (a): got code %q, want ok", got)
	} else {
		p1 := c.GetRelation("public", "a_p1")
		found := false
		for _, con := range c.ConstraintsOf(p1.OID) {
			if con.Type == ConstraintUnique && con.ConParentID != 0 {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected ALTER TABLE a ADD UNIQUE (a) to clone a linked unique constraint onto partition a_p1")
		}
	}
}
