package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// diffFunctions compares user-defined functions and procedures between two catalogs.
func diffFunctions(from, to *Catalog) []FunctionDiffEntry {
	// Build identity-based maps.
	// Identity: schema + name + formatted input arg types, e.g. "public.myfunc(integer,text)"
	fromMap := buildFuncIdentityMap(from)
	toMap := buildFuncIdentityMap(to)

	var entries []FunctionDiffEntry

	// Iterate from: key not in to → DiffDrop
	for id, fp := range fromMap {
		if _, ok := toMap[id]; !ok {
			entries = append(entries, FunctionDiffEntry{
				Action:     DiffDrop,
				SchemaName: from.schemas[fp.Schema.OID].Name,
				Name:       fp.Name,
				Identity:   id,
				From:       fp,
			})
		}
	}

	// Iterate to: key not in from → DiffAdd
	for id, tp := range toMap {
		if _, ok := fromMap[id]; !ok {
			entries = append(entries, FunctionDiffEntry{
				Action:     DiffAdd,
				SchemaName: to.schemas[tp.Schema.OID].Name,
				Name:       tp.Name,
				Identity:   id,
				To:         tp,
			})
		}
	}

	// Both exist → compare fields
	for id, fp := range fromMap {
		tp, ok := toMap[id]
		if !ok {
			continue
		}
		if funcChanged(from, to, fp, tp) {
			entries = append(entries, FunctionDiffEntry{
				Action:     DiffModify,
				SchemaName: from.schemas[fp.Schema.OID].Name,
				Name:       fp.Name,
				Identity:   id,
				From:       fp,
				To:         tp,
			})
		}
	}

	// Sort by identity for determinism.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Identity < entries[j].Identity
	})

	return entries
}

// buildFuncIdentityMap builds a map from identity string to UserProc.
func buildFuncIdentityMap(c *Catalog) map[string]*UserProc {
	m := make(map[string]*UserProc, len(c.userProcs))
	for _, up := range c.userProcs {
		id := funcIdentity(c, up)
		m[id] = up
	}
	return m
}

// funcIdentity returns the identity string for a function/procedure:
// "schemaname.funcname(type1,type2)"
func funcIdentity(c *Catalog, up *UserProc) string {
	schemaName := c.schemas[up.Schema.OID].Name

	var argStrs []string
	for _, aoid := range up.ArgTypes {
		argStrs = append(argStrs, c.FormatType(aoid, -1))
	}

	return fmt.Sprintf("%s.%s(%s)", schemaName, up.Name, strings.Join(argStrs, ","))
}

// funcChanged returns true if any tracked attribute differs between two UserProcs.
func funcChanged(from, to *Catalog, fp, tp *UserProc) bool {
	if fp.Body != tp.Body {
		return true
	}
	if fp.Volatile != tp.Volatile {
		return true
	}
	if fp.Parallel != tp.Parallel {
		return true
	}
	if fp.IsStrict != tp.IsStrict {
		return true
	}
	if fp.SecDef != tp.SecDef {
		return true
	}
	if fp.LeakProof != tp.LeakProof {
		return true
	}
	if fp.Language != tp.Language {
		return true
	}
	if from.FormatType(fp.RetType, -1) != to.FormatType(tp.RetType, -1) {
		return true
	}
	if fp.RetSet != tp.RetSet {
		return true
	}
	return false
}
