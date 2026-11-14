package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// EXPERIMENT (Snowflake ML experiment tracking), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:     []string{"EXPERIMENT"},
		create:    parseCreateExperiment,
		alter:     parseAlterExperiment,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-experiment
		// writes DROP EXPERIMENT <name> with no IF EXISTS, which the registry's
		// DROP always takes, so DROP is not registered.
		// MODIFY RUN <run> SET STATUS = '<status>' is the one SET.
		props: objectProps{"STATUS": propString},
		unset: []string{},
	})
}

// parseCreateExperiment reads what follows the name of
//
//	CREATE [ OR REPLACE ] EXPERIMENT [ IF NOT EXISTS ] <name>
//
// which is nothing.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-experiment
func parseCreateExperiment(p *Parser, s *ast.ObjectStmt) error {
	return nil
}

// parseAlterExperiment reads what follows the name of
//
//	ALTER EXPERIMENT <name> ADD RUN <run_name>
//	ALTER EXPERIMENT <name> COMMIT RUN <run_name>
//	  [ WITH STATUS = { 'FINISHED' | 'FAILED' | 'CANCELLED' } ]
//	ALTER EXPERIMENT <name> MODIFY RUN <run_name> SET STATUS = { 'RUNNING' | 'PAUSED' }
//	ALTER EXPERIMENT <name> DROP RUN <run_name>
//
// The page writes no IF EXISTS. WITH and SET are clauses of their own.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-experiment
func parseAlterExperiment(p *Parser, s *ast.ObjectStmt) error {
	if s.IfExists {
		return &ParseError{Position: s.Name.Loc.Start, End: s.Name.Loc.End, Message: "syntax error: ALTER EXPERIMENT takes no IF EXISTS"}
	}
	start := p.cur.Loc.Start
	var verb string
	for _, v := range []string{"ADD", "COMMIT", "MODIFY", "DROP"} {
		if p.atWords(v, "RUN") {
			verb = v
		}
	}
	if verb == "" {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume the verb
	p.advance() // consume RUN
	if err := p.parseAiNameClause(s, verb+" RUN", start); err != nil {
		return err
	}
	start = p.cur.Loc.Start
	switch verb {
	case "COMMIT":
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
			return p.parseExperimentStatus(s, "WITH", start)
		}
	case "MODIFY":
		if _, err := p.expect(kwSET); err != nil {
			return err
		}
		return p.parseExperimentStatus(s, "SET", start)
	}
	return nil
}

// parseExperimentStatus reads STATUS = '<status>' into a clause of keywords,
// begun at start and already consumed.
func (p *Parser) parseExperimentStatus(s *ast.ObjectStmt, keywords string, start int) error {
	if !p.curIsWord("STATUS") || !p.startsObjectProperty() {
		return p.syntaxErrorAtCur()
	}
	status, err := p.parseObjectProperty()
	if err != nil {
		return err
	}
	if status.Lit == nil || status.Lit.Kind != ast.LitString {
		return &ParseError{Position: status.Loc.Start, End: status.Loc.End, Message: "syntax error: a run status is a string"}
	}
	p.addObjectKeywordsClause(s, keywords, start).Options = []*ast.CopyOption{status}
	return nil
}
