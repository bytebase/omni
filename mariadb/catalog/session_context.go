package catalog

import "strings"

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

// renderWithSessionContext saves sql_mode, character_set_client, collation_connection and,
// when includeTimeZone is set, time_zone, sets them from ctx, runs stmtSQL, and restores the
// saved values in reverse order. Empty charset, collation, and time_zone fields are not framed.
func renderWithSessionContext(stmtSQL string, ctx SessionContext, includeTimeZone bool) string {
	var pre, post []string

	pre = append(pre, "SET @saved_sql_mode = @@sql_mode")
	pre = append(pre, "SET sql_mode = "+quoteStringLiteral(ctx.SQLMode))
	post = append(post, "SET sql_mode = @saved_sql_mode")

	if cs := strings.TrimSpace(ctx.CharacterSetClient); cs != "" {
		pre = append(pre, "SET @saved_cs_client = @@character_set_client")
		pre = append(pre, "SET character_set_client = "+cs)
		post = append(post, "SET character_set_client = @saved_cs_client")
	}
	if coll := strings.TrimSpace(ctx.CollationConnection); coll != "" {
		pre = append(pre, "SET @saved_coll_conn = @@collation_connection")
		pre = append(pre, "SET collation_connection = "+coll)
		post = append(post, "SET collation_connection = @saved_coll_conn")
	}
	if includeTimeZone {
		if tz := strings.TrimSpace(ctx.TimeZone); tz != "" {
			pre = append(pre, "SET @saved_time_zone = @@time_zone")
			pre = append(pre, "SET time_zone = "+quoteStringLiteral(tz))
			post = append(post, "SET time_zone = @saved_time_zone")
		}
	}

	// Restore in reverse order of the saves.
	reverse(post)

	parts := make([]string, 0, len(pre)+1+len(post))
	parts = append(parts, pre...)
	parts = append(parts, stmtSQL)
	parts = append(parts, post...)
	return strings.Join(parts, ";\n")
}

// reverse reverses a string slice in place.
func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// explicitTimestampDefaults runs a plan that creates, adds or changes a
// TIMESTAMP column with explicit_defaults_for_timestamp ON, restoring the
// session's setting after it. Every column is written with the default it
// has, and OFF, which 10.7 has by default, would give the first NOT NULL
// TIMESTAMP DEFAULT current_timestamp() ON UPDATE current_timestamp().
func explicitTimestampDefaults(ops []MigrationOp, diff *SchemaDiff) []MigrationOp {
	if !touchesTimestamps(diff) {
		return ops
	}
	set := MigrationOp{Type: OpAlterTable, SQL: "SET @saved_explicit_defaults_for_timestamp = @@explicit_defaults_for_timestamp, explicit_defaults_for_timestamp = 1", Phase: PhasePre}
	restore := MigrationOp{Type: OpAlterTable, SQL: "SET explicit_defaults_for_timestamp = @saved_explicit_defaults_for_timestamp", Phase: PhasePost}
	return append(append([]MigrationOp{set}, ops...), restore)
}

func touchesTimestamps(diff *SchemaDiff) bool {
	isTimestamp := func(c *Column) bool { return c != nil && toLower(c.DataType) == "timestamp" }
	for _, e := range diff.Tables {
		switch e.Action {
		case DiffAdd:
			for _, c := range e.To.Columns {
				if isTimestamp(c) {
					return true
				}
			}
		case DiffModify:
			for _, ce := range e.Columns {
				if ce.Action != DiffDrop && isTimestamp(ce.To) {
					return true
				}
			}
		}
	}
	return false
}
