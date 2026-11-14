package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// ONLINE FEATURE TABLE (Snowflake Feature Store), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:     []string{"ONLINE", "FEATURE", "TABLE"},
		create:    parseCreateOnlineFeatureTable,
		alter:     parseAlterOnlineFeatureTable,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-online-feature-table
		drop: true,
		props: objectProps{
			"TARGET_LAG":       propString,
			"WAREHOUSE":        propName,
			"REFRESH_MODE":     propOneOf("AUTO", "FULL", "INCREMENTAL"),
			"TIMESTAMP_COLUMN": propName,
			"COMMENT":          propString,
		},
		alterProps: objectProps{
			"COMMENT":    propString,
			"TARGET_LAG": propString,
			"WAREHOUSE":  propName,
		},
		unset: []string{}, // UNSET takes only TAG
	})
}

// parseCreateOnlineFeatureTable reads what follows the name of
//
//	CREATE [ OR REPLACE ] ONLINE FEATURE TABLE <name>
//	  PRIMARY KEY ( <col_name> [ , ... ] )
//	  TARGET_LAG = '<num> { seconds | minutes | hours | days }'
//	  WAREHOUSE = <warehouse_name>
//	  [ REFRESH_MODE = { AUTO | FULL | INCREMENTAL } ]
//	  [ TIMESTAMP_COLUMN = <col_name> ]
//	  [ [ WITH ] COMMENT = '<string_literal>' ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  FROM <source>
//
// The page writes no IF NOT EXISTS.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-online-feature-table
func parseCreateOnlineFeatureTable(p *Parser, s *ast.ObjectStmt) error {
	if s.IfNotExists {
		return &ParseError{Position: s.Name.Loc.Start, End: s.Name.Loc.End, Message: "syntax error: CREATE ONLINE FEATURE TABLE takes no IF NOT EXISTS"}
	}
	start := p.cur.Loc.Start
	if err := p.expectWords("PRIMARY", "KEY"); err != nil {
		return err
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	pk, err := p.parseCortexColumnList(s, "PRIMARY KEY", start)
	if err != nil {
		return err
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	pk.Loc.End = p.prev.Loc.End

	// The properties, where COMMENT may follow WITH as TAG may.
	props := &ast.ObjectClause{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	for {
		if p.cur.Type == kwWITH && isWordToken(p.peekNext(), "COMMENT") {
			p.advance() // consume WITH
		}
		switch {
		case p.startsWithTags() && (p.cur.Type == kwWITH || p.peekNext().Type == '('):
			tags, err := p.parseWithTags()
			if err != nil {
				return err
			}
			props.Tags = append(props.Tags, tags...)
			continue
		case p.startsObjectProperty():
			opt, err := p.parseObjectProperty()
			if err != nil {
				return err
			}
			props.Options = append(props.Options, opt)
			continue
		}
		break
	}
	if len(props.Options) == 0 {
		return p.syntaxErrorAtCur()
	}
	props.Loc.End = p.prev.Loc.End
	s.Clauses = append(s.Clauses, props)
	if err := p.requireObjectProperties(props, "TARGET_LAG", "WAREHOUSE"); err != nil {
		return err
	}

	start = p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	src, err := p.parseObjectName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "FROM", start).Names = []*ast.ObjectName{src}
	return nil
}

// parseAlterOnlineFeatureTable reads what follows the name of
//
//	ALTER ONLINE FEATURE TABLE [ IF EXISTS ] <name> { SUSPEND | RESUME }
//	ALTER ONLINE FEATURE TABLE [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER ONLINE FEATURE TABLE [ IF EXISTS ] <name> REFRESH
//	ALTER ONLINE FEATURE TABLE [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//	ALTER ONLINE FEATURE TABLE [ IF EXISTS ] <name> SET
//	  { TARGET_LAG = '<num> { seconds | minutes | hours | days }' [ WAREHOUSE = <warehouse_name> ]
//	  | WAREHOUSE = <warehouse_name> }
//	ALTER ONLINE FEATURE TABLE [ IF EXISTS ] <name> <tagAction>
//
// where tagAction is SET TAG <tag_name> = '<tag_value>' [ , ... ] or
// UNSET TAG <tag_name> [ , ... ].
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-online-feature-table
func parseAlterOnlineFeatureTable(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("SUSPEND"):
		p.addObjectKeywordsClause(s, "SUSPEND", start)
		return nil
	case p.acceptWords("RESUME"):
		p.addObjectKeywordsClause(s, "RESUME", start)
		return nil
	case p.acceptWords("REFRESH"):
		p.addObjectKeywordsClause(s, "REFRESH", start)
		return nil
	}
	if ok, err := p.parseObjectAlterCommon(s, actSet|actTags|actRename); ok || err != nil {
		return err
	}
	return p.syntaxErrorAtCur()
}
