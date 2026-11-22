package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Policy represents a row-level security policy.
type Policy struct {
	OID           uint32
	Name          string
	RelOID        uint32
	CmdType       string // "all", "select", "insert", "update", "delete"
	Permissive    bool
	Roles         []string
	UsingExpr     string       // deparsed USING expression
	CheckExpr     string       // deparsed WITH CHECK expression
	UsingAnalyzed AnalyzedExpr // analyzed USING expression (for Tier 2 deparse)
	CheckAnalyzed AnalyzedExpr // analyzed WITH CHECK expression (for Tier 2 deparse)
}

// CreatePolicy creates a new row-level security policy on a table.
//
// pg: src/backend/commands/policy.c — CreatePolicy
func (c *Catalog) CreatePolicy(stmt *nodes.CreatePolicyStmt) error {
	if stmt.Table == nil {
		return errInvalidParameterValue("CREATE POLICY requires a table name")
	}

	// Validate roles before opening the target relation, matching PostgreSQL's
	// policy_role_list_to_array in CreatePolicy.
	roles := c.extractPolicyRoles(stmt.Roles)
	if !c.snapshot {
		for _, r := range roles {
			if r != "" && r != "public" {
				if err := c.checkRole(r); err != nil {
					return err
				}
			}
		}
	}

	_, rel, err := c.findRelation(stmt.Table.Schemaname, stmt.Table.Relname)
	if err != nil {
		return err
	}

	// Policies can only be created on tables (regular or partitioned).
	// pg: src/backend/commands/policy.c — CreatePolicy (relation check)
	if rel.RelKind != 'r' && rel.RelKind != 'p' {
		return errWrongObjectType(rel.Name, "a table")
	}

	// Check for duplicate policy name.
	for _, p := range c.policiesByRel[rel.OID] {
		if p.Name == stmt.PolicyName {
			return &Error{
				Code:    CodeDuplicateObject,
				Message: fmt.Sprintf("policy %q for table %q already exists", stmt.PolicyName, rel.Name),
			}
		}
	}

	cmdType := stmt.CmdName
	if cmdType == "" {
		cmdType = "all"
	}

	// Validate USING/WITH CHECK per command type.
	// pg: src/backend/commands/policy.c — CreatePolicy (lines 596-612)
	if (cmdType == "select" || cmdType == "delete") && stmt.WithCheck != nil {
		return &Error{Code: CodeSyntaxError,
			Message: "WITH CHECK cannot be applied to SELECT or DELETE"}
	}
	if cmdType == "insert" && stmt.Qual != nil {
		return &Error{Code: CodeSyntaxError,
			Message: "only WITH CHECK expression allowed for INSERT"}
	}

	// Analyze and boolean-coerce USING/CHECK expressions before mutating catalog state.
	var usingAnalyzed, checkAnalyzed AnalyzedExpr
	var usingExpr, checkExpr string
	if stmt.Qual != nil {
		var err error
		usingAnalyzed, usingExpr, err = c.analyzePolicyExpr(stmt.Qual, rel)
		if err != nil {
			return err
		}
	}
	if stmt.WithCheck != nil {
		var err error
		checkAnalyzed, checkExpr, err = c.analyzePolicyExpr(stmt.WithCheck, rel)
		if err != nil {
			return err
		}
	}

	if c.snapshot {
		c.recordPolicyRoles(roles)
	}

	policy := &Policy{
		OID:           c.oidGen.Next(),
		Name:          stmt.PolicyName,
		RelOID:        rel.OID,
		CmdType:       cmdType,
		Permissive:    stmt.Permissive,
		Roles:         roles,
		UsingExpr:     usingExpr,
		CheckExpr:     checkExpr,
		UsingAnalyzed: usingAnalyzed,
		CheckAnalyzed: checkAnalyzed,
	}

	c.policies[policy.OID] = policy
	c.policiesByRel[rel.OID] = append(c.policiesByRel[rel.OID], policy)
	c.recordDependency('p', policy.OID, 0, 'r', rel.OID, 0, DepAuto)
	if policy.UsingAnalyzed != nil {
		c.recordDependencyOnSingleRelExpr('p', policy.OID, policy.UsingAnalyzed, rel.OID, DepNormal, DepNormal)
	}
	if policy.CheckAnalyzed != nil {
		c.recordDependencyOnSingleRelExpr('p', policy.OID, policy.CheckAnalyzed, rel.OID, DepNormal, DepNormal)
	}

	return nil
}

// AlterPolicy modifies an existing row-level security policy.
//
// pg: src/backend/commands/policy.c — AlterPolicy
func (c *Catalog) AlterPolicy(stmt *nodes.AlterPolicyStmt) error {
	if stmt.Table == nil {
		return errInvalidParameterValue("ALTER POLICY requires a table name")
	}

	// Validate roles before opening the target relation, matching PostgreSQL's
	// policy_role_list_to_array in AlterPolicy.
	var newRoles []string
	if stmt.Roles != nil {
		newRoles = c.extractPolicyRoles(stmt.Roles)
		if !c.snapshot {
			for _, r := range newRoles {
				if r != "" && r != "public" {
					if err := c.checkRole(r); err != nil {
						return err
					}
				}
			}
		}
	}

	_, rel, err := c.findRelation(stmt.Table.Schemaname, stmt.Table.Relname)
	if err != nil {
		return err
	}
	if rel.RelKind != 'r' && rel.RelKind != 'p' {
		return errWrongObjectType(rel.Name, "a table")
	}

	// Stage and validate USING/WITH CHECK expressions before policy lookup,
	// matching PostgreSQL's AlterPolicy (parse/coerce expressions on rel first).
	var newUsingAnalyzed, newCheckAnalyzed AnalyzedExpr
	var newUsingExpr, newCheckExpr string
	if stmt.Qual != nil {
		var err error
		newUsingAnalyzed, newUsingExpr, err = c.analyzePolicyExpr(stmt.Qual, rel)
		if err != nil {
			return err
		}
	}
	if stmt.WithCheck != nil {
		var err error
		newCheckAnalyzed, newCheckExpr, err = c.analyzePolicyExpr(stmt.WithCheck, rel)
		if err != nil {
			return err
		}
	}

	var policy *Policy
	for _, p := range c.policiesByRel[rel.OID] {
		if p.Name == stmt.PolicyName {
			policy = p
			break
		}
	}
	if policy == nil {
		return &Error{
			Code:    CodeUndefinedObject,
			Message: fmt.Sprintf("policy %q for table %q does not exist", stmt.PolicyName, rel.Name),
		}
	}

	// Validate USING/WITH CHECK per command type.
	// pg: src/backend/commands/policy.c — AlterPolicy (lines 898-915)
	if (policy.CmdType == "select" || policy.CmdType == "delete") && stmt.WithCheck != nil {
		return &Error{Code: CodeSyntaxError,
			Message: "only USING expression allowed for SELECT, DELETE"}
	}
	if policy.CmdType == "insert" && stmt.Qual != nil {
		return &Error{Code: CodeSyntaxError,
			Message: "only WITH CHECK expression allowed for INSERT"}
	}

	if stmt.Roles != nil {
		policy.Roles = newRoles
		if c.snapshot {
			c.recordPolicyRoles(policy.Roles)
		}
	}
	if stmt.Qual != nil {
		policy.UsingAnalyzed = newUsingAnalyzed
		policy.UsingExpr = newUsingExpr
	}
	if stmt.WithCheck != nil {
		policy.CheckAnalyzed = newCheckAnalyzed
		policy.CheckExpr = newCheckExpr
	}
	if stmt.Qual != nil || stmt.WithCheck != nil {
		c.removeDepsOf('p', policy.OID)
		c.recordDependency('p', policy.OID, 0, 'r', rel.OID, 0, DepAuto)
		if policy.UsingAnalyzed != nil {
			c.recordDependencyOnSingleRelExpr('p', policy.OID, policy.UsingAnalyzed, rel.OID, DepNormal, DepNormal)
		}
		if policy.CheckAnalyzed != nil {
			c.recordDependencyOnSingleRelExpr('p', policy.OID, policy.CheckAnalyzed, rel.OID, DepNormal, DepNormal)
		}
	}

	return nil
}

func (c *Catalog) analyzePolicyExpr(expr nodes.Node, rel *Relation) (AnalyzedExpr, string, error) {
	analyzed, err := c.AnalyzeStandaloneExpr(expr, rel)
	if err != nil {
		return nil, "", err
	}
	coerced, err := c.coerceToBoolean(analyzed, "POLICY")
	if err != nil {
		return nil, "", err
	}
	rte := c.buildRelationRTE(rel)
	return coerced, c.DeparseExpr(coerced, []*RangeTableEntry{rte}, false), nil
}

// extractRoleNames converts a list of RoleSpec nodes to role name strings.
// If the roles list is empty, returns ["public"] (PG default).
//
// pg: src/backend/commands/policy.c — policy role handling
// recordPolicyRoles records the cluster-wide roles a policy names, as a dump
// does not create them.
func (c *Catalog) recordPolicyRoles(roles []string) {
	for _, r := range roles {
		if r != "" && r != "public" {
			c.roles[r] = true
		}
	}
}

func (c *Catalog) extractPolicyRoles(roleList *nodes.List) []string {
	if roleList == nil || len(roleList.Items) == 0 {
		return []string{"public"}
	}
	var roles []string
	for _, item := range roleList.Items {
		rs, ok := item.(*nodes.RoleSpec)
		if !ok {
			continue
		}
		if nodes.RoleSpecType(rs.Roletype) == nodes.ROLESPEC_PUBLIC {
			roles = append(roles, "public")
		} else {
			roles = append(roles, c.ownerRoleName(rs))
		}
	}
	if len(roles) == 0 {
		return []string{"public"}
	}
	return roles
}

// removePolicy removes a single policy from all catalog maps and cleans its dependencies.
func (c *Catalog) removePolicy(pol *Policy) {
	if pol == nil {
		return
	}
	delete(c.policies, pol.OID)
	delete(c.extensionMembers, pol.OID)
	c.removeDepsOf('p', pol.OID)
	c.removeDepsOn('p', pol.OID)
	c.removeComments('p', pol.OID)
	plist := c.policiesByRel[pol.RelOID]
	for i, p := range plist {
		if p.OID == pol.OID {
			c.policiesByRel[pol.RelOID] = append(plist[:i], plist[i+1:]...)
			break
		}
	}
}

// removePoliciesForRelation removes all policies belonging to a relation.
func (c *Catalog) removePoliciesForRelation(relOID uint32) {
	for _, p := range c.policiesByRel[relOID] {
		delete(c.policies, p.OID)
		delete(c.extensionMembers, p.OID)
		c.removeDepsOf('p', p.OID)
		c.removeDepsOn('p', p.OID)
		c.removeComments('p', p.OID)
	}
	delete(c.policiesByRel, relOID)
}
