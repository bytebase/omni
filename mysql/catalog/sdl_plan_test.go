package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
)

func TestLoadSDLInDatabase_CallerRelativeLineNumbers(t *testing.T) {
	// Invalid token on line 3 of caller's SDL. Pre-creating the database directly
	// on the catalog avoids a synthetic SQL prologue, so line numbers stay exact.
	target := "CREATE TABLE `t` (\n" +
		"  `id` int NOT NULL,\n" +
		"  BOGUS TOKEN HERE,\n" +
		"  PRIMARY KEY (`id`)\n" +
		");\n"
	_, err := LoadSDLInDatabase(target, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "(line 3, column")
	require.NotContains(t, err.Error(), "(line 5, column")

	_, err = BuildSDLMigrationPlan("", target, SDLMigrationOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to load target schema")
	require.Contains(t, err.Error(), "(line 3, column")

	_, err = BuildSDLMigrationPlan(target, "", SDLMigrationOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to load source schema")
	require.Contains(t, err.Error(), "(line 3, column")
}

func TestLoadSDLInDatabase_FallbackBareTimestamp(t *testing.T) {
	const table = "CREATE TABLE `t` (`ts` timestamp) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n"
	// LOCK TABLES is rejected by LoadSDL, exercising the LoadSQL fallback.
	target := "LOCK TABLES `t` WRITE;\n" + table + "UNLOCK TABLES;\n"

	plan, err := BuildSDLMigrationPlan(table, target, SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan, "LoadSQL fallback must match LoadSDL bare TIMESTAMP canonicalization")
}

func TestLoadSDLInDatabase_EmptyAndCustomDatabase(t *testing.T) {
	c, err := LoadSDLInDatabase("", "custom_db")
	require.NoError(t, err)
	require.NotNil(t, c.GetDatabase("custom_db"))
	require.Equal(t, "custom_db", c.CurrentDatabase())

	c2, err := LoadSDLInDatabase("CREATE TABLE `t` (`id` int PRIMARY KEY);", "custom_db")
	require.NoError(t, err)
	db := c2.GetDatabase("custom_db")
	require.NotNil(t, db)
	require.NotNil(t, db.GetTable("t"))
}

func TestLoadSDLInDatabase_ExplicitDatabaseSwitch(t *testing.T) {
	sql := "CREATE TABLE `t1` (`id` int PRIMARY KEY);\n" +
		"CREATE DATABASE `other`;\n" +
		"USE `other`;\n" +
		"CREATE TABLE `t2` (`id` int PRIMARY KEY);\n"
	c, err := LoadSDLInDatabase(sql, DefaultSDLDatabase)
	require.NoError(t, err)
	require.NotNil(t, c.GetDatabase(DefaultSDLDatabase).GetTable("t1"))
	require.NotNil(t, c.GetDatabase("other").GetTable("t2"))
}

func TestBuildSDLMigrationPlan_EmptyTargetDropsEverything(t *testing.T) {
	source := "CREATE TABLE t (id INT PRIMARY KEY);\n"
	plan, err := BuildSDLMigrationPlan(source, "", SDLMigrationOptions{})
	require.NoError(t, err)
	require.NotNil(t, plan)
	sql := StripDatabaseContext(plan.SQL(), "")
	require.Contains(t, sql, "DROP TABLE")
	require.NotContains(t, sql, DefaultSDLDatabase)
}

func TestBuildSDLMigrationPlan_RejectsCTAS(t *testing.T) {
	source := "CREATE TABLE `t` (`id` int PRIMARY KEY) ENGINE=InnoDB;\n"

	// Bare CTAS must be rejected rather than falling back to LoadSQL and silently
	// omitting `t` (which would emit DROP TABLE `t`).
	plan, err := BuildSDLMigrationPlan(source, "CREATE TABLE t AS SELECT 1 AS id;\n", SDLMigrationOptions{})
	require.Error(t, err)
	require.ErrorContains(t, err, "CREATE TABLE ... AS SELECT")
	require.Nil(t, plan)

	// CTAS mixed with non-SDL statements (LOCK TABLES) must also be rejected.
	mixed := "LOCK TABLES `t` WRITE;\nCREATE TABLE t AS SELECT 1 AS id;\nUNLOCK TABLES;\n"
	plan, err = BuildSDLMigrationPlan(source, mixed, SDLMigrationOptions{})
	require.Error(t, err)
	require.ErrorContains(t, err, "CREATE TABLE ... AS SELECT")
	require.Nil(t, plan)
}

func TestBuildSDLMigrationPlan_NoOpReturnsNil(t *testing.T) {
	sdl := "CREATE TABLE `t` (`id` int NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n"
	plan, err := BuildSDLMigrationPlan(sdl, sdl, SDLMigrationOptions{})
	require.NoError(t, err)
	require.Nil(t, plan)
}

func TestBuildSDLMigrationPlan_SourceSessionContext(t *testing.T) {
	source := "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1;\n"
	target := "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 2;\n"

	ctxMap := &SessionContextMap{
		Functions: map[string]SessionContext{
			"f": {
				SQLMode:             "PIPES_AS_CONCAT,NO_BACKSLASH_ESCAPES",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_0900_ai_ci",
			},
		},
	}

	plan, err := BuildSDLMigrationPlan(source, target, SDLMigrationOptions{
		SourceSessionContext: ctxMap,
	})
	require.NoError(t, err)
	require.NotNil(t, plan)
	sql := StripDatabaseContext(plan.SQL(), "")
	require.Contains(t, sql, "SET sql_mode = 'PIPES_AS_CONCAT,NO_BACKSLASH_ESCAPES'")
	require.NotContains(t, sql, DefaultSDLDatabase)

	// Identical source and target with session context must still produce nil plan.
	noop, err := BuildSDLMigrationPlan(source, source, SDLMigrationOptions{
		SourceSessionContext: ctxMap,
	})
	require.NoError(t, err)
	require.Nil(t, noop)

	// Brand-new routine (absent from source) must be emitted bare.
	createPlan, err := BuildSDLMigrationPlan("", target, SDLMigrationOptions{
		SourceSessionContext: ctxMap,
	})
	require.NoError(t, err)
	require.NotNil(t, createPlan)
	require.NotContains(t, createPlan.SQL(), "SET sql_mode")
}

func TestSessionContextMapFromMetadata(t *testing.T) {
	require.Nil(t, SessionContextMapFromMetadata(nil))
	require.Nil(t, SessionContextMapFromMetadata(&metadata.DatabaseSchemaMetadata{}))

	meta := &metadata.DatabaseSchemaMetadata{
		Name: "app",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{{
				Name: "t",
				Triggers: []*metadata.TriggerMetadata{{
					Name:                "Trg_Audit",
					SqlMode:             "PIPES_AS_CONCAT",
					CharacterSetClient:  "utf8mb4",
					CollationConnection: "utf8mb4_0900_ai_ci",
				}},
			}},
			Functions: []*metadata.FunctionMetadata{{
				Name:                "Fn_Calc",
				SqlMode:             "STRICT_TRANS_TABLES",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_0900_ai_ci",
			}},
			Procedures: []*metadata.ProcedureMetadata{{
				Name:                "Proc_Run",
				SqlMode:             "NO_BACKSLASH_ESCAPES",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_0900_ai_ci",
			}},
			Events: []*metadata.EventMetadata{{
				Name:                "Evt_Daily",
				SqlMode:             "PIPES_AS_CONCAT",
				CharacterSetClient:  "utf8mb4",
				CollationConnection: "utf8mb4_0900_ai_ci",
				TimeZone:            "+08:00",
			}},
		}},
	}

	got := SessionContextMapFromMetadata(meta)
	require.NotNil(t, got)
	require.Equal(t, "STRICT_TRANS_TABLES", got.Functions["fn_calc"].SQLMode)
	require.Equal(t, "NO_BACKSLASH_ESCAPES", got.Procedures["proc_run"].SQLMode)
	require.Equal(t, "PIPES_AS_CONCAT", got.Triggers["trg_audit"].SQLMode)
	require.Equal(t, "+08:00", got.Events["evt_daily"].TimeZone)
}

func TestStripDatabaseContext(t *testing.T) {
	cases := []struct {
		name   string
		dbName string
		in     string
		want   string
	}{
		{
			name: "plain_qualifiers_stripped",
			in:   "ALTER TABLE `bbcatalog`.`t` ADD COLUMN `c` int;\nCREATE INDEX `i` ON `bbcatalog`.`t` (`c`)",
			want: "ALTER TABLE `t` ADD COLUMN `c` int;\nCREATE INDEX `i` ON `t` (`c`)",
		},
		{
			name: "single_quoted_literal_preserved",
			in:   "CREATE OR REPLACE VIEW `bbcatalog`.`v` AS select '`bbcatalog`.`x`' AS c from `bbcatalog`.`t`",
			want: "CREATE OR REPLACE VIEW `v` AS select '`bbcatalog`.`x`' AS c from `t`",
		},
		{
			name: "double_quoted_literal_preserved",
			in:   `select "` + "`bbcatalog`.`x`" + `" from ` + "`bbcatalog`.`t`",
			want: `select "` + "`bbcatalog`.`x`" + `" from ` + "`t`",
		},
		{
			name: "table_actually_named_bbcatalog",
			in:   "ALTER TABLE `bbcatalog`.`bbcatalog` ADD COLUMN `c` int",
			want: "ALTER TABLE `bbcatalog` ADD COLUMN `c` int",
		},
		{
			name: "three_part_reference",
			in:   "select `bbcatalog`.`t`.`c` from `bbcatalog`.`t`",
			want: "select `t`.`c` from `t`",
		},
		{
			name: "routine_body_comment_with_apostrophe_does_not_swallow_later_qualifier",
			in:   "CREATE FUNCTION `bbcatalog`.`f`() RETURNS int\nBEGIN\n  -- don't do this\n  RETURN 1;\nEND;\nALTER TABLE `bbcatalog`.`t` ADD COLUMN `c` int",
			want: "CREATE FUNCTION `f`() RETURNS int\nBEGIN\n  -- don't do this\n  RETURN 1;\nEND;\nALTER TABLE `t` ADD COLUMN `c` int",
		},
		{
			name: "hash_comment_opaque_and_executable_comment_scanned",
			in:   "CREATE TABLE `bbcatalog`.`t` (\n  `pt` point NOT NULL /*!80003 SRID 0 */\n); # trailing 'note\n/* block 'c' */ ALTER TABLE `bbcatalog`.`t` COMMENT ''",
			want: "CREATE TABLE `t` (\n  `pt` point NOT NULL /*!80003 SRID 0 */\n); # trailing 'note\n/* block 'c' */ ALTER TABLE `t` COMMENT ''",
		},
		{
			name: "no_qualifier_untouched",
			in:   "ALTER TABLE `t` ADD COLUMN `c` int DEFAULT 'bbcatalog'",
			want: "ALTER TABLE `t` ADD COLUMN `c` int DEFAULT 'bbcatalog'",
		},
		{
			name:   "custom_database_name",
			dbName: "custom_db",
			in:     "ALTER TABLE `custom_db`.`t` ADD COLUMN `c` int",
			want:   "ALTER TABLE `t` ADD COLUMN `c` int",
		},
		{
			name: "table_alias_named_bbcatalog_preserved_in_view",
			in:   "CREATE OR REPLACE VIEW `bbcatalog`.`v` AS select `bbcatalog`.`id` AS `id`,`u`.`id` AS `uid` from (`bbcatalog`.`t1` `bbcatalog` join `bbcatalog`.`t2` `u` on((`bbcatalog`.`id` = `u`.`id`)))",
			want: "CREATE OR REPLACE VIEW `v` AS select `bbcatalog`.`id` AS `id`,`u`.`id` AS `uid` from (`t1` `bbcatalog` join `t2` `u` on((`bbcatalog`.`id` = `u`.`id`)))",
		},
		{
			name: "table_named_bbcatalog_column_ref_preserved_and_3part_stripped",
			in:   "select `bbcatalog`.`id`, `bbcatalog`.`bbcatalog`.`name` from `bbcatalog`.`bbcatalog`",
			want: "select `bbcatalog`.`id`, `bbcatalog`.`name` from `bbcatalog`",
		},
		{
			name: "trigger_and_event_qualifiers_stripped_while_body_alias_preserved",
			in:   "CREATE TRIGGER `bbcatalog`.`trg` BEFORE INSERT ON `bbcatalog`.`t` FOR EACH ROW BEGIN UPDATE `bbcatalog`.`t` AS `bbcatalog` SET `bbcatalog`.`c` = 1; END;\nCREATE EVENT `bbcatalog`.`ev` ON SCHEDULE EVERY 1 DAY DO DELETE FROM `bbcatalog`.`t`;\nALTER EVENT `bbcatalog`.`ev` RENAME TO `bbcatalog`.`ev2`;\nDROP EVENT IF EXISTS `bbcatalog`.`ev2`",
			want: "CREATE TRIGGER `trg` BEFORE INSERT ON `t` FOR EACH ROW BEGIN UPDATE `t` AS `bbcatalog` SET `bbcatalog`.`c` = 1; END;\nCREATE EVENT `ev` ON SCHEDULE EVERY 1 DAY DO DELETE FROM `t`;\nALTER EVENT `ev` RENAME TO `ev2`;\nDROP EVENT IF EXISTS `ev2`",
		},
		{
			name: "function_call_qualifier_stripped",
			in:   "SELECT `bbcatalog`.`fn`(1)",
			want: "SELECT `fn`(1)",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, StripDatabaseContext(tc.in, tc.dbName))
		})
	}
}

func TestBuildSDLMigrationPlan_OfflineFixtures(t *testing.T) {
	readFixture := func(name string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join("testdata", "sdl_migration", name))
		require.NoError(t, err)
		return string(b)
	}

	selfIdempotent := []string{
		"representative_ddl.sql",
		"big_schema_user_sdl_core.sql",
		"multi_change_base.sql",
		"multi_change_target_80.sql",
		"drop_heavy_base.sql",
		"drop_heavy_target.sql",
		"diverge_schema_with_check.sql",
		"chain_s0.sql",
		"chain_s1_common.sql",
		"chain_s2_common.sql",
		"chain_s3_common.sql",
		"chain_s4_base.sql",
		"roundcube.sql",
		"employees.sql",
	}
	for _, name := range selfIdempotent {
		t.Run("self_idempotent/"+strings.TrimSuffix(name, ".sql"), func(t *testing.T) {
			sdl := readFixture(name)
			plan, err := BuildSDLMigrationPlan(sdl, sdl, SDLMigrationOptions{})
			require.NoError(t, err)
			require.Nil(t, plan, "self-diff of %s must be empty", name)
		})
	}

	t.Run("multi_change_plan_non_empty", func(t *testing.T) {
		plan, err := BuildSDLMigrationPlan(readFixture("multi_change_base.sql"), readFixture("multi_change_target_80.sql"), SDLMigrationOptions{})
		require.NoError(t, err)
		require.NotNil(t, plan)
		sql := StripDatabaseContext(plan.SQL(), "")
		require.NotEmpty(t, sql)
		require.NotContains(t, sql, DefaultSDLDatabase)
	})

	t.Run("drop_heavy_plan_contains_drops", func(t *testing.T) {
		plan, err := BuildSDLMigrationPlan(readFixture("drop_heavy_base.sql"), readFixture("drop_heavy_target.sql"), SDLMigrationOptions{})
		require.NoError(t, err)
		require.NotNil(t, plan)
		var hasDropTable, hasDropCol bool
		for _, op := range plan.Ops {
			if op.Type == OpDropTable {
				hasDropTable = true
			}
			if op.Type == OpDropColumn {
				hasDropCol = true
			}
		}
		require.True(t, hasDropTable, "expected OpDropTable in drop_heavy plan")
		require.True(t, hasDropCol, "expected OpDropColumn in drop_heavy plan")
	})
}
