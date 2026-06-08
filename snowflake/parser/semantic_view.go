package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Tag / semantic-view / dataset DDL — CREATE / ALTER (T4.9)
// ---------------------------------------------------------------------------
//
// Three governance / semantic-layer objects, all following the merged STAGE
// (T4.1) / FILE FORMAT (T4.2) / integration (T4.7) open-ended philosophy: the
// option vocabularies are large and version-growing, so every clause after a
// structural anchor is captured as an open-ended `KEY = <value>` CopyOption and
// the catalog/semantic layer (not the parser) validates legality.
//
//   - TAG: CREATE carries an ALLOWED_VALUES string list (the one anchor the docs
//     pin to first position) plus PROPAGATE / ON_CONFLICT / COMMENT options. The
//     legacy ANTLR create_tag rule modeled only `tag_allowed_values?
//     comment_clause?` and so lacks PROPAGATE / ON_CONFLICT — both present in the
//     official docs and the create-tag corpus (example_02). The docs (truth1)
//     win: trailing clauses are open-ended. ALTER TAG mirrors the docs' richer
//     surface (RENAME / ADD|DROP ALLOWED_VALUES / SET / UNSET multi-property /
//     SET|UNSET MASKING POLICY [FORCE] / UNSET DCM PROJECT), which exceeds the
//     legacy alter_tag_opts rule.
//   - SEMANTIC VIEW: the TABLES / RELATIONSHIPS / FACTS / DIMENSIONS / METRICS /
//     AI_VERIFIED_QUERIES sections are comma-separated definition lists whose
//     inner grammar is large and version-growing; each section's body is captured
//     as a balanced raw group rather than fully modeled. The scalar trailing
//     clauses (COMMENT, AI_SQL_GENERATION, AI_QUESTION_CATEGORIZATION) are
//     open-ended options; [WITH] TAG and COPY GRANTS are the remaining anchors.
//     The legacy create_semantic_view rule lacks the AI_* clauses and the
//     [WITH] TAG clause (docs win).
//   - DATASET: a newer ML object whose CREATE is, per both the docs and the
//     legacy grammar, just a name.

// ---------------------------------------------------------------------------
// CREATE TAG
// ---------------------------------------------------------------------------

// parseCreateTagStmt parses
//
//	CREATE [ OR REPLACE ] TAG [ IF NOT EXISTS ] <name>
//	  [ ALLOWED_VALUES '<v1>' [ , '<v2>' ... ] ]
//	  [ PROPAGATE = ... [ ON_CONFLICT = ... ] ]
//	  [ COMMENT = '<string_literal>' ]
//
// The CREATE keyword and the optional OR REPLACE modifier have already been
// consumed by parseCreateStmt; start is the Loc of the CREATE token, and cur is
// the TAG keyword.
func (p *Parser) parseCreateTagStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume TAG

	stmt := &ast.CreateTagStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseOptionalIfNotExists(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional ALLOWED_VALUES '<v>' [, ...] — the one anchor the docs require
	// first. PROPAGATE is NOT a starter for ALLOWED_VALUES, so it is captured as
	// an ordinary option below.
	if p.cur.Type == kwALLOWED_VALUES {
		p.advance() // consume ALLOWED_VALUES
		values, err := p.parseTagStringList()
		if err != nil {
			return nil, err
		}
		stmt.AllowedValues = values
	}

	// Trailing PROPAGATE / ON_CONFLICT / COMMENT clauses, each an open-ended
	// KEY = value option. The loop ends at a statement boundary or any token that
	// does not begin an option.
	for p.startsCopyOption() {
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
// ALTER TAG
// ---------------------------------------------------------------------------

// parseAlterTagStmt parses ALTER TAG [ IF EXISTS ] <name> <action>.
// The ALTER keyword has already been consumed; cur is the TAG keyword.
//
//	RENAME TO <new_name>
//	{ ADD | DROP } ALLOWED_VALUES '<v>' [ , ... ]
//	SET [ ALLOWED_VALUES '<v>' [ , ... ] ] [ PROPAGATE = ... [ ON_CONFLICT = ... ] ] [ COMMENT = '...' ]
//	SET MASKING POLICY <p> [ , MASKING POLICY <p2> ... ] [ FORCE ]
//	UNSET MASKING POLICY <p> [ , MASKING POLICY <p2> ... ]
//	UNSET { ALLOWED_VALUES | PROPAGATE | ON_CONFLICT | COMMENT | DCM PROJECT }
func (p *Parser) parseAlterTagStmt() (ast.Node, error) {
	altTok := p.advance() // consume TAG (anchors Loc.Start, ALTER convention)
	stmt := &ast.AlterTagStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	p.parseOptionalIfExists(&stmt.IfExists)

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

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
		stmt.Action = ast.AlterTagRename
		stmt.NewName = newName

	case kwADD, kwDROP:
		// { ADD | DROP } ALLOWED_VALUES '<v>' [ , ... ]
		add := p.cur.Type == kwADD
		p.advance() // consume ADD / DROP
		if _, err := p.expect(kwALLOWED_VALUES); err != nil {
			return nil, err
		}
		values, err := p.parseTagStringList()
		if err != nil {
			return nil, err
		}
		if add {
			stmt.Action = ast.AlterTagAddAllowedValues
		} else {
			stmt.Action = ast.AlterTagDropAllowedValues
		}
		stmt.AllowedValues = values

	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwMASKING {
			// SET MASKING POLICY <p> [ , MASKING POLICY <p2> ... ] [ FORCE ]
			policies, err := p.parseMaskingPolicyList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTagSetMaskingPolicy
			stmt.MaskingPolicies = policies
			if p.cur.Type == kwFORCE {
				p.advance() // consume FORCE
				stmt.Force = true
			}
			break
		}
		// SET [ ALLOWED_VALUES ... ] [ PROPAGATE ... ] [ COMMENT ... ]
		if err := p.parseAlterTagSet(stmt); err != nil {
			return nil, err
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwMASKING {
			// UNSET MASKING POLICY <p> [ , MASKING POLICY <p2> ... ]
			policies, err := p.parseMaskingPolicyList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterTagUnsetMaskingPolicy
			stmt.MaskingPolicies = policies
			break
		}
		// UNSET { ALLOWED_VALUES | PROPAGATE | ON_CONFLICT | COMMENT | DCM PROJECT }
		// — an open-ended property list (DCM PROJECT is two words), reusing the
		// STAGE UNSET-property reader.
		props, err := p.parseStageUnsetProps()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterTagUnset
		stmt.UnsetProps = props

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterTagSet parses the SET branch's options:
//
//	SET [ ALLOWED_VALUES '<v>' [ , ... ] ] [ PROPAGATE = ... [ ON_CONFLICT = ... ] ] [ COMMENT = '...' ]
//
// ALLOWED_VALUES is lifted into AllowedValues; every other clause is open-ended.
// SET with nothing settable is a syntax error. cur is positioned just past SET.
func (p *Parser) parseAlterTagSet(stmt *ast.AlterTagStmt) error {
	if p.cur.Type == kwALLOWED_VALUES {
		p.advance() // consume ALLOWED_VALUES
		values, err := p.parseTagStringList()
		if err != nil {
			return err
		}
		stmt.AllowedValues = values
	}

	var opts []*ast.CopyOption
	for p.startsCopyOption() {
		opt, err := p.parseCopyOption()
		if err != nil {
			return err
		}
		opts = append(opts, opt)
	}
	stmt.Options = opts

	if stmt.AllowedValues == nil && len(opts) == 0 {
		// SET with nothing settable is a syntax error.
		return p.syntaxErrorAtCur()
	}
	stmt.Action = ast.AlterTagSet
	return nil
}

// ---------------------------------------------------------------------------
// CREATE SEMANTIC VIEW
// ---------------------------------------------------------------------------

// parseCreateSemanticViewStmt parses
//
//	CREATE [ OR REPLACE ] SEMANTIC VIEW [ IF NOT EXISTS ] <name>
//	  TABLES ( ... )
//	  [ RELATIONSHIPS ( ... ) ] [ FACTS ( ... ) ] [ DIMENSIONS ( ... ) ]
//	  [ METRICS ( ... ) ] [ AI_VERIFIED_QUERIES ( ... ) ]
//	  [ COMMENT = '<string>' ] [ AI_SQL_GENERATION '<...>' ] [ AI_QUESTION_CATEGORIZATION '<...>' ]
//	  [ [ WITH ] TAG ( ... ) ] [ COPY GRANTS ]
//
// The CREATE / OR REPLACE tokens have already been consumed by parseCreateStmt;
// start is the Loc of the CREATE token, and cur is the SEMANTIC keyword.
func (p *Parser) parseCreateSemanticViewStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume SEMANTIC
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	stmt := &ast.CreateSemanticViewStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseOptionalIfNotExists(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Body clauses, in one order-tolerant loop: the parenthesized definition-list
	// sections (TABLES first per the docs, then RELATIONSHIPS / FACTS / DIMENSIONS
	// / METRICS / AI_VERIFIED_QUERIES), the scalar options (COMMENT /
	// AI_SQL_GENERATION / AI_QUESTION_CATEGORIZATION), the [WITH] TAG clause, and
	// COPY GRANTS. The docs fix a section order and place AI_VERIFIED_QUERIES after
	// the scalar AI_* options, so a single loop that accepts any of them in any
	// order (mirrors STAGE's order-tolerant tail) parses every documented ordering;
	// the semantic layer enforces the canonical order.
	for {
		if p.startsSemanticViewSection() {
			section, err := p.parseSemanticViewSection()
			if err != nil {
				return nil, err
			}
			stmt.Sections = append(stmt.Sections, section)
			continue
		}
		if p.startsSemanticViewTags() {
			tags, err := p.parseStageWithTags()
			if err != nil {
				return nil, err
			}
			stmt.Tags = append(stmt.Tags, tags...)
			continue
		}
		if p.startsCopyGrants() {
			p.advance() // consume COPY
			p.advance() // consume GRANTS
			stmt.CopyGrants = true
			continue
		}
		if p.startsSemanticViewOption() {
			opt, err := p.parseSemanticViewOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
			continue
		}
		break
	}
	if len(stmt.Sections) == 0 {
		// A SEMANTIC VIEW with no section is invalid (the docs require at least the
		// TABLES section). Rejecting here also prevents a bare section keyword such
		// as TABLES from being silently absorbed as the view name.
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// startsSemanticViewSection reports whether the current token begins a SEMANTIC
// VIEW section (TABLES / RELATIONSHIPS / FACTS / DIMENSIONS / METRICS /
// AI_VERIFIED_QUERIES). AI_VERIFIED_QUERIES is not a reserved keyword, so it
// lexes as an identifier and is matched by text.
func (p *Parser) startsSemanticViewSection() bool {
	switch p.cur.Type {
	case kwTABLES, kwRELATIONSHIPS, kwFACTS, kwDIMENSIONS, kwMETRICS:
		return true
	}
	return p.curIsWord("AI_VERIFIED_QUERIES")
}

// parseSemanticViewSection parses one `<KEYWORD> ( <body> )` section, capturing
// the body verbatim as a balanced raw group (the parentheses are not included in
// Body). cur is the section keyword.
func (p *Parser) parseSemanticViewSection() (*ast.SemanticViewSection, error) {
	kwTok := p.advance() // consume the section keyword
	section := &ast.SemanticViewSection{
		Keyword: strings.ToUpper(kwTok.Str),
		Loc:     ast.Loc{Start: kwTok.Loc.Start},
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Capture the inner body verbatim up to the matching ')'.
	bodyStart := p.cur.Loc.Start
	section.Body = p.captureBalancedRaw(bodyStart)
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	section.Loc.End = p.prev.Loc.End
	return section, nil
}

// startsSemanticViewTags reports whether the current position begins a [WITH]
// TAG clause (TAG directly, or WITH immediately followed by TAG).
func (p *Parser) startsSemanticViewTags() bool {
	if p.cur.Type == kwTAG {
		return true
	}
	return p.cur.Type == kwWITH && p.peekNext().Type == kwTAG
}

// startsSemanticViewOption reports whether the current token can begin a trailing
// scalar option (COMMENT / AI_SQL_GENERATION / AI_QUESTION_CATEGORIZATION). It is
// the COPY option predicate minus the structural keywords that anchor a trailing
// clause: WITH / TAG ([WITH] TAG) and COPY (COPY GRANTS).
func (p *Parser) startsSemanticViewOption() bool {
	switch p.cur.Type {
	case kwWITH, kwTAG, kwCOPY:
		return false
	}
	return p.startsCopyOption()
}

// parseSemanticViewOption parses one trailing scalar option. COMMENT uses the
// `KEY = value` shape; AI_SQL_GENERATION / AI_QUESTION_CATEGORIZATION use a bare
// `KEY '<string>'` shape (no '='), captured here as a CopyOption whose Lit holds
// the string. A `KEY = value` option falls through to the shared parseCopyOption.
func (p *Parser) parseSemanticViewOption() (*ast.CopyOption, error) {
	// AI_SQL_GENERATION / AI_QUESTION_CATEGORIZATION '<string>' — the value
	// follows the keyword directly with no '='. A string literal immediately after
	// the option name (with no intervening '=') marks this bare form; COMMENT = '…'
	// instead has '=' as the next token and so falls through to parseCopyOption.
	if p.peekNext().Type == tokString {
		nameTok := p.advance() // consume the option name
		valTok := p.advance()  // consume the string literal
		return &ast.CopyOption{
			Name: strings.ToUpper(nameTok.Str),
			Lit:  &ast.Literal{Kind: ast.LitString, Value: valTok.Str, Loc: valTok.Loc},
			Loc:  ast.Loc{Start: nameTok.Loc.Start, End: valTok.Loc.End},
		}, nil
	}
	return p.parseCopyOption()
}

// ---------------------------------------------------------------------------
// ALTER SEMANTIC VIEW
// ---------------------------------------------------------------------------

// parseAlterSemanticViewStmt parses ALTER SEMANTIC VIEW [ IF EXISTS ] <name>
// <action>. The ALTER keyword has already been consumed; cur is the SEMANTIC
// keyword.
//
//	RENAME TO <new_name>
//	SET COMMENT = '<string_literal>'
//	UNSET COMMENT
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
func (p *Parser) parseAlterSemanticViewStmt() (ast.Node, error) {
	startLoc := p.cur.Loc // SEMANTIC keyword anchors Loc.Start (ALTER convention)
	p.advance()           // consume SEMANTIC
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}
	stmt := &ast.AlterSemanticViewStmt{Loc: ast.Loc{Start: startLoc.Start}}

	p.parseOptionalIfExists(&stmt.IfExists)

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

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
		stmt.Action = ast.AlterSemanticViewRename
		stmt.NewName = newName

	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			// SET TAG <tag> = '<value>' [ , ... ]
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSemanticViewSetTag
			stmt.Tags = tags
		} else {
			// SET COMMENT = '<string>' — open-ended (only COMMENT is documented).
			var opts []*ast.CopyOption
			for p.startsCopyOption() {
				opt, err := p.parseCopyOption()
				if err != nil {
					return nil, err
				}
				opts = append(opts, opt)
			}
			if len(opts) == 0 {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Action = ast.AlterSemanticViewSet
			stmt.Options = opts
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			// UNSET TAG <tag> [ , ... ]
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSemanticViewUnsetTag
			stmt.UnsetTags = names
		} else {
			// UNSET COMMENT (and, permissively, any other property word).
			props, err := p.parseStageUnsetProps()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSemanticViewUnset
			stmt.UnsetProps = props
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE DATASET
// ---------------------------------------------------------------------------

// parseCreateDatasetStmt parses CREATE [ OR REPLACE ] DATASET [ IF NOT EXISTS ]
// <name>. The CREATE / OR REPLACE tokens have already been consumed; start is
// the Loc of the CREATE token, and cur is the DATASET keyword. DATASET carries
// no options (per both the docs and the legacy grammar).
//
// Divergence note: the official docs render IF NOT EXISTS *before* DATASET,
// whereas the legacy grammar and every other CREATE object in this engine place
// it *after* the object keyword. The post-keyword spelling is accepted here for
// consistency with the rest of the engine; the docs' pre-keyword spelling is a
// flagged divergence (likely a docs rendering quirk).
func (p *Parser) parseCreateDatasetStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume DATASET

	stmt := &ast.CreateDatasetStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseOptionalIfNotExists(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// parseTagStringList parses a comma-separated list of string literals:
// '<v1>' [ , '<v2>' ... ]. Consumes at least one. Used by ALLOWED_VALUES (CREATE
// TAG, ALTER TAG ADD/DROP/SET).
func (p *Parser) parseTagStringList() ([]string, error) {
	var values []string
	for {
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		values = append(values, tok.Str)
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return values, nil
}

// parseMaskingPolicyList parses a MASKING POLICY name list:
//
//	MASKING POLICY <p1> [ , MASKING POLICY <p2> ... ]
//
// Each entry repeats the MASKING POLICY keywords (per the docs and the legacy
// alter_tag_opts rule). cur is the first MASKING keyword. Consumes at least one.
func (p *Parser) parseMaskingPolicyList() ([]*ast.ObjectName, error) {
	var policies []*ast.ObjectName
	for {
		if _, err := p.expect(kwMASKING); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		policies = append(policies, name)
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return policies, nil
}
