package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// The storage, index and partition clauses OceanBase writes and the MySQL grammar
// this parser forked has no table option for.
// TestOceanBaseGrammarInventory holds these to a live OceanBase; this file pins what
// the parser builds from them, so a change that still parses but drops an option
// fails here.

func optionNames(t *testing.T, sql string) []string {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("%q: %v", sql, err)
	}
	ct, ok := list.Items[0].(*nodes.CreateTableStmt)
	if !ok {
		t.Fatalf("%q: not a CREATE TABLE", sql)
	}
	var names []string
	for _, o := range ct.Options {
		names = append(names, o.Name)
	}
	return names
}

func TestOceanBaseTableOptions(t *testing.T) {
	got := optionNames(t, "CREATE TABLE t (a INT PRIMARY KEY) ORGANIZATION INDEX "+
		"AUTO_INCREMENT_MODE = 'ORDER' REPLICA_NUM = 1 BLOCK_SIZE = 16384 "+
		"TABLET_SIZE = 134217728 PCTFREE = 0 USE_BLOOM_FILTER = FALSE "+
		"ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE")
	want := []string{"ORGANIZATION INDEX", "AUTO_INCREMENT_MODE", "REPLICA_NUM",
		"BLOCK_SIZE", "TABLET_SIZE", "PCTFREE", "USE_BLOOM_FILTER",
		"ENABLE_MACRO_BLOCK_BLOOM_FILTER"}
	if len(got) != len(want) {
		t.Fatalf("options = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("option %d = %q, want %q", i, got[i], want[i])
		}
	}
}

// DYNAMIC_PARTITION_POLICY takes a sub-list this parser keeps as written.
func TestOceanBaseDynamicPartitionPolicy(t *testing.T) {
	list, err := Parse("CREATE TABLE t (a INT) DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE, TIME_UNIT = 'DAY')")
	if err != nil {
		t.Fatal(err)
	}
	opts := list.Items[0].(*nodes.CreateTableStmt).Options
	if len(opts) != 1 || opts[0].Name != "DYNAMIC_PARTITION_POLICY" {
		t.Fatalf("options = %+v", opts)
	}
	if want := "(ENABLE = TRUE, TIME_UNIT = 'DAY')"; opts[0].Value != want {
		t.Errorf("value = %q, want %q", opts[0].Value, want)
	}
}

// TestOceanBaseRejectsMalformedDynamicPartitionPolicy pins that the policy's contents are
// parsed, not just counted for balanced parentheses. A parser that took any parenthesized
// text would pass DDL through review that the server then refuses.
func TestOceanBaseRejectsMalformedDynamicPartitionPolicy(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (a INT) DYNAMIC_PARTITION_POLICY = (ENABLE =)",
		"CREATE TABLE t (a INT) DYNAMIC_PARTITION_POLICY = (SELECT * FROM x)",
		"CREATE TABLE t (a INT) DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE,)",
		"CREATE TABLE t (a INT) DYNAMIC_PARTITION_POLICY = ()",
		"CREATE TABLE t (a INT) DYNAMIC_PARTITION_POLICY = (ENABLE = TRUE",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("parsed %q, but it is not a list of policy assignments", sql)
		}
	}
}

func TestOceanBaseIndexOptions(t *testing.T) {
	for _, tc := range []struct {
		name string
		sql  string
		want []string
	}{
		{"local", "CREATE TABLE t (a INT, KEY i (a) LOCAL)", []string{"LOCAL"}},
		{"global", "CREATE TABLE t (a INT, KEY i (a) GLOBAL)", []string{"GLOBAL"}},
		{"block size", "CREATE TABLE t (a INT, KEY i (a) BLOCK_SIZE 16384)", []string{"BLOCK_SIZE"}},
		{"block size then local", "CREATE TABLE t (a INT, KEY i (a) BLOCK_SIZE 16384 LOCAL)", []string{"BLOCK_SIZE", "LOCAL"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatal(err)
			}
			cons := list.Items[0].(*nodes.CreateTableStmt).Constraints
			if len(cons) != 1 {
				t.Fatalf("constraints = %d", len(cons))
			}
			var got []string
			for _, o := range cons[0].IndexOptions {
				got = append(got, o.Name)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("index options = %v, want %v", got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("option %d = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestOceanBaseSubpartitionTemplate(t *testing.T) {
	list, err := Parse("CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS(a) " +
		"SUBPARTITION BY KEY(b) SUBPARTITION TEMPLATE (SUBPARTITION p0, SUBPARTITION p1) " +
		"(PARTITION q0 VALUES LESS THAN (10))")
	if err != nil {
		t.Fatal(err)
	}
	part := list.Items[0].(*nodes.CreateTableStmt).Partitions
	if part == nil {
		t.Fatal("no partition clause")
	}
	if n := len(part.SubPartTemplate); n != 2 {
		t.Fatalf("subpartition template = %d entries, want 2", n)
	}
	if len(part.Partitions) != 1 {
		t.Errorf("partitions = %d, want 1", len(part.Partitions))
	}
}

// ALTER TABLE takes the index scope too.
func TestOceanBaseAlterIndexClauses(t *testing.T) {
	if _, err := Parse("ALTER TABLE t ADD UNIQUE INDEX uk (a) GLOBAL"); err != nil {
		t.Errorf("ADD UNIQUE INDEX ... GLOBAL: %v", err)
	}
}

// TestOceanBaseRejectsForeignGrammar keeps the parser half of foreignGrammar honest
// without an engine. The inventory asks OceanBase the same statements, but it skips
// without a container and CI never has one, so this is what guards against the parser
// growing permissive: accepting any of these would pass a statement through review that
// the server answers 1064 for. A customer reported the DROP INDEX form as rejected.
func TestOceanBaseRejectsForeignGrammar(t *testing.T) {
	if len(foreignGrammar) == 0 {
		t.Fatal("foreignGrammar is empty, so this proves nothing")
	}
	for _, c := range foreignGrammar {
		t.Run(c.name, func(t *testing.T) {
			if _, err := Parse(c.sql); err == nil {
				t.Errorf("parsed %q, but neither MySQL nor OceanBase takes it", c.sql)
			}
		})
	}
}
