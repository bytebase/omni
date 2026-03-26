package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Index represents an index on a relation.
type Index struct {
	OID            uint32
	Name           string
	Schema         *Schema
	RelOID         uint32
	Columns        []int16 // attnums (0 = expression column)
	IsUnique       bool
	IsPrimary      bool
	IsClustered    bool   // indisclustered
	IsReplicaIdent bool   // indisreplident
	ConstraintOID  uint32 // 0 if standalone
	AccessMethod   string // "btree" (default), "hash", "gist", "gin", "brin"
	NKeyColumns    int    // number of key columns (rest are INCLUDE); 0 = all are key
	WhereClause       string   // analyzed WHERE predicate text; empty = no partial index
	IndOption         []int16  // per-column flags: bit0=DESC, bit1=NULLS_FIRST
	NullsNotDistinct  bool     // NULLS NOT DISTINCT for unique indexes
	Exprs             []string // deparsed expressions for expression columns (attnum=0)
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

	// Resolve access method. Default to "btree" if empty.
	// pg: src/backend/commands/indexcmds.c — DefineIndex (amcanunique check)
	accessMethod := stmt.AccessMethod
	if accessMethod == "" {
		accessMethod = "btree"
	}

	// Validate AM exists in registry.
	am := c.accessMethodByName[accessMethod]
	if am == nil {
		return errUndefinedObject("access method", accessMethod)
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

	// Extract key column names/expressions and per-column options from IndexParams.
	var colNames []string
	var attnums []int16
	var indOption []int16
	var exprTexts []string // deparsed expressions for expression columns
	if stmt.IndexParams != nil {
		for _, item := range stmt.IndexParams.Items {
			elem, ok := item.(*nodes.IndexElem)
			if !ok {
				continue
			}
			if elem.Name == "" {
				// Expression index element: store attnum=0 (PG convention).
				attnums = append(attnums, 0)
				colNames = append(colNames, "expr")
				// Analyze and deparse the expression.
				if elem.Expr != nil {
					analyzed, aErr := c.AnalyzeStandaloneExpr(elem.Expr, rel)
					if aErr == nil {
						rte := c.buildRelationRTE(rel)
						exprTexts = append(exprTexts, c.DeparseExpr(analyzed, []*RangeTableEntry{rte}, true))
					} else {
						return fmt.Errorf("failed to analyze index expression: %w", aErr)
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
						Message: fmt.Sprintf("index creation on system columns is not supported")}
				}
				attnums = append(attnums, rel.Columns[idx].AttNum)
				colNames = append(colNames, elem.Name)
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
		}
	}

	nKeyColumns := len(attnums)

	// Process INCLUDE columns (IndexIncludingParams).
	if stmt.IndexIncludingParams != nil {
		for _, item := range stmt.IndexIncludingParams.Items {
			elem, ok := item.(*nodes.IndexElem)
			if !ok {
				continue
			}
			if elem.Name == "" {
				continue // INCLUDE does not support expressions
			}
			idx, ok := rel.colByName[elem.Name]
			if !ok {
				return errUndefinedColumn(elem.Name)
			}
			attnums = append(attnums, rel.Columns[idx].AttNum)
		}
	}

	if len(attnums) == 0 {
		return errInvalidParameterValue("index must specify at least one column")
	}

	// For unique indexes on partitioned tables, all partition key columns
	// must be included in the index.
	// pg: src/backend/commands/indexcmds.c — DefineIndex (line 709-756)
	if stmt.Unique && rel.RelKind == 'p' && rel.PartitionInfo != nil {
		for _, keyAttNum := range rel.PartitionInfo.KeyAttNums {
			if keyAttNum == 0 {
				continue // expression key, skip (can't validate without expression evaluation)
			}
			found := false
			for _, idxAttNum := range attnums {
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
	if stmt.WhereClause != nil {
		analyzed, aErr := c.AnalyzeStandaloneExpr(stmt.WhereClause, rel)
		if aErr == nil {
			rte := c.buildRelationRTE(rel)
			whereClause = c.DeparseExpr(analyzed, []*RangeTableEntry{rte}, true)
		} else {
			return fmt.Errorf("failed to analyze index WHERE clause: %w", aErr)
		}
	}

	// Auto-generate name if empty.
	name := stmt.Idxname
	if name == "" {
		// Use only key column names for index name generation.
		name = generateIndexName(rel.Name, colNames, stmt.Primary)
	}

	// Check name uniqueness in schema (both Relations and Indexes).
	if _, exists := schema.Relations[name]; exists {
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
		OID:              c.oidGen.Next(),
		Name:             name,
		Schema:           schema,
		RelOID:           rel.OID,
		Columns:          attnums,
		IsUnique:         stmt.Unique,
		IsPrimary:        stmt.Primary,
		AccessMethod:     accessMethod,
		NKeyColumns:      nKeyColumns,
		WhereClause:      whereClause,
		IndOption:        indOption,
		NullsNotDistinct: stmt.Nulls_not_distinct,
		Exprs:            exprTexts,
	}

	c.registerIndex(schema, idx)
	// pg: index_create() — whole-relation dep + per-column deps
	c.recordDependency('i', idx.OID, 0, 'r', rel.OID, 0, DepAuto)
	for _, attnum := range attnums {
		if attnum != 0 { // 0 = expression column, no column dep
			c.recordDependency('i', idx.OID, 0, 'r', rel.OID, int32(attnum), DepAuto)
		}
	}

	return nil
}

// createIndexInternal creates an index internally (shared by CreateIndex and constraint auto-creation).
func (c *Catalog) createIndexInternal(schema *Schema, rel *Relation, name string, attnums []int16, isUnique, isPrimary bool, constraintOID uint32) *Index {
	// Default IndOption: all zeros (ASC NULLS LAST for btree).
	indOption := make([]int16, len(attnums))

	idx := &Index{
		OID:           c.oidGen.Next(),
		Name:          name,
		Schema:        schema,
		RelOID:        rel.OID,
		Columns:       attnums,
		IsUnique:      isUnique,
		IsPrimary:     isPrimary,
		ConstraintOID: constraintOID,
		AccessMethod:  "btree",
		NKeyColumns:   len(attnums),
		IndOption:     indOption,
	}

	c.registerIndex(schema, idx)
	// pg: index_create() — whole-relation dep + per-column deps
	c.recordDependency('i', idx.OID, 0, 'r', rel.OID, 0, DepAuto)
	for _, attnum := range attnums {
		if attnum != 0 {
			c.recordDependency('i', idx.OID, 0, 'r', rel.OID, int32(attnum), DepAuto)
		}
	}

	return idx
}

// registerIndex adds an index to all catalog maps.
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
}

// removeIndexesForRelation removes all indexes belonging to a relation.
func (c *Catalog) removeIndexesForRelation(relOID uint32, schema *Schema) {
	for _, idx := range c.indexesByRel[relOID] {
		delete(schema.Indexes, idx.Name)
		delete(c.indexes, idx.OID)
		c.removeComments('i', idx.OID)
		c.removeDepsOf('i', idx.OID)
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
