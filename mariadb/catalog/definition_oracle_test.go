package catalog

import (
	"context"
	"database/sql"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/internal/testcontainer"
	"github.com/bytebase/omni/metadata"
)

// definitionLosses lists the objects of the fixture that a restored dump holds
// other than the server held them, and why. Mostly it is what Bytebase's sync
// leaves out of the snapshot (see syncLosses). An entry whose object restores
// equal fails, so a fix removes it.
var definitionLosses = map[string]string{
	"11.4.13/Verify no migration generated for identical schemas (tests primary key, index, check constraints, DEFAULT NULL fixes)/new/table t1":         lossHashBefore118,
	"11.4.13/Verify no migration generated for identical schemas (tests primary key, index, check constraints, DEFAULT NULL fixes)/new/table test_table": lossHashBefore118,
	"Alter table with column additions, type changes, and constraints/old/table products":                                                                lossOnUpdate,
	"Reverse: Alter table columns/new/table products":                                                                                                    lossOnUpdate,
	"Multiple triggers and audit logging/old/table inventory":                                                                                            lossOnUpdate,
	"Multiple triggers and audit logging/new/table inventory":                                                                                            lossOnUpdate,
	"Reverse: Drop triggers and views/old/table inventory":                                                                                               lossOnUpdate,
	"Reverse: Drop triggers and views/new/table inventory":                                                                                               lossOnUpdate,
	"MariaDB column types and defaults/new/table t":                                                                                                      lossOnUpdate,
	"MariaDB timestamps/new/table ts":                                                                                                                    lossOnUpdate,
	"Operations on partitioned tables/new/table sales_data":                                                                                              lossPartitions,
	"Re-adding a primary key and an overlapping secondary index on a partitioned table/new/table ads_member_asset_by_account_channel_d":                  lossPartitions,
	"Recreating a partitioned table from metadata/new/table ads_member_asset_by_account_channel_d":                                                       lossPartitions,
	"MariaDB partitioned tables/old/table pr":                                                                                                            lossPartitions,
	"MariaDB partitioned tables/new/table pr":                                                                                                            lossPartitions,
	"MariaDB partitioned tables/new/table ph":                                                                                                            lossPartitions,
	"MariaDB partitioned tables/new/table pk":                                                                                                            lossPartitions,
	"MariaDB partitioned tables/new/table pl":                                                                                                            lossPartitions,
	"MariaDB system-time partitioning/new/table sv":                                                                                                      lossPartitions,
	"MariaDB invisible columns/new/table inv":                                                                                                            lossInvisible,
	"MariaDB system-versioned tables/new/table hist":                                                                                                     lossVersioning,
	"MariaDB application-time periods/new/table booking":                                                                                                 lossPeriod,
	"MariaDB table options/new/table mem_t":                                                                                                              lossHash,
	"MariaDB index forms/new/table ix":                                                                                                                   lossIgnored,
	"MariaDB views/new/view v1":                                                                                                                          lossView,
	"MariaDB vector index/new/table emb":                                                                                                                 lossVector,
}

// restoreFailures lists the snapshots whose restore script the server refuses,
// and why.
var restoreFailures = map[string]string{
	"MariaDB routines/new": "sync: AGGREGATE, which the definition sync builds for a function leaves out, so FETCH GROUP NEXT ROW is 4105",
}

// TestOracle_DatabaseDefinitionRestoresSnapshot renders each snapshot synced
// from the pinned server, runs the restore script as the mariadb client would
// and, separately, the SDL, and compares what the server then holds with what
// it held for the schema the snapshot was taken of: every table, view,
// sequence, routine, trigger and event.
func TestOracle_DatabaseDefinitionRestoresSnapshot(t *testing.T) {
	seen := make(map[string]bool)
	forEachServer(t, func(t *testing.T, srv oracleServer, conn *sql.Conn, database string) {
		for _, tc := range readSnapshotCases(t, srv.fixture) {
			for _, side := range []struct {
				name   string
				schema string
				meta   map[string]any
			}{{"old", tc.OldSchema, tc.OldMetadata}, {"new", tc.NewSchema, tc.NewMetadata}} {
				t.Run(tc.Description+"/"+side.name, func(t *testing.T) {
					// STARTS CURRENT_TIMESTAMP depends on the clock, as no STARTS does.
					startsWritten := writtenStarts.MatchString(side.schema)
					resetDatabase(t, conn, database)
					applySchema(t, conn, side.schema)
					want := storedObjects(t, conn, startsWritten)

					meta := metadataFixture(t, side.meta)
					dump, err := DatabaseDefinition(meta, DefinitionOptions{PrintHeader: true})
					if err != nil {
						t.Fatal(err)
					}
					resetDatabase(t, conn, database)
					failure, failureListed := gapFor(restoreFailures, srv, tc.Description+"/"+side.name)
					for _, stmt := range splitLikeMySQLClient(dump) {
						if _, err := conn.ExecContext(context.Background(), stmt); err != nil {
							if failureListed {
								return
							}
							t.Fatalf("restore %q: %v\n%s", stmt, err, dump)
						}
					}
					if failureListed {
						t.Errorf("listed as failing to restore (%s) but restores: remove the entry", failure)
					}
					compareRestored(t, srv, tc.Description+"/"+side.name, want, storedObjects(t, conn, startsWritten), seen)

					sdl, err := DatabaseDefinition(meta, DefinitionOptions{SDLFormat: true})
					if err != nil {
						t.Fatal(err)
					}
					resetDatabase(t, conn, database)
					applySchema(t, conn, sdl)
					compareRestored(t, srv, tc.Description+"/"+side.name, want, storedObjects(t, conn, startsWritten), seen)
				})
			}
		}
	})
	for key := range definitionLosses {
		if !seen[key] {
			t.Errorf("definition loss %q names no object of the fixture", key)
		}
	}
}

func compareRestored(t *testing.T, srv oracleServer, prefix string, want, got map[string]string, seen map[string]bool) {
	t.Helper()
	var names []string
	for name := range want {
		names = append(names, name)
	}
	for name := range got {
		if _, ok := want[name]; !ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		key := prefix + "/" + name
		seen[key], seen[srv.name+"/"+key] = true, true
		loss, listed := gapFor(definitionLosses, srv, key)
		switch equal := got[name] == want[name]; {
		case !equal && !listed:
			t.Errorf("%s differs\n--- restored\n%s\n--- server\n%s", name, got[name], want[name])
		case equal && listed:
			t.Errorf("%s: listed as a definition loss (%s) but restores equal: remove the entry", name, loss)
		}
	}
}

// storedObjects is serverObjects with the current database's routines,
// triggers and events added, keyed "function f", "procedure p", "trigger t",
// "event e", each as its SHOW CREATE writes it without the definer.
func storedObjects(t *testing.T, conn *sql.Conn, startsWritten bool) map[string]string {
	t.Helper()
	out := serverObjects(t, conn)
	ctx := context.Background()
	list := func(query string) []string {
		rows, err := conn.QueryContext(ctx, query)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		var names []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				t.Fatal(err)
			}
			names = append(names, name)
		}
		return names
	}
	show := func(kind, name string, column int) string {
		rows, err := conn.QueryContext(ctx, "SHOW CREATE "+kind+" "+quoteIdent(name))
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		cols, _ := rows.Columns()
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if !rows.Next() {
			t.Fatalf("SHOW CREATE %s %s: no row", kind, name)
		}
		if err := rows.Scan(ptrs...); err != nil {
			t.Fatal(err)
		}
		// The sql_mode a routine, trigger or event was created under is part
		// of what it is, and SHOW CREATE answers it in column 1.
		return "sql_mode=" + vals[1].String + "\n" + viewDefiner.ReplaceAllString(vals[column].String, " ")
	}
	for _, name := range list("SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'FUNCTION'") {
		out["function "+name] = show("FUNCTION", name, 2)
	}
	for _, name := range list("SELECT ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_TYPE = 'PROCEDURE'") {
		out["procedure "+name] = show("PROCEDURE", name, 2)
	}
	// SHOW CREATE TRIGGER echoes the statement as it was written, so a
	// trigger is compared by what information_schema holds for it.
	rows, err := conn.QueryContext(ctx, "SELECT TRIGGER_NAME, CONCAT_WS('\\n', ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_ORDER, SQL_MODE, ACTION_STATEMENT) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE()")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var name, trigger string
		if err := rows.Scan(&name, &trigger); err != nil {
			t.Fatal(err)
		}
		out["trigger "+name] = trigger
	}
	rows.Close()
	for _, name := range list("SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE()") {
		event := strings.TrimSpace(show("EVENT", name, 3))
		if !startsWritten {
			// An event created without STARTS starts when it was created.
			event = eventStarts.ReplaceAllString(event, "")
		}
		out["event "+name] = event
	}
	return out
}

var (
	eventStarts   = regexp.MustCompile(` STARTS '[^']*'`)
	writtenStarts = regexp.MustCompile(`(?i)STARTS\s+'`)
)

// TestOracle_DatabaseDefinitionRestoresSessionAfterRoutines restores a dump
// whose function was created under NO_BACKSLASH_ESCAPES and whose view, written
// after it, holds a backslash escape, and requires the view to read as it was
// created.
func TestOracle_DatabaseDefinitionRestoresSessionAfterRoutines(t *testing.T) {
	conn, err := testcontainer.OpenDatabase(t).Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	database := testcontainer.DatabaseName(t)
	resetDatabase(t, conn, database)
	dump, err := DatabaseDefinition(&metadata.DatabaseSchemaMetadata{
		Name: database,
		Schemas: []*metadata.SchemaMetadata{{
			Functions: []*metadata.FunctionMetadata{{
				Name:                "f",
				Definition:          "CREATE FUNCTION `f`() RETURNS int(11)\nRETURN 1",
				CharacterSetClient:  "latin1",
				CollationConnection: "latin1_swedish_ci",
				SqlMode:             "NO_BACKSLASH_ESCAPES",
			}},
			Views: []*metadata.ViewMetadata{{Name: "v", Definition: "select 'a\\\\b' AS `x`"}},
		}},
	}, DefinitionOptions{PrintHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, stmt := range splitLikeMySQLClient(dump) {
		if _, err := conn.ExecContext(context.Background(), stmt); err != nil {
			t.Fatalf("restore %q: %v\n%s", stmt, err, dump)
		}
	}
	var x string
	if err := conn.QueryRowContext(context.Background(), "SELECT x FROM v").Scan(&x); err != nil {
		t.Fatal(err)
	}
	if x != `a\b` {
		t.Errorf("view reads %q, want %q:\n%s", x, `a\b`, dump)
	}
}

// TestOracle_DatabaseDefinitionRestoresTimestampsWithoutImplicitDefaults
// restores a NOT NULL TIMESTAMP with no default on a session with
// explicit_defaults_for_timestamp OFF, and requires it to stay without one.
func TestOracle_DatabaseDefinitionRestoresTimestampsWithoutImplicitDefaults(t *testing.T) {
	conn, err := testcontainer.OpenDatabase(t).Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	database := testcontainer.DatabaseName(t)
	resetDatabase(t, conn, database)
	dump, err := DatabaseDefinition(&metadata.DatabaseSchemaMetadata{
		Name: database,
		Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{{
			Name:    "t",
			Engine:  "InnoDB",
			Columns: []*metadata.ColumnMetadata{{Name: "ts", Type: "timestamp"}},
		}}}},
	}, DefinitionOptions{PrintHeader: true})
	if err != nil {
		t.Fatal(err)
	}
	stmts := append([]string{"SET explicit_defaults_for_timestamp = OFF"}, splitLikeMySQLClient(dump)...)
	for _, stmt := range stmts {
		if _, err := conn.ExecContext(context.Background(), stmt); err != nil {
			t.Fatalf("restore %q: %v", stmt, err)
		}
	}
	var name, got string
	if err := conn.QueryRowContext(context.Background(), "SHOW CREATE TABLE t").Scan(&name, &got); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "`ts` timestamp NOT NULL,") && !strings.Contains(got, "`ts` timestamp NOT NULL\n") {
		t.Errorf("ts took an implicit default:\n%s\n%s", got, dump)
	}
	var explicit int
	if err := conn.QueryRowContext(context.Background(), "SELECT @@explicit_defaults_for_timestamp").Scan(&explicit); err != nil || explicit != 0 {
		t.Errorf("the script left explicit_defaults_for_timestamp %d (%v), want the session's 0", explicit, err)
	}
}
