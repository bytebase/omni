package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Objects read by their documented shape (ast.ObjectStmt)
// ---------------------------------------------------------------------------
//
// Many of Snowflake's newer objects (compute pools, services, applications,
// Postgres instances, ...) share one statement shape: a name, a few structural
// clauses, and a bag of `KEY = value` properties. Each such object registers an
// objectForm in an init function of its own file: the keywords that name it,
// and a function that reads the clauses of its CREATE and of its ALTER after
// the name, by the form its docs page gives. parseCreateStmt and parseAlterStmt
// ask the registry before their own dispatch, so a new object needs no case in
// either switch. The helpers below read the clauses most objects share.

// objectForm is one registered object.
type objectForm struct {
	words  []string // the object's keywords, uppercased: COMPUTE, POOL
	kind   string   // words joined by spaces
	create func(p *Parser, s *ast.ObjectStmt) error
	alter  func(p *Parser, s *ast.ObjectStmt) error
	// orReplace and orAlter report whether CREATE OR REPLACE and CREATE OR
	// ALTER are documented for the object.
	orReplace, orAlter bool
	// drop reports whether DROP <kind> [ IF EXISTS ] <name> is documented,
	// and dropCascade whether a CASCADE may follow the name.
	drop, dropCascade bool
	// nameParts is how many parts the object's name may have: 1 for an
	// account-level object (a compute pool), 2 for one a database qualifies,
	// 0 for a schema object, which db.schema.name may qualify.
	nameParts int
	// noIfNotExists and noIfExists are set where the page writes no
	// CREATE ... IF NOT EXISTS or ALTER ... IF EXISTS, and dropNoIfExists where
	// its DROP page writes no IF EXISTS.
	noIfNotExists, noIfExists, dropNoIfExists bool
	// props, when set, are the properties the page documents and the value
	// each takes. parseObjectStmt then rejects a property of a CREATE's run
	// or a SET that props does not name or whose value has another shape,
	// and an UNSET of a name outside unset (or props, when unset is nil).
	props objectProps
	unset []string
	// alterProps, when set, are the properties ALTER ... SET documents where
	// they differ from CREATE's props.
	alterProps objectProps
}

// objectForms maps an object's first keyword to the objects it starts.
var objectForms = map[string][]*objectForm{}

// registerObject adds an object to the registry. It is called from init.
func registerObject(f *objectForm) {
	f.kind = strings.Join(f.words, " ")
	objectForms[f.words[0]] = append(objectForms[f.words[0]], f)
}

// matchObjectForm reports the registered object whose keywords start at cur,
// preferring the longest. It consumes nothing.
func (p *Parser) matchObjectForm() *objectForm {
	if p.cur.Type != tokIdent && p.cur.Type < 700 {
		return nil
	}
	var best *objectForm
	for _, f := range objectForms[strings.ToUpper(p.cur.Str)] {
		if (best == nil || len(f.words) > len(best.words)) && p.atWords(f.words...) {
			best = f
		}
	}
	return best
}

// isWordToken reports whether t is the word upper as SQL spells it: an
// unquoted identifier or a keyword whose uppercased text is upper. A string or
// a quoted identifier never stands for a keyword.
func isWordToken(t Token, upper string) bool {
	return (t.Type == tokIdent || t.Type >= 700) && strings.ToUpper(t.Str) == upper
}

// atWords reports whether the tokens from cur on spell words, each an
// unquoted identifier or a keyword (isWordToken). It consumes nothing.
func (p *Parser) atWords(words ...string) bool {
	if !isWordToken(p.cur, words[0]) {
		return false
	}
	if len(words) == 1 {
		return true
	}
	if !isWordToken(p.peekNext(), words[1]) {
		return false
	}
	if len(words) == 2 {
		return true
	}
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance()
	p.advance()
	for _, w := range words[2:] {
		if !isWordToken(p.cur, w) {
			return false
		}
		p.advance()
	}
	return true
}

// acceptWords consumes words if the tokens from cur spell them, and reports
// whether it did.
func (p *Parser) acceptWords(words ...string) bool {
	if !p.atWords(words...) {
		return false
	}
	for range words {
		p.advance()
	}
	return true
}

// expectWords consumes words, or fails at the first token that differs.
func (p *Parser) expectWords(words ...string) error {
	for _, w := range words {
		if !isWordToken(p.cur, w) {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	return nil
}

// parseObjectStmt reads a CREATE or ALTER of a registered object. cur is the
// object's first keyword; start is the CREATE or ALTER token's Loc.
func (p *Parser) parseObjectStmt(verb string, f *objectForm, start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	read := f.create
	if verb == "ALTER" {
		read = f.alter
	}
	if read == nil || (orReplace && !f.orReplace) || (orAlter && !f.orAlter) {
		return p.unsupported(verb)
	}
	for range f.words {
		p.advance()
	}
	s := &ast.ObjectStmt{
		Verb:      verb,
		Kind:      f.kind,
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}
	// IF NOT EXISTS excludes OR REPLACE (the docs write the two forms apart).
	if verb == "CREATE" && !orReplace && !f.noIfNotExists {
		if err := p.parseIfNotExists(&s.IfNotExists); err != nil {
			return nil, err
		}
	}
	if verb == "ALTER" && !f.noIfExists {
		if err := p.parseIfExists(&s.IfExists); err != nil {
			return nil, err
		}
	}
	name, err := p.parseFormName(f)
	if err != nil {
		return nil, err
	}
	s.Name = name
	if err := read(p, s); err != nil {
		return nil, err
	}
	if f.props != nil {
		if err := f.checkClauses(s); err != nil {
			return nil, err
		}
	}
	s.Loc.End = p.prev.Loc.End
	return s, nil
}

// parseObjectSimpleName reads a name a clause takes where the page writes a
// single <name>: one identifier part, with no database or schema.
func (p *Parser) parseObjectSimpleName() (*ast.ObjectName, error) {
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	if !name.Schema.IsEmpty() || !name.Database.IsEmpty() {
		return nil, &ParseError{
			Position: name.Loc.Start, End: name.Loc.End,
			Message: "syntax error: " + name.String() + " takes no qualifier here",
		}
	}
	return name, nil
}

// parseObjectSimpleNameList reads <name> [ , ... ] of single-part names.
func (p *Parser) parseObjectSimpleNameList() ([]*ast.ObjectName, error) {
	var names []*ast.ObjectName
	for {
		name, err := p.parseObjectSimpleName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.cur.Type != ',' {
			return names, nil
		}
		p.advance() // consume ','
	}
}

// parseFormName reads the name of an object of form f, holding it to the
// number of parts f.nameParts allows.
func (p *Parser) parseFormName(f *objectForm) (*ast.ObjectName, error) {
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	return p.checkFormNameParts(f, name)
}

// checkFormNameParts fails when name has more parts than f.nameParts.
func (p *Parser) checkFormNameParts(f *objectForm, name *ast.ObjectName) (*ast.ObjectName, error) {
	parts := 1
	if !name.Schema.IsEmpty() {
		parts++
	}
	if !name.Database.IsEmpty() {
		parts++
	}
	if f.nameParts > 0 && parts > f.nameParts {
		return nil, &ParseError{
			Position: name.Loc.Start, End: name.Loc.End,
			Message: fmt.Sprintf("syntax error: a %s's name has at most %d part(s)", strings.ToLower(f.kind), f.nameParts),
		}
	}
	return name, nil
}

// objectValue is the shape of value a documented property takes.
type objectValue struct {
	kind  objectValueKind
	words []string // valWord: the keywords it may be, uppercased (nil takes any name); valStringOr: those besides a string
}

type objectValueKind int

const (
	valNumber   objectValueKind = iota + 1 // <num>, optionally signed
	valString                              // '<string>'
	valWord                                // a one-part name, or one of words
	valQName                               // a name db.schema.name may qualify
	valStrings                             // ( '<string>' [ , ... ] )
	valNames                               // ( <name> [ , ... ] )
	valGroup                               // ( <key> = <value> [ , ... ] ), as the page spells it
	valStringOr                            // '<string>', or one of words
)

// objectProps maps a documented property's name to its value's shape.
type objectProps map[string]objectValue

var (
	propNumber  = objectValue{kind: valNumber}
	propString  = objectValue{kind: valString}
	propName    = objectValue{kind: valWord}
	propQName   = objectValue{kind: valQName}
	propBool    = objectValue{kind: valWord, words: []string{"TRUE", "FALSE"}}
	propStrings = objectValue{kind: valStrings}
	propNames   = objectValue{kind: valNames}
	propGroup   = objectValue{kind: valGroup}
)

// propOneOf is a property whose value is one of words.
func propOneOf(words ...string) objectValue { return objectValue{kind: valWord, words: words} }

// propStringOr is a property whose value is a '<string>' or one of words, as
// EVENT_TABLE = { '<name>' | NONE }.
func propStringOr(words ...string) objectValue { return objectValue{kind: valStringOr, words: words} }

// checkClauses holds a statement's property runs, SETs and UNSETs to f's
// documented properties.
func (f *objectForm) checkClauses(s *ast.ObjectStmt) error {
	props := f.props
	if s.Verb == "ALTER" && f.alterProps != nil {
		props = f.alterProps
	}
	for _, c := range s.Clauses {
		switch c.Keywords {
		case "", "SET":
			for _, o := range c.Options {
				v, ok := props[o.Name]
				if !ok || !v.matches(o) {
					return &ParseError{
						Position: o.Loc.Start, End: o.Loc.End,
						Message: "syntax error at or near " + o.Name,
					}
				}
			}
		case "UNSET":
			for _, k := range c.Keys {
				if !f.unsets(k) {
					return &ParseError{
						Position: c.Loc.Start, End: c.Loc.End,
						Message: "syntax error: " + s.Kind + " has no property " + k + " to unset",
					}
				}
			}
		}
	}
	return nil
}

// unsets reports whether the page documents UNSET of property k.
func (f *objectForm) unsets(k string) bool {
	if f.unset == nil {
		_, ok := f.props[k]
		return ok
	}
	for _, u := range f.unset {
		if u == k {
			return true
		}
	}
	return false
}

// matches reports whether option o holds a value of v's shape.
func (v objectValue) matches(o *ast.CopyOption) bool {
	switch v.kind {
	case valNumber:
		return o.Lit != nil && (o.Lit.Kind == ast.LitInt || o.Lit.Kind == ast.LitFloat)
	case valString:
		return o.Lit != nil && o.Lit.Kind == ast.LitString
	case valQName:
		return o.Words != "" && nameParts(o.Words) <= 3
	case valWord:
		if o.Words == "" || nameParts(o.Words) != 1 {
			return false
		}
		if v.words == nil {
			return true
		}
		for _, w := range v.words {
			if o.Words == w {
				return true
			}
		}
		return false
	case valStringOr:
		if o.Lit != nil && o.Lit.Kind == ast.LitString {
			return true
		}
		for _, w := range v.words {
			if o.Words == w {
				return true
			}
		}
		return false
	case valStrings:
		return o.List != nil && len(o.List) > 0 && allStrings(o.List)
	case valNames:
		if len(o.Group) == 0 {
			return false
		}
		for _, e := range o.Group {
			if !e.Bare {
				return false
			}
		}
		return true
	case valGroup:
		if len(o.Group) == 0 {
			return false
		}
		for _, e := range o.Group {
			if e.Bare {
				return false
			}
		}
		return true
	}
	return false
}

func allStrings(lits []*ast.Literal) bool {
	for _, l := range lits {
		if l.Kind != ast.LitString {
			return false
		}
	}
	return true
}

// startsObjectProperty reports whether cur begins a `KEY = value` property:
// a name word followed by '='.
// The key is an unquoted identifier or a keyword: a property's name is never
// quoted.
func (p *Parser) startsObjectProperty() bool {
	return (p.cur.Type == tokIdent || p.cur.Type >= 700) && p.peekNext().Type == '='
}

// parseObjectProperties reads a run of `KEY = value` properties and
// [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) clauses into one clause, and
// appends and returns it when it holds anything (nil otherwise). need says
// whether an empty run is a syntax error. keywords names the clause ("" for a
// CREATE's run).
func (p *Parser) parseObjectProperties(s *ast.ObjectStmt, keywords string, need bool) (*ast.ObjectClause, error) {
	c := &ast.ObjectClause{Keywords: keywords, Loc: ast.Loc{Start: p.cur.Loc.Start}}
	seen := map[string]bool{}
	for {
		tagStart := p.cur
		switch {
		case keywords == "" && p.startsWithTags() && (p.cur.Type == kwWITH || p.peekNext().Type == '('):
			// [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) is a CREATE's; an
			// ALTER sets tags with SET TAG <tag> = '<value>'.
			if c.Tags != nil {
				return nil, &ParseError{Position: tagStart.Loc.Start, End: tagStart.Loc.End, Message: "syntax error: TAG is set twice"}
			}
			p.match(kwWITH)
			p.advance() // consume TAG
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			tags, err := p.parseObjectTagList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			c.Tags = tags
		case keywords == "SET" && p.cur.Type == kwTAG && p.peekNext().Type != '(':
			// Some pages list TAG <tag> = '<value>' [ , ... ] among SET's
			// properties (ALTER COMPUTE POOL).
			if c.Tags != nil {
				return nil, &ParseError{Position: tagStart.Loc.Start, End: tagStart.Loc.End, Message: "syntax error: TAG is set twice"}
			}
			p.advance() // consume TAG
			tags, err := p.parseObjectTagList()
			if err != nil {
				return nil, err
			}
			c.Tags = tags
		case p.startsObjectProperty():
			opt, err := p.parseObjectProperty()
			if err != nil {
				return nil, err
			}
			if seen[opt.Name] {
				// A property is set once a run.
				return nil, &ParseError{
					Position: opt.Loc.Start, End: opt.Loc.End,
					Message: "syntax error: " + opt.Name + " is set twice",
				}
			}
			seen[opt.Name] = true
			c.Options = append(c.Options, opt)
		default:
			if len(c.Options) == 0 && len(c.Tags) == 0 {
				if need {
					return nil, p.syntaxErrorAtCur()
				}
				return nil, nil
			}
			c.Loc.End = p.prev.Loc.End
			s.Clauses = append(s.Clauses, c)
			return c, nil
		}
	}
}

// objectActions selects the shared ALTER actions an object's docs give.
type objectActions int

const (
	actSet      objectActions = 1 << iota // SET <property> [ ... ]
	actSetTag                             // SET TAG <tag> = '<value>' [ , ... ]
	actUnset                              // UNSET <property> [ , ... ]
	actUnsetTag                           // UNSET TAG <tag> [ , ... ]
	actRename                             // RENAME TO <new_name>
	// actTags is SET TAG and UNSET TAG.
	actTags = actSetTag | actUnsetTag
)

// parseObjectAlterCommon reads one of the shared ALTER actions allow names,
// when one starts at cur, and reports whether it did.
func (p *Parser) parseObjectAlterCommon(s *ast.ObjectStmt, allow objectActions) (bool, error) {
	start := p.cur.Loc.Start
	switch {
	case p.cur.Type == kwSET && p.peekNext().Type == kwTAG && allow&actSetTag != 0:
		p.advance() // consume SET
		p.advance() // consume TAG
		tags, err := p.parseObjectTagList()
		if err != nil {
			return true, err
		}
		s.Clauses = append(s.Clauses, &ast.ObjectClause{Keywords: "SET TAG", Tags: tags, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
		return true, nil
	case p.cur.Type == kwSET && allow&actSet != 0:
		p.advance() // consume SET
		c, err := p.parseObjectProperties(s, "SET", true)
		if err != nil {
			return true, err
		}
		c.Loc.Start = start
		return true, nil
	case p.cur.Type == kwUNSET && p.peekNext().Type == kwTAG && allow&actUnsetTag != 0:
		p.advance() // consume UNSET
		p.advance() // consume TAG
		names, err := p.parseObjectNameList()
		if err != nil {
			return true, err
		}
		s.Clauses = append(s.Clauses, &ast.ObjectClause{Keywords: "UNSET TAG", Names: names, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
		return true, nil
	case p.cur.Type == kwUNSET && allow&actUnset != 0:
		p.advance() // consume UNSET
		keys, err := p.parseObjectUnsetKeys()
		if err != nil {
			return true, err
		}
		s.Clauses = append(s.Clauses, &ast.ObjectClause{Keywords: "UNSET", Keys: keys, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
		return true, nil
	case p.cur.Type == kwRENAME && allow&actRename != 0:
		name, err := p.parseRenameTo()
		if err != nil {
			return true, err
		}
		s.Clauses = append(s.Clauses, &ast.ObjectClause{Keywords: "RENAME TO", Names: []*ast.ObjectName{name}, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
		return true, nil
	}
	return false, nil
}

// requireObjectProperties fails unless clause c sets every one of names, as
// a CREATE's required properties are.
func (p *Parser) requireObjectProperties(c *ast.ObjectClause, names ...string) error {
	for _, n := range names {
		found := false
		if c != nil {
			for _, o := range c.Options {
				if o.Name == n {
					found = true
					break
				}
			}
		}
		if !found {
			return &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "syntax error: " + n + " is required",
			}
		}
	}
	return nil
}

// parseObjectTagList reads <tag> = '<value>' [ , ... ] after SET TAG.
func (p *Parser) parseObjectTagList() ([]*ast.TagAssignment, error) {
	var tags []*ast.TagAssignment
	for {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		val, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &ast.TagAssignment{Name: name, Value: val.Str})
		if p.cur.Type != ',' {
			return tags, nil
		}
		p.advance() // consume ','
	}
}

// addObjectKeywordsClause appends a clause of keywords alone (SUSPEND, RESUME,
// STOP ALL, ...) once the caller has consumed them; start is where they began.
func (p *Parser) addObjectKeywordsClause(s *ast.ObjectStmt, keywords string, start int) *ast.ObjectClause {
	c := &ast.ObjectClause{Keywords: keywords, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}}
	s.Clauses = append(s.Clauses, c)
	return c
}

// parseObjectProperty reads one `KEY = value` property. cur is the key, which
// startsObjectProperty has seen followed by '='.
func (p *Parser) parseObjectProperty() (*ast.CopyOption, error) {
	key := p.advance()
	opt := &ast.CopyOption{Name: strings.ToUpper(key.Str), Loc: ast.Loc{Start: key.Loc.Start}}
	p.advance() // consume '='
	if err := p.parseObjectValue(opt); err != nil {
		return nil, err
	}
	opt.Loc.End = p.prev.Loc.End
	return opt, nil
}

// parseObjectValue reads a property's value into opt, in the shapes the docs
// give a property:
//
//	'<string>' | <number> | <word>         -> Lit, or Words
//	( )                                     -> an empty Group
//	( <literal> [ , ... ] )                 -> List
//	( <entry> [ , ... ] )                   -> Group, where an entry is
//	                                           <word> or '<string>' = <value>,
//	                                           a bare <word>, or ( ... )
//
// Unlike COPY's options, a list needs its commas and an entry its value.
func (p *Parser) parseObjectValue(opt *ast.CopyOption) error {
	switch p.cur.Type {
	case tokString, tokInt, tokFloat, tokReal:
		lit, err := p.parseCopyLiteral()
		if err != nil {
			return err
		}
		opt.Lit = lit
		return nil
	case '(':
		return p.parseObjectValueParen(opt)
	}
	if p.cur.Type == '-' && (p.peekNext().Type == tokInt || p.peekNext().Type == tokFloat || p.peekNext().Type == tokReal) {
		minus := p.advance()
		lit, err := p.parseCopyLiteral()
		if err != nil {
			return err
		}
		lit.Value = "-" + lit.Value
		lit.Ival = -lit.Ival
		lit.Loc.Start = minus.Loc.Start
		opt.Lit = lit
		return nil
	}
	words, err := p.parseObjectWord()
	if err != nil {
		return err
	}
	opt.Words = words
	return nil
}

// parseObjectWord reads a word value: a name, or a qualified name whose parts
// are joined by '.', each part an identifier or a keyword. Keywords and
// unquoted names are uppercased; a quoted part keeps its quotes, so "TRUE"
// never reads as the keyword TRUE.
func (p *Parser) parseObjectWord() (string, error) {
	var b strings.Builder
	for {
		switch {
		case p.cur.Type == tokQuotedIdent:
			tok := p.advance()
			b.WriteString(p.srcSlice(tok.Loc.Start, tok.Loc.End))
		case p.cur.Type == tokIdent || p.cur.Type >= 700:
			b.WriteString(strings.ToUpper(p.advance().Str))
		default:
			return "", p.syntaxErrorAtCur()
		}
		if p.cur.Type != '.' {
			return b.String(), nil
		}
		p.advance() // consume '.'
		b.WriteByte('.')
	}
}

// parseObjectValueParen reads a parenthesized property value. cur is '('.
func (p *Parser) parseObjectValueParen(opt *ast.CopyOption) error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	p.advance()                     // consume '('
	opt.Group = []*ast.CopyOption{} // "()" is recorded as an empty group
	if p.cur.Type == ')' {
		p.advance()
		return nil
	}
	// A list of literals: ( 's' [ , ... ] ).
	if p.curIsLiteral() && p.peekNext().Type != '=' {
		opt.Group = nil
		for {
			lit, err := p.parseCopyLiteral()
			if err != nil {
				return err
			}
			opt.List = append(opt.List, lit)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		_, err := p.expect(')')
		return err
	}
	for {
		entry := &ast.CopyOption{Loc: ast.Loc{Start: p.cur.Loc.Start}}
		switch {
		case p.cur.Type == '(':
			// A list of groups: ACCEPTED_TOKENS = ( ( ... ) , ( ... ) ).
			if err := p.parseObjectValueParen(entry); err != nil {
				return err
			}
		case (p.cur.Type == tokString || p.startsObjectProperty()) && p.peekNext().Type == '=':
			key := p.advance()
			entry.Name = strings.ToUpper(key.Str)
			if key.Type == tokString {
				// A string key (SECRETS = ( '<name>' = <secret> )) keeps its
				// quotes, so it stays apart from a word key.
				entry.Name = p.srcSlice(key.Loc.Start, key.Loc.End)
			}
			p.advance() // consume '='
			if err := p.parseObjectValue(entry); err != nil {
				return err
			}
		case p.isOptionWord(p.cur.Type) && !p.curIsLiteral():
			// A bare name in a list of names: ( a [ , ... ] ).
			words, err := p.parseObjectWord()
			if err != nil {
				return err
			}
			entry.Name = words
			entry.Bare = true
		default:
			return p.syntaxErrorAtCur()
		}
		entry.Loc.End = p.prev.Loc.End
		opt.Group = append(opt.Group, entry)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(')')
	return err
}

// parseObjectAsQuery reads AS <query> into a clause, as CORTEX SEARCH
// SERVICE's body is written. cur is AS.
func (p *Parser) parseObjectAsQuery(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwAS); err != nil {
		return err
	}
	q, err := p.parseQueryExpr()
	if err != nil {
		return err
	}
	s.Clauses = append(s.Clauses, &ast.ObjectClause{Keywords: "AS", Query: q, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
	return nil
}

// parseObjectUnsetKeys reads UNSET's <property> [ , ... ]: each an unquoted
// identifier or a keyword, never a number, string or quoted name.
func (p *Parser) parseObjectUnsetKeys() ([]string, error) {
	var keys []string
	for {
		if p.cur.Type != tokIdent && p.cur.Type < 700 {
			return nil, p.syntaxErrorAtCur()
		}
		keys = append(keys, strings.ToUpper(p.advance().Str))
		if p.cur.Type != ',' {
			return keys, nil
		}
		p.advance() // consume ','
	}
}

// nameParts counts the '.'-joined parts of a word value parseObjectWord read,
// leaving out dots inside a quoted part.
func nameParts(words string) int {
	n, quoted := 1, false
	for _, r := range words {
		switch {
		case r == '"':
			quoted = !quoted
		case r == '.' && !quoted:
			n++
		}
	}
	return n
}
