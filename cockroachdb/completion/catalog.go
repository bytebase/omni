package completion

import (
	"slices"

	metadatapb "github.com/bytebase/omni/metadata"
)

// catalog is the synced metadata of one database, as completion reads it.
type catalog struct {
	database   string
	searchPath []string
	schemas    []*schema
}

type schema struct {
	name      string
	relations []*relation
	functions []string
	types     []string // user-defined enum types
}

type relation struct {
	schema  string
	name    string
	kind    CandidateType
	columns []column
	indexes []string
}

type column struct {
	name   string
	hidden bool
}

// newCatalog indexes meta for completion. Unqualified names resolve through
// defaultSchema, then public, as CockroachDB's default search_path does.
func newCatalog(meta *metadatapb.DatabaseSchemaMetadata, defaultSchema string) *catalog {
	cat := &catalog{database: meta.GetName(), searchPath: []string{"public"}}
	if defaultSchema != "" && defaultSchema != "public" {
		cat.searchPath = []string{defaultSchema, "public"}
	}
	for _, s := range meta.GetSchemas() {
		sc := &schema{name: s.GetName()}
		for _, t := range s.GetTables() {
			rel := &relation{name: t.GetName(), kind: CandidateTable, columns: columnsOf(t.GetColumns())}
			for _, idx := range t.GetIndexes() {
				rel.indexes = append(rel.indexes, idx.GetName())
			}
			sc.relations = append(sc.relations, rel)
		}
		for _, t := range s.GetExternalTables() {
			sc.relations = append(sc.relations, &relation{name: t.GetName(), kind: CandidateTable, columns: columnsOf(t.GetColumns())})
		}
		for _, v := range s.GetViews() {
			sc.relations = append(sc.relations, &relation{name: v.GetName(), kind: CandidateView, columns: columnsOf(v.GetColumns())})
		}
		for _, mv := range s.GetMaterializedViews() {
			rel := &relation{name: mv.GetName(), kind: CandidateMaterializedView, columns: columnsOf(mv.GetColumns())}
			for _, idx := range mv.GetIndexes() {
				rel.indexes = append(rel.indexes, idx.GetName())
			}
			sc.relations = append(sc.relations, rel)
		}
		for _, seq := range s.GetSequences() {
			sc.relations = append(sc.relations, &relation{name: seq.GetName(), kind: CandidateSequence})
		}
		for _, fn := range s.GetFunctions() {
			sc.functions = append(sc.functions, fn.GetName())
		}
		for _, enum := range s.GetEnumTypes() {
			sc.types = append(sc.types, enum.GetName())
		}
		for _, rel := range sc.relations {
			rel.schema = sc.name
		}
		cat.schemas = append(cat.schemas, sc)
	}
	// CockroachDB's sync leaves a materialized view's columns empty; its
	// definition's select list names them.
	for _, s := range meta.GetSchemas() {
		for _, mv := range s.GetMaterializedViews() {
			rel := cat.relation([]string{s.GetName()}, mv.GetName())
			if rel == nil || len(rel.columns) > 0 {
				continue
			}
			for _, name := range definitionColumns(mv.GetDefinition(), cat) {
				rel.columns = append(rel.columns, column{name: name})
			}
		}
	}
	return cat
}

func columnsOf(cols []*metadatapb.ColumnMetadata) []column {
	out := make([]column, 0, len(cols))
	for _, c := range cols {
		out = append(out, column{name: c.GetName(), hidden: c.GetIsInvisible()})
	}
	return out
}

// schema returns the schema named name. Names compare exactly: the parser
// has already lowercased an unquoted name and kept a quoted one as written.
func (c *catalog) schema(name string) *schema {
	if c == nil {
		return nil
	}
	for _, s := range c.schemas {
		if s.name == name {
			return s
		}
	}
	return nil
}

// schemaOf resolves the qualifier of an object name as CockroachDB does:
// [schema] names a schema of the current database, or else the public
// schema of a database of that name; [database, schema] names a schema of
// that database. Only the current database's metadata is known.
func (c *catalog) schemaOf(qual []string) *schema {
	if c == nil {
		return nil
	}
	switch len(qual) {
	case 1:
		if s := c.schema(qual[0]); s != nil {
			return s
		}
		if qual[0] == c.database {
			return c.schema("public")
		}
	case 2:
		if qual[0] == c.database {
			return c.schema(qual[1])
		}
	}
	return nil
}

// relation resolves name qualified by qual, or, unqualified, in pg_catalog
// and then along the search path.
func (c *catalog) relation(qual []string, name string) *relation {
	if c == nil {
		return nil
	}
	if len(qual) > 0 {
		return c.schemaOf(qual).relation(name)
	}
	if slices.Contains(builtinRelations, name) {
		// pg_catalog comes before the search path, and the metadata does
		// not describe its relations.
		return nil
	}
	for _, sp := range c.searchPath {
		if rel := c.schema(sp).relation(name); rel != nil {
			return rel
		}
	}
	return nil
}

func (s *schema) relation(name string) *relation {
	if s == nil {
		return nil
	}
	for _, r := range s.relations {
		if r.name == name {
			return r
		}
	}
	return nil
}

// visibleColumns returns the columns SELECT * expands to.
func (r *relation) visibleColumns() []string {
	var out []string
	for _, col := range r.columns {
		if !col.hidden {
			out = append(out, col.name)
		}
	}
	return out
}
