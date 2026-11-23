package driver

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// aliasFunctionBodies are alias function bodies over x and y, each a shape the
// server renders differently from how it is written. TestDumpAliasFunctionRestores
// restores each one from the dump and requires the same statement back.
var aliasFunctionBodies = []string{
	"abs(x - y)",
	"abs((x - y))",
	"abs(((x - y)))",
	"abs(x + y - 2)",
	"abs(x - (y - 2))",
	"abs(x * 2 + y)",
	"abs(x / y)",
	"abs(x % y)",
	"abs(x DIV y)",
	"abs(x & 1 | y ^ 3)",
	"abs(~(x + y))",
	"abs(-x + y)",
	"abs(-(x - y))",
	"abs(- -x + y)",
	"abs(-x * y)",
	"abs(-1 * x + y)",
	"abs((-1) * x + y)",
	"abs(x - -y)",
	"abs(x * -1 + y)",
	"abs(x + - -1 + y)",
	"abs(CAST(x AS BIGINT) + y)",
	"abs(abs(x - 1) - y)",
	"abs(x + 1.50 + y)",
	"if(x + 1 > y, x - y, -y)",
	"if(x & 1 = 1, y, 0)",
	"if((x > 1) AND (y < 5 OR x = y), 1, 0)",
	"if(NOT x > y, 1, 0)",
	"if(x BETWEEN y - 1 AND y + 1, 1, 0)",
	"if(x IN (y - 1, y + 1), 1, 0)",
	"if(x IS NULL, -y, y)",
	"coalesce(CASE x WHEN y - 1 THEN 'a' END, 'b')",
	`concat('a''b', x, y)`,
	`concat('back\\slash', x, y)`,
	`concat('trailing\\', x, y)`,
	`concat('both\\''quote', x, y)`,
	`concat("dq'", x, y)`,
	"concat('tab\tand -- dashes', x, y)",
}

// TestDumpAliasFunctionRestores creates an alias function for every body,
// dumps the database, restores the dump into it emptied, and requires the
// second dump to be the first. The server states an alias function's body as
// it rendered it, and replaying that rendering renders it again, so the dump
// has to state a source that renders to it.
func TestDumpAliasFunctionRestores(t *testing.T) {
	ctx := context.Background()
	d := openTestDriver(t, "omni_dump_alias_function")
	var script strings.Builder
	for i, body := range aliasFunctionBodies {
		fmt.Fprintf(&script, "CREATE ALIAS FUNCTION alias_%02d(INT, INT) WITH PARAMETER(x, y) AS %s;\n", i, body)
	}
	execDumpScript(t, d.db, script.String())
	first := dumpScript(t, d)
	for i := range aliasFunctionBodies {
		if want := fmt.Sprintf("-- Function structure for alias_%02d(int, int)", i); !strings.Contains(first, want) {
			t.Errorf("the dump does not state %q:\n%s", want, first)
		}
	}

	for _, stmt := range []string{
		"DROP DATABASE " + quoteBacktick(d.databaseName),
		"CREATE DATABASE " + quoteBacktick(d.databaseName),
		"USE " + quoteBacktick(d.databaseName),
	} {
		if _, err := d.db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	execDumpScript(t, d.db, first)
	if diff := cmp.Diff(first, dumpScript(t, d)); diff != "" {
		t.Errorf("restoring the dump changes it (-dumped +restored):\n%s", diff)
	}
}

// TestAliasFunctionSource holds the source the dump states for what the server
// states, case by case; TestDumpAliasFunctionRestores holds it to the server.
func TestAliasFunctionSource(t *testing.T) {
	const head = "CREATE ALIAS FUNCTION f(int, int) WITH PARAMETER(x, y) AS "
	tests := []struct {
		name, stmt, want string
	}{
		{name: "arithmetic", stmt: head + "abs((`x` - `y`))", want: head + "abs(`x` - `y`)"},
		{name: "written parentheses", stmt: head + "abs(((`x` - `y`)))", want: head + "abs((`x` - `y`))"},
		{name: "nested arithmetic", stmt: head + "abs((((`x` & 1) | `y`) ^ 3))", want: head + "abs(`x` & 1 | `y` ^ 3)"},
		{name: "negation", stmt: head + "abs(((-1 * `x`) + `y`))", want: head + "abs(-`x` + `y`)"},
		{name: "negated parentheses", stmt: head + "abs((-1 * ((`x` - `y`))))", want: head + "abs(-(`x` - `y`))"},
		{name: "double negation", stmt: head + "abs((-1 * (-1 * `x`)))", want: head + "abs(- -`x`)"},
		{name: "negative literal", stmt: head + "abs((`x` - -1))", want: head + "abs(`x` - -1)"},
		{name: "comparison", stmt: head + "if(((`x` + 1) > `y`), 1, 0)", want: head + "if((`x` + 1 > `y`), 1, 0)"},
		{name: "backslash", stmt: head + `concat('a\b''c\', ` + "`x`, `y`)", want: head + `concat('a\\b''c\\', ` + "`x`, `y`)"},
		{name: "global", stmt: "CREATE GLOBAL ALIAS FUNCTION g(int) WITH PARAMETER(x) AS abs((`x` - 1))", want: "CREATE GLOBAL ALIAS FUNCTION g(int) WITH PARAMETER(x) AS abs(`x` - 1)"},
		{
			name: "not an alias",
			stmt: `CREATE FUNCTION j(int) RETURNS int PROPERTIES ("file"="file:///a\b.jar", "symbol"="J", "type"="JAVA_UDF")`,
			want: `CREATE FUNCTION j(int) RETURNS int PROPERTIES ("file"="file:///a\b.jar", "symbol"="J", "type"="JAVA_UDF")`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := aliasFunctionSource(test.stmt)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("aliasFunctionSource(%q)\n got %q\nwant %q", test.stmt, got, test.want)
			}
		})
	}
}
