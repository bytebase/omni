package catalog

import (
	"strings"
	"testing"
)

// Stored-form fidelity for temporal-unit / keyword-argument positions.
//
// Every case's `stored` string is a live-engine SHOW CREATE VIEW readback
// (oracle 8.0.32 :13306 + 5.7.25 :13307 — the two agree on every form here
// except WEIGHT_STRING's LEVEL suffix, which only 5.7 emits). LoadSDL of the
// stored form must regenerate EXACTLY the stored form: any drift (a
// backtick-quoted unit, a lost keyword form) produces CREATE VIEW DDL the
// engine rejects with error 1064 — the sys-schema dogfood bug this file
// pins: timestampdiff(`SECOND`,...) fails to apply while the engine stores
// timestampdiff(SECOND,...) bare.
func TestTemporalUnitStoredFormRoundTrip(t *testing.T) {
	cases := []struct{ id, stored string }{
		// TIMESTAMPDIFF: unit stays bare UPPERCASE, args joined without spaces.
		{"tsdiff-microsecond", "select timestampdiff(MICROSECOND,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-second", "select timestampdiff(SECOND,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-minute", "select timestampdiff(MINUTE,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-hour", "select timestampdiff(HOUR,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-day", "select timestampdiff(DAY,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-week", "select timestampdiff(WEEK,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-month", "select timestampdiff(MONTH,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-quarter", "select timestampdiff(QUARTER,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-year", "select timestampdiff(YEAR,'2020-01-01','2021-01-01') AS `x`"},
		// The sys.innodb_lock_waits shape: column refs + now() as arguments.
		{"tsdiff-column-args", "select timestampdiff(SECOND,`t`.`created`,now()) AS `x`"},
		// Nested interval arithmetic inside a timestampdiff argument.
		{"tsdiff-nested-interval", "select timestampdiff(SECOND,(now() - interval 1 day),now()) AS `x`"},
		// INTERVAL arithmetic: lowercase unit, value parenthesized when non-literal.
		{"interval-plus", "select ('2020-01-01' + interval 1 second) AS `x`"},
		{"interval-minus", "select ('2020-01-01' - interval 1 day) AS `x`"},
		{"interval-compound", "select ('2020-01-01' + interval '1:1' day_hour) AS `x`"},
		{"interval-expr-value", "select (now() + interval (1 + 1) day) AS `x`"},
		{"interval-neg-value", "select ('2020-01-01' + interval -(5) day) AS `x`"},
		// EXTRACT: lowercase unit, simple and compound.
		{"extract-microsecond", "select extract(microsecond from '2020-01-01 10:20:30') AS `x`"},
		{"extract-second", "select extract(second from '2020-01-01 10:20:30') AS `x`"},
		{"extract-minute", "select extract(minute from '2020-01-01 10:20:30') AS `x`"},
		{"extract-hour", "select extract(hour from '2020-01-01 10:20:30') AS `x`"},
		{"extract-day", "select extract(day from '2020-01-01 10:20:30') AS `x`"},
		{"extract-week", "select extract(week from '2020-01-01 10:20:30') AS `x`"},
		{"extract-month", "select extract(month from '2020-01-01 10:20:30') AS `x`"},
		{"extract-quarter", "select extract(quarter from '2020-01-01 10:20:30') AS `x`"},
		{"extract-year", "select extract(year from '2020-01-01 10:20:30') AS `x`"},
		{"extract-second-microsecond", "select extract(second_microsecond from '2020-01-01 10:20:30') AS `x`"},
		{"extract-minute-microsecond", "select extract(minute_microsecond from '2020-01-01 10:20:30') AS `x`"},
		{"extract-minute-second", "select extract(minute_second from '2020-01-01 10:20:30') AS `x`"},
		{"extract-hour-microsecond", "select extract(hour_microsecond from '2020-01-01 10:20:30') AS `x`"},
		{"extract-hour-second", "select extract(hour_second from '2020-01-01 10:20:30') AS `x`"},
		{"extract-hour-minute", "select extract(hour_minute from '2020-01-01 10:20:30') AS `x`"},
		{"extract-day-microsecond", "select extract(day_microsecond from '2020-01-01 10:20:30') AS `x`"},
		{"extract-day-second", "select extract(day_second from '2020-01-01 10:20:30') AS `x`"},
		{"extract-day-minute", "select extract(day_minute from '2020-01-01 10:20:30') AS `x`"},
		{"extract-day-hour", "select extract(day_hour from '2020-01-01 10:20:30') AS `x`"},
		{"extract-year-month", "select extract(year_month from '2020-01-01 10:20:30') AS `x`"},
		// GET_FORMAT: bare UPPERCASE type, and a space after its comma.
		{"get-format-date", "select get_format(DATE, 'USA') AS `x`"},
		{"get-format-time", "select get_format(TIME, 'USA') AS `x`"},
		{"get-format-datetime", "select get_format(DATETIME, 'USA') AS `x`"},
		// WEIGHT_STRING: plain, AS CHAR(n), the desugared AS BINARY cast form,
		// and the 5.7 LEVEL suffixes (5.7 readbacks always carry a level list).
		{"weight-string-plain", "select weight_string('ab') AS `x`"},
		{"weight-string-as-char", "select weight_string('ab' as char(4)) AS `x`"},
		{"weight-string-cast-binary", "select weight_string(cast('ab' as char(4) charset binary)) AS `x`"},
		{"weight-string-level", "select weight_string('ab' level 1) AS `x`"},
		{"weight-string-level-desc", "select weight_string('ab' level 1 desc) AS `x`"},
		{"weight-string-level-reverse", "select weight_string('ab' level 1 reverse) AS `x`"},
		{"weight-string-as-char-level", "select weight_string('ab' as char(4) level 1) AS `x`"},
		// TRIM: directional and the direction-less remstr FROM form.
		{"trim-both", "select trim(both 'x' from 'xax') AS `x`"},
		{"trim-leading", "select trim(leading 'x' from 'xax') AS `x`"},
		{"trim-trailing", "select trim(trailing 'x' from 'xax') AS `x`"},
		{"trim-remstr-from", "select trim(' x' from 'axa') AS `x`"},
		{"trim-plain", "select trim(' a ') AS `x`"},
		// SUBSTRING: engine stores the comma form under the substr name.
		{"substr-comma", "select substr('abcd',2,2) AS `x`"},
	}
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			sdl := "CREATE DATABASE d;\nUSE d;\n" +
				"CREATE TABLE `t` (`id` int NOT NULL, `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n" +
				"CREATE OR REPLACE VIEW `v` AS " + tc.stored + ";\n"
			cat, err := LoadSDL(sdl)
			if err != nil {
				t.Fatalf("LoadSDL(stored form) failed: %v", err)
			}
			v := cat.GetDatabase("d").Views["v"]
			if v == nil {
				t.Fatal("view not loaded")
			}
			if v.Definition != tc.stored {
				t.Errorf("stored form not a fixed point:\n  stored: %s\n  omni:   %s", tc.stored, v.Definition)
			}
		})
	}
}

// User-declared forms must canonicalize to the engine's stored form, so a
// no-op SDL diff against a live readback stays empty. Every `want` string is
// an oracle readback (8.0.32 + 5.7.25 agree): the date-arithmetic function
// family is stored as INTERVAL arithmetic, units are case-folded, SQL_TSI_
// synonyms are folded, and GET_FORMAT(TIMESTAMP,...) stores as DATETIME.
func TestTemporalUnitUserFormCanonicalization(t *testing.T) {
	cases := []struct{ id, user, want string }{
		{"timestampadd", "select TIMESTAMPADD(HOUR, 1, '2020-01-01') AS `x`", "select ('2020-01-01' + interval 1 hour) AS `x`"},
		{"timestampadd-sql-tsi", "select TIMESTAMPADD(SQL_TSI_HOUR, 1, '2020-01-01') AS `x`", "select ('2020-01-01' + interval 1 hour) AS `x`"},
		{"date-add", "select DATE_ADD('2020-01-01', INTERVAL 1 DAY) AS `x`", "select ('2020-01-01' + interval 1 day) AS `x`"},
		{"date-sub", "select DATE_SUB('2020-01-01', INTERVAL 1 DAY) AS `x`", "select ('2020-01-01' - interval 1 day) AS `x`"},
		{"date-add-compound", "select DATE_ADD('2020-01-01', INTERVAL '1:1' DAY_HOUR) AS `x`", "select ('2020-01-01' + interval '1:1' day_hour) AS `x`"},
		{"date-add-expr-count", "select DATE_ADD('2020-01-01', INTERVAL 2*3 HOUR) AS `x`", "select ('2020-01-01' + interval (2 * 3) hour) AS `x`"},
		{"adddate-interval", "select ADDDATE('2020-01-01', INTERVAL 1 DAY) AS `x`", "select ('2020-01-01' + interval 1 day) AS `x`"},
		{"adddate-days", "select ADDDATE('2020-01-01', 31) AS `x`", "select ('2020-01-01' + interval 31 day) AS `x`"},
		{"adddate-negative-days", "select ADDDATE('2020-01-01', -5) AS `x`", "select ('2020-01-01' + interval -(5) day) AS `x`"},
		{"subdate-interval", "select SUBDATE('2020-01-01', INTERVAL 1 DAY) AS `x`", "select ('2020-01-01' - interval 1 day) AS `x`"},
		{"subdate-days", "select SUBDATE('2020-01-01', 31) AS `x`", "select ('2020-01-01' - interval 31 day) AS `x`"},
		{"nested-date-add", "select DATE_ADD(DATE_ADD('2020-01-01', INTERVAL 1 DAY), INTERVAL 1 HOUR) AS `x`", "select (('2020-01-01' + interval 1 day) + interval 1 hour) AS `x`"},
		{"bare-interval", "select NOW() + INTERVAL 1 DAY AS `x`", "select (now() + interval 1 day) AS `x`"},
		{"interval-first-commutes", "select INTERVAL 1 DAY + NOW() AS `x`", "select (now() + interval 1 day) AS `x`"},
		{"interval-sql-tsi", "select NOW() + INTERVAL 1 SQL_TSI_DAY AS `x`", "select (now() + interval 1 day) AS `x`"},
		{"tsdiff-lowercase-unit", "select timestampdiff(second,'2020-01-01','2021-01-01') AS `x`", "select timestampdiff(SECOND,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-sql-tsi", "select timestampdiff(SQL_TSI_SECOND,'2020-01-01','2021-01-01') AS `x`", "select timestampdiff(SECOND,'2020-01-01','2021-01-01') AS `x`"},
		{"get-format-lowercase", "select get_format(date,'USA') AS `x`", "select get_format(DATE, 'USA') AS `x`"},
		{"get-format-timestamp", "select GET_FORMAT(TIMESTAMP,'USA') AS `x`", "select get_format(DATETIME, 'USA') AS `x`"},
		{"extract-uppercase-unit", "select EXTRACT(DAY_HOUR FROM '2020-01-01 10:20:30') AS `x`", "select extract(day_hour from '2020-01-01 10:20:30') AS `x`"},
		{"extract-sql-tsi", "select EXTRACT(SQL_TSI_DAY FROM '2020-01-01') AS `x`", "select extract(day from '2020-01-01') AS `x`"},
		{"weight-string-as-binary", "select WEIGHT_STRING('ab' AS BINARY(4)) AS `x`", "select weight_string(cast('ab' as char(4) charset binary)) AS `x`"},
		{"weight-string-level-asc", "select WEIGHT_STRING('ab' LEVEL 1 ASC) AS `x`", "select weight_string('ab' level 1) AS `x`"},
		{"substring-from-for", "select SUBSTRING('abcd' FROM 2 FOR 2) AS `x`", "select substr('abcd',2,2) AS `x`"},
		// A schema-qualified name is a stored function, not the builtin: its
		// arguments are ordinary expressions and the qualifier must survive
		// deparse (regression: the GET_FORMAT deparse special case must not
		// swallow `mydb`.). The unquoted qualified rendering below is omni's
		// long-standing form for ALL qualified functions — the engine stores
		// them backticked (oracle 8.0.32), a separate pre-existing gap.
		{"qualified-get-format", "select `mydb`.`get_format`('a','b') AS `x`", "select mydb.get_format('a','b') AS `x`"},
		{"qualified-timestampdiff", "select `mydb`.`timestampdiff`(`c`,'a','b') AS `x`", "select mydb.timestampdiff(`c`,'a','b') AS `x`"},
	}
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			sdl := "CREATE DATABASE d;\nUSE d;\nCREATE OR REPLACE VIEW `v` AS " + tc.user + ";\n"
			cat, err := LoadSDL(sdl)
			if err != nil {
				t.Fatalf("LoadSDL(user form) failed: %v", err)
			}
			v := cat.GetDatabase("d").Views["v"]
			if v == nil {
				t.Fatal("view not loaded")
			}
			if v.Definition != tc.want {
				t.Errorf("user form not canonicalized to the stored form:\n  user: %s\n  want: %s\n  omni: %s", tc.user, tc.want, v.Definition)
			}
		})
	}
}

// Generated-column bodies and functional indexes are rendered by their own
// printers (nodeToSQLGenerated / the index-expression path), which must emit
// the same engine-stored forms as view bodies. Every `stored` fragment below
// is a live SHOW CREATE TABLE readback (oracle 8.0.32): a lost case there
// renders `(?)` and the CREATE TABLE fails to apply.
func TestTemporalUnitGeneratedColumnRoundTrip(t *testing.T) {
	sdl := "CREATE DATABASE d;\nUSE d;\n" +
		"CREATE TABLE `fi` (\n" +
		"  `a` datetime DEFAULT NULL,\n" +
		"  `b` datetime DEFAULT NULL,\n" +
		"  KEY `idx` ((timestampdiff(SECOND,`a`,`b`)))\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n" +
		"CREATE TABLE `g` (\n" +
		"  `a` datetime DEFAULT NULL,\n" +
		"  `s` varchar(10) DEFAULT NULL,\n" +
		"  `c1` int GENERATED ALWAYS AS (timestampdiff(SECOND,`a`,`a`)) VIRTUAL,\n" +
		"  `c2` int GENERATED ALWAYS AS (extract(day_hour from `a`)) VIRTUAL,\n" +
		"  `c3` datetime GENERATED ALWAYS AS ((`a` + interval 1 day)) VIRTUAL,\n" +
		"  `c4` int GENERATED ALWAYS AS (abs(-(5))) VIRTUAL,\n" +
		"  `c5` varbinary(16) GENERATED ALWAYS AS (weight_string(`s` as char(4))) VIRTUAL\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n"
	cat, err := LoadSDL(sdl)
	if err != nil {
		t.Fatalf("LoadSDL failed: %v", err)
	}
	cat2, err := LoadSDL(sdl)
	if err != nil {
		t.Fatalf("LoadSDL (second) failed: %v", err)
	}
	if d := Diff(cat, cat2); !d.IsEmpty() {
		t.Errorf("self-diff not empty")
	}
	plan := GenerateMigration(New(), cat, Diff(New(), cat))
	sql := plan.SQL()
	if strings.Contains(sql, "(?)") {
		t.Errorf("generated plan contains the (?) placeholder — an expression node is missing from the generated-column renderer:\n%s", sql)
	}
	for _, want := range []string{
		"(timestampdiff(SECOND,`a`,`a`))",
		"(extract(day_hour from `a`))",
		"((`a` + interval 1 day))",
		"(abs(-(5)))",
		"(weight_string(`s` as char(4)))",
		"((timestampdiff(SECOND,`a`,`b`)))",
	} {
		if !strings.Contains(sql, want) {
			t.Errorf("generated plan missing the engine-stored fragment %q:\n%s", want, sql)
		}
	}
}

// The engine rejects malformed keyword arguments with error 1064; the parser
// must reject them too rather than mis-parse them as identifiers. The
// quoted-unit case is the inverse of the dogfood bug: timestampdiff(`SECOND`,...)
// is NOT valid MySQL (oracle 8.0.32 + 5.7.25).
func TestTemporalUnitParserRejections(t *testing.T) {
	cases := []struct{ id, body string }{
		{"tsdiff-quoted-unit", "select timestampdiff(`SECOND`,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-string-unit", "select timestampdiff('SECOND','2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-compound-unit", "select timestampdiff(DAY_HOUR,'2020-01-01','2021-01-01') AS `x`"},
		{"tsdiff-sql-tsi-microsecond", "select timestampdiff(SQL_TSI_MICROSECOND,'2020-01-01','2021-01-01') AS `x`"},
		{"get-format-year", "select get_format(YEAR,'USA') AS `x`"},
		{"get-format-string", "select get_format('DATE','USA') AS `x`"},
		{"extract-quoted-unit", "select extract(`day` from '2020-01-01') AS `x`"},
	}
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			sdl := "CREATE DATABASE d;\nUSE d;\nCREATE VIEW `v` AS " + tc.body + ";\n"
			if _, err := LoadSDL(sdl); err == nil {
				t.Errorf("expected parse error (engine rejects with 1064), got success for: %s", tc.body)
			} else if !strings.Contains(err.Error(), "line") {
				t.Logf("rejected (as expected): %v", err)
			}
		})
	}
}
