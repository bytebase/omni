// Package copyscan locates the inline data block that follows a
// COPY ... FROM STDIN statement in a SQL script. psql scripts and pg_dump
// plain-format output carry the data lines in the SQL stream, terminated
// by a line containing only "\.". Both the splitter and the parser consume
// the block with the same rules.
package copyscan

// SkipData consumes the inline COPY data block starting at position i
// (just past the COPY statement's semicolon): the remainder of that line,
// then data lines up to and including a line containing only "\." (psql
// recognizes the terminator only at the start of a line). Without a
// terminator the data runs to end of input, matching psql reading to EOF.
func SkipData(sql string, i int) int {
	// Finish the line the semicolon is on; data starts on the next line.
	for i < len(sql) && sql[i] != '\n' {
		i++
	}
	if i < len(sql) {
		i++ // consume the newline
	}
	for i < len(sql) {
		// i is at the start of a data line.
		if IsTerminatorLine(sql, i) {
			// Consume through the terminator line's newline (or EOF).
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
			if i < len(sql) {
				i++
			}
			return i
		}
		for i < len(sql) && sql[i] != '\n' {
			i++
		}
		if i < len(sql) {
			i++
		}
	}
	return i
}

// IsTerminatorLine reports whether the line starting at i consists of
// exactly "\." (optionally followed by a carriage return before the
// newline or end of input).
func IsTerminatorLine(sql string, i int) bool {
	if i+1 >= len(sql) || sql[i] != '\\' || sql[i+1] != '.' {
		return false
	}
	j := i + 2
	if j < len(sql) && sql[j] == '\r' {
		j++
	}
	return j >= len(sql) || sql[j] == '\n'
}
