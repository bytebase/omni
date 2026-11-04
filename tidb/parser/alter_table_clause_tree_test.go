package parser

import (
	"fmt"
	"testing"

	ast "github.com/bytebase/omni/tidb/ast"
)

// These tests assert what the ALTER TABLE clause readers build, not that the
// text parsed. The corpus records only a verdict, so a clause that parses into
// the wrong node, loses a value, or swallows the clause beside it reads as a
// closed gap there; the multi-clause cases below are the ones that catch a
// reader taking too much or too little.

// alterCmds parses one ALTER TABLE through the public entry point and returns
// its commands, so what a consumer of Parse sees is what is asserted.
func alterCmds(t *testing.T, sql string) []*ast.AlterTableCmd {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*ast.AlterTableStmt)
	if !ok {
		t.Fatalf("Parse(%q): %T, want *ast.AlterTableStmt", sql, list.Items[0])
	}
	return stmt.Commands
}

// alterOneCmd parses an ALTER TABLE that must hold exactly one command.
func alterOneCmd(t *testing.T, sql string) *ast.AlterTableCmd {
	t.Helper()
	cmds := alterCmds(t, sql)
	if len(cmds) != 1 {
		t.Fatalf("Parse(%q): %d commands, want 1", sql, len(cmds))
	}
	return cmds[0]
}

func TestAlterTableCacheTree(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want ast.AlterTableCmdType
	}{
		{"ALTER TABLE t CACHE", ast.ATCache},
		{"ALTER TABLE t NOCACHE", ast.ATNoCache},
	} {
		if got := alterOneCmd(t, c.sql).Type; got != c.want {
			t.Errorf("Parse(%q): Type = %d, want %d", c.sql, got, c.want)
		}
	}

	cmds := alterCmds(t, "ALTER TABLE t CACHE, NOCACHE")
	if len(cmds) != 2 {
		t.Fatalf("CACHE, NOCACHE: %d commands, want 2", len(cmds))
	}
	if cmds[0].Type != ast.ATCache || cmds[1].Type != ast.ATNoCache {
		t.Errorf("CACHE, NOCACHE: types = %d, %d, want %d, %d",
			cmds[0].Type, cmds[1].Type, ast.ATCache, ast.ATNoCache)
	}
}

func TestAlterTableWriteableTree(t *testing.T) {
	cmds := alterCmds(t, "ALTER TABLE t READ ONLY, READ WRITE")
	if len(cmds) != 2 {
		t.Fatalf("%d commands, want 2", len(cmds))
	}
	for i, want := range []bool{false, true} {
		if cmds[i].Type != ast.ATWriteable {
			t.Errorf("command %d: Type = %d, want ATWriteable", i, cmds[i].Type)
		}
		if cmds[i].Writeable != want {
			t.Errorf("command %d: Writeable = %v, want %v", i, cmds[i].Writeable, want)
		}
	}
}

func TestAlterTableCompactTree(t *testing.T) {
	for _, c := range []struct {
		sql     string
		names   []string
		tiflash bool
	}{
		{"alter table t compact", nil, false},
		{"alter table t compact tiflash replica", nil, true},
		{"alter table t compact partition p1,p2", []string{"p1", "p2"}, false},
		{"alter table t compact partition p1,p2 tiflash replica", []string{"p1", "p2"}, true},
	} {
		cmd := alterOneCmd(t, c.sql)
		if cmd.Type != ast.ATCompact {
			t.Errorf("Parse(%q): Type = %d, want ATCompact", c.sql, cmd.Type)
		}
		if got := cmd.CompactTiFlashReplica; got != c.tiflash {
			t.Errorf("Parse(%q): CompactTiFlashReplica = %v, want %v", c.sql, got, c.tiflash)
		}
		if len(cmd.PartitionNames) != len(c.names) {
			t.Errorf("Parse(%q): PartitionNames = %v, want %v", c.sql, cmd.PartitionNames, c.names)
			continue
		}
		for i, n := range c.names {
			if cmd.PartitionNames[i] != n {
				t.Errorf("Parse(%q): PartitionNames[%d] = %q, want %q", c.sql, i, cmd.PartitionNames[i], n)
			}
		}
	}
}

func TestAlterTableAttributesTree(t *testing.T) {
	for _, c := range []struct {
		sql         string
		wantType    ast.AlterTableCmdType
		wantDefault bool
		wantValue   string
	}{
		{"ALTER TABLE t ATTRIBUTES='str1,str2'", ast.ATAttributes, false, "str1,str2"},
		{"ALTER TABLE t ATTRIBUTES 'str1,str2'", ast.ATAttributes, false, "str1,str2"},
		{"ALTER TABLE t ATTRIBUTES=DeFaUlT", ast.ATAttributes, true, ""},
		{"ALTER TABLE t STATS_OPTIONS=\"str1,str2\"", ast.ATStatsOptions, false, "str1,str2"},
		{"ALTER TABLE t STATS_OPTIONS=DEFAULT", ast.ATStatsOptions, true, ""},
	} {
		cmd := alterOneCmd(t, c.sql)
		if cmd.Type != c.wantType {
			t.Errorf("Parse(%q): Type = %d, want %d", c.sql, cmd.Type, c.wantType)
		}
		if cmd.Attributes == nil {
			t.Errorf("Parse(%q): Attributes = nil", c.sql)
			continue
		}
		if cmd.Attributes.Default != c.wantDefault {
			t.Errorf("Parse(%q): Default = %v, want %v", c.sql, cmd.Attributes.Default, c.wantDefault)
		}
		if cmd.Attributes.Value != c.wantValue {
			t.Errorf("Parse(%q): Value = %q, want %q", c.sql, cmd.Attributes.Value, c.wantValue)
		}
	}

	// STATS_OPTIONS and ATTRIBUTES in one list, in that order.
	cmds := alterCmds(t, "ALTER TABLE t STATS_OPTIONS = 'a', ATTRIBUTES = 'b'")
	if len(cmds) != 2 {
		t.Fatalf("%d commands, want 2", len(cmds))
	}
	if cmds[0].Type != ast.ATStatsOptions || cmds[0].Attributes.Value != "a" {
		t.Errorf("command 0 = %d/%v, want ATStatsOptions/'a'", cmds[0].Type, cmds[0].Attributes)
	}
	if cmds[1].Type != ast.ATAttributes || cmds[1].Attributes.Value != "b" {
		t.Errorf("command 1 = %d/%v, want ATAttributes/'b'", cmds[1].Type, cmds[1].Attributes)
	}
}

func TestAlterTablePartitionClauseTree(t *testing.T) {
	cmd := alterOneCmd(t, "ALTER TABLE t PARTITION p ATTRIBUTES='str'")
	if cmd.Type != ast.ATPartitionAttributes {
		t.Fatalf("Type = %d, want ATPartitionAttributes", cmd.Type)
	}
	if len(cmd.PartitionNames) != 1 || cmd.PartitionNames[0] != "p" {
		t.Errorf("PartitionNames = %v, want [p]", cmd.PartitionNames)
	}
	if cmd.Attributes == nil || cmd.Attributes.Value != "str" {
		t.Errorf("Attributes = %v, want value \"str\"", cmd.Attributes)
	}

	// PARTITION name with the partition definition's options, where
	// STORAGE ENGINE is one option and not STORAGE followed by ENGINE.
	cmd = alterOneCmd(t, "ALTER TABLE t PARTITION p STORAGE ENGINE = x MAX_ROWS 5")
	if cmd.Type != ast.ATPartitionOptions {
		t.Fatalf("Type = %d, want ATPartitionOptions", cmd.Type)
	}
	if len(cmd.PartitionOptions) != 2 {
		t.Fatalf("PartitionOptions = %v, want 2", cmd.PartitionOptions)
	}
	if got := cmd.PartitionOptions[0]; got.Name != "STORAGE ENGINE" || got.Value != "x" {
		t.Errorf("PartitionOptions[0] = %q=%q, want \"STORAGE ENGINE\"=\"x\"", got.Name, got.Value)
	}
	if got := cmd.PartitionOptions[1]; got.Name != "MAX_ROWS" || got.Value != "5" {
		t.Errorf("PartitionOptions[1] = %q=%q, want \"MAX_ROWS\"=\"5\"", got.Name, got.Value)
	}

	// A bare PARTITION name carries no options at all.
	cmd = alterOneCmd(t, "ALTER TABLE t PARTITION p")
	if cmd.Type != ast.ATPartitionOptions || len(cmd.PartitionOptions) != 0 {
		t.Errorf("PARTITION p: Type = %d, options = %v, want ATPartitionOptions and none", cmd.Type, cmd.PartitionOptions)
	}

	// The clause closes the statement, so an action before it is its own
	// command and the clause is the last one.
	cmds := alterCmds(t, "ALTER TABLE t ADD COLUMN a INT PARTITION p ENGINE = x")
	if len(cmds) != 2 {
		t.Fatalf("%d commands, want 2", len(cmds))
	}
	if cmds[0].Type != ast.ATAddColumn || cmds[0].Column == nil || cmds[0].Column.Name != "a" {
		t.Errorf("command 0 = %d, want ATAddColumn of a", cmds[0].Type)
	}
	if cmds[1].Type != ast.ATPartitionOptions || len(cmds[1].PartitionOptions) != 1 {
		t.Errorf("command 1 = %d with %v, want ATPartitionOptions with one option", cmds[1].Type, cmds[1].PartitionOptions)
	}
}

// TestAlterTableMultiClauseTree is the case the brief singles out: five
// clauses, two separators and a table-option run in the middle. Each clause
// must be its own command, in the order written, with its own value.
func TestAlterTableMultiClauseTree(t *testing.T) {
	const sql = "ALTER TABLE d_n.t_n ALGORITHM = DEFAULT , MAX_ROWS 10, UNION ( d_n.t_n ) , ROW_FORMAT REDUNDANT, STATS_PERSISTENT = DEFAULT"
	cmds := alterCmds(t, sql)
	if len(cmds) != 5 {
		t.Fatalf("%d commands, want 5: %v", len(cmds), cmds)
	}
	if cmds[0].Type != ast.ATAlgorithm || cmds[0].Name != "DEFAULT" {
		t.Errorf("command 0 = %d/%q, want ATAlgorithm/DEFAULT", cmds[0].Type, cmds[0].Name)
	}
	for i, want := range []struct{ name, value string }{
		{"MAX_ROWS", "10"},
		{"UNION", "d_n.t_n"},
		{"ROW_FORMAT", "REDUNDANT"},
		{"STATS_PERSISTENT", "DEFAULT"},
	} {
		cmd := cmds[i+1]
		if cmd.Type != ast.ATTableOption {
			t.Errorf("command %d: Type = %d, want ATTableOption", i+1, cmd.Type)
			continue
		}
		if cmd.Option == nil || cmd.Option.Name != want.name || cmd.Option.Value != want.value {
			t.Errorf("command %d: Option = %v, want %s=%s", i+1, cmd.Option, want.name, want.value)
		}
	}

	// The same options chained with spaces instead of commas are still one
	// command each.
	cmds = alterCmds(t, "ALTER TABLE t ROW_FORMAT REDUNDANT STATS_PERSISTENT = DEFAULT KEY_BLOCK_SIZE 1")
	if len(cmds) != 3 {
		t.Fatalf("space-chained options: %d commands, want 3", len(cmds))
	}
	for i, want := range []string{"ROW_FORMAT", "STATS_PERSISTENT", "KEY_BLOCK_SIZE"} {
		if cmds[i].Option == nil || cmds[i].Option.Name != want {
			t.Errorf("command %d: Option = %v, want %s", i, cmds[i].Option, want)
		}
	}

	// LOCK, then a table-option run, then the closing REMOVE PARTITIONING,
	// which takes no comma of its own.
	cmds = alterCmds(t, "ALTER TABLE d_n.t_n LOCK DEFAULT , UNION = ( t_n , d_n.t_n ) REMOVE PARTITIONING")
	if len(cmds) != 3 {
		t.Fatalf("%d commands, want 3: %v", len(cmds), cmds)
	}
	if cmds[0].Type != ast.ATLock || cmds[0].Name != "DEFAULT" {
		t.Errorf("command 0 = %d/%q, want ATLock/DEFAULT", cmds[0].Type, cmds[0].Name)
	}
	if cmds[1].Option == nil || cmds[1].Option.Name != "UNION" || cmds[1].Option.Value != "t_n,d_n.t_n" {
		t.Errorf("command 1 = %v, want UNION=t_n,d_n.t_n", cmds[1].Option)
	}
	if cmds[2].Type != ast.ATRemovePartitioning {
		t.Errorf("command 2 = %d, want ATRemovePartitioning", cmds[2].Type)
	}

	// An empty UNION list is an option with no value, not a missing option.
	cmds = alterCmds(t, "ALTER TABLE t_n UNION ( ), KEY_BLOCK_SIZE = 1")
	if len(cmds) != 2 {
		t.Fatalf("empty UNION: %d commands, want 2", len(cmds))
	}
	if cmds[0].Option == nil || cmds[0].Option.Name != "UNION" || cmds[0].Option.Value != "" {
		t.Errorf("command 0 = %v, want UNION with an empty value", cmds[0].Option)
	}
	if cmds[1].Option == nil || cmds[1].Option.Name != "KEY_BLOCK_SIZE" {
		t.Errorf("command 1 = %v, want KEY_BLOCK_SIZE", cmds[1].Option)
	}
}

func TestAlterTableForceOptionTree(t *testing.T) {
	// FORCE alone is the action; FORCE before AUTO_INCREMENT is part of the
	// option, and the two must not be read as one another.
	cmd := alterOneCmd(t, "ALTER TABLE t FORCE")
	if cmd.Type != ast.ATForce {
		t.Errorf("FORCE: Type = %d, want ATForce", cmd.Type)
	}

	cmd = alterOneCmd(t, "ALTER TABLE t FORCE AUTO_INCREMENT = 3")
	if cmd.Type != ast.ATTableOption {
		t.Fatalf("FORCE AUTO_INCREMENT: Type = %d, want ATTableOption", cmd.Type)
	}
	if cmd.Option == nil || cmd.Option.Name != "AUTO_INCREMENT" || cmd.Option.Value != "3" || !cmd.Option.Force {
		t.Errorf("FORCE AUTO_INCREMENT: Option = %+v, want AUTO_INCREMENT=3 with Force", cmd.Option)
	}

	cmds := alterCmds(t, "alter table t auto_increment 30, force auto_random_base 40")
	if len(cmds) != 2 {
		t.Fatalf("%d commands, want 2", len(cmds))
	}
	if o := cmds[0].Option; o == nil || o.Name != "AUTO_INCREMENT" || o.Value != "30" || o.Force {
		t.Errorf("command 0 = %+v, want AUTO_INCREMENT=30 without Force", o)
	}
	if o := cmds[1].Option; o == nil || o.Name != "AUTO_RANDOM_BASE" || o.Value != "40" || !o.Force {
		t.Errorf("command 1 = %+v, want AUTO_RANDOM_BASE=40 with Force", o)
	}

	// FORCE, AUTO_INCREMENT with a comma between them is the action and the
	// plain option, not the FORCE-prefixed one.
	cmds = alterCmds(t, "ALTER TABLE t FORCE, AUTO_INCREMENT 3")
	if len(cmds) != 2 {
		t.Fatalf("%d commands, want 2", len(cmds))
	}
	if cmds[0].Type != ast.ATForce {
		t.Errorf("command 0 = %d, want ATForce", cmds[0].Type)
	}
	if o := cmds[1].Option; o == nil || o.Name != "AUTO_INCREMENT" || o.Force {
		t.Errorf("command 1 = %+v, want AUTO_INCREMENT without Force", o)
	}
}

func TestAlterTableIntervalPartitionTree(t *testing.T) {
	for _, c := range []struct {
		sql      string
		wantType ast.AlterTableCmdType
	}{
		{"ALTER TABLE t LAST PARTITION LESS THAN (1000)", ast.ATAddLastPartition},
		{"ALTER TABLE t FIRST PARTITION LESS THAN (1000)", ast.ATDropFirstPartition},
		{"ALTER TABLE t SPLIT MAXVALUE PARTITION LESS THAN (1000)", ast.ATReorganizeLastPartition},
		{"ALTER TABLE t MERGE FIRST PARTITION LESS THAN (1000)", ast.ATReorganizeFirstPartition},
	} {
		cmd := alterOneCmd(t, c.sql)
		if cmd.Type != c.wantType {
			t.Errorf("Parse(%q): Type = %d, want %d", c.sql, cmd.Type, c.wantType)
		}
		if cmd.PartitionBound == nil {
			t.Errorf("Parse(%q): PartitionBound = nil", c.sql)
		}
	}

	cmd := alterOneCmd(t, "ALTER TABLE t LAST PARTITION LESS THAN (1000) NO_WRITE_TO_BINLOG")
	if !cmd.NoWriteToBinlog {
		t.Errorf("LAST PARTITION ... NO_WRITE_TO_BINLOG: NoWriteToBinlog = false")
	}
	cmd = alterOneCmd(t, "ALTER TABLE t FIRST PARTITION LESS THAN (1000) IF EXISTS")
	if !cmd.IfExists {
		t.Errorf("FIRST PARTITION ... IF EXISTS: IfExists = false")
	}

	// LAST is an entry of the list, so an action may follow it.
	cmds := alterCmds(t, "ALTER TABLE t LAST PARTITION LESS THAN (1), ADD COLUMN a INT")
	if len(cmds) != 2 {
		t.Fatalf("%d commands, want 2", len(cmds))
	}
	if cmds[0].Type != ast.ATAddLastPartition || cmds[1].Type != ast.ATAddColumn {
		t.Errorf("types = %d, %d, want ATAddLastPartition, ATAddColumn", cmds[0].Type, cmds[1].Type)
	}
}

func TestAlterTableAnalyzePartitionTree(t *testing.T) {
	cmd := alterOneCmd(t, "ALTER TABLE t ANALYZE PARTITION p, q INDEX i, j WITH 4 BUCKETS, 3 TOPN, 0.5 SAMPLERATE, 2 CMSKETCH DEPTH, 7 CMSKETCH WIDTH, 9 SAMPLES")
	if cmd.Type != ast.ATAnalyzePartition {
		t.Fatalf("Type = %d, want ATAnalyzePartition", cmd.Type)
	}
	if got := cmd.PartitionNames; len(got) != 2 || got[0] != "p" || got[1] != "q" {
		t.Errorf("PartitionNames = %v, want [p q]", got)
	}
	if !cmd.AnalyzeIndexFlag {
		t.Errorf("AnalyzeIndexFlag = false, want true")
	}
	if got := cmd.AnalyzeIndexes; len(got) != 2 || got[0] != "i" || got[1] != "j" {
		t.Errorf("AnalyzeIndexes = %v, want [i j]", got)
	}
	// The WITH list is the one ANALYZE TABLE takes, with the same unit type.
	want := []struct {
		unit  ast.AnalyzeOptionType
		value string
	}{
		{ast.AnalyzeOptBuckets, "4"},
		{ast.AnalyzeOptTopN, "3"},
		{ast.AnalyzeOptSampleRate, "0.5"},
		{ast.AnalyzeOptCMSketchDepth, "2"},
		{ast.AnalyzeOptCMSketchWidth, "7"},
		{ast.AnalyzeOptSamples, "9"},
	}
	if len(cmd.AnalyzeOptions) != len(want) {
		t.Fatalf("AnalyzeOptions = %v, want %d entries", cmd.AnalyzeOptions, len(want))
	}
	for i, w := range want {
		got := cmd.AnalyzeOptions[i]
		if got.Type != w.unit || got.Value != w.value {
			t.Errorf("AnalyzeOptions[%d] = %s=%s, want %s=%s", i, got.Type, got.Value, w.unit, w.value)
		}
	}

	// INDEX with no names keeps the keyword without inventing a name.
	cmd = alterOneCmd(t, "ALTER TABLE t ANALYZE PARTITION p INDEX")
	if !cmd.AnalyzeIndexFlag || len(cmd.AnalyzeIndexes) != 0 {
		t.Errorf("INDEX with no names: flag = %v, names = %v, want true and none", cmd.AnalyzeIndexFlag, cmd.AnalyzeIndexes)
	}
	// No INDEX at all leaves the flag off, so a consumer can tell the two apart.
	cmd = alterOneCmd(t, "ALTER TABLE t ANALYZE PARTITION p")
	if cmd.AnalyzeIndexFlag {
		t.Errorf("no INDEX: AnalyzeIndexFlag = true, want false")
	}
}

func TestAlterTableIfExistsTree(t *testing.T) {
	for _, c := range []struct {
		sql         string
		wantType    ast.AlterTableCmdType
		ifNotExists bool
		ifExists    bool
	}{
		{"ALTER TABLE t ADD COLUMN IF NOT EXISTS a INT", ast.ATAddColumn, true, false},
		{"ALTER TABLE t ADD COLUMN a INT", ast.ATAddColumn, false, false},
		{"ALTER TABLE t ADD PARTITION IF NOT EXISTS PARTITIONS 5", ast.ATAddPartition, true, false},
		{"ALTER TABLE t MODIFY COLUMN IF EXISTS a INT", ast.ATModifyColumn, false, true},
		{"ALTER TABLE t CHANGE COLUMN IF EXISTS a b INT", ast.ATChangeColumn, false, true},
		{"ALTER TABLE t DROP INDEX IF EXISTS i", ast.ATDropIndex, false, true},
		{"ALTER TABLE t DROP PARTITION IF EXISTS p", ast.ATDropPartition, false, true},
	} {
		cmd := alterOneCmd(t, c.sql)
		if cmd.Type != c.wantType {
			t.Errorf("Parse(%q): Type = %d, want %d", c.sql, cmd.Type, c.wantType)
		}
		if cmd.IfNotExists != c.ifNotExists {
			t.Errorf("Parse(%q): IfNotExists = %v, want %v", c.sql, cmd.IfNotExists, c.ifNotExists)
		}
		if cmd.IfExists != c.ifExists {
			t.Errorf("Parse(%q): IfExists = %v, want %v", c.sql, cmd.IfExists, c.ifExists)
		}
	}

	// The IF NOT EXISTS of an index belongs to the constraint, not the action.
	cmd := alterOneCmd(t, "ALTER TABLE t ADD INDEX IF NOT EXISTS i (a)")
	if cmd.Constraint == nil || !cmd.Constraint.IfNotExists {
		t.Errorf("ADD INDEX IF NOT EXISTS: Constraint = %+v, want IfNotExists", cmd.Constraint)
	}
	cmd = alterOneCmd(t, "ALTER TABLE t ADD CONSTRAINT c FOREIGN KEY IF NOT EXISTS (a) REFERENCES u(b)")
	if cmd.Constraint == nil || !cmd.Constraint.IfNotExists || cmd.Constraint.Name != "c" {
		t.Errorf("FOREIGN KEY IF NOT EXISTS: Constraint = %+v, want name c with IfNotExists", cmd.Constraint)
	}
	// IF EXISTS must not be mistaken for the dropped column's name.
	cmd = alterOneCmd(t, "ALTER TABLE t DROP PARTITION IF EXISTS p1, p2")
	if got := cmd.PartitionNames; len(got) != 2 || got[0] != "p1" || got[1] != "p2" {
		t.Errorf("DROP PARTITION IF EXISTS p1, p2: PartitionNames = %v, want [p1 p2]", got)
	}
}

func TestAlterTableAddTableElementListTree(t *testing.T) {
	cmd := alterOneCmd(t, "ALTER TABLE t ADD COLUMN (a1 int, a2 int, CONSTRAINT b1 CHECK (a1>0), PRIMARY KEY (a1))")
	if cmd.Type != ast.ATAddColumn {
		t.Fatalf("Type = %d, want ATAddColumn", cmd.Type)
	}
	if len(cmd.Columns) != 2 || cmd.Columns[0].Name != "a1" || cmd.Columns[1].Name != "a2" {
		t.Errorf("Columns = %v, want a1 and a2", cmd.Columns)
	}
	if len(cmd.Constraints) != 2 {
		t.Fatalf("Constraints = %v, want 2", cmd.Constraints)
	}
	if cmd.Constraints[0].Type != ast.ConstrCheck || cmd.Constraints[0].Name != "b1" {
		t.Errorf("Constraints[0] = %+v, want a CHECK named b1", cmd.Constraints[0])
	}
	if cmd.Constraints[1].Type != ast.ConstrPrimaryKey {
		t.Errorf("Constraints[1] = %+v, want a PRIMARY KEY", cmd.Constraints[1])
	}

	// A list of constraints alone adds no column.
	cmd = alterOneCmd(t, "ALTER TABLE t ADD (PRIMARY KEY (a))")
	if len(cmd.Columns) != 0 || len(cmd.Constraints) != 1 {
		t.Errorf("ADD (PRIMARY KEY (a)): columns = %v, constraints = %v, want none and one", cmd.Columns, cmd.Constraints)
	}
}

func TestAlterTableRenameTree(t *testing.T) {
	for _, sql := range []string{
		"ALTER TABLE db.t RENAME to db1.t1",
		"ALTER TABLE db.t RENAME db1.t1",
		"ALTER TABLE db.t RENAME = db1.t1",
		"ALTER TABLE db.t RENAME as db1.t1",
	} {
		cmd := alterOneCmd(t, sql)
		if cmd.Type != ast.ATRenameTable {
			t.Errorf("Parse(%q): Type = %d, want ATRenameTable", sql, cmd.Type)
		}
		if cmd.NewTable == nil || cmd.NewTable.Schema != "db1" || cmd.NewTable.Name != "t1" {
			t.Errorf("Parse(%q): NewTable = %+v, want db1.t1", sql, cmd.NewTable)
		}
		if cmd.NewName != "db1.t1" {
			t.Errorf("Parse(%q): NewName = %q, want \"db1.t1\"", sql, cmd.NewName)
		}
	}

	cmd := alterOneCmd(t, "ALTER TABLE t RENAME = t1")
	if cmd.NewTable == nil || cmd.NewTable.Schema != "" || cmd.NewTable.Name != "t1" || cmd.NewName != "t1" {
		t.Errorf("RENAME = t1: NewTable = %+v, NewName = %q", cmd.NewTable, cmd.NewName)
	}
}

func TestAlterTableConvertCharsetTree(t *testing.T) {
	for _, c := range []struct {
		sql       string
		charset   string
		collation string
	}{
		{"ALTER TABLE t CONVERT TO CHARSET UTF8", "UTF8", ""},
		{"ALTER TABLE t CONVERT TO CHAR SET UTF8", "UTF8", ""},
		{"ALTER TABLE t CONVERT TO CHARACTER SET UTF8", "UTF8", ""},
		{"ALTER TABLE t CONVERT TO CHARSET UTF8 COLLATE UTF8_BIN", "UTF8", "UTF8_BIN"},
		{"alter table d_n.t_n convert to charset default", "DEFAULT", ""},
		{"alter table d_n.t_n convert to char set default", "DEFAULT", ""},
	} {
		cmd := alterOneCmd(t, c.sql)
		if cmd.Type != ast.ATConvertCharset {
			t.Errorf("Parse(%q): Type = %d, want ATConvertCharset", c.sql, cmd.Type)
		}
		if cmd.Name != c.charset {
			t.Errorf("Parse(%q): charset = %q, want %q", c.sql, cmd.Name, c.charset)
		}
		if cmd.NewName != c.collation {
			t.Errorf("Parse(%q): collation = %q, want %q", c.sql, cmd.NewName, c.collation)
		}
	}
}

func TestAlterTableIndexOptionGlobalLocalTree(t *testing.T) {
	// GLOBAL, LOCAL, CLUSTERED and NONCLUSTERED are index options, kept in
	// the option list in the order written, and the clustering is also
	// hoisted so a consumer need not read the list.
	cmd := alterOneCmd(t, "ALTER TABLE t ADD PRIMARY KEY (a) CLUSTERED GLOBAL")
	if cmd.Constraint == nil {
		t.Fatal("Constraint = nil")
	}
	var names []string
	for _, o := range cmd.Constraint.IndexOptions {
		names = append(names, o.Name)
	}
	if len(names) != 2 || names[0] != "CLUSTERED" || names[1] != "GLOBAL" {
		t.Errorf("IndexOptions = %v, want [CLUSTERED GLOBAL]", names)
	}
	if cmd.Constraint.Clustered == nil || !*cmd.Constraint.Clustered {
		t.Errorf("Clustered = %v, want true", cmd.Constraint.Clustered)
	}

	cmd = alterOneCmd(t, "ALTER TABLE t ADD INDEX (a) NONCLUSTERED")
	if cmd.Constraint == nil || cmd.Constraint.Clustered == nil || *cmd.Constraint.Clustered {
		t.Errorf("NONCLUSTERED: Clustered = %v, want false", cmd.Constraint.Clustered)
	}

	// An option list may hold both words; the last one decides, as the
	// grammar's own action does.
	cmd = alterOneCmd(t, "ALTER TABLE t ADD INDEX (a) GLOBAL LOCAL")
	names = nil
	for _, o := range cmd.Constraint.IndexOptions {
		names = append(names, o.Name)
	}
	if len(names) != 2 || names[0] != "GLOBAL" || names[1] != "LOCAL" {
		t.Errorf("IndexOptions = %v, want [GLOBAL LOCAL]", names)
	}
}

func TestAlterTableVectorIndexTree(t *testing.T) {
	cmd := alterOneCmd(t, "ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS idx USING HNSW ((VEC_COSINE_DISTANCE(a))) COMMENT 'a'")
	if cmd.Type != ast.ATAddConstraint {
		t.Fatalf("Type = %d, want ATAddConstraint", cmd.Type)
	}
	con := cmd.Constraint
	if con == nil {
		t.Fatal("Constraint = nil")
	}
	if con.Type != ast.ConstrVector {
		t.Errorf("Constraint.Type = %d, want ConstrVector", con.Type)
	}
	if !con.IfNotExists {
		t.Errorf("IfNotExists = false, want true")
	}
	if con.Name != "idx" {
		t.Errorf("Name = %q, want \"idx\"", con.Name)
	}
	if con.IndexType != "HNSW" {
		t.Errorf("IndexType = %q, want \"HNSW\"", con.IndexType)
	}
	if len(con.IndexColumns) != 1 || con.IndexColumns[0].Expr == nil {
		t.Errorf("IndexColumns = %v, want one expression key part", con.IndexColumns)
	}
	if len(con.IndexOptions) != 1 || con.IndexOptions[0].Name != "COMMENT" {
		t.Errorf("IndexOptions = %v, want one COMMENT", con.IndexOptions)
	}
}

func TestAlterTableForeignKeyKeyPartsTree(t *testing.T) {
	cmd := alterOneCmd(t, "alter table d_n.t_n add foreign key ident (ident, ident(1) asc) references t_n match partial on delete cascade")
	con := cmd.Constraint
	if con == nil || con.Type != ast.ConstrForeignKey {
		t.Fatalf("Constraint = %+v, want a FOREIGN KEY", con)
	}
	if con.IndexName != "ident" {
		t.Errorf("IndexName = %q, want \"ident\"", con.IndexName)
	}
	if len(con.IndexColumns) != 2 {
		t.Fatalf("IndexColumns = %v, want 2", con.IndexColumns)
	}
	if con.IndexColumns[1].Length != 1 {
		t.Errorf("IndexColumns[1].Length = %d, want 1", con.IndexColumns[1].Length)
	}
	if got := con.Columns; len(got) != 2 || got[0] != "ident" || got[1] != "ident" {
		t.Errorf("Columns = %v, want the two names", got)
	}
	if con.Match != "PARTIAL" {
		t.Errorf("Match = %q, want \"PARTIAL\"", con.Match)
	}
	if con.OnDelete != ast.RefActCascade {
		t.Errorf("OnDelete = %d, want RefActCascade", con.OnDelete)
	}
}

func TestAlterTableUpdateIndexesTree(t *testing.T) {
	// The repartition closes the statement, so it is the last command.
	cmds := alterCmds(t, "alter table t add column a int partition by hash(a) update indexes (idx_a global, idx_b local)")
	if len(cmds) != 2 {
		t.Fatalf("%d commands, want 2", len(cmds))
	}
	cmd := cmds[1]
	if cmd.Type != ast.ATPartitionBy || cmd.PartitionBy == nil {
		t.Fatalf("Type = %d, PartitionBy = %v, want ATPartitionBy with a clause", cmd.Type, cmd.PartitionBy)
	}
	ui := cmd.PartitionBy.UpdateIndexes
	if len(ui) != 2 {
		t.Fatalf("UpdateIndexes = %v, want 2", ui)
	}
	if ui[0].Name != "idx_a" || !ui[0].Global {
		t.Errorf("UpdateIndexes[0] = %+v, want idx_a GLOBAL", ui[0])
	}
	if ui[1].Name != "idx_b" || ui[1].Global {
		t.Errorf("UpdateIndexes[1] = %+v, want idx_b LOCAL", ui[1])
	}
}

func TestAlterTableColumnConstraintCheckTree(t *testing.T) {
	cmd := alterOneCmd(t, "ALTER TABLE t ADD a2 int CONSTRAINT ident CHECK (a2 > 1) NOT ENFORCED")
	col := cmd.Column
	if col == nil {
		t.Fatal("Column = nil")
	}
	var check *ast.ColumnConstraint
	for _, cc := range col.Constraints {
		if cc.Type == ast.ColConstrCheck {
			check = cc
		}
	}
	if check == nil {
		t.Fatalf("no CHECK among %v", col.Constraints)
	}
	if check.Name != "ident" {
		t.Errorf("CHECK name = %q, want \"ident\"", check.Name)
	}
	if !check.NotEnforced {
		t.Errorf("NotEnforced = false, want true")
	}
	if check.Expr == nil {
		t.Errorf("Expr = nil")
	}
}

func TestAlterTablePartitionValuesDefaultTree(t *testing.T) {
	cmd := alterOneCmd(t, "ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (DEFAULT, 1))")
	if len(cmd.PartitionDefs) != 1 {
		t.Fatalf("PartitionDefs = %v, want 1", cmd.PartitionDefs)
	}
	list, ok := cmd.PartitionDefs[0].Values.(*ast.List)
	if !ok {
		t.Fatalf("Values = %T, want *ast.List", cmd.PartitionDefs[0].Values)
	}
	if len(list.Items) != 2 {
		t.Fatalf("Values = %v, want 2 items", list.Items)
	}
	if s, ok := list.Items[0].(*ast.String); !ok || s.Str != "DEFAULT" {
		t.Errorf("Values[0] = %#v, want the DEFAULT marker", list.Items[0])
	}

	// The bare DEFAULT spelling carries the same marker.
	cmd = alterOneCmd(t, "ALTER TABLE t ADD PARTITION (PARTITION p DEFAULT)")
	if s, ok := cmd.PartitionDefs[0].Values.(*ast.String); !ok || s.Str != "DEFAULT" {
		t.Errorf("PARTITION p DEFAULT: Values = %#v, want the DEFAULT marker", cmd.PartitionDefs[0].Values)
	}
}

func TestAlterTableAlterColumnSetDefaultTree(t *testing.T) {
	for _, c := range []struct {
		sql  string
		kind string
	}{
		{"ALTER TABLE t ALTER COLUMN a SET DEFAULT 1", "*ast.IntLit"},
		{"ALTER TABLE t ALTER COLUMN a SET DEFAULT 'x'", "*ast.StringLit"},
		{"ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL", "*ast.NullLit"},
		{"ALTER TABLE t ALTER COLUMN a SET DEFAULT TRUE", "*ast.BoolLit"},
		{"ALTER TABLE t ALTER COLUMN a SET DEFAULT -1", "*ast.UnaryExpr"},
		{"ALTER TABLE t ALTER COLUMN a SET DEFAULT (1+1)", "*ast.BinaryExpr"},
	} {
		cmd := alterOneCmd(t, c.sql)
		if cmd.Type != ast.ATAlterColumnDefault {
			t.Errorf("Parse(%q): Type = %d, want ATAlterColumnDefault", c.sql, cmd.Type)
		}
		if cmd.DefaultExpr == nil {
			t.Errorf("Parse(%q): DefaultExpr = nil", c.sql)
			continue
		}
		if got := fmt.Sprintf("%T", cmd.DefaultExpr); got != c.kind {
			t.Errorf("Parse(%q): DefaultExpr = %s, want %s", c.sql, got, c.kind)
		}
	}
}
