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

	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
	obparser "github.com/bytebase/omni/oceanbase/parser"
	"gopkg.in/yaml.v3"
)

// migrationGap is why a case does not converge. refusal is part of the error
// OceanBase answers a statement of the plan with, when that is the gap; otherwise
// the plan applies and the stored schema differs.
type migrationGap struct {
	reason  string
	refusal string
}

// metadataMigrationGaps lists the cases of testdata/metadata_migration.yaml whose plan
// applied to OceanBase does not store what newSchema stores, and why. A case that
// converges while listed fails, so a fix removes its entry. "sync" names what
// Bytebase's OceanBase sync leaves out of the snapshot or reports wrongly, so no plan
// can restore it.
var metadataMigrationGaps = map[string]migrationGap{
	// OceanBase refuses MODIFY COLUMN on a table that has a CHECK constraint (1235).
	"Alter table with column additions, type changes, and constraints":                                                      {reason: "MODIFY COLUMN refused next to a CHECK", refusal: "Error 1235 (0A000): Alter column with check constraint not supported"},
	"Drop and recreate constraints with different definitions":                                                              {reason: "MODIFY COLUMN refused next to a CHECK", refusal: "Error 1235 (0A000): Alter column with check constraint not supported"},
	"Verify no migration generated for identical schemas (tests primary key, index, check constraints, DEFAULT NULL fixes)": {reason: "MODIFY COLUMN refused next to a CHECK", refusal: "Error 1235 (0A000): Alter column with check constraint not supported"},
	// sync: no CHECK constraints come back.
	"Basic table operations with columns, constraints, and indexes": {reason: "sync: CHECK constraints"},
	"Drop indexes and constraints from tables":                      {reason: "sync: CHECK constraints"},
	"Reverse: Basic table operations":                               {reason: "sync: CHECK constraints"},
	"Reverse: Add indexes and constraints":                          {reason: "sync: CHECK constraints"},
	// sync: no partitions come back.
	"Operations on partitioned tables":                                                  {reason: "sync: partitions and CHECK constraints"},
	"Re-adding a primary key and an overlapping secondary index on a partitioned table": {reason: "sync: partitions"},
	"Recreating a partitioned table from metadata":                                      {reason: "sync: partitions"},
	// sync: FULLTEXT, SPATIAL and functional keys name OceanBase's hidden columns, so
	// the table stands in.
	"FULLTEXT and SPATIAL indexes":               {reason: "sync: hidden key columns", refusal: "Error 1138 (22004): Invalid use of NULL value"},
	"Reverse: Drop FULLTEXT and SPATIAL indexes": {reason: "sync: hidden key columns"},
	"Reverse: Drop JSON columns and indexes":     {reason: "sync: hidden key columns"},
	// sync: a trigger's body is OceanBase's ACTION_STATEMENT, which repeats the header.
	"Reverse: Views and triggers":      {reason: "sync: trigger body", refusal: "near 'TRIGGER"},
	"Reverse: Drop triggers and views": {reason: "sync: trigger body", refusal: "near 'TRIGGER"},
	// sync: events come back without a definition, so the plan creates a stand-in.
	"Reverse: Drop events and advanced features": {reason: "sync: event definition", refusal: "near 'COMPLETION"},
	// sync: ON UPDATE CURRENT_TIMESTAMP is not read.
	"Reverse: Alter table columns": {reason: "sync: ON UPDATE and CHECK constraints"},
	// The catalog takes utf8mb4's default collation to be MySQL's utf8mb4_0900_ai_ci;
	// OceanBase's is utf8mb4_general_ci, so the plan leaves out the COLLATE the
	// target needs.
	"Reverse: Drop tables with foreign key dependencies": {reason: "default utf8mb4 collation, and sync: CHECK constraints"},
	// The plan writes count(*) as count(0) and function names in lower case;
	// OceanBase stores view and generation text as written.
	"Drop views and stored routines":  {reason: "count(*) written as count(0)"},
	"Reverse: Drop generated columns": {reason: "UPPER written as upper"},
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
	conn, err := testcontainer.OpenWith(t, testcontainer.Options{MultiStatements: true}).Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, tc := range cases {
		if tc.Source != "synced from OceanBase 4.3.5" {
			continue
		}
		t.Run(tc.Description, func(t *testing.T) {
			applyDB := applyDBName(t, "mm")
			from := loadSyncedFixtureInto(t, tc.Description, metadataFixture(t, tc.OldMetadata), applyDB)
			to := loadSyncedFixtureInto(t, tc.Description, metadataFixture(t, tc.NewMetadata), applyDB)
			plan := GenerateMigration(from, to, Diff(from, to))

			setUpMetadataMigrationDB(t, conn, applyDB, tc.OldSchema)
			gap, listed := metadataMigrationGaps[tc.Description]
			for _, op := range plan.Ops {
				if _, err := conn.ExecContext(context.Background(), op.SQL); err != nil {
					if listed && gap.refusal != "" && strings.Contains(err.Error(), gap.refusal) {
						return
					}
					t.Fatalf("apply %q: %v\nplan:\n%s", op.SQL, err, plan.SQL())
				}
			}
			setUpMetadataMigrationDB(t, conn, applyDB+"_to", tc.NewSchema)

			got := storedSchema(t, conn, applyDB)
			want := storedSchema(t, conn, applyDB+"_to")
			switch {
			case listed && gap.refusal != "":
				t.Errorf("listed as refused (%s) but the plan applies: correct the entry", gap.refusal)
			case got != want && !listed:
				t.Errorf("stored schema differs\nplan:\n%s\n--- got\n%s\n--- want\n%s", plan.SQL(), got, want)
			case got == want && listed:
				t.Errorf("listed as a gap (%s) but converges: remove the entry", gap.reason)
			}
		})
	}
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
	for _, seg := range obparser.Split(schema) {
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
