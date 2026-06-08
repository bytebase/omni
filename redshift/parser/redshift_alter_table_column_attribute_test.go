package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftAlterTableAddColumnEncodeParse(t *testing.T) {
	stmt := firstAlterTableStmt(t, "ALTER TABLE events ADD COLUMN event_data VARCHAR(MAX) ENCODE ZSTD;")
	if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
		t.Fatalf("expected one command, got %#v", stmt.Cmds)
	}
	cmd, ok := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if !ok {
		t.Fatalf("expected AlterTableCmd, got %T", stmt.Cmds.Items[0])
	}
	if cmd.Subtype != int(nodes.AT_AddColumn) {
		t.Fatalf("expected AT_AddColumn, got %v", cmd.Subtype)
	}
	col, ok := cmd.Def.(*nodes.ColumnDef)
	if !ok {
		t.Fatalf("expected ColumnDef, got %T", cmd.Def)
	}
	if col.Compression != "zstd" {
		t.Fatalf("expected compression zstd, got %q", col.Compression)
	}
}

func TestRedshiftAlterTableAlterColumnEncodeParse(t *testing.T) {
	stmt := firstAlterTableStmt(t, "ALTER TABLE events ALTER COLUMN event_data ENCODE LZO;")
	if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
		t.Fatalf("expected one command, got %#v", stmt.Cmds)
	}
	cmd, ok := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if !ok {
		t.Fatalf("expected AlterTableCmd, got %T", stmt.Cmds.Items[0])
	}
	if cmd.Subtype != int(nodes.AT_SetCompression) {
		t.Fatalf("expected AT_SetCompression, got %v", cmd.Subtype)
	}
	value, ok := cmd.Def.(*nodes.String)
	if !ok {
		t.Fatalf("expected String def, got %T", cmd.Def)
	}
	if value.Str != "lzo" {
		t.Fatalf("expected compression lzo, got %q", value.Str)
	}
}

func TestRedshiftAlterTableEncodeAutoParse(t *testing.T) {
	stmt := firstAlterTableStmt(t, "ALTER TABLE large_table ALTER ENCODE AUTO;")
	if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
		t.Fatalf("expected one command, got %#v", stmt.Cmds)
	}
	cmd, ok := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
	if !ok {
		t.Fatalf("expected AlterTableCmd, got %T", stmt.Cmds.Items[0])
	}
	if cmd.Subtype != int(nodes.AT_RedshiftAlterOption) {
		t.Fatalf("expected AT_RedshiftAlterOption, got %v", cmd.Subtype)
	}
	if cmd.Name != "encode" {
		t.Fatalf("expected option name encode, got %q", cmd.Name)
	}
	opts, ok := cmd.Def.(*nodes.List)
	if !ok || len(opts.Items) != 1 {
		t.Fatalf("expected one option, got %#v", cmd.Def)
	}
	elem, ok := opts.Items[0].(*nodes.DefElem)
	if !ok {
		t.Fatalf("expected DefElem, got %T", opts.Items[0])
	}
	if elem.Defname != "value" {
		t.Fatalf("expected value option, got %q", elem.Defname)
	}
	value, ok := elem.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected String arg, got %T", elem.Arg)
	}
	if value.Str != "auto" {
		t.Fatalf("expected encode auto, got %q", value.Str)
	}
}

func TestRedshiftAlterTableMaskingForDatasharesParse(t *testing.T) {
	tests := []struct {
		sql     string
		enabled bool
	}{
		{"ALTER TABLE shared_table MASKING ON FOR DATASHARES;", true},
		{"ALTER TABLE unmasked_table MASKING OFF FOR DATASHARES;", false},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			stmt := firstAlterTableStmt(t, tt.sql)
			if stmt.Cmds == nil || len(stmt.Cmds.Items) != 1 {
				t.Fatalf("expected one command, got %#v", stmt.Cmds)
			}
			cmd, ok := stmt.Cmds.Items[0].(*nodes.AlterTableCmd)
			if !ok {
				t.Fatalf("expected AlterTableCmd, got %T", stmt.Cmds.Items[0])
			}
			if cmd.Subtype != int(nodes.AT_RedshiftAlterOption) {
				t.Fatalf("expected AT_RedshiftAlterOption, got %v", cmd.Subtype)
			}
			if cmd.Name != "masking" {
				t.Fatalf("expected option name masking, got %q", cmd.Name)
			}
			opts, ok := cmd.Def.(*nodes.List)
			if !ok || len(opts.Items) != 2 {
				t.Fatalf("expected two options, got %#v", cmd.Def)
			}
			enabled := findDefElem(opts, "enabled")
			if enabled == nil {
				t.Fatalf("expected enabled option, got %#v", opts)
			}
			value, ok := enabled.Arg.(*nodes.Boolean)
			if !ok {
				t.Fatalf("expected Boolean enabled arg, got %T", enabled.Arg)
			}
			if value.Boolval != tt.enabled {
				t.Fatalf("expected enabled=%v, got %v", tt.enabled, value.Boolval)
			}
			if findDefElem(opts, "for_datashares") == nil {
				t.Fatalf("expected for_datashares option, got %#v", opts)
			}
		})
	}
}
