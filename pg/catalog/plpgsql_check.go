package catalog

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	plast "github.com/bytebase/omni/pg/plpgsql/ast"
	plparser "github.com/bytebase/omni/pg/plpgsql/parser"
)

// plpgsqlRoutine is what PL/pgSQL compiles a body against: the routine's
// parameters and result, or a DO block's nothing.
type plpgsqlRoutine struct {
	name    string
	kind    byte // 'f' function, 'p' procedure
	retType uint32
	retSet  bool
	hasOut  bool
	params  []string // every parameter's name, "" for an unnamed one
}

// checkPLpgSQL compiles a PL/pgSQL body as the validator does when
// check_function_bodies is on, and as DO does before it runs the block: the
// block parses, every SQL expression and statement in it parses, each
// variable assigned or selected into is declared and not CONSTANT, each
// declared type exists, each condition is one PL/pgSQL knows, and RETURN
// fits the routine. It does not look at the tables a statement names,
// which PL/pgSQL leaves to run time.
//
// pg: src/pl/plpgsql/src/pl_comp.c — do_compile
// pg: src/pl/plpgsql/src/pl_gram.y — make_return_stmt, make_return_next_stmt, check_assignable, read_into_target
func (c *Catalog) checkPLpgSQL(body string, r *plpgsqlRoutine) error {
	block, err := plparser.Parse(body)
	if err != nil {
		var pe *plparser.ParseError
		if errors.As(err, &pe) {
			return &Error{Code: CodeSyntaxError, Message: pe.Message}
		}
		return err
	}
	root := &plScope{vars: map[string]bool{"found": false}, labels: map[string]bool{}}
	if r.name != "" {
		root.labels[strings.ToLower(r.name)] = true
	}
	for i, name := range r.params {
		root.vars["$"+strconv.Itoa(i+1)] = false
		if name != "" {
			root.vars[name] = false
		}
	}
	switch r.retType {
	case TRIGGEROID:
		for _, v := range []string{"new", "old", "tg_name", "tg_when", "tg_level", "tg_op", "tg_relid", "tg_relname", "tg_table_name", "tg_table_schema", "tg_nargs", "tg_argv"} {
			root.vars[v] = false
		}
	case EVENTTRIGGEROID:
		root.vars["tg_event"], root.vars["tg_tag"] = false, false
	}
	ch := &plChecker{c: c, r: r}
	return ch.block(block, root)
}

type plScope struct {
	outer  *plScope
	vars   map[string]bool // name → declared CONSTANT
	labels map[string]bool
}

func (s *plScope) child() *plScope {
	return &plScope{outer: s, vars: map[string]bool{}, labels: map[string]bool{}}
}

// lookup reports whether name is a variable or label in scope, and whether
// the variable is CONSTANT.
func (s *plScope) lookup(name string) (found, constant bool) {
	for cur := s; cur != nil; cur = cur.outer {
		if c, ok := cur.vars[name]; ok {
			return true, c
		}
		if cur.labels[name] {
			return true, false
		}
	}
	return false, false
}

type plChecker struct {
	c *Catalog
	r *plpgsqlRoutine
}

func (ch *plChecker) block(b *plast.PLBlock, outer *plScope) error {
	s := outer.child()
	if b.Label != "" {
		s.labels[plIdent(b.Label)] = true
	}
	for _, d := range b.Declarations {
		switch v := d.(type) {
		case *plast.PLDeclare:
			if err := ch.declType(v.TypeName, s); err != nil {
				return err
			}
			if v.Default != "" {
				if err := ch.expr(v.Default); err != nil {
					return err
				}
			}
			s.vars[plIdent(v.Name)] = v.Constant
		case *plast.PLCursorDecl:
			s.vars[plIdent(v.Name)] = false
		case *plast.PLAliasDecl:
			s.vars[plIdent(v.Name)] = false
		}
	}
	if err := ch.stmts(b.Body, s); err != nil {
		return err
	}
	for _, e := range b.Exceptions {
		w, ok := e.(*plast.PLExceptionWhen)
		if !ok {
			continue
		}
		for _, cond := range w.Conditions {
			name := strings.ToLower(strings.TrimSpace(cond))
			if name == "others" || strings.HasPrefix(name, "sqlstate") {
				continue
			}
			if !plpgsqlConditions[name] {
				return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("unrecognized exception condition %q", name)}
			}
		}
		hs := s.child()
		hs.vars["sqlstate"], hs.vars["sqlerrm"] = false, false
		if err := ch.stmts(w.Body, hs); err != nil {
			return err
		}
	}
	return nil
}

func (ch *plChecker) stmts(list []plast.Node, s *plScope) error {
	for _, n := range list {
		if err := ch.stmt(n, s); err != nil {
			return err
		}
	}
	return nil
}

func (ch *plChecker) loop(label string, body []plast.Node, s *plScope, vars ...string) error {
	ls := s.child()
	if label != "" {
		ls.labels[plIdent(label)] = true
	}
	for _, v := range vars {
		if v != "" {
			ls.vars[plIdent(v)] = false
		}
	}
	return ch.stmts(body, ls)
}

func (ch *plChecker) stmt(n plast.Node, s *plScope) error {
	switch v := n.(type) {
	case *plast.PLBlock:
		return ch.block(v, s)
	case *plast.PLAssign:
		if err := ch.target(v.Target, s); err != nil {
			return err
		}
		return ch.expr(v.Expr)
	case *plast.PLExecSQL:
		for _, t := range v.Into {
			if err := ch.target(t, s); err != nil {
				return err
			}
		}
		return ch.sql(v.SQLText)
	case *plast.PLPerform:
		return ch.expr(v.Expr)
	case *plast.PLIf:
		if err := ch.expr(v.Condition); err != nil {
			return err
		}
		if err := ch.stmts(v.ThenBody, s); err != nil {
			return err
		}
		for _, e := range v.ElsIfs {
			if ei, ok := e.(*plast.PLElsIf); ok {
				if err := ch.expr(ei.Condition); err != nil {
					return err
				}
				if err := ch.stmts(ei.Body, s); err != nil {
					return err
				}
			}
		}
		return ch.stmts(v.ElseBody, s)
	case *plast.PLCase:
		if v.HasTest {
			if err := ch.expr(v.TestExpr); err != nil {
				return err
			}
		}
		for _, w := range v.Whens {
			if cw, ok := w.(*plast.PLCaseWhen); ok {
				if err := ch.stmts(cw.Body, s); err != nil {
					return err
				}
			}
		}
		return ch.stmts(v.ElseBody, s)
	case *plast.PLLoop:
		return ch.loop(v.Label, v.Body, s)
	case *plast.PLWhile:
		if err := ch.expr(v.Condition); err != nil {
			return err
		}
		return ch.loop(v.Label, v.Body, s)
	case *plast.PLForI:
		for _, e := range []string{v.Lower, v.Upper, v.Step} {
			if e != "" {
				if err := ch.expr(e); err != nil {
					return err
				}
			}
		}
		return ch.loop(v.Label, v.Body, s, v.Var)
	case *plast.PLForS:
		return ch.loop(v.Label, v.Body, s, v.Var)
	case *plast.PLForC:
		return ch.loop(v.Label, v.Body, s, v.Var)
	case *plast.PLForDynS:
		return ch.loop(v.Label, v.Body, s, v.Var)
	case *plast.PLForEachA:
		return ch.loop(v.Label, v.Body, s, v.Var)
	case *plast.PLReturn:
		return ch.ret(v.Expr)
	case *plast.PLReturnNext:
		if !ch.r.retSet {
			return &Error{Code: CodeDatatypeMismatch, Message: "cannot use RETURN NEXT in a non-SETOF function"}
		}
		if v.Expr != "" {
			return ch.expr(v.Expr)
		}
	case *plast.PLReturnQuery:
		if !ch.r.retSet {
			return &Error{Code: CodeDatatypeMismatch, Message: "cannot use RETURN QUERY in a non-SETOF function"}
		}
	case *plast.PLRaise:
		if v.CondName != "" && !plpgsqlConditions[strings.ToLower(v.CondName)] {
			return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("unrecognized exception condition %q", strings.ToLower(v.CondName))}
		}
	}
	return nil
}

// ret checks a RETURN against the routine.
//
// pg: src/pl/plpgsql/src/pl_gram.y — make_return_stmt
func (ch *plChecker) ret(expr string) error {
	r := ch.r
	switch {
	case r.retSet:
		if expr != "" {
			return &Error{Code: CodeDatatypeMismatch, Message: "RETURN cannot have a parameter in function returning set"}
		}
	case r.retType == VOIDOID || r.retType == EVENTTRIGGEROID:
		if expr != "" && r.kind == 'p' {
			return &Error{Code: CodeSyntaxError, Message: "RETURN cannot have a parameter in a procedure"}
		}
		if expr != "" {
			return &Error{Code: CodeDatatypeMismatch, Message: "RETURN cannot have a parameter in function returning void"}
		}
	case r.hasOut:
		if expr != "" {
			return &Error{Code: CodeDatatypeMismatch, Message: "RETURN cannot have a parameter in function with OUT parameters"}
		}
	case r.retType == TRIGGEROID:
	case expr == "":
		return &Error{Code: CodeSyntaxError, Message: `missing expression at or near ";"`}
	}
	return ch.expr(expr)
}

// target checks a variable assigned to or selected into: its first name is
// a variable or label in scope, and a variable is not CONSTANT.
//
// pg: src/pl/plpgsql/src/pl_gram.y — check_assignable, read_into_target
func (ch *plChecker) target(text string, s *plScope) error {
	root := plIdent(plRootName(text))
	if root == "" {
		return nil
	}
	found, constant := s.lookup(root)
	if !found {
		return &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("%q is not a known variable", root)}
	}
	if constant && !strings.ContainsAny(text, ".[") {
		return &Error{Code: CodeErrorInAssignment, Message: fmt.Sprintf("variable %q is declared CONSTANT", root)}
	}
	return nil
}

// expr checks that an expression parses, as PL/pgSQL's raw parse of it in
// a SELECT does.
//
// pg: src/pl/plpgsql/src/pl_gram.y — check_sql_expr
func (ch *plChecker) expr(text string) error {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	return ch.sql("SELECT " + text)
}

func (ch *plChecker) sql(text string) error {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	if _, err := parser.Parse(text); err != nil {
		var pe *parser.ParseError
		if errors.As(err, &pe) {
			return &Error{Code: CodeSyntaxError, Message: pe.Message}
		}
		return err
	}
	return nil
}

// declType checks that a declared type exists: a type name, a table's
// row type (table%ROWTYPE), or a column's type (table.column%TYPE).
//
// pg: src/pl/plpgsql/src/pl_comp.c — plpgsql_parse_wordrowtype, plpgsql_parse_cwordtype
func (ch *plChecker) declType(text string, s *plScope) error {
	t := strings.TrimSpace(text)
	upper := strings.ToUpper(t)
	switch {
	case t == "":
		return nil
	case strings.HasSuffix(upper, "%ROWTYPE"):
		names := plNames(t[:len(t)-len("%ROWTYPE")])
		if len(names) == 0 || len(names) > 2 {
			return nil
		}
		schema, rel := "", names[len(names)-1]
		if len(names) == 2 {
			schema = names[0]
		}
		if _, _, err := ch.c.findRelation(schema, rel); err != nil {
			if _, _, typErr := ch.c.resolveTypeName(&nodes.TypeName{Names: plNameList(names)}); typErr == nil {
				return nil
			}
			return err
		}
		return nil
	case strings.HasSuffix(upper, "%TYPE"):
		names := plNames(t[:len(t)-len("%TYPE")])
		if len(names) != 2 {
			return nil
		}
		if found, _ := s.lookup(names[0]); found {
			return nil
		}
		_, rel, err := ch.c.findRelation("", names[0])
		if err != nil || rel == nil {
			return nil
		}
		if _, ok := rel.colByName[names[1]]; !ok {
			return &Error{Code: CodeUndefinedColumn, Message: fmt.Sprintf("column %q of relation %q does not exist", names[1], names[0])}
		}
		return nil
	}
	list, err := parser.Parse("SELECT NULL::" + t)
	if err != nil || list == nil || len(list.Items) != 1 {
		return nil
	}
	var tn *nodes.TypeName
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		if tc, ok := n.(*nodes.TypeCast); ok && tn == nil {
			tn = tc.TypeName
			return false
		}
		return true
	})
	if tn == nil {
		return nil
	}
	if _, name := qualifiedName(tn.Names); strings.EqualFold(name, "record") || strings.EqualFold(name, "refcursor") {
		return nil
	}
	_, _, err = ch.c.resolveTypeName(tn)
	return err
}

// plIdent folds an identifier as PL/pgSQL's scanner does: a quoted one
// keeps its case, any other is lower-cased.
func plIdent(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
	}
	return strings.ToLower(s)
}

// plRootName returns the first name of a target such as rec.field or
// arr[1].
func plRootName(text string) string {
	text = strings.TrimSpace(text)
	inQuote := false
	for i, r := range text {
		switch {
		case r == '"':
			inQuote = !inQuote
		case !inQuote && (r == '.' || r == '[' || r == ' '):
			return text[:i]
		}
	}
	return text
}

// plNames splits a dotted name, folding each part.
func plNames(text string) []string {
	var out []string
	cur := strings.Builder{}
	inQuote := false
	for _, r := range strings.TrimSpace(text) {
		switch {
		case r == '"':
			inQuote = !inQuote
			cur.WriteRune(r)
		case r == '.' && !inQuote:
			out = append(out, plIdent(cur.String()))
			cur.Reset()
		default:
			cur.WriteRune(r)
		}
	}
	return append(out, plIdent(cur.String()))
}

func plNameList(names []string) *nodes.List {
	l := &nodes.List{}
	for _, n := range names {
		l.Items = append(l.Items, &nodes.String{Str: n})
	}
	return l
}

// checkDo compiles a DO block before it would run. What the block does when
// it runs, the objects it creates or the errors it raises, is not simulated.
//
// pg: src/backend/commands/functioncmds.c — ExecuteDoStmt
// pg: src/pl/plpgsql/src/pl_handler.c — plpgsql_inline_handler
func (c *Catalog) checkDo(s *nodes.DoStmt) error {
	language, body := "plpgsql", ""
	if s.Args != nil {
		for _, item := range s.Args.Items {
			de, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			switch de.Defname {
			case "language":
				language = strings.ToLower(defElemString(de))
			case "as":
				body = defElemString(de)
			}
		}
	}
	switch language {
	case "plpgsql":
		return c.checkPLpgSQL(body, &plpgsqlRoutine{kind: 'f', retType: VOIDOID})
	case "sql", "c", "internal":
		return &Error{Code: CodeFeatureNotSupported, Message: fmt.Sprintf("language %q does not support inline code execution", language)}
	case "plperl", "plperlu", "pltcl", "pltclu", "plpythonu", "plpython2u", "plpython3u":
		return nil
	}
	return &Error{Code: CodeUndefinedObject, Message: fmt.Sprintf("language %q does not exist", language)}
}
