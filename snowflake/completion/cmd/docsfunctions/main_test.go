package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFunctionRows(t *testing.T) {
	page := `<table><tbody>
<tr><td><a href="/f/abs">ABS</a></td><td>Returns the <code>absolute</code> value.</td><td>x</td></tr>
<tr><td><a href="/f/len">LENGTH, LEN</a></td><td>Returns the length.</td><td>x</td></tr>
<tr><td><a href="/f/seq">SEQ1 / SEQ2</a></td><td>Returns a sequence.</td><td>x</td></tr>
<tr><td><a href="/f/c">AI_COMPLETE (Single string)</a></td><td>Generates a response.</td><td>x</td></tr>
<tr><td><a href="/f/c2">AI_COMPLETE (Prompt object)</a></td><td>Generates again.</td><td>x</td></tr>
<tr><td><a href="/f/ct">CLASSIFY_TEXT (SNOWFLAKE.CORTEX)</a></td><td>Classifies.</td><td>x</td></tr>
<tr><td><a href="/f/avg">AVG (system data metric function)</a></td><td>Averages.</td><td>x</td></tr>
<tr><td><a href="/f/md5">MD5_NUMBER — Obsoleted</a></td><td>Gone.</td><td>x</td></tr>
<tr><td><a href="/f/w">SYSTEM$WHITELIST — Deprecated</a></td><td>Old.</td><td>x</td></tr>
<tr><td><a href="/f/b">[ NOT ] BETWEEN</a></td><td>Range.</td><td>x</td></tr>
<tr><td><a href="/f/t">TO_TIMESTAMP / TO_TIMESTAMP_*</a></td><td>Converts.</td><td>x</td></tr>
</tbody></table>`
	rows, err := functionRows(page)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range rows {
		got = append(got, r[0]+"="+r[1])
	}
	want := []string{
		"ABS=Returns the absolute value.",
		"LENGTH=Returns the length.",
		"LEN=Returns the length.",
		"SEQ1=Returns a sequence.",
		"SEQ2=Returns a sequence.",
		"AI_COMPLETE=Generates a response.",
		"SYSTEM$WHITELIST=Old.",
		"TO_TIMESTAMP=Converts.",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if _, err := functionRows("<html></html>"); err == nil {
		t.Error("a page without rows is not an error")
	}
}

func TestRenderKeepsPinnedDate(t *testing.T) {
	out := filepath.Join(t.TempDir(), "functions.tsv")
	rows := [][2]string{{"ABS", "Returns the absolute value."}}
	first, err := render(out, rows, "2026-01-02")
	if err != nil {
		t.Fatal(err)
	}
	pin, ok := readPin(first)
	if !ok || pin.date != "2026-01-02" {
		t.Fatalf("pin %+v in\n%s", pin, first)
	}
	if err := writeFile(out, first); err != nil {
		t.Fatal(err)
	}
	again, err := render(out, rows, "")
	if err != nil {
		t.Fatal(err)
	}
	if again != first {
		t.Errorf("unchanged rows moved the pin:\n%s", again)
	}
}
