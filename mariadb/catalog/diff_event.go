package catalog

import (
	"sort"
	"strings"
)

// diffEvents compares the events of two catalogs, keyed by (database, name). Two
// events are equal when eventCanonicalKey matches: schedule, ON COMPLETION,
// status, comment and body. The definer is not compared: it is an environment's
// account, not schema.
//
// The schedule is compared as text. A recurring schedule's STARTS clause is left
// out of the key: when none is written the server records the creation time, which
// the stored form does not tell apart from an explicit one, so a change to only an
// explicit STARTS is not reported. A non-literal schedule expression is compared
// as written, not evaluated.
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

// buildEventMap indexes every event of a catalog by (database, name),
// lower-cased.
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

// eventKey is the identity key of an event: (database, name), lower-cased.
type eventKey struct {
	db   string
	name string
}

// eventCanonicalKey returns the comparison key of an event: schedule, ON
// COMPLETION, status, comment and body. The definer is left out.
func eventCanonicalKey(e *Event, n *Normalizer) string {
	return encodeKeyFields(
		"schedule", canonicalizeSchedule(e.Schedule),
		"oncompletion", canonicalOnCompletion(e.OnCompletion),
		"status", canonicalEventStatus(e.Enable),
		"comment", n.CanonicalComment(e.Comment),
		"body", trimSQLSpace(e.Body),
	)
}

// canonicalOnCompletion reads an unspecified ON COMPLETION as NOT PRESERVE, the
// default.
func canonicalOnCompletion(v string) string {
	if eqFoldStr(strings.TrimSpace(v), "PRESERVE") {
		return "PRESERVE"
	}
	return "NOT PRESERVE"
}

// canonicalEventStatus folds the status to ENABLE, DISABLE or DISABLE ON SLAVE,
// reading an unspecified status as ENABLE, the default.
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

// trimSQLSpace trims surrounding whitespace from an opaque SQL fragment such as
// an event body.
func trimSQLSpace(s string) string {
	return strings.TrimSpace(s)
}

// collapseSpaces lower-cases nothing but collapses internal runs of ASCII whitespace to a single
// space and trims the ends, so "DISABLE   ON  SLAVE" and "DISABLE ON SLAVE" compare equal.
func collapseSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// canonicalizeSchedule folds an event's schedule text (e.g. "EVERY 1 HOUR STARTS
// '2026-01-01 00:00:00'" or "AT '2030-01-01 00:00:00'") into a comparison key:
//   - a trailing STARTS clause is dropped (see diffEvents);
//   - a trailing ENDS clause is kept;
//   - the EVERY interval goes through canonicalizeInterval;
//   - keywords are upper-cased and whitespace is collapsed.
//
// A one-shot `AT` schedule is kept whole. The clause order is
// `EVERY <interval> [STARTS <ts>] [ENDS <ts>]`, so a keyword scan suffices.
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
		body, _ = peelScheduleClause(body, "STARTS") // STARTS dropped; see diffEvents

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

// lastTopLevelKeyword returns the byte index of the last ` <keyword> `
// (case-insensitive, space-delimited) occurrence that is not inside a quoted
// region, or -1 if none. Quoting tracks both ' and " and treats a doubled quote as
// an escaped quote within the same region. Backslash escapes are not interpreted;
// a stored schedule timestamp does not contain one.
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

// canonicalizeInterval folds an EVERY interval (the `<value> <unit>` head of a
// recurring schedule, e.g. "1 HOUR", "'1 12:30:00' DAY_SECOND", "'05' MINUTE"):
//   - a single-field quoted value is unquoted and stripped of leading zeros:
//     "'05'" -> "5";
//   - a multi-field quoted value (containing ':' or a space inside the quotes)
//     keeps the quotes and separators but each numeric field loses its leading
//     zeros: "'1 12:30:00'" -> "'1 12:30:0'";
//   - the unit keyword is upper-cased.
//
// Units may be compound (HOUR_MINUTE, DAY_SECOND), so the unit is the last
// whitespace-separated token and the value is the rest. A bare numeric value is
// left as-is.
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

// canonicalizeIntervalValue strips leading zeros from every maximal digit run of
// an EVERY interval value, then unquotes a value that collapsed to a single bare
// integer. Only a number or a quoted interval literal ('1:05') is rewritten; an
// expression keeps its text, whose string literals are values of their own.
func canonicalizeIntervalValue(v string) string {
	quoted := len(v) >= 2 && (v[0] == '\'' || v[0] == '"') && v[len(v)-1] == v[0]
	inner := v
	if quoted {
		inner = v[1 : len(v)-1]
	}
	if strings.Trim(inner, "0123456789:.- ") != "" {
		return v
	}

	stripped := stripFieldLeadingZeros(inner)

	// A quoted value that is now a single bare integer is unquoted ('05' -> 5); a
	// multi-field value (still containing ':' or a space) keeps its quotes.
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
