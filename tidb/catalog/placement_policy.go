package catalog

import (
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
	delete(c.placementPolicies, key)
	return nil
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
