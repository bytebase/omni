package ddl

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Policy represents a row-level security policy.
type Policy struct {
	OID        uint32
	Name       string
	RelOID     uint32
	CmdType    string // "all", "select", "insert", "update", "delete"
	Permissive bool
	Roles      []string
	UsingExpr      string       // deparsed USING expression
	CheckExpr      string       // deparsed WITH CHECK expression
	UsingAnalyzed  AnalyzedExpr // analyzed USING expression (for Tier 2 deparse)
	CheckAnalyzed  AnalyzedExpr // analyzed WITH CHECK expression (for Tier 2 deparse)
}

// CreatePolicy creates a new row-level security policy on a table.
//
// pg: src/backend/commands/policy.c — CreatePolicy
func (c *Catalog) CreatePolicy(stmt *nodes.CreatePolicyStmt) error {
	if stmt.Table == nil {
		return errInvalidParameterValue("CREATE POLICY requires a table name")
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

	roles := extractRoleNames(stmt.Roles)
	usingExpr := deparseExprNode(stmt.Qual)
	checkExpr := deparseExprNode(stmt.WithCheck)

	policy := &Policy{
		OID:        c.oidGen.Next(),
		Name:       stmt.PolicyName,
		RelOID:     rel.OID,
		CmdType:    cmdType,
		Permissive: stmt.Permissive,
		Roles:      roles,
		UsingExpr:  usingExpr,
		CheckExpr:  checkExpr,
	}

	// Analyze USING/CHECK expressions using Tier 2 pipeline.
	if stmt.Qual != nil {
		if analyzed, err := c.AnalyzeStandaloneExpr(stmt.Qual, rel); err == nil && analyzed != nil {
			policy.UsingAnalyzed = analyzed
		}
	}
	if stmt.WithCheck != nil {
		if analyzed, err := c.AnalyzeStandaloneExpr(stmt.WithCheck, rel); err == nil && analyzed != nil {
			policy.CheckAnalyzed = analyzed
		}
	}

	c.policies[policy.OID] = policy
	c.policiesByRel[rel.OID] = append(c.policiesByRel[rel.OID], policy)
	c.recordDependency('p', policy.OID, 0, 'r', rel.OID, 0, DepAuto)

	return nil
}

// AlterPolicy modifies an existing row-level security policy.
//
// pg: src/backend/commands/policy.c — AlterPolicy
func (c *Catalog) AlterPolicy(stmt *nodes.AlterPolicyStmt) error {
	if stmt.Table == nil {
		return errInvalidParameterValue("ALTER POLICY requires a table name")
	}

	_, rel, err := c.findRelation(stmt.Table.Schemaname, stmt.Table.Relname)
	if err != nil {
		return err
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

	// Update fields if provided.
	if stmt.Roles != nil {
		policy.Roles = extractRoleNames(stmt.Roles)
	}
	if stmt.Qual != nil {
		policy.UsingExpr = deparseExprNode(stmt.Qual)
		if analyzed, err := c.AnalyzeStandaloneExpr(stmt.Qual, rel); err == nil && analyzed != nil {
			policy.UsingAnalyzed = analyzed
		}
	}
	if stmt.WithCheck != nil {
		policy.CheckExpr = deparseExprNode(stmt.WithCheck)
		if analyzed, err := c.AnalyzeStandaloneExpr(stmt.WithCheck, rel); err == nil && analyzed != nil {
			policy.CheckAnalyzed = analyzed
		}
	}

	return nil
}

// extractRoleNames converts a list of RoleSpec nodes to role name strings.
// If the roles list is empty, returns ["public"] (PG default).
//
// pg: src/backend/commands/policy.c — policy role handling
func extractRoleNames(roleList *nodes.List) []string {
	if roleList == nil || len(roleList.Items) == 0 {
		return []string{"public"}
	}
	var roles []string
	for _, item := range roleList.Items {
		rs, ok := item.(*nodes.RoleSpec)
		if !ok {
			continue
		}
		roles = append(roles, roleSpecName(rs))
	}
	if len(roles) == 0 {
		return []string{"public"}
	}
	return roles
}

// removePoliciesForRelation removes all policies belonging to a relation.
func (c *Catalog) removePoliciesForRelation(relOID uint32) {
	for _, p := range c.policiesByRel[relOID] {
		delete(c.policies, p.OID)
		c.removeDepsOf('p', p.OID)
	}
	delete(c.policiesByRel, relOID)
}
