package driver

import (
	"strings"
	"unicode"

	crrawanalysis "github.com/bytebase/omni/cockroachdb/analysis"
)

// sqlStatement is one statement of a script with its text and byte offsets.
// CockroachDB's parser supplies the text it read; the offsets are located in
// the original script so an execution log can point at it.
type sqlStatement struct {
	Text      string
	ByteStart int
	ByteEnd   int
}

// splitStatements splits script into statements the way Bytebase's CockroachDB
// driver does: through CockroachDB's own parser, with a trailing semicolon
// written on each statement.
func splitStatements(script string) ([]sqlStatement, error) {
	trimmed := strings.TrimRightFunc(script, isSpaceOrSemicolon) + ";"
	parsed, err := crrawanalysis.ParseScript(trimmed)
	if err != nil {
		return nil, err
	}
	statements := make([]sqlStatement, 0, len(parsed))
	offset := 0
	for _, stmt := range parsed {
		stmtSQL := stmt.Text()
		text := stmtSQL
		if !strings.HasSuffix(text, ";") {
			text += ";"
		}
		start, end := locateStatement(script, stmtSQL, offset)
		if end > len(script) {
			end = len(script)
		}
		statements = append(statements, sqlStatement{Text: text, ByteStart: start, ByteEnd: end})
		offset = end
	}
	return statements, nil
}

// locateStatement finds sql in script at or after offset, so a statement's
// offsets are where it was written. The parser may have rewritten trailing
// whitespace, so a statement it cannot find continues where the last one
// ended.
func locateStatement(script, sql string, offset int) (int, int) {
	if sql == "" || offset >= len(script) {
		return offset, offset
	}
	if i := strings.Index(script[offset:], sql); i >= 0 {
		start := offset + i
		return start, start + len(sql)
	}
	return offset, offset
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}

// isPlSQLBlock reports whether stmt is a PL/pgSQL anonymous code block (DO
// statement). Bytebase runs one as a single statement outside a transaction.
func isPlSQLBlock(stmt string) bool {
	return crrawanalysis.IsDoBlock(stmt)
}
