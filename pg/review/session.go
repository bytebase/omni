package review

import (
	"context"
	"errors"
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/catalog"
	"github.com/bytebase/omni/review"
)

// session is one target's walk: a catalog holding the target's schema
// that the change's DDL is applied to statement by statement, so every
// rule that reads a target sees the state the statement would run in.
type session struct {
	target review.Target
	cat    *catalog.Catalog
	owners *owners
}

// newSession builds the catalog: session user, the synced search path,
// then the schema. An object the snapshot cannot install as defined is
// replaced by a stand-in, which is not an error; only a done context is.
func newSession(ctx context.Context, target review.Target) (*session, error) {
	cat := catalog.New()
	if target.SessionUser != "" {
		cat.SetSessionUser(target.SessionUser)
	}
	if target.Schema != nil {
		if path := parseSearchPath(target.Schema.GetSearchPath()); len(path) > 0 {
			cat.SetSearchPath(path)
		}
		if _, err := cat.LoadMetadata(ctx, target.Schema, catalog.LoadMetadataOptions{Full: true}); err != nil {
			return nil, err
		}
	}
	return &session{target: target, cat: cat, owners: newOwners(target.Schema)}, nil
}

// walk applies the statements in order. WalkThrough reports the first
// statement the catalog rejects, and the walk stops there: the state
// after a failed statement is not the state the next one would run in.
// With no schema there is nothing to reject against, so the walk skips
// failures and keeps going, and WalkThrough reports nothing. Ownership
// findings do not stop the walk: the catalog applies the statement
// regardless, so the state stays consistent.
func (s *session) walk(ctx context.Context, stmts []statement, on map[review.Rule]bool, r *reporter) error {
	report := on[review.WalkThrough] && s.target.Schema != nil
	for i := range stmts {
		st := &stmts[i]
		if err := ctx.Err(); err != nil {
			return err
		}
		if !isUtility(st.node) {
			continue
		}
		if report && s.target.SessionUser != "" {
			for _, msg := range s.ownershipProblems(st.node) {
				r.report(review.WalkThrough, st.index, rangeOf(st.loc), msg)
			}
		}
		err := s.cat.ProcessUtility(st.node)
		s.cat.DrainWarnings()
		if err == nil || isUnsupported(err) || s.target.Schema == nil {
			continue
		}
		if report {
			r.report(review.WalkThrough, st.index, rangeOf(st.loc), catalogMessage(err))
		}
		return nil
	}
	return nil
}

// isUtility reports whether the catalog simulates the statement. DML and
// cursor statements change no schema and are not simulated.
func isUtility(n ast.Node) bool {
	switch n.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt,
		*ast.DeclareCursorStmt, *ast.FetchStmt, *ast.ClosePortalStmt, *ast.ReturnStmt:
		return false
	}
	return true
}

// isUnsupported reports the catalog's "unsupported utility statement"
// error: a statement kind it does not model, which is not a finding.
func isUnsupported(err error) bool {
	var catErr *catalog.Error
	return errors.As(err, &catErr) && catErr.Code == catalog.CodeInvalidParameterValue &&
		strings.HasPrefix(catErr.Message, "unsupported utility statement")
}

// catalogMessage is the finding text of a catalog error: the message
// PostgreSQL would print, without the severity and SQLSTATE decoration.
func catalogMessage(err error) string {
	var catErr *catalog.Error
	if errors.As(err, &catErr) {
		return catErr.Message
	}
	return err.Error()
}
