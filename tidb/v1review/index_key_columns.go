package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// keyColumnRule is what differs among the rules that judge the types of the
// columns keys are made of.
type keyColumnRule struct {
	code v1review.Code
	// typeOf writes a column definition's type as the rule judges and
	// prints it; nil is no type.
	typeOf func(*ast.DataType) string
	// attribute reports whether a column attribute of a column a CREATE
	// TABLE or an ADD COLUMN defines puts it in a key the rule judges.
	attribute func(ast.ColumnConstraintType) bool
	// changeAttribute is attribute for a column a MODIFY COLUMN or a CHANGE
	// COLUMN defines.
	changeAttribute func(ast.ColumnConstraintType) bool
	// constraint reports whether the rule judges the columns of a
	// table-level definition of that type.
	constraint func(ast.ConstraintType) bool
	// createIndex is whether the rule judges the columns of a CREATE INDEX.
	createIndex bool
	// lower is whether a type read from the schema is lower-cased before it
	// is judged and printed.
	lower bool
	// allowed reports whether a type passes.
	allowed func(string) bool
}

// judgeKeyColumns reports the key columns whose type rule does not allow,
// over a whole script: Statement is each finding's statement.
//
// A column definition, in a CREATE TABLE or in an ALTER TABLE's ADD COLUMN,
// MODIFY COLUMN or CHANGE COLUMN, gets one finding for each of its
// attributes the rule names, and its type is then remembered under the
// table's and the column's name exactly as written. CHANGE COLUMN first
// forgets the old name. A column inside the parenthesized ADD form is read;
// a key inside it is not.
//
// A table-level definition the rule names, in a CREATE TABLE or added by an
// ALTER TABLE, gets a finding for each plain column it lists whose type does
// not pass, as does a CREATE INDEX of any kind when the rule judges those.
// The type is the one remembered from the script; a column the script has
// not defined is looked up in schema, in the schema of the empty name,
// column names without regard to letter case, and skipped when schema does
// not hold it. When that lookup meets a schema with no snapshot the function
// returns v1review.ErrSchemaMissing.
//
// A CREATE TABLE finding anchors the column or key definition; every other
// finding anchors the whole statement.
func judgeKeyColumns(stmts []ast.Node, schema v1review.Schema, rule keyColumnRule) ([]v1review.Finding, error) {
	var findings []v1review.Finding
	defined := make(map[string]map[string]string)
	report := func(statement int, loc ast.Loc, table, column, columnType string) {
		if rule.allowed(columnType) {
			return
		}
		start, end := span(loc)
		findings = append(findings, v1review.Finding{Code: rule.code, Name: column, Table: table, Value: columnType, Statement: statement, Start: start, End: end})
	}
	define := func(statement int, loc ast.Loc, table string, column *ast.ColumnDef, attribute func(ast.ColumnConstraintType) bool) {
		columnType := rule.typeOf(column.TypeName)
		for _, c := range column.Constraints {
			if c != nil && attribute(c.Type) {
				report(statement, loc, table, column.Name, columnType)
			}
		}
		if defined[table] == nil {
			defined[table] = make(map[string]string)
		}
		defined[table][column.Name] = columnType
	}
	key := func(statement int, loc ast.Loc, table string, c *ast.Constraint) error {
		if !rule.constraint(c.Type) {
			return nil
		}
		return judgeListed(statement, loc, table, c.Columns, defined, schema, rule, report)
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, column := range n.Columns {
				if column != nil {
					define(i, column.Loc, n.Table.Name, column, rule.attribute)
				}
			}
			for _, c := range n.Constraints {
				if c == nil {
					continue
				}
				if err := key(i, c.Loc, n.Table.Name, c); err != nil {
					return nil, err
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			table := n.Table.Name
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATAddColumn:
					for _, column := range addedColumns(cmd) {
						define(i, n.Loc, table, column, rule.attribute)
					}
				case ast.ATAddConstraint, ast.ATAddIndex:
					if cmd.Constraint == nil {
						continue
					}
					if err := key(i, n.Loc, table, cmd.Constraint); err != nil {
						return nil, err
					}
				case ast.ATChangeColumn, ast.ATModifyColumn:
					if cmd.Column == nil {
						continue
					}
					if cmd.Type == ast.ATChangeColumn && cmd.Name != "" {
						delete(defined[table], cmd.Name)
					}
					define(i, n.Loc, table, cmd.Column, rule.changeAttribute)
				default:
				}
			}
		case *ast.CreateIndexStmt:
			if !rule.createIndex || n.Table == nil {
				continue
			}
			if err := judgeListed(i, n.Loc, n.Table.Name, keyPartColumns(n.Columns), defined, schema, rule, report); err != nil {
				return nil, err
			}
		}
	}
	return findings, nil
}

// judgeListed reports the listed columns of one key whose type does not
// pass, as judgeKeyColumns describes.
func judgeListed(statement int, loc ast.Loc, table string, columns []string, defined map[string]map[string]string, schema v1review.Schema, rule keyColumnRule, report func(int, ast.Loc, string, string, string)) error {
	for _, column := range columns {
		columnType, ok := defined[table][column]
		if !ok {
			if schema.Metadata == nil {
				return v1review.ErrSchemaMissing
			}
			found := schema.TableColumn("", table, column, false)
			if found == nil {
				continue
			}
			columnType = found.GetType()
			if rule.lower {
				columnType = strings.ToLower(columnType)
			}
		}
		report(statement, loc, table, column, columnType)
	}
	return nil
}

// lowerTypeName is a column type as most V1 index rules wrote it: the type's
// name in lower case, without length, scale or attributes, and empty for no
// type.
func lowerTypeName(t *ast.DataType) string {
	if t == nil {
		return ""
	}
	return strings.ToLower(t.Name)
}
