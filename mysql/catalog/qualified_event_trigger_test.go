package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// qualifiedNameScripts run with `a` current and tables a.t, b.t and b.only_b there. An event and a
// trigger belong to the database their name is qualified with, not to the current one.
var qualifiedNameScripts = map[string][]string{
	"events in each database": {
		"CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1",
		"CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2",
	},
	"DROP EVENT leaves the current database's": {
		"CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1",
		"CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2",
		"DROP EVENT b.e",
	},
	"ALTER EVENT finds it in its database": {
		"CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2",
		"ALTER EVENT b.e COMMENT 'x'",
	},
	"a missing event in another database": {
		"CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1",
		"DROP EVENT b.e",
		"ALTER EVENT b.e COMMENT 'x'",
	},
	"RENAME TO another database": {
		"CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2",
		"ALTER EVENT b.e RENAME TO a.f",
	},
	"RENAME TO an unqualified name goes to the current database": {
		"CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO SELECT 2",
		"ALTER EVENT b.e RENAME TO f",
	},
	"RENAME TO a name that is taken": {
		"CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1",
		"CREATE EVENT b.f ON SCHEDULE EVERY 1 DAY DO SELECT 2",
		"ALTER EVENT b.f RENAME TO a.e",
	},
	"a failed RENAME TO leaves the event as it was": {
		"CREATE EVENT a.e ON SCHEDULE EVERY 1 DAY DO SELECT 1",
		"CREATE EVENT a.f ON SCHEDULE EVERY 1 DAY DO SELECT 2",
		"ALTER EVENT a.e RENAME TO f COMMENT 'new'",
		"ALTER EVENT a.e RENAME TO nosuch.g COMMENT 'new'",
	},
	"an event in a missing database": {
		"CREATE EVENT nosuch.e ON SCHEDULE EVERY 1 DAY DO SELECT 1",
	},
	"a trigger and a table in different databases": {
		"CREATE TRIGGER b.tr BEFORE INSERT ON a.t FOR EACH ROW SET NEW.id = 1",
		"CREATE TRIGGER tr BEFORE INSERT ON b.t FOR EACH ROW SET NEW.id = 1",
	},
	"an unqualified table is in the trigger's database": {
		"CREATE TRIGGER b.tr BEFORE INSERT ON t FOR EACH ROW SET NEW.id = 1",
		"CREATE TRIGGER b.tr2 BEFORE INSERT ON only_b FOR EACH ROW SET NEW.id = 1",
	},
	"a trigger and a table in one database": {
		"CREATE TRIGGER b.tr BEFORE INSERT ON b.t FOR EACH ROW SET NEW.id = 1",
	},
	"a table the trigger's database does not have": {
		"CREATE TRIGGER tr BEFORE INSERT ON only_b FOR EACH ROW SET NEW.id = 1",
	},
	"a trigger in a missing database": {
		"CREATE TRIGGER nosuch.tr BEFORE INSERT ON t FOR EACH ROW SET NEW.id = 1",
	},
	"a missing trigger database and a table elsewhere": {
		"CREATE TRIGGER nosuch.tr BEFORE INSERT ON a.t FOR EACH ROW SET NEW.id = 1",
	},
	"a table in a missing database": {
		"CREATE TRIGGER a.tr BEFORE INSERT ON nosuch.t FOR EACH ROW SET NEW.id = 1",
	},
	"a wrong database and a missing table": {
		"CREATE TRIGGER b.tr BEFORE INSERT ON a.nosuch FOR EACH ROW SET NEW.id = 1",
	},
}

func qualifiedCatalogCode(c *Catalog, stmt string) int {
	results, err := c.Exec(stmt, nil)
	if err != nil {
		return -1
	}
	for _, r := range results {
		var cerr *Error
		if errors.As(r.Error, &cerr) {
			return cerr.Code
		}
		if r.Error != nil {
			return -1
		}
	}
	return 0
}

func qualifiedServerState(t *testing.T, conn *sql.DB) string {
	t.Helper()
	ctx := context.Background()
	var state []string
	rows, err := conn.QueryContext(ctx, "SELECT CONCAT('event ', EVENT_SCHEMA, '.', EVENT_NAME, ' ', EVENT_COMMENT) FROM information_schema.EVENTS WHERE EVENT_SCHEMA IN ('a','b') UNION ALL SELECT CONCAT('trigger ', TRIGGER_SCHEMA, '.', TRIGGER_NAME, ' on ', EVENT_OBJECT_TABLE) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA IN ('a','b')")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			t.Fatal(err)
		}
		state = append(state, s)
	}
	sort.Strings(state)
	return strings.Join(state, "; ")
}

func qualifiedCatalogState(c *Catalog) string {
	var state []string
	for _, name := range []string{"a", "b"} {
		db := c.GetDatabase(name)
		for _, e := range db.Events {
			state = append(state, fmt.Sprintf("event %s.%s %s", e.Database.Name, e.Name, e.Comment))
		}
		for _, tr := range db.Triggers {
			state = append(state, fmt.Sprintf("trigger %s.%s on %s", tr.Database.Name, tr.Name, tr.Table))
		}
	}
	sort.Strings(state)
	return strings.Join(state, "; ")
}

// TestQualifiedEventsAndTriggers runs each script on the server and on the catalog, and requires the
// same answer to every statement and the same events and triggers in the same databases.
func TestQualifiedEventsAndTriggers(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()
	conn := ctr.db
	const setup = "CREATE TABLE a.t (id INT); CREATE TABLE b.t (id INT); CREATE TABLE b.only_b (id INT)"
	for _, name := range slices.Sorted(maps.Keys(qualifiedNameScripts)) {
		script := qualifiedNameScripts[name]
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			for _, stmt := range append([]string{
				"DROP DATABASE IF EXISTS a", "DROP DATABASE IF EXISTS b", "CREATE DATABASE a", "CREATE DATABASE b",
				"USE a",
			}, strings.Split(setup, "; ")...) {
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					t.Fatalf("%s: %v", stmt, err)
				}
			}
			c := New()
			if _, err := c.Exec("CREATE DATABASE a; CREATE DATABASE b; USE a; "+setup, nil); err != nil {
				t.Fatal(err)
			}
			for _, stmt := range script {
				server := 0
				if _, err := conn.ExecContext(ctx, stmt); err != nil {
					var merr *gomysql.MySQLError
					if !errors.As(err, &merr) {
						t.Fatalf("%s: %v", stmt, err)
					}
					server = int(merr.Number)
				}
				if got := qualifiedCatalogCode(c, stmt); got != server {
					t.Errorf("%s: catalog answered %d, the server %d", stmt, got, server)
				}
			}
			if got, want := qualifiedCatalogState(c), qualifiedServerState(t, conn); got != want {
				t.Errorf("catalog holds [%s], the server [%s]", got, want)
			}
		})
	}
}
