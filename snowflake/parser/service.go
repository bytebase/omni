package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// SERVICE and EXECUTE [ INFERENCE ] JOB SERVICE (Snowpark Container Services),
// read as an ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// EXECUTE JOB SERVICE creates and runs a job service, so it is read as an
// ObjectStmt too: Verb "EXECUTE", Kind "JOB SERVICE" or "INFERENCE JOB
// SERVICE", and no Name (the job's name, when given, is its NAME property).
//
// A stage a FROM names (`FROM @<stage> SPECIFICATION_FILE = ...`, an
// inference job's `FROM @<stage>`) is the clause's first option, nameless,
// whose Words keep the reference verbatim, as CREATE STREAMLIT's FROM keeps
// one (streamlit.go).

func init() {
	registerObject(&objectForm{
		words:  []string{"SERVICE"},
		create: parseCreateService,
		alter:  parseAlterService,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-service
		//
		//	DROP SERVICE [ IF EXISTS ] <name> [ FORCE ]
		drop: true, dropForce: true,
		props: serviceProps,
	})
}

// serviceProps are the properties CREATE SERVICE and ALTER SERVICE ... SET
// document; ALTER SERVICE ... UNSET takes the same names.
var serviceProps = objectProps{
	"AUTO_SUSPEND_SECS":            propNumber,
	"EXTERNAL_ACCESS_INTEGRATIONS": propNames,
	"AUTO_RESUME":                  propBool,
	"MIN_INSTANCES":                propNumber,
	"MIN_READY_INSTANCES":          propNumber,
	"MAX_INSTANCES":                propNumber,
	"LOG_LEVEL":                    propString,
	"QUERY_WAREHOUSE":              propName,
	"COMMENT":                      propString,
}

// executeJobServiceForm and executeInferenceJobServiceForm hold the
// properties EXECUTE JOB SERVICE and EXECUTE INFERENCE JOB SERVICE document.
// They are not registered: EXECUTE is dispatched by parseExecuteStmt.
var (
	executeJobServiceForm = &objectForm{props: objectProps{
		"NAME":                         propQName,
		"ASYNC":                        propBool,
		"REPLICAS":                     propNumber,
		"QUERY_WAREHOUSE":              propName,
		"COMMENT":                      propString,
		"EXTERNAL_ACCESS_INTEGRATIONS": propNames,
	}}
	executeInferenceJobServiceForm = &objectForm{props: objectProps{
		"MODEL":    propQName,
		"VERSION":  propName,
		"FUNCTION": propString,
		"NAME":     propQName,
		"ASYNC":    propBool,
		"REPLICAS": propNumber,
	}}
)

// parseCreateService reads what follows the name of
//
//	CREATE SERVICE [ IF NOT EXISTS ] <name>
//	  IN COMPUTE POOL <compute_pool_name>
//	  [ <property> ... ]   (as the guide's examples place them)
//	  { fromSpecification | fromSpecificationTemplate }
//	  [ AUTO_SUSPEND_SECS = <num> ]
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <EAI_name> [ , ... ] ) ]
//	  [ AUTO_RESUME = { TRUE | FALSE } ] [ MIN_INSTANCES = <num> ]
//	  [ MIN_READY_INSTANCES = <num> ] [ MAX_INSTANCES = <num> ]
//	  [ LOG_LEVEL = '<log_level>' ] [ QUERY_WAREHOUSE = <warehouse_name> ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  [ COMMENT = '{string_literal}']
//
// https://docs.snowflake.com/en/sql-reference/sql/create-service
func parseCreateService(p *Parser, s *ast.ObjectStmt) error {
	if err := p.parseServiceComputePool(s); err != nil {
		return err
	}
	return p.parseServicePropertiesAndSpecification(s, true)
}

// parseAlterService reads what follows the name of
//
//	ALTER SERVICE [ IF EXISTS ] <name> { SUSPEND | RESUME }
//	ALTER SERVICE [ IF EXISTS ] <name> { fromSpecification | fromSpecificationTemplate }
//	ALTER SERVICE [ IF EXISTS ] <name> SET <property> [ ... ]
//	ALTER SERVICE [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//	ALTER SERVICE [ IF EXISTS ] <name> SET [ TAG <tag_name> = '<tag_value>' [ , ... ] ]
//	ALTER SERVICE [ IF EXISTS ] <name> RESTORE VOLUME <volume_name>
//	  INSTANCES <comma_separated_instance_ids> FROM SNAPSHOT <snapshot_name>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-service
func parseAlterService(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("SUSPEND"):
		p.addObjectKeywordsClause(s, "SUSPEND", start)
		return nil
	case p.acceptWords("RESUME"):
		p.addObjectKeywordsClause(s, "RESUME", start)
		return nil
	case p.cur.Type == kwFROM:
		return p.parseServiceSpecification(s, true)
	case p.atWords("RESTORE", "VOLUME"):
		return p.parseServiceRestoreVolume(s)
	}
	ok, err := p.parseObjectAlterCommon(s, actSet|actSetTag|actUnset)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	return checkServiceSet(s)
}

// parseServiceRestoreVolume reads
//
//	RESTORE VOLUME <volume_name> INSTANCES <comma_separated_instance_ids>
//	  FROM SNAPSHOT <snapshot_name>
//
// into three clauses: RESTORE VOLUME naming the volume, INSTANCES holding each
// instance id (an integer, as the page's example writes INSTANCES 0,2) as a
// nameless option, and FROM SNAPSHOT naming the snapshot. cur is RESTORE.
func (p *Parser) parseServiceRestoreVolume(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	p.advance() // consume RESTORE
	p.advance() // consume VOLUME
	volume, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "RESTORE VOLUME", start).Names = []*ast.ObjectName{volume}

	start = p.cur.Loc.Start
	if err := p.expectWords("INSTANCES"); err != nil {
		return err
	}
	var ids []*ast.CopyOption
	for {
		tok, err := p.expect(tokInt)
		if err != nil {
			return err
		}
		ids = append(ids, &ast.CopyOption{
			Lit: &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc},
			Loc: tok.Loc,
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	p.addObjectKeywordsClause(s, "INSTANCES", start).Options = ids

	start = p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	if err := p.expectWords("SNAPSHOT"); err != nil {
		return err
	}
	snapshot, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "FROM SNAPSHOT", start).Names = []*ast.ObjectName{snapshot}
	return nil
}

// checkServiceSet holds an ALTER's SET clause to its page: tags are set by a
// SET TAG of their own, never among SET's properties.
func checkServiceSet(s *ast.ObjectStmt) error {
	for _, c := range s.Clauses {
		if c.Keywords == "SET" && len(c.Tags) > 0 {
			return &ParseError{Position: c.Loc.Start, End: c.Loc.End, Message: "syntax error at or near TAG"}
		}
	}
	return nil
}

// parseExecuteJobService reads
//
//	EXECUTE JOB SERVICE IN COMPUTE POOL <compute_pool_name>
//	  [ <property> ... ]   (as the usage notes and examples place them)
//	  { fromSpecification | fromSpecificationTemplate }
//	  [ NAME = [<db>.<schema>.]<name> ] [ ASYNC = { TRUE | FALSE } ]
//	  [ REPLICAS = <num> ] [ QUERY_WAREHOUSE = <warehouse_name> ]
//	  [ COMMENT = '<string_literal>' ]
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <EAI_name> [ , ... ] ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/execute-job-service
//
//	EXECUTE INFERENCE JOB SERVICE IN COMPUTE POOL <compute_pool_name>
//	  WITH SPECIFICATION '<specification_text>'
//	  FROM { ( <query> ) | @<stage> }
//	  MODEL = [<db>.<schema>.]<model_name> [ VERSION = <version_name> ]
//	  [ FUNCTION = '<function_name>' ] [ NAME = [<db>.<schema>.]<name> ]
//	  [ ASYNC = { TRUE | FALSE } ] [ REPLICAS = <num> ]
//
// https://docs.snowflake.com/en/sql-reference/sql/execute-inference-job-service
//
// cur is the word after EXECUTE; start is the EXECUTE token's Loc.
func (p *Parser) parseExecuteJobService(start ast.Loc) (ast.Node, error) {
	s := &ast.ObjectStmt{Verb: "EXECUTE", Kind: "JOB SERVICE", Loc: ast.Loc{Start: start.Start}}
	inference := p.acceptWords("INFERENCE")
	if inference {
		s.Kind = "INFERENCE JOB SERVICE"
	}
	if err := p.expectWords("JOB", "SERVICE"); err != nil {
		return nil, err
	}
	if err := p.parseServiceComputePool(s); err != nil {
		return nil, err
	}
	if !inference {
		if err := p.parseServicePropertiesAndSpecification(s, false); err != nil {
			return nil, err
		}
		if err := executeJobServiceForm.checkClauses(s); err != nil {
			return nil, err
		}
		s.Loc.End = p.prev.Loc.End
		return s, nil
	}

	// WITH SPECIFICATION '<specification_text>'
	withStart := p.cur.Loc.Start
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if err := p.expectWords("SPECIFICATION"); err != nil {
		return nil, err
	}
	spec, err := p.parseServiceSpecText()
	if err != nil {
		return nil, err
	}
	p.addObjectKeywordsClause(s, "WITH SPECIFICATION", withStart).Options = []*ast.CopyOption{spec}

	// FROM { ( <query> ) | @[<namespace>.]<stage_name>[/<path>] }
	fromStart := p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	if p.cur.Type == '@' {
		stage, err := p.parseServiceStage()
		if err != nil {
			return nil, err
		}
		p.addObjectKeywordsClause(s, "FROM", fromStart).Options = []*ast.CopyOption{stage}
	} else {
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var q ast.Node
		if p.cur.Type == kwWITH {
			q, err = p.parseWithQueryExpr()
		} else {
			q, err = p.parseQueryExpr()
		}
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		p.addObjectKeywordsClause(s, "FROM", fromStart).Query = q
	}

	props, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return nil, err
	}
	if err := p.requireObjectProperties(props, "MODEL"); err != nil {
		return nil, err
	}
	if err := executeInferenceJobServiceForm.checkClauses(s); err != nil {
		return nil, err
	}
	s.Loc.End = p.prev.Loc.End
	return s, nil
}

// startsExecuteJobService reports whether cur, the word after EXECUTE, begins
// JOB SERVICE or INFERENCE JOB SERVICE. It consumes nothing.
func (p *Parser) startsExecuteJobService() bool {
	return p.atWords("JOB", "SERVICE") || p.atWords("INFERENCE", "JOB", "SERVICE")
}

// parseServicePropertiesAndSpecification reads a service's specification with
// a property run on either side of it. The Syntax blocks write the properties
// after the specification; the pages' usage notes and examples also write them
// before it (CREATE SERVICE ... IN COMPUTE POOL p MIN_INSTANCES=1 FROM ...,
// https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-services#examples;
// EXECUTE JOB SERVICE: "specify compute pool, followed by other properties",
// https://docs.snowflake.com/en/sql-reference/sql/execute-job-service#usage-notes).
// Each property, and the tags, are set once across both runs.
func (p *Parser) parseServicePropertiesAndSpecification(s *ast.ObjectStmt, bareFile bool) error {
	before, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	if err := p.parseServiceSpecification(s, bareFile); err != nil {
		return err
	}
	after, err := p.parseObjectProperties(s, "", false)
	if err != nil || before == nil || after == nil {
		return err
	}
	for _, o := range after.Options {
		for _, prev := range before.Options {
			if prev.Name == o.Name {
				return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error: " + o.Name + " is set twice"}
			}
		}
	}
	if len(before.Tags) > 0 && len(after.Tags) > 0 {
		return &ParseError{Position: after.Loc.Start, End: after.Loc.End, Message: "syntax error: TAG is set twice"}
	}
	return nil
}

// parseServiceComputePool reads IN COMPUTE POOL <compute_pool_name> into a
// clause.
func (p *Parser) parseServiceComputePool(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if err := p.expectWords("IN", "COMPUTE", "POOL"); err != nil {
		return err
	}
	name, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "IN COMPUTE POOL", start).Names = []*ast.ObjectName{name}
	return nil
}

// parseServiceSpecification reads a service's specification, as the SERVICE
// and JOB SERVICE pages give it:
//
//	fromSpecification ::=
//	  { FROM SPECIFICATION_FILE = '<yaml_file_path>'
//	  | FROM @<stage> SPECIFICATION_FILE = '<yaml_file_path>'
//	  | FROM SPECIFICATION <specification_text> }
//	fromSpecificationTemplate ::=
//	  { FROM SPECIFICATION_TEMPLATE_FILE = '<yaml_file_path>'
//	  | FROM @<stage> SPECIFICATION_TEMPLATE_FILE = '<yaml_file_path>'
//	  | FROM SPECIFICATION_TEMPLATE <specification_text> }
//	  USING ( <key> => <value> [ , ... ] )
//
// A *_FILE form is a FROM clause holding the property (SPECIFICATION_FILE =
// '<path>'), after the stage when one is named; an inline form is a FROM
// SPECIFICATION or FROM SPECIFICATION_TEMPLATE clause holding the text as a
// nameless option. USING is a clause of its own. bareFile says whether the page
// documents a *_FILE without a stage (CREATE and ALTER SERVICE do, for a Native
// App's service; EXECUTE JOB SERVICE does not).
func (p *Parser) parseServiceSpecification(s *ast.ObjectStmt, bareFile bool) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	var opts []*ast.CopyOption
	keywords := "FROM"
	switch {
	case p.cur.Type == '@' || (bareFile && p.startsServiceSpecFile()):
		if p.cur.Type == '@' {
			stage, err := p.parseServiceStage()
			if err != nil {
				return err
			}
			opts = append(opts, stage)
			if !p.startsServiceSpecFile() {
				return p.syntaxErrorAtCur()
			}
		}
		key := p.advance()
		p.advance() // consume '='
		path, err := p.expect(tokString)
		if err != nil {
			return err
		}
		opts = append(opts, &ast.CopyOption{
			Name: strings.ToUpper(key.Str),
			Lit:  &ast.Literal{Kind: ast.LitString, Value: path.Str, Loc: path.Loc},
			Loc:  ast.Loc{Start: key.Loc.Start, End: path.Loc.End},
		})
	case p.curIsWord("SPECIFICATION") || p.curIsWord("SPECIFICATION_TEMPLATE"):
		keywords += " " + strings.ToUpper(p.advance().Str)
		spec, err := p.parseServiceSpecText()
		if err != nil {
			return err
		}
		opts = append(opts, spec)
	default:
		return p.syntaxErrorAtCur()
	}
	p.addObjectKeywordsClause(s, keywords, start).Options = opts
	if keywords == "FROM SPECIFICATION_TEMPLATE" || opts[len(opts)-1].Name == "SPECIFICATION_TEMPLATE_FILE" {
		return p.parseServiceUsing(s)
	}
	return nil
}

// startsServiceSpecFile reports whether cur begins SPECIFICATION_FILE = or
// SPECIFICATION_TEMPLATE_FILE =. It consumes nothing.
func (p *Parser) startsServiceSpecFile() bool {
	return (isWordToken(p.cur, "SPECIFICATION_FILE") || isWordToken(p.cur, "SPECIFICATION_TEMPLATE_FILE")) &&
		p.peekNext().Type == '='
}

// parseServiceStage reads @[<namespace>.]<stage_name>[/<path>] into a
// nameless option whose Words keep the reference verbatim. cur is '@'; a '@'
// with no stage after it is a syntax error.
func (p *Parser) parseServiceStage() (*ast.CopyOption, error) {
	at := p.cur
	ref, err := p.parseStageOperand()
	if err != nil {
		return nil, err
	}
	if ref.Loc.End == at.Loc.End {
		return nil, syntaxErrorAtTok(at)
	}
	return &ast.CopyOption{Words: ref.Raw, Loc: ref.Loc}, nil
}

// parseServiceSpecText reads '<specification_text>' (the text may be a $$
// string) into a nameless option.
func (p *Parser) parseServiceSpecText() (*ast.CopyOption, error) {
	text, err := p.expect(tokString)
	if err != nil {
		return nil, err
	}
	return &ast.CopyOption{
		Lit: &ast.Literal{Kind: ast.LitString, Value: text.Str, Loc: text.Loc},
		Loc: text.Loc,
	}, nil
}

// parseServiceUsing reads a specification template's arguments,
//
//	USING ( <key> => <value> [ , <key> => <value> [ , ... ] ] )
//
// into a USING clause, one option per argument named by its key exactly as
// written: a template variable's name is case-sensitive ("Guidelines for
// passing values for specification variables",
// https://docs.snowflake.com/en/developer-guide/snowpark-container-services/working-with-services).
// The page does not say a key may not repeat, so a repeated key is read as
// written. A value is a string (' or $$), a number, or an alphanumeric word;
// JSON is written as a string.
func (p *Parser) parseServiceUsing(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwUSING); err != nil {
		return err
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	c := p.addObjectKeywordsClause(s, "USING", start)
	for {
		key, err := p.parseIdent()
		if err != nil {
			return err
		}
		if _, err := p.expect(tokAssoc); err != nil {
			return err
		}
		if p.cur.Type == '(' {
			return p.syntaxErrorAtCur()
		}
		// A template variable's name is case-sensitive: keep it as written.
		arg := &ast.CopyOption{Name: p.srcSlice(key.Loc.Start, key.Loc.End), Loc: ast.Loc{Start: key.Loc.Start}}

		if err := p.parseObjectValue(arg); err != nil {
			return err
		}
		arg.Loc.End = p.prev.Loc.End
		c.Options = append(c.Options, arg)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	c.Loc.End = p.prev.Loc.End
	return nil
}
