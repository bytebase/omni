package analysis

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// ExecutionFlags describes how a driver must execute a single statement.
type ExecutionFlags struct {
	// NonTransactional is true for statements that cannot run inside an
	// explicit transaction block (DROP DATABASE, CREATE/DROP INDEX CONCURRENTLY,
	// and VACUUM).
	NonTransactional bool
	// SessionRole is true for non-local SET ROLE statements that must be
	// replayed at session scope.
	SessionRole bool
	// Superuser is true for statements that require superuser privileges
	// (GRANT, GRANT ROLE, and CREATE EXTENSION).
	Superuser bool
}

// vacuumReg is the one statement classification CockroachDB's parser cannot
// answer: VACUUM has no tree node, so the anchored keyword is enough.
var vacuumReg = regexp.MustCompile(`(?i)^\s*VACUUM`)

// ClassifyExecution reports execution flags for a single CockroachDB statement.
// CockroachDB's own AST decides, so keywords inside string literals or comments
// do not affect classification. Unparseable or multi-statement input returns
// zero flags.
func ClassifyExecution(stmt string) ExecutionFlags {
	parsed, err := Parse(stmt)
	if err != nil || len(parsed) != 1 {
		return ExecutionFlags{}
	}
	switch node := parsed[0].AST.(type) {
	case *tree.DropDatabase:
		return ExecutionFlags{NonTransactional: true}
	case *tree.CreateIndex:
		return ExecutionFlags{NonTransactional: node.Concurrently}
	case *tree.DropIndex:
		return ExecutionFlags{NonTransactional: node.Concurrently}
	case *tree.SetVar:
		return ExecutionFlags{SessionRole: strings.EqualFold(node.Name, "role") && !node.Local}
	case *tree.Grant, *tree.GrantRole, *tree.CreateExtension:
		return ExecutionFlags{Superuser: true}
	}
	return ExecutionFlags{NonTransactional: vacuumReg.MatchString(stmt)}
}

// CreateDatabaseName parses script and returns the target database name from
// the first CREATE DATABASE statement it contains.
func CreateDatabaseName(script string) (string, error) {
	parsed, err := Parse(script)
	if err != nil {
		return "", err
	}
	for _, stmt := range parsed {
		if create, ok := stmt.AST.(*tree.CreateDatabase); ok && create.Name != "" {
			return string(create.Name), nil
		}
	}
	return "", errors.New("database name not found")
}

// IsDoBlock reports whether stmt parses as a single PL/pgSQL anonymous code
// block (DO statement).
func IsDoBlock(stmt string) bool {
	stmts, err := Parse(stmt)
	if err != nil || len(stmts) != 1 {
		return false
	}
	_, ok := stmts[0].AST.(*tree.DoBlock)
	return ok
}

// ReturnsRows parses statement and reports whether any statement in it returns
// result rows according to CockroachDB's parser.
func ReturnsRows(statement string) (bool, error) {
	parsed, err := Parse(statement)
	if err != nil {
		return false, err
	}
	for _, stmt := range parsed {
		if stmt.AST != nil && stmt.AST.StatementReturnType() == tree.Rows {
			return true, nil
		}
	}
	return false, nil
}

// IsSelectStatement reports whether statement parses as a single SELECT
// (including parenthesized SELECT) statement.
func IsSelectStatement(statement string) bool {
	parsed, err := Parse(statement)
	if err != nil || len(parsed) != 1 {
		return false
	}
	switch parsed[0].AST.(type) {
	case *tree.Select, *tree.ParenSelect:
		return true
	}
	return false
}

// IndexDefinition holds the parsed structure of a CockroachDB CREATE INDEX
// statement.
type IndexDefinition struct {
	Expressions []string
	Descending  []bool
	Unique      bool
}

// ParseIndexDefinition parses a single CREATE INDEX statement (such as
// pg_indexes.indexdef) and extracts its key expressions, sort directions, and
// uniqueness.
func ParseIndexDefinition(statement string) (*IndexDefinition, error) {
	nodes, err := Parse(statement)
	if err != nil {
		return nil, err
	}
	if len(nodes) != 1 {
		return nil, fmt.Errorf("invalid number of statement %v, expecting one", len(nodes))
	}
	node, ok := nodes[0].AST.(*tree.CreateIndex)
	if !ok {
		return nil, fmt.Errorf("invalid statement type %T, expecting CreateIndex", nodes[0].AST)
	}
	def := &IndexDefinition{
		Unique: node.Unique,
	}
	for _, indexElem := range node.Columns {
		switch {
		case indexElem.Column != "":
			def.Expressions = append(def.Expressions, indexElem.Column.String())
		case indexElem.Expr != nil:
			def.Expressions = append(def.Expressions, indexElem.Expr.String())
		case indexElem.OpClass != "":
			def.Expressions = append(def.Expressions, indexElem.OpClass.String())
		default:
			continue
		}
		def.Descending = append(def.Descending, indexElem.Direction == tree.Descending)
	}
	return def, nil
}
