package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// What MariaDB refuses of a query while it parses, past what the grammar says. The server
// adds each table to its list in an action, which is where a repeated name is 1066, a
// multi-table DELETE's targets are matched to the list with 1109, and a LIMIT under IN is
// 1235; a syntax error after the statement does not change the answer.
//
// Ref: st_select_lex::add_table_to_list, st_select_lex::add_item_to_list,
// multi_delete_set_locks_and_link_aux_tables and the IN and ALL/ANY arms of predicate in
// MariaDB 11.8 sql/sql_parse.cc and sql/sql_yacc.yy.

// fromKind tells what kind of entry a table list holds, which decides what it clashes with.
type fromKind int

const (
	fromTable fromKind = iota
	fromDerived
	fromJSONTable
)

type fromEntry struct {
	kind   fromKind
	schema string
	name   string
	alias  bool
	pos    int
}

// fromEntries flattens a table list into its tables, through joins and parenthesized
// groups and not into derived tables, which have lists of their own. The list can be as
// deep as the statement is long, so it keeps its own stack.
func fromEntries(refs []nodes.TableExpr) []fromEntry {
	var entries []fromEntry
	stack := make([]nodes.TableExpr, 0, len(refs))
	for i := len(refs) - 1; i >= 0; i-- {
		stack = append(stack, refs[i])
	}
	for len(stack) > 0 {
		ref := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		switch r := ref.(type) {
		case *nodes.JoinClause:
			stack = append(stack, r.Right, r.Left)
		case *nodes.TableRef:
			entry := fromEntry{kind: fromTable, schema: r.Schema, name: r.Name, pos: r.Loc.Start}
			if r.Alias != "" {
				entry.name = r.Alias
				entry.alias = true
			}
			entries = append(entries, entry)
		case *nodes.SubqueryExpr:
			entries = append(entries, fromEntry{kind: fromDerived, name: r.Alias, alias: true, pos: r.Loc.Start})
		case *nodes.JsonTableExpr:
			entries = append(entries, fromEntry{kind: fromJSONTable, name: r.Alias, alias: true, pos: r.Loc.Start})
		}
	}
	return entries
}

// refuseDuplicateTables is 1066, ER_NONUNIQ_TABLE: two entries of one list under one name,
// compared exactly, which is the answer the server gives unless lower_case_table_names has it
// fold them too. A table's database is part of its name, so `d.t, e.t` and `t, d.t` are
// two tables, and a derived table has none, so `t, (SELECT 1) t` is two entries as well. A
// JSON_TABLE clashes with anything of its alias.
func refuseDuplicateTables(entries []fromEntry) error {
	type key struct {
		schema string
		name   string
		kind   fromKind
	}
	seen := make(map[key]bool, len(entries))
	json := make(map[string]bool)
	names := make(map[string]bool, len(entries))
	for _, e := range entries {
		k := key{schema: e.schema, name: e.name, kind: e.kind}
		clash := names[e.name] && e.kind == fromJSONTable || json[e.name]
		if e.kind != fromJSONTable && seen[k] {
			clash = true
		}
		if clash {
			return &ParseError{Message: "Not unique table/alias: '" + e.name + "'", Position: e.pos}
		}
		seen[k] = true
		names[e.name] = true
		if e.kind == fromJSONTable {
			json[e.name] = true
		}
	}
	return nil
}

// refuseFromDuplicates holds a statement's table list to refuseDuplicateTables.
func refuseFromDuplicates(refs []nodes.TableExpr) error {
	return refuseDuplicateTables(fromEntries(refs))
}

// refuseDeleteTargets holds a multi-table DELETE's targets to its table list: 1066 for a
// target named twice, and 1109, ER_UNKNOWN_TABLE, for one the list has no entry for. A target
// matches an entry with an alias by the alias alone, written with no database, and an entry
// with none by its name and its database. A database left out is the current one, which the parser
// can't see, so a target or an entry with none may match either way and is left to the server:
// with d current, DELETE t FROM d.t and DELETE d.t FROM t run, and DELETE t FROM x.t answers 1109
// (mariadb:11.8.8). A JSON_TABLE is no target. The server folds
// case here when lower_case_table_names does, which the parser can't see, so case is left out
// and `DELETE T FROM t` reads whichever way the server is set.
func refuseDeleteTargets(targets, using []nodes.TableExpr) error {
	if len(using) == 0 {
		return nil
	}
	list := fromEntries(targets)
	if err := refuseDuplicateTables(list); err != nil {
		return err
	}
	from := fromEntries(using)
	for _, target := range list {
		found := false
		for _, e := range from {
			if e.kind == fromJSONTable || !strings.EqualFold(e.name, target.name) {
				continue
			}
			sameSchema := target.schema == "" || e.schema == "" || strings.EqualFold(e.schema, target.schema)
			if e.alias && target.schema == "" || !e.alias && sameSchema {
				found = true
				break
			}
		}
		if !found {
			return &ParseError{Message: "Unknown table '" + target.name + "' in MULTI DELETE", Position: target.pos}
		}
	}
	return nil
}

// refuseLimitedSubquery is 1235, ER_NOT_SUPPORTED_YET: a LIMIT on the query under IN, NOT IN,
// ANY, SOME or ALL, or on any branch of it. A query nested further down has its own say.
func refuseLimitedSubquery(sel *nodes.SelectStmt, pos int) error {
	stack := []*nodes.SelectStmt{sel}
	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if s == nil {
			continue
		}
		if s.Limit != nil && s.Limit.Count != nil {
			return &ParseError{
				Message:  "This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'",
				Position: pos,
			}
		}
		stack = append(stack, s.Left, s.Right, s.ParenSource)
	}
	return nil
}

// refuseSequenceCall holds NEXTVAL, LASTVAL and SETVAL to the arms the grammar gives them: the
// first argument is a table_ident, held to its checks, and SETVAL's second a number, signed
// to a 64-bit minimum. “ `nextval`(1) “ is an ordinary call and does not come here.
//
// Ref: function_call_conflict's NEXTVAL_SYM, LASTVAL_SYM and SETVAL_SYM arms in MariaDB 11.8
// sql/sql_yacc.yy.
func (p *Parser) refuseSequenceCall(fc *nodes.FuncCallExpr) error {
	switch fc.Name {
	case "NEXTVAL", "LASTVAL":
		if len(fc.Args) != 1 {
			return &ParseError{Message: "syntax error: " + fc.Name + " takes a sequence", Position: fc.Loc.Start}
		}
	case "SETVAL":
		if len(fc.Args) < 2 {
			return &ParseError{Message: "syntax error: SETVAL takes a sequence and a value", Position: fc.Loc.Start}
		}
		switch v := fc.Args[1].(type) {
		case *nodes.IntLit:
		case *nodes.UnaryExpr:
			// The magnitude is read off the written digits, since the literal's Value stops at
			// math.MaxInt64 and so cannot tell -9223372036854775808, which reads, from what does not.
			lit, ok := v.Operand.(*nodes.IntLit)
			if !ok || v.Op != nodes.UnaryMinus || !fitsNegatedInt64(p.inputText(lit.Loc.Start, lit.Loc.End)) {
				return &ParseError{Message: "syntax error: SETVAL takes a number", Position: fc.Loc.Start}
			}
		default:
			return &ParseError{Message: "syntax error: SETVAL takes a number", Position: fc.Loc.Start}
		}
	default:
		return nil
	}
	ref, ok := fc.Args[0].(*nodes.ColumnRef)
	if !ok || ref.Star || ref.Schema != "" {
		return &ParseError{Message: "syntax error: " + fc.Name + " takes a sequence", Position: fc.Loc.Start}
	}
	if ref.Table != "" {
		if err := refuseDatabaseName(ref.Table, fc.Loc.Start); err != nil {
			return err
		}
	}
	if identifierNameRefused(ref.Column) {
		return &ParseError{Message: "Incorrect table name '" + ref.Column + "'", Position: fc.Loc.Start}
	}
	return nil
}

// fitsNegatedInt64 reports whether the digits, with a minus sign, are an int64: at most
// 9223372036854775808.
func fitsNegatedInt64(digits string) bool {
	digits = strings.TrimLeft(digits, "0")
	const limit = "9223372036854775808"
	if len(digits) != len(limit) {
		return len(digits) < len(limit)
	}
	return digits <= limit
}
