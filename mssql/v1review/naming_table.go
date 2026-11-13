// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for SQL Server. The contract, and the split between
// these functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"regexp"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTable judges the table names top-level statements introduce
// against cfg: the table a plain CREATE TABLE creates, and the second
// argument of an EXEC sp_rename. Statements nested in a routine body or an
// IF or BEGIN block are not judged, nor are CREATE TABLE AS SELECT, CREATE
// TABLE AS CLONE OF, SELECT INTO, and RENAME OBJECT. A CREATE TABLE name is
// judged as written without its database or schema qualifier, a temporary
// table's # prefix included.
//
// Every unqualified sp_rename call is judged, whatever its @objtype, so a
// column, index, user data type, or database rename is judged too; a
// schema-qualified call such as sys.sp_rename is not. The judged name is
// the second argument read positionally even when arguments are named, so
// in sp_rename @newname = 'a', @objname = 'b' the name judged is b. Any
// literal there is judged, its text with one pair of enclosing single
// quotes removed; a variable is not judged.
//
// Lengths are in bytes. Every finding anchors the whole statement.
// NameMismatch and NameTooLong are the codes. Ported from Bytebase's V1
// rule NAMING_TABLE.
func NamingTable(stmts []ast.Node, cfg v1review.NamingConfig) ([]v1review.Finding, error) {
	pattern, err := cfg.Compile()
	if err != nil {
		return nil, err
	}
	if pattern == nil && cfg.MaxLength <= 0 {
		return nil, nil
	}
	naming := &tableNaming{pattern: pattern, max: cfg.MaxLength}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Name != nil {
				naming.check(n.Loc, n.Name.Object)
			}
		case *ast.ExecStmt:
			if name, ok := spRenameNewName(n); ok {
				naming.check(n.Loc, name)
			}
		}
	}
	return naming.findings, nil
}

// spRenameNewName returns the literal second argument of an unqualified
// sp_rename call.
func spRenameNewName(n *ast.ExecStmt) (string, bool) {
	if n.Name == nil || n.Name.Schema != "" || !strings.EqualFold(n.Name.Object, "sp_rename") {
		return "", false
	}
	if n.Args == nil || n.Args.Len() < 2 {
		return "", false
	}
	arg, ok := n.Args.Items[1].(*ast.ExecArg)
	if !ok || arg == nil {
		return "", false
	}
	lit, ok := arg.Value.(*ast.Literal)
	if !ok || lit == nil {
		return "", false
	}
	name := lit.Str
	if len(name) >= 2 && strings.HasPrefix(name, "'") && strings.HasSuffix(name, "'") {
		name = name[1 : len(name)-1]
	}
	return name, true
}

// tableNaming applies NamingTable's compiled config.
type tableNaming struct {
	pattern  *regexp.Regexp // nil requires nothing
	max      int            // 0 means no limit
	findings []v1review.Finding
}

func (c *tableNaming) check(loc ast.Loc, name string) {
	if name == "" {
		return
	}
	start, end := span(loc)
	if c.pattern != nil && !c.pattern.MatchString(name) {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Start: start, End: end})
	}
	if c.max > 0 && len(name) > c.max {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameTooLong, Name: name, Start: start, End: end})
	}
}

// span returns a Loc as a finding's byte range, 0,0 when the AST kept no
// start. A missing end collapses the range to its start.
func span(loc ast.Loc) (int, int) {
	if loc.Start < 0 {
		return 0, 0
	}
	return loc.Start, max(loc.End, loc.Start)
}
