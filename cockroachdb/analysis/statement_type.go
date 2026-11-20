package analysis

import (
	"slices"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser/statements"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// StatementType is the kind of change a statement makes. Its value is the
// name of the matching Bytebase storepb.StatementType.
type StatementType string

// StatementTypeUnspecified is a statement the classifier does not name.
const StatementTypeUnspecified StatementType = ""

const (
	AlterDatabase   StatementType = "ALTER_DATABASE"
	AlterIndex      StatementType = "ALTER_INDEX"
	AlterSequence   StatementType = "ALTER_SEQUENCE"
	AlterTable      StatementType = "ALTER_TABLE"
	AlterType       StatementType = "ALTER_TYPE"
	AlterView       StatementType = "ALTER_VIEW"
	Comment         StatementType = "COMMENT"
	CreateDatabase  StatementType = "CREATE_DATABASE"
	CreateExtension StatementType = "CREATE_EXTENSION"
	CreateFunction  StatementType = "CREATE_FUNCTION"
	CreateIndex     StatementType = "CREATE_INDEX"
	CreateSchema    StatementType = "CREATE_SCHEMA"
	CreateSequence  StatementType = "CREATE_SEQUENCE"
	CreateTable     StatementType = "CREATE_TABLE"
	CreateTrigger   StatementType = "CREATE_TRIGGER"
	CreateType      StatementType = "CREATE_TYPE"
	CreateView      StatementType = "CREATE_VIEW"
	Delete          StatementType = "DELETE"
	DropDatabase    StatementType = "DROP_DATABASE"
	DropFunction    StatementType = "DROP_FUNCTION"
	DropIndex       StatementType = "DROP_INDEX"
	DropSchema      StatementType = "DROP_SCHEMA"
	DropSequence    StatementType = "DROP_SEQUENCE"
	DropTable       StatementType = "DROP_TABLE"
	DropTrigger     StatementType = "DROP_TRIGGER"
	DropType        StatementType = "DROP_TYPE"
	DropView        StatementType = "DROP_VIEW"
	Execute         StatementType = "EXECUTE"
	Insert          StatementType = "INSERT"
	RenameIndex     StatementType = "RENAME_INDEX"
	RenameSchema    StatementType = "RENAME_SCHEMA"
	RenameSequence  StatementType = "RENAME_SEQUENCE"
	Truncate        StatementType = "TRUNCATE"
	Update          StatementType = "UPDATE"
)

// StatementTypesFromParsed returns the StatementTypes of opaque
// ParsedStatement values, keeping cockroachdb-parser types out of callers.
func StatementTypesFromParsed(stmts []ParsedStatement) []StatementType {
	return StatementTypes(unwrapStatements(stmts))
}

// StatementTypes returns the type of each statement that has one, in order,
// each followed by the types of the mutations in its CTEs and statement
// sources, each type once per statement. EXPLAIN ANALYZE takes the types of
// the statement it runs. A statement the classifier does not name adds none.
func StatementTypes(stmts statements.Statements) []StatementType {
	var statementTypes []StatementType
	for _, parsed := range stmts {
		var types []StatementType
		add := func(statementType StatementType) {
			if statementType != StatementTypeUnspecified && !slices.Contains(types, statementType) {
				types = append(types, statementType)
			}
		}
		stmt := parsed.AST
		// EXPLAIN ANALYZE executes the statement it explains.
		if explain, ok := stmt.(*tree.ExplainAnalyze); ok {
			stmt = explain.Statement
		}
		add(getStatementType(stmt))
		if createTable, ok := stmt.(*tree.CreateTable); ok && createTable.AsSource != nil {
			stmt = createTable.AsSource
		}
		for _, mutationType := range getMutationTypes(stmt) {
			add(mutationType)
		}
		statementTypes = append(statementTypes, types...)
	}
	return statementTypes
}

func getStatementType(stmt tree.Statement) StatementType {
	switch n := stmt.(type) {
	case *tree.CreateDatabase:
		return CreateDatabase
	case *tree.CreateTable:
		return CreateTable
	case *tree.CreateView:
		return CreateView
	case *tree.CreateIndex:
		return CreateIndex
	case *tree.CreateSequence:
		return CreateSequence
	case *tree.CreateSchema:
		return CreateSchema
	case *tree.CreateRoutine:
		// Procedures are functions here, as PostgreSQL classifies CREATE PROCEDURE.
		return CreateFunction
	case *tree.CreateTrigger:
		return CreateTrigger
	case *tree.CreateExtension:
		return CreateExtension
	case *tree.CreateType:
		return CreateType

	case *tree.DropDatabase:
		return DropDatabase
	case *tree.DropTable:
		return DropTable
	case *tree.DropView:
		// A materialized view holds data, so dropping one is as risky as dropping a table.
		if n.IsMaterialized {
			return DropTable
		}
		return DropView
	case *tree.DropIndex:
		return DropIndex
	case *tree.DropSequence:
		return DropSequence
	case *tree.DropSchema:
		return DropSchema
	case *tree.DropType:
		return DropType
	case *tree.DropTrigger:
		return DropTrigger
	case *tree.DropRoutine:
		return DropFunction
	case *tree.Truncate:
		return Truncate

	case *tree.RenameDatabase, *tree.AlterDatabaseOwner, *tree.ReparentDatabase,
		*tree.AlterDatabaseAddRegion, *tree.AlterDatabaseDropRegion, *tree.AlterDatabasePrimaryRegion,
		*tree.AlterDatabaseSurvivalGoal, *tree.AlterDatabasePlacement, *tree.AlterDatabaseAddSuperRegion,
		*tree.AlterDatabaseDropSuperRegion, *tree.AlterDatabaseAlterSuperRegion, *tree.AlterDatabaseSecondaryRegion,
		*tree.AlterDatabaseDropSecondaryRegion, *tree.AlterDatabaseSetZoneConfigExtension:
		return AlterDatabase
	case *tree.SetZoneConfig:
		switch {
		case n.Database != "":
			return AlterDatabase
		case n.TableOrIndex.Index != "":
			return AlterIndex
		case n.TableOrIndex.Table.ObjectName != "":
			return AlterTable
		default:
			return StatementTypeUnspecified
		}
	case *tree.AlterRoleSet:
		// ALTER DATABASE d SET parses as ALTER ROLE ALL IN DATABASE d SET, which sets the same defaults.
		if n.AllRoles && n.DatabaseName != "" {
			return AlterDatabase
		}
		return StatementTypeUnspecified
	case *tree.AlterTable, *tree.AlterTableLocality:
		return AlterTable
	case *tree.AlterTableSetSchema:
		return getRelationAlterType(n.IsView && !n.IsMaterialized, n.IsSequence)
	case *tree.AlterTableOwner:
		return getRelationAlterType(n.IsView && !n.IsMaterialized, n.IsSequence)
	case *tree.AlterIndex, *tree.AlterIndexVisible:
		return AlterIndex
	case *tree.AlterSequence:
		return AlterSequence
	case *tree.AlterType:
		return AlterType

	case *tree.RenameTable:
		switch {
		case n.IsSequence:
			return RenameSequence
		case n.IsView && !n.IsMaterialized:
			return AlterView
		default:
			return AlterTable
		}
	case *tree.RenameIndex:
		return RenameIndex
	case *tree.AlterSchema:
		if _, ok := n.Cmd.(*tree.AlterSchemaRename); ok {
			return RenameSchema
		}
		return StatementTypeUnspecified

	case *tree.CommentOnColumn, *tree.CommentOnConstraint, *tree.CommentOnDatabase, *tree.CommentOnIndex,
		*tree.CommentOnSchema, *tree.CommentOnTable, *tree.CommentOnType:
		return Comment

	case *tree.Insert:
		return Insert
	case *tree.Update:
		return Update
	case *tree.Delete:
		return Delete

	// A procedure, DO block or prepared statement runs code this classifier does not look into.
	case *tree.Call, *tree.DoBlock, *tree.Execute:
		return Execute
	default:
		return StatementTypeUnspecified
	}
}

func getRelationAlterType(isView bool, isSequence bool) StatementType {
	switch {
	case isSequence:
		return AlterSequence
	case isView:
		return AlterView
	default:
		return AlterTable
	}
}
