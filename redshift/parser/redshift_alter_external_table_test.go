package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftAlterExternalTableOptionsParse(t *testing.T) {
	tests := []string{
		"ALTER TABLE spectrum.sales SET TABLE PROPERTIES ('numRows'='1000000', 'compressionType'='gzip');",
		"ALTER TABLE spectrum.events SET LOCATION 's3://mybucket/data/events/2024/';",
		"ALTER TABLE spectrum.logs SET FILE FORMAT PARQUET;",
		"ALTER TABLE spectrum.sales ADD PARTITION (year=2024, month=1) LOCATION 's3://mybucket/sales/2024/01/';",
		"ALTER TABLE spectrum.events ADD IF NOT EXISTS PARTITION (date='2024-01-01') LOCATION 's3://mybucket/events/2024-01-01/';",
		"ALTER TABLE spectrum.sales DROP PARTITION (year=2023, month=12);",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
		})
	}
}

// TestRedshiftAlterExternalTablePartitionLocationParse pins the two
// partition clauses r_ALTER_TABLE gives an external table: PARTITION (...)
// SET LOCATION, and ADD PARTITION (...) LOCATION with the docs' comma list of
// locations.
func TestRedshiftAlterExternalTablePartitionLocationParse(t *testing.T) {
	stmt := firstAlterTableStmt(t, "ALTER TABLE spectrum.sales PARTITION(year = 2024) SET LOCATION 's3://mybucket/sales/2024/';")
	cmd, ok := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if !ok {
		t.Fatalf("expected AlterTableCmd, got %T", stmt.Cmds.Items[0])
	}
	if cmd.Name != "set_partition_location" {
		t.Fatalf("command = %q, want set_partition_location", cmd.Name)
	}
	assertAlterTableCmdOptionString(t, cmd, "location", "s3://mybucket/sales/2024/")
	assertAlterTableCmdOptionList(t, cmd, "partition", 1)

	add := firstAlterTableStmt(t, "ALTER TABLE spectrum.sales ADD PARTITION(year = 2024) LOCATION 's3://mybucket/sales/2024/' , 's3://mybucket/sales/2024b/';")
	addCmd, ok := add.Cmds.Items[0].(*nodes.AlterTableCmd)
	if !ok {
		t.Fatalf("expected AlterTableCmd, got %T", add.Cmds.Items[0])
	}
	if addCmd.Name != "add_partition" {
		t.Fatalf("command = %q, want add_partition", addCmd.Name)
	}
	// The docs write one location per string, so each string is an option.
	locations := 0
	for _, item := range alterTableCmdOptions(t, addCmd).Items {
		if elem, ok := item.(*nodes.DefElem); ok && elem.Defname == "location" {
			locations++
		}
	}
	if locations != 2 {
		t.Fatalf("locations = %d, want two", locations)
	}

	if _, err := Parse("ALTER TABLE spectrum.sales PARTITION(year = 2024)"); err == nil {
		t.Error("PARTITION (...) with no SET LOCATION parsed, want a syntax error")
	}
	if _, err := Parse("ALTER TABLE spectrum.sales ADD PARTITION(year = 2024) LOCATION 's3://b/a',"); err == nil {
		t.Error("ADD PARTITION with a trailing comma parsed, want a syntax error")
	}
}
