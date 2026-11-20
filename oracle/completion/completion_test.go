package completion

import (
	"slices"
	"strings"
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/catalog"
)

type candidateSpec struct {
	text string
	typ  CandidateType
	def  string
}

func TestCompletionGoldenCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Candidate
	}{
		{
			name: "multiple statements limits to current statement",
			input: `select count(1) from t1 where id; SELECT * FROM |
select * from QRTZ_J0B_DETAILS qjd where J0B_NAME like '%gray';`,
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "derived table star expansion",
			input: "SELECT x.| FROM (SELECT * FROM t2) x",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1"},
				{Type: CandidateColumn, Text: "C2"},
			},
		},
		{
			name:  "cte star expansion",
			input: "WITH x AS (SELECT * FROM t2) SELECT x.| FROM x;",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1"},
				{Type: CandidateColumn, Text: "C2"},
			},
		},
		{
			name:  "cte explicit column list",
			input: "WITH x(x1, x2) AS (SELECT * FROM t2) SELECT x.| FROM x;",
			want: []Candidate{
				{Type: CandidateColumn, Text: "X1"},
				{Type: CandidateColumn, Text: "X2"},
			},
		},
		{
			name:  "cte explicit columns unqualified select",
			input: "WITH x(x1, x2) AS (SELECT * FROM t2) SELECT | FROM x",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateTable, Text: "X"},
				{Type: CandidateColumn, Text: "X1"},
				{Type: CandidateColumn, Text: "X2"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "select aliases in group by",
			input: "SELECT c1 as eid, c2 as xid FROM t2 GROUP BY |",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "C2", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "EID"},
				{Type: CandidateColumn, Text: "XID"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "select aliases in order by",
			input: "SELECT c1 as eid, c2 as xid FROM t2 ORDER BY |",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "C2", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "EID"},
				{Type: CandidateColumn, Text: "XID"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name: "select columns after invalid first statement",
			input: `SELECT FROM basdkfjasldf;
SELECT | FROM t1`,
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "derived table alias and columns",
			input: "SELECT | FROM (SELECT c1 FROM t1) cc1",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "CC1"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "join left alias qualified",
			input: "SELECT cc1.| FROM t1 cc1 JOIN t2 on NOT cc1.c1 = t2.c1",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
			},
		},
		{
			name:  "join right table qualified in on",
			input: "SELECT * from t1 cc1 JOIN t2 on cc1.c1 = t2.|",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "C2", Definition: "SCHEMA1.T2 | int, NOT NULL"},
			},
		},
		{
			name:  "function argument qualified alias",
			input: "SELECT MAX(cc1.|) FROM t2 cc1",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "C2", Definition: "SCHEMA1.T2 | int, NOT NULL"},
			},
		},
		{
			name:  "select list physical alias offered as table",
			input: "SELECT | FROM t1 x",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateTable, Text: "X"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "schema qualified table in from",
			input: "SELECT | FROM schema1.t1",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := completeTestSQL(t, tt.input)
			var filtered []Candidate
			for _, c := range result {
				if c.Type == CandidateKeyword || c.Type == CandidateFunction {
					continue
				}
				filtered = append(filtered, c)
			}
			want := dedupAndSort(tt.want)
			if !slices.Equal(filtered, want) {
				t.Fatalf("got:\n%#v\nwant:\n%#v", filtered, want)
			}
		})
	}
}

func TestCompletionLongTailTableSources(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantColumns []string
	}{
		{
			name:        "json table explicit columns",
			input:       "SELECT jt.| FROM t1, JSON_TABLE(t1.c1, '$.rows[*]' COLUMNS (id NUMBER PATH '$.id', name VARCHAR2(100) PATH '$.name')) jt",
			wantColumns: []string{"ID", "NAME"},
		},
		{
			name:        "xmltable explicit columns",
			input:       "SELECT xt.| FROM t1, XMLTABLE('/root/row' PASSING t1.c1 COLUMNS id NUMBER PATH 'id', name VARCHAR2(100) PATH 'name') xt",
			wantColumns: []string{"ID", "NAME"},
		},
		{
			name:        "inline external declared columns",
			input:       "SELECT ext.| FROM EXTERNAL ((a NUMBER, b VARCHAR2(10)) TYPE ORACLE_LOADER DEFAULT DIRECTORY data_dir LOCATION ('x.csv')) ext",
			wantColumns: []string{"A", "B"},
		},
		{
			name:        "lateral inline view columns",
			input:       "SELECT l.| FROM t1, LATERAL (SELECT c1 FROM t2 WHERE t2.c1 = t1.c1) l",
			wantColumns: []string{"C1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := completeTestSQL(t, tt.input)
			var got []string
			for _, c := range result {
				if c.Type == CandidateColumn {
					got = append(got, c.Text)
				}
			}
			slices.Sort(got)
			if !slices.Equal(got, tt.wantColumns) {
				t.Fatalf("columns = %v, want %v", got, tt.wantColumns)
			}
		})
	}
}

func TestCompletionCoverageMatrix(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []candidateSpec
		notWant []candidateSpec
	}{
		{
			name:  "select list uses table columns",
			input: "SELECT | FROM t1",
			want:  cols("C1"),
		},
		{
			name:  "select list uses alias columns",
			input: "SELECT a.| FROM t1 a",
			want:  cols("C1"),
		},
		{
			name:  "select list uses schema qualified table columns",
			input: "SELECT | FROM schema1.t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "column prefix keeps matching column available",
			input: "SELECT C2| FROM t2",
			want:  cols("C2"),
		},
		{
			name:  "qualified column prefix keeps matching column available",
			input: "SELECT t2.C2| FROM t2",
			want:  cols("C2"),
		},
		{
			name:    "unknown qualifier offers no columns",
			input:   "SELECT missing.| FROM t1",
			notWant: cols("C1", "C2"),
		},
		{
			name:  "join on left alias columns",
			input: "SELECT * FROM t1 a JOIN t2 b ON a.|",
			want:  cols("C1"),
		},
		{
			name:  "join on right alias columns",
			input: "SELECT * FROM t1 a JOIN t2 b ON b.|",
			want:  cols("C1", "C2"),
		},
		{
			name:  "join where sees both table columns",
			input: "SELECT * FROM t1 a JOIN t2 b ON a.c1 = b.c1 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:    "left alias excludes right only column",
			input:   "SELECT a.| FROM t1 a JOIN t2 b ON a.c1 = b.c1",
			want:    cols("C1"),
			notWant: cols("C2"),
		},
		{
			name:  "join table reference offers tables",
			input: "SELECT * FROM t1 JOIN |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "using clause offers joined columns",
			input: "SELECT * FROM t1 JOIN t2 USING (|)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte explicit columns are available",
			input: "WITH x(x1, x2) AS (SELECT * FROM t2) SELECT x.| FROM x",
			want:  cols("X1", "X2"),
		},
		{
			name:  "cte star body expands",
			input: "WITH x AS (SELECT * FROM t2) SELECT x.| FROM x",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte table candidate is available without schema qualifier",
			input: "WITH x AS (SELECT * FROM t2) SELECT * FROM |",
			want:  tbls("X", "T1", "T2"),
		},
		{
			name:    "cte table candidate is excluded from schema qualified table reference",
			input:   "WITH x AS (SELECT * FROM t2) SELECT * FROM schema1.|",
			want:    tbls("T1", "T2"),
			notWant: tbls("X"),
		},
		{
			name:  "multiple ctes expose selected cte columns",
			input: "WITH a AS (SELECT c1 FROM t1), b AS (SELECT c1, c2 FROM t2) SELECT b.| FROM a JOIN b ON a.c1 = b.c1",
			want:  cols("C1", "C2"),
		},
		{
			name:  "derived table explicit projection columns",
			input: "SELECT x.| FROM (SELECT c1 FROM t1) x",
			want:  cols("C1"),
		},
		{
			name:  "derived table star body expands",
			input: "SELECT x.| FROM (SELECT * FROM t2) x",
			want:  cols("C1", "C2"),
		},
		{
			name:    "unqualified select from derived table exposes derived columns",
			input:   "SELECT | FROM (SELECT c1 FROM t2) x",
			want:    cols("C1"),
			notWant: cols("C2"),
		},
		{
			name:  "lateral alias columns",
			input: "SELECT l.| FROM t1, LATERAL (SELECT c1 FROM t2 WHERE t2.c1 = t1.c1) l",
			want:  cols("C1"),
		},
		{
			name:    "lateral body sees left outer reference",
			input:   "SELECT * FROM t1, LATERAL (SELECT t1.| FROM t2) l",
			want:    cols("C1"),
			notWant: cols("C2"),
		},
		{
			name:  "containers alias resolves base object columns",
			input: "SELECT c.| FROM CONTAINERS(schema1.t2) c",
			want:  cols("C1", "C2"),
		},
		{
			name:    "table function alias does not invent columns",
			input:   "SELECT x.| FROM TABLE(pkg.values()) x",
			notWant: cols("C1", "C2"),
		},
		{
			name:  "insert target table reference",
			input: "INSERT INTO |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "insert column list uses target table columns",
			input: "INSERT INTO t2 (|) VALUES (1)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "insert select uses source table columns",
			input: "INSERT INTO t1 SELECT | FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "update target table reference",
			input: "UPDATE | SET c1 = 1",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "update set uses target table columns",
			input: "UPDATE t2 SET |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "update where uses target table columns",
			input: "UPDATE t2 SET c1 = 1 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "delete target table reference",
			input: "DELETE FROM |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "delete where uses target table columns",
			input: "DELETE FROM t2 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "merge target table reference",
			input: "MERGE INTO |",
			want:  tbls("T1", "T2"),
		},
		{
			name:    "merge source table reference offers real tables",
			input:   "MERGE INTO t1 target USING |",
			want:    tbls("T1", "T2"),
			notWant: tbls("TARGET"),
		},
		{
			name:  "merge on uses target and source columns",
			input: "MERGE INTO t1 target USING t2 source ON |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "schema qualified table reference offers schema tables",
			input: "SELECT * FROM schema1.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "table prefix keeps matching table available",
			input: "SELECT * FROM t|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "quoted identifier qualifier resolves columns",
			input: `SELECT "T2".| FROM t2`,
			want:  cols("C1", "C2"),
		},
		{
			name:  "select alias is available in order by",
			input: "SELECT c1 AS c_alias FROM t1 ORDER BY |",
			want:  cols("C1", "C_ALIAS"),
		},
		{
			name:    "select alias is not available in where",
			input:   "SELECT c1 AS c_alias FROM t1 WHERE |",
			want:    cols("C1"),
			notWant: cols("C_ALIAS"),
		},
		{
			name:  "correlated subquery can qualify outer table",
			input: "SELECT * FROM t1 outer_t WHERE EXISTS (SELECT * FROM t2 WHERE outer_t.|)",
			want:  cols("C1"),
		},
		{
			name:  "correlated subquery can qualify inner table",
			input: "SELECT * FROM t1 outer_t WHERE EXISTS (SELECT * FROM t2 inner_t WHERE inner_t.|)",
			want:  cols("C1", "C2"),
		},
		{
			name:    "schema qualified alias excludes same table from other schema",
			input:   "SELECT a.| FROM schema1.t1 a JOIN schema4.t1 b ON a.c1 = b.c1",
			want:    cols("C1"),
			notWant: cols("SCHEMA4_ONLY"),
		},
		{
			name:    "pivot does not invent transformed columns",
			input:   "SELECT | FROM t2 PIVOT (COUNT(*) FOR c1 IN (1 AS one))",
			notWant: cols("ONE"),
		},
		{
			name:    "unpivot does not invent transformed columns",
			input:   "SELECT | FROM t2 UNPIVOT (val FOR col IN (c1 AS 'C1', c2 AS 'C2'))",
			notWant: cols("VAL", "COL"),
		},
		{
			name:    "model does not invent transformed columns",
			input:   "SELECT | FROM t2 MODEL DIMENSION BY (c1) MEASURES (c2) RULES (c2[1] = 1)",
			notWant: cols("VAL", "COL"),
		},
		{
			name:  "alter drop column offers table columns",
			input: "ALTER TABLE t2 DROP COLUMN |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "comment on schema qualified column offers columns",
			input: "COMMENT ON COLUMN schema1.t2.|",
			want:  cols("C1", "C2"),
		},
		{
			name:    "drop schema qualified view reference offers only views",
			input:   "DROP VIEW schema1.|",
			want:    vws("V1"),
			notWant: append(tbls("T1", "T2"), seqs("SEQ1", "USER_ID_SEQ")...),
		},
		{
			name:    "sequence schema qualified reference uses non default schema",
			input:   "DROP SEQUENCE schema2.|",
			want:    seqs("SCHEMA2_SEQ"),
			notWant: seqs("SEQ1", "USER_ID_SEQ"),
		},
		{
			name:  "create table datatype offers NUMBER",
			input: "CREATE TABLE t (c |)",
			want:  typs("NUMBER"),
		},
		{
			name:  "create table references offers tables",
			input: "CREATE TABLE child (parent_id NUMBER REFERENCES |)",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "quoted table alias columns",
			input: `SELECT "Alias".| FROM t2 "Alias"`,
			want:  cols("C1", "C2"),
		},
		{
			name:  "quoted cte alias columns",
			input: `WITH "X"(x1, x2) AS (SELECT * FROM t2) SELECT "X".| FROM "X"`,
			want:  cols("X1", "X2"),
		},
		{
			name:    "database link table keeps local table columns unavailable when metadata missing",
			input:   "SELECT remote_t.| FROM t2@remote remote_t",
			notWant: cols("C1", "C2"),
		},
		{
			name:    "unknown schema table reference offers no local tables",
			input:   "SELECT * FROM missing_schema.|",
			notWant: tbls("T1", "T2"),
		},
		{
			name:  "function offered in select list",
			input: "SELECT | FROM t1",
			want:  []candidateSpec{{text: "FN1", typ: CandidateFunction}, {text: "CAST", typ: CandidateFunction}},
		},
		{
			name:  "procedure offered in alter procedure",
			input: "ALTER PROCEDURE |",
			want:  []candidateSpec{{text: "PROC1", typ: CandidateProcedure}},
		},
		{
			name:  "sequence member nextval currval offered after sequence dot",
			input: "SELECT seq1.| FROM dual",
			want:  []candidateSpec{{text: "NEXTVAL", typ: CandidateKeyword}, {text: "CURRVAL", typ: CandidateKeyword}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := completeTestSQL(t, tt.input)
			for _, w := range tt.want {
				if !hasSpec(result, w) {
					t.Errorf("missing %s candidate %q in %v", w.typ, w.text, result)
				}
			}
			for _, nw := range tt.notWant {
				if hasSpec(result, nw) {
					t.Errorf("unexpected %s candidate %q in %v", nw.typ, nw.text, result)
				}
			}
		})
	}
}

func TestCompleteMaterializedViewsAndPackages(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "APP",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "APP",
				Tables: []*metadatapb.TableMetadata{
					{Name: "T1", Columns: []*metadatapb.ColumnMetadata{{Name: "C1", Type: "NUMBER"}}},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{Name: "MV1", Definition: "CREATE MATERIALIZED VIEW mv1 AS SELECT * FROM t1"},
				},
				Packages: []*metadatapb.PackageMetadata{
					{Name: "PKG1"},
				},
			},
		},
	}

	sql, pos := splitCaret(t, "SELECT * FROM |")
	res := Complete(sql, pos, meta)
	if !hasSpec(res, candidateSpec{text: "MV1", typ: CandidateMaterializedView}) {
		t.Fatalf("expected MV1 materialized view candidate, got %v", res)
	}
}

func TestCompleteCodexReviewFixes(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Name: "S1",
		PublicSynonyms: []*metadatapb.SynonymMetadata{
			{Name: "PUB_SYN", TargetSchema: "S1", TargetName: "T1"},
		},
		LinkedDatabases: []*metadatapb.LinkedDatabaseMetadata{
			{Name: "REMOTE_DB"},
		},
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "S1",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:     "T1",
						Columns:  []*metadatapb.ColumnMetadata{{Name: "S1_T1_COL", Type: "NUMBER"}},
						Triggers: []*metadatapb.TriggerMetadata{{Name: "TBL_TRG"}},
						CheckConstraints: []*metadatapb.CheckConstraintMetadata{
							{Name: "CHK_T1_POS", Expression: "S1_T1_COL > 0"},
						},
						ForeignKeys: []*metadatapb.ForeignKeyMetadata{
							{Name: "FK_T1_T2"},
						},
						Indexes: []*metadatapb.IndexMetadata{
							{Name: "PK_T1", Primary: true, IsConstraint: true},
						},
					},
					{
						Name:    "T2",
						Columns: []*metadatapb.ColumnMetadata{{Name: "S1_T2_COL", Type: "NUMBER"}},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name:     "V1",
						Triggers: []*metadatapb.TriggerMetadata{{Name: "VIEW_TRG"}},
					},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{
						Name:       "MV1",
						Definition: "CREATE MATERIALIZED VIEW MV1 AS SELECT S1_T1_COL FROM T1",
						Triggers:   []*metadatapb.TriggerMetadata{{Name: "MV_TRG"}},
						Indexes:    []*metadatapb.IndexMetadata{{Name: "MV_IDX"}},
					},
				},
				ExternalTables: []*metadatapb.ExternalTableMetadata{
					{
						Name:    "EXT1",
						Columns: []*metadatapb.ColumnMetadata{{Name: "EXT_COL", Type: "VARCHAR2(50)"}},
					},
				},
				Synonyms: []*metadatapb.SynonymMetadata{
					{Name: "PRIV_SYN", TargetSchema: "S1", TargetName: "T1"},
				},
				Packages: []*metadatapb.PackageMetadata{
					{
						Name: "PKG1",
						Definition: `CREATE OR REPLACE PACKAGE PKG1 AS
							PROCEDURE DO_WORK(p_id NUMBER);
							FUNCTION GET_COUNT RETURN NUMBER;
						END PKG1;`,
					},
				},
			},
			{
				Name: "S2",
				Tables: []*metadatapb.TableMetadata{
					{
						Name:    "T1",
						Columns: []*metadatapb.ColumnMetadata{{Name: "S2_T1_COL", Type: "NUMBER"}},
					},
					{
						Name:    "T2",
						Columns: []*metadatapb.ColumnMetadata{{Name: "S2_T2_COL", Type: "NUMBER"}},
					},
				},
			},
		},
	}

	t.Run("synonyms in relation completion", func(t *testing.T) {
		sql, pos := splitCaret(t, "SELECT * FROM |")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if !hasSpec(res, candidateSpec{text: "PRIV_SYN", typ: CandidateSynonym}) {
			t.Fatalf("expected PRIV_SYN synonym candidate, got %v", res)
		}
		if !hasSpec(res, candidateSpec{text: "PUB_SYN", typ: CandidateSynonym}) {
			t.Fatalf("expected PUB_SYN public synonym candidate, got %v", res)
		}
	})

	t.Run("database links after at sign", func(t *testing.T) {
		sql, pos := splitCaret(t, "SELECT * FROM t1@|")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if !hasSpec(res, candidateSpec{text: "REMOTE_DB", typ: CandidateDatabaseLink}) {
			t.Fatalf("expected REMOTE_DB link candidate, got %v", res)
		}
	})

	t.Run("triggers from tables views and mviews", func(t *testing.T) {
		sql, pos := splitCaret(t, "DROP TRIGGER |")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		for _, want := range []string{"TBL_TRG", "VIEW_TRG", "MV_TRG"} {
			if !hasSpec(res, candidateSpec{text: want, typ: CandidateTrigger}) {
				t.Fatalf("missing trigger %q in %v", want, res)
			}
		}
	})

	t.Run("schema table pairs do not cartesian product", func(t *testing.T) {
		sql, pos := splitCaret(t, "SELECT | FROM S1.T1, S2.T2")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if !hasSpec(res, candidateSpec{text: "S1_T1_COL", typ: CandidateColumn}) ||
			!hasSpec(res, candidateSpec{text: "S2_T2_COL", typ: CandidateColumn}) {
			t.Fatalf("expected S1_T1_COL and S2_T2_COL in %v", res)
		}
		if hasSpec(res, candidateSpec{text: "S1_T2_COL", typ: CandidateColumn}) ||
			hasSpec(res, candidateSpec{text: "S2_T1_COL", typ: CandidateColumn}) {
			t.Fatalf("unexpected cross-product columns in %v", res)
		}
	})

	t.Run("later cte hidden inside earlier cte body", func(t *testing.T) {
		sql, pos := splitCaret(t, "WITH first AS (SELECT * FROM |), later AS (SELECT * FROM t1) SELECT * FROM later")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if hasSpec(res, candidateSpec{text: "later", typ: CandidateTable}) || hasSpec(res, candidateSpec{text: "LATER", typ: CandidateTable}) {
			t.Fatalf("unexpected later CTE in earlier CTE body: %v", res)
		}
		if hasSpec(res, candidateSpec{text: "first", typ: CandidateTable}) || hasSpec(res, candidateSpec{text: "FIRST", typ: CandidateTable}) {
			t.Fatalf("unexpected self CTE before completion: %v", res)
		}
	})

	t.Run("expression operand not suggested as select alias", func(t *testing.T) {
		sql, pos := splitCaret(t, "SELECT S1_T2_COL + S1_T1_COL FROM t2 ORDER BY |")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if hasSpec(res, candidateSpec{text: "S1_T1_COL", typ: CandidateColumn}) {
			t.Fatalf("unexpected definitionless operand alias S1_T1_COL in %v", res)
		}
	})

	t.Run("synonym target columns resolved", func(t *testing.T) {
		sql, pos := splitCaret(t, "SELECT | FROM PRIV_SYN")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if !hasSpec(res, candidateSpec{text: "S1_T1_COL", typ: CandidateColumn}) {
			t.Fatalf("expected S1_T1_COL via PRIV_SYN, got %v", res)
		}
	})

	t.Run("materialized view indexes offered", func(t *testing.T) {
		sql, pos := splitCaret(t, "DROP INDEX |")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if !hasSpec(res, candidateSpec{text: "MV_IDX", typ: CandidateIndex}) {
			t.Fatalf("expected MV_IDX index candidate, got %v", res)
		}
	})

	t.Run("table constraints offered on drop constraint", func(t *testing.T) {
		sql, pos := splitCaret(t, "ALTER TABLE T1 DROP CONSTRAINT |")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		for _, want := range []string{"CHK_T1_POS", "FK_T1_T2", "PK_T1"} {
			if !hasSpec(res, candidateSpec{text: want, typ: CandidateConstraint}) {
				t.Fatalf("missing constraint %q in %v", want, res)
			}
		}
	})

	t.Run("package members offered after package dot", func(t *testing.T) {
		sql, pos := splitCaret(t, "BEGIN PKG1.|; END;")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if !hasSpec(res, candidateSpec{text: "DO_WORK", typ: CandidateProcedure}) {
			t.Fatalf("expected DO_WORK procedure, got %v", res)
		}
		if !hasSpec(res, candidateSpec{text: "GET_COUNT", typ: CandidateFunction}) {
			t.Fatalf("expected GET_COUNT function, got %v", res)
		}
	})

	t.Run("external tables offered in from clause", func(t *testing.T) {
		sql, pos := splitCaret(t, "SELECT * FROM |")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if !hasSpec(res, candidateSpec{text: "EXT1", typ: CandidateTable}) {
			t.Fatalf("expected EXT1 external table candidate, got %v", res)
		}
	})

	t.Run("materialized view columns resolved from definition", func(t *testing.T) {
		sql, pos := splitCaret(t, "SELECT | FROM MV1")
		res := Complete(sql, pos, meta, WithDefaultSchema("S1"))
		if !hasSpec(res, candidateSpec{text: "S1_T1_COL", typ: CandidateColumn}) {
			t.Fatalf("expected S1_T1_COL from MV1 definition, got %v", res)
		}
	})
}

func TestCompleteWithSingleCatalogAndNoOptions(t *testing.T) {
	cat, err := catalog.LoadSQL(`
		CREATE TABLE employees (
			emp_id NUMBER NOT NULL,
			emp_name VARCHAR2(100)
		);
		CREATE VIEW emp_v AS SELECT emp_id FROM employees;
		CREATE SEQUENCE emp_seq;
	`)
	if err != nil {
		t.Fatalf("LoadSQL failed: %v", err)
	}

	sql, pos := splitCaret(t, "SELECT | FROM employees")
	result := Complete(sql, pos, nil, WithCatalog(cat))
	if !hasSpec(result, candidateSpec{text: "EMP_ID", typ: CandidateColumn}) {
		t.Fatalf("expected EMP_ID column from catalog, got %v", result)
	}
	if !hasSpec(result, candidateSpec{text: "EMP_NAME", typ: CandidateColumn}) {
		t.Fatalf("expected EMP_NAME column from catalog, got %v", result)
	}

	sql2, pos2 := splitCaret(t, "SELECT * FROM |")
	result2 := Complete(sql2, pos2, cat.Metadata())
	if !hasSpec(result2, candidateSpec{text: "EMPLOYEES", typ: CandidateTable}) {
		t.Fatalf("expected EMPLOYEES table, got %v", result2)
	}
	if !hasSpec(result2, candidateSpec{text: "EMP_V", typ: CandidateView}) {
		t.Fatalf("expected EMP_V view, got %v", result2)
	}
	if !hasSpec(result2, candidateSpec{text: "EMP_SEQ", typ: CandidateSequence}) {
		t.Fatalf("expected EMP_SEQ sequence, got %v", result2)
	}
}

func TestCompleteOutOfRangeAndNilMetadata(t *testing.T) {
	cases := []struct {
		name string
		sql  string
		pos  int
	}{
		{name: "empty negative", sql: "", pos: -5},
		{name: "empty past end", sql: "", pos: 100},
		{name: "select past end", sql: "SELECT * FROM ", pos: 999},
		{name: "garbage", sql: "))) ??? ''' \"\"\"", pos: 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_ = Complete(tc.sql, tc.pos, nil)
		})
	}
}

func TestCompleteTrickySelectWithoutSemicolon(t *testing.T) {
	input := "SELECT * FROM broken(\nSELECT | FROM t2"
	result := completeTestSQL(t, input)
	if !hasSpec(result, candidateSpec{text: "C1", typ: CandidateColumn}) || !hasSpec(result, candidateSpec{text: "C2", typ: CandidateColumn}) {
		t.Fatalf("expected C1 and C2 after unseparated SELECT, got %v", result)
	}
}

func TestCompleteInsideQuotedPrefix(t *testing.T) {
	input := `SELECT * FROM "T|`
	result := completeTestSQL(t, input)
	if !hasSpec(result, candidateSpec{text: "T1", typ: CandidateTable}) || !hasSpec(result, candidateSpec{text: "T2", typ: CandidateTable}) {
		t.Fatalf("expected T1 and T2 inside open quote, got %v", result)
	}
}

func TestCompleteShadowedAliasRecursiveCTEAndQuotedCase(t *testing.T) {
	meta := &metadatapb.DatabaseSchemaMetadata{
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "HR",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "ALBUM",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "ALBUM_ID", Type: "NUMBER"},
							{Name: "TITLE", Type: "VARCHAR2(100)"},
						},
					},
					{
						Name: "MUSIC",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "MUSIC_ID", Type: "NUMBER"},
							{Name: "SONG_TITLE", Type: "VARCHAR2(100)"},
						},
					},
					{
						Name: "EMPLOYEES",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "UPPER_ID", Type: "NUMBER"},
						},
					},
				},
			},
			{
				Name: "hr",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "employees",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "lower_id", Type: "NUMBER"},
						},
					},
				},
			},
		},
	}

	t.Run("shadowed alias stops at nearest scope", func(t *testing.T) {
		sql, pos := splitCaret(t, "SELECT * FROM album x WHERE EXISTS (SELECT 1 FROM music x WHERE x.|)")
		res := Complete(sql, pos, meta, WithDefaultSchema("HR"))
		if !hasSpec(res, candidateSpec{text: "MUSIC_ID", typ: CandidateColumn}) || !hasSpec(res, candidateSpec{text: "SONG_TITLE", typ: CandidateColumn}) {
			t.Fatalf("expected inner x (MUSIC) columns, got %v", res)
		}
		if hasSpec(res, candidateSpec{text: "ALBUM_ID", typ: CandidateColumn}) || hasSpec(res, candidateSpec{text: "TITLE", typ: CandidateColumn}) {
			t.Fatalf("did not expect shadowed outer x (ALBUM) columns, got %v", res)
		}
	})

	t.Run("recursive CTE visible in recursive member but hidden in anchor", func(t *testing.T) {
		sql, pos := splitCaret(t, "WITH t(n) AS (SELECT 1 FROM dual UNION ALL SELECT n+1 FROM | WHERE n < 10) SELECT * FROM t")
		res := Complete(sql, pos, meta, WithDefaultSchema("HR"))
		if !hasSpec(res, candidateSpec{text: "T", typ: CandidateTable}) {
			t.Fatalf("expected recursive CTE T in recursive member FROM, got %v", res)
		}

		sqlCols, posCols := splitCaret(t, "WITH t(n) AS (SELECT 1 FROM dual UNION ALL SELECT | FROM t) SELECT * FROM t")
		resCols := Complete(sqlCols, posCols, meta, WithDefaultSchema("HR"))
		if !hasSpec(resCols, candidateSpec{text: "N", typ: CandidateColumn}) {
			t.Fatalf("expected recursive CTE column N in recursive member SELECT, got %v", resCols)
		}

		sqlAnchor, posAnchor := splitCaret(t, "WITH t(n) AS (SELECT 1 FROM | UNION ALL SELECT n+1 FROM t WHERE n < 10) SELECT * FROM t")
		resAnchor := Complete(sqlAnchor, posAnchor, meta, WithDefaultSchema("HR"))
		if hasSpec(resAnchor, candidateSpec{text: "T", typ: CandidateTable}) {
			t.Fatalf("did not expect CTE T inside its own anchor member, got %v", resAnchor)
		}
	})

	t.Run("quoted identifier preserves exact case", func(t *testing.T) {
		sqlMiss, posMiss := splitCaret(t, `SELECT | FROM "employees"`)
		resMiss := Complete(sqlMiss, posMiss, meta, WithDefaultSchema("HR"))
		if hasSpec(resMiss, candidateSpec{text: "UPPER_ID", typ: CandidateColumn}) {
			t.Fatalf("did not expect uppercase EMPLOYEES columns for quoted lowercase table, got %v", resMiss)
		}

		sqlHit, posHit := splitCaret(t, `SELECT | FROM "hr"."employees"`)
		resHit := Complete(sqlHit, posHit, meta, WithDefaultSchema("HR"))
		if !hasSpec(resHit, candidateSpec{text: "lower_id", typ: CandidateColumn}) {
			t.Fatalf("expected lowercase employees column for quoted schema and table, got %v", resHit)
		}
		if hasSpec(resHit, candidateSpec{text: "UPPER_ID", typ: CandidateColumn}) {
			t.Fatalf("did not expect uppercase HR.EMPLOYEES column, got %v", resHit)
		}

		sqlQualMiss, posQualMiss := splitCaret(t, `SELECT "employees".| FROM EMPLOYEES`)
		resQualMiss := Complete(sqlQualMiss, posQualMiss, meta, WithDefaultSchema("HR"))
		if hasSpec(resQualMiss, candidateSpec{text: "UPPER_ID", typ: CandidateColumn}) {
			t.Fatalf("did not expect uppercase EMPLOYEES column for lowercase quoted qualifier, got %v", resQualMiss)
		}

		sqlSchema, posSchema := splitCaret(t, `SELECT * FROM "hr".|`)
		resSchema := Complete(sqlSchema, posSchema, meta, WithDefaultSchema("HR"))
		if !hasSpec(resSchema, candidateSpec{text: "employees", typ: CandidateTable}) {
			t.Fatalf("expected lowercase hr.employees table, got %v", resSchema)
		}
		if hasSpec(resSchema, candidateSpec{text: "EMPLOYEES", typ: CandidateTable}) {
			t.Fatalf("did not expect uppercase HR.EMPLOYEES table for quoted hr., got %v", resSchema)
		}
	})

	t.Run("materialized view DDL offers materialized views", func(t *testing.T) {
		mvMeta := &metadatapb.DatabaseSchemaMetadata{
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "S1",
					Views: []*metadatapb.ViewMetadata{
						{Name: "V1"},
					},
					MaterializedViews: []*metadatapb.MaterializedViewMetadata{
						{Name: "MV1"},
					},
				},
			},
		}
		for _, query := range []string{"DROP MATERIALIZED VIEW |", "ALTER MATERIALIZED VIEW S1.|"} {
			sqlMV, posMV := splitCaret(t, query)
			resMV := Complete(sqlMV, posMV, mvMeta, WithDefaultSchema("S1"))
			if !hasSpec(resMV, candidateSpec{text: "MV1", typ: CandidateMaterializedView}) {
				t.Fatalf("%q: expected MV1 materialized view candidate, got %v", query, resMV)
			}
			if hasSpec(resMV, candidateSpec{text: "V1", typ: CandidateView}) {
				t.Fatalf("%q: did not expect regular view V1, got %v", query, resMV)
			}
		}
	})
}

func completeTestSQL(t *testing.T, input string) []Candidate {
	t.Helper()
	sql, pos := splitCaret(t, input)
	return Complete(
		sql,
		pos,
		nil,
		WithDefaultSchema("SCHEMA1"),
		WithSchemaNames([]string{"SCHEMA1", "SCHEMA2", "SCHEMA3"}),
		WithSchemaLoader(loadTestSchema),
	)
}

func splitCaret(t *testing.T, input string) (string, int) {
	t.Helper()
	idx := strings.Index(input, "|")
	if idx < 0 {
		t.Fatalf("missing '|' caret in %q", input)
	}
	return input[:idx] + input[idx+1:], idx
}

func hasSpec(candidates []Candidate, want candidateSpec) bool {
	for _, c := range candidates {
		if c.Type == want.typ && c.Text == want.text {
			if want.def == "" || c.Definition == want.def {
				return true
			}
		}
	}
	return false
}

func cols(texts ...string) []candidateSpec { return specs(CandidateColumn, texts...) }
func tbls(texts ...string) []candidateSpec { return specs(CandidateTable, texts...) }
func vws(texts ...string) []candidateSpec  { return specs(CandidateView, texts...) }
func seqs(texts ...string) []candidateSpec { return specs(CandidateSequence, texts...) }
func typs(texts ...string) []candidateSpec { return specs(CandidateType_, texts...) }

func specs(typ CandidateType, texts ...string) []candidateSpec {
	out := make([]candidateSpec, 0, len(texts))
	for _, t := range texts {
		out = append(out, candidateSpec{text: t, typ: typ})
	}
	return out
}

func loadTestSchema(name string) *metadatapb.DatabaseSchemaMetadata {
	switch name {
	case "SCHEMA1":
		return &metadatapb.DatabaseSchemaMetadata{
			Name: "SCHEMA1",
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "",
					Tables: []*metadatapb.TableMetadata{
						{
							Name: "T1",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "C1", Type: "int"},
							},
						},
						{
							Name: "T2",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "C1", Type: "int"},
								{Name: "C2", Type: "int"},
							},
						},
					},
					Views: []*metadatapb.ViewMetadata{
						{
							Name:       "V1",
							Definition: "CREATE VIEW v1 AS SELECT * FROM t1",
						},
					},
					Sequences: []*metadatapb.SequenceMetadata{
						{Name: "SEQ1"},
						{Name: "USER_ID_SEQ"},
					},
					Functions: []*metadatapb.FunctionMetadata{
						{Name: "FN1"},
					},
					Procedures: []*metadatapb.ProcedureMetadata{
						{Name: "PROC1"},
					},
					Packages: []*metadatapb.PackageMetadata{
						{Name: "PKG1"},
					},
				},
			},
		}
	case "SCHEMA2":
		return &metadatapb.DatabaseSchemaMetadata{
			Name: "SCHEMA2",
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "",
					Tables: []*metadatapb.TableMetadata{
						{
							Name:    "T1",
							Columns: []*metadatapb.ColumnMetadata{{Name: "C1", Type: "int"}},
						},
						{
							Name: "T2",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "C1", Type: "int"},
								{Name: "C2", Type: "int"},
							},
						},
					},
					Views:     []*metadatapb.ViewMetadata{{Name: "V1", Definition: "CREATE VIEW v1 AS SELECT * FROM t1"}},
					Sequences: []*metadatapb.SequenceMetadata{{Name: "SCHEMA2_SEQ"}},
				},
			},
		}
	case "SCHEMA3":
		return &metadatapb.DatabaseSchemaMetadata{
			Name: "SCHEMA3",
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "",
					Tables: []*metadatapb.TableMetadata{
						{
							Name:    "T1",
							Columns: []*metadatapb.ColumnMetadata{{Name: "C1", Type: "int"}},
						},
						{
							Name: "T2",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "C1", Type: "int"},
								{Name: "C2", Type: "int"},
							},
						},
					},
					Views: []*metadatapb.ViewMetadata{{Name: "V1", Definition: "CREATE VIEW v1 AS SELECT * FROM t1"}},
				},
			},
		}
	case "SCHEMA4":
		return &metadatapb.DatabaseSchemaMetadata{
			Name: "SCHEMA4",
			Schemas: []*metadatapb.SchemaMetadata{
				{
					Name: "",
					Tables: []*metadatapb.TableMetadata{
						{
							Name: "T1",
							Columns: []*metadatapb.ColumnMetadata{
								{Name: "C1", Type: "int"},
								{Name: "SCHEMA4_ONLY", Type: "int"},
							},
						},
					},
				},
			},
		}
	default:
		return nil
	}
}
