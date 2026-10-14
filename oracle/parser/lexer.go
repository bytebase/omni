// Package parser implements an Oracle PL/SQL lexer and recursive descent parser.
package parser

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// Token types returned by the lexer.
const (
	tokEOF = 0
)

// Non-keyword token types.
const (
	tokICONST = iota + 1000
	tokFCONST
	tokSCONST     // single-quoted string literal
	tokNCHARLIT   // N'...' national character literal
	tokIDENT      // unquoted identifier (lowercased)
	tokQIDENT     // "double-quoted" identifier (case-preserved)
	tokBIND       // :name or :1 bind variable
	tokCONCAT     // ||
	tokASSIGN     // :=
	tokASSOC      // =>
	tokDOTDOT     // ..
	tokEXPON      // **
	tokLESSEQ     // <=
	tokGREATEQ    // >=
	tokNOTEQ      // != or <> or ~= or ^=
	tokLABELOPEN  // <<
	tokLABELCLOSE // >>
	tokOP         // generic operator
	tokHINT       // /*+ hint text */
	tokDIRECTIVE  // $IF, $THEN, $ELSIF, $ELSE, $END, or $ERROR; never returned by NextToken
	tokINQUIRY    // $$name inquiry directive, a literal the compiler substitutes
)

// Token represents a lexical token.
type Token struct {
	Type int    // token type
	Str  string // string value
	Ival int64  // integer value
	Loc  int    // byte offset in source
	End  int    // exclusive end byte offset
}

// Lexer implements an Oracle SQL/PL/SQL lexer.
type Lexer struct {
	input string
	pos   int
	start int
	end   int // exclusive bound of the lexed range within input

	// ccDepth counts the conditional compilation $IF directives whose
	// selected branch holds the tokens being returned (see NextToken), and
	// ccLevels holds them, the innermost first. ccUndecided counts those
	// whose selection depends on an undecided one. Oracle 23ai sets no
	// bound on the nesting: 5000 levels compile.
	ccDepth     int
	ccLevels    *ccLevel
	ccUndecided int

	// ccFacts records what the selection so far made of the conditions the
	// text does not settle, for a later $IF with the same condition to
	// select alike (see ccFact). The lexer's copies share it.
	ccFacts *ccFacts

	// ccChoices and ccDecisions let the parser try each branch selection a
	// condition the text does not settle leaves open (see
	// Parser.parseStmtSelections). ccChoices says how many of an $IF's
	// undecided branches to pass over; ccDecisions, shared by the lexer's
	// copies, records each $IF that held an undecided branch, in order.
	ccChoices   *ccChoices
	ccDecisions *[]ccDecision

	// ccLenient is set on a lexer that must not stop at a misplaced
	// directive: a pre-scan that starts inside the text, at the token the
	// parser stands on, where the branch holding it may end in an $ELSIF,
	// $ELSE, or $END with no $IF before it; and Split, which is soft-fail.
	// Such a directive skips to the end of its construct, and an $ERROR is
	// skipped. The parser's own lexer checks the structure.
	ccLenient bool

	// ccSplit is set on Split's lexer, where a directive is stray only in
	// malformed input: Split drops it and lexes the code after it, an
	// $ELSIF or $ELSE branch or an $ERROR message too, to keep every
	// statement boundary.
	ccSplit bool

	// ccStray is set by Split outside a PL/SQL unit, where any directive is
	// stray, as Oracle preprocesses a unit alone: the lexer drops it and
	// lexes on, a branch's code and condition too, to keep every statement
	// boundary.
	ccStray bool

	// ccSkipEnd is where the last directive construct NextToken skipped
	// ends, for Split to keep directives after a unit's END with the unit.
	ccSkipEnd int

	// prevType is the type of the last token lexed, which says whether a
	// /*+ ... */ comment is a hint. ccRawDirective is set while nextRaw is
	// in the condition of an $IF or $ELSIF or the message of an $ERROR,
	// which leave prevType as the directive does.
	prevType       int
	ccRawDirective bool

	// literalLimit, when set, is the most bytes a string literal holds: the
	// parser's lexer fails past it, in code a directive skips too (PLS-00172
	// and ORA-01704 on Oracle 23ai past 32767), where Split's and
	// completion's lexers lex on.
	literalLimit int

	// numberLimit is set on the parser's lexer, which fails at a numeric
	// literal past NUMBER's range (numberOverflows) in the code it returns;
	// Oracle 23ai does not check the code a directive skips.
	numberLimit bool

	// cc, when set, records the preprocessor tokens the lexer meets. The
	// copies of the lexer that look ahead share it.
	cc *ccRecord

	// Error
	Err error
}

// ccRecord records where conditional compilation was used, for the checks
// the parser makes per statement: a directive or inquiry outside PL/SQL,
// and the condition of each $IF and $ELSIF the lexer evaluated.
type ccRecord struct {
	directives []int       // start offsets of $IF ... $ERROR and $$name
	conditions []nodes.Loc // the expressions between an evaluated $IF or $ELSIF and its $THEN
	binds      []int       // start offsets of the bind variables inside a directive
}

// noteDirective records a preprocessor token at loc. Offsets arrive in
// increasing order, whichever copy of the lexer meets one first.
func (r *ccRecord) noteDirective(loc int) {
	if r != nil && (len(r.directives) == 0 || loc > r.directives[len(r.directives)-1]) {
		r.directives = append(r.directives, loc)
	}
}

// noteBind records a bind variable inside a directive at loc.
func (r *ccRecord) noteBind(loc int) {
	if r != nil && (len(r.binds) == 0 || loc > r.binds[len(r.binds)-1]) {
		r.binds = append(r.binds, loc)
	}
}

// ccDecision is the selection the lexer made at an $IF one of whose
// conditions the text does not settle: it passed over that many undecided
// branches, and more says another selection could pass over one more.
type ccDecision struct {
	at     int
	passed int
	more   bool
}

// noteDecision records the selection made at the $IF at offset at, once
// whichever copy of the lexer makes it first.
func (l *Lexer) noteDecision(at, passed int, more bool) {
	if d := l.ccDecisions; d != nil && (len(*d) == 0 || at > (*d)[len(*d)-1].at) {
		*d = append(*d, ccDecision{at: at, passed: passed, more: more})
	}
}

// openDecisions returns the offsets of the $IF directives whose decisions
// the selections of the open levels depend on, the latest for each.
func (l *Lexer) openDecisions() []int {
	if l.ccUndecided == 0 {
		return nil
	}
	var open []int
	for level := l.ccLevels; level != nil; level = level.outer {
		if level.dep >= 0 {
			open = append(open, level.dep)
		}
	}
	return open
}

// ccFact is what a selection made of a condition the text does not settle:
// that it holds (is TRUE) or not, as decided at the $IF at offset by.
// Oracle 23ai compiles a unit under one environment, where a condition has
// one value: $IF $$x $THEN NULL; $ELSE garbage $END $IF $$x $THEN garbage
// $ELSE NULL; $END fails whatever $$x is.
type ccFact struct {
	holds bool
	by    int
}

// ccFacts holds the facts a selection made, by condition, identified by its
// tokens. A fact follows from the selection's choices and where its $IF
// stands, so the lexer's copies that look ahead share the table: a fact one
// records is the one the others would.
type ccFacts struct {
	byCond map[string]ccFact
}

// clone returns a copy of the facts, for a lexer that records no decisions
// to read and add to apart: one that recorded a fact without its decision
// would leave the parser no decision to vary.
func (f *ccFacts) clone() *ccFacts {
	if f == nil {
		return nil
	}
	c := &ccFacts{byCond: make(map[string]ccFact, len(f.byCond))}
	for cond, fact := range f.byCond {
		c.byCond[cond] = fact
	}
	return c
}

// fact returns the fact about the condition cond.
func (l *Lexer) fact(cond string) (ccFact, bool) {
	if l.ccFacts == nil {
		return ccFact{}, false
	}
	f, ok := l.ccFacts.byCond[cond]
	return f, ok
}

// noteFact records that the condition cond holds or not, as decided at the
// $IF at offset by, unless a fact about it is recorded.
func (l *Lexer) noteFact(cond string, holds bool, by int) {
	if cond == "" {
		return
	}
	if l.ccFacts == nil {
		l.ccFacts = &ccFacts{byCond: map[string]ccFact{}}
	}
	if _, ok := l.ccFacts.byCond[cond]; !ok {
		l.ccFacts.byCond[cond] = ccFact{holds: holds, by: by}
	}
}

// ccChoices is a branch selection to make: the $IF at offset at passes over
// pass undecided branches, those decided records keep the choices recorded
// there, in order of offset, and every other passes over none. A failed
// selection's decisions are recorded once, and each selection after it
// shares them.
type ccChoices struct {
	decided  []ccDecision
	at, pass int
}

// passFor returns how many undecided branches the $IF at offset at passes
// over.
func (c *ccChoices) passFor(at int) int {
	if c == nil {
		return 0
	}
	if at == c.at {
		return c.pass
	}
	i := sort.Search(len(c.decided), func(i int) bool { return c.decided[i].at >= at })
	if i < len(c.decided) && c.decided[i].at == at {
		return c.decided[i].passed
	}
	return 0
}

// noteCondition records the condition of an evaluated $IF or $ELSIF.
func (r *ccRecord) noteCondition(loc nodes.Loc) {
	if r != nil && (len(r.conditions) == 0 || loc.Start > r.conditions[len(r.conditions)-1].Start) {
		r.conditions = append(r.conditions, loc)
	}
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	return NewLexerRange(input, 0, len(input))
}

// NewLexerRange creates a lexer over input[start:end] that reports token
// offsets relative to input, so a segment of a larger script lexes in place
// with absolute positions and no copy of the text.
func NewLexerRange(input string, start, end int) *Lexer {
	if end > len(input) {
		end = len(input)
	}
	if start < 0 {
		start = 0
	}
	if start > end {
		start = end
	}
	return &Lexer{
		input: input,
		pos:   start,
		end:   end,
	}
}

// NextToken returns the next token from the input, after PL/SQL
// conditional compilation selected the branch of each $IF that Oracle
// compiles. Oracle 23ai decides by static values; where the text settles a
// condition (staticCondition), the lexer selects the branch Oracle does,
// and otherwise the first branch whose condition the text does not settle
// as FALSE or NULL, else the $ELSE branch. It returns the tokens of the
// selected branch and skips the rest, whose code Oracle 23ai does not
// parse, and the conditions, which the parser checks. It checks the
// directives' structure as Oracle 23ai does, in skipped branches too: a
// $THEN, $ELSIF, $ELSE, or $END out of place is PLS-00176, an $IF or
// $ERROR without its $END is PLS-00177, and an $ERROR in a branch
// certainly selected is PLS-00179. An inquiry directive $$name is an
// identifier.
func (l *Lexer) NextToken() Token {
	// A directive and the code it removes leave the token before it last:
	// SELECT $IF TRUE $THEN $END /*+ FULL(t) */ keeps its hint.
	restore, prevType := false, 0
	for {
		if restore {
			l.prevType, restore = prevType, false
		}
		tokPrev := l.prevType
		tok := l.nextRaw()
		if tok.Type == tokEOF {
			if l.ccDepth > 0 && l.Err == nil && !l.ccLenient {
				return l.failDirective(tok, "'$IF' preprocessor directive does not end properly")
			}
			return tok
		}
		if tok.Type != tokDIRECTIVE {
			if tok.Type == tokBIND && l.ccDepth > 0 {
				l.cc.noteBind(tok.Loc)
			}
			if (tok.Type == tokICONST || tok.Type == tokFCONST) && l.numberLimit && numberOverflows(tok.Str) {
				return l.failDirective(tok, "numeric overflow")
			}
			return tok
		}
		restore, prevType = true, tokPrev
		if l.ccStray && l.ccDepth == 0 {
			continue
		}
		before := *l
		ok := true
		switch tok.Str {
		case "$IF":
			ok = l.selectBranch(tok)
		case "$ELSIF", "$ELSE":
			// The selected branch ends; skip the rest of its $IF. No branch
			// follows an $ELSE.
			lenient := l.ccDepth == 0 && l.ccLenient
			if lenient && l.ccSplit {
				continue
			}
			if !lenient && (l.ccDepth == 0 || l.ccLevels.isElse) {
				return l.failDirective(tok, fmt.Sprintf("unexpected preprocessor token '%s'", tok.Str))
			}
			state := ccInElse
			if tok.Str == "$ELSIF" {
				state = ccInCondition
			}
			if _, ok = l.skipBranch(tok, state, false); ok && !lenient {
				l.closeLevel()
			}
		case "$END":
			if l.ccDepth == 0 {
				if l.ccLenient {
					continue
				}
				return l.failDirective(tok, "unexpected preprocessor token '$END'")
			}
			l.closeLevel()
		case "$ERROR":
			if l.ccDepth == 0 && l.ccSplit {
				continue
			}
			// An $ERROR in a selected branch fails it (PLS-00179); when the
			// selection was undecided, the parser tries the others.
			if !l.ccLenient {
				eof := l.failDirective(tok, "$ERROR directive in a selected branch")
				l.Err.(*directiveError).raised = true
				return eof
			}
			_, ok = l.skipBranch(tok, ccInError, false)
		default: // $THEN
			if l.ccLenient {
				continue
			}
			return l.failDirective(tok, "unexpected preprocessor token '$THEN'")
		}
		if !ok {
			if l.ccLenient {
				// Lex on after the malformed directive, as if it were not
				// there.
				*l = before
				continue
			}
			return Token{Type: tokEOF, Loc: l.pos, End: l.pos}
		}
		l.ccSkipEnd = l.pos
	}
}

// ccLevel is an open level of $IF: whether the branch it selected is an
// $ELSE branch, and the offset of the $IF whose decision its selection
// depends on, the latest of them, or -1 for none. Levels do not change once
// opened, so the lexer's copies share them.
type ccLevel struct {
	isElse bool
	dep    int
	outer  *ccLevel
}

// openLevel opens a level of $IF around the branch it selected, an $ELSE
// branch or not, depending on the decision at dep, or on none for -1.
func (l *Lexer) openLevel(isElse bool, dep int) bool {
	l.ccLevels = &ccLevel{isElse: isElse, dep: dep, outer: l.ccLevels}
	l.ccDepth++
	if dep >= 0 {
		l.ccUndecided++
	}
	return true
}

// closeLevel closes the innermost level of $IF.
func (l *Lexer) closeLevel() {
	if l.ccLevels.dep >= 0 {
		l.ccUndecided--
	}
	l.ccLevels = l.ccLevels.outer
	l.ccDepth--
}

// ccState is where a skipped conditional compilation construct stands.
type ccState int

const (
	ccInCondition ccState = iota // after $IF or $ELSIF, before $THEN
	ccInThen                     // in a $THEN branch
	ccInElse                     // in an $ELSE branch
	ccInError                    // after $ERROR, before its $END
)

// ccValue is what the text settles of a condition.
type ccValue int

const (
	ccUnknown ccValue = iota
	ccTrue
	ccFalse
)

// selectBranch selects the branch of the $IF at tok that Oracle compiles,
// skipping the branches before it, and opens its level; when it selects
// none, the input resumes after the $END. A branch whose condition the
// text does not settle is one Oracle may compile or pass over: the lexer
// takes the first such branch past the number ccChoices names, and records
// the decision for the parser to try the others; a condition the selection
// already decided (ccFact) is decided alike. It reports false after setting
// Err.
func (l *Lexer) selectBranch(tok Token) bool {
	at, pass := tok.Loc, l.ccChoices.passFor(tok.Loc)
	undecided, passed := false, 0
	// dep is the latest decision the selection depends on.
	dep := -1
	for {
		value, cond, ok := l.skipCondition(tok)
		if !ok {
			return false
		}
		if value == ccUnknown && cond != "" {
			if f, ok := l.fact(cond); ok {
				value = ccFalse
				if f.holds {
					value = ccTrue
				}
				dep = max(dep, f.by)
			}
		}
		switch value {
		case ccTrue:
			if undecided {
				l.noteDecision(at, passed, false)
				dep = at
			}
			return l.openLevel(false, dep)
		case ccUnknown:
			// A branch that opens with $ERROR raises whenever it is
			// selected, so no selection that compiles takes it: it is
			// passed over as if its condition were FALSE, with no
			// decision to vary. $IF DBMS_DB_VERSION.VER_LE_12 $THEN
			// $ERROR 'unsupported' $END $END costs no retry.
			if l.branchOpensWithError() {
				break
			}
			undecided = true
			if passed == pass {
				l.noteDecision(at, passed, true)
				l.noteFact(cond, true, at)
				return l.openLevel(false, at)
			}
			l.noteFact(cond, false, at)
			passed++
		}
		end, ok := l.skipBranch(tok, ccInThen, true)
		if !ok {
			return false
		}
		switch end.Str {
		case "$END":
			if undecided {
				l.noteDecision(at, passed, false)
			}
			return true
		case "$ELSE":
			if undecided {
				l.noteDecision(at, passed, false)
				dep = at
			}
			return l.openLevel(true, dep)
		}
		tok = end // $ELSIF
	}
}

// branchOpensWithError reports whether the branch after the $THEN just
// read opens with an $ERROR directive, comments aside.
func (l *Lexer) branchOpensWithError() bool {
	c := *l
	tok := c.nextRaw()
	return tok.Type == tokDIRECTIVE && tok.Str == "$ERROR"
}

// skipCondition skips the condition after the $IF or $ELSIF at tok up to
// its $THEN, recording it for the parser, and returns its value when the
// text settles it (staticCondition), and its tokens, which identify it
// among the conditions of a unit (ccFact). $$PLSQL_LINE differs from line
// to line, so its line is part of it. It reports false after setting Err.
func (l *Lexer) skipCondition(tok Token) (ccValue, string, bool) {
	cond := nodes.Loc{Start: -1, End: -1}
	tokens := 0
	var key []Token
	for {
		next := l.nextRaw()
		switch {
		case l.Err != nil:
			l.skippedError(next)
			return ccUnknown, "", false
		case next.Type == tokEOF:
			l.failDirective(tok, fmt.Sprintf("'%s' preprocessor directive does not end properly", tok.Str))
			return ccUnknown, "", false
		case next.Type == tokDIRECTIVE && next.Str == "$THEN":
			if tokens == 0 {
				// An empty condition is not a static boolean expression
				// (PLS-00174).
				l.failDirective(next, "a static boolean expression must be used")
				return ccUnknown, "", false
			}
			l.cc.noteCondition(cond)
			value, valid := staticCondition(l.input, cond)
			switch {
			case !valid && l.ccLenient:
				// Split lexes on past a condition that is no expression,
				// as past any malformed directive.
				l.failDirective(tok, "invalid preprocessor condition")
				return ccUnknown, "", false
			case !valid:
				// Oracle 23ai reports the condition (PLS-00174, PLS-00382)
				// before the branch it would select is parsed.
				at := NewLexerRange(l.input, cond.Start, cond.End).nextRaw()
				l.failDirective(at, "a static boolean expression must be used")
				return ccUnknown, "", false
			}
			return value, conditionKey(key), true
		case next.Type == tokDIRECTIVE:
			l.failDirective(next, fmt.Sprintf("unexpected preprocessor token '%s'", next.Str))
			return ccUnknown, "", false
		case next.Type == tokBIND:
			l.cc.noteBind(next.Loc)
		}
		if next.Type == tokINQUIRY && next.Str == "$$PLSQL_LINE" {
			// The line it stands on, by where that starts.
			next.Str += "@" + strconv.Itoa(strings.LastIndexByte(l.input[:next.Loc], '\n'))
		}
		key = append(key, next)
		if tokens == 0 {
			cond.Start = next.Loc
		}
		tokens++
		cond.End = next.End
	}
}

// conditionKey returns what identifies the condition of tokens among a
// unit's conditions: its tokens, without parentheses around the whole of
// it, so ($$x) is $$x.
func conditionKey(tokens []Token) string {
	for len(tokens) > 2 && tokens[0].Type == '(' && tokens[len(tokens)-1].Type == ')' {
		depth, outer := 0, true
		for i, tok := range tokens[:len(tokens)-1] {
			switch tok.Type {
			case '(':
				depth++
			case ')':
				depth--
			}
			if depth == 0 && i > 0 {
				// The first '(' closes before the end: (a) OR (b).
				outer = false
				break
			}
		}
		if !outer {
			break
		}
		tokens = tokens[1 : len(tokens)-1]
	}
	var key strings.Builder
	for _, tok := range tokens {
		fmt.Fprintf(&key, "%d:%s ", tok.Type, tok.Str)
	}
	return key.String()
}

// skipBranch skips code that is not selected, from the directive at tok in
// the given state, checking the structure of the constructs nested in it.
// It stops after the $END that closes the construct or, with stopAtBranch,
// at the $ELSIF, $ELSE, or $END that ends a $THEN branch, which it returns.
// It reports false after setting Err.
func (l *Lexer) skipBranch(tok Token, state ccState, stopAtBranch bool) (Token, bool) {
	stack := []ccState{state}
	for {
		next := l.nextRaw()
		if l.Err != nil {
			l.skippedError(next)
			return next, false
		}
		if next.Type == tokEOF {
			l.failDirective(tok, fmt.Sprintf("'%s' preprocessor directive does not end properly", tok.Str))
			return next, false
		}
		if next.Type != tokDIRECTIVE {
			if next.Type == tokBIND {
				l.cc.noteBind(next.Loc)
			}
			continue
		}
		top := &stack[len(stack)-1]
		ok := true
		switch next.Str {
		case "$THEN":
			ok = *top == ccInCondition
			*top = ccInThen
		case "$ELSIF":
			ok = *top == ccInThen
			*top = ccInCondition
		case "$ELSE":
			ok = *top == ccInThen
			*top = ccInElse
		case "$END":
			ok = *top != ccInCondition
		case "$IF":
			ok = *top == ccInThen || *top == ccInElse
		case "$ERROR":
			ok = *top == ccInThen || *top == ccInElse
		}
		if !ok {
			l.failDirective(next, fmt.Sprintf("unexpected preprocessor token '%s'", next.Str))
			return next, false
		}
		if stopAtBranch && len(stack) == 1 && (next.Str == "$ELSIF" || next.Str == "$ELSE" || next.Str == "$END") {
			return next, true
		}
		switch next.Str {
		case "$END":
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				return next, true
			}
		case "$IF":
			stack = append(stack, ccInCondition)
		case "$ERROR":
			stack = append(stack, ccInError)
		}
	}
}

// numberOverflows reports whether the numeric literal val is past the
// range of its type, as Oracle 23ai finds with PLS-00569 in PL/SQL and
// ORA-01426 in SQL. A NUMBER literal is when it comes to 1e126 or more once
// NUMBER rounds it to 40 digits: 1e126, and 9.99...95e125 with 41 digits or
// more; 1e-131 is 0. A BINARY_FLOAT (f) or BINARY_DOUBLE (d) literal is
// when IEEE rounding makes it infinite or, from a nonzero value, zero:
// 3.40282357e38f and 1e-46f, 1.7976931348623159e308d and 2e-324d. Only the
// leading digits of a NUMBER literal are read, however long.
func numberOverflows(val string) bool {
	if val == "" {
		return false
	}
	if suffix := val[len(val)-1]; strings.IndexByte("fFdD", suffix) >= 0 {
		bits := 64
		if suffix == 'f' || suffix == 'F' {
			bits = 32
		}
		f, _ := strconv.ParseFloat(val[:len(val)-1], bits)
		mant, _, _ := strings.Cut(strings.ToLower(val[:len(val)-1]), "e")
		return math.IsInf(f, 0) || f == 0 && strings.Trim(mant, "0.") != ""
	}
	mant, exp := val, 0
	if i := strings.IndexAny(val, "eE"); i >= 0 {
		mant = val[:i]
		e, err := strconv.Atoi(val[i+1:])
		if err != nil {
			// An exponent past int: past any bound either way.
			e = 1 << 40
			if strings.HasPrefix(val[i+1:], "-") {
				e = -e
			}
		}
		exp = max(min(e, 1<<40), -1<<40)
	}
	whole, frac, _ := strings.Cut(mant, ".")
	// digits run from the leading nonzero digit, at 10^top.
	digits := strings.TrimLeft(whole, "0")
	top := len(digits) - 1 + exp
	if digits == "" {
		digits = strings.TrimLeft(frac, "0")
		if digits == "" {
			return false
		}
		top = exp - (len(frac) - len(digits)) - 1
	} else {
		digits += frac
	}
	switch {
	case top > 125:
		return true
	case top < 125:
		return false
	}
	// NUMBER keeps the digits from 10^125 down to 10^86, and rounds up to
	// 10^126 only from 40 nines followed by a 5 on.
	return len(digits) > 40 && strings.TrimLeft(digits[:40], "9") == "" && digits[40] >= '5'
}

// directiveError is a lexical error at a token: a conditional compilation
// error at a preprocessor token, or a literal past its limit. raised says
// it is an $ERROR in a selected branch, which every selection of the
// branches holding it meets.
type directiveError struct {
	msg    string
	tok    Token
	raised bool
}

func (e *directiveError) Error() string { return e.msg }

// skippedError keeps a lexical error met in code a directive skips or a
// condition, at tok, where it is: it is the first error (a literal past
// its limit, an unterminated string), not the directive left open.
func (l *Lexer) skippedError(tok Token) {
	var directive *directiveError
	if !errors.As(l.Err, &directive) {
		l.Err = &directiveError{msg: l.Err.Error(), tok: tok}
	}
	l.pos = l.end
}

// failDirective sets Err for a conditional compilation error at tok and
// ends the input, returning the EOF token at tok.
func (l *Lexer) failDirective(tok Token, msg string) Token {
	l.Err = &directiveError{msg: msg, tok: tok}
	l.pos = l.end
	return Token{Type: tokEOF, Loc: tok.Loc, End: tok.Loc}
}

// nextRaw returns the next token before conditional compilation: directives
// are tokDIRECTIVE tokens.
func (l *Lexer) nextRaw() Token {
	if l.pos >= l.end {
		return Token{Type: tokEOF, Loc: l.pos, End: l.pos}
	}

	l.start = l.pos
	tok := l.lexInitial()
	// Set End to current position after token has been fully consumed.
	// For recursive calls (e.g., skipping comments), the inner NextToken
	// sets End correctly on its own returned token.
	tok.End = l.pos
	if (tok.Type == tokSCONST || tok.Type == tokNCHARLIT) && l.literalLimit > 0 && len(tok.Str) > l.literalLimit {
		return l.failDirective(tok, "string literal too long")
	}
	// A hint stands right after SELECT, INSERT, UPDATE, DELETE, or MERGE;
	// anywhere else it is a comment (Oracle 23ai runs BEGIN /*+ x */ NULL;
	// END; and SELECT 1 /*+ x */ FROM dual).
	if tok.Type == tokHINT {
		switch l.prevType {
		case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwMERGE:
		default:
			l.clearPrevType()
			return l.nextRaw()
		}
	}
	// A directive, with the condition or message after it, leaves the
	// token before it last, as the preprocessed text does: SELECT $IF TRUE
	// $THEN /*+ FULL(t) */ and SELECT $IF TRUE $THEN $END /*+ FULL(t) */
	// keep their hints.
	if tok.Type == tokDIRECTIVE {
		switch tok.Str {
		case "$IF", "$ELSIF", "$ERROR":
			l.ccRawDirective = true
		case "$THEN", "$END":
			l.ccRawDirective = false
		}
		return tok
	}
	if !l.ccRawDirective {
		l.prevType = tok.Type
	}
	return tok
}

// clearPrevType records that a comment came after the last token, so a
// hint after it is a comment too. In a directive's condition or message,
// which the preprocessed text leaves out, it does not: SELECT $IF /*c*/
// TRUE $THEN $END /*+ FULL(t) */ keeps its hint.
func (l *Lexer) clearPrevType() {
	if !l.ccRawDirective {
		l.prevType = 0
	}
}

// lexInitial handles tokens in the initial state.
func (l *Lexer) lexInitial() Token {
	l.skipWhitespace()
	if l.pos >= l.end {
		return Token{Type: tokEOF, Loc: l.pos}
	}

	l.start = l.pos
	ch := l.input[l.pos]

	// Line comments: --
	if ch == '-' && l.pos+1 < l.end && l.input[l.pos+1] == '-' {
		l.skipLineComment()
		l.clearPrevType() // a hint follows its keyword directly
		return l.nextRaw()
	}

	// Block comments: /* */ and hints: /*+ */
	if ch == '/' && l.pos+1 < l.end && l.input[l.pos+1] == '*' {
		return l.lexBlockCommentOrHint()
	}

	// National character literal: N'...'
	if (ch == 'N' || ch == 'n') && l.pos+1 < l.end && l.input[l.pos+1] == '\'' {
		l.pos++ // skip N
		str := l.lexSingleQuotedString()
		return Token{Type: tokNCHARLIT, Str: str, Loc: l.start}
	}

	// Q-quote mechanism: q'[delim]...[delim]'
	if (ch == 'Q' || ch == 'q') && l.pos+1 < l.end && l.input[l.pos+1] == '\'' {
		return l.lexQQuote()
	}

	// Single-quoted string: '...'
	if ch == '\'' {
		str := l.lexSingleQuotedString()
		return Token{Type: tokSCONST, Str: str, Loc: l.start}
	}

	// Double-quoted identifier: "..."
	if ch == '"' {
		return l.lexDoubleQuotedIdent()
	}

	// Bind variable: :name or :1
	if ch == ':' {
		if l.pos+1 < l.end {
			next := l.input[l.pos+1]
			// := assignment
			if next == '=' {
				l.pos += 2
				return Token{Type: tokASSIGN, Str: ":=", Loc: l.start}
			}
			// Bind variable
			if isIdentStart(next) || isDigit(next) {
				l.pos++ // skip :
				return l.lexBindVariable()
			}
		}
		// Single : is returned as itself
		l.pos++
		return Token{Type: int(ch), Str: string(ch), Loc: l.start}
	}

	// Two-character tokens
	if l.pos+1 < l.end {
		ch2 := l.input[l.pos : l.pos+2]
		switch ch2 {
		case "||":
			l.pos += 2
			return Token{Type: tokCONCAT, Str: "||", Loc: l.start}
		case "=>":
			l.pos += 2
			return Token{Type: tokASSOC, Str: "=>", Loc: l.start}
		case "..":
			l.pos += 2
			return Token{Type: tokDOTDOT, Str: "..", Loc: l.start}
		case "**":
			l.pos += 2
			return Token{Type: tokEXPON, Str: "**", Loc: l.start}
		case "<=":
			l.pos += 2
			return Token{Type: tokLESSEQ, Str: "<=", Loc: l.start}
		case ">=":
			l.pos += 2
			return Token{Type: tokGREATEQ, Str: ">=", Loc: l.start}
		case "!=":
			l.pos += 2
			return Token{Type: tokNOTEQ, Str: "!=", Loc: l.start}
		case "<>":
			l.pos += 2
			return Token{Type: tokNOTEQ, Str: "<>", Loc: l.start}
		case "~=":
			l.pos += 2
			return Token{Type: tokNOTEQ, Str: "~=", Loc: l.start}
		case "^=":
			l.pos += 2
			return Token{Type: tokNOTEQ, Str: "^=", Loc: l.start}
		case "<<":
			l.pos += 2
			return Token{Type: tokLABELOPEN, Str: "<<", Loc: l.start}
		case ">>":
			l.pos += 2
			return Token{Type: tokLABELCLOSE, Str: ">>", Loc: l.start}
		}
	}

	// Numbers
	if isDigit(ch) || (ch == '.' && l.pos+1 < l.end && isDigit(l.input[l.pos+1])) {
		return l.lexNumber()
	}

	// Identifiers and keywords
	if isIdentStart(ch) {
		return l.lexIdentOrKeyword()
	}

	// Conditional compilation: a selection or error directive, or an
	// inquiry directive $$name.
	if ch == '$' {
		if tok, ok := l.lexPreprocessorToken(); ok {
			return tok
		}
	}

	// Self-delimiting single characters
	l.pos++
	return Token{Type: int(ch), Str: string(ch), Loc: l.start}
}

// skipWhitespace skips whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.pos < l.end {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '\f' {
			l.pos++
		} else {
			break
		}
	}
}

// skipLineComment skips a -- comment to end of line.
func (l *Lexer) skipLineComment() {
	l.pos += 2 // skip --
	for l.pos < l.end {
		ch := l.input[l.pos]
		l.pos++
		if ch == '\n' {
			break
		}
	}
}

// lexBlockCommentOrHint handles /* ... */ comments and /*+ hints */.
func (l *Lexer) lexBlockCommentOrHint() Token {
	l.pos += 2 // skip /*

	// Check if it's a hint: /*+ ... */
	isHint := l.pos < l.end && l.input[l.pos] == '+'
	if isHint {
		l.pos++ // skip +
	}

	// Oracle block comments do not nest: the first */ terminates the
	// comment regardless of any /* sequences inside it. Engine-verified in
	// both directions on 11gR2: SELECT 1 /* a /* b */ FROM dual succeeds,
	// while appending a second */ is ORA-00936 — the shape that would be
	// legal under nesting semantics. (PostgreSQL nests; this is a genuine
	// per-engine divergence, so the PG lexer must keep its counting.)
	var buf strings.Builder
	terminated := false

	for l.pos < l.end {
		if l.input[l.pos] == '*' && l.pos+1 < l.end && l.input[l.pos+1] == '/' {
			l.pos += 2
			terminated = true
			break
		}
		buf.WriteByte(l.input[l.pos])
		l.pos++
	}

	if !terminated {
		l.Err = fmt.Errorf("unterminated block comment")
		return Token{Type: tokEOF, Loc: l.start}
	}

	if isHint {
		return Token{Type: tokHINT, Str: strings.TrimSpace(buf.String()), Loc: l.start}
	}

	// Regular comment -- skip and get next token. A hint follows its
	// keyword directly: SELECT /* c */ /*+ FULL(t) */ holds no hint.
	l.clearPrevType()
	return l.nextRaw()
}

// lexSingleQuotedString handles '...' strings with ” escape.
func (l *Lexer) lexSingleQuotedString() string {
	l.pos++ // skip opening '
	var buf strings.Builder

	for l.pos < l.end {
		ch := l.input[l.pos]
		if ch == '\'' {
			// Check for doubled quote ''
			if l.pos+1 < l.end && l.input[l.pos+1] == '\'' {
				buf.WriteByte('\'')
				l.pos += 2
				continue
			}
			// End of string
			l.pos++
			return buf.String()
		}
		buf.WriteByte(ch)
		l.pos++
	}

	l.Err = fmt.Errorf("unterminated string literal")
	return buf.String()
}

// lexQQuote handles Oracle's alternative quoting mechanism: q'[delim]...[delim]'
func (l *Lexer) lexQQuote() Token {
	l.pos += 2 // skip q'

	if l.pos >= l.end {
		l.Err = fmt.Errorf("unterminated q-quote string")
		return Token{Type: tokEOF, Loc: l.start}
	}

	openDelim := l.input[l.pos]
	l.pos++

	// Determine close delimiter
	var closeDelim byte
	switch openDelim {
	case '[':
		closeDelim = ']'
	case '{':
		closeDelim = '}'
	case '(':
		closeDelim = ')'
	case '<':
		closeDelim = '>'
	default:
		closeDelim = openDelim
	}

	var buf strings.Builder
	for l.pos < l.end {
		ch := l.input[l.pos]
		if ch == closeDelim && l.pos+1 < l.end && l.input[l.pos+1] == '\'' {
			l.pos += 2 // skip close delim and '
			return Token{Type: tokSCONST, Str: buf.String(), Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}

	l.Err = fmt.Errorf("unterminated q-quote string")
	return Token{Type: tokEOF, Loc: l.start}
}

// lexDoubleQuotedIdent handles "..." identifiers.
func (l *Lexer) lexDoubleQuotedIdent() Token {
	l.pos++ // skip opening "
	var buf strings.Builder

	for l.pos < l.end {
		ch := l.input[l.pos]
		if ch == '"' {
			// Check for "" (escaped quote)
			if l.pos+1 < l.end && l.input[l.pos+1] == '"' {
				buf.WriteByte('"')
				l.pos += 2
				continue
			}
			l.pos++ // skip closing "
			str := buf.String()
			if str == "" {
				l.Err = fmt.Errorf("zero-length delimited identifier")
				return Token{Type: tokEOF, Loc: l.start}
			}
			return Token{Type: tokQIDENT, Str: str, Loc: l.start}
		}
		buf.WriteByte(ch)
		l.pos++
	}

	l.Err = fmt.Errorf("unterminated quoted identifier")
	return Token{Type: tokEOF, Loc: l.start}
}

// lexBindVariable handles :name or :1 bind variables.
func (l *Lexer) lexBindVariable() Token {
	start := l.pos

	if isDigit(l.input[l.pos]) {
		// Numeric bind: :1, :2, etc.
		for l.pos < l.end && isDigit(l.input[l.pos]) {
			l.pos++
		}
	} else {
		// Named bind: :name
		for l.pos < l.end && isIdentCont(l.input[l.pos]) {
			l.pos++
		}
	}

	name := l.input[start:l.pos]
	return Token{Type: tokBIND, Str: name, Loc: l.start}
}

// lexNumber handles integer and floating-point literals.
func (l *Lexer) lexNumber() Token {
	start := l.pos
	isFloat := false

	// Integer part (may be empty for .5)
	for l.pos < l.end && isDigit(l.input[l.pos]) {
		l.pos++
	}

	// Decimal point
	if l.pos < l.end && l.input[l.pos] == '.' {
		// Check for .. (range operator)
		if l.pos+1 < l.end && l.input[l.pos+1] == '.' {
			goto done
		}
		l.pos++
		isFloat = true
		for l.pos < l.end && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}

	// Exponent
	if l.pos < l.end && (l.input[l.pos] == 'e' || l.input[l.pos] == 'E') {
		l.pos++
		if l.pos < l.end && (l.input[l.pos] == '+' || l.input[l.pos] == '-') {
			l.pos++
		}
		if l.pos >= l.end || !isDigit(l.input[l.pos]) {
			l.Err = fmt.Errorf("invalid numeric literal")
			return Token{Type: tokEOF, Loc: l.start}
		}
		for l.pos < l.end && isDigit(l.input[l.pos]) {
			l.pos++
		}
		isFloat = true
	}

	// Oracle allows d/D suffix for double, f/F for float (in some contexts)
	if l.pos < l.end {
		ch := l.input[l.pos]
		if ch == 'f' || ch == 'F' || ch == 'd' || ch == 'D' {
			l.pos++
			isFloat = true
		}
	}

done:
	numStr := l.input[start:l.pos]

	if isFloat {
		return Token{Type: tokFCONST, Str: numStr, Loc: l.start}
	}

	val, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return Token{Type: tokFCONST, Str: numStr, Loc: l.start}
	}

	return Token{Type: tokICONST, Ival: val, Str: numStr, Loc: l.start}
}

// lexIdentOrKeyword handles identifiers and keywords.
func (l *Lexer) lexIdentOrKeyword() Token {
	start := l.pos

	for l.pos < l.end && isIdentCont(l.input[l.pos]) {
		l.pos++
	}

	// Check for @dblink suffix on identifiers is handled at parser level
	ident := l.input[start:l.pos]

	// Check if it's a keyword (case-insensitive)
	upper := strings.ToUpper(ident)
	if kw, ok := oracleKeywords[upper]; ok {
		return Token{Type: kw, Str: upper, Loc: l.start}
	}

	// Return as identifier (Oracle unquoted identifiers are case-insensitive;
	// we store them uppercased per Oracle convention)
	return Token{Type: tokIDENT, Str: upper, Loc: l.start}
}

// lexPreprocessorToken lexes the PL/SQL preprocessor token at '$': a
// directive, as a tokDIRECTIVE token, or an inquiry directive $$name, as a
// tokINQUIRY token. Any other '$' is left to lex as itself.
func (l *Lexer) lexPreprocessorToken() (Token, bool) {
	start := l.pos
	if l.pos+2 < l.end && l.input[l.pos+1] == '$' && isIdentStart(l.input[l.pos+2]) {
		l.pos += 2
		for l.pos < l.end && isIdentCont(l.input[l.pos]) {
			l.pos++
		}
		l.cc.noteDirective(start)
		return Token{Type: tokINQUIRY, Str: "$$" + strings.ToUpper(l.input[start+2:l.pos]), Loc: start}, true
	}
	end := start + 1
	for end < l.end && isIdentCont(l.input[end]) {
		end++
	}
	switch word := strings.ToUpper(l.input[start+1 : end]); word {
	case "IF", "THEN", "ELSIF", "ELSE", "END", "ERROR":
		l.pos = end
		l.cc.noteDirective(start)
		return Token{Type: tokDIRECTIVE, Str: "$" + word, Loc: start}, true
	}
	return Token{}, false
}

// Character classification functions.

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isIdentStart(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_' || ch >= 128
}

func isIdentCont(ch byte) bool {
	return isIdentStart(ch) || isDigit(ch) || ch == '$' || ch == '#'
}

// oracleKeywords maps uppercase keyword strings to token types.
// Oracle keywords. This is a subset; more can be added as batches are implemented.
var oracleKeywords = map[string]int{
	// Reserved words
	"ACCESS":          kwACCESS,
	"ADD":             kwADD,
	"ALL":             kwALL,
	"ALTER":           kwALTER,
	"AND":             kwAND,
	"ANY":             kwANY,
	"AS":              kwAS,
	"ASC":             kwASC,
	"AUDIT":           kwAUDIT,
	"BETWEEN":         kwBETWEEN,
	"BY":              kwBY,
	"CASE":            kwCASE,
	"CAST":            kwCAST,
	"CHECK":           kwCHECK,
	"CLUSTER":         kwCLUSTER,
	"COLUMN":          kwCOLUMN,
	"COLUMN_VALUE":    kwCOLUMN_VALUE,
	"COMMENT":         kwCOMMENT,
	"COMMIT":          kwCOMMIT,
	"COMPRESS":        kwCOMPRESS,
	"CONNECT":         kwCONNECT,
	"CONSTRAINT":      kwCONSTRAINT,
	"CREATE":          kwCREATE,
	"CROSS":           kwCROSS,
	"CURRENT":         kwCURRENT,
	"DATE":            kwDATE,
	"DECIMAL":         kwDECIMAL,
	"DECODE":          kwDECODE,
	"DEFAULT":         kwDEFAULT,
	"DELETE":          kwDELETE,
	"DESC":            kwDESC,
	"DISTINCT":        kwDISTINCT,
	"DROP":            kwDROP,
	"ELSE":            kwELSE,
	"ELSIF":           kwELSIF,
	"END":             kwEND,
	"EXCEPTION":       kwEXCEPTION,
	"EXCLUSIVE":       kwEXCLUSIVE,
	"EXECUTE":         kwEXECUTE,
	"EXISTS":          kwEXISTS,
	"FETCH":           kwFETCH,
	"FILE":            kwFILE,
	"FIRST":           kwFIRST,
	"FLOAT":           kwFLOAT,
	"FOR":             kwFOR,
	"FORCE":           kwFORCE,
	"FROM":            kwFROM,
	"FULL":            kwFULL,
	"FUNCTION":        kwFUNCTION,
	"GOTO":            kwGOTO,
	"GRANT":           kwGRANT,
	"GROUP":           kwGROUP,
	"HAVING":          kwHAVING,
	"IDENTIFIED":      kwIDENTIFIED,
	"IF":              kwIF,
	"IMMEDIATE":       kwIMMEDIATE,
	"IN":              kwIN,
	"INCREMENT":       kwINCREMENT,
	"INDEX":           kwINDEX,
	"INITIAL":         kwINITIAL,
	"INNER":           kwINNER,
	"INSERT":          kwINSERT,
	"INTEGER":         kwINTEGER,
	"INTERSECT":       kwINTERSECT,
	"INTO":            kwINTO,
	"IS":              kwIS,
	"JOIN":            kwJOIN,
	"KEEP":            kwKEEP,
	"LAST":            kwLAST,
	"LEFT":            kwLEFT,
	"LEVEL":           kwLEVEL,
	"LIKE":            kwLIKE,
	"LIKE2":           kwLIKE2,
	"LIKE4":           kwLIKE4,
	"LIKEC":           kwLIKEC,
	"LOCK":            kwLOCK,
	"LONG":            kwLONG,
	"LOOP":            kwLOOP,
	"MATERIALIZED":    kwMATERIALIZED,
	"MAXEXTENTS":      kwMAXEXTENTS,
	"MERGE":           kwMERGE,
	"MINUS":           kwMINUS,
	"MLSLABEL":        kwMLSLABEL,
	"MODE":            kwMODE,
	"MODIFY":          kwMODIFY,
	"NATURAL":         kwNATURAL,
	"NESTED_TABLE_ID": kwNESTED_TABLE_ID,
	"NEXT":            kwNEXT,
	"NOAUDIT":         kwNOAUDIT,
	"NOCACHE":         kwNOCACHE,
	"NOCOMPRESS":      kwNOCOMPRESS,
	"NOCYCLE":         kwNOCYCLE,
	"NOMAXVALUE":      kwNOMAXVALUE,
	"NOMINVALUE":      kwNOMINVALUE,
	"NOORDER":         kwNOORDER,
	"NOT":             kwNOT,
	"NOWAIT":          kwNOWAIT,
	"NULL":            kwNULL,
	"NUMBER":          kwNUMBER,
	"OF":              kwOF,
	"OFFLINE":         kwOFFLINE,
	"OFFSET":          kwOFFSET,
	"ON":              kwON,
	"ONLINE":          kwONLINE,
	"ONLY":            kwONLY,
	"OPEN":            kwOPEN,
	"OPTION":          kwOPTION,
	"OR":              kwOR,
	"ORDER":           kwORDER,
	"OUTER":           kwOUTER,
	"OVER":            kwOVER,
	"PACKAGE":         kwPACKAGE,
	"PARALLEL":        kwPARALLEL,
	"PARTITION":       kwPARTITION,
	"PCTFREE":         kwPCTFREE,
	"PERCENT":         kwPERCENT,
	"PIVOT":           kwPIVOT,
	"PRIOR":           kwPRIOR,
	"PRIVILEGES":      kwPRIVILEGES,
	"PROCEDURE":       kwPROCEDURE,
	"PUBLIC":          kwPUBLIC,
	"PURGE":           kwPURGE,
	"RANGE":           kwRANGE,
	"RAW":             kwRAW,
	"RENAME":          kwRENAME,
	"REPLACE":         kwREPLACE,
	"RETURN":          kwRETURN,
	"RETURNING":       kwRETURNING,
	"REVERSE":         kwREVERSE,
	"REVOKE":          kwREVOKE,
	"RIGHT":           kwRIGHT,
	"ROLLBACK":        kwROLLBACK,
	"ROW":             kwROW,
	"ROWID":           kwROWID,
	"ROWNUM":          kwROWNUM,
	"ROWS":            kwROWS,
	"SAMPLE":          kwSAMPLE,
	"SAVEPOINT":       kwSAVEPOINT,
	"SELECT":          kwSELECT,
	"SEQUENCE":        kwSEQUENCE,
	"SESSION":         kwSESSION,
	"SET":             kwSET,
	"SIZE":            kwSIZE,
	"SMALLINT":        kwSMALLINT,
	"START":           kwSTART,
	"SUBPARTITION":    kwSUBPARTITION,
	"SYNONYM":         kwSYNONYM,
	"SYSDATE":         kwSYSDATE,
	"SYSTIMESTAMP":    kwSYSTIMESTAMP,
	"SYSTEM":          kwSYSTEM,
	"TABLE":           kwTABLE,
	"THEN":            kwTHEN,
	"TIES":            kwTIES,
	"TO":              kwTO,
	"TRIGGER":         kwTRIGGER,
	"TRUNCATE":        kwTRUNCATE,
	"TYPE":            kwTYPE,
	"UID":             kwUID,
	"UNION":           kwUNION,
	"UNIQUE":          kwUNIQUE,
	"UNPIVOT":         kwUNPIVOT,
	"UPDATE":          kwUPDATE,
	"USING":           kwUSING,
	"VALUES":          kwVALUES,
	"VARCHAR":         kwVARCHAR,
	"VARCHAR2":        kwVARCHAR2,
	"VIEW":            kwVIEW,
	"WAIT":            kwWAIT,
	"WHEN":            kwWHEN,
	"WHERE":           kwWHERE,
	"WHILE":           kwWHILE,
	"WITH":            kwWITH,
	"WORK":            kwWORK,

	// Additional Oracle keywords for various features
	"ANALYZE":       kwANALYZE,
	"BEGIN":         kwBEGIN,
	"BITMAP":        kwBITMAP,
	"BLOB":          kwBLOB,
	"BODY":          kwBODY,
	"BULK":          kwBULK,
	"CACHE":         kwCACHE,
	"CASCADE":       kwCASCADE,
	"CHAR":          kwCHAR,
	"CLOB":          kwCLOB,
	"CLOSE":         kwCLOSE,
	"COLLECT":       kwCOLLECT,
	"COMPOUND":      kwCOMPOUND,
	"CONSTRAINTS":   kwCONSTRAINTS,
	"CONTINUE":      kwCONTINUE,
	"CURSOR":        kwCURSOR,
	"CYCLE":         kwCYCLE,
	"DATABASE":      kwDATABASE,
	"DECLARE":       kwDECLARE,
	"DEFERRABLE":    kwDEFERRABLE,
	"DEFERRED":      kwDEFERRED,
	"DETERMINISTIC": kwDETERMINISTIC,
	"EACH":          kwEACH,
	"ENABLE":        kwENABLE,
	"ESCAPE":        kwESCAPE,
	"EXPLAIN":       kwEXPLAIN,
	"FLASHBACK":     kwFLASHBACK,
	"FOREIGN":       kwFOREIGN,
	"GENERATED":     kwGENERATED,
	"GLOBAL":        kwGLOBAL,
	"HASH":          kwHASH,
	"IDENTITY":      kwIDENTITY,
	"INCLUDE":       kwINCLUDE,
	"INITIALLY":     kwINITIALLY,
	"INSTEAD":       kwINSTEAD,
	"INTERVAL":      kwINTERVAL,
	"INVISIBLE":     kwINVISIBLE,
	"ISOLATION":     kwISOLATION,
	"KEY":           kwKEY,
	"LIMIT":         kwLIMIT,
	"LINK":          kwLINK,
	"LIST":          kwLIST,
	"LOCAL":         kwLOCAL,
	"LOG":           kwLOG,
	"LOGGING":       kwLOGGING,
	"MATCHED":       kwMATCHED,
	"MAXVALUE":      kwMAXVALUE,
	"MINVALUE":      kwMINVALUE,
	"MODEL":         kwMODEL,
	"MULTISET":      kwMULTISET,
	"NCHAR":         kwNCHAR,
	"NCLOB":         kwNCLOB,
	"NOLOGGING":     kwNOLOGGING,
	"NOPARALLEL":    kwNOPARALLEL,
	"NVARCHAR2":     kwNVARCHAR2,
	"NULLS":         kwNULLS,
	"PIPELINED":     kwPIPELINED,
	"PLAN":          kwPLAN,
	"PRIMARY":       kwPRIMARY,
	"PRIVATE":       kwPRIVATE,
	"RAISE":         kwRAISE,
	"READ":          kwREAD,
	"RECORD":        kwRECORD,
	"RECURSIVE":     kwRECURSIVE,
	"REF":           kwREF,
	"REFERENCES":    kwREFERENCES,
	"REFRESH":       kwREFRESH,
	"REJECT":        kwREJECT,
	"RESTRICT":      kwRESTRICT,
	"RESULT_CACHE":  kwRESULT_CACHE,
	"RESOURCE":      kwRESOURCE,
	"ROLE":          kwROLE,
	"ROWTYPE":       kwROWTYPE,
	"SKIP":          kwSKIP,
	"STORAGE":       kwSTORAGE,
	"SUBTYPE":       kwSUBTYPE,
	"TABLESPACE":    kwTABLESPACE,
	"TEMPORARY":     kwTEMPORARY,
	"TIMESTAMP":     kwTIMESTAMP,
	"TRANSACTION":   kwTRANSACTION,
	"USER":          kwUSER,
	"VALIDATE":      kwVALIDATE,
	"VARRAY":        kwVARRAY,
	"VIRTUAL":       kwVIRTUAL,
	"WRITE":         kwWRITE,
	"ZONE":          kwZONE,

	// Additional Oracle-specific
	"AFTER":               kwAFTER,
	"ALWAYS":              kwALWAYS,
	"BEFORE":              kwBEFORE,
	"CONNECT_BY_ROOT":     kwCONNECT_BY_ROOT,
	"CONSTANT":            kwCONSTANT,
	"DENSE_RANK":          kwDENSE_RANK,
	"DISABLE":             kwDISABLE,
	"ERRORS":              kwERRORS,
	"EXIT":                kwEXIT,
	"FORALL":              kwFORALL,
	"LOCKED":              kwLOCKED,
	"NOCOPY":              kwNOCOPY,
	"OUT":                 kwOUT,
	"PARALLEL_ENABLE":     kwPARALLEL_ENABLE,
	"PRAGMA":              kwPRAGMA,
	"RELY":                kwRELY,
	"REWRITE":             kwREWRITE,
	"SCN":                 kwSCN,
	"SEED":                kwSEED,
	"SYS_CONNECT_BY_PATH": kwSYS_CONNECT_BY_PATH,
	"UNDER":               kwUNDER,

	// SAMPLE / FLASHBACK query keywords
	"BLOCK":    kwBLOCK,
	"VERSIONS": kwVERSIONS,

	// MODEL clause keywords
	"AUTOMATIC":  kwAUTOMATIC,
	"DECREMENT":  kwDECREMENT,
	"DIMENSION":  kwDIMENSION,
	"ITERATE":    kwITERATE,
	"MAIN":       kwMAIN,
	"MEASURES":   kwMEASURES,
	"NAV":        kwNAV,
	"REFERENCE":  kwREFERENCE,
	"RULES":      kwRULES,
	"SEQUENTIAL": kwSEQUENTIAL,
	"UNTIL":      kwUNTIL,
	"UPDATED":    kwUPDATED,
	"UPSERT":     kwUPSERT,

	// Analytic function keywords
	"FOLLOWING": kwFOLLOWING,
	"GROUPS":    kwGROUPS,
	"PRECEDING": kwPRECEDING,
	"UNBOUNDED": kwUNBOUNDED,
	"WITHIN":    kwWITHIN,

	// Grouping extension keywords
	"CUBE":     kwCUBE,
	"GROUPING": kwGROUPING,
	"ROLLUP":   kwROLLUP,
	"SETS":     kwSETS,

	// LATERAL / XMLTABLE / JSON_TABLE keywords
	"COLUMNS":    kwCOLUMNS,
	"JSON_TABLE": kwJSON_TABLE,
	"LATERAL":    kwLATERAL,
	"NESTED":     kwNESTED,
	"ORDINALITY": kwORDINALITY,
	"PASSING":    kwPASSING,
	"PATH":       kwPATH,
	"XMLTABLE":   kwXMLTABLE,

	// Cursor / MULTISET / TREAT keywords
	"EXCEPT": kwEXCEPT,
	"TREAT":  kwTREAT,

	// XML expression keywords
	"CONTENT":      kwCONTENT,
	"NAME":         kwNAME,
	"XMLAGG":       kwXMLAGG,
	"XMLELEMENT":   kwXMLELEMENT,
	"XMLFOREST":    kwXMLFOREST,
	"XMLPARSE":     kwXMLPARSE,
	"XMLROOT":      kwXMLROOT,
	"XMLSERIALIZE": kwXMLSERIALIZE,

	// Utility statement keywords
	"CALL":  kwCALL,
	"SHARE": kwSHARE,

	// Session/audit keywords
	"ASSOCIATE":    kwASSOCIATE,
	"DISASSOCIATE": kwDISASSOCIATE,
	"STATISTICS":   kwSTATISTICS,
	"SUCCESSFUL":   kwSUCCESSFUL,
	"WHENEVER":     kwWHENEVER,

	// Admin DDL keywords
	"CONTEXT":   kwCONTEXT,
	"DIRECTORY": kwDIRECTORY,
	"JAVA":      kwJAVA,
	"LIBRARY":   kwLIBRARY,
	"PROFILE":   kwPROFILE,
	"SCHEMA":    kwSCHEMA,

	// PL/SQL additional keywords
	"PIPE": kwPIPE,

	// JSON expression keywords
	"FORMAT":          kwFORMAT,
	"JSON":            kwJSON,
	"JSON_ARRAY":      kwJSON_ARRAY,
	"JSON_EXISTS":     kwJSON_EXISTS,
	"JSON_MERGEPATCH": kwJSON_MERGEPATCH,
	"JSON_OBJECT":     kwJSON_OBJECT,
	"JSON_QUERY":      kwJSON_QUERY,
	"JSON_VALUE":      kwJSON_VALUE,
}

// Keyword token constants.
const (
	kwACCESS = iota + 2000
	kwADD
	kwAFTER
	kwALL
	kwALTER
	kwALWAYS
	kwANALYZE
	kwAND
	kwANY
	kwAS
	kwASC
	kwAUDIT
	kwBEFORE
	kwBEGIN
	kwBETWEEN
	kwBITMAP
	kwBLOB
	kwBODY
	kwBULK
	kwBY
	kwCACHE
	kwCASCADE
	kwCASE
	kwCAST
	kwCHAR
	kwCHECK
	kwCLOB
	kwCLOSE
	kwCLUSTER
	kwCOLLECT
	kwCOLUMN
	kwCOMMENT
	kwCOMMIT
	kwCOMPOUND
	kwCOMPRESS
	kwCONNECT
	kwCONNECT_BY_ROOT
	kwCONSTANT
	kwCONSTRAINT
	kwCONSTRAINTS
	kwCONTINUE
	kwCREATE
	kwCROSS
	kwCURRENT
	kwCURSOR
	kwCYCLE
	kwDATABASE
	kwDATE
	kwDECIMAL
	kwDECLARE
	kwDECODE
	kwDEFAULT
	kwDEFERRABLE
	kwDEFERRED
	kwDELETE
	kwDENSE_RANK
	kwDESC
	kwDETERMINISTIC
	kwDISABLE
	kwDISTINCT
	kwDROP
	kwEACH
	kwELSE
	kwELSIF
	kwENABLE
	kwEND
	kwERRORS
	kwESCAPE
	kwEXCEPTION
	kwEXCLUSIVE
	kwEXECUTE
	kwEXISTS
	kwEXIT
	kwEXPLAIN
	kwFETCH
	kwFILE
	kwFIRST
	kwFLASHBACK
	kwFLOAT
	kwFOR
	kwFORALL
	kwFORCE
	kwFOREIGN
	kwFROM
	kwFULL
	kwFUNCTION
	kwGENERATED
	kwGLOBAL
	kwGOTO
	kwGRANT
	kwGROUP
	kwHASH
	kwHAVING
	kwIDENTIFIED
	kwIDENTITY
	kwIF
	kwIMMEDIATE
	kwIN
	kwINCLUDE
	kwINCREMENT
	kwINDEX
	kwINITIAL
	kwINITIALLY
	kwINNER
	kwINSERT
	kwINSTEAD
	kwINTEGER
	kwINTERSECT
	kwINTERVAL
	kwINTO
	kwINVISIBLE
	kwIS
	kwISOLATION
	kwJOIN
	kwKEEP
	kwKEY
	kwLAST
	kwLEFT
	kwLEVEL
	kwLIKE
	kwLIKE2
	kwLIKE4
	kwLIKEC
	kwLIMIT
	kwLINK
	kwLIST
	kwLOCAL
	kwLOCK
	kwLOCKED
	kwLOG
	kwLOGGING
	kwLONG
	kwLOOP
	kwMATCHED
	kwMATERIALIZED
	kwMAXEXTENTS
	kwMAXVALUE
	kwMERGE
	kwMINUS
	kwMINVALUE
	kwMODE
	kwMODEL
	kwMODIFY
	kwMULTISET
	kwNATURAL
	kwNCHAR
	kwNCLOB
	kwNEXT
	kwNOAUDIT
	kwNOCACHE
	kwNOCOMPRESS
	kwNOCOPY
	kwNOCYCLE
	kwNOLOGGING
	kwNOMAXVALUE
	kwNOMINVALUE
	kwNOORDER
	kwNOPARALLEL
	kwNOT
	kwNOWAIT
	kwNULL
	kwNULLS
	kwNUMBER
	kwNVARCHAR2
	kwOF
	kwOFFLINE
	kwOFFSET
	kwON
	kwONLINE
	kwONLY
	kwOPEN
	kwOPTION
	kwOR
	kwORDER
	kwOUT
	kwOUTER
	kwOVER
	kwPACKAGE
	kwPARALLEL
	kwPARALLEL_ENABLE
	kwPARTITION
	kwPCTFREE
	kwPERCENT
	kwPIPELINED
	kwPIVOT
	kwPLAN
	kwPRAGMA
	kwPRIMARY
	kwPRIOR
	kwPRIVATE
	kwPRIVILEGES
	kwPROCEDURE
	kwPUBLIC
	kwPURGE
	kwRAISE
	kwRANGE
	kwRAW
	kwREAD
	kwRECORD
	kwRECURSIVE
	kwREF
	kwREFERENCES
	kwREFRESH
	kwREJECT
	kwRELY
	kwRENAME
	kwREPLACE
	kwRESTRICT
	kwRESULT_CACHE
	kwRETURN
	kwRETURNING
	kwREVERSE
	kwREVOKE
	kwREWRITE
	kwRIGHT
	kwROLE
	kwROLLBACK
	kwROW
	kwROWID
	kwROWNUM
	kwROWS
	kwROWTYPE
	kwSAMPLE
	kwSAVEPOINT
	kwSCN
	kwSEED
	kwSELECT
	kwSEQUENCE
	kwSESSION
	kwSET
	kwSIZE
	kwSKIP
	kwSMALLINT
	kwSTART
	kwSTORAGE
	kwSUBPARTITION
	kwSUBTYPE
	kwSYNONYM
	kwSYS_CONNECT_BY_PATH
	kwSYSDATE
	kwSYSTEM
	kwSYSTIMESTAMP
	kwTABLE
	kwTABLESPACE
	kwTEMPORARY
	kwTHEN
	kwTIES
	kwTIMESTAMP
	kwTO
	kwTRANSACTION
	kwTRIGGER
	kwTRUNCATE
	kwTYPE
	kwUNDER
	kwUNION
	kwUNIQUE
	kwUNPIVOT
	kwUPDATE
	kwUSER
	kwUSING
	kwVALIDATE
	kwVALUES
	kwVARCHAR
	kwVARCHAR2
	kwVARRAY
	kwVIEW
	kwVIRTUAL
	kwWAIT
	kwWHEN
	kwWHERE
	kwWHILE
	kwWITH
	kwWORK
	kwWRITE
	kwZONE

	// SAMPLE / FLASHBACK query keywords
	kwBLOCK
	kwVERSIONS

	// MODEL clause keywords
	kwAUTOMATIC
	kwDECREMENT
	kwDIMENSION
	kwITERATE
	kwMAIN
	kwMEASURES
	kwNAV
	kwREFERENCE
	kwRULES
	kwSEQUENTIAL
	kwUNTIL
	kwUPDATED
	kwUPSERT

	// Analytic function keywords
	kwFOLLOWING
	kwGROUPS
	kwPRECEDING
	kwUNBOUNDED
	kwWITHIN

	// Grouping extension keywords
	kwCUBE
	kwGROUPING
	kwROLLUP
	kwSETS

	// LATERAL / XMLTABLE / JSON_TABLE keywords
	kwCOLUMNS
	kwJSON_TABLE
	kwLATERAL
	kwNESTED
	kwORDINALITY
	kwPASSING
	kwPATH
	kwXMLTABLE

	// Cursor / MULTISET / TREAT keywords
	kwEXCEPT
	kwTREAT

	// XML expression keywords
	kwCONTENT
	kwNAME
	kwXMLAGG
	kwXMLELEMENT
	kwXMLFOREST
	kwXMLPARSE
	kwXMLROOT
	kwXMLSERIALIZE

	// Utility statement keywords
	kwCALL
	kwSHARE

	// Session/audit keywords
	kwASSOCIATE
	kwDISASSOCIATE
	kwSTATISTICS
	kwSUCCESSFUL
	kwWHENEVER

	// PL/SQL additional keywords
	kwPIPE

	// Admin DDL keywords
	kwCONTEXT
	kwDIRECTORY
	kwJAVA
	kwLIBRARY
	kwPROFILE
	kwSCHEMA

	// JSON expression keywords
	kwFORMAT
	kwJSON
	kwJSON_ARRAY
	kwJSON_EXISTS
	kwJSON_MERGEPATCH
	kwJSON_OBJECT
	kwJSON_QUERY
	kwJSON_VALUE
	kwCOLUMN_VALUE
	kwMLSLABEL
	kwNESTED_TABLE_ID
	kwRESOURCE
	kwUID
)
