package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/spanner/ast"
)

// --- PLACEMENT ---
//
// https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_placement
//
//	CREATE PLACEMENT placement_name
//	  [ OPTIONS ( instance_partition = 'partition_id'
//	      [, default_leader = 'leader_region_id' ]
//	      [, read_lease_regions = { 'region_name[, ...]' | null } ] ) ]
//
//	DROP PLACEMENT placement_name
//
// The pinned emulator also takes IF NOT EXISTS / IF EXISTS and a dotted name,
// and refuses a placement without a non-NULL instance_partition while it
// parses ("CREATE PLACEMENT statements require option `instance_partition` to
// be set"), so OPTIONS is required here. It does not know read_lease_regions
// yet; the documented option is accepted. The suite:spanner-ddl rows of
// testdata/compat.tsv hold its answers.

// placementOptions are the options CREATE PLACEMENT takes.
var placementOptions = map[string]spannerOptionKind{
	"instance_partition": spannerOptionString,
	"default_leader":     spannerOptionString,
	"read_lease_regions": spannerOptionString,
}

// parseCreatePlacement parses a CREATE PLACEMENT body after the shared CREATE
// prefix. cur is at the bare word `PLACEMENT`.
func (p *Parser) parseCreatePlacement(create Token, orReplace bool, scope string) (ast.Node, error) {
	if orReplace || scope != "" {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // PLACEMENT

	stmt := &ast.CreatePlacementStmt{}
	stmt.Loc.Start = create.Loc.Start

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if p.cur.Type != kwOPTIONS {
		return nil, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "syntax error: CREATE PLACEMENT requires OPTIONS (instance_partition = ...)",
		}
	}
	opts, err := p.parseDDLOptionsList()
	if err != nil {
		return nil, err
	}
	if err := p.checkSpannerOptions(opts, placementOptions); err != nil {
		return nil, err
	}
	if !hasNonNullOption(opts, "instance_partition") {
		return nil, &ParseError{
			Position: create.Loc.Start, End: p.prev.Loc.End,
			Message: "syntax error: CREATE PLACEMENT requires a non-NULL instance_partition option",
		}
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterPlacement parses an ALTER PLACEMENT body. The ALTER keyword has
// been consumed; cur is at the bare word `PLACEMENT`. The reference page has
// no ALTER PLACEMENT; the pinned emulator takes
//
//	ALTER PLACEMENT [IF EXISTS] name [SET OPTIONS ( placement_option [, ...] )]
//
// with CREATE PLACEMENT's options, none of them required, and a bare form.
func (p *Parser) parseAlterPlacement(alter Token) (ast.Node, error) {
	p.advance() // PLACEMENT

	stmt := &ast.AlterPlacementStmt{}
	stmt.Loc.Start = alter.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opts, err := p.parseSpannerSetOptions(placementOptions)
	if err != nil {
		return nil, err
	}
	stmt.SetOptions = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSpannerSetOptions parses an optional `SET OPTIONS ( … )` clause of a
// Spanner ALTER and checks its options against known. A SET followed by
// anything but OPTIONS is an error.
func (p *Parser) parseSpannerSetOptions(known map[string]spannerOptionKind) ([]*ast.OptionsEntry, error) {
	if p.cur.Type != kwSET {
		return nil, nil
	}
	p.advance() // SET
	if p.cur.Type != kwOPTIONS {
		return nil, p.syntaxErrorAtCur()
	}
	opts, err := p.parseDDLOptionsList()
	if err != nil {
		return nil, err
	}
	if err := p.checkSpannerOptions(opts, known); err != nil {
		return nil, err
	}
	return opts, nil
}

// parseDropPlacement parses a DROP PLACEMENT body. cur is at the bare word
// `PLACEMENT`.
func (p *Parser) parseDropPlacement(drop Token) (ast.Node, error) {
	p.advance() // PLACEMENT

	stmt := &ast.DropPlacementStmt{}
	stmt.Loc.Start = drop.Loc.Start

	ifExists, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	stmt.IfExists = ifExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// --- Spanner DDL option values ---

// spannerOptionKind is the value an option of a Spanner DDL statement takes.
type spannerOptionKind int

const (
	// spannerOptionString is a string literal or NULL.
	spannerOptionString spannerOptionKind = iota
	// spannerOptionStringArray is a non-empty `[ 'a' [, ...] ]` or NULL.
	spannerOptionStringArray
	// spannerOptionInt is an integer literal, optionally negated, or NULL.
	spannerOptionInt
)

// checkSpannerOptions holds an OPTIONS list to what Spanner's DDL parser takes
// for a statement whose options are listed in known. That parser reads an
// option's value as a literal, and refuses an unknown option name (matched
// case-sensitively), a repeated one, and a value of the wrong type while it
// parses: "Option: x is unknown.", "Duplicate option: x", "Unexpected value for
// option: x. Supported option values are strings and NULL."
func (p *Parser) checkSpannerOptions(opts []*ast.OptionsEntry, known map[string]spannerOptionKind) error {
	seen := make(map[string]bool, len(opts))
	for _, o := range opts {
		kind, ok := known[o.Name]
		if !ok {
			return optionError(o, fmt.Sprintf("syntax error: option %s is unknown", o.Name))
		}
		if seen[o.Name] {
			return optionError(o, fmt.Sprintf("syntax error: duplicate option %s", o.Name))
		}
		seen[o.Name] = true
		if o.Op != "=" || p.isBareProto(o.Value) || !spannerOptionValueIs(o.Value, kind) {
			return optionError(o, fmt.Sprintf("syntax error: unexpected value for option %s", o.Name))
		}
	}
	return nil
}

// isBareProto reports whether v is options_entry's bare PROTO keyword, which
// parseOptionsEntry holds as the string literal "PROTO", rather than a quoted
// string.
func (p *Parser) isBareProto(v ast.Node) bool {
	lit, ok := v.(*ast.Literal)
	if !ok || lit.Kind != ast.LitString {
		return false
	}
	s, e := absIndex(p, lit.Loc.Start), absIndex(p, lit.Loc.End)
	return s < e && strings.EqualFold(p.input[s:e], "PROTO")
}

func optionError(o *ast.OptionsEntry, msg string) *ParseError {
	return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: msg}
}

func spannerOptionValueIs(v ast.Node, kind spannerOptionKind) bool {
	if isLiteralKind(v, ast.LitNull) {
		return true
	}
	switch kind {
	case spannerOptionString:
		return isLiteralKind(v, ast.LitString)
	case spannerOptionStringArray:
		arr, ok := v.(*ast.ArrayExpr)
		if !ok || arr.HasArrayKeyword || arr.ElemType != nil || len(arr.Elements) == 0 {
			return false
		}
		for _, e := range arr.Elements {
			if !isLiteralKind(e, ast.LitString) {
				return false
			}
		}
		return true
	case spannerOptionInt:
		if u, ok := v.(*ast.UnaryExpr); ok && u.Op == ast.UnaryMinus {
			v = u.Expr
		}
		return isLiteralKind(v, ast.LitInt)
	}
	return false
}

func isLiteralKind(v ast.Node, kind ast.LiteralKind) bool {
	lit, ok := v.(*ast.Literal)
	return ok && lit.Kind == kind
}

// hasNonNullOption reports whether opts sets name to a value other than NULL.
func hasNonNullOption(opts []*ast.OptionsEntry, name string) bool {
	for _, o := range opts {
		if o.Name == name && !isLiteralKind(o.Value, ast.LitNull) {
			return true
		}
	}
	return false
}
