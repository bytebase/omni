package parser

import (
	"errors"
	"testing"
)

// TestOverAcceptsRefused holds the rows of compat.tsv that omni accepted where
// the pinned SQL Server refuses, and where the refusal is the grammar's and
// not a check of the parsed statement or a feature the server lacks: each
// statement is refused at the byte measured, and the form beside it, which
// SQL Server takes, still parses. The measured run is
// testdata/README.md's "Rows only omni accepts, measured".
func TestOverAcceptsRefused(t *testing.T) {
	refused := []struct {
		sql string
		at  int
	}{
		// ctasCreateTableStatement is ( columnNameList )? withTableOptions As
		// selectStatement (TSql160.g:523-535), and neither part is optional, so
		// a column name list with no WITH is no CTAS: SQL Server reads (HEAP) as
		// a column definition and answers Msg 173 at it, and ScriptDOM refuses
		// at AS, where the WITH belongs.
		{"CREATE TABLE a(HEAP) AS SELECT *", 21},
		{"CREATE TABLE a(a) AS SELECT *", 18},
		// withTableOptions is one option at least (TSql160.g:26710-26723):
		// WITH () is refused at the parenthesis that closes the empty list.
		{"CREATE TABLE a WITH() AS SELECT *", 20},
		{"CREATE TABLE a(a) WITH() AS SELECT *", 23},
	}
	for _, tc := range refused {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", tc.sql, err)
			continue
		}
		if pe.Position != tc.at {
			t.Errorf("Parse(%q): refused at %d (%s), want %d", tc.sql, pe.Position, pe.Message, tc.at)
		}
	}

	// accepted: the forms beside them. The server has no CTAS, so its forms
	// are the grammar's, which ScriptDOM parses; the rest it takes.
	accepted := []string{
		// AS selectStatement without a column list or WITH is the alternative
		// beside the CTAS rule (TSql160.g:26584).
		"CREATE TABLE a AS SELECT *",
		"CREATE TABLE a AS SELECT * FROM a",
		// The column list goes with the WITH options, and the options are the
		// ones a CTAS takes: HEAP and DISTRIBUTION are Synapse's.
		"CREATE TABLE a(a) WITH(HEAP) AS SELECT *",
		"CREATE TABLE a WITH(HEAP) AS SELECT *",
		"CREATE TABLE a(a INT)",
		"CREATE TABLE a(a INT) WITH(DATA_COMPRESSION = PAGE)",
		"CREATE TABLE a AS FILETABLE",
	}
	for _, sql := range accepted {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accept", sql, err)
		}
	}
}

// TestLineCommentOpening holds omni to what SQL Server reads where a line
// comment opens with ( *: the four characters --(* are no comment at all
// (Msg 102 near '--(*', the answer ScriptDOM's lexer gives too), and every
// other line comment is a comment, the same four characters later in one
// included.
func TestLineCommentOpening(t *testing.T) {
	refused := []struct {
		sql string
		at  int
	}{
		{"SELECT 1 --(* a *)", 9},
		{"SELECT 1 --(* a *)\n", 9},
		{"SELECT 1 --(*", 9},
		{"--(* a *)", 0},
	}
	for _, tc := range refused {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", tc.sql, err)
			continue
		}
		if pe.Position != tc.at {
			t.Errorf("Parse(%q): refused at %d (%s), want %d", tc.sql, pe.Position, pe.Message, tc.at)
		}
	}

	// accepted: the comment forms beside it, with the same four characters
	// where they do not open the comment.
	accepted := []string{
		"SELECT 1 --( *",
		"SELECT 1 --((* a *)",
		"SELECT 1 -- (* a *)",
		"SELECT 1 --a(* a *)",
		"SELECT 1 --x --(* a *)",
		"SELECT 1 /* (* a *) */",
	}
	for _, sql := range accepted {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accept", sql, err)
		}
	}
}

// TestKeywordChoiceSitesMeasured holds the forms beside the (every
// alternative) rows of tokens.tsv that closed gaps, each as the pinned SQL
// Server 2022 answered it under SET PARSEONLY ON (2026-11-26).
func TestKeywordChoiceSitesMeasured(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		// comparisonOperator (TSql160.g:31264) in an index filter, its two
		// characters written apart.
		{"CREATE INDEX a ON a(a) WHERE a ! = 1", true},
		{"CREATE INDEX a ON a(a) WHERE a < = 1", true},
		{"CREATE INDEX a ON a(a) WHERE a > = 1", true},
		{"CREATE INDEX a ON a(a) WHERE a < > 1", true},
		{"CREATE INDEX a ON a(a) WHERE a = 1 AND a ! > 1", true},
		{"ALTER TABLE a ADD INDEX a WHERE IDENTITYCOL ! < 1", true},
		{"ALTER TABLE a ADD INDEX a WHERE IDENTITYCOL ! 1", false},
		// openXmlSchemaItem's pattern is stringOrVariable (TSql160.g:20564);
		// OPENJSON's is a string only.
		{"SELECT * FROM OPENXML(@a, 'a') WITH(a a @a, b int 'b')", true},
		{"SELECT * FROM OPENXML(@a, 'a') WITH(a a @@a)", true},
		{"SELECT * FROM OPENXML(@a, 'a') WITH(a a 1)", false},
		{"SELECT * FROM OPENJSON(@a) WITH(a a @a)", false},
		// HOLDLOCK where a name goes in t( ... ), and no condition there.
		{"SELECT * FROM a(HOLDLOCK)", true},
		{"SELECT * FROM dbo.a(HOLDLOCK) AS b", true},
		{"SELECT * FROM a(NOLOCK, HOLDLOCK)", true},
		{"SELECT * FROM a(HOLDLOCK NOLOCK)", false},
		{"SELECT * FROM a(SPATIAL_WINDOW_MAX_CELLS = 1)", false},
		{"SELECT * FROM a(b IS NULL)", false},
		{"SELECT * FROM a(1, b = 1)", false},
		{"SELECT * FROM a(b IN (1))", false},
		{"SELECT * FROM a((SELECT 1))", true},
		{"SELECT * FROM a(DEFAULT)", true},
		// A text column's last part may be a pseudo-column (TSql160.g:32740).
		{"WRITETEXT a.$IDENTITY 0x00 NULL", true},
		{"READTEXT a.$ROWGUID 0x00 1 1", true},
		{"UPDATETEXT a.b.$IDENTITY 0x00 NULL NULL", true},
		{"WRITETEXT a.$a 0x00 NULL", false},
		{"WRITETEXT $IDENTITY 0x00 NULL", false},
	} {
		_, err := Parse(tc.sql)
		if got := err == nil; got != tc.accept {
			t.Errorf("Parse(%q): accept = %v (%v), SQL Server: %v", tc.sql, got, err, tc.accept)
		}
	}
}
