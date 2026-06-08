package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit (structural-gate) tests for the parser-ddl-spanner node: CREATE/ALTER/DROP
// of CHANGE STREAM, SEQUENCE, ROLE, LOCALITY GROUP, PROTO BUNDLE, plus the Spanner
// role-based GRANT/REVOKE forms. Accept/reject parity with the live Spanner
// emulator is in spanner_ddl_oracle_test.go.

// --- CHANGE STREAM ---

func changeStreamCreateOf(t *testing.T, sql string) *ast.CreateChangeStreamStmt {
	t.Helper()
	n := parseDDL(t, sql)
	s, ok := n.(*ast.CreateChangeStreamStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateChangeStreamStmt", sql, n)
	}
	return s
}

func TestCreateChangeStream_ForAll(t *testing.T) {
	s := changeStreamCreateOf(t, "CREATE CHANGE STREAM MyStream FOR ALL")
	if s.Name.String() != "MyStream" {
		t.Errorf("Name = %q, want MyStream", s.Name.String())
	}
	if !s.HasFor || !s.ForAll {
		t.Errorf("HasFor=%v ForAll=%v, want both true", s.HasFor, s.ForAll)
	}
	if len(s.ForTables) != 0 {
		t.Errorf("ForTables = %d, want 0", len(s.ForTables))
	}
}

func TestCreateChangeStream_ForTables(t *testing.T) {
	s := changeStreamCreateOf(t, "CREATE CHANGE STREAM s FOR Singers, Albums(Title, Budget), Tracks()")
	if !s.HasFor || s.ForAll {
		t.Fatalf("HasFor=%v ForAll=%v, want HasFor true / ForAll false", s.HasFor, s.ForAll)
	}
	if len(s.ForTables) != 3 {
		t.Fatalf("ForTables = %d, want 3", len(s.ForTables))
	}
	// Singers: whole table, no parens.
	if s.ForTables[0].Name.String() != "Singers" || s.ForTables[0].ExplicitColumns {
		t.Errorf("ForTables[0] = %+v, want Singers w/o explicit cols", s.ForTables[0])
	}
	// Albums(Title, Budget): explicit columns.
	if !s.ForTables[1].ExplicitColumns || len(s.ForTables[1].Columns) != 2 ||
		s.ForTables[1].Columns[0] != "Title" || s.ForTables[1].Columns[1] != "Budget" {
		t.Errorf("ForTables[1] = %+v, want Albums(Title, Budget)", s.ForTables[1])
	}
	// Tracks(): explicit, empty.
	if !s.ForTables[2].ExplicitColumns || len(s.ForTables[2].Columns) != 0 {
		t.Errorf("ForTables[2] = %+v, want Tracks() empty explicit", s.ForTables[2])
	}
}

func TestCreateChangeStream_Options(t *testing.T) {
	s := changeStreamCreateOf(t, "CREATE CHANGE STREAM s FOR ALL OPTIONS (retention_period = '7d', value_capture_type = 'NEW_ROW')")
	if len(s.Options) != 2 {
		t.Fatalf("Options = %d, want 2", len(s.Options))
	}
	if s.Options[0].Name != "retention_period" {
		t.Errorf("Options[0].Name = %q, want retention_period", s.Options[0].Name)
	}
}

func TestCreateChangeStream_OptionsOnly(t *testing.T) {
	// FOR clause omitted: HasFor must be false. (Emulator-accepted.)
	s := changeStreamCreateOf(t, "CREATE CHANGE STREAM s OPTIONS (retention_period = '7d')")
	if s.HasFor {
		t.Errorf("HasFor = true, want false (no FOR clause)")
	}
	if len(s.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(s.Options))
	}
}

func TestAlterChangeStream_SetForAll(t *testing.T) {
	n := parseDDL(t, "ALTER CHANGE STREAM s SET FOR ALL")
	a, ok := n.(*ast.AlterChangeStreamStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.AlterChangeStreamStmt", n)
	}
	if a.Action != ast.ChangeStreamSetFor || !a.ForAll {
		t.Errorf("Action=%v ForAll=%v, want SetFor/true", a.Action, a.ForAll)
	}
}

func TestAlterChangeStream_SetForTables(t *testing.T) {
	n := parseDDL(t, "ALTER CHANGE STREAM s SET FOR Singers, Albums")
	a := n.(*ast.AlterChangeStreamStmt)
	if a.Action != ast.ChangeStreamSetFor || a.ForAll || len(a.ForTables) != 2 {
		t.Errorf("got Action=%v ForAll=%v tables=%d, want SetFor/false/2", a.Action, a.ForAll, len(a.ForTables))
	}
}

func TestAlterChangeStream_DropForAll(t *testing.T) {
	n := parseDDL(t, "ALTER CHANGE STREAM s DROP FOR ALL")
	a := n.(*ast.AlterChangeStreamStmt)
	if a.Action != ast.ChangeStreamDropForAll {
		t.Errorf("Action = %v, want DropForAll", a.Action)
	}
}

func TestAlterChangeStream_SetOptions(t *testing.T) {
	n := parseDDL(t, "ALTER CHANGE STREAM s SET OPTIONS (retention_period = '14d')")
	a := n.(*ast.AlterChangeStreamStmt)
	if a.Action != ast.ChangeStreamSetOptions || len(a.Options) != 1 {
		t.Errorf("got Action=%v opts=%d, want SetOptions/1", a.Action, len(a.Options))
	}
}

func TestDropChangeStream(t *testing.T) {
	n := parseDDL(t, "DROP CHANGE STREAM MyStream")
	d, ok := n.(*ast.DropChangeStreamStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.DropChangeStreamStmt", n)
	}
	if d.Name.String() != "MyStream" {
		t.Errorf("Name = %q, want MyStream", d.Name.String())
	}
}

// --- SEQUENCE ---

func TestCreateSequence(t *testing.T) {
	n := parseDDL(t, "CREATE SEQUENCE MySeq OPTIONS (sequence_kind = 'bit_reversed_positive')")
	s, ok := n.(*ast.CreateSequenceStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.CreateSequenceStmt", n)
	}
	if s.Name.String() != "MySeq" || s.IfNotExists || len(s.Options) != 1 {
		t.Errorf("got Name=%q ifne=%v opts=%d, want MySeq/false/1", s.Name.String(), s.IfNotExists, len(s.Options))
	}
}

func TestCreateSequence_IfNotExistsNoOptions(t *testing.T) {
	n := parseDDL(t, "CREATE SEQUENCE IF NOT EXISTS s")
	s := n.(*ast.CreateSequenceStmt)
	if !s.IfNotExists || len(s.Options) != 0 {
		t.Errorf("got ifne=%v opts=%d, want true/0", s.IfNotExists, len(s.Options))
	}
}

func TestCreateSequence_Schema(t *testing.T) {
	n := parseDDL(t, "CREATE SEQUENCE myschema.MySeq")
	s := n.(*ast.CreateSequenceStmt)
	if s.Name.String() != "myschema.MySeq" {
		t.Errorf("Name = %q, want myschema.MySeq", s.Name.String())
	}
}

func TestAlterSequence(t *testing.T) {
	n := parseDDL(t, "ALTER SEQUENCE MySeq SET OPTIONS (start_with_counter = 9000)")
	a, ok := n.(*ast.AlterSequenceStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.AlterSequenceStmt", n)
	}
	if a.IfExists || len(a.SetOptions) != 1 {
		t.Errorf("got ife=%v opts=%d, want false/1", a.IfExists, len(a.SetOptions))
	}
}

func TestAlterSequence_IfExists(t *testing.T) {
	n := parseDDL(t, "ALTER SEQUENCE IF EXISTS s SET OPTIONS (x = 1)")
	a := n.(*ast.AlterSequenceStmt)
	if !a.IfExists {
		t.Errorf("IfExists = false, want true")
	}
}

func TestDropSequence(t *testing.T) {
	n := parseDDL(t, "DROP SEQUENCE MySeq")
	d, ok := n.(*ast.DropSequenceStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.DropSequenceStmt", n)
	}
	if d.IfExists {
		t.Errorf("IfExists = true, want false")
	}
}

func TestDropSequence_IfExists(t *testing.T) {
	n := parseDDL(t, "DROP SEQUENCE IF EXISTS MySeq")
	d := n.(*ast.DropSequenceStmt)
	if !d.IfExists {
		t.Errorf("IfExists = false, want true")
	}
}

// --- ROLE ---

func TestCreateRole(t *testing.T) {
	n := parseDDL(t, "CREATE ROLE analyst")
	r, ok := n.(*ast.CreateRoleStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.CreateRoleStmt", n)
	}
	if r.Name.String() != "analyst" {
		t.Errorf("Name = %q, want analyst", r.Name.String())
	}
}

func TestCreateRole_BacktickName(t *testing.T) {
	// A backtick-quoted name is a SINGLE identifier token (its body may contain a
	// dot); accepted, unlike an unquoted dotted name.
	n := parseDDL(t, "CREATE ROLE `dotted.name`")
	r := n.(*ast.CreateRoleStmt)
	if r.Name.String() != "dotted.name" || len(r.Name.Parts) != 1 {
		t.Errorf("Name = %+v, want single part 'dotted.name'", r.Name.Parts)
	}
}

func TestDropRole(t *testing.T) {
	n := parseDDL(t, "DROP ROLE analyst")
	r, ok := n.(*ast.DropRoleStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.DropRoleStmt", n)
	}
	if r.Name.String() != "analyst" {
		t.Errorf("Name = %q, want analyst", r.Name.String())
	}
}

func TestDropRole_DottedPathAccepted(t *testing.T) {
	// DROP ROLE accepts a dotted path (emulator-verified), unlike CREATE ROLE.
	n := parseDDL(t, "DROP ROLE a.b.c")
	r := n.(*ast.DropRoleStmt)
	if r.Name.String() != "a.b.c" {
		t.Errorf("Name = %q, want a.b.c", r.Name.String())
	}
}

// --- LOCALITY GROUP ---

func TestCreateLocalityGroup(t *testing.T) {
	n := parseDDL(t, "CREATE LOCALITY GROUP hot_data OPTIONS (storage = 'ssd')")
	lg, ok := n.(*ast.CreateLocalityGroupStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.CreateLocalityGroupStmt", n)
	}
	if lg.Name.String() != "hot_data" || len(lg.Options) != 1 {
		t.Errorf("got Name=%q opts=%d, want hot_data/1", lg.Name.String(), len(lg.Options))
	}
}

func TestCreateLocalityGroup_MultiOption(t *testing.T) {
	n := parseDDL(t, "CREATE LOCALITY GROUP cold_data OPTIONS (storage = 'hdd', ssd_to_hdd_spill_timespan = '30d')")
	lg := n.(*ast.CreateLocalityGroupStmt)
	if len(lg.Options) != 2 {
		t.Errorf("Options = %d, want 2", len(lg.Options))
	}
}

func TestAlterLocalityGroup(t *testing.T) {
	n := parseDDL(t, "ALTER LOCALITY GROUP cold_data SET OPTIONS (ssd_to_hdd_spill_timespan = '14d')")
	a, ok := n.(*ast.AlterLocalityGroupStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.AlterLocalityGroupStmt", n)
	}
	if len(a.SetOptions) != 1 {
		t.Errorf("SetOptions = %d, want 1", len(a.SetOptions))
	}
}

func TestAlterLocalityGroup_NoSetOptions(t *testing.T) {
	// The SET OPTIONS clause is optional for ALTER LOCALITY GROUP (emulator accepts
	// a bare form — unlike ALTER SEQUENCE / ALTER CHANGE STREAM).
	n := parseDDL(t, "ALTER LOCALITY GROUP g")
	a := n.(*ast.AlterLocalityGroupStmt)
	if len(a.SetOptions) != 0 {
		t.Errorf("SetOptions = %d, want 0", len(a.SetOptions))
	}
}

func TestDropLocalityGroup(t *testing.T) {
	n := parseDDL(t, "DROP LOCALITY GROUP cold_data")
	d, ok := n.(*ast.DropLocalityGroupStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.DropLocalityGroupStmt", n)
	}
	if d.Name.String() != "cold_data" {
		t.Errorf("Name = %q, want cold_data", d.Name.String())
	}
}

// --- PROTO BUNDLE ---

func TestCreateProtoBundle(t *testing.T) {
	n := parseDDL(t, "CREATE PROTO BUNDLE (`my.package.MyMessage`, `my.package.AnotherMessage`)")
	b, ok := n.(*ast.CreateProtoBundleStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.CreateProtoBundleStmt", n)
	}
	if len(b.Types) != 2 {
		t.Fatalf("Types = %d, want 2", len(b.Types))
	}
	if b.Types[0].String() != "my.package.MyMessage" {
		t.Errorf("Types[0] = %q, want my.package.MyMessage", b.Types[0].String())
	}
}

func TestAlterProtoBundle_Insert(t *testing.T) {
	n := parseDDL(t, "ALTER PROTO BUNDLE INSERT (`my.package.NewMessage`)")
	a, ok := n.(*ast.AlterProtoBundleStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.AlterProtoBundleStmt", n)
	}
	if len(a.Insert) != 1 || len(a.Update) != 0 || len(a.Delete) != 0 {
		t.Errorf("got insert=%d update=%d delete=%d, want 1/0/0", len(a.Insert), len(a.Update), len(a.Delete))
	}
}

func TestAlterProtoBundle_Combo(t *testing.T) {
	// All three groups, fixed order INSERT→UPDATE→DELETE, space-separated.
	n := parseDDL(t, "ALTER PROTO BUNDLE INSERT (`a.B`) UPDATE (`c.D`) DELETE (`e.F`)")
	a := n.(*ast.AlterProtoBundleStmt)
	if len(a.Insert) != 1 || len(a.Update) != 1 || len(a.Delete) != 1 {
		t.Errorf("got insert=%d update=%d delete=%d, want 1/1/1", len(a.Insert), len(a.Update), len(a.Delete))
	}
}

func TestAlterProtoBundle_DeleteMultiType(t *testing.T) {
	n := parseDDL(t, "ALTER PROTO BUNDLE DELETE (`a.B`, `c.D`)")
	a := n.(*ast.AlterProtoBundleStmt)
	if len(a.Delete) != 2 || len(a.Insert) != 0 {
		t.Errorf("got insert=%d delete=%d, want 0/2", len(a.Insert), len(a.Delete))
	}
}

// TestProtoBundle_TrailingComma guards the emulator-verified leniency: a PROTO
// BUNDLE type list ALLOWS a trailing comma (unlike a change-stream column list).
func TestProtoBundle_TrailingComma(t *testing.T) {
	b := parseDDL(t, "CREATE PROTO BUNDLE (`a.b.C`, `a.b.D`,)").(*ast.CreateProtoBundleStmt)
	if len(b.Types) != 2 {
		t.Errorf("CREATE Types = %d, want 2 (trailing comma allowed)", len(b.Types))
	}
	a := parseDDL(t, "ALTER PROTO BUNDLE INSERT (`a.b.C`,)").(*ast.AlterProtoBundleStmt)
	if len(a.Insert) != 1 {
		t.Errorf("ALTER Insert = %d, want 1 (trailing comma allowed)", len(a.Insert))
	}
	// An empty list (just `()`) still rejects (trailing comma needs a prior entry).
	assertReject(t, "CREATE PROTO BUNDLE ()")
	assertReject(t, "CREATE PROTO BUNDLE (,)")
}

// TestAlterProtoBundle_OrderingRejects guards the strict grammar: groups must be
// space-separated (no commas), at most once each, in INSERT→UPDATE→DELETE order.
func TestAlterProtoBundle_OrderingRejects(t *testing.T) {
	for _, sql := range []string{
		"ALTER PROTO BUNDLE INSERT (`a.B`), UPDATE (`c.D`)", // comma between groups
		"ALTER PROTO BUNDLE INSERT (`a.B`) INSERT (`c.D`)",  // repeated group
		"ALTER PROTO BUNDLE UPDATE (`a.B`) INSERT (`c.D`)",  // out of order
	} {
		assertReject(t, sql)
	}
}

// --- role-based GRANT / REVOKE (Spanner) ---

func TestGrant_ToRole(t *testing.T) {
	n := parseDDL(t, "GRANT SELECT ON TABLE Singers TO ROLE analyst")
	g, ok := n.(*ast.GrantStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.GrantStmt", n)
	}
	if len(g.Privileges) != 1 || g.Privileges[0].Name != "SELECT" {
		t.Errorf("Privileges = %+v, want [SELECT]", g.Privileges)
	}
	if len(g.ObjectType) != 1 || g.ObjectType[0] != "TABLE" {
		t.Errorf("ObjectType = %v, want [TABLE]", g.ObjectType)
	}
	if g.Path.String() != "Singers" {
		t.Errorf("Path = %q, want Singers", g.Path.String())
	}
	if len(g.Grantees) != 1 || g.Grantees[0].Kind != ast.GranteeRole || g.Grantees[0].Value != "analyst" {
		t.Errorf("Grantees = %+v, want [ROLE analyst]", g.Grantees)
	}
}

func TestGrant_ToRoleList(t *testing.T) {
	n := parseDDL(t, "GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE Albums TO ROLE editor, admin")
	g := n.(*ast.GrantStmt)
	if len(g.Privileges) != 4 {
		t.Errorf("Privileges = %d, want 4", len(g.Privileges))
	}
	if len(g.Grantees) != 2 || g.Grantees[0].Value != "editor" || g.Grantees[1].Value != "admin" {
		t.Errorf("Grantees = %+v, want [editor admin]", g.Grantees)
	}
	for i, gr := range g.Grantees {
		if gr.Kind != ast.GranteeRole {
			t.Errorf("Grantees[%d].Kind = %v, want ROLE", i, gr.Kind)
		}
	}
}

func TestGrant_ColumnLevelToRole(t *testing.T) {
	// Column-level privilege (DDL-034). The emulator rejects this ("does not yet
	// support column level access controls"), but it is valid GoogleSQL/Spanner
	// per the docs — the union parser accepts it (divergence #flagged).
	n := parseDDL(t, "GRANT SELECT(SingerId, FirstName) ON TABLE Singers TO ROLE read_only")
	g := n.(*ast.GrantStmt)
	if len(g.Privileges) != 1 || len(g.Privileges[0].Columns) != 2 {
		t.Errorf("Privileges = %+v, want SELECT with 2 columns", g.Privileges)
	}
	if g.Grantees[0].Kind != ast.GranteeRole {
		t.Errorf("Grantee kind = %v, want ROLE", g.Grantees[0].Kind)
	}
}

func TestGrant_ExecuteTableFunctionToRole(t *testing.T) {
	n := parseDDL(t, "GRANT EXECUTE ON TABLE FUNCTION get_singers TO ROLE analyst")
	g := n.(*ast.GrantStmt)
	if len(g.ObjectType) != 2 || g.ObjectType[0] != "TABLE" || g.ObjectType[1] != "FUNCTION" {
		t.Errorf("ObjectType = %v, want [TABLE FUNCTION]", g.ObjectType)
	}
}

func TestGrant_RoleToRole(t *testing.T) {
	n := parseDDL(t, "GRANT ROLE analyst TO ROLE senior_analyst")
	g, ok := n.(*ast.GrantStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.GrantStmt", n)
	}
	if len(g.Roles) != 1 || g.Roles[0].String() != "analyst" {
		t.Errorf("Roles = %+v, want [analyst]", g.Roles)
	}
	if len(g.Grantees) != 1 || g.Grantees[0].Kind != ast.GranteeRole || g.Grantees[0].Value != "senior_analyst" {
		t.Errorf("Grantees = %+v, want [ROLE senior_analyst]", g.Grantees)
	}
	if g.Privileges != nil || g.AllPrivileges {
		t.Errorf("role-grant must carry no privileges (privs=%v all=%v)", g.Privileges, g.AllPrivileges)
	}
}

func TestGrant_RoleListToRoleList(t *testing.T) {
	n := parseDDL(t, "GRANT ROLE a, b TO ROLE c, d")
	g := n.(*ast.GrantStmt)
	if len(g.Roles) != 2 || g.Roles[0].String() != "a" || g.Roles[1].String() != "b" {
		t.Errorf("Roles = %+v, want [a b]", g.Roles)
	}
	if len(g.Grantees) != 2 {
		t.Errorf("Grantees = %d, want 2", len(g.Grantees))
	}
}

func TestRevoke_FromRole(t *testing.T) {
	n := parseDDL(t, "REVOKE SELECT ON TABLE Singers FROM ROLE analyst")
	r, ok := n.(*ast.RevokeStmt)
	if !ok {
		t.Fatalf("statement is %T, want *ast.RevokeStmt", n)
	}
	if len(r.Grantees) != 1 || r.Grantees[0].Kind != ast.GranteeRole || r.Grantees[0].Value != "analyst" {
		t.Errorf("Grantees = %+v, want [ROLE analyst]", r.Grantees)
	}
}

func TestRevoke_RoleFromRole(t *testing.T) {
	n := parseDDL(t, "REVOKE ROLE analyst FROM ROLE senior_analyst")
	r := n.(*ast.RevokeStmt)
	if len(r.Roles) != 1 || r.Roles[0].String() != "analyst" {
		t.Errorf("Roles = %+v, want [analyst]", r.Roles)
	}
	if len(r.Grantees) != 1 || r.Grantees[0].Kind != ast.GranteeRole {
		t.Errorf("Grantees = %+v, want [ROLE senior_analyst]", r.Grantees)
	}
}

// Regression: the legacy ZetaSQL string-grantee GRANT must STILL parse (the
// union parser serves both dialects). This guards the parser-dcl behavior.
func TestGrant_LegacyStringGranteeStillWorks(t *testing.T) {
	n := parseDDL(t, "GRANT SELECT ON TABLE foo TO 'user@google.com'")
	g := n.(*ast.GrantStmt)
	if len(g.Grantees) != 1 || g.Grantees[0].Kind != ast.GranteeString || g.Grantees[0].Value != "user@google.com" {
		t.Errorf("Grantees = %+v, want [STRING user@google.com]", g.Grantees)
	}
}

// TestGrant_PrivilegeNamedRole guards the role-grant disambiguation: `GRANT ROLE
// ON …` is a LEGACY privilege grant whose privilege is named "ROLE"
// (privilege_name: identifier), NOT a Spanner role-to-role grant. The union must
// accept it (the emulator rejects it as a Spanner form — non-authoritative there).
// Discovered by the cross-model (Codex) review.
func TestGrant_PrivilegeNamedRole(t *testing.T) {
	n := parseDDL(t, "GRANT ROLE ON TABLE foo TO 'x'")
	g := n.(*ast.GrantStmt)
	if len(g.Roles) != 0 {
		t.Errorf("Roles = %v, want empty (this is a privilege grant, not role-to-role)", g.Roles)
	}
	if len(g.Privileges) != 1 || g.Privileges[0].Name != "ROLE" {
		t.Errorf("Privileges = %+v, want [ROLE]", g.Privileges)
	}
	if len(g.ObjectType) != 1 || g.ObjectType[0] != "TABLE" {
		t.Errorf("ObjectType = %v, want [TABLE]", g.ObjectType)
	}
	// `GRANT ROLE, SELECT ON …` (ROLE as the first privilege in a list).
	g2 := parseDDL(t, "GRANT ROLE, SELECT ON TABLE foo TO 'x'").(*ast.GrantStmt)
	if len(g2.Privileges) != 2 || g2.Privileges[0].Name != "ROLE" {
		t.Errorf("Privileges = %+v, want [ROLE SELECT]", g2.Privileges)
	}
}

// TestGrant_CommaSeparatedObjects guards the Spanner comma-separated object list
// (`ON TABLE t1, t2 TO ROLE r`). Discovered by the cross-model (Codex) review.
func TestGrant_CommaSeparatedObjects(t *testing.T) {
	g := parseDDL(t, "GRANT SELECT ON TABLE t1, t2, t3 TO ROLE r").(*ast.GrantStmt)
	if len(g.Paths) != 3 {
		t.Fatalf("Paths = %d, want 3", len(g.Paths))
	}
	if g.Paths[0].String() != "t1" || g.Paths[2].String() != "t3" {
		t.Errorf("Paths = %v, want [t1 t2 t3]", g.Paths)
	}
	if g.Path.String() != "t1" {
		t.Errorf("Path = %q, want t1 (== Paths[0])", g.Path.String())
	}
	if len(g.ObjectType) != 1 || g.ObjectType[0] != "TABLE" {
		t.Errorf("ObjectType = %v, want [TABLE]", g.ObjectType)
	}
	// REVOKE comma objects.
	r := parseDDL(t, "REVOKE SELECT ON TABLE t1, t2 FROM ROLE r").(*ast.RevokeStmt)
	if len(r.Paths) != 2 {
		t.Errorf("revoke Paths = %d, want 2", len(r.Paths))
	}
	// VIEW object type with a comma list.
	g3 := parseDDL(t, "GRANT SELECT ON VIEW v1, v2 TO ROLE r").(*ast.GrantStmt)
	if len(g3.ObjectType) != 1 || g3.ObjectType[0] != "VIEW" || len(g3.Paths) != 2 {
		t.Errorf("ObjectType=%v Paths=%v, want [VIEW] + 2 paths", g3.ObjectType, g3.Paths)
	}
}

// TestGrant_CommaObjectsRequireType guards that a comma-separated object list is
// the Spanner form, which REQUIRES a leading object-type keyword: `ON foo, bar`
// (no type) is invalid in BOTH dialects (Spanner needs target_type; legacy has no
// multi-object). A single object with no type stays valid (legacy). Surfaced while
// verifying the Codex comma-object finding against the oracle.
func TestGrant_CommaObjectsRequireType(t *testing.T) {
	// No type + comma list rejects.
	assertReject(t, "GRANT SELECT ON foo, bar TO ROLE r")
	assertReject(t, "REVOKE SELECT ON foo, bar FROM ROLE r")
	// Single object with no type still parses (legacy form).
	g := parseDDL(t, "GRANT SELECT ON foo TO 'x'").(*ast.GrantStmt)
	if len(g.ObjectType) != 0 || len(g.Paths) != 1 || g.Paths[0].String() != "foo" {
		t.Errorf("ObjectType=%v Paths=%v, want no type + [foo]", g.ObjectType, g.Paths)
	}
}

// TestGrant_RoleTargetStrict guards that a role-to-role grant/revoke REQUIRES a
// ROLE target (the emulator rejects a string/parameter target). Discovered by the
// cross-model (Codex) review.
func TestGrant_RoleTargetStrict(t *testing.T) {
	for _, sql := range []string{
		"GRANT ROLE analyst TO 'user'",       // string target
		"GRANT ROLE analyst TO @p",           // parameter target
		"REVOKE ROLE analyst FROM 'user'",    // string target (revoke)
		"GRANT ROLE a, b TO ROLE c, 'x'",     // mixed role + string target
	} {
		assertReject(t, sql)
	}
}

// --- owned divergence regressions ---

// TestInterleave_OnDeleteOptional is the permanent regression for divergence
// #112: a Spanner `INTERLEAVE IN PARENT p` with NO `ON DELETE` action must parse
// (defaulting to NO ACTION), matching the live emulator. The original parser-ddl
// port wrongly REQUIRED `ON DELETE`, over-rejecting the documented default form
// (DDL-007). A dangling `ON DELETE` (keyword present, no action) still rejects.
func TestInterleave_OnDeleteOptional(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE child (id INT64) PRIMARY KEY (id), INTERLEAVE IN PARENT parent")
	if ct.Interleave == nil || ct.Interleave.Parent.String() != "parent" {
		t.Fatalf("Interleave = %+v, want parent", ct.Interleave)
	}
	if ct.Interleave.OnDelete != ast.FKActionNone {
		t.Errorf("OnDelete = %v, want FKActionNone (no action specified)", ct.Interleave.OnDelete)
	}
	// Explicit actions still parse.
	ct2 := createTableOf(t, "CREATE TABLE c (id INT64) PRIMARY KEY (id), INTERLEAVE IN PARENT p ON DELETE CASCADE")
	if ct2.Interleave.OnDelete != ast.FKActionCascade {
		t.Errorf("OnDelete = %v, want FKActionCascade", ct2.Interleave.OnDelete)
	}
	// A dangling `ON DELETE` with no action is still a syntax error.
	assertReject(t, "CREATE TABLE c (id INT64) PRIMARY KEY (id), INTERLEAVE IN PARENT p ON DELETE")
}

// TestInlinePrimaryKey_Accepted is the regression for divergence #8: Spanner
// accepts an inline column `PRIMARY KEY` (`id INT64 PRIMARY KEY`). The union
// parser already accepts it (the BigQuery inline-PK production covers it); this
// guards that the inline-PK form keeps parsing and matches the emulator.
func TestInlinePrimaryKey_Accepted(t *testing.T) {
	// Just assert it parses without error (shape covered by parser-ddl tests).
	parseDDL(t, "CREATE TABLE tbl_inline (id INT64 PRIMARY KEY, name STRING(MAX))")
}

// --- negative tests (over-permissiveness guards) ---

func TestSpannerDDL_Rejects(t *testing.T) {
	for _, sql := range []string{
		// CHANGE STREAM: a dangling FOR, a bad keyword, missing name.
		"CREATE CHANGE STREAM",                  // missing name (+ at least nothing)
		"CREATE CHANGE STREAM s FOR",            // FOR with nothing after
		"CREATE CHANGE STREAM s FOR ALL ALL",    // trailing junk
		"CREATE CHANGE STREAM s FOR Singers(a,)", // trailing comma in column list (NOT allowed, unlike proto)
		"CREATE CHANGE STREAM s FOR Singers,",   // trailing comma in table list
		"ALTER CHANGE STREAM s SET",             // SET with nothing
		"ALTER CHANGE STREAM s DROP FOR",        // DROP FOR without ALL
		"ALTER CHANGE STREAM s SET FOR",         // SET FOR with nothing
		"DROP CHANGE STREAM",                    // missing name
		// SEQUENCE
		"ALTER SEQUENCE s",                      // missing SET OPTIONS
		"ALTER SEQUENCE s SET",                  // SET with nothing
		"CREATE SEQUENCE",                       // missing name
		// ROLE
		"CREATE ROLE",                           // missing name
		"DROP ROLE",                             // missing name
		// LOCALITY GROUP
		"CREATE LOCALITY GROUP",                 // missing name
		"ALTER LOCALITY GROUP g SET",            // SET with no OPTIONS
		"ALTER LOCALITY GROUP g RENAME TO h",    // not a valid locality-group action
		"DROP LOCALITY GROUP",                   // missing name
		// PROTO BUNDLE
		"CREATE PROTO BUNDLE",                   // missing ( )
		"CREATE PROTO BUNDLE ()",                // empty bundle (grammar requires >=1)
		"ALTER PROTO BUNDLE",                    // missing an action
		"ALTER PROTO BUNDLE INSERT",             // INSERT with no ( )
		// role GRANT/REVOKE
		"GRANT ROLE a TO ROLE",                  // TO ROLE with no role
		"GRANT SELECT ON TABLE t TO ROLE",       // TO ROLE with no role
		"GRANT ROLE TO ROLE r",                  // GRANT ROLE with no role
		"REVOKE ROLE a FROM ROLE",               // FROM ROLE with no role
		// role names are SINGLE identifiers (emulator rejects dotted role names).
		"CREATE ROLE a.b",                       // dotted CREATE ROLE name
		"GRANT SELECT ON TABLE t TO ROLE a.b",   // dotted grantee role
		"GRANT ROLE a.b TO ROLE c",              // dotted subject role
		"GRANT SELECT ON TABLE t TO ROLE r1, r2.x", // dotted role in list
		"REVOKE ROLE r1, r2.x FROM ROLE s",      // dotted role in revoke subject
	} {
		assertReject(t, sql)
	}
}
