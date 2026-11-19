package analysis

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/doris/ast"
	"github.com/bytebase/omni/doris/parser"
)

// ValidateQuery reports whether statement may run in the SQL editor's
// read-only mode, and whether it returns rows rather than only setting session
// state: (readOnly, returnsData, err).
//
// Decision is AST-based: each parsed top-level statement must be a SELECT,
// SHOW, DESCRIBE, EXPLAIN, HELP, or a session SET. Relying on leading-keyword classification
// alone would mis-accept CTE-prefixed DML such as `WITH x AS (...) UPDATE ...`,
// which Classify would tag as SELECT because `WITH` is its first token.
//
// A query that carries INTO OUTFILE is refused: it writes its result set to
// S3 or HDFS, so the data leaves the database past masking and past the
// export gate. isReadOnlyAST reads the clause off the tree, so the words in a
// comment or a quoted identifier are not mistaken for one (BOT-92).
//
// Syntax errors are surfaced up; that lets a caller reject malformed
// read-only SQL before execution.
//
// returnsData is false once a statement needs executing rather than reading,
// which is how a driver decides between a query and an ExecContext call.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	parsed, err := parseStatements(statement)
	if err != nil {
		return false, false, err
	}
	readOnly = true
	var hasExecute bool
	for _, node := range parsed {
		switch {
		case isReadOnlyAST(node):
		case readsRows(node):
			// A statement the editor does not admit as read-only may still
			// return rows; only its read-only verdict changes.
			readOnly = false
		case needsExecute(node):
			hasExecute = true
		default:
			return false, false, nil
		}
	}
	return readOnly, !hasExecute, nil
}

// readsRows reports whether a statement returns a result set without being a
// read-only query: ADMIN SHOW inspects the cluster and answers with rows.
func readsRows(node ast.Node) bool {
	admin, ok := node.(*ast.AdminStmt)
	return ok && strings.EqualFold(admin.Verb, "SHOW")
}

// parseStatements parses every non-empty statement of text, which a trailing
// semicolon and comments do not make a statement of their own.
func parseStatements(text string) ([]ast.Node, error) {
	if err := firstLexError(text); err != nil {
		return nil, err
	}
	var nodes []ast.Node
	for _, segment := range parser.Split(text) {
		if segment.Empty() || strings.TrimSpace(segment.Text) == "" {
			continue
		}
		file, err := parser.Parse(segment.Text)
		if err != nil {
			return nil, offsetParseError(err, segment.ByteStart)
		}
		if file != nil && len(file.Stmts) > 0 {
			nodes = append(nodes, file.Stmts[0])
		}
	}
	return nodes, nil
}

// offsetParseError moves the positions parse reported for the text at base
// into the coordinates of the script it was cut from.
func offsetParseError(err error, base int) error {
	if base == 0 {
		return err
	}
	errs := parser.AllErrors(err)
	if len(errs) == 0 {
		return err
	}
	for i := range errs {
		if errs[i].Position >= 0 {
			errs[i].Position += base
		}
		if errs[i].End >= 0 {
			errs[i].End += base
		}
	}
	return parser.ParseErrors(errs)
}

// needsExecute reports whether a statement is allowed in the editor but is not a
// read. A session variable assignment returns no rows, so the driver has to run it
// with ExecContext rather than read a result set; ValidateQuery signals that through
// its second return value. MySQL and Postgres gate SET the same way.
//
// Only variable assignment qualifies. SET PASSWORD and SET DEFAULT STORAGE VAULT are
// separate nodes and stay refused, as they are for Postgres: they administer the
// server rather than prepare a read. SET PROPERTY, which sets a user's properties,
// and SET LDAP_ADMIN_PASSWORD, which stores the cluster's LDAP bind password, are
// refused too.
func needsExecute(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.SetStmt:
		return n.Type != "PROPERTY" && !slices.ContainsFunc(n.Items, setsServerState)
	case *ast.UnsetStmt:
		return !strings.EqualFold(n.Scope, "GLOBAL")
	}
	return false
}

// setsServerState reports whether a SET item changes state outside the session,
// which the read-only mode does not allow: a global variable, a user's
// properties, or the cluster's LDAP bind password.
func setsServerState(item *ast.SetItem) bool {
	if item.UserVariable {
		return false
	}
	return strings.EqualFold(item.Scope, "GLOBAL") || strings.EqualFold(item.Name, "ldap_admin_password")
}

// isReadOnlyAST returns true when the given top-level AST node represents a
// read-only Doris statement (SELECT family, SHOW, DESCRIBE, EXPLAIN, HELP).
//
// Omni's parser currently has stub-shaped acceptance for some shapes — bare
// `SHOW`, `DESCRIBE`, `EXPLAIN` all produce a corresponding AST node with
// empty content rather than a parse error. We reject those here so the
// previous ANTLR helper's behaviour (rejecting incomplete forms) is
// preserved; the database would only reject them at execution time
// otherwise.
//
// Nil nodes are conservatively rejected — they indicate a parse path that
// produced no concrete statement, which shouldn't happen for valid read-only
// SQL after parseStatements succeeds.
func isReadOnlyAST(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.SelectStmt:
		// INTO OUTFILE exports to S3 or HDFS, past masking and the export gate.
		return n.Into == nil
	case *ast.ValuesStmt:
		// A VALUES row-constructor query writes nothing itself; the export
		// gate applies to it as to any query.
		return n.Into == nil
	case *ast.SetOpStmt:
		return n.Into == nil
	case *ast.GroupedQuery:
		// A repeated trailing clause group or a parenthesized WITH group —
		// (SELECT 1 LIMIT 1) LIMIT 2 — wraps a query; classify by the query.
		// The export gate applies to the group as well, since the clause
		// attaches to whichever node is the statement's outermost query.
		return n.Into == nil && isReadOnlyAST(n.Query)
	case *ast.ShowStmt:
		// Reject bare `SHOW` (Type is empty when the parser took the stub path
		// without seeing a recognised variant keyword).
		return n.Type != ""
	case *ast.ShowRoutineLoadStmt, *ast.ShowRoutineLoadTaskStmt,
		*ast.ShowJobStmt, *ast.ShowJobTaskStmt,
		*ast.ShowConstraintsStmt, *ast.ShowAnalyzeStmt, *ast.ShowStatsStmt:
		return true
	case *ast.DescribeStmt:
		// Reject bare `DESCRIBE` / `DESC` (Target nil means no table named).
		return n.Target != nil
	case *ast.ExplainStmt:
		// EXPLAIN requires an inner query; the inner query must itself be a
		// shape EXPLAIN is meaningful for (SELECT/DML/Show). This blocks
		// inputs like `EXPLAIN`, `EXPLAIN DROP TABLE t`, etc.
		return n.Query != nil && isExplainableInner(n.Query)
	case *ast.HelpStmt:
		return true
	default:
		return false
	}
}

// isExplainableInner reports whether `node` is a shape that EXPLAIN can
// legitimately wrap (SELECT family or DML). DDL inside EXPLAIN is rejected:
// Doris only supports EXPLAIN on query and DML statements.
//
// An export stays refused under EXPLAIN. The engine does not treat
// `EXPLAIN SELECT ... INTO OUTFILE ...` as a plain plan: it fails with the
// same "Please specify BROKER properties or check your local file path" as
// the bare export, where `EXPLAIN SELECT 1` returns a plan (container-
// verified). So the clause is not inert here, and the gate holds.
func isExplainableInner(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.SelectStmt:
		return n.Into == nil
	case *ast.ValuesStmt:
		// EXPLAIN VALUES(1) plans a query; the export gate applies as it does
		// to a SELECT.
		return n.Into == nil
	case *ast.SetOpStmt:
		return n.Into == nil
	case *ast.GroupedQuery:
		// EXPLAIN (SELECT 1 LIMIT 1) LIMIT 2 — engine-verified accept;
		// classify by the wrapped query.
		return n.Into == nil && isExplainableInner(n.Query)
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt,
		*ast.MergeStmt, *ast.TruncateTableStmt,
		*ast.DescribeStmt:
		// EXPLAIN a, EXPLAIN a ALL, EXPLAIN FUNCTION a() and the partition
		// form parse as an ExplainStmt wrapping a DescribeStmt; they are
		// read-only descriptions and stay explainable.
		return true
	}
	return false
}
