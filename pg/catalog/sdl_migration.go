package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
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
// Returns (nil, nil) when there are no remaining changes. When operations are
// filtered by opts.SkipSchemas, the filtered plan is validated against the
// source catalog (with the skipped schemas intact) so cross-schema dependency
// violations return an error instead of an unexecutable migration plan.
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
			var kept []MigrationOp
			skippedAny := false
			for _, op := range plan.Ops {
				if opTargetsSkippedSchema(op, skipSchemas) {
					skippedAny = true
					continue
				}
				kept = append(kept, op)
			}
			plan = &MigrationPlan{Ops: kept}
			if skippedAny && len(plan.Ops) > 0 {
				if err := validateFilteredPlan(from, diff, plan, skipSchemas); err != nil {
					return nil, err
				}
			}
		}
	}
	if len(plan.Ops) == 0 {
		return nil, nil
	}
	return plan, nil
}

// validateFilteredPlan replays the filtered migration plan against the source
// catalog (where skipped schemas remain untouched) and checks for references to
// skipped enum value additions, so cross-schema dependencies broken by
// filtering out skipped-schema operations fail instead of producing an
// unexecutable migration plan.
func validateFilteredPlan(from *Catalog, diff *SchemaDiff, plan *MigrationPlan, skipSchemas map[string]bool) error {
	planSQL := plan.SQL()

	from.allowImplicitForeignServer = true
	defer func() { from.allowImplicitForeignServer = false }()

	results, err := from.Exec(planSQL, nil)
	if err != nil {
		return fmt.Errorf("filtered migration plan is invalid after skipping schemas: %w", err)
	}
	for _, r := range results {
		if r.Error != nil {
			return fmt.Errorf("filtered migration plan is invalid after skipping schemas: %w", r.Error)
		}
	}
	return checkSkippedEnumValueAdditions(from, diff, planSQL, skipSchemas)
}

type skippedEnumAddition struct {
	schema string
	name   string
	values map[string]bool
}

// checkSkippedEnumValueAdditions rejects retained DDL statements that reference
// enum labels added only in a skipped schema (Catalog.Exec does not validate
// enum literal values inside column defaults or expression casts).
func checkSkippedEnumValueAdditions(c *Catalog, diff *SchemaDiff, planSQL string, skipSchemas map[string]bool) error {
	addedByOID := make(map[uint32]skippedEnumAddition)
	for _, e := range diff.Enums {
		if e.Action != DiffModify || !skipSchemas[e.SchemaName] {
			continue
		}
		existing := make(map[string]bool, len(e.FromValues))
		for _, v := range e.FromValues {
			existing[v] = true
		}
		added := make(map[string]bool)
		for _, v := range e.ToValues {
			if !existing[v] {
				added[v] = true
			}
		}
		if len(added) == 0 {
			continue
		}
		oid, _, err := c.ResolveType(TypeName{Schema: e.SchemaName, Name: e.Name, TypeMod: -1})
		if err != nil {
			continue
		}
		addedByOID[oid] = skippedEnumAddition{schema: e.SchemaName, name: e.Name, values: added}
	}
	if len(addedByOID) == 0 {
		return nil
	}

	list, err := pgparser.Parse(planSQL)
	if err != nil || list == nil {
		return nil
	}

	lookupOID := func(oid uint32) (skippedEnumAddition, bool) {
		if add, ok := addedByOID[oid]; ok {
			return add, true
		}
		if bt := c.typeByOID[oid]; bt != nil && bt.Category == 'A' {
			if add, ok := addedByOID[bt.Elem]; ok {
				return add, true
			}
		}
		return skippedEnumAddition{}, false
	}

	resolveTypeName := func(tn *nodes.TypeName) (skippedEnumAddition, bool) {
		if tn == nil {
			return skippedEnumAddition{}, false
		}
		oid, _, err := c.ResolveType(convertTypeNameToInternal(tn))
		if err != nil {
			return skippedEnumAddition{}, false
		}
		return lookupOID(oid)
	}

	findAddedLiteral := func(expr nodes.Node, add skippedEnumAddition) (string, bool) {
		switch e := expr.(type) {
		case *nodes.A_Const:
			if s, ok := e.Val.(*nodes.String); ok && add.values[s.Str] {
				return s.Str, true
			}
		case *nodes.A_ArrayExpr:
			if e.Elements != nil {
				for _, elem := range e.Elements.Items {
					if ac, ok := elem.(*nodes.A_Const); ok {
						if s, ok := ac.Val.(*nodes.String); ok && add.values[s.Str] {
							return s.Str, true
						}
					}
				}
			}
		}
		return "", false
	}

	var validationErr error
	for _, stmt := range list.Items {
		if validationErr != nil {
			break
		}
		nodes.Inspect(stmt, func(n nodes.Node) bool {
			if validationErr != nil {
				return false
			}
			switch node := n.(type) {
			case *nodes.TypeCast:
				if add, ok := resolveTypeName(node.TypeName); ok {
					if val, hit := findAddedLiteral(node.Arg, add); hit {
						validationErr = fmt.Errorf("filtered migration plan is invalid after skipping schemas: references enum value %q not present on %s", val, migrationQualifiedName(add.schema, add.name))
						return false
					}
				}
			case *nodes.ColumnDef:
				if add, ok := resolveTypeName(node.TypeName); ok && node.Constraints != nil {
					for _, item := range node.Constraints.Items {
						if con, ok := item.(*nodes.Constraint); ok && con.Contype == nodes.CONSTR_DEFAULT {
							if val, hit := findAddedLiteral(con.RawExpr, add); hit {
								validationErr = fmt.Errorf("filtered migration plan is invalid after skipping schemas: references enum value %q not present on %s", val, migrationQualifiedName(add.schema, add.name))
								return false
							}
						}
					}
				}
			case *nodes.CreateDomainStmt:
				if add, ok := resolveTypeName(node.Typname); ok && node.Constraints != nil {
					for _, item := range node.Constraints.Items {
						if con, ok := item.(*nodes.Constraint); ok && con.Contype == nodes.CONSTR_DEFAULT {
							if val, hit := findAddedLiteral(con.RawExpr, add); hit {
								validationErr = fmt.Errorf("filtered migration plan is invalid after skipping schemas: references enum value %q not present on %s", val, migrationQualifiedName(add.schema, add.name))
								return false
							}
						}
					}
				}
			case *nodes.AlterTableStmt:
				if node.Relation == nil || node.Cmds == nil {
					return true
				}
				rel := c.GetRelation(node.Relation.Schemaname, node.Relation.Relname)
				if rel == nil {
					return true
				}
				for _, item := range node.Cmds.Items {
					cmd, ok := item.(*nodes.AlterTableCmd)
					if !ok || nodes.AlterTableType(cmd.Subtype) != nodes.AT_ColumnDefault || cmd.Def == nil {
						continue
					}
					for _, col := range rel.Columns {
						if col.Name != cmd.Name {
							continue
						}
						if add, ok := lookupOID(col.TypeOID); ok {
							if val, hit := findAddedLiteral(cmd.Def, add); hit {
								validationErr = fmt.Errorf("filtered migration plan is invalid after skipping schemas: references enum value %q not present on %s", val, migrationQualifiedName(add.schema, add.name))
								return false
							}
						}
					}
				}
			}
			return true
		})
	}
	return validationErr
}
