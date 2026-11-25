package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// semanticViewMaterializationActions maps the verb before MATERIALIZATION in
// ALTER SEMANTIC VIEW to its action.
var semanticViewMaterializationActions = map[string]ast.AlterSemanticViewAction{
	"ADD":     ast.AlterSemanticViewAddMaterialization,
	"DROP":    ast.AlterSemanticViewDropMaterialization,
	"SUSPEND": ast.AlterSemanticViewSuspendMaterialization,
	"RESUME":  ast.AlterSemanticViewResumeMaterialization,
	"REFRESH": ast.AlterSemanticViewRefreshMaterialization,
}

// parseAlterSemanticViewMaterialization reads the materialization actions of
// ALTER SEMANTIC VIEW <name>, and reports whether one started at cur:
//
//	ADD MATERIALIZATION <materialization_name>
//	  WAREHOUSE = <warehouse_name>
//	  [ REFRESH_MODE = { AUTO | FULL | INCREMENTAL } ]
//	  [ IMMUTABLE WHERE ( <immutable_condition> ) ]
//	  AS
//	    DIMENSIONS <dimension_name> [ , ... ]
//	    METRICS <metric_name> [ , ... ]
//	    [ WHERE ( <filter_condition> ) ]
//	{ DROP | SUSPEND | RESUME | REFRESH } MATERIALIZATION <materialization_name>
//
// The page writes these forms without IF EXISTS.
// https://docs.snowflake.com/en/sql-reference/sql/alter-semantic-view
func (p *Parser) parseAlterSemanticViewMaterialization(stmt *ast.AlterSemanticViewStmt) (bool, error) {
	if !p.curIs(tokIdent) && p.cur.Type < 700 || !p.nextWord("MATERIALIZATION") {
		return false, nil
	}
	action, ok := semanticViewMaterializationActions[strings.ToUpper(p.cur.Str)]
	if !ok || !p.curWord(strings.ToUpper(p.cur.Str)) {
		return false, nil
	}
	if stmt.IfExists {
		return true, p.syntaxErrorAtCur()
	}
	p.advance() // consume the verb
	p.advance() // consume MATERIALIZATION
	stmt.Action = action
	name, err := p.parseObjectSimpleName()
	if err != nil {
		return true, err
	}
	stmt.Materialization = name
	if action != ast.AlterSemanticViewAddMaterialization {
		return true, nil
	}

	if err := p.expectWords("WAREHOUSE"); err != nil {
		return true, err
	}
	if _, err := p.expect('='); err != nil {
		return true, err
	}
	if stmt.Warehouse, err = p.parseObjectSimpleName(); err != nil {
		return true, err
	}
	if p.acceptWords("REFRESH_MODE") {
		if _, err := p.expect('='); err != nil {
			return true, err
		}
		switch {
		case p.curWord("AUTO"), p.curWord("FULL"), p.curWord("INCREMENTAL"):
			stmt.RefreshMode = strings.ToUpper(p.advance().Str)
		default:
			return true, p.syntaxErrorAtCur()
		}
	}
	if p.acceptWords("IMMUTABLE", "WHERE") {
		if stmt.ImmutableWhere, err = p.parseParenCondition(); err != nil {
			return true, err
		}
	}
	if err := p.expectWords("AS", "DIMENSIONS"); err != nil {
		return true, err
	}
	if stmt.Dimensions, err = p.parseSemanticMemberList(); err != nil {
		return true, err
	}
	if err := p.expectWords("METRICS"); err != nil {
		return true, err
	}
	if stmt.Metrics, err = p.parseSemanticMemberList(); err != nil {
		return true, err
	}
	if p.curIs(kwWHERE) {
		p.advance() // consume WHERE
		if stmt.Where, err = p.parseParenCondition(); err != nil {
			return true, err
		}
	}
	return true, nil
}

// parseParenCondition reads ( <condition> ). cur is '('.
func (p *Parser) parseParenCondition() (ast.Node, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cond, nil
}

// parseSemanticMemberList reads a comma list of dimension or metric names.
// The page's example writes them as <logical_table>.<name>, so a name has one
// or two parts.
func (p *Parser) parseSemanticMemberList() ([]*ast.ObjectName, error) {
	var names []*ast.ObjectName
	for {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if !name.Database.IsEmpty() {
			return nil, &ParseError{
				Position: name.Loc.Start, End: name.Loc.End,
				Message: "syntax error: " + name.String() + " takes at most <logical_table>.<name>",
			}
		}
		names = append(names, name)
		if !p.curIs(',') {
			return names, nil
		}
		p.advance() // consume ','
	}
}
