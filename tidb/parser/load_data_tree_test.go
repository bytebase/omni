package parser

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// What LOAD DATA and IMPORT INTO parse into, option by option.
//
// The corpus records only whether a statement parsed, so it would stay green
// on a reader that walked past every option it found. These cases pin the
// value each option carries instead: the bytes a hex or bit literal stands
// for, the string DEFINED NULL BY holds and the OPTIONALLY ENCLOSED that may
// follow it, FORMAT's name, each WITH option's name and value, and -- the
// one a corpus can never see -- that a user variable in a column list is a
// variable and not a column whose name begins with '@'.

func TestLoadDataTreeFieldTerminators(t *testing.T) {
	// X'6B6B' and B'110101101101011' are both the two bytes "kk": a bit
	// literal is padded on the left to fill whole bytes.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS TERMINATED BY X'6B6B'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_terminated "kk"}`)
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS TERMINATED BY B'110101101101011'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_terminated "kk"}`)
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS TERMINATED BY 0x6B6B",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_terminated "kk"}`)
	// Every option of a FIELDS and a LINES clause takes the same value.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS TERMINATED BY X'09' ENCLOSED BY X'22' ESCAPED BY X'5C' LINES STARTING BY X'3E' TERMINATED BY X'0A'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :lines_starting ">" :lines_terminated "\n" :fields_terminated "\t" :fields_enclosed "\"" :fields_escaped "\\"}`)
	// OPTIONALLY sets the flag as well as the value, and a plain ENCLOSED
	// item after it clears the flag again: v8.5.5 writes `"7"` for `INTO
	// OUTFILE ... FIELDS OPTIONALLY ENCLOSED BY '"' ENCLOSED BY '"'` and
	// `7` with the OPTIONALLY item alone. MySQL and MariaDB keep the flag,
	// which their own copy of this test records.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS OPTIONALLY ENCLOSED BY B'100010'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_enclosed "\"" :fields_optionally_enclosed true}`)
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS OPTIONALLY ENCLOSED BY 'a' ENCLOSED BY 'b'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_enclosed "b"}`)
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS ENCLOSED BY 'a' OPTIONALLY ENCLOSED BY 'b'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_enclosed "b" :fields_optionally_enclosed true}`)
}

func TestLoadDataTreeDefinedNullBy(t *testing.T) {
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS DEFINED NULL BY 'nil'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_defined_null "nil"}`)
	// A hex value is decoded here too: X'00' is one NUL byte.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS DEFINED NULL BY X'00'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_defined_null "\x00"}`)
	// DEFINED NULL BY '' is a written item, which is why the node carries a
	// flag beside the string.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS DEFINED NULL BY ''",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_defined_null ""}`)
	// The OPTIONALLY ENCLOSED tail belongs to this item and sets no
	// enclosing character.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS DEFINED NULL BY 'NULL' OPTIONALLY ENCLOSED",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_defined_null "NULL" :defined_null_optionally_enclosed true}`)
	// With a BY after it, the same words are the next item instead: the
	// tail is unset and ENCLOSED BY carries 'b'.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS DEFINED NULL BY 'a' OPTIONALLY ENCLOSED BY 'b'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_enclosed "b" :fields_optionally_enclosed true :fields_defined_null "a"}`)
	// A repeat leaves the last value, which is what the engine's own AST
	// keeps.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS DEFINED NULL BY 'a' DEFINED NULL BY 'b'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_defined_null "b"}`)
	// The tail belongs to its own item, so a repeat without it clears it.
	// Measured by loading a file on v8.5.5: with the tail an enclosed
	// `"a"` reads as NULL, and after this repeat it reads as the string.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS DEFINED NULL BY 'a' OPTIONALLY ENCLOSED DEFINED NULL BY 'b'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_defined_null "b"}`)
}

func TestLoadDataTreeFormat(t *testing.T) {
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' FORMAT 'delimited data' INTO TABLE t",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :format "delimited data" :table {TABLEREF :loc 61 :name t}}`)
	ParseAndCompare(t, "LOAD DATA INFILE '/x.sql' FORMAT 'sql file' REPLACE INTO TABLE t",
		`{LOAD_DATA :loc 0 :infile "/x.sql" :format "sql file" :replace true :table {TABLEREF :loc 63 :name t}}`)
	// FORMAT '' is a written clause, which is why the node carries a flag.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' FORMAT '' INTO TABLE t",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :format "" :table {TABLEREF :loc 47 :name t}}`)
}

func TestLoadDataTreeWithOptions(t *testing.T) {
	// A bare option has a name and no value; the name is kept as written.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t WITH detached",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :with {OPT :loc 44 :name detached}}`)
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t WITH threads=10, detached",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :with {OPT :loc 44 :name threads :val {INT_LIT :val 10 :loc 52}} {OPT :loc 56 :name detached}}`)
	// A signed value keeps its sign, and a string, a byte literal and NULL
	// each reach the node as themselves.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t WITH threads=-1",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :with {OPT :loc 44 :name threads :val {UNARY :op - :loc 52 :operand {INT_LIT :val 1 :loc 53}}}}`)
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t WITH character_set='utf8', threads=NULL",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :with {OPT :loc 44 :name character_set :val {STRING_LIT :val "utf8" :loc 58}} {OPT :loc 66 :name threads :val {NULL_LIT :loc 74}}}`)
}

func TestLoadDataTreeUserVariables(t *testing.T) {
	// @var1 is a variable, not a column named "@var1". A corpus that
	// records only whether this parsed cannot tell the two apart.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t (column1, @var1) SET column2 = @var1/100",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :columns {COLREF :loc 40 :col column1} {VAR :name var1 :loc 49} :set {ASSIGN :loc 60 :col {COLREF :loc 60 :col column2} :val {BINEXPR :op / :loc 70 :left {VAR :name var1 :loc 70} :right {INT_LIT :val 100 :loc 76}}}}`)
	ParseAndCompare(t, "IMPORT INTO t (a,@1) FROM '/file.csv'",
		`{IMPORT_INTO :loc 0 :table {TABLEREF :loc 12 :name t} :columns {COLREF :loc 15 :col a} {VAR :name 1 :loc 17} :path "/file.csv"}`)
	// A backquoted name is a column, even one spelled with an '@'. The
	// lexer strips the quotes, so only the source text tells the two
	// apart -- and getting it wrong would also reject the IMPORT INTO
	// below, which v8.5.5 takes.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t (`@v`)",
		"{LOAD_DATA :loc 0 :infile \"/x.csv\" :table {TABLEREF :loc 37 :name t} :columns {COLREF :loc 40 :col @v}}")
	ParseAndCompare(t, "IMPORT INTO t (`@v`) FROM SELECT 1",
		"{IMPORT_INTO :loc 0 :table {TABLEREF :loc 12 :name t} :columns {COLREF :loc 15 :col @v} :from {SELECT :loc 26 :targets ({INT_LIT :val 1 :loc 33})}}")
}

func TestImportIntoTree(t *testing.T) {
	ParseAndCompare(t, "IMPORT INTO t FROM '/file.csv'",
		`{IMPORT_INTO :loc 0 :table {TABLEREF :loc 12 :name t} :path "/file.csv"}`)
	ParseAndCompare(t, "IMPORT INTO d.t FROM '/file.csv' FORMAT 'sql file' WITH thread=1, detached",
		`{IMPORT_INTO :loc 0 :table {TABLEREF :loc 12 :schema d :name t} :path "/file.csv" :format "sql file" :with {OPT :loc 56 :name thread :val {INT_LIT :val 1 :loc 63}} {OPT :loc 66 :name detached}}`)
	// An empty pair of parentheses is a written list with no entries, so the
	// deparse shows no :columns.
	ParseAndCompare(t, "IMPORT INTO t () FROM '/file.csv'",
		`{IMPORT_INTO :loc 0 :table {TABLEREF :loc 12 :name t} :path "/file.csv"}`)
	ParseAndCompare(t, "IMPORT INTO t (a,b) SET c=1 FROM '/file.csv'",
		`{IMPORT_INTO :loc 0 :table {TABLEREF :loc 12 :name t} :columns {COLREF :loc 15 :col a} {COLREF :loc 17 :col b} :set {ASSIGN :loc 24 :col {COLREF :loc 24 :col c} :val {INT_LIT :val 1 :loc 26}} :path "/file.csv"}`)
}

// TestImportIntoTreeQuerySource pins that the query source lands on the node
// as a query and not as a path, which is what tells the two forms apart.
func TestImportIntoTreeQuerySource(t *testing.T) {
	for _, sql := range []string{
		"IMPORT INTO t FROM SELECT * FROM s",
		"IMPORT INTO t FROM (SELECT * FROM s)",
		"IMPORT INTO t FROM WITH cte AS (SELECT 1) SELECT * FROM cte",
	} {
		t.Run(sql, func(t *testing.T) {
			stmt := parseImportInto(t, sql)
			if stmt.Select == nil {
				t.Fatalf("Select is unset for %q", sql)
			}
			if stmt.Path != "" {
				t.Errorf("Path = %q, want empty: the source is a query", stmt.Path)
			}
		})
	}
	// A TABLE or VALUES target is the query expression it begins, so it
	// arrives in Select over its own source, as it does everywhere else.
	stmt := parseImportInto(t, "IMPORT INTO t FROM TABLE s WITH thread=2")
	if stmt.Select == nil {
		t.Fatal("Select is unset for FROM TABLE s")
	}
	if stmt.Select.TableSource == nil || stmt.Select.TableSource.Table == nil || stmt.Select.TableSource.Table.Name != "s" {
		t.Errorf("the source reads %#v, want the TABLE primary s", stmt.Select.TableSource)
	}
	if stmt.Path != "" {
		t.Errorf("Path = %q, want empty: the source is a query", stmt.Path)
	}
	// The clauses a bare TableStmt could not hold, each accepted on v8.5.5.
	for _, c := range []struct {
		sql   string
		check func(*testing.T, *ast.SelectStmt)
	}{
		{"IMPORT INTO t FROM TABLE t1 UNION TABLE t2", func(t *testing.T, sel *ast.SelectStmt) {
			if sel.SetOp != ast.SetOpUnion {
				t.Errorf("SetOp = %v, want UNION", sel.SetOp)
			}
		}},
		{"IMPORT INTO t FROM TABLE t1 ORDER BY a", func(t *testing.T, sel *ast.SelectStmt) {
			if len(sel.OrderBy) != 1 {
				t.Errorf("OrderBy has %d items, want 1", len(sel.OrderBy))
			}
		}},
		{"IMPORT INTO t FROM TABLE t1 LIMIT 1", func(t *testing.T, sel *ast.SelectStmt) {
			if sel.Limit == nil {
				t.Error("LIMIT was dropped")
			}
		}},
		{"IMPORT INTO t FROM VALUES ROW(1)", func(t *testing.T, sel *ast.SelectStmt) {
			if sel.ValuesSource == nil || len(sel.ValuesSource.Rows) != 1 {
				t.Errorf("the source reads %#v, want the VALUES primary with 1 row", sel.ValuesSource)
			}
		}},
	} {
		got := parseImportInto(t, c.sql)
		if got.Select == nil {
			t.Fatalf("Select is unset for %q", c.sql)
		}
		c.check(t, got.Select)
	}
	if len(stmt.Options) != 1 || stmt.Options[0].Name != "thread" {
		t.Errorf("Options = %v, want one named thread", stmt.Options)
	}
}

// parseImportInto parses sql and asserts it is one IMPORT INTO.
func parseImportInto(t *testing.T, sql string) *ast.ImportIntoStmt {
	t.Helper()
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", sql, result.Len())
	}
	stmt, ok := result.Items[0].(*ast.ImportIntoStmt)
	if !ok {
		t.Fatalf("Parse(%q) returned %T, want *ast.ImportIntoStmt", sql, result.Items[0])
	}
	return stmt
}
