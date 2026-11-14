package catalog

import "fmt"

// ResolveType resolves a TypeName to an OID and typmod.
func (c *Catalog) ResolveType(tn TypeName) (oid uint32, typmod int32, err error) {
	return c.lookupType(tn.Schema, resolveAlias(tn.Name), tn.Name, tn.IsArray, tn.TypeMod)
}

// lookupType finds the type named name, in schema or on the search path,
// as written: no alias is applied and an explicit schema is kept. shown is
// the name messages give.
func (c *Catalog) lookupType(schema, name, shown string, isArray bool, typmod int32) (uint32, int32, error) {
	tn := TypeName{Schema: schema, Name: shown, IsArray: isArray, TypeMod: typmod}
	var typ *BuiltinType

	if tn.Schema != "" {
		// Schema-qualified: look up schema, then type in that namespace.
		s := c.schemaByName[tn.Schema]
		if s == nil {
			return 0, -1, errUndefinedSchema(tn.Schema)
		}
		typ = c.typeByName[typeKey{ns: s.OID, name: name}]
	} else {
		// Unqualified: search path + implicit pg_catalog.
		for _, nsOID := range c.searchPathWithCatalog() {
			if t := c.typeByName[typeKey{ns: nsOID, name: name}]; t != nil {
				typ = t
				break
			}
		}
	}

	if typ == nil {
		return 0, -1, errUndefinedType(tn.Name)
	}

	// Handle array types.
	if tn.IsArray {
		if typ.Array == 0 {
			return 0, -1, errUndefinedType(tn.Name + "[]")
		}
		arrTyp := c.typeByOID[typ.Array]
		if arrTyp == nil {
			return 0, -1, errUndefinedType(tn.Name + "[]")
		}
		typ = arrTyp
	}

	// Validate typmod.
	typmod = tn.TypeMod
	if typmod != -1 && typ.ModIn == 0 {
		// pg: src/backend/parser/parse_type.c — typenameTypeMod
		return 0, -1, &Error{Code: CodeSyntaxError, Message: fmt.Sprintf("type modifier is not allowed for type %q", typ.TypeName)}
	}

	return typ.OID, typmod, nil
}
