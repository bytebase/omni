package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/mysql/internal/testcontainer"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
	"gopkg.in/yaml.v3"
)

// metadataMigrationGaps lists the cases of testdata/metadata_migration.yaml whose plan
// applied to MySQL does not store what newSchema stores, and why. A case that converges
// while listed fails, so a fix removes its entry.
var metadataMigrationGaps = map[string]string{
	// MySQL keeps an expression default DEFAULT (NULL) apart from DEFAULT NULL; the
	// catalog compares them equal, so the plan leaves the column as it was.
	"Drop and recreate constraints with different definitions":                                                              "DEFAULT (NULL) kept",
	"Circular foreign key dependencies and triggers":                                                                        "DEFAULT (NULL) kept",
	"Re-adding a primary key and an overlapping secondary index on a partitioned table":                                     "DEFAULT (NULL) kept",
	"Verify no migration generated for identical schemas (tests primary key, index, check constraints, DEFAULT NULL fixes)": "DEFAULT (NULL) kept",
}

// TestOracle_MetadataMigrationApplyCorrectness applies each fixture's plan, generated
// from snapshots synced from MySQL, to a database set up from oldSchema, and compares
// what MySQL then stores with what it stores for newSchema: SHOW CREATE TABLE and SHOW
// CREATE VIEW, and the information_schema rows of routines, their parameters,
// triggers, and events. Parameter lists are compared by their rows, since MySQL keeps
// a routine's parameter text as written.
func TestOracle_MetadataMigrationApplyCorrectness(t *testing.T) {
	content, err := os.ReadFile("testdata/metadata_migration.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var cases []metadataMigrationCase
	if err := yaml.Unmarshal(content, &cases); err != nil {
		t.Fatal(err)
	}
	// One connection, since the setup selects databases with USE.
	conn, err := testcontainer.OpenRoot(t, testcontainer.Options{MultiStatements: true}).Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, tc := range cases {
		if tc.Source != "synced from MySQL 8.0" {
			continue
		}
		t.Run(tc.Description, func(t *testing.T) {
			applyDB := applyDBName(t, "mm")
			from := loadMetadataFixtureInto(t, tc.OldMetadata, applyDB)
			to := loadMetadataFixtureInto(t, tc.NewMetadata, applyDB)
			plan := GenerateMigration(from, to, Diff(from, to))

			setUpMetadataMigrationDB(t, conn, applyDB, tc.OldSchema)
			for _, op := range plan.Ops {
				if _, err := conn.ExecContext(context.Background(), op.SQL); err != nil {
					t.Fatalf("apply %q: %v\nplan:\n%s", op.SQL, err, plan.SQL())
				}
			}
			setUpMetadataMigrationDB(t, conn, applyDB+"_to", tc.NewSchema)

			got := storedSchema(t, conn, applyDB)
			want := storedSchema(t, conn, applyDB+"_to")
			gap, listed := metadataMigrationGaps[tc.Description]
			switch {
			case got != want && !listed:
				t.Errorf("stored schema differs\nplan:\n%s\n--- got\n%s\n--- want\n%s", plan.SQL(), got, want)
			case got == want && listed:
				t.Errorf("listed as a gap (%s) but converges: remove the entry", gap)
			}
		})
	}
}

func loadMetadataFixtureInto(t *testing.T, m map[string]any, database string) *Catalog {
	t.Helper()
	meta := metadataFixture(t, m)
	c := New()
	if _, err := c.LoadMetadataInto(context.Background(), meta, database); err != nil {
		t.Fatal(err)
	}
	return c
}

// setUpMetadataMigrationDB creates database and runs schema in it, with foreign key
// checks off so tables may reference ones created after them.
func setUpMetadataMigrationDB(t *testing.T, conn *sql.Conn, database, schema string) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		"DROP DATABASE IF EXISTS " + database,
		"CREATE DATABASE " + database,
		"USE " + database,
		"SET foreign_key_checks = 0",
	}
	for _, seg := range mysqlparser.Split(schema) {
		if s := strings.TrimSpace(seg.Text); s != "" {
			stmts = append(stmts, s)
		}
	}
	stmts = append(stmts, "SET foreign_key_checks = 1")
	for _, s := range stmts {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Fatalf("set up %s: %q: %v", database, s, err)
		}
	}
}

var autoIncrementCounter = regexp.MustCompile(` AUTO_INCREMENT=\d+`)

// storedSchema renders what MySQL stores for a database's objects, with the database
// name taken out. Secondary keys are sorted: SHOW CREATE TABLE lists them in creation
// order, which does not change how the table behaves.
func storedSchema(t *testing.T, conn *sql.Conn, database string) string {
	t.Helper()
	ctx := context.Background()
	// SHOW CREATE VIEW leaves the current database's qualifier out.
	if _, err := conn.ExecContext(ctx, "USE "+database); err != nil {
		t.Fatal(err)
	}
	var out []string
	query := func(q string, args ...any) [][]string {
		rows, err := conn.QueryContext(ctx, q, args...)
		if err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		defer rows.Close()
		cols, _ := rows.Columns()
		var result [][]string
		for rows.Next() {
			vals := make([]sql.NullString, len(cols))
			ptrs := make([]any, len(cols))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			if err := rows.Scan(ptrs...); err != nil {
				t.Fatal(err)
			}
			row := make([]string, len(cols))
			for i, v := range vals {
				row[i] = v.String
			}
			result = append(result, row)
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		return result
	}
	tables := query("SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME", database)
	for _, tbl := range tables {
		kind := "TABLE"
		if tbl[1] == "VIEW" {
			kind = "VIEW"
		}
		row := query(fmt.Sprintf("SHOW CREATE %s `%s`.`%s`", kind, database, tbl[0]))
		out = append(out, sortSecondaryKeys(autoIncrementCounter.ReplaceAllString(row[0][1], "")))
	}
	for _, q := range []string{
		"SELECT ROUTINE_TYPE, ROUTINE_NAME, DTD_IDENTIFIER, ROUTINE_DEFINITION, IS_DETERMINISTIC, SQL_DATA_ACCESS, SECURITY_TYPE, ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = ? ORDER BY ROUTINE_TYPE, ROUTINE_NAME",
		"SELECT SPECIFIC_NAME, ORDINAL_POSITION, PARAMETER_MODE, PARAMETER_NAME, DTD_IDENTIFIER FROM information_schema.PARAMETERS WHERE SPECIFIC_SCHEMA = ? ORDER BY SPECIFIC_NAME, ORDINAL_POSITION",
		"SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_TIMING, ACTION_ORDER, ACTION_STATEMENT FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ? ORDER BY TRIGGER_NAME",
		"SELECT EVENT_NAME, EVENT_DEFINITION, EVENT_TYPE, INTERVAL_VALUE, INTERVAL_FIELD, STATUS, ON_COMPLETION FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ? ORDER BY EVENT_NAME",
	} {
		for _, row := range query(q, database) {
			out = append(out, strings.Join(row, " | "))
		}
	}
	slices.Sort(out)
	return strings.ReplaceAll(strings.Join(out, "\n"), "`"+database+"`", "`db`")
}

// sortSecondaryKeys sorts the KEY lines of a SHOW CREATE TABLE, keeping everything
// else in place.
func sortSecondaryKeys(ddl string) string {
	lines := strings.Split(ddl, "\n")
	var keys []string
	var idx []int
	for i, l := range lines {
		trimmed := strings.TrimSpace(l)
		for _, prefix := range []string{"KEY ", "UNIQUE KEY ", "FULLTEXT KEY ", "SPATIAL KEY "} {
			if strings.HasPrefix(trimmed, prefix) {
				keys = append(keys, strings.TrimSuffix(trimmed, ","))
				idx = append(idx, i)
				break
			}
		}
	}
	slices.Sort(keys)
	for n, i := range idx {
		comma := strings.HasSuffix(strings.TrimSpace(lines[i]), ",")
		lines[i] = "  " + keys[n]
		if comma {
			lines[i] += ","
		}
	}
	return strings.Join(lines, "\n")
}
