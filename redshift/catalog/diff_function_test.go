package catalog

import (
	"testing"
)

func TestDiffFunction(t *testing.T) {
	tests := []struct {
		name    string
		fromSQL string
		toSQL   string
		check   func(t *testing.T, entries []FunctionDiffEntry)
	}{
		{
			name:    "function added",
			fromSQL: `CREATE TABLE t1 (id int);`,
			toSQL: `CREATE TABLE t1 (id int);
				CREATE FUNCTION add_one(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x + 1';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffAdd {
					t.Fatalf("expected DiffAdd, got %d", entries[0].Action)
				}
				if entries[0].Name != "add_one" {
					t.Fatalf("expected name add_one, got %s", entries[0].Name)
				}
			},
		},
		{
			name:    "function dropped",
			fromSQL: `CREATE FUNCTION add_one(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x + 1';`,
			toSQL:   `SELECT 1;`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffDrop {
					t.Fatalf("expected DiffDrop, got %d", entries[0].Action)
				}
				if entries[0].Name != "add_one" {
					t.Fatalf("expected name add_one, got %s", entries[0].Name)
				}
			},
		},
		{
			name:    "function body changed",
			fromSQL: `CREATE FUNCTION add_one(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x + 1';`,
			toSQL:   `CREATE FUNCTION add_one(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x + 2';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
				if entries[0].From.Body != "SELECT x + 1" {
					t.Fatalf("unexpected from body: %s", entries[0].From.Body)
				}
				if entries[0].To.Body != "SELECT x + 2" {
					t.Fatalf("unexpected to body: %s", entries[0].To.Body)
				}
			},
		},
		{
			name:    "function volatility changed",
			fromSQL: `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql VOLATILE AS 'SELECT x';`,
			toSQL:   `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql IMMUTABLE AS 'SELECT x';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
				if entries[0].From.Volatile != 'v' {
					t.Fatalf("expected from volatile='v', got %c", entries[0].From.Volatile)
				}
				if entries[0].To.Volatile != 'i' {
					t.Fatalf("expected to volatile='i', got %c", entries[0].To.Volatile)
				}
			},
		},
		{
			name:    "function strictness changed",
			fromSQL: `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			toSQL:   `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql STRICT AS 'SELECT x';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
				if entries[0].From.IsStrict != false {
					t.Fatalf("expected from IsStrict=false")
				}
				if entries[0].To.IsStrict != true {
					t.Fatalf("expected to IsStrict=true")
				}
			},
		},
		{
			name:    "function security changed",
			fromSQL: `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			toSQL:   `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql SECURITY DEFINER AS 'SELECT x';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
				if entries[0].From.SecDef != false {
					t.Fatalf("expected from SecDef=false")
				}
				if entries[0].To.SecDef != true {
					t.Fatalf("expected to SecDef=true")
				}
			},
		},
		{
			name:    "function language changed",
			fromSQL: `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			toSQL:   `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE plpgsql AS 'BEGIN RETURN x; END';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
				if entries[0].From.Language != "sql" {
					t.Fatalf("expected from language=sql, got %s", entries[0].From.Language)
				}
				if entries[0].To.Language != "plpgsql" {
					t.Fatalf("expected to language=plpgsql, got %s", entries[0].To.Language)
				}
			},
		},
		{
			name:    "function return type changed",
			fromSQL: `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			toSQL:   `CREATE FUNCTION myfunc(x integer) RETURNS bigint LANGUAGE sql AS 'SELECT x::bigint';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
			},
		},
		{
			name:    "function parallel safety changed",
			fromSQL: `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			toSQL:   `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql PARALLEL SAFE AS 'SELECT x';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
				if entries[0].From.Parallel != 'u' {
					t.Fatalf("expected from parallel='u', got %c", entries[0].From.Parallel)
				}
				if entries[0].To.Parallel != 's' {
					t.Fatalf("expected to parallel='s', got %c", entries[0].To.Parallel)
				}
			},
		},
		{
			name:    "function leak-proof changed",
			fromSQL: `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			toSQL:   `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql LEAKPROOF AS 'SELECT x';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
				if entries[0].From.LeakProof != false {
					t.Fatalf("expected from LeakProof=false")
				}
				if entries[0].To.LeakProof != true {
					t.Fatalf("expected to LeakProof=true")
				}
			},
		},
		{
			name:    "function RETURNS SETOF changed",
			fromSQL: `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			toSQL:   `CREATE FUNCTION myfunc(x integer) RETURNS SETOF integer LANGUAGE sql AS 'SELECT x';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
				if entries[0].From.RetSet != false {
					t.Fatalf("expected from RetSet=false")
				}
				if entries[0].To.RetSet != true {
					t.Fatalf("expected to RetSet=true")
				}
			},
		},
		{
			name:    "procedure added and dropped",
			fromSQL: `CREATE PROCEDURE myproc(x integer) LANGUAGE sql AS 'SELECT 1';`,
			toSQL:   `CREATE PROCEDURE otherproc(y text) LANGUAGE sql AS 'SELECT 1';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 2 {
					t.Fatalf("expected 2 entries, got %d", len(entries))
				}
				var added, dropped bool
				for _, e := range entries {
					if e.Action == DiffAdd && e.Name == "otherproc" {
						added = true
						if e.To.Kind != 'p' {
							t.Fatalf("expected procedure Kind='p', got %c", e.To.Kind)
						}
					}
					if e.Action == DiffDrop && e.Name == "myproc" {
						dropped = true
						if e.From.Kind != 'p' {
							t.Fatalf("expected procedure Kind='p', got %c", e.From.Kind)
						}
					}
				}
				if !added {
					t.Fatal("expected otherproc to be added")
				}
				if !dropped {
					t.Fatal("expected myproc to be dropped")
				}
			},
		},
		{
			name: "overloaded functions are distinct",
			fromSQL: `
				CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';
				CREATE FUNCTION myfunc(x text) RETURNS text LANGUAGE sql AS 'SELECT x';`,
			toSQL: `
				CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';
				CREATE FUNCTION myfunc(x text) RETURNS text LANGUAGE sql AS 'SELECT x || x';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry (only text overload changed), got %d", len(entries))
				}
				if entries[0].Action != DiffModify {
					t.Fatalf("expected DiffModify, got %d", entries[0].Action)
				}
				// The modified one should be the text overload.
				if entries[0].To.Body != "SELECT x || x" {
					t.Fatalf("expected modified body for text overload, got %s", entries[0].To.Body)
				}
			},
		},
		{
			name: "function identity uses schema name and arg types",
			fromSQL: `
				CREATE SCHEMA s1;
				CREATE FUNCTION s1.myfunc(x integer, y text) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			toSQL: `
				CREATE SCHEMA s1;
				CREATE FUNCTION s1.myfunc(x integer, y text) RETURNS integer LANGUAGE sql AS 'SELECT x + 1';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 1 {
					t.Fatalf("expected 1 entry, got %d", len(entries))
				}
				if entries[0].SchemaName != "s1" {
					t.Fatalf("expected schema s1, got %s", entries[0].SchemaName)
				}
				id := entries[0].Identity
				if id != "s1.myfunc(integer,text)" {
					t.Fatalf("expected identity 's1.myfunc(integer,text)', got %s", id)
				}
			},
		},
		{
			name:    "function unchanged produces no entry",
			fromSQL: `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			toSQL:   `CREATE FUNCTION myfunc(x integer) RETURNS integer LANGUAGE sql AS 'SELECT x';`,
			check: func(t *testing.T, entries []FunctionDiffEntry) {
				if len(entries) != 0 {
					t.Fatalf("expected 0 entries, got %d", len(entries))
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from, err := LoadSQL(tt.fromSQL)
			if err != nil {
				t.Fatal(err)
			}
			to, err := LoadSQL(tt.toSQL)
			if err != nil {
				t.Fatal(err)
			}
			entries := diffFunctions(from, to)
			tt.check(t, entries)
		})
	}
}
