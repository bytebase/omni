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
		return c.ExecAlterOwnerStmt(s)
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
		case nodes.OBJECT_COLLATION:
			schemaName, name := qualifiedName(s.Defnames)
			schema, err := c.resolveTargetSchema(schemaName)
			if err != nil {
				return err
			}
			if name != "" {
				if schema.Collations != nil && schema.Collations[name] != 0 {
					if s.IfNotExists {
						c.addWarning(CodeWarningSkip, fmt.Sprintf("collation %q already exists, skipping", name))
						return nil
					}
					return errDuplicateObject("collation", name)
				}
				if schema.Collations == nil {
					schema.Collations = make(map[string]uint32)
				}
				schema.Collations[name] = c.oidGen.Next()
			}
			return nil
		default:
			// CREATE TEXT SEARCH/etc: no-op for pgddl
			return nil
		}
	case *nodes.RuleStmt:
		if s.Relation != nil && s.Rulename != "" {
			if _, rel, err := c.findRelation(s.Relation.Schemaname, s.Relation.Relname); err == nil && rel != nil {
				if rel.Rules == nil {
					rel.Rules = make(map[string]byte)
				}
				rel.Rules[s.Rulename] = 'O'
				rel.HasRules = true
			}
		}
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
		return c.CreateEventTriggerStmt(s)
	case *nodes.AlterEventTrigStmt:
		return c.AlterEventTrigger(s)
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
		return c.DefineTablespace(s)
	case *nodes.DropTableSpaceStmt:
		return c.RemoveTablespace(s)
	case *nodes.AlterTableSpaceOptionsStmt:
		return c.AlterTablespaceOptions(s)
	case *nodes.AlterTableMoveAllStmt:
		return c.MoveAllTablespaces(s)

	// ---- Roles/Databases ----
	case *nodes.CreateRoleStmt:
		return c.DefineRole(s)
	case *nodes.AlterRoleStmt:
		// ALTER ROLE: no-op for pgddl
		return nil
	case *nodes.AlterRoleSetStmt:
		// ALTER ROLE ... SET: no-op for pgddl
		return nil
	case *nodes.DropRoleStmt:
		return c.RemoveRoles(s)
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
		return c.DropOwned(s)
	case *nodes.ReassignOwnedStmt:
		return c.ReassignOwned(s)
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
		if s.Target == nodes.DISCARD_ALL {
			clear(c.preparedStmts)
		}
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
		if s.Name != "" && c.preparedStmts[s.Name] != nil {
			return &Error{
				Code:    CodeDuplicatePstatement,
				Message: fmt.Sprintf("prepared statement %q already exists", s.Name),
			}
		}
		var argTypes []uint32
		if s.Argtypes != nil {
			argTypes = make([]uint32, 0, len(s.Argtypes.Items))
			for _, item := range s.Argtypes.Items {
				tn, ok := item.(*nodes.TypeName)
				if !ok {
					continue
				}
				oid, _, err := c.resolveTypeName(tn)
				if err != nil {
					return err
				}
				argTypes = append(argTypes, oid)
			}
		}
		numParams := len(argTypes)
		if maxParam := maxParamRefInNode(s.Query); maxParam > numParams {
			numParams = maxParam
		}
		if s.Name != "" {
			if c.preparedStmts == nil {
				c.preparedStmts = make(map[string]*preparedStmt)
			}
			c.preparedStmts[s.Name] = &preparedStmt{
				query:     s.Query,
				argTypes:  argTypes,
				numParams: numParams,
			}
		}
		return nil
	case *nodes.ExecuteStmt:
		// EXECUTE: no-op for pgddl
		return nil
	case *nodes.DeallocateStmt:
		if s.IsAll || s.Name == "" {
			clear(c.preparedStmts)
		} else {
			delete(c.preparedStmts, s.Name)
		}
		return nil

	default:
		return errInvalidParameterValue(fmt.Sprintf("unsupported utility statement: %T", stmt))
	}
}

type maxParamVisitor struct {
	max int
}

func (v *maxParamVisitor) Visit(node nodes.Node) nodes.Visitor {
	if node == nil {
		return nil
	}
	if pr, ok := node.(*nodes.ParamRef); ok && pr.Number > v.max {
		v.max = pr.Number
	}
	return v
}

func maxParamRefInNode(node nodes.Node) int {
	if node == nil {
		return 0
	}
	var v maxParamVisitor
	nodes.Walk(&v, node)
	return v.max
}

// processVariableSet handles SET/RESET variable statements.
//
// pg: src/backend/utils/misc/guc.c — ExecSetVariableStmt
func (c *Catalog) processVariableSet(s *nodes.VariableSetStmt) error {
	if s.Kind == nodes.VAR_RESET_ALL {
		c.defaultTablespace = ""
		c.defaultTableAccessMethod = ""
	}
	switch s.Name {
	case "search_path":
		return c.processSetSearchPath(s)
	case "default_tablespace":
		return c.processSetDefaultTablespace(s)
	case "default_table_access_method":
		return c.processSetDefaultTableAccessMethod(s)
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

// processSetDefaultTablespace handles SET default_tablespace = ...
func (c *Catalog) processSetDefaultTablespace(s *nodes.VariableSetStmt) error {
	switch s.Kind {
	case nodes.VAR_SET_VALUE:
		if s.Args == nil {
			c.defaultTablespace = ""
			return nil
		}
		for _, arg := range s.Args.Items {
			ac, ok := arg.(*nodes.A_Const)
			if !ok {
				continue
			}
			sv, ok := ac.Val.(*nodes.String)
			if !ok {
				continue
			}
			if sv.Str != "" && c.checkTablespace(sv.Str) != nil {
				return errInvalidParameterValue(fmt.Sprintf("invalid value for parameter \"default_tablespace\": %q", sv.Str))
			}
			c.defaultTablespace = sv.Str
		}
	case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
		c.defaultTablespace = ""
	}
	return nil
}

// processSetDefaultTableAccessMethod handles SET default_table_access_method.
func (c *Catalog) processSetDefaultTableAccessMethod(s *nodes.VariableSetStmt) error {
	switch s.Kind {
	case nodes.VAR_SET_VALUE:
		if s.Args == nil {
			c.defaultTableAccessMethod = ""
			return nil
		}
		for _, arg := range s.Args.Items {
			ac, ok := arg.(*nodes.A_Const)
			if !ok {
				continue
			}
			sv, ok := ac.Val.(*nodes.String)
			if !ok {
				continue
			}
			if sv.Str != "" && !c.snapshot {
				am := c.accessMethodByName[sv.Str]
				if am == nil || am.Type != 't' {
					return errInvalidParameterValue(fmt.Sprintf("invalid value for parameter \"default_table_access_method\": %q", sv.Str))
				}
			}
			c.defaultTableAccessMethod = sv.Str
		}
	case nodes.VAR_SET_DEFAULT, nodes.VAR_RESET:
		c.defaultTableAccessMethod = ""
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
				if err := c.checkRole(roleName); err != nil {
					return err
				}
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
