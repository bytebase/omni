package catalog

import "testing"

// TestRoutineReadbackWithComment pins how a routine's COMMENT, which OceanBase's SHOW
// CREATE FUNCTION and SHOW CREATE PROCEDURE print unescaped, is written again. The
// readbacks are OceanBase 4.3.5.6's (TestOracle_RoutineDiffIdempotence's tape).
func TestRoutineReadbackWithComment(t *testing.T) {
	const head = "CREATE DEFINER = `root`@`%` FUNCTION `f`(a INT) RETURNS int(11)\n    DETERMINISTIC\n"
	for _, tc := range []struct {
		name, def, comment, want string
	}{
		{"quote", head + "    COMMENT 'it's here'\nRETURN a", "it's here", head + "    COMMENT 'it''s here'\nRETURN a"},
		{"backslash", head + "    COMMENT 'a\\b'\nRETURN a", `a\b`, head + "    COMMENT 'a\\\\b'\nRETURN a"},
		{"newline", "CREATE DEFINER = `root`@`%` PROCEDURE `p1`()\n    COMMENT 'x'y\\z\nw'\nBEGIN SELECT 1; END", "x'y\\z\nw", "CREATE DEFINER = `root`@`%` PROCEDURE `p1`()\n    COMMENT 'x''y\\\\z\nw'\nBEGIN SELECT 1; END"},
		{"plain", head + "    COMMENT 'hi'\nRETURN a", "hi", head + "    COMMENT 'hi'\nRETURN a"},
		{"no comment", head + "RETURN a", "", head + "RETURN a"},
		{"not printed", head + "RETURN a", "it's", head + "RETURN a"},
		// Only the characteristic is written again: the body follows it, and its text is
		// the body's own.
		{"body line", head + "    COMMENT 'it's'\nBEGIN\n    SELECT 1 FROM t WHERE c = '\n    COMMENT 'it's'\n';\nEND", "it's",
			head + "    COMMENT 'it''s'\nBEGIN\n    SELECT 1 FROM t WHERE c = '\n    COMMENT 'it's'\n';\nEND"},
	} {
		if got := RoutineReadbackWithComment(tc.def, tc.comment); got != tc.want {
			t.Errorf("%s: RoutineReadbackWithComment(%q, %q) =\n%q\nwant\n%q", tc.name, tc.def, tc.comment, got, tc.want)
		}
	}
}
