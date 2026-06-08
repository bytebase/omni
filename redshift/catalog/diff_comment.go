package catalog

import (
	"fmt"
	"sort"
)

// commentNormKey is a name-based key for matching comments across catalogs
// (where OIDs differ).
type commentNormKey struct {
	ObjType     byte
	Description string // name-based description of the object
	SubID       int16
}

// diffComments compares comments between two catalogs and returns diff entries.
func diffComments(from, to *Catalog) []CommentDiffEntry {
	fromMap := buildCommentNormMap(from)
	toMap := buildCommentNormMap(to)

	var entries []CommentDiffEntry

	// from-only → Drop
	for nk, text := range fromMap {
		if _, ok := toMap[nk]; !ok {
			entries = append(entries, CommentDiffEntry{
				Action:         DiffDrop,
				ObjType:        nk.ObjType,
				ObjDescription: nk.Description,
				SubID:          nk.SubID,
				From:           text,
			})
		}
	}

	// to-only → Add
	for nk, text := range toMap {
		if _, ok := fromMap[nk]; !ok {
			entries = append(entries, CommentDiffEntry{
				Action:         DiffAdd,
				ObjType:        nk.ObjType,
				ObjDescription: nk.Description,
				SubID:          nk.SubID,
				To:             text,
			})
		}
	}

	// both → Modify if text differs
	for nk, fromText := range fromMap {
		toText, ok := toMap[nk]
		if !ok {
			continue
		}
		if fromText != toText {
			entries = append(entries, CommentDiffEntry{
				Action:         DiffModify,
				ObjType:        nk.ObjType,
				ObjDescription: nk.Description,
				SubID:          nk.SubID,
				From:           fromText,
				To:             toText,
			})
		}
	}

	// Sort for determinism: by ObjType, then Description, then SubID.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].ObjType != entries[j].ObjType {
			return entries[i].ObjType < entries[j].ObjType
		}
		if entries[i].ObjDescription != entries[j].ObjDescription {
			return entries[i].ObjDescription < entries[j].ObjDescription
		}
		return entries[i].SubID < entries[j].SubID
	})

	return entries
}

// buildCommentNormMap builds a map from normalized (name-based) comment keys
// to comment text for all comments in a catalog.
func buildCommentNormMap(c *Catalog) map[commentNormKey]string {
	m := make(map[commentNormKey]string, len(c.comments))
	for ck, text := range c.comments {
		desc := resolveCommentDescription(c, ck)
		if desc == "" {
			continue // object not found — skip
		}
		nk := commentNormKey{
			ObjType:     ck.ObjType,
			Description: desc,
			SubID:       ck.SubID,
		}
		m[nk] = text
	}
	return m
}

// resolveCommentDescription converts an OID-based comment key into a
// name-based description string for cross-catalog comparison.
func resolveCommentDescription(c *Catalog, ck commentKey) string {
	switch ck.ObjType {
	case 'r': // relation (table, view, matview, foreign table)
		rel := c.relationByOID[ck.ObjOID]
		if rel == nil || rel.Schema == nil {
			return ""
		}
		return fmt.Sprintf("%s.%s", rel.Schema.Name, rel.Name)

	case 'i': // index
		idx := c.indexes[ck.ObjOID]
		if idx == nil {
			return ""
		}
		// Find the schema for this index.
		for _, s := range c.schemas {
			for name, si := range s.Indexes {
				if si.OID == ck.ObjOID {
					return fmt.Sprintf("%s.%s", s.Name, name)
				}
			}
		}
		return ""

	case 'f': // function/procedure
		up := c.userProcs[ck.ObjOID]
		if up == nil {
			return ""
		}
		return funcIdentity(c, up)

	case 'n': // schema
		for _, s := range c.schemas {
			if s.OID == ck.ObjOID {
				return s.Name
			}
		}
		return ""

	case 't': // type (enum/domain)
		bt := c.typeByOID[ck.ObjOID]
		if bt == nil {
			return ""
		}
		ns := c.schemas[bt.Namespace]
		if ns == nil {
			return ""
		}
		return fmt.Sprintf("%s.%s", ns.Name, bt.TypeName)

	case 's': // sequence
		seq := c.sequenceByOID[ck.ObjOID]
		if seq == nil {
			return ""
		}
		// Find the schema for this sequence.
		for _, s := range c.schemas {
			for name, ss := range s.Sequences {
				if ss.OID == ck.ObjOID {
					return fmt.Sprintf("%s.%s", s.Name, name)
				}
			}
		}
		return ""

	case 'c': // constraint
		con := c.constraints[ck.ObjOID]
		if con == nil {
			return ""
		}
		// Build "schema.table.constraint" for uniqueness.
		rel := c.relationByOID[con.RelOID]
		if rel == nil || rel.Schema == nil {
			return con.Name
		}
		return fmt.Sprintf("%s.%s.%s", rel.Schema.Name, rel.Name, con.Name)

	case 'g': // trigger
		trig := c.triggers[ck.ObjOID]
		if trig == nil {
			return ""
		}
		rel := c.relationByOID[trig.RelOID]
		if rel == nil || rel.Schema == nil {
			return trig.Name
		}
		return fmt.Sprintf("%s.%s.%s", rel.Schema.Name, rel.Name, trig.Name)

	case 'E': // event trigger
		evt := c.eventTriggers[ck.ObjOID]
		if evt == nil {
			return ""
		}
		return evt.Name

	case 'd': // domain constraint
		// Domain constraints are stored on DomainType.Constraints.
		for _, dt := range c.domainTypes {
			for _, dc := range dt.Constraints {
				if dc.OID == ck.ObjOID {
					bt := c.typeByOID[dt.TypeOID]
					if bt == nil {
						return dc.Name
					}
					ns := c.schemas[bt.Namespace]
					if ns == nil {
						return dc.Name
					}
					return fmt.Sprintf("%s.%s.%s", ns.Name, bt.TypeName, dc.Name)
				}
			}
		}
		return ""

	case 'p': // policy
		for _, p := range c.policies {
			if p.OID == ck.ObjOID {
				rel := c.relationByOID[p.RelOID]
				if rel == nil || rel.Schema == nil {
					return p.Name
				}
				return fmt.Sprintf("%s.%s.%s", rel.Schema.Name, rel.Name, p.Name)
			}
		}
		return ""

	default:
		return ""
	}
}
