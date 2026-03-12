package ddl

import "fmt"

// CoercionPathway describes how a type cast is performed.
type CoercionPathway int

const (
	CoercionNone    CoercionPathway = iota // no coercion possible
	CoercionRelabel                        // binary compatible (no-op cast)
	CoercionFunc                           // cast via function
	CoercionIO                             // cast via I/O conversion
)

// FindCoercionPathway returns the coercion pathway from source to target
// for the given context ('i'=implicit, 'a'=assignment, 'e'=explicit).
// The second return value is the cast function OID (0 for relabel/none).
func (c *Catalog) FindCoercionPathway(source, target uint32, context byte) (CoercionPathway, uint32) {
	if source == target {
		return CoercionRelabel, 0
	}

	cast := c.castIndex[castKey{source: source, target: target}]
	if cast != nil {
		if !contextPermits(cast.Context, context) {
			return CoercionNone, 0
		}

		switch cast.Method {
		case 'b':
			return CoercionRelabel, 0
		case 'f':
			return CoercionFunc, cast.Func
		case 'i':
			return CoercionIO, 0
		}
	}

	// Domain unwrapping: if source or target is a domain, try the base type.
	// pg: src/backend/parser/parse_coerce.c — find_coercion_pathway (domain handling)
	if cast == nil {
		baseSource := c.getBaseType(source)
		baseTarget := c.getBaseType(target)
		if baseSource != source || baseTarget != target {
			p, f := c.FindCoercionPathway(baseSource, baseTarget, context)
			if p != CoercionNone {
				return p, f
			}
		}
	}

	// Array coercion: if both source and target are array types, check element coercion.
	// pg: src/backend/parser/parse_coerce.c — find_coercion_pathway (array coercion)
	if cast == nil {
		sourceType := c.typeByOID[source]
		targetType := c.typeByOID[target]
		if sourceType != nil && targetType != nil &&
			sourceType.Category == 'A' && targetType.Category == 'A' &&
			sourceType.Elem != 0 && targetType.Elem != 0 {
			elemPath, _ := c.FindCoercionPathway(sourceType.Elem, targetType.Elem, context)
			if elemPath != CoercionNone {
				return CoercionFunc, 0 // array coercion via element casting
			}
		}
	}

	// No pg_cast entry found. Try automatic I/O coercion fallback.
	// pg: src/backend/parser/parse_coerce.c — find_coercion_pathway
	//
	// Rules:
	// 1. Assignment or stricter context + target is string category → CoerceViaIO
	// 2. Explicit context + source is string category → CoerceViaIO
	if cast == nil {
		targetType := c.typeByOID[target]
		sourceType := c.typeByOID[source]
		if (context == 'a' || context == 'e') && targetType != nil && targetType.Category == 'S' {
			return CoercionIO, 0
		}
		if context == 'e' && sourceType != nil && sourceType.Category == 'S' {
			return CoercionIO, 0
		}
	}

	return CoercionNone, 0
}

// getBaseType unwraps domain chains to find the underlying base type.
// pg: src/backend/utils/cache/lsyscache.c — getBaseType
func (c *Catalog) getBaseType(typeOID uint32) uint32 {
	for i := 0; i < 32; i++ { // safety limit to prevent infinite loops
		dom := c.domainTypes[typeOID]
		if dom == nil {
			return typeOID
		}
		typeOID = dom.BaseTypeOID
	}
	return typeOID
}

// CanCoerce returns true if source can be coerced to target in the given context.
func (c *Catalog) CanCoerce(source, target uint32, context byte) bool {
	p, _ := c.FindCoercionPathway(source, target, context)
	return p != CoercionNone
}

// IsBinaryCoercible returns true if source can be cast to target as a no-op relabel.
//
// pg: src/backend/parser/parse_coerce.c — IsBinaryCoercibleWithCast
func (c *Catalog) IsBinaryCoercible(source, target uint32) bool {
	if source == target {
		return true
	}

	// Anything is coercible to ANY or ANYELEMENT or ANYCOMPATIBLE.
	if target == ANYOID || target == ANYELEMENTOID || target == ANYCOMPATIBLEOID {
		return true
	}

	// Domain unwrap: reduce source to its base type.
	// pg: IsBinaryCoercibleWithCast — getBaseType(srctype)
	if bt := c.typeByOID[source]; bt != nil && bt.Type == 'd' && bt.BaseType != 0 {
		source = c.getBaseType(source)
		if source == target {
			return true
		}
	}

	// Array type coercible to ANY[COMPATIBLE]ARRAY.
	if target == ANYARRAYOID || target == ANYCOMPATIBLEARRAYOID {
		if c.typeIsArray(source) {
			return true
		}
	}

	// Non-array type coercible to ANY[COMPATIBLE]NONARRAY.
	if target == ANYNONARRAYOID || target == ANYCOMPATIBLENONARRAYOID {
		if !c.typeIsArray(source) {
			return true
		}
	}

	// Enum type coercible to ANYENUM.
	if target == ANYENUMOID {
		if c.typeIsEnum(source) {
			return true
		}
	}

	// Range type coercible to ANY[COMPATIBLE]RANGE.
	if target == ANYRANGEOID || target == ANYCOMPATIBLERANGEOID {
		if c.typeIsRange(source) {
			return true
		}
	}

	// Multirange type coercible to ANY[COMPATIBLE]MULTIRANGE.
	if target == ANYMULTIRANGEOID || target == ANYCOMPATIBLEMULTIRANGEOID {
		if c.typeIsMultirange(source) {
			return true
		}
	}

	// Composite type coercible to RECORD.
	if target == RECORDOID {
		if bt := c.typeByOID[source]; bt != nil && bt.Type == 'c' {
			return true
		}
	}

	// Fall back to pg_cast lookup.
	cast := c.castIndex[castKey{source: source, target: target}]
	return cast != nil && cast.Method == 'b'
}

// typeIsEnum returns true if the given OID is an enum type.
//
// pg: src/backend/utils/cache/lsyscache.c — type_is_enum
func (c *Catalog) typeIsEnum(typeOID uint32) bool {
	bt := c.typeByOID[typeOID]
	return bt != nil && bt.Type == 'e'
}

// typeIsArray returns true if the given OID is an array type.
//
// pg: src/backend/utils/cache/lsyscache.c — type_is_array
func (c *Catalog) typeIsArray(typeOID uint32) bool {
	bt := c.typeByOID[typeOID]
	return bt != nil && bt.Elem != 0 && bt.Len == -1
}

// typeIsRange returns true if the given OID is a range type.
//
// pg: src/backend/utils/cache/lsyscache.c — type_is_range
func (c *Catalog) typeIsRange(typeOID uint32) bool {
	_, ok := c.rangeTypes[typeOID]
	return ok
}

// typeIsMultirange returns true if the given OID is a multirange type.
//
// pg: src/backend/utils/cache/lsyscache.c — type_is_multirange
func (c *Catalog) typeIsMultirange(typeOID uint32) bool {
	for _, rt := range c.rangeTypes {
		if rt.MultirangeOID == typeOID {
			return true
		}
	}
	return false
}

// selectCommonTypeCore is the shared algorithm for common type selection.
// If context is non-empty, category mismatches produce an error mentioning it.
// If context is empty, category mismatches return (0, nil) for silent failure.
//
// (pgddl helper — shared core for select_common_type and select_common_type_from_oids)
func (c *Catalog) selectCommonTypeCore(typeOIDs []uint32, context string) (uint32, error) {
	if len(typeOIDs) == 0 {
		return TEXTOID, nil
	}

	ptype := typeOIDs[0]
	startIdx := 1

	// If all input types are valid and exactly the same, just pick that type.
	// This is the only way that we will resolve the result as being a domain
	// type; otherwise domains are smashed to their base types for comparison.
	if ptype != UNKNOWNOID {
		allSame := true
		for i := 1; i < len(typeOIDs); i++ {
			if typeOIDs[i] != ptype {
				startIdx = i
				allSame = false
				break
			}
		}
		if allSame {
			return ptype, nil
		}
	}

	// Set up for the full algorithm. Unwrap domains.
	ptype = c.getBaseType(ptype)
	pt := c.typeByOID[ptype]
	if pt == nil {
		return 0, fmt.Errorf("type OID %d not found", ptype)
	}
	pcategory := pt.Category
	pispreferred := pt.IsPreferred

	for i := startIdx; i < len(typeOIDs); i++ {
		ntype := c.getBaseType(typeOIDs[i])

		// Move on if no new information.
		if ntype == UNKNOWNOID || ntype == ptype {
			continue
		}

		nt := c.typeByOID[ntype]
		if nt == nil {
			return 0, fmt.Errorf("type OID %d not found", ntype)
		}
		ncategory := nt.Category
		nispreferred := nt.IsPreferred

		if ptype == UNKNOWNOID {
			// So far only unknowns, so take anything.
			ptype = ntype
			pcategory = ncategory
			pispreferred = nispreferred
		} else if ncategory != pcategory {
			// Both types in different categories — not much hope.
			if context == "" {
				return 0, nil // silent failure
			}
			return 0, fmt.Errorf("%s types %s and %s cannot be matched",
				context, c.typeName(ptype), c.typeName(ntype))
		} else if !pispreferred &&
			c.CanCoerce(ptype, ntype, 'i') &&
			!c.CanCoerce(ntype, ptype, 'i') {
			// Take new type if can coerce to it implicitly but not the
			// other way; but if we have a preferred type, stay on it.
			ptype = ntype
			pcategory = ncategory
			pispreferred = nispreferred
		}
	}

	// If all inputs were UNKNOWN, resolve as TEXT.
	if ptype == UNKNOWNOID {
		ptype = TEXTOID
	}

	return ptype, nil
}

// selectCommonType picks a common type from a set of analyzed expressions.
// The context string is used in error messages (e.g. "CASE", "UNION", "ARRAY").
// If context is empty, category mismatches return (0, nil) for silent failure
// (matching PG's behavior when context is NULL).
//
// pg: src/backend/parser/parse_coerce.c — select_common_type
func (c *Catalog) selectCommonType(exprs []AnalyzedExpr, context string) (uint32, error) {
	typeOIDs := make([]uint32, len(exprs))
	for i, e := range exprs {
		typeOIDs[i] = e.exprType()
	}
	return c.selectCommonTypeCore(typeOIDs, context)
}

// selectCommonTypeFromOIDs picks a common type from a set of type OIDs.
// If noerror is true, category mismatches return (0, nil) silently.
// If noerror is false, an error is returned with context "argument".
//
// pg: src/backend/parser/parse_coerce.c — select_common_type_from_oids
func (c *Catalog) selectCommonTypeFromOIDs(typeOIDs []uint32, noerror bool) (uint32, error) {
	context := "argument"
	if noerror {
		context = ""
	}
	return c.selectCommonTypeCore(typeOIDs, context)
}

// contextPermits checks whether a cast with context castCtx is allowed
// in the requested context reqCtx.
// Cast contexts: 'i'=implicit, 'a'=assignment, 'e'=explicit only.
// Request contexts: 'i'=implicit allows only 'i'; 'a'=assignment allows 'i'+'a'; 'e'=explicit allows all.
func contextPermits(castCtx, reqCtx byte) bool {
	switch reqCtx {
	case 'e':
		return true // explicit permits all
	case 'a':
		return castCtx == 'i' || castCtx == 'a'
	case 'i':
		return castCtx == 'i'
	default:
		return false
	}
}

// getBaseElementType returns the element type OID for an array type,
// unwrapping domains first. Returns 0 if the type is not an array.
//
// pg: src/backend/utils/cache/lsyscache.c — get_base_element_type
func (c *Catalog) getBaseElementType(typeOID uint32) uint32 {
	base := c.getBaseType(typeOID)
	bt := c.typeByOID[base]
	if bt != nil && bt.Elem != 0 && bt.Category == 'A' {
		return bt.Elem
	}
	return 0
}

// typeFuncClass classifies a type's function-return behavior.
type typeFuncClass int

const (
	typeFuncScalar    typeFuncClass = iota // scalar type
	typeFuncComposite                      // composite (row) type
	typeFuncRecord                         // RECORD pseudo-type
	typeFuncOther                          // other pseudo-type
)

// getTypeFuncClass classifies a type for function-in-FROM processing.
//
// pg: src/backend/utils/fmgr/funcapi.c — get_type_func_class
func (c *Catalog) getTypeFuncClass(typeOID uint32) typeFuncClass {
	bt := c.typeByOID[typeOID]
	if bt == nil {
		return typeFuncOther
	}
	if bt.Type == 'c' && bt.RelID != 0 {
		return typeFuncComposite
	}
	if typeOID == RECORDOID {
		return typeFuncRecord
	}
	if bt.Type == 'p' {
		return typeFuncOther
	}
	return typeFuncScalar
}

// preferredTypeForCategory returns the preferred type OID for the same
// category as the given type, or 0 if none found.
//
// pg: src/backend/parser/parse_coerce.c — find_typmod_coercion_function (category lookup)
func (c *Catalog) preferredTypeForCategory(typeOID uint32) uint32 {
	bt := c.typeByOID[typeOID]
	if bt == nil {
		return 0
	}
	for _, t := range c.typeByOID {
		if t.Category == bt.Category && t.IsPreferred {
			return t.OID
		}
	}
	return 0
}
