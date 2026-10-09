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
	// SegmentJavaSource is a CREATE JAVA ... AS statement: SQL up to AS, then
	// Java source text. It executes like SegmentSQL, but its text must not be
	// read with the SQL lexer past AS: Java's "--" decrement, "//" comments,
	// and char literals lex as SQL comments and strings.
	SegmentJavaSource
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
				prefixEmpty := onlyIgnorableSQLPlusPrefix(sql, stmtStart, tok.Loc)
				if !prefixEmpty && !(cmd.flush && tok.Type == '/') {
					state.observe(tok)
					continue
				}
				if cmd.flush {
					if prefixEmpty {
						if tok.Type == '/' && len(segments) > 0 {
							stmtStart = lineEndBeforeBreak(sql, tok.End)
						} else {
							stmtStart = lineEndAfterBreak(sql, tok.End)
						}
					} else {
						segments = appendSegment(segments, sql, stmtStart, trimRightSpace(sql, tok.Loc))
						stmtStart = lineEndBeforeBreak(sql, tok.End)
					}
					lexer.pos = lineEndAfterBreak(sql, tok.End)
					state.reset()
					continue
				} else {
					if prefixEmpty || cmd.terminatesBufferedSQL {
						lineEnd := lineEndBeforeBreak(sql, tok.End)
						nextStart := lineEndAfterBreak(sql, tok.End)
						commandStart := stmtStart
						if !prefixEmpty {
							segments = appendSegment(segments, sql, stmtStart, trimRightSpace(sql, tok.Loc))
							commandStart = lineStartOffset(sql, tok.Loc)
						}
						segments = appendSegmentWithKind(segments, sql, commandStart, lineEnd, SegmentSQLPlusCommand)
						stmtStart = nextStart
						lexer.pos = lineEndAfterBreak(sql, tok.End)
						state.reset()
						continue
					}
				}
			}
		}

		state.observe(tok)

		if !state.inPLSQL && state.javaHead.observe(tok) {
			// tok is the AS of CREATE JAVA ... AS. The Java source after it
			// runs to the next line holding only "/"; its ';' ends nothing.
			end, next := javaSourceEnd(sql, tok.End)
			segments = appendSegmentWithKind(segments, sql, stmtStart, end, SegmentJavaSource)
			stmtStart = slashNextSegmentStart(sql, next)
			lexer.pos = next
			state.reset()
			continue
		}

		switch tok.Type {
		case ';':
			if state.inPLSQL {
				if state.plsqlCanEndAtSemicolon() {
					end := tok.End
					if !state.endPending && !state.callSpecKeepsSemicolon() {
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
			if (!state.inPLSQL || state.plsqlCanEndAtSlashDelimiter()) && isSlashDelimiterLine(sql, tok.Loc, tok.End) {
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

	// A stored unit's head runs to the IS|AS that opens its implementation;
	// headDepth tracks the parameter list's parentheses before it. A nested
	// subprogram frame is pushed at its IS|AS.
	headDepth int
	isAs      bool
	// afterIsAs and callSpecWord look at the two tokens after IS|AS for a
	// call spec, as isCallSpecStartTokens decides it.
	afterIsAs    bool
	callSpecWord string
	// callSpec marks an implementation that is a call spec: it has no END,
	// and its ';' ends it.
	callSpec bool
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

	javaHead javaSourceHead
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
		if n := len(s.frames); n > 1 && s.frames[n-1].callSpec {
			// A nested call spec has no END; its ';' closes it.
			s.frames = s.frames[:n-1]
		}
		return
	}

	if s.pendingSubprogram {
		if tok.Type == kwIS || tok.Type == kwAS {
			s.pushFrame(splitPLSQLSubprogram, false)
			top := &s.frames[len(s.frames)-1]
			top.isAs = true
			top.afterIsAs = true
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
	if !top.bodyStarted {
		switch top.kind {
		case splitPLSQLStoredUnit, splitPLSQLSubprogram:
			s.observeSubprogramHead(top, tok)
		case splitPLSQLTrigger:
			if len(s.frames) == 1 && tok.Type == kwCALL {
				s.callSpecStarted = true
			}
		}
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

// observeSubprogramHead follows a stored unit or nested subprogram up to its
// implementation. A call spec right after IS|AS replaces BEGIN ... END and
// ends at its ';'. WRAPPED replaces IS|AS in a wrapped stored unit. Both are
// recognized only outside the parameter list, so a parameter or variable
// named LANGUAGE, EXTERNAL, or WRAPPED does not end the unit early.
func (s *splitState) observeSubprogramHead(top *splitPLSQLFrame, tok Token) {
	switch {
	case top.callSpecWord != "":
		word := top.callSpecWord
		top.callSpecWord = ""
		if isCallSpecStartTokens(Token{Type: tokIDENT, Str: word}, tok) {
			top.callSpec = true
			if len(s.frames) == 1 {
				s.callSpecStarted = true
			}
		}
	case top.afterIsAs:
		top.afterIsAs = false
		if tok.Type == tokIDENT && (tok.Str == "LANGUAGE" || tok.Str == "EXTERNAL") {
			top.callSpecWord = tok.Str
		}
	case !top.isAs:
		switch {
		case tok.Type == '(':
			top.headDepth++
		case tok.Type == ')':
			top.headDepth--
		case top.headDepth > 0:
		case tok.Type == kwIS || tok.Type == kwAS:
			top.isAs = true
			top.afterIsAs = true
		case tok.Type == tokIDENT && tok.Str == "WRAPPED" && len(s.frames) == 1:
			s.callSpecStarted = true
		}
	}
}

// callSpecKeepsSemicolon reports whether the stored unit being split is
// implemented by a call spec, whose ';' belongs to the unit: Oracle compiles
// it with PLS-00103 without one. A trigger's CALL routine and a wrapped unit
// end before their ';' instead, which Oracle would reject in the text.
func (s *splitState) callSpecKeepsSemicolon() bool {
	return len(s.frames) == 1 && s.frames[0].callSpec
}

func (s *splitState) plsqlCanEndAtSemicolon() bool {
	if s.endPending {
		return s.closedOutermost
	}
	if s.callSpecStarted {
		return true
	}
	return false
}

func (s *splitState) plsqlCanEndAtSlashDelimiter() bool {
	return s.callSpecStarted
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

type sqlPlusCommand struct {
	flush                 bool
	terminatesBufferedSQL bool
}

func sqlPlusCommandAtLineStart(sql string, tok Token) (sqlPlusCommand, bool) {
	lineStart := lineStartOffset(sql, tok.Loc)
	i := skipHorizontalSpace(sql, lineStart)
	if i != tok.Loc {
		return sqlPlusCommand{}, false
	}

	if tok.Type == '/' && isSlashDelimiterLine(sql, tok.Loc, tok.End) {
		return sqlPlusCommand{flush: true, terminatesBufferedSQL: true}, true
	}
	if tok.Type == '@' || tok.Type == '!' {
		return sqlPlusCommand{terminatesBufferedSQL: true}, true
	}

	word := splitTokenWord(tok)
	if word == "" {
		return sqlPlusCommand{}, false
	}
	if sqlPlusCommandWordIsQualifiedIdentifier(sql, tok.End) {
		return sqlPlusCommand{}, false
	}
	if (word == "R" || word == "RUN") && !onlyHorizontalSpaceUntilLineEnd(sql, tok.End) {
		return sqlPlusCommand{}, false
	}
	if isOracleSetStatement(word, sql, tok.End) {
		return sqlPlusCommand{}, false
	}
	if isSQLPlusFlushCommand(word) {
		return sqlPlusCommand{flush: true, terminatesBufferedSQL: true}, true
	}
	if isSQLPlusLineCommand(word) {
		return sqlPlusCommand{terminatesBufferedSQL: isSQLPlusLineCommandThatTerminatesBufferedSQL(word, sql, tok.End)}, true
	}
	return sqlPlusCommand{}, false
}

func splitTokenWord(tok Token) string {
	if tok.Type == tokIDENT || tok.Type >= 2000 {
		return tok.Str
	}
	return ""
}

func sqlPlusCommandWordIsQualifiedIdentifier(sql string, pos int) bool {
	return pos < len(sql) && sql[pos] == '.'
}

func onlyHorizontalSpaceUntilLineEnd(sql string, pos int) bool {
	for pos < len(sql) {
		switch sql[pos] {
		case ' ', '\t', '\f':
			pos++
		case '\r', '\n':
			return true
		default:
			return false
		}
	}
	return true
}

func isOracleSetStatement(word, sql string, pos int) bool {
	if word != "SET" {
		return false
	}
	next := nextWordOnLine(sql, pos)
	switch next {
	case "TRANSACTION", "ROLE", "CONSTRAINT", "CONSTRAINTS":
		return true
	default:
		return false
	}
}

func nextWordOnLine(sql string, pos int) string {
	pos = skipHorizontalSpace(sql, pos)
	start := pos
	for pos < len(sql) {
		c := sql[pos]
		if c == '\n' || c == '\r' || !(c == '_' || c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z') {
			break
		}
		pos++
	}
	if pos == start {
		return ""
	}

	buf := make([]byte, pos-start)
	for i := start; i < pos; i++ {
		c := sql[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		buf[i-start] = c
	}
	return string(buf)
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

func isSQLPlusLineCommandThatTerminatesBufferedSQL(word, sql string, pos int) bool {
	switch word {
	case "ACC", "ACCEPT",
		"BRE", "BREAK", "BTI", "BTITLE",
		"COL", "COLUMN", "COMP", "COMPUTE",
		"DEF", "DEFINE",
		"HO", "HOST",
		"PRI", "PRINT", "PRO", "PROMPT",
		"REM", "REMARK",
		"SHO", "SHOW", "SPO", "SPOOL",
		"TTI", "TTITLE",
		"UNDEF", "UNDEFINE",
		"VAR", "VARIABLE",
		"WHENEVER":
		return true
	case "CONN", "CONNECT":
		return isSQLPlusConnectCommandThatTerminatesBufferedSQL(sql, pos)
	case "STA", "START":
		return nextWordOnLine(sql, pos) != "WITH"
	case "SET":
		return isSQLPlusSetCommandThatTerminatesBufferedSQL(sql, pos)
	default:
		return false
	}
}

func isSQLPlusConnectCommandThatTerminatesBufferedSQL(sql string, pos int) bool {
	switch nextWordOnLine(sql, pos) {
	case "BY", "TO":
		return false
	default:
		return true
	}
}

func isSQLPlusSetCommandThatTerminatesBufferedSQL(sql string, pos int) bool {
	switch nextWordOnLine(sql, pos) {
	case "APPINFO", "ARRAYSIZE", "AUTOCOMMIT", "AUTOPRINT", "AUTORECOVERY", "AUTOTRACE",
		"BLOCKTERMINATOR", "CMDSEP", "COLINVISIBLE", "COLSEP", "CONCAT", "COPYCOMMIT",
		"COPYTYPECHECK", "DEF", "DEFINE", "DESCRIBE", "ECHO", "EDITFILE", "EMBEDDED", "ERRORLOGGING",
		"ESCAPE", "ESCCHAR", "EXITCOMMIT", "FEEDBACK", "FLAGGER", "FLUSH", "HEADING",
		"HEADSEP", "INSTANCE", "LINESIZE", "LOBOFFSET", "LOGSOURCE", "LONG", "LONGCHUNKSIZE",
		"MARKUP", "NEWPAGE", "NULL", "NUMFORMAT", "NUMWIDTH", "PAGESIZE", "PAUSE",
		"RECSEP", "RECSEPCHAR", "SCAN", "SERVEROUT", "SERVEROUTPUT", "SHIFTINOUT", "SHOWMODE", "SQLBLANKLINES",
		"SQLCASE", "SQLCONTINUE", "SQLNUMBER", "SQLPLUSCOMPATIBILITY", "SQLPREFIX",
		"SQLPROMPT", "SQLTERMINATOR", "SUFFIX", "TAB", "TERMOUT", "TIME", "TIMING",
		"TRIMOUT", "TRIMSPOOL", "UNDERLINE", "VERIFY", "WRAP":
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

// splitSkipBlockComment advances past an Oracle block comment. Oracle block
// comments do not nest — the first */ terminates the comment (engine-verified;
// see lexBlockCommentOrHint) — so the scan must stay byte-identical to the
// lexer's rule or the splitter and parser would disagree on statement
// boundaries around comments containing a stray /*.
func splitSkipBlockComment(sql string, i int) (int, bool) {
	i += 2
	for i < len(sql) {
		if sql[i] == '*' && i+1 < len(sql) && sql[i+1] == '/' {
			return i + 2, true
		}
		i++
	}
	return i, false
}
