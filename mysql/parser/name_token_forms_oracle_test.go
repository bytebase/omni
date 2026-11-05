package parser

import (
	"errors"
	"sort"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// serverParses asks the oracle through PREPARE whether sql parses, and reads an error
// beyond 1064 by where it was raised. The server refuses a few statements while it parses
// with the code of the check that refused them, as FOR CHANNEL's line feed answers 1525.
// The server reads all of the text before it looks anything up, so a statement refused
// after parsing answers 1064 once a syntax error follows it as a second statement, and one
// refused while parsing answers its own code again. A system variable it looks up while
// parsing (1193) counts as read.
func (o *parserOracle) serverParses(t *testing.T, sql string) bool {
	t.Helper()
	conn, err := o.db.Conn(o.ctx)
	if err != nil {
		t.Fatalf("taking a connection for %q: %v", sql, err)
	}
	defer conn.Close()
	prepare := func(text string) uint16 {
		if _, err := conn.ExecContext(o.ctx, "SET @omni_name_form = ?", text); err != nil {
			t.Fatalf("binding %q: %v", text, err)
		}
		_, err := conn.ExecContext(o.ctx, "PREPARE omni_name_form_stmt FROM @omni_name_form")
		if err == nil {
			_, _ = conn.ExecContext(o.ctx, "DEALLOCATE PREPARE omni_name_form_stmt")
			return 0
		}
		var myErr *gomysql.MySQLError
		if !errors.As(err, &myErr) {
			t.Fatalf("no verdict for %q: %v", text, err)
		}
		return myErr.Number
	}
	switch code := prepare(sql); code {
	case 0, 1193:
		return true
	case 1064:
		return false
	default:
		return prepare(sql+"\n; SELECT FROM") != code
	}
}

// TestServerParsesControls holds serverParses to statements whose answer is known.
func TestServerParsesControls(t *testing.T) {
	o := startParserOracle(t)
	for _, c := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT 1 FROM nosuch", true},                          // 1146, after parsing
		{"SELECT 1 FROM nosuch -- a comment", true},             // the tail survives the comment
		{"SELEKT 1", false},                                     // 1064
		{"FLUSH RELAY LOGS FOR CHANNEL 'a\nb'", false},          // 1525, while parsing
		{"SELECT @@nosuch", true},                               // 1193, a name it looks up
		{"FLUSH RELAY LOGS FOR CHANNEL 'c' FOR CHANNEL", false}, // 1064
	} {
		if got := o.serverParses(t, c.sql); got != c.accept {
			t.Errorf("%q: serverParses=%v, want %v", c.sql, got, c.accept)
		}
	}
}

// TestNameTokenFormsOracle asserts omni agrees with a live MySQL 8.0 on
// nameTokenFormCases and flushFormCases (container test).
func TestNameTokenFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, cases := range [][]parserFormCase{nameTokenFormCases, flushFormCases} {
		for _, tc := range cases {
			t.Run(tc.sql, func(t *testing.T) {
				if got := o.serverParses(t, tc.sql); got != tc.accept {
					t.Fatalf("MySQL accepts=%v, want %v for %q", got, tc.accept, tc.sql)
				}
				_, err := Parse(tc.sql)
				if omniAccepts := err == nil; omniAccepts != tc.accept {
					t.Errorf("DIVERGENCE omni accepts=%v, MySQL accepts=%v for %q: %v",
						omniAccepts, tc.accept, tc.sql, err)
				}
			})
		}
	}
}

// TestIntroducerWordsMatchMySQL sweeps every character set and collation the server
// lists, and every name charsetIntroducerNames holds, through a name position as an
// underscored word, and holds omni's verdict to the server's. A word the server reads as
// an introducer is a 1064 there and nowhere else, so this catches a set that is missing a
// character set or holds one the server lacks. After a separating dot no word introduces,
// which each introducer is checked for too.
func TestIntroducerWordsMatchMySQL(t *testing.T) {
	o := startParserOracle(t)
	words := map[string]bool{"utf8": true, "filename": true, "default": true, "foo": true}
	for name := range charsetIntroducerNames {
		words[name] = true
	}
	rows, err := o.db.QueryContext(o.ctx,
		"SELECT CHARACTER_SET_NAME FROM information_schema.CHARACTER_SETS "+
			"UNION SELECT COLLATION_NAME FROM information_schema.COLLATIONS")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		words[name] = true
	}
	if err := rows.Close(); err != nil {
		t.Fatal(err)
	}
	sorted := make([]string, 0, len(words))
	for w := range words {
		sorted = append(sorted, w)
	}
	sort.Strings(sorted)

	introducers := 0
	for _, w := range sorted {
		sqls := []string{"CREATE TABLE _" + w + " (b int)"}
		if charsetIntroducerNames[w] {
			sqls = append(sqls, "SELECT * FROM test._"+w)
		}
		for _, sql := range sqls {
			server := o.serverParses(t, sql)
			_, err := Parse(sql)
			if omni := err == nil; omni != server {
				t.Errorf("DIVERGENCE omni accepts=%v, MySQL accepts=%v for %q: %v", omni, server, sql, err)
			}
			if !server {
				introducers++
			}
		}
	}
	// The 41 character sets of SHOW CHARACTER SET, and utf8.
	if introducers != 42 {
		t.Errorf("MySQL refused %d words as introducers, want 42", introducers)
	}
}
