package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestParseQualifiedTypeMultiComponent locks down the parseGenericType
// fix that allows arbitrary-depth qualified type names. Before the fix,
// parseGenericType consumed exactly one dot, so 3-component names like
// `pg_catalog.int4` worked but `db.schema.mytype` failed in every type
// position.
//
// The 11 positions covered here were identified by codex review of the
// follow-up plan. Each exercises parseGenericType via a different call
// chain; the same fix unblocks all of them simultaneously.
//
// Plan: docs/plans/2026-04-14-pg-followups.md
func TestParseQualifiedTypeMultiComponent(t *testing.T) {
	cases := []struct {
		name      string
		sql       string
		wantNames []string // expected TypeName.Names list
	}{
		// ----------------------- 2-component (regression) -----------------------
		{
			name:      "2-component CAST (regression sanity)",
			sql:       `SELECT CAST(NULL AS pg_catalog.int4)`,
			wantNames: []string{"pg_catalog", "int4"},
		},
		{
			name:      "2-component CREATE TABLE column (regression sanity)",
			sql:       `CREATE TABLE t (c pg_catalog.int4)`,
			wantNames: []string{"pg_catalog", "int4"},
		},

		// ----------------------- 3-component, 11 positions ----------------------
		{
			name:      "3-component CAST",
			sql:       `SELECT CAST(NULL AS db.schema.mytype)`,
			wantNames: []string{"db", "schema", "mytype"},
		},
		{
			name:      "3-component TYPECAST",
			sql:       `SELECT 1::db.schema.mytype`,
			wantNames: []string{"db", "schema", "mytype"},
		},
		{
			name:      "3-component CREATE TABLE column",
			sql:       `CREATE TABLE t (c db.schema.mytype)`,
			wantNames: []string{"db", "schema", "mytype"},
		},
		{
			name:      "3-component ALTER TABLE",
			sql:       `ALTER TABLE t ALTER COLUMN c TYPE db.schema.mytype`,
			wantNames: []string{"db", "schema", "mytype"},
		},
		{
			name:      "3-component CREATE FUNCTION param",
			sql:       `CREATE FUNCTION f(x db.schema.mytype) RETURNS int AS 'select 1' LANGUAGE sql`,
			wantNames: []string{"db", "schema", "mytype"},
		},
		{
			name:      "3-component CREATE FUNCTION return",
			sql:       `CREATE FUNCTION f() RETURNS db.schema.mytype AS 'select 1' LANGUAGE sql`,
			wantNames: []string{"db", "schema", "mytype"},
		},
		{
			name:      "3-component RETURNS TABLE column",
			sql:       `CREATE FUNCTION f() RETURNS TABLE (x db.schema.mytype) AS 'select 1' LANGUAGE sql`,
			wantNames: []string{"db", "schema", "mytype"},
		},
		// Note: CREATE OPERATOR LEFTARG uses a dedicated test
		// (TestParseCreateOperatorQualifiedLeftarg below) because its
		// AST has multiple TypeName-shaped fields and the generic
		// findFirstTypeName walker stops on the FUNCTION value
		// before reaching LEFTARG.
		{
			name:      "3-component CREATE SEQUENCE AS",
			sql:       `CREATE SEQUENCE s AS db.schema.mytype`,
			wantNames: []string{"db", "schema", "mytype"},
		},
		{
			name:      "3-component XMLSERIALIZE",
			sql:       `SELECT XMLSERIALIZE(CONTENT '<a/>' AS db.schema.mytype)`,
			wantNames: []string{"db", "schema", "mytype"},
		},
		{
			name:      "3-component JSON_SERIALIZE RETURNING",
			sql:       `SELECT JSON_SERIALIZE('{}' RETURNING db.schema.mytype)`,
			wantNames: []string{"db", "schema", "mytype"},
		},

		// --------------------- 4-component (parser accepts) ---------------------
		// PG accepts any depth at parse time (recursive `attrs` rule). The
		// catalog rejects 4+ at name-resolution time — see
		// pg/catalog/nodeutil.go:typeNameParts. omni's parser should accept
		// the syntax to match PG; the catalog is the layer that decides what's
		// semantically valid.
		{
			name:      "4-component CREATE TABLE column",
			sql:       `CREATE TABLE t (c a.b.c.d)`,
			wantNames: []string{"a", "b", "c", "d"},
		},
		{
			name:      "5-component CAST",
			sql:       `SELECT CAST(NULL AS a.b.c.d.e)`,
			wantNames: []string{"a", "b", "c", "d", "e"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			tn := findFirstTypeName(t, stmts)
			if tn == nil {
				t.Fatalf("no TypeName found in parsed AST")
			}
			assertTypeNameNames(t, tn, tc.wantNames)
		})
	}
}

// TestParseCreateOperatorQualifiedLeftarg covers the CREATE OPERATOR
// LEFTARG / RIGHTARG case, which goes through parseDefArg → parseFuncType
// → parseTypename → parseSimpleTypename → parseGenericType. The generic
// walker in TestParseQualifiedTypeMultiComponent can't be reused here
// because CREATE OPERATOR has multiple TypeName-shaped fields (FUNCTION,
// LEFTARG, RIGHTARG); we drill into the LEFTARG DefElem by name.
func TestParseCreateOperatorQualifiedLeftarg(t *testing.T) {
	cases := []struct {
		name           string
		sql            string
		wantLeftargFQN []string
	}{
		{
			name:           "qualified LEFTARG with unqualified RIGHTARG",
			sql:            `CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = db.schema.mytype, RIGHTARG = int4)`,
			wantLeftargFQN: []string{"db", "schema", "mytype"},
		},
		{
			name:           "qualified LEFTARG with qualified RIGHTARG",
			sql:            `CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = db.schema.l, RIGHTARG = db.schema.r)`,
			wantLeftargFQN: []string{"db", "schema", "l"},
		},
		{
			name:           "2-component LEFTARG (regression)",
			sql:            `CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = pg_catalog.int4, RIGHTARG = int4)`,
			wantLeftargFQN: []string{"pg_catalog", "int4"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			if stmts == nil || len(stmts.Items) == 0 {
				t.Fatalf("no statements parsed")
			}
			raw, ok := stmts.Items[0].(*nodes.RawStmt)
			if !ok {
				t.Fatalf("expected RawStmt, got %T", stmts.Items[0])
			}
			def, ok := raw.Stmt.(*nodes.DefineStmt)
			if !ok {
				t.Fatalf("expected DefineStmt, got %T", raw.Stmt)
			}
			if def.Definition == nil {
				t.Fatalf("DefineStmt has nil Definition")
			}

			// Find the LEFTARG DefElem.
			var leftarg *nodes.DefElem
			for _, item := range def.Definition.Items {
				de, ok := item.(*nodes.DefElem)
				if !ok {
					continue
				}
				if de.Defname == "leftarg" {
					leftarg = de
					break
				}
			}
			if leftarg == nil {
				t.Fatalf("no LEFTARG DefElem found")
			}
			tn, ok := leftarg.Arg.(*nodes.TypeName)
			if !ok {
				t.Fatalf("LEFTARG.Arg: expected *TypeName, got %T", leftarg.Arg)
			}
			assertTypeNameNames(t, tn, tc.wantLeftargFQN)
		})
	}
}

// findFirstTypeName walks the AST and returns the first *nodes.TypeName
// it encounters. This is a heuristic that works for all test cases here
// because each test SQL has exactly one significant TypeName at the
// position under test (and the AST walk visits it before any nested
// TypeNames in unrelated subtrees).
func findFirstTypeName(t *testing.T, stmts *nodes.List) *nodes.TypeName {
	t.Helper()
	if stmts == nil {
		return nil
	}
	var found *nodes.TypeName
	var walk func(n nodes.Node)
	walk = func(n nodes.Node) {
		if n == nil || found != nil {
			return
		}
		if tn, ok := n.(*nodes.TypeName); ok {
			found = tn
			return
		}
		// Use reflection-free walking: switch on the common wrapper types.
		switch v := n.(type) {
		case *nodes.RawStmt:
			walk(v.Stmt)
		case *nodes.List:
			for _, item := range v.Items {
				walk(item)
				if found != nil {
					return
				}
			}
		case *nodes.CreateStmt:
			if v.TableElts != nil {
				walk(v.TableElts)
			}
		case *nodes.ColumnDef:
			if v.TypeName != nil {
				walk(v.TypeName)
			}
		case *nodes.CreateFunctionStmt:
			if v.Parameters != nil {
				walk(v.Parameters)
			}
			if found == nil && v.ReturnType != nil {
				walk(v.ReturnType)
			}
		case *nodes.FunctionParameter:
			if v.ArgType != nil {
				walk(v.ArgType)
			}
		case *nodes.AlterTableStmt:
			if v.Cmds != nil {
				walk(v.Cmds)
			}
		case *nodes.AlterTableCmd:
			if v.Def != nil {
				walk(v.Def)
			}
		case *nodes.SelectStmt:
			if v.TargetList != nil {
				walk(v.TargetList)
			}
		case *nodes.ResTarget:
			if v.Val != nil {
				walk(v.Val)
			}
		case *nodes.TypeCast:
			if v.TypeName != nil {
				walk(v.TypeName)
			}
		case *nodes.DefineStmt:
			if v.Definition != nil {
				walk(v.Definition)
			}
		case *nodes.DefElem:
			if v.Arg != nil {
				walk(v.Arg)
			}
		case *nodes.CreateSeqStmt:
			if v.Options != nil {
				walk(v.Options)
			}
		case *nodes.XmlSerialize:
			if v.TypeName != nil {
				walk(v.TypeName)
			}
		case *nodes.JsonSerializeExpr:
			if v.Output != nil && v.Output.TypeName != nil {
				walk(v.Output.TypeName)
			}
		case *nodes.FuncCall:
			// JSON_SERIALIZE may be desugared to FuncCall; walk args
			if v.Args != nil {
				walk(v.Args)
			}
		}
	}
	walk(stmts)
	return found
}
