package catalog

import (
	"sort"
	"strings"
)

// MySQL SDL diff — triggers (the full BEFORE/AFTER × INSERT/UPDATE/DELETE surface).
//
// Triggers are DATABASE-level objects in MySQL (keyed by (database, name), each bound to a
// table via Trigger.Table), unlike PG where a trigger hangs off a relation. So diffTriggers
// mirrors diffTables' whole-catalog shape — build a (database, name) map from each side and
// report Add / Drop / Modify — rather than PG's per-relation diffTriggers. It is wired into
// DiffWithNormalizer (diff.go) and reported via SchemaDiff.Triggers.
//
// IDEMPOTENCE is the point of this node. A trigger has no ALTER form in MySQL, and the synced
// (stored) schema is read back from information_schema.triggers / SHOW CREATE TRIGGER. The
// canonical key the diff compares is the set of fields that survive that round-trip identically:
//
//   - Timing (BEFORE/AFTER)     — the parser upper-cases it on both the user and the readback
//     side, so it compares directly.
//   - Event (INSERT/UPDATE/DELETE) — likewise upper-cased by the parser.
//   - Table (the owning table)  — compared case-insensitively (MySQL identifier folding).
//   - Body (the trigger action statement) — MySQL stores the body VERBATIM (byte-for-byte the
//     text the user wrote: whitespace, newlines, case, and any BEGIN…END block are preserved in
//     both information_schema.triggers.ACTION_STATEMENT and SHOW CREATE TRIGGER). Verified on
//     live 5.7.25 + 8.0.32. So body equality is an exact compare after trimming surrounding
//     whitespace (the loader already TrimSpaces it in createTrigger).
//
// DELIBERATELY NOT COMPARED (each would phantom-diff a no-op trigger against its own readback —
// see the normalization flags in the PR body):
//
//   - Definer — the user form usually omits DEFINER (the loader defaults it to `root`@`%`),
//     while the readback always carries the server's DEFINER; the two spellings also differ
//     (the SHOW-CREATE backtick form vs the parser's bare form). DEFINER is an ownership/security
//     attribute, not schema shape, and is conventionally ignored in a declarative schema diff.
//   - Order (FOLLOWS/PRECEDES) — the relative ordering of multiple triggers on the same
//     table+timing+event. SHOW CREATE TRIGGER does NOT echo FOLLOWS/PRECEDES, and the loader
//     stores triggers in an unordered map without modelling the resolved information_schema
//     ACTION_ORDER as a dependency. A trigger declared with `FOLLOWS x` therefore loads with a
//     non-nil Order that its own readback (no FOLLOWS) lacks. Comparing Order would phantom-diff
//     every ordered trigger. Trigger ordering is FLAGGED as not-diffed (see the PR flag
//     trigger-order-not-modelled); the omission has no observable schema divergence for the
//     overwhelmingly common single-trigger-per-event case, and the loader limitation is the
//     proper home for a fix.
//   - CharacterSetClient / CollationConnection / DatabaseCollation — session/connection charset
//     metadata captured at create time, not part of the trigger's declarative shape; the user
//     SDL form carries the loader's session defaults, the readback carries the creating
//     session's. Ignored, same rationale as DEFINER.
//
// implemented by omni:triggers breadth node
func diffTriggers(from, to *Catalog, n *Normalizer) []TriggerDiffEntry {
	fromMap := buildTriggerMap(from)
	toMap := buildTriggerMap(to)

	var result []TriggerDiffEntry

	// Dropped: in from but not in to.
	for key, fromTrig := range fromMap {
		if _, ok := toMap[key]; !ok {
			result = append(result, TriggerDiffEntry{
				Action:   DiffDrop,
				Database: key.db,
				Name:     key.name,
				From:     fromTrig,
			})
		}
	}

	// Added or modified: in to.
	for key, toTrig := range toMap {
		fromTrig, ok := fromMap[key]
		if !ok {
			result = append(result, TriggerDiffEntry{
				Action:   DiffAdd,
				Database: key.db,
				Name:     key.name,
				To:       toTrig,
			})
			continue
		}
		// Both exist — a change is rendered as DROP + CREATE (no ALTER TRIGGER in MySQL).
		if triggersChanged(fromTrig, toTrig, n) {
			result = append(result, TriggerDiffEntry{
				Action:   DiffModify,
				Database: key.db,
				Name:     key.name,
				From:     fromTrig,
				To:       toTrig,
			})
		}
	}

	// Determinism: sort by (database, name, action) so the diff — and the DDL generate-core
	// derives from it — is stable across runs.
	sort.Slice(result, func(i, j int) bool {
		if result[i].Database != result[j].Database {
			return result[i].Database < result[j].Database
		}
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// buildTriggerMap indexes every trigger in a catalog by (database, name), both lower-cased.
// MySQL has no schemas/namespaces, so the database is the qualifying scope.
func buildTriggerMap(c *Catalog) map[triggerKey]*Trigger {
	m := make(map[triggerKey]*Trigger)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		for name, trig := range db.Triggers {
			if trig == nil {
				continue // defensive; the loader never stores nil
			}
			m[triggerKey{db: toLower(db.Name), name: name}] = trig
		}
	}
	return m
}

// triggerKey is the identity key for a trigger: (database, name) lower-cased.
type triggerKey struct {
	db   string
	name string
}

// triggersChanged reports whether two same-identity triggers differ in any compared field.
// The compared set is Timing, Event, Table, and canonical Body (see the diffTriggers doc for
// what is deliberately excluded and why). Any difference means DROP + CREATE.
func triggersChanged(from, to *Trigger, n *Normalizer) bool {
	if !strings.EqualFold(from.Timing, to.Timing) {
		return true
	}
	if !strings.EqualFold(from.Event, to.Event) {
		return true
	}
	if !strings.EqualFold(from.Table, to.Table) {
		return true
	}
	return canonicalTriggerBody(from.Body, n) != canonicalTriggerBody(to.Body, n)
}

// canonicalTriggerBody reduces a trigger action statement to its canonical comparison form.
// MySQL preserves the body verbatim on storage (verified on live 5.7.25 + 8.0.32), so the
// user form and the engine readback are already byte-identical for the body; the only
// canonicalization needed is trimming surrounding whitespace (which the loader's createTrigger
// already applies, so this is belt-and-suspenders for any caller that builds a Trigger directly).
// The Normalizer is threaded for signature symmetry with the other Canonical* helpers and to
// leave a single seam should a future version diverge in stored-body form.
func canonicalTriggerBody(body string, _ *Normalizer) string {
	return strings.TrimSpace(body)
}
