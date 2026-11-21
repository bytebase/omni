package catalog

import (
	"fmt"
	"sort"
)

// SDLMigrationOptions configures BuildSDLMigrationPlan.
type SDLMigrationOptions struct {
	// SkipSchemas lists schema names excluded from the generated migration plan
	// (for example, backup schemas).
	SkipSchemas []string
}

// LoadSDLOrSQL parses a PostgreSQL schema text and returns a Catalog reflecting
// it. LoadSDL is the canonical, dependency-aware path. LoadSQL is tried only as
// a fallback for raw dumps containing non-SDL statements (for example SET or
// SELECT) that LoadSDL legitimately rejects. When both fail, the LoadSDL error
// is returned because it is more diagnostic than LoadSQL's order-dependent
// failures (such as a forward foreign-key reference).
func LoadSDLOrSQL(text string) (*Catalog, error) {
	c, sdlErr := LoadSDL(text)
	if sdlErr == nil {
		return c, nil
	}
	if c, err := LoadSQL(text); err == nil {
		return c, nil
	}
	return nil, sdlErr
}

// BuildSDLMigrationPlan loads sourceSDL and targetSDL into catalogs via
// LoadSDLOrSQL, diffs them, generates a MigrationPlan, and filters out any
// operations targeting schemas listed in opts.SkipSchemas.
//
// Returns (nil, nil) when there are no remaining changes. Returns an error when
// filtering opts.SkipSchemas would drop prerequisite operations required by
// retained operations (for example, a view in a skipped schema depending on an
// altered column or dropped table in a non-skipped schema, or a retained object
// depending on an object or column newly created in a skipped schema).
func BuildSDLMigrationPlan(sourceSDL, targetSDL string, opts SDLMigrationOptions) (*MigrationPlan, error) {
	from, err := LoadSDLOrSQL(sourceSDL)
	if err != nil {
		return nil, fmt.Errorf("failed to load source schema: %w", err)
	}
	to, err := LoadSDLOrSQL(targetSDL)
	if err != nil {
		return nil, fmt.Errorf("failed to load target schema: %w", err)
	}
	diff := Diff(from, to)
	if diff.IsEmpty() {
		return nil, nil
	}
	plan := GenerateMigration(from, to, diff)
	if len(opts.SkipSchemas) > 0 {
		skipSchemas := make(map[string]bool, len(opts.SkipSchemas))
		for _, s := range opts.SkipSchemas {
			if s != "" {
				skipSchemas[s] = true
			}
		}
		if len(skipSchemas) > 0 {
			filtered, err := filterSkippedSchemaOps(from, to, diff, plan, skipSchemas)
			if err != nil {
				return nil, err
			}
			plan = filtered
		}
	}
	if len(plan.Ops) == 0 {
		return nil, nil
	}
	return plan, nil
}

type depObjectKey struct {
	objType byte
	objOID  uint32
}

type depSubKey struct {
	objType byte
	objOID  uint32
	subID   int32
}

// normalizeDepObjType maps MigrationOp.ObjType codes onto the object-type bytes
// recorded in Catalog.deps (triggers use 'T' in MigrationOp vs 'g' in deps;
// sequences use 'S' in MigrationOp vs 's' in deps).
func normalizeDepObjType(objType byte) byte {
	switch objType {
	case 'T':
		return 'g'
	case 'S':
		return 's'
	default:
		return objType
	}
}

// isCreateMigrationOp reports whether opType creates a new object whose absence
// would break retained operations that depend on it. Alterations of pre-existing
// objects (such as OpAddColumn or OpAlterTable) are handled separately at
// column granularity rather than treating the whole object as missing.
func isCreateMigrationOp(opType MigrationOpType) bool {
	switch opType {
	case OpCreateSchema,
		OpCreateTable,
		OpCreateView,
		OpCreateType,
		OpCreateFunction,
		OpCreateSequence,
		OpCreateIndex,
		OpCreateTrigger,
		OpCreateEventTrigger,
		OpCreateExtension,
		OpCreateCollation,
		OpCreateAccessMethod,
		OpCreatePolicy,
		OpCreateRule,
		OpAddConstraint:
		return true
	default:
		return false
	}
}

// filterSkippedSchemaOps removes operations targeting skipSchemas from plan,
// returning an error if doing so would break cross-schema dependencies of
// retained operations.
func filterSkippedSchemaOps(from, to *Catalog, diff *SchemaDiff, plan *MigrationPlan, skipSchemas map[string]bool) (*MigrationPlan, error) {
	var kept, skipped []MigrationOp
	for _, op := range plan.Ops {
		if opTargetsSkippedSchema(op, skipSchemas) {
			skipped = append(skipped, op)
		} else {
			kept = append(kept, op)
		}
	}
	if len(kept) == 0 || len(skipped) == 0 {
		return &MigrationPlan{Ops: kept}, nil
	}

	if err := checkSkippedFromDependencies(from, to, diff, kept, skipSchemas); err != nil {
		return nil, err
	}
	if err := checkSkippedToDependencies(to, diff, kept, skipped, skipSchemas); err != nil {
		return nil, err
	}
	return &MigrationPlan{Ops: kept}, nil
}

// checkSkippedFromDependencies rejects plans where a retained drop or column
// type/drop alteration in a non-skipped schema has dependent objects in a
// skipped schema whose prerequisite drop was filtered out.
func checkSkippedFromDependencies(from, to *Catalog, diff *SchemaDiff, kept []MigrationOp, skipSchemas map[string]bool) error {
	var fullDropSeeds []depObjectKey
	for _, op := range kept {
		if op.Type == OpDropColumn {
			// Checked at column granularity below via diff.
			continue
		}
		if (op.Phase == PhasePre || op.Phase == PhaseLateDrop) && op.ObjType != 0 && op.ObjOID != 0 {
			fullDropSeeds = append(fullDropSeeds, depObjectKey{
				objType: normalizeDepObjType(op.ObjType),
				objOID:  op.ObjOID,
			})
		}
	}

	var colAlterSeeds []depSubKey
	var colDropSeeds []depSubKey
	collectRelColSeeds := func(rels []RelationDiffEntry) {
		for _, rel := range rels {
			if rel.Action != DiffModify || skipSchemas[rel.SchemaName] || rel.From == nil {
				continue
			}
			for _, col := range rel.Columns {
				switch col.Action {
				case DiffDrop:
					if col.From != nil {
						colDropSeeds = append(colDropSeeds, depSubKey{
							objType: 'r',
							objOID:  rel.From.OID,
							subID:   int32(col.From.AttNum),
						})
					}
				case DiffModify:
					if col.From == nil || col.To == nil {
						continue
					}
					oldType := from.FormatType(col.From.TypeOID, col.From.TypeMod)
					newType := to.FormatType(col.To.TypeOID, col.To.TypeMod)
					if oldType != newType || columnCollationKey(from, col.From) != columnCollationKey(to, col.To) {
						colAlterSeeds = append(colAlterSeeds, depSubKey{
							objType: 'r',
							objOID:  rel.From.OID,
							subID:   int32(col.From.AttNum),
						})
					}
				}
			}
		}
	}
	collectRelColSeeds(diff.Relations)
	collectRelColSeeds(diff.ForeignTables)

	if len(colAlterSeeds) > 0 {
		deps := transitiveDependents(from, nil, colAlterSeeds)
		if schemaName, objName := firstSkippedFromDependent(from, deps, skipSchemas, true); schemaName != "" {
			return fmt.Errorf("cannot skip schema %q: dependent object %s blocks column alteration in non-skipped schema", schemaName, migrationQualifiedName(schemaName, objName))
		}
	}

	if len(colDropSeeds) > 0 {
		deps := transitiveDependents(from, nil, colDropSeeds)
		if schemaName, objName := firstSkippedFromDependent(from, deps, skipSchemas, false); schemaName != "" {
			return fmt.Errorf("cannot skip schema %q: dependent object %s blocks column drop in non-skipped schema", schemaName, migrationQualifiedName(schemaName, objName))
		}
	}

	if len(fullDropSeeds) > 0 {
		deps := transitiveDependents(from, fullDropSeeds, nil)
		if schemaName, objName := firstSkippedFromDependent(from, deps, skipSchemas, false); schemaName != "" {
			return fmt.Errorf("cannot skip schema %q: dependent object %s blocks drop in non-skipped schema", schemaName, migrationQualifiedName(schemaName, objName))
		}
	}
	return nil
}

// checkSkippedToDependencies rejects plans where a retained CREATE/ALTER
// operation depends on an object or column newly created in a skipped schema.
func checkSkippedToDependencies(to *Catalog, diff *SchemaDiff, kept, skipped []MigrationOp, skipSchemas map[string]bool) error {
	var createdWholeSeeds []depObjectKey
	for _, op := range skipped {
		if !isCreateMigrationOp(op.Type) || op.ObjType == 0 || op.ObjOID == 0 {
			continue
		}
		key := depObjectKey{
			objType: normalizeDepObjType(op.ObjType),
			objOID:  op.ObjOID,
		}
		createdWholeSeeds = append(createdWholeSeeds, key)
		// Also seed the auto-created array type so column types like skipped.mytype[] match.
		if key.objType == 't' {
			if bt := to.typeByOID[key.objOID]; bt != nil && bt.Array != 0 {
				createdWholeSeeds = append(createdWholeSeeds, depObjectKey{objType: 't', objOID: bt.Array})
			}
		}
	}

	var createdColSeeds []depSubKey
	collectSkippedAddedCols := func(rels []RelationDiffEntry) {
		for _, rel := range rels {
			if rel.Action != DiffModify || !skipSchemas[rel.SchemaName] || rel.To == nil {
				continue
			}
			for _, col := range rel.Columns {
				if col.Action == DiffAdd && col.To != nil {
					createdColSeeds = append(createdColSeeds, depSubKey{
						objType: 'r',
						objOID:  rel.To.OID,
						subID:   int32(col.To.AttNum),
					})
				}
			}
		}
	}
	collectSkippedAddedCols(diff.Relations)
	collectSkippedAddedCols(diff.ForeignTables)

	if len(createdWholeSeeds) == 0 && len(createdColSeeds) == 0 {
		return nil
	}

	dependents := transitiveDependents(to, createdWholeSeeds, createdColSeeds)
	// Lift constraint and index dependents onto their owning relation so
	// inline table constraints/indexes created within OpCreateTable match.
	lifted := make([]depObjectKey, 0)
	for key := range dependents {
		switch key.objType {
		case 'c':
			if con := to.constraints[key.objOID]; con != nil && con.RelOID != 0 {
				lifted = append(lifted, depObjectKey{objType: 'r', objOID: con.RelOID})
			}
		case 'i':
			if idx := to.indexes[key.objOID]; idx != nil && idx.RelOID != 0 {
				lifted = append(lifted, depObjectKey{objType: 'r', objOID: idx.RelOID})
			}
		}
	}
	for _, key := range lifted {
		dependents[key] = true
	}

	for _, op := range kept {
		if (op.Phase != PhaseMain && op.Phase != PhasePost) || op.ObjType == 0 || op.ObjOID == 0 {
			continue
		}
		key := depObjectKey{
			objType: normalizeDepObjType(op.ObjType),
			objOID:  op.ObjOID,
		}
		if dependents[key] {
			return fmt.Errorf("cannot skip schema: %s %s depends on an object created in a skipped schema", op.Type, migrationQualifiedName(op.SchemaName, op.ObjectName))
		}
	}
	return nil
}

// transitiveDependents returns all catalog objects that directly or
// transitively depend on any whole-object key in wholeSeeds or any specific
// sub-object (e.g. table column AttNum) in subSeeds.
func transitiveDependents(c *Catalog, wholeSeeds []depObjectKey, subSeeds []depSubKey) map[depObjectKey]bool {
	visited := make(map[depObjectKey]bool, len(wholeSeeds))
	queue := make([]depObjectKey, 0, len(wholeSeeds)+len(subSeeds))
	dependents := make(map[depObjectKey]bool)

	for _, s := range wholeSeeds {
		if !visited[s] {
			visited[s] = true
			queue = append(queue, s)
		}
	}

	// First hop for sub-object seeds: match only dependencies pointing at the
	// exact RefSubID (for example, a specific column AttNum).
	if len(subSeeds) > 0 {
		subSet := make(map[depSubKey]bool, len(subSeeds))
		for _, sub := range subSeeds {
			subSet[sub] = true
		}
		for _, d := range c.deps {
			if d.RefSubID <= 0 {
				continue
			}
			if !subSet[depSubKey{objType: d.RefType, objOID: d.RefOID, subID: d.RefSubID}] {
				continue
			}
			dep := depObjectKey{objType: d.ObjType, objOID: d.ObjOID}
			if !visited[dep] {
				visited[dep] = true
				dependents[dep] = true
				queue = append(queue, dep)
			}
		}
	}

	for len(queue) > 0 {
		ref := queue[0]
		queue = queue[1:]
		for _, d := range c.deps {
			if d.RefType != ref.objType || d.RefOID != ref.objOID {
				continue
			}
			dep := depObjectKey{objType: d.ObjType, objOID: d.ObjOID}
			if !visited[dep] {
				visited[dep] = true
				dependents[dep] = true
				queue = append(queue, dep)
			}
		}
	}
	return dependents
}

func firstSkippedFromDependent(c *Catalog, dependents map[depObjectKey]bool, skipSchemas map[string]bool, viewsAndRulesOnly bool) (string, string) {
	type match struct{ schema, name string }
	var matches []match
	for key := range dependents {
		if s, n := fromDependentInSkippedSchema(c, key, skipSchemas, viewsAndRulesOnly); s != "" {
			matches = append(matches, match{schema: s, name: n})
		}
	}
	if len(matches) == 0 {
		return "", ""
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].schema != matches[j].schema {
			return matches[i].schema < matches[j].schema
		}
		return matches[i].name < matches[j].name
	})
	return matches[0].schema, matches[0].name
}

// fromDependentInSkippedSchema resolves a catalog dependency key to its
// (schema, name) when it belongs to skipSchemas. When viewsAndRulesOnly is
// true, only views, materialized views, and rules are reported (matching the
// objects PostgreSQL refuses to leave intact across column type changes).
func fromDependentInSkippedSchema(c *Catalog, key depObjectKey, skipSchemas map[string]bool, viewsAndRulesOnly bool) (string, string) {
	switch key.objType {
	case 'r':
		rel := c.GetRelationByOID(key.objOID)
		if rel == nil || rel.Schema == nil || !skipSchemas[rel.Schema.Name] {
			return "", ""
		}
		if viewsAndRulesOnly && rel.RelKind != 'v' && rel.RelKind != 'm' {
			return "", ""
		}
		return rel.Schema.Name, rel.Name
	case 'R':
		rule := c.rules[key.objOID]
		if rule == nil {
			return "", ""
		}
		owner := c.GetRelationByOID(rule.OwnerRelOID)
		if owner == nil || owner.Schema == nil || !skipSchemas[owner.Schema.Name] {
			return "", ""
		}
		return owner.Schema.Name, rule.Name
	case 'c':
		if viewsAndRulesOnly {
			return "", ""
		}
		con := c.constraints[key.objOID]
		if con == nil {
			return "", ""
		}
		rel := c.GetRelationByOID(con.RelOID)
		if rel == nil || rel.Schema == nil || !skipSchemas[rel.Schema.Name] {
			return "", ""
		}
		return rel.Schema.Name, con.Name
	case 'f':
		if viewsAndRulesOnly {
			return "", ""
		}
		up := c.userProcs[key.objOID]
		if up == nil || up.Schema == nil || !skipSchemas[up.Schema.Name] {
			return "", ""
		}
		return up.Schema.Name, up.Name
	case 't':
		if viewsAndRulesOnly {
			return "", ""
		}
		bt := c.typeByOID[key.objOID]
		if bt == nil {
			return "", ""
		}
		if s := c.schemas[bt.Namespace]; s != nil && skipSchemas[s.Name] {
			return s.Name, bt.TypeName
		}
	case 's':
		if viewsAndRulesOnly {
			return "", ""
		}
		seq := c.sequenceByOID[key.objOID]
		if seq == nil || seq.Schema == nil || !skipSchemas[seq.Schema.Name] {
			return "", ""
		}
		return seq.Schema.Name, seq.Name
	case 'i':
		if viewsAndRulesOnly {
			return "", ""
		}
		idx := c.indexes[key.objOID]
		if idx == nil {
			return "", ""
		}
		rel := c.GetRelationByOID(idx.RelOID)
		if rel == nil || rel.Schema == nil || !skipSchemas[rel.Schema.Name] {
			return "", ""
		}
		return rel.Schema.Name, idx.Name
	case 'g':
		if viewsAndRulesOnly {
			return "", ""
		}
		trig := c.triggers[key.objOID]
		if trig == nil {
			return "", ""
		}
		rel := c.GetRelationByOID(trig.RelOID)
		if rel == nil || rel.Schema == nil || !skipSchemas[rel.Schema.Name] {
			return "", ""
		}
		return rel.Schema.Name, trig.Name
	}
	return "", ""
}
