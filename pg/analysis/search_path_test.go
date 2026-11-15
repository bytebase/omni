package analysis

import "testing"

func TestWithSearchPath(t *testing.T) {
	statement := WithSearchPath("DELETE FROM t;", []string{"app", `we"ird`})
	if want := "SET LOCAL search_path TO \"app\", \"we\"\"ird\";\nDELETE FROM t;"; statement != want {
		t.Fatalf("WithSearchPath = %q, want %q", statement, want)
	}
	setup, rest := SplitSearchPath(statement)
	if setup != `SET LOCAL search_path TO "app", "we""ird"` || rest != "DELETE FROM t;" {
		t.Fatalf("SplitSearchPath = (%q, %q)", setup, rest)
	}
	setup, rest = SplitSearchPath("DELETE FROM t;")
	if setup != "" || rest != "DELETE FROM t;" {
		t.Fatalf("SplitSearchPath plain = (%q, %q)", setup, rest)
	}
}
