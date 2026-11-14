package grammar

import (
	"fmt"
	"os"
	"strings"
)

// Forms a tokens.tsv row can give besides a token's text.
const (
	// Empty: the token is written as nothing, such as a hint the lexer
	// reads from a comment.
	Empty = "(empty)"
	// Unreachable: the lexer never produces the token in the mode the judge
	// runs in, so an alternative that needs it gets no generated row and
	// needs a hand-added one.
	Unreachable = "(unreachable)"
	// Spaced: a literal the grammar writes in quotes, named here by its
	// text, is written with a space on each side, where it would otherwise
	// sit tight against its neighbours. OceanBase's lexer reads a word, a
	// dot and a word as one token, so a dot that is a token of its own must
	// stand apart.
	Spaced = "(spaced)"
	// Reserved: a word, named by itself, that the engine never reads as an
	// identifier. A recognizer matches it only where the grammar spells it,
	// never for a token class such as a name.
	Reserved = "(reserved)"
	// Qualified, at the start of a row's note: the token is a name a
	// recognizer also matches qualified, two or three of them joined by
	// dots, as Snowflake writes db.schema.table wherever it takes a name.
	Qualified = "(qualified)"
	// Path, at the start of a row's note: the token's sample is a variable,
	// and a recognizer reads its sigil as the start of a path, one token
	// through any ~ or % after the sigil and every name character, dot and
	// slash after that, as Snowflake writes stages: @~, @%t, @db.s.a/dir; a
	// % may also follow a dot, as a namespace's table stage, @db.s.%t, does.
	// A string that starts with the sigil, '@a/dir with spaces', is one too.
	Path = "(path)"
)

// Tokens is an engine's testdata/tokens.tsv: the text generated rows write
// for each terminal whose grammar declaration does not spell it.
//
// A row may name a rule instead. Its form is then the text every other
// rule's row writes where that rule stands, in place of the rule's shortest
// derivation. It is for a name the engine's parser looks up while it
// parses, where the grammar gives the name a rule of its own: MariaDB's
// charset_name is any identifier, and an unknown one answers before a later
// syntax error, so the rows that pass through it would all read as
// accepted. It is also for the plainest expression, where the grammar's
// order does not make it a name: PartiQL lists CURRENT_USER first among its
// one-token terms, so exprTerm is sampled as a name. The rule's own
// alternatives keep their rows, written from their own symbols.
type Tokens struct {
	forms     map[string]string
	qualified map[string]bool
	path      map[string]bool
}

// ReadTokens reads a tokens.tsv. Each row is token, form, and a note citing
// the vendor's lexer or lexical docs; the note may be empty.
func ReadTokens(path string) (*Tokens, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	t := &Tokens{forms: map[string]string{}, qualified: map[string]bool{}, path: map[string]bool{}}
	for i, line := range strings.Split(string(b), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) < 2 || len(f) > 3 {
			return nil, fmt.Errorf("%s:%d: want token, form, and an optional note, tab-separated", path, i+1)
		}
		name, form := f[0], f[1]
		if name == "" || form == "" {
			return nil, fmt.Errorf("%s:%d: empty token or form; write %s for a token written as nothing", path, i+1, Empty)
		}
		if _, dup := t.forms[name]; dup {
			return nil, fmt.Errorf("%s:%d: %s appears twice", path, i+1, name)
		}
		t.forms[name] = form
		if len(f) == 3 && strings.HasPrefix(f[2], Qualified) {
			t.qualified[name] = true
		}
		if len(f) == 3 && strings.HasPrefix(f[2], Path) {
			t.path[name] = true
		}
	}
	return t, nil
}
