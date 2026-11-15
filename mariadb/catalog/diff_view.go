package catalog

import (
	"sort"
	"strings"
)

// diffViews compares the views of two catalogs, keyed by (database, name). Two
// views are equal when SHOW CREATE VIEW gives the same text apart from the
// definer, which is an environment's account and not schema (see canonicalView).
func diffViews(from, to *Catalog, n *Normalizer) []ViewDiffEntry {
	fromMap := buildViewMap(from)
	toMap := buildViewMap(to)

	var result []ViewDiffEntry

	// Dropped: in from but not in to.
	for key, fromView := range fromMap {
		if _, ok := toMap[key]; !ok {
			result = append(result, ViewDiffEntry{
				Action:   DiffDrop,
				Database: key.db,
				Name:     key.name,
				From:     fromView,
			})
		}
	}

	// Added or modified: in to.
	for key, toView := range toMap {
		fromView, ok := fromMap[key]
		if !ok {
			result = append(result, ViewDiffEntry{
				Action:   DiffAdd,
				Database: key.db,
				Name:     key.name,
				To:       toView,
			})
			continue
		}
		if viewsChanged(fromView, toView, n) {
			result = append(result, ViewDiffEntry{
				Action:   DiffModify,
				Database: key.db,
				Name:     key.name,
				From:     fromView,
				To:       toView,
			})
		}
	}

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

// buildViewMap indexes every view of a catalog by (database, name), lower-cased.
func buildViewMap(c *Catalog) map[tableKey]*View {
	m := make(map[tableKey]*View)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		for name, v := range db.Views {
			if v == nil {
				continue
			}
			m[tableKey{db: toLower(db.Name), name: name}] = v
		}
	}
	return m
}

// viewsChanged reports whether two same-identity views differ by canonicalView.
func viewsChanged(a, b *View, n *Normalizer) bool {
	return canonicalView(a, n) != canonicalView(b, n)
}

// canonicalView is the comparison key of a view: SHOW CREATE VIEW without the
// definer. The catalog stores references into the view's own database without the
// database, so a body that writes them and one that does not give the same key.
func canonicalView(v *View, _ *Normalizer) string {
	// SHOW CREATE VIEW, less the definer, which is an environment's account
	// rather than schema. MariaDB folds a column list into the body's aliases,
	// so a view written with one and its synced form compare equal.
	definer := "DEFINER=" + formatDefiner(v.Definer) + " "
	return strings.Replace(showCreateView(v), definer, "", 1)
}
