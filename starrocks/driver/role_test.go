package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"errors"
	"fmt"
	"io"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestInstanceRolesRefused holds what SyncInstance reads from a connection the
// server refuses the account statements: no roles, rather than a failed sync.
// The refusal is the one the pinned server states to an account without GRANT
// on SYSTEM. A failure the server did not state as a refusal still fails.
func TestInstanceRolesRefused(t *testing.T) {
	refusal := &gomysql.MySQLError{Number: errAccessDenied, Message: "Access denied; you need (at least one of) the GRANT privilege(s) on SYSTEM for this operation."}
	copy(refusal.SQLState[:], "42000")
	tests := []struct {
		name    string
		fail    string
		err     error
		wantErr bool
	}{
		{name: "users refused", fail: instanceUsersQuery, err: refusal},
		{name: "roles refused", fail: instanceRolesQuery, err: refusal},
		{name: "grants refused", fail: "SHOW GRANTS FOR 'root'@'%'", err: refusal},
		{name: "users unanswered", fail: instanceUsersQuery, err: errors.New("connection reset"), wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := sql.OpenDB(roleStubConnector{answer: func(query string) ([]string, [][]sqldriver.Value, error) {
				if query == test.fail {
					return nil, nil, test.err
				}
				switch query {
				case instanceUsersQuery:
					return []string{"User"}, [][]sqldriver.Value{{"'root'@'%'"}}, nil
				case instanceRolesQuery:
					return []string{"Name", "Builtin", "Comment"}, [][]sqldriver.Value{{"root", "true", nil}}, nil
				case "SHOW GRANTS FOR 'root'@'%'", "SHOW GRANTS FOR ROLE `root`":
					return []string{"UserIdentity", "Catalog", "Grants"}, nil, nil
				}
				return nil, nil, fmt.Errorf("unexpected statement: %q", query)
			}})
			t.Cleanup(func() { _ = db.Close() })
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

// roleStubConnector answers every query with answer, so a test reaches a
// refusal without a server.
type roleStubConnector struct {
	answer func(query string) ([]string, [][]sqldriver.Value, error)
}

func (c roleStubConnector) Connect(context.Context) (sqldriver.Conn, error) {
	return roleStubConn(c), nil
}

func (roleStubConnector) Driver() sqldriver.Driver { return nil }

type roleStubConn roleStubConnector

func (roleStubConn) Prepare(string) (sqldriver.Stmt, error) {
	return nil, errors.New("prepare is not supported")
}
func (roleStubConn) Close() error { return nil }
func (roleStubConn) Begin() (sqldriver.Tx, error) {
	return nil, errors.New("transactions are not supported")
}

func (c roleStubConn) QueryContext(_ context.Context, query string, _ []sqldriver.NamedValue) (sqldriver.Rows, error) {
	columns, rows, err := c.answer(query)
	if err != nil {
		return nil, err
	}
	return &roleStubRows{columns: columns, rows: rows}, nil
}

type roleStubRows struct {
	columns []string
	rows    [][]sqldriver.Value
	next    int
}

func (r *roleStubRows) Columns() []string { return r.columns }
func (*roleStubRows) Close() error        { return nil }
func (r *roleStubRows) Next(dest []sqldriver.Value) error {
	if r.next >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.next])
	r.next++
	return nil
}
