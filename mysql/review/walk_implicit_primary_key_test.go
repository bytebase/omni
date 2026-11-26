package review

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/review"
)

// TestWalkImplicitPrimaryKeyVisible holds the walk to mysql/catalog's
// invisiblePrimaryKeyCases (implicit_primary_key_test.go), whose MySQL
// answers were measured on mysql:8.0.46: 3522 for an invisible primary
// key, explicit or the UNIQUE key MySQL makes one. The walk, like the
// catalog, leaves expression key parts alone.
func TestWalkImplicitPrimaryKeyVisible(t *testing.T) {
	target := review.Target{Schema: walkSchema()}
	for _, tc := range []struct {
		name    string
		sql     string
		refused bool
	}{
		{"explicit", "CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE)", true},
		{"explicit_last_option_wins", "CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE VISIBLE)", false},
		{"sole_unique", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE)", true},
		{"constraint_unique", "CREATE TABLE t (a INT NOT NULL, CONSTRAINT u UNIQUE (a) INVISIBLE)", true},
		{"composite", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a, b) INVISIBLE)", true},
		{"last_option_wins_invisible", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) VISIBLE INVISIBLE)", true},
		{"last_option_wins_visible", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE VISIBLE)", false},
		{"first_of_two", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, UNIQUE KEY v (b))", true},
		{"second_of_two", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)", false},
		{"plain_key_first", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, KEY k (b), UNIQUE KEY u (a) INVISIBLE)", true},
		{"nullable", "CREATE TABLE t (a INT, UNIQUE KEY u (a) INVISIBLE)", false},
		{"partly_nullable", "CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY u (a, b) INVISIBLE)", false},
		{"nullable_key_first", "CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)", true},
		{"explicit_key", "CREATE TABLE t (a INT NOT NULL, b INT, PRIMARY KEY (b), UNIQUE KEY u (a) INVISIBLE)", false},
		{"explicit_key_last", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, PRIMARY KEY (b))", false},
		{"column_primary_key", "CREATE TABLE t (a INT PRIMARY KEY, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", false},
		{"column_unique_first", "CREATE TABLE t (a INT NOT NULL UNIQUE, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", false},
		{"serial_first", "CREATE TABLE t (a SERIAL, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", false},
		{"column_unique_after", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, b INT NOT NULL UNIQUE)", true},
		{"prefix", "CREATE TABLE t (a VARCHAR(10) NOT NULL, UNIQUE KEY u (a(3)) INVISIBLE)", true},
		{"prefix_first_visible", "CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(3)), UNIQUE KEY v (b) INVISIBLE)", false},
		{"auto_increment", "CREATE TABLE t (a INT AUTO_INCREMENT, UNIQUE KEY u (a) INVISIBLE)", true},
		{"auto_increment_unique_first", "CREATE TABLE t (a INT AUTO_INCREMENT UNIQUE, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", false},
		{"stored_generated", "CREATE TABLE t (a INT NOT NULL, g INT AS (a + 1) STORED NOT NULL, UNIQUE KEY u (g) INVISIBLE)", true},
		{"expression_first_visible", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY f ((a + 1)), UNIQUE KEY v (b) INVISIBLE)", false},
		{"expression_not_null", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u ((a + 1)) INVISIBLE)", false},
		// mysql/catalog's TestCreateTableGeneratedPrimaryKeyTakesInvisibleUnique.
		{"generated_primary_key", "SET SESSION sql_generate_invisible_primary_key = ON; CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE)", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sql := strings.Replace(tc.sql, "CREATE TABLE t ", "CREATE TABLE ipk ", 1)
			if tc.refused {
				oneFinding(t, sql, 0, "A primary key can't be invisible.", target)
			} else {
				noFinding(t, sql, target)
			}
		})
	}
}
