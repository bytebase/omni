package catalog

import (
	"strings"
	"testing"
)

// In-process, hermetic tests for the event differ + its canonicalization. They assert the
// SchemaDiff.Events content (right DiffAction per event) for representative change forms and lock
// in the canonical-key normalization decisions (schedule STARTS-stripping, interval re-rendering,
// status/ON COMPLETION defaults, DEFINER ignore) without a live engine. The oracle-backed
// idempotence + apply-correctness proofs live in migration_event_oracle_test.go
// (correctness-protocol.md gates 1 & 2).

// findEvent locates an event entry in a SchemaDiff by name (lower-cased).
func findEvent(d *SchemaDiff, name string) *EventDiffEntry {
	for i := range d.Events {
		if strings.EqualFold(d.Events[i].Name, name) {
			return &d.Events[i]
		}
	}
	return nil
}

const evDBDDL = "CREATE DATABASE app DEFAULT CHARSET=utf8mb4;\nUSE app;\n"

func diffEventsOnly(t *testing.T, fromSQL, toSQL string) *SchemaDiff {
	t.Helper()
	from := loadCat(t, evDBDDL+fromSQL)
	to := loadCat(t, evDBDDL+toSQL)
	return DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
}

// ---- canonicalizeSchedule: the headline normalization -----------------------

func TestCanonicalizeSchedule(t *testing.T) {
	cases := []struct {
		name string
		a, b string // two schedules that must canonicalize EQUAL
	}{
		// STARTS auto-injected (synced) vs omitted (user SDL) → equal.
		{"every-starts-stripped",
			"EVERY 1 HOUR STARTS '2026-06-26 12:32:56'",
			"EVERY 1 HOUR"},
		// Two different STARTS on the same EVERY schedule → equal (STARTS dropped entirely).
		{"every-different-starts-equal",
			"EVERY 1 HOUR STARTS '2025-01-01 00:00:00'",
			"EVERY 1 HOUR STARTS '2025-06-01 00:00:00'"},
		// ENDS kept; STARTS between interval and ENDS stripped.
		{"every-ends-kept-starts-stripped",
			"EVERY 1 DAY STARTS '2026-06-26 12:32:56' ENDS '2030-01-01 00:00:00'",
			"EVERY 1 DAY ENDS '2030-01-01 00:00:00'"},
		// Quoted single-field interval unquoted + leading-zero stripped.
		{"interval-single-field-unquote",
			"EVERY '05' MINUTE STARTS '2026-06-26 12:32:56'",
			"EVERY 5 MINUTE"},
		// Quoted multi-field interval: per-field leading zeros stripped, quotes kept.
		{"interval-multi-field-leading-zeros",
			"EVERY '2:03:04' HOUR_SECOND",
			"EVERY '2:3:4' HOUR_SECOND"},
		{"interval-day-second-leading-zeros",
			"EVERY '1 12:30:00' DAY_SECOND",
			"EVERY '1 12:30:0' DAY_SECOND"},
		// Whitespace / case insensitivity.
		{"whitespace-case",
			"every   1    hour",
			"EVERY 1 HOUR"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ka := canonicalizeSchedule(c.a)
			kb := canonicalizeSchedule(c.b)
			if ka != kb {
				t.Errorf("schedules must canonicalize equal:\n  a=%q -> %q\n  b=%q -> %q", c.a, ka, c.b, kb)
			}
		})
	}
}

func TestCanonicalizeSchedule_Distinct(t *testing.T) {
	// Schedules that must NOT collapse to the same key (real differences).
	cases := []struct {
		name string
		a, b string
	}{
		{"different-interval-value", "EVERY 1 HOUR", "EVERY 2 HOUR"},
		{"different-interval-unit", "EVERY 1 HOUR", "EVERY 1 DAY"},
		{"every-vs-at", "EVERY 1 HOUR", "AT '2030-01-01 00:00:00'"},
		{"different-at", "AT '2030-01-01 00:00:00'", "AT '2031-01-01 00:00:00'"},
		{"ends-presence", "EVERY 1 DAY ENDS '2030-01-01 00:00:00'", "EVERY 1 DAY"},
		{"different-ends", "EVERY 1 DAY ENDS '2030-01-01 00:00:00'", "EVERY 1 DAY ENDS '2031-01-01 00:00:00'"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if canonicalizeSchedule(c.a) == canonicalizeSchedule(c.b) {
				t.Errorf("schedules must be distinct but both -> %q (a=%q b=%q)",
					canonicalizeSchedule(c.a), c.a, c.b)
			}
		})
	}
}

func TestCanonicalizeSchedule_QuoteAwarePeel(t *testing.T) {
	// A timestamp literal that pathologically contains the bare word STARTS/ENDS must NOT be
	// mis-split by the clause peeler (quote-aware). The whole quoted value stays in the schedule.
	got := canonicalizeSchedule("EVERY 1 DAY STARTS '2025-01-01 ENDS HERE 00:00:00'")
	// STARTS is stripped (auto/explicit indistinguishable), but the ENDS inside the quote must not
	// be promoted to a real ENDS clause — so there is no " ENDS " token in the canonical form.
	if strings.Contains(got, "ENDS") {
		t.Errorf("quoted ENDS must not become a clause; got %q", got)
	}
	if got != "EVERY 1 DAY" {
		t.Errorf("want canonical %q, got %q", "EVERY 1 DAY", got)
	}

	// A real ENDS following a quoted STARTS that contains 'ENDS' must still be peeled.
	got2 := canonicalizeSchedule("EVERY 1 DAY STARTS '2025-01-01 ENDS X 00:00:00' ENDS '2030-01-01 00:00:00'")
	if got2 != "EVERY 1 DAY ENDS '2030-01-01 00:00:00'" {
		t.Errorf("want real trailing ENDS kept, got %q", got2)
	}
}

func TestCanonicalizeSchedule_UTF8Safe(t *testing.T) {
	// Regression: lastTopLevelKeyword must not panic or mis-split when the schedule text contains a
	// rune whose ToUpper changes byte length ('ı' U+0131 -> "I", 'ſ' U+017F -> "S"). Folding with a
	// per-byte ASCII upcase keeps index space and slice space aligned. These must simply not panic
	// and must not invent a clause from the non-ASCII bytes.
	inputs := []string{
		"EVERY 1 DAY ENDS '2030' " + strings.Repeat("ı", 30),
		"EVERY 1 DAY STARTS 'ıııııııı' ENDS '2030-01-01 00:00:00'",
		"EVERY 1 HOUR ſſſ",
		"AT 'ı2030-01-01 00:00:00ſ'",
	}
	for _, in := range inputs {
		// Must not panic.
		got := canonicalizeSchedule(in)
		_ = got
		// scheduleHasEnds must also not panic on the same inputs (it shares the scan).
		_ = scheduleHasEnds(in)
	}
	// The real trailing ENDS after a non-ASCII STARTS literal must still be detected.
	if !scheduleHasEnds("EVERY 1 DAY STARTS 'ıııı' ENDS '2030-01-01 00:00:00'") {
		t.Errorf("real trailing ENDS after non-ASCII STARTS must be detected")
	}
	if got := canonicalizeSchedule("EVERY 1 DAY STARTS 'ıııı' ENDS '2030-01-01 00:00:00'"); got != "EVERY 1 DAY ENDS '2030-01-01 00:00:00'" {
		t.Errorf("want ENDS kept after non-ASCII STARTS, got %q", got)
	}
}

func TestCanonicalEventStatusAndCompletion(t *testing.T) {
	// "" defaults to ENABLE / NOT PRESERVE; explicit forms preserved; DISABLE ON SLAVE distinct.
	if canonicalEventStatus("") != "ENABLE" {
		t.Errorf("empty status must canonicalize to ENABLE, got %q", canonicalEventStatus(""))
	}
	if canonicalEventStatus("enable") != "ENABLE" {
		t.Errorf("enable -> ENABLE")
	}
	if canonicalEventStatus("DISABLE") != "DISABLE" {
		t.Errorf("DISABLE preserved")
	}
	if canonicalEventStatus("disable on slave") != "DISABLE ON SLAVE" {
		t.Errorf("disable on slave -> DISABLE ON SLAVE, got %q", canonicalEventStatus("disable on slave"))
	}
	if canonicalEventStatus("DISABLE") == canonicalEventStatus("DISABLE ON SLAVE") {
		t.Errorf("DISABLE and DISABLE ON SLAVE must be distinct")
	}
	if canonicalOnCompletion("") != "NOT PRESERVE" {
		t.Errorf("empty ON COMPLETION must canonicalize to NOT PRESERVE")
	}
	if canonicalOnCompletion("PRESERVE") != "PRESERVE" {
		t.Errorf("PRESERVE preserved")
	}
	if canonicalOnCompletion("NOT PRESERVE") != "NOT PRESERVE" {
		t.Errorf("NOT PRESERVE preserved")
	}
}

// ---- structural diff (Add / Drop / Modify / no-op) --------------------------

func TestDiffEvents_SelfEmpty(t *testing.T) {
	schemas := []string{
		"CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1;",
		"CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2025-01-01 00:00:00' ENDS '2026-01-01 00:00:00' ON COMPLETION PRESERVE DISABLE COMMENT 'c' DO SET @q = 1;",
		"CREATE EVENT e ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO SET @y = 2;",
	}
	for _, s := range schemas {
		c := loadCat(t, evDBDDL+s)
		if d := DiffWithNormalizer(c, c, NormalizerFor(MySQL80)); !d.IsEmpty() {
			t.Errorf("self-diff not empty for %q: events=%d", s, len(d.Events))
		}
	}
}

func TestDiffEvents_StartsInjectionNoOp(t *testing.T) {
	// The headline idempotence case in-process: user form (no STARTS) vs synced form (auto STARTS)
	// must diff EMPTY.
	from := loadCat(t, evDBDDL+"CREATE EVENT e ON SCHEDULE EVERY 1 HOUR STARTS '2026-06-26 12:32:56' DO SET @x = 1;")
	to := loadCat(t, evDBDDL+"CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1;")
	if d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80)); !d.IsEmpty() {
		t.Errorf("STARTS-injection no-op must diff empty, got %d event changes", len(d.Events))
	}
}

func TestDiffEvents_Add(t *testing.T) {
	d := diffEventsOnly(t, "", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1;")
	e := findEvent(d, "e")
	if e == nil || e.Action != DiffAdd {
		t.Fatalf("want ADD event e, got %+v", d.Events)
	}
}

func TestDiffEvents_Drop(t *testing.T) {
	d := diffEventsOnly(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x = 1;", "")
	e := findEvent(d, "e")
	if e == nil || e.Action != DiffDrop {
		t.Fatalf("want DROP event e, got %+v", d.Events)
	}
}

func TestDiffEvents_Modify(t *testing.T) {
	cases := []struct {
		name     string
		from, to string
	}{
		{"schedule", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=1;", "CREATE EVENT e ON SCHEDULE EVERY 2 HOUR DO SET @x=1;"},
		{"status", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ENABLE DO SET @x=1;", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x=1;"},
		{"completion", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE DO SET @x=1;", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR ON COMPLETION PRESERVE DO SET @x=1;"},
		{"comment", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'a' DO SET @x=1;", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'b' DO SET @x=1;"},
		{"body", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=1;", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=2;"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := diffEventsOnly(t, c.from, c.to)
			e := findEvent(d, "e")
			if e == nil || e.Action != DiffModify {
				t.Fatalf("want MODIFY event e for %s, got %+v", c.name, d.Events)
			}
		})
	}
}

func TestDiffEvents_DefinerIgnored(t *testing.T) {
	// Two events differing ONLY in DEFINER must NOT diff (DEFINER is ignore-in-diff).
	from := loadCat(t, evDBDDL+"CREATE DEFINER=`a`@`%` EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=1;")
	to := loadCat(t, evDBDDL+"CREATE DEFINER=`b`@`localhost` EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=1;")
	if d := DiffWithNormalizer(from, to, NormalizerFor(MySQL80)); !d.IsEmpty() {
		t.Errorf("DEFINER-only change must not diff, got %d event changes", len(d.Events))
	}
}

// ---- generate: the rendered DDL shape ---------------------------------------

func TestGenerateEventDDL_Shapes(t *testing.T) {
	// ADD → CREATE EVENT; DROP → DROP EVENT IF EXISTS; MODIFY → ALTER EVENT. COMMENT is always
	// rendered (COMMENT '' when empty) so an ALTER that clears a comment resets it.
	add := genEventPlan(t, "", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=1;")
	mustContain(t, add, "CREATE EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x=1")

	drop := genEventPlan(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=1;", "")
	mustContain(t, drop, "DROP EVENT IF EXISTS `e`")

	mod := genEventPlan(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=1;", "CREATE EVENT e ON SCHEDULE EVERY 2 HOUR DISABLE DO SET @x=1;")
	mustContain(t, mod, "ALTER EVENT `e` ON SCHEDULE EVERY 2 HOUR ON COMPLETION NOT PRESERVE DISABLE COMMENT '' DO SET @x=1")

	// MODIFY that CLEARS a comment must render COMMENT '' (regression for the comment-removal bug).
	clr := genEventPlan(t, "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR COMMENT 'old' DO SET @x=1;", "CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=1;")
	mustContain(t, clr, "ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE ENABLE COMMENT '' DO SET @x=1")
}

func TestGenerateEventDDL_CommentEscaping(t *testing.T) {
	// Identifier + comment escaping must survive into the rendered DDL.
	add := genEventPlan(t, "", "CREATE EVENT `weird``name` ON SCHEDULE EVERY 1 HOUR COMMENT 'it''s a test' DO SET @x=1;")
	mustContain(t, add, "CREATE EVENT `weird``name`")
	mustContain(t, add, "COMMENT 'it''s a test'")
}

func TestGenerateEventDDL_EmptyOnNoOp(t *testing.T) {
	from := loadCat(t, evDBDDL+"CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DO SET @x=1;")
	diff := DiffWithNormalizer(from, from, NormalizerFor(MySQL80))
	plan := GenerateMigrationWithNormalizer(from, from, diff, NormalizerFor(MySQL80))
	if plan.SQL() != "" {
		t.Errorf("no-op plan must be empty, got:\n%s", plan.SQL())
	}
}

func TestGenerateEventDDL_DropPhaseBeforeCreate(t *testing.T) {
	// A plan that drops one event and creates another: the DROP must sort before the CREATE
	// (PhasePre < PhaseMain).
	from := loadCat(t, evDBDDL+"CREATE EVENT old ON SCHEDULE EVERY 1 HOUR DO SET @x=1;")
	to := loadCat(t, evDBDDL+"CREATE EVENT new ON SCHEDULE EVERY 1 HOUR DO SET @x=1;")
	diff := DiffWithNormalizer(from, to, NormalizerFor(MySQL80))
	plan := GenerateMigrationWithNormalizer(from, to, diff, NormalizerFor(MySQL80))
	sql := plan.SQL()
	di := strings.Index(sql, "DROP EVENT")
	ci := strings.Index(sql, "CREATE EVENT")
	if di < 0 || ci < 0 || di > ci {
		t.Errorf("DROP must precede CREATE in:\n%s", sql)
	}
}

func genEventPlan(t *testing.T, fromSQL, toSQL string) string {
	t.Helper()
	var from, to *Catalog
	if strings.TrimSpace(fromSQL) == "" {
		from = New()
	} else {
		from = loadCat(t, evDBDDL+fromSQL)
	}
	if strings.TrimSpace(toSQL) == "" {
		to = New()
	} else {
		to = loadCat(t, evDBDDL+toSQL)
	}
	n := NormalizerFor(MySQL80)
	diff := DiffWithNormalizer(from, to, n)
	return GenerateMigrationWithNormalizer(from, to, diff, n).SQL()
}

func mustContain(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("expected to contain:\n  %q\ngot:\n  %q", needle, haystack)
	}
}
