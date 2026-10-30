package analysis

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Name lookups.
//
// The resolver compares relation names, column names, and USING keys with
// strings.EqualFold, and where several entries match it takes the first one
// in list order. Scanning the list for every lookup made a select list, FROM
// list, USING list, or BY NAME arm of n entries cost O(n^2). The indexes here
// key each list on foldKey and give the same first match. An index over a
// list that grows while lookups are made (the leaf relations of a FROM
// clause) or that may be read only partly (a relation's projection) is
// extended on demand, so it never costs more than the scans it replaces.

// foldKey maps two names to the same key exactly when strings.EqualFold
// reports them equal: that compares rune by rune, matching runes of one
// simple case-folding orbit, so each rune becomes the least of its orbit.
func foldKey(name string) string {
	ascii := true
	for i := 0; i < len(name); i++ {
		if name[i] >= utf8.RuneSelf {
			ascii = false
			break
		}
	}
	if ascii {
		return strings.ToUpper(name)
	}
	var b strings.Builder
	for _, r := range name {
		least := r
		for f := unicode.SimpleFold(r); f != r; f = unicode.SimpleFold(f) {
			if f < least {
				least = f
			}
		}
		b.WriteRune(least)
	}
	return b.String()
}

// foldSet is a set of names compared under strings.EqualFold. The nil set is
// empty.
type foldSet map[string]struct{}

func newFoldSet(names []string) foldSet {
	if len(names) == 0 {
		return nil
	}
	s := make(foldSet, len(names))
	for _, n := range names {
		s[foldKey(n)] = struct{}{}
	}
	return s
}

func (s foldSet) has(name string) bool {
	_, ok := s[foldKey(name)]
	return ok
}

// leafSet is the leaf FROM relations of the current SELECT, in FROM order
// (see spanWalker.leaves), with the indexes behind the column resolver's
// lookups. Relations are only ever appended, and a relation's projection is
// never changed once built, so each index covers a prefix of rels (or of the
// columns) and is extended when a lookup misses.
type leafSet struct {
	rels []*relation

	// byName maps foldKey(rel.name) to the first relation of rels[:named]
	// with that name.
	byName map[string]*relation
	named  int

	// bases holds the folded table parts of each base relation in
	// rels[:based].
	bases map[baseKey]struct{}
	based int

	// baseCount counts the base relations; firstBase is the first one.
	baseCount int
	firstBase *relation

	// columns maps a folded column name to the first column with that name
	// among the non-star columns of the derived (non-base) relations, in FROM
	// order, up to the cursor (colRel, colPos).
	columns        map[string]*projColumn
	colRel, colPos int
}

// baseKey is a base table reference with each part folded by foldKey.
type baseKey struct {
	catalog, database, schema, table string
}

func newBaseKey(ref ColumnRef) baseKey {
	return baseKey{foldKey(ref.Catalog), foldKey(ref.Database), foldKey(ref.Schema), foldKey(ref.Table)}
}

func (s *leafSet) add(rel *relation) {
	s.rels = append(s.rels, rel)
	if rel.isBase {
		if s.baseCount == 0 {
			s.firstBase = rel
		}
		s.baseCount++
	}
}

// find returns the first relation whose reference name matches name
// (case-insensitive, the GoogleSQL identifier rule), or nil.
func (s *leafSet) find(name string) *relation {
	key := foldKey(name)
	if rel, ok := s.byName[key]; ok {
		return rel
	}
	for s.named < len(s.rels) {
		rel := s.rels[s.named]
		s.named++
		if rel.name == "" {
			continue
		}
		k := foldKey(rel.name)
		if _, ok := s.byName[k]; ok {
			continue
		}
		if s.byName == nil {
			s.byName = make(map[string]*relation)
		}
		s.byName[k] = rel
		if k == key {
			return rel
		}
	}
	return nil
}

// hasBase reports whether some base relation's table matches ref's
// Catalog/Database/Schema/Table, each compared case-insensitively.
func (s *leafSet) hasBase(ref ColumnRef) bool {
	for ; s.based < len(s.rels); s.based++ {
		if rel := s.rels[s.based]; rel.isBase {
			if s.bases == nil {
				s.bases = make(map[baseKey]struct{})
			}
			s.bases[newBaseKey(rel.baseRef)] = struct{}{}
		}
	}
	_, ok := s.bases[newBaseKey(ref)]
	return ok
}

// soleBase returns the single base relation when exactly one exists, else
// nil. With several base relations a bare column is ambiguous (omni cannot
// enumerate base-table columns metadata-free), so the caller keeps the bare
// ref and the consumer matches additively by name.
func (s *leafSet) soleBase() *relation {
	if s.baseCount == 1 {
		return s.firstBase
	}
	return nil
}

// derivedColumn returns the sources of the first non-star column named
// column (case-insensitive) of the derived relations, in FROM order: what
// projectionColumnSources without base stars finds scanning each non-base
// relation in turn.
func (s *leafSet) derivedColumn(column string) ([]ColumnRef, bool) {
	key := foldKey(column)
	if c, ok := s.columns[key]; ok {
		return c.sources, true
	}
	for ; s.colRel < len(s.rels); s.colRel, s.colPos = s.colRel+1, 0 {
		rel := s.rels[s.colRel]
		if rel.isBase {
			continue
		}
		for s.colPos < len(rel.columns) {
			c := &rel.columns[s.colPos]
			s.colPos++
			if c.baseStar != nil {
				continue
			}
			k := foldKey(c.name)
			if _, ok := s.columns[k]; ok {
				continue
			}
			if s.columns == nil {
				s.columns = make(map[string]*projColumn)
			}
			s.columns[k] = c
			if k == key {
				return c.sources, true
			}
		}
	}
	return nil, false
}

// relIndex indexes a relation's projection for projectionColumnSources over
// columns[:scanned].
type relIndex struct {
	first   map[string]int  // folded name → first non-star column with that name
	stars   []int           // base-star columns, in order
	excepts map[int]foldSet // a base-star column → its except list, built on first use
	scanned int
}

// projectionColumnSources returns the base lineage of the column named `column`
// in a relation's projection, or nil if the relation has no such named column.
//
// A concrete projection column matches by name. A baseStar element (a whole base
// table whose columns omni cannot enumerate) matches ANY column name — but ONLY
// when allowBaseStar is set: a baseStar match is sound for a QUALIFIED reference
// (`rel.col` names the relation, so `col` must be one of its columns) but NOT for
// an UNQUALIFIED reference (omni cannot know which of several base relations — or
// join arms — owns a bare `col`; that is left to the metadata-aware consumer to
// match by name). The first column that matches wins.
func projectionColumnSources(rel *relation, column string, allowBaseStar bool) ([]ColumnRef, bool) {
	if rel.index == nil {
		rel.index = &relIndex{}
	}
	x := rel.index
	key := foldKey(column)
	match := -1
	if i, ok := x.first[key]; ok {
		match = i
	}
	if allowBaseStar {
		for _, p := range x.stars {
			if match >= 0 && p > match {
				break
			}
			if !x.starExcepts(rel, p, key) {
				match = p
				break
			}
		}
	}
	for match < 0 && x.scanned < len(rel.columns) {
		p := x.scanned
		x.scanned++
		c := &rel.columns[p]
		if c.baseStar != nil {
			x.stars = append(x.stars, p)
			if allowBaseStar && !x.starExcepts(rel, p, key) {
				match = p
			}
			continue
		}
		k := foldKey(c.name)
		if _, ok := x.first[k]; !ok {
			if x.first == nil {
				x.first = make(map[string]int)
			}
			x.first[k] = p
		}
		if k == key {
			match = p
		}
	}
	if match < 0 {
		return nil, false
	}
	c := rel.columns[match]
	if c.baseStar != nil {
		ref := *c.baseStar
		ref.Column = column
		return []ColumnRef{ref}, true
	}
	// Found — even when its lineage is empty (e.g. a CTE column projected from
	// a literal like `1 AS n`); return found=true so the caller does NOT apply
	// the unresolved-reference fallback (which would wrongly fabricate a bare
	// column ref). Empty-but-resolved lineage is the legacy result.
	return c.sources, true
}

// starExcepts reports whether the base-star column p excludes the column
// whose folded name is key: a USING-coalesced star excludes its key columns,
// which are owned by the coalesced concrete column earlier in the projection.
func (x *relIndex) starExcepts(rel *relation, p int, key string) bool {
	except := rel.columns[p].baseStarExcept
	if len(except) == 0 {
		return false
	}
	set, ok := x.excepts[p]
	if !ok {
		set = newFoldSet(except)
		if x.excepts == nil {
			x.excepts = make(map[int]foldSet)
		}
		x.excepts[p] = set
	}
	_, found := set[key]
	return found
}
