// Package metacmd recognizes psql metacommand lines in SQL scripts. Since
// the CVE-2025-8714 fix (August 2025 point releases), pg_dump plain-format
// output brackets the dump in \restrict / \unrestrict commands, so any
// modern dump contains metacommand lines that are not SQL. A metacommand
// line is a backslash followed by a letter at the start of a line, running
// to end of line — the same recognition psql itself uses. "\." is not a
// letter and stays the COPY inline-data terminator.
package metacmd

// IsLineStart reports whether position i begins a psql metacommand line:
// i is at the start of a line (or of the input) and the bytes are a
// backslash followed by an ASCII letter. Callers must be in top-level scan
// state (not inside strings, comments, dollar-quotes, or COPY data), where
// a backslash can never begin a SQL statement.
func IsLineStart(sql string, i int, lineStart bool) bool {
	if !lineStart || i+1 >= len(sql) || sql[i] != '\\' {
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
