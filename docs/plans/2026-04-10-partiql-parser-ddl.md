# PartiQL Parser-DDL Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add CREATE TABLE, CREATE INDEX, DROP TABLE, and DROP INDEX parsing + introduce ParseStatement as the statement-level entry point.

**Architecture:** `ddl.go` (parseCreateCommand + parseDropCommand + parsePathSimple), `parser.go` (ParseStatement dispatcher), `parser_test.go` (TestParser_StmtGoldens harness + DDL error cases), `testdata/parser-ddl/` (8 golden pairs).

**Tech Stack:** Same as parser-foundation. AST nodes pre-exist: `CreateTableStmt{Name *VarRef, Loc}`, `CreateIndexStmt{Table *VarRef, Paths []*PathExpr, Loc}`, `DropTableStmt{Name *VarRef, Loc}`, `DropIndexStmt{Index *VarRef, Table *VarRef, Loc}`.

**IMPORTANT constraints:**
- Do NOT remove INSERT/UPDATE/DELETE stubs from `parseSelectExpr` — TestParser_AWSCorpus calls ParseExpr and depends on those stubs firing for DML corpus files.
- Use `.partiql` + `.golden` extensions (matching the foundation pattern, NOT `.sql` + `.json`).
- Use `ast.NodeToString` for golden output (NOT JSON).
- Use plain `testing.T` assertions (NOT testify).
- The actual PartiQL DDL grammar has NO `IF NOT EXISTS`, `IF EXISTS`, `UNIQUE`, or column definitions. Match the grammar exactly.

---

## Grammar reference (PartiQLParser.g4 lines 73-86)

```
ddl
    : createCommand
    | dropCommand
    ;

createCommand
    : CREATE TABLE symbolPrimitive
    | CREATE INDEX ON symbolPrimitive PAREN_LEFT pathSimple ( COMMA pathSimple )* PAREN_RIGHT
    ;

dropCommand
    : DROP TABLE target=symbolPrimitive
    | DROP INDEX target=symbolPrimitive ON on=symbolPrimitive
    ;
```

## AST nodes (verified in partiql/ast/stmts.go)

```
CreateTableStmt { Name *VarRef; Loc Loc }
CreateIndexStmt { Table *VarRef; Paths []*PathExpr; Loc Loc }
DropTableStmt   { Name *VarRef; Loc Loc }
DropIndexStmt   { Index *VarRef; Table *VarRef; Loc Loc }
```

## outfuncs.go format (verified)

```
CreateTableStmt{Name:VarRef{Name:t}}
CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:name} Steps:[]}]}
DropTableStmt{Name:VarRef{Name:t}}
DropIndexStmt{Index:VarRef{Name:idx} Table:VarRef{Name:t}}
```

---

### Task 1: parser.go — Add ParseStatement

**Files:**
- Modify: `partiql/parser/parser.go` (add ParseStatement)

Add `ParseStatement` to parser.go after `ParseExpr`:

```go
// ParseStatement parses a single PartiQL statement and asserts EOF.
// Dispatches based on the first token:
//   - CREATE/DROP        → DDL parsers (ddl.go)
//   - INSERT/UPDATE/etc  → deferred to parser-dml (DAG node 6)
//   - EXEC/EXECUTE       → deferred to parse-entry (DAG node 8)
//   - Everything else    → DQL: parse as expression via parseExprTop.
//     If the result implements StmtNode (e.g., SelectStmt after node 5
//     lands), return it directly. Otherwise defer to parse-entry
//     (DAG node 8) which will add bare-expression-as-statement support.
//
// Note: The INSERT/UPDATE/DELETE stubs in parseSelectExpr (expr.go)
// are intentionally KEPT. They serve as "you tried to use DML as an
// expression" error markers for ParseExpr callers. ParseStatement
// provides the statement-level dispatch for real DML handling.
func (p *Parser) ParseStatement() (ast.StmtNode, error) {
	if err := p.checkLexerErr(); err != nil {
		return nil, err
	}
	var stmt ast.StmtNode
	var err error
	switch p.cur.Type {
	case tokCREATE:
		stmt, err = p.parseCreateCommand()
	case tokDROP:
		stmt, err = p.parseDropCommand()
	case tokINSERT:
		return nil, p.deferredFeature("INSERT", "parser-dml (DAG node 6)")
	case tokUPDATE:
		return nil, p.deferredFeature("UPDATE", "parser-dml (DAG node 6)")
	case tokDELETE:
		return nil, p.deferredFeature("DELETE", "parser-dml (DAG node 6)")
	case tokREPLACE:
		return nil, p.deferredFeature("REPLACE", "parser-dml (DAG node 6)")
	case tokUPSERT:
		return nil, p.deferredFeature("UPSERT", "parser-dml (DAG node 6)")
	case tokREMOVE:
		return nil, p.deferredFeature("REMOVE", "parser-dml (DAG node 6)")
	case tokEXEC, tokEXECUTE:
		return nil, p.deferredFeature("EXEC", "parse-entry (DAG node 8)")
	default:
		// DQL fallback: the grammar says dql : expr, so any expression is
		// a valid DQL statement. Parse with parseExprTop (no EOF check).
		var expr ast.ExprNode
		expr, err = p.parseExprTop()
		if err != nil {
			return nil, err
		}
		// If the result implements StmtNode (e.g., SelectStmt after
		// node 5 lands), use it directly.
		if sn, ok := expr.(ast.StmtNode); ok {
			stmt = sn
		} else {
			// Bare expressions as statements (e.g., `1 + 2;`) are
			// grammar-legal but need an ExprStmt wrapper that doesn't
			// exist yet. Defer to parse-entry (node 8).
			return nil, p.deferredFeature("bare expression as statement", "parse-entry (DAG node 8)")
		}
	}
	if err != nil {
		return nil, err
	}
	if p.cur.Type != tokEOF {
		return nil, &ParseError{
			Message: fmt.Sprintf("unexpected token %q after statement", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
	return stmt, nil
}
```

Also need to verify token constants exist: tokREPLACE, tokUPSERT, tokREMOVE, tokEXEC, tokEXECUTE, tokSET. If any are missing, the switch case will fail to compile.

Tests: the existing suite must still pass. ParseStatement is additive — no existing code calls it yet.

```bash
go build ./partiql/parser/...
go test ./partiql/parser/...
```

Commit:
```
feat(partiql/parser): add ParseStatement dispatcher

Introduces ParseStatement as the statement-level entry point.
Dispatches CREATE/DROP to DDL parsers (stubbed until Task 2),
DML keywords to deferred-feature stubs for parser-dml (node 6),
EXEC to a stub for parse-entry (node 8), and everything else to
the DQL expression fallback via parseExprTop.

The INSERT/UPDATE/DELETE stubs in parseSelectExpr are intentionally
kept — they serve as expression-level error markers for ParseExpr
callers. ParseStatement provides the parallel statement-level
dispatch.
```

---

### Task 2: ddl.go — parseCreateCommand + parseDropCommand + parsePathSimple

**Files:**
- Create: `partiql/parser/ddl.go`

Create ddl.go with the 4 DDL parsers and the pathSimple helper. Match the grammar EXACTLY — no IF NOT EXISTS, no UNIQUE, no column definitions:

```go
package parser

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
)

// parseCreateCommand handles CREATE TABLE and CREATE INDEX.
//
// Grammar: createCommand (lines 78-81):
//   CREATE TABLE symbolPrimitive
//   CREATE INDEX ON symbolPrimitive ( pathSimple (, pathSimple)* )
func (p *Parser) parseCreateCommand() (ast.StmtNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume CREATE

	switch p.cur.Type {
	case tokTABLE:
		p.advance() // consume TABLE
		name, cs, loc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.CreateTableStmt{
			Name: &ast.VarRef{Name: name, CaseSensitive: cs, Loc: loc},
			Loc:  ast.Loc{Start: start, End: loc.End},
		}, nil

	case tokINDEX:
		p.advance() // consume INDEX
		if _, err := p.expect(tokON); err != nil {
			return nil, err
		}
		tableName, tableCS, tableLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokPAREN_LEFT); err != nil {
			return nil, err
		}
		var paths []*ast.PathExpr
		for {
			path, err := p.parsePathSimple()
			if err != nil {
				return nil, err
			}
			paths = append(paths, path)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume ,
		}
		rp, err := p.expect(tokPAREN_RIGHT)
		if err != nil {
			return nil, err
		}
		return &ast.CreateIndexStmt{
			Table: &ast.VarRef{Name: tableName, CaseSensitive: tableCS, Loc: tableLoc},
			Paths: paths,
			Loc:   ast.Loc{Start: start, End: rp.Loc.End},
		}, nil
	}

	return nil, &ParseError{
		Message: fmt.Sprintf("expected TABLE or INDEX after CREATE, got %q", p.cur.Str),
		Loc:     p.cur.Loc,
	}
}

// parseDropCommand handles DROP TABLE and DROP INDEX.
//
// Grammar: dropCommand (lines 83-86):
//   DROP TABLE target=symbolPrimitive
//   DROP INDEX target=symbolPrimitive ON on=symbolPrimitive
func (p *Parser) parseDropCommand() (ast.StmtNode, error) {
	start := p.cur.Loc.Start
	p.advance() // consume DROP

	switch p.cur.Type {
	case tokTABLE:
		p.advance() // consume TABLE
		name, cs, loc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.DropTableStmt{
			Name: &ast.VarRef{Name: name, CaseSensitive: cs, Loc: loc},
			Loc:  ast.Loc{Start: start, End: loc.End},
		}, nil

	case tokINDEX:
		p.advance() // consume INDEX
		indexName, indexCS, indexLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokON); err != nil {
			return nil, err
		}
		tableName, tableCS, tableLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		return &ast.DropIndexStmt{
			Index: &ast.VarRef{Name: indexName, CaseSensitive: indexCS, Loc: indexLoc},
			Table: &ast.VarRef{Name: tableName, CaseSensitive: tableCS, Loc: tableLoc},
			Loc:   ast.Loc{Start: start, End: tableLoc.End},
		}, nil
	}

	return nil, &ParseError{
		Message: fmt.Sprintf("expected TABLE or INDEX after DROP, got %q", p.cur.Str),
		Loc:     p.cur.Loc,
	}
}

// parsePathSimple parses a simplified path for DDL/DML targets.
// Unlike expression-level paths, pathSimple only allows literal and
// symbol bracket keys and dot-symbol steps — no expression indices,
// no wildcards.
//
// Grammar: pathSimple (lines 110-117):
//   symbolPrimitive pathSimpleSteps*
//
// Returns a *PathExpr with a VarRef root and DotStep/IndexStep steps.
func (p *Parser) parsePathSimple() (*ast.PathExpr, error) {
	rootName, rootCS, rootLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, err
	}
	root := &ast.VarRef{Name: rootName, CaseSensitive: rootCS, Loc: rootLoc}
	var steps []ast.PathStep
	endLoc := rootLoc.End
	for {
		switch p.cur.Type {
		case tokPERIOD:
			stepStart := p.cur.Loc.Start
			p.advance() // consume .
			name, cs, nameLoc, err := p.parseSymbolPrimitive()
			if err != nil {
				return nil, err
			}
			steps = append(steps, &ast.DotStep{
				Field: name, CaseSensitive: cs,
				Loc: ast.Loc{Start: stepStart, End: nameLoc.End},
			})
			endLoc = nameLoc.End
		case tokBRACKET_LEFT:
			stepStart := p.cur.Loc.Start
			p.advance() // consume [
			var idx ast.ExprNode
			switch p.cur.Type {
			case tokSCONST, tokICONST, tokFCONST:
				idx, err = p.parseLiteral()
				if err != nil {
					return nil, err
				}
			case tokIDENT, tokIDENT_QUOTED:
				n, cs, loc, err := p.parseSymbolPrimitive()
				if err != nil {
					return nil, err
				}
				idx = &ast.VarRef{Name: n, CaseSensitive: cs, Loc: loc}
			default:
				return nil, &ParseError{
					Message: fmt.Sprintf("expected literal or identifier in path bracket, got %q", p.cur.Str),
					Loc:     p.cur.Loc,
				}
			}
			rp, err := p.expect(tokBRACKET_RIGHT)
			if err != nil {
				return nil, err
			}
			steps = append(steps, &ast.IndexStep{
				Index: idx,
				Loc:   ast.Loc{Start: stepStart, End: rp.Loc.End},
			})
			endLoc = rp.Loc.End
		default:
			goto done
		}
	}
done:
	return &ast.PathExpr{
		Root: root, Steps: steps,
		Loc: ast.Loc{Start: rootLoc.Start, End: endLoc},
	}, nil
}
```

Tests:
```bash
go build ./partiql/parser/...
go test ./partiql/parser/...
```

Commit:
```
feat(partiql/parser): ddl.go with CREATE/DROP TABLE/INDEX

Adds parseCreateCommand (CREATE TABLE, CREATE INDEX ON ... (...)),
parseDropCommand (DROP TABLE, DROP INDEX ... ON ...), and
parsePathSimple (simplified DML/DDL path: symbol + dot/bracket
steps without expression indices or wildcards).

Matches PartiQLParser.g4 lines 73-86 and 110-117 exactly.
PartiQL DDL has no column definitions, IF NOT EXISTS, or UNIQUE
modifiers.
```

---

### Task 3: TestParser_StmtGoldens harness + 8 golden pairs + error cases

**Files:**
- Modify: `partiql/parser/parser_test.go` (add TestParser_StmtGoldens + DDL errors)
- Create: `partiql/parser/testdata/parser-ddl/` (8 .partiql + 8 .golden files)

Add TestParser_StmtGoldens to parser_test.go. This mirrors TestParser_Goldens but calls ParseStatement and walks `testdata/parser-ddl/`:

```go
// TestParser_StmtGoldens walks testdata/parser-ddl/*.partiql and
// compares ParseStatement output against .golden files.
func TestParser_StmtGoldens(t *testing.T) {
	files, err := filepath.Glob("testdata/parser-ddl/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no DDL golden inputs found under testdata/parser-ddl/")
	}
	for _, inPath := range files {
		name := strings.TrimSuffix(filepath.Base(inPath), ".partiql")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatal(err)
			}
			p := NewParser(string(input))
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got := ast.NodeToString(stmt)
			goldenPath := strings.TrimSuffix(inPath, ".partiql") + ".golden"
			if *update {
				if err := os.WriteFile(goldenPath, []byte(got+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("golden file missing: %s (run with -update to create)", goldenPath)
			}
			if got+"\n" != string(want) {
				t.Errorf("AST mismatch\ngot:\n%s\nwant:\n%s", got, string(want))
			}
		})
	}
}
```

**8 golden pairs** under `testdata/parser-ddl/`:

| Input file | Content | Expected golden |
|-----------|---------|----------------|
| `create_table.partiql` | `CREATE TABLE t` | `CreateTableStmt{Name:VarRef{Name:t}}` |
| `create_table_quoted.partiql` | `CREATE TABLE "MyTable"` | `CreateTableStmt{Name:VarRef{Name:MyTable CaseSensitive:true}}` |
| `create_index.partiql` | `CREATE INDEX ON t (name)` | `CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:name} Steps:[]}]}` |
| `create_index_multi.partiql` | `CREATE INDEX ON t (a, b.c)` | `CreateIndexStmt{Table:VarRef{Name:t} Paths:[PathExpr{Root:VarRef{Name:a} Steps:[]} PathExpr{Root:VarRef{Name:b} Steps:[DotStep{Field:c}]}]}` |
| `drop_table.partiql` | `DROP TABLE t` | `DropTableStmt{Name:VarRef{Name:t}}` |
| `drop_table_quoted.partiql` | `DROP TABLE "MyTable"` | `DropTableStmt{Name:VarRef{Name:MyTable CaseSensitive:true}}` |
| `drop_index.partiql` | `DROP INDEX idx ON t` | `DropIndexStmt{Index:VarRef{Name:idx} Table:VarRef{Name:t}}` |
| `drop_index_quoted.partiql` | `DROP INDEX "myIdx" ON "myTable"` | `DropIndexStmt{Index:VarRef{Name:myIdx CaseSensitive:true} Table:VarRef{Name:myTable CaseSensitive:true}}` |

Each `.golden` file ends with exactly one `\n`.

**Error cases** (add to TestParser_Errors):

```go
// DDL error cases (Task 3 of parser-ddl)
{"create_no_keyword", "CREATE", "expected TABLE or INDEX after CREATE"},
{"create_bad_keyword", "CREATE SLAB t", "expected TABLE or INDEX after CREATE"},
{"drop_no_keyword", "DROP", "expected TABLE or INDEX after DROP"},
{"drop_bad_keyword", "DROP SLAB t", "expected TABLE or INDEX after DROP"},
{"drop_index_missing_on", "DROP INDEX idx t", "expected ON"},
{"create_index_missing_paren", "CREATE INDEX ON t name", "expected PAREN_LEFT"},
```

Note: the error cases test ParseStatement. Currently TestParser_Errors uses ParseExpr. Add a parallel error test section that uses ParseStatement for DDL errors:

```go
// TestParser_DDLErrors tests DDL-specific error cases via ParseStatement.
func TestParser_DDLErrors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// ... the 6 cases above
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseStatement()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}
```

Tests:
```bash
go test -v -run "TestParser_StmtGoldens|TestParser_DDLErrors" ./partiql/parser/...
go test ./partiql/parser/...
```

Commit:
```
test(partiql/parser): DDL golden tests + error cases

Adds TestParser_StmtGoldens walking testdata/parser-ddl/ with
ParseStatement entry point. 8 golden pairs covering CREATE TABLE,
CREATE INDEX, DROP TABLE, DROP INDEX with both bare and quoted
identifiers + multi-key index paths.

Adds TestParser_DDLErrors with 6 cases for DDL syntax errors.
```

---

### Task 4: Final verification + commit

```bash
go test -v ./partiql/parser/...
go vet ./partiql/parser/...
gofmt -l partiql/parser/
grep -n $'\xe2\x80\x9d' partiql/parser/*.go
```

All existing tests must still pass. DDL adds ~8 golden sub-tests + ~6 error sub-tests. Total new: ~14.

Commit (allow-empty if no changes):
```
chore(partiql/parser): final verification pass for parser-ddl (DAG node 7)
```
