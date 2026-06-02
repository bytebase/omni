package catalog

import (
	"sort"
	"strings"

	"github.com/bytebase/omni/trino/ast"
)

// Database is one Trino catalog: the second level of the
// catalog -> schema -> table/view hierarchy. It owns a set of schemas keyed
// by normalized name.
//
// The Go type is named Database (rather than Catalog) to keep the package's
// internal hierarchy distinct from the top-level Catalog container and to
// mirror the mysql/tidb catalog packages; user-facing API uses Trino's
// "catalog" terminology (see Catalog.GetCatalog / Catalog.EnsureCatalog).
type Database struct {
	// Name is the normalized catalog name.
	Name string
	// schemas maps a normalized schema name to its Schema.
	schemas map[string]*Schema
}

// newDatabase creates an empty catalog with the given (already-normalized)
// name.
func newDatabase(name string) *Database {
	return &Database{Name: name, schemas: make(map[string]*Schema)}
}

// EnsureSchema returns the named schema, creating and registering an empty one
// if it does not already exist. name must already be normalized (see the
// package doc on Normalize); pass ast.Identifier.Normalize() output or a
// Normalize() result, not raw source text.
func (d *Database) EnsureSchema(name string) *Schema {
	if s, ok := d.schemas[name]; ok {
		return s
	}
	s := newSchema(name)
	d.schemas[name] = s
	return s
}

// GetSchema returns the named schema, or nil if it is not registered. name
// must already be normalized.
func (d *Database) GetSchema(name string) *Schema {
	return d.schemas[name]
}

// HasSchema reports whether the named (normalized) schema is registered.
func (d *Database) HasSchema(name string) bool {
	_, ok := d.schemas[name]
	return ok
}

// Schemas returns the normalized names of all schemas in this catalog in
// sorted order.
func (d *Database) Schemas() []string {
	result := make([]string, 0, len(d.schemas))
	for name := range d.schemas {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

// Normalize folds a raw Trino identifier to the canonical form used for name
// resolution and as the key in this catalog. It mirrors the legacy bytebase
// NormalizeTrinoIdentifier and ast.Identifier.Normalize: a double-quoted
// identifier ("...") has its quotes stripped and its case preserved (quoted
// identifiers are case-sensitive in Trino), while any other (unquoted)
// identifier is folded to lower case.
//
// Use Normalize once, at the boundary, when you hold raw source text. Do not
// call it on a string that is already normalized — an already-folded name like
// "MyView" (the case-preserved result of normalizing the quoted "MyView")
// contains no surrounding quotes, so a second Normalize would wrongly lower it.
// Callers that have a parsed ast.Identifier should prefer the *Ident methods
// (e.g. Catalog.EnsureCatalogIdent) or pass id.Normalize() directly, since the
// AST node carries the quoting metadata this string helper has to re-derive
// from literal quote characters. In particular this string form only strips
// the outer quotes; it does not collapse an embedded doubled-quote escape
// (`""`), because a raw metadata/identifier name is not a re-quoted source
// token. Pass DatabaseMetadata names (which arrive unquoted) and unquoted
// source text here; route quoted source tokens through the AST.
//
// The empty quoted identifier `""` normalizes to the empty string.
func Normalize(ident string) string {
	if len(ident) >= 2 && strings.HasPrefix(ident, `"`) && strings.HasSuffix(ident, `"`) {
		return ident[1 : len(ident)-1]
	}
	return strings.ToLower(ident)
}

// normalizeIdent returns the catalog key for a parsed identifier, deferring to
// the AST node's own quote-aware normalization. A nil identifier yields "".
func normalizeIdent(id *ast.Identifier) string {
	if id == nil {
		return ""
	}
	return id.Normalize()
}
