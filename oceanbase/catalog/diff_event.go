package catalog

import (
	"sort"
	"strings"
)

// MySQL SDL diff — scheduled events.
//
// Wired into DiffWithNormalizer (diff.go) so the SchemaDiff reports event changes via
// SchemaDiff.Events. Mirrors diffTables' shape: whole-catalog, taking from/to *Catalog and the
// version-fixing Normalizer. Events are database-level objects keyed by (database, lower(name)) —
// the same key eventcmds.go stores under db.Events — so the differ uses a buildEventMap index
// exactly like diffTables/buildTableMap, then sorts the result for determinism.
//
// Equality is decided by a single canonical comparison key (eventCanonicalKey) — the event
// analog of CanonicalColumn — so an event whose surface form differs from its synced
// SHOW CREATE EVENT readback but whose canonical form is identical produces no phantom diff.
// The canonicalization is grounded in what the engine actually stores (probed on live 5.7.25 +
// 8.0.32); the high-risk surfaces are:
//
//   - SCHEDULE. A recurring `EVERY <interval>` event ALWAYS gets a `STARTS '<create-time>'`
//     clause auto-injected by the engine when the user omits STARTS, and that injected timestamp
//     is (a) non-deterministic (it is NOW() at apply time) and (b) indistinguishable in the
//     stored form from a user-written STARTS — both `information_schema.EVENTS.STARTS` and
//     `SHOW CREATE EVENT` render them identically. So STARTS is dropped from the schedule key
//     (canonicalizeSchedule); see the FLAG below. ENDS, by contrast, is never auto-injected
//     (NULL when absent), so it is meaningful and kept. The EVERY interval itself is re-rendered
//     by the engine — quoted multi-field intervals lose per-field leading zeros
//     (`'12:30:00'`→`'12:30:0'`) and single-field quoted intervals are unquoted (`'05'`→`5`) — so
//     the interval literal is canonicalized to match.
//   - STATUS. ENABLE is the default and is always shown explicitly by SHOW CREATE EVENT; the
//     catalog stores "" for an unspecified status, which canonicalizes to ENABLE. DISABLE ON
//     SLAVE stays distinct from DISABLE (information_schema reports it as SLAVESIDE_DISABLED).
//   - ON COMPLETION. NOT PRESERVE is the default and is always shown explicitly; "" canonicalizes
//     to NOT PRESERVE.
//   - DEFINER. Always resolves to the executing user (`root`@`%` on the oracle box) regardless of
//     whether the user wrote it (verified on live 8.0: `DEFINER=CURRENT_USER` and an omitted
//     definer both store as “ `root`@`%` “), so it is NOT part of the schema's logical identity
//     and is dropped from the key (ignore-in-diff).
//
// FLAG 1 (recorded to the state document): because the engine stores an auto-injected STARTS
// identically to an explicit one, a release that changes ONLY an explicit STARTS on a recurring
// event is not detected as a diff. There is no signal in the stored form to distinguish the two,
// and the idempotence gate (user `EVERY 1 HOUR` ≡ stored `EVERY 1 HOUR STARTS '<now>'`) forces
// STARTS to be normalized out. This is the most authoritative behavior given the evidence.
//
// FLAG 2 (recorded to the state document): the schedule is compared as text, but MySQL EVALUATES a
// non-literal schedule expression to a literal at create time — `AT CURRENT_TIMESTAMP + INTERVAL 1
// DAY` is stored as `AT '<evaluated-ts>'`, and `EVERY 60*60 SECOND` as `EVERY 3600 SECOND`
// (verified on live 8.0.32). So a HAND-AUTHORED SDL that writes such an expression would phantom-
// diff against its synced literal form. This does not affect the real release path: bytebase syncs
// events via SHOW CREATE EVENT (oracle.md), which always yields the evaluated literal, so both
// sides of a release-path comparison are already literals. Offline expression evaluation is
// infeasible (it needs the server clock / NOW()), so this is flagged, not normalized — the same
// class of limitation as FLAG 1.
//
// implemented by omni:events breadth node
func diffEvents(from, to *Catalog, n *Normalizer) []EventDiffEntry {
	fromMap := buildEventMap(from)
	toMap := buildEventMap(to)

	var result []EventDiffEntry

	// Dropped: in from but not in to.
	for key, fromEvent := range fromMap {
		if _, ok := toMap[key]; !ok {
			result = append(result, EventDiffEntry{
				Action:   DiffDrop,
				Database: key.db,
				Name:     fromEvent.Name,
				From:     fromEvent,
			})
		}
	}

	// Added or modified: in to.
	for key, toEvent := range toMap {
		fromEvent, ok := fromMap[key]
		if !ok {
			result = append(result, EventDiffEntry{
				Action:   DiffAdd,
				Database: key.db,
				Name:     toEvent.Name,
				To:       toEvent,
			})
			continue
		}
		if eventCanonicalKey(fromEvent, n) != eventCanonicalKey(toEvent, n) {
			result = append(result, EventDiffEntry{
				Action:   DiffModify,
				Database: key.db,
				Name:     toEvent.Name,
				From:     fromEvent,
				To:       toEvent,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Database != result[j].Database {
			return result[i].Database < result[j].Database
		}
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// buildEventMap indexes every event in a catalog by (database, name), both lower-cased — the
// identity key db.Events is already keyed by. Mirrors buildTableMap.
func buildEventMap(c *Catalog) map[eventKey]*Event {
	m := make(map[eventKey]*Event)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		for name, e := range db.Events {
			if e == nil {
				continue
			}
			m[eventKey{db: toLower(db.Name), name: name}] = e
		}
	}
	return m
}

// eventKey is the identity key for an event: (database, name) lower-cased. MySQL has no
// schemas/namespaces, so the database is the qualifying scope.
type eventKey struct {
	db   string
	name string
}

// eventCanonicalKey returns a single stable comparison key for an event, folding every
// canonicalized aspect (schedule, ON COMPLETION, status, comment, body) into one string. DEFINER
// is intentionally excluded (ignore-in-diff). Equal keys mean no change. This is the event analog
// of Normalizer.CanonicalColumn.
func eventCanonicalKey(e *Event, n *Normalizer) string {
	return encodeKeyFields(
		"schedule", canonicalizeSchedule(e.Schedule),
		"oncompletion", canonicalOnCompletion(e.OnCompletion),
		"status", canonicalEventStatus(e.Enable),
		"comment", n.CanonicalComment(e.Comment),
		"body", trimSQLSpace(e.Body),
	)
}

// canonicalOnCompletion folds the ON COMPLETION clause to its canonical stored form. MySQL's
// default is NOT PRESERVE, which SHOW CREATE EVENT always renders explicitly; the catalog stores
// "" for an unspecified clause, so "" and "NOT PRESERVE" compare equal.
func canonicalOnCompletion(v string) string {
	if eqFoldStr(strings.TrimSpace(v), "PRESERVE") {
		return "PRESERVE"
	}
	return "NOT PRESERVE"
}

// canonicalEventStatus folds the ENABLE/DISABLE/DISABLE ON SLAVE status to its canonical stored
// form. ENABLE is the default and is always shown explicitly by SHOW CREATE EVENT; the catalog
// stores "" for an unspecified status, so "" and "ENABLE" compare equal.
func canonicalEventStatus(v string) string {
	switch {
	case eqFoldStr(strings.TrimSpace(v), "DISABLE"):
		return "DISABLE"
	case eqFoldStr(collapseSpaces(v), "DISABLE ON SLAVE"):
		return "DISABLE ON SLAVE"
	default:
		return "ENABLE"
	}
}

// trimSQLSpace trims surrounding ASCII whitespace from an opaque SQL fragment (event body, etc.).
// eventcmds already TrimSpace's the body on load, but a stored readback may carry different outer
// spacing, so trim again for a stable key.
func trimSQLSpace(s string) string {
	return strings.TrimSpace(s)
}

// collapseSpaces lower-cases nothing but collapses internal runs of ASCII whitespace to a single
// space and trims the ends, so "DISABLE   ON  SLAVE" and "DISABLE ON SLAVE" compare equal.
func collapseSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// canonicalizeSchedule folds an event's ON SCHEDULE clause (the raw schedule text the catalog
// stores in Event.Schedule, e.g. "EVERY 1 HOUR STARTS '2026-01-01 00:00:00'" or "AT '2030-01-01
// 00:00:00'") to a version-independent canonical key. It is the single owner of the schedule
// normalization decisions documented above:
//
//   - splits off a trailing STARTS '<ts>' clause and DROPS it (auto-injected, indistinguishable
//     from explicit — see the FLAG in the package doc);
//   - keeps a trailing ENDS '<ts>' clause (never auto-injected → meaningful);
//   - canonicalizes the EVERY interval literal so the engine's re-rendering of quoted intervals
//     (per-field leading-zero stripping; single-field unquoting) round-trips empty;
//   - upper-cases keywords and collapses whitespace so surface spacing/case never diffs.
//
// A one-shot `AT <expr>` schedule is kept whole (the engine evaluates the expression to a literal
// at create time; the literal is the schedule's identity). The clause structure MySQL guarantees
// is `EVERY <interval> [STARTS <ts>] [ENDS <ts>]` — STARTS always precedes ENDS — so a simple
// keyword scan suffices; we never need to re-parse the interval grammar.
func canonicalizeSchedule(raw string) string {
	s := collapseSpaces(raw)
	if s == "" {
		return ""
	}

	upper := strings.ToUpper(s)
	if strings.HasPrefix(upper, "EVERY") {
		body := strings.TrimSpace(s[len("EVERY"):])

		// Peel a trailing ENDS '<ts>' first (it comes after STARTS), then a STARTS '<ts>'. Because
		// STARTS precedes ENDS in the grammar, the remaining head after both peels is the interval.
		var ends string
		body, ends = peelScheduleClause(body, "ENDS")
		body, _ = peelScheduleClause(body, "STARTS") // STARTS dropped (auto-injected/explicit indistinguishable)

		interval := canonicalizeInterval(strings.TrimSpace(body))
		out := "EVERY " + interval
		if ends != "" {
			out += " ENDS " + ends
		}
		return out
	}

	if strings.HasPrefix(upper, "AT") {
		at := strings.TrimSpace(s[len("AT"):])
		return "AT " + at
	}

	// Unknown shape (defensive): return the whitespace/clause-collapsed form so identical inputs
	// still compare equal.
	return s
}

// peelScheduleClause splits a schedule body at the LAST top-level ` <keyword> ` boundary
// (case-insensitive), returning the head before it and the value after it (trimmed). If the
// keyword is absent, head == body and value == "". The match is quote-aware: occurrences inside a
// single- or double-quoted region are skipped, so a keyword that appears inside a timestamp
// literal (pathological, but defensively handled) never causes a mis-split. The keyword must be
// space-delimited so it does not match inside an interval unit name (HOUR, DAY, ...).
func peelScheduleClause(body, keyword string) (head, value string) {
	idx := lastTopLevelKeyword(body, keyword)
	if idx < 0 {
		return body, ""
	}
	needleLen := len(keyword) + 2 // surrounding spaces
	head = strings.TrimSpace(body[:idx])
	value = strings.TrimSpace(body[idx+needleLen:])
	return head, value
}

// lastTopLevelKeyword returns the byte index of the last ` <keyword> ` (case-insensitive,
// space-delimited) occurrence that is NOT inside a quoted region, or -1 if none. Quoting tracks
// both ' and " and treats a doubled quote (”) as an escaped quote within the same region — which
// is exactly how MySQL's SHOW CREATE EVENT renders quotes inside a stored schedule literal (the
// only quoting form that appears here). Backslash escapes are not interpreted; they do not occur in
// canonical stored-schedule timestamps, so this is sufficient for the inputs this scan sees.
func lastTopLevelKeyword(s, keyword string) int {
	needle := " " + keyword + " "
	last := -1
	var quote byte // 0 when outside any quote
	for i := 0; i < len(s); i++ {
		c := s[i]
		if quote != 0 {
			if c == quote {
				// A doubled quote is an escaped quote, not a close.
				if i+1 < len(s) && s[i+1] == quote {
					i++
					continue
				}
				quote = 0
			}
			continue
		}
		if c == '\'' || c == '"' {
			quote = c
			continue
		}
		// Fold case with a per-byte ASCII upcase (asciiEqualFold), NOT strings.ToUpper: ToUpper can
		// change a string's byte length (e.g. 'ı' U+0131 -> "I", 'ſ' U+017F -> "S"), which would
		// desync the upper-cased string from the s-derived index i and slice out of bounds (panic)
		// or misalign (silent mis-split). Schedule keywords are pure ASCII, so an ASCII fold is
		// correct and panic-free on any input.
		if i+len(needle) <= len(s) && asciiEqualFold(s[i:i+len(needle)], needle) {
			last = i
		}
	}
	return last
}

// asciiEqualFold reports whether a and b are equal under ASCII case folding (A–Z ↔ a–z only). It
// never changes byte length, so callers can match in the original byte-index space. b is the ASCII
// needle.
func asciiEqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if asciiUpper(a[i]) != asciiUpper(b[i]) {
			return false
		}
	}
	return true
}

// asciiUpper upper-cases a single ASCII byte; non-letters and non-ASCII bytes pass through.
func asciiUpper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - ('a' - 'A')
	}
	return c
}

// canonicalizeInterval canonicalizes an EVERY interval (the `<value> <unit>` head of a recurring
// schedule, e.g. "1 HOUR", "'1 12:30:00' DAY_SECOND", "'05' MINUTE") to match the engine's stored
// re-rendering:
//   - a single-field quoted value is unquoted and stripped of leading zeros: "'05'" -> "5";
//   - a multi-field quoted value (containing ':' or a space inside the quotes) keeps the quotes
//     and separators but strips each numeric field's leading zeros: "'1 12:30:00'" -> "'1 12:30:0'";
//   - the unit keyword is upper-cased.
//
// The value is whatever precedes the final unit token. Units may be compound (HOUR_MINUTE,
// DAY_SECOND), so the unit is taken as the last whitespace-separated token and the value is the
// rest. A bare numeric value ("100") is left as-is (no quotes, no leading zeros to strip in the
// forms the engine produces).
func canonicalizeInterval(s string) string {
	if s == "" {
		return ""
	}
	fields := strings.Fields(s)
	if len(fields) < 2 {
		// No unit token to split on; return upper-cased as a defensive fallback.
		return strings.ToUpper(s)
	}
	unit := strings.ToUpper(fields[len(fields)-1])
	value := strings.TrimSpace(strings.Join(fields[:len(fields)-1], " "))
	return canonicalizeIntervalValue(value) + " " + unit
}

// canonicalizeIntervalValue canonicalizes the numeric value portion of an EVERY interval to the
// engine's stored form (see canonicalizeInterval). It strips leading zeros from every maximal
// digit run, then unquotes a value that collapsed to a single bare integer.
func canonicalizeIntervalValue(v string) string {
	quoted := len(v) >= 2 && (v[0] == '\'' || v[0] == '"') && v[len(v)-1] == v[0]
	inner := v
	if quoted {
		inner = v[1 : len(v)-1]
	}

	stripped := stripFieldLeadingZeros(inner)

	// A quoted value that is now a single bare integer (no separators) is stored unquoted by the
	// engine ('05' -> 5). A multi-field value (still containing ':' or space) keeps its quotes.
	if quoted {
		if isAllDigits(stripped) {
			return stripped
		}
		return "'" + stripped + "'"
	}
	return stripped
}

// stripFieldLeadingZeros strips leading zeros from each maximal run of ASCII digits in s,
// preserving every non-digit separator (':', ' ', '-', '.') verbatim. A run that is all zeros
// collapses to a single "0".
func stripFieldLeadingZeros(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		c := s[i]
		if c < '0' || c > '9' {
			b.WriteByte(c)
			i++
			continue
		}
		// Consume a digit run.
		j := i
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
		}
		run := s[i:j]
		trimmed := strings.TrimLeft(run, "0")
		if trimmed == "" {
			trimmed = "0"
		}
		b.WriteString(trimmed)
		i = j
	}
	return b.String()
}

// isAllDigits reports whether s is non-empty and every byte is an ASCII digit.
func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
