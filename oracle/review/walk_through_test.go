package review

import (
	"context"
	"math/rand"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/review/walk/walktest"
)

// walkSchema is the synced SHOP schema the unit tests walk on, as
// Bytebase's sync lists it: tables with keys, constraints, and an index, a
// view, a materialized view, a sequence, and routines. With triggers, a
// table TRG has one.
func walkSchema(triggers bool) *metadata.DatabaseSchemaMetadata {
	cols := func(names ...string) []*metadata.ColumnMetadata {
		var out []*metadata.ColumnMetadata
		for _, n := range names {
			out = append(out, &metadata.ColumnMetadata{Name: n})
		}
		return out
	}
	s := &metadata.SchemaMetadata{
		Tables: []*metadata.TableMetadata{
			{Name: "P", Columns: cols("ID", "NOTE"), Indexes: []*metadata.IndexMetadata{{Name: "PK_P", Primary: true, Unique: true, IsConstraint: true}}},
			{
				Name:             "T",
				Columns:          cols("ID", "CODE", "P_ID", "N"),
				Indexes:          []*metadata.IndexMetadata{{Name: "PK_T", Primary: true, Unique: true, IsConstraint: true}, {Name: "UK_CODE", Unique: true, IsConstraint: true}, {Name: "IX_N"}},
				ForeignKeys:      []*metadata.ForeignKeyMetadata{{Name: "FK_P", ReferencedTable: "P"}},
				CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "CK_N", Expression: "n > 0"}, {Name: "NN_CODE", Expression: `"CODE" IS NOT NULL`}},
			},
			{Name: "NOKEY", Columns: cols("A")},
			// A table with a valid-time period: its hidden columns are not
			// listed, and its check constraint names them.
			{Name: "PERIODIC", Columns: cols("A"), CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "VT8cf288", Expression: `"VT_START" < "VT_END"`}}},
			// An object table, whose hidden columns are not listed.
			{Name: "OT", Columns: cols("X", "Y"), SkipDump: true},
			// A table with a nested table column, and its storage table, of
			// no listed column.
			{Name: "PARENT", Columns: cols("A", "C")},
			{Name: "CHILD"},
			// A table with a domain index.
			{Name: "DOC", Columns: cols("ID", "BODY"), Indexes: []*metadata.IndexMetadata{{Name: "DOC_CTX", Type: "DOMAIN"}}},
			// A materialized view log's table, of P.
			{Name: "MLOG$_P", Columns: cols("ID", "SNAPTIME$$"), Indexes: []*metadata.IndexMetadata{{Name: "I_MLOG$_P"}}},
		},
		Views:             []*metadata.ViewMetadata{{Name: "PV"}},
		MaterializedViews: []*metadata.MaterializedViewMetadata{{Name: "MV", Indexes: []*metadata.IndexMetadata{{Name: "MV_IX"}}}},
		Sequences:         []*metadata.SequenceMetadata{{Name: "SQ"}},
		Functions:         []*metadata.FunctionMetadata{{Name: "F"}},
		Procedures:        []*metadata.ProcedureMetadata{{Name: "PR"}},
		Packages:          []*metadata.PackageMetadata{{Name: "PKG"}},
		Synonyms:          []*metadata.SynonymMetadata{{Name: "SYN", TargetSchema: "SHOP", TargetName: "T"}},
	}
	if triggers {
		s.Tables = append(s.Tables, &metadata.TableMetadata{Name: "TRG", Columns: cols("A"), Triggers: []*metadata.TriggerMetadata{{Name: "TR"}}})
	}
	return &metadata.DatabaseSchemaMetadata{
		Name:    "SHOP",
		Schemas: []*metadata.SchemaMetadata{s},
		PublicSynonyms: []*metadata.SynonymMetadata{
			{Name: "DUAL", TargetSchema: "SYS", TargetName: "DUAL"},
			{Name: "PUBT", TargetSchema: "SYSTEM", TargetName: "PUBT"},
		},
	}
}

func walkFindings(t *testing.T, sql string, target review.Target) []review.Finding {
	t.Helper()
	result, err := Review(context.Background(), sql, review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{target})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range result.Findings {
		if f.Rule == review.Syntax {
			t.Fatalf("syntax: %s\n%s", f.Message, sql)
		}
		if err := review.CheckCodes(f); err != nil {
			t.Errorf("%v\n%s", err, sql)
		}
	}
	return result.Findings
}

// oneFinding asserts the change's one WALK_THROUGH finding.
func oneFinding(t *testing.T, sql string, statement int, message string, target review.Target) {
	t.Helper()
	got := walkFindings(t, sql, target)
	if len(got) != 1 || got[0].Statement != statement || got[0].Message != message {
		t.Errorf("%s\nfindings = %+v, want statement %d: %q", sql, got, statement, message)
	}
}

func noFinding(t *testing.T, sql string, target review.Target) {
	t.Helper()
	if got := walkFindings(t, sql, target); len(got) != 0 {
		t.Errorf("%s\nfindings = %+v, want none", sql, got)
	}
}

func TestWalkThroughOracle(t *testing.T) {
	shop := review.Target{Schema: walkSchema(false)}
	// Names a statement creates.
	oneFinding(t, "CREATE TABLE x (a NUMBER);\nCREATE TABLE t (a NUMBER);", 1, "Table `T` already exists.", shop)
	oneFinding(t, "CREATE SEQUENCE pv;", 0, "The name `PV` is already taken.", shop)
	oneFinding(t, "CREATE OR REPLACE VIEW t AS SELECT 1 a FROM dual;", 0, "The name `T` is already taken.", shop)
	noFinding(t, "CREATE OR REPLACE VIEW pv AS SELECT 1 a FROM dual;\nCREATE OR REPLACE FUNCTION f RETURN NUMBER AS BEGIN RETURN 1; END;", shop)
	noFinding(t, "CREATE FUNCTION IF NOT EXISTS f RETURN NUMBER IS BEGIN RETURN 1; END;\n/\nCREATE FUNCTION IF NOT EXISTS t RETURN NUMBER IS BEGIN RETURN 1; END;", shop)
	oneFinding(t, "CREATE FUNCTION IF NOT EXISTS g RETURN NUMBER IS BEGIN RETURN 1; END;\n/\nCREATE TABLE g (a NUMBER);", 1, "The name `G` is already taken.", shop)
	oneFinding(t, "CREATE OR REPLACE PROCEDURE f AS BEGIN NULL; END;", 0, "The name `F` is already taken.", shop)
	noFinding(t, "CREATE TABLE IF NOT EXISTS t (a NUMBER, b NUMBER CONSTRAINT ck_n CHECK (b > 0));", shop)
	oneFinding(t, "CREATE TABLE IF NOT EXISTS t (a NUMBER, a NUMBER);", 0, "This statement declares column `A` twice.", shop)
	// Tables, which ALTER, DROP, and TRUNCATE take by no synonym.
	oneFinding(t, "ALTER TABLE nosuch ADD (x NUMBER);", 0, "Table `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE pv ADD (x NUMBER);", 0, "`PV` is a view, not a table.", shop)
	noFinding(t, "ALTER TABLE mv LOGGING;", shop)
	oneFinding(t, "DROP TABLE pv;", 0, "`PV` is a view, not a table.", shop)
	oneFinding(t, "DROP TABLE mv;", 0, "`MV` is a materialized view, not a table.", shop)
	oneFinding(t, "DROP TABLE IF EXISTS mv;", 0, "`MV` is a materialized view, not a table.", shop)
	noFinding(t, "DROP TABLE IF EXISTS pv;\nDROP TABLE IF EXISTS nosuch;", shop)
	oneFinding(t, "TRUNCATE TABLE pv;", 0, "`PV` is a view, not a table.", shop)
	oneFinding(t, "TRUNCATE TABLE sq;", 0, "`SQ` is a sequence, not a table.", shop)
	oneFinding(t, "TRUNCATE TABLE pkg;", 0, "`PKG` is a package, not a table.", shop)
	oneFinding(t, "TRUNCATE TABLE nosuch;", 0, "Table `NOSUCH` doesn't exist.", shop)
	noFinding(t, "TRUNCATE TABLE mv;", shop)
	// Columns, keys, and constraints.
	oneFinding(t, "ALTER TABLE t ADD (code NUMBER);", 0, "Column `CODE` already exists in `T`.", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN zz;", 0, "Column `ZZ` doesn't exist in `T`.", shop)
	oneFinding(t, "ALTER TABLE t DROP (n, zz);", 0, "Column `ZZ` doesn't exist in `T`.", shop)
	oneFinding(t, "ALTER TABLE t RENAME COLUMN n TO code;", 0, "Column `CODE` already exists in `T`.", shop)
	oneFinding(t, "ALTER TABLE t RENAME COLUMN n TO n;", 0, "Column `N` already exists in `T`.", shop)
	oneFinding(t, "ALTER TABLE t MODIFY (zz NUMBER);", 0, "Column `ZZ` doesn't exist in `T`.", shop)
	oneFinding(t, "ALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);", 0, "Constraint `CK_N` already exists.", shop)
	oneFinding(t, "ALTER TABLE nokey ADD (b NUMBER CONSTRAINT nn_code NOT NULL);", 0, "Constraint `NN_CODE` already exists.", shop)
	oneFinding(t, "ALTER TABLE t ADD PRIMARY KEY (n);", 0, "`T` already has a primary key.", shop)
	oneFinding(t, "ALTER TABLE nokey DROP PRIMARY KEY;", 0, "`NOKEY` has no primary key.", shop)
	oneFinding(t, "ALTER TABLE p DROP CONSTRAINT ck_n;", 0, "Constraint `CK_N` doesn't exist in `P`.", shop)
	oneFinding(t, "ALTER TABLE t DISABLE CONSTRAINT ck_p;\nALTER TABLE t DROP CONSTRAINT ck_n;\nALTER TABLE t ENABLE CONSTRAINT ck_n;", 2, "Constraint `CK_N` doesn't exist in `T`.", shop)
	oneFinding(t, "ALTER TABLE t RENAME CONSTRAINT ck_n TO fk_p;", 0, "Constraint `FK_P` already exists.", shop)
	noFinding(t, "ALTER TABLE t DROP CONSTRAINT nosuch;\nALTER TABLE t ADD CONSTRAINT nosuch CHECK (n > 0);", shop)
	// A trailing ENABLE or DISABLE sees the constraints of the statement's
	// ADD and MODIFY clauses, and MODIFY CONSTRAINT does not.
	noFinding(t, "CREATE TABLE x (a NUMBER);\nALTER TABLE x ADD (z NUMBER CONSTRAINT c5 CHECK (z > 0)) DISABLE CONSTRAINT c5;\nALTER TABLE x ADD (CONSTRAINT c9 CHECK (a > 0)) ENABLE CONSTRAINT c9;\nALTER TABLE x MODIFY (a CONSTRAINT nn_a NOT NULL) DISABLE CONSTRAINT nn_a;", shop)
	oneFinding(t, "CREATE TABLE x (a NUMBER);\nALTER TABLE x ADD (z NUMBER CONSTRAINT c8 CHECK (z > 0)) MODIFY CONSTRAINT c8 DISABLE;", 1, "Constraint `C8` doesn't exist in `X`.", shop)
	oneFinding(t, "CREATE TABLE x (a NUMBER);\nALTER TABLE x ADD (z NUMBER CONSTRAINT c5 CHECK (z > 0)) DISABLE CONSTRAINT c6;", 1, "Constraint `C6` doesn't exist in `X`.", shop)
	// The text alone.
	oneFinding(t, "CREATE TABLE x (a NUMBER, a NUMBER);", 0, "This statement declares column `A` twice.", shop)
	oneFinding(t, "CREATE TABLE x (a NUMBER CONSTRAINT c CHECK (a > 0), b NUMBER CONSTRAINT c NOT NULL);", 0, "This statement declares constraint `C` twice.", shop)
	oneFinding(t, "CREATE TABLE x (a NUMBER PRIMARY KEY, b NUMBER, PRIMARY KEY (b));", 0, "`X` already has a primary key.", shop)
	oneFinding(t, "ALTER TABLE t ADD (x NUMBER) MODIFY (x NUMBER);", 0, "This statement names column `X` twice.", shop)
	oneFinding(t, "ALTER TABLE t DROP CONSTRAINT ck_n DROP CONSTRAINT ck_n;", 0, "This statement drops constraint `CK_N` twice.", shop)
	oneFinding(t, "CREATE TABLE "+strings.Repeat("a", 129)+" (a NUMBER);", 0, "The name `"+strings.Repeat("A", 32)+"…` is longer than Oracle's 128-byte limit.", shop)
	noFinding(t, "CREATE TABLE "+strings.Repeat("a", 128)+" (a NUMBER);", shop)
	// Indexes and triggers: the sync lists them in part.
	oneFinding(t, "CREATE INDEX ix_n ON nokey (a);", 0, "The name `IX_N` is already taken.", shop)
	oneFinding(t, "CREATE INDEX mv_ix ON nokey (a);", 0, "The name `MV_IX` is already taken.", shop)
	noFinding(t, "CREATE INDEX pk_t ON nokey (a);\nDROP INDEX nosuch;\nCREATE INDEX x ON nosuch (a);", shop)
	oneFinding(t, "CREATE INDEX x ON pv (id);", 0, "`PV` is a view, not a table.", shop)
	oneFinding(t, "CREATE INDEX i ON nokey (a);\nDROP INDEX i;\nDROP INDEX i;", 2, "Index `I` doesn't exist.", shop)
	oneFinding(t, "CREATE INDEX i ON nokey (a);\nALTER INDEX i RENAME TO ix_n;", 1, "The name `IX_N` is already taken.", shop)
	oneFinding(t, "CREATE TRIGGER tr2 BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;\n/\nCREATE TRIGGER tr2 BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;", 1, "Trigger `TR2` already exists.", shop)
	oneFinding(t, "CREATE TRIGGER tr2 BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;\n/\nCREATE OR REPLACE TRIGGER tr2 BEFORE INSERT ON p FOR EACH ROW BEGIN NULL; END;", 1, "Trigger `TR2` already exists on another table.", shop)
	// Views, sequences, synonyms, and routines.
	oneFinding(t, "DROP VIEW t;", 0, "`T` is a table, not a view.", shop)
	oneFinding(t, "DROP SEQUENCE t;", 0, "`T` is a table, not a sequence.", shop)
	oneFinding(t, "ALTER SEQUENCE nosuch INCREMENT BY 2;", 0, "Sequence `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "DROP SYNONYM t;", 0, "`T` isn't a private synonym.", shop)
	noFinding(t, "DROP SYNONYM nosuch;\nDROP PUBLIC SYNONYM t;\nDROP PROCEDURE nosuch2;", shop)
	oneFinding(t, "DROP SYNONYM nosuch;\nDROP PROCEDURE nosuch;", 1, "Procedure `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "DROP FUNCTION pr;", 0, "`PR` isn't a function.", shop)
	oneFinding(t, "CREATE SYNONYM s FOR t;\nDROP SYNONYM s;\nDROP SYNONYM s;", 2, "Private synonym `S` doesn't exist.", shop)
	// RENAME.
	oneFinding(t, "RENAME pv TO t;", 0, "The name `T` is already taken.", shop)
	oneFinding(t, "RENAME f TO g;", 0, "`F` is a function, not a table, view, sequence, or private synonym.", shop)
	oneFinding(t, "RENAME mv TO x;", 0, "`MV` is a materialized view, not a table, view, sequence, or private synonym.", shop)
	oneFinding(t, "RENAME t TO t;", 0, "The name `T` is already taken.", shop)
	// Grammar-derived simulation fixes (TestWalkSimulation):
	// 1) ALTER TABLE ... MODIFY VARRAY ... (SHRINK SPACE [CASCADE]) is accepted
	//    by Oracle 23 even when the table is missing or another object kind.
	noFinding(t, "ALTER TABLE nosuch MODIFY VARRAY a (SHRINK SPACE);", shop)
	noFinding(t, "ALTER TABLE nosuch MODIFY VARRAY a (SHRINK SPACE CASCADE);", shop)
	noFinding(t, "ALTER TABLE pv MODIFY VARRAY a (SHRINK SPACE);", shop)
	noFinding(t, "ALTER TABLE sq MODIFY VARRAY a (SHRINK SPACE CASCADE);", shop)
	oneFinding(t, "ALTER TABLE nosuch MODIFY VARRAY a (CACHE);", 0, "Table `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE nosuch MODIFY VARRAY shrink (CACHE);", 0, "Table `NOSUCH` doesn't exist.", shop)
	// 2) ALTER TABLE ... MOVE and ALTER TABLE ... MODIFY PARTITION BY ...
	//    follow private and public synonyms in Oracle 23, while still failing
	//    on missing names or non-table/non-synonym objects.
	noFinding(t, "ALTER TABLE syn MOVE;", shop)
	noFinding(t, "ALTER TABLE syn MOVE ONLINE;", shop)
	noFinding(t, "ALTER TABLE syn MODIFY PARTITION BY HASH (id) PARTITIONS 2;", shop)
	noFinding(t, "ALTER TABLE pubt MOVE;", shop)
	oneFinding(t, "ALTER TABLE nosuch MOVE;", 0, "Table `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE pv MOVE;", 0, "`PV` is a view, not a table.", shop)
	oneFinding(t, "ALTER TABLE sq MODIFY PARTITION BY HASH (a) PARTITIONS 2;", 0, "`SQ` is a sequence, not a table.", shop)
	oneFinding(t, "ALTER TABLE f MOVE;", 0, "`F` is a function, not a table.", shop)
	// 3) CREATE SYNONYM IF NOT EXISTS is a no-op when the name is taken.
	noFinding(t, "CREATE SYNONYM IF NOT EXISTS t FOR p;", shop)
	noFinding(t, "CREATE SYNONYM IF NOT EXISTS syn FOR p;", shop)
	oneFinding(t, "CREATE SYNONYM IF NOT EXISTS s2 FOR p;\nCREATE TABLE s2 (a NUMBER);", 1, "The name `S2` is already taken.", shop)
	// 4) ALTER TABLE with unfollowed clauses still checks table existence.
	oneFinding(t, "ALTER TABLE nosuch ADD SUPPLEMENTAL LOG GROUP lg (a) ALWAYS;", 0, "Table `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE pv DROP SUPPLEMENTAL LOG GROUP lg;", 0, "`PV` is a view, not a table.", shop)
	// 5) CREATE/DROP MATERIALIZED VIEW and special CREATE TABLE name checks.
	oneFinding(t, "CREATE MATERIALIZED VIEW t AS SELECT 1 a FROM dual;", 0, "The name `T` is already taken.", shop)
	oneFinding(t, "CREATE MATERIALIZED VIEW mv AS SELECT 1 a FROM dual;", 0, "Materialized view `MV` already exists.", shop)
	noFinding(t, "CREATE MATERIALIZED VIEW IF NOT EXISTS t AS SELECT 1 a FROM dual;", shop)
	noFinding(t, "CREATE MATERIALIZED VIEW t ON PREBUILT TABLE AS SELECT 1 a FROM dual;", shop)
	oneFinding(t, "DROP MATERIALIZED VIEW nosuch;", 0, "Materialized view `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "DROP MATERIALIZED VIEW t;", 0, "`T` is a table, not a materialized view.", shop)
	oneFinding(t, "DROP MATERIALIZED VIEW pv;", 0, "`PV` is a view, not a materialized view.", shop)
	noFinding(t, "DROP MATERIALIZED VIEW IF EXISTS t;", shop)
	noFinding(t, "DROP MATERIALIZED VIEW mv;", shop)
	oneFinding(t, "CREATE TABLE t OF ty;", 0, "Table `T` already exists.", shop)
	noFinding(t, "CREATE TABLE IF NOT EXISTS t OF ty;", shop)
	// 6) ALTER PROCEDURE/FUNCTION/PACKAGE/TYPE/SYNONYM and CREATE VECTOR INDEX.
	oneFinding(t, "ALTER PROCEDURE t COMPILE;", 0, "`T` isn't a procedure.", shop)
	oneFinding(t, "ALTER FUNCTION pr COMPILE;", 0, "`PR` isn't a function.", shop)
	noFinding(t, "ALTER FUNCTION IF EXISTS pr COMPILE;", shop)
	noFinding(t, "ALTER FUNCTION f COMPILE;", shop)
	oneFinding(t, "ALTER PACKAGE f COMPILE;", 0, "`F` isn't a package.", shop)
	noFinding(t, "ALTER PACKAGE t COMPILE BODY;", shop)
	oneFinding(t, "ALTER TYPE t COMPILE;", 0, "`T` isn't a type.", shop)
	noFinding(t, "ALTER TYPE t COMPILE BODY;", shop)
	oneFinding(t, "ALTER SYNONYM t COMPILE;", 0, "`T` isn't a private synonym.", shop)
	noFinding(t, "ALTER SYNONYM IF EXISTS t EDITIONABLE;", shop)
	noFinding(t, "ALTER SYNONYM syn COMPILE;", shop)
	oneFinding(t, "CREATE VECTOR INDEX ix_n ON nokey (a) ORGANIZATION INMEMORY NEIGHBOR GRAPH;", 0, "The name `IX_N` is already taken.", shop)
	oneFinding(t, "CREATE VECTOR INDEX vix ON pv (id) ORGANIZATION INMEMORY NEIGHBOR GRAPH;", 0, "`PV` is a view, not a table.", shop)
	noFinding(t, "CREATE VECTOR INDEX IF NOT EXISTS ix_n ON nokey (a) ORGANIZATION INMEMORY NEIGHBOR GRAPH;", shop)
	oneFinding(t, "CREATE VECTOR INDEX vix ON mv (a) ORGANIZATION INMEMORY NEIGHBOR GRAPH;", 0, "`MV` is a materialized view, not a table.", shop)
	oneFinding(t, "CREATE VECTOR INDEX vix ON nosuch (a) ORGANIZATION INMEMORY NEIGHBOR GRAPH;", 0, "Table `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "CREATE VECTOR INDEX vix ON nokey (zz) ORGANIZATION INMEMORY NEIGHBOR GRAPH;", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	noFinding(t, "CREATE VECTOR INDEX vix ON syn (id) ORGANIZATION INMEMORY NEIGHBOR GRAPH;", shop)
	noFinding(t, "CREATE VECTOR INDEX vix ON nokey (a) ORGANIZATION INMEMORY NEIGHBOR GRAPH;", shop)
	// 7) ALTER MATERIALIZED VIEW: non-partition and mview-only partition clauses
	//    require a materialized view; table-like partition clauses allow tables
	//    and synonyms while rejecting missing names and non-table kinds.
	noFinding(t, "ALTER MATERIALIZED VIEW mv COMPILE;", shop)
	noFinding(t, "ALTER MATERIALIZED VIEW mv REFRESH COMPLETE;", shop)
	noFinding(t, "ALTER MATERIALIZED VIEW IF EXISTS nosuch COMPILE;", shop)
	noFinding(t, "ALTER MATERIALIZED VIEW IF EXISTS t COMPILE;", shop)
	oneFinding(t, "ALTER MATERIALIZED VIEW nosuch COMPILE;", 0, "Materialized view `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "ALTER MATERIALIZED VIEW t COMPILE;", 0, "`T` is a table, not a materialized view.", shop)
	oneFinding(t, "ALTER MATERIALIZED VIEW pv REFRESH FAST;", 0, "`PV` is a view, not a materialized view.", shop)
	oneFinding(t, "ALTER MATERIALIZED VIEW t SET PARTITIONING AUTOMATIC;", 0, "`T` is a table, not a materialized view.", shop)
	oneFinding(t, "ALTER MATERIALIZED VIEW syn MODIFY DEFAULT ATTRIBUTES COMPRESS;", 0, "`SYN` is a synonym, not a materialized view.", shop)
	noFinding(t, "ALTER MATERIALIZED VIEW t MODIFY PARTITION p1 LOGGING;", shop)
	noFinding(t, "ALTER MATERIALIZED VIEW syn ADD PARTITION;", shop)
	noFinding(t, "ALTER MATERIALIZED VIEW IF EXISTS nosuch ADD PARTITION;", shop)
	oneFinding(t, "ALTER MATERIALIZED VIEW nosuch ADD PARTITION;", 0, "Table `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "ALTER MATERIALIZED VIEW pv ADD PARTITION;", 0, "`PV` is a view, not a table.", shop)
	oneFinding(t, "ALTER MATERIALIZED VIEW sq ADD PARTITION;", 0, "`SQ` is a sequence, not a table.", shop)
	oneFinding(t, "ALTER MATERIALIZED VIEW f ADD PARTITION;", 0, "`F` is a function, not a table.", shop)
	// 8) CREATE MATERIALIZED VIEW ON PREBUILT TABLE, MV column refs, and
	//    missing SELECT aliases (ORA-00998).
	noFinding(t, "CREATE MATERIALIZED VIEW t ON PREBUILT TABLE AS SELECT id FROM p;", shop)
	noFinding(t, "CREATE MATERIALIZED VIEW IF NOT EXISTS nosuch ON PREBUILT TABLE AS SELECT 1 a FROM dual;", shop)
	oneFinding(t, "CREATE MATERIALIZED VIEW nosuch ON PREBUILT TABLE AS SELECT 1 a FROM dual;", 0, "Table `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "CREATE MATERIALIZED VIEW mv ON PREBUILT TABLE AS SELECT 1 a FROM dual;", 0, "`MV` is a materialized view, not a table.", shop)
	oneFinding(t, "CREATE MATERIALIZED VIEW pv ON PREBUILT TABLE AS SELECT 1 a FROM dual;", 0, "`PV` is a view, not a table.", shop)
	oneFinding(t, "CREATE MATERIALIZED VIEW mv2 NESTED TABLE missing_col STORE AS nt_s AS SELECT 1 a FROM dual;", 0, "Column `MISSING_COL` doesn't exist in `MV2`.", shop)
	oneFinding(t, "CREATE MATERIALIZED VIEW mv2 PARTITION BY RANGE (missing_col) (PARTITION p1 VALUES LESS THAN (10)) AS SELECT 1 a FROM dual;", 0, "Column `MISSING_COL` doesn't exist in `MV2`.", shop)
	noFinding(t, "CREATE MATERIALIZED VIEW mv2 NESTED TABLE a STORE AS nt_s AS SELECT 1 a FROM dual;", shop)
	noFinding(t, "CREATE MATERIALIZED VIEW mv2 (a) PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS THAN (10)) AS SELECT 1 FROM dual;", shop)
	noFinding(t, "CREATE MATERIALIZED VIEW mv2 NESTED TABLE a STORE AS nt_s AS SELECT * FROM t;", shop)
	oneFinding(t, "CREATE VIEW v2 AS SELECT 1 FROM dual;", 0, "Select expression in `V2` must have a column alias.", shop)
	oneFinding(t, "CREATE TABLE x2 AS SELECT 1 + 1 FROM dual;", 0, "Select expression in `X2` must have a column alias.", shop)
	noFinding(t, "CREATE FORCE VIEW v2 AS SELECT 1 FROM dual;", shop)
	noFinding(t, "CREATE VIEW v2 (a) AS SELECT 1 FROM dual;", shop)
	noFinding(t, "CREATE VIEW v2 AS SELECT sysdate, user, (id) FROM p;", shop)
	noFinding(t, "CREATE VIEW v2 AS SELECT * FROM p;", shop)
	// 9) CREATE/ALTER TABLE missing datatypes, constraint/LOB/VARRAY/nested/partition
	//    column checks, and MODIFY PARTITION BY on materialized views.
	oneFinding(t, "CREATE TABLE x2 (a);", 0, "Column `A` must specify a datatype.", shop)
	oneFinding(t, "CREATE TABLE x2 (a GENERATED ALWAYS AS IDENTITY);", 0, "Column `A` must specify a datatype.", shop)
	noFinding(t, "CREATE TABLE x2 (a NUMBER, b AS (a + 1));", shop)
	noFinding(t, "CREATE TABLE x2 (a REFERENCES p (id));", shop)
	noFinding(t, "CREATE TABLE x2 (a DOMAIN d);", shop)
	oneFinding(t, "ALTER TABLE nokey ADD (b);", 0, "Column `B` must specify a datatype.", shop)
	noFinding(t, "ALTER TABLE nokey ADD (b AS (a + 1));", shop)
	noFinding(t, "ALTER TABLE nokey ADD (b NUMBER) ADD CONSTRAINT pk_nk PRIMARY KEY (b);", shop)
	oneFinding(t, "ALTER TABLE nokey ADD CONSTRAINT uq_zz UNIQUE (zz);", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "ALTER TABLE nokey DROP UNIQUE (zz);", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "ALTER TABLE nokey ENABLE UNIQUE (zz);", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "ALTER TABLE nokey MODIFY NESTED TABLE zz RETURN AS LOCATOR;", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "ALTER TABLE nokey MODIFY LOB (zz) (RETENTION);", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "ALTER TABLE nokey MODIFY VARRAY zz (RETENTION MAX);", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "ALTER TABLE nosuch MODIFY VARRAY a (RETENTION MAX);", 0, "Table `NOSUCH` doesn't exist.", shop)
	oneFinding(t, "ALTER TABLE nokey MODIFY COLUMN zz SUBSTITUTABLE AT ALL LEVELS;", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "ALTER TABLE mv MODIFY PARTITION BY HASH (id) PARTITIONS 2;", 0, "`MV` is a materialized view, not a table.", shop)
	oneFinding(t, "ALTER TABLE nokey MODIFY PARTITION BY RANGE (zz) (PARTITION p1 VALUES LESS THAN (10));", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	noFinding(t, "ALTER TABLE nokey MODIFY PARTITION BY HASH (a) PARTITIONS 2;", shop)
	// 10) CREATE INDEX ON CLUSTER wrong-kind and bare index column checks.
	oneFinding(t, "CREATE INDEX i ON CLUSTER t;", 0, "`T` is a table, not a cluster.", shop)
	oneFinding(t, "CREATE INDEX i ON CLUSTER pv;", 0, "`PV` is a view, not a cluster.", shop)
	noFinding(t, "CREATE INDEX i ON CLUSTER nosuch_cluster;", shop)
	oneFinding(t, "CREATE INDEX i ON nokey (zz);", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "CREATE INDEX i ON nokey (a, a);", 0, "This statement lists column `A` twice.", shop)
	noFinding(t, "CREATE INDEX i ON nokey (a);", shop)
	// 11) COMMENT ON TABLE, COLUMN, MATERIALIZED VIEW, and unsynced categories.
	oneFinding(t, "COMMENT ON TABLE nosuch IS 'x';", 0, "`NOSUCH` doesn't exist.", shop)
	oneFinding(t, "COMMENT ON TABLE mv IS 'x';", 0, "`MV` is a materialized view, not a table or view.", shop)
	oneFinding(t, "COMMENT ON TABLE sq IS 'x';", 0, "`SQ` is a sequence, not a table or view.", shop)
	noFinding(t, "COMMENT ON TABLE t IS 'x';\nCOMMENT ON TABLE pv IS 'x';\nCOMMENT ON TABLE syn IS 'x';", shop)
	oneFinding(t, "COMMENT ON MATERIALIZED VIEW nosuch IS 'x';", 0, "`NOSUCH` doesn't exist.", shop)
	noFinding(t, "COMMENT ON MATERIALIZED VIEW mv IS 'x';\nCOMMENT ON MATERIALIZED VIEW t IS 'x';", shop)
	oneFinding(t, "COMMENT ON COLUMN t.zz IS 'x';", 0, "Column `ZZ` doesn't exist in `T`.", shop)
	oneFinding(t, "COMMENT ON COLUMN nosuch.a IS 'x';", 0, "`NOSUCH` doesn't exist.", shop)
	noFinding(t, "COMMENT ON COLUMN t.n IS 'x';\nCOMMENT ON COLUMN pv.zz IS 'x';\nCOMMENT ON INDEXTYPE nosuch IS 'x';", shop)
	// 12) DML target columns, explicit counts, LOG ERRORS, MERGE USING,
	//     unaliased self-joins, and WITH READ ONLY inline view targets.
	oneFinding(t, "INSERT INTO nokey (zz) VALUES (1);", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "INSERT INTO nokey (a, a) VALUES (1, 2);", 0, "This statement names column `A` twice.", shop)
	oneFinding(t, "INSERT INTO nokey (a) VALUES (1, 2);", 0, "INSERT column count does not match VALUES count.", shop)
	oneFinding(t, "INSERT INTO p (id, note) SELECT 1 FROM dual;", 0, "INSERT column count does not match SELECT list count.", shop)
	noFinding(t, "INSERT INTO p (id, note) SELECT * FROM p;", shop)
	oneFinding(t, "UPDATE nokey SET zz = 1;", 0, "Column `ZZ` doesn't exist in `NOKEY`.", shop)
	oneFinding(t, "UPDATE nokey SET a = 1 FROM nokey;", 0, "Ambiguous unqualified column in UPDATE with unaliased self-join on `NOKEY`.", shop)
	noFinding(t, "UPDATE nokey x SET x.a = 1 FROM nokey y WHERE x.a = y.a;", shop)
	oneFinding(t, "MERGE INTO nokey USING nosuch ON (nokey.a = nosuch.a) WHEN MATCHED THEN UPDATE SET nokey.a = 1;", 0, "`NOSUCH` doesn't exist.", shop)
	oneFinding(t, "MERGE INTO nokey USING nokey ON (1 = 1) WHEN NOT MATCHED THEN INSERT VALUES (1);", 0, "Ambiguous unaliased self-join in MERGE on `NOKEY`.", shop)
	noFinding(t, "MERGE INTO nokey t1 USING nokey t2 ON (t1.a = t2.a) WHEN MATCHED THEN UPDATE SET t1.a = 1;", shop)
	oneFinding(t, "DELETE FROM nokey LOG ERRORS;", 0, "`ERR$_NOKEY` doesn't exist.", shop)
	oneFinding(t, "UPDATE nokey SET a = 1 LOG ERRORS INTO nosuch_err;", 0, "`NOSUCH_ERR` doesn't exist.", shop)
	noFinding(t, "CREATE TABLE err$_nokey (a NUMBER);\nDELETE FROM nokey LOG ERRORS;", shop)
	oneFinding(t, "UPDATE (SELECT id FROM p WITH READ ONLY) SET id = 1;", 0, "Cannot modify a subquery declared WITH READ ONLY.", shop)
	oneFinding(t, "DELETE (SELECT id FROM p WITH READ ONLY);", 0, "Cannot modify a subquery declared WITH READ ONLY.", shop)
	noFinding(t, "DELETE nokey FROM (SELECT 1 FROM dual WITH READ ONLY);", shop)
}

func TestWalkThroughOracleEffects(t *testing.T) {
	shop := review.Target{Schema: walkSchema(false)}
	// After IF [NOT] EXISTS on a name the walk does not know, the name is
	// there, or gone from its kind's space; an ALTER ... IF EXISTS leaves
	// it as it was.
	oneFinding(t, "CREATE TABLE IF NOT EXISTS other.x (a NUMBER);\nCREATE TABLE other.x (a NUMBER);", 1, "The name `OTHER.X` is already taken.", shop)
	oneFinding(t, "DROP TABLE IF EXISTS other.x;\nALTER TABLE other.x ADD (a NUMBER);", 1, "Table `OTHER.X` doesn't exist.", shop)
	oneFinding(t, "CREATE INDEX IF NOT EXISTS i ON nokey (a);\nCREATE INDEX i ON p (note);", 1, "The name `I` is already taken.", shop)
	oneFinding(t, "DROP INDEX IF EXISTS i;\nDROP INDEX i;", 1, "Index `I` doesn't exist.", shop)
	noFinding(t, "ALTER SEQUENCE IF EXISTS other.sq INCREMENT BY 2;\nDROP SEQUENCE other.sq;", shop)
	noFinding(t, "DROP SYNONYM IF EXISTS s;\nDROP SYNONYM s;", shop)
	// A dropped table takes its indexes, constraint names, and columns.
	noFinding(t, "DROP TABLE t;\nCREATE INDEX ix_n ON nokey (a);\nALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nCREATE TABLE t (z NUMBER);\nALTER TABLE t ADD (code NUMBER);", shop)
	oneFinding(t, "DROP TABLE t;\nALTER TABLE t ADD (x NUMBER);", 1, "Table `T` doesn't exist.", shop)
	oneFinding(t, "DROP TABLE t;\nCREATE TABLE t (z NUMBER);\nALTER TABLE t DROP COLUMN code;", 2, "Column `CODE` doesn't exist in `T`.", shop)
	// A renamed table keeps its columns and leaves its old name.
	oneFinding(t, "ALTER TABLE t RENAME TO t2;\nALTER TABLE t2 DROP COLUMN zz;", 1, "Column `ZZ` doesn't exist in `T2`.", shop)
	oneFinding(t, "ALTER TABLE t RENAME TO t2;\nALTER TABLE t ADD (x NUMBER);", 1, "Table `T` doesn't exist.", shop)
	oneFinding(t, "RENAME t TO t2;\nCREATE TABLE t2 (a NUMBER);", 1, "Table `T2` already exists.", shop)
	noFinding(t, "ALTER TABLE t RENAME TO t2;\nDROP TABLE t2;\nCREATE INDEX ix_n ON nokey (a);", shop)
	// An index or trigger may be in another schema than its table.
	noFinding(t, "CREATE TABLE other.o (a NUMBER);\nCREATE INDEX oi ON other.o (a);\nCREATE TRIGGER otr BEFORE INSERT ON other.o FOR EACH ROW BEGIN NULL; END;\n/\nDROP TABLE other.o;\nCREATE INDEX oi ON nokey (a);\nCREATE TRIGGER otr BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;", shop)
	oneFinding(t, "CREATE TABLE other.o (a NUMBER);\nCREATE INDEX oi ON other.o (a);\nCREATE INDEX oi ON nokey (a);", 2, "The name `OI` is already taken.", shop)
	// ALTER TABLE sees the table as it was.
	noFinding(t, "ALTER TABLE t RENAME COLUMN n TO m;\nALTER TABLE t ADD (n NUMBER);\nALTER TABLE t DROP (m);", shop)
	oneFinding(t, "ALTER TABLE t DROP COLUMN n;\nALTER TABLE t DROP COLUMN n;", 1, "Column `N` doesn't exist in `T`.", shop)
	// A column's constraints may go with it.
	noFinding(t, "ALTER TABLE t DROP COLUMN n;\nALTER TABLE t DROP CONSTRAINT ck_n;", shop)
	noFinding(t, "ALTER TABLE t MODIFY (code NULL);\nALTER TABLE nokey ADD CONSTRAINT nn_code CHECK (a > 0);", shop)
	// A nested table's storage table goes with its parent, or its column.
	oneFinding(t, "CREATE TABLE child (a NUMBER);", 0, "Table `CHILD` already exists.", shop)
	noFinding(t, "DROP TABLE parent;\nCREATE TABLE child (a NUMBER);", shop)
	noFinding(t, "ALTER TABLE parent DROP COLUMN c;\nCREATE TABLE child (a NUMBER);", shop)
	noFinding(t, "ALTER TABLE parent SET UNUSED COLUMN c;\nALTER TABLE parent DROP UNUSED COLUMNS;\nCREATE TABLE child (a NUMBER);", shop)
	noFinding(t, "ALTER TABLE nokey ADD (c nt_ty) NESTED TABLE c STORE AS s2;\nDROP TABLE s2;", shop)
	oneFinding(t, "DROP TABLE other.o;\nCREATE TABLE child (a NUMBER);", 1, "Table `CHILD` already exists.", shop)
	// A materialized view log's table goes with its master.
	oneFinding(t, "CREATE TABLE mlog$_p (a NUMBER);", 0, "Table `MLOG$_P` already exists.", shop)
	noFinding(t, "DROP TABLE p;\nCREATE TABLE mlog$_p (a NUMBER);\nCREATE INDEX i_mlog$_p ON nokey (a);", shop)
	oneFinding(t, "DROP TABLE other.o;\nCREATE INDEX i_mlog$_p ON nokey (a);", 1, "The name `I_MLOG$_P` is already taken.", shop)
	// DROP TYPE ... FORCE takes the columns of the type, and their
	// constraints and indexes.
	oneFinding(t, "DROP TYPE ty;\nALTER TABLE t ADD (n NUMBER);", 1, "Column `N` already exists in `T`.", shop)
	noFinding(t, "DROP TYPE ty FORCE;\nALTER TABLE t ADD (n NUMBER);\nALTER TABLE nokey ADD CONSTRAINT ck_n CHECK (a > 0);\nCREATE INDEX ix_n ON nokey (a);", shop)
	// Columns the sync may not list: a period's, an object table's.
	noFinding(t, "ALTER TABLE periodic DROP COLUMN vt_start;\nALTER TABLE ot DROP COLUMN xmldata;\nALTER TABLE nokey DROP COLUMN ora_archive_state;", shop)
	oneFinding(t, "ALTER TABLE ot ADD (x NUMBER);", 0, "Column `X` already exists in `OT`.", shop)
	// A key's index takes a name the walk does not know.
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT uq UNIQUE (a);\nCREATE INDEX uq ON p (note);", shop)
	noFinding(t, "ALTER TABLE nokey ADD CONSTRAINT uq UNIQUE (a) USING INDEX (CREATE INDEX ux ON nokey (a));\nCREATE INDEX i ON p (note);\nDROP INDEX i;\nDROP INDEX ux;", shop)
	// A view's triggers and constraints go with it.
	noFinding(t, "CREATE VIEW v2 AS SELECT 1 a FROM dual WITH CHECK OPTION CONSTRAINT vck;\nDROP VIEW v2;\nALTER TABLE nokey ADD CONSTRAINT vck CHECK (a > 0);", shop)
	oneFinding(t, "CREATE VIEW v2 AS SELECT 1 a FROM dual WITH CHECK OPTION CONSTRAINT vck;\nALTER TABLE nokey ADD CONSTRAINT vck CHECK (a > 0);", 1, "Constraint `VCK` already exists.", shop)
	noFinding(t, "CREATE VIEW v2 AS SELECT 1 a FROM dual WITH CHECK OPTION CONSTRAINT vck;\nCREATE OR REPLACE VIEW v2 AS SELECT 1 a FROM dual WITH CHECK OPTION CONSTRAINT vck;", shop)
}

func TestWalkThroughOracleSession(t *testing.T) {
	shop := review.Target{Schema: walkSchema(false)}
	// Another schema is unknown, until the change makes its objects.
	noFinding(t, "ALTER SESSION SET CURRENT_SCHEMA = other;\nCREATE TABLE t (a NUMBER);", shop)
	oneFinding(t, "ALTER SESSION SET CURRENT_SCHEMA = other;\nCREATE TABLE x (a NUMBER);\nCREATE TABLE other.x (a NUMBER);", 2, "Table `OTHER.X` already exists.", shop)
	oneFinding(t, "ALTER SESSION SET CURRENT_SCHEMA = other;\nCREATE TABLE shop.t (a NUMBER);", 1, "Table `SHOP.T` already exists.", shop)
	oneFinding(t, "ALTER SESSION SET CURRENT_SCHEMA = other;\nALTER SESSION SET CURRENT_SCHEMA = shop;\nCREATE TABLE t (a NUMBER);", 2, "Table `T` already exists.", shop)
	oneFinding(t, "ALTER SESSION SET NLS_DATE_FORMAT = 'YYYY';\nCREATE TABLE t (a NUMBER);", 1, "Table `T` already exists.", shop)
	// What the walk does not follow ends it.
	noFinding(t, "ALTER SESSION SET EDITION = e1;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "BEGIN NULL; END;\n/\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CALL pr();\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "FLASHBACK TABLE t TO BEFORE DROP;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CREATE MATERIALIZED VIEW mv2 AS SELECT 1 a FROM dual;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CREATE TRIGGER st AFTER LOGON ON SCHEMA BEGIN NULL; END;\n/\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CREATE TRIGGER st AFTER CREATE ON shop.SCHEMA BEGIN NULL; END;\n/\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CREATE TRIGGER st BEFORE DDL ON other.SCHEMA BEGIN NULL; END;\n/\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CREATE TRIGGER st AFTER STARTUP ON DATABASE BEGIN NULL; END;\n/\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CREATE TABLE \"é\" (a NUMBER);\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "ALTER TABLE t ADD SUPPLEMENTAL LOG GROUP lg (n) ALWAYS;\nCREATE TABLE t (a NUMBER);", shop)
	// A table's column list carries one too, and the log group ends the walk
	// as it does after the ALTER above.
	noFinding(t, "CREATE TABLE t (a NUMBER, SUPPLEMENTAL LOG GROUP lg (a) ALWAYS);\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "INSERT INTO t@lnk VALUES (1);\nCREATE TABLE t (a NUMBER);", shop)
	// A string is no name.
	oneFinding(t, "INSERT INTO nokey VALUES ('é');\nCREATE TABLE t (a NUMBER);", 1, "Table `T` already exists.", shop)
}

// TestWalkThroughOracleUserCode pins where a statement may run code that
// changes a schema, which ends the walk after its checks (spec §6).
func TestWalkThroughOracleUserCode(t *testing.T) {
	shop := review.Target{Schema: walkSchema(false)}
	triggered := review.Target{Schema: walkSchema(true)}
	// A query that calls a function of a schema.
	noFinding(t, "SELECT f(1) FROM dual;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "SELECT upper(shop.f(1)) FROM dual;\nCREATE TABLE t (a NUMBER);", shop)
	oneFinding(t, "SELECT upper('a'), nvl(1, 2), count(*) FROM dual;\nCREATE TABLE t (a NUMBER);", 1, "Table `T` already exists.", shop)
	noFinding(t, "CREATE TABLE x AS SELECT f(1) a FROM dual;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CREATE INDEX i ON nokey (f(a));\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CREATE TABLE x AS SELECT f(1) a FROM dual;\nCREATE TABLE x (a NUMBER);", shop)
	// A name that may call a function without parentheses: a bare name no
	// table in scope has as a column, or a qualifier that is no table in
	// scope. A column, a sequence, and a builtin Oracle takes first are not.
	noFinding(t, "SELECT f FROM dual;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "INSERT INTO nokey (a) VALUES (f);\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "UPDATE t SET n = f WHERE id = 1;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "DELETE FROM t WHERE id = 1 AND f = 1;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "INSERT INTO nokey (a) VALUES (pkg.x);\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "INSERT INTO nokey (a) SELECT a FROM nokey WHERE shop.f = 1;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "UPDATE t SET n = 1 WHERE id IN (SELECT id FROM p);\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "INSERT INTO nokey (a) VALUES (sq2.nextval);\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "SELECT x FROM pv;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "CREATE INDEX i ON nokey (a + b);\nCREATE TABLE t (a NUMBER);", shop)
	oneFinding(t, "CREATE INDEX i ON syn (p_id, b);\nCREATE INDEX j ON nokey (a + 1);\nCREATE TABLE t (a NUMBER);", 2, "Table `T` already exists.", shop)
	oneFinding(t, "UPDATE t x SET x.n = 1 WHERE x.id = 1 AND code = 'a' RETURNING n INTO :v;\nCREATE TABLE t (a NUMBER);", 1, "Table `T` already exists.", shop)
	oneFinding(t, "INSERT INTO nokey (a) VALUES (sq.nextval + 1);\nINSERT INTO nokey (a) SELECT p.id FROM p JOIN t ON t.p_id = p.id WHERE note = current_date;\nCREATE TABLE t (a NUMBER);", 2, "Table `T` already exists.", shop)
	oneFinding(t, "MERGE INTO nokey k USING p ON (k.a = p.id) WHEN NOT MATCHED THEN INSERT (a) VALUES (p.id);\nDELETE FROM shop.t WHERE shop.t.id = 1;\nCREATE INDEX i ON nokey (a);\nCREATE TABLE t (a NUMBER);", 3, "Table `T` already exists.", shop)
	oneFinding(t, "CREATE TABLE x (a NUMBER);\nUPDATE x SET a = 1 WHERE a = 2;\nCREATE TABLE x (b NUMBER);", 2, "Table `X` already exists.", shop)
	// A DELETE's FROM tables are in scope, as an UPDATE's are.
	oneFinding(t, "DELETE nokey FROM p WHERE nokey.a = p.id;\nCREATE TABLE t (a NUMBER);", 1, "Table `T` already exists.", shop)
	// A quoted name of a function of a syntax of its own is a plain call.
	noFinding(t, "SELECT \"EXTRACT\"(1) FROM dual;\nCREATE TABLE t (a NUMBER);", shop)
	oneFinding(t, "SELECT extract(year FROM sysdate) FROM dual;\nCREATE TABLE t (a NUMBER);", 1, "Table `T` already exists.", shop)
	// A domain index's indextype runs code of its own on the table's DDL
	// and DML, and to build or drop the index.
	for _, q := range []string{
		"CREATE INDEX x ON nokey (a) INDEXTYPE IS ctxsys.context;",
		"DROP INDEX doc_ctx;",
		"ALTER INDEX doc_ctx REBUILD;",
		"DROP TABLE doc;",
		"TRUNCATE TABLE doc;",
		"ALTER TABLE doc ADD (z NUMBER);",
		"RENAME doc TO doc2;",
		"INSERT INTO doc (id) VALUES (1);",
	} {
		noFinding(t, q+"\nCREATE TABLE t (a NUMBER);", shop)
	}
	oneFinding(t, "ALTER TABLE doc ADD (body NUMBER);", 0, "Column `BODY` already exists in `DOC`.", shop)
	// DML that may fire a trigger: any in the schema or the change, or a
	// target the walk does not know as a table.
	oneFinding(t, "INSERT INTO nokey VALUES (1);\nUPDATE t SET n = 1;\nDELETE FROM p;\nCREATE TABLE t (a NUMBER);", 3, "Table `T` already exists.", shop)
	noFinding(t, "INSERT INTO nokey VALUES (1);\nCREATE TABLE t (a NUMBER);", triggered)
	noFinding(t, "CREATE TRIGGER tr2 BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;\n/\nINSERT INTO nokey VALUES (1);\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "INSERT INTO pv VALUES (1, 'a');\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "INSERT INTO pubt VALUES (1, 1);\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "UPDATE (SELECT * FROM t) SET n = 1;\nCREATE TABLE t (a NUMBER);", shop)
	noFinding(t, "INSERT INTO nokey VALUES (f(1));\nCREATE TABLE t (a NUMBER);", shop)
	// With another spec of the plan, only the change's own facts count, and
	// DML ends the walk.
	other := review.Target{Schema: walkSchema(false), OtherSpecs: true}
	noFinding(t, "CREATE TABLE t (a NUMBER);", other)
	oneFinding(t, "CREATE TABLE x (a NUMBER);\nCREATE TABLE x (a NUMBER);", 1, "Table `X` already exists.", other)
	noFinding(t, "CREATE TABLE x (a NUMBER);\nINSERT INTO x VALUES (1);\nCREATE TABLE x (a NUMBER);", other)
}

func TestWalkThroughOracleDMLTargets(t *testing.T) {
	shop := review.Target{Schema: walkSchema(false)}
	// Nothing holds the name: the schema's tables, views, and synonyms, and
	// for an unqualified name the public synonyms, show it.
	oneFinding(t, "INSERT INTO nosuch VALUES (1);", 0, "`NOSUCH` doesn't exist.", shop)
	oneFinding(t, "UPDATE shop.nosuch SET a = 1;", 0, "`SHOP.NOSUCH` doesn't exist.", shop)
	oneFinding(t, "MERGE INTO nosuch k USING p ON (k.a = p.id) WHEN MATCHED THEN UPDATE SET k.a = 1;", 0, "`NOSUCH` doesn't exist.", shop)
	oneFinding(t, "INSERT ALL INTO nokey VALUES (1) INTO nosuch VALUES (1) SELECT 1 FROM dual;", 0, "`NOSUCH` doesn't exist.", shop)
	oneFinding(t, "DROP TABLE nokey;\nINSERT INTO nokey VALUES (1);", 1, "`NOKEY` doesn't exist.", shop)
	// A qualified name skips the public synonyms.
	noFinding(t, "INSERT INTO pubt VALUES (1, 1);", shop)
	oneFinding(t, "INSERT INTO shop.pubt VALUES (1, 1);", 0, "`SHOP.PUBT` doesn't exist.", shop)
	noFinding(t, "INSERT INTO \"PUBLIC\".pubt VALUES (1, 1);", shop)
	// Synonyms, the synced ones and the change's own.
	noFinding(t, "INSERT INTO syn (id) VALUES (1);", shop)
	oneFinding(t, "DROP SYNONYM syn;\nINSERT INTO syn (id) VALUES (1);", 1, "`SYN` doesn't exist.", shop)
	noFinding(t, "CREATE SYNONYM s2 FOR t;\nINSERT INTO s2 (id) VALUES (1);", shop)
	noFinding(t, "CREATE PUBLIC SYNONYM ps FOR t;\nINSERT INTO ps (id) VALUES (1);", shop)
	oneFinding(t, "DROP PUBLIC SYNONYM pubt;\nINSERT INTO pubt VALUES (1, 1);", 1, "`PUBT` doesn't exist.", shop)
	noFinding(t, "DROP TABLE nokey;\nCREATE PUBLIC SYNONYM nokey FOR p;\nINSERT INTO nokey (id) VALUES (1);", shop)
	// The sync leaves out a synonym whose object is gone or out of reach,
	// which a later object or grant may restore: after one, no synonym's
	// absence is trusted, but a name the change itself frees.
	noFinding(t, "CREATE TABLE nosuch2 (a NUMBER);\nINSERT INTO s_gone VALUES (1);", shop)
	noFinding(t, "CREATE VIEW v2 AS SELECT 1 a FROM dual;\nINSERT INTO s_gone VALUES (1);", shop)
	noFinding(t, "CREATE SYNONYM s2 FOR t;\nINSERT INTO s_gone VALUES (1);", shop)
	noFinding(t, "CREATE PUBLIC SYNONYM ps FOR t;\nINSERT INTO s_gone VALUES (1);", shop)
	noFinding(t, "GRANT SELECT ON t TO PUBLIC;\nINSERT INTO s_gone VALUES (1);", shop)
	noFinding(t, "RENAME nokey TO nokey2;\nINSERT INTO s_gone VALUES (1);", shop)
	noFinding(t, "ALTER TABLE nokey RENAME TO nokey2;\nINSERT INTO s_gone VALUES (1);", shop)
	oneFinding(t, "RENAME syn TO syn2;\nINSERT INTO shop.syn (id) VALUES (1);", 1, "`SHOP.SYN` doesn't exist.", shop)
	oneFinding(t, "CREATE INDEX i ON nokey (a);\nINSERT INTO s_gone VALUES (1);", 1, "`S_GONE` doesn't exist.", shop)
	// An object of another kind refuses DML with an error of its own.
	oneFinding(t, "DELETE FROM sq;", 0, "`SQ` is a sequence, not a table or view.", shop)
	oneFinding(t, "INSERT INTO pkg VALUES (1);", 0, "`PKG` is a package, not a table or view.", shop)
	noFinding(t, "INSERT INTO pv VALUES (1, 'a');", shop)
	// What the walk does not know.
	noFinding(t, "ALTER SESSION SET CURRENT_SCHEMA = other;\nINSERT INTO nosuch VALUES (1);", shop)
	noFinding(t, "INSERT INTO other.nosuch VALUES (1);", shop)
	noFinding(t, "INSERT INTO nosuch VALUES (1);", review.Target{Schema: walkSchema(false), OtherSpecs: true})
}

func TestWalkThroughOracleTargets(t *testing.T) {
	// A target without the database's name is not walked.
	noFinding(t, "CREATE TABLE x (a NUMBER);\nCREATE TABLE x (a NUMBER);", review.Target{Schema: &metadata.DatabaseSchemaMetadata{}})
	// OtherSpecs is part of the group key.
	schema := walkSchema(false)
	result, err := Review(context.Background(), "CREATE TABLE t (a NUMBER);", review.Options{Rules: []review.Rule{review.WalkThrough}}, []review.Target{{Schema: schema}, {Schema: schema, OtherSpecs: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 1 || len(result.Findings[0].Targets) != 1 || result.Findings[0].Targets[0] != 0 {
		t.Errorf("findings = %+v, want one on target 0", result.Findings)
	}
}

// TestWalkThroughOracleSnapshot pins what the snapshot trusts (spec §3).
func TestWalkThroughOracleSnapshot(t *testing.T) {
	s := newOracleSnapshot(walkSchema(true))
	key := func(space, name string, path ...string) walk.Key {
		return walk.Key{Space: space, Parent: walk.Path(append([]string{"SHOP"}, path...)...), Name: name}
	}
	tests := []struct {
		key  walk.Key
		fact walk.Fact
		kind string
	}{
		{key(spaceObject, "T"), walk.Present, kindTable},
		{key(spaceObject, "MV"), walk.Present, kindMView},
		{key(spaceObject, "PKG"), walk.Present, kindPackage},
		{key(spaceObject, "NOSUCH"), walk.Unknown, ""},
		{key(spaceTable, "MV"), walk.Present, kindMView},
		{key(spaceTable, "PV"), walk.Absent, ""},
		{key(spaceTable, "NOSUCH"), walk.Absent, ""},
		{key(spaceTable, "SYS_IOT_OVER_1"), walk.Unknown, ""},
		{key(spaceView, "PV"), walk.Present, kindView},
		{key(spaceView, "T"), walk.Absent, ""},
		{key(spaceSequence, "SQ"), walk.Present, kindSequence},
		{key(spaceSequence, "ISEQ$$_1"), walk.Unknown, ""},
		{key(spaceColumn, "CODE", "T"), walk.Present, ""},
		{key(spaceColumn, "ZZ", "T"), walk.Absent, ""},
		{key(spaceColumn, "SYS_NC00001$", "T"), walk.Unknown, ""},
		{key(spaceColumn, "ORA_ARCHIVE_STATE", "T"), walk.Unknown, ""},
		{key(spaceColumn, "ZZ", "PERIODIC"), walk.Unknown, ""},
		{key(spaceColumn, "ZZ", "OT"), walk.Unknown, ""},
		{key(spaceColumn, "ZZ", "MV"), walk.Unknown, ""},
		{key(spacePrimary, "", "T"), walk.Present, ""},
		{key(spacePrimary, "", "NOKEY"), walk.Absent, ""},
		{key(spacePrimary, "", "MV"), walk.Unknown, ""},
		{key(spaceTableConstraint, "UK_CODE", "T"), walk.Present, ""},
		{key(spaceTableConstraint, "NOSUCH", "T"), walk.Unknown, ""},
		{key(spaceConstraint, "FK_P"), walk.Present, tablePath("SHOP", "T")},
		{key(spaceConstraint, "NOSUCH"), walk.Unknown, ""},
		{key(spaceIndex, "IX_N"), walk.Present, tablePath("SHOP", "T")},
		{key(spaceIndex, "MV_IX"), walk.Present, tablePath("SHOP", "MV")},
		{key(spaceIndex, "PK_T"), walk.Unknown, ""},
		{key(spaceIndex, "NOSUCH"), walk.Unknown, ""},
		{key(spaceTrigger, "TR"), walk.Present, tablePath("SHOP", "TRG")},
		{key(spaceTrigger, "NOSUCH"), walk.Unknown, ""},
		{walk.Key{Space: spaceTable, Parent: walk.Path("OTHER"), Name: "T"}, walk.Unknown, ""},
	}
	for _, tt := range tests {
		f, kind := s.Lookup(tt.key)
		if f != tt.fact || kind != tt.kind {
			t.Errorf("%s %q in %q = %v %q, want %v %q", tt.key.Space, tt.key.Name, walk.Names(tt.key.Parent), f, kind, tt.fact, tt.kind)
		}
	}
	if strings.Join(s.storage, ",") != "CHILD" || !s.domainTables["DOC"] || !s.domainIndexes["DOC_CTX"] {
		t.Errorf("storage = %q, domain tables = %v, domain indexes = %v", s.storage, s.domainTables, s.domainIndexes)
	}
	if got := quotedNames(`"A" > 0 AND b = 'say "B"' AND "C" < 'it''s'`); strings.Join(got, ",") != "A,C" {
		t.Errorf("quotedNames = %q", got)
	}
}

// walkForms are the statement forms of the Oracle spec's §5 table, each on
// walkSchema. A form the table lists must lower; one that ends the walk by
// design is Last.
var walkForms = []walktest.Form{
	{SQL: "CREATE TABLE x (a NUMBER CONSTRAINT pk_x PRIMARY KEY, b NUMBER CONSTRAINT nn_b NOT NULL, c NUMBER UNIQUE, CONSTRAINT ck_x CHECK (a > 0), CONSTRAINT fk_x FOREIGN KEY (b) REFERENCES p (id));"},
	{SQL: "CREATE GLOBAL TEMPORARY TABLE g (a NUMBER) ON COMMIT PRESERVE ROWS;"},
	{SQL: "CREATE TABLE IF NOT EXISTS t (a NUMBER);"},
	{SQL: "CREATE TABLE x (a NUMBER, b CLOB) LOB (b) STORE AS lseg PARTITION BY HASH (a) PARTITIONS 2;"},
	{SQL: "CREATE TABLE x AS SELECT * FROM t;"},
	{SQL: "ALTER TABLE t ADD (x NUMBER CONSTRAINT nn_x NOT NULL, y NUMBER);"},
	{SQL: "ALTER TABLE t ADD z NUMBER;"},
	{SQL: "ALTER TABLE t MODIFY (n NUMBER(10), code NULL);"},
	{SQL: "ALTER TABLE t DROP COLUMN n;"},
	{SQL: "ALTER TABLE t DROP (n, p_id) CASCADE CONSTRAINTS;"},
	{SQL: "ALTER TABLE t SET UNUSED COLUMN n;"},
	{SQL: "ALTER TABLE t SET UNUSED (n, p_id);"},
	{SQL: "ALTER TABLE t RENAME COLUMN n TO m;"},
	{SQL: "ALTER TABLE nokey ADD CONSTRAINT pk_nokey PRIMARY KEY (a);"},
	{SQL: "ALTER TABLE nokey ADD (CONSTRAINT ck_a CHECK (a > 0), CONSTRAINT uq_a UNIQUE (a));"},
	{SQL: "ALTER TABLE nokey ADD CONSTRAINT fk_a FOREIGN KEY (a) REFERENCES p (id);"},
	{SQL: "ALTER TABLE t DROP CONSTRAINT ck_n;"},
	{SQL: "ALTER TABLE p DROP PRIMARY KEY CASCADE;"},
	{SQL: "ALTER TABLE t DROP UNIQUE (code);"},
	{SQL: "ALTER TABLE t RENAME CONSTRAINT ck_n TO ck_m;"},
	{SQL: "ALTER TABLE t DISABLE CONSTRAINT ck_n;\nALTER TABLE t ENABLE CONSTRAINT ck_n;\nALTER TABLE t MODIFY CONSTRAINT ck_n DISABLE;"},
	{SQL: "ALTER TABLE t DISABLE PRIMARY KEY;\nALTER TABLE t ENABLE UNIQUE (code);\nALTER TABLE t MODIFY PRIMARY KEY ENABLE;"},
	{SQL: "ALTER TABLE t RENAME TO t2;"},
	{SQL: "ALTER TABLE t LOGGING;\nALTER TABLE t READ ONLY;\nALTER TABLE t PARALLEL 4;\nALTER TABLE t ENABLE ROW MOVEMENT;\nALTER TABLE t MOVE;\nALTER TABLE t SHRINK SPACE;"},
	{SQL: "ALTER TABLE t ENABLE ALL TRIGGERS;\nALTER TABLE t DISABLE TABLE LOCK;\nALTER TABLE t DROP UNUSED COLUMNS;"},
	{SQL: "ALTER TABLE t ADD PARTITION p2 VALUES LESS THAN (10);\nALTER TABLE t TRUNCATE PARTITION p1;"},
	{SQL: "DROP TABLE t CASCADE CONSTRAINTS PURGE;"},
	{SQL: "DROP TABLE IF EXISTS nosuch;"},
	{SQL: "TRUNCATE TABLE t;"},
	{SQL: "CREATE UNIQUE INDEX i ON nokey (a);\nCREATE BITMAP INDEX j ON t (n);\nCREATE INDEX IF NOT EXISTS k ON t (code);"},
	{SQL: "CREATE INDEX i ON nokey (upper(a));"},
	{SQL: "DROP INDEX ix_n;\nDROP INDEX IF EXISTS nosuch;"},
	{SQL: "ALTER INDEX ix_n RENAME TO ix_m;"},
	{SQL: "ALTER INDEX ix_n REBUILD;"},
	{SQL: "CREATE VIEW v AS SELECT 1 a FROM dual;\nCREATE OR REPLACE FORCE VIEW pv AS SELECT 1 a FROM dual WITH READ ONLY CONSTRAINT pv_ro;"},
	{SQL: "CREATE VIEW v (a, CONSTRAINT v_pk PRIMARY KEY (a) RELY DISABLE NOVALIDATE) AS SELECT 1 a FROM dual;"},
	{SQL: "ALTER VIEW pv COMPILE;\nALTER VIEW pv ADD CONSTRAINT pv_pk PRIMARY KEY (id) DISABLE NOVALIDATE;\nALTER VIEW pv DROP CONSTRAINT pv_pk;"},
	{SQL: "DROP VIEW pv CASCADE CONSTRAINTS;"},
	{SQL: "CREATE SEQUENCE s;\nCREATE SEQUENCE IF NOT EXISTS sq;"},
	{SQL: "ALTER SEQUENCE sq INCREMENT BY 2;\nDROP SEQUENCE sq;\nDROP SEQUENCE IF EXISTS sq;"},
	{SQL: "CREATE SYNONYM s FOR t;\nCREATE OR REPLACE SYNONYM s FOR p;\nDROP SYNONYM s FORCE;"},
	{SQL: "CREATE PUBLIC SYNONYM ps FOR t;\nDROP PUBLIC SYNONYM ps;"},
	{SQL: "CREATE OR REPLACE TYPE ty AS OBJECT (a NUMBER);\n/\nCREATE OR REPLACE TYPE BODY ty AS END;\n/\nDROP TYPE BODY ty;\nDROP TYPE ty FORCE;"},
	{SQL: "CREATE OR REPLACE PROCEDURE p2 AS BEGIN NULL; END;\n/\nCREATE OR REPLACE FUNCTION f2 RETURN NUMBER AS BEGIN RETURN 1; END;\n/\nDROP PROCEDURE p2;\nDROP FUNCTION IF EXISTS f2;"},
	{SQL: "CREATE OR REPLACE PACKAGE pk2 AS END;\n/\nCREATE OR REPLACE PACKAGE BODY pk2 AS END;\n/\nDROP PACKAGE BODY pk2;\nDROP PACKAGE pk2;"},
	{SQL: "CREATE OR REPLACE TRIGGER tr2 BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;\n/\nCREATE TRIGGER tr3 INSTEAD OF INSERT ON pv FOR EACH ROW BEGIN NULL; END;"},
	{SQL: "CREATE TRIGGER tr2 BEFORE INSERT ON nokey FOR EACH ROW BEGIN NULL; END;\n/\nALTER TRIGGER tr2 DISABLE;\nALTER TRIGGER tr2 RENAME TO tr3;\nDROP TRIGGER tr3;\nDROP TRIGGER IF EXISTS tr3;"},
	{SQL: "RENAME pv TO pv2;\nRENAME sq TO sq2;\nRENAME t TO t2;"},
	{SQL: "COMMENT ON TABLE t IS 'x';\nCOMMENT ON COLUMN t.n IS 'x';\nGRANT SELECT ON t TO PUBLIC;\nREVOKE SELECT ON t FROM PUBLIC;\nLOCK TABLE t IN EXCLUSIVE MODE;"},
	{SQL: "ANALYZE TABLE t COMPUTE STATISTICS;\nEXPLAIN PLAN FOR SELECT 1 FROM dual;\nPURGE RECYCLEBIN;\nAUDIT SELECT ON t;\nNOAUDIT SELECT ON t;"},
	{SQL: "COMMIT;\nSAVEPOINT a;\nROLLBACK TO SAVEPOINT a;\nROLLBACK;\nSET TRANSACTION READ ONLY;\nSET CONSTRAINTS ALL DEFERRED;"},
	{SQL: "ALTER SESSION SET CURRENT_SCHEMA = other;\nALTER SESSION SET NLS_DATE_FORMAT = 'YYYY';\nALTER SESSION ENABLE PARALLEL DML;"},
	{SQL: "INSERT INTO nokey VALUES (1);\nUPDATE t SET n = 1;\nDELETE FROM p WHERE id = 1;\nMERGE INTO nokey USING p ON (nokey.a = p.id) WHEN MATCHED THEN UPDATE SET nokey.a = 1;"},
	{SQL: "INSERT ALL INTO nokey (a) VALUES (1) INTO p (id) VALUES (2) SELECT * FROM dual;"},
	{SQL: "SELECT upper(note), count(*) FROM p GROUP BY note;"},
	{SQL: "SELECT av.sales FROM ANALYTIC VIEW (USING sales_av) av;"},
	// A derived analytic view's own columns: its ADD MEASURES names and a
	// base measure's fact column, which the synced relation lacks.
	{SQL: "SELECT av.lag_sales, lag_sales FROM ANALYTIC VIEW (USING sales_av ADD MEASURES(lag_sales AS(1))) av;"},
	{SQL: "WITH av ANALYTIC VIEW AS (USING sales_av ADD MEASURES(lag_sales AS(1))) SELECT lag_sales FROM av;"},
	// The element's columns are its USING relation's, under the element's
	// name or the FROM alias.
	{SQL: "SELECT note FROM p;"},
	{SQL: "WITH av ANALYTIC VIEW AS (USING p) SELECT av.note FROM av;"},
	{SQL: "WITH av ANALYTIC VIEW AS (USING p) SELECT note FROM av x;"},
	// An INSERT's own query reads the element the same way.
	{SQL: "INSERT INTO nokey WITH av ANALYTIC VIEW AS (USING p ADD MEASURES(m AS(1))) SELECT m FROM av;"},
	// A name the statement qualifies is a table, not the element.
	{SQL: "WITH av ANALYTIC VIEW AS (USING p ADD MEASURES(f AS(1))) SELECT f FROM shop.p;", Last: true},
	// Ends the walk after its checks: it may run code of a schema.
	{SQL: "SELECT f(1) FROM dual;", Last: true},
	// An unused element's ADD MEASURES is no column of the query: F may be
	// a schema function.
	{SQL: "WITH av ANALYTIC VIEW AS (USING sales_av ADD MEASURES(f AS(1))) SELECT f FROM p;", Last: true},
	{SQL: "SELECT m FROM ANALYTIC VIEW(USING av ADD MEASURES(m AS(pkg.f())));", Last: true},
	{SQL: "INSERT INTO pv VALUES (1, 'a');", Last: true},
	{SQL: "CREATE TABLE x AS SELECT f(1) a FROM dual;", Last: true},
	{SQL: "CREATE INDEX i ON nokey (f(a));", Last: true},
}

func TestWalkThroughOracleLowersEveryForm(t *testing.T) {
	walktest.Coverage(t, walkForms, func(t *testing.T, sql string) (n, lowered int, last bool) {
		t.Helper()
		_, stmts, syntax := parseReview(sql)
		if syntax != nil {
			t.Fatalf("syntax: %s\n%s", syntax.Message, sql)
		}
		w := newWalker(sql, stmts, review.Target{Schema: walkSchema(false)})
		l := walk.NewLedger(w.snap)
		for i := range stmts {
			s := w.lower(i, l)
			if s == nil {
				return len(stmts), i, false
			}
			if s.Last {
				return len(stmts), i + 1, true
			}
			if s.Apply != nil {
				s.Apply(l)
			}
		}
		return len(stmts), len(stmts), false
	})
}

// TestWalkSweepChangesParse holds the sweep's templates to the parser: a
// change of the sweep with a syntax error would fail the sweep.
func TestWalkSweepChangesParse(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for range 2000 {
		change := walkSweepChange(rng, "SHOP1", "OTHER1", 20)
		if _, _, syntax := parseReview(change); syntax != nil {
			t.Fatalf("syntax: %s\n%s", syntax.Message, change)
		}
	}
}
