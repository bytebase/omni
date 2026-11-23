package analysis

// CandidateType classifies a SQL or shell completion candidate.
type CandidateType int

const (
	// CandidateTypeNone is the zero value for an unclassified candidate.
	CandidateTypeNone CandidateType = iota

	CandidateKeyword
	CandidateCatalog
	CandidateDatabase
	CandidateSchema
	CandidateTable
	CandidateView
	CandidateMaterializedView
	CandidateColumn
	CandidateFunction
	CandidateProcedure
	CandidatePackage
	CandidateSequence
	CandidateIndex
	CandidateTrigger
	CandidateEvent
	CandidateType_
	CandidateVariable
	CandidateCharset
	CandidateEngine
	CandidatePolicy
	CandidateExtension
	CandidateLogin
	CandidateUser
	CandidateRole
	CandidateSynonym
	CandidateDatabaseLink
	CandidateConstraint

	// MongoDB shell candidate kinds.
	CandidateCollection
	CandidateMethod
	CandidateCursorMethod
	CandidateAggStage
	CandidateQueryOperator
	CandidateBSONHelper
	CandidateShowTarget
	CandidateDbMethod
	CandidateRsMethod
	CandidateShMethod
)

// String returns the lower-case kind name of the candidate type.
func (t CandidateType) String() string {
	switch t {
	case CandidateTypeNone:
		return "none"
	case CandidateKeyword:
		return "keyword"
	case CandidateCatalog:
		return "catalog"
	case CandidateDatabase:
		return "database"
	case CandidateSchema:
		return "schema"
	case CandidateTable:
		return "table"
	case CandidateView:
		return "view"
	case CandidateMaterializedView:
		return "materialized_view"
	case CandidateColumn:
		return "column"
	case CandidateFunction:
		return "function"
	case CandidateProcedure:
		return "procedure"
	case CandidatePackage:
		return "package"
	case CandidateSequence:
		return "sequence"
	case CandidateIndex:
		return "index"
	case CandidateTrigger:
		return "trigger"
	case CandidateEvent:
		return "event"
	case CandidateType_:
		return "type"
	case CandidateVariable:
		return "variable"
	case CandidateCharset:
		return "charset"
	case CandidateEngine:
		return "engine"
	case CandidatePolicy:
		return "policy"
	case CandidateExtension:
		return "extension"
	case CandidateLogin:
		return "login"
	case CandidateUser:
		return "user"
	case CandidateRole:
		return "role"
	case CandidateSynonym:
		return "synonym"
	case CandidateDatabaseLink:
		return "database_link"
	case CandidateConstraint:
		return "constraint"
	case CandidateCollection:
		return "collection"
	case CandidateMethod:
		return "method"
	case CandidateCursorMethod:
		return "cursor_method"
	case CandidateAggStage:
		return "agg_stage"
	case CandidateQueryOperator:
		return "query_operator"
	case CandidateBSONHelper:
		return "bson_helper"
	case CandidateShowTarget:
		return "show_target"
	case CandidateDbMethod:
		return "db_method"
	case CandidateRsMethod:
		return "rs_method"
	case CandidateShMethod:
		return "sh_method"
	default:
		return "unknown"
	}
}
