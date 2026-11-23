package driver

import (
	"context"
	sqldriver "database/sql/driver"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/driver"
)

// TestInstanceRolesRefused holds what SyncInstance reads from a connection the
// server refuses the account statements: no roles, rather than a failed sync,
// which is what the pinned server states to an account without the GRANT
// privilege. A failure the server did not state as a refusal still fails.
func TestInstanceRolesRefused(t *testing.T) {
	refusal := mysqlError(1227, "42000", "errCode = 2, detailMessage = Access denied; you need (at least one of) the (GRANT) privilege(s) for this operation")
	tests := []struct {
		name    string
		users   error
		roles   error
		wantErr bool
	}{
		{name: "users refused", users: refusal},
		{name: "roles refused", roles: refusal},
		{name: "users unanswered", users: mysqlError(1105, "HY000", "errCode = 2, detailMessage = query timeout"), wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := stubDatabase(t, func(query string) ([]string, [][]sqldriver.Value, error) {
				switch query {
				case instanceUsersQuery:
					if test.users != nil {
						return nil, nil, test.users
					}
					return []string{"UserIdentity", "Comment", "Password", "Roles", "GlobalPrivs"}, [][]sqldriver.Value{{"'root'@'%'", "ROOT", "No", "operator", "Node_priv,Admin_priv"}}, nil
				case instanceRolesQuery:
					if test.roles != nil {
						return nil, nil, test.roles
					}
					return []string{"Name", "Comment", "Users", "GlobalPrivs"}, [][]sqldriver.Value{{"operator", nil, "'root'@'%'", "Node_priv,Admin_priv"}}, nil
				}
				return nil, nil, fmt.Errorf("unexpected statement: %q", query)
			})
			roles, err := (&Driver{db: db}).getInstanceRoles(context.Background())
			if test.wantErr {
				if err == nil {
					t.Fatalf("getInstanceRoles = %v, want an error", roles)
				}
				return
			}
			if err != nil || roles != nil {
				t.Fatalf("getInstanceRoles = %v, %v; want no roles and no error", roles, err)
			}
		})
	}
}

// TestInstanceRolesColumns holds how a row becomes a role: the name column,
// then every privilege column the row states, by name, in the server's order,
// whatever columns a version adds; the columns that describe the account
// rather than its privileges are left out, and so is an empty one.
func TestInstanceRolesColumns(t *testing.T) {
	db := stubDatabase(t, func(query string) ([]string, [][]sqldriver.Value, error) {
		switch query {
		case instanceUsersQuery:
			return []string{"UserIdentity", "Comment", "Password", "Roles", "GlobalPrivs", "NewerPrivs"}, [][]sqldriver.Value{
				{"'b'@'%'", "two", "Yes", "", nil, "x: Usage_priv"},
				{"'a'@'%'", "", "No", "r", "Admin_priv", nil},
			}, nil
		case instanceRolesQuery:
			return []string{"Name", "Comment", "Users", "GlobalPrivs"}, [][]sqldriver.Value{
				{"r", "a role", "'a'@'%'", nil},
			}, nil
		}
		return nil, nil, fmt.Errorf("unexpected statement: %q", query)
	})
	roles, err := (&Driver{db: db}).getInstanceRoles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	attribute := func(s string) *string { return &s }
	want := []*driver.Role{
		{Name: "'a'@'%'", Attribute: attribute("Roles: r\nGlobalPrivs: Admin_priv")},
		{Name: "'b'@'%'", Attribute: attribute("NewerPrivs: x: Usage_priv")},
		{Name: "r"},
	}
	if diff := cmp.Diff(want, roles); diff != "" {
		t.Errorf("getInstanceRoles (-want +got):\n%s", diff)
	}
}
