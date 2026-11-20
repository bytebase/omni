package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// Constraint represents a table constraint.
//
// pg: src/include/catalog/pg_constraint.h
type Constraint struct {
	OID            uint32
	Name           string
	Type           ConstraintType
	RelOID         uint32
	Namespace      uint32       // connamespace (schema OID)
	Columns        []int16      // key attnums
	IncludeColumns []int16      // PK/UNIQUE/EXCLUDE INCLUDE attnums
	FRelOID        uint32       // FK: referenced relation OID
	FColumns       []int16      // FK: referenced attnums
	FKUpdAction    byte         // FK: 'a'=NO ACTION, 'r'=RESTRICT, 'c'=CASCADE, 'n'=SET NULL, 'd'=SET DEFAULT
	FKDelAction    byte         // FK: same codes
	FKMatchType    byte         // FK: 's'=SIMPLE, 'f'=FULL, 'p'=PARTIAL
	Deferrable     bool         // condeferrable
	Deferred       bool         // condeferred
	Validated      bool         // convalidated
	CheckExpr      string       // CHECK: opaque expression
	CheckAnalyzed  AnalyzedExpr // CHECK: analyzed form (for Tier 2 deparse)
	IndexOID       uint32       // PK/UNIQUE/EXCLUDE: backing index OID
	ExclOps        []string     // EXCLUDE only: operator names per column

	// Inheritance/partition fields.
	// pg: src/include/catalog/pg_constraint.h lines 89-108
	ConParentID  uint32 // parent constraint OID (partitions)
	ConIsLocal   bool   // locally defined (not only inherited)
	ConInhCount  int16  // inheritance count
	ConNoInherit bool   // cannot be inherited

	// FK operator arrays.
	// pg: src/include/catalog/pg_constraint.h lines 123-145
	PFEqOp       []uint32 // FK: PK=FK equality operators
	PPEqOp       []uint32 // FK: PK=PK equality operators
	FFEqOp       []uint32 // FK: FK=FK equality operators
	FKDelSetCols []int16  // FK: subset of columns for ON DELETE SET NULL/DEFAULT
}

// AddConstraint adds a constraint to an existing relation, bypassing DDL text.
func (c *Catalog) AddConstraint(schemaName, tableName string, def ConstraintDef) error {
	schema, rel, err := c.findRelation(schemaName, tableName)
	if err != nil {
		return err
	}
	return c.addConstraint(schema, rel, def)
}

// addConstraint adds a constraint to a relation during CREATE TABLE.
func (c *Catalog) addConstraint(schema *Schema, rel *Relation, def ConstraintDef) error {
	switch def.Type {
	case ConstraintPK:
		return c.addPKConstraint(schema, rel, def)
	case ConstraintUnique:
		return c.addUniqueConstraint(schema, rel, def)
	case ConstraintFK:
		return c.addFKConstraint(schema, rel, def)
	case ConstraintCheck:
		return c.addCheckConstraint(rel, def)
	case ConstraintExclude:
		return c.addExcludeConstraint(schema, rel, def)
	default:
		return fmt.Errorf("unsupported constraint type %c", def.Type)
	}
}

func (c *Catalog) addPKConstraint(schema *Schema, rel *Relation, def ConstraintDef) error {
	attnums, existingIdx, err := c.resolveIndexConstraintAttnums(schema, rel, def)
	if err != nil {
		return err
	}

	// Verify no existing PK on this table.
	for _, con := range c.consByRel[rel.OID] {
		if con.Type == ConstraintPK {
			reused, err := c.reusePartitionIndexConstraint(schema, rel, con, def, attnums, existingIdx)
			if reused || err != nil {
				return err
			}
			return errDuplicatePKey(rel.Name)
		}
	}

	return c.addIndexBackedConstraint(schema, rel, def, ConstraintPK, attnums, existingIdx)
}

func (c *Catalog) addUniqueConstraint(schema *Schema, rel *Relation, def ConstraintDef) error {
	attnums, existingIdx, err := c.resolveIndexConstraintAttnums(schema, rel, def)
	if err != nil {
		return err
	}

	for _, con := range c.consByRel[rel.OID] {
		if con.Type != ConstraintUnique {
			continue
		}
		reused, err := c.reusePartitionIndexConstraint(schema, rel, con, def, attnums, existingIdx)
		if reused || err != nil {
			return err
		}
	}

	return c.addIndexBackedConstraint(schema, rel, def, ConstraintUnique, attnums, existingIdx)
}

func (c *Catalog) addIndexBackedConstraint(schema *Schema, rel *Relation, def ConstraintDef, typ ConstraintType, attnums []int16, existingIdx *Index) error {
	var includeAttnums []int16
	if existingIdx == nil && len(def.IncludeColumns) > 0 {
		// pg: src/backend/parser/parse_utilcmd.c — transformIndexConstraint
		// Checks duplicates within cstmt->including (indexIncludingParams) separately
		// from cstmt->keys (indexParams); key+INCLUDE overlap such as UNIQUE (a) INCLUDE (a)
		// is accepted by PostgreSQL 17 (see testdata/simulate.tsv:701, TestSimulation.tape:17979).
		seenInc := make(map[string]bool, len(def.IncludeColumns))
		for _, incName := range def.IncludeColumns {
			if seenInc[incName] {
				return errDuplicateColumn(incName)
			}
			seenInc[incName] = true
		}
		var err error
		includeAttnums, err = c.resolveColumnNames(rel, def.IncludeColumns)
		if err != nil {
			return err
		}
	} else if existingIdx != nil && existingIdx.NKeyColumns > 0 && existingIdx.NKeyColumns < len(existingIdx.Columns) {
		includeAttnums = append([]int16(nil), existingIdx.Columns[existingIdx.NKeyColumns:]...)
	}

	name := def.Name
	if name == "" {
		if existingIdx != nil {
			name = existingIdx.Name
		} else {
			allRawCols := make([]string, 0, len(def.Columns)+len(def.IncludeColumns))
			allRawCols = append(allRawCols, def.Columns...)
			allRawCols = append(allRawCols, def.IncludeColumns...)
			chosenCols := chooseIndexColumnNames(allRawCols)
			name = c.chooseConstraintAndIndexName(schema, rel, chosenCols, typ)
		}
	}

	effectiveTablespace := c.effectiveIndexTablespace(def.IndexSpace, rel)
	if existingIdx == nil {
		if err := c.checkIndexTablespace(effectiveTablespace, rel); err != nil {
			return err
		}
	}

	conOID := c.oidGen.Next()

	isPrimary := typ == ConstraintPK
	if isPrimary {
		for _, attnum := range attnums {
			rel.Columns[attnum-1].NotNull = true
		}
	}

	idx := existingIdx
	if idx == nil {
		idxName := name
		if idxName == "" {
			idxName = generateIndexName(rel.Name, def.Columns, isPrimary)
		}
		allAttnums := make([]int16, 0, len(attnums)+len(includeAttnums))
		allAttnums = append(allAttnums, attnums...)
		allAttnums = append(allAttnums, includeAttnums...)
		idx = c.createIndexInternal(schema, rel, idxName, allAttnums, len(attnums), true, isPrimary, conOID)
		idx.Tablespace = normalizeTablespace(effectiveTablespace)
	} else {
		idx.IsUnique = true
		idx.IsPrimary = isPrimary
		idx.ConstraintOID = conOID
		if err := c.renameAttachedIndex(schema, idx, name); err != nil {
			return err
		}
	}

	con := &Constraint{
		OID:            conOID,
		Name:           name,
		Type:           typ,
		RelOID:         rel.OID,
		Namespace:      schema.OID,
		Columns:        attnums,
		IncludeColumns: includeAttnums,
		IndexOID:       idx.OID,
		Deferrable:     def.Deferrable,
		Deferred:       def.Deferred,
		Validated:      true,
		ConIsLocal:     true,
	}

	c.registerConstraint(rel.OID, con)
	// Constraint depends on owning relation columns (auto-dropped).
	// pg: pg_constraint.c — CreateConstraintEntry lines 254-268
	for _, attnum := range attnums {
		c.recordDependency('c', con.OID, 0, 'r', rel.OID, int32(attnum), DepAuto)
	}
	for _, attnum := range includeAttnums {
		c.recordDependency('c', con.OID, 0, 'r', rel.OID, int32(attnum), DepAuto)
	}
	// Index depends on its constraint (internal).
	c.recordDependency('i', idx.OID, 0, 'c', con.OID, 0, DepInternal)

	return nil
}

func (c *Catalog) resolveIndexConstraintAttnums(schema *Schema, rel *Relation, def ConstraintDef) ([]int16, *Index, error) {
	if def.IndexName == "" {
		seen := make(map[string]bool, len(def.Columns))
		for _, name := range def.Columns {
			if seen[name] {
				return nil, nil, errDuplicateColumn(name)
			}
			seen[name] = true
		}
		attnums, err := c.resolveColumnNames(rel, def.Columns)
		return attnums, nil, err
	}

	idx, err := c.findNamedIndex(schema, def.IndexName)
	if err != nil {
		return nil, nil, err
	}
	if idx.RelOID != rel.OID {
		return nil, nil, errInvalidObjectDefinition(fmt.Sprintf("index %q does not belong to table %q", def.IndexName, rel.Name))
	}
	if !idx.IsUnique {
		return nil, nil, errInvalidObjectDefinition(fmt.Sprintf("index %q is not unique", def.IndexName))
	}
	attnums := indexConstraintAttnums(idx)
	if len(attnums) == 0 {
		return nil, nil, errInvalidParameterValue("constraint index must specify at least one column")
	}
	for _, attnum := range attnums {
		if attnum == 0 {
			return nil, nil, errInvalidObjectDefinition("constraints cannot use expression indexes")
		}
	}
	return attnums, idx, nil
}

func indexConstraintAttnums(idx *Index) []int16 {
	if idx == nil {
		return nil
	}
	n := idx.NKeyColumns
	if n <= 0 || n > len(idx.Columns) {
		n = len(idx.Columns)
	}
	out := make([]int16, n)
	copy(out, idx.Columns[:n])
	return out
}

func (c *Catalog) reusePartitionIndexConstraint(schema *Schema, rel *Relation, con *Constraint, def ConstraintDef, attnums []int16, existingIdx *Index) (bool, error) {
	if rel.PartitionOf == 0 || con.ConParentID == 0 {
		return false, nil
	}
	if !sameAttnums(con.Columns, attnums) {
		return false, nil
	}
	if existingIdx != nil && con.IndexOID != existingIdx.OID {
		return false, nil
	}

	name := def.Name
	if name == "" && existingIdx != nil {
		name = existingIdx.Name
	}
	if name != "" {
		for _, existing := range c.consByRel[rel.OID] {
			if existing.OID != con.OID && existing.Name == name {
				return true, errDuplicateObject("constraint", name)
			}
		}
		con.Name = name
	}
	con.ConIsLocal = true
	con.Validated = true

	if con.IndexOID != 0 {
		idx := c.indexes[con.IndexOID]
		if idx != nil {
			idx.IsPrimary = con.Type == ConstraintPK
			idx.IsUnique = true
			idx.ConstraintOID = con.OID

			idxName := name
			if idxName != "" {
				if err := c.renameAttachedIndex(schema, idx, idxName); err != nil {
					return true, err
				}
			}
		}
	}

	// The cloned partition PK already enforced this in the normal path, but
	// keep the invariant here for partitions loaded from partial metadata.
	for _, attnum := range attnums {
		if int(attnum) > 0 && int(attnum) <= len(rel.Columns) {
			rel.Columns[attnum-1].NotNull = true
		}
	}

	return true, nil
}

func (c *Catalog) renameAttachedIndex(schema *Schema, idx *Index, name string) error {
	if idx == nil || name == "" || idx.Name == name {
		return nil
	}
	if existing := schema.Indexes[name]; existing != nil && existing.OID != idx.OID {
		return errDuplicateObject("index", name)
	}
	delete(schema.Indexes, idx.Name)
	idx.Name = name
	schema.Indexes[idx.Name] = idx
	return nil
}

func sameAttnums(a, b []int16) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (c *Catalog) addFKConstraint(schema *Schema, rel *Relation, def ConstraintDef) error {
	// Resolve local columns.
	localAttnums, err := c.resolveColumnNames(rel, def.Columns)
	if err != nil {
		return err
	}

	// Find referenced table.
	_, refRel, err := c.findRelation(def.RefSchema, def.RefTable)
	if err != nil {
		return err
	}

	// Persistence cross-reference: permanent tables cannot reference temporary tables.
	// pg: src/backend/commands/tablecmds.c — ATAddForeignKeyConstraint (persistence check)
	relPersistence := rel.Persistence
	if relPersistence == 0 {
		relPersistence = 'p'
	}
	refPersistence := refRel.Persistence
	if refPersistence == 0 {
		refPersistence = 'p'
	}
	if relPersistence != 't' && refPersistence == 't' {
		return errInvalidFK("constraints on permanent tables may reference only permanent tables")
	}

	// Resolve referenced columns (or use PK if RefColumns empty).
	var refAttnums []int16
	if len(def.RefColumns) == 0 {
		// Find PK of referenced table.
		pk := c.findPKConstraint(refRel.OID)
		if pk == nil {
			return errInvalidFK(fmt.Sprintf("there is no primary key for referenced table %q", def.RefTable))
		}
		refAttnums = pk.Columns
		if len(localAttnums) != len(refAttnums) {
			return errInvalidFK("number of referencing and referenced columns for foreign key disagree")
		}
	} else {
		refAttnums, err = c.resolveColumnNames(refRel, def.RefColumns)
		if err != nil {
			return err
		}
		if len(localAttnums) != len(refAttnums) {
			return errInvalidFK("number of referencing and referenced columns for foreign key disagree")
		}
	}

	// Verify PK/UNIQUE constraint or non-partial unique index on referenced columns.
	// pg: src/backend/commands/tablecmds.c — transformFkeyCheckAttrs
	// PostgreSQL accepts a unique index as a FK referential target even when
	// no explicit UNIQUE constraint exists. pg_dump relies on this and emits
	// CREATE UNIQUE INDEX as a separate statement after CREATE TABLE.
	if c.lookupFKSupportingIndex(refRel.OID, refAttnums) == 0 {
		return errInvalidFK(fmt.Sprintf("there is no unique constraint matching given keys for referenced table %q", def.RefTable))
	}

	// Verify type compatibility.
	for i := range localAttnums {
		localCol := rel.Columns[localAttnums[i]-1]
		refCol := refRel.Columns[refAttnums[i]-1]
		if c.findFKEqOp(refCol.TypeOID, localCol.TypeOID) == 0 {
			localType := c.typeByOID[localCol.TypeOID]
			refType := c.typeByOID[refCol.TypeOID]
			localName := "unknown"
			refName := "unknown"
			if localType != nil {
				localName = localType.TypeName
			}
			if refType != nil {
				refName = refType.TypeName
			}
			return errDatatypeMismatch(fmt.Sprintf(
				"foreign key constraint %q: column %q is type %s but references column %q of type %s",
				def.Name, rel.Columns[localAttnums[i]-1].Name, localName,
				refRel.Columns[refAttnums[i]-1].Name, refName,
			))
		}
	}

	name := def.Name
	if name == "" {
		name = generateConstraintName(rel.Name, def.Columns, ConstraintFK)
	}

	// Normalize FK action/match defaults.
	fkUpd := def.FKUpdAction
	if fkUpd == 0 {
		fkUpd = 'a'
	}
	fkDel := def.FKDelAction
	if fkDel == 0 {
		fkDel = 'a'
	}
	fkMatch := def.FKMatchType
	if fkMatch == 0 {
		fkMatch = 's'
	}

	// Build FK equality operator arrays.
	// pg: src/backend/commands/tablecmds.c — ATAddForeignKeyConstraint
	var pfEqOp, ppEqOp, ffEqOp []uint32
	for i := range localAttnums {
		localCol := rel.Columns[localAttnums[i]-1]
		refCol := refRel.Columns[refAttnums[i]-1]
		pfOp := c.findFKEqOp(refCol.TypeOID, localCol.TypeOID)
		ppOp := c.findEqualityOp(refCol.TypeOID, refCol.TypeOID)
		ffOp := c.findEqualityOp(localCol.TypeOID, localCol.TypeOID)
		pfEqOp = append(pfEqOp, pfOp)
		ppEqOp = append(ppEqOp, ppOp)
		ffEqOp = append(ffEqOp, ffOp)
	}

	con := &Constraint{
		OID:         c.oidGen.Next(),
		Name:        name,
		Type:        ConstraintFK,
		RelOID:      rel.OID,
		Namespace:   schema.OID,
		Columns:     localAttnums,
		FRelOID:     refRel.OID,
		FColumns:    refAttnums,
		FKUpdAction: fkUpd,
		FKDelAction: fkDel,
		FKMatchType: fkMatch,
		Deferrable:  def.Deferrable,
		Deferred:    def.Deferred,
		Validated:   !def.SkipValidation,
		ConIsLocal:  true,
		PFEqOp:      pfEqOp,
		PPEqOp:      ppEqOp,
		FFEqOp:      ffEqOp,
	}

	c.registerConstraint(rel.OID, con)
	// pg: pg_constraint.c — CreateConstraintEntry: local columns
	for _, attnum := range localAttnums {
		c.recordDependency('c', con.OID, 0, 'r', rel.OID, int32(attnum), DepAuto)
	}
	// FK depends on referenced table columns (normal dep: blocks DROP without CASCADE).
	// pg: pg_constraint.c — CreateConstraintEntry: FK referenced columns
	for _, attnum := range refAttnums {
		c.recordDependency('c', con.OID, 0, 'r', refRel.OID, int32(attnum), DepNormal)
	}

	// FK depends on the supporting unique/PK index of the referenced table.
	// pg: src/backend/commands/tablecmds.c — ATAddForeignKeyConstraint
	if refIdxOID := c.lookupFKSupportingIndex(refRel.OID, refAttnums); refIdxOID != 0 {
		c.recordDependency('c', con.OID, 0, 'r', refIdxOID, 0, DepNormal)
	}
	seenFKOp := map[uint32]bool{}
	for _, opGroups := range [][]uint32{pfEqOp, ppEqOp, ffEqOp} {
		for _, opOID := range opGroups {
			if opOID >= FirstNormalObjectId && !seenFKOp[opOID] {
				seenFKOp[opOID] = true
				c.recordDependency('c', con.OID, 0, 'O', opOID, 0, DepNormal)
			}
		}
	}

	return nil
}

func (c *Catalog) addCheckConstraint(rel *Relation, def ConstraintDef) error {
	// Analyze the CHECK expression if we have the raw AST node.
	// pg: src/backend/commands/tablecmds.c — cookConstraint (CHECK analysis)
	var analyzed AnalyzedExpr
	if def.RawCheckExpr != nil {
		if rawExprContainsSubLink(def.RawCheckExpr) {
			return &Error{Code: CodeFeatureNotSupported,
				Message: "cannot use subquery in check constraint"}
		}
		a, err := c.AnalyzeStandaloneExpr(def.RawCheckExpr, rel)
		if err != nil {
			if !c.installingSnapshot {
				return err
			}
		} else if a != nil {
			coerced, cerr := c.coerceToBoolean(a, "CHECK")
			if cerr != nil {
				if !c.installingSnapshot {
					return cerr
				}
			} else {
				analyzed = coerced
			}
		}
	}

	// Extract referenced columns for conkey tracking: the columns the
	// analyzed expression reads, or a text scan when there is none.
	// pg: src/backend/catalog/heap.c — StoreRelCheck (pull_var_clause)
	var cols []int16
	if analyzed != nil {
		cols = c.exprColumns(analyzed, rel.OID)
	} else {
		cols = extractCheckExprColumns(rel, def.CheckExpr)
	}

	// pg: src/backend/catalog/heap.c — ChooseConstraintName
	// PG includes column name when exactly one column is referenced:
	//   1 column: {table}_{col}_check
	//   0 or 2+ columns: {table}_check
	// With numeric suffix for uniqueness: _check, _check1, _check2, etc.
	name := def.Name
	if name == "" {
		baseName := rel.Name + "_check"
		if len(cols) == 1 {
			// Find the single referenced column name.
			for _, col := range rel.Columns {
				if col.AttNum == cols[0] {
					baseName = rel.Name + "_" + col.Name + "_check"
					break
				}
			}
		}
		name = baseName
		// Deduplicate: check for existing constraints with the same name.
		for suffix := 1; ; suffix++ {
			conflict := false
			for _, con := range c.consByRel[rel.OID] {
				if con.Name == name {
					conflict = true
					break
				}
			}
			if !conflict {
				break
			}
			name = fmt.Sprintf("%s%d", baseName, suffix)
		}
	}

	con := &Constraint{
		OID:        c.oidGen.Next(),
		Name:       name,
		Type:       ConstraintCheck,
		RelOID:     rel.OID,
		CheckExpr:  def.CheckExpr,
		Columns:    cols,
		Validated:  !def.SkipValidation,
		ConIsLocal: true,
	}

	if analyzed != nil {
		con.CheckAnalyzed = analyzed
		rte := c.buildRelationRTE(rel)
		con.CheckExpr = c.DeparseExpr(analyzed, []*RangeTableEntry{rte}, false)
	}

	c.registerConstraint(rel.OID, con)

	// pg: src/backend/catalog/pg_constraint.c — CreateConstraintEntry lines 254-268
	if len(con.Columns) > 0 {
		for _, attnum := range con.Columns {
			c.recordDependency('c', con.OID, 0, 'r', rel.OID, int32(attnum), DepAuto)
		}
	} else {
		c.recordDependency('c', con.OID, 0, 'r', rel.OID, 0, DepAuto)
	}

	// pg: src/backend/catalog/pg_constraint.c:368-376
	// Lines 374-376: recordDependencyOnSingleRelExpr records deps on objects
	// in the CHECK expression. Column refs to the same relation use selfBehavior
	// (DepNormal); external refs (functions, operators) also use DepNormal.
	if con.CheckAnalyzed != nil {
		c.recordDependencyOnSingleRelExpr('c', con.OID, con.CheckAnalyzed, rel.OID,
			DepNormal, DepNormal)
	}

	return nil
}

// extractCheckExprColumns finds column references in a CHECK expression text.
// Returns sorted attnums of columns whose names appear in the expression.
//
// pg: recordDependencyOnSingleRelExpr walks the expression tree;
// since we store text, we scan for column name matches.
func extractCheckExprColumns(rel *Relation, exprText string) []int16 {
	if exprText == "" {
		return nil
	}
	var cols []int16
	for _, col := range rel.Columns {
		if strings.Contains(exprText, col.Name) {
			cols = append(cols, col.AttNum)
		}
	}
	sort.Slice(cols, func(i, j int) bool { return cols[i] < cols[j] })
	return cols
}

// addExcludeConstraint adds an EXCLUDE constraint with a backing index.
//
// pg: src/backend/commands/tablecmds.c — DefineRelation (exclusion constraint handling via DefineIndex)
// opClassRefsFrom turns stored operator-class names back into references with
// an unqualified schema and no type, for objects built without an IndexElem
// list.
func opClassRefsFrom(opClasses []string) []opClassRef {
	refs := make([]opClassRef, 0, len(opClasses))
	for _, name := range opClasses {
		if name == "" {
			refs = append(refs, opClassRef{})
			continue
		}
		schema, refName := "", name
		if i := strings.LastIndexByte(name, '.'); i >= 0 {
			schema, refName = name[:i], name[i+1:]
		}
		refs = append(refs, opClassRef{schema: schema, name: refName})
	}
	return refs
}

func (c *Catalog) addExcludeConstraint(schema *Schema, rel *Relation, def ConstraintDef) error {
	attnums, err := c.resolveColumnNames(rel, def.Columns)
	if err != nil {
		return err
	}
	var includeAttnums []int16
	if len(def.IncludeColumns) > 0 {
		seenInc := make(map[string]bool, len(def.IncludeColumns))
		for _, incName := range def.IncludeColumns {
			if seenInc[incName] {
				return errDuplicateColumn(incName)
			}
			seenInc[incName] = true
		}
		includeAttnums, err = c.resolveColumnNames(rel, def.IncludeColumns)
		if err != nil {
			return err
		}
	}

	// Determine access method (default GiST for exclusion constraints).
	accessMethod := def.AccessMethod
	if accessMethod == "" {
		accessMethod = "gist"
	}

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
	// Validate that built-in AMs support amgettuple (required for exclusion constraints).
	// Only btree, hash, gist, spgist support it; gin and brin do not.
	// pg: src/backend/commands/indexcmds.c — DefineIndex (line 875-879)
	if am.OID < FirstNormalObjectId {
		switch accessMethod {
		case "btree", "hash", "gist", "spgist":
		default:
			return &Error{Code: CodeFeatureNotSupported,
				Message: fmt.Sprintf("access method %q does not support exclusion constraints", accessMethod)}
		}
	}

	effectiveTablespace := c.effectiveIndexTablespace(def.IndexSpace, rel)
	if err := c.checkIndexTablespace(effectiveTablespace, rel); err != nil {
		return err
	}

	// Validate exclusion operator classes and operators before registering the
	// backing index or constraint.
	// pg: src/backend/commands/indexcmds.c — DefineIndex (ResolveIndexExclusionOperator)
	exclOpClassRefs := c.indexOpClassRefs(rel, attnums, len(attnums), def.ExclOpClasses)
	resolvedExclOpOIDs := make([]uint32, len(def.ExclOps))
	for i, opName := range def.ExclOps {
		if i < len(exclOpClassRefs) && exclOpClassRefs[i].name != "" {
			if err := c.validateIndexOpClass(exclOpClassRefs[i], am, nil); err != nil {
				return err
			}
		}
		if opName == "" || i >= len(attnums) {
			continue
		}
		ci := indexOfAttNum(rel, attnums[i])
		if ci < 0 {
			continue
		}
		colType := rel.Columns[ci].TypeOID
		opSchema := ""
		if i < len(def.ExclOpSchemas) {
			opSchema = def.ExclOpSchemas[i]
		}
		if i < len(exclOpClassRefs) && opSchema == "" {
			if opcOID := c.resolveOpClassOID(exclOpClassRefs[i], am.OID); opcOID >= FirstNormalObjectId {
				if opc := c.opClasses[opcOID]; opc != nil && opc.Operators[opName] {
					if opOID := c.resolveVisibleOperatorOID(opName, colType, colType); opOID != 0 {
						resolvedExclOpOIDs[i] = opOID
					}
					continue
				}
			}
		}
		opOID, err := c.resolveExclusionOperatorOID(opSchema, opName, colType)
		if err != nil {
			return err
		}
		resolvedExclOpOIDs[i] = opOID
	}

	name := def.Name
	allRawCols := make([]string, 0, len(def.Columns)+len(def.IncludeColumns))
	allRawCols = append(allRawCols, def.Columns...)
	allRawCols = append(allRawCols, def.IncludeColumns...)
	chosenCols := chooseIndexColumnNames(allRawCols)
	if name == "" {
		name = c.chooseConstraintAndIndexName(schema, rel, chosenCols, ConstraintExclude)
	}

	conOID := c.oidGen.Next()

	// Create backing index.
	idxName := name
	indOption := make([]int16, len(attnums))
	var storedOpClasses []string
	for _, oc := range def.ExclOpClasses {
		if oc != "" {
			storedOpClasses = append([]string(nil), def.ExclOpClasses...)
			break
		}
	}
	allAttnums := make([]int16, 0, len(attnums)+len(includeAttnums))
	allAttnums = append(allAttnums, attnums...)
	allAttnums = append(allAttnums, includeAttnums...)
	idx := &Index{
		OID:           c.oidGen.Next(),
		Name:          idxName,
		Schema:        schema,
		RelOID:        rel.OID,
		Columns:       allAttnums,
		ColumnNames:   chosenCols,
		IsUnique:      false, // EXCLUDE indexes are not unique
		IsPrimary:     false,
		ConstraintOID: conOID,
		AccessMethod:  accessMethod,
		Tablespace:    normalizeTablespace(effectiveTablespace),
		OwnerName:     rel.OwnerName,
		NKeyColumns:   len(attnums),
		IndOption:     indOption,
		OpClasses:     storedOpClasses,
	}
	c.registerIndex(schema, idx)
	c.recordDependency('i', idx.OID, 0, 'r', rel.OID, 0, DepAuto)
	c.recordColumnDependencies('i', idx.OID, rel.OID, allAttnums, DepAuto)
	if c.accessMethods[am.OID] != nil {
		c.recordDependency('i', idx.OID, 0, 'm', am.OID, 0, DepNormal)
	}
	c.recordOpClassDependencies('i', idx.OID, c.indexOpClassRefs(rel, allAttnums, len(attnums), def.ExclOpClasses), idx.AccessMethod)
	for _, opOID := range resolvedExclOpOIDs {
		if opOID >= FirstNormalObjectId {
			c.recordDependency('c', conOID, 0, 'O', opOID, 0, DepNormal)
		}
	}

	con := &Constraint{
		OID:            conOID,
		Name:           name,
		Type:           ConstraintExclude,
		RelOID:         rel.OID,
		Namespace:      schema.OID,
		Columns:        attnums,
		IncludeColumns: includeAttnums,
		IndexOID:       idx.OID,
		ExclOps:        def.ExclOps,
		Deferrable:     def.Deferrable,
		Deferred:       def.Deferred,
		Validated:      true,
		ConIsLocal:     true,
	}

	c.registerConstraint(rel.OID, con)
	// pg: pg_constraint.c — CreateConstraintEntry lines 254-268
	for _, attnum := range allAttnums {
		c.recordDependency('c', con.OID, 0, 'r', rel.OID, int32(attnum), DepAuto)
	}
	c.recordDependency('i', idx.OID, 0, 'c', con.OID, 0, DepInternal)

	return nil
}

// resolveExclusionOperatorOID resolves a binary exclusion operator for (colType, colType).
//
// pg: src/backend/parser/parse_oper.c — compatible_oper
func (c *Catalog) resolveExclusionOperatorOID(opSchema, opName string, colType uint32) (uint32, error) {
	if opSchema != "" {
		s := c.schemaByName[opSchema]
		if s == nil {
			return 0, errUndefinedSchema(opSchema)
		}
		if c.snapshot || c.installingSnapshot {
			return 0, nil
		}
		for _, op := range c.operByKey[operKey{name: opName, left: colType, right: colType}] {
			opNS := uint32(PGCatalogNamespace)
			if op.OID >= FirstNormalObjectId {
				opNS = c.operatorSchemas[op.OID]
			}
			if opNS == s.OID {
				return op.OID, nil
			}
		}
		var coercible []funcCandidate
		inputTypes := []uint32{colType, colType}
		for key, ops := range c.operByKey {
			if key.name != opName {
				continue
			}
			for _, op := range ops {
				if op.Kind == 'l' {
					continue
				}
				opNS := uint32(PGCatalogNamespace)
				if op.OID >= FirstNormalObjectId {
					opNS = c.operatorSchemas[op.OID]
				}
				if opNS == s.OID && c.canCoerceArgs(inputTypes, []uint32{op.Left, op.Right}) {
					coercible = append(coercible, funcCandidate{argTypes: []uint32{op.Left, op.Right}, op: op})
				}
			}
		}
		if len(coercible) == 1 {
			return coercible[0].op.OID, nil
		}
		if len(coercible) > 1 {
			if best := c.funcSelectCandidate(inputTypes, coercible); best != nil {
				return best.op.OID, nil
			}
		}
		return 0, errUndefinedFunction(opName, []uint32{colType, colType})
	}

	if opOID := c.resolveVisibleOperatorOID(opName, colType, colType); opOID != 0 {
		return opOID, nil
	}
	if c.snapshot || c.installingSnapshot {
		return 0, nil
	}
	op, _, _, err := c.resolveOpWithCoercionFull(opName, colType, colType)
	if err != nil {
		return 0, err
	}
	return op.OID, nil
}

// chooseConstraintAndIndexName picks a free name in both the relation/index
// namespace and the relation's constraint list, matching ChooseRelationName.
func (c *Catalog) chooseConstraintAndIndexName(schema *Schema, rel *Relation, columns []string, ctype ConstraintType) string {
	var name2, label string
	switch ctype {
	case ConstraintPK:
		label = "pkey"
	case ConstraintUnique:
		name2 = strings.Join(columns, "_")
		label = "key"
	case ConstraintExclude:
		name2 = strings.Join(columns, "_")
		label = "excl"
	default:
		return generateConstraintName(rel.Name, columns, ctype)
	}
	base := rel.Name
	if name2 != "" {
		base += "_" + name2
	}
	if label != "" {
		base += "_" + label
	}
	cand := base
	for k := 1; ; k++ {
		_, relHeld := relationNamespaceKind(schema, cand)
		conHeld := false
		for _, con := range c.consByRel[rel.OID] {
			if con.Name == cand {
				conHeld = true
				break
			}
		}
		if !relHeld && !conHeld {
			return cand
		}
		cand = fmt.Sprintf("%s%d", base, k)
	}
}

func (c *Catalog) findExactEqualityOp(leftType, rightType uint32) uint32 {
	ops := c.LookupOperatorExact("=", leftType, rightType)
	if len(ops) > 0 {
		return ops[0].OID
	}
	return 0
}

func (c *Catalog) findFKEqOp(refType, localType uint32) uint32 {
	if op := c.findExactEqualityOp(refType, localType); op != 0 {
		return op
	}

	refBase := c.getBaseType(refType)
	localBase := c.getBaseType(localType)
	if refBase != refType || localBase != localType {
		if op := c.findExactEqualityOp(refBase, localBase); op != 0 {
			return op
		}
	}

	// PG can implement a FK when the referencing column can be implicitly
	// coerced to the referenced column's equality type, even without a direct
	// cross-type equality operator (for example varchar -> text or int4 -> numeric).
	if c.CanCoerce(localBase, refBase, 'i') {
		if op := c.findExactEqualityOp(refBase, refBase); op != 0 {
			return op
		}
	}

	refBT := c.typeByOID[refBase]
	localBT := c.typeByOID[localBase]
	if refBT != nil && localBT != nil && refBT.Category == 'S' && localBT.Category == 'S' {
		if op := c.findExactEqualityOp(TEXTOID, TEXTOID); op != 0 {
			return op
		}
	}

	return 0
}

// findEqualityOp looks up the equality operator for two types.
// Returns the operator OID, or 0 if not found.
//
// (pgddl helper — PG uses lookup_type_cache for this)
func (c *Catalog) findEqualityOp(leftType, rightType uint32) uint32 {
	ops := c.LookupOperatorExact("=", leftType, rightType)
	if len(ops) > 0 {
		return ops[0].OID
	}
	// Try same-type if cross-type not found.
	if leftType != rightType {
		ops = c.LookupOperatorExact("=", leftType, leftType)
		if len(ops) > 0 {
			return ops[0].OID
		}
	}
	return 0
}

// registerConstraint adds a constraint to catalog maps.
func (c *Catalog) registerConstraint(relOID uint32, con *Constraint) {
	c.constraints[con.OID] = con
	c.consByRel[relOID] = append(c.consByRel[relOID], con)
}

// resolveVisibleOperatorOID finds the operator matching (name, left, right)
// in search-path order, returning 0 when none matches.
func (c *Catalog) resolveVisibleOperatorOID(name string, left, right uint32) uint32 {
	ops := c.operByKey[operKey{name: name, left: left, right: right}]
	if len(ops) == 0 {
		return 0
	}
	for _, nsOID := range c.searchPathWithCatalog() {
		for _, op := range ops {
			opNS := uint32(PGCatalogNamespace)
			if op.OID >= FirstNormalObjectId {
				opNS = c.operatorSchemas[op.OID]
			}
			if opNS == nsOID {
				return op.OID
			}
		}
	}
	return ops[0].OID
}

// removeConstraint removes a single constraint and its backing index (if any).
func (c *Catalog) removeConstraint(schema *Schema, con *Constraint) {
	// Remove backing index if present.
	if con.IndexOID != 0 {
		if idx := c.indexes[con.IndexOID]; idx != nil {
			c.removeIndex(schema, idx.Name, idx)
		}
	}

	delete(c.constraints, con.OID)
	c.removeComments('c', con.OID)
	c.removeDepsOf('c', con.OID)
	c.removeDepsOn('c', con.OID)

	// Remove from consByRel.
	list := c.consByRel[con.RelOID]
	for i, x := range list {
		if x.OID == con.OID {
			c.consByRel[con.RelOID] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

// removeConstraintsForRelation removes all constraints belonging to a relation.
func (c *Catalog) removeConstraintsForRelation(relOID uint32, schema *Schema) {
	for _, con := range c.consByRel[relOID] {
		if con.IndexOID != 0 {
			if idx := c.indexes[con.IndexOID]; idx != nil {
				delete(schema.Indexes, idx.Name)
				delete(c.indexes, idx.OID)
				c.removeComments('i', idx.OID)
				c.removeDepsOf('i', idx.OID)
				c.removeDepsOn('i', idx.OID)
				c.removeDepsOn('r', idx.OID)
			}
		}
		delete(c.constraints, con.OID)
		c.removeComments('c', con.OID)
		c.removeDepsOf('c', con.OID)
		c.removeDepsOn('c', con.OID)
	}
	delete(c.consByRel, relOID)
}

// findPKConstraint finds the PK constraint on a relation, or nil.
func (c *Catalog) findPKConstraint(relOID uint32) *Constraint {
	for _, con := range c.consByRel[relOID] {
		if con.Type == ConstraintPK {
			return con
		}
	}
	return nil
}

// lookupUniqueConstraint finds a PK/UNIQUE constraint covering exactly the given columns.
func (c *Catalog) lookupUniqueConstraint(relOID uint32, cols []int16) *Constraint {
	for _, con := range c.consByRel[relOID] {
		if con.Type != ConstraintPK && con.Type != ConstraintUnique {
			continue
		}
		if len(con.Columns) != len(cols) {
			continue
		}
		match := true
		for i := range cols {
			if con.Columns[i] != cols[i] {
				match = false
				break
			}
		}
		if match {
			return con
		}
	}
	return nil
}

// lookupFKSupportingIndex finds an index suitable to back a foreign-key
// reference to (relOID, cols). It returns the index OID, or 0 if none qualifies.
//
// Resolution order:
//  1. A PK or explicit UNIQUE constraint over exactly the referenced columns —
//     the constraint's backing index OID is returned.
//  2. A standalone unique index over exactly the referenced columns. PostgreSQL
//     accepts a unique index as a FK target if it is unique, non-partial,
//     contains no expression key columns, and has no INCLUDE columns. pg_dump
//     relies on this when it emits CREATE UNIQUE INDEX as a separate statement
//     after the owning CREATE TABLE.
//
// pg: src/backend/commands/tablecmds.c — transformFkeyCheckAttrs
func (c *Catalog) lookupFKSupportingIndex(relOID uint32, cols []int16) uint32 {
	if con := c.lookupUniqueConstraint(relOID, cols); con != nil {
		return con.IndexOID
	}

	for _, idx := range c.indexesByRel[relOID] {
		if !idx.IsUnique {
			continue
		}
		// Skip partial indexes — a partial unique index does not enforce
		// uniqueness over the whole relation.
		if idx.WhereClause != "" {
			continue
		}
		// Skip indexes with deparsed expression columns.
		if len(idx.Exprs) > 0 {
			continue
		}
		// Compare key columns only (excluding INCLUDE columns).
		nKey := idx.NKeyColumns
		if nKey == 0 {
			nKey = len(idx.Columns)
		}
		if nKey != len(cols) {
			continue
		}
		// Reject indexes that have any expression key columns (attnum=0).
		hasExprKey := false
		for i := 0; i < nKey; i++ {
			if idx.Columns[i] == 0 {
				hasExprKey = true
				break
			}
		}
		if hasExprKey {
			continue
		}
		match := true
		for i := 0; i < nKey; i++ {
			if idx.Columns[i] != cols[i] {
				match = false
				break
			}
		}
		if match {
			return idx.OID
		}
	}
	return 0
}

// dropDependents drops all objects that have a DepNormal dependency on the given referent.
// This is used for CASCADE behavior (e.g., dropping FK constraints from referencing tables).
// The referent's entries on itself (a generated column reading another column
// of its table) are skipped: the referent goes with them.
func (c *Catalog) dropDependents(refType byte, refOID uint32) {
	deps := c.findNormalDependents(refType, refOID)
	for _, dep := range deps {
		if dep.ObjType == refType && dep.ObjOID == refOID {
			continue
		}
		c.dropDependentObject(dep.ObjType, dep.ObjOID)
	}
}

// dropDependentObject drops one dependent found by dropDependents or
// relationDependents, with the objects that depend on it in turn.
func (c *Catalog) dropDependentObject(objType byte, objOID uint32) {
	switch objType {
	case 'c':
		con := c.constraints[objOID]
		if con == nil {
			return
		}
		c.dropDependents('c', con.OID)
		if con.IndexOID != 0 {
			c.dropDependents('r', con.IndexOID)
			c.dropDependents('i', con.IndexOID)
		}
		// Find the schema of the owning relation to remove indexes.
		if ownerRel := c.relationByOID[con.RelOID]; ownerRel != nil && c.constraints[con.OID] != nil {
			c.removeConstraint(ownerRel.Schema, con)
		}
	case 'r':
		rel := c.relationByOID[objOID]
		if rel == nil {
			return
		}
		c.cascadeDropRelation(rel)
	case 'i':
		idx := c.indexes[objOID]
		if idx == nil {
			return
		}
		// An index that backs a constraint goes with it.
		// pg: src/backend/catalog/dependency.c — findDependentObjects
		// (DEPENDENCY_INTERNAL: delete the owning object instead)
		if con := c.constraints[idx.ConstraintOID]; con != nil {
			c.dropDependentObject('c', con.OID)
			return
		}
		var childIdxOIDs []uint32
		for oid, child := range c.indexes {
			if child.ParentIndexOID == idx.OID && oid != idx.OID {
				childIdxOIDs = append(childIdxOIDs, oid)
			}
		}
		sort.Slice(childIdxOIDs, func(i, j int) bool { return childIdxOIDs[i] < childIdxOIDs[j] })
		for _, childOID := range childIdxOIDs {
			c.dropDependentObject('i', childOID)
		}
		c.dropDependents('r', idx.OID)
		c.dropDependents('i', idx.OID)
		if c.indexes[idx.OID] == idx {
			c.removeIndex(idx.Schema, idx.Name, idx)
		}
	case 's':
		seq := c.sequenceByOID[objOID]
		if seq == nil {
			return
		}
		c.dropDependents('s', seq.OID)
		if c.sequenceByOID[seq.OID] == seq {
			c.removeSequence(seq.Schema, seq)
		}
	case 'f':
		up := c.userProcs[objOID]
		if up == nil {
			return
		}
		c.dropDependents('f', up.OID)
		if c.userProcs[up.OID] != nil {
			c.removeFunction(up.OID, up.Name)
		}
	case 'g':
		trig := c.triggers[objOID]
		if trig == nil {
			return
		}
		c.removeTrigger(trig)
	case 'p':
		pol := c.policies[objOID]
		if pol == nil {
			return
		}
		c.removePolicy(pol)
	case 'R':
		rule := c.rules[objOID]
		if rule == nil {
			return
		}
		c.removeRule(c.relationByOID[rule.OwnerRelOID], rule)
	case 'E':
		evt := c.eventTriggers[objOID]
		if evt == nil {
			return
		}
		c.removeEventTrigger(evt)
	case 'O':
		if c.operByOID[objOID] != nil {
			c.dropDependents('O', objOID)
			c.removeOperator(objOID)
		}
	case 'm':
		c.cascadeDropAccessMethod(objOID)
	case 'q':
		c.cascadeDropOpClass(objOID)
	case 'Q':
		c.cascadeDropOpFamily(objOID)
	case 't':
		bt := c.typeByOID[objOID]
		if bt == nil {
			return
		}
		if bt.Type == 'c' && bt.RelID != 0 {
			if rel := c.relationByOID[bt.RelID]; rel != nil && rel.RelKind == 'c' {
				c.cascadeDropRelation(rel)
			}
			return
		}
		for _, relOID := range c.typeRelationDependents(bt) {
			if depRel := c.relationByOID[relOID]; depRel != nil {
				c.cascadeDropRelation(depRel)
			}
		}
		c.dropTypeDependents(bt)
		if c.typeByOID[bt.OID] != nil {
			c.deleteUserType(bt)
		}
	}
}

// cascadeDropRelation drops rel after recursively cascading through its
// auto/internal child relations, normal dependents, and row-type dependents.
func (c *Catalog) cascadeDropRelation(rel *Relation) {
	c.cascadeDropRelationSeen(rel, map[uint32]bool{})
}

func (c *Catalog) cascadeDropRelationSeen(rel *Relation, seen map[uint32]bool) {
	if rel == nil || c.relationByOID[rel.OID] == nil || seen[rel.OID] {
		return
	}
	seen[rel.OID] = true
	for _, d := range c.deps {
		if d.ObjType == 'r' && d.RefType == 'r' && d.RefOID == rel.OID && (d.DepType == DepAuto || d.DepType == DepInternal) {
			if child := c.relationByOID[d.ObjOID]; child != nil && child.OID != rel.OID {
				c.cascadeDropRelationSeen(child, seen)
			}
		}
	}
	if hasNormalDependent(c.relationDependents(rel.OID, 0)) {
		c.dropDependents('r', rel.OID)
	}
	if rel.RowTypeOID != 0 {
		if bt := c.typeByOID[rel.RowTypeOID]; bt != nil {
			for _, depRelOID := range c.typeRelationDependents(bt) {
				if depRel := c.relationByOID[depRelOID]; depRel != nil && depRel.OID != rel.OID {
					c.cascadeDropRelationSeen(depRel, seen)
				}
			}
			c.dropTypeDependents(bt)
		}
	}
	if c.relationByOID[rel.OID] != nil {
		c.removeRelation(rel.Schema, rel.Name, rel)
	}
}

// generateConstraintName generates a default constraint name.
func generateConstraintName(tableName string, columns []string, ctype ConstraintType) string {
	switch ctype {
	case ConstraintPK:
		return tableName + "_pkey"
	case ConstraintUnique:
		return tableName + "_" + strings.Join(columns, "_") + "_key"
	case ConstraintFK:
		return tableName + "_" + strings.Join(columns, "_") + "_fkey"
	case ConstraintCheck:
		if len(columns) > 0 {
			return tableName + "_" + strings.Join(columns, "_") + "_check"
		}
		return tableName + "_check"
	case ConstraintExclude:
		return tableName + "_" + strings.Join(columns, "_") + "_excl"
	default:
		return tableName + "_con"
	}
}

// indexOfAttNum returns the column index in rel whose AttNum is attnum, or -1.
func indexOfAttNum(rel *Relation, attnum int16) int {
	for i, col := range rel.Columns {
		if col.AttNum == attnum {
			return i
		}
	}
	return -1
}
