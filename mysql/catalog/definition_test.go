package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

func TestTableDefinitionWithCheckConstraint(t *testing.T) {
	table := &metadata.TableMetadata{
		Name: "t1",
		Columns: []*metadata.ColumnMetadata{
			{Name: "type", Type: "varchar(10)", Nullable: true, Default: "NULL"},
			{Name: "amount", Type: "decimal(10,2)", Nullable: true, Default: "NULL"},
		},
		CheckConstraints: []*metadata.CheckConstraintMetadata{
			{
				Name: "c1",
				// The expression as synced from MySQL after unescaping the
				// information_schema escaped form ((`type` = _utf8mb4\'A\')...).
				Expression: "(((`type` = _utf8mb4'A') and (`amount` is not null)))",
			},
		},
		Engine:    "InnoDB",
		Charset:   "utf8mb4",
		Collation: "utf8mb4_general_ci",
	}

	got, err := TableDefinition(table)
	require.NoError(t, err)

	want := "--\n" +
		"-- Table structure for `t1`\n" +
		"--\n" +
		"CREATE TABLE `t1` (\n" +
		"  `type` varchar(10) DEFAULT NULL,\n" +
		"  `amount` decimal(10,2) DEFAULT NULL,\n" +
		"  CONSTRAINT `c1` CHECK (((`type` = _utf8mb4'A') and (`amount` is not null)))\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;\n\n"
	require.Equal(t, want, got)
	require.NotContains(t, got, `\'`)
}

func TestTableDefinitionWithSystemVersioning(t *testing.T) {
	table := &metadata.TableMetadata{
		Name: "explicit_period",
		Columns: []*metadata.ColumnMetadata{
			{Name: "id", Type: "int", Nullable: true, Default: "NULL"},
			{Name: "valid_from", Type: "timestamp(6)", Generation: &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_STORED, Expression: "ROW START"}},
			{Name: "valid_to", Type: "timestamp(6)", Generation: &metadata.GenerationMetadata{Type: metadata.GenerationMetadata_TYPE_STORED, Expression: "ROW END"}},
		},
		Engine:        "InnoDB",
		Charset:       "utf8mb4",
		Collation:     "utf8mb4_general_ci",
		Comment:       "versioned",
		CreateOptions: "WITH SYSTEM VERSIONING",
	}

	got, err := TableDefinition(table)
	require.NoError(t, err)

	want := "--\n" +
		"-- Table structure for `explicit_period`\n" +
		"--\n" +
		"CREATE TABLE `explicit_period` (\n" +
		"  `id` int DEFAULT NULL,\n" +
		"  `valid_from` timestamp(6) GENERATED ALWAYS AS ROW START NOT NULL,\n" +
		"  `valid_to` timestamp(6) GENERATED ALWAYS AS ROW END NOT NULL,\n" +
		"  PERIOD FOR SYSTEM_TIME (`valid_from`, `valid_to`)\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='versioned' WITH SYSTEM VERSIONING;\n\n"
	require.Equal(t, want, got)

	// The schema-text diff path loads a dump with this catalog. Loading the table
	// unversioned would let the diff recreate it without history, so it must fail.
	c := New()
	_, err = c.Exec("CREATE DATABASE tmp; USE tmp;", nil)
	require.NoError(t, err)
	results, err := c.Exec(got, &ExecOptions{ContinueOnError: true})
	require.True(t, err != nil || slices.ContainsFunc(results, func(r ExecResult) bool { return r.Error != nil }))
}

func TestFunctionDefinitionWithMultilineParameters(t *testing.T) {
	function := &metadata.FunctionMetadata{
		Name: "f1",
		// The definition as synced from SHOW CREATE FUNCTION: the parameter
		// list keeps its original multi-line formatting, and the CHARSET and
		// COLLATE attributes of the RETURNS clause have been stripped.
		Definition: "CREATE FUNCTION `f1`(\n" +
			"    p_a BIGINT,\n" +
			"    p_b VARCHAR(36)\n" +
			") RETURNS char(36)\n" +
			"    NO SQL\n" +
			"    DETERMINISTIC\n" +
			"BEGIN\n" +
			"    RETURN p_b;\n" +
			"END",
		CharacterSetClient:  "utf8mb4",
		CollationConnection: "utf8mb4_unicode_ci",
		SqlMode:             "ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES",
	}

	got, err := FunctionDefinition(function)
	require.NoError(t, err)

	want := "--\n" +
		"-- Function structure for `f1`\n" +
		"--\n" +
		"SET character_set_client = utf8mb4;\n" +
		"SET character_set_results = utf8mb4;\n" +
		"SET collation_connection = utf8mb4_unicode_ci;\n" +
		"SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES';\n" +
		"DELIMITER ;;\n" +
		"CREATE FUNCTION `f1`(\n" +
		"    p_a BIGINT,\n" +
		"    p_b VARCHAR(36)\n" +
		") RETURNS char(36)\n" +
		"    NO SQL\n" +
		"    DETERMINISTIC\n" +
		"BEGIN\n" +
		"    RETURN p_b;\n" +
		"END;;\n" +
		"DELIMITER ;\n\n"
	require.Equal(t, want, got)
}

func TestDatabaseDefinitionSkipsEventWithoutDefinition(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{
		Name:    "db",
		Schemas: []*metadata.SchemaMetadata{{Events: []*metadata.EventMetadata{{Name: "e"}}}},
	}
	for _, sdl := range []bool{false, true} {
		definition, err := DatabaseDefinition(meta, DefinitionOptions{SDLFormat: sdl})
		require.NoError(t, err)
		require.Empty(t, definition, "SDLFormat=%t", sdl)
	}
}

func TestRenderColumnDefaultBitLiteral(t *testing.T) {
	testCases := []struct {
		name       string
		columnType string
		// def is the value stored in ColumnMetadata.Default. For a BIT column the MySQL
		// sync stores the bit literal QUOTE()-escaped, e.g. b'0' becomes 'b\'0\''.
		def  string
		want string
	}{
		{name: "bit1 quoted-escaped", columnType: "bit(1)", def: `'b\'0\''`, want: "b'0'"},
		{name: "bit8 quoted-escaped", columnType: "bit(8)", def: `'b\'101\''`, want: "b'101'"},
		{name: "bit no width", columnType: "bit", def: `'b\'1\''`, want: "b'1'"},
		{name: "bit already bare literal", columnType: "bit(4)", def: "b'1010'", want: "b'1010'"},
		{name: "bit hex 0x form", columnType: "bit(8)", def: "0xFF", want: "0xFF"},
		{name: "bit hex x'' quoted-escaped", columnType: "bit(8)", def: `'x\'1f\''`, want: "x'1f'"},
		// Non-bit columns and non-literal defaults are emitted verbatim.
		{name: "varchar string default untouched", columnType: "varchar(10)", def: "'hello'", want: "'hello'"},
		{name: "varchar string that looks bit-ish but not bit column", columnType: "varchar(10)", def: `'b\'0\''`, want: `'b\'0\''`},
		{name: "int numeric default untouched", columnType: "int", def: "42", want: "42"},
		// A BIT column whose default is not a bit/hex literal (defensive) is left verbatim.
		{name: "bit with non-literal default", columnType: "bit(8)", def: "'abc'", want: "'abc'"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderColumnDefault(&metadata.ColumnMetadata{Type: tc.columnType, Default: tc.def})
			require.Equal(t, tc.want, got)
		})
	}
}

func TestPrintColumnClauseBitDefaultUnquoted(t *testing.T) {
	table := &metadata.TableMetadata{Name: "t", Charset: "utf8mb4"}
	var buf strings.Builder
	col := &metadata.ColumnMetadata{Name: "b1", Type: "bit(1)", Nullable: false, Default: `'b\'0\''`}
	require.NoError(t, printColumnClause(&buf, col, table))
	require.Equal(t, "  `b1` bit(1) NOT NULL DEFAULT b'0'", buf.String())

	// A BIT column with no default (b64 BIT(64)) must stay fine: NULL default emitted as
	// DEFAULT NULL, never as a bit literal.
	var buf2 strings.Builder
	col2 := &metadata.ColumnMetadata{Name: "b64", Type: "bit(64)", Nullable: true, Default: "NULL"}
	require.NoError(t, printColumnClause(&buf2, col2, table))
	require.Equal(t, "  `b64` bit(64) DEFAULT NULL", buf2.String())
}

func TestNormalizeFunctionalIndexExpr(t *testing.T) {
	testCases := []struct {
		name string
		expr string
		want string
	}{
		{
			name: "multi-valued json index with utf8mb4 introducer and escaped quotes",
			// The form synced from information_schema.STATISTICS.EXPRESSION.
			expr: "(cast(json_extract(`tags`,_utf8mb4\\'$.ids\\') as unsigned array))",
			want: "(cast(json_extract(`tags`,'$.ids') as unsigned array))",
		},
		{
			name: "latin1 introducer",
			expr: "(json_extract(`c`,_latin1\\'$.a\\'))",
			want: "(json_extract(`c`,'$.a'))",
		},
		{
			name: "no introducer no escaping is unchanged",
			expr: "(json_extract(`tags`,'$.ids'))",
			want: "(json_extract(`tags`,'$.ids'))",
		},
		{
			name: "introducer-like underscore inside string literal is preserved",
			// _utf8mb4 appears INSIDE the literal here; only a real introducer (before the
			// opening quote) must be stripped.
			expr: "(concat(`c`,_utf8mb4\\'_utf8mb4 is text\\'))",
			want: "(concat(`c`,'_utf8mb4 is text'))",
		},
		{
			name: "leading-underscore identifier not before a quote is kept",
			expr: "(`_weird`)",
			want: "(`_weird`)",
		},
		{
			name: "backticked identifier with underscore is untouched",
			expr: "(json_extract(`my_col`,_utf8mb4\\'$.k\\'))",
			want: "(json_extract(`my_col`,'$.k'))",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, normalizeFunctionalIndexExpr(tc.expr))
		})
	}
}

func TestPrintIndexClauseFunctionalExprNormalized(t *testing.T) {
	var buf strings.Builder
	index := &metadata.IndexMetadata{
		Name:        "idx_tags",
		Visible:     true,
		Expressions: []string{"(cast(json_extract(`tags`,_utf8mb4\\'$.ids\\') as unsigned array))"},
	}
	require.NoError(t, printIndexClause(&buf, index, "InnoDB"))
	require.Equal(t, ",\n  KEY `idx_tags` ((cast(json_extract(`tags`,'$.ids') as unsigned array)))", buf.String())
	require.NotContains(t, buf.String(), "_utf8mb4")
	require.NotContains(t, buf.String(), `\'`)
}

// TestStripViewBodyDatabaseQualifier pins the token-scanner rewrite (X6). The failure
// shapes were live-verified: (a) `db`.`db`.`col` (a table named like the database, the
// 5.7 three-part form) was double-stripped down to `col`; (b) a backtick alias
// containing a quote (`it's`) flipped the old scanner's literal state so NO subsequent
// qualifier was stripped; (c) a genuine string literal containing the qualifier bytes
// was corrupted.
func TestStripViewBodyDatabaseQualifier(t *testing.T) {
	cases := []struct {
		name   string
		body   string
		dbName string
		want   string
	}{
		{
			name:   "table_named_like_database_keeps_table_qualifier",
			body:   "select `db`.`db`.`col` AS `col` from `db`.`db`",
			dbName: "db",
			want:   "select `db`.`col` AS `col` from `db`",
		},
		{
			name: "alias_with_quote_does_not_flip_literal_state",
			// 5.7 derived-table shape: all three qualifiers must strip even though the
			// alias `it's` contains a single quote.
			body:   "select `x`.`a` AS `it's` from (select `mydb`.`t`.`a` AS `a` from `mydb`.`t` where (`mydb`.`t`.`a` > 0)) `x`",
			dbName: "mydb",
			want:   "select `x`.`a` AS `it's` from (select `t`.`a` AS `a` from `t` where (`t`.`a` > 0)) `x`",
		},
		{
			name:   "control_alias_without_quote",
			body:   "select `x`.`a` AS `plain` from (select `mydb`.`t`.`a` AS `a` from `mydb`.`t`) `x`",
			dbName: "mydb",
			want:   "select `x`.`a` AS `plain` from (select `t`.`a` AS `a` from `t`) `x`",
		},
		{
			name:   "cross_database_reference_preserved",
			body:   "select `otherdb`.`t`.`c` from `otherdb`.`t` join `mydb`.`u` on (`otherdb`.`t`.`id` = `mydb`.`u`.`id`)",
			dbName: "mydb",
			want:   "select `otherdb`.`t`.`c` from `otherdb`.`t` join `u` on (`otherdb`.`t`.`id` = `u`.`id`)",
		},
		{
			name: "db_named_segment_inside_foreign_reference_preserved",
			// `mydb` here is the TABLE part of a cross-database reference — a
			// continuation after ".", never a qualifier.
			body:   "select `otherdb`.`mydb`.`c` from `otherdb`.`mydb`",
			dbName: "mydb",
			want:   "select `otherdb`.`mydb`.`c` from `otherdb`.`mydb`",
		},
		{
			name: "quote_flip_then_literal_corruption",
			// The live-verified composed corruption: with the old scanner, the quote
			// inside alias `it's` flipped the literal state, so the REAL literal's
			// interior was treated as top-level (its `mydb`. stripped — data change)
			// while the genuine trailing qualifier survived.
			body:   "select `it's` AS a, '`mydb`.`t` and  spaces' AS lit from `mydb`.`t`",
			dbName: "mydb",
			want:   "select `it's` AS a, '`mydb`.`t` and  spaces' AS lit from `t`",
		},
		{
			name:   "string_literal_containing_qualifier_untouched",
			body:   "select '`mydb`.`t`' AS lit, `mydb`.`t`.`c` from `mydb`.`t`",
			dbName: "mydb",
			want:   "select '`mydb`.`t`' AS lit, `t`.`c` from `t`",
		},
		{
			name:   "double_quoted_literal_untouched",
			body:   `select "` + "`mydb`.`t`" + `" AS lit from ` + "`mydb`.`t`",
			dbName: "mydb",
			want:   `select "` + "`mydb`.`t`" + `" AS lit from ` + "`t`",
		},
		{
			name:   "qualifier_not_followed_by_identifier_preserved",
			body:   "select `mydb`.* from `mydb`.`t`",
			dbName: "mydb",
			want:   "select `mydb`.* from `t`",
		},
		{
			name:   "alias_spelled_like_database_kept",
			body:   "select `db`.`id` AS `id` from (`t` `db` join `u` on((`db`.`id` = `u`.`id`)))",
			dbName: "db",
			want:   "select `db`.`id` AS `id` from (`t` `db` join `u` on((`db`.`id` = `u`.`id`)))",
		},
		{
			name:   "alias_spelled_like_database_with_qualified_tables",
			body:   "select `db`.`id` AS `id` from (`db`.`t` `db` join `db`.`u` on((`db`.`id` = `db`.`u`.`id`)))",
			dbName: "db",
			want:   "select `db`.`id` AS `id` from (`t` `db` join `u` on((`db`.`id` = `u`.`id`)))",
		},
		{
			name:   "function_qualifier_stripped",
			body:   "select `mydb`.`f`(`mydb`.`t`.`a`) AS `x` from `mydb`.`t`",
			dbName: "mydb",
			want:   "select `f`(`t`.`a`) AS `x` from `t`",
		},
		{
			name:   "check_option_suffix",
			body:   "select `mydb`.`t`.`a` AS `a` from `mydb`.`t` WITH CASCADED CHECK OPTION",
			dbName: "mydb",
			want:   "select `t`.`a` AS `a` from `t` WITH CASCADED CHECK OPTION",
		},
		{
			name:   "unparsable_body_unchanged",
			body:   "select `mydb`.`t`.`a` from `mydb`.`t` where",
			dbName: "mydb",
			want:   "select `mydb`.`t`.`a` from `mydb`.`t` where",
		},
		{
			name:   "empty_db_name_no_op",
			body:   "select `t`.`c` from `t`",
			dbName: "",
			want:   "select `t`.`c` from `t`",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, stripViewBodyDatabaseQualifier(tc.body, tc.dbName))
		})
	}
}

// TestStripLeadingDefiner pins the definer-clause parser (X7): quoted accounts may
// legally contain spaces (`my user`@`%`), '@', and doubled/escaped quotes — the old
// cut-at-first-space produced corrupt "CREATE user`@`%` …" output.
func TestStripLeadingDefiner(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "plain_account",
			in:   "CREATE DEFINER=`root`@`localhost` EVENT `e` ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
			want: "CREATE EVENT `e` ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
		},
		{
			name: "quoted_user_with_space",
			in:   "CREATE DEFINER=`my user`@`%` EVENT `e` ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
			want: "CREATE EVENT `e` ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
		},
		{
			name: "quoted_host_with_space_function",
			in:   "CREATE DEFINER=`a``b`@`local host` FUNCTION `f`() RETURNS int DETERMINISTIC RETURN 1",
			want: "CREATE FUNCTION `f`() RETURNS int DETERMINISTIC RETURN 1",
		},
		{
			name: "single_quoted_account_procedure",
			in:   "CREATE DEFINER='my user'@'%' PROCEDURE `p`() BEGIN END",
			want: "CREATE PROCEDURE `p`() BEGIN END",
		},
		{
			name: "double_quoted_account",
			in:   `CREATE DEFINER="my user"@"%" EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t`,
			want: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
		},
		{
			name: "unquoted_account",
			in:   "CREATE DEFINER=root@localhost EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
			want: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
		},
		{
			name: "current_user",
			in:   "CREATE DEFINER=CURRENT_USER EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
			want: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
		},
		{
			name: "current_user_parens",
			in:   "CREATE DEFINER=CURRENT_USER() EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
			want: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
		},
		{
			name: "spaces_around_equals",
			in:   "CREATE DEFINER = `my user`@`%` EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
			want: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
		},
		{
			name: "no_definer_unchanged",
			in:   "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
			want: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO DELETE FROM t",
		},
		{
			name: "unterminated_quote_left_unchanged",
			in:   "CREATE DEFINER=`broken EVENT e DO SELECT 1",
			want: "CREATE DEFINER=`broken EVENT e DO SELECT 1",
		},
		{
			name: "not_a_create_statement_unchanged",
			in:   "ALTER DEFINER=`root`@`%` EVENT e COMMENT 'x'",
			want: "ALTER DEFINER=`root`@`%` EVENT e COMMENT 'x'",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, stripLeadingDefiner(tc.in))
		})
	}
}

// TestPrintColumnClauseSRIDPresence pins the presence semantics of the SRID attribute
// (X5): explicit SRID 0 must be emitted (it is a valid spatial reference system,
// distinct from "no SRID"), and an unset SRID emits nothing.
func TestPrintColumnClauseSRIDPresence(t *testing.T) {
	srid := func(v uint32) *uint32 { return &v }
	table := &metadata.TableMetadata{Name: "t"}

	render := func(col *metadata.ColumnMetadata) string {
		var buf strings.Builder
		require.NoError(t, printColumnClause(&buf, col, table))
		return buf.String()
	}

	require.Equal(t, "  `pt` point NOT NULL /*!80003 SRID 0 */",
		render(&metadata.ColumnMetadata{Name: "pt", Type: "point", Nullable: false, Srid: srid(0)}))
	require.Equal(t, "  `pt` point NOT NULL /*!80003 SRID 4326 */",
		render(&metadata.ColumnMetadata{Name: "pt", Type: "point", Nullable: false, Srid: srid(4326)}))
	// Custom SRSs above int32 range must render unmangled.
	require.Equal(t, "  `pt` point NOT NULL /*!80003 SRID 3000000000 */",
		render(&metadata.ColumnMetadata{Name: "pt", Type: "point", Nullable: false, Srid: srid(3000000000)}))
	require.Equal(t, "  `pt` point NOT NULL",
		render(&metadata.ColumnMetadata{Name: "pt", Type: "point", Nullable: false}))
}

// TestPrintColumnClauseInvisibleCommentOrder pins the SDL dumper's INVISIBLE/COMMENT order
// (BYT-9830). MySQL's canonical SHOW CREATE emits INVISIBLE before COMMENT for both regular
// and generated columns (verified against 8.0.32), so the dumped form must match — otherwise
// the dump diverges from SHOW CREATE and from the migration generator's generated-column path.
func TestPrintColumnClauseInvisibleCommentOrder(t *testing.T) {
	srid := func(v uint32) *uint32 { return &v }
	table := &metadata.TableMetadata{Name: "t"}

	render := func(col *metadata.ColumnMetadata) string {
		var buf strings.Builder
		require.NoError(t, printColumnClause(&buf, col, table))
		return buf.String()
	}

	// Regular INVISIBLE + COMMENT column: INVISIBLE precedes COMMENT.
	require.Equal(t, "  `a` int DEFAULT NULL /*!80023 INVISIBLE */ COMMENT 'c'",
		render(&metadata.ColumnMetadata{Name: "a", Type: "int", Nullable: true, Default: "NULL", IsInvisible: true, Comment: "c"}))

	// Generated spatial INVISIBLE + COMMENT column: same canonical INVISIBLE-before-COMMENT order.
	require.Equal(t,
		"  `loc` point GENERATED ALWAYS AS (st_srid(point(`lng`,`lat`),4326)) STORED NOT NULL /*!80003 SRID 4326 */ /*!80023 INVISIBLE */ COMMENT 'geo'",
		render(&metadata.ColumnMetadata{Name: "loc", Type: "point", Nullable: false, Srid: srid(4326), IsInvisible: true, Comment: "geo", Generation: &metadata.GenerationMetadata{
			Type:       metadata.GenerationMetadata_TYPE_STORED,
			Expression: "st_srid(point(`lng`,`lat`),4326)",
		}}))
}

// TestDatabaseDefinitionRestoresSnapshot restores a synced snapshot into empty
// databases twice, from the restore script as the mysql client runs it and from the
// SDL, and compares every restored object with the source.
func TestDatabaseDefinitionRestoresSnapshot(t *testing.T) {
	ctx := context.Background()
	conn := testcontainer.OpenEmpty(t)
	exec := func(t *testing.T, stmts ...string) {
		t.Helper()
		for _, stmt := range stmts {
			_, err := conn.ExecContext(ctx, stmt)
			require.NoError(t, err, stmt)
		}
	}
	const function = "CREATE FUNCTION `f``n`(x int) RETURNS int\n" +
		"    DETERMINISTIC\n" +
		"BEGIN\n" +
		"  DECLARE y int;\n" +
		"  SET y = x + 1;\n" +
		"  RETURN y;\n" +
		"END"
	exec(t,
		"CREATE DATABASE shared",
		"CREATE TABLE shared.parent (id int PRIMARY KEY)",
		"CREATE DATABASE src",
		"USE src",
		"CREATE TABLE `a``b` ("+
			"id int NOT NULL, "+
			"`c``d` varchar(20) COMMENT 'it''s a \\\\ path', "+
			"parent_id int, "+
			"PRIMARY KEY (id), "+
			"KEY `k``1` (`c``d`), "+
			"KEY parent_id (parent_id), "+
			"CONSTRAINT `fk``x` FOREIGN KEY (parent_id) REFERENCES shared.parent (id), "+
			"CONSTRAINT `ck``1` CHECK (id > 0)"+
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8 COMMENT 'table''s \\\\ comment'",
		"CREATE TABLE u (id int PRIMARY KEY)",
		"CREATE TABLE mem (a int NOT NULL, b int, PRIMARY KEY (a) USING BTREE, KEY k_hash (b), KEY k_btree (a, b) USING BTREE) ENGINE=MEMORY",
		"CREATE TRIGGER `tr``1` BEFORE INSERT ON `a``b` FOR EACH ROW SET NEW.`c``d` = 'x'",
		function,
		"CREATE VIEW v_fn AS SELECT `f``n`(1) AS n",
		// The alias is spelled like the database.
		"CREATE VIEW v_alias AS SELECT src.id AS id FROM `a``b` AS src JOIN u ON src.id = u.id",
	)

	// The sync reads comments through QUOTE() and strips the outer quotes.
	queryString := func(t *testing.T, query string, args ...any) string {
		t.Helper()
		var value string
		require.NoError(t, conn.QueryRowContext(ctx, query, args...).Scan(&value), query)
		return value
	}
	unquote := func(s string) string { return s[1 : len(s)-1] }
	columnComment := unquote(queryString(t, "SELECT QUOTE(COLUMN_COMMENT) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = 'a`b' AND COLUMN_NAME = 'c`d'"))
	tableComment := unquote(queryString(t, "SELECT QUOTE(TABLE_COMMENT) FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = 'a`b'"))
	viewDefinition := func(t *testing.T, view string) string {
		return queryString(t, "SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = 'src' AND TABLE_NAME = ?", view)
	}
	const sqlMode = "ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION"
	intColumn := func(name string, nullable bool) *metadata.ColumnMetadata {
		column := &metadata.ColumnMetadata{Name: name, Type: "int", Nullable: nullable}
		if nullable {
			column.Default = "NULL"
		}
		return column
	}
	meta := &metadata.DatabaseSchemaMetadata{
		Name: "src",
		Schemas: []*metadata.SchemaMetadata{{
			Tables: []*metadata.TableMetadata{
				{
					Name: "a`b",
					Columns: []*metadata.ColumnMetadata{
						intColumn("id", false),
						{Name: "c`d", Type: "varchar(20)", Nullable: true, Default: "NULL", CharacterSet: "utf8mb4", Collation: "utf8mb4_0900_ai_ci", Comment: columnComment},
						intColumn("parent_id", true),
					},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Visible: true, Expressions: []string{"id"}},
						{Name: "k`1", Type: "BTREE", Visible: true, Expressions: []string{"c`d"}},
						{Name: "parent_id", Type: "BTREE", Visible: true, Expressions: []string{"parent_id"}},
					},
					ForeignKeys: []*metadata.ForeignKeyMetadata{{
						Name: "fk`x", Columns: []string{"parent_id"},
						ReferencedSchema: "shared", ReferencedTable: "parent", ReferencedColumns: []string{"id"},
						OnDelete: "NO ACTION", OnUpdate: "NO ACTION", MatchType: "NONE",
					}},
					CheckConstraints: []*metadata.CheckConstraintMetadata{{Name: "ck`1", Expression: "(`id` > 0)"}},
					Triggers: []*metadata.TriggerMetadata{{
						Name: "tr`1", Timing: "BEFORE", Event: "INSERT", Body: "SET NEW.`c``d` = 'x'",
						CharacterSetClient: "utf8mb4", CollationConnection: "utf8mb4_0900_ai_ci", SqlMode: sqlMode,
					}},
					Engine:        "InnoDB",
					Charset:       "utf8mb4",
					Collation:     "utf8mb4_0900_ai_ci",
					CreateOptions: "row_format=COMPRESSED key_block_size=8",
					Comment:       tableComment,
				},
				{
					Name:      "u",
					Columns:   []*metadata.ColumnMetadata{intColumn("id", false)},
					Indexes:   []*metadata.IndexMetadata{{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Visible: true, Expressions: []string{"id"}}},
					Engine:    "InnoDB",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
				},
				{
					Name:    "mem",
					Columns: []*metadata.ColumnMetadata{intColumn("a", false), intColumn("b", true)},
					Indexes: []*metadata.IndexMetadata{
						{Name: "PRIMARY", Type: "BTREE", Primary: true, Unique: true, Visible: true, Expressions: []string{"a"}},
						{Name: "k_hash", Type: "HASH", Visible: true, Expressions: []string{"b"}},
						{Name: "k_btree", Type: "BTREE", Visible: true, Expressions: []string{"a", "b"}},
					},
					Engine:    "MEMORY",
					Charset:   "utf8mb4",
					Collation: "utf8mb4_0900_ai_ci",
				},
			},
			Views: []*metadata.ViewMetadata{
				{Name: "v_alias", Definition: viewDefinition(t, "v_alias"), Columns: []*metadata.ColumnMetadata{intColumn("id", false)}},
				{Name: "v_fn", Definition: viewDefinition(t, "v_fn"), Columns: []*metadata.ColumnMetadata{intColumn("n", true)}},
			},
			Functions: []*metadata.FunctionMetadata{{
				Name: "f`n", Definition: function,
				CharacterSetClient: "utf8mb4", CollationConnection: "utf8mb4_0900_ai_ci", SqlMode: sqlMode,
			}},
		}},
	}

	script, err := DatabaseDefinition(meta, DefinitionOptions{PrintHeader: true})
	require.NoError(t, err)
	exec(t, "CREATE DATABASE dump", "USE dump")
	exec(t, splitLikeMySQLClient(script)...)

	sdl, err := DatabaseDefinition(meta, DefinitionOptions{SDLFormat: true})
	require.NoError(t, err)
	exec(t, "CREATE DATABASE `sdl`", "USE `sdl`")
	for _, segment := range mysqlparser.Split(sdl) {
		if !segment.Empty() {
			exec(t, segment.Text)
		}
	}

	showCreate := func(t *testing.T, kind, db, name string, column int) string {
		t.Helper()
		rows, err := conn.QueryContext(ctx, fmt.Sprintf("SHOW CREATE %s %s.%s", kind, quoteIdent(db), quoteIdent(name)))
		require.NoError(t, err)
		defer rows.Close()
		columns, err := rows.Columns()
		require.NoError(t, err)
		values := make([]sql.RawBytes, len(columns))
		dest := make([]any, len(values))
		for i := range values {
			dest[i] = &values[i]
		}
		require.True(t, rows.Next(), "SHOW CREATE %s %s.%s", kind, db, name)
		require.NoError(t, rows.Scan(dest...))
		return string(values[column])
	}
	triggers := func(t *testing.T, db string) string {
		return queryString(t, "SELECT GROUP_CONCAT(CONCAT_WS(' ', TRIGGER_NAME, ACTION_TIMING, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT) ORDER BY TRIGGER_NAME) FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = ?", db)
	}
	// The SDL orders keys by name, so its tables are compared line by line.
	keyLines := func(create string) []string {
		lines := strings.Split(create, "\n")
		for i, line := range lines {
			lines[i] = strings.TrimSuffix(line, ",")
		}
		slices.Sort(lines)
		return lines
	}
	for _, db := range []string{"dump", "sdl"} {
		t.Run(db, func(t *testing.T) {
			for _, table := range []string{"a`b", "u", "mem"} {
				want, got := showCreate(t, "TABLE", "src", table, 1), showCreate(t, "TABLE", db, table, 1)
				if db == "sdl" {
					require.Equal(t, keyLines(want), keyLines(got), "table %s", table)
				} else {
					require.Equal(t, want, got, "table %s", table)
				}
			}
			// information_schema qualifies view bodies with their database, and the
			// restore script keeps the source's qualifiers, so both are removed.
			for _, view := range []string{"v_alias", "v_fn"} {
				got := queryString(t, "SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?", db, view)
				want := stripViewBodyDatabaseQualifier(viewDefinition(t, view), "src")
				require.Equal(t, want, stripViewBodyDatabaseQualifier(stripViewBodyDatabaseQualifier(got, "src"), db), "view %s", view)
			}
			require.Equal(t, "2", queryString(t, fmt.Sprintf("SELECT n FROM %s.v_fn", quoteIdent(db))))
			require.Equal(t, showCreate(t, "FUNCTION", "src", "f`n", 2), showCreate(t, "FUNCTION", db, "f`n", 2))
			require.Equal(t, triggers(t, "src"), triggers(t, db))
		})
	}

	// A diff loads the SDL into a database of its own and must read it as equal to
	// itself.
	const sdlDatabase = "CREATE DATABASE `bbcatalog`;\nUSE `bbcatalog`;\n"
	from, err := LoadSDL(sdlDatabase + sdl)
	require.NoError(t, err)
	to, err := LoadSDL(sdlDatabase + sdl)
	require.NoError(t, err)
	require.True(t, Diff(from, to).IsEmpty())
}

// splitLikeMySQLClient splits a restore script into the statements the mysql client
// sends: it ends one at the active delimiter outside strings and comments, and a
// DELIMITER line changes the delimiter. Unlike mysqlparser.Split it does not know
// BEGIN ... END, so a routine body needs a delimiter its semicolons do not match.
func splitLikeMySQLClient(script string) []string {
	const directive = "DELIMITER "
	var stmts []string
	var cur strings.Builder
	flush := func() {
		if stmt := strings.TrimSpace(cur.String()); stmt != "" {
			stmts = append(stmts, stmt)
		}
		cur.Reset()
	}
	delimiter := ";"
	var quote byte
	for i := 0; i < len(script); {
		c := script[i]
		if quote != 0 {
			cur.WriteByte(c)
			switch {
			case c == '\\' && quote != '`' && i+1 < len(script):
				cur.WriteByte(script[i+1])
				i += 2
				continue
			case c == quote:
				quote = 0
			}
			i++
			continue
		}
		lineEnd := strings.IndexByte(script[i:], '\n')
		if lineEnd < 0 {
			lineEnd = len(script) - i
		}
		switch {
		case strings.TrimSpace(cur.String()) == "" && len(script)-i >= len(directive) && strings.EqualFold(script[i:i+len(directive)], directive):
			delimiter = strings.TrimSpace(script[i+len(directive) : i+lineEnd])
			cur.Reset()
			i += lineEnd
		case strings.HasPrefix(script[i:], delimiter):
			flush()
			i += len(delimiter)
		case c == '#' || strings.HasPrefix(script[i:], "--") && (i+2 == len(script) || script[i+2] <= ' '):
			i += lineEnd
		case strings.HasPrefix(script[i:], "/*") && !strings.HasPrefix(script[i:], "/*!"):
			end := strings.Index(script[i+2:], "*/")
			if end < 0 {
				i = len(script)
			} else {
				i += end + 4
			}
		default:
			if c == '\'' || c == '"' || c == '`' {
				quote = c
			}
			cur.WriteByte(c)
			i++
		}
	}
	flush()
	return stmts
}

func TestTableDefinitionQualifiesOnlyForeignDatabases(t *testing.T) {
	table := &metadata.TableMetadata{
		Name:    "child",
		Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int"}, {Name: "b", Type: "int"}},
		ForeignKeys: []*metadata.ForeignKeyMetadata{
			{Name: "fk_local", Columns: []string{"a"}, ReferencedSchema: "db", ReferencedTable: "p", ReferencedColumns: []string{"id"}},
			{Name: "fk_other", Columns: []string{"b"}, ReferencedSchema: "other", ReferencedTable: "p", ReferencedColumns: []string{"id"}},
		},
	}
	meta := &metadata.DatabaseSchemaMetadata{Name: "db", Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{table}}}}
	for _, sdl := range []bool{false, true} {
		got, err := DatabaseDefinition(meta, DefinitionOptions{SDLFormat: sdl})
		require.NoError(t, err)
		require.Contains(t, got, "FOREIGN KEY (`a`) REFERENCES `p` (`id`)", "SDLFormat=%t", sdl)
		require.Contains(t, got, "FOREIGN KEY (`b`) REFERENCES `other`.`p` (`id`)", "SDLFormat=%t", sdl)
	}
}

func TestIndexUsingClause(t *testing.T) {
	for _, tc := range []struct {
		indexType, engine, want string
	}{
		{"BTREE", "InnoDB", ""},
		{"BTREE", "", ""},
		{"HASH", "InnoDB", " USING HASH"},
		{"HASH", "MEMORY", ""},
		{"BTREE", "MEMORY", " USING BTREE"},
		{"btree", "heap", " USING BTREE"},
		{"FULLTEXT", "InnoDB", ""},
		{"SPATIAL", "MEMORY", ""},
	} {
		require.Equal(t, tc.want, indexUsingClause(tc.indexType, tc.engine), "%s on %s", tc.indexType, tc.engine)
	}
}
