package catalog

import (
	"fmt"
	"sort"
	"strings"
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
// depending on an object newly created in a skipped schema).
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
			filtered, err := filterSkippedSchemaOps(from, to, plan, skipSchemas)
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

// filterSkippedSchemaOps removes operations targeting skipSchemas from plan,
// returning an error if doing so would break cross-schema dependencies of
// retained operations.
func filterSkippedSchemaOps(from, to *Catalog, plan *MigrationPlan, skipSchemas map[string]bool) (*MigrationPlan, error) {
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

	if err := checkSkippedFromDependencies(from, kept, skipSchemas); err != nil {
		return nil, err
	}
	if err := checkSkippedToDependencies(to, kept, skipped); err != nil {
		return nil, err
	}
	return &MigrationPlan{Ops: kept}, nil
}

// checkSkippedFromDependencies rejects plans where a retained drop or column
// type/drop alteration in a non-skipped schema has dependent objects in a
// skipped schema whose prerequisite drop was filtered out.
func checkSkippedFromDependencies(from *Catalog, kept []MigrationOp, skipSchemas map[string]bool) error {
	var fullDropSeeds []depObjectKey
	var colAlterSeeds []depObjectKey
	for _, op := range kept {
		switch {
		case op.Type == OpDropColumn:
			if op.ObjType == 'r' && op.ObjOID != 0 {
				colAlterSeeds = append(colAlterSeeds, depObjectKey{objType: 'r', objOID: op.ObjOID})
			}
		case op.Type == OpAlterColumn && strings.Contains(op.SQL, " TYPE "):
			if rel := from.GetRelation(op.SchemaName, op.ObjectName); rel != nil {
				colAlterSeeds = append(colAlterSeeds, depObjectKey{objType: 'r', objOID: rel.OID})
			}
		case (op.Phase == PhasePre || op.Phase == PhaseLateDrop) && op.ObjType != 0 && op.ObjOID != 0:
			fullDropSeeds = append(fullDropSeeds, depObjectKey{objType: op.ObjType, objOID: op.ObjOID})
		}
	}

	if len(colAlterSeeds) > 0 {
		if schemaName, objName := firstSkippedFromDependent(from, transitiveDependents(from, colAlterSeeds), skipSchemas, true); schemaName != "" {
			return fmt.Errorf("cannot skip schema %q: dependent object %s blocks column alteration in non-skipped schema", schemaName, migrationQualifiedName(schemaName, objName))
		}
	}

	if len(fullDropSeeds) > 0 {
		if schemaName, objName := firstSkippedFromDependent(from, transitiveDependents(from, fullDropSeeds), skipSchemas, false); schemaName != "" {
			return fmt.Errorf("cannot skip schema %q: dependent object %s blocks drop in non-skipped schema", schemaName, migrationQualifiedName(schemaName, objName))
		}
	}
	return nil
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

// checkSkippedToDependencies rejects plans where a retained CREATE/ALTER
// operation depends on an object created only in a skipped schema.
func checkSkippedToDependencies(to *Catalog, kept, skipped []MigrationOp) error {
	var createdSeeds []depObjectKey
	for _, op := range skipped {
		if op.Phase == PhaseMain && op.ObjType != 0 && op.ObjOID != 0 {
			createdSeeds = append(createdSeeds, depObjectKey{objType: op.ObjType, objOID: op.ObjOID})
		}
	}
	if len(createdSeeds) == 0 {
		return nil
	}
	dependents := transitiveDependents(to, createdSeeds)
	for _, op := range kept {
		if (op.Phase != PhaseMain && op.Phase != PhasePost) || op.ObjType == 0 || op.ObjOID == 0 {
			continue
		}
		if dependents[depObjectKey{objType: op.ObjType, objOID: op.ObjOID}] {
			return fmt.Errorf("cannot skip schema: %s %s depends on an object created in a skipped schema", op.Type, migrationQualifiedName(op.SchemaName, op.ObjectName))
		}
	}
	return nil
}

// transitiveDependents returns all catalog objects that directly or
// transitively depend on any key in seeds, excluding the seeds themselves.
func transitiveDependents(c *Catalog, seeds []depObjectKey) map[depObjectKey]bool {
	visited := make(map[depObjectKey]bool, len(seeds))
	queue := make([]depObjectKey, 0, len(seeds))
	for _, s := range seeds {
		if !visited[s] {
			visited[s] = true
			queue = append(queue, s)
		}
	}
	dependents := make(map[depObjectKey]bool)
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

// fromDependentInSkippedSchema resolves a catalog dependency key to its
// (schema, name) when it belongs to skipSchemas. When viewsAndRulesOnly is
// true, only views, materialized views, and rules are reported (matching the
// objects PostgreSQL refuses to leave intact across column type/drop changes).
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
	case 'S':
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
	case 'T':
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
