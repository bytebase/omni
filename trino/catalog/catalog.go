// Package catalog provides in-memory metadata context for Trino
// auto-completion and query-span/lineage analysis.
//
// Unlike the flat partiql (DynamoDB tables) and mongo (collections) catalogs,
// Trino has a three-level object namespace — catalog -> schema -> table/view
// -> column. A fully-qualified Trino name is catalog.schema.table; shorter
// names resolve against the session's current catalog and schema (set via
// USE catalog.schema). This package models that hierarchy and exposes the
// read API the completion package needs (catalog/schema/table/view/column
// candidates and qualified-name resolution).
//
// # Identifier normalization (the load-bearing contract)
//
// Trino folds unquoted identifiers to lower case and preserves the case of
// double-quoted identifiers; quoted names are therefore case-sensitive. This
// catalog stores and looks up every name by its already-normalized form and
// performs no folding of its own — that keeps it from double-folding a name
// like "MyView" (the case-preserved result of normalizing a quoted "MyView").
//
// Callers normalize once, at the boundary:
//   - With a parsed ast.Identifier / ast.QualifiedName, use the *Ident methods
//     or pass id.Normalize() / qn.NormalizedParts() — the AST node carries the
//     quoting metadata.
//   - With raw source text (e.g. names from a DatabaseMetadata snapshot or a
//     test fixture), call the package Normalize function first.
//
// The package depends only on trino/ast for that normalization contract; it
// does not depend on the parser.
package catalog

import (
	"sort"

	"github.com/bytebase/omni/trino/ast"
)

// Catalog is the in-memory Trino catalog: the top of the
// catalog -> schema -> table/view -> column hierarchy. It also carries the
// session's current catalog and schema, used to resolve unqualified and
// partially-qualified names the way Trino does.
//
// Catalog is not safe for concurrent mutation; build it once (e.g. from a
// DatabaseMetadata snapshot) and then read from it.
type Catalog struct {
	// databases maps a normalized catalog name to its Database. Trino calls
	// the top-level namespace a "catalog"; we name the Go type Database to
	// stay consistent with the mysql/tidb catalog packages, but every
	// user-facing accessor uses Trino's "catalog" terminology.
	databases map[string]*Database

	// currentCatalog and currentSchema hold the session context set by
	// USE catalog.schema (or USE schema). They default to empty, meaning
	// "no session catalog/schema selected"; resolution of unqualified names
	// then fails rather than guessing.
	currentCatalog string
	currentSchema  string
}

// New creates an empty Catalog with no session catalog or schema selected.
func New() *Catalog {
	return &Catalog{databases: make(map[string]*Database)}
}

// SetCurrentCatalog records the session's current catalog (as set by
// USE catalog.schema). name must already be normalized.
func (c *Catalog) SetCurrentCatalog(name string) { c.currentCatalog = name }

// CurrentCatalog returns the session's current catalog name (normalized), or
// "" if none is selected.
func (c *Catalog) CurrentCatalog() string { return c.currentCatalog }

// SetCurrentSchema records the session's current schema. name must already be
// normalized.
func (c *Catalog) SetCurrentSchema(name string) { c.currentSchema = name }

// CurrentSchema returns the session's current schema name (normalized), or
// "" if none is selected.
func (c *Catalog) CurrentSchema() string { return c.currentSchema }

// EnsureCatalog returns the named catalog, creating and registering an empty
// one if it does not already exist. name must already be normalized (see the
// package doc); use EnsureCatalogIdent when you hold a parsed ast.Identifier.
func (c *Catalog) EnsureCatalog(name string) *Database {
	if db, ok := c.databases[name]; ok {
		return db
	}
	db := newDatabase(name)
	c.databases[name] = db
	return db
}

// EnsureCatalogIdent is EnsureCatalog keyed by a parsed identifier, using the
// AST node's quote-aware normalization.
func (c *Catalog) EnsureCatalogIdent(id *ast.Identifier) *Database {
	return c.EnsureCatalog(normalizeIdent(id))
}

// GetCatalog returns the named catalog, or nil if it is not registered. name
// must already be normalized.
func (c *Catalog) GetCatalog(name string) *Database {
	return c.databases[name]
}

// HasCatalog reports whether the named (normalized) catalog is registered.
func (c *Catalog) HasCatalog(name string) bool {
	_, ok := c.databases[name]
	return ok
}

// Catalogs returns the normalized names of all registered catalogs in sorted
// order. The slice is freshly allocated and safe for the caller to retain.
func (c *Catalog) Catalogs() []string {
	result := make([]string, 0, len(c.databases))
	for name := range c.databases {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}
