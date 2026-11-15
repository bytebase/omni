package catalog

import "sort"

// SequenceDiffEntry is one sequence that differs between two catalogs.
type SequenceDiffEntry struct {
	Action   DiffAction
	Database string
	Name     string
	From     *Sequence // nil for DiffAdd
	To       *Sequence // nil for DiffDrop
}

// diffSequences compares two catalogs' sequences by database and name, and a
// sequence both hold by what SHOW CREATE SEQUENCE writes for it.
func diffSequences(from, to *Catalog) []SequenceDiffEntry {
	fromMap, toMap := sequencesOf(from), sequencesOf(to)
	var result []SequenceDiffEntry
	for key, f := range fromMap {
		if _, ok := toMap[key]; !ok {
			result = append(result, SequenceDiffEntry{Action: DiffDrop, Database: key.db, Name: f.Name, From: f})
		}
	}
	for key, t := range toMap {
		f, ok := fromMap[key]
		switch {
		case !ok:
			result = append(result, SequenceDiffEntry{Action: DiffAdd, Database: key.db, Name: t.Name, To: t})
		case showCreateSequence(f) != showCreateSequence(t):
			result = append(result, SequenceDiffEntry{Action: DiffModify, Database: key.db, Name: t.Name, From: f, To: t})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		a, b := result[i], result[j]
		if a.Database != b.Database {
			return a.Database < b.Database
		}
		if x, y := toLower(a.Name), toLower(b.Name); x != y {
			return x < y
		}
		return a.Action < b.Action
	})
	return result
}

func sequencesOf(c *Catalog) map[eventKey]*Sequence {
	out := make(map[eventKey]*Sequence)
	if c == nil {
		return out
	}
	for _, db := range c.Databases() {
		for name, seq := range db.Sequences {
			out[eventKey{db: toLower(db.Name), name: name}] = seq
		}
	}
	return out
}
