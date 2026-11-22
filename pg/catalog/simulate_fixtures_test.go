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

	// ALTER TABLE ... ADD UNIQUE on a partitioned table links child constraint via ConParentID,
	// and dropping it through the parent removes both parent and child constraints + indexes.
	if got, c := execCode(t, setup, "ALTER TABLE a ADD CONSTRAINT uq_a UNIQUE (a)"); got != "" {
		t.Fatalf("ALTER TABLE a ADD CONSTRAINT uq_a UNIQUE (a): got code %q, want ok", got)
	} else {
		p1 := c.GetRelation("public", "a_p1")
		found := false
		for _, con := range c.ConstraintsOf(p1.OID) {
			if con.Type == ConstraintUnique && con.ConParentID != 0 {
				found = true
			}
		}
		if !found {
			t.Fatalf("expected ALTER TABLE a ADD CONSTRAINT uq_a UNIQUE (a) to clone a linked unique constraint onto partition a_p1")
		}
		if results, err := c.Exec("ALTER TABLE a DROP CONSTRAINT uq_a", nil); err != nil {
			t.Fatalf("DROP CONSTRAINT uq_a: %v", err)
		} else if len(results) > 0 && results[0].Error != nil {
			t.Fatalf("DROP CONSTRAINT uq_a: %v", results[0].Error)
		}
		for _, con := range c.ConstraintsOf(p1.OID) {
			if con.Type == ConstraintUnique {
				t.Fatalf("expected DROP CONSTRAINT uq_a to remove child unique constraint, still have %q", con.Name)
			}
		}
		for _, idx := range c.IndexesOf(p1.OID) {
			if !idx.IsPrimary {
				t.Fatalf("expected DROP CONSTRAINT uq_a to remove child backing index, still have %q", idx.Name)
			}
		}
	}

	// Adding a PRIMARY KEY to a partitioned table when a partition already has
	// an incompatible primary key must fail with 42P16 and roll back cleanly.
	pkConflictSetup := `
CREATE TABLE p (a int, b int) PARTITION BY RANGE (a);
CREATE TABLE ch PARTITION OF p FOR VALUES FROM (0) TO (10);
ALTER TABLE ch ADD PRIMARY KEY (b);
`
	if got, c := execCode(t, pkConflictSetup, "ALTER TABLE p ADD PRIMARY KEY (a)"); got != CodeDuplicatePKey {
		t.Fatalf("incompatible child PK: got code %q, want %q", got, CodeDuplicatePKey)
	} else {
		p := c.GetRelation("public", "p")
		ch := c.GetRelation("public", "ch")
		if len(c.ConstraintsOf(p.OID)) != 0 {
			t.Fatalf("expected rollback on parent p, still have constraints %+v", c.ConstraintsOf(p.OID))
		}
		if len(c.ConstraintsOf(ch.OID)) != 1 {
			t.Fatalf("expected child ch to keep its single PK, got %d constraints", len(c.ConstraintsOf(ch.OID)))
		}
	}

	// Adding a PRIMARY KEY to a partitioned table when a partition already has a
	// compatible primary key attaches it instead of duplicating PKs, and dropping
	// it through the parent removes both.
	pkReuseSetup := `
CREATE TABLE p (a int, b int) PARTITION BY RANGE (a);
CREATE TABLE ch PARTITION OF p FOR VALUES FROM (0) TO (10);
ALTER TABLE ch ADD PRIMARY KEY (a);
`
	if got, c := execCode(t, pkReuseSetup, "ALTER TABLE p ADD PRIMARY KEY (a)"); got != "" {
		t.Fatalf("compatible child PK: got code %q, want ok", got)
	} else {
		ch := c.GetRelation("public", "ch")
		var pks []*Constraint
		for _, con := range c.ConstraintsOf(ch.OID) {
			if con.Type == ConstraintPK {
				pks = append(pks, con)
			}
		}
		if len(pks) != 1 || pks[0].ConParentID == 0 {
			t.Fatalf("expected single attached PK on ch, got %+v", pks)
		}
		if results, err := c.Exec("ALTER TABLE p DROP CONSTRAINT p_pkey", nil); err != nil {
			t.Fatalf("DROP CONSTRAINT p_pkey: %v", err)
		} else if len(results) > 0 && results[0].Error != nil {
			t.Fatalf("DROP CONSTRAINT p_pkey: %v", results[0].Error)
		}
		if len(c.ConstraintsOf(ch.OID)) != 0 {
			t.Fatalf("expected DROP CONSTRAINT p_pkey to drop child PK, still have %+v", c.ConstraintsOf(ch.OID))
		}
	}
}

func TestWave3SimulationFixes(t *testing.T) {
	fkSetup := `
CREATE TABLE p (b int PRIMARY KEY, d int);
CREATE TABLE a (b int PRIMARY KEY, a int, d int, CONSTRAINT a FOREIGN KEY (a) REFERENCES p (b), CONSTRAINT a_chk CHECK (d > 0));
CREATE TABLE ch (b int PRIMARY KEY, a_b int REFERENCES a (b));
`
	domainSetup := `CREATE DOMAIN a AS int DEFAULT 1 NOT NULL CONSTRAINT a CHECK (VALUE > 0);`
	enumSetup := `CREATE TYPE a AS ENUM ('x', 'y');`
	compSetup := `CREATE TYPE a AS (b int, a int, d int);`
	funcSetup := `
CREATE TABLE a (b int PRIMARY KEY, a int, d int);
CREATE FUNCTION a() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END; $$;
CREATE FUNCTION "authorization"() RETURNS int LANGUAGE sql AS 'SELECT 1';
`
	trigSetup := funcSetup + `CREATE TRIGGER a BEFORE INSERT ON a FOR EACH ROW EXECUTE FUNCTION a();`
	polSetup := `
CREATE TABLE a (b int PRIMARY KEY, a boolean, d int);
ALTER TABLE a ENABLE ROW LEVEL SECURITY;
CREATE POLICY a ON a USING (a);
`
	seqSetup := `CREATE SEQUENCE a;`
	viewSetup := `CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE VIEW a AS SELECT b, a, d FROM t;`
	mvSetup := `CREATE TABLE t (b int PRIMARY KEY, a int, d int); CREATE MATERIALIZED VIEW a AS SELECT b, a, d FROM t;`
	partSetup := `CREATE TABLE a (b int, a int, d int, PRIMARY KEY (b, a)) PARTITION BY RANGE (a); CREATE TABLE a_p1 PARTITION OF a FOR VALUES FROM (0) TO (100);`

	cases := []struct {
		name  string
		setup string
		stmt  string
		want  string
	}{
		// Domains & types
		{"domain add constraint not null", domainSetup, "ALTER DOMAIN a ADD CONSTRAINT a NOT NULL", ""},
		{"domain add not null", domainSetup, "ALTER DOMAIN a ADD NOT NULL", ""},
		{"domain set default colref", domainSetup, "ALTER DOMAIN a SET DEFAULT a", CodeFeatureNotSupported},
		{"domain add check bad col", domainSetup, "ALTER DOMAIN a ADD CHECK (a)", CodeUndefinedColumn},
		{"domain rename constraint missing domain", "", "ALTER DOMAIN a RENAME CONSTRAINT a TO b", CodeUndefinedObject},
		{"domain rename constraint non-domain", enumSetup, "ALTER DOMAIN a RENAME CONSTRAINT a TO b", CodeWrongObjectType},
		{"domain rename constraint duplicate", domainSetup, "ALTER DOMAIN a RENAME CONSTRAINT a TO a", CodeDuplicateObject},
		{"comment on domain constraint", domainSetup, "COMMENT ON CONSTRAINT a ON DOMAIN a IS 's'", ""},
		{"comment on domain constraint non-domain", enumSetup, "COMMENT ON CONSTRAINT a ON DOMAIN a IS 's'", CodeWrongObjectType},
		{"create domain dup before constraint", enumSetup, "CREATE DOMAIN a a PRIMARY KEY", CodeDuplicateObject},
		{"create domain missing base before constraint", "", "CREATE DOMAIN a a DEFERRABLE", CodeUndefinedObject},
		{"create enum dup label", "", "CREATE TYPE a AS ENUM ('s', 's')", CodeUniqueViolation},
		{"create range unknown attr", "", "CREATE TYPE a AS RANGE (a)", CodeSyntaxError},
		{"create type cross-db", "", "CREATE TYPE a.a.a", CodeFeatureNotSupported},
		{"create type full without shell", "", "CREATE TYPE a(a)", CodeDuplicateObject},
		{"create composite dup type", enumSetup, "CREATE TYPE a AS ()", CodeDuplicateObject},
		{"create composite dup col before type", "", "CREATE TYPE a AS (a a, a a)", CodeDuplicateColumn},
		{"alter type set missing", "", "ALTER TYPE a SET (a = 1)", CodeUndefinedObject},
		{"alter type set unknown attr", enumSetup, "ALTER TYPE a SET (a = 1)", CodeSyntaxError},
		{"alter type rename table rowtype", fkSetup, "ALTER TYPE a RENAME TO a", CodeWrongObjectType},
		{"alter type rename attr on partitioned", partSetup, "ALTER TYPE a RENAME ATTRIBUTE a TO a", CodeInvalidTableDefinition},

		// DROP TYPE / DOMAIN
		{"drop builtin type", "", "DROP TYPE INT", CodeDependentObjects},
		{"drop type bad typmod", "", "DROP TYPE BIT(a)", CodeInvalidTextRepresentation},
		{"drop type table rowtype", fkSetup, "DROP TYPE a", CodeDependentObjects},
		{"drop type if exists table rowtype", fkSetup, "DROP TYPE IF EXISTS a", CodeDependentObjects},
		{"drop domain on enum", enumSetup, "DROP DOMAIN a", CodeWrongObjectType},
		{"drop domain if exists on enum", enumSetup, "DROP DOMAIN IF EXISTS a", CodeWrongObjectType},

		// Table/View/Matview vs standalone type & FK/Identity
		{"create table if not exists vs enum", enumSetup, "CREATE TABLE IF NOT EXISTS a ()", CodeDuplicateObject},
		{"create matview if not exists vs domain", domainSetup, "CREATE MATERIALIZED VIEW IF NOT EXISTS a AS SELECT", CodeDuplicateObject},
		{"create view vs enum", enumSetup, "CREATE VIEW a AS SELECT", CodeDuplicateObject},
		{"create table of composite dup col", compSetup, "CREATE TABLE a OF a (a, a)", CodeDuplicateColumn},
		{"truncate fk referenced", fkSetup, "TRUNCATE a", CodeFeatureNotSupported},
		{"alter table set unlogged with fk", fkSetup, "ALTER TABLE a SET UNLOGGED", CodeInvalidTableDefinition},
		{"alter table add identity non-int", polSetup, "ALTER TABLE a ALTER a ADD GENERATED ALWAYS AS IDENTITY", CodeInvalidParameterValue},
		{"alter index rename composite", compSetup, "ALTER INDEX a RENAME TO a", CodeWrongObjectType},

		// Functions & Procedures
		{"create function dup window", "", "CREATE FUNCTION a() WINDOW WINDOW", CodeSyntaxError},
		{"alter function bad parallel", funcSetup, "ALTER FUNCTION AUTHORIZATION PARALLEL a", CodeSyntaxError},
		{"alter function dup strict", funcSetup, "ALTER FUNCTION AUTHORIZATION STRICT STRICT", CodeSyntaxError},
		{"alter function reset unknown guc", funcSetup, "ALTER FUNCTION AUTHORIZATION RESET a", CodeUndefinedObject},
		{"drop procedure if exists on func", funcSetup, "DROP PROCEDURE IF EXISTS AUTHORIZATION", ""},
		{"drop procedure on func no parens", funcSetup, "DROP PROCEDURE AUTHORIZATION", CodeUndefinedFunction},
		{"comment on procedure on func no parens", funcSetup, "COMMENT ON PROCEDURE AUTHORIZATION IS 's'", CodeUndefinedFunction},
		{"alter function owner missing role first", "", "ALTER FUNCTION AUTHORIZATION OWNER TO a", CodeUndefinedObject},
		{"alter domain owner missing role first", enumSetup, "ALTER DOMAIN a OWNER TO a", CodeUndefinedObject},

		// Triggers & Policies & table_openrv
		{"create trigger on matview", mvSetup, "CREATE TRIGGER a BEFORE INSERT ON a EXECUTE FUNCTION a()", CodeWrongObjectType},
		{"create trigger on composite", compSetup, "CREATE TRIGGER a BEFORE INSERT ON a EXECUTE FUNCTION a()", CodeWrongObjectType},
		{"create trigger referencing row on view", viewSetup, "CREATE TRIGGER a BEFORE INSERT ON a REFERENCING NEW ROW a EXECUTE FUNCTION a()", CodeFeatureNotSupported},
		{"alter trigger rename same name", trigSetup, "ALTER TRIGGER a ON a RENAME TO a", ""},
		{"alter trigger rename on matview", mvSetup, "ALTER TRIGGER a ON a RENAME TO a", CodeWrongObjectType},
		{"drop policy if exists on sequence", seqSetup, "DROP POLICY IF EXISTS a ON a", ""},
		{"drop policy if exists on composite", compSetup, "DROP POLICY IF EXISTS a ON a", CodeWrongObjectType},
		{"drop trigger on sequence", seqSetup, "DROP TRIGGER a ON a", CodeUndefinedObject},
		{"comment on constraint on sequence", seqSetup, "COMMENT ON CONSTRAINT a ON a IS 's'", CodeUndefinedObject},
		{"comment on policy on composite", compSetup, "COMMENT ON POLICY a ON a IS 's'", CodeWrongObjectType},
		{"create policy unknown role first", "", "CREATE POLICY a ON a TO a", CodeUndefinedObject},
		{"alter policy on view", viewSetup, "ALTER POLICY a ON a", CodeWrongObjectType},
		{"alter policy bad expr before lookup", fkSetup, "ALTER POLICY a ON a USING (a)", CodeDatatypeMismatch},
		{"alter policy if exists rename missing table", "", "ALTER POLICY IF EXISTS a ON a RENAME TO a", CodeUndefinedTable},
		{"alter policy if exists rename missing policy", fkSetup, "ALTER POLICY IF EXISTS a ON a RENAME TO a", CodeUndefinedObject},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, _ := execCode(t, tc.setup, tc.stmt); got != tc.want {
				t.Errorf("%s: got code %q, want %q", tc.stmt, got, tc.want)
			}
		})
	}

	// Verify SQL-parsed ALTER DOMAIN DROP NOT NULL and SET NOT NULL mutate NotNull accurately.
	if got, c := execCode(t, domainSetup, "ALTER DOMAIN a DROP NOT NULL"); got != "" {
		t.Fatalf("ALTER DOMAIN DROP NOT NULL failed: %s", got)
	} else if oid, _, err := c.ResolveType(TypeName{Name: "a", TypeMod: -1}); err != nil {
		t.Fatal(err)
	} else if c.DomainInfo(oid).NotNull {
		t.Fatalf("expected domain NotNull=false after DROP NOT NULL")
	}

	if got, c := execCode(t, "CREATE DOMAIN a AS int;", "ALTER DOMAIN a SET NOT NULL"); got != "" {
		t.Fatalf("ALTER DOMAIN SET NOT NULL failed: %s", got)
	} else if oid, _, err := c.ResolveType(TypeName{Name: "a", TypeMod: -1}); err != nil {
		t.Fatal(err)
	} else if !c.DomainInfo(oid).NotNull {
		t.Fatalf("expected domain NotNull=true after SET NOT NULL")
	}

	if got, c := execCode(t, domainSetup, "ALTER DOMAIN a RENAME CONSTRAINT a TO a_renamed"); got != "" {
		t.Fatalf("ALTER DOMAIN RENAME CONSTRAINT failed: %s", got)
	} else if oid, _, err := c.ResolveType(TypeName{Name: "a", TypeMod: -1}); err != nil {
		t.Fatal(err)
	} else if cons := c.DomainInfo(oid).Constraints; len(cons) != 1 || cons[0].Name != "a_renamed" {
		t.Fatalf("expected renamed constraint a_renamed, got %+v", cons)
	}
}
