package driver

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytebase/omni/driver"
)

// getInstanceRoles reads the account's users with the roles granted to each.
//
// Both views are in SNOWFLAKE.ACCOUNT_USAGE, so the login needs a role that can
// read it, and GRANTS_TO_USERS lags the account by up to an hour: a grant made
// just now may be missing. The attribute is the roles joined, which is what
// Bytebase's instance list showed for a user.
func (d *Driver) getInstanceRoles(ctx context.Context) ([]*driver.Role, error) {
	// Ordered by grantee and role: a user's roles are one attribute, and the
	// order the view happens to answer in is not the account's.
	const grantQuery = `
		SELECT
			GRANTEE_NAME,
			ROLE
		FROM SNOWFLAKE.ACCOUNT_USAGE.GRANTS_TO_USERS
		ORDER BY GRANTEE_NAME, ROLE`
	grantRows, err := d.db.QueryContext(ctx, grantQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", grantQuery, err)
	}
	defer grantRows.Close()
	grants := make(map[string][]string)
	for grantRows.Next() {
		var name, role string
		if err := grantRows.Scan(&name, &role); err != nil {
			return nil, err
		}
		grants[name] = append(grants[name], role)
	}
	if err := grantRows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", grantQuery, err)
	}

	// A user that was deleted and created again has an entry per life, so only
	// the live one is a role here.
	const userQuery = `
		SELECT
			name
		FROM SNOWFLAKE.ACCOUNT_USAGE.USERS
		WHERE DELETED_ON IS NULL
		ORDER BY name`
	rows, err := d.db.QueryContext(ctx, userQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", userQuery, err)
	}
	defer rows.Close()

	var roles []*driver.Role
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		attribute := strings.Join(grants[name], ", ")
		roles = append(roles, &driver.Role{
			Name:      name,
			Attribute: &attribute,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute query %q: %w", userQuery, err)
	}
	return roles, nil
}
