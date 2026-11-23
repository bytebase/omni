package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	mariadbparser "github.com/bytebase/omni/mariadb/parser"
)

// metadataSnapshotCase is one schema pair of testdata/metadata_migration, with
// the snapshots Bytebase's MariaDB sync took of each on the server the file is
// named for.
type metadataSnapshotCase struct {
	Description string         `yaml:"description"`
	OldSchema   string         `yaml:"oldSchema"`
	NewSchema   string         `yaml:"newSchema"`
	Source      string         `yaml:"source"`
	OldMetadata map[string]any `yaml:"oldMetadata"`
	NewMetadata map[string]any `yaml:"newMetadata"`
}

// pinnedSnapshots is the fixture synced from the pinned server.
const pinnedSnapshots = "testdata/metadata_migration/11.8.8.yaml"

func readSnapshotCases(t *testing.T, path string) []metadataSnapshotCase {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var cases []metadataSnapshotCase
	if err := yaml.Unmarshal(content, &cases); err != nil {
		t.Fatal(err)
	}
	return cases
}

// syncLosses lists the objects of the fixture that load other than the server
// holds them, and why: mostly what Bytebase's sync leaves out of the snapshot,
// which the catalog then cannot install. An entry whose object loads equal to
// the server fails, so a fix removes it.
// The reasons an object of the fixture loads other than the server holds it.
const (
	lossOnUpdate      = "sync: ON UPDATE, which sync reads from EXTRA case-sensitively while MariaDB writes on update current_timestamp() in lower case"
	lossPartitions    = "sync: partitioning, which sync reads on MySQL instances only"
	lossInvisible     = "sync: INVISIBLE, which sync reads on stock MySQL only"
	lossVersioning    = "sync: a column's WITHOUT SYSTEM VERSIONING"
	lossPeriod        = "sync: an application-time PERIOD FOR and WITHOUT OVERLAPS, of which only the implicit check and the key's columns remain"
	lossHash          = "sync: whether a MEMORY key wrote USING HASH, its engine's default"
	lossIgnored       = "sync: IGNORED, which sync reports as visible"
	lossView          = "sync: a view's ALGORITHM and SQL SECURITY"
	lossVector        = "sync: a vector key's M and DISTANCE"
	gapCompressed     = "omni: the parser does not read a column's COMPRESSED attribute"
	lossHashBefore118 = "sync: USING HASH on an InnoDB key, which information_schema reports as BTREE before 11.8"
	gapInetCast1011   = "omni: the plan converts a column between INET6 and INET4 in place, which 10.11 refuses (4078) and 11.4 does"
	gapQuoting        = "omni: SHOW CREATE VIEW backquotes a stored function's name, which the catalog does not yet tell from a builtin's"
)

var syncLosses = map[string]string{
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
	"MariaDB compressed columns/new/table cc":                                                                                                            gapCompressed,
	"MariaDB dependent views and functions/old/view dv1":                                                                                                 gapQuoting,
	"MariaDB dependent views and functions/new/view dv1":                                                                                                 gapQuoting,
}

// TestOracle_LoadMetadataReproducesStoredSchema loads each snapshot synced from
// the pinned server and compares what the catalog then shows for every table,
// view and sequence with what the server shows for the schema the snapshot
// was taken of.
func TestOracle_LoadMetadataReproducesStoredSchema(t *testing.T) {
	seen := make(map[string]bool)
	forEachServer(t, func(t *testing.T, srv oracleServer, conn *sql.Conn, database string) {
		for _, tc := range readSnapshotCases(t, srv.fixture) {
			for _, side := range []struct {
				name   string
				schema string
				meta   map[string]any
			}{{"old", tc.OldSchema, tc.OldMetadata}, {"new", tc.NewSchema, tc.NewMetadata}} {
				t.Run(tc.Description+"/"+side.name, func(t *testing.T) {
					resetDatabase(t, conn, database)
					applySchema(t, conn, side.schema)
					want := serverObjects(t, conn)

					c := newCatalogAt(srv.version)
					report, err := c.LoadMetadataInto(context.Background(), metadataFixture(t, side.meta), database)
					if err != nil {
						t.Fatal(err)
					}
					for key, err := range report.Degraded {
						t.Logf("degraded %s: %v", key, err)
					}
					for key, err := range report.Missing {
						t.Logf("missing %s: %v", key, err)
					}
					got := catalogObjects(c, database)
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
						key := tc.Description + "/" + side.name + "/" + name
						seen[key], seen[srv.name+"/"+key] = true, true
						loss, listed := gapFor(syncLosses, srv, key)
						switch equal := got[name] == want[name]; {
						case !equal && !listed:
							t.Errorf("%s differs\n--- catalog\n%s\n--- server\n%s", name, got[name], want[name])
						case equal && listed:
							t.Errorf("%s: listed as a sync loss (%s) but loads equal: remove the entry", name, loss)
						}
					}
				})
			}
		}
	})
	for key := range syncLosses {
		if !seen[key] {
			t.Errorf("sync loss %q names no object of the fixture", key)
		}
	}
}

func resetDatabase(t *testing.T, conn *sql.Conn, database string) {
	t.Helper()
	for _, s := range []string{
		"DROP DATABASE IF EXISTS " + quoteIdent(database),
		"CREATE DATABASE " + quoteIdent(database),
		"USE " + quoteIdent(database),
	} {
		if _, err := conn.ExecContext(context.Background(), s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
}

// applySchema runs schema with foreign key checks off. A statement that fails
// runs again after the rest, as long as a pass runs one: a dump may list a
// view before the view it reads.
func applySchema(t *testing.T, conn *sql.Conn, schema string) {
	t.Helper()
	ctx := context.Background()
	if _, err := conn.ExecContext(ctx, "SET foreign_key_checks = 0"); err != nil {
		t.Fatal(err)
	}
	var pending []string
	for _, seg := range mariadbparser.Split(schema) {
		if s := strings.TrimSpace(seg.Text); s != "" {
			pending = append(pending, s)
		}
	}
	for len(pending) > 0 {
		var failed []string
		var lastErr error
		for _, s := range pending {
			if _, err := conn.ExecContext(ctx, s); err != nil {
				failed = append(failed, s)
				lastErr = fmt.Errorf("%q: %w", s, err)
			}
		}
		if len(failed) == len(pending) {
			t.Fatalf("apply schema: %v", lastErr)
		}
		pending = failed
	}
	if _, err := conn.ExecContext(ctx, "SET foreign_key_checks = 1"); err != nil {
		t.Fatal(err)
	}
}

var viewDefiner = regexp.MustCompile(` DEFINER=\S+ `)

// serverObjects is what the server shows for each table, view and sequence of
// the current database, keyed "table t", "view v", "sequence s".
func serverObjects(t *testing.T, conn *sql.Conn) map[string]string {
	t.Helper()
	ctx := context.Background()
	rows, err := conn.QueryContext(ctx, "SHOW FULL TABLES")
	if err != nil {
		t.Fatal(err)
	}
	type object struct{ name, kind string }
	var objects []object
	for rows.Next() {
		var name, kind string
		if err := rows.Scan(&name, &kind); err != nil {
			t.Fatal(err)
		}
		objects = append(objects, object{name, kind})
	}
	rows.Close()
	out := make(map[string]string)
	for _, o := range objects {
		var stmt, key string
		switch o.kind {
		case "VIEW":
			stmt, key = "SHOW CREATE VIEW ", "view "
		case "SEQUENCE":
			stmt, key = "SHOW CREATE SEQUENCE ", "sequence "
		default:
			stmt, key = "SHOW CREATE TABLE ", "table "
		}
		r, err := conn.QueryContext(ctx, stmt+quoteIdent(o.name))
		if err != nil {
			t.Fatal(err)
		}
		cols, _ := r.Columns()
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if !r.Next() {
			t.Fatalf("%s%s: no row", stmt, o.name)
		}
		if err := r.Scan(ptrs...); err != nil {
			t.Fatal(err)
		}
		r.Close()
		out[key+o.name] = sortTableKeys(viewDefiner.ReplaceAllString(vals[1].String, " "))
	}
	return out
}

// sortTableKeys sorts the secondary keys and the constraints of a SHOW CREATE
// TABLE text. The server lists them in the order they were created, which a
// snapshot does not keep: sync lists them by name.
func sortTableKeys(ddl string) string {
	lines := strings.Split(ddl, "\n")
	end := slices.IndexFunc(lines, func(line string) bool { return strings.HasPrefix(line, ")") })
	if end < 2 || !strings.HasPrefix(lines[0], "CREATE TABLE") {
		return ddl
	}
	body := lines[1:end]
	var plain, keys []string
	for _, line := range body {
		line = strings.TrimSuffix(line, ",")
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "KEY "), strings.HasPrefix(trimmed, "UNIQUE KEY "),
			strings.HasPrefix(trimmed, "FULLTEXT KEY "), strings.HasPrefix(trimmed, "SPATIAL KEY "),
			strings.HasPrefix(trimmed, "VECTOR KEY "), strings.HasPrefix(trimmed, "CONSTRAINT "):
			keys = append(keys, line)
		default:
			plain = append(plain, line)
		}
	}
	sort.Strings(keys)
	body = append(plain, keys...)
	return lines[0] + "\n" + strings.Join(body, ",\n") + "\n" + strings.Join(lines[end:], "\n")
}

// catalogObjects is what the catalog shows for the same objects.
func catalogObjects(c *Catalog, database string) map[string]string {
	db := c.GetDatabase(database)
	out := make(map[string]string)
	if db == nil {
		return out
	}
	for _, tbl := range db.Tables {
		out["table "+tbl.Name] = sortTableKeys(c.ShowCreateTable(database, tbl.Name))
	}
	for _, v := range db.Views {
		out["view "+v.Name] = viewDefiner.ReplaceAllString(c.ShowCreateView(database, v.Name), " ")
	}
	for _, seq := range db.Sequences {
		out["sequence "+seq.Name] = c.ShowCreateSequence(database, seq.Name)
	}
	return out
}
