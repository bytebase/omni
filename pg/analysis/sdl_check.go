package analysis

import (
	"errors"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
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
// Detail is the error text; Start is the unmatched parenthesis or the
// parser's error position, -1 when the parser gave none. Both checks report
// it alone.
const SDLSyntaxError SDLCode = "SYNTAX_ERROR"

// SDLFinding is one problem an SDL check found. Code names it, and each code's
// comment lists the fields it sets; the others are zero.
type SDLFinding struct {
	Code SDLCode
	// Start is the byte offset, in the text the finding belongs to, of the
	// node it is about, -1 when that node has no location. CheckSDLStyle's
	// findings about a statement's object set -1 and are anchored at
	// StatementStart.
	Start int
	// StatementStart is the byte offset of the first non-space byte of the
	// CheckSDLStyle statement holding the node, where any leading comment
	// starts. CheckSDLIntegrity leaves it 0.
	StatementStart int

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
	// Start is the definition's byte offset in its file, -1 when unknown.
	Start int
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
// function finding is anchored at its statement; a constraint finding at the
// constraint, and a column finding at its first disallowed constraint. A
// function's schema is the first part of a name of two or more parts, so in
// a.b.c it is "a".
//
// Unbalanced parentheses outside literals and comments, or a parse error,
// yield one SDLSyntaxError and nothing else.
func CheckSDLStyle(sql string) []SDLFinding {
	stmts, syntaxError := parseSDL(sql)
	if syntaxError != nil {
		return []SDLFinding{*syntaxError}
	}

	var findings []SDLFinding
	for _, stmt := range stmts {
		if stmt.AST == nil {
			continue
		}
		anchor := sdlStatementAnchor(sql, stmt)

		switch n := stmt.AST.(type) {
		case *ast.CreateStmt:
			findings = checkSDLCreateStmt(n, anchor, findings)
		case *ast.IndexStmt:
			findings = checkSDLIndexStmt(n, anchor, findings)
		case *ast.ViewStmt:
			if n.View != nil && n.View.Schemaname == "" {
				findings = append(findings, anchor.with(SDLFinding{Code: SDLViewWithoutSchema, Name: n.View.Relname}))
			}
		case *ast.CreateSeqStmt:
			if n.Sequence != nil && n.Sequence.Schemaname == "" {
				findings = append(findings, anchor.with(SDLFinding{Code: SDLSequenceWithoutSchema, Name: n.Sequence.Relname}))
			}
		case *ast.CreateFunctionStmt:
			parts := sdlListStrings(n.Funcname)
			if len(parts) < 2 || parts[0] == "" {
				name := ""
				if len(parts) > 0 {
					name = parts[len(parts)-1]
				}
				findings = append(findings, anchor.with(SDLFinding{Code: SDLFunctionWithoutSchema, Name: name}))
			}
		case *ast.AlterSeqStmt:
			if n.Sequence != nil && n.Sequence.Schemaname == "" {
				findings = append(findings, anchor.with(SDLFinding{Code: SDLAlterSequenceWithoutSchema, Name: n.Sequence.Relname}))
			}
		default:
		}
	}
	return findings
}

// parseSDL rejects unbalanced parentheses, which the parser may recover from
// silently, then parses. A failure is returned as an SDLSyntaxError finding.
func parseSDL(sql string) ([]omnipg.Statement, *SDLFinding) {
	if pos, err := validateSDLParenBalance(sql); err != nil {
		return nil, &SDLFinding{Code: SDLSyntaxError, Start: pos, Detail: err.Error()}
	}
	stmts, err := omnipg.Parse(sql)
	if err != nil {
		start := -1
		var parseErr *parser.ParseError
		if errors.As(err, &parseErr) {
			start = parseErr.Position
		}
		return nil, &SDLFinding{Code: SDLSyntaxError, Start: start, Detail: err.Error()}
	}
	return stmts, nil
}

// sdlAnchor is the statement a CheckSDLStyle finding belongs to.
type sdlAnchor struct {
	statementStart int
}

func sdlStatementAnchor(sql string, stmt omnipg.Statement) sdlAnchor {
	start := stmt.ByteStart
	for start < stmt.ByteEnd && sdlIsSpace(sql[start]) {
		start++
	}
	return sdlAnchor{statementStart: start}
}

func sdlIsSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// with anchors f at the statement itself.
func (a sdlAnchor) with(f SDLFinding) SDLFinding {
	return a.at(-1, f)
}

// at anchors f at the node starting at byte offset start, or at the
// statement when start is -1.
func (a sdlAnchor) at(start int, f SDLFinding) SDLFinding {
	f.Start = start
	f.StatementStart = a.statementStart
	return f
}

func checkSDLCreateStmt(n *ast.CreateStmt, stmt sdlAnchor, findings []SDLFinding) []SDLFinding {
	tableName := ""
	schemaName := ""
	if n.Relation != nil {
		tableName = n.Relation.Relname
		schemaName = n.Relation.Schemaname
	}
	if schemaName == "" {
		findings = append(findings, stmt.with(SDLFinding{Code: SDLTableWithoutSchema, Table: tableName}))
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
		firstStart := -1
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
					firstStart = c.Loc.Start
				}
				disallowed = append(disallowed, typeName)
				if c.Contype == ast.CONSTR_FOREIGN && c.Pktable != nil && c.Pktable.Schemaname == "" {
					findings = append(findings, stmt.at(c.Loc.Start, SDLFinding{
						Code:     SDLColumnForeignKeyWithoutSchema,
						Table:    tableName,
						Column:   col.Colname,
						RefTable: c.Pktable.Relname,
					}))
				}
			}
		}
		if len(disallowed) > 0 {
			findings = append(findings, stmt.at(firstStart, SDLFinding{
				Code:            SDLColumnConstraint,
				Table:           tableName,
				Column:          col.Colname,
				ConstraintTypes: disallowed,
			}))
		}
	}

	for _, c := range tableConstraints {
		if c.Conname == "" {
			findings = append(findings, stmt.at(c.Loc.Start, SDLFinding{
				Code:           SDLUnnamedConstraint,
				Table:          tableName,
				ConstraintType: sdlTableConstraintType(c.Contype),
			}))
		}
		if c.Contype == ast.CONSTR_FOREIGN && c.Pktable != nil && c.Pktable.Schemaname == "" {
			findings = append(findings, stmt.at(c.Loc.Start, SDLFinding{
				Code:     SDLForeignKeyWithoutSchema,
				Table:    tableName,
				RefTable: c.Pktable.Relname,
			}))
		}
	}
	return findings
}

func checkSDLIndexStmt(n *ast.IndexStmt, stmt sdlAnchor, findings []SDLFinding) []SDLFinding {
	if n.Idxname == "" {
		findings = append(findings, stmt.with(SDLFinding{Code: SDLUnnamedIndex}))
	}
	if n.Relation != nil && n.Relation.Schemaname == "" {
		findings = append(findings, stmt.with(SDLFinding{
			Code:  SDLIndexWithoutSchema,
			Name:  n.Idxname,
			Table: n.Relation.Relname,
		}))
	}
	return findings
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
// strings, dollar-quoted strings and comments, with the byte offset of the
// first unmatched ")" or the innermost unclosed "(". Block comments do not
// nest, and a single-quoted string ignores backslash escapes.
func validateSDLParenBalance(sql string) (int, error) {
	var open []int
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
			open = append(open, i)
			i++
		case ch == ')':
			if len(open) == 0 {
				return i, errors.New("syntax error: unmatched closing parenthesis")
			}
			open = open[:len(open)-1]
			i++
		default:
			i++
		}
	}
	if len(open) != 0 {
		return open[len(open)-1], errors.New("syntax error: unmatched opening parenthesis")
	}
	return 0, nil
}
