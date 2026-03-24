package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// grantKey is the identity key for a grant: (objType, resolved object name, grantee, privilege, sorted columns).
type grantKey struct {
	objType   byte
	objName   string
	grantee   string
	privilege string
	columns   string // sorted, comma-joined
}

// diffGrants compares grants between two catalogs and returns diff entries.
func diffGrants(from, to *Catalog) []GrantDiffEntry {
	fromMap := buildGrantMap(from)
	toMap := buildGrantMap(to)

	var entries []GrantDiffEntry

	// Iterate from: key not in to → DiffDrop (revoked)
	for k, fg := range fromMap {
		if _, ok := toMap[k]; !ok {
			entries = append(entries, GrantDiffEntry{
				Action: DiffDrop,
				From:   fg,
			})
		}
	}

	// Iterate to: key not in from → DiffAdd (granted)
	for k, tg := range toMap {
		if _, ok := fromMap[k]; !ok {
			entries = append(entries, GrantDiffEntry{
				Action: DiffAdd,
				To:     tg,
			})
		}
	}

	// Both exist → compare WithGrant
	for k, fg := range fromMap {
		tg, ok := toMap[k]
		if !ok {
			continue
		}
		if fg.WithGrant != tg.WithGrant {
			entries = append(entries, GrantDiffEntry{
				Action: DiffModify,
				From:   fg,
				To:     tg,
			})
		}
	}

	// Sort for determinism: by objType, objName, grantee, privilege, columns, then action.
	sort.Slice(entries, func(i, j int) bool {
		ki := grantSortKey(entries[i])
		kj := grantSortKey(entries[j])
		return ki < kj
	})

	return entries
}

// grantSortKey returns a string for sorting grant diff entries deterministically.
func grantSortKey(e GrantDiffEntry) string {
	g := e.To
	if e.Action == DiffDrop {
		g = e.From
	}
	cols := make([]string, len(g.Columns))
	copy(cols, g.Columns)
	sort.Strings(cols)
	return fmt.Sprintf("%c|%s|%s|%s|%d", g.ObjType, g.Grantee, g.Privilege, strings.Join(cols, ","), e.Action)
}

// buildGrantMap builds a map from identity key to Grant for all grants in a catalog.
func buildGrantMap(c *Catalog) map[grantKey]Grant {
	m := make(map[grantKey]Grant, len(c.grants))
	for _, g := range c.grants {
		k := makeGrantKey(c, g)
		m[k] = g
	}
	return m
}

// makeGrantKey resolves the OID to a name and builds the identity key.
func makeGrantKey(c *Catalog, g Grant) grantKey {
	objName := resolveGrantObjName(c, g.ObjType, g.ObjOID)
	cols := make([]string, len(g.Columns))
	copy(cols, g.Columns)
	sort.Strings(cols)
	return grantKey{
		objType:   g.ObjType,
		objName:   objName,
		grantee:   g.Grantee,
		privilege: g.Privilege,
		columns:   strings.Join(cols, ","),
	}
}

// resolveGrantObjName resolves a grant's ObjOID to a stable name string.
func resolveGrantObjName(c *Catalog, objType byte, objOID uint32) string {
	switch objType {
	case 'r':
		rel := c.GetRelationByOID(objOID)
		if rel != nil && rel.Schema != nil {
			return rel.Schema.Name + "." + rel.Name
		}
	case 's':
		seq := c.GetSequenceByOID(objOID)
		if seq != nil && seq.Schema != nil {
			return seq.Schema.Name + "." + seq.Name
		}
	case 'f':
		up := c.userProcs[objOID]
		if up != nil {
			return funcIdentity(c, up)
		}
	case 'n':
		s := c.schemas[objOID]
		if s != nil {
			return s.Name
		}
	}
	return fmt.Sprintf("oid:%d", objOID)
}
