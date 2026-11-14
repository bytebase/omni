package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// TASK GRAPH, read as an ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// A task graph is named after TASK GRAPH, and a task may itself be named
// GRAPH (EXECUTE TASK graph RETRY LAST). So the form is not in the registry,
// whose match on TASK GRAPH would take every such task: the TASK cases of
// CREATE, ALTER, DROP and EXECUTE ask parseTaskOrGraph, which reads the
// statement as a task graph's and, when that does not read to the statement's
// end, as a task's.

// taskGraphForm is TASK GRAPH. It is not registered (see above).
var taskGraphForm = &objectForm{
	words:  []string{"TASK", "GRAPH"},
	kind:   "TASK GRAPH",
	create: parseCreateTaskGraph,
	alter:  parseAlterTaskGraph,
	// https://docs.snowflake.com/en/sql-reference/sql/create-task-graph
	orReplace: true,
	orAlter:   true,
	// https://docs.snowflake.com/en/sql-reference/sql/drop-task-graph
	drop:  true,
	props: taskGraphProps,
}

// taskGraphProps are the graph-wide properties CREATE TASK GRAPH documents,
// which ALTER TASK GRAPH ... SET and UNSET take.
var taskGraphProps = objectProps{
	"SCHEDULE":                        propString,
	"ROOTS":                           propNames,
	"FINALIZER":                       propQName,
	"OVERLAP_POLICY":                  propOneOf("NO_OVERLAP", "ALLOW_CHILD_OVERLAP", "ALLOW_ALL_OVERLAP"),
	"ERROR_INTEGRATION":               propName,
	"SUCCESS_INTEGRATION":             propName,
	"CONFIG":                          propString,
	"TASK_AUTO_RETRY_ATTEMPTS":        propNumber,
	"SUSPEND_TASK_AFTER_NUM_FAILURES": propNumber,
	"COMMENT":                         propString,
}

// taskGraphDefaults are the properties TASK_DEFAULTS ( ... ) supports, as
// CREATE TASK GRAPH lists them, with the values CREATE TASK gives each.
var taskGraphDefaults = objectProps{
	"WAREHOUSE": propName,
	"USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE": propString,
	"TARGET_COMPLETION_INTERVAL":               propString,
	"SERVERLESS_TASK_MIN_STATEMENT_SIZE":       propString,
	"SERVERLESS_TASK_MAX_STATEMENT_SIZE":       propString,
	"USER_TASK_TIMEOUT_MS":                     propNumber,
	"COMMENT":                                  propString,
	"LOG_LEVEL":                                propString,
	"TRACE_LEVEL":                              propString,
}

// startsTaskGraph reports whether cur, TASK, is followed by the word GRAPH.
func (p *Parser) startsTaskGraph() bool {
	return p.cur.Type == kwTASK && isWordToken(p.peekNext(), "GRAPH")
}

// parseTaskOrGraph reads a statement whose object starts TASK GRAPH: by graph
// when that reads to the statement's end, else by task. When neither does, it
// returns the graph's error.
func (p *Parser) parseTaskOrGraph(graph, task func() (ast.Node, error)) (ast.Node, error) {
	snap := p.snapshot()
	node, gerr := graph()
	if gerr == nil {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			return node, nil
		}
		gerr = p.syntaxErrorAtCur()
	}
	p.restore(snap)
	if node, err := task(); err == nil {
		return node, nil
	}
	return nil, gerr
}

// parseCreateTaskGraph reads what follows the name of
//
//	CREATE [ OR REPLACE ] TASK GRAPH [ IF NOT EXISTS ] <name>
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  [ WITH CONTACT ( <purpose> = <contact_name> [ , ... ] ) ]
//	  [ SCHEDULE = '...' ] [ ROOTS = ( <root_task_name> ) ] [ FINALIZER = <task_name> ]
//	  [ OVERLAP_POLICY = { NO_OVERLAP | ALLOW_CHILD_OVERLAP | ALLOW_ALL_OVERLAP } ]
//	  [ ERROR_INTEGRATION = <integration_name> ] [ SUCCESS_INTEGRATION = <integration_name> ]
//	  [ CONFIG = <configuration_string> ] [ TASK_AUTO_RETRY_ATTEMPTS = <num> ]
//	  [ SUSPEND_TASK_AFTER_NUM_FAILURES = <num> ] [ COMMENT = '<string>' ]
//	  [ TASK_DEFAULTS ( <default_property> = <value> [ , ... ] ) ]
//	  [ WHEN <condition> ]
//
//	CREATE OR ALTER TASK GRAPH <name> ...   (no TAG, CONTACT or IF NOT EXISTS)
//
// https://docs.snowflake.com/en/sql-reference/sql/create-task-graph
func parseCreateTaskGraph(p *Parser, s *ast.ObjectStmt) error {
	if s.OrAlter && s.IfNotExists {
		return &ParseError{Position: s.Name.Loc.Start, End: s.Name.Loc.End, Message: "syntax error: CREATE OR ALTER TASK GRAPH takes no IF NOT EXISTS"}
	}
	seen := map[string]bool{}
	tags := false
	for {
		start := p.cur.Loc.Start
		switch {
		case p.cur.Type == kwWITH && isWordToken(p.peekNext(), "CONTACT"):
			if s.OrAlter || seen["WITH CONTACT"] {
				return p.syntaxErrorAtCur()
			}
			seen["WITH CONTACT"] = true
			p.advance() // consume WITH
			p.advance() // consume CONTACT
			group, err := p.parseTaskGraphGroup(false)
			if err != nil {
				return err
			}
			if err := checkTaskGraphContacts(group); err != nil {
				return err
			}
			p.addObjectKeywordsClause(s, "WITH CONTACT", start).Options = []*ast.CopyOption{group}
		case isWordToken(p.cur, "TASK_DEFAULTS") && p.peekNext().Type == '(':
			if seen["TASK_DEFAULTS"] {
				return p.syntaxErrorAtCur()
			}
			seen["TASK_DEFAULTS"] = true
			p.advance() // consume TASK_DEFAULTS
			group, err := p.parseTaskGraphDefaults(false)
			if err != nil {
				return err
			}
			p.addObjectKeywordsClause(s, "TASK_DEFAULTS", start).Options = []*ast.CopyOption{group}
		case p.cur.Type == kwWHEN:
			// WHEN is the last clause.
			p.advance() // consume WHEN
			cond, err := p.parseExpr()
			if err != nil {
				return err
			}
			s.Clauses = append(s.Clauses, &ast.ObjectClause{Keywords: "WHEN", Query: cond, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
			return nil
		default:
			c, err := p.parseObjectProperties(s, "", false)
			if err != nil {
				return err
			}
			if c == nil {
				return nil
			}
			// The runs between the clauses above are one property list: each
			// property, and TAG, comes once.
			if c.Tags != nil {
				if tags || s.OrAlter {
					return &ParseError{Position: c.Loc.Start, End: c.Loc.End, Message: "syntax error: TAG is not allowed here"}
				}
				tags = true
			}
			for _, o := range c.Options {
				if seen[o.Name] {
					return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error: " + o.Name + " is set twice"}
				}
				seen[o.Name] = true
				if err := checkTaskGraphProperty(o); err != nil {
					return err
				}
			}
		}
	}
}

// parseAlterTaskGraph reads what follows the name of
//
//	ALTER TASK GRAPH [ IF EXISTS ] <name> SET <property> = <value> [ , ... ]
//	ALTER TASK GRAPH [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//	ALTER TASK GRAPH [ IF EXISTS ] <name> SET TASK_DEFAULTS ( <default_property> = <value> [ , ... ] )
//	ALTER TASK GRAPH [ IF EXISTS ] <name> UNSET TASK_DEFAULTS ( <default_property> [ , ... ] )
//	ALTER TASK GRAPH [ IF EXISTS ] <name> MODIFY WHEN <boolean_expr>
//	ALTER TASK GRAPH [ IF EXISTS ] <name> REMOVE WHEN
//	ALTER TASK GRAPH [ IF EXISTS ] <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER TASK GRAPH [ IF EXISTS ] <name> UNSET TAG <tag_name> [ , ... ]
//	ALTER TASK GRAPH [ IF EXISTS ] <name> SET CONTACT <purpose> = <contact_name> [ , ... ]
//	ALTER TASK GRAPH [ IF EXISTS ] <name> UNSET CONTACT <purpose> [ , ... ]
//	ALTER TASK GRAPH [ IF EXISTS ] <name> { COMMIT | REVERT | RESUME | SUSPEND }
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-task-graph
func parseAlterTaskGraph(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	for _, w := range []string{"COMMIT", "REVERT", "RESUME", "SUSPEND"} {
		if p.acceptWords(w) {
			p.addObjectKeywordsClause(s, w, start)
			return nil
		}
	}
	switch {
	case p.acceptWords("MODIFY", "WHEN"):
		cond, err := p.parseExpr()
		if err != nil {
			return err
		}
		s.Clauses = append(s.Clauses, &ast.ObjectClause{Keywords: "MODIFY WHEN", Query: cond, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
		return nil
	case p.acceptWords("REMOVE", "WHEN"):
		p.addObjectKeywordsClause(s, "REMOVE WHEN", start)
		return nil
	case p.cur.Type == kwSET && isWordToken(p.peekNext(), "TASK_DEFAULTS"):
		p.advance() // consume SET
		p.advance() // consume TASK_DEFAULTS
		group, err := p.parseTaskGraphDefaults(false)
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "SET TASK_DEFAULTS", start).Options = []*ast.CopyOption{group}
		return nil
	case p.cur.Type == kwUNSET && isWordToken(p.peekNext(), "TASK_DEFAULTS"):
		p.advance() // consume UNSET
		p.advance() // consume TASK_DEFAULTS
		group, err := p.parseTaskGraphDefaults(true)
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "UNSET TASK_DEFAULTS", start).Options = []*ast.CopyOption{group}
		return nil
	case p.cur.Type == kwSET && isWordToken(p.peekNext(), "CONTACT"):
		p.advance() // consume SET
		p.advance() // consume CONTACT
		var contacts []*ast.CopyOption
		for {
			e, err := p.parseTaskGraphEntry(false)
			if err != nil {
				return err
			}
			contacts = append(contacts, e)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		if err := checkTaskGraphContacts(&ast.CopyOption{Group: contacts}); err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "SET CONTACT", start).Options = contacts
		return nil
	case p.cur.Type == kwUNSET && isWordToken(p.peekNext(), "CONTACT"):
		p.advance() // consume UNSET
		p.advance() // consume CONTACT
		keys, err := p.parseObjectUnsetKeys()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "UNSET CONTACT", start).Keys = keys
		return nil
	case p.cur.Type == kwSET && p.peekNext().Type != kwTAG:
		// SET <property> = <value> [ , ... ]: the page separates them with
		// commas.
		p.advance() // consume SET
		c := &ast.ObjectClause{Keywords: "SET", Loc: ast.Loc{Start: start}}
		seen := map[string]bool{}
		for {
			if !p.startsObjectProperty() {
				return p.syntaxErrorAtCur()
			}
			o, err := p.parseObjectProperty()
			if err != nil {
				return err
			}
			if seen[o.Name] {
				return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error: " + o.Name + " is set twice"}
			}
			seen[o.Name] = true
			if err := checkTaskGraphProperty(o); err != nil {
				return err
			}
			c.Options = append(c.Options, o)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		c.Loc.End = p.prev.Loc.End
		s.Clauses = append(s.Clauses, c)
		return nil
	}
	ok, err := p.parseObjectAlterCommon(s, actUnset|actTags)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}

// checkTaskGraphProperty holds a property to the value shape the page gives
// beyond objectProps: ROOTS names exactly one task.
func checkTaskGraphProperty(o *ast.CopyOption) error {
	if o.Name == "ROOTS" && len(o.Group) != 1 {
		return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error: ROOTS names exactly one root task"}
	}
	return nil
}

// checkTaskGraphContacts holds WITH CONTACT and SET CONTACT to
// <purpose> = <contact_name> entries. Like a tag list, it may name a purpose
// more than once.
func checkTaskGraphContacts(group *ast.CopyOption) error {
	for _, e := range group.Group {
		if !propName.matches(e) {
			return &ParseError{Position: e.Loc.Start, End: e.Loc.End, Message: "syntax error: a contact is <purpose> = <contact_name>"}
		}
	}
	return nil
}

// parseTaskGraphDefaults reads TASK_DEFAULTS' ( <default_property> = <value>
// [ , ... ] ), or with names ( <default_property> [ , ... ] ), holding each
// property to the ones the page lists, once. cur is '('.
func (p *Parser) parseTaskGraphDefaults(names bool) (*ast.CopyOption, error) {
	group, err := p.parseTaskGraphGroup(names)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, e := range group.Group {
		v, ok := taskGraphDefaults[e.Name]
		if !ok || (!names && !v.matches(e)) {
			return nil, &ParseError{Position: e.Loc.Start, End: e.Loc.End, Message: "syntax error at or near " + e.Name}
		}
		if seen[e.Name] {
			return nil, &ParseError{Position: e.Loc.Start, End: e.Loc.End, Message: "syntax error: " + e.Name + " is set twice"}
		}
		seen[e.Name] = true
	}
	return group, nil
}

// parseTaskGraphGroup reads ( <entry> [ , ... ] ) into a nameless option
// whose Group holds the entries. cur is '('.
func (p *Parser) parseTaskGraphGroup(names bool) (*ast.CopyOption, error) {
	open, err := p.expect('(')
	if err != nil {
		return nil, err
	}
	group := &ast.CopyOption{Loc: ast.Loc{Start: open.Loc.Start}}
	for {
		e, err := p.parseTaskGraphEntry(names)
		if err != nil {
			return nil, err
		}
		group.Group = append(group.Group, e)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	group.Loc.End = p.prev.Loc.End
	return group, nil
}

// parseTaskGraphEntry reads one <name> = <value> entry, or with names a bare
// <name>. The name is an unquoted identifier or keyword.
func (p *Parser) parseTaskGraphEntry(names bool) (*ast.CopyOption, error) {
	if !isNameWordTok(p.cur) {
		return nil, p.syntaxErrorAtCur()
	}
	if names {
		key := p.advance()
		return &ast.CopyOption{Name: strings.ToUpper(key.Str), Bare: true, Loc: key.Loc}, nil
	}
	if p.peekNext().Type != '=' {
		p.advance()
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseObjectProperty()
}
