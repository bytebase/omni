package parser

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestCreateOrReplaceOracle pins OceanBase's verdict on createOrReplaceCases
// and checks omni against it. The verdict comes from PREPARE, which does not
// run the statement. The expected code pins the engine, so an image upgrade that
// changes a verdict fails here instead of passing silently. For a 1064, the
// engine's message must quote from the token where the case's near starts,
// which places omni's error at the same token as OceanBase's.
func TestCreateOrReplaceOracle(t *testing.T) {
	o := startParserOracle(t)
	if _, err := o.db.ExecContext(o.ctx, "CREATE TABLE t (id INT PRIMARY KEY, n INT)"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	for _, tc := range createOrReplaceCases {
		var code uint16
		_, err := o.db.ExecContext(o.ctx, "PREPARE parse_probe FROM '"+strings.ReplaceAll(tc.sql, "'", "''")+"'")
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Errorf("PREPARE %q: %v", tc.sql, err)
				continue
			}
			code = myErr.Number
		}
		if code != tc.code {
			t.Errorf("OceanBase code %d, want %d for %q: %v", code, tc.code, tc.sql, err)
			continue
		}
		if code == 1064 {
			if quoted, ok := quotedNear(err.Error()); !ok || !nearMatches(tc.sql, tc.near, quoted) {
				t.Errorf("OceanBase's error for %q does not quote from %q: %v", tc.sql, tc.near, err)
			}
		}
		_, perr := Parse(tc.sql)
		if omniAccepts := perr == nil; omniAccepts != parsesInEngine(code) {
			t.Errorf("DIVERGENCE omni accepts=%v, OceanBase code %d for %q: %v", omniAccepts, code, tc.sql, perr)
		}
	}
}

var nearQuote = regexp.MustCompile(`(?s)near '(.*)' at line \d+$`)

// quotedNear returns the text an OceanBase 1064 quotes after "near".
func quotedNear(msg string) (string, bool) {
	m := nearQuote.FindStringSubmatch(msg)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// nearMatches reports whether OceanBase's quote starts where near does: the
// quote is a prefix of near, often the one token. At the end of the input near
// is empty and OceanBase quotes the statement's last token.
func nearMatches(sql, near, quoted string) bool {
	if near == "" {
		return quoted != "" && strings.HasSuffix(strings.TrimSpace(sql), quoted)
	}
	return quoted != "" && strings.HasPrefix(near, quoted)
}
