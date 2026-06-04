package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-dml DAG node (with update_delete.go and
// merge.go): it implements Trino's INSERT and TRUNCATE data-manipulation
// statements as hand-written recursive-descent parsers over the token stream.
//
// The statement structs live here in package parser (matching the Trino
// convention for parser-built node types — Expr in expr.go, DataType in
// datatypes.go, the SELECT layer in select.go). They satisfy ast.Node so
// parseStmt can return them in the File statement list; like select.go's
// QueryStmt they carry the placeholder tag ast.T_Invalid because the ast
// node-tag set is closed to the ast-core node (File / Identifier /
// QualifiedName) and is extended only by the utility/dcl-tcl statement nodes —
// adding DML tags would collide with the concurrently-developed DDL node that
// also extends trino/ast/nodetags.go. Downstream consumers (analysis,
// completion) reach the concrete statement type by a Go type switch, not by
// the tag, exactly as they do for QueryStmt.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	statement
//	    : ...
//	    | INSERT_ INTO_ qualifiedName columnAliases? rootQuery # insertInto
//	    | TRUNCATE_ TABLE_ qualifiedName                       # truncateTable
//	    ;
//	rootQuery   : withFunction? query ;
//	columnAliases : LPAREN_ identifier (COMMA_ identifier)* RPAREN_ ;
//
// The implementation is adjudicated against the live Trino 481 oracle, not the
// literal legacy grammar. Oracle-confirmed facts (Trino 481) baked in:
//
//	D-INS1 (INSERT requires a source query). `INSERT INTO t` and
//	   `INSERT INTO t (a, b)` are SYNTAX_ERRORs: the rootQuery is mandatory.
//	   The query may be SELECT / VALUES / TABLE / WITH … / a parenthesized
//	   query — exactly the queryStmt entry the SELECT node already parses.
//	D-INS2 (target is at most catalog.schema.table). `INSERT INTO a.b.c.d …`
//	   is a SYNTAX_ERROR; the target qualifiedName is bounded to 3 parts.
//	D-TR1  (TRUNCATE requires the TABLE keyword). `TRUNCATE orders` is a
//	   SYNTAX_ERROR; only `TRUNCATE TABLE qualifiedName` is accepted, and the
//	   target is likewise bounded to 3 parts.
//
// FLAGGED divergence (ledger #5, "DML branch @branch"): Trino 481 accepts an
// optional `@ branch_name` after the target table in INSERT/DELETE/UPDATE/MERGE
// (e.g. `INSERT INTO cities @ audit VALUES (1, 'x')`). omni REJECTS it because
// the merged lexer node emits no token for '@' (it records an "unrecognized
// character" lex error), so the form cannot reach this parser intact. Adding
// '@' to the token model is a lexer-node concern outside this node's writes
// scope; until that lands, `@branch` stays a flagged, omni-rejects/Trino-accepts
// divergence and is excluded from this node's differential corpus. See the
// migration divergence ledger.

// ---------------------------------------------------------------------------
// INSERT
// ---------------------------------------------------------------------------

// InsertStmt is an `INSERT INTO qualifiedName [(col, …)] rootQuery` statement
// (the insertInto alternative). Target is the destination table (1-3 parts).
// Columns is the optional explicit column list (nil when absent). Source is the
// query supplying the rows — a SELECT, VALUES, TABLE, or WITH query (never nil;
// the source is mandatory, D-INS1).
type InsertStmt struct {
	Target  *ast.QualifiedName
	Columns []*ast.Identifier // explicit column list, nil when absent
	Source  *Query
	Loc     ast.Loc
}

// Tag implements ast.Node. See the file header for why DML statements use the
// placeholder T_Invalid rather than a dedicated tag.
func (n *InsertStmt) Tag() ast.NodeTag { return ast.T_Invalid }

// Span returns the source byte range.
func (n *InsertStmt) Span() ast.Loc { return n.Loc }

// ---------------------------------------------------------------------------
// TRUNCATE
// ---------------------------------------------------------------------------

// TruncateStmt is a `TRUNCATE TABLE qualifiedName` statement (the truncateTable
// alternative). Target is the table to truncate (1-3 parts). Trino's TRUNCATE
// has no column / partition / options clause — it is exactly these three
// keywords plus a name.
type TruncateStmt struct {
	Target *ast.QualifiedName
	Loc    ast.Loc
}

// Tag implements ast.Node.
func (n *TruncateStmt) Tag() ast.NodeTag { return ast.T_Invalid }

// Span returns the source byte range.
func (n *TruncateStmt) Span() ast.Loc { return n.Loc }

// Compile-time assertions that the statement types satisfy ast.Node.
var (
	_ ast.Node = (*InsertStmt)(nil)
	_ ast.Node = (*TruncateStmt)(nil)
)

// ---------------------------------------------------------------------------
// Parsers
// ---------------------------------------------------------------------------

// parseInsertStmt parses `INSERT INTO qualifiedName [(col, …)] rootQuery`. On
// entry cur is the INSERT keyword (not yet consumed).
//
// The optional `( identifier, … )` column list is disambiguated from a
// parenthesized source query (`INSERT INTO t (SELECT …)`): a '(' begins the
// column list only when it is followed by an identifier and then a ',' or ')'
// (i.e. `(a)` or `(a, b)`). A '(' followed by SELECT/WITH/TABLE/VALUES — or by
// an identifier that is NOT immediately closed/continued, which only a query
// can be — is the parenthesized source query, parsed by parseQuery.
func (p *Parser) parseInsertStmt() (ast.Node, error) {
	insertTok := p.advance() // consume INSERT
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	// D-INS2: target is at most catalog.schema.table.
	target, err := p.parseBoundedQualifiedName(3, "insert target")
	if err != nil {
		return nil, err
	}

	stmt := &InsertStmt{Target: target, Loc: ast.Loc{Start: insertTok.Loc.Start, End: target.Loc.End}}

	// Optional explicit column list. Only treat a leading '(' as the column
	// list when it is unambiguously `( ident [, …] )`; otherwise it opens a
	// parenthesized source query.
	if p.cur.Kind == int('(') && p.columnListFollows() {
		cols, _, err := p.parseColumnAliases()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// D-INS1: the source query is mandatory. parseQuery handles SELECT / VALUES
	// / TABLE / WITH … / a parenthesized query (and reports a leading
	// WITH FUNCTION inline routine as unsupported, deferred to parser-routines).
	src, err := p.parseQuery()
	if err != nil {
		return nil, err
	}
	stmt.Source = src
	stmt.Loc.End = src.Loc.End
	return stmt, nil
}

// columnListFollows reports whether the current '(' opens an INSERT column
// list (`( ident [, ident]* )`) rather than a parenthesized source query.
// It must NOT consume input. The decision uses two-token lookahead plus the
// fact that a column list's second token is always a non-keyword identifier:
//
//   - `(` then SELECT/WITH/TABLE/VALUES/'(' → a query, not a column list.
//   - `(` then an identifier → a column list ONLY if a one-element `(a)` or a
//     multi-element `(a, …)` list is syntactically possible. We can see only
//     the token after '(', so we accept the identifier case here and let the
//     subsequent parseColumnAliases enforce the `, ident` / `)` shape. A query
//     whose first token is an identifier necessarily continues with something
//     other than an immediate ')'/',' (e.g. `(a + b)` — a VALUES-less query is
//     not legal anyway), so this never misclassifies a real Trino query: the
//     only `( identifier …` query primaries are a subquery (which starts with
//     SELECT/WITH/TABLE/VALUES, handled above) — a bare `( expr )` is not a
//     valid INSERT source on its own.
//
// The leading '(' is the current token; we look at the token after it.
func (p *Parser) columnListFollows() bool {
	switch p.peekNext().Kind {
	case kwSELECT, kwWITH, kwTABLE, kwVALUES, int('('):
		// A parenthesized source query, never a column list.
		return false
	default:
		// `( identifier …` — treat as a column list. isIdentifierStart covers
		// unquoted, quoted, back-quoted and non-reserved-keyword identifiers.
		return isIdentifierStart(p.peekNext().Kind)
	}
}

// parseTruncateStmt parses `TRUNCATE TABLE qualifiedName`. On entry cur is the
// TRUNCATE keyword. D-TR1: the TABLE keyword is mandatory; the target is at
// most catalog.schema.table.
func (p *Parser) parseTruncateStmt() (ast.Node, error) {
	truncTok := p.advance() // consume TRUNCATE
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	target, err := p.parseBoundedQualifiedName(3, "truncate target")
	if err != nil {
		return nil, err
	}
	return &TruncateStmt{
		Target: target,
		Loc:    ast.Loc{Start: truncTok.Loc.Start, End: target.Loc.End},
	}, nil
}
