package completion

import (
	"cmp"
	"context"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/clickhouse/parser"
	metadatapb "github.com/bytebase/omni/metadata"
)

// SessionConfig configures session-aware ClickHouse completion.
type SessionConfig struct {
	// DefaultDatabase is the current database, which resolves names with
	// no database.
	DefaultDatabase string
	// LoadDatabase fetches a database's metadata: DefaultDatabase's, and
	// another database's when a name in the statement is qualified by it.
	// When nil or failing, completion offers keywords only.
	LoadDatabase func(ctx context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error)
	// ListDatabaseNames lists the instance's databases. When nil or
	// failing, DefaultDatabase is the one database offered.
	ListDatabaseNames func(ctx context.Context) ([]string, error)
}

// CompleteSession completes the statement of a script at the 1-based line
// and 0-based UTF-16 code-unit caretOffset, as an editor reports a caret.
func CompleteSession(ctx context.Context, cfg SessionConfig, statement string, caretLine, caretOffset int) []Candidate {
	stmt, pos := statementAtCaret(statement, lineOffsetToBytePos(statement, caretLine, caretOffset))
	got := complete(ctx, stmt, pos, newCatalog(ctx, cfg))

	type key struct {
		typ  CandidateType
		text string
	}
	seen := make(map[key]bool, len(got))
	result := make([]Candidate, 0, len(got))
	for _, c := range got {
		k := key{c.Type, c.Text}
		if seen[k] {
			continue
		}
		seen[k] = true
		result = append(result, c)
	}
	slices.SortFunc(result, func(a, b Candidate) int {
		return cmp.Or(cmp.Compare(a.Type, b.Type), cmp.Compare(a.Text, b.Text))
	})
	return result
}

// bareIdentifier is an identifier ClickHouse reads unquoted.
var bareIdentifier = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// QuoteIdentifierIfNeeded backquotes name unless ClickHouse reads it back
// bare: a name with other characters is quoted, and so is a keyword the
// grammar does not read as a column and a table name, as it reads id. TRUE
// and NULL parse there, but as literals, so they are quoted too.
func QuoteIdentifierIfNeeded(name string) string {
	if bareIdentifier.MatchString(name) && readsAsName(name) {
		return name
	}
	return "`" + strings.ReplaceAll(strings.ReplaceAll(name, `\`, `\\`), "`", "``") + "`"
}

// readsAsName reports whether the parser reads name, unquoted, as a column
// in a select list and a WHERE clause, and as a table.
func readsAsName(name string) bool {
	stmts, err := parser.Parse("SELECT " + name + ", " + name + " FROM " + name + " WHERE " + name)
	if err != nil || len(stmts) != 1 {
		return false
	}
	sel, ok := stmts[0].(*parser.SelectStmt)
	if u, isUnion := stmts[0].(*parser.SelectUnionStmt); isUnion && len(u.Items) == 1 {
		sel, ok = u.Items[0], true
	}
	if !ok || len(sel.Columns) != 2 {
		return false
	}
	for _, e := range append(slices.Clone(sel.Columns), sel.Where) {
		if id, ok := e.(*parser.IdentExpr); !ok || len(id.Parts) != 1 || id.Parts[0] != name {
			return false
		}
	}
	return true
}

// lineOffsetToBytePos converts a 1-based line and 0-based UTF-16 code-unit
// column into a byte offset into s.
func lineOffsetToBytePos(s string, line, offset int) int {
	currentLine := 1
	col := 0
	for i := 0; i < len(s); {
		if currentLine == line && col >= offset {
			return i
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == '\n' {
			if currentLine == line {
				return i
			}
			currentLine++
			col = 0
		} else if r <= 0xFFFF {
			col++
		} else {
			col += 2
		}
		i += size
	}
	return len(s)
}

// statementAtCaret isolates the statement holding bytePos, so neither
// earlier statements nor broken later ones affect completion or its cost.
// A segment ends with its semicolon and any comment on the same line, so a
// caret after the semicolon is in the text after it.
func statementAtCaret(script string, bytePos int) (string, int) {
	bytePos = max(0, min(bytePos, len(script)))
	for _, seg := range parser.Split(script) {
		if bytePos < seg.ByteStart || bytePos > seg.ByteEnd {
			continue
		}
		rel := bytePos - seg.ByteStart
		for _, tok := range significantTokens(seg.Text) {
			if tok.Type == parser.TokSemicolon && tok.End <= rel {
				return seg.Text[tok.End:], rel - tok.End
			}
		}
		return seg.Text, rel
	}
	return "", 0
}
