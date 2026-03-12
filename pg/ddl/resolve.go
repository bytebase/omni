package ddl

// ResolveType resolves a TypeName to an OID and typmod.
func (c *Catalog) ResolveType(tn TypeName) (oid uint32, typmod int32, err error) {
	name := resolveAlias(tn.Name)

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
		return 0, -1, errInvalidParameterValue("type " + typ.TypeName + " does not support type modifiers")
	}

	return typ.OID, typmod, nil
}
