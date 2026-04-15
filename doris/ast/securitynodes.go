package ast

// This file holds security DDL AST node types (T5.5):
//   - ROW POLICY
//   - ENCRYPTION KEY
//   - DICTIONARY
//   - ROLE
//   - USER / SET PASSWORD

// ---------------------------------------------------------------------------
// ROW POLICY
// ---------------------------------------------------------------------------

// CreateRowPolicyStmt represents:
//
//	CREATE ROW POLICY [IF NOT EXISTS] name ON table_name
//	    [AS {RESTRICTIVE | PERMISSIVE}]
//	    TO user_or_role
//	    USING (expr)
type CreateRowPolicyStmt struct {
	Name        string
	IfNotExists bool
	Type        string      // "RESTRICTIVE" or "PERMISSIVE"; empty = default
	On          *ObjectName // ON table_name
	To          string      // TO user_or_role
	Using       string      // USING (expr) — stored as raw text
	Loc         Loc
}

// Tag implements Node.
func (n *CreateRowPolicyStmt) Tag() NodeTag { return T_CreateRowPolicyStmt }

var _ Node = (*CreateRowPolicyStmt)(nil)

// DropRowPolicyStmt represents:
//
//	DROP ROW POLICY name ON table_name
type DropRowPolicyStmt struct {
	Name string
	On   *ObjectName
	Loc  Loc
}

// Tag implements Node.
func (n *DropRowPolicyStmt) Tag() NodeTag { return T_DropRowPolicyStmt }

var _ Node = (*DropRowPolicyStmt)(nil)

// ---------------------------------------------------------------------------
// ENCRYPTION KEY
// ---------------------------------------------------------------------------

// CreateEncryptKeyStmt represents:
//
//	CREATE ENCRYPTKEY [IF NOT EXISTS] name AS 'key_value'
type CreateEncryptKeyStmt struct {
	Name        *ObjectName
	IfNotExists bool
	Key         string // AS 'key_value'
	Loc         Loc
}

// Tag implements Node.
func (n *CreateEncryptKeyStmt) Tag() NodeTag { return T_CreateEncryptKeyStmt }

var _ Node = (*CreateEncryptKeyStmt)(nil)

// DropEncryptKeyStmt represents:
//
//	DROP ENCRYPTKEY [IF EXISTS] name
type DropEncryptKeyStmt struct {
	Name     *ObjectName
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropEncryptKeyStmt) Tag() NodeTag { return T_DropEncryptKeyStmt }

var _ Node = (*DropEncryptKeyStmt)(nil)

// ---------------------------------------------------------------------------
// DICTIONARY
// ---------------------------------------------------------------------------

// DictionaryColumn represents one column entry in a CREATE DICTIONARY column list:
//
//	col_name KEY | VALUE
type DictionaryColumn struct {
	Name string
	Role string // "KEY" or "VALUE"
	Loc  Loc
}

// Tag implements Node.
func (n *DictionaryColumn) Tag() NodeTag { return T_DictionaryColumn }

var _ Node = (*DictionaryColumn)(nil)

// CreateDictionaryStmt represents:
//
//	CREATE DICTIONARY [IF NOT EXISTS] name
//	    USING table_name
//	    (col1 KEY, col2 VALUE, ...)
//	    LAYOUT(HASH_MAP|IP_TRIE|...)
//	    PROPERTIES(...)
type CreateDictionaryStmt struct {
	Name        *ObjectName
	IfNotExists bool
	UsingTable  *ObjectName
	Columns     []*DictionaryColumn
	Layout      string      // HASH_MAP, IP_TRIE, etc.
	Properties  []*Property
	Loc         Loc
}

// Tag implements Node.
func (n *CreateDictionaryStmt) Tag() NodeTag { return T_CreateDictionaryStmt }

var _ Node = (*CreateDictionaryStmt)(nil)

// AlterDictionaryStmt represents:
//
//	ALTER DICTIONARY name PROPERTIES(...)
type AlterDictionaryStmt struct {
	Name       *ObjectName
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *AlterDictionaryStmt) Tag() NodeTag { return T_AlterDictionaryStmt }

var _ Node = (*AlterDictionaryStmt)(nil)

// DropDictionaryStmt represents:
//
//	DROP DICTIONARY [IF EXISTS] name
type DropDictionaryStmt struct {
	Name     *ObjectName
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropDictionaryStmt) Tag() NodeTag { return T_DropDictionaryStmt }

var _ Node = (*DropDictionaryStmt)(nil)

// RefreshDictionaryStmt represents:
//
//	REFRESH DICTIONARY name
type RefreshDictionaryStmt struct {
	Name *ObjectName
	Loc  Loc
}

// Tag implements Node.
func (n *RefreshDictionaryStmt) Tag() NodeTag { return T_RefreshDictionaryStmt }

var _ Node = (*RefreshDictionaryStmt)(nil)

// ---------------------------------------------------------------------------
// ROLE
// ---------------------------------------------------------------------------

// CreateRoleStmt represents:
//
//	CREATE ROLE [IF NOT EXISTS] name [COMMENT 'text']
type CreateRoleStmt struct {
	Name        string
	IfNotExists bool
	Comment     string
	Loc         Loc
}

// Tag implements Node.
func (n *CreateRoleStmt) Tag() NodeTag { return T_CreateRoleStmt }

var _ Node = (*CreateRoleStmt)(nil)

// AlterRoleStmt represents:
//
//	ALTER ROLE name COMMENT 'text'
type AlterRoleStmt struct {
	Name    string
	Comment string
	Loc     Loc
}

// Tag implements Node.
func (n *AlterRoleStmt) Tag() NodeTag { return T_AlterRoleStmt }

var _ Node = (*AlterRoleStmt)(nil)

// DropRoleStmt represents:
//
//	DROP ROLE [IF EXISTS] name
type DropRoleStmt struct {
	Name     string
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropRoleStmt) Tag() NodeTag { return T_DropRoleStmt }

var _ Node = (*DropRoleStmt)(nil)

// ---------------------------------------------------------------------------
// USER
// ---------------------------------------------------------------------------

// UserIdentity represents the 'user'@'host' form used in user statements.
// Host defaults to '%' when omitted.
type UserIdentity struct {
	Username string
	Host     string // '%' when not specified
	Loc      Loc
}

// Tag implements Node.
func (n *UserIdentity) Tag() NodeTag { return T_UserIdentity }

var _ Node = (*UserIdentity)(nil)

// CreateUserStmt represents:
//
//	CREATE USER [IF NOT EXISTS] 'user'@'host'
//	    [IDENTIFIED BY 'password' | IDENTIFIED BY PASSWORD 'hash']
//	    [DEFAULT ROLE 'role']
//	    [password policy options...]
//	    [COMMENT 'text']
type CreateUserStmt struct {
	Name         *UserIdentity
	IfNotExists  bool
	Password     string // IDENTIFIED BY 'password'
	PasswordHash string // IDENTIFIED BY PASSWORD 'hash'
	DefaultRole  string
	Comment      string
	// Password policy options (stored as raw token text for forward-compat)
	PasswordExpire  bool
	PasswordExpireInterval int    // INTERVAL n DAY; 0 = not set
	FailedLoginAttempts int
	PasswordLockTime   int    // lock time in days; 0 = not set
	PasswordHistory    int    // HISTORY n; 0 = not set
	Loc          Loc
}

// Tag implements Node.
func (n *CreateUserStmt) Tag() NodeTag { return T_CreateUserStmt }

var _ Node = (*CreateUserStmt)(nil)

// AlterUserStmt represents:
//
//	ALTER USER [IF EXISTS] 'user'@'host'
//	    [IDENTIFIED BY 'password']
//	    [password policy options...]
//	    [ACCOUNT_LOCK | ACCOUNT_UNLOCK]
//	    [COMMENT 'text']
type AlterUserStmt struct {
	Name         *UserIdentity
	IfExists     bool
	Password     string
	PasswordHash string
	Comment      string
	AccountLock   bool // ACCOUNT_LOCK
	AccountUnlock bool // ACCOUNT_UNLOCK
	// Password policy options
	FailedLoginAttempts int
	PasswordLockTime    int
	Loc                 Loc
}

// Tag implements Node.
func (n *AlterUserStmt) Tag() NodeTag { return T_AlterUserStmt }

var _ Node = (*AlterUserStmt)(nil)

// DropUserStmt represents:
//
//	DROP USER [IF EXISTS] 'user'@'host'
type DropUserStmt struct {
	Name     *UserIdentity
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropUserStmt) Tag() NodeTag { return T_DropUserStmt }

var _ Node = (*DropUserStmt)(nil)

// SetPasswordStmt represents:
//
//	SET PASSWORD [FOR 'user'@'host'] = 'password'
//	SET PASSWORD [FOR 'user'@'host'] = PASSWORD('cleartext')
type SetPasswordStmt struct {
	For      *UserIdentity // nil when no FOR clause
	Password string        // final password value (cleartext or hash)
	IsHash   bool          // true when the RHS was a bare string (hash), false when PASSWORD(...)
	Loc      Loc
}

// Tag implements Node.
func (n *SetPasswordStmt) Tag() NodeTag { return T_SetPasswordStmt }

var _ Node = (*SetPasswordStmt)(nil)
