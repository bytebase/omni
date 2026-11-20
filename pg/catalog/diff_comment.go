package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// commentNormKey is a name-based key for matching comments across catalogs
// (where OIDs differ). Storing schema, object, and sub-object separately avoids
// ambiguity when identifiers contain periods.
type commentNormKey struct {
	ObjType     byte
	Description string // human/migration-op description of the object
	SchemaName  string
	ObjectName  string
	SubObjName  string
	SubID       int16
}

func (k commentNormKey) parts() []string {
	switch {
	case k.SubObjName != "":
		return []string{k.SchemaName, k.ObjectName, k.SubObjName}
	case k.SchemaName != "" && k.ObjectName != "":
		return []string{k.SchemaName, k.ObjectName}
	case k.ObjectName != "":
		return []string{k.ObjectName}
	case k.SchemaName != "":
		return []string{k.SchemaName}
	default:
		return nil
	}
}

// diffComments compares comments between two catalogs and returns diff entries.
func diffComments(from, to *Catalog) []CommentDiffEntry {
	fromMap := buildCommentNormMap(from, to, true)
	toMap := buildCommentNormMap(to, from, false)

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
				objParts:       nk.parts(),
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
				objParts:       nk.parts(),
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
				objParts:       nk.parts(),
			})
		}
	}

	// Sort for determinism: by ObjType, then Description, then parts, then SubID.
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
func buildCommentNormMap(c, other *Catalog, isFrom bool) map[commentNormKey]string {
	m := make(map[commentNormKey]string, len(c.comments))
	for ck, text := range c.comments {
		if base, ok := c.extensionComments[ck]; !c.diffsMemberSetting(ck.ObjOID, other, isFrom, ok && base == text) {
			continue
		}
		nk, ok := resolveCommentNormKey(c, ck)
		if !ok {
			continue // object not found — skip
		}
		m[nk] = text
	}
	return m
}

// resolveCommentNormKey converts an OID-based comment key into a structured
// name-based key for cross-catalog comparison.
func resolveCommentNormKey(c *Catalog, ck commentKey) (commentNormKey, bool) {
	base := commentNormKey{ObjType: ck.ObjType, SubID: ck.SubID}
	switch ck.ObjType {
	case 'r': // relation (table, view, matview, foreign table)
		rel := c.relationByOID[ck.ObjOID]
		if rel == nil || rel.Schema == nil {
			return commentNormKey{}, false
		}
		base.SchemaName = rel.Schema.Name
		base.ObjectName = rel.Name
		base.Description = fmt.Sprintf("%s.%s", rel.Schema.Name, rel.Name)
		return base, true

	case 'i': // index
		idx := c.indexes[ck.ObjOID]
		if idx == nil {
			return commentNormKey{}, false
		}
		for _, s := range c.schemas {
			for name, si := range s.Indexes {
				if si.OID == ck.ObjOID {
					base.SchemaName = s.Name
					base.ObjectName = name
					base.Description = fmt.Sprintf("%s.%s", s.Name, name)
					return base, true
				}
			}
		}
		return commentNormKey{}, false

	case 'f': // function/procedure
		up := c.userProcs[ck.ObjOID]
		if up == nil {
			return commentNormKey{}, false
		}
		desc := funcIdentity(c, up)
		base.Description = desc
		if up.Schema != nil {
			base.SchemaName = up.Schema.Name
			base.ObjectName = strings.TrimPrefix(desc, up.Schema.Name+".")
		}
		return base, desc != ""

	case 'n': // schema
		for _, s := range c.schemas {
			if s.OID == ck.ObjOID {
				base.SchemaName = s.Name
				base.Description = s.Name
				return base, true
			}
		}
		return commentNormKey{}, false

	case 't': // type (enum/domain/range)
		bt := c.typeByOID[ck.ObjOID]
		if bt == nil {
			return commentNormKey{}, false
		}
		ns := c.schemas[bt.Namespace]
		if ns == nil {
			return commentNormKey{}, false
		}
		base.SchemaName = ns.Name
		base.ObjectName = bt.TypeName
		base.Description = fmt.Sprintf("%s.%s", ns.Name, bt.TypeName)
		return base, true

	case 's': // sequence
		seq := c.sequenceByOID[ck.ObjOID]
		if seq == nil {
			return commentNormKey{}, false
		}
		for _, s := range c.schemas {
			for name, ss := range s.Sequences {
				if ss.OID == ck.ObjOID {
					base.SchemaName = s.Name
					base.ObjectName = name
					base.Description = fmt.Sprintf("%s.%s", s.Name, name)
					return base, true
				}
			}
		}
		return commentNormKey{}, false

	case 'c': // constraint
		con := c.constraints[ck.ObjOID]
		if con == nil {
			return commentNormKey{}, false
		}
		rel := c.relationByOID[con.RelOID]
		if rel == nil || rel.Schema == nil {
			base.Description = con.Name
			return base, true
		}
		base.SchemaName = rel.Schema.Name
		base.ObjectName = rel.Name
		base.SubObjName = con.Name
		base.Description = fmt.Sprintf("%s.%s.%s", rel.Schema.Name, rel.Name, con.Name)
		return base, true

	case 'g': // trigger
		trig := c.triggers[ck.ObjOID]
		if trig == nil {
			return commentNormKey{}, false
		}
		rel := c.relationByOID[trig.RelOID]
		if rel == nil || rel.Schema == nil {
			base.Description = trig.Name
			return base, true
		}
		base.SchemaName = rel.Schema.Name
		base.ObjectName = rel.Name
		base.SubObjName = trig.Name
		base.Description = fmt.Sprintf("%s.%s.%s", rel.Schema.Name, rel.Name, trig.Name)
		return base, true

	case 'E': // event trigger
		evt := c.eventTriggers[ck.ObjOID]
		if evt == nil {
			return commentNormKey{}, false
		}
		base.ObjectName = evt.Name
		base.Description = evt.Name
		return base, true

	case 'd': // domain constraint
		for _, dt := range c.domainTypes {
			for _, dc := range dt.Constraints {
				if dc.OID == ck.ObjOID {
					bt := c.typeByOID[dt.TypeOID]
					if bt == nil {
						base.Description = dc.Name
						return base, true
					}
					ns := c.schemas[bt.Namespace]
					if ns == nil {
						base.Description = dc.Name
						return base, true
					}
					base.SchemaName = ns.Name
					base.ObjectName = bt.TypeName
					base.SubObjName = dc.Name
					base.Description = fmt.Sprintf("%s.%s.%s", ns.Name, bt.TypeName, dc.Name)
					return base, true
				}
			}
		}
		return commentNormKey{}, false

	case 'p': // policy
		for _, p := range c.policies {
			if p.OID == ck.ObjOID {
				rel := c.relationByOID[p.RelOID]
				if rel == nil || rel.Schema == nil {
					base.Description = p.Name
					return base, true
				}
				base.SchemaName = rel.Schema.Name
				base.ObjectName = rel.Name
				base.SubObjName = p.Name
				base.Description = fmt.Sprintf("%s.%s.%s", rel.Schema.Name, rel.Name, p.Name)
				return base, true
			}
		}
		return commentNormKey{}, false

	default:
		return commentNormKey{}, false
	}
}
