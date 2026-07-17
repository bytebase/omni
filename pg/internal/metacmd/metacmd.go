// Package metacmd recognizes psql metacommand lines in SQL scripts. Since
// the CVE-2025-8714 fix (August 2025 point releases), pg_dump plain-format
// output brackets the dump in \restrict / \unrestrict commands, so any
// modern dump contains metacommand lines that are not SQL. A metacommand
// line is a backslash followed by a letter at the start of a line, running
// to end of line — the same recognition psql itself uses. "\." is not a
// letter and stays the COPY inline-data terminator.
package metacmd

// IsMetaCommand reports whether position i begins a psql metacommand: a
// backslash followed by an ASCII letter, at ANY top-level position. psql
// recognizes backslash commands anywhere outside quotes — not only at line
// starts (engine-verified: SELECT 1; \echo MIDLINE executes) — and a
// top-level backslash is never valid SQL (scan.l's operator charset has no
// backslash), so this can never consume legal statements. Being
// position-context-free is what makes re-splitting stable: a segment's own
// text always reproduces the decision made on the full script. Callers
// must be in top-level scan state (not inside strings, comments,
// dollar-quotes, or COPY data).
func IsMetaCommand(sql string, i int) bool {
	if i+1 >= len(sql) || sql[i] != '\\' {
		return false
	}
	b := sql[i+1]
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

// SkipLine consumes the metacommand line starting at position i, returning
// the position just past its newline (or end of input).
func SkipLine(sql string, i int) int {
	for i < len(sql) && sql[i] != '\n' {
		i++
	}
	if i < len(sql) {
		i++
	}
	return i
}
