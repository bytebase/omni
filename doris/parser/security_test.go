package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseOne(t *testing.T, sql string) ast.Node {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	return file.Stmts[0]
}

func parseCreateEncryptKeyStmt(t *testing.T, sql string) *ast.CreateEncryptKeyStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.CreateEncryptKeyStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateEncryptKeyStmt, got %T", n)
	}
	return stmt
}

func parseDropEncryptKeyStmt(t *testing.T, sql string) *ast.DropEncryptKeyStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.DropEncryptKeyStmt)
	if !ok {
		t.Fatalf("expected *ast.DropEncryptKeyStmt, got %T", n)
	}
	return stmt
}

func parseCreateRoleStmt(t *testing.T, sql string) *ast.CreateRoleStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.CreateRoleStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateRoleStmt, got %T", n)
	}
	return stmt
}

func parseDropRoleStmt(t *testing.T, sql string) *ast.DropRoleStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.DropRoleStmt)
	if !ok {
		t.Fatalf("expected *ast.DropRoleStmt, got %T", n)
	}
	return stmt
}

func parseCreateUserStmt(t *testing.T, sql string) *ast.CreateUserStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.CreateUserStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateUserStmt, got %T", n)
	}
	return stmt
}

func parseAlterUserStmt(t *testing.T, sql string) *ast.AlterUserStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.AlterUserStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterUserStmt, got %T", n)
	}
	return stmt
}

func parseDropUserStmt(t *testing.T, sql string) *ast.DropUserStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.DropUserStmt)
	if !ok {
		t.Fatalf("expected *ast.DropUserStmt, got %T", n)
	}
	return stmt
}

func parseSetPasswordStmt(t *testing.T, sql string) *ast.SetPasswordStmt {
	t.Helper()
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.SetPasswordStmt)
	if !ok {
		t.Fatalf("expected *ast.SetPasswordStmt, got %T", n)
	}
	return stmt
}

// ---------------------------------------------------------------------------
// ENCRYPTION KEY — legacy corpus: security_encryptkey.sql
// ---------------------------------------------------------------------------

func TestCreateEncryptKey_Simple(t *testing.T) {
	// CREATE ENCRYPTKEY my_key AS "ABCD123456789";
	stmt := parseCreateEncryptKeyStmt(t, `CREATE ENCRYPTKEY my_key AS "ABCD123456789"`)
	if stmt.Name.String() != "my_key" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_key")
	}
	if stmt.Key != "ABCD123456789" {
		t.Errorf("Key = %q, want %q", stmt.Key, "ABCD123456789")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
}

func TestCreateEncryptKey_Qualified(t *testing.T) {
	// CREATE ENCRYPTKEY testdb.test_key AS "ABCD123456789";
	stmt := parseCreateEncryptKeyStmt(t, `CREATE ENCRYPTKEY testdb.test_key AS "ABCD123456789"`)
	if stmt.Name.String() != "testdb.test_key" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "testdb.test_key")
	}
	if stmt.Key != "ABCD123456789" {
		t.Errorf("Key = %q, want %q", stmt.Key, "ABCD123456789")
	}
}

func TestDropEncryptKey_Simple(t *testing.T) {
	// DROP ENCRYPTKEY my_key;
	stmt := parseDropEncryptKeyStmt(t, `DROP ENCRYPTKEY my_key`)
	if stmt.Name.String() != "my_key" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_key")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropEncryptKey_IfExists(t *testing.T) {
	// DROP ENCRYPTKEY IF EXISTS testdb.my_key
	stmt := parseDropEncryptKeyStmt(t, `DROP ENCRYPTKEY IF EXISTS testdb.my_key`)
	if stmt.Name.String() != "testdb.my_key" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "testdb.my_key")
	}
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

// ---------------------------------------------------------------------------
// ROLE — legacy corpus: account_role.sql
// ---------------------------------------------------------------------------

func TestCreateRole_Simple(t *testing.T) {
	// CREATE ROLE role1;
	stmt := parseCreateRoleStmt(t, `CREATE ROLE role1`)
	if stmt.Name != "role1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "role1")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if stmt.Comment != "" {
		t.Errorf("Comment = %q, want empty", stmt.Comment)
	}
}

func TestCreateRole_WithComment(t *testing.T) {
	// CREATE ROLE role2 COMMENT "this is my first role";
	stmt := parseCreateRoleStmt(t, `CREATE ROLE role2 COMMENT "this is my first role"`)
	if stmt.Name != "role2" {
		t.Errorf("Name = %q, want %q", stmt.Name, "role2")
	}
	if stmt.Comment != "this is my first role" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "this is my first role")
	}
}

func TestDropRole_Simple(t *testing.T) {
	// DROP ROLE role1;
	stmt := parseDropRoleStmt(t, `DROP ROLE role1`)
	if stmt.Name != "role1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "role1")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropRole_IfExists(t *testing.T) {
	// DROP ROLE IF EXISTS role1
	stmt := parseDropRoleStmt(t, `DROP ROLE IF EXISTS role1`)
	if stmt.Name != "role1" {
		t.Errorf("Name = %q, want %q", stmt.Name, "role1")
	}
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

// ---------------------------------------------------------------------------
// USER — CREATE USER — legacy corpus: account_create_user.sql
// ---------------------------------------------------------------------------

func TestCreateUser_Simple(t *testing.T) {
	// CREATE USER 'jack';
	stmt := parseCreateUserStmt(t, `CREATE USER 'jack'`)
	if stmt.Name.Username != "jack" {
		t.Errorf("Username = %q, want %q", stmt.Name.Username, "jack")
	}
	if stmt.Name.Host != "%" {
		t.Errorf("Host = %q, want %%", stmt.Name.Host)
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
}

func TestCreateUser_WithHost(t *testing.T) {
	// CREATE USER jack@'172.10.1.10' IDENTIFIED BY '123456';
	stmt := parseCreateUserStmt(t, `CREATE USER jack@'172.10.1.10' IDENTIFIED BY '123456'`)
	if stmt.Name.Username != "jack" {
		t.Errorf("Username = %q, want %q", stmt.Name.Username, "jack")
	}
	if stmt.Name.Host != "172.10.1.10" {
		t.Errorf("Host = %q, want %q", stmt.Name.Host, "172.10.1.10")
	}
	if stmt.Password != "123456" {
		t.Errorf("Password = %q, want %q", stmt.Password, "123456")
	}
}

func TestCreateUser_IdentifiedByPasswordHash(t *testing.T) {
	// CREATE USER jack@'172.10.1.10' IDENTIFIED BY PASSWORD '*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9';
	stmt := parseCreateUserStmt(t, `CREATE USER jack@'172.10.1.10' IDENTIFIED BY PASSWORD '*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9'`)
	if stmt.Name.Username != "jack" {
		t.Errorf("Username = %q, want %q", stmt.Name.Username, "jack")
	}
	if stmt.PasswordHash != "*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9" {
		t.Errorf("PasswordHash = %q, unexpected", stmt.PasswordHash)
	}
	if stmt.Password != "" {
		t.Errorf("Password should be empty, got %q", stmt.Password)
	}
}

func TestCreateUser_DefaultRole(t *testing.T) {
	// CREATE USER 'jack'@'192.168.%' DEFAULT ROLE 'example_role';
	stmt := parseCreateUserStmt(t, `CREATE USER 'jack'@'192.168.%' DEFAULT ROLE 'example_role'`)
	if stmt.Name.Username != "jack" {
		t.Errorf("Username = %q, want %q", stmt.Name.Username, "jack")
	}
	if stmt.Name.Host != "192.168.%" {
		t.Errorf("Host = %q, want %q", stmt.Name.Host, "192.168.%")
	}
	if stmt.DefaultRole != "example_role" {
		t.Errorf("DefaultRole = %q, want %q", stmt.DefaultRole, "example_role")
	}
}

func TestCreateUser_PasswordAndDefaultRole(t *testing.T) {
	// CREATE USER 'jack'@'%' IDENTIFIED BY '12345' DEFAULT ROLE 'my_role';
	stmt := parseCreateUserStmt(t, `CREATE USER 'jack'@'%' IDENTIFIED BY '12345' DEFAULT ROLE 'my_role'`)
	if stmt.Password != "12345" {
		t.Errorf("Password = %q, want %q", stmt.Password, "12345")
	}
	if stmt.DefaultRole != "my_role" {
		t.Errorf("DefaultRole = %q, want %q", stmt.DefaultRole, "my_role")
	}
}

func TestCreateUser_PasswordPolicy(t *testing.T) {
	// CREATE USER 'jack' IDENTIFIED BY '12345' PASSWORD_EXPIRE INTERVAL 10 DAY FAILED_LOGIN_ATTEMPTS 3 PASSWORD_LOCK_TIME 1 DAY;
	stmt := parseCreateUserStmt(t, `CREATE USER 'jack' IDENTIFIED BY '12345' PASSWORD_EXPIRE INTERVAL 10 DAY FAILED_LOGIN_ATTEMPTS 3 PASSWORD_LOCK_TIME 1 DAY`)
	if stmt.Password != "12345" {
		t.Errorf("Password = %q, want %q", stmt.Password, "12345")
	}
	if !stmt.PasswordExpire {
		t.Error("PasswordExpire should be true")
	}
	if stmt.PasswordExpireInterval != 10 {
		t.Errorf("PasswordExpireInterval = %d, want 10", stmt.PasswordExpireInterval)
	}
	if stmt.FailedLoginAttempts != 3 {
		t.Errorf("FailedLoginAttempts = %d, want 3", stmt.FailedLoginAttempts)
	}
	if stmt.PasswordLockTime != 1 {
		t.Errorf("PasswordLockTime = %d, want 1", stmt.PasswordLockTime)
	}
}

func TestCreateUser_PasswordHistory(t *testing.T) {
	// CREATE USER 'jack' IDENTIFIED BY '12345' PASSWORD_HISTORY 8;
	stmt := parseCreateUserStmt(t, `CREATE USER 'jack' IDENTIFIED BY '12345' PASSWORD_HISTORY 8`)
	if stmt.PasswordHistory != 8 {
		t.Errorf("PasswordHistory = %d, want 8", stmt.PasswordHistory)
	}
}

func TestCreateUser_Comment(t *testing.T) {
	// CREATE USER 'jack' COMMENT "this is my first user"
	stmt := parseCreateUserStmt(t, `CREATE USER 'jack' COMMENT "this is my first user"`)
	if stmt.Comment != "this is my first user" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "this is my first user")
	}
}

func TestCreateUser_DomainHost(t *testing.T) {
	// CREATE USER 'jack'@'example_domain' IDENTIFIED BY '12345';
	stmt := parseCreateUserStmt(t, `CREATE USER 'jack'@'example_domain' IDENTIFIED BY '12345'`)
	if stmt.Name.Host != "example_domain" {
		t.Errorf("Host = %q, want %q", stmt.Name.Host, "example_domain")
	}
}

// ---------------------------------------------------------------------------
// USER — ALTER USER — legacy corpus: account_alter_user.sql
// ---------------------------------------------------------------------------

func TestAlterUser_IdentifiedBy(t *testing.T) {
	// ALTER USER jack@'%' IDENTIFIED BY "12345";
	stmt := parseAlterUserStmt(t, `ALTER USER jack@'%' IDENTIFIED BY "12345"`)
	if stmt.Name.Username != "jack" {
		t.Errorf("Username = %q, want %q", stmt.Name.Username, "jack")
	}
	if stmt.Name.Host != "%" {
		t.Errorf("Host = %q, want %%", stmt.Name.Host)
	}
	if stmt.Password != "12345" {
		t.Errorf("Password = %q, want %q", stmt.Password, "12345")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestAlterUser_PasswordPolicy(t *testing.T) {
	// ALTER USER jack@'%' FAILED_LOGIN_ATTEMPTS 3 PASSWORD_LOCK_TIME 1 DAY;
	stmt := parseAlterUserStmt(t, `ALTER USER jack@'%' FAILED_LOGIN_ATTEMPTS 3 PASSWORD_LOCK_TIME 1 DAY`)
	if stmt.FailedLoginAttempts != 3 {
		t.Errorf("FailedLoginAttempts = %d, want 3", stmt.FailedLoginAttempts)
	}
	if stmt.PasswordLockTime != 1 {
		t.Errorf("PasswordLockTime = %d, want 1", stmt.PasswordLockTime)
	}
}

func TestAlterUser_AccountUnlock(t *testing.T) {
	// ALTER USER jack@'%' ACCOUNT_UNLOCK;
	stmt := parseAlterUserStmt(t, `ALTER USER jack@'%' ACCOUNT_UNLOCK`)
	if !stmt.AccountUnlock {
		t.Error("AccountUnlock should be true")
	}
	if stmt.AccountLock {
		t.Error("AccountLock should be false")
	}
}

func TestAlterUser_Comment(t *testing.T) {
	// ALTER USER jack@'%' COMMENT "this is my first user"
	stmt := parseAlterUserStmt(t, `ALTER USER jack@'%' COMMENT "this is my first user"`)
	if stmt.Comment != "this is my first user" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "this is my first user")
	}
}

// ---------------------------------------------------------------------------
// USER — DROP USER — legacy corpus: account_drop_user.sql
// ---------------------------------------------------------------------------

func TestDropUser_Simple(t *testing.T) {
	// DROP USER 'jack'@'192.%'
	stmt := parseDropUserStmt(t, `DROP USER 'jack'@'192.%'`)
	if stmt.Name.Username != "jack" {
		t.Errorf("Username = %q, want %q", stmt.Name.Username, "jack")
	}
	if stmt.Name.Host != "192.%" {
		t.Errorf("Host = %q, want %q", stmt.Name.Host, "192.%")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropUser_IfExists(t *testing.T) {
	stmt := parseDropUserStmt(t, `DROP USER IF EXISTS 'jack'@'%'`)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name.Username != "jack" {
		t.Errorf("Username = %q, want %q", stmt.Name.Username, "jack")
	}
}

// ---------------------------------------------------------------------------
// SET PASSWORD — legacy corpus: account_set_password.sql
// ---------------------------------------------------------------------------

func TestSetPassword_HashNoFor(t *testing.T) {
	// SET PASSWORD = '*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9';
	stmt := parseSetPasswordStmt(t, `SET PASSWORD = '*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9'`)
	if stmt.For != nil {
		t.Error("For should be nil")
	}
	if stmt.Password != "*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9" {
		t.Errorf("Password = %q, unexpected", stmt.Password)
	}
	if !stmt.IsHash {
		t.Error("IsHash should be true")
	}
}

func TestSetPassword_PasswordFuncNoFor(t *testing.T) {
	// SET PASSWORD = PASSWORD('123456');
	stmt := parseSetPasswordStmt(t, `SET PASSWORD = PASSWORD('123456')`)
	if stmt.For != nil {
		t.Error("For should be nil")
	}
	if stmt.Password != "123456" {
		t.Errorf("Password = %q, want %q", stmt.Password, "123456")
	}
	if stmt.IsHash {
		t.Error("IsHash should be false for PASSWORD(...) form")
	}
}

func TestSetPassword_PasswordFuncWithFor(t *testing.T) {
	// SET PASSWORD FOR 'jack'@'192.%' = PASSWORD('123456');
	stmt := parseSetPasswordStmt(t, `SET PASSWORD FOR 'jack'@'192.%' = PASSWORD('123456')`)
	if stmt.For == nil {
		t.Fatal("For should not be nil")
	}
	if stmt.For.Username != "jack" {
		t.Errorf("For.Username = %q, want %q", stmt.For.Username, "jack")
	}
	if stmt.For.Host != "192.%" {
		t.Errorf("For.Host = %q, want %q", stmt.For.Host, "192.%")
	}
	if stmt.Password != "123456" {
		t.Errorf("Password = %q, want %q", stmt.Password, "123456")
	}
}

func TestSetPassword_HashWithFor(t *testing.T) {
	// SET PASSWORD FOR 'jack'@'domain' = '*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9'
	stmt := parseSetPasswordStmt(t, `SET PASSWORD FOR 'jack'@'domain' = '*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9'`)
	if stmt.For == nil {
		t.Fatal("For should not be nil")
	}
	if stmt.For.Username != "jack" {
		t.Errorf("For.Username = %q, want %q", stmt.For.Username, "jack")
	}
	if stmt.For.Host != "domain" {
		t.Errorf("For.Host = %q, want %q", stmt.For.Host, "domain")
	}
	if !stmt.IsHash {
		t.Error("IsHash should be true")
	}
}

// ---------------------------------------------------------------------------
// ROW POLICY
// ---------------------------------------------------------------------------

func TestCreateRowPolicy_Basic(t *testing.T) {
	sql := `CREATE ROW POLICY test_policy ON test_table TO test_user USING (k1 = 1)`
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.CreateRowPolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateRowPolicyStmt, got %T", n)
	}
	if stmt.Name != "test_policy" {
		t.Errorf("Name = %q, want %q", stmt.Name, "test_policy")
	}
	if stmt.On.String() != "test_table" {
		t.Errorf("On = %q, want %q", stmt.On.String(), "test_table")
	}
	if stmt.To != "test_user" {
		t.Errorf("To = %q, want %q", stmt.To, "test_user")
	}
	if stmt.Type != "" {
		t.Errorf("Type = %q, want empty", stmt.Type)
	}
}

func TestCreateRowPolicy_Restrictive(t *testing.T) {
	sql := `CREATE ROW POLICY IF NOT EXISTS p1 ON db.tbl AS RESTRICTIVE TO role1 USING (col > 0)`
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.CreateRowPolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateRowPolicyStmt, got %T", n)
	}
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Type != "RESTRICTIVE" {
		t.Errorf("Type = %q, want RESTRICTIVE", stmt.Type)
	}
	if stmt.On.String() != "db.tbl" {
		t.Errorf("On = %q, want %q", stmt.On.String(), "db.tbl")
	}
}

func TestCreateRowPolicy_Permissive(t *testing.T) {
	sql := `CREATE ROW POLICY p2 ON tbl AS PERMISSIVE TO user1 USING (a = 1)`
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.CreateRowPolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateRowPolicyStmt, got %T", n)
	}
	if stmt.Type != "PERMISSIVE" {
		t.Errorf("Type = %q, want PERMISSIVE", stmt.Type)
	}
}

func TestDropRowPolicy_Basic(t *testing.T) {
	sql := `DROP ROW POLICY test_policy ON test_table`
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.DropRowPolicyStmt)
	if !ok {
		t.Fatalf("expected *ast.DropRowPolicyStmt, got %T", n)
	}
	if stmt.Name != "test_policy" {
		t.Errorf("Name = %q, want %q", stmt.Name, "test_policy")
	}
	if stmt.On.String() != "test_table" {
		t.Errorf("On = %q, want %q", stmt.On.String(), "test_table")
	}
}

// ---------------------------------------------------------------------------
// DICTIONARY
// ---------------------------------------------------------------------------

func TestCreateDictionary_Basic(t *testing.T) {
	sql := `CREATE DICTIONARY IF NOT EXISTS my_dict USING my_table (id KEY, name VALUE) LAYOUT(HASH_MAP) PROPERTIES("read_timeout" = "3000")`
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.CreateDictionaryStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateDictionaryStmt, got %T", n)
	}
	if stmt.Name.String() != "my_dict" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_dict")
	}
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.UsingTable.String() != "my_table" {
		t.Errorf("UsingTable = %q, want %q", stmt.UsingTable.String(), "my_table")
	}
	if len(stmt.Columns) != 2 {
		t.Fatalf("Columns: got %d, want 2", len(stmt.Columns))
	}
	if stmt.Columns[0].Name != "id" || stmt.Columns[0].Role != "KEY" {
		t.Errorf("Columns[0] = {%q, %q}, want {id, KEY}", stmt.Columns[0].Name, stmt.Columns[0].Role)
	}
	if stmt.Columns[1].Name != "name" || stmt.Columns[1].Role != "VALUE" {
		t.Errorf("Columns[1] = {%q, %q}, want {name, VALUE}", stmt.Columns[1].Name, stmt.Columns[1].Role)
	}
	if stmt.Layout != "hash_map" {
		t.Errorf("Layout = %q, want %q", stmt.Layout, "hash_map")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("Properties: got %d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "read_timeout" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "read_timeout")
	}
}

func TestAlterDictionary_Basic(t *testing.T) {
	sql := `ALTER DICTIONARY my_dict PROPERTIES("write_timeout" = "5000")`
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.AlterDictionaryStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterDictionaryStmt, got %T", n)
	}
	if stmt.Name.String() != "my_dict" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_dict")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("Properties: got %d, want 1", len(stmt.Properties))
	}
}

func TestDropDictionary_Basic(t *testing.T) {
	sql := `DROP DICTIONARY my_dict`
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.DropDictionaryStmt)
	if !ok {
		t.Fatalf("expected *ast.DropDictionaryStmt, got %T", n)
	}
	if stmt.Name.String() != "my_dict" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_dict")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
}

func TestDropDictionary_IfExists(t *testing.T) {
	sql := `DROP DICTIONARY IF EXISTS db.my_dict`
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.DropDictionaryStmt)
	if !ok {
		t.Fatalf("expected *ast.DropDictionaryStmt, got %T", n)
	}
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name.String() != "db.my_dict" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "db.my_dict")
	}
}

func TestRefreshDictionary_Basic(t *testing.T) {
	sql := `REFRESH DICTIONARY my_dict`
	n := parseOne(t, sql)
	stmt, ok := n.(*ast.RefreshDictionaryStmt)
	if !ok {
		t.Fatalf("expected *ast.RefreshDictionaryStmt, got %T", n)
	}
	if stmt.Name.String() != "my_dict" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "my_dict")
	}
}

// ---------------------------------------------------------------------------
// Node tags
// ---------------------------------------------------------------------------

func TestSecurityNodeTags(t *testing.T) {
	tests := []struct {
		node ast.Node
		want ast.NodeTag
	}{
		{&ast.CreateRowPolicyStmt{}, ast.T_CreateRowPolicyStmt},
		{&ast.DropRowPolicyStmt{}, ast.T_DropRowPolicyStmt},
		{&ast.CreateEncryptKeyStmt{}, ast.T_CreateEncryptKeyStmt},
		{&ast.DropEncryptKeyStmt{}, ast.T_DropEncryptKeyStmt},
		{&ast.DictionaryColumn{}, ast.T_DictionaryColumn},
		{&ast.CreateDictionaryStmt{}, ast.T_CreateDictionaryStmt},
		{&ast.AlterDictionaryStmt{}, ast.T_AlterDictionaryStmt},
		{&ast.DropDictionaryStmt{}, ast.T_DropDictionaryStmt},
		{&ast.RefreshDictionaryStmt{}, ast.T_RefreshDictionaryStmt},
		{&ast.CreateRoleStmt{}, ast.T_CreateRoleStmt},
		{&ast.AlterRoleStmt{}, ast.T_AlterRoleStmt},
		{&ast.DropRoleStmt{}, ast.T_DropRoleStmt},
		{&ast.UserIdentity{}, ast.T_UserIdentity},
		{&ast.CreateUserStmt{}, ast.T_CreateUserStmt},
		{&ast.AlterUserStmt{}, ast.T_AlterUserStmt},
		{&ast.DropUserStmt{}, ast.T_DropUserStmt},
		{&ast.SetPasswordStmt{}, ast.T_SetPasswordStmt},
	}
	for _, tt := range tests {
		if tt.node.Tag() != tt.want {
			t.Errorf("%T.Tag() = %v, want %v", tt.node, tt.node.Tag(), tt.want)
		}
	}
}

func TestSecurityNodeTagStrings(t *testing.T) {
	tests := []struct {
		tag  ast.NodeTag
		want string
	}{
		{ast.T_CreateRowPolicyStmt, "CreateRowPolicyStmt"},
		{ast.T_DropRowPolicyStmt, "DropRowPolicyStmt"},
		{ast.T_CreateEncryptKeyStmt, "CreateEncryptKeyStmt"},
		{ast.T_DropEncryptKeyStmt, "DropEncryptKeyStmt"},
		{ast.T_DictionaryColumn, "DictionaryColumn"},
		{ast.T_CreateDictionaryStmt, "CreateDictionaryStmt"},
		{ast.T_AlterDictionaryStmt, "AlterDictionaryStmt"},
		{ast.T_DropDictionaryStmt, "DropDictionaryStmt"},
		{ast.T_RefreshDictionaryStmt, "RefreshDictionaryStmt"},
		{ast.T_CreateRoleStmt, "CreateRoleStmt"},
		{ast.T_AlterRoleStmt, "AlterRoleStmt"},
		{ast.T_DropRoleStmt, "DropRoleStmt"},
		{ast.T_UserIdentity, "UserIdentity"},
		{ast.T_CreateUserStmt, "CreateUserStmt"},
		{ast.T_AlterUserStmt, "AlterUserStmt"},
		{ast.T_DropUserStmt, "DropUserStmt"},
		{ast.T_SetPasswordStmt, "SetPasswordStmt"},
	}
	for _, tt := range tests {
		if tt.tag.String() != tt.want {
			t.Errorf("NodeTag(%d).String() = %q, want %q", tt.tag, tt.tag.String(), tt.want)
		}
	}
}
