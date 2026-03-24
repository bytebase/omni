package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// Grant represents a recorded privilege grant.
type Grant struct {
	ObjType   byte   // 'r'=relation, 's'=sequence, 'f'=function, 'n'=schema
	ObjOID    uint32 // OID of the object
	Grantee   string // role name ("" = PUBLIC)
	Privilege string // SELECT, INSERT, UPDATE, DELETE, REFERENCES, TRIGGER, EXECUTE, USAGE, CREATE, ALL
	Columns   []string
	WithGrant bool
}

// ExecGrantStmt processes GRANT and REVOKE statements.
// For pgddl, grants are stored for catalog completeness but do not affect DDL semantics.
//
// pg: src/backend/catalog/aclchk.c — ExecGrantStmt_oids
func (c *Catalog) ExecGrantStmt(stmt *nodes.GrantStmt) error {
	// Extract privileges.
	privs := c.extractPrivileges(stmt.Privileges)

	// Extract grantees.
	grantees := c.extractGrantees(stmt.Grantees)

	// Determine object type code.
	objTypeCode := grantObjTypeCode(nodes.ObjectType(stmt.Objtype))

	if stmt.IsGrant {
		return c.execGrant(stmt, objTypeCode, privs, grantees)
	}
	return c.execRevoke(stmt, objTypeCode, privs, grantees)
}

func (c *Catalog) execGrant(stmt *nodes.GrantStmt, objTypeCode byte, privs []grantPriv, grantees []string) error {
	// PG forbids WITH GRANT OPTION to PUBLIC.
	// pg: src/backend/catalog/aclchk.c — aclcheck (line 208-211)
	if stmt.GrantOption {
		for _, g := range grantees {
			if g == "" { // empty string represents PUBLIC
				return &Error{Code: CodeInvalidGrantOperation,
					Message: "grant options can only be granted to roles"}
			}
		}
	}

	// Validate privileges are valid for the object type.
	// pg: src/backend/catalog/aclchk.c — ExecGrant_Relation, ExecGrant_Sequence, ExecGrant_Function, ExecGrant_Namespace
	for _, priv := range privs {
		if !isValidPrivilege(priv.name, objTypeCode) {
			return &Error{Code: CodeInvalidGrantOperation,
				Message: fmt.Sprintf("invalid privilege type %s for %s", priv.name, grantObjTypeName(objTypeCode))}
		}
	}

	// Resolve each object and record grants.
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		objOID, err := c.resolveGrantObject(objTypeCode, obj)
		if err != nil {
			return err
		}
		for _, grantee := range grantees {
			for _, priv := range privs {
				c.grants = append(c.grants, Grant{
					ObjType:   objTypeCode,
					ObjOID:    objOID,
					Grantee:   grantee,
					Privilege: priv.name,
					Columns:   priv.columns,
					WithGrant: stmt.GrantOption,
				})
			}
		}
	}
	return nil
}

func (c *Catalog) execRevoke(stmt *nodes.GrantStmt, objTypeCode byte, privs []grantPriv, grantees []string) error {
	// Resolve each object and remove matching grants.
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		objOID, err := c.resolveGrantObject(objTypeCode, obj)
		if err != nil {
			return err
		}
		c.revokeGrants(objTypeCode, objOID, grantees, privs)
	}
	return nil
}

// revokeGrants removes matching grants from the catalog.
func (c *Catalog) revokeGrants(objType byte, objOID uint32, grantees []string, privs []grantPriv) {
	n := 0
	for _, g := range c.grants {
		if g.ObjType == objType && g.ObjOID == objOID && matchGrantee(g.Grantee, grantees) && matchPriv(g.Privilege, privs) {
			continue // remove
		}
		c.grants[n] = g
		n++
	}
	c.grants = c.grants[:n]
}

func matchGrantee(grantee string, grantees []string) bool {
	for _, g := range grantees {
		if g == grantee {
			return true
		}
	}
	return false
}

func matchPriv(priv string, privs []grantPriv) bool {
	for _, p := range privs {
		if p.name == priv || p.name == "ALL" {
			return true
		}
	}
	return false
}

// grantPriv holds a privilege name and optional column list.
type grantPriv struct {
	name    string
	columns []string
}

func (c *Catalog) extractPrivileges(privList *nodes.List) []grantPriv {
	if privList == nil {
		// nil means ALL PRIVILEGES.
		return []grantPriv{{name: "ALL"}}
	}
	var result []grantPriv
	for _, item := range privList.Items {
		ap, ok := item.(*nodes.AccessPriv)
		if !ok {
			continue
		}
		name := ap.PrivName
		if name == "" {
			name = "ALL"
		}
		var cols []string
		if ap.Cols != nil {
			for _, col := range ap.Cols.Items {
				cols = append(cols, stringVal(col))
			}
		}
		result = append(result, grantPriv{name: name, columns: cols})
	}
	return result
}

func (c *Catalog) extractGrantees(granteeList *nodes.List) []string {
	if granteeList == nil {
		return nil
	}
	var result []string
	for _, item := range granteeList.Items {
		rs, ok := item.(*nodes.RoleSpec)
		if !ok {
			continue
		}
		result = append(result, roleSpecName(rs))
	}
	return result
}

// roleSpecName returns the role name from a RoleSpec node.
func roleSpecName(rs *nodes.RoleSpec) string {
	switch nodes.RoleSpecType(rs.Roletype) {
	case nodes.ROLESPEC_CSTRING:
		return rs.Rolename
	case nodes.ROLESPEC_CURRENT_ROLE, nodes.ROLESPEC_CURRENT_USER:
		return "CURRENT_USER"
	case nodes.ROLESPEC_SESSION_USER:
		return "SESSION_USER"
	case nodes.ROLESPEC_PUBLIC:
		return ""
	default:
		return rs.Rolename
	}
}

func grantObjTypeCode(objType nodes.ObjectType) byte {
	switch objType {
	case nodes.OBJECT_TABLE, nodes.OBJECT_VIEW, nodes.OBJECT_MATVIEW:
		return 'r'
	case nodes.OBJECT_SEQUENCE:
		return 's'
	case nodes.OBJECT_FUNCTION, nodes.OBJECT_PROCEDURE, nodes.OBJECT_ROUTINE:
		return 'f'
	case nodes.OBJECT_SCHEMA:
		return 'n'
	case nodes.OBJECT_TYPE, nodes.OBJECT_DOMAIN:
		return 'T'
	default:
		return 'r'
	}
}

func (c *Catalog) resolveGrantObject(objTypeCode byte, obj nodes.Node) (uint32, error) {
	switch objTypeCode {
	case 'r':
		// Table/view: obj is a *RangeVar.
		rv, ok := obj.(*nodes.RangeVar)
		if !ok {
			return 0, errInvalidParameterValue("expected RangeVar for GRANT on relation")
		}
		_, rel, err := c.findRelation(rv.Schemaname, rv.Relname)
		if err != nil {
			return 0, err
		}
		return rel.OID, nil
	case 's':
		// Sequence: obj is a *RangeVar.
		rv, ok := obj.(*nodes.RangeVar)
		if !ok {
			return 0, errInvalidParameterValue("expected RangeVar for GRANT on sequence")
		}
		seq, err := c.findSequence(rv.Schemaname, rv.Relname)
		if err != nil {
			return 0, err
		}
		return seq.OID, nil
	case 'f':
		// Function: obj is *ObjectWithArgs.
		owa, ok := obj.(*nodes.ObjectWithArgs)
		if !ok {
			return 0, errInvalidParameterValue("expected ObjectWithArgs for GRANT on function")
		}
		schemaName, funcName := qualifiedName(owa.Objname)
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return 0, err
		}
		candidates := c.findUserProcsByName(schema, funcName)
		if len(candidates) == 0 {
			return 0, errUndefinedObject("function", funcName)
		}
		return candidates[0].OID, nil
	case 'n':
		// Schema: obj may be a *String or a *List wrapping a single *String.
		// pg grammar: GRANT ... ON SCHEMA name_list produces List nodes.
		var name string
		if lst, ok := obj.(*nodes.List); ok && len(lst.Items) > 0 {
			name = stringVal(lst.Items[len(lst.Items)-1])
		} else {
			name = stringVal(obj)
		}
		s := c.schemaByName[name]
		if s == nil {
			return 0, errUndefinedSchema(name)
		}
		return s.OID, nil
	case 'T':
		// Type/domain: obj is a *nodes.List of name parts (like DROP TYPE).
		schemaName, typeName := extractDropObjectName(obj)
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return 0, err
		}
		bt := c.typeByName[typeKey{ns: schema.OID, name: typeName}]
		if bt == nil {
			return 0, errUndefinedType(typeName)
		}
		return bt.OID, nil
	default:
		return 0, errInvalidParameterValue("unsupported grant object type")
	}
}

// removeGrantsForObject removes all grants on the given object.
func (c *Catalog) removeGrantsForObject(objType byte, objOID uint32) {
	n := 0
	for _, g := range c.grants {
		if g.ObjType == objType && g.ObjOID == objOID {
			continue
		}
		c.grants[n] = g
		n++
	}
	c.grants = c.grants[:n]
}

// isValidPrivilege checks whether a privilege name is valid for a given object type code.
//
// pg: src/backend/catalog/aclchk.c — ExecGrant_Relation, ExecGrant_Sequence, etc. (privilege validation)
func isValidPrivilege(priv string, objTypeCode byte) bool {
	p := strings.ToUpper(priv)
	if p == "ALL" {
		return true
	}
	switch objTypeCode {
	case 'r': // relation (table/view/matview)
		switch p {
		case "SELECT", "INSERT", "UPDATE", "DELETE", "REFERENCES", "TRIGGER", "TRUNCATE":
			return true
		}
	case 's': // sequence
		switch p {
		case "USAGE", "SELECT", "UPDATE":
			return true
		}
	case 'f': // function/procedure
		switch p {
		case "EXECUTE":
			return true
		}
	case 'n': // schema
		switch p {
		case "CREATE", "USAGE":
			return true
		}
	case 'T': // type/domain
		switch p {
		case "USAGE":
			return true
		}
	}
	return false
}

// grantObjTypeName returns a human-readable name for a grant object type code.
func grantObjTypeName(objTypeCode byte) string {
	switch objTypeCode {
	case 'r':
		return "relation"
	case 's':
		return "sequence"
	case 'f':
		return "function"
	case 'n':
		return "schema"
	case 'T':
		return "type"
	default:
		return "object"
	}
}
