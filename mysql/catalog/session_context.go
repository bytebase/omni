package catalog

import "strings"

// MySQL SDL — object session context (sql_mode / charset / collation / time_zone).
//
// A routine, trigger, or event carries the session context it was created under:
// sql_mode, character_set_client, collation_connection, and — for events — the session
// time_zone. A declarative rollout must not silently drop that context when it re-emits the
// object (e.g. a routine authored under PIPES_AS_CONCAT re-created under the server default
// treats `||` as OR). Two emission paths must preserve it, both verified on live 8.0.32:
//
//   - DROP+CREATE recreate (routines/triggers always for a body/param change — MySQL has no
//     CREATE OR REPLACE; events for the ENDS-removal case). The fresh CREATE adopts the
//     SESSION context, so it must run under the OLD object's context.
//   - ALTER EVENT … DO <body> (the common event body/schedule change). Empirically, ALTER
//     EVENT RE-STAMPS sql_mode, character_set_client, and collation_connection from the
//     session (it does NOT re-stamp time_zone), so a bare ALTER under the deploy session
//     rewrites the event's stored modes. It too must run under the OLD event's context.
//
// NOT re-stamped, so left bare: ALTER FUNCTION / ALTER PROCEDURE (the characteristic-only
// routine modify — COMMENT / SQL SECURITY / DATA ACCESS). Verified on live 8.0.32: an
// ALTER FUNCTION COMMENT under a default-mode session leaves ROUTINES.SQL_MODE unchanged.
//
// The declarative SDL text is BARE (a clean export carries no `SET sql_mode` framing), so
// omni cannot recover the old context from the source SDL. Instead the context enters
// out-of-band as STRUCTURED data: bytebase passes the synced metadata's per-object context
// (RoutineMetadata.SqlMode, TriggerMetadata.SqlMode, EventMetadata.TimeZone, …) via
// ApplySessionContext, which stamps it onto the already-loaded catalog's objects. The SDL
// text stays 100% bare; the context lives only on the in-memory objects.
//
// The context is deliberately NOT part of any object's declarative identity — the diff keys
// (canonicalRoutine, triggersChanged, eventCanonicalKey) all exclude it, so a mode-only
// difference never triggers a phantom recreate. It is consumed only by the generators, which
// wrap the emitted CREATE/ALTER in a concat-safe save/restore sourced from the OLD (From)
// object (renderWithSessionContext), matching mysqldump's framing, so one object's context
// can never leak into the next statement of a multi-statement migration.

// SessionContext is the session state captured when a routine/trigger/event was created.
// Every field is the value as stored in the synced metadata; empty strings are legal and
// are preserved verbatim (an authored empty sql_mode round-trips as SET sql_mode=”). A
// PARTIAL context (some axes empty) preserves the populated axes and leaves the session
// default for the rest — synced information_schema metadata always carries sql_mode /
// charset / collation for a real object, so this only arises for a hand-built context.
// TimeZone applies only to events and is ignored for routines and triggers.
//
// The database default collation is deliberately NOT modeled: it is a property of the SCHEMA
// (CREATE DATABASE … DEFAULT COLLATE), not a session variable, so it is not preservable
// per-object. Verified on live 8.0.32 — setting `SET collation_database = <x>` before a
// CREATE leaves information_schema.ROUTINES.DATABASE_COLLATION at the schema default, not
// <x>. A DEFAULT-COLLATE change is its own database-level schema diff.
type SessionContext struct {
	SQLMode             string
	CharacterSetClient  string
	CollationConnection string
	TimeZone            string // events only
}

// SessionContextMap carries the per-object session context for a whole schema, keyed by
// LOWER-CASED object name within each object kind (matching the catalog's identity
// folding). It is the structured out-of-band input bytebase hands to ApplySessionContext;
// the SDL text itself never carries session context. A nil or missing entry simply leaves
// that object's context unset (a bare recreate).
//
// Name-only keying is sufficient because the MySQL declarative-rollout path is
// SINGLE-DATABASE: bytebase loads both the source (dumped) and target (user) SDL under one
// synthetic database (bbcatalog), so within a diff a name is unambiguous. ApplySessionContext
// applies the map across every database in the catalog; a hypothetical multi-database catalog
// with same-named objects in different databases is out of scope for this path.
type SessionContextMap struct {
	Functions  map[string]SessionContext
	Procedures map[string]SessionContext
	Triggers   map[string]SessionContext
	Events     map[string]SessionContext
}

// ApplySessionContext stamps the supplied per-object session context onto every matching
// object already loaded in the catalog, across all databases. It is called on the SOURCE
// (from/current) catalog — the one whose objects supply the ORIGINAL context a recreate
// must restore — after LoadSDL/LoadSQL and before Diff.
//
// Matching is by lower-cased object name within each kind; the MySQL release path is
// single-database (see SessionContextMap), so the name alone is the identity. An object
// present in the catalog but absent from the map is left with no context (HasSessionContext
// stays false → bare recreate); a map entry with no matching object is ignored.
//
// Marking HasSessionContext lets the generator distinguish "context applied, sql_mode empty"
// (emit SET sql_mode=”) from "no context" (bare CREATE), so an authored empty sql_mode
// round-trips.
func (c *Catalog) ApplySessionContext(m SessionContextMap) {
	if c == nil {
		return
	}
	for _, db := range c.Databases() {
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

// renderWithSessionContext wraps a single CREATE or ALTER statement in a concat-safe
// save/restore of the object's original session context, mirroring mysqldump's framing:
//
//	SET @saved_sql_mode = @@sql_mode; SET sql_mode = '<mode>';
//	SET @saved_cs_client = @@character_set_client; SET character_set_client = <cs>;
//	SET @saved_coll_conn = @@collation_connection; SET collation_connection = <coll>;
//	SET @saved_time_zone = @@time_zone; SET time_zone = '<tz>';   -- events only
//	<CREATE|ALTER …>;
//	SET time_zone = @saved_time_zone;
//	SET collation_connection = @saved_coll_conn;
//	SET character_set_client = @saved_cs_client;
//	SET sql_mode = @saved_sql_mode;
//
// The restores run in reverse order and re-capture each live @@-value immediately before
// overwriting it, so applying a multi-object migration never leaks one object's context into
// the next statement even though the @saved_* names are fixed (each block fully saves →
// sets → runs → restores before the next block runs; blocks never nest). Each clause is
// joined with the same ";\n" separator MigrationPlan.SQL() uses to join ops, so the whole
// framed block is a well-formed statement sequence inside one MigrationOp.SQL.
//
// sql_mode is emitted UNCONDITIONALLY when framing (even when empty) so an authored empty
// sql_mode (”) forces empty modes rather than inheriting the server default — the framing
// is only ever invoked with a known-present context. character_set_client,
// collation_connection, and time_zone are emitted only when non-empty: their values are
// bare identifiers / literals and an empty value would produce invalid SQL, and a synced
// object without one of them simply keeps the session default for that axis.
//
// includeTimeZone gates the time_zone axis (events set it true; routines and triggers false
// — time_zone has no effect on a routine or trigger body's stored form).
func renderWithSessionContext(stmtSQL string, ctx SessionContext, includeTimeZone bool) string {
	var pre, post []string

	// sql_mode: always framed.
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

	// Restores run in reverse of the saves so nested contexts unwind cleanly.
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
