package completion

import (
	"context"
	"slices"

	metadatapb "github.com/bytebase/omni/metadata"
)

// catalog is the synced metadata completion reads: the current database,
// and others loaded when a statement names them.
type catalog struct {
	ctx       context.Context
	cfg       SessionConfig
	current   string
	databases map[string]*database // nil when a load failed
	names     []string
	listed    bool
}

type database struct {
	name      string
	relations []*relation
}

type relation struct {
	name    string
	kind    CandidateType
	comment string
	columns []column
}

type column struct {
	name, typ, comment string
}

func newCatalog(ctx context.Context, cfg SessionConfig) *catalog {
	return &catalog{ctx: ctx, cfg: cfg, current: cfg.DefaultDatabase, databases: map[string]*database{}}
}

// database returns the database named name, the current one when name is
// empty, or nil when it cannot be loaded.
func (c *catalog) database(name string) *database {
	if name == "" {
		name = c.current
	}
	if name == "" || c.cfg.LoadDatabase == nil {
		return nil
	}
	if db, ok := c.databases[name]; ok {
		return db
	}
	meta, err := c.cfg.LoadDatabase(c.ctx, name)
	var db *database
	if err == nil && meta != nil {
		db = newDatabase(name, meta)
	}
	c.databases[name] = db
	return db
}

// newDatabase indexes a database's metadata. ClickHouse has no schemas: its
// sync puts every object in one unnamed schema.
func newDatabase(name string, meta *metadatapb.DatabaseSchemaMetadata) *database {
	db := &database{name: name}
	for _, s := range meta.GetSchemas() {
		for _, t := range s.GetTables() {
			db.relations = append(db.relations, &relation{name: t.GetName(), kind: CandidateTable, comment: t.GetComment(), columns: columnsOf(t.GetColumns())})
		}
		for _, v := range s.GetViews() {
			db.relations = append(db.relations, &relation{name: v.GetName(), kind: CandidateView, comment: v.GetComment(), columns: columnsOf(v.GetColumns())})
		}
		for _, mv := range s.GetMaterializedViews() {
			db.relations = append(db.relations, &relation{name: mv.GetName(), kind: CandidateMaterializedView, comment: mv.GetComment(), columns: columnsOf(mv.GetColumns())})
		}
	}
	return db
}

func columnsOf(cols []*metadatapb.ColumnMetadata) []column {
	out := make([]column, 0, len(cols))
	for _, c := range cols {
		out = append(out, column{name: c.GetName(), typ: c.GetType(), comment: c.GetComment()})
	}
	return out
}

// relation returns the table or view named name in database db, the current
// one when db is empty. ClickHouse names are case-sensitive.
func (c *catalog) relation(db, name string) *relation {
	d := c.database(db)
	if d == nil {
		return nil
	}
	for _, r := range d.relations {
		if r.name == name {
			return r
		}
	}
	return nil
}

// databaseNames lists the instance's databases, or the current one when
// they cannot be listed.
func (c *catalog) databaseNames() []string {
	if !c.listed {
		c.listed = true
		if c.cfg.ListDatabaseNames != nil {
			if names, err := c.cfg.ListDatabaseNames(c.ctx); err == nil {
				c.names = names
			}
		}
		if len(c.names) == 0 && c.current != "" {
			c.names = []string{c.current}
		}
	}
	return slices.Clone(c.names)
}
