package catalog

import (
	"sort"
	"strings"
)

// diffTriggers compares the triggers of two catalogs, keyed by (database, name).
// Two triggers are equal when their timing, event, table and body match; the body
// is compared as written, after trimming surrounding whitespace.
//
// Not compared: the definer, an environment's account and not schema; the
// FOLLOWS/PRECEDES order; and the charset and collation of the session that
// created the trigger.
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
		// A change is a drop and a create; there is no ALTER TRIGGER.
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

// buildTriggerMap indexes every trigger of a catalog by (database, name),
// lower-cased.
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

// triggersChanged reports whether two same-identity triggers differ in Timing,
// Event, Table, Body or a known sql_mode. Any difference is a drop and a create.
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
	if from.HasSessionContext && to.HasSessionContext && !sameSQLMode(from.SQLMode, to.SQLMode) {
		return true
	}
	return canonicalTriggerBody(from.Body, n) != canonicalTriggerBody(to.Body, n)
}

// canonicalTriggerBody trims the surrounding whitespace of a trigger body, which
// is otherwise compared as written.
func canonicalTriggerBody(body string, _ *Normalizer) string {
	return strings.TrimSpace(body)
}
