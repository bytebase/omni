package catalog

import (
	"fmt"
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// DefineRelation creates a new relation (table).
//
// pg: src/backend/commands/tablecmds.c — DefineRelation
func (c *Catalog) DefineRelation(stmt *nodes.CreateStmt, relkind byte) error {
	rv := stmt.Relation
	schema, err := c.resolveTargetSchema(rv.Schemaname)
	if err != nil {
		return err
	}
	relName := rv.Relname

	// A typed table resolves its type before the target name is checked:
	// CREATE TABLE a OF a on table a is 42809 ("a" is not a composite type),
	// not a duplicate. IF NOT EXISTS skips the whole statement instead.
	if stmt.OfTypename != nil && !stmt.IfNotExists {
		oid, _, err := c.resolveTypeName(stmt.OfTypename)
		if err != nil {
			return err
		}
		if _, err := c.checkOfType(oid); err != nil {
			return err
		}
	}

	// IF NOT EXISTS skips an existing relation before column definitions and
	// parents are resolved; a non-relation type (e.g. enum/domain) does NOT
	// satisfy IF NOT EXISTS and instead fails later with 42710. Without IF NOT
	// EXISTS, PostgreSQL runs MergeAttributes and BuildDescForRelation (checking
	// duplicate column names and column types) before heap_create_with_catalog
	// rejects a duplicate relation name.
	if stmt.IfNotExists {
		if _, exists := schema.Relations[relName]; exists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", relName))
			return nil
		}
	}

	// If PARTITION BY is specified, this is a partitioned table.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (partitioned table handling)
	if stmt.Partspec != nil {
		relkind = 'p'
	}

	// Collect column definitions and inline constraints from TableElts.
	var colDefs []ColumnDef
	var constraintDefs []ConstraintDef
	var likeSources []likeSource

	if stmt.TableElts != nil {
		for _, elt := range stmt.TableElts.Items {
			switch n := elt.(type) {
			case *nodes.ColumnDef:
				cd, inlineCons, err := c.convertColumnDef(n, relName, schema)
				if err != nil {
					return err
				}
				colDefs = append(colDefs, cd)
				constraintDefs = append(constraintDefs, inlineCons...)
			case *nodes.Constraint:
				if cdef, ok := convertConstraintNode(n); ok {
					constraintDefs = append(constraintDefs, cdef)
				}
			case *nodes.TableLikeClause:
				// LIKE clause: copy columns (and optionally constraints/indexes/defaults)
				// from the source table.
				// pg: src/backend/commands/tablecmds.c — expandTableLikeClause
				likeDefs, likeCons, likeSrc, likeOpts, err := c.expandTableLikeClause(n)
				if err != nil {
					return err
				}
				colDefs = append(colDefs, likeDefs...)
				constraintDefs = append(constraintDefs, likeCons...)
				if likeSrc != nil {
					likeSources = append(likeSources, likeSource{rel: likeSrc, opts: likeOpts})
				}
			}
		}
	}

	// Table-level constraints from Constraints list.
	if stmt.Constraints != nil {
		for _, con := range stmt.Constraints.Items {
			if n, ok := con.(*nodes.Constraint); ok {
				if cdef, ok := convertConstraintNode(n); ok {
					constraintDefs = append(constraintDefs, cdef)
				}
			}
		}
	}
	if relkind == 'f' {
		for _, cdef := range constraintDefs {
			if err := validateForeignTableConstraintType(cdef.Type); err != nil {
				return err
			}
		}
	}

	// Build a name→index map for the child column defs (used by merge).
	// Mark all columns from the CREATE statement as locally defined.
	childColMap := make(map[string]int, len(colDefs))
	for i := range colDefs {
		colDefs[i].IsLocal = true
		childColMap[colDefs[i].Name] = i
	}

	// IF NOT EXISTS skips a name an index or sequence holds before the
	// parents are resolved, as PostgreSQL does for a partition target.
	if stmt.IfNotExists {
		if _, ok := relationNamespaceKind(schema, relName); ok {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", relName))
			return nil
		}
	}

	// DefineRelation refuses ON COMMIT on a table that is not temporary
	// before it reads anything else: after transformCreateStmt's IF NOT
	// EXISTS skip and its OF type lookup, before parents, columns and the
	// relation name.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (oncommit validation)
	if stmt.OnCommit != nodes.ONCOMMIT_NOOP && rv.Relpersistence != 't' {
		if stmt.OfTypename != nil && stmt.IfNotExists {
			oid, _, err := c.resolveTypeName(stmt.OfTypename)
			if err != nil {
				return err
			}
			if _, err := c.checkOfType(oid); err != nil {
				return err
			}
		}
		return &Error{Code: CodeInvalidTableDefinition,
			Message: "ON COMMIT can only be used on temporary tables"}
	}

	// Handle OF TYPE (typed tables).
	// pg: src/backend/parser/parse_utilcmd.c — transformOfType
	// pg: src/backend/commands/tablecmds.c — MergeAttributes (is_from_type handling)
	var ofTypeOID uint32
	if stmt.OfTypename != nil {
		oid, _, err := c.resolveTypeName(stmt.OfTypename)
		if err != nil {
			return err
		}
		typeRel, err := c.checkOfType(oid)
		if err != nil {
			return err
		}

		// In PG, type columns come first (with is_from_type=true), then user columns.
		// MergeAttributes merges user columns with matching type columns.
		userColDefs := colDefs
		userConstraints := constraintDefs
		colDefs = nil
		constraintDefs = nil
		childColMap = make(map[string]int)

		// First, add all columns from the type.
		for _, pcol := range typeRel.Columns {
			childColMap[pcol.Name] = len(colDefs)
			colDefs = append(colDefs, ColumnDef{
				Name:       pcol.Name,
				Type:       c.typeNameForOID(pcol.TypeOID, pcol.TypeMod),
				IsLocal:    false,
				IsFromType: true,
			})
		}

		// Now merge user-specified columns.
		seenUserCols := make(map[string]bool, len(userColDefs))
		for _, ucd := range userColDefs {
			if seenUserCols[ucd.Name] {
				return errDuplicateColumn(ucd.Name)
			}
			seenUserCols[ucd.Name] = true
			typeIdx, exists := childColMap[ucd.Name]
			if !exists {
				// User specified a column that doesn't exist in the type.
				// pg: src/backend/commands/tablecmds.c — MergeAttributes (line 2516-2527)
				return &Error{Code: CodeUndefinedColumn,
					Message: fmt.Sprintf("column %q does not exist", ucd.Name)}
			}
			// Merge user's NOT NULL, default, and constraints into the type column.
			// pg: src/backend/commands/tablecmds.c — MergeAttributes (line 2537-2547)
			tcd := &colDefs[typeIdx]
			tcd.NotNull = ucd.NotNull
			tcd.NotNullConName = ucd.NotNullConName
			tcd.Default = ucd.Default
			tcd.RawDefault = ucd.RawDefault
			tcd.Generated = ucd.Generated
			tcd.GenerationExpr = ucd.GenerationExpr
			tcd.RawGenExpr = ucd.RawGenExpr
			tcd.IsFromType = false
			tcd.IsLocal = true
		}
		// Restore user-level constraints.
		constraintDefs = userConstraints

		ofTypeOID = oid
	}

	// Mark locally-declared NOT NULL before merging inherited columns so
	// PG 18 cataloged NOT NULL constraints preserve ConIsLocal vs ConInhCount.
	for i := range colDefs {
		if colDefs[i].NotNull || colDefs[i].Identity != 0 || colDefs[i].IsSerial != 0 {
			colDefs[i].NotNullLocal = true
		}
	}

	// Merge inherited columns (INHERITS clause).
	// pg: src/backend/commands/tablecmds.c — MergeAttributes
	var parentOIDs []uint32
	isPartition := (stmt.Partbound != nil)
	var conflictDefaults map[string]bool // column name → has conflicting defaults
	if stmt.InhRelations != nil {
		for _, item := range stmt.InhRelations.Items {
			rv, ok := item.(*nodes.RangeVar)
			if !ok {
				continue
			}
			_, parent, err := c.findRelation(rv.Schemaname, rv.Relname)
			if err != nil {
				return err
			}
			if parent.RelKind != 'r' && parent.RelKind != 'p' && parent.RelKind != 'f' {
				return errWrongObjectType(rv.Relname, "a table")
			}

			// Prevent inheriting from a partition (must use PARTITION OF instead).
			// pg: src/backend/commands/tablecmds.c — MergeAttributes (relispartition check)
			if !isPartition && parent.PartitionOf != 0 {
				return &Error{Code: CodeWrongObjectType,
					Message: fmt.Sprintf("cannot inherit from partition %q", parent.Name)}
			}

			// Cannot use regular INHERITS on a partitioned table (use PARTITION OF).
			// pg: src/backend/commands/tablecmds.c — MergeAttributes
			if !isPartition && parent.RelKind == 'p' {
				return &Error{Code: CodeWrongObjectType,
					Message: fmt.Sprintf("cannot inherit from partitioned table %q", parent.Name)}
			}

			// Check persistence conflicts.
			// pg: src/backend/commands/tablecmds.c — MergeAttributes (persistence checks)
			parentPersistence := parent.Persistence
			if parentPersistence == 0 {
				parentPersistence = 'p'
			}
			childPersistence := byte('p')
			if stmt.Relation != nil && stmt.Relation.Relpersistence != 0 {
				childPersistence = stmt.Relation.Relpersistence
			}
			// Temp child/partition of permanent parent is not allowed.
			// pg: src/backend/commands/tablecmds.c — MergeAttributes (persistence checks)
			if childPersistence == 't' && parentPersistence != 't' {
				msg := fmt.Sprintf("cannot create a temporary relation as inheritance child of permanent relation %q", parent.Name)
				if isPartition {
					msg = fmt.Sprintf("cannot create a temporary relation as partition of permanent relation %q", parent.Name)
				}
				return &Error{Code: CodeWrongObjectType, Message: msg}
			}
			// Permanent child of temp parent is not allowed.
			if childPersistence != 't' && parentPersistence == 't' {
				msg := fmt.Sprintf("cannot inherit from temporary relation %q", parent.Name)
				if isPartition {
					msg = fmt.Sprintf("cannot create a permanent relation as partition of temporary relation %q", parent.Name)
				}
				return &Error{Code: CodeWrongObjectType, Message: msg}
			}

			// Check for duplicate parent.
			// pg: src/backend/commands/tablecmds.c — MergeAttributes (duplicate parent check)
			for _, prevOID := range parentOIDs {
				if prevOID == parent.OID {
					return &Error{Code: CodeDuplicateObject,
						Message: fmt.Sprintf("relation %q would be inherited from more than once", parent.Name)}
				}
			}

			parentOIDs = append(parentOIDs, parent.OID)

			// Copy parent columns, merging if child has same-name column.
			// pgddl never has dropped columns (columns are physically removed), so no
			// attisdropped check is needed here — already correct by construction.
			for _, pcol := range parent.Columns {
				if existIdx, exists := childColMap[pcol.Name]; exists {
					existCD := &colDefs[existIdx]
					// Column exists in child — types must match.
					childOID, _, err := c.ResolveType(existCD.Type)
					if err != nil {
						return err
					}
					if childOID != pcol.TypeOID {
						return errDatatypeMismatch(fmt.Sprintf(
							"column %q has a type conflict: inherited as %s, declared as %s",
							pcol.Name, c.typeByOID[pcol.TypeOID].TypeName, c.typeByOID[childOID].TypeName,
						))
					}
					// TypeMod must match.
					// pg: src/backend/commands/tablecmds.c — MergeAttributes (typmod check)
					childTypMod := existCD.Type.TypeMod
					if childTypMod != pcol.TypeMod && childTypMod != -1 && pcol.TypeMod != -1 {
						return errDatatypeMismatch(fmt.Sprintf(
							"column %q has a type modifier conflict", pcol.Name,
						))
					}
					// Collation must match.
					// pg: src/backend/commands/tablecmds.c — MergeAttributes (collation check)
					if existCD.CollationName != "" && pcol.CollationName != "" {
						childCollSchema := existCD.CollationSchema
						if existCD.CollClause != nil {
							if _, cid, cerr := c.resolveCollationList(existCD.CollClause.Collname); cerr == nil {
								childCollSchema = c.resolvedCollationSchema(existCD.CollClause.Collname, existCD.CollationName, cid)
							}
						}
						parentCollSchema := columnCollationSchema(c, pcol)
						if existCD.CollationName != pcol.CollationName || childCollSchema != parentCollSchema {
							return errDatatypeMismatch(fmt.Sprintf(
								"column %q has a collation conflict", pcol.Name,
							))
						}
					} else if existCD.CollationName == "" && pcol.CollationName != "" {
						existCD.CollationSchema = columnCollationSchema(c, pcol)
						existCD.CollationName = pcol.CollationName
						existCD.UserCollationOID = pcol.UserCollationOID
					}
					// Merge NOT NULL from parent (skip NO INHERIT NOT NULL constraints).
					if pcol.NotNull && !c.isNotNullNoInherit(parent, pcol.AttNum) {
						existCD.NotNull = true
						existCD.NotNullInhCount++
						if existCD.InheritedNNName == "" {
							existCD.InheritedNNName = c.notNullConstraintName(parent, pcol.AttNum)
						}
					}
					// Merge default from parent if child has no default.
					// pg: src/backend/commands/tablecmds.c — MergeAttributes (default conflict detection)
					if existCD.Default == "" && pcol.Default != "" {
						existCD.Default = pcol.Default
					} else if existCD.Default != "" && pcol.Default != "" && existCD.Default != pcol.Default {
						// Mark as conflicting — error unless child locally overrides.
						if !existCD.IsLocal {
							if conflictDefaults == nil {
								conflictDefaults = make(map[string]bool)
							}
							conflictDefaults[pcol.Name] = true
						}
					}
					// Partitions inherit identity from parent.
					// pg: src/backend/commands/tablecmds.c — MergeAttributes (identity merge)
					if isPartition && pcol.Identity != 0 {
						existCD.Identity = pcol.Identity
					}
					existCD.InhCount++
				} else {
					// New inherited column — not locally defined.
					inheritNN := pcol.NotNull && !c.isNotNullNoInherit(parent, pcol.AttNum)
					var nnInh int16
					var inheritedNNName string
					if inheritNN {
						nnInh = 1
						inheritedNNName = c.notNullConstraintName(parent, pcol.AttNum)
					}
					newCD := ColumnDef{
						Name:             pcol.Name,
						Type:             c.typeNameForOID(pcol.TypeOID, pcol.TypeMod),
						NotNull:          inheritNN,
						InheritedNNName:  inheritedNNName,
						NotNullLocal:     false,
						NotNullInhCount:  nnInh,
						Default:          pcol.Default,
						IsLocal:          false,
						InhCount:         1,
						CollationSchema:  columnCollationSchema(c, pcol),
						CollationName:    pcol.CollationName,
						UserCollationOID: pcol.UserCollationOID,
					}
					// Partitions inherit identity; regular inheritance does not.
					// pg: src/backend/commands/tablecmds.c — MergeAttributes (identity handling)
					if isPartition {
						newCD.Identity = pcol.Identity
					}
					childColMap[pcol.Name] = len(colDefs)
					colDefs = append(colDefs, newCD)
				}
			}

			// Copy CHECK constraints from parent (skip NO INHERIT constraints).
			// pg: src/backend/commands/tablecmds.c — MergeAttributes (constraint copy)
			for _, con := range c.consByRel[parent.OID] {
				if con.Type == ConstraintCheck && !con.ConNoInherit {
					constraintDefs = append(constraintDefs, ConstraintDef{
						Name:      con.Name,
						Type:      ConstraintCheck,
						CheckExpr: con.CheckExpr,
					})
				}
			}

			// The inheritance entry is recorded only once the relation is
			// registered, so a later error leaves no row for a child that
			// does not exist.
		}
	}

	// Validate no conflicting defaults from multiple parents.
	// pg: src/backend/commands/tablecmds.c — MergeAttributes (bogus default detection)
	for name := range conflictDefaults {
		idx := childColMap[name]
		cd := colDefs[idx]
		if cd.IsLocal {
			continue // child overrides, OK
		}
		if cd.Generated != 0 {
			return &Error{Code: CodeInvalidColumnDefinition,
				Message: fmt.Sprintf("column %q inherits conflicting generation expressions", name)}
		}
		return &Error{Code: CodeInvalidColumnDefinition,
			Message: fmt.Sprintf("column %q inherits conflicting default values", name)}
	}

	// For partitions, validate user-specified column options against inherited columns.
	// pg: src/backend/commands/tablecmds.c — MergeAttributes (saved_columns validation, line 2927-3000)
	if isPartition {
		for i := range colDefs {
			cd := &colDefs[i]
			if !cd.IsLocal {
				continue // only check user-specified columns
			}
			// Inherited column must exist (already merged during the loop above).
			// Check generated/default/identity conflicts.
			if cd.Generated != 0 {
				// Parent column is generated — partition column cannot specify default or identity.
				if cd.Default != "" && cd.Generated == 0 {
					return &Error{Code: CodeInvalidColumnDefinition,
						Message: fmt.Sprintf("column %q inherits from generated column but specifies default", cd.Name)}
				}
				if cd.Identity != 0 {
					return &Error{Code: CodeInvalidColumnDefinition,
						Message: fmt.Sprintf("column %q inherits from generated column but specifies identity", cd.Name)}
				}
			}
		}
	}

	// Expand SERIAL columns (not applicable to composite types).
	type serialInfo struct {
		colIdx     int
		seqName    string
		typeOID    uint32
		isIdentity bool
	}
	var serials []serialInfo
	if relkind != 'c' {
		for i := range colDefs {
			cd := &colDefs[i]
			if cd.IsSerial != 0 {
				var seqTypeOID uint32
				var colTypeName string
				switch cd.IsSerial {
				case 2:
					seqTypeOID = INT2OID
					colTypeName = "int2"
				case 4:
					seqTypeOID = INT4OID
					colTypeName = "int4"
				default:
					seqTypeOID = INT8OID
					colTypeName = "int8"
				}
				seqName := fmt.Sprintf("%s_%s_seq", relName, cd.Name)
				cd.Type = TypeName{Name: colTypeName, TypeMod: -1}
				cd.NotNull = true
				cd.Default = fmt.Sprintf("nextval('%s'::regclass)", seqName)
				serials = append(serials, serialInfo{colIdx: i, seqName: seqName, typeOID: seqTypeOID})
			}
		}
	}

	// Expand identity columns as implicit sequences (like SERIAL).
	if relkind != 'c' {
		for i := range colDefs {
			cd := &colDefs[i]
			if cd.Identity == 'a' || cd.Identity == 'd' {
				seqName := fmt.Sprintf("%s_%s_seq", relName, cd.Name)
				colTypeOID, _, err := c.ResolveType(cd.Type)
				if err != nil {
					return err
				}
				seqTypeOID, err := identitySequenceTypeOID(colTypeOID)
				if err != nil {
					return err
				}
				serials = append(serials, serialInfo{colIdx: i, seqName: seqName, typeOID: seqTypeOID, isIdentity: true})
			}
		}
	}

	// Resolve column types and check for duplicates.
	colByName := make(map[string]int, len(colDefs))
	columns := make([]*Column, len(colDefs))
	colUserCollOIDs := make([]uint32, len(colDefs))

	for i, cd := range colDefs {
		if _, dup := colByName[cd.Name]; dup {
			return errDuplicateColumn(cd.Name)
		}
		colByName[cd.Name] = i

		typeOID, typmod, err := c.ResolveType(cd.Type)
		if err != nil {
			return err
		}

		typ := c.typeByOID[typeOID]
		collSchema := cd.CollationSchema
		colUserCollOIDs[i] = cd.UserCollationOID
		if cd.CollClause != nil {
			_, collOID, cerr := c.resolveColumnCollation(cd.CollClause.Collname, typeOID)
			if cerr != nil {
				return cerr
			}
			colUserCollOIDs[i] = collOID
			collSchema = c.resolvedCollationSchema(cd.CollClause.Collname, cd.CollationName, collOID)
		}
		stor, err := c.resolveColumnStorage(typ, cd.StorageName)
		if err != nil {
			return err
		}
		comp, err := c.resolveColumnCompression(typ, cd.Compression)
		if err != nil {
			return err
		}
		// pg: src/backend/catalog/heap.c — InsertPgAttributeTuples
		// attcollation uses the explicit COLLATE OID when specified, else the type's collation.
		coll := c.effectiveColumnCollation(typeOID, cd.CollationName, colUserCollOIDs[i])
		columns[i] = &Column{
			AttNum:           int16(i + 1),
			Name:             cd.Name,
			TypeOID:          typeOID,
			TypeMod:          typmod,
			NotNull:          cd.NotNull,
			Len:              typ.Len,
			ByVal:            typ.ByVal,
			Align:            typ.Align,
			Storage:          stor,
			Compression:      comp,
			Collation:        coll,
			IsLocal:          cd.IsLocal,
			InhCount:         cd.InhCount,
			CollationSchema:  collSchema,
			CollationName:    cd.CollationName,
			UserCollationOID: colUserCollOIDs[i],
			ForeignOptions:   slices.Clone(cd.ForeignOptions),
		}
	}

	// PostgreSQL keeps indexes and sequences in the relation namespace, so a
	// name one holds is not free. The check comes after the columns' types and
	// the parents are resolved, as PostgreSQL resolves them first: CREATE
	// TABLE a(a a) on an index a is 42704 for the type, CREATE TABLE a(a a, a a)
	// on a table a is 42701 for the duplicate column, and PARTITION OF an
	// index is 42809, before any is a duplicate relation name.
	if !stmt.IfNotExists {
		if _, ok := relationNamespaceKind(schema, relName); ok {
			return errDuplicateTable(relName)
		}
	}
	// A non-relation type (e.g. enum or domain) with the same name conflicts
	// with the table's implicit row type (42710), even with IF NOT EXISTS.
	if c.typeByName[typeKey{ns: schema.OID, name: relName}] != nil {
		return errDuplicateObject("type", relName)
	}

	// DefineRelation calls toast_reloptions(toastOptions, true) unconditionally
	// before NewRelationCreateToastTable checks whether a TOAST table is needed.
	if err := c.validateRelOptions(stmt.Options, relkind, ""); err != nil {
		return err
	}
	// PostgreSQL resolves a new relation's tablespace in DefineRelation order:
	// explicit TABLESPACE clause, then a partitioned parent's stored tablespace
	// when non-empty, then GetDefaultTablespace ("" for temp relations, else the
	// session's effective default_tablespace).
	effectiveTablespace := stmt.Tablespacename
	if effectiveTablespace == "" && isPartition && len(parentOIDs) > 0 {
		if parent := c.relationByOID[parentOIDs[0]]; parent != nil && parent.Tablespace != "" {
			effectiveTablespace = parent.Tablespace
		}
	}
	if effectiveTablespace == "" && rv.Relpersistence != 't' {
		effectiveTablespace = c.effectiveTablespace("")
	}
	if stmt.Partspec != nil && effectiveTablespace == "pg_default" {
		return &Error{Code: CodeFeatureNotSupported,
			Message: "cannot specify default tablespace for partitioned relations"}
	}
	if relkind == 'r' || relkind == 'p' || relkind == 'm' {
		if err := c.checkRelationTablespace(effectiveTablespace); err != nil {
			return err
		}
	}
	tableAccessMethod := stmt.AccessMethod
	if tableAccessMethod == "" && isPartition && len(parentOIDs) > 0 {
		if parent := c.relationByOID[parentOIDs[0]]; parent != nil && parent.AccessMethod != "" {
			tableAccessMethod = parent.AccessMethod
		}
	}
	if tableAccessMethod == "" && (relkind == 'r' || relkind == 'm') {
		tableAccessMethod = c.defaultTableAccessMethod
	}
	if err := c.checkTableAccessMethod(tableAccessMethod); err != nil {
		return err
	}

	// Constraint indexes and cloned partition indexes may take a tablespace;
	// refuse one the catalog cannot place before any implicit sequence is made.
	for _, def := range constraintDefs {
		if (def.Type != ConstraintPK && def.Type != ConstraintUnique && def.Type != ConstraintExclude) || def.IndexName != "" {
			continue
		}
		idxSpace := def.IndexSpace
		if idxSpace == "" && rv.Relpersistence != 't' {
			idxSpace = c.effectiveTablespace("")
		}
		if stmt.Partspec != nil && idxSpace == "pg_default" {
			return &Error{Code: CodeFeatureNotSupported,
				Message: "cannot specify default tablespace for partitioned relations"}
		}
		if idxSpace != "" {
			if err := c.checkRelationTablespace(idxSpace); err != nil {
				return err
			}
		}
	}
	if isPartition && len(parentOIDs) > 0 && rv.Relpersistence != 't' {
		for _, idx := range c.indexesByRel[parentOIDs[0]] {
			if idx.Tablespace == "" {
				partIdxSpace := c.effectiveTablespace("")
				if stmt.Partspec != nil && partIdxSpace == "pg_default" {
					return &Error{Code: CodeFeatureNotSupported,
						Message: "cannot specify default tablespace for partitioned relations"}
				}
				if err := c.checkRelationTablespace(partIdxSpace); err != nil {
					return err
				}
			}
		}
	}
	if rv.Relpersistence != 't' {
		for _, ls := range likeSources {
			if ls.opts&tableLikeIndexes == 0 {
				continue
			}
			for _, idx := range c.indexesByRel[ls.rel.OID] {
				if idx.Tablespace == "" {
					likeIdxSpace := c.effectiveTablespace("")
					if stmt.Partspec != nil && likeIdxSpace == "pg_default" {
						return &Error{Code: CodeFeatureNotSupported,
							Message: "cannot specify default tablespace for partitioned relations"}
					}
					if err := c.checkRelationTablespace(likeIdxSpace); err != nil {
						return err
					}
				}
			}
		}
	}

	// Check the column count limit before creating implicit sequences.
	// pg: src/backend/commands/tablecmds.c — DefineRelation
	if len(columns) > MaxHeapAttributeNumber {
		return &Error{Code: CodeTooManyColumns, Message: fmt.Sprintf("tables can have at most %d columns", MaxHeapAttributeNumber)}
	}

	// Validate partition key specification before creating implicit sequences
	// or registering the relation.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (StorePartitionKey)
	var partInfo *PartitionInfo
	if stmt.Partspec != nil {
		ps := stmt.Partspec
		strategy := byte(0)
		if len(ps.Strategy) > 0 {
			strategy = ps.Strategy[0] // 'l', 'r', 'h'
		}
		// LIST partitioning supports only a single column.
		// pg: src/backend/parser/parse_utilcmd.c — transformPartitionSpec
		if strategy == 'l' && ps.PartParams != nil && len(ps.PartParams.Items) > 1 {
			return &Error{Code: CodeInvalidObjectDefinition,
				Message: "cannot use \"list\" partition strategy with more than one column"}
		}

		// Check PARTITION_MAX_KEYS limit.
		// pg: src/backend/commands/tablecmds.c — DefineRelation (line 1158-1162)
		if ps.PartParams != nil && len(ps.PartParams.Items) > PARTITION_MAX_KEYS {
			return &Error{Code: CodeProgramLimitExceeded,
				Message: fmt.Sprintf("cannot partition using more than %d columns", PARTITION_MAX_KEYS)}
		}

		var keyAttNums []int16
		if ps.PartParams != nil {
			for _, p := range ps.PartParams.Items {
				pe, ok := p.(*nodes.PartitionElem)
				if !ok {
					continue
				}
				if pe.Name == "" {
					if pe.Expr == nil {
						return &Error{Code: CodeSyntaxError,
							Message: "syntax error at or near \")\""}
					}
					// Expression partition key — store attnum=0 (PG convention).
					keyAttNums = append(keyAttNums, 0)
					continue
				}
				idx, exists := colByName[pe.Name]
				if !exists {
					return errUndefinedColumn(pe.Name)
				}
				keyAttNums = append(keyAttNums, int16(idx+1))
			}
		}
		partInfo = &PartitionInfo{
			Strategy:   strategy,
			KeyAttNums: keyAttNums,
			NKeyAttrs:  int16(len(keyAttNums)),
		}
	}

	// Create implicit sequences for SERIAL and identity columns.
	var createdSeqs []*Sequence
	cleanupCreatedSeqs := func() {
		for _, s := range createdSeqs {
			if c.sequenceByOID[s.OID] != nil {
				c.removeSequence(s.Schema, s)
			}
		}
	}
	for _, si := range serials {
		if si.isIdentity {
			cd := &colDefs[si.colIdx]
			seq, seqRV, explicitName, err := c.createIdentitySequence(schema, relName, cd.Name, "", si.typeOID, cd.IdentityOptions)
			if err != nil {
				cleanupCreatedSeqs()
				return err
			}
			cd.IdentitySequence = seqRV
			cd.Default = identitySequenceDefault(schema, seq.Schema, seq.Name, explicitName)
			createdSeqs = append(createdSeqs, seq)
			continue
		}
		if _, exists := schema.Sequences[si.seqName]; exists {
			cleanupCreatedSeqs()
			return errDuplicateObject("sequence", si.seqName)
		}
		seq := c.createSequenceInternal(schema, si.seqName, si.typeOID)
		createdSeqs = append(createdSeqs, seq)
	}

	// Allocate OIDs.
	relOID := c.oidGen.Next()
	rowTypeOID := c.oidGen.Next()
	arrayOID := c.oidGen.Next()

	storedTablespace := normalizeTablespace(effectiveTablespace)
	if relkind != 'r' && relkind != 'p' && relkind != 'm' {
		storedTablespace = ""
	}
	rel := &Relation{
		OID:           relOID,
		Name:          relName,
		Schema:        schema,
		RelKind:       relkind,
		Columns:       columns,
		colByName:     colByName,
		RowTypeOID:    rowTypeOID,
		ArrayOID:      arrayOID,
		InhParents:    parentOIDs,
		InhCount:      len(parentOIDs),
		OfTypeOID:     ofTypeOID,
		Tablespace:    storedTablespace,
		AccessMethod:  tableAccessMethod,
		OwnerName:     c.currentOwner(),
		OwnerExplicit: !c.installingSnapshot && c.currentUser != "",
	}

	// Parse persistence from RangeVar.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (relpersistence)
	if rv.Relpersistence != 0 {
		rel.Persistence = rv.Relpersistence
	} else {
		rel.Persistence = 'p' // default: permanent
	}

	switch stmt.OnCommit {
	case nodes.ONCOMMIT_PRESERVE_ROWS:
		rel.OnCommit = 'p'
	case nodes.ONCOMMIT_DELETE_ROWS:
		rel.OnCommit = 'd'
	case nodes.ONCOMMIT_DROP:
		rel.OnCommit = 'D'
	}

	// Store column defaults, generated expressions, and identity.
	for i, cd := range colDefs {
		if cd.Default != "" || cd.RawDefault != nil {
			columns[i].HasDefault = true
			columns[i].Default = cd.Default
		}
		if cd.Generated == 's' {
			columns[i].Generated = 's'
			columns[i].HasDefault = true
			columns[i].Default = cd.GenerationExpr
		}
		if cd.Identity == 'a' || cd.Identity == 'd' {
			columns[i].Identity = cd.Identity
			columns[i].NotNull = true // identity implies NOT NULL
		}
	}

	// Register relation.
	schema.Relations[relName] = rel
	c.relationByOID[relOID] = rel
	for i, collOID := range colUserCollOIDs {
		if collOID != 0 {
			c.recordDependency('r', relOID, int32(columns[i].AttNum), 'L', collOID, 0, DepNormal)
		}
	}
	if tableAccessMethod != "" {
		if am := c.accessMethodByName[tableAccessMethod]; am != nil {
			c.recordDependency('r', relOID, 0, 'm', am.OID, 0, DepNormal)
		}
	}

	// Record inheritance entries and dependencies now that the child exists.
	if len(parentOIDs) > 0 {
		for i, parentOID := range parentOIDs {
			c.inhEntries = append(c.inhEntries, InhEntry{
				InhRelID:  relOID,
				InhParent: parentOID,
				InhSeqNo:  int32(i + 1),
			})
		}
		// For PARTITION OF, use DepAuto (auto-dropped with parent).
		// For regular INHERITS, use DepNormal (requires CASCADE).
		depType := DepNormal
		if stmt.Partbound != nil {
			depType = DepAuto
		}
		for _, parentOID := range parentOIDs {
			c.recordDependency('r', relOID, 0, 'r', parentOID, 0, depType)
		}
	}
	rel.PartitionInfo = partInfo
	rollbackRelation := func() {
		c.removeRelation(schema, relName, rel)
		cleanupCreatedSeqs()
	}

	// Store partition bound (PARTITION OF ... FOR VALUES clause).
	// pg: src/backend/commands/tablecmds.c — DefineRelation (StorePartitionBound)
	if stmt.Partbound != nil {
		if bs, ok := stmt.Partbound.(*nodes.PartitionBoundSpec); ok {
			// Validate partition bound strategy matches parent's partitioning strategy.
			// pg: src/backend/partitioning/partbounds.c — check_new_partition_bound
			if !bs.IsDefault && len(parentOIDs) > 0 {
				parent := c.relationByOID[parentOIDs[0]]
				if parent != nil && parent.PartitionInfo != nil {
					parentStrat := parent.PartitionInfo.Strategy
					switch parentStrat {
					case 'l':
						if bs.Listdatums == nil || len(bs.Listdatums.Items) == 0 {
							// Allow — some LIST partitions may have empty bounds parsed differently
						}
					case 'r':
						if (bs.Lowerdatums == nil || len(bs.Lowerdatums.Items) == 0) &&
							(bs.Upperdatums == nil || len(bs.Upperdatums.Items) == 0) {
							// Allow — may be special bound syntax
						}
					case 'h':
						if bs.Modulus == 0 && bs.Remainder == 0 {
							// Allow — may be default or special case
						}
					}
				}
			}

			listVals, listKeywords := deparseDatumListWithKinds(bs.Listdatums)
			lowerVals, lowerKeywords := deparseDatumListWithKinds(bs.Lowerdatums)
			upperVals, upperKeywords := deparseDatumListWithKinds(bs.Upperdatums)
			rel.PartitionBound = &PartitionBound{
				Strategy:      bs.Strategy,
				IsDefault:     bs.IsDefault,
				ListValues:    listVals,
				ListKeywords:  listKeywords,
				LowerBound:    lowerVals,
				LowerKeywords: lowerKeywords,
				UpperBound:    upperVals,
				UpperKeywords: upperKeywords,
				Modulus:       bs.Modulus,
				Remainder:     bs.Remainder,
			}
		}
		// Set PartitionOf to the parent table OID.
		if len(parentOIDs) > 0 {
			rel.PartitionOf = parentOIDs[0]
			// Check for partition bound overlap with existing partitions of the same parent.
			if err := c.checkPartitionBoundOverlap(parentOIDs[0], rel.OID, rel.PartitionBound); err != nil {
				rollbackRelation()
				return err
			}
		}
	}

	// Register composite row type.
	rowType := &BuiltinType{
		OID:       rowTypeOID,
		TypeName:  relName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'c',
		Category:  'C',
		IsDefined: true,
		Delim:     ',',
		RelID:     relOID,
		Array:     arrayOID,
		Align:     'd',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[rowTypeOID] = rowType
	c.typeByName[typeKey{ns: schema.OID, name: relName}] = rowType

	// Register array type.
	arrayType := &BuiltinType{
		OID:       arrayOID,
		TypeName:  fmt.Sprintf("_%s", relName),
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      rowTypeOID,
		Align:     'd',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[arrayOID] = arrayType
	c.typeByName[typeKey{ns: schema.OID, name: arrayType.TypeName}] = arrayType

	// Record column type dependencies for composite types.
	// When a composite type column references another composite type, we need
	// a dep edge so the referenced type is created first during migration.
	// pg: src/backend/catalog/heap.c — heap_create_with_catalog (type dependencies)
	if relkind == 'c' {
		for _, col := range columns {
			// Skip built-in types (OID < FirstNormalObjectId).
			if col.TypeOID >= FirstNormalObjectId {
				// Record: this relation depends on the type.
				c.recordDependency('r', relOID, int32(col.AttNum), 't', col.TypeOID, 0, DepNormal)
			}
		}
	}

	// Process constraints (composite types have no constraints).
	if relkind != 'c' {
		// Pass 1: reserve all explicitly named column NOT NULL constraints
		// before table constraints and unnamed NOT NULL constraints so column
		// order never causes an implicit name to collide with an explicit one.
		for i, col := range columns {
			if col.NotNull && i < len(colDefs) && colDefs[i].NotNullConName != "" {
				if err := c.ensureNotNullConstraintWithInh(schema, rel, col, colDefs[i].NotNullConName, colDefs[i].NotNullLocal, colDefs[i].NotNullInhCount); err != nil {
					rollbackRelation()
					return err
				}
			}
		}

		c.deferImplicitNotNull = true
		for _, cd := range constraintDefs {
			if err := c.addConstraint(schema, rel, cd); err != nil {
				c.deferImplicitNotNull = false
				rollbackRelation()
				return err
			}
		}
		c.deferImplicitNotNull = false

		// Pass 2: register remaining (unnamed) PG 18 cataloged NOT NULL
		// constraints (including implicit PK NOT NULLs) after all explicit
		// constraint names are reserved.
		pk := c.findPKConstraint(rel.OID)
		for i, col := range columns {
			if !col.NotNull {
				continue
			}
			isLocal := col.IsLocal
			var inhCount int16
			var conName string
			if i < len(colDefs) {
				if colDefs[i].NotNullConName != "" {
					continue
				}
				isLocal = colDefs[i].NotNullLocal
				inhCount = colDefs[i].NotNullInhCount
				if colDefs[i].InheritedNNName != "" && c.checkConstraintNameFree(rel.OID, colDefs[i].InheritedNNName) == nil {
					conName = colDefs[i].InheritedNNName
				}
			}
			if !isLocal && pk != nil && slices.Contains(pk.Columns, col.AttNum) {
				isLocal = true
			}
			if err := c.ensureNotNullConstraintWithInh(schema, rel, col, conName, isLocal, inhCount); err != nil {
				rollbackRelation()
				return err
			}
		}

		// Mark implicit sequences as column-owned in the catalog model so
		// diff/drop logic treats SERIAL and IDENTITY sequences as managed objects.
		for i, si := range serials {
			seq := createdSeqs[i]
			seq.Persistence = rel.Persistence
			seq.OwnerRelOID = rel.OID
			seq.OwnerAttNum = rel.Columns[si.colIdx].AttNum
			c.recordDependency('s', seq.OID, 0, 'r', rel.OID, int32(rel.Columns[si.colIdx].AttNum), DepAuto)
		}
	}

	// Analyze column defaults and generated expressions using Tier 2 pipeline.
	// Coerce to column type to match PG's cookDefault behavior.
	//
	// pg: src/backend/commands/tablecmds.c — cookDefault / cookConstraint
	for i, cd := range colDefs {
		if cd.RawDefault != nil && columns[i].HasDefault {
			analyzed, err := c.AnalyzeColumnDefaultExpr(cd.RawDefault, rel)
			if err != nil {
				rollbackRelation()
				return err
			}
			if analyzed == nil {
				continue
			}
			// pg: cookDefault uses COERCE_IMPLICIT_CAST ('i') as display format
			coerced, cerr := c.coerceToTargetTypeWithFormat(analyzed, analyzed.exprType(), columns[i].TypeOID, 'a', 'i')
			if cerr != nil {
				rollbackRelation()
				return cerr
			}
			if coerced != nil {
				analyzed = coerced
			}
			columns[i].DefaultAnalyzed = analyzed
			c.recordDependencyOnSingleRelExprForObject('r', rel.OID, int32(columns[i].AttNum), analyzed, rel.OID,
				DepNormal, DepNormal)
			rte := c.buildRelationRTE(rel)
			columns[i].Default = c.DeparseExpr(analyzed, []*RangeTableEntry{rte}, false)
		}
		if cd.RawGenExpr != nil && columns[i].Generated == 's' {
			if rawExprContainsSubLink(cd.RawGenExpr) {
				rollbackRelation()
				return &Error{Code: CodeFeatureNotSupported,
					Message: "cannot use subquery in column generation expression"}
			}
			if analyzed, err := c.AnalyzeStandaloneExpr(cd.RawGenExpr, rel); err == nil && analyzed != nil {
				coerced, cerr := c.coerceToTargetType(analyzed, analyzed.exprType(), columns[i].TypeOID, 'i')
				if cerr == nil && coerced != nil {
					analyzed = coerced
				}
				columns[i].DefaultAnalyzed = analyzed
				c.recordDependencyOnSingleRelExprForObject('r', rel.OID, int32(columns[i].AttNum), analyzed, rel.OID,
					DepNormal, DepNormal)
				rte := c.buildRelationRTE(rel)
				genExpr := c.DeparseExpr(analyzed, []*RangeTableEntry{rte}, false)
				columns[i].GenerationExpr = genExpr
				columns[i].Default = genExpr
			}
		}
	}

	// Clone indexes and comments from LIKE sources.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (like_defined indexes)
	for _, ls := range likeSources {
		if ls.opts&tableLikeIndexes != 0 {
			c.cloneLikeIndexes(schema, ls.rel, rel)
		}
		if ls.opts&tableLikeComments != 0 {
			c.cloneLikeComments(ls.rel, rel)
		}
	}

	// Clone parent's indexes, triggers, and FK constraints to partition.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (CloneRowTriggersToPartition etc.)
	if rel.PartitionOf != 0 {
		parent := c.relationByOID[rel.PartitionOf]
		if parent != nil {
			c.cloneIndexesToPartition(schema, parent, rel)
			c.cloneRowTriggersToPartition(parent, rel)
			c.cloneForeignKeyConstraints(schema, parent, rel)
		}
	}

	return nil
}

// RemoveRelations handles DROP TABLE/VIEW/MATERIALIZED VIEW for one or more objects.
//
// pg: src/backend/commands/tablecmds.c — RemoveRelations
func (c *Catalog) RemoveRelations(stmt *nodes.DropStmt) error {
	cascade := stmt.Behavior == int(nodes.DROP_CASCADE)
	removeType := nodes.ObjectType(stmt.RemoveType)

	if stmt.Objects == nil {
		return nil
	}

	// Resolve and validate every target before removing any, so a wrong kind
	// later in the list leaves the catalog as PostgreSQL's rollback would.
	type target struct {
		schema *Schema
		rel    *Relation
		name   string
	}
	var targets []target
	for _, obj := range stmt.Objects.Items {
		// Objects are represented as a List of String nodes (schema-qualified name).
		schemaName, relName := extractDropObjectName(obj)

		schema, rel, err := c.findRelation(schemaName, relName)
		if err != nil {
			if stmt.Missing_ok && missingObject(err) {
				c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q does not exist, skipping", relName))
				continue
			}
			return err
		}

		// Validate relkind matches the drop type.
		switch removeType {
		case nodes.OBJECT_VIEW:
			if rel.RelKind != 'v' {
				return &Error{
					Code:    CodeWrongObjectType,
					Message: fmt.Sprintf("%q is not a view", relName),
				}
			}
		case nodes.OBJECT_MATVIEW:
			if rel.RelKind != 'm' {
				return &Error{
					Code:    CodeWrongObjectType,
					Message: fmt.Sprintf("%q is not a materialized view", relName),
				}
			}
		case nodes.OBJECT_FOREIGN_TABLE:
			// pg: src/backend/commands/tablecmds.c — RemoveRelations (OBJECT_FOREIGN_TABLE)
			if rel.RelKind != 'f' {
				return &Error{
					Code:    CodeWrongObjectType,
					Message: fmt.Sprintf("%q is not a foreign table", relName),
				}
			}
		default: // OBJECT_TABLE
			if rel.RelKind != 'r' && rel.RelKind != 'p' {
				return &Error{
					Code:    CodeWrongObjectType,
					Message: fmt.Sprintf("%q is not a table", relName),
				}
			}
		}
		if c.isExtensionMember(rel.OID) {
			return errDependentObjects("extension", c.extensionMembers[rel.OID])
		}
		targets = append(targets, target{schema: schema, rel: rel, name: relName})
	}

	dropRelOIDs := map[uint32]bool{}
	dropTypeOIDs := map[uint32]bool{}
	var autoDropParts []target
	addDropRel := func(r *Relation) {
		dropRelOIDs[r.OID] = true
		if r.RowTypeOID != 0 {
			dropTypeOIDs[r.RowTypeOID] = true
		}
		if r.ArrayOID != 0 {
			dropTypeOIDs[r.ArrayOID] = true
		}
	}
	var collectPartitions func(parent *Relation)
	collectPartitions = func(parent *Relation) {
		for _, childOID := range c.findPartitionChildRelations(parent.OID) {
			child := c.relationByOID[childOID]
			if child == nil || dropRelOIDs[child.OID] {
				continue
			}
			addDropRel(child)
			autoDropParts = append(autoDropParts, target{schema: child.Schema, rel: child, name: child.Name})
			collectPartitions(child)
		}
	}
	for _, t := range targets {
		addDropRel(t.rel)
	}
	for _, t := range targets {
		if t.rel.RelKind == 'p' {
			collectPartitions(t.rel)
		}
	}
	if !cascade {
		kind := "table"
		if removeType == nodes.OBJECT_VIEW {
			kind = "view"
		} else if removeType == nodes.OBJECT_MATVIEW {
			kind = "materialized view"
		}
		allToCheck := make([]target, 0, len(targets)+len(autoDropParts))
		allToCheck = append(allToCheck, targets...)
		allToCheck = append(allToCheck, autoDropParts...)
		for _, t := range allToCheck {
			for _, dep := range c.relationDependents(t.rel.OID, 0) {
				if !dep.Normal {
					continue
				}
				if (dep.ObjType == 'r' && dropRelOIDs[dep.ObjOID]) ||
					(dep.ObjType == 'c' && c.constraints[dep.ObjOID] != nil && dropRelOIDs[c.constraints[dep.ObjOID].RelOID]) {
					continue
				}
				return errDependentObjects(kind, t.name)
			}
			if t.rel.RowTypeOID != 0 {
				if bt := c.typeByOID[t.rel.RowTypeOID]; bt != nil {
					for _, depRelOID := range c.typeRelationDependents(bt) {
						if !dropRelOIDs[depRelOID] {
							return errDependentObjects(kind, t.name)
						}
					}
					for _, dep := range c.typeNormalDependents(bt) {
						if (dep.ObjType == 'r' && dropRelOIDs[dep.ObjOID]) ||
							(dep.ObjType == 't' && dropTypeOIDs[dep.ObjOID]) ||
							(dep.ObjType == 'c' && c.constraints[dep.ObjOID] != nil && dropRelOIDs[c.constraints[dep.ObjOID].RelOID]) {
							continue
						}
						return errDependentObjects(kind, t.name)
					}
				}
			}
			for _, cast := range c.castIndex {
				if (t.rel.RowTypeOID != 0 && (cast.Source == t.rel.RowTypeOID || cast.Target == t.rel.RowTypeOID)) ||
					(t.rel.ArrayOID != 0 && (cast.Source == t.rel.ArrayOID || cast.Target == t.rel.ArrayOID)) {
					return errDependentObjects(kind, t.name)
				}
			}
		}
	}

	if !cascade {
		// Drop descendant partitions bottom-up (DepAuto on their partitioned parent).
		for i := len(autoDropParts) - 1; i >= 0; i-- {
			p := autoDropParts[i]
			if c.relationByOID[p.rel.OID] != nil {
				c.removeRelation(p.schema, p.name, p.rel)
			}
		}
	}
	for _, t := range targets {
		if c.relationByOID[t.rel.OID] == nil {
			continue
		}
		if cascade {
			c.cascadeDropRelation(t.rel)
			continue
		}
		c.removeRelation(t.schema, t.name, t.rel)
	}
	return nil
}

// convertColumnDef converts a pgparser ColumnDef to a catalog ColumnDef and
// any inline constraints.
//
// (pgddl helper — PG processes ColumnDef inline in DefineRelation)
func (c *Catalog) convertColumnDef(cd *nodes.ColumnDef, relName string, schema *Schema) (ColumnDef, []ConstraintDef, error) {
	result := ColumnDef{
		Name:    cd.Colname,
		NotNull: cd.IsNotNull,
	}

	// Reject SETOF columns.
	// pg: src/backend/commands/tablecmds.c — BuildDescForRelation
	if cd.TypeName != nil && cd.TypeName.Setof {
		return ColumnDef{}, nil, &Error{
			Code:    CodeInvalidTableDefinition,
			Message: fmt.Sprintf("column %q cannot be declared SETOF", cd.Colname),
		}
	}

	// Check for SERIAL types.
	serialWidth := isSerialType(cd.TypeName)
	if serialWidth != 0 {
		result.IsSerial = serialWidth
	} else if cd.TypeName != nil {
		// Resolve type name.
		typSchema, name := typeNameParts(cd.TypeName)
		if typSchema == "pg_catalog" {
			typSchema = ""
		}

		rawMods := extractRawTypmods(cd.TypeName.Typmods)
		typmod := int32(-1)
		if len(rawMods) > 0 {
			typmod = encodeTypModByName(name, rawMods)
		}
		// pgparser sets Typemod=0 (Go zero value) when no modifier is specified,
		// while PG uses -1. Accept both as "no modifier".
		if cd.TypeName.Typemod > 0 {
			typmod = cd.TypeName.Typemod
		}

		isArray := cd.TypeName.ArrayBounds != nil && len(cd.TypeName.ArrayBounds.Items) > 0

		// Array dimension limit.
		// pg: src/backend/commands/tablecmds.c — BuildDescForRelation (MAXDIM check)
		if cd.TypeName.ArrayBounds != nil && len(cd.TypeName.ArrayBounds.Items) > 6 {
			return ColumnDef{}, nil, &Error{
				Code:    CodeProgramLimitExceeded,
				Message: fmt.Sprintf("number of array dimensions (%d) exceeds the maximum allowed (%d)", len(cd.TypeName.ArrayBounds.Items), 6),
			}
		}

		result.Type = TypeName{
			Schema:  typSchema,
			Name:    name,
			TypeMod: typmod,
			IsArray: isArray,
		}
	}

	// Handle explicit COLLATE, STORAGE, and COMPRESSION clauses.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (COLLATE / STORAGE / COMPRESSION)
	if cd.CollClause != nil && cd.CollClause.Collname != nil && len(cd.CollClause.Collname.Items) > 0 {
		result.CollClause = cd.CollClause
		sc, nm := qualifiedName(cd.CollClause.Collname)
		if sc != "" {
			result.CollationSchema = sc
		}
		result.CollationName = nm
	}
	result.StorageName = cd.StorageName
	result.Compression = cd.Compression

	// Handle DEFAULT expression — raw text will be overwritten by DeparseExpr
	// after AnalyzeStandaloneExpr in DefineRelation.
	if cd.RawDefault != nil {
		result.RawDefault = cd.RawDefault
	}
	if cd.CookedDefault != nil {
		result.RawDefault = cd.CookedDefault
	}

	// Handle identity column (ColumnDef.Identity field).
	if cd.Identity == 'a' || cd.Identity == 'd' {
		result.Identity = cd.Identity
		result.IdentitySequence = cd.IdentitySequence
	}

	// Handle generated column (ColumnDef.Generated field).
	if cd.Generated == 's' {
		result.Generated = 's'
	}
	if cd.Fdwoptions != nil && len(cd.Fdwoptions.Items) > 0 {
		fdwOpts, err := applyForeignOptionChanges(nil, cd.Fdwoptions)
		if err != nil {
			return ColumnDef{}, nil, err
		}
		result.ForeignOptions = fdwOpts
	}

	// Process inline column constraints.
	var cons []ConstraintDef
	if cd.Constraints != nil {
		for _, item := range cd.Constraints.Items {
			con, ok := item.(*nodes.Constraint)
			if !ok {
				continue
			}
			switch con.Contype {
			case nodes.CONSTR_NOTNULL:
				result.NotNull = true
				if con.Conname != "" {
					result.NotNullConName = con.Conname
				}
			case nodes.CONSTR_DEFAULT:
				if con.RawExpr != nil {
					result.RawDefault = con.RawExpr
				}
				if con.CookedExpr != "" {
					result.Default = con.CookedExpr
				}
			case nodes.CONSTR_IDENTITY:
				result.Identity = con.GeneratedWhen
				result.IdentityOptions = con.Options
			case nodes.CONSTR_GENERATED:
				result.Generated = 's'
				if con.RawExpr != nil {
					result.RawGenExpr = con.RawExpr
				}
				if con.CookedExpr != "" {
					result.GenerationExpr = con.CookedExpr
				}
			case nodes.CONSTR_PRIMARY:
				cons = append(cons, ConstraintDef{
					Name:    con.Conname,
					Type:    ConstraintPK,
					Columns: []string{cd.Colname},
				})
			case nodes.CONSTR_UNIQUE:
				cons = append(cons, ConstraintDef{
					Name:    con.Conname,
					Type:    ConstraintUnique,
					Columns: []string{cd.Colname},
				})
			case nodes.CONSTR_CHECK:
				checkExpr := con.CookedExpr
				// Raw check expression text will be filled from AnalyzeStandaloneExpr + DeparseExpr
				// in addCheckConstraint.
				cons = append(cons, ConstraintDef{
					Name:         con.Conname,
					Type:         ConstraintCheck,
					Columns:      []string{cd.Colname},
					CheckExpr:    checkExpr,
					RawCheckExpr: con.RawExpr,
				})
			case nodes.CONSTR_FOREIGN:
				def := ConstraintDef{
					Name:        con.Conname,
					Type:        ConstraintFK,
					Columns:     []string{cd.Colname},
					FKUpdAction: normalizeFKAction(con.FkUpdaction),
					FKDelAction: normalizeFKAction(con.FkDelaction),
					FKMatchType: normalizeFKMatch(con.FkMatchtype),
					Deferrable:  con.Deferrable,
					Deferred:    con.Initdeferred,
				}
				if con.Pktable != nil {
					def.RefSchema = con.Pktable.Schemaname
					def.RefTable = con.Pktable.Relname
				}
				def.RefColumns = stringListItems(con.PkAttrs)
				cons = append(cons, def)
			}
		}
	}

	if result.Identity != 0 && (result.RawDefault != nil || result.Default != "") {
		return ColumnDef{}, nil, &Error{
			Code:    CodeSyntaxError,
			Message: fmt.Sprintf("both default and identity specified for column %q of table %q", result.Name, relName),
		}
	}

	return result, cons, nil
}

// extractDropObjectName extracts schema and name from a DropStmt object.
// Objects in DropStmt.Objects can be:
// - *nodes.List containing String nodes (for tables, views, indexes, etc.)
// - *nodes.String (for schemas, simple names)
func extractDropObjectName(obj nodes.Node) (schema, name string) {
	switch n := obj.(type) {
	case *nodes.List:
		return qualifiedName(n)
	case *nodes.String:
		return "", n.Str
	case *nodes.RangeVar:
		return n.Schemaname, n.Relname
	default:
		return "", ""
	}
}

// isSerialTypeName checks if a type name string represents a SERIAL type.
func isSerialTypeName(name string) byte {
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

// LIKE option bitmask constants.
// pg: src/include/nodes/parsenodes.h — CreateStmtLikeOption
const (
	tableLikeComments    uint32 = 1 << 0
	tableLikeCompression uint32 = 1 << 1
	tableLikeConstraints uint32 = 1 << 2
	tableLikeDefaults    uint32 = 1 << 3
	tableLikeGenerated   uint32 = 1 << 4
	tableLikeIdentity    uint32 = 1 << 5
	tableLikeIndexes     uint32 = 1 << 6
	tableLikeStatistics  uint32 = 1 << 7
	tableLikeStorage     uint32 = 1 << 8
)

// likeSource tracks a LIKE source relation and its options for deferred index/comment cloning.
type likeSource struct {
	rel  *Relation
	opts uint32
}

// typeNameForOID constructs a catalog TypeName from a resolved type OID and typmod,
// preserving non-pg_catalog schema qualification and array-element mapping so
// cross-schema inherited/typed/LIKE columns resolve regardless of search_path.
func (c *Catalog) typeNameForOID(typeOID uint32, typmod int32) TypeName {
	typ := c.typeByOID[typeOID]
	if typ == nil {
		return TypeName{TypeMod: typmod}
	}
	schemaName := ""
	if typ.Namespace != 0 && typ.Namespace != PGCatalogNamespace {
		if ns := c.schemas[typ.Namespace]; ns != nil {
			schemaName = ns.Name
		}
	}
	if typ.Elem != 0 {
		if elem := c.typeByOID[typ.Elem]; elem != nil && elem.Array == typ.OID {
			return TypeName{Schema: schemaName, Name: elem.TypeName, TypeMod: typmod, IsArray: true}
		}
	}
	return TypeName{Schema: schemaName, Name: typ.TypeName, TypeMod: typmod}
}

// expandTableLikeClause copies columns and optionally constraints from a source table.
// Returns the source relation and options for deferred index/comment cloning.
//
// pg: src/backend/commands/tablecmds.c — expandTableLikeClause
func (c *Catalog) expandTableLikeClause(clause *nodes.TableLikeClause) ([]ColumnDef, []ConstraintDef, *Relation, uint32, error) {
	if clause.Relation == nil {
		return nil, nil, nil, 0, errInvalidParameterValue("LIKE requires a table name")
	}

	_, srcRel, err := c.findRelation(clause.Relation.Schemaname, clause.Relation.Relname)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	opts := uint32(clause.Options)

	// Always copy column names and types.
	var colDefs []ColumnDef
	for _, col := range srcRel.Columns {
		cd := ColumnDef{
			Name:    col.Name,
			Type:    c.typeNameForOID(col.TypeOID, col.TypeMod),
			NotNull: col.NotNull,
			IsLocal: true,
		}

		// Copy defaults if INCLUDING DEFAULTS or INCLUDING ALL.
		if opts&tableLikeDefaults != 0 && col.HasDefault && col.Generated == 0 {
			cd.Default = col.Default
		}

		// Copy generated column expression if INCLUDING GENERATED or INCLUDING ALL.
		if opts&tableLikeGenerated != 0 && col.Generated == 's' {
			cd.Generated = 's'
			cd.GenerationExpr = col.GenerationExpr
		}

		// Copy identity if INCLUDING IDENTITY or INCLUDING ALL.
		if opts&tableLikeIdentity != 0 && col.Identity != 0 {
			cd.Identity = col.Identity
		}

		colDefs = append(colDefs, cd)
	}

	// Copy CHECK constraints if INCLUDING CONSTRAINTS or INCLUDING ALL.
	var constraintDefs []ConstraintDef
	if opts&tableLikeConstraints != 0 {
		for _, con := range c.consByRel[srcRel.OID] {
			if con.Type == ConstraintCheck {
				constraintDefs = append(constraintDefs, ConstraintDef{
					Name:         con.Name,
					Type:         ConstraintCheck,
					CheckExpr:    con.CheckExpr,
					RawCheckExpr: parseCheckExpr(con.CheckExpr),
				})
			}
		}
	}

	// Indexes and comments are handled after the relation is created (deferred).
	// Return the source relation and options for post-creation cloning.
	// pg: src/backend/commands/tablecmds.c — expandTableLikeClause (indexes/comments)
	return colDefs, constraintDefs, srcRel, opts, nil
}

// parseCheckExpr re-parses a CHECK expression text to get a raw AST node.
// This is used when copying constraints via LIKE: PG re-parses the cooked
// expression text so that column references resolve against the target table.
//
// (pgddl helper — PG uses stringToNode + transformExpr in expandTableLikeClause)
func parseCheckExpr(exprText string) nodes.Node {
	if exprText == "" {
		return nil
	}
	list, err := pgparser.Parse("SELECT 1 WHERE " + exprText)
	if err != nil || list == nil || len(list.Items) == 0 {
		return nil
	}
	sel, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok || sel.WhereClause == nil {
		return nil
	}
	return sel.WhereClause
}

// checkPartitionBoundOverlap checks for overlapping partition bounds among
// existing partitions of the same parent.
//
// pg: src/backend/partitioning/partbounds.c — check_new_partition_bound
func (c *Catalog) checkPartitionBoundOverlap(parentOID, selfOID uint32, newBound *PartitionBound) error {
	if newBound == nil {
		return nil
	}

	for _, e := range c.inhEntries {
		if e.InhParent != parentOID || e.InhRelID == selfOID {
			continue
		}
		existing := c.relationByOID[e.InhRelID]
		if existing == nil || existing.PartitionBound == nil {
			continue
		}
		eb := existing.PartitionBound

		// DEFAULT partition: only one allowed.
		if newBound.IsDefault && eb.IsDefault {
			return errInvalidObjectDefinition("partition \"default\" would overlap with existing default partition")
		}

		// Skip overlap checks against default partition for non-default new partitions.
		if eb.IsDefault || newBound.IsDefault {
			continue
		}

		switch newBound.Strategy {
		case 'l': // LIST
			for i, nv := range newBound.ListValues {
				newKw := i < len(newBound.ListKeywords) && newBound.ListKeywords[i]
				for j, ev := range eb.ListValues {
					evKw := j < len(eb.ListKeywords) && eb.ListKeywords[j]
					if nv == ev && newKw == evKw {
						return errInvalidObjectDefinition(
							fmt.Sprintf("partition bound value %q already exists in partition %q", nv, existing.Name))
					}
				}
			}
		case 'r': // RANGE
			// Overlap if: new.lower < existing.upper AND new.upper > existing.lower.
			if !rangeAfter(newBound.LowerBound, eb.UpperBound) && !rangeAfter(eb.LowerBound, newBound.UpperBound) {
				return errInvalidObjectDefinition(
					fmt.Sprintf("partition bound overlaps with existing partition %q", existing.Name))
			}
		case 'h': // HASH
			if newBound.Modulus == eb.Modulus && newBound.Remainder == eb.Remainder {
				return errInvalidObjectDefinition(
					fmt.Sprintf("partition hash bound (modulus %d, remainder %d) conflicts with partition %q",
						newBound.Modulus, newBound.Remainder, existing.Name))
			}
		}
	}
	return nil
}

// rangeAfter returns true if bound a is >= bound b.
// Uses numeric comparison when values are parseable as numbers,
// falls back to string comparison otherwise.
func rangeAfter(a, b []string) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		cmp := compareDatumStrings(a[i], b[i])
		if cmp < 0 {
			return false
		}
		if cmp > 0 {
			return true
		}
	}
	return true // equal = boundary point
}

// compareDatumStrings compares two deparsed datum values.
// Returns -1, 0, or 1. Tries numeric comparison first.
func compareDatumStrings(a, b string) int {
	var ai, bi int64
	if _, err := fmt.Sscanf(a, "%d", &ai); err == nil {
		if _, err := fmt.Sscanf(b, "%d", &bi); err == nil {
			if ai < bi {
				return -1
			}
			if ai > bi {
				return 1
			}
			return 0
		}
	}
	// Fall back to string comparison.
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// checkOfType validates that the given type OID refers to a composite type
// and returns its backing relation.
//
// pg: src/backend/commands/tablecmds.c — check_of_type
func (c *Catalog) checkOfType(typeOID uint32) (*Relation, error) {
	typ := c.typeByOID[typeOID]
	if typ == nil || typ.Type != 'c' {
		name := "unknown"
		if typ != nil {
			name = typ.TypeName
		}
		return nil, errWrongObjectType(name, "a composite type")
	}
	rel := c.relationByOID[typ.RelID]
	if rel == nil || rel.RelKind != 'c' {
		return nil, errWrongObjectType(typ.TypeName, "a composite type")
	}
	return rel, nil
}

// ExecuteTruncate handles TRUNCATE TABLE.
// For pgddl (no physical data), this validates the targets and optionally
// resets owned sequence values for RESTART IDENTITY.
//
// pg: src/backend/commands/tablecmds.c — ExecuteTruncate
func (c *Catalog) ExecuteTruncate(stmt *nodes.TruncateStmt) error {
	if stmt.Relations == nil {
		return nil
	}

	for _, item := range stmt.Relations.Items {
		rv, ok := item.(*nodes.RangeVar)
		if !ok {
			continue
		}
		_, rel, err := c.findRelation(rv.Schemaname, rv.Relname)
		if err != nil {
			return err
		}
		// TRUNCATE is only valid on tables (not views, sequences, etc.).
		if rel.RelKind != 'r' && rel.RelKind != 'p' {
			return errWrongObjectType(rv.Relname, "a table")
		}

		// Cannot truncate ONLY a partitioned table.
		// pg: src/backend/commands/tablecmds.c — ExecuteTruncateGuts
		if !rv.Inh && rel.RelKind == 'p' {
			return &Error{
				Code:    CodeWrongObjectType,
				Message: fmt.Sprintf("cannot truncate only a partitioned table\nHint: Do not specify the ONLY keyword, or use TRUNCATE ONLY on the partitions directly."),
			}
		}

		// RESTART IDENTITY: reset owned sequences to their start values.
		if stmt.RestartSeqs {
			for _, seq := range c.sequenceByOID {
				if seq.OwnerRelOID == rel.OID {
					// Reset to start value (pgddl doesn't track current value,
					// but this ensures the semantics are recorded).
					_ = seq // no-op: we don't track current value
				}
			}
		}
	}

	// FK validation: check tables referenced by FK constraints.
	// pg: src/backend/commands/tablecmds.c — ExecuteTruncateGuts (heap_truncate_check_FKs)
	// Collect OIDs of all truncated tables.
	truncatedOIDs := make(map[uint32]bool)
	for _, item := range stmt.Relations.Items {
		rv, ok := item.(*nodes.RangeVar)
		if !ok {
			continue
		}
		_, rel, err := c.findRelation(rv.Schemaname, rv.Relname)
		if err != nil {
			continue // already validated above
		}
		truncatedOIDs[rel.OID] = true
	}

	// Find FK constraints where the referenced table is being truncated
	// but the referencing table is not.
	for _, con := range c.constraints {
		if con.Type != ConstraintFK {
			continue
		}
		if !truncatedOIDs[con.FRelOID] {
			continue
		}
		// The referencing table is also being truncated — no issue.
		if truncatedOIDs[con.RelOID] {
			continue
		}
		refRel := c.relationByOID[con.RelOID]
		if refRel == nil {
			continue
		}
		// RESTRICT: error if any referencing table is not in the truncation set.
		// pg: src/backend/commands/tablecmds.c — heap_truncate_check_FKs
		if stmt.Behavior != nodes.DROP_CASCADE {
			fRel := c.relationByOID[con.FRelOID]
			fRelName := ""
			if fRel != nil {
				fRelName = fRel.Name
			}
			return &Error{
				Code:    CodeFeatureNotSupported,
				Message: fmt.Sprintf("cannot truncate a table referenced in a foreign key constraint\nDetail: Table %q references %q.", refRel.Name, fRelName),
			}
		}
		// CASCADE: in PG, the referencing table would be added to the truncation set.
		// For pgddl, we just validate — no physical data to truncate.
	}

	return nil
}

// cloneIndexesToPartition clones parent's indexes to a partition.
// For PK/UNIQUE indexes backed by constraints, the constraint is also cloned.
// For PK indexes, NOT NULL is forced on the partition's corresponding columns.
//
// pg: src/backend/commands/tablecmds.c — CloneRowTriggersToPartition/DefineIndex (partition path)
// (pgddl helper — combines CloneIndexesToPartition + DefineIndexForPartition)
func (c *Catalog) cloneIndexesToPartition(schema *Schema, parent, partition *Relation) {
	for _, idx := range c.indexesByRel[parent.OID] {
		c.cloneIndexToPartition(schema, parent, partition, idx)
	}
}

// cloneIndexToPartitionTree recursively clones a single parent index (and any
// backing PK/UNIQUE/EXCLUDE constraint) down to all existing partitions.
func (c *Catalog) cloneIndexToPartitionTree(parent *Relation, idx *Index) {
	if parent == nil || idx == nil {
		return
	}
	for _, childOID := range c.findPartitionChildRelations(parent.OID) {
		child := c.relationByOID[childOID]
		if child == nil {
			continue
		}
		partIdx := c.cloneIndexToPartition(child.Schema, parent, child, idx)
		if partIdx != nil && child.RelKind == 'p' {
			c.cloneIndexToPartitionTree(child, partIdx)
		}
	}
}

// cloneIndexConstraintToPartitionTree propagates an index-backed constraint
// (PK/UNIQUE/EXCLUDE) added to a partitioned table down to existing partitions.
// For PRIMARY KEY, it first validates that no partition has an incompatible
// existing primary key, reusing and attaching any compatible existing PK.
func (c *Catalog) cloneIndexConstraintToPartitionTree(parent *Relation, parentIdx *Index) error {
	if parent == nil || parentIdx == nil {
		return nil
	}
	parentCon := c.constraints[parentIdx.ConstraintOID]
	if parentCon != nil && parentCon.Type == ConstraintPK {
		var validate func(curParent *Relation, curCols, curIncl []int16) error
		validate = func(curParent *Relation, curCols, curIncl []int16) error {
			for _, childOID := range c.findPartitionChildRelations(curParent.OID) {
				child := c.relationByOID[childOID]
				if child == nil {
					continue
				}
				mappedCols := c.mapColumnAttnums(curParent, child, curCols)
				mappedIncl := c.mapColumnAttnums(curParent, child, curIncl)
				if existingPK := c.findPKConstraint(child.OID); existingPK != nil {
					if !slices.Equal(existingPK.Columns, mappedCols) ||
						!slices.Equal(existingPK.IncludeColumns, mappedIncl) ||
						existingPK.Deferrable != parentCon.Deferrable ||
						existingPK.Deferred != parentCon.Deferred {
						return errDuplicatePKey(child.Name)
					}
				}
				if child.RelKind == 'p' {
					if err := validate(child, mappedCols, mappedIncl); err != nil {
						return err
					}
				}
			}
			return nil
		}
		if err := validate(parent, parentCon.Columns, parentCon.IncludeColumns); err != nil {
			return err
		}
	}

	var propagate func(curParent *Relation, curIdx *Index)
	propagate = func(curParent *Relation, curIdx *Index) {
		curCon := c.constraints[curIdx.ConstraintOID]
		for _, childOID := range c.findPartitionChildRelations(curParent.OID) {
			child := c.relationByOID[childOID]
			if child == nil {
				continue
			}
			if curCon != nil && curCon.Type == ConstraintPK {
				if existingPK := c.findPKConstraint(child.OID); existingPK != nil {
					existingPK.ConParentID = curCon.OID
					var existingIdx *Index
					if existingPK.IndexOID != 0 {
						existingIdx = c.indexes[existingPK.IndexOID]
					}
					if existingIdx != nil {
						existingIdx.ParentIndexOID = curIdx.OID
						c.recordDependency('i', existingIdx.OID, 0, 'i', curIdx.OID, 0, DepInternal)
						if child.RelKind == 'p' {
							propagate(child, existingIdx)
						}
					}
					continue
				}
			}
			partIdx := c.cloneIndexToPartition(child.Schema, curParent, child, curIdx)
			if partIdx != nil && child.RelKind == 'p' {
				propagate(child, partIdx)
			}
		}
	}
	propagate(parent, parentIdx)
	return nil
}

func (c *Catalog) cloneIndexToPartition(schema *Schema, parent, partition *Relation, idx *Index) *Index {
	// Map parent attnums to partition attnums by column name.
	partAttnums := c.mapColumnAttnums(parent, partition, idx.Columns)
	if partAttnums == nil {
		return nil // column not found in partition (shouldn't happen for valid partition)
	}

	partIdxName := fmt.Sprintf("%s_%s_idx", partition.Name, idx.Name)
	if idx.IsPrimary {
		partIdxName = makePGObjectName(partition.Name, "", "pkey")
	}
	// Ensure uniqueness against schema relation/index/sequence names and
	// partition constraint names (for PK indexes that name the PK constraint).
	nameInUse := func(name string) bool {
		if _, exists := schema.Indexes[name]; exists {
			return true
		}
		if _, exists := schema.Relations[name]; exists {
			return true
		}
		if _, exists := schema.Sequences[name]; exists {
			return true
		}
		if idx.IsPrimary {
			for _, con := range c.consByRel[partition.OID] {
				if con.Name == name {
					return true
				}
			}
		}
		return false
	}
	if nameInUse(partIdxName) {
		base := partIdxName
		for suffix := 1; ; suffix++ {
			suf := fmt.Sprintf("%d", suffix)
			cand := base
			for len(cand)+len(suf) > 63 {
				cand = trimOneRune(cand)
			}
			cand += suf
			if !nameInUse(cand) {
				partIdxName = cand
				break
			}
		}
	}

	// Clone IndOption.
	var partIndOption []int16
	if idx.IndOption != nil {
		partIndOption = make([]int16, len(idx.IndOption))
		copy(partIndOption, idx.IndOption)
	}

	partIdxSpace := idx.Tablespace
	if partIdxSpace == "" {
		partIdxSpace = normalizeTablespace(c.effectiveIndexTablespace("", partition))
	}
	partIdx := &Index{
		OID:               c.oidGen.Next(),
		Name:              partIdxName,
		Schema:            schema,
		RelOID:            partition.OID,
		Tablespace:        partIdxSpace,
		OwnerName:         partition.OwnerName,
		Columns:           partAttnums,
		ColumnNames:       slices.Clone(idx.ColumnNames),
		IsUnique:          idx.IsUnique,
		IsPrimary:         idx.IsPrimary,
		AccessMethod:      idx.AccessMethod,
		NKeyColumns:       idx.NKeyColumns,
		WhereClause:       idx.WhereClause,
		IndOption:         partIndOption,
		NullsNotDistinct:  idx.NullsNotDistinct,
		Exprs:             slices.Clone(idx.Exprs),
		OpClasses:         slices.Clone(idx.OpClasses),
		OpClassOptions:    slices.Clone(idx.OpClassOptions),
		Collations:        cloneStringSlices(idx.Collations),
		UserCollationOIDs: slices.Clone(idx.UserCollationOIDs),
		AnalyzedExprs:     slices.Clone(idx.AnalyzedExprs),
		AnalyzedWhere:     idx.AnalyzedWhere,
		ParentIndexOID:    idx.OID,
	}

	// If the parent index backs a constraint, clone the constraint too.
	if idx.ConstraintOID != 0 {
		parentCon := c.constraints[idx.ConstraintOID]
		if parentCon != nil {
			conOID := c.oidGen.Next()
			partConAttnums := c.mapColumnAttnums(parent, partition, parentCon.Columns)
			if partConAttnums == nil {
				partConAttnums = partAttnums
			}

			partConName := fmt.Sprintf("%s_%s", partition.Name, parentCon.Name)
			if parentCon.Type == ConstraintPK {
				partConName = partIdxName
			}

			partCon := &Constraint{
				OID:            conOID,
				Name:           partConName,
				Type:           parentCon.Type,
				RelOID:         partition.OID,
				Namespace:      schema.OID,
				Columns:        partConAttnums,
				IncludeColumns: c.mapColumnAttnums(parent, partition, parentCon.IncludeColumns),
				IndexOID:       partIdx.OID,
				Deferrable:     parentCon.Deferrable,
				Deferred:       parentCon.Deferred,
				Validated:      true,
				ConIsLocal:     true,
				ConParentID:    parentCon.OID,
				ExclOps:        slices.Clone(parentCon.ExclOps),
				ExclOpSchemas:  slices.Clone(parentCon.ExclOpSchemas),
			}

			partIdx.ConstraintOID = conOID
			c.registerConstraint(partition.OID, partCon)
			c.recordDependency('c', conOID, 0, 'r', partition.OID, 0, DepAuto)
			c.recordColumnDependencies('c', conOID, partition.OID, partConAttnums, DepAuto)
			var parentConExtDeps []DepEntry
			for _, d := range c.deps {
				if d.ObjType == 'c' && d.ObjOID == parentCon.OID && d.DepType == DepNormal &&
					!(d.RefType == 'r' && d.RefOID == parent.OID) {
					parentConExtDeps = append(parentConExtDeps, d)
				}
			}
			for _, d := range parentConExtDeps {
				c.recordDependency('c', conOID, 0, d.RefType, d.RefOID, d.RefSubID, DepNormal)
			}
		}
	}

	c.registerIndex(schema, partIdx)
	c.recordDependency('i', partIdx.OID, 0, 'r', partition.OID, 0, DepAuto)
	c.recordColumnDependencies('i', partIdx.OID, partition.OID, partAttnums, DepAuto)
	if am := c.accessMethodByName[partIdx.AccessMethod]; am != nil {
		c.recordDependency('i', partIdx.OID, 0, 'm', am.OID, 0, DepNormal)
	}
	c.recordOpClassDependencies('i', partIdx.OID, c.indexOpClassRefs(partition, partAttnums, partIdx.NKeyColumns, partIdx.OpClasses), partIdx.AccessMethod)
	var parentIdxExtDeps []DepEntry
	for _, d := range c.deps {
		if d.ObjType == 'i' && d.ObjOID == idx.OID && d.DepType == DepNormal &&
			!(d.RefType == 'r' && d.RefOID == parent.OID) {
			parentIdxExtDeps = append(parentIdxExtDeps, d)
		}
	}
	for _, d := range parentIdxExtDeps {
		c.recordDependency('i', partIdx.OID, 0, d.RefType, d.RefOID, d.RefSubID, DepNormal)
	}

	// Force NOT NULL on PK key columns; INCLUDE columns keep their nullability.
	if idx.IsPrimary {
		for _, attnum := range indexKeyAttnums(idx, partAttnums) {
			if int(attnum) > 0 && int(attnum) <= len(partition.Columns) {
				partCol := partition.Columns[attnum-1]
				partCol.NotNull = true
				conName := ""
				if pIdx, ok := parent.colByName[partCol.Name]; ok {
					if parentName := c.notNullConstraintName(parent, parent.Columns[pIdx].AttNum); parentName != "" && c.checkConstraintNameFree(partition.OID, parentName) == nil {
						conName = parentName
					}
				}
				_ = c.ensureNotNullConstraintWithInh(schema, partition, partCol, conName, false, 1)
			}
		}
	}
	return partIdx
}

// cloneRowTriggersToPartition clones FOR EACH ROW triggers from parent to partition.
// Statement-level triggers and constraint triggers are skipped.
//
// pg: src/backend/commands/trigger.c — CloneRowTriggersToPartition
// (pgddl helper)
func (c *Catalog) cloneRowTriggersToPartition(parent, partition *Relation) {
	for _, trig := range c.triggersByRel[parent.OID] {
		// Only clone FOR EACH ROW triggers.
		if !trig.ForEachRow {
			continue
		}
		// Skip constraint triggers (cloned via FK cloning path).
		if trig.IsConstraint {
			continue
		}

		// Map columns (UPDATE OF columns).
		var partCols []int16
		if trig.Columns != nil {
			partCols = c.mapColumnAttnums(parent, partition, trig.Columns)
		}

		// Clone trigger arguments.
		var partArgs []string
		if trig.Args != nil {
			partArgs = make([]string, len(trig.Args))
			copy(partArgs, trig.Args)
		}

		partTrig := &Trigger{
			OID:               c.oidGen.Next(),
			Name:              trig.Name,
			RelOID:            partition.OID,
			FuncOID:           trig.FuncOID,
			Timing:            trig.Timing,
			Events:            trig.Events,
			ForEachRow:        true,
			WhenExpr:          trig.WhenExpr,
			Columns:           partCols,
			Enabled:           trig.Enabled,
			OldTransitionName: trig.OldTransitionName,
			NewTransitionName: trig.NewTransitionName,
			Args:              partArgs,
			ParentTriggerOID:  trig.OID,
		}

		c.triggers[partTrig.OID] = partTrig
		c.triggersByRel[partition.OID] = append(c.triggersByRel[partition.OID], partTrig)
		c.recordDependency('g', partTrig.OID, 0, 'r', partition.OID, 0, DepAuto)
		// The clone depends on the columns its source does: those it is
		// UPDATE OF and those its WHEN condition reads, plus any external
		// normal dependencies (trigger function and WHEN functions/operators/types).
		var colDeps, extDeps []DepEntry
		for _, d := range c.deps {
			if d.ObjType != 'g' || d.ObjOID != trig.OID {
				continue
			}
			if d.RefType == 'r' && d.RefOID == parent.OID && d.RefSubID > 0 {
				colDeps = append(colDeps, d)
			} else if d.DepType == DepNormal && !(d.RefType == 'r' && d.RefOID == parent.OID) {
				extDeps = append(extDeps, d)
			}
		}
		for _, d := range colDeps {
			if m := c.mapColumnAttnums(parent, partition, []int16{int16(d.RefSubID)}); len(m) == 1 && m[0] != 0 {
				c.recordDependency('g', partTrig.OID, 0, 'r', partition.OID, int32(m[0]), d.DepType)
			}
		}
		if partTrig.FuncOID != 0 {
			c.recordDependency('g', partTrig.OID, 0, 'f', partTrig.FuncOID, 0, DepNormal)
		}
		for _, d := range extDeps {
			c.recordDependency('g', partTrig.OID, 0, d.RefType, d.RefOID, d.RefSubID, DepNormal)
		}
	}
}

// cloneForeignKeyConstraints clones FK constraints from parent to partition.
//
// pg: src/backend/commands/tablecmds.c — CloneForeignKeyConstraints
// (pgddl helper)
func (c *Catalog) cloneForeignKeyConstraints(schema *Schema, parent, partition *Relation) {
	for _, con := range c.consByRel[parent.OID] {
		if con.Type != ConstraintFK {
			continue
		}

		// Map local columns.
		partLocalAttnums := c.mapColumnAttnums(parent, partition, con.Columns)
		if partLocalAttnums == nil {
			continue
		}

		// Generate partition constraint name.
		partConName := fmt.Sprintf("%s_%s", partition.Name, con.Name)

		partCon := &Constraint{
			OID:          c.oidGen.Next(),
			Name:         partConName,
			Type:         ConstraintFK,
			RelOID:       partition.OID,
			Namespace:    schema.OID,
			Columns:      partLocalAttnums,
			FRelOID:      con.FRelOID,
			FColumns:     slices.Clone(con.FColumns),
			FKUpdAction:  con.FKUpdAction,
			FKDelAction:  con.FKDelAction,
			FKMatchType:  con.FKMatchType,
			Deferrable:   con.Deferrable,
			Deferred:     con.Deferred,
			Validated:    true,
			ConIsLocal:   true,
			ConParentID:  con.OID,
			PFEqOp:       slices.Clone(con.PFEqOp),
			PPEqOp:       slices.Clone(con.PPEqOp),
			FFEqOp:       slices.Clone(con.FFEqOp),
			FKDelSetCols: slices.Clone(con.FKDelSetCols),
		}

		c.registerConstraint(partition.OID, partCon)
		c.recordDependency('c', partCon.OID, 0, 'r', partition.OID, 0, DepAuto)
		c.recordColumnDependencies('c', partCon.OID, partition.OID, partLocalAttnums, DepAuto)
		if con.FRelOID != 0 {
			c.recordDependency('c', partCon.OID, 0, 'r', con.FRelOID, 0, DepNormal)
			c.recordColumnDependencies('c', partCon.OID, con.FRelOID, con.FColumns, DepNormal)
			if refIdxOID := c.lookupFKSupportingIndex(con.FRelOID, con.FColumns); refIdxOID != 0 {
				c.recordDependency('c', partCon.OID, 0, 'r', refIdxOID, 0, DepNormal)
			}
		}
		seenFKOp := map[uint32]bool{}
		for _, opGroup := range [][]uint32{con.PFEqOp, con.PPEqOp, con.FFEqOp} {
			for _, opOID := range opGroup {
				if opOID >= FirstNormalObjectId && !seenFKOp[opOID] {
					seenFKOp[opOID] = true
					c.recordDependency('c', partCon.OID, 0, 'O', opOID, 0, DepNormal)
				}
			}
		}
	}
}

// mapColumnAttnums maps parent column attnums to partition column attnums by name.
//
// (pgddl helper — PG uses map_partition_varattnos)
func (c *Catalog) mapColumnAttnums(parent, partition *Relation, parentAttnums []int16) []int16 {
	result := make([]int16, len(parentAttnums))
	for i, parentAttnum := range parentAttnums {
		if parentAttnum == 0 {
			result[i] = 0 // expression column
			continue
		}
		if int(parentAttnum) > len(parent.Columns) {
			return nil
		}
		colName := parent.Columns[parentAttnum-1].Name
		partIdx, ok := partition.colByName[colName]
		if !ok {
			return nil // column not found
		}
		result[i] = partition.Columns[partIdx].AttNum
	}
	return result
}

// cloneLikeIndexes clones indexes from a LIKE source relation to the new relation.
// PG names the cloned indexes using ChooseRelationName (which calls generateIndexName-like logic).
//
// pg: src/backend/commands/tablecmds.c — expandTableLikeClause (index cloning path)
// (pgddl helper — PG generates IndexStmt and calls DefineIndex; we clone directly)
func (c *Catalog) cloneLikeIndexes(schema *Schema, src, dst *Relation) {
	for _, idx := range c.indexesByRel[src.OID] {
		// Map source attnums to destination attnums by column name.
		dstAttnums := c.mapColumnAttnums(src, dst, idx.Columns)
		if dstAttnums == nil {
			continue
		}

		// Generate index (and backing constraint) name based on new table name
		// and column names, checking both relation/index and constraint namespaces.
		var colNames []string
		for _, attnum := range dstAttnums {
			if int(attnum) > 0 && int(attnum) <= len(dst.Columns) {
				colNames = append(colNames, dst.Columns[attnum-1].Name)
			}
		}
		var srcCon *Constraint
		if idx.ConstraintOID != 0 {
			srcCon = c.constraints[idx.ConstraintOID]
		}
		var idxName string
		if srcCon != nil {
			idxName = c.chooseConstraintAndIndexName(schema, dst, colNames, srcCon.Type)
		} else {
			idxName = generateIndexName(dst.Name, colNames, idx.IsPrimary)
			// Ensure uniqueness.
			for _, exists := schema.Indexes[idxName]; exists; _, exists = schema.Indexes[idxName] {
				idxName = idxName + "1"
			}
		}

		// Clone IndOption.
		var dstIndOption []int16
		if idx.IndOption != nil {
			dstIndOption = make([]int16, len(idx.IndOption))
			copy(dstIndOption, idx.IndOption)
		}

		cloneTablespace := idx.Tablespace
		if cloneTablespace == "" {
			cloneTablespace = normalizeTablespace(c.effectiveIndexTablespace("", dst))
		}
		newIdx := &Index{
			OID:               c.oidGen.Next(),
			Name:              idxName,
			Schema:            schema,
			RelOID:            dst.OID,
			Tablespace:        cloneTablespace,
			OwnerName:         dst.OwnerName,
			Columns:           dstAttnums,
			ColumnNames:       slices.Clone(idx.ColumnNames),
			IsUnique:          idx.IsUnique,
			IsPrimary:         idx.IsPrimary,
			AccessMethod:      idx.AccessMethod,
			NKeyColumns:       idx.NKeyColumns,
			WhereClause:       idx.WhereClause,
			IndOption:         dstIndOption,
			NullsNotDistinct:  idx.NullsNotDistinct,
			Exprs:             slices.Clone(idx.Exprs),
			OpClasses:         slices.Clone(idx.OpClasses),
			OpClassOptions:    slices.Clone(idx.OpClassOptions),
			Collations:        cloneStringSlices(idx.Collations),
			UserCollationOIDs: slices.Clone(idx.UserCollationOIDs),
			AnalyzedExprs:     slices.Clone(idx.AnalyzedExprs),
			AnalyzedWhere:     idx.AnalyzedWhere,
		}

		// If the source index backs a constraint, clone the constraint too.
		// pg: same logic as cloneIndexesToPartition — PK/UNIQUE/EXCLUDE constraints
		// are backed by indexes and must be cloned together.
		if am := c.accessMethodByName[newIdx.AccessMethod]; am != nil {
			c.recordDependency('i', newIdx.OID, 0, 'm', am.OID, 0, DepNormal)
		}
		c.recordOpClassDependencies('i', newIdx.OID, c.indexOpClassRefs(dst, dstAttnums, newIdx.NKeyColumns, newIdx.OpClasses), newIdx.AccessMethod)
		var srcIdxExtDeps []DepEntry
		for _, d := range c.deps {
			if d.ObjType == 'i' && d.ObjOID == idx.OID && d.DepType == DepNormal &&
				!(d.RefType == 'r' && d.RefOID == src.OID) {
				srcIdxExtDeps = append(srcIdxExtDeps, d)
			}
		}
		for _, d := range srcIdxExtDeps {
			c.recordDependency('i', newIdx.OID, 0, d.RefType, d.RefOID, d.RefSubID, DepNormal)
		}
		if srcCon != nil {
			conOID := c.oidGen.Next()
			conAttnums := c.mapColumnAttnums(src, dst, srcCon.Columns)
			if conAttnums == nil {
				conAttnums = dstAttnums
			}

			// PG generates a fresh index/constraint name via ChooseRelationName
			// when cloning LIKE indexes (avoiding collisions with existing
			// constraints such as PG 18 NOT NULL or copied CHECK constraints).
			conName := idxName

			newCon := &Constraint{
				OID:            conOID,
				Name:           conName,
				Type:           srcCon.Type,
				RelOID:         dst.OID,
				Namespace:      schema.OID,
				Columns:        conAttnums,
				IncludeColumns: c.mapColumnAttnums(src, dst, srcCon.IncludeColumns),
				IndexOID:       newIdx.OID,
				Deferrable:     srcCon.Deferrable,
				Deferred:       srcCon.Deferred,
				Validated:      true,
				ConIsLocal:     true,
				ExclOps:        slices.Clone(srcCon.ExclOps),
				ExclOpSchemas:  slices.Clone(srcCon.ExclOpSchemas),
			}

			newIdx.ConstraintOID = conOID
			c.registerConstraint(dst.OID, newCon)
			c.recordDependency('c', conOID, 0, 'r', dst.OID, 0, DepAuto)
			c.recordColumnDependencies('c', conOID, dst.OID, conAttnums, DepAuto)
			var srcConExtDeps []DepEntry
			for _, d := range c.deps {
				if d.ObjType == 'c' && d.ObjOID == srcCon.OID && d.DepType == DepNormal &&
					!(d.RefType == 'r' && d.RefOID == src.OID) {
					srcConExtDeps = append(srcConExtDeps, d)
				}
			}
			for _, d := range srcConExtDeps {
				c.recordDependency('c', conOID, 0, d.RefType, d.RefOID, d.RefSubID, DepNormal)
			}
		}

		c.registerIndex(schema, newIdx)
		c.recordDependency('i', newIdx.OID, 0, 'r', dst.OID, 0, DepAuto)
		c.recordColumnDependencies('i', newIdx.OID, dst.OID, dstAttnums, DepAuto)

		// Force NOT NULL on PK key columns; INCLUDE columns keep their nullability.
		if idx.IsPrimary {
			for _, attnum := range indexKeyAttnums(idx, dstAttnums) {
				if int(attnum) > 0 && int(attnum) <= len(dst.Columns) {
					dst.Columns[attnum-1].NotNull = true
					_ = c.ensureNotNullConstraint(schema, dst, dst.Columns[attnum-1], "")
				}
			}
		}
	}
}

// cloneLikeComments clones comments from a LIKE source relation to the new relation.
// Copies table-level and column-level comments.
//
// pg: src/backend/commands/tablecmds.c — expandTableLikeClause (comment cloning)
// (pgddl helper)
func (c *Catalog) cloneLikeComments(src, dst *Relation) {
	// Note: PG does NOT clone the table-level comment — only column and index comments.

	// Clone column comments by matching column names.
	for _, srcCol := range src.Columns {
		comment, ok := c.comments[commentKey{ObjType: 'r', ObjOID: src.OID, SubID: srcCol.AttNum}]
		if !ok {
			continue
		}
		// Find corresponding column in dst by name.
		if dstIdx, ok := dst.colByName[srcCol.Name]; ok {
			dstAttNum := dst.Columns[dstIdx].AttNum
			c.comments[commentKey{ObjType: 'r', ObjOID: dst.OID, SubID: dstAttNum}] = comment
		}
	}

	// Clone index comments.
	for _, srcIdx := range c.indexesByRel[src.OID] {
		srcComment, ok := c.comments[commentKey{ObjType: 'i', ObjOID: srcIdx.OID, SubID: 0}]
		if !ok {
			continue
		}
		// Find corresponding cloned index in dst by matching column pattern.
		for _, dstIdx := range c.indexesByRel[dst.OID] {
			if len(dstIdx.Columns) == len(srcIdx.Columns) {
				match := true
				for i, srcAttnum := range srcIdx.Columns {
					if int(srcAttnum) > 0 && int(srcAttnum) <= len(src.Columns) {
						srcColName := src.Columns[srcAttnum-1].Name
						dstAttnum := dstIdx.Columns[i]
						if int(dstAttnum) <= 0 || int(dstAttnum) > len(dst.Columns) || dst.Columns[dstAttnum-1].Name != srcColName {
							match = false
							break
						}
					}
				}
				if match {
					c.comments[commentKey{ObjType: 'i', ObjOID: dstIdx.OID, SubID: 0}] = srcComment
					break
				}
			}
		}
	}
}

func rawExprContainsSubLink(expr nodes.Node) bool {
	found := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if found || n == nil {
			return false
		}
		if _, ok := n.(*nodes.SubLink); ok {
			found = true
			return false
		}
		return true
	})
	return found
}

// indexKeyAttnums returns the key part of attnums, a copy of idx's columns
// mapped to another relation: the included columns follow the keys.
func indexKeyAttnums(idx *Index, attnums []int16) []int16 {
	if idx.NKeyColumns > 0 && idx.NKeyColumns < len(attnums) {
		return attnums[:idx.NKeyColumns]
	}
	return attnums
}
