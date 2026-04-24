package pg

import "github.com/bytebase/omni/pg/parser"

// CompletionContext is the parser-native context for SQL completion.
type CompletionContext = parser.CompletionContext

// ScopeSnapshot describes the SELECT-level relation scope visible at a cursor.
type ScopeSnapshot = parser.ScopeSnapshot

// RangeReference is a syntax-level FROM/JOIN reference.
type RangeReference = parser.RangeReference

// RangeReferenceKind classifies a range-table entry visible to completion.
type RangeReferenceKind = parser.RangeReferenceKind

const (
	RangeReferenceRelation  = parser.RangeReferenceRelation
	RangeReferenceSubquery  = parser.RangeReferenceSubquery
	RangeReferenceFunction  = parser.RangeReferenceFunction
	RangeReferenceJoinAlias = parser.RangeReferenceJoinAlias
	RangeReferenceCTE       = parser.RangeReferenceCTE
)

// CollectCompletion returns completion candidates plus a best-effort visible
// relation scope at cursorOffset. Ordinary Parse remains strict.
func CollectCompletion(sql string, cursorOffset int) *CompletionContext {
	return parser.CollectCompletion(sql, cursorOffset)
}
