package analysis

import (
	"errors"
	"fmt"
	"go/constant"
	"strconv"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/scanner"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// RewriteResultLimit rewrites one SELECT so it returns at most limitCount
// rows. A statement that is not a SELECT, VALUES or TABLE is returned
// unchanged, and so is one whose integer LIMIT is at most limitCount.
//
//   - Another integer, ALL or NULL after LIMIT or FETCH FIRST/NEXT is
//     replaced in the text by limitCount.
//   - With no LIMIT, "LIMIT <limitCount>" goes before a FOR UPDATE/SHARE
//     clause, or after the statement's last token.
//   - Any other LIMIT, an expression or placeholder, becomes
//     least(<it>, limitCount), rendered from the AST.
//
// The rewrite never wraps the statement, which CockroachDB refuses for a
// query with AS OF SYSTEM TIME. An error means empty input, a parse error or
// more than one statement.
func RewriteResultLimit(statement string, limitCount int) (string, error) {
	parsed, err := Parse(statement)
	if err != nil {
		return "", err
	}
	if len(parsed) != 1 {
		return "", fmt.Errorf("expected exactly one statement, got %d", len(parsed))
	}
	var sel *tree.Select
	switch n := parsed[0].AST.(type) {
	case *tree.Select:
		sel = n
	case *tree.ParenSelect:
		sel = &tree.Select{Select: n}
	default:
		return statement, nil
	}
	limitText := strconv.Itoa(limitCount)
	if sel.Limit == nil || (sel.Limit.Count == nil && !sel.Limit.LimitAll) {
		pos, beforeLocking := limitInsertPosition(statement, len(sel.Locking) > 0)
		if pos < 0 {
			return "", errors.New("cannot locate the end of the statement")
		}
		if beforeLocking {
			return statement[:pos] + "LIMIT " + limitText + " " + statement[pos:], nil
		}
		return statement[:pos] + " LIMIT " + limitText + statement[pos:], nil
	}
	if count, ok := integerLimit(sel.Limit); ok && count >= 0 && count <= int64(limitCount) {
		return statement, nil
	}
	if start, end, ok := limitCountSpan(statement); ok {
		return statement[:start] + limitText + statement[end:], nil
	}
	sel.Limit.Count = &tree.FuncExpr{
		Func:  tree.WrapFunction("least"),
		Exprs: tree.Exprs{sel.Limit.Count, tree.NewNumVal(constant.MakeInt64(int64(limitCount)), limitText, false)},
	}
	return tree.AsString(sel), nil
}

func integerLimit(limit *tree.Limit) (int64, bool) {
	num, ok := limit.Count.(*tree.NumVal)
	if limit.LimitAll || !ok {
		return 0, false
	}
	value, err := num.AsInt64()
	return value, err == nil
}

type limitToken struct {
	id  int32
	pos int
	end int
}

// topLevelTokens scans statement and returns its tokens outside parentheses
// and brackets, with each token's byte span, up to its first top-level ";".
func topLevelTokens(statement string) ([]limitToken, bool) {
	var sc scanner.SQLScanner
	sc.Init(statement)
	defer sc.Cleanup()
	var tokens []limitToken
	var lval scanSym
	depth := 0
	for {
		sc.Scan(&lval)
		switch lval.id {
		case 0:
			return tokens, true
		case lexbase.ERROR:
			return nil, false
		case '(', '[':
			depth++
		case ')', ']':
			depth--
		case ';':
			if depth == 0 {
				return tokens, true
			}
		}
		if depth == 0 {
			tokens = append(tokens, limitToken{id: lval.id, pos: int(lval.pos), end: sc.Pos()})
		}
	}
}

// limitInsertPosition returns where a LIMIT clause goes: the start of the
// top-level locking clause when the statement has one, or the end of the last
// token.
func limitInsertPosition(statement string, locking bool) (int, bool) {
	tokens, ok := topLevelTokens(statement)
	if !ok || len(tokens) == 0 {
		return -1, false
	}
	if locking {
		// FOR also names COLLATION FOR and may be an alias, so only FOR UPDATE,
		// FOR NO KEY UPDATE, FOR SHARE and FOR KEY SHARE start the clause.
		for i := 0; i+1 < len(tokens); i++ {
			if tokens[i].id != lexbase.FOR {
				continue
			}
			switch tokens[i+1].id {
			case lexbase.UPDATE, lexbase.SHARE, lexbase.NO, lexbase.KEY:
				return tokens[i].pos, true
			}
		}
	}
	return tokens[len(tokens)-1].end, false
}

// limitCountSpan returns the byte span of the row count of the top-level
// LIMIT or FETCH FIRST/NEXT clause when it is one integer, ALL or NULL token.
func limitCountSpan(statement string) (int, int, bool) {
	tokens, ok := topLevelTokens(statement)
	if !ok {
		return 0, 0, false
	}
	for i := len(tokens) - 1; i > 0; i-- {
		switch tokens[i].id {
		case lexbase.ICONST, lexbase.ALL, lexbase.NULL:
		default:
			continue
		}
		previous := tokens[i-1].id
		if previous == lexbase.LIMIT || (i > 1 && tokens[i-2].id == lexbase.FETCH && (previous == lexbase.FIRST || previous == lexbase.NEXT)) {
			// LIMIT 1 + 1 starts with an integer too; only a lone token is the count.
			if i+1 < len(tokens) && !endsLimitCount(tokens[i+1].id) {
				return 0, 0, false
			}
			return tokens[i].pos, tokens[i].end, true
		}
	}
	return 0, 0, false
}

// endsLimitCount reports whether a token may follow a LIMIT or FETCH count.
func endsLimitCount(id int32) bool {
	switch id {
	case lexbase.OFFSET, lexbase.FOR, lexbase.ROW, lexbase.ROWS:
		return true
	}
	return false
}
