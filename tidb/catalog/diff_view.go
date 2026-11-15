package catalog

import "strings"

// diffViews compares the views of two catalogs by what SHOW CREATE VIEW
// writes, less the definer, which is an environment's account rather than
// schema.
func diffViews(from, to *Catalog) []ViewDiffEntry {
	return diffObjects(viewsOf(from), viewsOf(to),
		func(a, b *View) bool { return canonicalView(a) != canonicalView(b) },
		func(action DiffAction, key objectKey, f, t *View) ViewDiffEntry {
			e := ViewDiffEntry{Action: action, Database: key.db, From: f, To: t}
			if t != nil {
				e.Name = t.Name
			} else {
				e.Name = f.Name
			}
			return e
		},
		func(e ViewDiffEntry) (string, string, DiffAction) { return e.Database, toLower(e.Name), e.Action })
}

func viewsOf(c *Catalog) map[objectKey]*View {
	m := make(map[objectKey]*View)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		for name, v := range db.Views {
			if v != nil {
				m[objectKey{db: toLower(db.Name), name: name}] = v
			}
		}
	}
	return m
}

func canonicalView(v *View) string {
	definer := v.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}
	return strings.Replace(showCreateView(v), " DEFINER="+formatDefiner(definer), "", 1)
}

// diffSequences compares the sequences of two catalogs by what SHOW CREATE
// SEQUENCE writes.
func diffSequences(from, to *Catalog) []SequenceDiffEntry {
	return diffObjects(sequencesOf(from), sequencesOf(to),
		func(a, b *Sequence) bool { return showCreateSequence(a) != showCreateSequence(b) },
		func(action DiffAction, key objectKey, f, t *Sequence) SequenceDiffEntry {
			e := SequenceDiffEntry{Action: action, Database: key.db, From: f, To: t}
			if t != nil {
				e.Name = t.Name
			} else {
				e.Name = f.Name
			}
			return e
		},
		func(e SequenceDiffEntry) (string, string, DiffAction) { return e.Database, toLower(e.Name), e.Action })
}

func sequencesOf(c *Catalog) map[objectKey]*Sequence {
	m := make(map[objectKey]*Sequence)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		for name, s := range db.Sequences {
			if s != nil {
				m[objectKey{db: toLower(db.Name), name: name}] = s
			}
		}
	}
	return m
}
