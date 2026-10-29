package mssql

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
)

type limitEdit struct {
	start int
	end   int
	text  string
}

// limitSource is the statement being rewritten, with its tokens lexed once
// on first use. topInsertPosition used to lex the text before every SELECT
// it adds TOP to, so a chain of n UNION ALL operands lexed O(n²) tokens.
type limitSource struct {
	sql    string
	tokens []parser.Token
	lexed  bool
}

func (s *limitSource) lexTokens() []parser.Token {
	if !s.lexed {
		s.tokens = parser.Tokenize(s.sql)
		s.lexed = true
	}
	return s.tokens
}

// StatementWithResultLimit rewrites a single T-SQL SELECT statement so it
// returns at most limit rows.
func StatementWithResultLimit(statement string, limit int) (string, error) {
	if limit <= 0 {
		return "", fmt.Errorf("limit must be positive")
	}
	if strings.TrimSpace(statement) == "" {
		return "", fmt.Errorf("empty statement")
	}

	stmts, err := Parse(statement)
	if err != nil {
		return "", err
	}
	if len(stmts) != 1 {
		return "", fmt.Errorf("expected exactly 1 statement, got %d", len(stmts))
	}

	stmt, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok {
		return normalizeStatementTerminator(statement), nil
	}

	var edits []limitEdit
	rewriteSelectWithLimit(&edits, &limitSource{sql: statement}, stmt, limit)
	rewritten := applyLimitEdits(statement, edits)
	return normalizeStatementTerminator(rewritten), nil
}

func rewriteSelectWithLimit(edits *[]limitEdit, src *limitSource, stmt *ast.SelectStmt, limit int) {
	if stmt == nil {
		return
	}
	if stmt.Op != ast.SetOpNone {
		if hasTopClause(stmt) {
			rewriteExistingTopClauses(edits, stmt, limit)
			return
		}
		if ordered := findOrderedSelect(stmt); ordered != nil {
			rewriteOrderedSelect(edits, src.sql, ordered, limit)
			return
		}
		addTopToSelectLeaves(edits, src, stmt, limit)
		return
	}
	if stmt.Top != nil {
		rewriteTopClause(edits, stmt.Top, limit)
		return
	}
	if stmt.OrderByClause != nil {
		rewriteOrderedSelect(edits, src.sql, stmt, limit)
		return
	}
	addTopClause(edits, src, stmt, limit)
}

func rewriteExistingTopClauses(edits *[]limitEdit, stmt *ast.SelectStmt, limit int) {
	if stmt == nil {
		return
	}
	if stmt.Op != ast.SetOpNone {
		rewriteExistingTopClauses(edits, stmt.Larg, limit)
		rewriteExistingTopClauses(edits, stmt.Rarg, limit)
		return
	}
	if stmt.Top != nil {
		rewriteTopClause(edits, stmt.Top, limit)
	}
}

func addTopToSelectLeaves(edits *[]limitEdit, src *limitSource, stmt *ast.SelectStmt, limit int) {
	if stmt == nil {
		return
	}
	if stmt.Op != ast.SetOpNone {
		addTopToSelectLeaves(edits, src, stmt.Larg, limit)
		addTopToSelectLeaves(edits, src, stmt.Rarg, limit)
		return
	}
	addTopClause(edits, src, stmt, limit)
}

func rewriteTopClause(edits *[]limitEdit, top *ast.TopClause, limit int) {
	if top == nil || !top.Loc.IsValid() {
		return
	}
	topLimit := extractInt(top.Count)
	if !top.Percent && !top.WithTies && topLimit > 0 && topLimit <= limit {
		return
	}
	*edits = append(*edits, limitEdit{
		start: top.Loc.Start,
		end:   top.Loc.End,
		text:  fmt.Sprintf("TOP %d", limit),
	})
}

func addTopClause(edits *[]limitEdit, src *limitSource, stmt *ast.SelectStmt, limit int) {
	if stmt == nil || stmt.TargetList == nil {
		return
	}
	insertPos := topInsertPosition(src, stmt)
	if insertPos < 0 {
		return
	}
	*edits = append(*edits, limitEdit{
		start: insertPos,
		end:   insertPos,
		text:  fmt.Sprintf(" TOP %d", limit),
	})
}

// topInsertPosition returns the offset just past the SELECT, ALL, or DISTINCT
// keyword that TOP follows: the last such token before the target list.
func topInsertPosition(src *limitSource, stmt *ast.SelectStmt) int {
	targetLoc := ast.NodeLoc(stmt.TargetList)
	if !targetLoc.IsValid() {
		return -1
	}
	want := "SELECT"
	if stmt.Distinct {
		want = "DISTINCT"
	} else if stmt.All {
		want = "ALL"
	}
	// The target list starts a token, so the statement's tokens that end by
	// its start are exactly the tokens of the text before it.
	tokens := src.lexTokens()
	n := sort.Search(len(tokens), func(i int) bool { return tokens[i].End > targetLoc.Start })
	for i := n - 1; i >= 0; i-- {
		if parser.TokenName(tokens[i].Type) == want {
			return tokens[i].End
		}
	}
	return -1
}

func rewriteOrderedSelect(edits *[]limitEdit, sql string, stmt *ast.SelectStmt, limit int) {
	if stmt.FetchClause != nil {
		fetchLimit := extractInt(stmt.FetchClause.Count)
		if fetchLimit > 0 && fetchLimit <= limit {
			return
		}
		loc := ast.NodeLoc(stmt.FetchClause.Count)
		if loc.IsValid() {
			*edits = append(*edits, limitEdit{start: loc.Start, end: loc.End, text: fmt.Sprintf("%d", limit)})
		}
		return
	}
	if stmt.OffsetClause != nil {
		loc := ast.NodeLoc(stmt.OffsetClause)
		if loc.IsValid() {
			pos := offsetRowsEnd(sql, loc.End)
			*edits = append(*edits, limitEdit{start: pos, end: pos, text: fmt.Sprintf(" FETCH NEXT %d ROWS ONLY", limit)})
		}
		return
	}
	loc := ast.NodeLoc(stmt.OrderByClause)
	if loc.IsValid() {
		*edits = append(*edits, limitEdit{
			start: loc.End,
			end:   loc.End,
			text:  fmt.Sprintf(" OFFSET 0 ROWS FETCH NEXT %d ROWS ONLY", limit),
		})
	}
}

func hasTopClause(stmt *ast.SelectStmt) bool {
	if stmt == nil {
		return false
	}
	if stmt.Op != ast.SetOpNone {
		return hasTopClause(stmt.Larg) || hasTopClause(stmt.Rarg)
	}
	return stmt.Top != nil
}

func findOrderedSelect(stmt *ast.SelectStmt) *ast.SelectStmt {
	if stmt == nil {
		return nil
	}
	if stmt.OrderByClause != nil {
		return stmt
	}
	if ordered := findOrderedSelect(stmt.Larg); ordered != nil {
		return ordered
	}
	return findOrderedSelect(stmt.Rarg)
}

func extractInt(node ast.ExprNode) int {
	lit, ok := node.(*ast.Literal)
	if !ok || lit.Type != ast.LitInteger {
		return 0
	}
	return int(lit.Ival)
}

func offsetRowsEnd(sql string, pos int) int {
	pos = skipSpaces(sql, pos)
	for _, keyword := range []string{"ROWS", "ROW"} {
		if end, ok := consumeKeyword(sql, pos, keyword); ok {
			return end
		}
	}
	return pos
}

func skipSpaces(sql string, pos int) int {
	for pos < len(sql) {
		r, size := utf8.DecodeRuneInString(sql[pos:])
		if !unicode.IsSpace(r) {
			break
		}
		pos += size
	}
	return pos
}

func consumeKeyword(sql string, pos int, keyword string) (int, bool) {
	if len(sql[pos:]) < len(keyword) || !strings.EqualFold(sql[pos:pos+len(keyword)], keyword) {
		return pos, false
	}
	end := pos + len(keyword)
	if end < len(sql) {
		r := rune(sql[end])
		if r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
			return pos, false
		}
	}
	return end, true
}

// applyLimitEdits applies the edits in one pass. The rewrites touch disjoint
// ranges of the statement; splicing each into the whole string in turn, as
// this used to, copied the statement once per edit.
func applyLimitEdits(sql string, edits []limitEdit) string {
	sort.SliceStable(edits, func(i, j int) bool {
		return edits[i].start < edits[j].start
	})
	var b strings.Builder
	last := 0
	for _, edit := range edits {
		if edit.start < last || edit.end < edit.start || edit.end > len(sql) {
			continue
		}
		b.WriteString(sql[last:edit.start])
		b.WriteString(edit.text)
		last = edit.end
	}
	b.WriteString(sql[last:])
	return b.String()
}

func normalizeStatementTerminator(sql string) string {
	return strings.TrimRightFunc(sql, func(r rune) bool {
		return unicode.IsSpace(r) || r == ';'
	}) + ";"
}
