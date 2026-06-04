package deparse

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/trino/internal/trinooracle"
)

// connectOracle dials the live Trino oracle for the round-trip gate, mirroring
// the skip discipline the parser nodes use: skipped (not failed) in -short mode
// or when no Trino server is reachable, so the suite stays green offline while
// exercising the real parser when a server is up. Start one with:
//
//	docker run -d --name trino-oracle -p 18080:8080 trinodb/trino:latest
func connectOracle(t *testing.T) *trinooracle.Oracle {
	t.Helper()
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := trinooracle.Connect("")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ver, err := o.Ping(ctx)
	if err != nil {
		t.Skipf("trino oracle not reachable (start a Trino server and set %s if needed): %v", trinooracle.EnvURL, err)
	}
	t.Logf("connected to Trino %s", ver)
	return o
}

// TestGeneratedSDLParsesOnTrino is the deparse round-trip / structural gate
// required by the correctness protocol for a metadata->SDL emitter: every
// CREATE TABLE the deparser emits must be accepted by the real Trino 481 parser.
// This proves the output is genuine, parseable Trino DDL rather than merely
// byte-matching a (possibly buggy) legacy expectation.
//
// We assert acceptance, not successful execution: Trino accepts the syntax but
// then fails semantically (e.g. SCHEMA_NOT_FOUND in the memory catalog), which
// the oracle classifies as Accepted=true. We deliberately do not feed the
// special-characters legacy case here — `test.table` parses as a three-part
// name and `id-field` is not a legal column type — because that case validates
// the *format*, not Trino-acceptability; its acceptance is covered by the
// golden test, and Trino-shaped variants are checked below.
func TestGeneratedSDLParsesOnTrino(t *testing.T) {
	o := connectOracle(t)

	tables := []*TableMetadata{
		{
			Name: "testtable",
			Columns: []*ColumnMetadata{
				{Name: "id", Type: "bigint", Nullable: false},
				{Name: "name", Type: "varchar", Nullable: true},
			},
		},
		{
			// Empty-column table: Trino rejects "CREATE TABLE t ()" (a table
			// needs at least one column or a column-less form is invalid), so
			// this exercises the negative side of the round trip below, not here.
			Name: "single_col",
			Columns: []*ColumnMetadata{
				{Name: "c", Type: "integer", Nullable: false},
			},
		},
		{
			// A spread of common Trino types, including parameterized and
			// nested ones, to make sure the verbatim type text round-trips.
			Name: "typed",
			Columns: []*ColumnMetadata{
				{Name: "a", Type: "varchar(10)", Nullable: false},
				{Name: "b", Type: "decimal(10,2)", Nullable: true},
				{Name: "c", Type: "array(integer)", Nullable: true},
				{Name: "d", Type: "map(varchar, bigint)", Nullable: true},
				{Name: "e", Type: "row(x integer, y varchar)", Nullable: true},
				{Name: "f", Type: "timestamp(3) with time zone", Nullable: true},
			},
		},
		{
			// A quoted identifier whose normalized form is case-sensitive: the
			// double quotes the deparser always adds are exactly what preserves
			// "MixedCase" and a reserved word like "table" used as a column.
			Name: "MixedCase",
			Columns: []*ColumnMetadata{
				{Name: "MixedCol", Type: "bigint", Nullable: false},
				{Name: "table", Type: "varchar", Nullable: true},
			},
		},
	}

	for _, table := range tables {
		t.Run(table.Name, func(t *testing.T) {
			ddl, err := GetDatabaseDefinition(oneTable("memory", "default", table))
			if err != nil {
				t.Fatalf("GetDatabaseDefinition: %v", err)
			}
			for _, stmt := range splitStatements(ddl) {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				res, err := o.CheckSyntax(ctx, stmt)
				cancel()
				if err != nil {
					t.Fatalf("oracle CheckSyntax(%q): %v", stmt, err)
				}
				if !res.Accepted {
					t.Errorf("Trino rejected generated SDL as a syntax error:\n%s\nerrorName=%q code=%d msg=%s",
						stmt, res.ErrorName, res.ErrorCode, res.Message)
				}
			}
		})
	}
}

// TestEmptyColumnTableRejectedByTrino documents an honest limitation of the
// legacy format that the round-trip surfaces: a table with no columns renders
// to `CREATE TABLE ... ()`, which Trino's parser rejects as a SYNTAX_ERROR. The
// deparser reproduces the legacy byte-output (proven by the golden test), but
// that particular output is not executable Trino. We assert the rejection so
// the limitation is recorded, not hidden; it is flagged in the divergence
// ledger as an inherited-legacy issue, not a regression introduced here.
func TestEmptyColumnTableRejectedByTrino(t *testing.T) {
	o := connectOracle(t)

	ddl, err := GetDatabaseDefinition(oneTable("memory", "default", &TableMetadata{
		Name:    "empty_table",
		Columns: nil,
	}))
	if err != nil {
		t.Fatalf("GetDatabaseDefinition: %v", err)
	}
	stmts := splitStatements(ddl)
	if len(stmts) != 1 {
		t.Fatalf("expected one statement, got %d: %q", len(stmts), ddl)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := o.CheckSyntax(ctx, stmts[0])
	if err != nil {
		t.Fatalf("oracle CheckSyntax(%q): %v", stmts[0], err)
	}
	if res.Accepted {
		t.Errorf("expected Trino to reject the column-less CREATE TABLE, but it was accepted:\n%s", stmts[0])
	}
	if res.ErrorName != "SYNTAX_ERROR" {
		t.Errorf("expected SYNTAX_ERROR for column-less CREATE TABLE, got errorName=%q (msg=%s)", res.ErrorName, res.Message)
	}
}

// splitStatements breaks the multi-statement deparse output into individual
// trimmed CREATE TABLE statements, dropping the trailing blank lines. The
// deparser terminates each statement with ";" and separates them with "\n\n",
// so splitting on ";" and trimming is sufficient for the simple, generated DDL
// the deparser emits.
func splitStatements(ddl string) []string {
	var out []string
	for _, part := range strings.Split(ddl, ";") {
		s := strings.TrimSpace(part)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
