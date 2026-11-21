package catalog

import "fmt"

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
// Returns (nil, nil) when there are no remaining changes.
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
			plan = plan.Filter(func(op MigrationOp) bool {
				return !opTargetsSkippedSchema(op, skipSchemas)
			})
		}
	}
	if len(plan.Ops) == 0 {
		return nil, nil
	}
	return plan, nil
}
