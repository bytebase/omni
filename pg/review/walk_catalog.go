package review

import (
	"context"
	"errors"
	"math"
	"slices"
	"strings"

	"github.com/bytebase/omni/pg/ast"
	pgcatalog "github.com/bytebase/omni/pg/catalog"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// catalogCodes are the errors pg/catalog finds in a statement that hold
// whatever else the database holds, once the synced schema is loaded in
// full: each is pinned by pg/catalog's simulate suites on PostgreSQL 17 and
// 18. Left out are the errors an object the sync does not list can cause
// or cure (a missing schema, relation, type, role, language or event
// trigger: 3F000, 42P01, 42704; a coercion a cast could allow, or a call
// an aggregate or a domain could answer, curable), privileges (42501), and
// 0A000, which the catalog also raises for what it does not model.
var catalogCodes = map[string]review.FindingCode{
	"22005": review.UnsupportedForm, // error in assignment
	"22007": review.UnsupportedForm, // invalid datetime format
	"22P02": review.UnsupportedForm, // invalid text representation
	"2BP01": review.UnsupportedForm, // dependent objects still exist
	"42601": review.UnsupportedForm, // syntax error, in a routine's body
	"42611": review.UnsupportedForm, // invalid column definition
	"42702": review.UnsupportedForm, // ambiguous column
	"42703": review.ObjectNotFound,  // undefined column, of a synced table
	"42710": review.DuplicateName,   // duplicate object
	"42723": review.DuplicateName,   // duplicate function
	"42725": review.UnsupportedForm, // ambiguous function
	"42803": review.UnsupportedForm, // grouping error: an aggregate where none may be
	"42804": review.UnsupportedForm, // datatype mismatch, of a RETURN
	"42809": review.WrongObjectKind, // wrong object type
	"42883": review.ObjectNotFound,  // undefined function, by its signature
	"42P02": review.UnsupportedForm, // undefined parameter
	"42P13": review.UnsupportedForm, // invalid function definition
	"42P17": review.UnsupportedForm, // invalid object definition
	"42P20": review.UnsupportedForm, // windowing error: a window function where none may be
	"P0001": review.UnsupportedForm, // RAISE EXCEPTION, in a DO block the catalog runs
}

// catalogRefusal replays steps, the part of the change the walk reached,
// in a pg/catalog loaded from the synced schema, and returns a violation at
// the first statement the catalog refuses with one of catalogCodes, when
// catalogJudges it, in the session the walk knows for each step. It is
// asked only when the walk found nothing. The replay ends, accepting, at a schema that does not
// load in full, at any other refusal, and after a CALL or DO, whose effects
// the catalog does not follow, or a change to event triggers.
func catalogRefusal(ctx context.Context, stmts []statement, steps []step, sessions []catalogSession, target review.Target, snap *pgSnapshot) *walk.Violation {
	if len(steps) == 0 || target.Schema == nil || target.OtherSpecs {
		return nil
	}
	c := pgcatalog.New()
	if user := sessions[0].user; user != "" {
		c.SetSessionUser(user)
	}
	c.SetSearchPath(sessions[0].path)
	// Full loads procedures, which a call resolves among the functions.
	report, err := c.LoadMetadata(ctx, target.Schema, pgcatalog.LoadMetadataOptions{Full: true})
	if err != nil || len(report.Degraded) > 0 || len(report.Missing) > 0 {
		return nil
	}
	// PostgreSQL checks a routine's body at CREATE unless the change turns
	// check_function_bodies off; a snapshot's catalog checks none unless
	// told to.
	if _, err := c.Exec("SET check_function_bodies = on", nil); err != nil {
		return nil
	}
	for i, st := range steps {
		// The walk's session, which follows what the catalog does not, as
		// set_config does.
		c.SetSearchPath(sessions[i].path)
		if user := sessions[i].user; user != "" {
			c.SetRole(user)
		}
		node := stmts[st.stmt].node
		switch v := node.(type) {
		case *ast.SelectStmt:
			if v.IntoClause != nil {
				// SELECT INTO creates a table, which the replay does not.
				return nil
			}
			continue
		case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
			// DML changes no schema the walk followed it past.
			continue
		case *ast.VariableSetStmt:
			// The catalog keeps a SET LOCAL past the end of its transaction.
			if v.IsLocal {
				return nil
			}
		}
		err := processUtility(c, node)
		var refusal *pgcatalog.Error
		if errors.As(err, &refusal) {
			code, ok := catalogCodes[refusal.Code]
			if !ok || !catalogJudges(node) || curable(refusal, node, snap, sessions[i].path) || namesUntrustedRoutine(node) || setsOwnSettings(node) {
				return nil
			}
			message := "PostgreSQL refuses this statement: " + refusal.Message + " (SQLSTATE " + refusal.Code + ")."
			return &walk.Violation{
				Statement: i,
				Check:     walk.Check{Note: &walkNote{rng: rangeOf(stmts[st.stmt].loc), absent: message, present: message, fail: code}},
				Got:       walk.Present,
			}
		}
		if err != nil {
			return nil
		}
		switch node.(type) {
		case *ast.CallStmt, *ast.DoStmt, *ast.CreateEventTrigStmt, *ast.AlterEventTrigStmt:
			// User code ran, or an event trigger may fire on what follows.
			return nil
		}
	}
	return nil
}

// catalogSession is the search path and current user the walk knows a
// step to run in.
type catalogSession struct {
	path []string
	user string
}

// namesUntrustedRoutine reports whether node names a routine the sync may
// leave out (trustedRoutineName), which a call or a definition may resolve
// to without the catalog knowing it.
func namesUntrustedRoutine(node ast.Node) bool {
	untrusted := false
	check := func(names *ast.List) {
		if parts := nameParts(names); len(parts) > 0 && !trustedRoutineName(parts[len(parts)-1]) {
			untrusted = true
		}
	}
	ast.Inspect(node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.FuncCall:
			check(v.Funcname)
		case *ast.ObjectWithArgs:
			check(v.Objname)
		case *ast.CreateFunctionStmt:
			check(v.Funcname)
		case *ast.CreateTrigStmt:
			check(v.Funcname)
		case *ast.CreateEventTrigStmt:
			check(v.Funcname)
		}
		return !untrusted
	})
	return untrusted
}

// setsOwnSettings reports whether node defines or changes a routine with
// its own SET clause, under which PostgreSQL checks the body and the
// catalog does not.
func setsOwnSettings(node ast.Node) bool {
	var options *ast.List
	switch v := node.(type) {
	case *ast.CreateFunctionStmt:
		options = v.Options
	case *ast.AlterFunctionStmt:
		options = v.Actions
	default:
		return false
	}
	if options == nil {
		return false
	}
	for _, item := range options.Items {
		if de, ok := item.(*ast.DefElem); ok && strings.EqualFold(de.Defname, "set") {
			return true
		}
	}
	return false
}

// curable reports whether an object the sync does not list could cure the
// refusal of node, run on path. A DO block or a called procedure runs
// statements whose names the walk cannot check against the sync, under a
// search path of their own, so only RAISE EXCEPTION, a DO block's syntax,
// which PL/pgSQL compiles before it runs, and a CALL's failure to resolve
// its procedure hold. Otherwise: a user routine named as a
// built-in aggregate or window function, found first where the path puts
// pg_catalog after its schema; a cast, for a coercion the catalog finds
// none for (a RETURN that does not fit its routine's kind needs no cast);
// or, for a missing function, anything but a routine node names by its
// signature, which the sync lists unless it may have lost routines: an
// operator, a call in an expression (an aggregate, or a domain's
// function-style cast) or in a body, a CALL's argument an implicit cast
// may carry, and an aggregate that DROP ROUTINE or renaming a function or
// routine finds by its signature.
func curable(err *pgcatalog.Error, node ast.Node, snap *pgSnapshot, path []string) bool {
	switch v := node.(type) {
	case *ast.DoStmt:
		return err.Code != "P0001" && err.Code != "42601"
	case *ast.CallStmt:
		if err.Code != "P0001" && !resolvesCall(err, v) {
			return true
		}
	}
	switch err.Code {
	case "42803", "42P20":
		return slices.Index(path, "pg_catalog") > 0
	case "42846":
		return true
	case "42804":
		return !strings.Contains(err.Message, "RETURN")
	case "42P13":
		return strings.Contains(err.Message, " instead of ") || strings.Contains(err.Message, "Actual return type")
	case "42883":
		switch v := node.(type) {
		case *ast.RenameStmt:
			if v.RenameType != ast.OBJECT_PROCEDURE {
				return true
			}
		case *ast.DropStmt:
			if ast.ObjectType(v.RemoveType) == ast.OBJECT_ROUTINE {
				return true
			}
		case *ast.CallStmt:
			if v.Funccall != nil && !untypedArgs(v.Funccall.Args) {
				// An implicit cast the sync does not list may let an
				// argument of a known type reach the procedure.
				return true
			}
		}
		return snap == nil || snap.lostRoutines || !namesAny(err.Message, signatureNames(node))
	}
	return false
}

// resolvesCall reports whether err is call's failure to resolve the
// procedure it names, rather than a refusal of a statement the procedure's
// body runs: it must be checkCall's wording, with the types of call's own
// arguments, each a literal.
func resolvesCall(err *pgcatalog.Error, call *ast.CallStmt) bool {
	if call.Funccall == nil {
		return false
	}
	parts := nameParts(call.Funccall.Funcname)
	if len(parts) == 0 {
		return false
	}
	types, ok := literalArgTypes(call.Funccall.Args)
	if !ok {
		return false
	}
	signature := parts[len(parts)-1] + "(" + strings.Join(types, ", ") + ")"
	switch err.Code {
	case "42883":
		return err.Message == "procedure "+signature+" does not exist"
	case "42809":
		return err.Message == signature+" is not a procedure"
	}
	return false
}

// literalArgTypes returns the type PostgreSQL gives each argument, when
// each is a literal: a quoted string or NULL is unknown, an integer that
// fits integer is integer.
func literalArgTypes(args *ast.List) ([]string, bool) {
	var types []string
	if args == nil {
		return types, true
	}
	for _, item := range args.Items {
		if named, ok := item.(*ast.NamedArgExpr); ok {
			item = named.Arg
		}
		c, ok := item.(*ast.A_Const)
		if !ok {
			return nil, false
		}
		switch v := c.Val.(type) {
		case *ast.String:
			types = append(types, "unknown")
		case *ast.Integer:
			if v.Ival < math.MinInt32 || v.Ival > math.MaxInt32 {
				return nil, false
			}
			types = append(types, "integer")
		default:
			if !c.Isnull {
				return nil, false
			}
			types = append(types, "unknown")
		}
	}
	return types, true
}

// namesAny reports whether a refusal's message names one of the routines.
func namesAny(message string, names []string) bool {
	for _, name := range names {
		if strings.Contains(message, name+"(") || strings.Contains(message, `"`+name+`"`) {
			return true
		}
	}
	return false
}

// untypedArgs reports whether every argument is a quoted literal or NULL,
// which takes any parameter's type with no cast.
func untypedArgs(args *ast.List) bool {
	if args == nil {
		return true
	}
	for _, item := range args.Items {
		if named, ok := item.(*ast.NamedArgExpr); ok {
			item = named.Arg
		}
		c, ok := item.(*ast.A_Const)
		if !ok {
			return false
		}
		if _, str := c.Val.(*ast.String); !str && !c.Isnull {
			return false
		}
	}
	return true
}

// signatureNames lists the routines node names by signature, which the
// catalog looks up as such: the procedure a CALL names, a routine an ALTER,
// DROP or RENAME names, and a trigger's or event trigger's function.
func signatureNames(node ast.Node) []string {
	var names []string
	add := func(list *ast.List) {
		if parts := nameParts(list); len(parts) > 0 {
			names = append(names, parts[len(parts)-1])
		}
	}
	switch v := node.(type) {
	case *ast.CallStmt:
		if v.Funccall != nil {
			add(v.Funccall.Funcname)
		}
	case *ast.CreateTrigStmt:
		add(v.Funcname)
	case *ast.CreateEventTrigStmt:
		add(v.Funcname)
	case *ast.AlterFunctionStmt, *ast.DropStmt, *ast.RenameStmt:
		ast.Inspect(node, func(n ast.Node) bool {
			if owa, ok := n.(*ast.ObjectWithArgs); ok {
				add(owa.Objname)
			}
			return true
		})
	}
	return names
}

// catalogJudges reports whether the catalog's refusal of node is reported:
// a view's query; what defines, changes, calls or drops a routine or a
// trigger; dropping a table, view or type, or a column or its type, that a
// routine may depend on; which pg/catalog's queries and programs suites
// pin. Any other statement the catalog follows, and its refusal ends the
// replay.
func catalogJudges(node ast.Node) bool {
	switch v := node.(type) {
	case *ast.ViewStmt, *ast.CreateFunctionStmt, *ast.AlterFunctionStmt, *ast.CallStmt, *ast.DoStmt,
		*ast.CreateTrigStmt, *ast.CreateEventTrigStmt, *ast.AlterEventTrigStmt:
		return true
	case *ast.DropStmt:
		switch ast.ObjectType(v.RemoveType) {
		case ast.OBJECT_FUNCTION, ast.OBJECT_PROCEDURE, ast.OBJECT_ROUTINE, ast.OBJECT_TRIGGER, ast.OBJECT_EVENT_TRIGGER,
			ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_TYPE:
			return true
		}
	case *ast.AlterTableStmt:
		if v.Cmds == nil {
			return false
		}
		for _, item := range v.Cmds.Items {
			cmd, ok := item.(*ast.AlterTableCmd)
			if !ok || cmd.Subtype != int(ast.AT_DropColumn) && cmd.Subtype != int(ast.AT_AlterColumnType) {
				return false
			}
		}
		return true
	case *ast.RenameStmt:
		switch v.RenameType {
		case ast.OBJECT_FUNCTION, ast.OBJECT_PROCEDURE, ast.OBJECT_ROUTINE:
			return true
		}
	}
	return false
}

// processUtility runs one statement in the catalog. A panic in the catalog
// is no refusal.
func processUtility(c *pgcatalog.Catalog, node ast.Node) (err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("pg/catalog panicked")
		}
	}()
	return c.ProcessUtility(node)
}
