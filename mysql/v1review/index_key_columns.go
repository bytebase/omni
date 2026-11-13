package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// keyColumnRule is what differs among the rules that judge the types of the
// columns keys are made of.
type keyColumnRule struct {
	code v1review.Code
	// typeOf writes a column definition's type as the rule judges and
	// prints it.
	typeOf func(*ast.DataType) string
	// attribute reports whether a column attribute puts its column in a
	// key the rule judges.
	attribute func(ast.ColumnConstraintType) bool
	// constraint reports whether the rule judges the columns of a
	// table-level definition of that type.
	constraint func(ast.ConstraintType) bool
	// createIndex is whether the rule judges the columns of a CREATE INDEX
	// that is neither FULLTEXT nor SPATIAL.
	createIndex bool
	// lower is whether a type found for a listed column is lower-cased
	// before it is judged and printed.
	lower bool
	// allowed reports whether a type passes.
	allowed func(string) bool
}

// judgeKeyColumns reports the key columns whose type rule does not allow,
// over a whole script: Statement is each finding's statement.
//
// A column definition with a data type, in a CREATE TABLE or in an ALTER
// TABLE's ADD COLUMN, MODIFY COLUMN or CHANGE COLUMN, gets one finding for
// each of its attributes the rule names, and its type is then remembered
// under the table's and the column's name exactly as written. CHANGE COLUMN
// first forgets the old name. A column inside the parenthesized ADD form is
// read; a key inside it is not.
//
// A table-level definition the rule names, in a CREATE TABLE or added by an
// ALTER TABLE, gets a finding for each plain column it lists whose type does
// not pass, as does a CREATE INDEX when the rule judges those. The type is
// the one remembered from the script; a column the script has not defined is
// looked up in schema, in the schema of the empty name, without regard to
// letter case, and skipped when schema does not hold it. When that lookup
// meets a schema with no snapshot the function returns
// v1review.ErrSchemaMissing.
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
		findings = append(findings, v1review.Finding{Code: rule.code, Name: column, Table: table, Type: columnType, Statement: statement, Start: start, End: end})
	}
	define := func(statement int, loc ast.Loc, table string, column *ast.ColumnDef) {
		if column == nil || column.TypeName == nil {
			return
		}
		columnType := rule.typeOf(column.TypeName)
		for _, attribute := range column.Constraints {
			if rule.attribute(attribute.Type) {
				report(statement, loc, table, column.Name, columnType)
			}
		}
		if defined[table] == nil {
			defined[table] = make(map[string]string)
		}
		defined[table][column.Name] = columnType
	}
	key := func(statement int, loc ast.Loc, table string, columns []string) error {
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
			}
			if rule.lower {
				columnType = strings.ToLower(columnType)
			}
			report(statement, loc, table, column, columnType)
		}
		return nil
	}
	for i, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, column := range n.Columns {
				if column != nil {
					define(i, column.Loc, n.Table.Name, column)
				}
			}
			for _, c := range n.Constraints {
				if c == nil || !rule.constraint(c.Type) {
					continue
				}
				if err := key(i, c.Loc, n.Table.Name, c.Columns); err != nil {
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
					if cmd.Column != nil {
						define(i, n.Loc, table, cmd.Column)
						continue
					}
					for _, column := range cmd.Columns {
						define(i, n.Loc, table, column)
					}
				case ast.ATModifyColumn:
					define(i, n.Loc, table, cmd.Column)
				case ast.ATChangeColumn:
					delete(defined[table], cmd.Name)
					define(i, n.Loc, table, cmd.Column)
				case ast.ATAddConstraint, ast.ATAddIndex:
					if cmd.Constraint == nil || !rule.constraint(cmd.Constraint.Type) {
						continue
					}
					if err := key(i, n.Loc, table, cmd.Constraint.Columns); err != nil {
						return nil, err
					}
				default:
				}
			}
		case *ast.CreateIndexStmt:
			if !rule.createIndex || n.Table == nil || n.Fulltext || n.Spatial {
				continue
			}
			if err := key(i, n.Loc, n.Table.Name, keyPartColumns(n.Columns)); err != nil {
				return nil, err
			}
		}
	}
	return findings, nil
}

// lowerTypeName is a column type as most V1 index rules wrote it: the type's
// name in lower case, without length, scale or attributes.
func lowerTypeName(t *ast.DataType) string {
	return strings.ToLower(t.Name)
}
