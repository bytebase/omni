package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// COMMENT ON / TRUNCATE (T6.3)
// ---------------------------------------------------------------------------
//
// Like GRANT/REVOKE and SHOW, COMMENT's object-type surface grows over time, so
// the object type is captured as an open-ended uppercased token run rather than
// a fixed enum. The COLUMN form is special-cased to match the legacy grammar's
// distinct `COMMENT ... ON COLUMN full_column_name` rule.
//
// NOTE on ABORT: the legacy grammar has no top-level ABORT statement — ABORT
// appears only as the tail of ALTER WAREHOUSE ("... ABORT ALL QUERIES"), which
// is owned by the warehouse-alter node, and as a non-reserved keyword. parseStmt
// has no kwABORT dispatch case. ABORT is therefore intentionally NOT implemented
// here; see the PR body for the rationale.

// ---------------------------------------------------------------------------
// COMMENT
// ---------------------------------------------------------------------------

// parseCommentStmt parses a COMMENT statement:
//
//	COMMENT [IF EXISTS] ON <object_type> <name> [ ( <signature> ) ] IS '<string>'
//	COMMENT [IF EXISTS] ON COLUMN <column_name> IS '<string>'
//
// The COMMENT keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseCommentStmt() (ast.Node, error) {
	commentTok := p.advance() // consume COMMENT
	start := commentTok.Loc

	stmt := &ast.CommentStmt{Loc: ast.Loc{Start: start.Start}}

	// Optional IF EXISTS.
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		stmt.IfExists = true
	}

	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}

	if p.cur.Type == kwCOLUMN {
		// COMMENT ON COLUMN <column_name> IS '<string>'. The legacy
		// full_column_name rule allows up to four dotted parts
		// (db.schema.table.column), which exceeds ObjectName's 3-part shape, so a
		// ColumnRef holds the parts.
		p.advance() // consume COLUMN
		stmt.IsColumn = true
		stmt.ObjectType = "COLUMN"
		col, err := p.parseColumnRef()
		if err != nil {
			return nil, err
		}
		stmt.Column = col
	} else {
		// COMMENT ON <object_type> <name> [ ( signature ) ] IS '<string>'
		objType, name, err := p.parseTypedObject()
		if err != nil {
			return nil, err
		}
		stmt.ObjectType = objType
		stmt.Name = name

		// Optional FUNCTION/PROCEDURE argument-type signature.
		if p.cur.Type == '(' {
			sig, err := p.parseGrantSignature()
			if err != nil {
				return nil, err
			}
			stmt.Signature = sig
		}
	}

	if _, err := p.expect(kwIS); err != nil {
		return nil, err
	}
	str, err := p.expect(tokString)
	if err != nil {
		return nil, err
	}
	stmt.Comment = str.Str

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseTypedObject parses an open-ended "<object_type> <name>" pair where the
// type may be multiple words (TABLE, MASKING POLICY, ROW ACCESS POLICY, ...)
// and the name is the final space-separated unit. It mirrors the technique in
// grant_revoke.go's parseGrantTargetObject: read name units left-to-right,
// shifting each previous unit into the type, so the LAST unit is the object name
// and all preceding units form the (multi-word) type. This needs no fixed type
// vocabulary, so new object types parse without code changes.
//
// The run terminates at IS, '(' (a signature), or any non-name token. Requires
// at least a type word and a name word.
func (p *Parser) parseTypedObject() (string, *ast.ObjectName, error) {
	var typeWords []string
	name, err := p.parseObjectName()
	if err != nil {
		return "", nil, err
	}

	for p.startsTypedObjectNameUnit() {
		typeWords = append(typeWords, strings.ToUpper(name.String()))
		name, err = p.parseObjectName()
		if err != nil {
			return "", nil, err
		}
	}

	if len(typeWords) == 0 {
		return "", nil, &ParseError{
			Loc: name.Loc,
			Msg: "expected object type before object name",
		}
	}
	return strings.Join(typeWords, " "), name, nil
}

// startsTypedObjectNameUnit reports whether the current token can begin another
// space-separated name unit inside a "<type> <name>" run — i.e. it is a name
// word and not a clause terminator (IS, '(', or EOF).
func (p *Parser) startsTypedObjectNameUnit() bool {
	switch p.cur.Type {
	case kwIS, '(', tokEOF:
		return false
	}
	return p.isObjectTypeWord(p.cur.Type)
}

// parseColumnRef parses a dotted column reference of 1 to 4 parts, mirroring the
// legacy full_column_name rule ([db.][schema.][table.]column). Each part is
// parsed with parseNamePart, which accepts keywords (e.g. reserved words used as
// names) as well as bare and quoted identifiers. A reference with more than four
// parts is rejected (no documented Snowflake column reference exceeds
// db.schema.table.column).
func (p *Parser) parseColumnRef() (*ast.ColumnRef, error) {
	first, err := p.parseNamePart()
	if err != nil {
		return nil, err
	}
	parts := []ast.Ident{first}
	for p.cur.Type == '.' {
		dotLoc := p.cur.Loc
		p.advance() // consume '.'
		if len(parts) >= 4 {
			return nil, &ParseError{
				Loc: dotLoc,
				Msg: "column reference has more than 4 parts (max is db.schema.table.column)",
			}
		}
		part, err := p.parseNamePart()
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return &ast.ColumnRef{
		Parts: parts,
		Loc:   ast.Loc{Start: parts[0].Loc.Start, End: parts[len(parts)-1].Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// TRUNCATE
// ---------------------------------------------------------------------------

// parseTruncateStmt parses a TRUNCATE statement:
//
//	TRUNCATE [TABLE] [IF EXISTS] <name>
//	TRUNCATE MATERIALIZED VIEW <name>
//
// The TRUNCATE keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseTruncateStmt() (ast.Node, error) {
	truncateTok := p.advance() // consume TRUNCATE
	start := truncateTok.Loc

	stmt := &ast.TruncateStmt{Loc: ast.Loc{Start: start.Start}}

	switch p.cur.Type {
	case kwMATERIALIZED:
		// TRUNCATE MATERIALIZED VIEW <name> — no IF EXISTS in the legacy grammar.
		p.advance() // consume MATERIALIZED
		if _, err := p.expect(kwVIEW); err != nil {
			return nil, err
		}
		stmt.MaterializedView = true
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil

	case kwTABLE:
		// TRUNCATE TABLE [IF EXISTS] ...
		p.advance() // consume TABLE
		stmt.IfExists = p.tryIfExists()

	default:
		// TRUNCATE [IF EXISTS] ... — TABLE keyword omitted.
		stmt.IfExists = p.tryIfExists()
	}

	// TRUNCATE [TABLE] [IF EXISTS] ERROR_TABLE( <base_table_name> ) — documented
	// by Snowflake (absent from the legacy grammar). ERROR_TABLE is not a
	// keyword; it is recognized as an identifier whose text is ERROR_TABLE and
	// which is immediately followed by '('.
	if p.startsErrorTable() {
		p.advance() // consume ERROR_TABLE
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.ErrorTable = true
		stmt.Name = name
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// startsErrorTable reports whether the current token begins the ERROR_TABLE(...)
// form: an unquoted identifier spelled "ERROR_TABLE" (case-insensitive) followed
// by '('.
func (p *Parser) startsErrorTable() bool {
	return p.cur.Type == tokIdent &&
		strings.EqualFold(p.cur.Str, "ERROR_TABLE") &&
		p.peekNext().Type == '('
}

// tryIfExists consumes an IF EXISTS pair if present and reports whether it was.
func (p *Parser) tryIfExists() bool {
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		return true
	}
	return false
}
