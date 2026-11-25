package completion

import (
	"crypto/sha256"
	"encoding/hex"
	"slices"
	"strings"
	"testing"
)

// TestFunctionsPin holds functions.tsv to its pin: the digest in its header
// is the SHA-256 of its rows, as cmd/docsfunctions wrote them.
func TestFunctionsPin(t *testing.T) {
	var pin string
	var rows strings.Builder
	for _, line := range strings.SplitAfter(functionsTSV, "\n") {
		if strings.HasPrefix(line, "# pin: ") {
			pin = line
		}
		if line != "" && !strings.HasPrefix(line, "#") {
			rows.WriteString(line)
		}
	}
	fields := strings.Fields(strings.TrimPrefix(pin, "# pin: "))
	if len(fields) != 3 {
		t.Fatalf("functions.tsv has no pin header: %q", pin)
	}
	sum := sha256.Sum256([]byte(rows.String()))
	if got := "sha256:" + hex.EncodeToString(sum[:]); got != fields[1] {
		t.Errorf("functions.tsv rows hash to %s, pin says %s: run cmd/docsfunctions rather than editing the file", got, fields[1])
	}
}

func TestBuiltinFunctions(t *testing.T) {
	var names []string
	for _, f := range builtinFunctions() {
		if f.name == "" || f.comment == "" {
			t.Errorf("incomplete row %+v", f)
		}
		names = append(names, f.name)
	}
	for _, want := range []string{"COALESCE", "FLATTEN", "PARSE_JSON", "IFF", "LEN", "LENGTH", "SYSTEM$TYPEOF"} {
		if !slices.Contains(names, want) {
			t.Errorf("missing %s", want)
		}
	}
	for _, absent := range []string{"CLASSIFY_TEXT", "MD5_NUMBER"} {
		if slices.Contains(names, absent) {
			t.Errorf("unexpected %s", absent)
		}
	}
}
