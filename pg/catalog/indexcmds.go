package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Index represents an index on a relation.
type Index struct {
	OID              uint32
	Name             string
	Schema           *Schema
	RelOID           uint32
	Columns          []int16 // attnums (0 = expression column)
	IsUnique         bool
	IsPrimary        bool
	IsClustered      bool     // indisclustered
	IsReplicaIdent   bool     // indisreplident
	ConstraintOID    uint32   // 0 if standalone
	AccessMethod     string   // "btree" (default), "hash", "gist", "gin", "brin"
	Tablespace       string   // tablespace the index was created in, or "" for the default
	OwnerName        string   // role that owns the index
	NKeyColumns      int      // number of key columns (rest are INCLUDE); 0 = all are key
	WhereClause      string   // analyzed WHERE predicate text; empty = no partial index
	IndOption        []int16  // per-column flags: bit0=DESC, bit1=NULLS_FIRST
	NullsNotDistinct bool     // NULLS NOT DISTINCT for unique indexes
	Exprs            []string // deparsed expressions for expression columns (attnum=0)
	ColumnNames      []string // chosen per-column index attribute names (key + INCLUDE)
	// OpClasses holds, per key column, the operator class the statement
	// named, dotted when qualified, or "" for the default.
	OpClasses []string
	// Collations holds, per key column, the explicit COLLATE identifier parts
	// the statement named, or nil when omitted.
	Collations [][]string
	// UserCollationOIDs holds, per key column, the resolved user collation OID
	// (0 when omitted or a built-in collation).
	UserCollationOIDs []uint32
	AnalyzedExprs     []AnalyzedExpr
	AnalyzedWhere     AnalyzedExpr
	// ParentIndexOID is the partitioned table's index this partition's index
	// is attached to, or 0. Such an index comes and goes with its parent.
	ParentIndexOID uint32
}

// DefineIndex creates a new index on a relation.
//
// pg: src/backend/commands/indexcmds.c — DefineIndex
func (c *Catalog) DefineIndex(stmt *nodes.IndexStmt) error {
	rv := stmt.Relation
	schema, rel, err := c.findRelation(rv.Schemaname, rv.Relname)
	if err != nil {
		return err
	}

	// Validate relation kind: indexes allowed on tables, partitioned tables, and matviews.
	// pg: src/backend/commands/indexcmds.c — DefineIndex (relkind check)
	if rel.RelKind != 'r' && rel.RelKind != 'p' && rel.RelKind != 'm' {
		return errWrongObjectType(rv.Relname, "a table or materialized view")
	}

	effectiveTablespace := c.effectiveIndexTablespace(stmt.TableSpace, rel)
	if err := c.checkIndexTablespace(effectiveTablespace, rel); err != nil {
		return err
	}

	// Resolve access method. Default to "btree" if empty.
	// pg: src/backend/commands/indexcmds.c — DefineIndex (amcanunique check)
	accessMethod := stmt.AccessMethod
	if accessMethod == "" {
		accessMethod = "btree"
	}

	// Validate AM exists in registry and is an index access method.
	am := c.accessMethodByName[accessMethod]
	if am == nil {
		if !c.snapshot {
			return errUndefinedObject("access method", accessMethod)
		}
		am = &AccessMethod{Name: accessMethod, Type: 'i', OID: c.oidGen.Next()}
	}
	if am.Type != 'i' {
		return &Error{Code: CodeInternalError,
			Message: fmt.Sprintf("access method %q is not of type INDEX", accessMethod)}
	}

	// Validate AM capabilities (only for known built-in AMs; custom AMs pass through).
	isBuiltinAM := am.OID < FirstNormalObjectId
	if stmt.Unique && isBuiltinAM {
		// Only btree and hash support unique indexes.
		// pg: src/backend/commands/indexcmds.c — DefineIndex (line 395)
		if accessMethod != "btree" && accessMethod != "hash" {
			return &Error{Code: CodeFeatureNotSupported,
				Message: fmt.Sprintf("access method %q does not support unique indexes", accessMethod)}
		}
	}

	// INCLUDE columns are only supported by btree, gist, and spgist.
	// pg: src/backend/commands/indexcmds.c — DefineIndex (amcaninclude check)
	if stmt.IndexIncludingParams != nil && len(stmt.IndexIncludingParams.Items) > 0 && isBuiltinAM {
		if accessMethod != "btree" && accessMethod != "gist" && accessMethod != "spgist" {
			return &Error{Code: CodeFeatureNotSupported,
				Message: fmt.Sprintf("access method %q does not support included columns", accessMethod)}
		}
	}

	// Hash indexes don't support multi-column.
	// pg: src/backend/commands/indexcmds.c — DefineIndex (amcanmulticol check)
	if accessMethod == "hash" && stmt.IndexParams != nil && len(stmt.IndexParams.Items) > 1 {
		return &Error{Code: CodeFeatureNotSupported,
			Message: "access method \"hash\" does not support multicolumn indexes"}
	}

	// Check INDEX_MAX_KEYS limit.
	// pg: src/backend/commands/indexcmds.c — DefineIndex (INDEX_MAX_KEYS check)
	nIndexParams := 0
	if stmt.IndexParams != nil {
		nIndexParams = len(stmt.IndexParams.Items)
	}
	nIncludeParams := 0
	if stmt.IndexIncludingParams != nil {
		nIncludeParams = len(stmt.IndexIncludingParams.Items)
	}
	if nIndexParams+nIncludeParams > INDEX_MAX_KEYS {
		return &Error{Code: CodeProgramLimitExceeded,
			Message: fmt.Sprintf("cannot use more than %d columns in an index", INDEX_MAX_KEYS)}
	}

	if err := c.validateRelOptions(stmt.Options, 'i', accessMethod); err != nil {
		return err
	}

	// Extract key column names/expressions and per-column options from IndexParams.
	var rawColNames []string
	var attnums []int16
	var indOption []int16
	var exprTexts []string // deparsed expressions for expression columns
	var analyzedExprs []AnalyzedExpr
	var opClasses []string
	var collations [][]string
	var opClassRefs []opClassRef
	var userCollOIDs []uint32
	var keyUserCollOIDs []uint32
	namedOpClass := false
	namedCollation := false
	if stmt.IndexParams != nil {
		for _, item := range stmt.IndexParams.Items {
			elem, ok := item.(*nodes.IndexElem)
			if !ok {
				continue
			}
			var opClassTypeOID uint32
			if elem.Name == "" {
				if elem.Expr == nil {
					attnums = append(attnums, 0)
					cname := elem.Indexcolname
					if cname == "" {
						cname = "expr"
					}
					rawColNames = append(rawColNames, cname)
				} else {
					analyzed, aErr := c.AnalyzeStandaloneExpr(elem.Expr, rel)
					if aErr != nil {
						return aErr
					}
					// Parenthesized simple column reference ((col)) is treated by
					// PostgreSQL's ComputeIndexAttrs as a plain attribute reference,
					// not an expression column. Whole-row Vars (AttNum == 0) remain
					// expression columns.
					if v, isVar := analyzed.(*VarExpr); isVar && v.LevelsUp == 0 && v.RangeIdx == 0 && v.AttNum != 0 {
						if v.AttNum < 0 {
							return &Error{Code: CodeFeatureNotSupported,
								Message: "index creation on system columns is not supported"}
						}
						ci := indexOfAttNum(rel, v.AttNum)
						if ci < 0 {
							return errUndefinedColumn("")
						}
						attnums = append(attnums, v.AttNum)
						cname := elem.Indexcolname
						if cname == "" {
							cname = rel.Columns[ci].Name
						}
						rawColNames = append(rawColNames, cname)
						opClassTypeOID = rel.Columns[ci].TypeOID
					} else {
						if !c.installingSnapshot {
							if err := c.checkIndexExprImmutable(analyzed); err != nil {
								return err
							}
						}
						attnums = append(attnums, 0)
						cname := elem.Indexcolname
						if cname == "" {
							cname = figureIndexColname(elem.Expr)
						}
						rawColNames = append(rawColNames, cname)
						rte := c.buildRelationRTE(rel)
						exprTexts = append(exprTexts, c.DeparseExpr(analyzed, []*RangeTableEntry{rte}, true))
						analyzedExprs = append(analyzedExprs, analyzed)
						if analyzed != nil {
							opClassTypeOID = analyzed.exprType()
						}
					}
				}
			} else {
				idx, ok := rel.colByName[elem.Name]
				if !ok {
					return errUndefinedColumn(elem.Name)
				}
				// System columns cannot be indexed.
				// pg: src/backend/commands/indexcmds.c — DefineIndex
				if rel.Columns[idx].AttNum < 0 {
					return &Error{Code: CodeFeatureNotSupported,
						Message: "index creation on system columns is not supported"}
				}
				attnums = append(attnums, rel.Columns[idx].AttNum)
				cname := elem.Indexcolname
				if cname == "" {
					cname = elem.Name
				}
				rawColNames = append(rawColNames, cname)
				opClassTypeOID = rel.Columns[idx].TypeOID
			}

			var collParts []string
			var elemUserCollOID uint32
			if elem.Collation != nil && len(elem.Collation.Items) > 0 {
				_, collOID, err := c.resolveColumnCollation(elem.Collation, opClassTypeOID)
				if err != nil {
					return err
				}
				if collOID != 0 {
					userCollOIDs = append(userCollOIDs, collOID)
					elemUserCollOID = collOID
				}
				for _, part := range elem.Collation.Items {
					collParts = append(collParts, stringVal(part))
				}
			}

			// Compute indoption flags: bit0=DESC, bit1=NULLS_FIRST.
			// pg: src/backend/commands/indexcmds.c — DefineIndex (indoption computation)
			var opt int16
			if elem.Ordering == nodes.SORTBY_DESC {
				opt |= 1 // INDOPTION_DESC
				// PG default: DESC implies NULLS FIRST unless explicitly NULLS LAST.
				if elem.NullsOrdering != nodes.SORTBY_NULLS_LAST {
					opt |= 2 // INDOPTION_NULLS_FIRST
				}
			} else {
				// ASC: only set NULLS_FIRST if explicitly requested.
				if elem.NullsOrdering == nodes.SORTBY_NULLS_FIRST {
					opt |= 2 // INDOPTION_NULLS_FIRST
				}
			}
			indOption = append(indOption, opt)

			var opClass []string
			if elem.Opclass != nil {
				for _, part := range elem.Opclass.Items {
					opClass = append(opClass, stringVal(part))
				}
			}
			if len(opClass) > 0 {
				ref := opClassRef{
					schema:  strings.Join(opClass[:len(opClass)-1], "."),
					name:    opClass[len(opClass)-1],
					typeOID: opClassTypeOID,
				}
				if err := c.validateIndexOpClass(ref, am, elem.Opclassopts); err != nil {
					return err
				}
				opClassRefs = append(opClassRefs, ref)
			} else {
				// No opclass named: PostgreSQL still picks the default one for
				// the access method and column/expression type, and depends on it.
				opClassRefs = append(opClassRefs, opClassRef{typeOID: opClassTypeOID})
			}
			opClasses = append(opClasses, strings.Join(opClass, "."))
			namedOpClass = namedOpClass || len(opClass) > 0
			collations = append(collations, collParts)
			keyUserCollOIDs = append(keyUserCollOIDs, elemUserCollOID)
			namedCollation = namedCollation || len(collParts) > 0
		}
	}
	if !namedOpClass {
		opClasses = nil
	}
	if !namedCollation {
		collations = nil
		keyUserCollOIDs = nil
	}

	nKeyColumns := len(attnums)

	// Process INCLUDE columns (IndexIncludingParams).
	// pg: src/backend/commands/indexcmds.c — ComputeIndexAttrs
	if stmt.IndexIncludingParams != nil {
		for _, item := range stmt.IndexIncludingParams.Items {
			elem, ok := item.(*nodes.IndexElem)
			if !ok {
				continue
			}
			if elem.Name == "" {
				return &Error{Code: CodeFeatureNotSupported,
					Message: "expressions are not supported in included columns"}
			}
			if elem.Collation != nil && len(elem.Collation.Items) > 0 {
				return errInvalidObjectDefinition("including column does not support a collation")
			}
			if elem.Opclass != nil && len(elem.Opclass.Items) > 0 {
				return errInvalidObjectDefinition("including column does not support an operator class")
			}
			if elem.Ordering != nodes.SORTBY_DEFAULT {
				return errInvalidObjectDefinition("including column does not support ASC/DESC options")
			}
			if elem.NullsOrdering != nodes.SORTBY_NULLS_DEFAULT {
				return errInvalidObjectDefinition("including column does not support NULLS FIRST/LAST options")
			}
			idx, ok := rel.colByName[elem.Name]
			if !ok {
				return errUndefinedColumn(elem.Name)
			}
			if rel.Columns[idx].AttNum < 0 {
				return &Error{Code: CodeFeatureNotSupported,
					Message: "index creation on system columns is not supported"}
			}
			attnums = append(attnums, rel.Columns[idx].AttNum)
			cname := elem.Indexcolname
			if cname == "" {
				cname = elem.Name
			}
			rawColNames = append(rawColNames, cname)
		}
	}

	if len(attnums) == 0 {
		return errInvalidParameterValue("index must specify at least one column")
	}

	chosenColNames := chooseIndexColumnNames(rawColNames)

	// For unique indexes on partitioned tables, all partition key columns
	// must be included in the index key (INCLUDE columns do not enforce uniqueness).
	// pg: src/backend/commands/indexcmds.c — DefineIndex (line 709-756)
	if stmt.Unique && rel.RelKind == 'p' && rel.PartitionInfo != nil {
		for _, keyAttNum := range rel.PartitionInfo.KeyAttNums {
			if keyAttNum == 0 {
				continue // expression key, skip (can't validate without expression evaluation)
			}
			found := false
			for _, idxAttNum := range attnums[:nKeyColumns] {
				if idxAttNum == keyAttNum {
					found = true
					break
				}
			}
			if !found {
				// Find column name for error message.
				colName := "unknown"
				for _, col := range rel.Columns {
					if col.AttNum == keyAttNum {
						colName = col.Name
						break
					}
				}
				return &Error{Code: CodeFeatureNotSupported,
					Message: fmt.Sprintf("unique constraint on partitioned table must include all partitioning columns\n  Detail: %s is not included in the constraint definition", colName)}
			}
		}
	}

	// Analyze and deparse WHERE clause.
	whereClause := ""
	var analyzedWhere AnalyzedExpr
	if stmt.WhereClause != nil {
		analyzed, aErr := c.AnalyzeStandaloneExpr(stmt.WhereClause, rel)
		if aErr != nil {
			return aErr
		}
		if !c.installingSnapshot {
			coerced, cerr := c.coerceToBoolean(analyzed, "WHERE")
			if cerr != nil {
				return cerr
			}
			if err := c.checkIndexExprImmutable(coerced); err != nil {
				return &Error{Code: CodeInvalidObjectDefinition,
					Message: "functions in index predicate must be marked IMMUTABLE"}
			}
			analyzed = coerced
		}
		rte := c.buildRelationRTE(rel)
		whereClause = c.DeparseExpr(analyzed, []*RangeTableEntry{rte}, true)
		analyzedWhere = analyzed
	}

	// Auto-generate name if empty (ChooseRelationName over all chosen column names).
	name := stmt.Idxname
	if name == "" {
		suffix := "idx"
		cols := chosenColNames
		if stmt.Primary {
			suffix = "pkey"
			cols = nil
		}
		name = c.chooseRelationName(schema, rel.Name, strings.Join(cols, "_"), suffix)
	}

	// Check name uniqueness in schema (Relations, Indexes, Sequences, Types).
	if _, exists := schema.Relations[name]; exists {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", name))
			return nil
		}
		return errDuplicateTable(name)
	}
	if _, exists := schema.Sequences[name]; exists {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", name))
			return nil
		}
		return errDuplicateTable(name)
	}
	if _, exists := schema.Indexes[name]; exists {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", name))
			return nil
		}
		return errDuplicateObject("index", name)
	}
	idx := &Index{
		OID:               c.oidGen.Next(),
		Name:              name,
		Schema:            schema,
		RelOID:            rel.OID,
		Columns:           attnums,
		ColumnNames:       chosenColNames,
		IsUnique:          stmt.Unique,
		IsPrimary:         stmt.Primary,
		AccessMethod:      accessMethod,
		Tablespace:        normalizeTablespace(effectiveTablespace),
		OwnerName:         rel.OwnerName,
		NKeyColumns:       nKeyColumns,
		WhereClause:       whereClause,
		IndOption:         indOption,
		NullsNotDistinct:  stmt.Nulls_not_distinct,
		Exprs:             exprTexts,
		OpClasses:         opClasses,
		Collations:        collations,
		UserCollationOIDs: keyUserCollOIDs,
		AnalyzedExprs:     analyzedExprs,
		AnalyzedWhere:     analyzedWhere,
	}

	c.registerIndex(schema, idx)
	// pg: index_create() — whole-relation dep + per-column deps
	c.recordDependency('i', idx.OID, 0, 'r', rel.OID, 0, DepAuto)
	for _, attnum := range attnums {
		if attnum != 0 { // 0 = expression column, no column dep
			c.recordDependency('i', idx.OID, 0, 'r', rel.OID, int32(attnum), DepAuto)
		}
	}
	// pg: index_create() — expressions and predicate depend AUTO on the
	// columns they read and NORMAL on everything else they reference.
	for _, e := range analyzedExprs {
		c.recordDependencyOnSingleRelExpr('i', idx.OID, e, rel.OID, DepNormal, DepAuto)
	}
	if analyzedWhere != nil {
		c.recordDependencyOnSingleRelExpr('i', idx.OID, analyzedWhere, rel.OID, DepNormal, DepAuto)
	}
	c.recordOpClassDependencies('i', idx.OID, opClassRefs, accessMethod)
	for _, collOID := range userCollOIDs {
		c.recordDependency('i', idx.OID, 0, 'L', collOID, 0, DepNormal)
	}
	if am := c.accessMethodByName[accessMethod]; am != nil {
		c.recordDependency('i', idx.OID, 0, 'm', am.OID, 0, DepNormal)
	}

	return nil
}

// createIndexInternal creates an index internally (shared by CreateIndex and constraint auto-creation).
func (c *Catalog) createIndexInternal(schema *Schema, rel *Relation, name string, attnums []int16, nKeyCols int, isUnique, isPrimary bool, constraintOID uint32) *Index {
	if nKeyCols <= 0 || nKeyCols > len(attnums) {
		nKeyCols = len(attnums)
	}
	// Default IndOption: all zeros (ASC NULLS LAST for btree) on key columns.
	indOption := make([]int16, nKeyCols)
	rawNames := make([]string, 0, len(attnums))
	for _, attnum := range attnums {
		if ci := indexOfAttNum(rel, attnum); ci >= 0 {
			rawNames = append(rawNames, rel.Columns[ci].Name)
		} else {
			rawNames = append(rawNames, "expr")
		}
	}

	idx := &Index{
		OID:           c.oidGen.Next(),
		Name:          name,
		Schema:        schema,
		RelOID:        rel.OID,
		Columns:       attnums,
		ColumnNames:   chooseIndexColumnNames(rawNames),
		IsUnique:      isUnique,
		IsPrimary:     isPrimary,
		ConstraintOID: constraintOID,
		AccessMethod:  "btree",
		OwnerName:     rel.OwnerName,
		NKeyColumns:   nKeyCols,
		IndOption:     indOption,
	}

	c.registerIndex(schema, idx)
	if am := c.accessMethodByName[idx.AccessMethod]; am != nil {
		c.recordDependency('i', idx.OID, 0, 'm', am.OID, 0, DepNormal)
	}
	c.recordOpClassDependencies('i', idx.OID, c.indexOpClassRefs(rel, attnums, nKeyCols, nil), idx.AccessMethod)
	// pg: index_create() — whole-relation dep + per-column deps
	c.recordDependency('i', idx.OID, 0, 'r', rel.OID, 0, DepAuto)
	for _, attnum := range attnums {
		if attnum != 0 {
			c.recordDependency('i', idx.OID, 0, 'r', rel.OID, int32(attnum), DepAuto)
		}
	}

	return idx
}

// indexOpClassRefs builds per-key-column operator-class references (including
// column types for default or overloaded operator-class lookup) for indexes
// constructed from resolved column attnums.
func (c *Catalog) indexOpClassRefs(rel *Relation, attnums []int16, nKeyColumns int, opClasses []string) []opClassRef {
	nKey := nKeyColumns
	if nKey <= 0 || nKey > len(attnums) {
		nKey = len(attnums)
	}
	refs := make([]opClassRef, 0, nKey)
	for i := 0; i < nKey; i++ {
		var ref opClassRef
		if i < len(opClasses) && opClasses[i] != "" {
			ref.schema, ref.name = "", opClasses[i]
			if dot := strings.LastIndexByte(opClasses[i], '.'); dot >= 0 {
				ref.schema, ref.name = opClasses[i][:dot], opClasses[i][dot+1:]
			}
		}
		if rel != nil && attnums[i] > 0 {
			if ci := indexOfAttNum(rel, attnums[i]); ci >= 0 {
				ref.typeOID = rel.Columns[ci].TypeOID
			}
		}
		refs = append(refs, ref)
	}
	return refs
}

// registerIndex adds an index to all catalog maps.
// opClassRef is one operator-class reference an index column made: the schema
// it named (empty when unqualified), its name, and the indexed type it must
// apply to (0 for an expression column).
type opClassRef struct {
	schema  string
	name    string
	typeOID uint32
}

// resolveOpClassOID finds the operator class a column reference names, the way
// PostgreSQL does: by schema (or the first on the search path), access method,
// and indexed type. It returns 0 when nothing matches.
//
// pg: src/backend/commands/opclasscmds.c — GetIndexOpClass
func (c *Catalog) resolveOpClassOID(ref opClassRef, amOID uint32) uint32 {
	// An index with no opclass for a column takes the default one for its
	// access method and type.
	// pg: src/backend/commands/indexcmds.c — GetDefaultOpClass
	if ref.name == "" {
		if ref.typeOID == 0 {
			return 0
		}
		if opc := c.opClassByKey[opClassKey{amOID: amOID, typeOID: ref.typeOID}]; opc != nil {
			return opc.OID
		}
		return 0
	}
	var best, byType uint32
	consider := func(oid uint32, opc *OpClass) {
		if opc.Name != ref.name || opc.AMOID != amOID {
			return
		}
		if ref.typeOID != 0 && opc.TypeOID == ref.typeOID {
			if byType == 0 || oid < byType {
				byType = oid
			}
			return
		}
		if best == 0 || oid < best {
			best = oid
		}
	}
	if ref.schema != "" {
		s := c.schemaByName[ref.schema]
		if s == nil {
			return 0
		}
		for oid, opc := range c.opClasses {
			if opc.Schema == s.OID {
				consider(oid, opc)
			}
		}
	} else {
		// An unqualified name resolves in the first schema on the search path
		// that holds it, as GetIndexOpClass does through OpclassnameGetOpcid.
		// Without a known column type (an expression element), a name match in
		// that schema is enough: it is the visible class.
		for _, nsOID := range c.searchPathWithCatalog() {
			schemaBest, schemaType := uint32(0), uint32(0)
			for oid, opc := range c.opClasses {
				if opc.Schema != nsOID || opc.Name != ref.name || opc.AMOID != amOID {
					continue
				}
				if schemaBest == 0 || oid < schemaBest {
					schemaBest = oid
				}
				if ref.typeOID != 0 && opc.TypeOID == ref.typeOID && (schemaType == 0 || oid < schemaType) {
					schemaType = oid
				}
			}
			if schemaType != 0 {
				return schemaType
			}
			if schemaBest != 0 {
				return schemaBest
			}
		}
	}
	if byType != 0 {
		return byType
	}
	return best
}

// recordOpClassDependencies records a NORMAL dependency of an index on each
// operator class its columns name, so dropping the operator class (or the
// extension that owns it) reports the index as a dependent.
//
// pg: src/backend/catalog/index.c — index_create (opclass dependencies)
func (c *Catalog) recordOpClassDependencies(objType byte, objOID uint32, refs []opClassRef, amName string) {
	am := c.accessMethodByName[amName]
	if am == nil {
		return
	}
	seen := map[uint32]bool{}
	for _, ref := range refs {
		oid := c.resolveOpClassOID(ref, am.OID)
		if oid == 0 || seen[oid] {
			continue
		}
		seen[oid] = true
		c.recordDependency(objType, objOID, 0, 'o', oid, 0, DepNormal)
	}
}

func (c *Catalog) registerIndex(schema *Schema, idx *Index) {
	c.indexes[idx.OID] = idx
	c.indexesByRel[idx.RelOID] = append(c.indexesByRel[idx.RelOID], idx)
	schema.Indexes[idx.Name] = idx
}

// removeIndex removes an index from all catalog maps and cleans up dependencies.
func (c *Catalog) removeIndex(schema *Schema, name string, idx *Index) {
	delete(schema.Indexes, name)
	delete(c.indexes, idx.OID)

	// Remove from indexesByRel.
	list := c.indexesByRel[idx.RelOID]
	for i, x := range list {
		if x.OID == idx.OID {
			c.indexesByRel[idx.RelOID] = append(list[:i], list[i+1:]...)
			break
		}
	}

	c.removeComments('i', idx.OID)
	c.removeDepsOf('i', idx.OID)
	c.removeDepsOn('i', idx.OID)
	c.removeDepsOn('r', idx.OID)
}

// removeIndexesForRelation removes all indexes belonging to a relation.
func (c *Catalog) removeIndexesForRelation(relOID uint32, schema *Schema) {
	for _, idx := range c.indexesByRel[relOID] {
		delete(schema.Indexes, idx.Name)
		delete(c.indexes, idx.OID)
		c.removeComments('i', idx.OID)
		c.removeDepsOf('i', idx.OID)
		c.removeDepsOn('i', idx.OID)
		c.removeDepsOn('r', idx.OID)
	}
	delete(c.indexesByRel, relOID)
}

// resolveColumnNames resolves column name strings to attnums.
func (c *Catalog) resolveColumnNames(rel *Relation, names []string) ([]int16, error) {
	attnums := make([]int16, len(names))
	for i, name := range names {
		idx, ok := rel.colByName[name]
		if !ok {
			return nil, errUndefinedColumn(name)
		}
		attnums[i] = rel.Columns[idx].AttNum
	}
	return attnums, nil
}

// generateIndexName generates a default index name.
func generateIndexName(tableName string, columns []string, isPrimary bool) string {
	if isPrimary {
		return tableName + "_pkey"
	}
	return tableName + "_" + strings.Join(columns, "_") + "_idx"
}

// chooseIndexColumnNames implements PostgreSQL's ChooseIndexColumnNames,
// appending 1, 2, ... when two index columns (across key and INCLUDE columns)
// share the same base name.
//
// pg: src/backend/commands/indexcmds.c — ChooseIndexColumnNames
func chooseIndexColumnNames(raw []string) []string {
	out := make([]string, 0, len(raw))
	used := make(map[string]bool, len(raw))
	for _, base := range raw {
		if base == "" {
			base = "expr"
		}
		name := base
		for k := 1; used[name]; k++ {
			name = fmt.Sprintf("%s%d", base, k)
		}
		used[name] = true
		out = append(out, name)
	}
	return out
}

// figureIndexColname derives PostgreSQL's default column name for an index
// expression (e.g., whole-row (a) -> "a", lower(x) -> "lower", else "expr").
//
// pg: src/backend/parser/parse_target.c — FigureIndexColname
func figureIndexColname(node nodes.Node) string {
	switch n := node.(type) {
	case *nodes.ColumnRef:
		if n.Fields != nil && len(n.Fields.Items) > 0 {
			last := n.Fields.Items[len(n.Fields.Items)-1]
			if _, isStar := last.(*nodes.A_Star); isStar && len(n.Fields.Items) >= 2 {
				if s := stringVal(n.Fields.Items[len(n.Fields.Items)-2]); s != "" {
					return s
				}
			}
			if s := stringVal(last); s != "" {
				return s
			}
		}
	case *nodes.FuncCall:
		if n.Funcname != nil && len(n.Funcname.Items) > 0 {
			if s := stringVal(n.Funcname.Items[len(n.Funcname.Items)-1]); s != "" {
				return s
			}
		}
	case *nodes.TypeCast:
		return figureIndexColname(n.Arg)
	case *nodes.CollateClause:
		return figureIndexColname(n.Arg)
	}
	return "expr"
}

// chooseRelationName picks an unused relation/type name in schema starting from
// name1[_name2]_label and appending 1, 2, ... on collision.
//
// pg: src/backend/commands/indexcmds.c — ChooseRelationName
func (c *Catalog) chooseRelationName(schema *Schema, name1, name2, label string) string {
	base := name1
	if name2 != "" {
		base += "_" + name2
	}
	if label != "" {
		base += "_" + label
	}
	cand := base
	for k := 1; ; k++ {
		if _, relHeld := relationNamespaceKind(schema, cand); !relHeld {
			return cand
		}
		cand = fmt.Sprintf("%s%d", base, k)
	}
}

// checkIndexExprImmutable verifies that an analyzed index expression contains
// no mutable (STABLE / VOLATILE) functions, operators, or SQLValueFunctions.
//
// pg: src/backend/commands/indexcmds.c — ComputeIndexAttrs (contain_mutable_functions)
func (c *Catalog) checkIndexExprImmutable(expr AnalyzedExpr) error {
	if expr == nil {
		return nil
	}
	mutable := false
	WalkAnalyzedExpr(expr, func(e AnalyzedExpr) bool {
		if mutable {
			return false
		}
		switch v := e.(type) {
		case *SQLValueFuncExpr:
			mutable = true
			return false
		case *FuncCallExpr:
			if c.isProcMutable(v.FuncOID) {
				mutable = true
				return false
			}
		case *OpExpr:
			if op := c.operByOID[v.OpOID]; op != nil && c.isProcMutable(op.Code) {
				mutable = true
				return false
			}
		case *DistinctExprQ:
			if op := c.operByOID[v.OpOID]; op != nil && c.isProcMutable(op.Code) {
				mutable = true
				return false
			}
		case *NullIfExprQ:
			if op := c.operByOID[v.OpOID]; op != nil && c.isProcMutable(op.Code) {
				mutable = true
				return false
			}
		case *ScalarArrayOpExpr:
			if op := c.operByOID[v.OpOID]; op != nil && c.isProcMutable(op.Code) {
				mutable = true
				return false
			}
		case *CoerceViaIOExpr:
			var inType uint32
			if v.Arg != nil {
				inType = v.Arg.exprType()
			}
			if c.isTypeIOFuncMutable(inType, false) || c.isTypeIOFuncMutable(v.ResultType, true) {
				mutable = true
				return false
			}
		}
		return true
	})
	if mutable {
		return errInvalidObjectDefinition("functions in index expression must be marked IMMUTABLE")
	}
	return nil
}

func (c *Catalog) isProcMutable(procOID uint32) bool {
	if procOID == 0 {
		return false
	}
	bp := c.procByOID[procOID]
	return bp != nil && bp.Volatile != 0 && bp.Volatile != 'i'
}

func (c *Catalog) isTypeIOFuncMutable(typeOID uint32, isInput bool) bool {
	if typeOID == 0 {
		return false
	}
	for dt := c.domainTypes[typeOID]; dt != nil; dt = c.domainTypes[typeOID] {
		typeOID = dt.BaseTypeOID
	}
	bt := c.typeByOID[typeOID]
	if bt == nil {
		return false
	}
	fnOID := bt.Output
	if isInput {
		fnOID = bt.Input
	}
	if c.isProcMutable(fnOID) {
		return true
	}
	if bt.Elem != 0 {
		if c.isTypeIOFuncMutable(bt.Elem, isInput) {
			return true
		}
	}
	if bt.RelID != 0 {
		if rel := c.relationByOID[bt.RelID]; rel != nil {
			for _, col := range rel.Columns {
				if c.isTypeIOFuncMutable(col.TypeOID, isInput) {
					return true
				}
			}
		}
	}
	return false
}

// builtinOpClasses maps (amName, opcName) to the input type OIDs registered in
// PostgreSQL 17's pg_catalog.pg_opclass.
var builtinOpClasses = map[[2]string][]uint32{
	{"brin", "bit_minmax_ops"}:               {1560},
	{"brin", "box_inclusion_ops"}:            {603},
	{"brin", "bpchar_bloom_ops"}:             {1042},
	{"brin", "bpchar_minmax_ops"}:            {1042},
	{"brin", "bytea_bloom_ops"}:              {17},
	{"brin", "bytea_minmax_ops"}:             {17},
	{"brin", "char_bloom_ops"}:               {18},
	{"brin", "char_minmax_ops"}:              {18},
	{"brin", "date_bloom_ops"}:               {1082},
	{"brin", "date_minmax_multi_ops"}:        {1082},
	{"brin", "date_minmax_ops"}:              {1082},
	{"brin", "float4_bloom_ops"}:             {700},
	{"brin", "float4_minmax_multi_ops"}:      {700},
	{"brin", "float4_minmax_ops"}:            {700},
	{"brin", "float8_bloom_ops"}:             {701},
	{"brin", "float8_minmax_multi_ops"}:      {701},
	{"brin", "float8_minmax_ops"}:            {701},
	{"brin", "inet_bloom_ops"}:               {869},
	{"brin", "inet_inclusion_ops"}:           {869},
	{"brin", "inet_minmax_multi_ops"}:        {869},
	{"brin", "inet_minmax_ops"}:              {869},
	{"brin", "int2_bloom_ops"}:               {21},
	{"brin", "int2_minmax_multi_ops"}:        {21},
	{"brin", "int2_minmax_ops"}:              {21},
	{"brin", "int4_bloom_ops"}:               {23},
	{"brin", "int4_minmax_multi_ops"}:        {23},
	{"brin", "int4_minmax_ops"}:              {23},
	{"brin", "int8_bloom_ops"}:               {20},
	{"brin", "int8_minmax_multi_ops"}:        {20},
	{"brin", "int8_minmax_ops"}:              {20},
	{"brin", "interval_bloom_ops"}:           {1186},
	{"brin", "interval_minmax_multi_ops"}:    {1186},
	{"brin", "interval_minmax_ops"}:          {1186},
	{"brin", "macaddr8_bloom_ops"}:           {774},
	{"brin", "macaddr8_minmax_multi_ops"}:    {774},
	{"brin", "macaddr8_minmax_ops"}:          {774},
	{"brin", "macaddr_bloom_ops"}:            {829},
	{"brin", "macaddr_minmax_multi_ops"}:     {829},
	{"brin", "macaddr_minmax_ops"}:           {829},
	{"brin", "name_bloom_ops"}:               {19},
	{"brin", "name_minmax_ops"}:              {19},
	{"brin", "numeric_bloom_ops"}:            {1700},
	{"brin", "numeric_minmax_multi_ops"}:     {1700},
	{"brin", "numeric_minmax_ops"}:           {1700},
	{"brin", "oid_bloom_ops"}:                {26},
	{"brin", "oid_minmax_multi_ops"}:         {26},
	{"brin", "oid_minmax_ops"}:               {26},
	{"brin", "pg_lsn_bloom_ops"}:             {3220},
	{"brin", "pg_lsn_minmax_multi_ops"}:      {3220},
	{"brin", "pg_lsn_minmax_ops"}:            {3220},
	{"brin", "range_inclusion_ops"}:          {3831},
	{"brin", "text_bloom_ops"}:               {25},
	{"brin", "text_minmax_ops"}:              {25},
	{"brin", "tid_bloom_ops"}:                {27},
	{"brin", "tid_minmax_multi_ops"}:         {27},
	{"brin", "tid_minmax_ops"}:               {27},
	{"brin", "time_bloom_ops"}:               {1083},
	{"brin", "time_minmax_multi_ops"}:        {1083},
	{"brin", "time_minmax_ops"}:              {1083},
	{"brin", "timestamp_bloom_ops"}:          {1114},
	{"brin", "timestamp_minmax_multi_ops"}:   {1114},
	{"brin", "timestamp_minmax_ops"}:         {1114},
	{"brin", "timestamptz_bloom_ops"}:        {1184},
	{"brin", "timestamptz_minmax_multi_ops"}: {1184},
	{"brin", "timestamptz_minmax_ops"}:       {1184},
	{"brin", "timetz_bloom_ops"}:             {1266},
	{"brin", "timetz_minmax_multi_ops"}:      {1266},
	{"brin", "timetz_minmax_ops"}:            {1266},
	{"brin", "uuid_bloom_ops"}:               {2950},
	{"brin", "uuid_minmax_multi_ops"}:        {2950},
	{"brin", "uuid_minmax_ops"}:              {2950},
	{"brin", "varbit_minmax_ops"}:            {1562},
	{"btree", "array_ops"}:                   {2277},
	{"btree", "bit_ops"}:                     {1560},
	{"btree", "bool_ops"}:                    {16},
	{"btree", "bpchar_ops"}:                  {1042},
	{"btree", "bpchar_pattern_ops"}:          {1042},
	{"btree", "bytea_ops"}:                   {17},
	{"btree", "char_ops"}:                    {18},
	{"btree", "cidr_ops"}:                    {869},
	{"btree", "date_ops"}:                    {1082},
	{"btree", "enum_ops"}:                    {3500},
	{"btree", "float4_ops"}:                  {700},
	{"btree", "float8_ops"}:                  {701},
	{"btree", "inet_ops"}:                    {869},
	{"btree", "int2_ops"}:                    {21},
	{"btree", "int4_ops"}:                    {23},
	{"btree", "int8_ops"}:                    {20},
	{"btree", "interval_ops"}:                {1186},
	{"btree", "jsonb_ops"}:                   {3802},
	{"btree", "macaddr8_ops"}:                {774},
	{"btree", "macaddr_ops"}:                 {829},
	{"btree", "money_ops"}:                   {790},
	{"btree", "multirange_ops"}:              {4537},
	{"btree", "name_ops"}:                    {19},
	{"btree", "numeric_ops"}:                 {1700},
	{"btree", "oid_ops"}:                     {26},
	{"btree", "oidvector_ops"}:               {30},
	{"btree", "pg_lsn_ops"}:                  {3220},
	{"btree", "range_ops"}:                   {3831},
	{"btree", "record_image_ops"}:            {2249},
	{"btree", "record_ops"}:                  {2249},
	{"btree", "text_ops"}:                    {25},
	{"btree", "text_pattern_ops"}:            {25},
	{"btree", "tid_ops"}:                     {27},
	{"btree", "time_ops"}:                    {1083},
	{"btree", "timestamp_ops"}:               {1114},
	{"btree", "timestamptz_ops"}:             {1184},
	{"btree", "timetz_ops"}:                  {1266},
	{"btree", "tsquery_ops"}:                 {3615},
	{"btree", "tsvector_ops"}:                {3614},
	{"btree", "uuid_ops"}:                    {2950},
	{"btree", "varbit_ops"}:                  {1562},
	{"btree", "varchar_ops"}:                 {25},
	{"btree", "varchar_pattern_ops"}:         {25},
	{"btree", "xid8_ops"}:                    {5069},
	{"gin", "array_ops"}:                     {2277},
	{"gin", "jsonb_ops"}:                     {3802},
	{"gin", "jsonb_path_ops"}:                {3802},
	{"gin", "tsvector_ops"}:                  {3614},
	{"gist", "box_ops"}:                      {603},
	{"gist", "circle_ops"}:                   {718},
	{"gist", "inet_ops"}:                     {869},
	{"gist", "multirange_ops"}:               {4537},
	{"gist", "point_ops"}:                    {600},
	{"gist", "poly_ops"}:                     {604},
	{"gist", "range_ops"}:                    {3831},
	{"gist", "tsquery_ops"}:                  {3615},
	{"gist", "tsvector_ops"}:                 {3614},
	{"hash", "aclitem_ops"}:                  {1033},
	{"hash", "array_ops"}:                    {2277},
	{"hash", "bool_ops"}:                     {16},
	{"hash", "bpchar_ops"}:                   {1042},
	{"hash", "bpchar_pattern_ops"}:           {1042},
	{"hash", "bytea_ops"}:                    {17},
	{"hash", "char_ops"}:                     {18},
	{"hash", "cid_ops"}:                      {29},
	{"hash", "cidr_ops"}:                     {869},
	{"hash", "date_ops"}:                     {1082},
	{"hash", "enum_ops"}:                     {3500},
	{"hash", "float4_ops"}:                   {700},
	{"hash", "float8_ops"}:                   {701},
	{"hash", "inet_ops"}:                     {869},
	{"hash", "int2_ops"}:                     {21},
	{"hash", "int4_ops"}:                     {23},
	{"hash", "int8_ops"}:                     {20},
	{"hash", "interval_ops"}:                 {1186},
	{"hash", "jsonb_ops"}:                    {3802},
	{"hash", "macaddr8_ops"}:                 {774},
	{"hash", "macaddr_ops"}:                  {829},
	{"hash", "multirange_ops"}:               {4537},
	{"hash", "name_ops"}:                     {19},
	{"hash", "numeric_ops"}:                  {1700},
	{"hash", "oid_ops"}:                      {26},
	{"hash", "oidvector_ops"}:                {30},
	{"hash", "pg_lsn_ops"}:                   {3220},
	{"hash", "range_ops"}:                    {3831},
	{"hash", "record_ops"}:                   {2249},
	{"hash", "text_ops"}:                     {25},
	{"hash", "text_pattern_ops"}:             {25},
	{"hash", "tid_ops"}:                      {27},
	{"hash", "time_ops"}:                     {1083},
	{"hash", "timestamp_ops"}:                {1114},
	{"hash", "timestamptz_ops"}:              {1184},
	{"hash", "timetz_ops"}:                   {1266},
	{"hash", "uuid_ops"}:                     {2950},
	{"hash", "varchar_ops"}:                  {25},
	{"hash", "varchar_pattern_ops"}:          {25},
	{"hash", "xid8_ops"}:                     {5069},
	{"hash", "xid_ops"}:                      {28},
	{"spgist", "box_ops"}:                    {603},
	{"spgist", "inet_ops"}:                   {869},
	{"spgist", "kd_point_ops"}:               {600},
	{"spgist", "poly_ops"}:                   {604},
	{"spgist", "quad_point_ops"}:             {600},
	{"spgist", "range_ops"}:                  {3831},
	{"spgist", "text_ops"}:                   {25},
}

func (c *Catalog) opClassAcceptsType(attrTypeOID, opcTypeOID uint32) bool {
	if attrTypeOID == 0 || opcTypeOID == 0 {
		return true
	}
	return c.IsBinaryCoercible(attrTypeOID, opcTypeOID)
}

func builtinOpClassAcceptsOption(amName, opcName, optKey string) bool {
	switch amName {
	case "brin":
		if strings.HasSuffix(opcName, "_minmax_multi_ops") {
			return optKey == "values_per_range"
		}
		if strings.HasSuffix(opcName, "_bloom_ops") {
			return optKey == "n_distinct_per_range" || optKey == "false_positive_rate"
		}
	case "gist":
		if opcName == "tsvector_ops" {
			return optKey == "siglen"
		}
	}
	return false
}

// validateIndexOpClass validates an explicit operator class (and any operator
// class parameters) on an index element.
//
// pg: src/backend/commands/indexcmds.c — ResolveOpClass
func (c *Catalog) validateIndexOpClass(ref opClassRef, am *AccessMethod, opclassOpts *nodes.List) error {
	if am == nil {
		return nil
	}
	if ref.schema != "" {
		if c.schemaByName[ref.schema] == nil {
			return errUndefinedSchema(ref.schema)
		}
	}
	if c.snapshot || c.installingSnapshot {
		return nil
	}
	if am.OID >= FirstNormalObjectId {
		hasAny := false
		for _, opc := range c.opClasses {
			if opc.AMOID == am.OID {
				hasAny = true
				break
			}
		}
		if !hasAny {
			return nil
		}
	}

	collectInNamespace := func(nsOID uint32) ([]uint32, bool) {
		var types []uint32
		hasUserOpc := false
		for _, opc := range c.opClasses {
			if opc.Schema == nsOID && opc.AMOID == am.OID && opc.Name == ref.name {
				types = append(types, opc.TypeOID)
				hasUserOpc = true
			}
		}
		if nsOID == PGCatalogNamespace && am.OID < FirstNormalObjectId {
			types = append(types, builtinOpClasses[[2]string{am.Name, ref.name}]...)
		}
		return types, hasUserOpc
	}

	var matchedTypes []uint32
	hasUserOpc := false
	if ref.schema != "" {
		s := c.schemaByName[ref.schema]
		matchedTypes, hasUserOpc = collectInNamespace(s.OID)
	} else {
		for _, nsOID := range c.searchPathWithCatalog() {
			matchedTypes, hasUserOpc = collectInNamespace(nsOID)
			if len(matchedTypes) > 0 {
				break
			}
		}
	}
	if len(matchedTypes) == 0 {
		return errUndefinedObject("operator class", fmt.Sprintf("%s for access method %s", ref.name, am.Name))
	}
	if ref.typeOID != 0 {
		okType := false
		for _, opcType := range matchedTypes {
			if c.opClassAcceptsType(ref.typeOID, opcType) {
				okType = true
				break
			}
		}
		if !okType {
			return errDatatypeMismatch(fmt.Sprintf("operator class %q does not accept data type %s",
				ref.name, c.formatType(ref.typeOID, -1)))
		}
	}
	if opclassOpts != nil && len(opclassOpts.Items) > 0 {
		seen := make(map[string]bool, len(opclassOpts.Items))
		for _, item := range opclassOpts.Items {
			d, ok := item.(*nodes.DefElem)
			if !ok {
				continue
			}
			if d.Defnamespace != "" {
				return errInvalidParameterValue(fmt.Sprintf("unrecognized parameter namespace %q", d.Defnamespace))
			}
			key := strings.ToLower(d.Defname)
			if seen[key] {
				return errInvalidParameterValue(fmt.Sprintf("parameter %q specified more than once", d.Defname))
			}
			seen[key] = true
			if !hasUserOpc && am.OID < FirstNormalObjectId && !builtinOpClassAcceptsOption(am.Name, ref.name, key) {
				return errInvalidParameterValue(fmt.Sprintf("unrecognized parameter %q", d.Defname))
			}
		}
	}
	return nil
}
