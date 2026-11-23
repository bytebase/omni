package driver

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"io"
	"strings"
	"testing"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/google/go-cmp/cmp"
)

// TestInstanceRolesOlderReleases holds the sync to what an older release
// answers: 24.8 and earlier have no system.users.valid_until (UNKNOWN_IDENTIFIER),
// and 23.3 and earlier no format_display_secrets_in_show_and_select
// (UNKNOWN_SETTING). The answers are the codes those releases gave for these
// statements; the pinned server has both, so no tape can hold them.
func TestInstanceRolesOlderReleases(t *testing.T) {
	db := sql.OpenDB(fakeConnector{answer: func(query string) ([][]string, error) {
		switch {
		case strings.Contains(query, "apply_to_list"):
			return nil, nil
		case strings.Contains(query, "valid_until"):
			return nil, &clickhouse.Exception{Code: unknownIdentifier, Message: "Missing columns: 'valid_until'"}
		case strings.Contains(query, "FROM system.users"):
			return [][]string{{"default", ""}}, nil
		case strings.Contains(query, "FROM system.roles"):
			return nil, nil
		case strings.Contains(query, "SETTINGS format_display_secrets_in_show_and_select"):
			return nil, &clickhouse.Exception{Code: unknownSetting, Message: "Unknown setting format_display_secrets_in_show_and_select"}
		case strings.HasPrefix(query, "SHOW CREATE USER "):
			return [][]string{{"CREATE USER default IDENTIFIED WITH plaintext_password"}}, nil
		case strings.HasPrefix(query, "SHOW GRANTS FOR "):
			return [][]string{{"GRANT SHOW ON *.* TO default"}}, nil
		}
		t.Fatalf("unexpected query %q", query)
		return nil, nil
	}})
	defer func() { _ = db.Close() }()

	roles, err := (&Driver{db: db}).getInstanceRoles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 1 || roles[0].Name != "default" || roles[0].ValidUntil != nil || roles[0].Attribute == nil {
		t.Fatalf("roles = %+v", roles)
	}
	want := "CREATE USER default IDENTIFIED WITH plaintext_password\nGRANT SHOW ON *.* TO default"
	if diff := cmp.Diff(want, *roles[0].Attribute); diff != "" {
		t.Errorf("attribute differs (-want +got):\n%s", diff)
	}
}

// fakeConnector answers every query with rows of string columns, or an error,
// as answer says.
type fakeConnector struct {
	answer func(query string) ([][]string, error)
}

func (c fakeConnector) Connect(context.Context) (sqldriver.Conn, error) { return fakeConn(c), nil }
func (c fakeConnector) Driver() sqldriver.Driver                        { return nil }

type fakeConn fakeConnector

func (fakeConn) Prepare(string) (sqldriver.Stmt, error) { return nil, sqldriver.ErrSkip }
func (fakeConn) Close() error                           { return nil }
func (fakeConn) Begin() (sqldriver.Tx, error)           { return nil, sqldriver.ErrSkip }

func (c fakeConn) QueryContext(_ context.Context, query string, _ []sqldriver.NamedValue) (sqldriver.Rows, error) {
	values, err := c.answer(query)
	if err != nil {
		return nil, err
	}
	columns := 1
	if len(values) > 0 {
		columns = len(values[0])
	}
	return &fakeRows{columns: columns, values: values}, nil
}

type fakeRows struct {
	columns int
	values  [][]string
}

func (r *fakeRows) Columns() []string { return make([]string, r.columns) }
func (*fakeRows) Close() error        { return nil }

func (r *fakeRows) Next(dest []sqldriver.Value) error {
	if len(r.values) == 0 {
		return io.EOF
	}
	for i, v := range r.values[0] {
		dest[i] = v
	}
	r.values = r.values[1:]
	return nil
}
