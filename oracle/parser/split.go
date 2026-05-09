package parser

// Segment represents a portion of Oracle SQL script text delimited by a
// top-level semicolon or SQL*Plus-style slash command.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
}

// Empty returns true if the segment contains only whitespace, semicolons, and
// comments.
func (s Segment) Empty() bool {
	i := 0
	for i < len(s.Text) {
		switch s.Text[i] {
		case ' ', '\t', '\n', '\r', '\f', ';':
			i++
			continue
		case '-':
			if i+1 < len(s.Text) && s.Text[i+1] == '-' {
				i = splitSkipLineComment(s.Text, i)
				continue
			}
		case '/':
			if i+1 < len(s.Text) && s.Text[i+1] == '*' {
				next, ok := splitSkipBlockComment(s.Text, i)
				if !ok {
					return false
				}
				i = next
				continue
			}
		}
		return false
	}
	return true
}

// Split splits an Oracle SQL script into executable segments. It is deliberately
// lexical and soft-fail: invalid or incomplete SQL still returns best-effort
// statement boundaries.
func Split(sql string) []Segment {
	if sql == "" {
		return nil
	}

	lexer := NewLexer(sql)
	stmtStart := 0
	var segments []Segment
	state := splitState{}

	for {
		tok := lexer.NextToken()
		if tok.Type == tokEOF {
			break
		}

		if !state.inPLSQL {
			if cmd, ok := sqlPlusCommandAtLineStart(sql, tok); ok {
				if cmd.flush {
					if onlyIgnorableSQLPlusPrefix(sql, stmtStart, tok.Loc) {
						if tok.Type == '/' && len(segments) > 0 {
							stmtStart = lineEndBeforeBreak(sql, tok.End)
						} else {
							stmtStart = lineEndAfterBreak(sql, tok.End)
						}
					} else {
						segments = appendSegment(segments, sql, stmtStart, trimRightSpace(sql, tok.Loc))
						stmtStart = lineEndBeforeBreak(sql, tok.End)
					}
				} else if onlyIgnorableSQLPlusPrefix(sql, stmtStart, tok.Loc) {
					stmtStart = lineEndAfterBreak(sql, tok.End)
				}
				lexer.pos = lineEndAfterBreak(sql, tok.End)
				state.reset()
				continue
			}
		}

		state.observe(tok)

		switch tok.Type {
		case ';':
			if state.inPLSQL {
				if state.plsqlCanEndAtSemicolon() {
					segments = appendSegment(segments, sql, stmtStart, tok.End)
					stmtStart = tok.End
					state.reset()
				} else {
					state.afterPLSQLSemicolon()
				}
			} else {
				segments = appendSegment(segments, sql, stmtStart, tok.Loc)
				stmtStart = tok.End
			}
		case '/':
			if isSlashDelimiterLine(sql, tok.Loc, tok.End) {
				segments = appendSegment(segments, sql, stmtStart, trimRightSpace(sql, tok.Loc))
				stmtStart = slashNextSegmentStart(sql, tok.End)
				state.reset()
			}
		}

		if lexer.Err != nil {
			break
		}
	}

	segments = appendSegment(segments, sql, stmtStart, len(sql))
	if len(segments) == 0 {
		return nil
	}
	return segments
}

type splitPLSQLKind int

const (
	splitPLSQLNone splitPLSQLKind = iota
	splitPLSQLBlock
	splitPLSQLStoredUnit
	splitPLSQLPackage
)

type splitState struct {
	inPLSQL bool
	kind    splitPLSQLKind
	depth   int

	pendingCreate     bool
	pendingCreateType bool

	endPending       bool
	endFollowerWords int
}

func (s *splitState) reset() {
	*s = splitState{}
}

func (s *splitState) observe(tok Token) {
	if tok.Type == tokEOF {
		return
	}

	if !s.inPLSQL {
		s.observeTopLevel(tok)
		return
	}

	s.observePLSQL(tok)
}

func (s *splitState) observeTopLevel(tok Token) {
	if s.pendingCreateType {
		if tok.Type == kwBODY {
			s.startPLSQL(splitPLSQLPackage)
		}
		s.pendingCreateType = false
	}

	switch tok.Type {
	case kwCREATE:
		s.pendingCreate = true
	case kwOR, kwREPLACE:
		// CREATE OR REPLACE keeps looking for the created object type.
	case tokIDENT:
		if tok.Str != "EDITIONABLE" && tok.Str != "NONEDITIONABLE" {
			s.pendingCreate = false
		}
	case kwPROCEDURE, kwFUNCTION, kwTRIGGER:
		if s.pendingCreate {
			s.startPLSQL(splitPLSQLStoredUnit)
		}
	case kwPACKAGE:
		if s.pendingCreate {
			s.startPLSQL(splitPLSQLPackage)
		}
	case kwTYPE:
		if s.pendingCreate {
			s.pendingCreateType = true
		}
	case kwDECLARE:
		s.startPLSQL(splitPLSQLBlock)
	case kwBEGIN:
		s.startPLSQL(splitPLSQLBlock)
		s.depth = 1
	default:
		if tok.Type != tokHINT {
			s.pendingCreate = false
		}
	}
}

func (s *splitState) startPLSQL(kind splitPLSQLKind) {
	s.inPLSQL = true
	s.kind = kind
	s.pendingCreate = false
	s.pendingCreateType = false
	s.endPending = false
	s.endFollowerWords = 0
}

func (s *splitState) observePLSQL(tok Token) {
	if s.endPending {
		if tok.Type != ';' && isSplitWordToken(tok) {
			s.endFollowerWords++
		}
		if tok.Type != ';' {
			return
		}
	}

	switch tok.Type {
	case kwBEGIN, kwIF, kwCASE:
		s.depth++
	case kwLOOP:
		if !s.endPending {
			s.depth++
		}
	case kwEND:
		if s.depth > 0 {
			s.depth--
		}
		if s.depth == 0 {
			s.endPending = true
			s.endFollowerWords = 0
		}
	}
}

func (s *splitState) plsqlCanEndAtSemicolon() bool {
	if !s.endPending {
		return false
	}
	if s.kind == splitPLSQLPackage {
		return false
	}
	return s.endFollowerWords <= 1
}

func (s *splitState) afterPLSQLSemicolon() {
	s.endPending = false
	s.endFollowerWords = 0
}

func isSplitWordToken(tok Token) bool {
	return tok.Type == tokIDENT || tok.Type == tokQIDENT || tok.Type >= 2000
}

type sqlPlusCommand struct {
	flush bool
}

func sqlPlusCommandAtLineStart(sql string, tok Token) (sqlPlusCommand, bool) {
	lineStart := lineStartOffset(sql, tok.Loc)
	i := skipHorizontalSpace(sql, lineStart)
	if i != tok.Loc {
		return sqlPlusCommand{}, false
	}

	if tok.Type == '/' && isSlashDelimiterLine(sql, tok.Loc, tok.End) {
		return sqlPlusCommand{flush: true}, true
	}
	if tok.Type == '@' || tok.Type == '!' {
		return sqlPlusCommand{}, true
	}

	word := splitTokenWord(tok)
	if word == "" {
		return sqlPlusCommand{}, false
	}
	if isSQLPlusFlushCommand(word) {
		return sqlPlusCommand{flush: true}, true
	}
	if isSQLPlusLineCommand(word) {
		return sqlPlusCommand{}, true
	}
	return sqlPlusCommand{}, false
}

func splitTokenWord(tok Token) string {
	if tok.Type == tokIDENT || tok.Type >= 2000 {
		return tok.Str
	}
	return ""
}

func isSQLPlusFlushCommand(word string) bool {
	switch word {
	case "RUN", "R":
		return true
	default:
		return false
	}
}

func isSQLPlusLineCommand(word string) bool {
	switch word {
	case "ACC", "ACCEPT",
		"APP", "APPEND", "ARG", "ARGUMENT",
		"ARCHIVE", "ATTRIBUTE",
		"BRE", "BREAK", "BTI", "BTITLE",
		"CHANGE", "C", "CL", "CLEAR", "COL", "COLUMN", "COMP", "COMPUTE", "CONFIG", "CONN", "CONNECT", "COPY",
		"DEF", "DEFINE", "DEL", "DESC", "DESCRIBE", "DISC", "DISCONNECT",
		"ED", "EDIT", "EXEC", "EXECUTE", "EXIT",
		"GET", "HELP", "HISTORY", "HO", "HOST",
		"INPUT",
		"L", "LIST",
		"OERR",
		"PASSW", "PASSWORD", "PAU", "PAUSE", "PING", "PRI", "PRINT", "PRO", "PROMPT",
		"QUIT",
		"RECOVER", "REM", "REMARK", "REPF", "REPFOOTER", "REPH", "REPHEADER",
		"SAVE", "SET", "SHO", "SHOW", "SHUTDOWN", "SPO", "SPOOL", "STA", "START", "STARTUP", "STORE",
		"TIMI", "TIMING", "TTI", "TTITLE",
		"UNDEF", "UNDEFINE",
		"VAR", "VARIABLE",
		"WHENEVER",
		"XQUERY":
		return true
	default:
		return false
	}
}

func onlyIgnorableSQLPlusPrefix(sql string, start, end int) bool {
	seg := Segment{Text: sql[start:end], ByteStart: start, ByteEnd: end}
	return seg.Empty()
}

func lineStartOffset(sql string, pos int) int {
	if pos > len(sql) {
		pos = len(sql)
	}
	for pos > 0 && sql[pos-1] != '\n' && sql[pos-1] != '\r' {
		pos--
	}
	return pos
}

func skipHorizontalSpace(sql string, i int) int {
	for i < len(sql) {
		switch sql[i] {
		case ' ', '\t', '\f':
			i++
		default:
			return i
		}
	}
	return i
}

func lineEndBeforeBreak(sql string, pos int) int {
	for pos < len(sql) && sql[pos] != '\n' && sql[pos] != '\r' {
		pos++
	}
	return pos
}

func lineEndAfterBreak(sql string, pos int) int {
	pos = lineEndBeforeBreak(sql, pos)
	if pos < len(sql) && sql[pos] == '\r' {
		pos++
	}
	if pos < len(sql) && sql[pos] == '\n' {
		pos++
	}
	return pos
}

func appendSegment(segments []Segment, sql string, start, end int) []Segment {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	if end > len(sql) {
		end = len(sql)
	}
	seg := Segment{
		Text:      sql[start:end],
		ByteStart: start,
		ByteEnd:   end,
	}
	if seg.Empty() {
		return segments
	}
	return append(segments, seg)
}

func isSlashDelimiterLine(sql string, loc, end int) bool {
	if loc < 0 || loc >= len(sql) || sql[loc] != '/' {
		return false
	}
	for i := loc - 1; i >= 0 && sql[i] != '\n' && sql[i] != '\r'; i-- {
		if sql[i] != ' ' && sql[i] != '\t' && sql[i] != '\f' {
			return false
		}
	}
	for i := end; i < len(sql) && sql[i] != '\n' && sql[i] != '\r'; i++ {
		if sql[i] != ' ' && sql[i] != '\t' && sql[i] != '\f' {
			return false
		}
	}
	return true
}

func trimRightSpace(sql string, end int) int {
	for end > 0 {
		switch sql[end-1] {
		case ' ', '\t', '\n', '\r', '\f':
			end--
		default:
			return end
		}
	}
	return end
}

func slashNextSegmentStart(sql string, start int) int {
	i := start
	for i < len(sql) {
		switch sql[i] {
		case ' ', '\t', '\f':
			i++
		case '\r', '\n':
			return i
		default:
			return i
		}
	}
	return i
}

func splitSkipLineComment(sql string, i int) int {
	i += 2
	for i < len(sql) {
		i++
		if sql[i-1] == '\n' {
			break
		}
	}
	return i
}

func splitSkipBlockComment(sql string, i int) (int, bool) {
	i += 2
	depth := 1
	for i < len(sql) && depth > 0 {
		switch {
		case sql[i] == '/' && i+1 < len(sql) && sql[i+1] == '*':
			depth++
			i += 2
		case sql[i] == '*' && i+1 < len(sql) && sql[i+1] == '/':
			depth--
			i += 2
		default:
			i++
		}
	}
	return i, depth == 0
}
