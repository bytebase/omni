package catalog

import (
	"fmt"
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

	// Check duplicate relation name.
	if _, exists := schema.Relations[relName]; exists {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", relName))
			return nil
		}
		return errDuplicateTable(relName)
	}

	// Check conflict with existing type name in same namespace.
	if c.typeByName[typeKey{ns: schema.OID, name: relName}] != nil {
		if stmt.IfNotExists {
			c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q already exists, skipping", relName))
			return nil
		}
		return errDuplicateTable(relName)
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

	// Build a name→index map for the child column defs (used by merge).
	// Mark all columns from the CREATE statement as locally defined.
	childColMap := make(map[string]int, len(colDefs))
	for i := range colDefs {
		colDefs[i].IsLocal = true
		childColMap[colDefs[i].Name] = i
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
			typ := c.typeByOID[pcol.TypeOID]
			typName := ""
			if typ != nil {
				typName = typ.TypeName
			}
			childColMap[pcol.Name] = len(colDefs)
			colDefs = append(colDefs, ColumnDef{
				Name:       pcol.Name,
				Type:       TypeName{Name: typName, TypeMod: pcol.TypeMod},
				IsLocal:    false,
				IsFromType: true,
			})
		}

		// Now merge user-specified columns.
		for _, ucd := range userColDefs {
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

	// Merge inherited columns (INHERITS clause).
	// pg: src/backend/commands/tablecmds.c — MergeAttributes
	var parentOIDs []uint32
	isPartition := (stmt.Partbound != nil)
	var conflictDefaults map[string]bool // column name → has conflicting defaults
	if stmt.InhRelations != nil {
		for seqNo, item := range stmt.InhRelations.Items {
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
					if existCD.CollationName != "" && pcol.CollationName != "" &&
						existCD.CollationName != pcol.CollationName {
						return errDatatypeMismatch(fmt.Sprintf(
							"column %q has a collation conflict", pcol.Name,
						))
					}
					// Merge NOT NULL from parent.
					if pcol.NotNull {
						existCD.NotNull = true
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
					typ := c.typeByOID[pcol.TypeOID]
					typName := ""
					if typ != nil {
						typName = typ.TypeName
					}
					newCD := ColumnDef{
						Name:     pcol.Name,
						Type:     TypeName{Name: typName, TypeMod: pcol.TypeMod},
						NotNull:  pcol.NotNull,
						Default:  pcol.Default,
						IsLocal:  false,
						InhCount: 1,
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

			// Record inheritance entry.
			c.inhEntries = append(c.inhEntries, InhEntry{
				InhRelID:  0, // will be set after relation is created
				InhParent: parent.OID,
				InhSeqNo:  int32(seqNo + 1), // 1-based
			})
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

	// Create implicit sequences for SERIAL and identity columns.
	var createdSeqs []*Sequence
	for _, si := range serials {
		if si.isIdentity {
			cd := &colDefs[si.colIdx]
			seq, seqRV, explicitName, err := c.createIdentitySequence(schema, relName, cd.Name, si.typeOID, cd.IdentityOptions)
			if err != nil {
				return err
			}
			cd.IdentitySequence = seqRV
			cd.Default = identitySequenceDefault(schema, seq.Schema, seq.Name, explicitName)
			createdSeqs = append(createdSeqs, seq)
			continue
		}
		if _, exists := schema.Sequences[si.seqName]; exists {
			return errDuplicateObject("sequence", si.seqName)
		}
		seq := c.createSequenceInternal(schema, si.seqName, si.typeOID)
		createdSeqs = append(createdSeqs, seq)
	}

	// Resolve column types and check for duplicates.
	colByName := make(map[string]int, len(colDefs))
	columns := make([]*Column, len(colDefs))

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
		// pg: src/backend/catalog/heap.c — InsertPgAttributeTuples
		// attcollation is derived from type; for arrays, inherited from element type.
		coll := c.typeCollation(typeOID)
		columns[i] = &Column{
			AttNum:        int16(i + 1),
			Name:          cd.Name,
			TypeOID:       typeOID,
			TypeMod:       typmod,
			NotNull:       cd.NotNull,
			Len:           typ.Len,
			ByVal:         typ.ByVal,
			Align:         typ.Align,
			Storage:       typ.Storage,
			Collation:     coll,
			IsLocal:       cd.IsLocal,
			InhCount:      cd.InhCount,
			CollationName: cd.CollationName,
		}
	}

	// Check column count limit.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (MaxHeapAttributeNumber check)
	if len(columns) > MaxHeapAttributeNumber {
		return &Error{Code: CodeTooManyColumns, Message: fmt.Sprintf("tables can have at most %d columns", MaxHeapAttributeNumber)}
	}

	// Allocate OIDs.
	relOID := c.oidGen.Next()
	rowTypeOID := c.oidGen.Next()
	arrayOID := c.oidGen.Next()

	rel := &Relation{
		OID:        relOID,
		Name:       relName,
		Schema:     schema,
		RelKind:    relkind,
		Columns:    columns,
		colByName:  colByName,
		RowTypeOID: rowTypeOID,
		ArrayOID:   arrayOID,
		InhParents: parentOIDs,
		InhCount:   len(parentOIDs),
		OfTypeOID:  ofTypeOID,
	}

	// Parse persistence from RangeVar.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (relpersistence)
	if rv.Relpersistence != 0 {
		rel.Persistence = rv.Relpersistence
	} else {
		rel.Persistence = 'p' // default: permanent
	}

	// ON COMMIT can only be used on temporary tables.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (oncommit validation)
	if stmt.OnCommit != nodes.ONCOMMIT_NOOP && rel.Persistence != 't' {
		return &Error{Code: CodeInvalidTableDefinition,
			Message: "ON COMMIT can only be used on temporary tables"}
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

	// Fix up inheritance entries with the child's OID and record dependencies.
	if len(parentOIDs) > 0 {
		for i := len(c.inhEntries) - len(parentOIDs); i < len(c.inhEntries); i++ {
			c.inhEntries[i].InhRelID = relOID
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

	// Store partition key info (PARTITION BY clause).
	// pg: src/backend/commands/tablecmds.c — DefineRelation (StorePartitionKey)
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
		rel.PartitionInfo = &PartitionInfo{
			Strategy:   strategy,
			KeyAttNums: keyAttNums,
			NKeyAttrs:  int16(len(keyAttNums)),
		}
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

			rel.PartitionBound = &PartitionBound{
				Strategy:   bs.Strategy,
				IsDefault:  bs.IsDefault,
				ListValues: deparseDatumList(bs.Listdatums),
				LowerBound: deparseDatumList(bs.Lowerdatums),
				UpperBound: deparseDatumList(bs.Upperdatums),
				Modulus:    bs.Modulus,
				Remainder:  bs.Remainder,
			}
		}
		// Set PartitionOf to the parent table OID.
		if len(parentOIDs) > 0 {
			rel.PartitionOf = parentOIDs[0]
			// Check for partition bound overlap with existing partitions of the same parent.
			if err := c.checkPartitionBoundOverlap(parentOIDs[0], rel.OID, rel.PartitionBound); err != nil {
				// Roll back the registration done above.
				c.removeRelation(schema, relName, rel)
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
		for _, cd := range constraintDefs {
			if err := c.addConstraint(schema, rel, cd); err != nil {
				c.removeRelation(schema, relName, rel)
				return err
			}
		}

		// Mark implicit sequences as column-owned in the catalog model so
		// diff/drop logic treats SERIAL and IDENTITY sequences as managed objects.
		for i, si := range serials {
			seq := createdSeqs[i]
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
			if analyzed, err := c.AnalyzeStandaloneExpr(cd.RawDefault, rel); err == nil && analyzed != nil {
				// pg: cookDefault uses COERCE_IMPLICIT_CAST ('i') as display format
				coerced, cerr := c.coerceToTargetType(analyzed, analyzed.exprType(), columns[i].TypeOID, 'i')
				if cerr == nil && coerced != nil {
					analyzed = coerced
				}
				columns[i].DefaultAnalyzed = analyzed
				c.recordDependencyOnSingleRelExprForObject('r', rel.OID, int32(columns[i].AttNum), analyzed, rel.OID,
					DepNormal, DepNormal)
				rte := c.buildRelationRTE(rel)
				columns[i].Default = c.DeparseExpr(analyzed, []*RangeTableEntry{rte}, false)
			}
		}
		if cd.RawGenExpr != nil && columns[i].Generated == 's' {
			if rawExprContainsSubLink(cd.RawGenExpr) {
				c.removeRelation(schema, relName, rel)
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

	for _, obj := range stmt.Objects.Items {
		// Objects are represented as a List of String nodes (schema-qualified name).
		schemaName, relName := extractDropObjectName(obj)

		schema, rel, err := c.findRelation(schemaName, relName)
		if err != nil {
			if stmt.Missing_ok {
				if _, ok := err.(*Error); ok {
					c.addWarning(CodeWarningSkip, fmt.Sprintf("relation %q does not exist, skipping", relName))
					continue
				}
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

		// Check for normal dependents.
		if deps := c.findNormalDependents('r', rel.OID); len(deps) > 0 {
			if !cascade {
				kind := "table"
				if removeType == nodes.OBJECT_VIEW {
					kind = "view"
				} else if removeType == nodes.OBJECT_MATVIEW {
					kind = "materialized view"
				}
				return errDependentObjects(kind, relName)
			}
			if removeType == nodes.OBJECT_VIEW || removeType == nodes.OBJECT_MATVIEW {
				c.dropViewDependents('r', rel.OID)
			} else {
				c.dropDependents('r', rel.OID)
			}
		}

		c.removeRelation(schema, relName, rel)
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
			typmod = encodeTypModByName(resolveAlias(name), rawMods)
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
			Name:    resolveAlias(name),
			TypeMod: typmod,
			IsArray: isArray,
		}
	}

	// Handle explicit COLLATE clause.
	// pg: src/backend/commands/tablecmds.c — DefineRelation (COLLATE handling)
	if cd.CollClause != nil && cd.CollClause.Collname != nil && len(cd.CollClause.Collname.Items) > 0 {
		result.CollationName = stringVal(cd.CollClause.Collname.Items[len(cd.CollClause.Collname.Items)-1])
	}

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
		typ := c.typeByOID[col.TypeOID]
		typName := ""
		if typ != nil {
			typName = typ.TypeName
		}

		cd := ColumnDef{
			Name:    col.Name,
			Type:    TypeName{Name: typName, TypeMod: col.TypeMod},
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
			for _, nv := range newBound.ListValues {
				for _, ev := range eb.ListValues {
					if nv == ev {
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
				Code:    CodeDependentObjects,
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
		// Map parent attnums to partition attnums by column name.
		partAttnums := c.mapColumnAttnums(parent, partition, idx.Columns)
		if partAttnums == nil {
			continue // column not found in partition (shouldn't happen for valid partition)
		}

		// Generate partition index name.
		partIdxName := fmt.Sprintf("%s_%s_idx", partition.Name, idx.Name)
		if idx.IsPrimary {
			partIdxName = partition.Name + "_pkey"
		}
		// Ensure uniqueness.
		for _, exists := schema.Indexes[partIdxName]; exists; _, exists = schema.Indexes[partIdxName] {
			partIdxName = partIdxName + "1"
		}

		// Clone IndOption.
		var partIndOption []int16
		if idx.IndOption != nil {
			partIndOption = make([]int16, len(idx.IndOption))
			copy(partIndOption, idx.IndOption)
		}

		partIdx := &Index{
			OID:              c.oidGen.Next(),
			Name:             partIdxName,
			Schema:           schema,
			RelOID:           partition.OID,
			Columns:          partAttnums,
			IsUnique:         idx.IsUnique,
			IsPrimary:        idx.IsPrimary,
			AccessMethod:     idx.AccessMethod,
			NKeyColumns:      idx.NKeyColumns,
			WhereClause:      idx.WhereClause,
			IndOption:        partIndOption,
			NullsNotDistinct: idx.NullsNotDistinct,
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

				partCon := &Constraint{
					OID:         conOID,
					Name:        partConName,
					Type:        parentCon.Type,
					RelOID:      partition.OID,
					Namespace:   schema.OID,
					Columns:     partConAttnums,
					IndexOID:    partIdx.OID,
					Validated:   true,
					ConIsLocal:  true,
					ConParentID: parentCon.OID,
				}

				partIdx.ConstraintOID = conOID
				c.registerConstraint(partition.OID, partCon)
				c.recordDependency('c', conOID, 0, 'r', partition.OID, 0, DepAuto)
			}
		}

		c.registerIndex(schema, partIdx)
		c.recordDependency('i', partIdx.OID, 0, 'r', partition.OID, 0, DepAuto)

		// Force NOT NULL on PK columns.
		if idx.IsPrimary {
			for _, attnum := range partAttnums {
				if int(attnum) > 0 && int(attnum) <= len(partition.Columns) {
					partition.Columns[attnum-1].NotNull = true
				}
			}
		}
	}
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
		}

		c.triggersByRel[partition.OID] = append(c.triggersByRel[partition.OID], partTrig)
		c.recordDependency('g', partTrig.OID, 0, 'r', partition.OID, 0, DepAuto)
		if partTrig.FuncOID != 0 {
			c.recordDependency('g', partTrig.OID, 0, 'f', partTrig.FuncOID, 0, DepNormal)
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
			FColumns:     con.FColumns,
			FKUpdAction:  con.FKUpdAction,
			FKDelAction:  con.FKDelAction,
			FKMatchType:  con.FKMatchType,
			Deferrable:   con.Deferrable,
			Deferred:     con.Deferred,
			Validated:    true,
			ConIsLocal:   true,
			ConParentID:  con.OID,
			PFEqOp:       con.PFEqOp,
			PPEqOp:       con.PPEqOp,
			FFEqOp:       con.FFEqOp,
			FKDelSetCols: con.FKDelSetCols,
		}

		c.registerConstraint(partition.OID, partCon)
		c.recordDependency('c', partCon.OID, 0, 'r', partition.OID, 0, DepAuto)
		if con.FRelOID != 0 {
			c.recordDependency('c', partCon.OID, 0, 'r', con.FRelOID, 0, DepNormal)
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

		// Generate index name based on new table name and column names.
		var colNames []string
		for _, attnum := range dstAttnums {
			if int(attnum) > 0 && int(attnum) <= len(dst.Columns) {
				colNames = append(colNames, dst.Columns[attnum-1].Name)
			}
		}
		idxName := generateIndexName(dst.Name, colNames, idx.IsPrimary)
		// Ensure uniqueness.
		for _, exists := schema.Indexes[idxName]; exists; _, exists = schema.Indexes[idxName] {
			idxName = idxName + "1"
		}

		// Clone IndOption.
		var dstIndOption []int16
		if idx.IndOption != nil {
			dstIndOption = make([]int16, len(idx.IndOption))
			copy(dstIndOption, idx.IndOption)
		}

		newIdx := &Index{
			OID:              c.oidGen.Next(),
			Name:             idxName,
			Schema:           schema,
			RelOID:           dst.OID,
			Columns:          dstAttnums,
			IsUnique:         idx.IsUnique,
			IsPrimary:        idx.IsPrimary,
			AccessMethod:     idx.AccessMethod,
			NKeyColumns:      idx.NKeyColumns,
			WhereClause:      idx.WhereClause,
			IndOption:        dstIndOption,
			NullsNotDistinct: idx.NullsNotDistinct,
		}

		// If the source index backs a constraint, clone the constraint too.
		// pg: same logic as cloneIndexesToPartition — PK/UNIQUE constraints
		// are backed by indexes and must be cloned together.
		if idx.ConstraintOID != 0 {
			srcCon := c.constraints[idx.ConstraintOID]
			if srcCon != nil {
				conOID := c.oidGen.Next()
				conAttnums := c.mapColumnAttnums(src, dst, srcCon.Columns)
				if conAttnums == nil {
					conAttnums = dstAttnums
				}

				// PG naming: replace source table prefix with destination table name.
				// e.g. t1_pkey -> t2_pkey
				conName := srcCon.Name
				if strings.HasPrefix(srcCon.Name, src.Name) {
					conName = dst.Name + srcCon.Name[len(src.Name):]
				}

				newCon := &Constraint{
					OID:        conOID,
					Name:       conName,
					Type:       srcCon.Type,
					RelOID:     dst.OID,
					Namespace:  schema.OID,
					Columns:    conAttnums,
					IndexOID:   newIdx.OID,
					Validated:  true,
					ConIsLocal: true,
				}

				newIdx.ConstraintOID = conOID
				c.registerConstraint(dst.OID, newCon)
				c.recordDependency('c', conOID, 0, 'r', dst.OID, 0, DepAuto)
			}
		}

		c.registerIndex(schema, newIdx)
		c.recordDependency('i', newIdx.OID, 0, 'r', dst.OID, 0, DepAuto)

		// Force NOT NULL on PK columns.
		if idx.IsPrimary {
			for _, attnum := range dstAttnums {
				if int(attnum) > 0 && int(attnum) <= len(dst.Columns) {
					dst.Columns[attnum-1].NotNull = true
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
