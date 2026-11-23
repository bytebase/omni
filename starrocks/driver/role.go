package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/driver"
)

// StarRocks lists its accounts and roles by name and states each one's
// privileges as GRANT statements, one row each, under a Grants column. Unlike
// its parent Doris, it has no statement that states every account's privileges
// at once.
const (
	instanceUsersQuery = "SHOW USERS"
	instanceRolesQuery = "SHOW ROLES"
	grantsColumn       = "Grants"
)

// The client error numbers StarRocks refuses a statement for want of a
// privilege under: 5203 is its own, which SHOW USERS, SHOW ROLES and SHOW
// GRANTS FOR state to an account without GRANT on SYSTEM; the other three are
// MySQL's, kept for the statements that still answer with them.
const (
	errDatabaseAccessDenied = 1044
	errObjectAccessDenied   = 1142
	errOperationDenied      = 1227
	errAccessDenied         = 5203
)

// getInstanceRoles returns every account, as 'user'@'host', then every role, by
// name, each with the GRANT statements the server states for it, one per line.
// The statements need GRANT on SYSTEM; a connection refused them reads no roles
// rather than failing the sync, as SyncInstance reads the databases it can see.
func (d *Driver) getInstanceRoles(ctx context.Context) ([]*driver.Role, error) {
	roles, err := d.readInstanceRoles(ctx)
	if err != nil && privilegeRefused(ctx, err) {
		return nil, nil
	}
	return roles, err
}

func (d *Driver) readInstanceRoles(ctx context.Context) ([]*driver.Role, error) {
	users, err := d.firstColumn(ctx, instanceUsersQuery)
	if err != nil {
		return nil, err
	}
	roleNames, err := d.firstColumn(ctx, instanceRolesQuery)
	if err != nil {
		return nil, err
	}
	sort.Strings(users)
	sort.Strings(roleNames)
	roles := make([]*driver.Role, 0, len(users)+len(roleNames))
	for _, user := range users {
		// SHOW USERS states the identity quoted, as SHOW GRANTS FOR takes it;
		// StarRocks refuses a user or role name with a quote in it.
		role, err := d.readGrants(ctx, user, "SHOW GRANTS FOR "+user)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	for _, name := range roleNames {
		role, err := d.readGrants(ctx, name, "SHOW GRANTS FOR ROLE "+quoteBacktick(name))
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// readGrants returns one account or role with the GRANT statements the query
// states for it, in the server's order.
func (d *Driver) readGrants(ctx context.Context, name, query string) (*driver.Role, error) {
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	grants := -1
	for i, column := range columns {
		if strings.EqualFold(column, grantsColumn) {
			grants = i
		}
	}
	if grants < 0 {
		return nil, fmt.Errorf("%s: no %s column in %v", query, grantsColumn, columns)
	}
	var lines []string
	for rows.Next() {
		values, err := scanRow(rows, len(columns))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", query, err)
		}
		if values[grants].Valid && values[grants].String != "" {
			lines = append(lines, values[grants].String)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	role := &driver.Role{Name: name}
	if len(lines) > 0 {
		attribute := strings.Join(lines, "\n")
		role.Attribute = &attribute
	}
	return role, nil
}

// firstColumn reads the first column of every row the query returns.
func (d *Driver) firstColumn(ctx context.Context, query string) ([]string, error) {
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var out []string
	for rows.Next() {
		values, err := scanRow(rows, len(columns))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", query, err)
		}
		out = append(out, values[0].String)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	return out, nil
}

func scanRow(rows *sql.Rows, width int) ([]sql.NullString, error) {
	values := make([]sql.NullString, width)
	dest := make([]any, width)
	for i := range values {
		dest[i] = &values[i]
	}
	return values, rows.Scan(dest...)
}

// privilegeRefused reports whether the server refused a statement for want of
// a privilege. Every other failure, including an ended context, fails the
// sync.
func privilegeRefused(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	var mysqlErr *gomysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return false
	}
	switch mysqlErr.Number {
	case errDatabaseAccessDenied, errObjectAccessDenied, errOperationDenied, errAccessDenied:
		return true
	}
	return false
}
