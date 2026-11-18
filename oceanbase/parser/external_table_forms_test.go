package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// This file holds the external-table forms to the grammar testdata/compat.tsv
// records: the tbl_name alternative that names a location (url_expr), the
// SOURCE(...) and FILES(...) table functions with their property lists, the
// statement that reads them (CREATE EXTERNAL TABLE, DELETE FROM, LOAD DATA ...
// FROM url_spec), and ALTER EXTERNAL TABLE's two own actions.

// TestExternalTableForms pins which external-table statements parse.
func TestExternalTableForms(t *testing.T) {
	for _, sql := range []string{
		// CREATE EXTERNAL TABLE: special_table_type's EXTERNAL alternative, whose
		// external options are table_option's.
		"CREATE EXTERNAL TABLE a (b INT) LOCATION 's' FORMAT(TYPE = 'csv')",
		"CREATE EXTERNAL TABLE a (b INT) LOCATION 's' FORMAT(TYPE = 'csv') PATTERN 'p'",
		"CREATE EXTERNAL TABLE a (b INT) PROPERTIES(TYPE = 's')",
		"CREATE EXTERNAL TABLE a SELECT a",
		// ALTER EXTERNAL TABLE: ALTER TABLE's actions, and its own two.
		"ALTER EXTERNAL TABLE a",
		"ALTER EXTERNAL TABLE a ADD COLUMN b INT",
		"ALTER EXTERNAL TABLE a ADD PARTITION() LOCATION 's'",
		"ALTER EXTERNAL TABLE a ADD PARTITION(ds = '2024') LOCATION 's'",
		"ALTER EXTERNAL TABLE a ADD PARTITION(ds '2024', ds '2025') LOCATION 's'",
		"ALTER EXTERNAL TABLE a ADD PARTITION(, ds = '2024') LOCATION 's'",
		"ALTER EXTERNAL TABLE a DROP PARTITION LOCATION 's'",
		// DROP PARTITION of a partition named LOCATION is the generic action.
		"ALTER EXTERNAL TABLE a DROP PARTITION LOCATION",
		// url_expr as tbl_name: SOURCE(...), FILES(...) and a quoted path, in
		// DELETE and in FROM.
		"DELETE FROM SOURCE(TYPE = 's')",
		"DELETE FROM SOURCE(TYPE = 's') AS a",
		"DELETE FROM SOURCE(TYPE = 's' TYPE = 's')",
		"DELETE FROM SOURCE(COMPRESSION_CODE = a)",
		"DELETE FROM SOURCE(COMPRESSION = 'zstd', COLLECT_STATISTICS_ON_CREATE = TRUE)",
		"DELETE FROM 's'(FORMAT(TYPE = 's'))",
		"DELETE FROM 's'(FORMAT(TYPE = 's') PATTERN 'p')",
		"DELETE FROM FILES(LOCATION 's', FORMAT(TYPE = 's'), PATTERN 'p')",
		"DELETE FROM FILES(LOCATION 's' FORMAT(TYPE = 's'))",
		"SELECT a FROM SOURCE(TYPE = 's')",
		"SELECT a FROM 's'(FORMAT(TYPE = 's'))",
		"SELECT a FROM FILES(LOCATION 's', FORMAT(TYPE = 's'))",
		// LOAD DATA ... FROM url_spec, with the query form and opt_on_error.
		"LOAD DATA FROM SOURCE(TYPE = 's') INTO TABLE a",
		"LOAD DATA FROM SOURCE(TYPE = 's') INTO TABLE a LOG ERRORS",
		"LOAD DATA FROM SOURCE(TYPE = 's') INTO TABLE a LOG ERRORS REJECT LIMIT 1",
		"LOAD DATA FROM SOURCE(TYPE = 's') INTO TABLE a LOG ERRORS REJECT LIMIT UNLIMITED",
		"LOAD DATA FROM FILES(LOCATION 's', FORMAT(TYPE = 's')) INTO TABLE a",
		"LOAD DATA FROM(SELECT a FROM SOURCE(TYPE = 's')) INTO TABLE a",
		"LOAD DATA REPLACE FROM SOURCE(TYPE = 's') INTO TABLE a PARTITION (p0) (b) LOG ERRORS",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		// The alias of opt_as_alias is written with AS.
		"DELETE FROM SOURCE(TYPE = 's') a",
		// external_properties' COMP_EQ is required, and its value's shape follows
		// its word.
		"DELETE FROM SOURCE(TYPE 's')",
		"DELETE FROM SOURCE(TYPE = a)",
		"DELETE FROM SOURCE(TYPE =)",
		"DELETE FROM SOURCE(ZZZ = 's')",
		"DELETE FROM SOURCE()",
		"DELETE FROM SOURCE(TYPE = 's'",
		// format_expr and location_expr are required in their places.
		"SELECT a FROM 's'()",
		"SELECT a FROM FILES(LOCATION 's')",
		"SELECT a FROM FILES('s' FORMAT(TYPE = 's'))",
		"DELETE FROM 's'(FORMAT(TYPE 's'))",
		// ALTER EXTERNAL TABLE's partition actions need both halves and LOCATION.
		"ALTER EXTERNAL TABLE a ADD PARTITION(a) LOCATION 's'",
		"ALTER EXTERNAL TABLE a ADD PARTITION('s') LOCATION 's'",
		"ALTER EXTERNAL TABLE a ADD PARTITION(a 's')",
		"ALTER EXTERNAL TABLE a DROP PARTITION 's'",
		// url_spec takes the table functions, the query form needs its SELECT,
		// and INTO TABLE is required.
		"LOAD DATA FROM 's'(FORMAT(TYPE = 's')) INTO TABLE a",
		"LOAD DATA FROM(a FROM SOURCE(TYPE = 's')) INTO TABLE a",
		"LOAD DATA FROM SOURCE(TYPE = 's') TABLE a",
		"LOAD DATA FROM SOURCE(TYPE = 's') INTO a",
		"LOAD DATA FROM SOURCE(TYPE = 's') INTO TABLE a LOG",
		"LOAD DATA FROM SOURCE(TYPE = 's') INTO TABLE a LOG ERRORS REJECT",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}

// TestExternalTableTree pins what an external table reference parses into: its
// function, path, properties, format, pattern, alias and Location, in DELETE
// and in ALTER EXTERNAL TABLE.
func TestExternalTableTree(t *testing.T) {
	const sql = "DELETE FROM FILES(LOCATION 's', FORMAT(TYPE = 'csv'), PATTERN 'p') AS f"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	stmt, ok := list.Items[0].(*nodes.DeleteStmt)
	if !ok || len(stmt.Tables) != 1 {
		t.Fatalf("got %#v, want one DeleteStmt with one table", list.Items[0])
	}
	ref, ok := stmt.Tables[0].(*nodes.ExternalTableRef)
	if !ok {
		t.Fatalf("got %T, want *nodes.ExternalTableRef", stmt.Tables[0])
	}
	if ref.Function != "FILES" || ref.Path != "s" || ref.Pattern != "p" || ref.Alias != "f" {
		t.Errorf("ref = %#v, want FILES of s with pattern p and alias f", ref)
	}
	if len(ref.Format) != 1 || ref.Format[0].Name != "TYPE" {
		t.Errorf("format = %#v, want one TYPE entry", ref.Format)
	}
	if ref.Loc.Start != len("DELETE FROM ") || ref.Loc.End != len(sql) {
		t.Errorf("Loc = %+v, want the reference through its alias", ref.Loc)
	}

	list, err = Parse("ALTER EXTERNAL TABLE a ADD PARTITION(ds = '2024') LOCATION 'sink'")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	alter, ok := list.Items[0].(*nodes.AlterTableStmt)
	if !ok || !alter.External || len(alter.Commands) != 1 {
		t.Fatalf("got %#v, want one external AlterTableStmt with one action", list.Items[0])
	}
	cmd := alter.Commands[0]
	if cmd.Type != nodes.ATAddExternalPartition || cmd.Location != "sink" {
		t.Errorf("cmd = %#v, want ATAddExternalPartition to sink", cmd)
	}
	if len(cmd.PartitionValues) != 1 || cmd.PartitionValues[0].Column.Column != "ds" {
		t.Errorf("partitions = %#v, want one for ds", cmd.PartitionValues)
	}
}
