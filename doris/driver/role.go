package driver

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/bytebase/omni/driver"
)

// Doris states every account's privileges in one statement and every role's in
// another, one row each, with a column per kind of object a privilege is on.
// The columns change across versions (3.x added the compute group and storage
// vault ones), so they are read by name. SHOW GRANTS FOR states the same row as
// SHOW ALL GRANTS, not GRANT statements as MySQL does, so there is no statement
// form to read.
const (
	instanceUsersQuery = "SHOW ALL GRANTS"
	instanceRolesQuery = "SHOW ROLES"
)

// privilegeTable is how one of the two statements above names its rows: the
// column holding the account or role, and the columns that are not privileges.
type privilegeTable struct {
	query  string
	name   string
	ignore []string
}

var (
	// A user's Roles column names the roles granted to it; Comment and
	// Password (whether one is set) describe the account, not its privileges.
	instanceUsers = privilegeTable{query: instanceUsersQuery, name: "UserIdentity", ignore: []string{"Comment", "Password"}}
	// A role's Users column lists its members, which their own rows name.
	instanceRoles = privilegeTable{query: instanceRolesQuery, name: "Name", ignore: []string{"Comment", "Users"}}
)

// getInstanceRoles returns every account, as 'user'@'host', then every role, by
// name, each with its privileges as "Column: value" lines in the server's own
// terms. Both statements need the GRANT privilege; a connection refused them
// reads no roles rather than failing the sync, as SyncInstance reads the
// databases it can see.
func (d *Driver) getInstanceRoles(ctx context.Context) ([]*driver.Role, error) {
	var roles []*driver.Role
	for _, table := range []privilegeTable{instanceUsers, instanceRoles} {
		rows, err := d.readPrivileges(ctx, table)
		if err != nil {
			if probeRefused(ctx, "", err) {
				return nil, nil
			}
			return nil, err
		}
		roles = append(roles, rows...)
	}
	return roles, nil
}

// readPrivileges reads one privilege statement's rows, sorted by name. A
// privilege column the row leaves empty is not stated.
func (d *Driver) readPrivileges(ctx context.Context, table privilegeTable) ([]*driver.Role, error) {
	rows, err := d.db.QueryContext(ctx, table.query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", table.query, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	nameColumn := -1
	for i, column := range columns {
		if strings.EqualFold(column, table.name) {
			nameColumn = i
		}
	}
	if nameColumn < 0 {
		return nil, fmt.Errorf("%s: no %s column in %v", table.query, table.name, columns)
	}
	var roles []*driver.Role
	for rows.Next() {
		values := make([]sql.NullString, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("%s: %w", table.query, err)
		}
		var lines []string
		for i, column := range columns {
			if i == nameColumn || !values[i].Valid || values[i].String == "" || ignoredColumn(table.ignore, column) {
				continue
			}
			lines = append(lines, column+": "+values[i].String)
		}
		role := &driver.Role{Name: values[nameColumn].String}
		if len(lines) > 0 {
			attribute := strings.Join(lines, "\n")
			role.Attribute = &attribute
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", table.query, err)
	}
	sort.Slice(roles, func(i, j int) bool { return roles[i].Name < roles[j].Name })
	return roles, nil
}

func ignoredColumn(ignore []string, column string) bool {
	for _, name := range ignore {
		if strings.EqualFold(name, column) {
			return true
		}
	}
	return false
}
