package driver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bytebase/omni/driver"
)

// getInstanceRoles reads the cluster's users with their attributes written as
// psql's \du writes them. Redshift keeps its users in pg_user, which has no
// role inheritance, role creation or replication to report: those are written
// as the fixed values psql shows for a Redshift user.
func (d *Driver) getInstanceRoles(ctx context.Context) ([]*driver.Role, error) {
	// Reference: https://sourcegraph.com/github.com/postgres/postgres@REL_14_0/-/blob/src/bin/psql/describe.c?L3792
	const query = `
	SELECT
		u.usename AS rolename,
		u.usesuper AS rolsuper,
		true AS rolinherit,
		false AS rolcreaterole,
		u.usecreatedb AS rolcreatedb,
		true AS rolcanlogin,
		-1 AS rolconnlimit,
		u.valuntil as rolvaliduntil
	FROM pg_catalog.pg_user u
	ORDER BY 1;
	`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute roles query: %w", err)
	}
	defer rows.Close()

	var roles []*driver.Role
	for rows.Next() {
		var role string
		var super, inherit, createRole, createDB, canLogin bool
		var connectionLimit int32
		var validUntil sql.NullString
		if err := rows.Scan(
			&role,
			&super,
			&inherit,
			&createRole,
			&createDB,
			&canLogin,
			&connectionLimit,
			&validUntil,
		); err != nil {
			return nil, err
		}

		var attributes []string
		if super {
			attributes = append(attributes, "Superuser")
		}
		if !inherit {
			attributes = append(attributes, "No inheritance")
		}
		if createRole {
			attributes = append(attributes, "Create role")
		}
		if createDB {
			attributes = append(attributes, "Create DB")
		}
		if !canLogin {
			attributes = append(attributes, "Cannot login")
		}
		if connectionLimit >= 0 {
			switch connectionLimit {
			case 0:
				attributes = append(attributes, "No connections")
			case 1:
				attributes = append(attributes, "1 connection")
			default:
				attributes = append(attributes, fmt.Sprintf("%d connections", connectionLimit))
			}
		}
		if validUntil.Valid {
			attributes = append(attributes, fmt.Sprintf("Password valid until %s", validUntil.String))
		}
		attribute := strings.Join(attributes, " ")
		roles = append(roles, &driver.Role{
			Name:      role,
			Attribute: &attribute,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}
