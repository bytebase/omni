package completion

import (
	"context"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
)

// maxDatabaseLoads bounds how many databases besides the current one a
// completion loads for names qualified by another database.
const maxDatabaseLoads = 2

// catalog holds the session's databases, loading another database the first
// time a name is qualified by it. Names are compared as Snowflake resolves
// them: exactly, after an unquoted identifier is folded to upper case.
type catalog struct {
	ctx           context.Context
	load          func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
	current       string
	defaultSchema string
	databases     map[string]*database // nil value: not found
	loads         int
}

type database struct {
	name    string
	schemas []*schema
}

type schema struct {
	name      string
	relations []*relation
	functions []routine
}

type relation struct {
	schema, name string
	kind         CandidateType
	comment      string
	columns      []column
}

type column struct {
	name, typ, comment string
}

type routine struct {
	name, definition, comment string
}

func newCatalog(ctx context.Context, cfg SessionConfig) *catalog {
	c := &catalog{
		ctx:           ctx,
		load:          cfg.LoadDatabase,
		current:       cfg.DefaultDatabase,
		defaultSchema: cfg.DefaultSchema,
		databases:     map[string]*database{},
	}
	if c.defaultSchema == "" {
		c.defaultSchema = "PUBLIC"
	}
	if c.load == nil || c.current == "" {
		return nil
	}
	if c.database(c.current) == nil {
		return nil
	}
	return c
}

// database returns the database named name, loading it on first use.
func (c *catalog) database(name string) *database {
	if c == nil {
		return nil
	}
	if db, ok := c.databases[name]; ok {
		return db
	}
	if name != c.current {
		if c.loads >= maxDatabaseLoads {
			return nil
		}
		c.loads++
	}
	meta, err := c.load(c.ctx, name)
	var db *database
	if err == nil && meta != nil {
		db = newDatabase(name, meta)
	}
	c.databases[name] = db
	return db
}

func newDatabase(name string, meta *metadatapb.DatabaseSchemaMetadata) *database {
	db := &database{name: name}
	for _, s := range meta.GetSchemas() {
		sc := &schema{name: s.GetName()}
		add := func(name string, kind CandidateType, comment string, cols []*metadatapb.ColumnMetadata) {
			r := &relation{schema: sc.name, name: name, kind: kind, comment: comment}
			for _, col := range cols {
				r.columns = append(r.columns, column{name: col.GetName(), typ: col.GetType(), comment: col.GetComment()})
			}
			sc.relations = append(sc.relations, r)
		}
		for _, t := range s.GetTables() {
			add(t.GetName(), CandidateTable, t.GetComment(), t.GetColumns())
		}
		for _, t := range s.GetExternalTables() {
			add(t.GetName(), CandidateTable, "", t.GetColumns())
		}
		for _, v := range s.GetViews() {
			add(v.GetName(), CandidateView, v.GetComment(), v.GetColumns())
		}
		for _, v := range s.GetMaterializedViews() {
			add(v.GetName(), CandidateMaterializedView, v.GetComment(), v.GetColumns())
		}
		for _, f := range s.GetFunctions() {
			sc.functions = append(sc.functions, routine{name: f.GetName(), definition: f.GetSignature(), comment: f.GetComment()})
		}
		db.schemas = append(db.schemas, sc)
	}
	return db
}

func (db *database) schema(name string) *schema {
	if db == nil {
		return nil
	}
	for _, s := range db.schemas {
		if s.name == name {
			return s
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

// searchPath returns the schemas of the current database an unqualified
// name is looked up in. A query resolves it through SEARCH_PATH, whose
// default is $current, $public; DDL and DML use the current schema alone
// (https://docs.snowflake.com/en/sql-reference/name-resolution). A session
// that sets another SEARCH_PATH is not seen.
func (c *catalog) searchPath(query bool) []string {
	if query && c.defaultSchema != "PUBLIC" {
		return []string{c.defaultSchema, "PUBLIC"}
	}
	return []string{c.defaultSchema}
}

// qualify returns the database, schema and name a table name of up to three
// resolved parts stands for. A name with no database is in the current one;
// a name with no schema is in the first schema of the search path that has
// it, or the current schema when none does.
func (c *catalog) qualify(parts []string, query bool) []string {
	switch {
	case c == nil || len(parts) == 0 || len(parts) >= 3:
		return parts
	case len(parts) == 2:
		return []string{c.current, parts[0], parts[1]}
	}
	db := c.database(c.current)
	for _, s := range c.searchPath(query) {
		if db.schema(s).relation(parts[0]) != nil {
			return []string{c.current, s, parts[0]}
		}
	}
	return []string{c.current, c.defaultSchema, parts[0]}
}

// relationAt returns the relation a qualified name names.
func (c *catalog) relationAt(full []string) *relation {
	if c == nil || len(full) != 3 {
		return nil
	}
	return c.database(full[0]).schema(full[1]).relation(full[2])
}

// schemaAt returns the schema a qualifier of up to two resolved parts names:
// a schema of the current database, or database.schema.
func (c *catalog) schemaAt(parts []string) *schema {
	if c == nil {
		return nil
	}
	switch len(parts) {
	case 1:
		return c.database(c.current).schema(parts[0])
	case 2:
		return c.database(parts[0]).schema(parts[1])
	}
	return nil
}

// qualifiedName renders a relation's name for a Definition.
func (r *relation) qualifiedName() string {
	return strings.Join([]string{r.schema, r.name}, ".")
}
