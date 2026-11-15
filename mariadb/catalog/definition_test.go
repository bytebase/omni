package catalog

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/metadata"
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
				// MariaDB's CHECK_CLAUSE, which carries no outer parentheses.
				Expression: "`type` = 'A' and `amount` is not null",
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
		"  CONSTRAINT `c1` CHECK (`type` = 'A' and `amount` is not null)\n" +
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

	// The catalog loads the dump versioned, as MariaDB does.
	c := New()
	_, err = c.Exec("CREATE DATABASE tmp; USE tmp;", nil)
	require.NoError(t, err)
	results, err := c.Exec(got, &ExecOptions{ContinueOnError: true})
	require.NoError(t, err)
	require.False(t, slices.ContainsFunc(results, func(r ExecResult) bool { return r.Error != nil }))
	tbl := c.GetDatabase("tmp").GetTable("explicit_period")
	require.True(t, tbl.SystemVersioned)
	require.Equal(t, "valid_from", tbl.PeriodStartCol)
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

func TestRenderColumnDefault(t *testing.T) {
	testCases := []struct {
		name       string
		columnType string
		// def is ColumnMetadata.Default as MariaDB sync records it: QUOTE() of
		// information_schema's COLUMN_DEFAULT, as captured from mariadb:11.8.8.
		def  string
		want string
	}{
		{name: "string", columnType: "varchar(10)", def: `'\'x\''`, want: "'x'"},
		{name: "null", columnType: "int", def: "'NULL'", want: "NULL"},
		{name: "bare null", columnType: "int", def: "NULL", want: "NULL"},
		{name: "decimal", columnType: "decimal(10,2)", def: "'1.50'", want: "1.50"},
		{name: "zerofill", columnType: "int(5) unsigned zerofill", def: "'00003'", want: "00003"},
		{name: "bit", columnType: "bit(3)", def: `'b\'101\''`, want: "b'101'"},
		{name: "binary", columnType: "varbinary(8)", def: `'x\'0102\''`, want: "x'0102'"},
		{name: "function", columnType: "date", def: "'curdate()'", want: "curdate()"},
		{name: "function with strings", columnType: "varchar(20)", def: `'concat(\'a\',\'b\')'`, want: "concat('a','b')"},
		{name: "sequence", columnType: "int", def: "'nextval(`testdb`.`s1`)'", want: "nextval(`testdb`.`s1`)"},
		{name: "expression", columnType: "int", def: "'`a` + 1'", want: "(`a` + 1)"},
		{name: "current timestamp stays bare", columnType: "datetime(3)", def: "current_timestamp(3)", want: "current_timestamp(3)"},
		// A quoted word that is no expression of COLUMN_DEFAULT's, which
		// backquotes every column, is a string default written by hand.
		{name: "hand-written string", columnType: "varchar(10)", def: "'hello'", want: "'hello'"},
		{name: "hand-written string with spaces", columnType: "varchar(20)", def: "'hello world'", want: "'hello world'"},
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
	require.NoError(t, printColumnClause(&buf, col, table, nil, ""))
	require.Equal(t, "  `b1` bit(1) NOT NULL DEFAULT b'0'", buf.String())

	// A BIT column with no default (b64 BIT(64)) must stay fine: NULL default emitted as
	// DEFAULT NULL, never as a bit literal.
	var buf2 strings.Builder
	col2 := &metadata.ColumnMetadata{Name: "b64", Type: "bit(64)", Nullable: true, Default: "NULL"}
	require.NoError(t, printColumnClause(&buf2, col2, table, nil, ""))
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
	require.NoError(t, printIndexClause(&buf, index, &metadata.TableMetadata{Engine: "InnoDB"}))
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
			name:   "function_qualifier_kept",
			body:   "select `mydb`.`f`(`mydb`.`t`.`a`) AS `x` from `mydb`.`t`",
			dbName: "mydb",
			want:   "select `mydb`.`f`(`t`.`a`) AS `x` from `t`",
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
// TestPrintColumnClauseInvisibleCommentOrder pins INVISIBLE after the default
// and before COMMENT, where mariadb:11.8.8's SHOW CREATE TABLE writes it.
func TestPrintColumnClauseInvisibleCommentOrder(t *testing.T) {
	table := &metadata.TableMetadata{Name: "t"}
	render := func(col *metadata.ColumnMetadata) string {
		var buf strings.Builder
		require.NoError(t, printColumnClause(&buf, col, table, nil, ""))
		return buf.String()
	}
	require.Equal(t, "  `a` int DEFAULT NULL INVISIBLE COMMENT 'c'",
		render(&metadata.ColumnMetadata{Name: "a", Type: "int", Nullable: true, Default: "NULL", IsInvisible: true, Comment: "c"}))
	require.Equal(t,
		"  `g` int GENERATED ALWAYS AS (`a` + 1) STORED INVISIBLE COMMENT 'gen'",
		render(&metadata.ColumnMetadata{Name: "g", Type: "int", Nullable: true, IsInvisible: true, Comment: "gen", Generation: &metadata.GenerationMetadata{
			Type:       metadata.GenerationMetadata_TYPE_STORED,
			Expression: "`a` + 1",
		}}))
}

// splitLikeMySQLClient splits a restore script into the statements the mysql client
// sends: it ends one at the active delimiter outside strings and comments, and a
// DELIMITER line changes the delimiter. Unlike mariadbparser.Split it does not know
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
		case strings.HasPrefix(script[i:], "/*") && !strings.HasPrefix(script[i:], "/*!") && !strings.HasPrefix(script[i:], "/*M!"):
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

func TestRequalifySequences(t *testing.T) {
	for _, tc := range []struct{ expr, from, to, want string }{
		{"nextval(`app`.`s`)", "app", "", "nextval(`s`)"},
		{"nextval(`app`.`s`)", "app", "other", "nextval(`other`.`s`)"},
		{"setval(`app`.`s`,5,1,0)", "app", "", "setval(`s`,5,1,0)"},
		// A string that spells the qualifier is a value, not a sequence.
		{"'`app`.item'", "app", "", "'`app`.item'"},
		{"concat('`app`.', nextval(`app`.`s`))", "app", "", "concat('`app`.', nextval(`s`))"},
		{"nextval(`elsewhere`.`s`)", "app", "", "nextval(`elsewhere`.`s`)"},
	} {
		require.Equal(t, tc.want, requalifySequences(tc.expr, tc.from, tc.to), tc.expr)
	}
}

func TestMultiFileDefinitionWritesSequences(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{Name: "app", Schemas: []*metadata.SchemaMetadata{{
		Sequences: []*metadata.SequenceMetadata{{Name: "s", Start: "1", MinValue: "1", MaxValue: "9223372036854775806", Increment: "1", CacheSize: "1000"}},
		Tables:    []*metadata.TableMetadata{{Name: "t", Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int", Default: "'nextval(`app`.`s`)'"}}}},
	}}}
	files, err := MultiFileDefinition(meta)
	require.NoError(t, err)
	require.Len(t, files, 2)
	require.Equal(t, "CREATE SEQUENCE `s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle;\n\n", files[0].Content)
	require.Contains(t, files[1].Content, "DEFAULT nextval(`s`)")
}

func TestFunctionDefinitionRestoresEmptySQLMode(t *testing.T) {
	got, err := FunctionDefinition(&metadata.FunctionMetadata{
		Name:                "f",
		Definition:          "CREATE FUNCTION `f`() RETURNS int(11)\nRETURN 1",
		CharacterSetClient:  "utf8mb4",
		CollationConnection: "utf8mb4_general_ci",
	})
	require.NoError(t, err)
	require.Contains(t, got, "SET sql_mode = '';\n")
}

func TestPrintColumnClauseNullableTimestamp(t *testing.T) {
	var buf strings.Builder
	require.NoError(t, printColumnClause(&buf, &metadata.ColumnMetadata{Name: "ts", Type: "timestamp(6)", Nullable: true, Default: "NULL"}, &metadata.TableMetadata{}, nil, ""))
	require.Equal(t, "  `ts` timestamp(6) NULL DEFAULT NULL", buf.String())
}

func TestForeignKeyClauseKeepsNoAction(t *testing.T) {
	var buf strings.Builder
	require.NoError(t, printForeignKeyClause(&buf, &metadata.ForeignKeyMetadata{Name: "f", Columns: []string{"a"}, ReferencedTable: "p", ReferencedColumns: []string{"id"}, OnDelete: "NO ACTION", OnUpdate: "RESTRICT"}, ""))
	require.Contains(t, buf.String(), "ON DELETE NO ACTION")
	require.NotContains(t, buf.String(), "ON UPDATE")
}

func TestTableDefinitionKeepsParenthesizedColumnKeys(t *testing.T) {
	got, err := TableDefinition(&metadata.TableMetadata{
		Name:    "t",
		Columns: []*metadata.ColumnMetadata{{Name: "(n)", Type: "int(11)"}, {Name: "n", Type: "int(11)"}},
		Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: []string{"(n)"}, Visible: true},
			{Name: "ix", Type: "BTREE", Expressions: []string{"(n)", "n"}, Visible: true},
		},
	})
	require.NoError(t, err)
	require.Contains(t, got, "PRIMARY KEY (`(n)`)")
	require.Contains(t, got, "KEY `ix` (`(n)`, `n`)")
}

func TestTableDefinitionWritesPrimaryKeyComment(t *testing.T) {
	got, err := TableDefinition(&metadata.TableMetadata{
		Name:    "t",
		Columns: []*metadata.ColumnMetadata{{Name: "id", Type: "int(11)"}},
		Indexes: []*metadata.IndexMetadata{{Name: "PRIMARY", Primary: true, Unique: true, Type: "BTREE", Expressions: []string{"id"}, Comment: "identity", Visible: true}},
	})
	require.NoError(t, err)
	require.Contains(t, got, "PRIMARY KEY (`id`) COMMENT 'identity'")
}

func TestSequenceDefinitionWritesTheSyncedComment(t *testing.T) {
	// Sync stores a comment as QUOTE() escapes it, without the outer quotes.
	got, err := SequenceDefinition(&metadata.SequenceMetadata{Name: "s", Start: "1", Increment: "1", MinValue: "1", MaxValue: "9223372036854775806", CacheSize: "1000", Comment: `owner\'s counter`})
	require.NoError(t, err)
	require.Contains(t, got, `COMMENT='owner\'s counter'`)
	c := New()
	_, err = c.LoadMetadata(context.Background(), &metadata.DatabaseSchemaMetadata{Name: "d", Schemas: []*metadata.SchemaMetadata{{Sequences: []*metadata.SequenceMetadata{{Name: "s", Comment: `owner\'s counter`}}}}})
	require.NoError(t, err)
	require.Equal(t, "owner's counter", c.GetDatabase("d").Sequences["s"].Comment)
}
