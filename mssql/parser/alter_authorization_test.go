package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestAlterAuthorizationNames holds ALTER AUTHORIZATION's securable: ON
// [ class :: ] entity_name, with the class a class_type of its docs, and an
// entity_name of dot-separated parts.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-authorization-transact-sql
func TestAlterAuthorizationNames(t *testing.T) {
	tests := []struct{ sql, name string }{
		{"ALTER AUTHORIZATION ON dbo.MyTable TO newOwner", "dbo.MyTable"},
		{"ALTER AUTHORIZATION ON a TO a", "a"},
		{"ALTER AUTHORIZATION ON OBJECT::dbo.MyTable TO newOwner", "OBJECT::dbo.MyTable"},
		{"ALTER AUTHORIZATION ON SCHEMA::Sales TO dbo", "SCHEMA::Sales"},
		{"ALTER AUTHORIZATION ON DATABASE::mydb TO sa", "DATABASE::mydb"},
		{"ALTER AUTHORIZATION ON TYPE::dbo.MyType TO newOwner", "TYPE::dbo.MyType"},
		{"ALTER AUTHORIZATION ON ROLE::db_datareader TO dbo", "ROLE::db_datareader"},
		{"ALTER AUTHORIZATION ON ASSEMBLY::MyAssembly TO newOwner", "ASSEMBLY::MyAssembly"},
		{"ALTER AUTHORIZATION ON ASYMMETRIC KEY::k TO a", "ASYMMETRIC KEY::k"},
		{"ALTER AUTHORIZATION ON FULLTEXT CATALOG::c TO a", "FULLTEXT CATALOG::c"},
		{"ALTER AUTHORIZATION ON FULLTEXT STOPLIST::s TO a", "FULLTEXT STOPLIST::s"},
		{"ALTER AUTHORIZATION ON MESSAGE TYPE::m TO a", "MESSAGE TYPE::m"},
		{"ALTER AUTHORIZATION ON REMOTE SERVICE BINDING::r TO a", "REMOTE SERVICE BINDING::r"},
		{"ALTER AUTHORIZATION ON SEARCH PROPERTY LIST::p TO a", "SEARCH PROPERTY LIST::p"},
		{"ALTER AUTHORIZATION ON SERVER ROLE::r TO a", "SERVER ROLE::r"},
		{"ALTER AUTHORIZATION ON AVAILABILITY GROUP::g TO a", "AVAILABILITY GROUP::g"},
		{"ALTER AUTHORIZATION ON APPLICATION ROLE::r TO a", "APPLICATION ROLE::r"},
		{"ALTER AUTHORIZATION ON LOGIN::l TO a", "LOGIN::l"},
		{"ALTER AUTHORIZATION ON USER::u TO a", "USER::u"},
		{"ALTER AUTHORIZATION ON SERVER::s TO a", "SERVER::s"},
		{"ALTER AUTHORIZATION ON DATABASE SCOPED CREDENTIAL::c TO a", "DATABASE SCOPED CREDENTIAL::c"},
		{"ALTER AUTHORIZATION ON SYMMETRIC KEY::k TO a", "SYMMETRIC KEY::k"},
		{"ALTER AUTHORIZATION ON XML SCHEMA COLLECTION::dbo.MyXmlSchema TO newOwner", "XML SCHEMA COLLECTION::dbo.MyXmlSchema"},
		// A name with an empty part, as identifierList (TSql160.g:30325) has.
		{"ALTER AUTHORIZATION ON object::dbo..a TO b", "OBJECT::dbo..a"},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			list := ParseAndCheck(t, tt.sql)
			if list.Len() != 1 {
				t.Fatalf("Parse(%q): %d statements, want 1", tt.sql, list.Len())
			}
			stmt, ok := list.Items[0].(*ast.SecurityStmt)
			if !ok {
				t.Fatalf("Parse(%q): %T, want *ast.SecurityStmt", tt.sql, list.Items[0])
			}
			if stmt.Name != tt.name {
				t.Errorf("Parse(%q): Name = %q, want %q", tt.sql, stmt.Name, tt.name)
			}
			if stmt.Options == nil || len(stmt.Options.Items) != 1 {
				t.Fatalf("Parse(%q): Options = %v, want the principal", tt.sql, stmt.Options)
			}
			principal, ok := stmt.Options.Items[0].(*ast.String)
			if !ok {
				t.Fatalf("Parse(%q): principal = %T, want *ast.String", tt.sql, stmt.Options.Items[0])
			}
			if principal.Str == "" {
				t.Errorf("Parse(%q): empty principal", tt.sql)
			}
			checkLocation(t, tt.sql, "SecurityStmt", stmt.Loc)
		})
	}

	// TO SCHEMA OWNER is one principal.
	list := ParseAndCheck(t, "ALTER AUTHORIZATION ON OBJECT::dbo.MyTable TO SCHEMA OWNER")
	stmt := list.Items[0].(*ast.SecurityStmt)
	if got := stmt.Options.Items[0].(*ast.String).Str; got != "SCHEMA OWNER" {
		t.Errorf("principal = %q, want %q", got, "SCHEMA OWNER")
	}
}

// TestAlterAuthorizationRefusals holds the statements SQL Server refuses. Each
// was measured under SET PARSEONLY ON (Msg 102), and each is a row of
// testdata/compat.tsv:
//
//   - ON needs an entity name, and TO needs a principal.
//   - A securable's class is one of the class_type words before `::`;
//     ParseSecurityObjectKind (TSql160.g:13924) refuses any other, so `a :: a`
//     is no securable.
//   - An entity_name is a multiPartIdentifier, so a second word is not one.
func TestAlterAuthorizationRefusals(t *testing.T) {
	for _, sql := range []string{
		"ALTER AUTHORIZATION ON",
		"ALTER AUTHORIZATION ON a",
		"ALTER AUTHORIZATION ON a TO",
		"ALTER AUTHORIZATION ON TO a",
		"ALTER AUTHORIZATION ON a a",
		"ALTER AUTHORIZATION ON a a TO a",
		"ALTER AUTHORIZATION ON a zzz TO a",
		"ALTER AUTHORIZATION ON zzz a TO a",
		"ALTER AUTHORIZATION ON a :: a TO a",
		"ALTER AUTHORIZATION ON a :: zzz a TO a",
		"ALTER AUTHORIZATION ON a TO SCHEMA",
		"ALTER AUTHORIZATION ON TRIGGER::a TO b",
		"ALTER AUTHORIZATION ON AUDIT::a TO b",
		"ALTER AUTHORIZATION ON EXTERNAL DATA SOURCE::a TO b",
		"ALTER AUTHORIZATION ON FULLTEXT INDEX::a TO b",
		"ALTER AUTHORIZATION ON XML SCHEMA::a TO b",
		"ALTER AUTHORIZATION ON a TO a zzz",
		"ALTER AUTHORIZATION a TO a",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): accepted, want rejected", sql)
		} else if _, ok := err.(*ParseError); !ok {
			t.Errorf("Parse(%q): %T, want *ParseError", sql, err)
		}
	}
}
