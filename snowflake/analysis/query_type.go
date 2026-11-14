package analysis

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// QueryType is the permission class of a statement. The values mirror
// Bytebase's base.QueryType.
type QueryType int

const (
	// QueryTypeUnknown is a statement not deliberately mapped. Bytebase's SQL
	// service denies it before execution.
	QueryTypeUnknown QueryType = iota
	// Select is a read-only query.
	Select
	// Explain is an EXPLAIN statement.
	Explain
	// SelectInfoSchema reads system metadata: SHOW and DESCRIBE.
	SelectInfoSchema
	// DDL changes schema.
	DDL
	// DML changes table data.
	DML
)

// String returns the name of the query type.
func (q QueryType) String() string {
	switch q {
	case Select:
		return "SELECT"
	case Explain:
		return "EXPLAIN"
	case SelectInfoSchema:
		return "SELECT_INFO_SCHEMA"
	case DDL:
		return "DDL"
	case DML:
		return "DML"
	default:
		return "UNKNOWN"
	}
}

// QueryTypeOf classifies one parsed statement:
//
//   - SELECT and set operations are Select.
//   - INSERT (single and multi-table), UPDATE, DELETE, MERGE, CALL,
//     EXECUTE IMMEDIATE, EXECUTE TASK and COPY INTO <table> are DML; a stored
//     procedure can change data.
//   - EXPLAIN is Explain whatever it plans.
//   - SHOW and DESCRIBE are SelectInfoSchema, except that a result pipe
//     (SHOW ... ->> <query>) and stmt ->> <query> take the piped query's type.
//   - USE and SET are Select.
//   - COMMENT and every CREATE, ALTER, DROP and UNDROP node are DDL.
//   - Everything else, including nil, is QueryTypeUnknown: TRUNCATE,
//     GRANT/REVOKE, COMMIT/ROLLBACK, PUT/GET, LIST/REMOVE, UNSET, scripting
//     blocks. Mapping them to DDL would let DDL holders run account and stage
//     administration that was never deliberately allowed.
func QueryTypeOf(node ast.Node) QueryType {
	if node == nil {
		return QueryTypeUnknown
	}
	switch n := node.(type) {
	case *ast.SelectStmt, *ast.SetOperationStmt:
		return Select

	case *ast.InsertStmt, *ast.InsertMultiStmt,
		*ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt,
		*ast.CallStmt, *ast.ExecuteImmediateStmt, *ast.ExecuteTaskStmt:
		return DML

	case *ast.ResultScanStmt:
		// The source statement's read-only-ness is checked by ValidateQuery.
		return QueryTypeOf(n.Query)
	case *ast.ExplainStmt:
		return Explain
	case *ast.ShowStmt:
		// Classify a piped SHOW by its query, so the query is checked and
		// masked as one instead of passing as metadata access.
		if n.Pipe != nil {
			return QueryTypeOf(n.Pipe)
		}
		return SelectInfoSchema
	case *ast.DescribeStmt:
		return SelectInfoSchema

	case *ast.UseStmt:
		return Select

	case *ast.SetStmt:
		return Select

	case *ast.CopyIntoTableStmt:
		return DML

	case *ast.CommentStmt:
		return DDL
	}

	// Every omni DDL node is named Create*, Alter*, Drop* or Undrop*.
	tag := node.Tag().String()
	for _, prefix := range []string{"Create", "Alter", "Drop", "Undrop"} {
		if strings.HasPrefix(tag, prefix) {
			return DDL
		}
	}

	return QueryTypeUnknown
}
