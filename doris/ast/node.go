// Package ast defines Doris parse-tree node types.
//
// The package mirrors omni's snowflake/ast and mysql/ast conventions:
// every concrete node implements Node by exposing a Tag() NodeTag method,
// every node carries a Loc field for source-position tracking, and
// walk_generated.go dispatches walker traversal via a type switch.
package ast

// Node is the interface implemented by every Doris parse-tree node.
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
// it is NOT part of the Node interface (mirroring snowflake/ast). Helpers
// in loc.go extract Loc from any Node.
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
