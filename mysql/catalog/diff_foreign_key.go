package catalog

import (
	"sort"
	"strings"
)

// MySQL SDL diff — foreign keys.
//
// This is the MySQL analog of PG's foreign-key differ. It compares the FOREIGN KEY
// constraint set of two versions of a table and reports the per-FK changes via
// TableDiffEntry.ForeignKeys. MySQL keeps FK + CHECK on the constraint path (show.go:121);
// PRIMARY KEY / UNIQUE / secondary indexes live on the index path and are owned by the index
// node. So this node owns ONLY Table.Constraints entries of type ConForeignKey.
//
// Identity is the constraint name, lower-cased. MySQL auto-names an unnamed FK
// `<table>_ibfk_<N>` at load time (tablecmds.go), so by the time a catalog reaches the diff
// both the user form and the engine readback carry a concrete name — matching by name is the
// stored-form identity. Equality is decided by a canonical key (canonicalForeignKey) that
// folds the referencing columns, the referenced database+table+columns, and the ON DELETE/ON
// UPDATE actions through the NORMALIZER's version-aware FK-action canonicalization
// (CanonicalFKAction): RESTRICT, NO ACTION, and an absent clause are the same referential
// behavior and must collapse onto one key. This is required because SHOW CREATE echoes the
// "default" action differently per version (verified against both live engines):
//   - 8.0 OMITS NO ACTION (the implicit default) but ECHOES RESTRICT verbatim;
//   - 5.7 OMITS RESTRICT (its implicit default) but ECHOES NO ACTION verbatim.
//
// So a user `ON DELETE RESTRICT` and the engine readback (which on 5.7 drops it, on 8.0 keeps
// it) must produce the same canonical key on both versions, else the FK phantom-diffs forever.
// A MySQL FK cannot be altered in place, so a changed FK is a DROP followed by an ADD — emitted
// in PhasePost by the generator (migration_foreign_key.go).
//
// FK-implicit backing indexes are NOT this differ's concern here — they are diffed (and
// excluded) by the index node (diff_index.go fkImplicitIndexNames). This differ only reports
// the FK constraint change; the generator (migration_foreign_key.go) owns the lifecycle of the
// leftover backing index that MySQL leaves after a DROP FOREIGN KEY.
func diffForeignKeys(from, to *Table, n *Normalizer) []ForeignKeyDiffEntry {
	fromMap := foreignKeyMap(from)
	toMap := foreignKeyMap(to)

	var result []ForeignKeyDiffEntry

	// Dropped: in from but not in to.
	for name, fromFK := range fromMap {
		if _, ok := toMap[name]; !ok {
			result = append(result, ForeignKeyDiffEntry{
				Action: DiffDrop,
				Name:   fromFK.Name,
				From:   fromFK,
			})
		}
	}

	// Added or modified: in to.
	for name, toFK := range toMap {
		fromFK, ok := fromMap[name]
		if !ok {
			result = append(result, ForeignKeyDiffEntry{
				Action: DiffAdd,
				Name:   toFK.Name,
				To:     toFK,
			})
			continue
		}
		if foreignKeysChanged(fromFK, toFK, n) {
			result = append(result, ForeignKeyDiffEntry{
				Action: DiffModify,
				Name:   toFK.Name,
				From:   fromFK,
				To:     toFK,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if a, b := toLower(result[i].Name), toLower(result[j].Name); a != b {
			return a < b
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// foreignKeyMap indexes a table's FOREIGN KEY constraints by lower-cased constraint name.
// Only ConForeignKey constraints are included; PK/UNIQUE (index path) and CHECK (the check
// node) are not this node's concern.
func foreignKeyMap(t *Table) map[string]*Constraint {
	m := map[string]*Constraint{}
	if t == nil {
		return m
	}
	for _, con := range t.Constraints {
		if con == nil || con.Type != ConForeignKey {
			continue
		}
		m[toLower(con.Name)] = con
	}
	return m
}

// foreignKeysChanged reports whether two same-name foreign keys differ, comparing their
// canonical keys. canonicalForeignKey folds every stored-form-significant attribute into one
// collision-free key, so this differ never re-implements a per-attribute comparison.
func foreignKeysChanged(a, b *Constraint, n *Normalizer) bool {
	return canonicalForeignKey(a, n) != canonicalForeignKey(b, n)
}

// canonicalForeignKey returns a single stable comparison key for a foreign key, folding the
// attributes that survive into MySQL's stored form:
//   - the ON DELETE / ON UPDATE actions via canonicalFKAction, which collapses RESTRICT / NO
//     ACTION / absent onto one "default" key (see that function for the per-version echo
//     evidence).
//
// The referencing columns, referenced database/table, and referenced columns are compared
// case-insensitively and ORDER-SENSITIVELY (a composite FK's column order is significant). The
// referenced database is folded to lower case but kept distinct from "" — a same-database FK
// (SHOW CREATE omits the db qualifier → RefDatabase == "") and an explicitly-qualified one are
// the same FK only when both resolve to the same schema; since the synced readback never
// qualifies a same-db FK, "" is the canonical same-db form and is preserved verbatim. The
// constraint NAME is the identity key (handled by the caller's map), not part of this content
// key.
//
// The Normalizer is threaded for symmetry with the other differs and so the action
// canonicalization can move to normalize-core later (it is version-independent for the
// EQUALITY question — all of RESTRICT/NO ACTION/absent fold together regardless of version —
// even though the per-version stored ECHO differs, which only the renderer in show.go cares
// about).
func canonicalForeignKey(con *Constraint, _ *Normalizer) string {
	return encodeKeyFields(
		"cols", lowerJoin(con.Columns),
		"refdb", toLower(con.RefDatabase),
		"reftbl", toLower(con.RefTable),
		"refcols", lowerJoin(con.RefColumns),
		"ondelete", canonicalFKAction(con.OnDelete),
		"onupdate", canonicalFKAction(con.OnUpdate),
	)
}

// canonicalFKAction folds a foreign-key referential action onto a canonical key for the
// EQUALITY comparison. The three forms RESTRICT, NO ACTION, and an absent clause ("") are the
// same referential behavior in MySQL (RESTRICT and NO ACTION are synonyms; no clause defaults
// to NO ACTION), so they collapse to one "default" key. CASCADE, SET NULL, and SET DEFAULT are
// distinct and kept verbatim (upper-cased).
//
// This collapse is REQUIRED for idempotence because SHOW CREATE echoes the default action
// differently per version (verified against both live engines, 5.7 :13307 / 8.0 :13306):
//   - 8.0 OMITS NO ACTION but ECHOES RESTRICT verbatim;
//   - 5.7 OMITS RESTRICT but ECHOES NO ACTION verbatim.
//
// So `ON DELETE RESTRICT` written by the user reloads as "RESTRICT", but its 5.7 readback (which
// drops the clause) reloads as "NO ACTION" (the loader's default) — and on 8.0 the readback
// keeps "RESTRICT". Both must compare equal to the user form on both versions, which only a
// version-INDEPENDENT default-collapse achieves. (The renderer's per-version echo is show.go's
// concern, not the equality key's.)
//
// FLAG: this is FK-action canonicalization and arguably belongs in normalize-core alongside
// CanonicalColumn / CanonicalComment. It is kept local to the FK node (no other node needs it)
// pending a normalize-core CanonicalFKAction; if one is added, route through it. This mirrors
// the index node keeping canonicalIndexKeyBlockSize local with the same FLAG.
func canonicalFKAction(action string) string {
	switch strings.ToUpper(strings.TrimSpace(action)) {
	case "", "RESTRICT", "NO ACTION":
		return "DEFAULT"
	default:
		return strings.ToUpper(strings.TrimSpace(action))
	}
}

// lowerJoin lower-cases each element and joins with a comma, preserving order (FK column order
// is significant). Used for both the referencing and referenced column lists.
func lowerJoin(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	lowered := make([]string, len(parts))
	for i, p := range parts {
		lowered[i] = toLower(p)
	}
	return strings.Join(lowered, ",")
}
