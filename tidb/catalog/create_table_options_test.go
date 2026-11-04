package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
	tidbparser "github.com/bytebase/omni/tidb/parser"
)

// What the CREATE TABLE forms this engine newly parses leave in the catalog.
//
// The corpus gate in tidb/parser records whether a statement parses and nothing
// about what it parsed into, so a form can go from reject to accept there and
// still misreport forever: a SQL_TSI_YEAR column that lands as an unknown type,
// an ENUM whose hex member comes out as the text `0x61`, a global temporary
// table recorded as an ordinary one. Each want below is TiDB v8.5.5's own SHOW
// CREATE TABLE, read on 2026-11-04.
//
// tidb/catalog's differential tests sit behind TIDB_CATALOG_PARITY, so this runs
// unconditionally instead.
func TestCreateTableOptionsCatalog(t *testing.T) {
	for _, tc := range []struct {
		name  string
		ddl   string
		table string
		want  string
	}{
		{
			// SQL_TSI_YEAR is YEAR, not an unknown type, and the numeric
			// modifiers YEAR now takes are not part of what it reports: the
			// server writes `year` for `year unsigned zerofill`.
			name:  "SQL_TSI_YEAR is the YEAR type",
			ddl:   "CREATE TABLE y (a SQL_TSI_YEAR, b SQL_TSI_YEAR(4), c YEAR ZEROFILL, d SQL_TSI_YEAR(4) UNSIGNED ZEROFILL)",
			table: "y",
			want: "CREATE TABLE `y` (\n" +
				"  `a` year DEFAULT NULL,\n" +
				"  `b` year DEFAULT NULL,\n" +
				"  `c` year DEFAULT NULL,\n" +
				"  `d` year DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin",
		},
		{
			// A hex or bit member is the bytes it holds, so 0x61 is the member
			// `a` and 0x6162 is `ab` — not the text that wrote them.
			name:  "an ENUM or SET member written as a literal is its bytes",
			ddl:   "CREATE TABLE e (c1 ENUM(0x61, 'b'), c2 SET(0b01100001, 'b'), c3 ENUM(x'61'), c4 ENUM(b'01100001'), c5 ENUM(0x6162))",
			table: "e",
			want: "CREATE TABLE `e` (\n" +
				"  `c1` enum('a','b') COLLATE utf8mb4_bin DEFAULT NULL,\n" +
				"  `c2` set('a','b') COLLATE utf8mb4_bin DEFAULT NULL,\n" +
				"  `c3` enum('a') COLLATE utf8mb4_bin DEFAULT NULL,\n" +
				"  `c4` enum('a') COLLATE utf8mb4_bin DEFAULT NULL,\n" +
				"  `c5` enum('ab') COLLATE utf8mb4_bin DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin",
		},
		{
			// The bytes a literal member holds, not the text that wrote them,
			// and the two spellings are not read alike: an unquoted literal
			// pads its digits to whole bytes while a quoted hex one must
			// already fill them. Reading the contents of a quoted hex literal
			// as though `0b` were a prefix decoded c1 as the member `a`.
			name:  "a literal member keeps every byte it holds",
			ddl:   "CREATE TABLE b (c1 ENUM(x'0b61'), c2 ENUM(x'0b'), c3 ENUM(x''), c4 ENUM(b'1'), c5 ENUM(0x6))",
			table: "b",
			want: "CREATE TABLE `b` (\n" +
				"  `c1` enum('\v" + "a') COLLATE utf8mb4_bin DEFAULT NULL,\n" +
				"  `c2` enum('\v') COLLATE utf8mb4_bin DEFAULT NULL,\n" +
				"  `c3` enum('') COLLATE utf8mb4_bin DEFAULT NULL,\n" +
				"  `c4` enum('\x01') COLLATE utf8mb4_bin DEFAULT NULL,\n" +
				"  `c5` enum('\x06') COLLATE utf8mb4_bin DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin",
		},
		{
			// SERIAL DEFAULT VALUE is NOT NULL AUTO_INCREMENT UNIQUE, so the
			// column is not null and auto-incrementing and the table gains the
			// unique key — which is also what makes it the table's implicit
			// primary key candidate.
			name:  "SERIAL DEFAULT VALUE is NOT NULL AUTO_INCREMENT UNIQUE",
			ddl:   "CREATE TABLE s (a INT SERIAL DEFAULT VALUE)",
			table: "s",
			want: "CREATE TABLE `s` (\n" +
				"  `a` int NOT NULL AUTO_INCREMENT,\n" +
				"  UNIQUE KEY `a` (`a`)\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin",
		},
		{
			// A global temporary table is not an ordinary one and not an
			// ordinary temporary one: both words come back, and so does the ON
			// COMMIT clause, last on the options line.
			name:  "a global temporary table reports as one",
			ddl:   "CREATE GLOBAL TEMPORARY TABLE g (a INT) ON COMMIT DELETE ROWS",
			table: "g",
			want: "CREATE GLOBAL TEMPORARY TABLE `g` (\n" +
				"  `a` int DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ON COMMIT DELETE ROWS",
		},
		{
			name:  "ON COMMIT comes after the other options",
			ddl:   "CREATE GLOBAL TEMPORARY TABLE g (a INT) COMMENT 'hi' ON COMMIT DELETE ROWS",
			table: "g",
			want: "CREATE GLOBAL TEMPORARY TABLE `g` (\n" +
				"  `a` int DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='hi' ON COMMIT DELETE ROWS",
		},
		{
			// ON COMMIT goes last, after the partition clause, because that is
			// where the grammar takes it. Writing it before the partitions gave
			// a definition this parser then refused, which the re-parse below
			// is here to catch. The server refuses the two together with 1562,
			// so this state is one only the catalog reaches -- it does not
			// model 1562 for an ordinary temporary table either -- but the
			// order still has to be the grammar's.
			name:  "ON COMMIT comes after the partition clause",
			ddl:   "CREATE GLOBAL TEMPORARY TABLE g (a INT) PARTITION BY HASH(a) PARTITIONS 10 ON COMMIT DELETE ROWS",
			table: "g",
			want: "CREATE GLOBAL TEMPORARY TABLE `g` (\n" +
				"  `a` int DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin\n" +
				"/*!50100 PARTITION BY HASH (`a`)\n" +
				"PARTITIONS 10 */ ON COMMIT DELETE ROWS",
		},
		{
			// PRESERVE ROWS is grammar v8.5.5 parses and then refuses with
			// 8200, so the catalog keeps the clause it was given rather than
			// rewriting it to the one TiDB can run.
			name:  "ON COMMIT PRESERVE ROWS is kept as written",
			ddl:   "CREATE GLOBAL TEMPORARY TABLE g (a INT) ON COMMIT PRESERVE ROWS",
			table: "g",
			want: "CREATE GLOBAL TEMPORARY TABLE `g` (\n" +
				"  `a` int DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ON COMMIT PRESERVE ROWS",
		},
		{
			// An ordinary temporary table is untouched by the clause's arrival.
			name:  "an ordinary temporary table still reports as one",
			ddl:   "CREATE TEMPORARY TABLE tt (a INT)",
			table: "tt",
			want: "CREATE TEMPORARY TABLE `tt` (\n" +
				"  `a` int DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin",
		},
		{
			// A name that begins with a digit is the name, backquoted on the
			// way out whether or not it was on the way in.
			name:  "a digit-leading name is a name",
			ddl:   "CREATE TABLE `123` (123a1 INT, 1_2 INT, 12e INT)",
			table: "123",
			want: "CREATE TABLE `123` (\n" +
				"  `123a1` int DEFAULT NULL,\n" +
				"  `1_2` int DEFAULT NULL,\n" +
				"  `12e` int DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin",
		},
		{
			// TiDB prints none of the statistics options, nor CHECKSUM or
			// TABLE_CHECKSUM, so a table that carries them reports as a table
			// that does not — which is what the catalog must do as well, and
			// what it would stop doing if someone added a case for one.
			name:  "the statistics options and the checksum do not report",
			ddl:   "CREATE TABLE st (a INT) STATS_BUCKETS=1 STATS_TOPN=2 STATS_SAMPLE_RATE=0.5 STATS_COL_CHOICE='all' STATS_COL_LIST='a' STATS_AUTO_RECALC=1 TABLE_CHECKSUM=1",
			table: "st",
			want: "CREATE TABLE `st` (\n" +
				"  `a` int DEFAULT NULL\n" +
				") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := wtSetup(t)
			wtExec(t, c, tc.ddl)
			got := c.ShowCreateTable("testdb", tc.table)
			if got != tc.want {
				t.Errorf("SHOW CREATE TABLE after %s:\ngot:\n%s\nwant:\n%s", tc.ddl, got, tc.want)
			}
			// What the catalog writes has to be something it can read back.
			// Rendering a clause in a position the grammar does not take it is
			// invisible to a comparison against a string someone wrote by hand,
			// and this is what caught ON COMMIT landing ahead of the partitions.
			if _, err := tidbparser.Parse(got); err != nil {
				t.Errorf("the definition this catalog wrote does not parse: %v\n%s", err, got)
			}
		})
	}
}

// TestGlobalTemporaryTableAST pins the two fields the statement carries, since
// SHOW CREATE TABLE reads them through the catalog and a reader of the AST —
// Bytebase's adapter, review, deparse — reads them directly.
func TestGlobalTemporaryTableAST(t *testing.T) {
	for _, tc := range []struct {
		sql          string
		wantGlobal   bool
		wantTemp     bool
		wantOnCommit string
	}{
		// Temporary stays false: v8.5.5 lists a global temporary table in SHOW
		// TABLES as a BASE TABLE, refuses a base table of the same name with
		// 1050, drops it with a plain DROP TABLE and answers 1051 to DROP
		// TEMPORARY TABLE for it, so it is not in the temporary namespace.
		{"CREATE GLOBAL TEMPORARY TABLE t (a INT) ON COMMIT DELETE ROWS", true, false, "DELETE ROWS"},
		{"CREATE GLOBAL TEMPORARY TABLE t (a INT) ON COMMIT PRESERVE ROWS", true, false, "PRESERVE ROWS"},
		{"CREATE GLOBAL TEMPORARY TABLE t LIKE t1 ON COMMIT DELETE ROWS", true, false, "DELETE ROWS"},
		{"CREATE TEMPORARY TABLE t (a INT)", false, true, ""},
		{"CREATE TABLE t (a INT)", false, false, ""},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseCreateTableForTest(t, tc.sql)
			if stmt.Global != tc.wantGlobal {
				t.Errorf("Global = %v, want %v", stmt.Global, tc.wantGlobal)
			}
			if stmt.Temporary != tc.wantTemp {
				t.Errorf("Temporary = %v, want %v", stmt.Temporary, tc.wantTemp)
			}
			if stmt.OnCommit != tc.wantOnCommit {
				t.Errorf("OnCommit = %q, want %q", stmt.OnCommit, tc.wantOnCommit)
			}
		})
	}
}

// TestDropGlobalTemporaryTableAST pins DROP's half of the same pair.
func TestDropGlobalTemporaryTableAST(t *testing.T) {
	for _, tc := range []struct {
		sql        string
		wantGlobal bool
		wantTemp   bool
	}{
		{"DROP GLOBAL TEMPORARY TABLE t", true, false},
		{"DROP TEMPORARY TABLE t", false, true},
		{"DROP TABLE t", false, false},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseDropTableForTest(t, tc.sql)
			if stmt.Global != tc.wantGlobal {
				t.Errorf("Global = %v, want %v", stmt.Global, tc.wantGlobal)
			}
			if stmt.Temporary != tc.wantTemp {
				t.Errorf("Temporary = %v, want %v", stmt.Temporary, tc.wantTemp)
			}
		})
	}
}

// TestSerialDefaultValueAST pins what SERIAL DEFAULT VALUE puts on the column,
// which is the three attributes it is another name for rather than an attribute
// of its own.
func TestSerialDefaultValueAST(t *testing.T) {
	stmt := parseCreateTableForTest(t, "CREATE TABLE t (a INT SERIAL DEFAULT VALUE)")
	if len(stmt.Columns) != 1 {
		t.Fatalf("got %d columns, want 1", len(stmt.Columns))
	}
	col := stmt.Columns[0]
	if !col.AutoIncrement {
		t.Error("AutoIncrement = false, want true")
	}
	var notNull, autoInc, unique int
	for _, con := range col.Constraints {
		switch con.Type {
		case nodes.ColConstrNotNull:
			notNull++
		case nodes.ColConstrAutoIncrement:
			autoInc++
		case nodes.ColConstrUnique:
			unique++
		}
	}
	if notNull != 1 || autoInc != 1 || unique != 1 {
		t.Errorf("constraints: NOT NULL=%d AUTO_INCREMENT=%d UNIQUE=%d, want 1 each",
			notNull, autoInc, unique)
	}
}

// TestEnumLiteralMembersAST pins that the members reach the AST decoded, so a
// consumer that is not the catalog — deparse, review — sees the same bytes the
// engine stored.
func TestEnumLiteralMembersAST(t *testing.T) {
	stmt := parseCreateTableForTest(t,
		"CREATE TABLE t (c1 ENUM(0x61, 'b'), c2 SET(0b01100001), c3 ENUM(x'6162'), "+
			"c4 ENUM(b'0110000101100010'), c5 ENUM(x'0b61'), c6 ENUM(0x6), c7 ENUM(b'1'), c8 ENUM(x''))")
	want := [][]string{
		{"a", "b"}, {"a"}, {"ab"}, {"ab"},
		// The contents of a quoted hex literal may begin with `0b`, which is
		// two bytes and not a prefix; an unquoted literal pads its digits to
		// whole bytes; a quoted literal with nothing in it is the empty string.
		{"\x0ba"}, {"\x06"}, {"\x01"}, {""},
	}
	if len(stmt.Columns) != len(want) {
		t.Fatalf("got %d columns, want %d", len(stmt.Columns), len(want))
	}
	for i, col := range stmt.Columns {
		got := col.TypeName.EnumValues
		if len(got) != len(want[i]) {
			t.Errorf("column %d: members %q, want %q", i, got, want[i])
			continue
		}
		for j := range got {
			if got[j] != want[i][j] {
				t.Errorf("column %d member %d = %q, want %q", i, j, got[j], want[i][j])
			}
		}
	}
}

// parseCreateTableForTest parses one CREATE TABLE and returns its statement.
func parseCreateTableForTest(t *testing.T, sql string) *nodes.CreateTableStmt {
	t.Helper()
	list, err := tidbparser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) gave %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*nodes.CreateTableStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave %T, want *nodes.CreateTableStmt", sql, list.Items[0])
	}
	return stmt
}

// parseDropTableForTest parses one DROP TABLE and returns its statement.
func parseDropTableForTest(t *testing.T, sql string) *nodes.DropTableStmt {
	t.Helper()
	list, err := tidbparser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) gave %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*nodes.DropTableStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave %T, want *nodes.DropTableStmt", sql, list.Items[0])
	}
	return stmt
}
