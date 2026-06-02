package catalog

import "github.com/bytebase/omni/trino/ast"

// ResolvedRelation is the outcome of resolving a qualified name to a relation.
// Exactly one of Table or View is non-nil when Found is true; both are nil
// when Found is false. The Catalog/Schema fields record the (normalized)
// container the name resolved against, which completion uses to offer sibling
// candidates.
type ResolvedRelation struct {
	Catalog string
	Schema  string
	Table   *Table
	View    *View
	Found   bool
}

// ColumnNames returns the resolved relation's column names in order, or nil if
// nothing was found.
func (r ResolvedRelation) ColumnNames() []string {
	switch {
	case r.Table != nil:
		return r.Table.ColumnNames()
	case r.View != nil:
		return r.View.ColumnNames()
	default:
		return nil
	}
}

// ResolveRelation resolves a dotted name to a table or view, applying Trino's
// name-resolution rules against the session's current catalog and schema:
//
//   - 1 part  [object]                  -> currentCatalog.currentSchema.object
//   - 2 parts [schema, object]          -> currentCatalog.schema.object
//   - 3 parts [catalog, schema, object] -> catalog.schema.object
//
// parts must already be normalized — pass ast.QualifiedName.NormalizedParts()
// (or ResolveQualifiedName, which does that for you). Resolution prefers a
// table over a view when both exist under the resolved name (the same
// precedence as Schema.Relations). A name with zero or more than three parts,
// or one whose required session context is unset, yields Found == false.
func (c *Catalog) ResolveRelation(parts []string) ResolvedRelation {
	catalogName, schemaName, objectName, ok := c.resolveContainer(parts)
	if !ok {
		return ResolvedRelation{}
	}
	db := c.GetCatalog(catalogName)
	if db == nil {
		return ResolvedRelation{}
	}
	sch := db.GetSchema(schemaName)
	if sch == nil {
		return ResolvedRelation{}
	}
	if t := sch.GetTable(objectName); t != nil {
		return ResolvedRelation{Catalog: db.Name, Schema: sch.Name, Table: t, Found: true}
	}
	if v := sch.GetView(objectName); v != nil {
		return ResolvedRelation{Catalog: db.Name, Schema: sch.Name, View: v, Found: true}
	}
	return ResolvedRelation{}
}

// ResolveQualifiedName resolves a parsed qualified name, normalizing each
// component via the AST node's quote-aware rules before applying the same
// resolution as ResolveRelation. A nil name yields Found == false.
func (c *Catalog) ResolveQualifiedName(qn *ast.QualifiedName) ResolvedRelation {
	if qn == nil {
		return ResolvedRelation{}
	}
	return c.ResolveRelation(qn.NormalizedParts())
}

// resolveContainer maps the (normalized) dotted name parts to a fully-qualified
// (catalog, schema, object) triple using the session context, returning
// ok == false when the part count is unsupported or the required session
// context is missing.
func (c *Catalog) resolveContainer(parts []string) (catalogName, schemaName, objectName string, ok bool) {
	switch len(parts) {
	case 1:
		if c.currentCatalog == "" || c.currentSchema == "" {
			return "", "", "", false
		}
		return c.currentCatalog, c.currentSchema, parts[0], true
	case 2:
		if c.currentCatalog == "" {
			return "", "", "", false
		}
		return c.currentCatalog, parts[0], parts[1], true
	case 3:
		return parts[0], parts[1], parts[2], true
	default:
		return "", "", "", false
	}
}
