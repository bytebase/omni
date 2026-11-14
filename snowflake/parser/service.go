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
// The docs also give `FROM @<stage> SPECIFICATION_FILE = ...` (and, for an
// inference job, `FROM @<stage>`). ObjectClause has no field for a stage
// location yet, so those forms are not read.

func init() {
	registerObject(&objectForm{
		words:  []string{"SERVICE"},
		create: parseCreateService,
		alter:  parseAlterService,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-service
		// documents a trailing FORCE, which DropStmt cannot record yet.
		drop:  true,
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
		"NAME":                         propName,
		"ASYNC":                        propBool,
		"REPLICAS":                     propNumber,
		"QUERY_WAREHOUSE":              propName,
		"COMMENT":                      propString,
		"EXTERNAL_ACCESS_INTEGRATIONS": propNames,
	}}
	executeInferenceJobServiceForm = &objectForm{props: objectProps{
		"MODEL":    propName,
		"VERSION":  propName,
		"FUNCTION": propString,
		"NAME":     propName,
		"ASYNC":    propBool,
		"REPLICAS": propNumber,
	}}
)

// parseCreateService reads what follows the name of
//
//	CREATE SERVICE [ IF NOT EXISTS ] <name>
//	  IN COMPUTE POOL <compute_pool_name>
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
	if err := p.parseServiceSpecification(s, true); err != nil {
		return err
	}
	_, err := p.parseObjectProperties(s, "", false)
	return err
}

// parseAlterService reads what follows the name of
//
//	ALTER SERVICE [ IF EXISTS ] <name> { SUSPEND | RESUME }
//	ALTER SERVICE [ IF EXISTS ] <name> { fromSpecification | fromSpecificationTemplate }
//	ALTER SERVICE [ IF EXISTS ] <name> SET <property> [ ... ]
//	ALTER SERVICE [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//	ALTER SERVICE [ IF EXISTS ] <name> SET [ TAG <tag_name> = '<tag_value>' [ , ... ] ]
//
// RESTORE VOLUME <volume_name> INSTANCES <comma_separated_instance_ids>
// FROM SNAPSHOT <snapshot_name> is documented too, but its instance ids are
// numbers no ObjectClause field holds, so it is not read yet.
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
		if err := p.parseServiceSpecification(s, false); err != nil {
			return nil, err
		}
		if _, err := p.parseObjectProperties(s, "", false); err != nil {
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

	// FROM ( <query> )
	fromStart := p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
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
// '<path>'); an inline form is a FROM SPECIFICATION or FROM
// SPECIFICATION_TEMPLATE clause holding the text as a nameless option. USING
// is a clause of its own. bareFile says whether
// the page documents a *_FILE without a stage (CREATE and ALTER SERVICE do, for
// a Native App's service; EXECUTE JOB SERVICE does not). The @<stage> forms are
// not read yet (see the top of this file).
func (p *Parser) parseServiceSpecification(s *ast.ObjectStmt, bareFile bool) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	var spec *ast.CopyOption
	keywords := "FROM"
	switch {
	case bareFile && (p.curIsWord("SPECIFICATION_FILE") || p.curIsWord("SPECIFICATION_TEMPLATE_FILE")) &&
		p.peekNext().Type == '=':
		key := p.advance()
		p.advance() // consume '='
		path, err := p.expect(tokString)
		if err != nil {
			return err
		}
		spec = &ast.CopyOption{
			Name: strings.ToUpper(key.Str),
			Lit:  &ast.Literal{Kind: ast.LitString, Value: path.Str, Loc: path.Loc},
			Loc:  ast.Loc{Start: key.Loc.Start, End: path.Loc.End},
		}
	case p.curIsWord("SPECIFICATION") || p.curIsWord("SPECIFICATION_TEMPLATE"):
		keywords += " " + strings.ToUpper(p.advance().Str)
		var err error
		if spec, err = p.parseServiceSpecText(); err != nil {
			return err
		}
	default:
		return p.syntaxErrorAtCur()
	}
	p.addObjectKeywordsClause(s, keywords, start).Options = []*ast.CopyOption{spec}
	if keywords == "FROM SPECIFICATION_TEMPLATE" || spec.Name == "SPECIFICATION_TEMPLATE_FILE" {
		return p.parseServiceUsing(s)
	}
	return nil
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
// into a USING clause, one option per argument named by its key (uppercased,
// or a quoted key with its quotes). The page does not say a key may not repeat,
// so a repeated key is read as written. A value is a string (' or
// $$), a number, or an alphanumeric word; JSON is written as a string.
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
		arg := &ast.CopyOption{Name: strings.ToUpper(key.Name), Loc: ast.Loc{Start: key.Loc.Start}}
		if key.Quoted {
			arg.Name = p.srcSlice(key.Loc.Start, key.Loc.End)
		}

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
