// Package ast defines GoogleSQL parse-tree node types.
//
// GoogleSQL is the SQL dialect shared by Google BigQuery and Google Cloud
// Spanner; one omni parser (and therefore one AST) serves both engines, with
// dialect-specific deltas handled above this layer. The grammar is a
// hand-port of Google's open-source ZetaSQL reference, so the node shapes here
// follow ZetaSQL's tree.
//
// The package is structured to mirror omni's pg/ast, mysql/ast, and
// snowflake/ast conventions: every concrete node implements Node by exposing a
// Tag() NodeTag method, every node carries a Loc field for source-position
// tracking, and walk_generated.go dispatches walker traversal via a generated
// type switch.
package ast

// Node is the interface implemented by every GoogleSQL parse-tree node.
//
// Tag returns a unique NodeTag identifying the concrete type. Use Tag for
// fast switch dispatch in hot paths; use a Go type assertion when you need
// to access the concrete fields.
type Node interface {
	Tag() NodeTag
}

// Loc represents a source byte range. -1 means "unknown" for either field.
//
// Loc is a value type embedded as a plain field on every concrete node;
// it is NOT part of the Node interface (mirroring pg/ast). Helpers in loc.go
// extract Loc from any Node.
type Loc struct {
	Start int // inclusive start byte offset (-1 if unknown)
	End   int // exclusive end byte offset   (-1 if unknown)
}

// NoLoc returns a Loc with both Start and End set to -1 (unknown).
func NoLoc() Loc {
	return Loc{Start: -1, End: -1}
}

// IsValid reports whether both Start and End are non-negative.
func (l Loc) IsValid() bool {
	return l.Start >= 0 && l.End >= 0
}

// Contains reports whether l fully contains other (inclusive on both ends).
// Returns false if either Loc is invalid.
func (l Loc) Contains(other Loc) bool {
	if !l.IsValid() || !other.IsValid() {
		return false
	}
	return l.Start <= other.Start && other.End <= l.End
}

// Merge returns the smallest Loc that contains both l and other.
// If either side is invalid, returns the other side. If both are invalid,
// returns NoLoc().
func (l Loc) Merge(other Loc) Loc {
	if !l.IsValid() {
		return other
	}
	if !other.IsValid() {
		return l
	}
	out := l
	if other.Start < out.Start {
		out.Start = other.Start
	}
	if other.End > out.End {
		out.End = other.End
	}
	return out
}
