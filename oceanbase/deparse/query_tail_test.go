package deparse

import (
	"regexp"
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

var queryTailLocFree = regexp.MustCompile(` :loc -?[0-9]+`)

// queryTail prints the locking clauses and INTOs of a query expression, leaving out
// the select lists, where a deparse adds its aliases.
func queryTail(s *ast.SelectStmt) string {
	if s == nil {
		return ""
	}
	out := ""
	if s.ForUpdate != nil {
		out += ast.NodeToString(s.ForUpdate)
	}
	if s.Into != nil {
		out += ast.NodeToString(s.Into)
	}
	return out + "(" + queryTail(s.Left) + ")(" + queryTail(s.Right) + ")(" + queryTail(s.ParenSource) + ")"
}

// TestDeparseKeepsQueryTail: the locking clause and INTO after a query expression
// survive a deparse, on every shape that holds them, and what comes out parses back to
// the same clauses.
func TestDeparseKeepsQueryTail(t *testing.T) {
	for _, c := range []struct{ input, expected string }{
		// VALUES
		{"VALUES ROW(1) INTO @x", "values row(1) into @`x`"},
		{"VALUES ROW(1) INTO @x, @y", "values row(1) into @`x`,@`y`"},
		{"VALUES ROW(1) INTO OUTFILE 'f'", "values row(1) into outfile 'f'"},
		{"VALUES ROW(1) INTO DUMPFILE 'f'", "values row(1) into dumpfile 'f'"},
		{"VALUES ROW(1) FOR UPDATE", "values row(1) for update"},
		{"VALUES ROW(1) LOCK IN SHARE MODE", "values row(1) lock in share mode"},
		{"VALUES ROW(1) FOR UPDATE NOWAIT", "values row(1) for update nowait"},
		{"VALUES ROW(1) ORDER BY 1 LIMIT 1 FOR UPDATE INTO @x", "values row(1) order by 1 limit 1 for update into @`x`"},
		// A parenthesized query, with an ORDER BY or LIMIT before the clauses.
		{"(SELECT 1) LIMIT 1 FOR UPDATE", "(select 1 AS `1`) limit 1 for update"},
		{"(SELECT 1) ORDER BY 1 INTO @x", "(select 1 AS `1`) order by 1 into @`x`"},
		{"(VALUES ROW(1)) LIMIT 1 FOR UPDATE INTO @x", "(values row(1)) limit 1 for update into @`x`"},
		// A set operation.
		{"(SELECT 1) UNION (SELECT 2) FOR UPDATE", "(select 1 AS `1`) union (select 2 AS `2`) for update"},
		{"(SELECT 1) UNION (SELECT 2) INTO @x", "(select 1 AS `1`) union (select 2 AS `2`) into @`x`"},
		{"VALUES ROW(1) UNION VALUES ROW(2) INTO @x", "values row(1) union values row(2) into @`x`"},
		{"VALUES ROW(1) UNION VALUES ROW(2) FOR UPDATE", "values row(1) union values row(2) for update"},
		// A SELECT.
		{"SELECT 1 INTO @x", "select 1 AS `1` into @`x`"},
		{"SELECT a FROM t FOR UPDATE INTO @x", "select `a` AS `a` from `t` for update into @`x`"},
		// Every target of INTO, and every option of OUTFILE.
		{"SELECT 1 INTO local_v", "select 1 AS `1` into `local_v`"},
		{"SELECT 1 INTO @x, local_v", "select 1 AS `1` into @`x`,`local_v`"},
		{
			"VALUES ROW(1) INTO OUTFILE 'f' CHARACTER SET utf8mb4 FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '\"' ESCAPED BY '\\\\' LINES STARTING BY 'a' TERMINATED BY 'b'",
			"values row(1) into outfile 'f' character set utf8mb4 fields terminated by ',' optionally enclosed by '\"' escaped by '\\\\' lines starting by 'a' terminated by 'b'",
		},
		{
			"VALUES ROW(1) INTO OUTFILE S3 's3://b/p' FORMAT CSV HEADER FIELDS TERMINATED BY ',' LINES TERMINATED BY 'x' MANIFEST ON OVERWRITE OFF ENCRYPTION SSE_KMS 'k'",
			"values row(1) into outfile s3 's3://b/p' format csv header fields terminated by ',' lines terminated by 'x' manifest on overwrite off encryption sse_kms 'k'",
		},
		{"VALUES ROW(1) INTO OUTFILE 'it''s'", "values row(1) into outfile 'it\\'s'"},
	} {
		t.Run(c.input, func(t *testing.T) {
			sel := parseSelect(t, c.input)
			got := DeparseSelect(sel)
			if got != c.expected {
				t.Fatalf("DeparseSelect(%q)\n got: %s\nwant: %s", c.input, got, c.expected)
			}
			again := parseSelect(t, got)
			was := queryTailLocFree.ReplaceAllString(queryTail(sel), "")
			now := queryTailLocFree.ReplaceAllString(queryTail(again), "")
			if was != now {
				t.Errorf("the clauses do not read back the same:\n was: %s\n now: %s", was, now)
			}
			// It holds its shape, too.
			if second := DeparseSelect(again); second != got {
				t.Errorf("deparse is not stable:\n first: %s\nsecond: %s", got, second)
			}
		})
	}
}

// TestParseMarksBareIntoVariable: only DECLAREd names are bare, and the deparse depends
// on telling them from @user variables.
func TestParseMarksBareIntoVariable(t *testing.T) {
	list, err := parser.Parse("SELECT 1 INTO @u, v")
	if err != nil {
		t.Fatal(err)
	}
	vars := list.Items[0].(*ast.SelectStmt).Into.Vars
	if len(vars) != 2 || vars[0].Bare || !vars[1].Bare {
		t.Errorf("Vars = %+v", vars)
	}
}
