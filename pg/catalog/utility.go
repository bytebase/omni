package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/pg/ast"
)

// ProcessUtility dispatches a utility (DDL) statement to the appropriate handler.
//
// pg: src/backend/tcop/utility.c — standard_ProcessUtility
func (c *Catalog) ProcessUtility(stmt nodes.Node) error {
	switch s := stmt.(type) {
	case *nodes.CreateSchemaStmt:
		return c.CreateSchemaCommand(s)
	case *nodes.CreateStmt:
		return c.DefineRelation(s, 'r')
	case *nodes.ViewStmt:
		return c.DefineView(s)
	case *nodes.IndexStmt:
		return c.DefineIndex(s)
	case *nodes.CreateSeqStmt:
		return c.DefineSequence(s)
	case *nodes.AlterSeqStmt:
		return c.AlterSequenceStmt(s)
	case *nodes.CreateEnumStmt:
		return c.DefineEnum(s)
	case *nodes.AlterEnumStmt:
		return c.AlterEnumStmt(s)
	case *nodes.CreateDomainStmt:
		return c.DefineDomain(s)
	case *nodes.AlterDomainStmt:
		return c.AlterDomainStmt(s)
	case *nodes.CreateFunctionStmt:
		return c.CreateFunctionStmt(s)
	case *nodes.CreateTrigStmt:
		return c.CreateTriggerStmt(s)
	case *nodes.CommentStmt:
		return c.CommentObject(s)
	case *nodes.CompositeTypeStmt:
		return c.DefineCompositeType(s)
	case *nodes.CreateRangeStmt:
		return c.DefineRange(s)
	case *nodes.CreateTableAsStmt:
		objType := nodes.ObjectType(s.Objtype)
		if objType == nodes.OBJECT_MATVIEW || objType == nodes.OBJECT_TABLE {
			return c.ExecCreateTableAs(s)
		}
		return errInvalidParameterValue(fmt.Sprintf("unsupported CreateTableAsStmt objtype: %d", s.Objtype))
	case *nodes.RefreshMatViewStmt:
		return c.ExecRefreshMatView(s)
	case *nodes.DropStmt:
		return c.RemoveObjects(s)
	case *nodes.AlterTableStmt:
		return c.AlterTableStmt(s)
	case *nodes.RenameStmt:
		return c.ExecRenameStmt(s)
	case *nodes.GrantStmt:
		return c.ExecGrantStmt(s)
	case *nodes.CreatePolicyStmt:
		return c.CreatePolicy(s)
	case *nodes.AlterPolicyStmt:
		return c.AlterPolicy(s)
	case *nodes.AlterFunctionStmt:
		return c.AlterFunction(s)
	case *nodes.AlterObjectSchemaStmt:
		return c.ExecAlterObjectSchemaStmt(s)
	case *nodes.AlterOwnerStmt:
		return nil // no-op: pgddl has no ownership tracking
	case *nodes.TruncateStmt:
		return c.ExecuteTruncate(s)
	case *nodes.AlterTypeStmt:
		// ALTER TYPE ... SET (...) / DROP ATTRIBUTE: no-op for pgddl
		// pg: src/backend/commands/typecmds.c — AlterType
		return nil
	case *nodes.DefineStmt:
		// pg: src/backend/tcop/utility.c — ProcessUtilitySlow (DefineStmt cases)
		switch s.Kind {
		case nodes.OBJECT_TYPE:
			return c.DefineType(s)
		case nodes.OBJECT_OPERATOR:
			return c.DefineOperator(s)
		case nodes.OBJECT_AGGREGATE:
			return c.DefineAggregate(s)
		default:
			// CREATE COLLATION/TEXT SEARCH/etc: no-op for pgddl
			return nil
		}
	case *nodes.RuleStmt:
		// CREATE RULE: no-op for pgddl
		// pg: src/backend/rewrite/rewriteDefine.c — DefineRule
		return nil
	case *nodes.CreateCastStmt:
		return c.CreateCast(s)
	case *nodes.CreateForeignTableStmt:
		// CREATE FOREIGN TABLE: delegate to DefineRelation with relkind='f'
		// pg: src/backend/commands/tablecmds.c — DefineRelation (foreign table)
		return c.DefineRelation(&s.Base, 'f')
	case *nodes.CreateStatsStmt:
		// CREATE STATISTICS: no-op for pgddl (planner hint)
		return nil
	case *nodes.AlterStatsStmt:
		// ALTER STATISTICS: no-op for pgddl (planner hint)
		return nil
	case *nodes.AlterDefaultPrivilegesStmt:
		// ALTER DEFAULT PRIVILEGES: no-op for pgddl (no ACL tracking)
		return nil
	case *nodes.AlterTSDictionaryStmt:
		// ALTER TEXT SEARCH DICTIONARY: no-op for pgddl
		return nil
	case *nodes.AlterTSConfigurationStmt:
		// ALTER TEXT SEARCH CONFIGURATION: no-op for pgddl
		return nil
	case *nodes.CreateExtensionStmt:
		return c.CreateExtension(s)
	case *nodes.ReindexStmt:
		// REINDEX: no-op for pgddl (physical index rebuild)
		// pg: src/backend/commands/indexcmds.c — ReindexTable etc.
		return nil
	case *nodes.CreateOpClassStmt:
		return c.DefineOpClass(s)
	case *nodes.CreateOpFamilyStmt:
		return c.CreateOpFamily(s)
	case *nodes.AlterOpFamilyStmt:
		return c.AlterOpFamily(s)
	case *nodes.CreateConversionStmt:
		// CREATE CONVERSION: no-op for pgddl
		return nil
	case *nodes.CreateAmStmt:
		return c.CreateAccessMethod(s)
	case *nodes.CreatePublicationStmt:
		// CREATE PUBLICATION: no-op for pgddl
		return nil
	case *nodes.AlterPublicationStmt:
		// ALTER PUBLICATION: no-op for pgddl
		return nil
	case *nodes.CreateSubscriptionStmt:
		// CREATE SUBSCRIPTION: no-op for pgddl
		return nil
	case *nodes.AlterSubscriptionStmt:
		// ALTER SUBSCRIPTION: no-op for pgddl
		return nil
	case *nodes.DropSubscriptionStmt:
		// DROP SUBSCRIPTION: no-op for pgddl
		return nil
	case *nodes.CreateTransformStmt:
		// CREATE TRANSFORM: no-op for pgddl
		return nil
	case *nodes.CreateEventTrigStmt:
		// CREATE EVENT TRIGGER: no-op for pgddl
		return nil
	case *nodes.AlterEventTrigStmt:
		// ALTER EVENT TRIGGER: no-op for pgddl
		return nil
	case *nodes.VariableSetStmt:
		// pg: src/backend/tcop/utility.c — standard_ProcessUtility (VariableSetStmt)
		return c.processVariableSet(s)
	case *nodes.TransactionStmt:
		// BEGIN/COMMIT/ROLLBACK: no-op for pgddl
		return nil
	// ---- Foreign data wrappers ----
	case *nodes.AlterExtensionStmt:
		// ALTER EXTENSION: no-op for pgddl
		return nil
	case *nodes.AlterExtensionContentsStmt:
		// ALTER EXTENSION ... ADD/DROP: no-op for pgddl
		return nil
	case *nodes.CreateFdwStmt:
		// CREATE FOREIGN DATA WRAPPER: no-op for pgddl
		return nil
	case *nodes.AlterFdwStmt:
		// ALTER FOREIGN DATA WRAPPER: no-op for pgddl
		return nil
	case *nodes.CreateForeignServerStmt:
		// CREATE SERVER: no-op for pgddl
		return nil
	case *nodes.AlterForeignServerStmt:
		// ALTER SERVER: no-op for pgddl
		return nil
	case *nodes.CreateUserMappingStmt:
		// CREATE USER MAPPING: no-op for pgddl
		return nil
	case *nodes.AlterUserMappingStmt:
		// ALTER USER MAPPING: no-op for pgddl
		return nil
	case *nodes.DropUserMappingStmt:
		// DROP USER MAPPING: no-op for pgddl
		return nil
	case *nodes.ImportForeignSchemaStmt:
		// IMPORT FOREIGN SCHEMA: no-op for pgddl
		return nil

	// ---- Operators/Collations ----
	case *nodes.AlterCollationStmt:
		// ALTER COLLATION: no-op for pgddl
		return nil
	case *nodes.AlterOperatorStmt:
		// ALTER OPERATOR: no-op for pgddl
		return nil

	// ---- Tablespaces ----
	case *nodes.CreateTableSpaceStmt:
		// CREATE TABLESPACE: no-op for pgddl (physical storage)
		return nil
	case *nodes.DropTableSpaceStmt:
		// DROP TABLESPACE: no-op for pgddl (physical storage)
		return nil
	case *nodes.AlterTableSpaceOptionsStmt:
		// ALTER TABLESPACE ... SET/RESET: no-op for pgddl
		return nil
	case *nodes.AlterTableMoveAllStmt:
		// ALTER TABLE ALL IN TABLESPACE ... SET TABLESPACE: no-op for pgddl
		return nil

	// ---- Roles/Databases ----
	case *nodes.CreateRoleStmt:
		// CREATE ROLE/USER/GROUP: no-op for pgddl
		return nil
	case *nodes.AlterRoleStmt:
		// ALTER ROLE: no-op for pgddl
		return nil
	case *nodes.AlterRoleSetStmt:
		// ALTER ROLE ... SET: no-op for pgddl
		return nil
	case *nodes.DropRoleStmt:
		// DROP ROLE/USER/GROUP: no-op for pgddl
		return nil
	case *nodes.GrantRoleStmt:
		// GRANT role TO role: no-op for pgddl
		return nil
	case *nodes.CreatedbStmt:
		// CREATE DATABASE: no-op for pgddl
		return nil
	case *nodes.AlterDatabaseStmt:
		// ALTER DATABASE: no-op for pgddl
		return nil
	case *nodes.AlterDatabaseSetStmt:
		// ALTER DATABASE ... SET: no-op for pgddl
		return nil
	case *nodes.AlterDatabaseRefreshCollStmt:
		// ALTER DATABASE ... REFRESH COLLATION VERSION: no-op for pgddl
		return nil
	case *nodes.DropdbStmt:
		// DROP DATABASE: no-op for pgddl
		return nil

	// ---- Misc DDL ----
	case *nodes.AlterObjectDependsStmt:
		// ALTER ... DEPENDS ON EXTENSION: no-op for pgddl
		return nil
	case *nodes.CreatePLangStmt:
		// CREATE LANGUAGE: no-op for pgddl
		return nil
	case *nodes.AlterSystemStmt:
		// ALTER SYSTEM: no-op for pgddl
		return nil
	case *nodes.DropOwnedStmt:
		// DROP OWNED: no-op for pgddl
		return nil
	case *nodes.ReassignOwnedStmt:
		// REASSIGN OWNED: no-op for pgddl
		return nil
	case *nodes.SecLabelStmt:
		// SECURITY LABEL: no-op for pgddl
		return nil

	// ---- Non-DDL utility statements ----
	case *nodes.CopyStmt:
		// COPY: no-op for pgddl
		return nil
	case *nodes.DoStmt:
		// DO: no-op for pgddl
		return nil
	case *nodes.ExplainStmt:
		// EXPLAIN: no-op for pgddl
		return nil
	case *nodes.VacuumStmt:
		// VACUUM/ANALYZE: no-op for pgddl
		return nil
	case *nodes.ClusterStmt:
		// CLUSTER: no-op for pgddl
		return nil
	case *nodes.CheckPointStmt:
		// CHECKPOINT: no-op for pgddl
		return nil
	case *nodes.DiscardStmt:
		// DISCARD: no-op for pgddl
		return nil
	case *nodes.LockStmt:
		// LOCK TABLE: no-op for pgddl
		return nil
	case *nodes.ListenStmt:
		// LISTEN: no-op for pgddl
		return nil
	case *nodes.UnlistenStmt:
		// UNLISTEN: no-op for pgddl
		return nil
	case *nodes.NotifyStmt:
		// NOTIFY: no-op for pgddl
		return nil
	case *nodes.LoadStmt:
		// LOAD: no-op for pgddl
		return nil
	case *nodes.ConstraintsSetStmt:
		// SET CONSTRAINTS: no-op for pgddl
		return nil
	case *nodes.VariableShowStmt:
		// SHOW: no-op for pgddl
		return nil
	case *nodes.CallStmt:
		// CALL: no-op for pgddl
		return nil
	case *nodes.PrepareStmt:
		// PREPARE: no-op for pgddl
		return nil
	case *nodes.ExecuteStmt:
		// EXECUTE: no-op for pgddl
		return nil
	case *nodes.DeallocateStmt:
		// DEALLOCATE: no-op for pgddl
		return nil

	default:
		return errInvalidParameterValue(fmt.Sprintf("unsupported utility statement: %T", stmt))
	}
}

// processVariableSet handles SET/RESET variable statements.
//
// pg: src/backend/utils/misc/guc.c — ExecSetVariableStmt
func (c *Catalog) processVariableSet(s *nodes.VariableSetStmt) error {
	switch s.Name {
	case "search_path":
		return c.processSetSearchPath(s)
	case "role":
		return c.processSetRole(s)
	}
	// All other SET/RESET commands: no-op for pgddl.
	return nil
}

// processSetSearchPath handles SET search_path = ...
//
// pg: src/backend/catalog/namespace.c — preprocessNamespacePath
func (c *Catalog) processSetSearchPath(s *nodes.VariableSetStmt) error {
	switch s.Kind {
	case nodes.VAR_SET_VALUE:
		if s.Args == nil {
			return nil
		}
		schemas := make([]string, 0, len(s.Args.Items))
		for _, arg := range s.Args.Items {
			ac, ok := arg.(*nodes.A_Const)
			if !ok {
				continue
			}
			sv, ok := ac.Val.(*nodes.String)
			if !ok {
				continue
			}
			schemas = append(schemas, sv.Str)
		}
		c.SetSearchPath(schemas)
	case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
		// RESET search_path / SET search_path TO DEFAULT
		// pg: GUC default for search_path is '"$user", public'
		c.SetSearchPath([]string{"$user", "public"})
	}
	return nil
}

// processSetRole handles SET ROLE, RESET ROLE, SET ROLE NONE.
//
// pg: src/backend/commands/variable.c — check_role, assign_role
// pg: src/backend/utils/init/miscinit.c — SetCurrentRoleId
func (c *Catalog) processSetRole(s *nodes.VariableSetStmt) error {
	switch s.Kind {
	case nodes.VAR_SET_VALUE:
		if s.Args == nil {
			return nil
		}
		// Extract the role name from the first argument.
		// SET ROLE rolename is parsed as SET role = 'rolename'.
		for _, arg := range s.Args.Items {
			ac, ok := arg.(*nodes.A_Const)
			if !ok {
				continue
			}
			sv, ok := ac.Val.(*nodes.String)
			if !ok {
				continue
			}
			roleName := sv.Str
			// pg: check_role — "none" maps to InvalidOid → reset
			if roleName == "none" {
				c.ResetRole()
			} else {
				c.SetRole(roleName)
			}
			return nil
		}
	case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
		// RESET ROLE / SET ROLE TO DEFAULT → revert to session user.
		// pg: SetCurrentRoleId(InvalidOid, ...) → falls back to SessionUserId
		c.ResetRole()
	}
	return nil
}
