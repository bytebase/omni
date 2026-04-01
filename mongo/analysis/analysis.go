// Package analysis extracts structural information from parsed MongoDB statements.
package analysis

import "github.com/bytebase/omni/mongo/ast"

// Operation classifies a MongoDB statement by its analytical behavior.
type Operation int

const (
	OpUnknown   Operation = iota
	OpFind                // find
	OpFindOne             // findOne
	OpAggregate           // aggregate
	OpCount               // countDocuments, estimatedDocumentCount, count
	OpDistinct            // distinct
	OpRead                // getIndexes, stats, storageSize, etc.
	OpWrite               // insertOne, updateMany, deleteOne, etc.
	OpAdmin               // createIndex, drop, renameCollection, etc.
	OpInfo                // show dbs, getCollectionNames, serverStatus, etc.
	OpExplain             // .explain() wrapped
)

// IsRead returns true for read operations (find, findOne, aggregate, count, distinct, read).
func (op Operation) IsRead() bool {
	switch op {
	case OpFind, OpFindOne, OpAggregate, OpCount, OpDistinct, OpRead:
		return true
	}
	return false
}

// IsWrite returns true for write operations.
func (op Operation) IsWrite() bool { return op == OpWrite }

// IsAdmin returns true for administrative operations.
func (op Operation) IsAdmin() bool { return op == OpAdmin }

// IsInfo returns true for informational/metadata operations.
func (op Operation) IsInfo() bool { return op == OpInfo }

// String returns a human-readable name for the operation.
func (op Operation) String() string {
	switch op {
	case OpFind:
		return "find"
	case OpFindOne:
		return "findOne"
	case OpAggregate:
		return "aggregate"
	case OpCount:
		return "count"
	case OpDistinct:
		return "distinct"
	case OpRead:
		return "read"
	case OpWrite:
		return "write"
	case OpAdmin:
		return "admin"
	case OpInfo:
		return "info"
	case OpExplain:
		return "explain"
	default:
		return "unknown"
	}
}

// StatementAnalysis is the result of analyzing a parsed MongoDB statement.
type StatementAnalysis struct {
	Operation        Operation
	MethodName       string   // original method name (e.g. "countDocuments" when Operation is OpCount)
	Collection       string   // target collection; empty for db-level/show commands
	PredicateFields  []string // sorted dot-paths from query filter document
	PipelineStages   []string // aggregate only: stage names in order
	ShapePreserving  bool     // true when all pipeline stages preserve document structure
	UnsupportedStage string   // first non-shape-preserving, non-join stage
	Joins            []JoinInfo
}

// JoinInfo records a join extracted from a $lookup or $graphLookup pipeline stage.
type JoinInfo struct {
	Collection string // the "from" field
	AsField    string // the "as" field
}

// Analyze extracts structural information from a single parsed MongoDB AST node.
// Returns nil for unrecognizable nodes.
func Analyze(node ast.Node) *StatementAnalysis {
	_ = node
	return nil
}
