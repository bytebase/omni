package ast

import "strings"

// This file holds the foundational Trino parse-tree node types. The ast-core
// node ships the File root container plus the two identifier nodes the whole
// grammar pivots on (Trino's `identifier` and `qualifiedName` grammar rules);
// later migration nodes (lexer, parser, types, expressions, statements)
// populate the rest.

// File is the root node of a parsed Trino source file. It holds the
// top-level statement list and the byte range covering the entire file.
type File struct {
	Stmts []Node
	Loc   Loc
}

// Tag implements Node.
func (n *File) Tag() NodeTag { return T_File }

// Compile-time assertion that *File satisfies Node.
var _ Node = (*File)(nil)

// ---------------------------------------------------------------------------
// Identifier nodes
// ---------------------------------------------------------------------------

// Identifier represents a single Trino identifier — one component of a name.
// It corresponds to the legacy grammar's `identifier` rule, whose four
// alternatives are: an unquoted identifier, a `"double-quoted"` identifier,
// a “ `backtick-quoted` “ identifier, and a digit-leading identifier.
//
// Value stores the identifier's source text with the surrounding quote
// characters already stripped and any doubled-quote escape (`""` or “ “ “)
// collapsed to a single quote. Quoting metadata is preserved separately
// because Trino's identifier resolution is case-sensitive for quoted names
// and case-insensitive (folded to lower case) for unquoted ones — see
// Normalize.
type Identifier struct {
	// Value is the unquoted identifier text (quotes stripped, escapes
	// collapsed). Case is preserved exactly as written in the source.
	Value string
	// Quoted reports whether the source identifier was delimited
	// (double-quoted or backtick-quoted). Digit-leading and bare
	// identifiers are not quoted.
	Quoted bool
	// QuoteRune is the delimiter rune used when Quoted is true: '"' for a
	// standard double-quoted identifier or '`' for a backtick-quoted one.
	// It is 0 when Quoted is false.
	QuoteRune rune
	Loc       Loc
}

// Tag implements Node.
func (n *Identifier) Tag() NodeTag { return T_Identifier }

// Compile-time assertion that *Identifier satisfies Node.
var _ Node = (*Identifier)(nil)

// Normalize returns the identifier in the canonical form Trino uses for
// name resolution: quoted identifiers keep their exact case, while unquoted
// identifiers are folded to lower case (Trino identifiers are
// case-insensitive unless quoted; note Trino folds to lower case, unlike
// Snowflake which folds to upper). This mirrors the legacy bytebase helper
// NormalizeTrinoIdentifier so downstream lineage/completion consumers can
// compare names consistently.
//
// Use Normalize for name comparison; use String for source-faithful
// rendering (deparse, error messages). This mirrors the
// snowflake/ast.Ident.Normalize / String split.
func (n *Identifier) Normalize() string {
	if n.Quoted {
		return n.Value
	}
	return strings.ToLower(n.Value)
}

// String returns the source-faithful form of the identifier, re-adding the
// original delimiter when it was quoted and escaping any embedded delimiter
// by doubling it (`"` -> `""`, “ ` “ -> “ “ “). Unquoted identifiers
// are returned verbatim with case preserved. Mirrors
// snowflake/ast.Ident.String; useful for deparse and error messages.
func (n *Identifier) String() string {
	if !n.Quoted {
		return n.Value
	}
	q := n.QuoteRune
	if q == 0 {
		// Quoted with no recorded delimiter: default to the standard
		// double-quote so output is still valid Trino syntax.
		q = '"'
	}
	qs := string(q)
	return qs + strings.ReplaceAll(n.Value, qs, qs+qs) + qs
}

// ---------------------------------------------------------------------------
// QualifiedName
// ---------------------------------------------------------------------------

// QualifiedName represents a dot-separated chain of identifiers, matching the
// legacy grammar's `qualifiedName` rule (`identifier ('.' identifier)*`).
// In Trino a three-part name is catalog.schema.table; shorter names are
// resolved against the session's current catalog/schema.
//
// Parts holds the component identifiers in source order, preserving each
// component's quoting metadata so callers can normalize per-part (catalog
// and schema may be quoted independently of the object name).
type QualifiedName struct {
	Parts []*Identifier
	Loc   Loc
}

// Tag implements Node.
func (n *QualifiedName) Tag() NodeTag { return T_QualifiedName }

// Compile-time assertion that *QualifiedName satisfies Node.
var _ Node = (*QualifiedName)(nil)

// NormalizedParts returns the per-component normalized names (see
// Identifier.Normalize). Nil component entries are skipped.
func (n *QualifiedName) NormalizedParts() []string {
	if n == nil {
		return nil
	}
	parts := make([]string, 0, len(n.Parts))
	for _, p := range n.Parts {
		if p == nil {
			continue
		}
		parts = append(parts, p.Normalize())
	}
	return parts
}

// Normalize returns the canonical dotted form with each component normalized
// per Trino resolution rules (e.g., "catalog.schema.table"). Use this for
// name comparison. Nil components are skipped. Mirrors
// snowflake/ast.ObjectName.Normalize.
func (n *QualifiedName) Normalize() string {
	return strings.Join(n.NormalizedParts(), ".")
}

// String returns the source-faithful dotted form, re-quoting each component
// that was originally quoted and preserving case (see Identifier.String).
// Nil components are skipped. Use this for deparse and error messages;
// use Normalize for name comparison. Mirrors snowflake/ast.ObjectName.String.
func (n *QualifiedName) String() string {
	if n == nil {
		return ""
	}
	parts := make([]string, 0, len(n.Parts))
	for _, p := range n.Parts {
		if p == nil {
			continue
		}
		parts = append(parts, p.String())
	}
	return strings.Join(parts, ".")
}
