package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/bytebase/omni/driver"
)

// accessDenied is ClickHouse's ACCESS_DENIED: the login lacks a grant the
// statement needs, such as SHOW ACCESS for an access entity or SELECT on the
// system table that lists one.
var accessDenied int32 = 497

func isAccessDenied(err error) bool {
	var engineError *clickhouse.Exception
	return errors.As(err, &engineError) && engineError.Code == accessDenied
}

// getInstanceRoles reads the server's users and roles in name order. Each
// carries, as its Attribute, the access statements the server states for it,
// one per line: its CREATE USER or CREATE ROLE, its grants as SHOW GRANTS
// states them, and the CREATE QUOTA and CREATE SETTINGS PROFILE of every quota
// and settings profile assigned to it by name or by ALL. A user's VALID UNTIL is
// its ValidUntil. Row policies are a table's, so a database's dump states them.
//
// A login that may read system.users but not every other access entity gets the
// users without their attributes, which is what the sync read before it read
// more; a sync that worked for a login keeps working.
func (d *Driver) getInstanceRoles(ctx context.Context) ([]*driver.Role, error) {
	users, err := d.instanceUsers(ctx)
	if err != nil {
		return nil, err
	}
	roles, err := d.describeAccess(ctx, users)
	if isAccessDenied(err) {
		for _, user := range users {
			user.ValidUntil, user.Attribute = nil, nil
		}
		return users, nil
	}
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// instanceUsers reads every user and its VALID UNTIL. A user with several
// authentication methods holds one per method, and the server zeroes a method
// that does not expire: the user expires when the last method does, and never
// when one of them does not.
func (d *Driver) instanceUsers(ctx context.Context) ([]*driver.Role, error) {
	const query = `
		SELECT
			name,
			if(empty(valid_until) OR arrayExists(t -> toUnixTimestamp(t) = 0, valid_until), '', toString(arrayMax(valid_until)))
		FROM system.users
		ORDER BY name`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	defer rows.Close()
	var users []*driver.Role
	for rows.Next() {
		var name, validUntil string
		if err := rows.Scan(&name, &validUntil); err != nil {
			return nil, err
		}
		user := &driver.Role{Name: name}
		if validUntil != "" {
			user.ValidUntil = &validUntil
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	return users, nil
}

// accessKind is one kind of access entity a role's attribute states: the
// system table that lists it and the SHOW CREATE that states one.
type accessKind struct {
	table string
	show  string
}

var (
	accessQuotas   = accessKind{table: "system.quotas", show: "SHOW CREATE QUOTA "}
	accessProfiles = accessKind{table: "system.settings_profiles", show: "SHOW CREATE SETTINGS PROFILE "}
)

// describeAccess adds the server's roles to users and states every one's
// attribute. Its reads need the SHOW ACCESS grants, and the first one the login
// lacks is returned as the server's ACCESS_DENIED.
func (d *Driver) describeAccess(ctx context.Context, users []*driver.Role) ([]*driver.Role, error) {
	roleNames, err := d.accessColumn(ctx, "SELECT name FROM system.roles ORDER BY name")
	if err != nil {
		return nil, err
	}
	type entity struct {
		role     *driver.Role
		show     string
		settings string
	}
	var entities []entity
	for _, user := range users {
		// A user's authentication is a secret the server states in full when
		// the session asks for it, and an attribute is stored by the caller:
		// this statement asks it not to, whatever the connection's settings.
		entities = append(entities, entity{role: user, show: "SHOW CREATE USER ", settings: " SETTINGS format_display_secrets_in_show_and_select = 0"})
	}
	for _, name := range roleNames {
		entities = append(entities, entity{role: &driver.Role{Name: name}, show: "SHOW CREATE ROLE "})
	}
	// A user and a role may share a name, and the user sorts first.
	sort.SliceStable(entities, func(i, j int) bool { return entities[i].role.Name < entities[j].role.Name })

	quotas, err := d.assignedAccess(ctx, accessQuotas)
	if err != nil {
		return nil, err
	}
	profiles, err := d.assignedAccess(ctx, accessProfiles)
	if err != nil {
		return nil, err
	}
	roles := make([]*driver.Role, 0, len(entities))
	for _, e := range entities {
		name := dumpIdentifier(e.role.Name)
		lines, err := d.accessColumn(ctx, e.show+name+e.settings)
		if err != nil {
			return nil, err
		}
		// SHOW GRANTS names a user before a role of the same name, so a role
		// that shares a user's name states the user's grants; the server has
		// no statement that names the role alone.
		grants, err := d.accessColumn(ctx, "SHOW GRANTS FOR "+name)
		if err != nil {
			return nil, err
		}
		lines = append(lines, grants...)
		lines = append(lines, quotas[e.role.Name]...)
		lines = append(lines, profiles[e.role.Name]...)
		attribute := strings.Join(lines, "\n")
		e.role.Attribute = &attribute
		roles = append(roles, e.role)
	}
	return roles, nil
}

// assignedAccess states every entity of kind assigned to a user or role, keyed
// by the grantee's name, in the entity's name order. An entity is assigned by
// its TO list, or by TO ALL without an EXCEPT that names the grantee; one
// assigned to nobody is stated for nobody. The membership is read on the
// server, since a tape cannot hold the arrays that list it.
func (d *Driver) assignedAccess(ctx context.Context, kind accessKind) (map[string][]string, error) {
	query := `
		SELECT
			grantee.name,
			entity.name
		FROM ` + kind.table + ` AS entity
		CROSS JOIN (SELECT name FROM system.users UNION DISTINCT SELECT name FROM system.roles) AS grantee
		WHERE has(entity.apply_to_list, grantee.name)
			OR (entity.apply_to_all AND NOT has(entity.apply_to_except, grantee.name))
		ORDER BY entity.name, grantee.name`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	type assignment struct{ grantee, entity string }
	var assignments []assignment
	for rows.Next() {
		var a assignment
		if err := rows.Scan(&a.grantee, &a.entity); err != nil {
			_ = rows.Close()
			return nil, err
		}
		assignments = append(assignments, a)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, formatErrorWithQuery(err, query)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	// Each entity is stated once, in the order the rows name them.
	stated := map[string][]string{}
	assigned := map[string][]string{}
	for _, a := range assignments {
		statement, ok := stated[a.entity]
		if !ok {
			if statement, err = d.accessColumn(ctx, kind.show+dumpIdentifier(a.entity)); err != nil {
				return nil, err
			}
			stated[a.entity] = statement
		}
		assigned[a.grantee] = append(assigned[a.grantee], statement...)
	}
	return assigned, nil
}

// accessColumn runs a statement that answers one string column and returns its
// rows.
func (d *Driver) accessColumn(ctx context.Context, query string) ([]string, error) {
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	defer rows.Close()
	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		return nil, formatErrorWithQuery(err, query)
	}
	return values, nil
}

// dumpRowPolicy is one row policy a dump writes: the target its header names
// and the statement that restores it.
type dumpRowPolicy struct {
	target string
	stmt   string
}

// dumpRowPolicies reads the row policies on database's tables, and on the
// database itself (ON db.*), as the server states them, in table and name
// order. A row policy is an access entity: it outlives the table and the
// database it names, so a restore meets the policy it dumped, and each
// statement replaces it (CREATE ROW POLICY OR REPLACE). A policy's TO list names
// users and roles, which the server a dump is restored to must have.
//
// A login that may not read row policies gets denied and no policies, and the
// dump says so in a comment instead of failing, as a dump of every object but
// the policies did not fail before.
func dumpRowPolicies(ctx context.Context, conn *sql.Conn, database string) (policies []dumpRowPolicy, denied bool, err error) {
	const query = `
		SELECT
			short_name,
			table
		FROM system.row_policies
		WHERE database = ?
		ORDER BY table, short_name`
	rows, err := conn.QueryContext(ctx, query, database)
	if err != nil {
		return nil, isAccessDenied(err), deniedOr(formatErrorWithQuery(err, query))
	}
	var targets []string
	for rows.Next() {
		var name, table string
		if err := rows.Scan(&name, &table); err != nil {
			_ = rows.Close()
			return nil, false, err
		}
		targets = append(targets, dumpIdentifier(name)+" ON "+dumpRowPolicyTable(database, table))
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, false, formatErrorWithQuery(err, query)
	}
	if err := rows.Close(); err != nil {
		return nil, false, err
	}
	for _, target := range targets {
		show := "SHOW CREATE ROW POLICY " + target
		var stated string
		if err := conn.QueryRowContext(ctx, show).Scan(&stated); err != nil {
			return nil, isAccessDenied(err), deniedOr(formatErrorWithQuery(err, show))
		}
		rest, ok := strings.CutPrefix(stated, "CREATE ROW POLICY ")
		if !ok {
			return nil, false, fmt.Errorf("%s stated %q", show, stated)
		}
		policies = append(policies, dumpRowPolicy{target: target, stmt: "CREATE ROW POLICY OR REPLACE " + rest})
	}
	return policies, false, nil
}

// deniedOr is nil for ACCESS_DENIED, which a caller reports as denied, and err
// otherwise.
func deniedOr(err error) error {
	if isAccessDenied(err) {
		return nil
	}
	return err
}

// dumpRowPolicyTable is the ON target of a row policy: a table, or every table
// of the database when the policy names none.
func dumpRowPolicyTable(database, table string) string {
	if table == "" {
		return dumpIdentifier(database) + ".*"
	}
	return dumpIdentifier(database) + "." + dumpIdentifier(table)
}
