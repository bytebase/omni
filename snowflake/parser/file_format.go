package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// File-format DDL — CREATE / ALTER FILE FORMAT (T4.2)
// ---------------------------------------------------------------------------
//
// A named file format carries a TYPE ({ CSV | JSON | AVRO | ORC | PARQUET |
// XML }) plus a large, type-dependent and version-growing vocabulary of
// formatTypeOptions. Rather than mirror the legacy ANTLR grammar's
// format_type_options rule — already stale relative to the docs (it lacks
// USE_VECTORIZED_SCANNER, USE_LOGICAL_TYPE, MULTI_LINE, the per-type COMPRESSION
// variants, ...) — every option that follows the TYPE clause is parsed as an
// open-ended `KEY = <value>` pair (ast.CopyOption), reusing the merged COPY
// (T5.2) machinery (parseCopyOption / startsCopyOption). The `TYPE = <fmt>`
// clause is the one structural anchor: it selects which options are legal, so it
// is lifted out into a dedicated field. COMMENT is captured open-ended in
// Options. The catalog/semantic layer, not the parser, validates that an option
// is real and legal for the chosen TYPE. This mirrors the merged STAGE (T4.1) /
// COPY (T5.2) open-ended approach.
//
// FILE FORMAT lexes either as one FILE_FORMAT keyword token or as two separate
// FILE and FORMAT tokens (depending on intervening whitespace); both spellings
// are accepted, mirroring the DROP FILE FORMAT path (drop.go).

// ---------------------------------------------------------------------------
// CREATE FILE FORMAT
// ---------------------------------------------------------------------------

// parseCreateFileFormatStmt parses the body of a
//
//	CREATE [ OR REPLACE ] [ { TEMP | TEMPORARY | VOLATILE } ]
//	  FILE FORMAT [ IF NOT EXISTS ] <name>
//	  [ TYPE = { CSV | JSON | AVRO | ORC | PARQUET | XML } ]
//	  [ formatTypeOptions ]
//	  [ COMMENT = '<string_literal>' ]
//
// statement. The CREATE keyword and the optional OR REPLACE / TEMPORARY
// modifiers have already been consumed by parseCreateStmt; start is the Loc of
// the CREATE token, and cur is the FILE_FORMAT keyword (or the FILE keyword of
// the two-token spelling). temporary is set when any of TEMP / TEMPORARY /
// VOLATILE preceded FILE FORMAT.
func (p *Parser) parseCreateFileFormatStmt(start ast.Loc, orReplace, orAlter, temporary bool) (ast.Node, error) {
	if err := p.consumeFileFormatKeyword(); err != nil {
		return nil, err
	}

	stmt := &ast.CreateFileFormatStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Temporary: temporary,
		Loc:       ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
			stmt.IfNotExists = true
		}
	}

	// File-format name.
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Trailing clauses: the optional TYPE anchor, the open-ended formatTypeOptions,
	// and the COMMENT clause, in any order. The docs list TYPE first, but the loop
	// accepts TYPE anywhere and routes it to the dedicated Type field; every other
	// clause is an open-ended option. The loop ends at a statement boundary or any
	// token that does not begin an option.
	for p.startsCopyOption() {
		// TYPE = <fmt> is the structural anchor; capture it into Type, not Options.
		if p.cur.Type == kwTYPE {
			if err := p.parseFileFormatType(&stmt.Type, &stmt.TypeLoc); err != nil {
				return nil, err
			}
			continue
		}
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		stmt.Options = append(stmt.Options, opt)
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER FILE FORMAT
// ---------------------------------------------------------------------------

// parseAlterFileFormatStmt parses
//
//	ALTER FILE FORMAT [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER FILE FORMAT [ IF EXISTS ] <name> SET { [ formatTypeOptions ] [ COMMENT = '...' ] }
//
// The ALTER keyword has already been consumed; cur is the FILE_FORMAT keyword
// (or the FILE keyword of the two-token spelling). The SET form is
// unparenthesized (per the docs and the legacy ANTLR grammar). ALTER FILE FORMAT
// supports no UNSET and no SET TAG.
func (p *Parser) parseAlterFileFormatStmt() (ast.Node, error) {
	startLoc := p.cur.Loc // FILE_FORMAT / FILE keyword anchors Loc.Start (ALTER convention)
	if err := p.consumeFileFormatKeyword(); err != nil {
		return nil, err
	}
	stmt := &ast.AlterFileFormatStmt{Loc: ast.Loc{Start: startLoc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// File-format name.
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Action branch.
	switch p.cur.Type {
	case kwRENAME:
		// RENAME TO <new_name>
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterFileFormatRename
		stmt.NewName = newName

	case kwSET:
		// SET <formatTypeOptions> [ COMMENT = '...' ] — unparenthesized, open-ended
		// KEY = value params. A TYPE clause is not settable post-creation, so it is
		// captured (open-ended) like any other option rather than lifted out.
		p.advance() // consume SET
		var opts []*ast.CopyOption
		for p.startsCopyOption() {
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			opts = append(opts, opt)
		}
		if len(opts) == 0 {
			// SET with nothing settable is a syntax error.
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterFileFormatSet
		stmt.Options = opts

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared file-format helpers
// ---------------------------------------------------------------------------

// consumeFileFormatKeyword consumes the FILE FORMAT object keyword in either of
// its two lexings: the single FILE_FORMAT keyword token, or the two separate
// FILE and FORMAT tokens. On entry cur is the FILE_FORMAT or FILE keyword
// (callers dispatch on exactly those). Mirrors the DROP FILE FORMAT handling.
func (p *Parser) consumeFileFormatKeyword() error {
	switch p.cur.Type {
	case kwFILE_FORMAT:
		p.advance() // consume FILE_FORMAT (single keyword token)
		return nil
	case kwFILE:
		p.advance() // consume FILE
		if _, err := p.expect(kwFORMAT); err != nil {
			return err
		}
		return nil
	default:
		return p.syntaxErrorAtCur()
	}
}

// parseFileFormatType parses the structural `TYPE = { CSV | JSON | AVRO | ORC |
// PARQUET | XML }` clause, storing the uppercased format word into *typ and the
// value's source span into *loc. The format value is accepted as either a bare
// keyword (the documented form, TYPE = CSV) or a quoted string (TYPE = 'CSV'),
// which some real-world DDL uses; both normalize to the uppercased word. On
// entry cur is the TYPE keyword. The parser does not validate that the value is
// one of the six documented types — an unknown TYPE is left to the
// catalog/semantic layer, matching the open-ended option philosophy.
func (p *Parser) parseFileFormatType(typ *string, loc *ast.Loc) error {
	p.advance() // consume TYPE
	if _, err := p.expect('='); err != nil {
		return err
	}

	// Quoted form: TYPE = 'CSV'.
	if p.cur.Type == tokString {
		tok := p.advance()
		*typ = strings.ToUpper(tok.Str)
		*loc = tok.Loc
		return nil
	}

	// Bare-word form: TYPE = CSV (and any other keyword/identifier). Reuses the
	// single-token word-value reader so a value is always required. The value's
	// start is captured before reading, since parseValueWord may consume dotted
	// continuations.
	start := p.cur.Loc.Start
	words, end, err := p.parseValueWord()
	if err != nil {
		return err
	}
	*typ = words
	*loc = ast.Loc{Start: start, End: end}
	return nil
}
