package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"testing"

	mariadbparser "github.com/bytebase/omni/mariadb/parser"
)

// migrationGaps lists the migrations, by "<case>/<direction>/<source>", whose
// plan applied to the server does not converge, and why. The source is "text"
// for catalogs loaded from the schema texts and "snapshot" for ones loaded
// from the synced snapshots. An entry that converges fails, so a fix removes
// it.
var migrationGaps = map[string]string{
	"11.4.13/Verify no migration generated for identical schemas (tests primary key, index, check constraints, DEFAULT NULL fixes)/forward/snapshot":  lossHashBefore118,
	"11.4.13/Verify no migration generated for identical schemas (tests primary key, index, check constraints, DEFAULT NULL fixes)/backward/snapshot": lossHashBefore118,
	"10.11.19/MariaDB INET4 and UUID keys/forward/text":                                                   gapInetCast1011,
	"10.11.19/MariaDB INET4 and UUID keys/forward/snapshot":                                               gapInetCast1011,
	"10.11.19/MariaDB INET4 and UUID keys/backward/text":                                                  gapInetCast1011,
	"10.11.19/MariaDB INET4 and UUID keys/backward/snapshot":                                              gapInetCast1011,
	"Alter table with column additions, type changes, and constraints/backward/snapshot":                  lossOnUpdate,
	"Operations on partitioned tables/forward/snapshot":                                                   lossPartitions,
	"Operations on partitioned tables/backward/snapshot":                                                  lossPartitions,
	"Re-adding a primary key and an overlapping secondary index on a partitioned table/forward/snapshot":  lossPartitions,
	"Re-adding a primary key and an overlapping secondary index on a partitioned table/backward/snapshot": lossPartitions,
	"Recreating a partitioned table from metadata/forward/snapshot":                                       lossPartitions,
	"Reverse: Alter table columns/forward/snapshot":                                                       lossOnUpdate,
	"MariaDB column types and defaults/forward/snapshot":                                                  lossOnUpdate,
	"MariaDB invisible columns/forward/snapshot":                                                          lossInvisible,
	"MariaDB invisible columns/backward/snapshot":                                                         lossInvisible,
	"MariaDB system-versioned tables/forward/snapshot":                                                    lossVersioning,
	"MariaDB application-time periods/forward/snapshot":                                                   lossPeriod,
	"MariaDB application-time periods/backward/snapshot":                                                  lossPeriod,
	"MariaDB table options/forward/snapshot":                                                              lossHash,
	"MariaDB index forms/forward/snapshot":                                                                lossIgnored,
	"MariaDB partitioned tables/forward/snapshot":                                                         lossPartitions,
	"MariaDB partitioned tables/backward/snapshot":                                                        lossPartitions,
	"MariaDB system-time partitioning/forward/snapshot":                                                   lossPartitions,
	"MariaDB system-time partitioning/backward/snapshot":                                                  lossPartitions,
	"MariaDB views/forward/snapshot":                                                                      lossView,
	"MariaDB routines/forward/text":                                                                       gapForRange,
	"MariaDB routines/forward/snapshot":                                                                   lossAggregate + "; " + gapForRange,
	"MariaDB routines/backward/text":                                                                      gapForRange,
	"MariaDB routines/backward/snapshot":                                                                  lossAggregate + "; " + gapForRange,
	"MariaDB timestamps/forward/snapshot":                                                                 lossOnUpdate,
	"MariaDB vector index/forward/snapshot":                                                               lossVector,
	"MariaDB compressed columns/forward/text":                                                             gapCompressed,
	"MariaDB compressed columns/forward/snapshot":                                                         gapCompressed,
	"MariaDB compressed columns/backward/text":                                                            gapCompressed,
	"MariaDB compressed columns/backward/snapshot":                                                        gapCompressed,
	"MariaDB sequences/forward/text":                                                                      gapSequenceRestart,
	"MariaDB sequences/forward/snapshot":                                                                  gapSequenceRestart,
	"pair sequence limits leaving out the start/forward/text":                                             gapSequenceRestart,
	"MariaDB column renames and reorders/forward/text":                                                    gapColumnOrder,
	"MariaDB column renames and reorders/forward/snapshot":                                                gapColumnOrder,
	"MariaDB column renames and reorders/backward/text":                                                   gapColumnOrder,
	"MariaDB column renames and reorders/backward/snapshot":                                               gapColumnOrder,
}

const (
	gapSequenceRestart = "plan: the new limits leave out the sequence's next value, which ALTER SEQUENCE refuses (4085); the plan warns rather than restart a sequence whose next value it cannot see"
	lossAggregate      = "sync: the AGGREGATE of an aggregate function"
	gapColumnOrder     = "omni: the diff does not compare column order, so a plan moves no column"
)

// TestOracle_MigrationConverges applies each plan GenerateMigration writes
// from schema A to schema B to a database holding A, and requires the server
// then to hold what it holds for B: every table, view, sequence, routine,
// trigger and event. The schema pairs are the snapshot fixture's, in both
// directions, and diffPairs'.
func TestOracle_MigrationConverges(t *testing.T) {
	seen := make(map[string]bool)
	forEachServer(t, func(t *testing.T, srv oracleServer, conn *sql.Conn, database string) {
		type side struct {
			schema string
			meta   map[string]any
		}
		type migration struct {
			name     string
			from, to side
		}
		var migrations []migration
		for _, tc := range readSnapshotCases(t, srv.fixture) {
			old, new := side{tc.OldSchema, tc.OldMetadata}, side{tc.NewSchema, tc.NewMetadata}
			migrations = append(migrations,
				migration{tc.Description + "/forward", old, new},
				migration{tc.Description + "/backward", new, old})
		}
		for _, p := range diffPairs {
			if _, refused := refusedPairs[srv.name+"/"+p.name]; refused {
				continue // TestOracle_DiffMatchesServer holds the server to refusing it
			}
			a, b := side{schema: p.a + ";"}, side{schema: p.b + ";"}
			migrations = append(migrations,
				migration{"pair " + p.name + "/forward", a, b},
				migration{"pair " + p.name + "/backward", b, a})
		}

		otherDatabase := database + "_other"
		if _, err := conn.ExecContext(context.Background(), "CREATE DATABASE IF NOT EXISTS "+quoteIdent(otherDatabase)); err != nil {
			t.Fatal(err)
		}
		for _, m := range migrations {
			t.Run(m.name, func(t *testing.T) {
				startsWritten := writtenStarts.MatchString(m.from.schema + m.to.schema)
				resetDatabase(t, conn, database)
				applySchema(t, conn, m.to.schema)
				want := migrationObjects(t, conn, startsWritten)

				sources := []string{"text"}
				if m.from.meta != nil {
					sources = append(sources, "snapshot")
				}
				for _, source := range sources {
					key := m.name + "/" + source
					seen[key], seen[srv.name+"/"+key] = true, true
					gap, listed := gapFor(migrationGaps, srv, key)
					var from, to *Catalog
					var err error
					if source == "text" {
						from, err = schemaTextCatalogIn(m.from.schema, database, srv.version)
						if err == nil {
							to, err = schemaTextCatalogIn(m.to.schema, database, srv.version)
						}
					} else {
						from, to = newCatalogAt(srv.version), newCatalogAt(srv.version)
						if _, err = from.LoadMetadataInto(context.Background(), metadataFixture(t, m.from.meta), database); err == nil {
							_, err = to.LoadMetadataInto(context.Background(), metadataFixture(t, m.to.meta), database)
						}
					}
					if err != nil {
						if !listed {
							t.Errorf("%s: load: %v", source, err)
						}
						continue
					}
					plan := GenerateMigration(from, to, Diff(from, to))

					resetDatabase(t, conn, database)
					applySchema(t, conn, m.from.schema)
					// The plan names each object with its database, so it runs
					// from another one, and with explicit_defaults_for_timestamp
					// OFF, which the plan sets for itself where it matters.
					useDatabase(t, conn, otherDatabase)
					setSession(t, conn, "explicit_defaults_for_timestamp = OFF")
					applyErr := applyPlan(conn, plan)
					setSession(t, conn, "explicit_defaults_for_timestamp = @@global.explicit_defaults_for_timestamp")
					useDatabase(t, conn, database)
					var got map[string]string
					if applyErr == nil {
						got = migrationObjects(t, conn, startsWritten)
					}
					converged := applyErr == nil && sameObjects(got, want)
					switch {
					case !converged && !listed:
						t.Errorf("%s: does not converge: %v\nplan:\n%s\n%s", source, applyErr, plan.SQL(), objectsDiff(got, want))
					case converged && listed:
						t.Errorf("%s: listed as a gap (%s) but converges: remove the entry", source, gap)
					}
				}
			})
		}
	})
	for key := range migrationGaps {
		if !seen[key] {
			t.Errorf("migration gap %q names no migration", key)
		}
	}
}

// schemaTextCatalogIn is schemaTextCatalog into the database named database.
func schemaTextCatalogIn(schema, database string, version Version) (*Catalog, error) {
	prefix := "CREATE DATABASE " + quoteIdent(database) + ";\nUSE " + quoteIdent(database) + ";\n"
	c, err := LoadSDLWithVersion(prefix+schema, version)
	if err == nil || !strings.Contains(err.Error(), "SDL does not allow") {
		return c, err
	}
	c = newCatalogAt(version)
	c.SetForeignKeyChecks(false)
	results, err := c.Exec(prefix+schema, &ExecOptions{ContinueOnError: true})
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		if r.Error != nil {
			return nil, fmt.Errorf("line %d: %w", r.Line, r.Error)
		}
	}
	return c, nil
}

// applyPlan runs each statement of the plan in order. An op may hold several
// statements, such as a routine framed by its session context.
func applyPlan(conn *sql.Conn, plan *MigrationPlan) error {
	for _, op := range plan.Ops {
		for _, seg := range mariadbparser.Split(op.SQL) {
			stmt := strings.TrimSpace(seg.Text)
			if stmt == "" {
				continue
			}
			if _, err := conn.ExecContext(context.Background(), stmt); err != nil {
				return fmt.Errorf("%q: %w", stmt, err)
			}
		}
	}
	return nil
}

func sameObjects(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func objectsDiff(got, want map[string]string) string {
	var names []string
	for k := range want {
		names = append(names, k)
	}
	for k := range got {
		if _, ok := want[k]; !ok {
			names = append(names, k)
		}
	}
	sort.Strings(names)
	var b strings.Builder
	for _, k := range names {
		if got[k] != want[k] {
			fmt.Fprintf(&b, "%s\n--- got\n%s\n--- want\n%s\n", k, got[k], want[k])
		}
	}
	return b.String()
}

// migrationObjects is storedObjects with each routine stated by what the
// server holds for it rather than by SHOW CREATE, which echoes the statement as
// it was written: a plan writes the parameters, the characteristics and their
// order its own way.
func migrationObjects(t *testing.T, conn *sql.Conn, startsWritten bool) map[string]string {
	t.Helper()
	out := storedObjects(t, conn, startsWritten)
	for k := range out {
		if strings.HasPrefix(k, "function ") || strings.HasPrefix(k, "procedure ") {
			delete(out, k)
		}
	}
	rows, err := conn.QueryContext(context.Background(), `
		SELECT LOWER(r.ROUTINE_TYPE), r.ROUTINE_NAME, CONCAT_WS('\n',
			r.DTD_IDENTIFIER, r.IS_DETERMINISTIC, r.SQL_DATA_ACCESS, r.SECURITY_TYPE,
			r.ROUTINE_COMMENT, r.SQL_MODE, r.ROUTINE_DEFINITION,
			(SELECT GROUP_CONCAT(CONCAT_WS(' ', p.PARAMETER_MODE, LOWER(p.PARAMETER_NAME), p.DTD_IDENTIFIER) ORDER BY p.ORDINAL_POSITION SEPARATOR ', ')
			 FROM information_schema.PARAMETERS p
			 WHERE p.SPECIFIC_SCHEMA = r.ROUTINE_SCHEMA AND p.SPECIFIC_NAME = r.SPECIFIC_NAME AND p.ROUTINE_TYPE = r.ROUTINE_TYPE AND p.ORDINAL_POSITION > 0))
		FROM information_schema.ROUTINES r WHERE r.ROUTINE_SCHEMA = DATABASE()`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var kind, name string
		var body sql.NullString
		if err := rows.Scan(&kind, &name, &body); err != nil {
			t.Fatal(err)
		}
		out[kind+" "+name] = body.String
	}
	return out
}

func useDatabase(t *testing.T, conn *sql.Conn, database string) {
	t.Helper()
	if _, err := conn.ExecContext(context.Background(), "USE "+quoteIdent(database)); err != nil {
		t.Fatal(err)
	}
}

func setSession(t *testing.T, conn *sql.Conn, assignment string) {
	t.Helper()
	if _, err := conn.ExecContext(context.Background(), "SET "+assignment); err != nil {
		t.Fatal(err)
	}
}
