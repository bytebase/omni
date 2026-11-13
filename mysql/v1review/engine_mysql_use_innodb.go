package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// EngineMySQLUseInnoDB judges the storage engine top-level statements
// choose, which must be InnoDB in any letter case, and reports at most one
// finding a statement. Statements nested in a routine body are not judged.
//
// A CREATE TABLE is judged by its first ENGINE option, when that has a
// value. An ALTER TABLE is judged by its ENGINE table options, an empty
// value included. A SET is judged by each assignment to a variable named
// default_storage_engine, after any scope word, whose value is a string or
// a bare word; a variable written with @@ is not judged.
//
// Name is the engine as written. A CREATE TABLE's finding has Search set to
// it: the V1 advisor put the finding where the statement text first holds
// that text, which may be earlier than the option. Every range is the whole
// statement's. Ported from Bytebase's V1 rule ENGINE_MYSQL_USE_INNODB.
func EngineMySQLUseInnoDB(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			for _, opt := range n.Options {
				if !strings.EqualFold(opt.Name, "ENGINE") {
					continue
				}
				if opt.Value != "" && !isInnoDB(opt.Value) {
					findings = append(findings, notInnoDB(n.Loc, opt.Value, opt.Value))
				}
				break
			}
		case *ast.AlterTableStmt:
			for _, cmd := range n.Commands {
				if cmd.Type == ast.ATTableOption && cmd.Option != nil &&
					strings.EqualFold(cmd.Option.Name, "ENGINE") && !isInnoDB(cmd.Option.Value) {
					findings = append(findings, notInnoDB(n.Loc, cmd.Option.Value, ""))
					break
				}
			}
		case *ast.SetStmt:
			if engine, ok := setsOtherEngine(n); ok {
				findings = append(findings, notInnoDB(n.Loc, engine, ""))
			}
		}
	}
	return findings
}

// setsOtherEngine returns the first engine other than InnoDB a SET assigns
// to default_storage_engine.
func setsOtherEngine(n *ast.SetStmt) (string, bool) {
	for _, assign := range n.Assignments {
		if assign.Column == nil || strings.ToLower(assign.Column.Column) != "default_storage_engine" {
			continue
		}
		switch value := assign.Value.(type) {
		case *ast.StringLit:
			if !isInnoDB(value.Value) {
				return value.Value, true
			}
		case *ast.ColumnRef:
			if !isInnoDB(value.Column) {
				return value.Column, true
			}
		}
	}
	return "", false
}

func isInnoDB(engine string) bool {
	return strings.ToLower(engine) == "innodb"
}

func notInnoDB(loc ast.Loc, engine, search string) v1review.Finding {
	start, end := span(loc)
	return v1review.Finding{Code: v1review.EngineNotInnoDB, Name: engine, Search: search, Start: start, End: end}
}
