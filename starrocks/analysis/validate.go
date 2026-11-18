package analysis

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
	"github.com/bytebase/omni/starrocks/parser"
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
// server rather than prepare a read. Also refused: SET ... AS DEFAULT STORAGE VOLUME,
// which sets the cluster's default volume; SET DEFAULT ROLE ... TO user, which
// changes users' default roles; and SET LDAP_ADMIN_PASSWORD, which stores the
// cluster's LDAP bind password. SET ROLE only activates roles for the session.
func needsExecute(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.SetStmt:
		return n.Type != "DEFAULT STORAGE VOLUME" && !slices.ContainsFunc(n.Items, setsServerState)
	case *ast.UnsetStmt:
		return !strings.EqualFold(n.Scope, "GLOBAL")
	}
	return false
}

// setsServerState reports whether a SET item is one of the server-changing forms
// omni parses as an item: SET DEFAULT ROLE, and LDAP_ADMIN_PASSWORD, which StarRocks
// lexes as a non-reserved word and so reads as a variable name.
func setsServerState(item *ast.SetItem) bool {
	if item.UserVariable {
		return false
	}
	return strings.EqualFold(item.Scope, "GLOBAL") || item.Name == "default role" || strings.EqualFold(item.Name, "ldap_admin_password")
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
	case *ast.SelectStmt, *ast.SetOpStmt, *ast.ParenSelect:
		return isReadOnlyQuery(node)
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
	case *ast.DescStorageVolumeStmt:
		// DESC STORAGE VOLUME answers with the volume's settings, like DESC
		// over a table.
		return n.Name != ""
	case *ast.OptimizerTraceStmt:
		// TRACE only plans the query it wraps; it never runs it. A trace of a
		// query that exports through INTO OUTFILE is still not a read.
		return n.Query != nil && isReadOnlyQuery(n.Query)
	case *ast.ExplainStmt:
		// EXPLAIN requires an inner query; the inner query must itself be a
		// shape EXPLAIN is meaningful for (SELECT family or DML). This blocks
		// inputs like `EXPLAIN`, `EXPLAIN DROP TABLE t`, etc.
		if n.Query == nil || !isExplainableInner(n.Query) {
			return false
		}
		// EXPLAIN ANALYZE runs the statement rather than only planning it, so
		// over DML it writes. StarRocks aborts the INSERT transaction by
		// default, but the run still creates any partitions the rows need and
		// those outlive the abort. The plan-only modifiers never execute, so
		// they stay readable over DML.
		if n.Type == "ANALYZE" {
			return isReadOnlyQuery(n.Query)
		}
		return true
	case *ast.HelpStmt:
		return true
	default:
		return false
	}
}

// isExplainableInner reports whether `node` is a shape that EXPLAIN can
// legitimately wrap (SELECT family or DML). DDL inside EXPLAIN is rejected:
// StarRocks only supports EXPLAIN on query and DML statements. Queries are
// classified by isReadOnlyQuery, the one export gate.
func isExplainableInner(node ast.Node) bool {
	if isReadOnlyQuery(node) {
		return true
	}
	switch node.(type) {
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt,
		*ast.TruncateTableStmt:
		return true
	}
	return false
}

// isReadOnlyQuery reports whether `node` is a SELECT-family query that reads
// and never writes. INTO OUTFILE exports to S3 or HDFS, past masking and past
// the export gate, so a query carrying one is not a read.
//
// This is the single export gate for the SELECT family. Both the top level
// (isReadOnlyAST) and the EXPLAIN path (isExplainableInner) go through it, so
// the editor cannot accept by way of EXPLAIN what it refuses when asked
// directly.
func isReadOnlyQuery(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.SelectStmt:
		return n.Into == nil
	case *ast.SetOpStmt:
		// The parser attaches INTO to an arm rather than the set operation.
		return !setOpWritesOutfile(n)
	case *ast.ParenSelect:
		// (SELECT ...) / ((SELECT ...)) — classify by the wrapped query so the
		// gate still applies inside the parens. ParenSelect only ever wraps a
		// SELECT, SetOpStmt or another ParenSelect, so no DML hides here.
		return isReadOnlyQuery(n.Sel)
	}
	return false
}

// setOpWritesOutfile walks the arms: the parser attaches INTO to an arm rather
// than to the set operation, and arms nest.
func setOpWritesOutfile(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.SelectStmt:
		return n.Into != nil
	case *ast.SetOpStmt:
		return setOpWritesOutfile(n.Left) || setOpWritesOutfile(n.Right)
	case *ast.ParenSelect:
		// A parenthesized arm — SELECT 1 UNION (SELECT ... INTO OUTFILE ...) —
		// must not slip the export gate.
		return setOpWritesOutfile(n.Sel)
	default:
		return false
	}
}
