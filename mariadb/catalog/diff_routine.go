package catalog

import (
	"sort"
	"strings"
)

// diffFunctions compares the stored functions of two catalogs, keyed by
// (database, name). Two routines are equal when canonicalRoutine gives the same
// key: parameters, RETURNS type, body, and the characteristics DETERMINISTIC, DATA
// ACCESS, SQL SECURITY and COMMENT. The definer is not compared: it is an
// environment's account, not schema.
func diffFunctions(from, to *Catalog, n *Normalizer) []RoutineDiffEntry {
	return diffRoutineMaps(from, to, n, false)
}

// diffProcedures is diffFunctions for stored procedures.
func diffProcedures(from, to *Catalog, n *Normalizer) []RoutineDiffEntry {
	return diffRoutineMaps(from, to, n, true)
}

// diffRoutineMaps compares the functions or, per isProcedure, the procedures of
// every database of the two catalogs, matched by lower-cased (database, name).
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

// routineRef pairs a routine with its database's declared name, which the diff
// entry carries; the map is keyed by the lower-cased (database, name).
type routineRef struct {
	dbName  string
	routine *Routine
}

// routineEntries indexes every function (or procedure, per isProcedure) of every
// database of the catalog by lower-cased "database\x00name". The NUL separator
// keeps a routine `b` in database `a` from colliding with a routine in database
// `a\x00b`.
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

// routinesChanged reports whether two same-identity routines differ by
// canonicalRoutine.
func routinesChanged(a, b *Routine, n *Normalizer) bool {
	return canonicalRoutine(a, n) != canonicalRoutine(b, n)
}

// canonicalRoutine is the comparison key of a routine: its parameters in order,
// the RETURNS type of a function, the characteristics with defaults folded
// (canonicalRoutineCharacteristics), and the body, which is compared as written.
// The definer is left out. routineAlterSuffices says which differences an ALTER
// can apply.
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

// canonicalRoutineParam encodes one parameter: its direction, with an omitted
// direction read as IN; its lower-cased name; and its type through
// canonicalRoutineType.
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

// canonicalRoutineType canonicalizes a parameter or RETURNS type: the base type
// through canonicalRoutineTypeText, which keeps ENUM and SET member case, and any
// CHARSET/COLLATE suffix, lower-cased and appended.
func canonicalRoutineType(typeStr string, n *Normalizer) string {
	base, suffix := splitRoutineTypeCharset(typeStr)
	if strings.TrimSpace(base) == "" {
		return suffix
	}
	canonical := canonicalRoutineTypeText(base)
	if suffix != "" {
		canonical += " " + suffix
	}
	return canonical
}

// canonicalRoutineReturns returns the canonical RETURNS type of a function, or ""
// for a procedure. The CHARSET/COLLATE suffix the loader appends to a string type
// takes part in the key.
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

// canonicalRoutineCharacteristics builds the characteristics key, reading a
// missing value as its default so a form that states a default and one that omits
// it compare equal:
//   - DETERMINISTIC: "YES" | "NO"; missing is "NO".
//   - DATA ACCESS: CONTAINS SQL | NO SQL | READS SQL DATA | MODIFIES SQL DATA;
//     missing is "CONTAINS SQL".
//   - SQL SECURITY: "DEFINER" | "INVOKER"; missing is "DEFINER".
//   - COMMENT: the comment text, through CanonicalComment.
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

// routineAlterSuffices reports whether a routine change can be applied with ALTER
// FUNCTION/PROCEDURE rather than a drop and create. ALTER changes only SQL
// SECURITY, COMMENT and DATA ACCESS, so it suffices when the routines are equal
// once those three are masked out of both keys.
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
