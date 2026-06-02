package ast

// NodeLoc returns the source location of n, or NoLoc() if n is nil or its
// concrete type carries no Loc field. Every concrete node type added under
// trino/ast must add a case here.
//
// n may be either an untyped nil Node or a typed-nil pointer (e.g. a
// (*Identifier)(nil) boxed in a non-nil interface, which occurs when an
// optional child field is left unset); both yield NoLoc() rather than a
// panic. The pattern matches doris/ast.NodeLoc and snowflake/ast.NodeLoc
// but additionally hardens the per-case dereference against typed nils.
func NodeLoc(n Node) Loc {
	if n == nil {
		return NoLoc()
	}
	switch v := n.(type) {
	case *File:
		if v == nil {
			return NoLoc()
		}
		return v.Loc
	case *Identifier:
		if v == nil {
			return NoLoc()
		}
		return v.Loc
	case *QualifiedName:
		if v == nil {
			return NoLoc()
		}
		return v.Loc
	default:
		return NoLoc()
	}
}

// SpanNodes returns the smallest Loc that covers every node in nodes.
// Nil entries and nodes whose Loc is invalid (per Loc.IsValid) are skipped.
// Returns NoLoc() when no node has a valid Loc (including the empty-args
// case). Skipping invalid locs explicitly keeps a half-set Loc (e.g. a
// node with End == -1 from an upstream parser slip) from leaking into the
// result, since Loc.Merge returns the other side verbatim when one side is
// invalid.
func SpanNodes(nodes ...Node) Loc {
	out := NoLoc()
	for _, n := range nodes {
		if n == nil {
			continue
		}
		loc := NodeLoc(n)
		if !loc.IsValid() {
			continue
		}
		out = out.Merge(loc)
	}
	return out
}
