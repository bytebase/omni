package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// PlacementPolicy is a first-class catalog object representing a
// TiDB placement policy definition (CREATE PLACEMENT POLICY ...).
// Options are stored in source order; duplicate keywords are
// preserved (matches TiDB's own parser behavior — duplicate detection
// is a semantic concern, not a parse concern).
type PlacementPolicy struct {
	Name    string
	Options []*PlacementPolicyOption
}

// PlacementPolicyOption mirrors the parser AST node but lives in the
// catalog layer so downstream consumers (diff, deparse) don't need
// to import the parser's ast package just to inspect a policy.
type PlacementPolicyOption struct {
	Name     string // uppercased keyword (e.g. "PRIMARY_REGION")
	Value    string // string literal content OR numeric text
	IntValue uint64
	IsInt    bool
}

// GetPlacementPolicy returns the named policy, or nil if not found.
// Policy names are case-insensitive per TiDB's CIStr convention.
func (c *Catalog) GetPlacementPolicy(name string) *PlacementPolicy {
	if c.placementPolicies == nil {
		return nil
	}
	return c.placementPolicies[toLower(name)]
}

// PlacementPolicies returns all defined placement policies, in
// non-deterministic order (map iteration).
func (c *Catalog) PlacementPolicies() []*PlacementPolicy {
	out := make([]*PlacementPolicy, 0, len(c.placementPolicies))
	for _, p := range c.placementPolicies {
		out = append(out, p)
	}
	return out
}

func (c *Catalog) createPlacementPolicy(stmt *nodes.CreatePlacementPolicyStmt) error {
	if c.placementPolicies == nil {
		c.placementPolicies = make(map[string]*PlacementPolicy)
	}
	key := toLower(stmt.Name)
	if _, exists := c.placementPolicies[key]; exists {
		if stmt.OrReplace {
			// OR REPLACE overwrites; no error.
		} else if stmt.IfNotExists {
			return nil
		} else {
			return errDupPlacementPolicy(stmt.Name)
		}
	}
	c.placementPolicies[key] = &PlacementPolicy{
		Name:    stmt.Name,
		Options: convertPolicyOptions(stmt.Options),
	}
	return nil
}

func (c *Catalog) alterPlacementPolicy(stmt *nodes.AlterPlacementPolicyStmt) error {
	key := toLower(stmt.Name)
	existing, ok := c.placementPolicies[key]
	if !ok {
		if stmt.IfExists {
			return nil
		}
		return errUnknownPlacementPolicy(stmt.Name)
	}
	// ALTER replaces the option list wholesale in TiDB — it's not a
	// merge. Mirror that here so the catalog agrees with TiDB state.
	existing.Options = convertPolicyOptions(stmt.Options)
	return nil
}

func (c *Catalog) dropPlacementPolicy(stmt *nodes.DropPlacementPolicyStmt) error {
	key := toLower(stmt.Name)
	if _, ok := c.placementPolicies[key]; !ok {
		if stmt.IfExists {
			return nil
		}
		return errUnknownPlacementPolicy(stmt.Name)
	}
	// TiDB rejects DROP when any database or table still references
	// the policy (error 8240 "Placement policy 'X' is still in use").
	// The catalog mirrors TiDB's execution-time validation since it
	// has both the database and table reference info indexed.
	if ref := c.findPolicyReference(stmt.Name); ref != "" {
		return errPlacementPolicyInUse(stmt.Name, ref)
	}
	delete(c.placementPolicies, key)
	return nil
}

// findPolicyReference returns a human-readable description of the
// first database or table referencing the named policy, or "" if
// there are no references. Case-insensitive comparison matches
// TiDB's CIStr policy-name convention.
func (c *Catalog) findPolicyReference(name string) string {
	target := toLower(name)
	for _, db := range c.databases {
		if toLower(db.PlacementPolicy) == target {
			return fmt.Sprintf("database %q", db.Name)
		}
		for _, t := range db.Tables {
			if toLower(t.PlacementPolicy) == target {
				return fmt.Sprintf("table %q.%q", db.Name, t.Name)
			}
		}
	}
	return ""
}

// validatePolicyRef returns an error if the named policy is not
// registered in the catalog. Empty names (no PLACEMENT POLICY clause)
// are treated as "no reference" and return nil. Used by the table
// and database option wiring to reject references to unknown policies
// (TiDB error 8237 parity).
//
// The name "default" (case-insensitive) is a TiDB sentinel meaning
// "clear placement / inherit" and short-circuits catalog validation —
// see pkg/ddl/placement_policy.go's defaultPlacementPolicyName check
// in upstream. All four grammar arms that produce the policy reference
// (`='default'`, `=default`, `=DEFAULT`, `SET DEFAULT`) collapse to the
// same StrValue, so a single string-compare suffices.
func (c *Catalog) validatePolicyRef(name string) error {
	if name == "" || isDefaultPolicyName(name) {
		return nil
	}
	if c.GetPlacementPolicy(name) == nil {
		return errUnknownPlacementPolicy(name)
	}
	return nil
}

// isDefaultPolicyName reports whether the given policy reference is
// the special-cased "default" sentinel (case-insensitive). Callers
// must treat this as a clear-policy operation, not a named reference.
func isDefaultPolicyName(name string) bool {
	return toLower(name) == "default"
}

// resolvePolicyRef returns the stored form of a policy reference. The
// special-cased "default" sentinel collapses to the empty string so
// downstream code sees "no policy" rather than a pseudo-policy named
// "default". Matches TiDB's in-memory representation, where SET
// PLACEMENT POLICY = default clears the ref on the table/database.
func resolvePolicyRef(name string) string {
	if isDefaultPolicyName(name) {
		return ""
	}
	return name
}

func convertPolicyOptions(src []*nodes.PlacementPolicyOption) []*PlacementPolicyOption {
	out := make([]*PlacementPolicyOption, len(src))
	for i, s := range src {
		out[i] = &PlacementPolicyOption{
			Name:     s.Name,
			Value:    s.Value,
			IntValue: s.IntValue,
			IsInt:    s.IsInt,
		}
	}
	return out
}
