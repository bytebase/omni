package catalog

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

// jsonQueryStoredForms are the clauses of JSON_QUERY this package has to print as the server
// stores them. Each has RETURNING JSON or ASIS, because OceanBase 4.3.5.6 answers 5443 for a
// JSON_QUERY with neither. DOT ON MISMATCH is not one: the server stores `t`.`j`.a for the
// whole call, which does not read back.
var jsonQueryStoredForms = []string{
	"RETURNING JSON", "RETURNING JSON TRUNCATE", "RETURNING JSON ALLOW SCALARS", "RETURNING JSON DISALLOW SCALARS",
	"RETURNING JSON PRETTY", "RETURNING JSON ASCII", "RETURNING JSON ASIS", "RETURNING JSON MULTIVALUE",
	"RETURNING JSON WITHOUT WRAPPER", "RETURNING JSON WITHOUT ARRAY WRAPPER", "RETURNING JSON WITH WRAPPER",
	"RETURNING JSON WITH ARRAY WRAPPER", "RETURNING JSON WITH UNCONDITIONAL WRAPPER", "RETURNING JSON WITH CONDITIONAL WRAPPER",
	"RETURNING JSON WITH UNCONDITIONAL ARRAY WRAPPER", "RETURNING JSON WITH CONDITIONAL ARRAY WRAPPER",
	"RETURNING JSON ERROR ON ERROR", "RETURNING JSON NULL ON ERROR", "RETURNING JSON EMPTY ON ERROR",
	"RETURNING JSON EMPTY ARRAY ON ERROR", "RETURNING JSON EMPTY OBJECT ON ERROR",
	"RETURNING JSON ERROR ON EMPTY", "RETURNING JSON NULL ON EMPTY", "RETURNING JSON EMPTY ON EMPTY",
	"RETURNING JSON EMPTY ARRAY ON EMPTY", "RETURNING JSON EMPTY OBJECT ON EMPTY",
	"RETURNING JSON ERROR ON MISMATCH", "RETURNING JSON NULL ON MISMATCH",
	"RETURNING JSON ERROR ON ERROR NULL ON EMPTY", "RETURNING JSON NULL ON EMPTY ERROR ON ERROR",
	"RETURNING JSON ERROR ON ERROR ERROR ON MISMATCH", "RETURNING JSON NULL ON EMPTY ERROR ON MISMATCH",
	"RETURNING JSON NULL ON EMPTY ERROR ON ERROR ERROR ON MISMATCH",
	"RETURNING JSON TRUNCATE ALLOW SCALARS PRETTY ASCII WITH CONDITIONAL ARRAY WRAPPER ASIS EMPTY OBJECT ON ERROR EMPTY ARRAY ON EMPTY ERROR ON MISMATCH MULTIVALUE",
	"ASIS", "ASIS ERROR ON MISMATCH", "ASIS MULTIVALUE",
	"RETURNING CHAR(10) ASIS", "RETURNING CHAR(10) CHARSET ascii ASIS", "RETURNING CHAR(10) BINARY ASIS",
	"RETURNING CHAR(0) ASIS", "RETURNING NCHAR(3) ASIS", "RETURNING BINARY(4) ASIS",
	"RETURNING SIGNED ASIS", "RETURNING SIGNED INTEGER ASIS", "RETURNING UNSIGNED INTEGER ASIS",
	"RETURNING DECIMAL ASIS", "RETURNING DECIMAL(5) ASIS", "RETURNING DECIMAL(5,2) ASIS", "RETURNING NUMBER(5,2) ASIS",
	"RETURNING FLOAT ASIS", "RETURNING FLOAT(5) ASIS", "RETURNING DOUBLE ASIS",
	"RETURNING DATE ASIS", "RETURNING DATETIME ASIS", "RETURNING DATETIME(3) ASIS", "RETURNING TIME ASIS", "RETURNING TIME(2) ASIS",
	// What a multi-valued index reads back as.
	"RETURNING UNSIGNED WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH",
}

// TestOracle_JSONQueryStoredForms: a generated column and a view hold JSON_QUERY the way
// OceanBase stores it, so no clause goes missing and what this package writes is the text
// the engine would.
func TestOracle_JSONQueryStoredForms(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	stored := regexp.MustCompile(`json_query\(.*\)`)

	for _, clauses := range jsonQueryStoredForms {
		t.Run(clauses, gapped(func(t testing.TB) {
			table := "CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_QUERY(a, '$.b' " + clauses + ")))"
			server, err := storeAndRead(ctx, o, "zz_json_query", table, "g")
			if err != nil {
				t.Fatalf("the server refused %q: %v", table, err)
			}
			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d; "+table, nil); err != nil {
				t.Fatal(err)
			}
			got, want := stored.FindString(c.ShowCreateTable("d", "g")), stored.FindString(server)
			if got != want {
				t.Errorf("generated column:\nomni:   %s\nserver: %s", got, want)
			}

			viewCall := "JSON_QUERY(a, '$.b' " + clauses + ")"
			if _, err := o.db.ExecContext(ctx, "CREATE VIEW v AS SELECT "+viewCall+" AS x FROM g"); err != nil {
				t.Fatal(err)
			}
			var name, serverView, charset, collation string
			row := o.db.QueryRowContext(ctx, "SHOW CREATE VIEW v")
			if err := row.Scan(&name, &serverView, &charset, &collation); err != nil {
				t.Fatal(err)
			}
			// The view qualifies the column and the call does not: the clauses are what is
			// compared.
			call := deparsedCall(t, c, viewCall)
			if clausesOf := call[strings.Index(call, " RETURNING "):]; !strings.Contains(serverView, clausesOf) {
				t.Errorf("view:\nomni:   %s\nserver: %s", call, serverView)
			}
		}))
	}
}

// TestJSONQueryStoredFormsReadBack: what this package writes for a form reads back to the
// same text, in a generated column and in a query.
func TestJSONQueryStoredFormsReadBack(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d", nil); err != nil {
		t.Fatal(err)
	}
	for _, clauses := range append([]string{"", "RETURNING JSON DOT ON MISMATCH"}, jsonQueryStoredForms...) {
		t.Run(clauses, func(t *testing.T) {
			call := "JSON_QUERY('{}', '$.b' " + clauses + ")"
			first := deparsedCall(t, c, call)
			if again := deparsedCall(t, c, first); again != first {
				t.Errorf("deparsed %q, which reads back as %q", first, again)
			}
			_, query, second := roundTrip(t, c, "SELECT "+call)
			if again := DeparseQuery(second); again != query {
				t.Errorf("deparsed %q, which reads back as %q", query, again)
			}
			if clausesOf := first[strings.Index(first, " RETURNING "):]; !strings.Contains(query, clausesOf) {
				t.Errorf("the query %q has lost the clauses of %q", query, first)
			}
		})
	}
}

// TestJSONQueryReturningIsTheResultType: RETURNING says what the call answers.
func TestJSONQueryReturningIsTheResultType(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d", nil); err != nil {
		t.Fatal(err)
	}
	q, _, _ := roundTrip(t, c, "SELECT JSON_QUERY('{}', '$.a' RETURNING JSON) AS x")
	call, ok := q.TargetList[0].Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("target = %T", q.TargetList[0].Expr)
	}
	if call.ResultType == nil || call.ResultType.BaseType != BaseTypeJSON {
		t.Errorf("ResultType = %+v, want JSON", call.ResultType)
	}
}

// TestOracle_JSONQueryFunctionalIndex: a functional index over JSON_QUERY is taken or
// refused as the server takes it, with the server's code. Written ASIS, with a RETURNING
// that is no JSON and without MULTIVALUE, it is a multi-valued index; anything else answers
// JSON (3753), or nothing the server will return (5443).
func TestOracle_JSONQueryFunctionalIndex(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	for _, clauses := range []string{
		"", "RETURNING JSON", "RETURNING UNSIGNED", "RETURNING JSON MULTIVALUE", "RETURNING UNSIGNED MULTIVALUE",
		"ASIS", "RETURNING UNSIGNED ASIS", "RETURNING CHAR(40) ASIS", "RETURNING JSON ASIS",
		"ASIS MULTIVALUE", "RETURNING UNSIGNED ASIS MULTIVALUE",
		"RETURNING UNSIGNED WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH",
		"RETURNING DATETIME(3) WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH",
		"RETURNING UNSIGNED WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH MULTIVALUE",
	} {
		t.Run(clauses, gapped(func(t testing.TB) {
			table := "CREATE TABLE f (id INT PRIMARY KEY, a JSON, KEY k ((JSON_QUERY(a, '$.x' " + clauses + "))))"
			serverCode := 0
			if _, err := storeAndRead(ctx, o, "zz_json_query_index", table, "f"); err != nil {
				var merr *gomysql.MySQLError
				if !errors.As(err, &merr) {
					t.Fatal(err)
				}
				serverCode = int(merr.Number)
			}
			c := New()
			results, err := c.Exec("CREATE DATABASE d; USE d; "+table, nil)
			if err != nil {
				t.Fatal(err)
			}
			catalogCode := 0
			var cerr *Error
			if errors.As(results[len(results)-1].Error, &cerr) {
				catalogCode = cerr.Code
			}
			if catalogCode != serverCode {
				t.Errorf("catalog answered %d, the server %d", catalogCode, serverCode)
			}
		}))
	}
}

// TestMultiValuedKeyStoredForm: a key written as a cast to an array is compared as the
// JSON_QUERY call the server stores, and an operand the server was not asked about is
// left as it was written.
func TestMultiValuedKeyStoredForm(t *testing.T) {
	const tail = " WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)"
	for key, want := range map[string]string{
		"CAST(j->'$.ids' AS UNSIGNED ARRAY)":                       "json_query(`j`,'$.ids' RETURNING unsigned" + tail,
		"CAST(j->>'$.ids' AS UNSIGNED ARRAY)":                      "json_query(`j`,'$.ids' RETURNING unsigned" + tail,
		"CAST(JSON_EXTRACT(j,'$.ids') AS SIGNED ARRAY)":            "json_query(`j`,'$.ids' RETURNING signed" + tail,
		"CAST(t.j->'$.ids' AS CHAR(40) ARRAY)":                     "json_query(`j`,'$.ids' RETURNING char(40) charset utf8mb4" + tail,
		"CAST(j AS DECIMAL(10,2) ARRAY)":                           "json_query(`j`,'$' RETURNING number(10,2)" + tail,
		"CAST((j) AS DATETIME ARRAY)":                              "json_query(`j`,'$' RETURNING datetime(0)" + tail,
		"CAST(j->'$.ids' AS UNSIGNED)":                             "",
		"CAST(JSON_EXTRACT(j,'$.a','$.b') AS UNSIGNED ARRAY)":      "",
		"CAST(JSON_EXTRACT(j,CONCAT('$.','a')) AS UNSIGNED ARRAY)": "",
		"CAST(d.JSON_EXTRACT(j,'$.a') AS UNSIGNED ARRAY)":          "",
		"CAST(JSON_MERGE_PATCH(j,k) AS UNSIGNED ARRAY)":            "",
		"json_query(`j`,'$.ids' RETURNING unsigned" + tail:         "",
		"j + 1": "",
	} {
		t.Run(key, func(t *testing.T) {
			list, err := parser.Parse("SELECT " + key)
			if err != nil {
				t.Fatal(err)
			}
			expr := list.Items[0].(*nodes.SelectStmt).TargetList[0]
			if rt, ok := expr.(*nodes.ResTarget); ok {
				expr = rt.Val
			}
			got := ""
			if stored := multiValuedKeyStoredForm(expr); stored != nil {
				got = nodeToSQL(stored)
			}
			if got != want {
				t.Errorf("stored form is %q, want %q", got, want)
			}
		})
	}
}

// TestMultiValuedIndexGoesInWithItsTable: a new table's multi-valued index is part of its
// CREATE TABLE, in either spelling, and its other indexes are added after. OceanBase builds a
// multi-valued index with its table only (TestOracle_IndexMigrationApplyCorrectness).
func TestMultiValuedIndexGoesInWithItsTable(t *testing.T) {
	const stored = "json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)"
	for name, key := range map[string]string{
		"written": "CAST(j->'$.ids' AS UNSIGNED ARRAY)",
		"stored":  stored,
	} {
		t.Run(name, func(t *testing.T) {
			to, err := LoadSQL("CREATE DATABASE d; USE d; CREATE TABLE t (id INT PRIMARY KEY, a INT, j JSON, KEY ka (a), KEY mv ((" + key + ")))")
			if err != nil {
				t.Fatal(err)
			}
			n := NormalizerFor(MySQL80)
			from := New()
			plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
			var create, adds []string
			for _, op := range plan.Ops {
				switch op.Type {
				case OpCreateTable:
					create = append(create, op.SQL)
				case OpAddIndex:
					adds = append(adds, op.SQL)
				}
			}
			if len(create) != 1 || !strings.Contains(create[0], "KEY `mv` ((") {
				t.Errorf("CREATE TABLE does not carry the multi-valued index:\n%s", strings.Join(create, "\n"))
			}
			if len(adds) != 1 || !strings.Contains(adds[0], "ADD KEY `ka`") {
				t.Errorf("want one ADD, of the plain index, got:\n%s", strings.Join(adds, "\n"))
			}
		})
	}
}

// TestMultiValuedKeyEqualsItsStoredText: the two spellings of a multi-valued key are one
// index to the diff, and another path or type is another index.
func TestMultiValuedKeyEqualsItsStoredText(t *testing.T) {
	load := func(key string) *Catalog {
		t.Helper()
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((" + key + ")))")
		if err != nil {
			t.Fatal(err)
		}
		return c
	}
	n := NormalizerFor(MySQL80)
	written := load("CAST(j->'$.ids' AS UNSIGNED ARRAY)")
	const tail = " WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)"
	for key, same := range map[string]bool{
		"json_query(`j`,'$.ids' RETURNING unsigned" + tail:                                                                   true,
		"JSON_QUERY(j, '$.ids' RETURNING UNSIGNED WITHOUT ARRAY WRAPPER ASIS ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH)": true,
		"CAST(JSON_EXTRACT(j,'$.ids') AS UNSIGNED INTEGER ARRAY)":                                                            true,
		"json_query(`j`,'$.other' RETURNING unsigned" + tail:                                                                 false,
		"json_query(`j`,'$.ids' RETURNING signed" + tail:                                                                     false,
		"json_query(`j`,'$.ids' RETURNING unsigned WITH WRAPPER asis error on error null on empty null on mismatch)":         false,
		"CAST(j->'$.ids' AS SIGNED ARRAY)":                                                                                   false,
	} {
		t.Run(key, func(t *testing.T) {
			other := load(key)
			forward, back := DiffWithNormalizer(written, other, n), DiffWithNormalizer(other, written, n)
			if forward.IsEmpty() != same || back.IsEmpty() != same {
				t.Errorf("the diff is empty forward %v and back %v, want %v", forward.IsEmpty(), back.IsEmpty(), same)
			}
		})
	}
}

// TestLoadMetadataInstallsMultiValuedIndex: a synced multi-valued index, whose key OceanBase
// reports as the JSON_QUERY call it stores, installs as that expression, and is the index
// its table's DDL writes as a cast to an array.
func TestLoadMetadataInstallsMultiValuedIndex(t *testing.T) {
	const stored = "json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)"
	synced, report := loadMySQLSnapshot(t, snapshotWithTables(&metadata.TableMetadata{
		Name: "t",
		Columns: []*metadata.ColumnMetadata{
			{Name: "id", Type: "int"},
			{Name: "j", Type: "json", Nullable: true},
		},
		Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Expressions: []string{"id"}},
			{Name: "mv", Type: "BTREE", Expressions: []string{stored}, Visible: true},
		},
	}))
	requireReport(t, report, nil, nil)
	key := func(c *Catalog) *Index {
		t.Helper()
		for _, idx := range requireTable(t, c, "t").Indexes {
			if idx.Name == "mv" {
				return idx
			}
		}
		t.Fatal("index mv not installed")
		return nil
	}
	mv := key(synced)
	if len(mv.Columns) != 1 || mv.Columns[0].ExprNode == nil || !isMultiValuedIndex(mv) {
		t.Fatalf("mv key parts = %+v, want the one JSON_QUERY expression", mv.Columns)
	}
	written, err := LoadSQL("CREATE DATABASE " + snapshotDB + "; USE " + snapshotDB + "; CREATE TABLE t (id INT PRIMARY KEY, j JSON, KEY mv ((CAST(j->'$.ids' AS UNSIGNED ARRAY))))")
	if err != nil {
		t.Fatal(err)
	}
	n := NormalizerFor(MySQL80)
	if got, want := canonicalIndex(mv, n), canonicalIndex(key(written), n); got != want {
		t.Errorf("the synced index is not the written one:\nsynced:  %s\nwritten: %s", got, want)
	}
}
