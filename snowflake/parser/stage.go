package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Stage DDL — CREATE / ALTER STAGE (T4.1)
// ---------------------------------------------------------------------------
//
// Snowflake's stage grammar carries large, version-growing option vocabularies
// for both internal and external stages: the external cloud params (URL,
// STORAGE_INTEGRATION, CREDENTIALS, ENCRYPTION, ENDPOINT, AWS_ACCESS_POINT_ARN,
// USE_PRIVATELINK_ENDPOINT, ...), the directory-table params (DIRECTORY = (...)),
// the file-format params (FILE_FORMAT = (FORMAT_NAME = ... | TYPE = ... ...)) and
// the copy options (COPY_OPTIONS = (...)). Rather than mirror the legacy ANTLR
// grammar's finite, already-stale enumerations (its external_stage_params and
// stage_encryption_opts_* rules lack AWS_ACCESS_POINT_ARN, USE_PRIVATELINK_ENDPOINT,
// ENDPOINT, the s3compat:// form, ... all of which appear in the docs corpus),
// every param that follows the stage name is parsed as an open-ended
// `KEY = <value>` pair (ast.CopyOption), reusing the merged COPY (T5.2)
// machinery (parseCopyOption / parseCopyOptionParen). Only the structurally
// distinct WITH TAG and COMMENT clauses, and the ALTER action keywords
// (RENAME / SET / UNSET / REFRESH), anchor the grammar. The catalog/semantic
// layer, not the parser, validates that an option is real and legal. This
// mirrors the open-ended COPY / GRANT / SHOW approaches already merged.

// ---------------------------------------------------------------------------
// CREATE STAGE
// ---------------------------------------------------------------------------

// parseCreateStageStmt parses the body of a
//
//	CREATE [OR REPLACE] [TEMP|TEMPORARY] STAGE [IF NOT EXISTS] <name>
//	  <stageParams> [DIRECTORY=(...)] [FILE_FORMAT=(...)] [COPY_OPTIONS=(...)]
//	  [COMMENT='<string>'] [[WITH] TAG (<tag>='<value>' [,...])]
//
// statement. The CREATE keyword and the optional OR REPLACE / TEMPORARY
// modifiers have already been consumed by parseCreateStmt; start is the Loc of
// the CREATE token, and cur is the STAGE keyword.
func (p *Parser) parseCreateStageStmt(start ast.Loc, orReplace, temporary bool) (ast.Node, error) {
	p.advance() // consume STAGE

	stmt := &ast.CreateStageStmt{
		OrReplace: orReplace,
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

	// Stage name.
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Trailing clauses: open-ended cloud / directory / file-format / copy /
	// comment params (each a KEY = value pair) interleaved with the [WITH] TAG
	// clause. The docs place COMMENT before [WITH] TAG, while the legacy ANTLR
	// grammar lists `with_tags? comment_clause?` (TAG before COMMENT); accepting
	// either order (a) follows the docs, (b) regresses neither, and (c) avoids
	// silently dropping a COMMENT that trails the TAG clause. The loop ends at a
	// statement boundary or any token that begins neither an option nor a tag
	// clause.
	for {
		if p.startsStageTags() {
			tags, err := p.parseStageWithTags()
			if err != nil {
				return nil, err
			}
			stmt.Tags = append(stmt.Tags, tags...)
			continue
		}
		if p.startsStageOption() {
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
			continue
		}
		break
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER STAGE
// ---------------------------------------------------------------------------

// parseAlterStageStmt parses ALTER STAGE [IF EXISTS] <name> <action>.
// The ALTER keyword has already been consumed; cur is the STAGE keyword.
//
//	RENAME TO <new_name>
//	SET <options>
//	SET TAG <tag> = '<value>' [, ...]
//	UNSET TAG <tag> [, ...]
//	UNSET <property> [, ...]            (e.g. UNSET DCM PROJECT)
//	REFRESH [ SUBPATH = '<relative-path>' ]
func (p *Parser) parseAlterStageStmt() (ast.Node, error) {
	altTok := p.advance() // consume STAGE
	stmt := &ast.AlterStageStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Stage name.
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
		stmt.Action = ast.AlterStageRename
		stmt.NewName = newName

	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			// SET TAG <tag> = '<value>' [, ...]
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterStageSetTag
			stmt.Tags = tags
		} else {
			// SET <options> — open-ended KEY = value params.
			opts, err := p.parseStageOptions()
			if err != nil {
				return nil, err
			}
			if len(opts) == 0 {
				// SET with nothing settable is a syntax error.
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Action = ast.AlterStageSet
			stmt.Options = opts
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			// UNSET TAG <tag> [, ...]
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterStageUnsetTag
			stmt.UnsetTags = names
		} else {
			// UNSET <property> [, ...] (e.g. DCM PROJECT). Each property is an
			// open-ended run of name words (DCM PROJECT is two words); names are
			// not enumerated.
			props, err := p.parseStageUnsetProps()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterStageUnset
			stmt.UnsetProps = props
		}

	case kwREFRESH:
		// REFRESH [ SUBPATH = '<relative-path>' ]
		p.advance() // consume REFRESH
		stmt.Action = ast.AlterStageRefresh
		if p.curIsWord("SUBPATH") {
			p.advance() // consume SUBPATH
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return nil, err
			}
			s := tok.Str
			stmt.Subpath = &s
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared stage helpers
// ---------------------------------------------------------------------------

// parseStageOptions parses a run of zero or more stage params, each an
// open-ended `KEY = <value>` pair (reusing parseCopyOption). The run continues
// while the current token begins another option and terminates at a statement
// boundary (';' / EOF), at the WITH / TAG clause, or at any token that does not
// begin an option. WITH and TAG are explicitly excluded as option starters
// because they introduce the trailing WITH TAG clause, not a stage param.
func (p *Parser) parseStageOptions() ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for p.startsStageOption() {
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

// startsStageOption reports whether the current token can begin a stage param.
// It is the COPY option predicate minus WITH and TAG, which anchor the trailing
// [WITH] TAG clause and must not be swallowed as an option name.
func (p *Parser) startsStageOption() bool {
	switch p.cur.Type {
	case kwWITH, kwTAG:
		return false
	}
	return p.startsCopyOption()
}

// startsStageTags reports whether the current position begins a [WITH] TAG
// clause: either the TAG keyword directly, or WITH immediately followed by TAG
// (a WITH not followed by TAG is not a tag clause and is left for the caller).
func (p *Parser) startsStageTags() bool {
	if p.cur.Type == kwTAG {
		return true
	}
	return p.cur.Type == kwWITH && p.peekNext().Type == kwTAG
}

// parseStageWithTags parses a [WITH] TAG (...) clause and returns the tag
// assignments. It accepts both the WITH TAG and the bare TAG spellings
// (matching the docs' "[ WITH ] TAG (...)"). The caller must have confirmed
// startsStageTags; calling it otherwise yields a syntax error from
// parseTagAssignments.
func (p *Parser) parseStageWithTags() ([]*ast.TagAssignment, error) {
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH (peekNext == TAG guaranteed by startsStageTags)
	}
	return p.parseTagAssignments()
}

// parseTagSetList parses a SET TAG <tag> = '<value>' [, ...] list. Unlike the
// parenthesized WITH TAG (...) form, the SET TAG form is an unparenthesized
// comma-separated assignment list (per CREATE/ALTER ... SET TAG). The TAG
// keyword is consumed here.
func (p *Parser) parseTagSetList() ([]*ast.TagAssignment, error) {
	if _, err := p.expect(kwTAG); err != nil {
		return nil, err
	}

	var tags []*ast.TagAssignment
	for {
		tagName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		valueTok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &ast.TagAssignment{Name: tagName, Value: valueTok.Str})

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return tags, nil
}

// parseUnsetTagNameList parses an UNSET TAG <tag> [, ...] name list. Unlike the
// DATABASE/SCHEMA UNSET TAG (...) form, ALTER STAGE's UNSET TAG list is
// unparenthesized (per the docs). The TAG keyword is consumed here.
func (p *Parser) parseUnsetTagNameList() ([]*ast.ObjectName, error) {
	if _, err := p.expect(kwTAG); err != nil {
		return nil, err
	}

	var names []*ast.ObjectName
	for {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return names, nil
}

// parseStageUnsetProps parses an UNSET <property> [, ...] list, where each
// property is an open-ended run of one or more name words (e.g. the two-word
// DCM PROJECT, or a single COMMENT / URL / CREDENTIALS / ENCRYPTION /
// STORAGE_INTEGRATION). Property names are not enumerated; consecutive
// space-separated name words within one property are joined with a single
// space, and a comma separates properties. Consumes at least one property.
func (p *Parser) parseStageUnsetProps() ([]string, error) {
	var props []string
	for {
		prop, err := p.parseStageUnsetProp()
		if err != nil {
			return nil, err
		}
		props = append(props, prop)

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return props, nil
}

// parseStageUnsetProp parses a single UNSET property: one or more consecutive
// name words (identifier or keyword), joined with single spaces and uppercased.
// Multi-word properties such as DCM PROJECT are captured whole. The run stops at
// a comma, a statement boundary, or any non-word token.
func (p *Parser) parseStageUnsetProp() (string, error) {
	if !p.isOptionWord(p.cur.Type) {
		return "", p.syntaxErrorAtCur()
	}
	var b strings.Builder
	b.WriteString(strings.ToUpper(p.advance().Str))
	for p.isOptionWord(p.cur.Type) {
		b.WriteByte(' ')
		b.WriteString(strings.ToUpper(p.advance().Str))
	}
	return b.String(), nil
}
