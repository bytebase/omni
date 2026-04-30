package mssql

import "github.com/bytebase/omni/mssql/parser"

// CompletionContext is the parser-native context for SQL completion.
type CompletionContext = parser.CompletionContext

// CompletionIntent describes the object class and qualifier implied by the
// grammar position at the cursor.
type CompletionIntent = parser.CompletionIntent

// ObjectKind classifies catalog object completion intent.
type ObjectKind = parser.ObjectKind

// MultipartName is a partially typed MSSQL object name.
type MultipartName = parser.MultipartName

// ScopeSnapshot describes relation scope visible at a cursor.
type ScopeSnapshot = parser.ScopeSnapshot

// RangeReference is a syntax-level table expression reference.
type RangeReference = parser.RangeReference

// RangeReferenceKind classifies a range-table entry visible to completion.
type RangeReferenceKind = parser.RangeReferenceKind

const (
	ObjectKindUnknown   = parser.ObjectKindUnknown
	ObjectKindDatabase  = parser.ObjectKindDatabase
	ObjectKindSchema    = parser.ObjectKindSchema
	ObjectKindTable     = parser.ObjectKindTable
	ObjectKindView      = parser.ObjectKindView
	ObjectKindSequence  = parser.ObjectKindSequence
	ObjectKindProcedure = parser.ObjectKindProcedure
	ObjectKindFunction  = parser.ObjectKindFunction
	ObjectKindType      = parser.ObjectKindType
	ObjectKindColumn    = parser.ObjectKindColumn
)

const (
	RangeReferenceRelation      = parser.RangeReferenceRelation
	RangeReferenceSubquery      = parser.RangeReferenceSubquery
	RangeReferenceFunction      = parser.RangeReferenceFunction
	RangeReferenceJoinAlias     = parser.RangeReferenceJoinAlias
	RangeReferenceCTE           = parser.RangeReferenceCTE
	RangeReferenceValues        = parser.RangeReferenceValues
	RangeReferenceTableVariable = parser.RangeReferenceTableVariable
	RangeReferenceDMLTarget     = parser.RangeReferenceDMLTarget
	RangeReferenceMergeTarget   = parser.RangeReferenceMergeTarget
	RangeReferenceMergeSource   = parser.RangeReferenceMergeSource
)

// CollectCompletion returns completion candidates plus a best-effort visible
// relation scope at cursorOffset. Ordinary Parse remains strict.
func CollectCompletion(sql string, cursorOffset int) *CompletionContext {
	return parser.CollectCompletion(sql, cursorOffset)
}
