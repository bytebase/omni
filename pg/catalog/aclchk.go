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
	Grantor   string // role that issued the grant
	Privilege string // SELECT, INSERT, UPDATE, DELETE, REFERENCES, TRIGGER, EXECUTE, USAGE, CREATE, ALL
	Columns   []string
	WithGrant bool
	// FromExtension marks a grant an extension's script made: part of the
	// extension's baseline.
	FromExtension bool
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
	// A dump grants to cluster-wide roles it does not create, so a grantee is
	// recorded as observed rather than required — but only once the statement
	// has succeeded, as PostgreSQL's rollback would undo an earlier record.

	// Determine object type code.
	objTypeCode := grantObjTypeCode(nodes.ObjectType(stmt.Objtype))

	grantor := c.currentOwner()
	if stmt.Grantor != nil {
		grantor = c.ownerRoleName(stmt.Grantor)
		if err := c.checkRole(grantor); err != nil {
			return err
		}
	}
	var execErr error
	if stmt.IsGrant {
		execErr = c.execGrant(stmt, objTypeCode, privs, grantees, grantor)
	} else {
		execErr = c.execRevoke(stmt, objTypeCode, privs, grantees, grantor)
	}
	if execErr != nil {
		return execErr
	}
	if stmt.IsGrant {
		for _, grantee := range grantees {
			if grantee != "" {
				c.roles[grantee] = true
			}
		}
	}
	return nil
}

func (c *Catalog) execGrant(stmt *nodes.GrantStmt, objTypeCode byte, privs []grantPriv, grantees []string, grantor string) error {
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

	// Resolve all target objects before mutating c.grants.
	if stmt.Objects == nil {
		return nil
	}
	objOIDs := make([]uint32, 0, len(stmt.Objects.Items))
	for _, obj := range stmt.Objects.Items {
		objOID, err := c.resolveGrantObject(stmt.Objtype, objTypeCode, obj)
		if err != nil {
			return err
		}
		objOIDs = append(objOIDs, objOID)
	}
	for _, objOID := range objOIDs {
		for _, grantee := range grantees {
			for _, priv := range privs {
				c.grants = append(c.grants, Grant{
					ObjType:   objTypeCode,
					ObjOID:    objOID,
					Grantee:   grantee,
					Grantor:   grantor,
					Privilege: priv.name,
					Columns:   priv.columns,
					WithGrant: stmt.GrantOption,
				})
			}
		}
	}
	return nil
}

func (c *Catalog) execRevoke(stmt *nodes.GrantStmt, objTypeCode byte, privs []grantPriv, grantees []string, grantor string) error {
	for _, priv := range privs {
		if !isValidPrivilege(priv.name, objTypeCode) {
			return &Error{Code: CodeInvalidGrantOperation,
				Message: fmt.Sprintf("invalid privilege type %s for %s", priv.name, grantObjTypeName(objTypeCode))}
		}
	}
	if stmt.Objects == nil {
		return nil
	}
	objOIDs := make([]uint32, 0, len(stmt.Objects.Items))
	for _, obj := range stmt.Objects.Items {
		objOID, err := c.resolveGrantObject(stmt.Objtype, objTypeCode, obj)
		if err != nil {
			return err
		}
		objOIDs = append(objOIDs, objOID)
	}
	if len(c.grants) == 0 {
		return nil
	}

	targetOIDs := make(map[uint32]bool, len(objOIDs))
	for _, oid := range objOIDs {
		targetOIDs[oid] = true
	}
	cascade := stmt.Behavior == nodes.DROP_CASCADE
	grantPrivCovers := func(held, passed string) bool {
		h := strings.ToUpper(held)
		p := strings.ToUpper(passed)
		return h == "ALL" || p == "ALL" || h == p
	}

	trial := make([]Grant, len(c.grants))
	copy(trial, c.grants)
	removed := make([]bool, len(trial))
	lostGrantOpt := make([]bool, len(trial))
	for i, g := range trial {
		if g.ObjType == objTypeCode && targetOIDs[g.ObjOID] && g.Grantor == grantor &&
			matchGrantee(g.Grantee, grantees) && matchPriv(g.Privilege, privs) {
			if g.WithGrant {
				lostGrantOpt[i] = true
			}
			if stmt.GrantOption {
				trial[i].WithGrant = false
			} else {
				removed[i] = true
			}
		}
	}

	for changed := true; changed; {
		changed = false
		for i, g := range trial {
			if removed[i] || g.Grantor == "" || c.objectOwnerName(g.ObjType, g.ObjOID) == g.Grantor {
				continue
			}
			lost := false
			stillHas := false
			for j, src := range trial {
				if src.ObjType != g.ObjType || src.ObjOID != g.ObjOID || src.Grantee != g.Grantor {
					continue
				}
				if !grantPrivCovers(c.grants[j].Privilege, g.Privilege) {
					continue
				}
				if lostGrantOpt[j] {
					lost = true
				}
				if !removed[j] && src.WithGrant {
					stillHas = true
				}
			}
			if lost && !stillHas {
				if !cascade {
					return errDependentObjects("grant", g.Grantor)
				}
				if g.WithGrant {
					lostGrantOpt[i] = true
				}
				removed[i] = true
				changed = true
			}
		}
	}

	kept := c.grants[:0]
	for i, g := range trial {
		if !removed[i] {
			kept = append(kept, g)
		}
	}
	c.grants = kept
	return nil
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
		if nodes.RoleSpecType(rs.Roletype) == nodes.ROLESPEC_PUBLIC {
			result = append(result, "")
		} else {
			result = append(result, c.ownerRoleName(rs))
		}
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
	case nodes.OBJECT_TABLE, nodes.OBJECT_VIEW, nodes.OBJECT_MATVIEW, nodes.OBJECT_FOREIGN_TABLE:
		return 'r'
	case nodes.OBJECT_SEQUENCE:
		return 's'
	case nodes.OBJECT_FUNCTION, nodes.OBJECT_PROCEDURE, nodes.OBJECT_ROUTINE, nodes.OBJECT_AGGREGATE:
		return 'f'
	case nodes.OBJECT_SCHEMA:
		return 'n'
	case nodes.OBJECT_TYPE, nodes.OBJECT_DOMAIN:
		return 'T'
	default:
		return 'r'
	}
}

func (c *Catalog) resolveGrantObject(objType nodes.ObjectType, objTypeCode byte, obj nodes.Node) (uint32, error) {
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
		// Function/procedure/routine/aggregate: obj is *ObjectWithArgs.
		owa, ok := obj.(*nodes.ObjectWithArgs)
		if !ok {
			return 0, errInvalidParameterValue("expected ObjectWithArgs for GRANT on function")
		}
		return c.resolveGrantRoutine(objType, owa)
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
		var bt *BuiltinType
		if schemaName != "" {
			schema, err := c.resolveTargetSchema(schemaName)
			if err != nil {
				return 0, err
			}
			bt = c.typeByName[typeKey{ns: schema.OID, name: typeName}]
		} else {
			for _, nsOID := range c.searchPathWithCatalog() {
				if found := c.typeByName[typeKey{ns: nsOID, name: typeName}]; found != nil {
					bt = found
					break
				}
			}
		}
		if bt == nil {
			return 0, errUndefinedType(typeName)
		}
		if objType == nodes.OBJECT_DOMAIN && bt.Type != 'd' {
			return 0, errWrongObjectType(typeName, "a domain")
		}
		return bt.OID, nil
	default:
		return 0, errInvalidParameterValue("unsupported grant object type")
	}
}

func (c *Catalog) resolveGrantRoutine(objType nodes.ObjectType, owa *nodes.ObjectWithArgs) (uint32, error) {
	kindStr := "function"
	switch objType {
	case nodes.OBJECT_PROCEDURE:
		kindStr = "procedure"
	case nodes.OBJECT_AGGREGATE:
		kindStr = "aggregate"
	case nodes.OBJECT_ROUTINE:
		kindStr = "routine"
	}
	schemaName, funcName := qualifiedName(owa.Objname)
	argOIDs, argErr := c.routineArgOIDs(owa)
	if argErr != nil {
		return 0, argErr
	}
	var schemasToSearch []uint32
	if schemaName != "" {
		schema, err := c.resolveTargetSchema(schemaName)
		if err != nil {
			return 0, err
		}
		schemasToSearch = []uint32{schema.OID}
	} else {
		schemasToSearch = c.searchPathWithCatalog()
	}

	checkRoutineKind := func(kind byte) error {
		switch objType {
		case nodes.OBJECT_PROCEDURE:
			if kind != 'p' {
				return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not a procedure", funcName)}
			}
		case nodes.OBJECT_FUNCTION:
			if kind == 'p' {
				return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not a function", funcName)}
			}
		case nodes.OBJECT_AGGREGATE:
			if kind != PROKIND_AGGREGATE {
				return &Error{Code: CodeWrongObjectType, Message: fmt.Sprintf("%q is not an aggregate", funcName)}
			}
		}
		return nil
	}
	procSchemaOID := func(p *BuiltinProc) (uint32, bool) {
		if p.OID < FirstNormalObjectId {
			return PGCatalogNamespace, true
		}
		if up := c.userProcs[p.OID]; up != nil && up.Schema != nil {
			return up.Schema.OID, true
		}
		return 0, false
	}

	if argOIDs != nil {
		for _, nsOID := range schemasToSearch {
			for _, p := range c.procByName[funcName] {
				pNS, ok := procSchemaOID(p)
				if !ok || pNS != nsOID || int(p.NArgs) != len(argOIDs) {
					continue
				}
				match := true
				for i, a := range argOIDs {
					if p.ArgTypes[i] != a {
						match = false
						break
					}
				}
				if match {
					if err := checkRoutineKind(p.Kind); err != nil {
						return 0, err
					}
					return p.OID, nil
				}
			}
		}
		return 0, &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("%s %s does not exist", kindStr, funcName)}
	}

	for _, nsOID := range schemasToSearch {
		var matches []*BuiltinProc
		for _, p := range c.procByName[funcName] {
			pNS, ok := procSchemaOID(p)
			if !ok || pNS != nsOID {
				continue
			}
			if checkRoutineKind(p.Kind) == nil {
				matches = append(matches, p)
			}
		}
		if len(matches) > 1 {
			return 0, errAmbiguousFunction(funcName)
		}
		if len(matches) == 1 {
			return matches[0].OID, nil
		}
	}
	return 0, &Error{Code: CodeUndefinedFunction, Message: fmt.Sprintf("%s %s does not exist", kindStr, funcName)}
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
