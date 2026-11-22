package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

// priorBackupFinding is a prior-backup finding as the tests compare it: every
// field a finding can set, the range as the text it anchors.
type priorBackupFinding struct {
	Code      v1review.Code
	Name      string
	Schema    string
	Table     string
	Value     string
	Statement int
	Text      string
}

// priorBackupFacts is a prior backup that can take 100 bytes and finds its
// database.
func priorBackupFacts(size int) v1review.PriorBackupFacts {
	return v1review.PriorBackupFacts{TotalSize: size, MaxSize: 100, Database: "BBDATAARCHIVE", DatabaseExists: true}
}

// priorBackupSchema is the pre-change snapshot of database DB: one CLOB table
// and two LONG tables, held under the empty schema name the Oracle sync stores
// the connection schema's tables under.
func priorBackupSchema() v1review.Schema {
	return v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{
		Name: "DB",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{Name: "T", Columns: []*metadata.ColumnMetadata{{Name: "ID", Type: "NUMBER"}}},
				{Name: "T_LONG", Columns: []*metadata.ColumnMetadata{
					{Name: "ID", Type: "NUMBER"},
					{Name: "PAYLOAD", Type: "LONG"},
				}},
				{Name: "T_LONG_RAW", Columns: []*metadata.ColumnMetadata{
					{Name: "ID", Type: "NUMBER"},
					{Name: "BLOB_ISH", Type: "LONG RAW"},
				}},
				{Name: "T_CLOB", Columns: []*metadata.ColumnMetadata{
					{Name: "ID", Type: "NUMBER"},
					{Name: "DOC", Type: "CLOB"},
				}},
			},
		}},
	}}
}

// judgePriorBackup parses sql the way Bytebase does and runs PriorBackupCheck
// over its statements.
func judgePriorBackup(t *testing.T, sql string, schema v1review.Schema, facts v1review.PriorBackupFacts) []priorBackupFinding {
	t.Helper()
	var findings []priorBackupFinding
	for _, f := range PriorBackupCheck(parseStmts(t, sql), schema, facts) {
		findings = append(findings, priorBackupFinding{f.Code, f.Name, f.Schema, f.Table, f.Value, f.Statement, sql[f.Start:f.End]})
	}
	return findings
}

func TestPriorBackupCheck(t *testing.T) {
	tooLarge := func() priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupScriptTooLarge, Value: "100"}
	}
	ddl := func() priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupDDLInBatch}
	}
	missingDatabase := func() priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupDatabaseMissing, Name: "BBDATAARCHIVE"}
	}
	mixedDML := func(table string) priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupMixedDMLOnTable, Table: table}
	}
	longColumn := func(schema, table string, columns string) priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupUncopyableColumn, Schema: schema, Table: table, Value: columns}
	}

	ddlCases := []struct {
		name string
		ddl  string
	}{
		{name: "create table", ddl: "CREATE TABLE U (A NUMBER)"},
		{name: "alter table", ddl: "ALTER TABLE T ADD B NUMBER"},
		{name: "drop table", ddl: "DROP TABLE U"},
		{name: "create index", ddl: "CREATE INDEX I ON T (ID)"},
		{name: "create view", ddl: "CREATE VIEW V AS SELECT 1 FROM DUAL"},
		{name: "truncate", ddl: "TRUNCATE TABLE T"},
		{name: "comment", ddl: "COMMENT ON TABLE T IS 'x'"},
	}
	for _, tc := range ddlCases {
		t.Run(tc.name, func(t *testing.T) {
			got := judgePriorBackup(t, tc.ddl+";\nDELETE FROM T WHERE ID = 1;", priorBackupSchema(), priorBackupFacts(10))
			if diff := cmp.Diff([]priorBackupFinding{ddl()}, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}

	for _, tt := range []struct {
		name   string
		sql    string
		schema v1review.Schema
		facts  v1review.PriorBackupFacts
		want   []priorBackupFinding
	}{
		{name: "dml only", sql: "UPDATE T SET ID = 1 WHERE ID = 2;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{name: "nothing", sql: "SELECT 1 FROM DUAL;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{name: "size over the limit", sql: "UPDATE T SET ID = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(101), want: []priorBackupFinding{tooLarge()}},
		{name: "size at the limit", sql: "UPDATE T SET ID = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(100)},
		{
			name: "database missing ends the check",
			sql:  "CREATE TABLE U (A NUMBER);\nUPDATE T_LONG SET ID = 1 WHERE ID = 2;\nDELETE FROM T WHERE ID = 3;",
			// The size finding stays: V1 reported it before the lookup.
			schema: priorBackupSchema(),
			facts: v1review.PriorBackupFacts{
				TotalSize: 101, MaxSize: 100, Database: "BBDATAARCHIVE",
			},
			want: []priorBackupFinding{tooLarge(), missingDatabase()},
		},
		{
			name:   "ddl is reported once for the batch",
			sql:    "CREATE TABLE U (A NUMBER);\nALTER TABLE U ADD B NUMBER;\nDELETE FROM T WHERE ID = 1;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{ddl()},
		},
		{
			// V1 reported the DDL whether or not the batch changes rows.
			name:   "ddl alone is reported as mixed",
			sql:    "CREATE TABLE U (A NUMBER);\nDROP TABLE U;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{ddl()},
		},
		{name: "create user is not ddl", sql: "CREATE USER U IDENTIFIED BY P;\nDELETE FROM T WHERE ID = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{name: "grant is not ddl", sql: "GRANT SELECT ON T TO U;\nDELETE FROM T WHERE ID = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{name: "create sequence is not ddl", sql: "CREATE SEQUENCE S START WITH 1;\nDELETE FROM T WHERE ID = 1;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{
			name:   "mixed dml on one table",
			sql:    "UPDATE T SET ID = 1 WHERE ID = 2;\nDELETE FROM T WHERE ID = 3;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("DB.T")},
		},
		{name: "one kind of dml twice", sql: "UPDATE T SET ID = 1 WHERE ID = 2;\nUPDATE T SET ID = 3 WHERE ID = 4;", schema: priorBackupSchema(), facts: priorBackupFacts(10)},
		{
			name:   "an insert is not grouped",
			sql:    "INSERT INTO T (ID) VALUES (1);\nUPDATE T SET ID = 2 WHERE ID = 3;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "a qualified and an unqualified reference are one table",
			sql:    "UPDATE DB.T SET ID = 1 WHERE ID = 2;\nDELETE FROM T WHERE ID = 3;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("DB.T")},
		},
		{
			name:   "two owners are two tables",
			sql:    "UPDATE DB.T SET ID = 1 WHERE ID = 2;\nDELETE FROM OTHER.T WHERE ID = 3;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:  "without a database name",
			sql:   "UPDATE T SET ID = 1 WHERE ID = 2;\nDELETE FROM T WHERE ID = 3;",
			facts: priorBackupFacts(10),
			want:  []priorBackupFinding{mixedDML(".T")},
		},
		{
			name:   "mixed tables in statement order",
			sql:    "UPDATE T SET ID = 1 WHERE ID = 2;\nDELETE FROM U WHERE ID = 3;\nUPDATE U SET ID = 4 WHERE ID = 5;\nDELETE FROM T WHERE ID = 6;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{mixedDML("DB.T"), mixedDML("DB.U")},
		},
		{
			name:   "long column warns",
			sql:    "UPDATE T_LONG SET ID = 1 WHERE ID = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{longColumn("DB", "T_LONG", `"PAYLOAD"`)},
		},
		{
			name:   "long raw column warns",
			sql:    "DELETE FROM T_LONG_RAW WHERE ID = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{longColumn("DB", "T_LONG_RAW", `"BLOB_ISH"`)},
		},
		{
			name:   "a clob column does not warn",
			sql:    "UPDATE T_CLOB SET ID = 1 WHERE ID = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "one warning per table, not per statement",
			sql:    "UPDATE T_LONG SET ID = 1 WHERE ID = 2;\nUPDATE T_LONG SET ID = 3 WHERE ID = 4;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{longColumn("DB", "T_LONG", `"PAYLOAD"`)},
		},
		{
			name:   "a qualified dml still finds the current schema's tables",
			sql:    "UPDATE DB.T_LONG SET ID = 1 WHERE ID = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want:   []priorBackupFinding{longColumn("DB", "T_LONG", `"PAYLOAD"`)},
		},
		{
			name:   "another owner stays silent",
			sql:    "UPDATE OTHER_OWNER.T_LONG SET ID = 1 WHERE ID = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "a table the snapshot does not hold stays silent",
			sql:    "UPDATE NO_SUCH_TABLE SET ID = 1 WHERE ID = 2;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
		},
		{
			name:   "every finding in order",
			sql:    "CREATE TABLE U (A NUMBER);\nUPDATE T_LONG SET ID = 1 WHERE ID = 2;\nDELETE FROM T_LONG WHERE ID = 3;",
			schema: priorBackupSchema(),
			facts:  priorBackupFacts(10),
			want: []priorBackupFinding{
				ddl(),
				mixedDML("DB.T_LONG"),
				longColumn("DB", "T_LONG", `"PAYLOAD"`),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgePriorBackup(t, tt.sql, tt.schema, tt.facts)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

// TestPriorBackupCheckIdentifierCase pins the name comparison the LONG lookup
// makes: in case-sensitive mode an unquoted reference (catalogued upper-case)
// and a quoted one (exact case) are different tables, so a fold-blind compare
// would borrow the wrong table's columns.
func TestPriorBackupCheckIdentifierCase(t *testing.T) {
	db := &metadata.DatabaseSchemaMetadata{
		Name: "DB",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{Name: "T_LONG", Columns: []*metadata.ColumnMetadata{
					{Name: "ID", Type: "NUMBER"},
					{Name: "PAYLOAD", Type: "LONG"},
				}},
				{Name: "t_long", Columns: []*metadata.ColumnMetadata{
					{Name: "ID", Type: "NUMBER"},
					{Name: "NOTE", Type: "VARCHAR2"},
				}},
				{Name: "t_raw", Columns: []*metadata.ColumnMetadata{
					{Name: "ID", Type: "NUMBER"},
					{Name: "RAWCOL", Type: "LONG RAW"},
				}},
			},
		}},
	}
	sensitive := v1review.Schema{Metadata: db, ObjectCaseSensitive: true}
	insensitive := v1review.Schema{Metadata: db}

	for _, tt := range []struct {
		name   string
		sql    string
		schema v1review.Schema
		want   []priorBackupFinding
	}{
		{
			name:   "a quoted reference keeps its case",
			sql:    `UPDATE "t_long" SET ID = 1 WHERE ID = 2;`,
			schema: sensitive,
		},
		{
			name:   "a quoted table with a long column warns",
			sql:    `UPDATE "t_raw" SET ID = 1 WHERE ID = 2;`,
			schema: sensitive,
			want:   []priorBackupFinding{{Code: v1review.PriorBackupUncopyableColumn, Schema: "DB", Table: "t_raw", Value: `"RAWCOL"`}},
		},
		{
			name:   "an unquoted uppercase reference finds the uppercase table",
			sql:    "UPDATE T_LONG SET ID = 1 WHERE ID = 2;",
			schema: sensitive,
			want:   []priorBackupFinding{{Code: v1review.PriorBackupUncopyableColumn, Schema: "DB", Table: "T_LONG", Value: `"PAYLOAD"`}},
		},
		{
			name:   "an unquoted lowercase reference folds upper",
			sql:    "UPDATE t_long SET ID = 1 WHERE ID = 2;",
			schema: sensitive,
			want:   []priorBackupFinding{{Code: v1review.PriorBackupUncopyableColumn, Schema: "DB", Table: "T_LONG", Value: `"PAYLOAD"`}},
		},
		{
			name:   "a case-insensitive instance folds",
			sql:    `UPDATE "t_long" SET ID = 1 WHERE ID = 2;`,
			schema: insensitive,
			want:   []priorBackupFinding{{Code: v1review.PriorBackupUncopyableColumn, Schema: "DB", Table: "t_long", Value: `"PAYLOAD"`}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgePriorBackup(t, tt.sql, tt.schema, priorBackupFacts(10))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

// TestPriorBackupCheckWithoutSnapshot pins that a rule given no snapshot
// reports no LONG column rather than panicking.
func TestPriorBackupCheckWithoutSnapshot(t *testing.T) {
	if got := judgePriorBackup(t, "UPDATE T_LONG SET ID = 1 WHERE ID = 2;", v1review.Schema{}, priorBackupFacts(10)); len(got) != 0 {
		t.Errorf("findings = %+v, want none", got)
	}
}
