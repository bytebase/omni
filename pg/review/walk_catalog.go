package review

import (
	"context"
	"errors"
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
// or cure (a missing name: 3F000, 42P01, 42703, 42704, 42883; a coercion
// a cast could allow, castCurable), privileges (42501), and 0A000, which
// the catalog also raises for what it does not model.
var catalogCodes = map[string]review.FindingCode{
	"22005": review.UnsupportedForm, // error in assignment
	"22007": review.UnsupportedForm, // invalid datetime format
	"22P02": review.UnsupportedForm, // invalid text representation
	"2BP01": review.UnsupportedForm, // dependent objects still exist
	"42601": review.UnsupportedForm, // syntax error, in a routine's body
	"42702": review.UnsupportedForm, // ambiguous column
	"42710": review.DuplicateName,   // duplicate object
	"42723": review.DuplicateName,   // duplicate function
	"42725": review.UnsupportedForm, // ambiguous function
	"42804": review.UnsupportedForm, // datatype mismatch, of a RETURN
	"42809": review.WrongObjectKind, // wrong object type
	"42P02": review.UnsupportedForm, // undefined parameter
	"42P13": review.UnsupportedForm, // invalid function definition
	"42P17": review.UnsupportedForm, // invalid object definition
}

// catalogRefusal replays the change, in the order plan runs it, in a
// pg/catalog loaded from the synced schema, and returns a violation at the
// first statement the catalog refuses with one of catalogCodes, when
// catalogJudges it. It is asked only when the walk found nothing and knows
// the search path. The replay ends, accepting, at a schema that does not
// load in full, at any other refusal, at a rollback, and after a statement
// that may run user code (DML, CALL, DO), whose effects the catalog does
// not follow.
func catalogRefusal(ctx context.Context, stmts []statement, plan *run, target review.Target, path []string, user string) *walk.Violation {
	if target.Schema == nil || target.OtherSpecs || len(path) == 0 || user == "" && slices.Contains(path, "$user") {
		return nil
	}
	c := pgcatalog.New()
	if user != "" {
		c.SetSessionUser(user)
	}
	c.SetSearchPath(path)
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
	for i, st := range plan.steps {
		if st.again {
			return nil
		}
		node := stmts[st.stmt].node
		switch v := node.(type) {
		case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
			return nil
		case *ast.TransactionStmt:
			// The catalog does not undo what a rollback undoes.
			switch v.Kind {
			case ast.TRANS_STMT_ROLLBACK, ast.TRANS_STMT_ROLLBACK_TO, ast.TRANS_STMT_PREPARE, ast.TRANS_STMT_ROLLBACK_PREPARED:
				return nil
			}
		}
		err := processUtility(c, node)
		var refusal *pgcatalog.Error
		if errors.As(err, &refusal) {
			code, ok := catalogCodes[refusal.Code]
			if !ok || !catalogJudges(node) || castCurable(refusal) {
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
		case *ast.CallStmt, *ast.DoStmt:
			return nil
		}
	}
	return nil
}

// castCurable reports whether a cast the sync does not list could cure the
// refusal: a coercion the catalog finds no cast for. A RETURN that does
// not fit its routine's kind needs no cast.
func castCurable(err *pgcatalog.Error) bool {
	switch err.Code {
	case "42846":
		return true
	case "42804":
		return !strings.Contains(err.Message, "RETURN")
	case "42P13":
		return strings.Contains(err.Message, " instead of ") || strings.Contains(err.Message, "Actual return type")
	}
	return false
}

// catalogJudges reports whether the catalog's refusal of node is reported:
// a view's query, and what defines, changes, calls or drops a routine or a
// trigger, which pg/catalog's queries and programs suites pin. Any other
// statement the catalog follows, and its refusal ends the replay.
func catalogJudges(node ast.Node) bool {
	switch v := node.(type) {
	case *ast.ViewStmt, *ast.CreateFunctionStmt, *ast.AlterFunctionStmt, *ast.CallStmt, *ast.DoStmt,
		*ast.CreateTrigStmt, *ast.CreateEventTrigStmt, *ast.AlterEventTrigStmt:
		return true
	case *ast.DropStmt:
		switch ast.ObjectType(v.RemoveType) {
		case ast.OBJECT_FUNCTION, ast.OBJECT_PROCEDURE, ast.OBJECT_ROUTINE, ast.OBJECT_TRIGGER, ast.OBJECT_EVENT_TRIGGER:
			return true
		}
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
