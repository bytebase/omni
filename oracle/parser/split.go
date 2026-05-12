package parser

// Segment represents a portion of Oracle SQL script text delimited by a
// top-level semicolon or SQL*Plus-style slash command.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
	Kind      SegmentKind
}

// SegmentKind classifies the kind of source text represented by a Segment.
type SegmentKind int

const (
	SegmentSQL SegmentKind = iota
	SegmentSQLPlusCommand
)

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

// Split splits an Oracle SQL script into source segments. It is deliberately
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
				} else {
					lineEnd := lineEndBeforeBreak(sql, tok.End)
					nextStart := lineEndAfterBreak(sql, tok.End)
					commandStart := stmtStart
					if !onlyIgnorableSQLPlusPrefix(sql, stmtStart, tok.Loc) {
						segments = appendSegment(segments, sql, stmtStart, trimRightSpace(sql, tok.Loc))
						commandStart = lineStartOffset(sql, tok.Loc)
					}
					segments = appendSegmentWithKind(segments, sql, commandStart, lineEnd, SegmentSQLPlusCommand)
					stmtStart = nextStart
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
					end := tok.End
					if !state.endPending {
						end = tok.Loc
					}
					segments = appendSegment(segments, sql, stmtStart, end)
					stmtStart = tok.End
					state.reset()
				} else {
					state.afterPLSQLSemicolon()
				}
			} else {
				segments = appendSegment(segments, sql, stmtStart, tok.Loc)
				stmtStart = tok.End
				state.reset()
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
	splitPLSQLTypeBody
	splitPLSQLTrigger
	splitPLSQLSubprogram
	splitPLSQLCase
)

type splitPLSQLFrame struct {
	kind        splitPLSQLKind
	bodyStarted bool
	compound    bool
}

type splitState struct {
	inPLSQL bool
	frames  []splitPLSQLFrame

	pendingCreate     bool
	pendingCreateType bool
	topLevelTokens    int
	pendingSubprogram bool
	pendingCaseEnd    bool
	callSpecStarted   bool

	endPending      bool
	closedOutermost bool
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
	if tok.Type == tokHINT {
		return
	}

	if s.topLevelTokens == 0 {
		switch tok.Type {
		case kwBEGIN:
			s.startPLSQL(splitPLSQLBlock)
			s.frames[0].bodyStarted = true
			return
		case kwDECLARE, tokLABELOPEN:
			s.startPLSQL(splitPLSQLBlock)
			return
		}
	}

	if s.pendingCreateType {
		if tok.Type == kwBODY {
			s.startPLSQL(splitPLSQLTypeBody)
			return
		}
		s.pendingCreateType = false
	}

	switch tok.Type {
	case kwCREATE:
		s.pendingCreate = true
		s.topLevelTokens++
	case kwOR, kwREPLACE:
		// CREATE OR REPLACE keeps looking for the created object type.
		s.topLevelTokens++
	case kwIF, kwNOT, kwEXISTS, kwAND:
		// CREATE IF NOT EXISTS and CREATE OR REPLACE AND COMPILE/RESOLVE JAVA
		// keep looking for the created object type.
		if !s.pendingCreate {
			s.pendingCreate = false
		}
		s.topLevelTokens++
	case tokIDENT:
		switch tok.Str {
		case "EDITIONABLE", "NONEDITIONABLE", "RESOLVE", "COMPILE", "NOFORCE":
			// CREATE modifiers before the object type.
		default:
			s.pendingCreate = false
		}
		s.topLevelTokens++
	case kwPROCEDURE, kwFUNCTION, kwTRIGGER:
		if s.pendingCreate {
			if tok.Type == kwTRIGGER {
				s.startPLSQL(splitPLSQLTrigger)
			} else {
				s.startPLSQL(splitPLSQLStoredUnit)
			}
			return
		}
		s.pendingCreate = false
		s.topLevelTokens++
	case kwPACKAGE:
		if s.pendingCreate {
			s.startPLSQL(splitPLSQLPackage)
			return
		}
		s.pendingCreate = false
		s.topLevelTokens++
	case kwTYPE:
		if s.pendingCreate {
			s.pendingCreateType = true
		}
		s.topLevelTokens++
	default:
		s.pendingCreate = false
		s.topLevelTokens++
	}
}

func (s *splitState) startPLSQL(kind splitPLSQLKind) {
	s.inPLSQL = true
	s.frames = []splitPLSQLFrame{{kind: kind}}
	s.pendingCreate = false
	s.pendingCreateType = false
	s.topLevelTokens = 0
	s.pendingSubprogram = false
	s.pendingCaseEnd = false
	s.endPending = false
	s.closedOutermost = false
}

func (s *splitState) observePLSQL(tok Token) {
	if s.pendingCaseEnd {
		s.pendingCaseEnd = false
		if tok.Type == kwCASE {
			return
		}
	}

	if s.endPending {
		if tok.Type != ';' {
			return
		}
		return
	}

	if tok.Type == ';' {
		s.pendingSubprogram = false
		return
	}

	if s.pendingSubprogram {
		if tok.Type == kwIS || tok.Type == kwAS {
			s.pushFrame(splitPLSQLSubprogram, false)
			s.pendingSubprogram = false
			return
		}
	}

	if len(s.frames) == 0 {
		return
	}

	top := &s.frames[len(s.frames)-1]

	if top.kind == splitPLSQLTrigger && tok.Type == kwCOMPOUND {
		top.compound = true
		return
	}

	if !top.bodyStarted && (top.kind == splitPLSQLStoredUnit || top.kind == splitPLSQLTrigger) && splitStartsCallSpec(tok) {
		s.callSpecStarted = true
	}

	if s.canStartNestedSubprogram(tok) {
		s.pendingSubprogram = true
		return
	}

	switch tok.Type {
	case kwBEGIN:
		s.observePLSQLBegin()
	case kwIF:
		if s.inExecutablePLSQL() {
			s.pushFrame(splitPLSQLBlock, true)
		}
	case kwCASE:
		s.pushFrame(splitPLSQLCase, true)
	case kwLOOP:
		if s.inExecutablePLSQL() {
			s.pushFrame(splitPLSQLBlock, true)
		}
	case kwEND:
		s.closePLSQLFrame()
	}
}

func (s *splitState) plsqlCanEndAtSemicolon() bool {
	if s.endPending {
		return s.closedOutermost
	}
	if len(s.frames) == 1 {
		top := s.frames[0]
		if (top.kind == splitPLSQLStoredUnit || top.kind == splitPLSQLTrigger) && !top.bodyStarted && s.callSpecStarted {
			return true
		}
	}
	return false
}

func (s *splitState) afterPLSQLSemicolon() {
	s.endPending = false
	s.closedOutermost = false
	s.pendingSubprogram = false
	s.pendingCaseEnd = false
	s.callSpecStarted = false
}

func (s *splitState) pushFrame(kind splitPLSQLKind, bodyStarted bool) {
	s.frames = append(s.frames, splitPLSQLFrame{kind: kind, bodyStarted: bodyStarted})
}

func (s *splitState) observePLSQLBegin() {
	if len(s.frames) == 0 {
		return
	}
	top := &s.frames[len(s.frames)-1]
	if top.kind == splitPLSQLTrigger && top.compound {
		s.pushFrame(splitPLSQLBlock, true)
		return
	}
	if !top.bodyStarted {
		top.bodyStarted = true
		return
	}
	s.pushFrame(splitPLSQLBlock, true)
}

func (s *splitState) closePLSQLFrame() {
	closedKind := splitPLSQLNone
	if len(s.frames) > 0 {
		closedKind = s.frames[len(s.frames)-1].kind
		s.frames = s.frames[:len(s.frames)-1]
	}
	if closedKind == splitPLSQLCase {
		s.pendingCaseEnd = true
		s.endPending = false
		s.closedOutermost = false
		s.pendingSubprogram = false
		return
	}
	s.endPending = true
	s.closedOutermost = len(s.frames) == 0
	s.pendingSubprogram = false
}

func (s *splitState) inExecutablePLSQL() bool {
	if len(s.frames) == 0 {
		return false
	}
	return s.frames[len(s.frames)-1].bodyStarted
}

func (s *splitState) canStartNestedSubprogram(tok Token) bool {
	if tok.Type != kwPROCEDURE && tok.Type != kwFUNCTION {
		return false
	}
	if len(s.frames) == 0 {
		return false
	}
	top := s.frames[len(s.frames)-1]
	switch top.kind {
	case splitPLSQLPackage, splitPLSQLTypeBody:
		return !top.bodyStarted
	case splitPLSQLStoredUnit, splitPLSQLSubprogram, splitPLSQLBlock:
		return !top.bodyStarted
	default:
		return false
	}
}

func splitStartsCallSpec(tok Token) bool {
	if tok.Type == kwCALL {
		return true
	}
	return tok.Type == tokIDENT && (tok.Str == "LANGUAGE" || tok.Str == "EXTERNAL")
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
	if word == "START" && nextLineWordEquals(sql, tok.End, "WITH") {
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

func nextLineWordEquals(sql string, pos int, want string) bool {
	pos = skipHorizontalSpace(sql, pos)
	if pos >= len(sql) || sql[pos] == '\n' || sql[pos] == '\r' {
		return false
	}

	end := pos
	for end < len(sql) {
		c := sql[end]
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' || c == '$' || c == '#') {
			break
		}
		end++
	}
	if end-pos != len(want) {
		return false
	}
	for i := 0; i < len(want); i++ {
		c := sql[pos+i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		if c != want[i] {
			return false
		}
	}
	return true
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
	return appendSegmentWithKind(segments, sql, start, end, SegmentSQL)
}

func appendSegmentWithKind(segments []Segment, sql string, start, end int, kind SegmentKind) []Segment {
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
		Kind:      kind,
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
