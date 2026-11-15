package analysis

import (
	"errors"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
)

// SDLCode names the condition an SDL check found.
type SDLCode string

// Conditions CheckSDLStyle reports.
const (
	// SDLTableWithoutSchema: CREATE TABLE names no schema. Table.
	SDLTableWithoutSchema SDLCode = "TABLE_WITHOUT_SCHEMA"
	// SDLColumnForeignKeyWithoutSchema: a column-level REFERENCES names no
	// schema. Table, Column, RefTable.
	SDLColumnForeignKeyWithoutSchema SDLCode = "COLUMN_FOREIGN_KEY_WITHOUT_SCHEMA"
	// SDLColumnConstraint: a column carries PRIMARY KEY, UNIQUE, CHECK or
	// FOREIGN KEY. Table, Column, ConstraintTypes in declaration order.
	SDLColumnConstraint SDLCode = "COLUMN_CONSTRAINT"
	// SDLUnnamedConstraint: a table constraint has no name. Table,
	// ConstraintType ("PRIMARY KEY", "UNIQUE", "CHECK", "FOREIGN KEY", or
	// "CONSTRAINT" for any other kind).
	SDLUnnamedConstraint SDLCode = "UNNAMED_CONSTRAINT"
	// SDLForeignKeyWithoutSchema: a table-level FOREIGN KEY references a
	// table without a schema. Table, RefTable.
	SDLForeignKeyWithoutSchema SDLCode = "FOREIGN_KEY_WITHOUT_SCHEMA"
	// SDLUnnamedIndex: CREATE INDEX has no index name.
	SDLUnnamedIndex SDLCode = "UNNAMED_INDEX"
	// SDLIndexWithoutSchema: CREATE INDEX names its table without a schema.
	// Name (the index, empty when unnamed), Table.
	SDLIndexWithoutSchema SDLCode = "INDEX_WITHOUT_SCHEMA"
	// SDLViewWithoutSchema: CREATE VIEW names no schema. Name.
	SDLViewWithoutSchema SDLCode = "VIEW_WITHOUT_SCHEMA"
	// SDLSequenceWithoutSchema: CREATE SEQUENCE names no schema. Name.
	SDLSequenceWithoutSchema SDLCode = "SEQUENCE_WITHOUT_SCHEMA"
	// SDLFunctionWithoutSchema: CREATE FUNCTION or PROCEDURE has a one-part
	// name. Name.
	SDLFunctionWithoutSchema SDLCode = "FUNCTION_WITHOUT_SCHEMA"
	// SDLAlterSequenceWithoutSchema: ALTER SEQUENCE names no schema. Name.
	SDLAlterSequenceWithoutSchema SDLCode = "ALTER_SEQUENCE_WITHOUT_SCHEMA"
)

// SDLSyntaxError: the text has unbalanced parentheses or does not parse.
// Detail is the error text; Line is 0. Both checks report it alone.
const SDLSyntaxError SDLCode = "SYNTAX_ERROR"

// SDLFinding is one problem an SDL check found. Code names it, and each code's
// comment lists the fields it sets; the others are zero.
type SDLFinding struct {
	Code SDLCode
	// Line is the 1-based line in the checked text the finding is reported
	// at, 0 when it has none.
	Line int

	Schema string
	Table  string
	// Name is the object the finding is about when it is not the table:
	// an index, view, sequence, function or constraint.
	Name           string
	Column         string
	ConstraintType string
	// ConstraintTypes are the disallowed constraint kinds of an
	// SDLColumnConstraint column.
	ConstraintTypes []string

	RefSchema string
	RefTable  string
	RefColumn string
	// Columns are a foreign key's own columns (SDLForeignKeyTableNotFound) or
	// the referenced table's columns (SDLForeignKeyColumnNotFound).
	Columns []string

	// First is the earlier definition a duplicate collides with.
	First *SDLDefinition
	// PrimaryKeys are the PRIMARY KEY constraints of an
	// SDLMultiplePrimaryKeys table, in declaration order.
	PrimaryKeys []SDLDefinition

	Detail string
}

// SDLDefinition locates an object definition a finding cites.
type SDLDefinition struct {
	// Path is the file the definition is in; set only for cross-file
	// duplicates.
	Path string
	// Name is the constraint name, empty when unnamed.
	Name           string
	Table          string
	ConstraintType string
	Line           int
}

// CheckSDLStyle checks one declarative schema (SDL) file for the conventions
// Bytebase's SDL release check requires: every created table, index table,
// view, sequence and function, and every ALTER SEQUENCE, names its schema;
// columns carry no PRIMARY KEY, UNIQUE, CHECK or FOREIGN KEY constraint;
// TableElts constraints are named and their foreign keys name the referenced
// table's schema; indexes are named. A CreateStmt's separate Constraints list
// is not checked.
//
// Only top-level statements are judged. A table, index, view, sequence or
// function finding is reported at its statement's first line (where any
// leading comment starts); a constraint finding at the constraint's own line,
// and a column finding at its first disallowed constraint. A function's schema
// is the first part of a name of two or more parts, so in a.b.c it is "a".
//
// Unbalanced parentheses outside literals and comments, or a parse error,
// yield one SDLSyntaxError and nothing else.
func CheckSDLStyle(sql string) []SDLFinding {
	stmts, err := parseSDL(sql)
	if err != nil {
		return []SDLFinding{{Code: SDLSyntaxError, Detail: err.Error()}}
	}

	var findings []SDLFinding
	for _, stmt := range stmts {
		if stmt.AST == nil {
			continue
		}
		line := stmt.Start.Line

		switch n := stmt.AST.(type) {
		case *ast.CreateStmt:
			findings = checkSDLCreateStmt(sql, n, line, findings)
		case *ast.IndexStmt:
			findings = checkSDLIndexStmt(n, line, findings)
		case *ast.ViewStmt:
			if n.View != nil && n.View.Schemaname == "" {
				findings = append(findings, SDLFinding{Code: SDLViewWithoutSchema, Line: line, Name: n.View.Relname})
			}
		case *ast.CreateSeqStmt:
			if n.Sequence != nil && n.Sequence.Schemaname == "" {
				findings = append(findings, SDLFinding{Code: SDLSequenceWithoutSchema, Line: line, Name: n.Sequence.Relname})
			}
		case *ast.CreateFunctionStmt:
			parts := sdlListStrings(n.Funcname)
			if len(parts) < 2 || parts[0] == "" {
				name := ""
				if len(parts) > 0 {
					name = parts[len(parts)-1]
				}
				findings = append(findings, SDLFinding{Code: SDLFunctionWithoutSchema, Line: line, Name: name})
			}
		case *ast.AlterSeqStmt:
			if n.Sequence != nil && n.Sequence.Schemaname == "" {
				findings = append(findings, SDLFinding{Code: SDLAlterSequenceWithoutSchema, Line: line, Name: n.Sequence.Relname})
			}
		default:
		}
	}
	return findings
}

// parseSDL rejects unbalanced parentheses, which the parser may recover from
// silently, then parses.
func parseSDL(sql string) ([]omnipg.Statement, error) {
	if err := validateSDLParenBalance(sql); err != nil {
		return nil, err
	}
	return omnipg.Parse(sql)
}

func checkSDLCreateStmt(sql string, n *ast.CreateStmt, stmtLine int, findings []SDLFinding) []SDLFinding {
	tableName := ""
	schemaName := ""
	if n.Relation != nil {
		tableName = n.Relation.Relname
		schemaName = n.Relation.Schemaname
	}
	if schemaName == "" {
		findings = append(findings, SDLFinding{Code: SDLTableWithoutSchema, Line: stmtLine, Table: tableName})
	}

	var cols []*ast.ColumnDef
	var tableConstraints []*ast.Constraint
	if n.TableElts != nil {
		for _, item := range n.TableElts.Items {
			switch elem := item.(type) {
			case *ast.ColumnDef:
				cols = append(cols, elem)
			case *ast.Constraint:
				tableConstraints = append(tableConstraints, elem)
			default:
			}
		}
	}

	for _, col := range cols {
		var disallowed []string
		firstLine := 0
		if col.Constraints != nil {
			for _, item := range col.Constraints.Items {
				c, ok := item.(*ast.Constraint)
				if !ok {
					continue
				}
				typeName := sdlColumnConstraintType(c.Contype)
				if typeName == "" {
					continue
				}
				if disallowed == nil {
					firstLine = sdlNodeLine(sql, c.Loc.Start, stmtLine)
				}
				disallowed = append(disallowed, typeName)
				if c.Contype == ast.CONSTR_FOREIGN && c.Pktable != nil && c.Pktable.Schemaname == "" {
					findings = append(findings, SDLFinding{
						Code:     SDLColumnForeignKeyWithoutSchema,
						Line:     sdlNodeLine(sql, c.Loc.Start, stmtLine),
						Table:    tableName,
						Column:   col.Colname,
						RefTable: c.Pktable.Relname,
					})
				}
			}
		}
		if len(disallowed) > 0 {
			findings = append(findings, SDLFinding{
				Code:            SDLColumnConstraint,
				Line:            firstLine,
				Table:           tableName,
				Column:          col.Colname,
				ConstraintTypes: disallowed,
			})
		}
	}

	for _, c := range tableConstraints {
		if c.Conname == "" {
			findings = append(findings, SDLFinding{
				Code:           SDLUnnamedConstraint,
				Line:           sdlNodeLine(sql, c.Loc.Start, stmtLine),
				Table:          tableName,
				ConstraintType: sdlTableConstraintType(c.Contype),
			})
		}
		if c.Contype == ast.CONSTR_FOREIGN && c.Pktable != nil && c.Pktable.Schemaname == "" {
			findings = append(findings, SDLFinding{
				Code:     SDLForeignKeyWithoutSchema,
				Line:     sdlNodeLine(sql, c.Loc.Start, stmtLine),
				Table:    tableName,
				RefTable: c.Pktable.Relname,
			})
		}
	}
	return findings
}

func checkSDLIndexStmt(n *ast.IndexStmt, stmtLine int, findings []SDLFinding) []SDLFinding {
	if n.Idxname == "" {
		findings = append(findings, SDLFinding{Code: SDLUnnamedIndex, Line: stmtLine})
	}
	if n.Relation != nil && n.Relation.Schemaname == "" {
		findings = append(findings, SDLFinding{
			Code:  SDLIndexWithoutSchema,
			Line:  stmtLine,
			Name:  n.Idxname,
			Table: n.Relation.Relname,
		})
	}
	return findings
}

// sdlNodeLine is the line of a node starting at byte offset start, or
// defaultLine when the node has no location.
func sdlNodeLine(sql string, start int, defaultLine int) int {
	if start >= 0 {
		return sdlLineAt(sql, start)
	}
	return defaultLine
}

// sdlLineAt is the 1-based line of byte offset in text, clamped to the text.
// "\r", "\n" and "\r\n" each end a line.
func sdlLineAt(text string, offset int) int {
	offset = max(0, min(offset, len(text)))
	line := 1
	for i := 0; i < offset; i++ {
		switch text[i] {
		case '\r':
			line++
		case '\n':
			if i == 0 || text[i-1] != '\r' {
				line++
			}
		default:
		}
	}
	return line
}

// sdlColumnConstraintType names the column-level constraint kinds SDL
// disallows, and is empty for the allowed ones.
func sdlColumnConstraintType(contype ast.ConstrType) string {
	switch contype {
	case ast.CONSTR_PRIMARY:
		return "PRIMARY KEY"
	case ast.CONSTR_UNIQUE:
		return "UNIQUE"
	case ast.CONSTR_CHECK:
		return "CHECK"
	case ast.CONSTR_FOREIGN:
		return "FOREIGN KEY"
	default:
		return ""
	}
}

func sdlTableConstraintType(contype ast.ConstrType) string {
	if name := sdlColumnConstraintType(contype); name != "" {
		return name
	}
	return "CONSTRAINT"
}

func sdlListStrings(list *ast.List) []string {
	if list == nil {
		return nil
	}
	var result []string
	for _, item := range list.Items {
		if s, ok := item.(*ast.String); ok {
			result = append(result, s.Str)
		}
	}
	return result
}

// validateSDLParenBalance reports unmatched parentheses outside single-quoted
// strings, dollar-quoted strings and comments. Block comments do not nest, and
// a single-quoted string ignores backslash escapes.
func validateSDLParenBalance(sql string) error {
	depth := 0
	i := 0
	for i < len(sql) {
		ch := sql[i]
		switch {
		case ch == '-' && i+1 < len(sql) && sql[i+1] == '-':
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
		case ch == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i += 2
			for i+1 < len(sql) && (sql[i] != '*' || sql[i+1] != '/') {
				i++
			}
			if i+1 < len(sql) {
				i += 2
			}
		case ch == '\'':
			i++
			for i < len(sql) {
				if sql[i] == '\'' {
					if i+1 < len(sql) && sql[i+1] == '\'' {
						i += 2
						continue
					}
					break
				}
				i++
			}
			i++
		case ch == '$':
			j := i + 1
			for j < len(sql) && (sql[j] == '_' || (sql[j] >= 'a' && sql[j] <= 'z') || (sql[j] >= 'A' && sql[j] <= 'Z') || (sql[j] >= '0' && sql[j] <= '9')) {
				j++
			}
			if j < len(sql) && sql[j] == '$' {
				tag := sql[i : j+1]
				i = j + 1
				for i+len(tag) <= len(sql) {
					if sql[i] == '$' && sql[i:i+len(tag)] == tag {
						i += len(tag)
						break
					}
					i++
				}
			} else {
				i++
			}
		case ch == '(':
			depth++
			i++
		case ch == ')':
			depth--
			if depth < 0 {
				return errors.New("syntax error: unmatched closing parenthesis")
			}
			i++
		default:
			i++
		}
	}
	if depth != 0 {
		return errors.New("syntax error: unmatched opening parenthesis")
	}
	return nil
}
