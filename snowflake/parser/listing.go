package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Listings, read as ast.ObjectStmt: CREATE EXTERNAL LISTING, CREATE
// ORGANIZATION LISTING, and ALTER and DROP LISTING, which act on both
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"EXTERNAL", "LISTING"},
		create: parseCreateExternalListing,
		// https://docs.snowflake.com/en/sql-reference/sql/create-listing
		nameParts: 1,
		props:     externalListingProps,
		unset:     []string{},
	})
	registerObject(&objectForm{
		words:  []string{"ORGANIZATION", "LISTING"},
		create: parseCreateOrganizationListing,
		// https://docs.snowflake.com/en/sql-reference/sql/create-organization-listing
		nameParts: 1,
		props:     objectProps{"PUBLISH": propBool},
		unset:     []string{},
	})
	registerObject(&objectForm{
		words: []string{"LISTING"},
		alter: parseAlterListing,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-listing
		drop:           true,
		dropNoIfExists: true,
		nameParts:      1,
		props:          externalListingProps,
		unset:          []string{},
	})
}

// externalListingProps are the properties CREATE EXTERNAL LISTING's AS form
// and ALTER LISTING's AS form document.
var externalListingProps = objectProps{
	"PUBLISH": propBool,
	"REVIEW":  propBool,
	"COMMENT": propString,
}

// parseCreateExternalListing reads what follows the name of
//
//	CREATE EXTERNAL LISTING [ IF NOT EXISTS ] <name>
//	  [ { SHARE <share_name> | APPLICATION PACKAGE <package_name> } ]
//	  AS '<yaml_manifest_string>'
//	  [ PUBLISH = { TRUE | FALSE } ] [ REVIEW = { TRUE | FALSE } ] [ COMMENT = '<string>' ]
//	CREATE EXTERNAL LISTING [ IF NOT EXISTS ] <name>
//	  [ { SHARE <share_name> | APPLICATION PACKAGE <package_name> } ]
//	  FROM '<yaml_manifest_stage_location>'
//	  [ PUBLISH = { TRUE | FALSE } ] [ REVIEW = { TRUE | FALSE } ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-listing
func parseCreateExternalListing(p *Parser, s *ast.ObjectStmt) error {
	return p.parseCreateListing(s, true)
}

// parseCreateOrganizationListing reads what follows the name of
//
//	CREATE ORGANIZATION LISTING [ IF NOT EXISTS ] <name>
//	  [ { SHARE <share_name> | APPLICATION PACKAGE <package_name> } ]
//	  { AS '<yaml_manifest_string>' | FROM '<yaml_manifest_stage_location>' }
//	  [ PUBLISH = { TRUE | FALSE } ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-organization-listing
func parseCreateOrganizationListing(p *Parser, s *ast.ObjectStmt) error {
	return p.parseCreateListing(s, false)
}

// parseCreateListing reads a listing's source, its manifest and properties.
// external says whether it is an EXTERNAL LISTING, whose FROM form documents
// no COMMENT; an ORGANIZATION LISTING documents PUBLISH alone (its form's
// props).
func (p *Parser) parseCreateListing(s *ast.ObjectStmt, external bool) error {
	start := p.cur.Loc.Start
	switch {
	case p.cur.Type == kwSHARE:
		p.advance() // consume SHARE
		name, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "SHARE", start).Names = []*ast.ObjectName{name}
	case p.acceptWords("APPLICATION", "PACKAGE"):
		name, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "APPLICATION PACKAGE", start).Names = []*ast.ObjectName{name}
	}
	from := p.cur.Type == kwFROM
	switch p.cur.Type {
	case kwAS:
		if err := p.parseObjectAsString(s); err != nil {
			return err
		}
	case kwFROM:
		if err := p.parseObjectKeywordString(s, "FROM"); err != nil {
			return err
		}
	default:
		return p.syntaxErrorAtCur()
	}
	props, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	if err := rejectObjectTags(s); err != nil {
		return err
	}
	if external && from && props != nil {
		for _, o := range props.Options {
			if o.Name == "COMMENT" {
				return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error: FROM takes no COMMENT"}
			}
		}
	}
	return nil
}

// parseObjectKeywordString reads <keyword> '<string>' into a clause holding
// the string as a nameless option. cur is the keyword.
func (p *Parser) parseObjectKeywordString(s *ast.ObjectStmt, keyword string) error {
	start := p.cur.Loc.Start
	p.advance() // consume the keyword
	if p.cur.Type != tokString {
		return p.syntaxErrorAtCur()
	}
	lit, err := p.parseCopyLiteral()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, keyword, start).Options = []*ast.CopyOption{{Lit: lit, Loc: lit.Loc}}
	return nil
}

// parseAlterListing reads what follows the name of
//
//	ALTER LISTING [ IF EXISTS ] <name> [ { PUBLISH | UNPUBLISH | REVIEW } ]
//	ALTER LISTING [ IF EXISTS ] <name> AS '<yaml_manifest_string>'
//	  [ PUBLISH = { TRUE | FALSE } ] [ REVIEW = { TRUE | FALSE } ] [ COMMENT = '<string>' ]
//	ALTER LISTING <name> ADD VERSION [ [ IF NOT EXISTS ] <version_name> ]
//	  FROM <yaml_manifest_stage_location> [ COMMENT = '<string>' ]
//	ALTER LISTING [ IF EXISTS ] <name> { ADD | REMOVE } TARGETS <manifest>
//	ALTER LISTING [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER LISTING [ IF EXISTS ] <name> SET COMMENT = '<string>'
//
// The page's examples write the stage location as @<stage>[/<path>] and the
// manifest as a string ($$...$$).
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-listing
func parseAlterListing(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("PUBLISH"):
		p.addObjectKeywordsClause(s, "PUBLISH", start)
		return nil
	case p.acceptWords("UNPUBLISH"):
		p.addObjectKeywordsClause(s, "UNPUBLISH", start)
		return nil
	case p.acceptWords("REVIEW"):
		p.addObjectKeywordsClause(s, "REVIEW", start)
		return nil
	case p.cur.Type == kwAS:
		if err := p.parseObjectAsString(s); err != nil {
			return err
		}
		if _, err := p.parseObjectProperties(s, "", false); err != nil {
			return err
		}
		return rejectObjectTags(s)
	case p.atWords("ADD", "TARGETS"), p.atWords("REMOVE", "TARGETS"):
		verb := p.advance().Str
		p.advance() // consume TARGETS
		if p.cur.Type != tokString {
			return p.syntaxErrorAtCur()
		}
		lit, err := p.parseCopyLiteral()
		if err != nil {
			return err
		}
		c := p.addObjectKeywordsClause(s, strings.ToUpper(verb)+" TARGETS", start)
		c.Options = []*ast.CopyOption{{Lit: lit, Loc: lit.Loc}}
		return nil
	case p.cur.Type == kwRENAME:
		return p.parseObjectRenameSimple(s)
	case p.cur.Type == kwSET:
		// SET COMMENT = '<string>' alone.
		p.advance() // consume SET
		if !p.curIsWord("COMMENT") || !p.startsObjectProperty() {
			return p.syntaxErrorAtCur()
		}
		opt, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		s.Clauses = append(s.Clauses, &ast.ObjectClause{
			Keywords: "SET", Options: []*ast.CopyOption{opt},
			Loc: ast.Loc{Start: start, End: p.prev.Loc.End},
		})
		return nil
	case p.atWords("ADD", "VERSION"):
		if s.IfExists {
			// The page writes no IF EXISTS on ADD VERSION.
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume ADD
		p.advance() // consume VERSION
		if err := p.parseObjectAddVersion(s, "ADD VERSION", start); err != nil {
			return err
		}
		if err := p.parseObjectFromStage(s); err != nil {
			return err
		}
		props, err := p.parseObjectProperties(s, "", false)
		if err != nil {
			return err
		}
		if err := rejectObjectTags(s); err != nil {
			return err
		}
		if props != nil && (len(props.Options) != 1 || props.Options[0].Name != "COMMENT") {
			return &ParseError{Position: props.Loc.Start, End: props.Loc.End, Message: "syntax error: ADD VERSION takes COMMENT alone"}
		}
		return nil
	case p.cur.Type == tokEOF || p.cur.Type == ';':
		// ALTER LISTING <name> with no action: the page writes the action
		// in brackets.
		return nil
	}
	return p.syntaxErrorAtCur()
}

// startsListingRevision reports whether cur begins DESC LISTING's
// REVISION = ...: the word REVISION followed by '='. It consumes nothing.
func (p *Parser) startsListingRevision() bool {
	return isWordToken(p.cur, "REVISION") && p.peekNext().Type == '='
}

// parseDescribeListingRevision reads
//
//	{ DESC | DESCRIBE } LISTING <name> [ REVISION = { DRAFT | PUBLISHED } ]
//
// from REVISION on.
//
// https://docs.snowflake.com/en/sql-reference/sql/desc-listing
func (p *Parser) parseDescribeListingRevision(stmt *ast.DescribeStmt) error {
	p.advance() // consume REVISION
	p.advance() // consume '='
	for _, w := range []string{"DRAFT", "PUBLISHED"} {
		if p.acceptWords(w) {
			stmt.Revision = w
			return nil
		}
	}
	return p.syntaxErrorAtCur()
}

// parseShowAvailableListingsFilters reads the filters of
//
//	SHOW [ TERSE ] AVAILABLE LISTINGS [ LIMIT <rows> ]
//	  [ IS_IMPORTED = TRUE ] [ IS_ORGANIZATION = TRUE ] [ IS_SHARED_WITH_ME = TRUE ]
//
// each at most once, in the page's order.
//
// https://docs.snowflake.com/en/sql-reference/sql/show-available-listings
func (p *Parser) parseShowAvailableListingsFilters(stmt *ast.ShowStmt) error {
	for _, f := range []string{"IS_IMPORTED", "IS_ORGANIZATION", "IS_SHARED_WITH_ME"} {
		if !isWordToken(p.cur, f) || p.peekNext().Type != '=' {
			continue
		}
		p.advance() // consume the filter
		p.advance() // consume '='
		if err := p.expectWords("TRUE"); err != nil {
			return err
		}
		stmt.TrueFilters = append(stmt.TrueFilters, f)
	}
	return nil
}
