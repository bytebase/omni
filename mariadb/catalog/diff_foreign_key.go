package catalog

import (
	"sort"
	"strconv"
	"strings"
)

// diffForeignKeys compares the FOREIGN KEY constraints of two versions of a
// table. Identity is the lower-cased constraint name; the catalog names an unnamed
// key when it loads it. Two same-name keys are equal when their columns,
// referenced table and columns, and ON DELETE and ON UPDATE actions match
// (canonicalForeignKey).
//
// The index that backs a foreign key is not reported here; the index differ leaves
// it out (fkImplicitIndexNames).
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

// foreignKeyMap indexes a table's FOREIGN KEY constraints by lower-cased name.
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

// foreignKeysChanged reports whether two same-name foreign keys differ by
// canonicalForeignKey.
func foreignKeysChanged(a, b *Constraint, n *Normalizer) bool {
	return canonicalForeignKey(a, n) != canonicalForeignKey(b, n)
}

// canonicalForeignKey is the comparison key of a foreign key: its referencing
// columns, referenced database, table and columns, all lower-cased and in order
// because a composite key's column order matters, and its ON DELETE and ON UPDATE
// actions. The name is the identity key, not part of this key.
func canonicalForeignKey(con *Constraint, _ *Normalizer) string {
	return encodeKeyFields(
		"cols", lowerJoin(con.Columns),
		"refdb", canonicalRefDatabase(con),
		"reftbl", toLower(con.RefTable),
		"refcols", lowerJoin(con.RefColumns),
		"ondelete", canonicalFKAction(con.OnDelete),
		"onupdate", canonicalFKAction(con.OnUpdate),
	)
}

// canonicalRefDatabase returns the lower-cased referenced database, or "" when it
// is the key's own database, so a key that names its own database and one that
// does not compare equal.
func canonicalRefDatabase(con *Constraint) string {
	refDB := toLower(con.RefDatabase)
	if refDB == "" {
		return ""
	}
	if con.Table != nil && con.Table.Database != nil && refDB == toLower(con.Table.Database.Name) {
		return ""
	}
	return refDB
}

// canonicalFKAction folds an absent clause onto RESTRICT, which is what MariaDB
// stores for it and leaves out of SHOW CREATE TABLE. NO ACTION is stored and
// written apart, though InnoDB enforces it alike. Measured on mariadb:11.8.8.
func canonicalFKAction(action string) string {
	switch a := strings.ToUpper(strings.TrimSpace(action)); a {
	case "":
		return "RESTRICT"
	default:
		return a
	}
}

// lowerJoin lower-cases each element and encodes the list in order, each element
// length-prefixed so that a name containing the delimiter cannot make two
// different lists collide.
func lowerJoin(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	var b strings.Builder
	for _, p := range parts {
		lower := toLower(p)
		b.WriteString(strconv.Itoa(len(lower)))
		b.WriteByte(':')
		b.WriteString(lower)
		b.WriteByte('|')
	}
	return b.String()
}
