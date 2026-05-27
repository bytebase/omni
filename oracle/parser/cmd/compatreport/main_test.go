package main

import (
	"strings"
	"testing"
)

func TestSummarizeRows(t *testing.T) {
	rows := []compatRow{
		{Fields: map[string]string{"id": "compat_select_001", "phase": "A", "family": "select", "expect": "accept", "reference_class": "syntax"}},
		{Fields: map[string]string{"id": "compat_ddl_001", "phase": "B", "family": "ddl", "expect": "reject", "reference_class": "syntax"}},
		{Fields: map[string]string{"id": "compat_admin_001", "phase": "E", "family": "admin", "expect": "unsafe", "reference_class": "privilege_catalog"}},
	}

	summary := summarizeRows(rows)

	if summary.Total != 3 {
		t.Fatalf("Total=%d, want 3", summary.Total)
	}
	if got := summary.Phases["A"]; got != 1 {
		t.Fatalf("phase A=%d, want 1", got)
	}
	if got := summary.Families["admin"]; got != 1 {
		t.Fatalf("family admin=%d, want 1", got)
	}
	if got := summary.Expectations["unsafe"]; got != 1 {
		t.Fatalf("unsafe=%d, want 1", got)
	}
	if got := summary.ReferenceClasses["privilege_catalog"]; got != 1 {
		t.Fatalf("privilege_catalog=%d, want 1", got)
	}
}

func TestReplaceMarkedSection(t *testing.T) {
	doc := strings.Join([]string{
		"before",
		reportStartMarker,
		"old report",
		reportEndMarker,
		"after",
		"",
	}, "\n")

	updated, err := replaceMarkedSection(doc, "new report\n")
	if err != nil {
		t.Fatalf("replaceMarkedSection returned error: %v", err)
	}

	want := strings.Join([]string{
		"before",
		reportStartMarker,
		"new report",
		reportEndMarker,
		"after",
		"",
	}, "\n")
	if updated != want {
		t.Fatalf("updated doc mismatch\ngot:\n%s\nwant:\n%s", updated, want)
	}
}

func TestReplaceMarkedSectionRequiresMarkers(t *testing.T) {
	if _, err := replaceMarkedSection("no markers\n", "new report\n"); err == nil {
		t.Fatal("replaceMarkedSection succeeded without markers")
	}
}

func TestParseOutcomeLine(t *testing.T) {
	log := `
    compat_oracle_ref_test.go:53: Oracle compatibility reference report: match_accept=361,match_reject=182,oracle_privilege_catalog=45,oracle_semantic_catalog=2,unsafe=18
    compat_oracle_ref_test.go:54: Oracle compatibility reference report by family: admin/match_accept=1,utility/match_reject=5
    compat_oracle_ref_test.go:60: Oracle privileged compatibility reference report: match_accept=40,oracle_privilege_catalog=5
    compat_oracle_ref_test.go:61: Oracle privileged compatibility reference report by family: admin/match_accept=38,utility/oracle_privilege_catalog=5
`

	outcomes, ok := parseOutcomeLine(log, "Oracle compatibility reference report:")
	if !ok {
		t.Fatal("parseOutcomeLine did not find reference report")
	}
	if got := outcomes["match_accept"]; got != 361 {
		t.Fatalf("match_accept=%d, want 361", got)
	}
	if got := outcomes["oracle_privilege_catalog"]; got != 45 {
		t.Fatalf("oracle_privilege_catalog=%d, want 45", got)
	}

	byFamily, ok := parseOutcomeLine(log, "Oracle compatibility reference report by family:")
	if !ok {
		t.Fatal("parseOutcomeLine did not find by-family report")
	}
	if got := byFamily["admin/match_accept"]; got != 1 {
		t.Fatalf("admin/match_accept=%d, want 1", got)
	}

	privileged, ok := parseOutcomeLine(log, "Oracle privileged compatibility reference report:")
	if !ok {
		t.Fatal("parseOutcomeLine did not find privileged reference report")
	}
	if got := privileged["match_accept"]; got != 40 {
		t.Fatalf("privileged match_accept=%d, want 40", got)
	}
}

func TestRenderMarkdown(t *testing.T) {
	report := reportData{
		Manifest: manifestSummary{
			Total:            3,
			Phases:           map[string]int{"A": 1, "B": 2},
			Families:         map[string]int{"ddl": 2, "select": 1},
			Expectations:     map[string]int{"accept": 2, "reject": 1},
			ReferenceClasses: map[string]int{"syntax": 3},
		},
		Local: map[string]int{"local_accept_match": 2, "local_reject_match": 1},
		Reference: map[string]int{
			"match_accept":             2,
			"oracle_privilege_catalog": 1,
		},
		ReferenceByFamily: map[string]int{
			"ddl/match_accept":     1,
			"select/match_reject":  1,
			"utility/match_accept": 1,
		},
		PrivilegedReference: map[string]int{
			"match_accept":             1,
			"oracle_privilege_catalog": 1,
		},
		PrivilegedReferenceByFamily: map[string]int{
			"admin/match_accept":               1,
			"utility/oracle_privilege_catalog": 1,
		},
	}

	md := renderMarkdown(report)
	for _, want := range []string{
		"| total rows | 3 |",
		"| phase A rows | 1 |",
		"| ddl | 2 |",
		"| local accept match | 2 |",
		"| oracle privilege catalog | 1 |",
		"| ddl | match accept | 1 |",
		"## Oracle Privileged Reference Report",
		"| admin | match accept | 1 |",
	} {
		if !strings.Contains(md, want) {
			t.Fatalf("renderMarkdown missing %q\n%s", want, md)
		}
	}
}
