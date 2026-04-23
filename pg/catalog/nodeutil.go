package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// qualifiedName extracts schema and name from a list of String nodes.
// If the list has one element, schema is empty. If two, first is schema.
// Used for CreateEnumStmt.TypeName, DropStmt.Objects, CreateDomainStmt.Domainname, etc.
func qualifiedName(names *nodes.List) (schema, name string) {
	if names == nil {
		return "", ""
	}
	items := names.Items
	switch len(items) {
	case 1:
		return "", stringVal(items[0])
	case 2:
		return stringVal(items[0]), stringVal(items[1])
	case 3:
		// catalog.schema.name — ignore catalog
		return stringVal(items[1]), stringVal(items[2])
	default:
		return "", ""
	}
}

// stringVal extracts the string value from a *nodes.String node.
func stringVal(n nodes.Node) string {
	if s, ok := n.(*nodes.String); ok {
		return s.Str
	}
	return ""
}

// intVal extracts the integer value from a *nodes.Integer node.
func intVal(n nodes.Node) int64 {
	if i, ok := n.(*nodes.Integer); ok {
		return i.Ival
	}
	return 0
}

// deparseAConst converts an A_Const node to its string representation.
func deparseAConst(n nodes.Node) string {
	switch v := n.(type) {
	case *nodes.A_Const:
		if v.Isnull {
			return "NULL"
		}
		return deparseAConst(v.Val)
	case *nodes.Integer:
		return fmt.Sprintf("%d", v.Ival)
	case *nodes.Float:
		return v.Fval
	case *nodes.String:
		return v.Str
	case *nodes.Boolean:
		if v.Boolval {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", n)
	}
}

// deparseDatumList converts a list of datum nodes to strings.
func deparseDatumList(list *nodes.List) []string {
	if list == nil {
		return nil
	}
	result := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, deparseAConst(item))
	}
	return result
}

// stringListItems extracts all string values from a List of String nodes.
func stringListItems(list *nodes.List) []string {
	if list == nil {
		return nil
	}
	result := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, stringVal(item))
	}
	return result
}

// extractRawTypmods extracts all integer typmod values from a pgparser Typmods list.
// (pgddl helper — PG extracts these in each type's typmodin function)
func extractRawTypmods(typmods *nodes.List) []int32 {
	if typmods == nil || len(typmods.Items) == 0 {
		return nil
	}
	var mods []int32
	for _, item := range typmods.Items {
		switch v := item.(type) {
		case *nodes.A_Const:
			if v.Val != nil {
				mods = append(mods, int32(intVal(v.Val)))
			}
		case *nodes.Integer:
			mods = append(mods, int32(v.Ival))
		}
	}
	return mods
}

// encodeTypModByName applies PG's type-specific typmod encoding.
// pg: each type's typmodin function (e.g. varchartypmodin, numerictypmodin)
func encodeTypModByName(typName string, rawMods []int32) int32 {
	if len(rawMods) == 0 {
		return -1
	}
	switch typName {
	case "varchar", "bpchar":
		return rawMods[0] + 4 // VARHDRSZ
	case "numeric":
		scale := int32(0)
		if len(rawMods) > 1 {
			scale = rawMods[1]
		}
		return ((rawMods[0] << 16) | scale) + 4 // VARHDRSZ
	default:
		return rawMods[0]
	}
}

// resolveTypeName resolves a pgparser TypeName to an OID and typmod.
//
// (pgddl helper — combines PG's LookupTypeNameOid + typenameTypMod logic)
func (c *Catalog) resolveTypeName(tn *nodes.TypeName) (uint32, int32, error) {
	if tn == nil {
		return 0, -1, fmt.Errorf("NULL type name")
	}

	// Extract schema and type name from Names list.
	// The parser represents types as qualified names: e.g., ["pg_catalog", "int4"].
	schema, name := typeNameParts(tn)

	// Strip pg_catalog prefix since our search path includes it implicitly.
	if schema == "pg_catalog" {
		schema = ""
	}

	name = resolveAlias(name)

	// Compute typmod from Typmods list.
	rawMods := extractRawTypmods(tn.Typmods)
	typmod := int32(-1)
	if len(rawMods) > 0 {
		typmod = rawMods[0] // pass raw first mod for ResolveType validation
	}
	// pgparser sets Typemod=0 (Go zero value) when no modifier is specified,
	// while PG uses -1. Accept both as "no modifier".
	if tn.Typemod > 0 {
		typmod = tn.Typemod
	}

	isArray := tn.ArrayBounds != nil && len(tn.ArrayBounds.Items) > 0

	// Use the existing type resolution logic.
	oid, tm, err := c.ResolveType(TypeName{
		Schema:  schema,
		Name:    name,
		TypeMod: typmod,
		IsArray: isArray,
	})
	if err != nil {
		return 0, -1, err
	}

	// Encode typmod in PG format using type-specific encoding.
	// pg: each type's typmodin function (varchartypmodin, numerictypmodin, etc.)
	if len(rawMods) > 0 {
		tm = encodeTypModByName(name, rawMods)
	}

	return oid, tm, nil
}

// typeNameParts extracts schema and type name from a nodes.TypeName.Names list.
//
// PG accepts qualified type names with up to 3 components in its grammar:
//
//   - 1 component:  name                  → ("", name)
//   - 2 components: schema.name           → (schema, name)
//   - 3 components: catalog.schema.name   → (schema, name)  (catalog is dropped)
//
// For 3-component names, PG validates at name-resolution time that the
// catalog component matches the current database name. If it doesn't,
// PG raises "cross-database references are not implemented". omni has
// no concept of "current database", so we cannot perform that
// validation — we unconditionally drop the catalog component.
//
// Trade-off (deliberate, codex review of commit a47e68d):
//
//   - Permissive (current behavior): `mydb.pg_catalog.int4` resolves
//     correctly when mydb is the current db (the common case for
//     pg_dump output). `otherdb.pg_catalog.int4` is also accepted as
//     local pg_catalog.int4 — a false positive that omni cannot
//     distinguish from the valid case.
//   - Strict (reject all 3-component): would break pg_dump round-trip
//     for any dump that emits 3-part type names with the current db
//     as the catalog (the most common shape).
//
// We chose permissive because the strict alternative blocks a real
// workflow (dump-restore) to catch a syntax form that virtually no
// user writes by hand. Cross-db type references are not implemented
// in PG itself, so the false-positive case represents SQL that
// wouldn't have worked anywhere — only the error message is missing.
// If omni gains a current-database concept, this branch should
// validate and reject mismatches.
//
// 4+ components: PG rejects these at name-resolution time as "improper
// qualified name". omni returns ("", "") so the downstream resolver fails
// to find a type with empty name. We deliberately do NOT fall back to
// `("", lastItem)`: that would silently let `CREATE TABLE t (c
// a.b.c.int4)` resolve as the local `int4` type, turning invalid SQL
// into a successful parse with the wrong AST. A typed error from
// typeNameParts would give a more informative error message than
// `type "" does not exist`, but it would require a signature change at
// all 8 call sites and is out of scope for this fix.
//
// History: prior to this fix, typeNameParts treated len > 2 as
// ("", lastItem) — silently dropping ALL qualification beyond the last
// component. That broke 3-component names in three currently-reachable
// flows: %TYPE references, parseAnyName-driven CREATE TABLE OF / ALTER
// TABLE OF / COMMENT ON CONSTRAINT ON DOMAIN, and (after the paired
// parseGenericType fix) every type position. See
// docs/plans/2026-04-14-pg-followups.md for the full audit.
func typeNameParts(tn *nodes.TypeName) (schema, name string) {
	if tn.Names == nil {
		return "", ""
	}
	items := tn.Names.Items
	switch len(items) {
	case 0:
		return "", ""
	case 1:
		return "", stringVal(items[0])
	case 2:
		return stringVal(items[0]), stringVal(items[1])
	case 3:
		// catalog.schema.name — drop the catalog prefix. PG would
		// validate it matches the current database, but we don't
		// track that here.
		return stringVal(items[1]), stringVal(items[2])
	default:
		// 4+ components: PG rejects these at name resolution time as
		// "improper qualified name". Return empty so the downstream
		// resolver fails. Falling back to ("", lastItem) would silently
		// resolve `CREATE TABLE t (c a.b.c.int4)` as the local `int4`
		// type — a silent success bug that codex caught during the
		// implementation review of this commit's predecessor.
		return "", ""
	}
}

// nodeConstraintType converts a pgparser ConstrType to the catalog ConstraintType.
func nodeConstraintType(ct nodes.ConstrType) (ConstraintType, bool) {
	switch ct {
	case nodes.CONSTR_PRIMARY:
		return ConstraintPK, true
	case nodes.CONSTR_UNIQUE:
		return ConstraintUnique, true
	case nodes.CONSTR_FOREIGN:
		return ConstraintFK, true
	case nodes.CONSTR_CHECK:
		return ConstraintCheck, true
	case nodes.CONSTR_EXCLUSION:
		return ConstraintExclude, true
	default:
		return 0, false
	}
}

// convertConstraintNode converts a pgparser Constraint node to a catalog ConstraintDef.
func convertConstraintNode(con *nodes.Constraint) (ConstraintDef, bool) {
	ct, ok := nodeConstraintType(con.Contype)
	if !ok {
		return ConstraintDef{}, false
	}

	def := ConstraintDef{
		Name: con.Conname,
		Type: ct,
	}

	// Extract deferrable flags (applies to PK, UNIQUE, FK).
	def.Deferrable = con.Deferrable
	def.Deferred = con.Initdeferred
	def.SkipValidation = con.SkipValidation

	switch ct {
	case ConstraintPK, ConstraintUnique:
		def.Columns = stringListItems(con.Keys)
		def.IndexName = con.Indexname // user-specified USING INDEX name
	case ConstraintFK:
		def.Columns = stringListItems(con.FkAttrs)
		def.RefColumns = stringListItems(con.PkAttrs)
		if con.Pktable != nil {
			def.RefSchema = con.Pktable.Schemaname
			def.RefTable = con.Pktable.Relname
		}
		def.FKUpdAction = normalizeFKAction(con.FkUpdaction)
		def.FKDelAction = normalizeFKAction(con.FkDelaction)
		def.FKMatchType = normalizeFKMatch(con.FkMatchtype)
	case ConstraintCheck:
		def.RawCheckExpr = con.RawExpr
		def.CheckExpr = con.CookedExpr
		// Raw check expression text will be filled from AnalyzeStandaloneExpr + DeparseExpr
		// in addCheckConstraint.
	case ConstraintExclude:
		def.AccessMethod = con.AccessMethod
		if def.AccessMethod == "" {
			def.AccessMethod = "gist"
		}
		// Extract exclusion columns and operators from con.Exclusions.
		// pgparser output: list of pair-Lists, each [IndexElem, opList].
		// Manual construction: flat alternating [IndexElem, opList, IndexElem, opList, ...].
		if con.Exclusions != nil {
			items := con.Exclusions.Items
			// Detect format: if first item is a List (not IndexElem), use nested pair format.
			if len(items) > 0 {
				if _, isPair := items[0].(*nodes.List); isPair {
					// Nested pair format (pgparser output).
					for _, item := range items {
						pair, ok := item.(*nodes.List)
						if !ok || len(pair.Items) < 2 {
							continue
						}
						if elem, ok := pair.Items[0].(*nodes.IndexElem); ok && elem.Name != "" {
							def.Columns = append(def.Columns, elem.Name)
						}
						if opList, ok := pair.Items[1].(*nodes.List); ok {
							opName := ""
							for _, n := range opList.Items {
								opName = stringVal(n)
							}
							def.ExclOps = append(def.ExclOps, opName)
						}
					}
				} else {
					// Flat alternating format (manual construction).
					for i := 0; i < len(items)-1; i += 2 {
						if elem, ok := items[i].(*nodes.IndexElem); ok && elem.Name != "" {
							def.Columns = append(def.Columns, elem.Name)
						}
						if opList, ok := items[i+1].(*nodes.List); ok {
							opName := ""
							for _, n := range opList.Items {
								opName = stringVal(n)
							}
							def.ExclOps = append(def.ExclOps, opName)
						}
					}
				}
			}
		}
	}

	return def, true
}

// aliasName extracts the alias name from a pgparser Alias.
func aliasName(a *nodes.Alias) string {
	if a == nil {
		return ""
	}
	return a.Aliasname
}

// normalizeFKAction maps pgparser FK action bytes to PG pg_constraint format.
// PG uses: 'a'=NO ACTION, 'r'=RESTRICT, 'c'=CASCADE, 'n'=SET NULL, 'd'=SET DEFAULT.
// pgparser uses the same chars. Default (zero) maps to 'a' (NO ACTION).
func normalizeFKAction(action byte) byte {
	switch action {
	case 'a', 'r', 'c', 'n', 'd':
		return action
	default:
		return 'a' // NO ACTION
	}
}

// normalizeFKMatch maps pgparser FK match type to PG format.
// 's'=SIMPLE, 'f'=FULL, 'p'=PARTIAL. Default maps to 's'.
func normalizeFKMatch(matchType byte) byte {
	switch matchType {
	case 's', 'f', 'p':
		return matchType
	default:
		return 's' // SIMPLE
	}
}

// convertTypeNameToInternal converts a pgparser TypeName to a catalog TypeName.
func convertTypeNameToInternal(tn *nodes.TypeName) TypeName {
	if tn == nil {
		return TypeName{TypeMod: -1}
	}
	schema, name := typeNameParts(tn)
	if schema == "pg_catalog" {
		schema = ""
	}
	typmod := int32(-1)
	if tn.Typmods != nil && len(tn.Typmods.Items) > 0 {
		if ac, ok := tn.Typmods.Items[0].(*nodes.A_Const); ok {
			if ac.Val != nil {
				typmod = int32(intVal(ac.Val))
			}
		} else if i, ok := tn.Typmods.Items[0].(*nodes.Integer); ok {
			typmod = int32(i.Ival)
		}
	}
	// pgparser sets Typemod=0 (Go zero value) when no modifier is specified,
	// while PG uses -1. Accept both as "no modifier".
	if tn.Typemod > 0 {
		typmod = tn.Typemod
	}
	isArray := tn.ArrayBounds != nil && len(tn.ArrayBounds.Items) > 0
	return TypeName{
		Schema:  schema,
		Name:    resolveAlias(name),
		TypeMod: typmod,
		IsArray: isArray,
	}
}

// defElemString extracts a string value from a DefElem.Arg.
func defElemString(d *nodes.DefElem) string {
	if d.Arg == nil {
		return ""
	}
	switch v := d.Arg.(type) {
	case *nodes.String:
		return v.Str
	case *nodes.TypeName:
		_, n := typeNameParts(v)
		return n
	default:
		return fmt.Sprintf("%v", d.Arg)
	}
}

// defElemInt extracts an int64 value from a DefElem.Arg.
func defElemInt(d *nodes.DefElem) (int64, bool) {
	if d.Arg == nil {
		return 0, false
	}
	switch v := d.Arg.(type) {
	case *nodes.Integer:
		return v.Ival, true
	case *nodes.Float:
		// Float node stores large integers as strings.
		var n int64
		fmt.Sscanf(v.Fval, "%d", &n)
		return n, true
	case *nodes.A_Const:
		switch val := v.Val.(type) {
		case *nodes.Integer:
			return val.Ival, true
		case *nodes.Float:
			var n int64
			fmt.Sscanf(val.Fval, "%d", &n)
			return n, true
		}
	default:
		return 0, false
	}
	return 0, false
}

// defElemBool extracts a boolean value from a DefElem.Arg.
func defElemBool(d *nodes.DefElem) bool {
	if d.Arg == nil {
		// A DefElem with nil Arg typically means TRUE (e.g., "CYCLE" without explicit value).
		return true
	}
	switch v := d.Arg.(type) {
	case *nodes.Boolean:
		return v.Boolval
	case *nodes.Integer:
		return v.Ival != 0
	case *nodes.String:
		return strings.ToLower(v.Str) == "true" || v.Str == "1"
	default:
		return false
	}
}

// isSerialType checks if a type name represents a SERIAL type and returns
// the serial width (2=smallserial, 4=serial, 8=bigserial), or 0 if not serial.
func isSerialType(tn *nodes.TypeName) byte {
	if tn == nil || tn.Names == nil {
		return 0
	}
	_, name := typeNameParts(tn)
	switch strings.ToLower(name) {
	case "smallserial", "serial2":
		return 2
	case "serial", "serial4":
		return 4
	case "bigserial", "serial8":
		return 8
	default:
		return 0
	}
}
