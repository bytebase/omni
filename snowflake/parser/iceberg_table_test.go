package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// CREATE / DROP / UNDROP ICEBERG TABLE
// ---------------------------------------------------------------------------

func TestCreateIcebergTable_SnowflakeCatalog(t *testing.T) {
	sql := `CREATE OR REPLACE TRANSIENT ICEBERG TABLE db.s.t (id INT NOT NULL, ts TIMESTAMP, CONSTRAINT pk PRIMARY KEY (id))
  PARTITION BY (id, BUCKET(16, id), TRUNCATE(4, id), YEAR(ts), MONTH(ts), DAY(ts), HOUR(ts))
  PATH_LAYOUT = HIERARCHICAL
  CLUSTER BY (id)
  EXTERNAL_VOLUME = 'vol' CATALOG = 'SNOWFLAKE' BASE_LOCATION = 'loc/'
  TARGET_FILE_SIZE = '64MB'
  COPY GRANTS COPY TAGS COPY = (COMMENTS, TAGS)
  COMMENT = 'c'
  WITH ROW ACCESS POLICY rap ON (id)
  WITH AGGREGATION POLICY ap
  WITH TAG (t1 = 'v')
  WITH CONTACT (steward = alice)`
	stmt, errs := testParseCreateTable(sql)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Iceberg || !stmt.OrReplace || !stmt.Transient {
		t.Errorf("flags: Iceberg=%v OrReplace=%v Transient=%v", stmt.Iceberg, stmt.OrReplace, stmt.Transient)
	}
	if len(stmt.Columns) != 2 || len(stmt.Constraints) != 1 {
		t.Fatalf("columns=%d constraints=%d", len(stmt.Columns), len(stmt.Constraints))
	}
	if len(stmt.PartitionBy) != 7 {
		t.Fatalf("PartitionBy = %d items, want 7", len(stmt.PartitionBy))
	}
	if _, ok := stmt.PartitionBy[0].(*ast.ColumnRef); !ok {
		t.Errorf("PartitionBy[0] = %T, want *ast.ColumnRef", stmt.PartitionBy[0])
	}
	bucket, ok := stmt.PartitionBy[1].(*ast.FuncCallExpr)
	if !ok || bucket.Name.Name.Name != "BUCKET" || len(bucket.Args) != 2 {
		t.Errorf("PartitionBy[1] = %#v, want BUCKET(16, id)", stmt.PartitionBy[1])
	}
	var names []string
	for _, o := range stmt.Options {
		names = append(names, o.Name)
	}
	want := []string{"PATH_LAYOUT", "EXTERNAL_VOLUME", "CATALOG", "BASE_LOCATION", "TARGET_FILE_SIZE", "COPY"}
	if len(names) != len(want) {
		t.Fatalf("options = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("options = %v, want %v", names, want)
		}
	}
	if !stmt.CopyGrants || !stmt.CopyTags {
		t.Errorf("CopyGrants=%v CopyTags=%v", stmt.CopyGrants, stmt.CopyTags)
	}
	if stmt.Comment == nil || *stmt.Comment != "c" {
		t.Errorf("Comment = %v", stmt.Comment)
	}
	if stmt.RowPolicy == nil || len(stmt.RowPolicy.Columns) != 1 {
		t.Errorf("RowPolicy = %+v", stmt.RowPolicy)
	}
	if stmt.AggregationPolicy == nil || stmt.AggregationPolicy.Name.Name != "ap" {
		t.Errorf("AggregationPolicy = %+v", stmt.AggregationPolicy)
	}
	if len(stmt.Tags) != 1 || len(stmt.Contacts) != 1 || stmt.Contacts[0].Name != "STEWARD" {
		t.Errorf("Tags=%d Contacts=%+v", len(stmt.Tags), stmt.Contacts)
	}
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(sql) {
		t.Errorf("Loc = %+v, want {0 %d}", stmt.Loc, len(sql))
	}
}

func TestCreateIcebergTable_ExternalCatalogs(t *testing.T) {
	for _, sql := range []string{
		// AWS Glue / REST catalog
		`CREATE ICEBERG TABLE IF NOT EXISTS t EXTERNAL_VOLUME = 'v' CATALOG = 'glue' CATALOG_TABLE_NAME = 'n' CATALOG_NAMESPACE = 'ns' AUTO_REFRESH = TRUE`,
		// Iceberg files
		`CREATE ICEBERG TABLE t CATALOG = 'objcat' METADATA_FILE_PATH = 'p/v1.metadata.json' REPLACE_INVALID_CHARACTERS = FALSE`,
		// Delta files
		`CREATE ICEBERG TABLE t BASE_LOCATION = 'delta/' WITH TAG (t1 = 'v') WITH CONTACT (support = bob, steward = alice)`,
		// Parquet files
		`CREATE ICEBERG TABLE t (a INT, b INT AS (a + 1)) PARTITION BY (a) BASE_LOCATION = 'p/' ENABLE_SCHEMA_INFERENCE = TRUE PATTERN = '.*'`,
		// REST: PURGE before a CTAS
		`CREATE OR REPLACE ICEBERG TABLE t (a, b INT) PURGE AS SELECT 1, 2`,
		// CTAS with properties before AS
		`CREATE ICEBERG TABLE t CLUSTER BY (a) CATALOG = 'SNOWFLAKE' ICEBERG_VERSION = 3 AS SELECT 1 AS a`,
		// LIKE and CLONE
		`CREATE TRANSIENT ICEBERG TABLE t LIKE s CLUSTER BY (a) COPY GRANTS`,
		`CREATE ICEBERG TABLE t CLONE s AT (OFFSET => -60) COPY GRANTS`,
		// a table named iceberg is still a plain table
		`CREATE TABLE iceberg (a INT)`,
		// ROW ACCESS / AGGREGATION POLICY without WITH
		`CREATE ICEBERG TABLE t (a INT) ROW ACCESS POLICY r ON (a) AGGREGATION POLICY p`,
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
}

func TestCreateIcebergTable_Rejects(t *testing.T) {
	for _, sql := range []string{
		`CREATE ICEBERG TABLE t PARTITION BY (SQRT(a))`,
		`CREATE ICEBERG TABLE t PARTITION BY (BUCKET(a))`,
		`CREATE ICEBERG TABLE t PARTITION BY (YEAR(1, a))`,
		`CREATE ICEBERG TABLE t PARTITION BY ()`,
		`CREATE ICEBERG TABLE t COPY = (ROLES)`,
		`CREATE ICEBERG TABLE t COMMENT 'c'`,
		`CREATE ICEBERG TABLE t BASE_LOCATION = 's' zzz`,
		`CREATE ICEBERG TABLE t LIKE s AS SELECT 1`,
		`CREATE ICEBERG TABLE t ROW ACCESS POLICY r ON ()`,
		`CREATE ICEBERG TABLE t ROW ACCESS POLICY r ON (a,)`,
		`CREATE ICEBERG TABLE t WITH CONTACT ()`,
		`CREATE OR ALTER ICEBERG TABLE t (a INT)`,
		`CREATE TEMPORARY ICEBERG TABLE t (a INT)`,
		`CREATE ICEBERG t (a INT)`,
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want an error", sql)
		}
	}
}

func TestStartsAggregationPolicy(t *testing.T) {
	for sql, want := range map[string]bool{
		"AGGREGATION POLICY p": true,
		"AGGREGATION = 1":      false,
		"POLICY p":             false,
	} {
		p := &Parser{lexer: NewLexer(sql), input: sql}
		p.advance()
		if got := p.startsAggregationPolicy(); got != want {
			t.Errorf("startsAggregationPolicy(%q) = %v, want %v", sql, got, want)
		}
	}
}

func TestDropIcebergTable(t *testing.T) {
	tests := []struct {
		sql              string
		iceberg, purge   bool
		noPurge, cascade bool
	}{
		{sql: "DROP ICEBERG TABLE t", iceberg: true},
		{sql: "DROP ICEBERG TABLE IF EXISTS db.s.t CASCADE PURGE", iceberg: true, purge: true, cascade: true},
		{sql: "DROP TABLE t NOPURGE", noPurge: true},
		{sql: "DROP TABLE t PURGE", purge: true},
	}
	for _, tt := range tests {
		node, err := Parse(tt.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tt.sql, err)
		}
		stmt := node.Stmts[0].(*ast.DropStmt)
		if stmt.Kind != ast.DropTable || stmt.Iceberg != tt.iceberg || stmt.Purge != tt.purge ||
			stmt.NoPurge != tt.noPurge || stmt.Cascade != tt.cascade {
			t.Errorf("%q: got %+v", tt.sql, stmt)
		}
		if stmt.Loc.End != len(tt.sql) {
			t.Errorf("%q: Loc.End = %d, want %d", tt.sql, stmt.Loc.End, len(tt.sql))
		}
	}
	for _, sql := range []string{"DROP ICEBERG TABLE t PURGE NOPURGE", "DROP ICEBERG t"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want an error", sql)
		}
	}
}

func TestUndropIcebergTable(t *testing.T) {
	node, err := Parse("UNDROP ICEBERG TABLE db.s.t")
	if err != nil {
		t.Fatal(err)
	}
	stmt := node.Stmts[0].(*ast.UndropStmt)
	if stmt.Kind != ast.UndropTable || !stmt.Iceberg || stmt.Name.Name.Name != "t" {
		t.Errorf("got %+v", stmt)
	}
}
