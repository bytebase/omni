package analysis

// StatementType is the kind of change a statement makes. The values mirror
// Bytebase's storepb.StatementType names so callers across engines map them 1:1.
type StatementType int

const (
	// StatementTypeUnspecified is a statement the classifier does not name.
	StatementTypeUnspecified StatementType = iota

	CreateDatabase
	CreateTable
	CreateView
	CreateIndex
	CreateSequence
	CreateSchema
	CreateFunction
	CreateTrigger
	CreateProcedure
	CreateEvent
	CreateExtension
	CreateType

	DropDatabase
	DropTable
	DropView
	DropIndex
	DropSequence
	DropSchema
	DropFunction
	DropTrigger
	DropProcedure
	DropEvent
	DropExtension
	DropType

	AlterDatabase
	AlterTable
	AlterView
	AlterSequence
	AlterEvent
	AlterType
	AlterIndex

	Truncate
	Rename
	RenameIndex
	RenameSchema
	RenameSequence
	Comment

	Insert
	Update
	Delete
	Merge

	Execute
)

// String returns the canonical storepb.StatementType name of the statement type.
func (s StatementType) String() string {
	switch s {
	case CreateDatabase:
		return "CREATE_DATABASE"
	case CreateTable:
		return "CREATE_TABLE"
	case CreateView:
		return "CREATE_VIEW"
	case CreateIndex:
		return "CREATE_INDEX"
	case CreateSequence:
		return "CREATE_SEQUENCE"
	case CreateSchema:
		return "CREATE_SCHEMA"
	case CreateFunction:
		return "CREATE_FUNCTION"
	case CreateTrigger:
		return "CREATE_TRIGGER"
	case CreateProcedure:
		return "CREATE_PROCEDURE"
	case CreateEvent:
		return "CREATE_EVENT"
	case CreateExtension:
		return "CREATE_EXTENSION"
	case CreateType:
		return "CREATE_TYPE"
	case DropDatabase:
		return "DROP_DATABASE"
	case DropTable:
		return "DROP_TABLE"
	case DropView:
		return "DROP_VIEW"
	case DropIndex:
		return "DROP_INDEX"
	case DropSequence:
		return "DROP_SEQUENCE"
	case DropSchema:
		return "DROP_SCHEMA"
	case DropFunction:
		return "DROP_FUNCTION"
	case DropTrigger:
		return "DROP_TRIGGER"
	case DropProcedure:
		return "DROP_PROCEDURE"
	case DropEvent:
		return "DROP_EVENT"
	case DropExtension:
		return "DROP_EXTENSION"
	case DropType:
		return "DROP_TYPE"
	case AlterDatabase:
		return "ALTER_DATABASE"
	case AlterTable:
		return "ALTER_TABLE"
	case AlterView:
		return "ALTER_VIEW"
	case AlterSequence:
		return "ALTER_SEQUENCE"
	case AlterEvent:
		return "ALTER_EVENT"
	case AlterType:
		return "ALTER_TYPE"
	case AlterIndex:
		return "ALTER_INDEX"
	case Truncate:
		return "TRUNCATE"
	case Rename:
		return "RENAME"
	case RenameIndex:
		return "RENAME_INDEX"
	case RenameSchema:
		return "RENAME_SCHEMA"
	case RenameSequence:
		return "RENAME_SEQUENCE"
	case Comment:
		return "COMMENT"
	case Insert:
		return "INSERT"
	case Update:
		return "UPDATE"
	case Delete:
		return "DELETE"
	case Merge:
		return "MERGE"
	case Execute:
		return "EXECUTE"
	default:
		return "STATEMENT_TYPE_UNSPECIFIED"
	}
}
