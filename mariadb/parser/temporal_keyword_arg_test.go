package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestTemporalKeywordArgAST pins the node the unit selector becomes. Without
// the dedicated production the generic argument path builds a ColumnRef, which
// parses but means something else: a consumer resolving column lineage then
// looks for a column named SECOND, fails, and loses the whole statement's
// lineage. That is what this asserts against, not just that the node exists.
func TestTemporalKeywordArgAST(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		keyword string
	}{
		{"SELECT TIMESTAMPDIFF(SECOND, a, b) FROM t;", "SECOND"},
		{"SELECT TIMESTAMPADD(MINUTE, 5, a) FROM t;", "MINUTE"},
		{"SELECT GET_FORMAT(DATE, 'USA') FROM t;", "DATE"},
		// SQL_TSI_ synonyms and lower case both store the canonical name.
		{"SELECT TIMESTAMPDIFF(sql_tsi_second, a, b) FROM t;", "SECOND"},
		{"SELECT TIMESTAMPADD(SQL_TSI_YEAR, 1, a) FROM t;", "YEAR"},
		// GET_FORMAT's TIMESTAMP is stored as DATETIME.
		{"SELECT GET_FORMAT(TIMESTAMP, 'USA') FROM t;", "DATETIME"},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		var keywords []string
		var columns []string
		ast.Walk(keywordArgCollector(func(n ast.Node) {
			switch node := n.(type) {
			case *ast.KeywordArg:
				keywords = append(keywords, node.Keyword)
			case *ast.ColumnRef:
				columns = append(columns, node.Column)
			}
		}), list.Items[0])

		if len(keywords) != 1 || keywords[0] != tc.keyword {
			t.Errorf("%q: keyword args = %v, want [%s]", tc.sql, keywords, tc.keyword)
		}
		for _, column := range columns {
			if column == tc.keyword {
				t.Errorf("%q: the unit became a column reference %q", tc.sql, column)
			}
		}
	}
}

// keywordArgCollector adapts a visit function to ast.Visitor.
type keywordArgCollector func(ast.Node)

func (c keywordArgCollector) Visit(n ast.Node) ast.Visitor {
	c(n)
	return c
}
