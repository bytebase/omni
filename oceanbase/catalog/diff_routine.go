package catalog

import (
	"sort"
	"strings"
)

// MySQL SDL diff — stored routines (functions + procedures).
//
// This is the MySQL analog of PG's diffFunctions (pg/catalog/diff_function.go). MySQL keeps
// functions and procedures in two separate maps on a Database (db.Functions / db.Procedures),
// so the dispatcher (DiffWithNormalizer in diff.go) packs them into two separate SchemaDiff
// slices (SchemaDiff.Functions / SchemaDiff.Procedures), both element type RoutineDiffEntry
// (distinguished by RoutineDiffEntry.IsProcedure). This node fills BOTH via the shared
// diffRoutineMaps engine; diffFunctions and diffProcedures are thin selectors over it.
//
// Identity is (database, name) lower-cased — within a kind. A function and a procedure of the
// same name are different objects (different maps), so there is no cross-kind collision: the
// dispatcher calls diffFunctions and diffProcedures independently.
//
// Equality is decided by canonicalRoutine, which folds the stored-form-significant attributes:
// the parameter list (direction/name/type), the RETURNS type (functions), the routine BODY, and
// the in-scope characteristics (DETERMINISTIC, DATA ACCESS, SQL SECURITY, COMMENT). Two routines
// whose surface DDL differs but whose stored form is identical produce no phantom diff. Every
// canonicalization-sensitive decision is oracle-backed (see canonicalRoutine).
//
// DEFINER is NOT part of the comparison: it is an environment attribute (the creating account),
// not schema-as-code, and the loader defaults it to `root`@`%` when absent (routinecmds.go).
// Diffing it would emit spurious DROP+CREATE on every cross-deployment release. (Verified: a
// no-op release must not re-create a routine merely because the synced DEFINER differs.)
func diffFunctions(from, to *Catalog, n *Normalizer) []RoutineDiffEntry {
	return diffRoutineMaps(from, to, n, false)
}

// diffProcedures is the stored-procedure half of the routine breadth node; see diffFunctions.
func diffProcedures(from, to *Catalog, n *Normalizer) []RoutineDiffEntry {
	return diffRoutineMaps(from, to, n, true)
}

// diffRoutineMaps compares the routine maps (functions or procedures, selected by isProcedure)
// across all databases of the two catalogs and returns the per-routine changes. It mirrors the
// table differ's whole-catalog shape: walk every database present on either side, match routines
// by lower-cased (database, name), and emit Add/Drop/Modify.
func diffRoutineMaps(from, to *Catalog, n *Normalizer, isProcedure bool) []RoutineDiffEntry {
	fromMap := routineEntries(from, isProcedure)
	toMap := routineEntries(to, isProcedure)

	var result []RoutineDiffEntry

	// Dropped: in from but not in to.
	for key, fr := range fromMap {
		if _, ok := toMap[key]; !ok {
			result = append(result, RoutineDiffEntry{
				Action:      DiffDrop,
				Database:    fr.dbName,
				Name:        fr.routine.Name,
				IsProcedure: isProcedure,
				From:        fr.routine,
			})
		}
	}

	// Added or modified: in to.
	for key, tr := range toMap {
		fr, ok := fromMap[key]
		if !ok {
			result = append(result, RoutineDiffEntry{
				Action:      DiffAdd,
				Database:    tr.dbName,
				Name:        tr.routine.Name,
				IsProcedure: isProcedure,
				To:          tr.routine,
			})
			continue
		}
		if routinesChanged(fr.routine, tr.routine, n) {
			result = append(result, RoutineDiffEntry{
				Action:      DiffModify,
				Database:    tr.dbName,
				Name:        tr.routine.Name,
				IsProcedure: isProcedure,
				From:        fr.routine,
				To:          tr.routine,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if a, b := toLower(result[i].Database), toLower(result[j].Database); a != b {
			return a < b
		}
		if a, b := toLower(result[i].Name), toLower(result[j].Name); a != b {
			return a < b
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// routineRef bundles a routine with the original-case name of the database it lives in, so the
// diff entry can carry the declared database casing (matching tableIdent's case discipline) while
// the map is keyed by the lower-cased (database, name) identity.
type routineRef struct {
	dbName  string
	routine *Routine
}

// routineEntries indexes every function (or procedure, per isProcedure) in every database of the
// catalog by lower-cased "database\x00name". The NUL separator keeps a routine `b` in database
// `a` from colliding with a routine in database `a\x00b` (database/name both lower-cased).
func routineEntries(c *Catalog, isProcedure bool) map[string]routineRef {
	m := make(map[string]routineRef)
	if c == nil {
		return m
	}
	for _, db := range c.Databases() {
		routineMap := db.Functions
		if isProcedure {
			routineMap = db.Procedures
		}
		for _, r := range routineMap {
			if r == nil {
				continue
			}
			key := toLower(db.Name) + "\x00" + toLower(r.Name)
			m[key] = routineRef{dbName: db.Name, routine: r}
		}
	}
	return m
}

// routinesChanged reports whether two same-identity routines differ, comparing their canonical
// keys (the MySQL analog of PG's per-object change check). canonicalRoutine folds every
// stored-form-significant attribute into one collision-free key, so this differ never
// re-implements a per-attribute comparison.
func routinesChanged(a, b *Routine, n *Normalizer) bool {
	return canonicalRoutine(a, n) != canonicalRoutine(b, n)
}

// canonicalRoutine returns a single stable comparison key for a routine, folding the attributes
// that survive into MySQL's stored form. The comparison is whole-routine: MySQL has no
// CREATE OR REPLACE and no ALTER for a routine's body/params/RETURNS, so any change to those is a
// DROP+CREATE (migration_routine.go decides DROP+CREATE vs the characteristic-only ALTER).
//
// Folded fields and their oracle grounding (probed on live 5.7 :13307 and 8.0 :13306):
//   - params: direction + name + type, in order. The type strings are taken verbatim from the
//     loader (routinecmds.go formatParamType), which renders the same form the engine stores —
//     so a 5.7 `int(11)` vs 8.0 `int` divergence is already baked into each side's loaded value
//     and never cross-compared (Diff is always single-version). Order is significant.
//   - returns (functions only): the RETURNS type string, likewise loader-rendered.
//   - body: the routine body is stored BYTE-FOR-BYTE verbatim by the engine (internal
//     whitespace, tabs, blank lines all preserved — verified), and the loader stores
//     strings.TrimSpace(BodyText), so two routines with the same source body have identical Body.
//     Compared as an opaque string.
//   - characteristics: DEFAULT-FOLDED via canonicalRoutineCharacteristics — a missing or
//     explicitly-default characteristic collapses to its canonical default, because SHOW CREATE
//     DROPS default characteristics (CONTAINS SQL, SQL SECURITY DEFINER, NOT DETERMINISTIC are all
//     omitted from the readback — verified). Without this fold a user form that spells out a
//     default would phantom-diff against the readback that omits it.
//
// NOT folded (deliberately): DEFINER (environment, see diffFunctions) and LANGUAGE (always SQL,
// never echoed distinctly by SHOW CREATE — verified — so it carries no stored-form signal).
func canonicalRoutine(r *Routine, n *Normalizer) string {
	if r == nil {
		return ""
	}
	params := make([]string, 0, len(r.Params))
	for _, p := range r.Params {
		params = append(params, canonicalRoutineParam(p, n))
	}
	return encodeKeyFields(
		"isproc", boolKey(r.IsProcedure),
		"params", strings.Join(params, ","),
		"returns", canonicalRoutineReturns(r, n),
		"chars", canonicalRoutineCharacteristics(r, n),
		"body", r.Body,
	)
}

// canonicalRoutineParam encodes one routine parameter: direction (procedures only), name
// (lower-cased — MySQL parameter names are case-insensitive identifiers), and type. The type is
// Direction is folded to upper-case; an empty direction (function param, or a procedure IN that
// the user omitted) collapses to "IN" — MySQL's default parameter direction — so `p(x INT)` and
// `p(IN x INT)` compare equal.
//
// The TYPE is canonicalized the SAME way as the RETURNS type (canonicalRoutineType), NOT by a flat
// strings.ToUpper of the whole string: upper-casing the entire type would fold ENUM/SET member
// VALUES too (collapsing ENUM('a','b') and ENUM('A','b')), making a legitimate member-case change
// invisible to the differ — the comparison-side mirror of the render-side corruption. Routing
// through parseColumnType + canonicalType (→ canonicalEnumSuffix) keeps the keyword case-stable
// while preserving member-value case, matching the column path.
func canonicalRoutineParam(p *RoutineParam, n *Normalizer) string {
	if p == nil {
		return ""
	}
	dir := strings.ToUpper(strings.TrimSpace(p.Direction))
	if dir == "" {
		dir = "IN"
	}
	return encodeKeyFields(
		"dir", dir,
		"name", toLower(p.Name),
		"type", canonicalRoutineType(p.TypeName, n),
	)
}

// canonicalRoutineType canonicalizes a routine parameter/return type string: the base type goes
// through normalize-core's parseColumnType + canonicalType (folding int display width, aliases,
// decimal precision, and — crucially — preserving ENUM/SET member-value case via
// canonicalEnumSuffix), and any "CHARSET <cs> [COLLATE <coll>]" suffix is split off and re-appended
// lower-cased. This is the single type-key builder shared by params and RETURNS, so both compare a
// type the same case-correct way.
func canonicalRoutineType(typeStr string, n *Normalizer) string {
	base, suffix := splitRoutineTypeCharset(typeStr)
	if strings.TrimSpace(base) == "" {
		return suffix
	}
	canonical := n.canonicalType(parseColumnType(base))
	if suffix != "" {
		canonical += " " + suffix
	}
	return canonical
}

// canonicalRoutineReturns returns the canonical RETURNS-type key for a function, "" for a
// procedure (procedures have no RETURNS).
//
// Unlike the PARAMETER list — which SHOW CREATE stores VERBATIM as the user typed it, identically
// on 5.7 and 8.0 (verified: a param `b INT(11)` reads back `b INT(11)` on both, so params need no
// version folding) — the RETURNS type IS rendered in the engine's version-canonical form: 5.7
// shows `RETURNS int(11)`, 8.0 shows `RETURNS int` (the same integer-display-width divergence
// columns have). So a user form `RETURNS INT` must canonicalize equal to a 5.7 readback's
// `RETURNS int(11)`. The base type is therefore routed through normalize-core's type canonicalizer
// (parseColumnType + canonicalType — the SAME int-width / alias / decimal-precision rules
// CanonicalColumnType uses), NOT compared as a raw string.
//
// The CHARSET/COLLATE suffix (the loader appends "CHARSET <cs>" for string return types, matching
// SHOW CREATE) is split off, lower-cased, and re-appended so it still participates in the key
// (an 8.0 `CHARSET utf8mb4` return differs from a 5.7 `CHARSET latin1` one). Both sides are loaded
// the same way per version, so an identical function self-compares equal.
func canonicalRoutineReturns(r *Routine, n *Normalizer) string {
	if r.IsProcedure || strings.TrimSpace(r.Returns) == "" {
		return ""
	}
	return canonicalRoutineType(r.Returns, n)
}

// splitRoutineTypeCharset separates a routine type string into its base type and a normalized
// "CHARSET <cs> [COLLATE <coll>]" suffix (lower-cased keyword, lower-cased charset/collation). The
// loader renders string return types as e.g. "varchar(100) CHARSET utf8mb4"; the base
// ("varchar(100)") goes through the type canonicalizer while the charset rides along as a stable
// key fragment. A type with no CHARSET/COLLATE returns (type, "").
//
// The CHARSET/COLLATE clause always follows the base type and any parenthesized argument list, so
// the search starts AFTER the last ')'. This prevents mis-splitting a type whose ENUM/SET member
// literal contains the substring " charset "/" collate " (e.g. enum('x charset y')) — the member
// text lives inside the parentheses, before the search window.
func splitRoutineTypeCharset(typeStr string) (base, suffix string) {
	s := strings.TrimSpace(typeStr)
	// Search for the suffix keyword only past any member/argument list.
	searchFrom := strings.LastIndexByte(s, ')') + 1 // 0 when there is no ')'
	low := strings.ToLower(s)
	rel := strings.Index(low[searchFrom:], " charset ")
	if rel < 0 {
		// COLLATE may appear without an explicit CHARSET in rare forms; handle it too.
		rel = strings.Index(low[searchFrom:], " collate ")
		if rel < 0 {
			return s, ""
		}
	}
	idx := searchFrom + rel
	base = strings.TrimSpace(s[:idx])
	// Normalize the suffix: collapse internal whitespace and lower-case keywords + values.
	suffix = strings.Join(strings.Fields(strings.ToLower(s[idx:])), " ")
	return base, suffix
}

// canonicalRoutineCharacteristics builds the default-folded characteristics key. The four
// in-scope characteristics are each collapsed to a canonical value, defaulting a missing key to
// MySQL's default (verified: SHOW CREATE omits defaults), so a user form that states a default
// and the readback that drops it produce the same key:
//   - DETERMINISTIC: "YES" | "NO"; missing → "NO" (NOT DETERMINISTIC is the default).
//   - DATA ACCESS: one of CONTAINS SQL | NO SQL | READS SQL DATA | MODIFIES SQL DATA; missing →
//     "CONTAINS SQL" (the default, which SHOW CREATE never echoes).
//   - SQL SECURITY: "DEFINER" | "INVOKER" (upper-cased); missing → "DEFINER" (the default).
//   - COMMENT: the comment content (decoded by the loader; routed through CanonicalComment);
//     missing → "".
//
// LANGUAGE is intentionally absent (always SQL, no stored-form signal).
func canonicalRoutineCharacteristics(r *Routine, n *Normalizer) string {
	det := r.Characteristics["DETERMINISTIC"]
	if det == "" {
		det = "NO"
	}
	dataAccess := r.Characteristics["DATA ACCESS"]
	if dataAccess == "" {
		dataAccess = "CONTAINS SQL"
	}
	security := strings.ToUpper(strings.TrimSpace(r.Characteristics["SQL SECURITY"]))
	if security == "" {
		security = "DEFINER"
	}
	comment := n.CanonicalComment(r.Characteristics["COMMENT"])
	return encodeKeyFields(
		"det", strings.ToUpper(det),
		"data", strings.ToUpper(strings.TrimSpace(dataAccess)),
		"sec", security,
		"comment", comment,
	)
}

// routineAlterSuffices reports whether a routine MODIFY can be applied with ALTER
// FUNCTION/PROCEDURE alone, rather than DROP+CREATE. ALTER can change ONLY the characteristics
// SQL SECURITY, COMMENT, and DATA ACCESS (verified against both live engines: ALTER accepts
// COMMENT / SQL SECURITY / {CONTAINS SQL|NO SQL|READS SQL DATA|MODIFIES SQL DATA} / LANGUAGE SQL,
// and REJECTS DETERMINISTIC with a syntax error). So ALTER suffices iff the routines are equal
// EXCEPT for those three characteristics — i.e. params, RETURNS, body, and DETERMINISTIC are all
// unchanged. Anything else requires DROP+CREATE (no CREATE OR REPLACE in MySQL).
//
// Implemented as: equal once SQL SECURITY / COMMENT / DATA ACCESS are masked out of BOTH canonical
// keys. If the masked keys are equal, every remaining (DROP+CREATE-only) field already matched, so
// the only differences live in the ALTER-able set.
func routineAlterSuffices(from, to *Routine, n *Normalizer) bool {
	return canonicalRoutineExceptAlterable(from, n) == canonicalRoutineExceptAlterable(to, n)
}

// canonicalRoutineExceptAlterable is canonicalRoutine with the ALTER-able characteristics
// (SQL SECURITY, COMMENT, DATA ACCESS) removed, leaving only the DROP+CREATE-significant fields
// (params, RETURNS, body, DETERMINISTIC). Two routines that differ only in ALTER-able
// characteristics share this key.
func canonicalRoutineExceptAlterable(r *Routine, n *Normalizer) string {
	if r == nil {
		return ""
	}
	params := make([]string, 0, len(r.Params))
	for _, p := range r.Params {
		params = append(params, canonicalRoutineParam(p, n))
	}
	det := r.Characteristics["DETERMINISTIC"]
	if det == "" {
		det = "NO"
	}
	return encodeKeyFields(
		"isproc", boolKey(r.IsProcedure),
		"params", strings.Join(params, ","),
		"returns", canonicalRoutineReturns(r, n),
		"det", strings.ToUpper(det),
		"body", r.Body,
	)
}

// boolKey renders a bool as a stable single-char key fragment.
func boolKey(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
