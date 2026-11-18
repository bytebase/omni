package catalog

import (
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// Grant represents a recorded privilege grant.
type Grant struct {
	ObjType   byte   // 'r'=relation, 's'=sequence, 'f'=function, 'n'=schema
	ObjOID    uint32 // OID of the object
	Grantee   string // role name ("" = PUBLIC)
	Privilege string // SELECT, INSERT, UPDATE, DELETE, REFERENCES, TRIGGER, EXECUTE, USAGE, CREATE, ALL
	Columns   []string
	WithGrant bool

	// Target is the AST object type the grant was written with, where the
	// catalog's code cannot say it: an external table is a relation ('r') and
	// an external schema a schema ('n'), and each renders with ON EXTERNAL.
	Target nodes.ObjectType
	// GranteeKind is the grantee alternative the grant was written with —
	// "iam_role", "namespace", "account", "account_via_data_catalog",
	// "admins_of_account", "datashare", "rls_policy" — or "" for a role name.
	GranteeKind string
	// GranteeName is the grantee's own name, without its kind's prefix.
	GranteeName string
}

// grantee is one TO or FROM name a GRANT or REVOKE lists.
type grantee struct {
	name            string // the identity the catalog matches on: the kind-prefixed name, "" for PUBLIC
	kind            string // the alternative it was written with, "" for a role name
	value           string // its name without the kind's prefix
	withGrantOption bool   // the grantee itself carried WITH GRANT OPTION
}

// redshiftGranteeKind is the kind a RoleSpec names, "" for a database role.
func redshiftGranteeKind(rs *nodes.RoleSpec) string {
	switch nodes.RoleSpecType(rs.Roletype) {
	case nodes.ROLESPEC_IAM_ROLE:
		return "iam_role"
	case nodes.ROLESPEC_IAM_ROLE_DEFAULT:
		return "iam_role_default"
	case nodes.ROLESPEC_NAMESPACE:
		return "namespace"
	case nodes.ROLESPEC_ACCOUNT:
		return "account"
	case nodes.ROLESPEC_ACCOUNT_VIA_DATA_CATALOG:
		return "account_via_data_catalog"
	case nodes.ROLESPEC_ADMINS_OF_ACCOUNT:
		return "admins_of_account"
	case nodes.ROLESPEC_DATASHARE:
		return "datashare"
	case nodes.ROLESPEC_RLS_POLICY:
		return "rls_policy"
	}
	return ""
}

// newGrantee reads a grantee from its RoleSpec.
func newGrantee(rs *nodes.RoleSpec) grantee {
	g := grantee{name: roleSpecName(rs), kind: redshiftGranteeKind(rs), withGrantOption: rs.WithGrantOption}
	if g.kind != "" {
		g.value = strings.TrimPrefix(rs.Rolename, g.kind+":")
	}
	return g
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

	if stmt.Targtype == nodes.ACL_TARGET_SCOPED {
		// A scoped permission covers objects that do not exist yet: it names no
		// object this catalog resolves and records nothing.
		return nil
	}

	// Determine object type code.
	objTypeCode := grantObjTypeCode(nodes.ObjectType(stmt.Objtype))
	if objTypeCode == 'x' {
		// A Redshift object the catalog does not model — an external table or
		// schema, a datashare, a COPY job, a model — records nothing: its grants
		// name no object this catalog resolves, and no DDL statement reads them.
		// A grant option to PUBLIC is still refused, as for every object type.
		return checkGrantOptionNotPublic(stmt, grantees)
	}
	if objTypeCode == 'c' {
		// A COPY template is no catalog object: its grants name nothing this
		// catalog resolves, and no DDL statement reads them. What a grant to
		// PUBLIC may not carry, and which privileges the object takes, are still
		// checked, as for every object type.
		if err := checkGrantOptionNotPublic(stmt, grantees); err != nil {
			return err
		}
		for _, priv := range privs {
			if !isValidPrivilege(priv.name, objTypeCode) {
				return &Error{Code: CodeInvalidGrantOperation,
					Message: fmt.Sprintf("invalid privilege type %s for %s", priv.name, grantObjTypeName(objTypeCode))}
			}
		}
		return nil
	}

	if stmt.IsGrant {
		return c.execGrant(stmt, objTypeCode, privs, grantees)
	}
	return c.execRevoke(stmt, objTypeCode, privs, grantees)
}

// checkGrantOptionNotPublic refuses WITH GRANT OPTION to PUBLIC.
// pg: src/backend/catalog/aclchk.c — aclcheck (line 208-211)
func checkGrantOptionNotPublic(stmt *nodes.GrantStmt, grantees []grantee) error {
	if !stmt.IsGrant {
		return nil
	}
	// The option is the grantee's where the parser recorded it per grantee; a
	// statement-wide flag handed in directly applies to them all.
	perGrantee := false
	for _, g := range grantees {
		perGrantee = perGrantee || g.withGrantOption
	}
	for _, g := range grantees {
		if g.name != "" { // not PUBLIC
			continue
		}
		if g.withGrantOption || (stmt.GrantOption && !perGrantee) {
			return &Error{Code: CodeInvalidGrantOperation,
				Message: "grant options can only be granted to roles"}
		}
	}
	return nil
}

func (c *Catalog) execGrant(stmt *nodes.GrantStmt, objTypeCode byte, privs []grantPriv, grantees []grantee) error {
	if err := checkGrantOptionNotPublic(stmt, grantees); err != nil {
		return err
	}

	// Validate privileges are valid for the object type.
	// pg: src/backend/catalog/aclchk.c — ExecGrant_Relation, ExecGrant_Sequence, ExecGrant_Function, ExecGrant_Namespace
	for _, priv := range privs {
		if !validGrantPrivilege(nodes.ObjectType(stmt.Objtype), priv.name, objTypeCode) {
			return &Error{Code: CodeInvalidGrantOperation,
				Message: fmt.Sprintf("invalid privilege type %s for %s", priv.name, grantObjTypeName(objTypeCode))}
		}
	}

	// Resolve each object and record grants.
	if stmt.Objects == nil {
		return nil
	}
	perGrantee := false
	for _, grantee := range grantees {
		perGrantee = perGrantee || grantee.withGrantOption
	}
	for _, obj := range stmt.Objects.Items {
		objOID, err := c.resolveGrantObject(objTypeCode, obj)
		if err != nil {
			return err
		}
		for _, grantee := range grantees {
			// The option belongs to the grantee it follows where the parser
			// recorded it there, and to the statement otherwise.
			withGrant := grantee.withGrantOption || (stmt.GrantOption && !perGrantee)
			for _, priv := range privs {
				c.grants = append(c.grants, Grant{
					ObjType:     objTypeCode,
					ObjOID:      objOID,
					Grantee:     grantee.name,
					Privilege:   priv.name,
					Columns:     priv.columns,
					WithGrant:   withGrant,
					Target:      nodes.ObjectType(stmt.Objtype),
					GranteeKind: grantee.kind,
					GranteeName: grantee.value,
				})
			}
		}
	}
	return nil
}

func (c *Catalog) execRevoke(stmt *nodes.GrantStmt, objTypeCode byte, privs []grantPriv, grantees []grantee) error {
	// Resolve each object and remove matching grants.
	if stmt.Objects == nil {
		return nil
	}
	for _, obj := range stmt.Objects.Items {
		objOID, err := c.resolveGrantObject(objTypeCode, obj)
		if err != nil {
			return err
		}
		c.revokeGrants(nodes.ObjectType(stmt.Objtype), objTypeCode, objOID, grantees, privs)
	}
	return nil
}

// revokeGrants removes matching grants from the catalog. A grant matches on its
// target and grantee kind too: a relation grant to a role named "iam_role:r" is
// not the external-table grant to IAM_ROLE 'r' of the same name.
func (c *Catalog) revokeGrants(target nodes.ObjectType, objType byte, objOID uint32, grantees []grantee, privs []grantPriv) {
	n := 0
	for _, g := range c.grants {
		if g.ObjType != objType || g.ObjOID != objOID || !sameGrantTarget(g.Target, target) || !matchGrantee(g, grantees) {
			c.grants[n] = g
			n++
			continue
		}
		kept, keep := revokeOnGrant(g, privs)
		if keep {
			c.grants[n] = kept
			n++
		}
	}
	c.grants = c.grants[:n]
}

// sameGrantTarget reports whether two grant targets name the same object: a
// ROUTINE is a function or a procedure, according to PostgreSQL's GRANT, which
// this catalog inherits.
func sameGrantTarget(a, b nodes.ObjectType) bool {
	if a == b {
		return true
	}
	routine := func(x nodes.ObjectType) bool {
		return x == nodes.OBJECT_FUNCTION || x == nodes.OBJECT_PROCEDURE
	}
	return (a == nodes.OBJECT_ROUTINE && routine(b)) || (b == nodes.OBJECT_ROUTINE && routine(a))
}

// matchGrantee matches a recorded grant's grantee against a list of grantees.
func matchGrantee(g Grant, grantees []grantee) bool {
	for _, want := range grantees {
		if g.Grantee == want.name && g.GranteeKind == want.kind {
			return true
		}
	}
	return false
}

// revokeOnGrant applies a REVOKE's privileges to one recorded grant. A privilege
// revoke drops the grant; a column revoke takes the columns it names out of a
// larger column grant and keeps the rest; ALL PRIVILEGES takes everything.
func revokeOnGrant(g Grant, privs []grantPriv) (Grant, bool) {
	for _, p := range privs {
		if p.name != "ALL" && p.name != g.Privilege {
			continue
		}
		if p.name == "ALL" || len(p.columns) == 0 || len(g.Columns) == 0 {
			return g, false
		}
		rest := subtractGrantColumns(g.Columns, p.columns)
		if len(rest) == 0 {
			return g, false
		}
		if len(rest) == len(g.Columns) {
			continue // the revoke names no column of this grant
		}
		g.Columns = rest
		return g, true
	}
	return g, true
}

// subtractGrantColumns returns cols without the columns in a.
func subtractGrantColumns(cols, a []string) []string {
	drop := make(map[string]bool, len(a))
	for _, c := range a {
		drop[c] = true
	}
	rest := make([]string, 0, len(cols))
	for _, c := range cols {
		if !drop[c] {
			rest = append(rest, c)
		}
	}
	return rest
}

// sameGrantColumns reports whether two column lists name the same columns.
func sameGrantColumns(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	x := append([]string(nil), a...)
	y := append([]string(nil), b...)
	sort.Strings(x)
	sort.Strings(y)
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
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

func (c *Catalog) extractGrantees(granteeList *nodes.List) []grantee {
	if granteeList == nil {
		return nil
	}
	var result []grantee
	for _, item := range granteeList.Items {
		rs, ok := item.(*nodes.RoleSpec)
		if !ok {
			continue
		}
		result = append(result, newGrantee(rs))
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
	case nodes.OBJECT_TEMPLATE:
		// A COPY template's privileges are ALTER, DROP and USAGE; its grants are
		// recorded nowhere (see ExecGrantStmt).
		return 'c'
	case nodes.OBJECT_FOREIGN_TABLE:
		// An external table is a relation in the loaded metadata, so its grants
		// are recorded against it; redshiftOnlyPrivileges holds the privileges
		// r_GRANT and r_REVOKE list for it.
		return 'r'
	case nodes.OBJECT_EXTERNAL_SCHEMA:
		// An external schema is a schema in the loaded metadata, as above.
		return 'n'
	case nodes.OBJECT_DATASHARE, nodes.OBJECT_COPY_JOB, nodes.OBJECT_MODEL:
		// Redshift's own object types: the catalog resolves none of them, so
		// PostgreSQL's ACL validation has no privilege set to hold them to.
		return 'x'
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

// redshiftOnlyPrivileges is the privilege set r_GRANT and r_REVOKE list for a
// Redshift object the catalog resolves as a PostgreSQL one: an external table
// is a relation and an external schema a schema.
var redshiftOnlyPrivileges = map[nodes.ObjectType]map[string]bool{
	nodes.OBJECT_FOREIGN_TABLE: {
		"SELECT": true, "ALTER": true, "DROP": true, "DELETE": true, "INSERT": true,
	},
	nodes.OBJECT_EXTERNAL_SCHEMA: {
		"CREATE": true, "ALTER": true, "DROP": true,
	},
}

// validGrantPrivilege holds a privilege to the object type's own set:
// Redshift's for an external table or schema, PostgreSQL's otherwise.
func validGrantPrivilege(objtype nodes.ObjectType, priv string, objTypeCode byte) bool {
	if set, ok := redshiftOnlyPrivileges[objtype]; ok {
		p := strings.ToUpper(priv)
		return p == "ALL" || set[p]
	}
	return isValidPrivilege(priv, objTypeCode)
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
	case 'c': // COPY template
		switch p {
		case "ALTER", "DROP", "USAGE":
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
	case 'c':
		return "template"
	default:
		return "object"
	}
}
