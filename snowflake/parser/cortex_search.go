package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// CORTEX SEARCH SERVICE, read as an ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// Column lists (ON, TEXT INDEXES, VECTOR INDEXES, ATTRIBUTES) are clauses
// whose Names hold one single-part name per column; a parenthesized one
// (PRIMARY KEY ( ... ), SET ATTRIBUTES ( ... )) is a clause whose one
// nameless option lists the columns as bare names.

func init() {
	registerObject(&objectForm{
		words:     []string{"CORTEX", "SEARCH", "SERVICE"},
		create:    parseCreateCortexSearch,
		alter:     parseAlterCortexSearch,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-cortex-search
		// writes no IF EXISTS.
		drop:           true,
		dropNoIfExists: true,
		// CREATE's properties, then ALTER ... SET's. ALTER's
		// AUTO_SUSPEND = NULL has no shape here yet.
		props: objectProps{
			"WAREHOUSE":                      propName,
			"TARGET_LAG":                     propString,
			"EMBEDDING_MODEL":                propString,
			"REFRESH_MODE":                   propOneOf("FULL", "INCREMENTAL"),
			"INITIALIZE":                     propOneOf("ON_CREATE", "ON_SCHEDULE"),
			"FULL_INDEX_BUILD_INTERVAL_DAYS": propNumber,
			"REQUEST_LOGGING":                propBool,
			"AUTO_SUSPEND":                   propNumber,
			"COMMENT":                        propString,
		},
		alterProps: objectProps{
			"TARGET_LAG":                     propString,
			"WAREHOUSE":                      propName,
			"PRIMARY KEY":                    propNames,
			"FULL_INDEX_BUILD_INTERVAL_DAYS": propNumber,
			"REQUEST_LOGGING":                propBool,
			"AUTO_SUSPEND":                   propNumber,
			"COMMENT":                        propString,
		},
		unset: []string{}, // UNSET takes only PRIMARY KEY, ATTRIBUTES and TAG
	})
}

// aiCheckOptions holds opts to props, as objectForm.checkClauses holds a
// CREATE's run and a SET, for the ML/AI clauses that it does not see (a WITH
// run, a version's SET), and rejects a property set twice, as
// parseObjectProperties does. A nil props checks only the latter.
func aiCheckOptions(opts []*ast.CopyOption, props objectProps) error {
	for i, o := range opts {
		for _, prev := range opts[:i] {
			if prev.Name == o.Name {
				return &ParseError{
					Position: o.Loc.Start, End: o.Loc.End,
					Message: "syntax error: " + o.Name + " is set twice",
				}
			}
		}
		if props == nil {
			continue
		}
		if v, ok := props[o.Name]; !ok || !v.matches(o) {
			return &ParseError{
				Position: o.Loc.Start, End: o.Loc.End,
				Message: "syntax error at or near " + o.Name,
			}
		}
	}
	return nil
}

// parseAiParenNames reads ( <col_name> [ , ... ] ) into a clause of keywords,
// begun at start and already consumed, whose one nameless option lists the
// columns as bare names. cur is '('.
func (p *Parser) parseAiParenNames(s *ast.ObjectStmt, keywords string, start int) error {
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	cols := &ast.CopyOption{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	if err := p.parseObjectValueParen(cols); err != nil {
		return err
	}
	cols.Loc.End = p.prev.Loc.End
	if !propNames.matches(cols) {
		return &ParseError{Position: cols.Loc.Start, End: cols.Loc.End, Message: "syntax error: " + keywords + " lists column names"}
	}
	p.addObjectKeywordsClause(s, keywords, start).Options = []*ast.CopyOption{cols}
	return nil
}

// parseCreateCortexSearch reads what follows the name of
//
//	CREATE [ OR REPLACE ] CORTEX SEARCH SERVICE [ IF NOT EXISTS ] <name>
//	  ON <search_column>
//	  [ PRIMARY KEY ( <col_name> [ , ... ] ) ]
//	  ATTRIBUTES <col_name> [ , ... ]
//	  WAREHOUSE = <warehouse_name>
//	  TARGET_LAG = '<num> { seconds | minutes | hours | days }'
//	  [ EMBEDDING_MODEL = <embedding_model_name> ]
//	  [ REFRESH_MODE = { FULL | INCREMENTAL } ]
//	  [ INITIALIZE = { ON_CREATE | ON_SCHEDULE } ]
//	  [ FULL_INDEX_BUILD_INTERVAL_DAYS = <num> ]
//	  [ REQUEST_LOGGING = { TRUE | FALSE } ]
//	  [ AUTO_SUSPEND = <num_seconds> ]
//	  [ COMMENT = '<comment>' ]
//	AS <query>
//
//	CREATE [ OR REPLACE ] CORTEX SEARCH SERVICE <name>
//	  TEXT INDEXES <text_column_name> [ , ... ]
//	  VECTOR INDEXES <column_specification> [ , ... ]
//	  [ PRIMARY KEY ( <col_name> [ , ... ] ) ]
//	  ATTRIBUTES <col_name> [ , ... ]
//	  WAREHOUSE = <warehouse_name>
//	  TARGET_LAG = '<num> { seconds | minutes | hours | days }'
//	  [ ... the properties above, but EMBEDDING_MODEL ... ]
//	AS <query>
//
// where a column specification is
// <vector_column_name> [ ( { model | query_model } = '<model_name>' ) ].
//
// https://docs.snowflake.com/en/sql-reference/sql/create-cortex-search
func parseCreateCortexSearch(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.cur.Type == kwON:
		p.advance() // consume ON
		col, err := p.parseAiSingleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "ON", start).Names = []*ast.ObjectName{col}
	case p.acceptWords("TEXT", "INDEXES"):
		if s.IfNotExists {
			// The multi-index form writes no IF NOT EXISTS.
			return &ParseError{Position: start, End: p.prev.Loc.End, Message: "syntax error: IF NOT EXISTS is not documented with TEXT INDEXES"}
		}
		if _, err := p.parseCortexColumnList(s, "TEXT INDEXES", start); err != nil {
			return err
		}
		if err := p.parseCortexVectorIndexes(s); err != nil {
			return err
		}
	default:
		return p.syntaxErrorAtCur()
	}
	if start := p.cur.Loc.Start; p.acceptWords("PRIMARY", "KEY") {
		if err := p.parseAiParenNames(s, "PRIMARY KEY", start); err != nil {
			return err
		}
	}
	start = p.cur.Loc.Start
	if err := p.expectWords("ATTRIBUTES"); err != nil {
		return err
	}
	if _, err := p.parseCortexColumnList(s, "ATTRIBUTES", start); err != nil {
		return err
	}
	props, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return err
	}
	if err := p.requireObjectProperties(props, "WAREHOUSE", "TARGET_LAG"); err != nil {
		return err
	}
	return p.parseObjectAsQuery(s)
}

// parseCortexVectorIndexes reads
// VECTOR INDEXES <vector_column_name> [ ( { model | query_model } = '<model_name>' ) ] [ , ... ]
// into a VECTOR INDEXES clause: Names lists the columns, and Options holds,
// for each column given a model, an option named by the column whose Group is
// that one setting.
func (p *Parser) parseCortexVectorIndexes(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if err := p.expectWords("VECTOR", "INDEXES"); err != nil {
		return err
	}
	c := p.addObjectKeywordsClause(s, "VECTOR INDEXES", start)
	for {
		col, err := p.parseAiSingleName()
		if err != nil {
			return err
		}
		c.Names = append(c.Names, col)
		if p.cur.Type == '(' {
			opt := &ast.CopyOption{Name: col.Name.Normalize(), Loc: ast.Loc{Start: col.Loc.Start}}
			p.advance() // consume '('
			if !p.curIsWord("MODEL") && !p.curIsWord("QUERY_MODEL") || !p.startsObjectProperty() {
				return p.syntaxErrorAtCur()
			}
			setting, err := p.parseObjectProperty()
			if err != nil {
				return err
			}
			if setting.Lit == nil || setting.Lit.Kind != ast.LitString {
				return &ParseError{Position: setting.Loc.Start, End: setting.Loc.End, Message: "syntax error: a model name is a string"}
			}
			if _, err := p.expect(')'); err != nil {
				return err
			}
			opt.Group = []*ast.CopyOption{setting}
			opt.Loc.End = p.prev.Loc.End
			c.Options = append(c.Options, opt)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	c.Loc.End = p.prev.Loc.End
	return nil
}

// parseAlterCortexSearch reads what follows the name of
//
//	ALTER CORTEX SEARCH SERVICE [ IF EXISTS ] <name> { SUSPEND | RESUME } [ { INDEXING | SERVING } ]
//	ALTER CORTEX SEARCH SERVICE [ IF EXISTS ] <name> REFRESH
//	ALTER CORTEX SEARCH SERVICE [ IF EXISTS ] <name> SET
//	  [ TARGET_LAG = '<num> { seconds | minutes | hours | days }' ]
//	  [ WAREHOUSE = <warehouse_name> ]
//	  [ PRIMARY KEY = ( <col_name> [ , ... ] ) ]
//	  [ FULL_INDEX_BUILD_INTERVAL_DAYS = <num> ]
//	  [ REQUEST_LOGGING = { TRUE | FALSE } ]
//	  [ AUTO_SUSPEND = { <num_seconds> | NULL } ]
//	  [ COMMENT = '<string_literal>' ]
//	ALTER CORTEX SEARCH SERVICE [ IF EXISTS ] <name> UNSET [ PRIMARY KEY ]
//	ALTER CORTEX SEARCH SERVICE [ IF EXISTS ] <name> SET ATTRIBUTES ( <col_name> [ , ... ] )
//	ALTER CORTEX SEARCH SERVICE [ IF EXISTS ] <name> UNSET ATTRIBUTES
//	ALTER CORTEX SEARCH SERVICE <name> ADD SCORING PROFILE [ IF NOT EXISTS ] <profile_name> <scoring_profile>
//	ALTER CORTEX SEARCH SERVICE <name> DROP SCORING PROFILE [ IF EXISTS ] <profile_name>
//	ALTER CORTEX SEARCH SERVICE [ IF EXISTS ] <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER CORTEX SEARCH SERVICE [ IF EXISTS ] <name> UNSET TAG <tag_name> [ , ... ]
//
// The scoring profile is a JSON string. The page writes no UNSET of another
// property.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-cortex-search
func parseAlterCortexSearch(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if ok, err := p.parseObjectAlterCommon(s, actTags); ok || err != nil {
		return err
	}
	switch {
	case p.atWords("SUSPEND") || p.atWords("RESUME"):
		c := p.addObjectKeywordsClause(s, strings.ToUpper(p.advance().Str), start)
		c.Loc.End = p.prev.Loc.End
		if p.atWords("INDEXING") || p.atWords("SERVING") {
			c.Keywords += " " + strings.ToUpper(p.advance().Str)
			c.Loc.End = p.prev.Loc.End
		}
	case p.acceptWords("REFRESH"):
		p.addObjectKeywordsClause(s, "REFRESH", start)
	case p.acceptWords("SET", "ATTRIBUTES"):
		return p.parseAiParenNames(s, "SET ATTRIBUTES", start)
	case p.cur.Type == kwSET:
		return p.parseCortexSet(s)
	case p.acceptWords("UNSET", "PRIMARY", "KEY"):
		p.addObjectKeywordsClause(s, "UNSET PRIMARY KEY", start)
	case p.acceptWords("UNSET", "ATTRIBUTES"):
		p.addObjectKeywordsClause(s, "UNSET ATTRIBUTES", start)
	case !s.IfExists && p.acceptWords("ADD", "SCORING", "PROFILE"):
		c := p.addObjectKeywordsClause(s, "ADD SCORING PROFILE", start)
		if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
			if err := p.parseIfNotExists(new(bool)); err != nil {
				return err
			}
			c.Keywords = "ADD SCORING PROFILE IF NOT EXISTS"
		}
		name, err := p.parseAiSingleName()
		if err != nil {
			return err
		}
		profile, err := p.expect(tokString)
		if err != nil {
			return err
		}
		c.Names = []*ast.ObjectName{name}
		c.Options = []*ast.CopyOption{{Lit: &ast.Literal{Kind: ast.LitString, Value: profile.Str, Loc: profile.Loc}, Loc: profile.Loc}}
		c.Loc.End = p.prev.Loc.End
	case !s.IfExists && p.acceptWords("DROP", "SCORING", "PROFILE"):
		c := p.addObjectKeywordsClause(s, "DROP SCORING PROFILE", start)
		if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			c.Keywords = "DROP SCORING PROFILE IF EXISTS"
		}
		name, err := p.parseAiSingleName()
		if err != nil {
			return err
		}
		c.Names = []*ast.ObjectName{name}
		c.Loc.End = p.prev.Loc.End
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseCortexSet reads ALTER's SET <property> [ ... ], where a property is
// `KEY = value` or PRIMARY KEY = ( <col_name> [ , ... ] ), the latter an
// option named PRIMARY KEY whose Group holds the columns as bare names.
// cur is SET.
func (p *Parser) parseCortexSet(s *ast.ObjectStmt) error {
	c := &ast.ObjectClause{Keywords: "SET", Loc: ast.Loc{Start: p.cur.Loc.Start}}
	p.advance() // consume SET
	for {
		switch {
		case p.atWords("PRIMARY", "KEY"):
			key := p.advance()
			p.advance() // consume KEY
			if _, err := p.expect('='); err != nil {
				return err
			}
			opt := &ast.CopyOption{Name: "PRIMARY KEY", Loc: ast.Loc{Start: key.Loc.Start}}
			if p.cur.Type != '(' {
				return p.syntaxErrorAtCur()
			}
			if err := p.parseObjectValueParen(opt); err != nil {
				return err
			}
			for _, e := range opt.Group {
				if !e.Bare {
					return &ParseError{Position: e.Loc.Start, End: e.Loc.End, Message: "syntax error: PRIMARY KEY lists column names"}
				}
			}
			if len(opt.Group) == 0 {
				return &ParseError{Position: opt.Loc.Start, End: p.prev.Loc.End, Message: "syntax error: PRIMARY KEY lists column names"}
			}
			opt.Loc.End = p.prev.Loc.End
			c.Options = append(c.Options, opt)
		case p.startsObjectProperty():
			opt, err := p.parseObjectProperty()
			if err != nil {
				return err
			}
			c.Options = append(c.Options, opt)
		default:
			if len(c.Options) == 0 {
				return p.syntaxErrorAtCur()
			}
			if err := aiCheckOptions(c.Options, nil); err != nil {
				return err
			}
			c.Loc.End = p.prev.Loc.End
			s.Clauses = append(s.Clauses, c)
			return nil
		}
	}
}

// parseCortexColumnList reads <col_name> [ , ... ] into a clause of keywords
// begun at start, and appends and returns it.
func (p *Parser) parseCortexColumnList(s *ast.ObjectStmt, keywords string, start int) (*ast.ObjectClause, error) {
	c := p.addObjectKeywordsClause(s, keywords, start)
	for {
		col, err := p.parseAiSingleName()
		if err != nil {
			return nil, err
		}
		c.Names = append(c.Names, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	c.Loc.End = p.prev.Loc.End
	return c, nil
}

// parseAiSingleName reads a name the page writes as a single <name> (a
// column, profile, version or run).
func (p *Parser) parseAiSingleName() (*ast.ObjectName, error) {
	return p.parseObjectSimpleName()
}
