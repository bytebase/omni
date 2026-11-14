package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ALTER ICEBERG TABLE (https://docs.snowflake.com/en/sql-reference/sql/alter-iceberg-table
// and its REFRESH, CONVERT TO MANAGED, OPTIMIZE, partition-evolution and
// version-tag pages) parses to an AlterTableStmt with Iceberg set.
func TestAlterIcebergTable(t *testing.T) {
	stmt := requireAlterTable(t, "ALTER ICEBERG TABLE IF EXISTS db.s.t ALTER COLUMN c SET DATA TYPE INT RENAME FIELDS")
	if !stmt.Iceberg || !stmt.IfExists || stmt.Loc.Start != 6 { // ALTER statements start at the object keyword
		t.Errorf("stmt = %+v", stmt)
	}
	if a := requireAction(t, stmt, 0); a.Kind != ast.AlterTableAlterColumn || !a.ColumnAlters[0].RenameFields {
		t.Errorf("action = %+v", a)
	}

	cases := map[string]ast.AlterTableActionKind{
		"ALTER ICEBERG TABLE t REFRESH 'metadata/v2.json'":                                             ast.AlterTableRefresh,
		"ALTER ICEBERG TABLE t CONVERT TO MANAGED BASE_LOCATION = 'p' DROP_FROM_SOURCE_CATALOG = TRUE": ast.AlterTableConvertToManaged,
		"ALTER ICEBERG TABLE t OPTIMIZE (REMOVE_POSITION_DELETES)":                                     ast.AlterTableOptimize,
		"ALTER ICEBERG TABLE t ADD PARTITION BY (BUCKET(10, c))":                                       ast.AlterTableAddPartitionBy,
		"ALTER ICEBERG TABLE t DROP PARTITION BY (c)":                                                  ast.AlterTableDropPartitionBy,
		"ALTER ICEBERG TABLE t REPLACE PARTITION BY (DAY(c)) WITH (HOUR(c))":                           ast.AlterTableReplacePartitionBy,
		"ALTER ICEBERG TABLE t CREATE OR REPLACE VERSION_TAG v AS OF VERSION 12 RETAIN 3 DAYS":         ast.AlterTableCreateVersionTag,
		"ALTER ICEBERG TABLE t CREATE VERSION_TAG IF NOT EXISTS v AS OF TIMESTAMP 1700000000":          ast.AlterTableCreateVersionTag,
		"ALTER ICEBERG TABLE t ALTER VERSION_TAG v AS OF VERSION 3":                                    ast.AlterTableAlterVersionTag,
		"ALTER ICEBERG TABLE t REPLACE VERSION_TAG v AS OF VERSION 3 RETAIN 1 HOURS":                   ast.AlterTableReplaceVersionTag,
		"ALTER ICEBERG TABLE t DROP VERSION_TAG IF EXISTS v":                                           ast.AlterTableDropVersionTag,
		"ALTER ICEBERG TABLE t UPGRADE TO ICEBERG VERSION 3":                                           ast.AlterTableUpgradeIcebergVersion,
		"ALTER ICEBERG TABLE t SET TARGET_FILE_SIZE = '64MB' CONTACT (steward = alice)":                ast.AlterTableSet,
		"ALTER ICEBERG TABLE t UNSET REPLACE_INVALID_CHARACTERS, PATH_LAYOUT":                          ast.AlterTableUnset,
		"ALTER ICEBERG TABLE t SUSPEND RECLUSTER":                                                      ast.AlterTableSuspendRecluster,
		"ALTER ICEBERG TABLE t DROP version_tag":                                                       ast.AlterTableDropColumn,
		"ALTER ICEBERG TABLE t ADD partition INT":                                                      ast.AlterTableAddColumn,
	}
	for sql, kind := range cases {
		stmt := requireAlterTable(t, sql)
		if a := requireAction(t, stmt, 0); a.Kind != kind || !stmt.Iceberg {
			t.Errorf("%q: kind = %v, want %v", sql, a.Kind, kind)
		}
	}

	stmt = requireAlterTable(t, "ALTER ICEBERG TABLE t ALTER c SET WRITE DEFAULT 1, d DROP WRITE DEFAULT")
	if cas := requireAction(t, stmt, 0).ColumnAlters; len(cas) != 2 ||
		cas[0].Kind != ast.ColumnAlterSetWriteDefault || cas[1].Kind != ast.ColumnAlterDropWriteDefault {
		t.Errorf("column alters = %+v", cas)
	}
	stmt = requireAlterTable(t, "ALTER ICEBERG TABLE t ADD PARTITION BY (TRUNCATE(4, c))")
	if call, ok := requireAction(t, stmt, 0).Partition.(*ast.FuncCallExpr); !ok || call.Name.Name.Name != "TRUNCATE" || len(call.Args) != 2 {
		t.Errorf("partition = %+v", requireAction(t, stmt, 0).Partition)
	}

	for _, sql := range []string{
		// Forms the ALTER ICEBERG TABLE page does not have.
		"ALTER ICEBERG TABLE t RENAME TO u",
		"ALTER ICEBERG TABLE t SWAP WITH u",
		"ALTER ICEBERG TABLE t RECLUSTER",
		"ALTER ICEBERG TABLE t ADD DATA METRIC FUNCTION m ON (a)",
		"ALTER ICEBERG TABLE t ADD PRIMARY KEY (a)",
		// Exactly one partition expression; known transforms only.
		"ALTER ICEBERG TABLE t ADD PARTITION BY (a, b)",
		"ALTER ICEBERG TABLE t ADD PARTITION BY (WEEK(a))",
		"ALTER ICEBERG TABLE t OPTIMIZE (COMPACT)",
		"ALTER ICEBERG TABLE t SET TARGET_FILE_SIZE = 16 MB",
		"ALTER ICEBERG TABLE t SET PATH_LAYOUT = 'FLAT'",
		"ALTER ICEBERG TABLE t CONVERT TO MANAGED NOT_A_PROPERTY = 12",
		"ALTER ICEBERG TABLE t CONVERT TO MANAGED ERROR_LOGGING = TRUE",
		// Plain-table properties ALTER ICEBERG TABLE's SET does not take.
		"ALTER ICEBERG TABLE t SET STAGE_FILE_FORMAT = (TYPE = CSV)",
		"ALTER ICEBERG TABLE t SET STAGE_COPY_OPTIONS = (ON_ERROR = CONTINUE)",
		"ALTER ICEBERG TABLE t SET CHANGE_TRACKING = TRUE",
		"ALTER ICEBERG TABLE t SET DATA_RETENTION_TIME_IN_DAYS = 1 COMMENT = 'c'",
		"ALTER ICEBERG TABLE t SET LOG_EVENT_LEVEL = TRACE",
		"ALTER ICEBERG TABLE t SET TARGET_FILE_SIZE = '20MB'",
		"ALTER ICEBERG TABLE t CONVERT TO MANAGED DROP_FROM_SOURCE_CATALOG = 'yes'",
		// Iceberg-only forms on a plain table.
		"ALTER TABLE t REFRESH",
		"ALTER TABLE t UPGRADE TO ICEBERG VERSION 3",
		"ALTER TABLE t ALTER c SET WRITE DEFAULT 1",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
}
