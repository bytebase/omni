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
// time a name is qualified by it. A nil catalog has no objects.
type catalog struct {
	ctx       context.Context
	load      func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
	list      func(ctx context.Context) ([]string, error)
	current   string
	databases map[string]*metadatapb.DatabaseSchemaMetadata // nil value: not found
	loads     int
}

func newCatalog(ctx context.Context, cfg SessionConfig) *catalog {
	if cfg.LoadDatabase == nil && cfg.ListDatabaseNames == nil {
		return nil
	}
	return &catalog{
		ctx:       ctx,
		load:      cfg.LoadDatabase,
		list:      cfg.ListDatabaseNames,
		current:   cfg.DefaultDatabase,
		databases: map[string]*metadatapb.DatabaseSchemaMetadata{},
	}
}

func (c *catalog) currentDatabase() *metadatapb.DatabaseSchemaMetadata {
	if c == nil || c.current == "" {
		return nil
	}
	return c.database(c.current)
}

// database returns the database named name, loading it on first use.
func (c *catalog) database(name string) *metadatapb.DatabaseSchemaMetadata {
	if c == nil || c.load == nil || name == "" {
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
	db, err := c.load(c.ctx, name)
	if err != nil {
		db = nil
	}
	c.databases[name] = db
	return db
}

// databaseNames returns the instance's databases as candidates.
func (c *catalog) databaseNames() []Candidate {
	if c == nil || c.list == nil {
		return nil
	}
	names, err := c.list(c.ctx)
	if err != nil {
		return nil
	}
	out := make([]Candidate, 0, len(names))
	for _, n := range names {
		out = append(out, Candidate{Text: n, Type: CandidateDatabase})
	}
	return out
}

// relations returns the tables, views and materialized views of db.
func (c *catalog) relations(db *metadatapb.DatabaseSchemaMetadata) []Candidate {
	var out []Candidate
	for _, s := range db.GetSchemas() {
		for _, t := range s.GetTables() {
			out = append(out, Candidate{Text: t.GetName(), Type: CandidateTable, Definition: db.GetName() + "." + t.GetName(), Comment: t.GetComment()})
		}
		for _, v := range s.GetViews() {
			out = append(out, Candidate{Text: v.GetName(), Type: CandidateView, Definition: db.GetName() + "." + v.GetName(), Comment: v.GetComment()})
		}
		for _, v := range s.GetMaterializedViews() {
			out = append(out, Candidate{Text: v.GetName(), Type: CandidateMaterializedView, Definition: db.GetName() + "." + v.GetName(), Comment: v.GetComment()})
		}
	}
	return out
}

// columns returns the columns of the relation named by parts: table,
// db.table, or internal.db.table. A name that is one of ctes is the CTE and
// has no columns in the catalog.
func (c *catalog) columns(parts []string, ctes []string) []Candidate {
	var db *metadatapb.DatabaseSchemaMetadata
	switch len(parts) {
	case 1:
		for _, cte := range ctes {
			if cte == parts[0] {
				return nil
			}
		}
		db = c.currentDatabase()
	case 2:
		db = c.database(parts[0])
	case 3:
		if !strings.EqualFold(parts[0], "internal") {
			return nil
		}
		db = c.database(parts[1])
	default:
		return nil
	}
	name := parts[len(parts)-1]
	var out []Candidate
	add := func(cols []*metadatapb.ColumnMetadata) {
		for _, col := range cols {
			out = append(out, Candidate{Text: col.GetName(), Type: CandidateColumn, Definition: name + " | " + col.GetType(), Comment: col.GetComment()})
		}
	}
	for _, s := range db.GetSchemas() {
		for _, t := range s.GetTables() {
			if t.GetName() == name {
				add(t.GetColumns())
			}
		}
		for _, v := range s.GetViews() {
			if v.GetName() == name {
				add(v.GetColumns())
			}
		}
		for _, v := range s.GetMaterializedViews() {
			if v.GetName() == name {
				add(v.GetColumns())
			}
		}
	}
	return out
}
