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
		// The CREATE page writes no IF NOT EXISTS.
		noIfNotExists: true,
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
	start := p.cur.Loc.Start
	if err := p.expectWords("PRIMARY", "KEY"); err != nil {
		return err
	}
	if err := p.parseAiParenNames(s, "PRIMARY KEY", start); err != nil {
		return err
	}

	// The properties, where COMMENT may follow WITH as TAG may.
	props := &ast.ObjectClause{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	for {
		if p.curIs(kwWITH) && p.nextWord("COMMENT") {
			p.advance() // consume WITH
		}
		switch {
		case p.startsWithTags() && (p.curIs(kwWITH) || p.nextIs('(')):
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
	if err := aiCheckOptions(props.Options, nil); err != nil {
		return err
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
	// The page writes the source as a single <source>.
	src, err := p.parseObjectSimpleName()
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
