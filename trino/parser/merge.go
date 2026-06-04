package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-dml DAG node (with insert.go and
// update_delete.go): it implements Trino's MERGE statement and its WHEN clauses
// as hand-written recursive-descent parsers over the token stream.
//
// The statement structs live here in package parser and satisfy ast.Node with
// the placeholder tag ast.T_Invalid; see the insert.go file header for why DML
// statement nodes use T_Invalid rather than a dedicated trino/ast/nodetags.go
// tag.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	statement
//	    : ...
//	    | MERGE_ INTO_ qualifiedName (AS_? identifier)? USING_ relation
//	          ON_ expression mergeCase+                                            # merge
//	    ;
//	mergeCase
//	    : WHEN_ MATCHED_ (AND_ condition = expression)? THEN_ UPDATE_ SET_
//	          targets += identifier EQ_ values += expression
//	          (COMMA_ targets += identifier EQ_ values += expression)*             # mergeUpdate
//	    | WHEN_ MATCHED_ (AND_ condition = expression)? THEN_ DELETE_              # mergeDelete
//	    | WHEN_ NOT_ MATCHED_ (AND_ condition = expression)? THEN_ INSERT_
//	          (LPAREN_ targets += identifier (COMMA_ targets += identifier)* RPAREN_)?
//	          VALUES_ LPAREN_ values += expression (COMMA_ values += expression)* RPAREN_ # mergeInsert
//	    ;
//
// The implementation is adjudicated against the live Trino 481 oracle, not the
// literal legacy grammar. Oracle-confirmed facts (Trino 481) baked in:
//
//	D-MRG1 (at least one WHEN clause is required). `MERGE INTO t u USING s ON
//	   t.c = s.c` with no WHEN clause is a SYNTAX_ERROR (mergeCase+).
//	D-MRG2 (UPDATE SET is unparenthesized `col = expr` pairs). The Trino 481
//	   docs render the syntax as `THEN UPDATE SET ( column = expr [, …] )` but
//	   the oracle REJECTS the parenthesized form: `WHEN MATCHED THEN UPDATE SET
//	   (p = s.p)` is a SYNTAX_ERROR, while `WHEN MATCHED THEN UPDATE SET p = s.p`
//	   is accepted. The legacy grammar and the docs' own examples agree with the
//	   oracle (no parens), so the parenthesized doc syntax is treated as a doc
//	   error and NOT implemented (recorded as a flagged docs-vs-oracle item).
//	D-MRG3 (INSERT VALUES requires at least one expression). `THEN INSERT VALUES
//	   ()` is a SYNTAX_ERROR; the optional column list `( col, … )` may precede
//	   VALUES but the VALUES list itself is non-empty.
//	D-MRG4 (the target is at most catalog.schema.table; the USING source is a
//	   full `relation`). A 4-part MERGE target (`a.b.c.d`) is a SYNTAX_ERROR, but
//	   a 4-part USING-source table is accepted (the source goes through the
//	   relation rule, which does not bound the name) — so the target is parsed
//	   with the 3-part bound and the source with parseRelation, matching the FROM
//	   clause. A join source (`USING a JOIN b ON …`) is accepted; a top-level
//	   alias appended to a join (`USING a JOIN b ON … s`) is a SYNTAX_ERROR,
//	   which falls out naturally because parseRelation binds the alias to the
//	   inner sampledRelation, not the whole join.
//
// FLAGGED divergence (ledger #5): `@ branch_name` after the target table
// (`MERGE INTO accounts @ audit t USING …`) is accepted by Trino 481 but
// rejected by omni because the lexer emits no '@' token. See the insert.go
// file header.

// ---------------------------------------------------------------------------
// MERGE WHEN-clause AST
// ---------------------------------------------------------------------------

// MergeCaseKind classifies a MergeWhenClause (the three mergeCase
// alternatives).
type MergeCaseKind int

const (
	// MergeUpdate is `WHEN MATCHED [AND cond] THEN UPDATE SET col = expr, …`.
	MergeUpdate MergeCaseKind = iota
	// MergeDelete is `WHEN MATCHED [AND cond] THEN DELETE`.
	MergeDelete
	// MergeInsert is `WHEN NOT MATCHED [AND cond] THEN INSERT [(col, …)]
	// VALUES (expr, …)`.
	MergeInsert
)

// MergeWhenClause is one WHEN clause of a MERGE (the mergeCase rule). Kind
// selects the alternative. Condition is the optional `AND <expression>` guard
// (nil when absent). The remaining fields are populated per Kind:
//
//   - MergeUpdate: Assignments holds the `col = expr` SET pairs (non-empty).
//   - MergeDelete: no further fields.
//   - MergeInsert: Columns is the optional explicit column list (nil when the
//     `(col, …)` is absent); Values holds the VALUES expressions (non-empty).
type MergeWhenClause struct {
	Kind        MergeCaseKind
	Condition   Expr               // optional AND <condition>, nil when absent
	Assignments []UpdateAssignment // MergeUpdate: the SET col = expr pairs
	Columns     []*ast.Identifier  // MergeInsert: optional INSERT column list, nil when absent
	Values      []Expr             // MergeInsert: the VALUES expressions
	Loc         ast.Loc
}

// ---------------------------------------------------------------------------
// MERGE statement AST
// ---------------------------------------------------------------------------

// MergeStmt is a
//
//	MERGE INTO qualifiedName [[AS] alias] USING relation ON expression
//	    whenClause+
//
// statement (the merge alternative). Target is the destination table (1-3
// parts). Alias is the optional target alias (nil when absent). Source is the
// USING relation (a table, subquery, join, …). On is the merge-condition
// expression. Whens is the non-empty list of WHEN clauses.
type MergeStmt struct {
	Target *ast.QualifiedName
	Alias  *ast.Identifier // optional [AS] target alias, nil when absent
	Source Relation
	On     Expr
	Whens  []MergeWhenClause
	Loc    ast.Loc
}

// Tag implements ast.Node.
func (n *MergeStmt) Tag() ast.NodeTag { return ast.T_Invalid }

// Span returns the source byte range.
func (n *MergeStmt) Span() ast.Loc { return n.Loc }

// Compile-time assertion that *MergeStmt satisfies ast.Node.
var _ ast.Node = (*MergeStmt)(nil)

// ---------------------------------------------------------------------------
// Parsers
// ---------------------------------------------------------------------------

// parseMergeStmt parses
//
//	MERGE INTO qualifiedName [[AS] alias] USING relation ON expression
//	    mergeCase+
//
// On entry cur is the MERGE keyword.
func (p *Parser) parseMergeStmt() (ast.Node, error) {
	mergeTok := p.advance() // consume MERGE
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}

	// D-MRG4: target is at most catalog.schema.table.
	target, err := p.parseBoundedQualifiedName(3, "merge target")
	if err != nil {
		return nil, err
	}

	stmt := &MergeStmt{Target: target, Loc: ast.Loc{Start: mergeTok.Loc.Start, End: target.Loc.End}}

	// Optional target alias: `AS identifier` or a bare identifier. USING is a
	// reserved keyword (so isIdentifierStart(USING) is false): a bare identifier
	// here is therefore unambiguously the alias, never the USING clause.
	if p.cur.Kind == kwAS {
		p.advance() // consume AS — an identifier must follow
		alias, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Alias = alias
		stmt.Loc.End = alias.Loc.End
	} else if isIdentifierStart(p.cur.Kind) {
		alias := identFromToken(p.advance())
		stmt.Alias = alias
		stmt.Loc.End = alias.Loc.End
	}

	if _, err := p.expect(kwUSING); err != nil {
		return nil, err
	}

	// D-MRG4: the USING source is a full relation (table / subquery / join,
	// with optional alias and TABLESAMPLE), parsed exactly like a FROM item.
	source, err := p.parseRelation()
	if err != nil {
		return nil, err
	}
	stmt.Source = source
	stmt.Loc.End = source.Span().End

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	on, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	stmt.On = on
	stmt.Loc.End = on.Span().End

	// D-MRG1: at least one WHEN clause (mergeCase+).
	for p.cur.Kind == kwWHEN {
		clause, err := p.parseMergeWhenClause()
		if err != nil {
			return nil, err
		}
		stmt.Whens = append(stmt.Whens, clause)
		stmt.Loc.End = clause.Loc.End
	}
	if len(stmt.Whens) == 0 {
		// No WHEN clause at all: report the missing clause at the current token.
		return nil, p.syntaxErrorAtCur()
	}

	return stmt, nil
}

// parseMergeWhenClause parses one mergeCase. On entry cur is the WHEN keyword.
// The clause shape is decided by MATCHED vs NOT MATCHED and then by the THEN
// action keyword (UPDATE / DELETE for MATCHED; INSERT for NOT MATCHED).
func (p *Parser) parseMergeWhenClause() (MergeWhenClause, error) {
	whenTok := p.advance() // consume WHEN

	notMatched := false
	if p.cur.Kind == kwNOT {
		p.advance() // consume NOT
		notMatched = true
	}
	if _, err := p.expect(kwMATCHED); err != nil {
		return MergeWhenClause{}, err
	}

	clause := MergeWhenClause{Loc: ast.Loc{Start: whenTok.Loc.Start, End: p.prev.Loc.End}}

	// Optional `AND <condition>` guard.
	if p.cur.Kind == kwAND {
		p.advance() // consume AND
		cond, err := p.parseExpr()
		if err != nil {
			return MergeWhenClause{}, err
		}
		clause.Condition = cond
		clause.Loc.End = cond.Span().End
	}

	if _, err := p.expect(kwTHEN); err != nil {
		return MergeWhenClause{}, err
	}

	if notMatched {
		// WHEN NOT MATCHED … THEN INSERT … VALUES (…)
		if err := p.parseMergeInsertAction(&clause); err != nil {
			return MergeWhenClause{}, err
		}
		return clause, nil
	}

	// WHEN MATCHED … THEN (UPDATE SET … | DELETE)
	switch p.cur.Kind {
	case kwUPDATE:
		if err := p.parseMergeUpdateAction(&clause); err != nil {
			return MergeWhenClause{}, err
		}
	case kwDELETE:
		clause.Kind = MergeDelete
		delTok := p.advance() // consume DELETE
		clause.Loc.End = delTok.Loc.End
	default:
		// MATCHED clauses allow only UPDATE or DELETE.
		return MergeWhenClause{}, p.syntaxErrorAtCur()
	}
	return clause, nil
}

// parseMergeUpdateAction parses `UPDATE SET col = expr (, col = expr)*` into
// clause (a WHEN MATCHED … THEN UPDATE). On entry cur is the UPDATE keyword.
// D-MRG2: the SET list is unparenthesized `identifier = expression` pairs,
// reusing parseUpdateAssignment so MERGE and UPDATE share one assignment shape.
func (p *Parser) parseMergeUpdateAction(clause *MergeWhenClause) error {
	clause.Kind = MergeUpdate
	p.advance() // consume UPDATE
	if _, err := p.expect(kwSET); err != nil {
		return err
	}
	for {
		asgn, err := p.parseUpdateAssignment()
		if err != nil {
			return err
		}
		clause.Assignments = append(clause.Assignments, asgn)
		clause.Loc.End = asgn.Loc.End
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}
	return nil
}

// parseMergeInsertAction parses `INSERT [(col, …)] VALUES (expr, …)` into clause
// (a WHEN NOT MATCHED … THEN INSERT). On entry cur is the action keyword: it
// MUST be INSERT — a WHEN NOT MATCHED clause's only legal action is INSERT, so
// `THEN DELETE`/`THEN UPDATE`/`THEN <other>` is a SYNTAX_ERROR (Trino 481
// rejects `WHEN NOT MATCHED THEN DELETE VALUES (1)`). D-MRG3: the VALUES list
// is non-empty; the column list is optional.
func (p *Parser) parseMergeInsertAction(clause *MergeWhenClause) error {
	clause.Kind = MergeInsert
	if _, err := p.expect(kwINSERT); err != nil {
		return err
	}

	// Optional explicit column list `( identifier (, identifier)* )`.
	if p.cur.Kind == int('(') {
		cols, _, err := p.parseColumnAliases()
		if err != nil {
			return err
		}
		clause.Columns = cols
	}

	if _, err := p.expect(kwVALUES); err != nil {
		return err
	}
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	for {
		val, err := p.parseExpr()
		if err != nil {
			return err
		}
		clause.Values = append(clause.Values, val)
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return err
	}
	clause.Loc.End = closeTok.Loc.End
	return nil
}
