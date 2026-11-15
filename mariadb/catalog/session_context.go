package catalog

// A routine, trigger, or event carries the session context it was created
// under: sql_mode, character_set_client, collation_connection and, for an
// event, time_zone. The sql_mode decides how the body parses, so it is part of
// what the object is. Schema text carries none of it; LoadMetadata stamps it
// from the synced metadata, and ApplySessionContext does so for a catalog
// loaded from text. It is not part of an object's identity.

// SessionContext is the session state an object was created under, as the
// synced metadata records it. An empty field is kept as it is. TimeZone applies
// to events only.
type SessionContext struct {
	SQLMode             string
	CharacterSetClient  string
	CollationConnection string
	TimeZone            string // events only
}

// SessionContextMap holds each object's session context by lower-cased name
// within its kind. A name is unambiguous because a diff loads both sides into
// one database.
type SessionContextMap struct {
	Functions  map[string]SessionContext
	Procedures map[string]SessionContext
	Triggers   map[string]SessionContext
	Events     map[string]SessionContext
}

// ApplySessionContext stamps each object of the catalog that m names, in any
// database, with its context and marks it HasSessionContext, so an empty
// sql_mode is told apart from an unknown one. An object m does not name keeps
// none.
func (c *Catalog) ApplySessionContext(m SessionContextMap) {
	if c == nil {
		return
	}
	for _, db := range c.Databases() {
		applyDatabaseSessionContext(db, m)
	}
}

// applyDatabaseSessionContext stamps the routines, triggers, and events of one
// database with their entries in m.
func applyDatabaseSessionContext(db *Database, m SessionContextMap) {
	applyRoutineContext(db.Functions, m.Functions)
	applyRoutineContext(db.Procedures, m.Procedures)
	for key, tr := range db.Triggers {
		if tr == nil {
			continue
		}
		if ctx, ok := lookupContext(m.Triggers, key, tr.Name); ok {
			tr.SQLMode = ctx.SQLMode
			tr.CharacterSetClient = ctx.CharacterSetClient
			tr.CollationConnection = ctx.CollationConnection
			tr.HasSessionContext = true
		}
	}
	for key, e := range db.Events {
		if e == nil {
			continue
		}
		if ctx, ok := lookupContext(m.Events, key, e.Name); ok {
			e.SQLMode = ctx.SQLMode
			e.CharacterSetClient = ctx.CharacterSetClient
			e.CollationConnection = ctx.CollationConnection
			e.TimeZone = ctx.TimeZone
			e.HasSessionContext = true
		}
	}
}

// applyRoutineContext stamps context onto a function or procedure map (both are
// map[string]*Routine keyed by lower-cased name).
func applyRoutineContext(routines map[string]*Routine, ctxByName map[string]SessionContext) {
	for key, r := range routines {
		if r == nil {
			continue
		}
		if ctx, ok := lookupContext(ctxByName, key, r.Name); ok {
			r.SQLMode = ctx.SQLMode
			r.CharacterSetClient = ctx.CharacterSetClient
			r.CollationConnection = ctx.CollationConnection
			r.HasSessionContext = true
		}
	}
}

// lookupContext finds an object's context by the catalog map key (already lower-cased) and
// falls back to the lower-cased object name, so callers may key the input map either way.
func lookupContext(m map[string]SessionContext, key, name string) (SessionContext, bool) {
	if m == nil {
		return SessionContext{}, false
	}
	if ctx, ok := m[key]; ok {
		return ctx, true
	}
	if ctx, ok := m[toLower(name)]; ok {
		return ctx, true
	}
	return SessionContext{}, false
}
