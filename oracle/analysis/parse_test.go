package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

func TestParse(t *testing.T) {
	list, err := Parse("SELECT 1 FROM DUAL;\nSELECT A FROM T")
	require.NoError(t, err)
	require.Len(t, list.Items, 2)
	second, ok := list.Items[1].(*ast.RawStmt)
	require.True(t, ok)
	require.Equal(t, 20, second.Stmt.(*ast.SelectStmt).Loc.Start)

	list, err = Parse("-- only a comment")
	require.NoError(t, err)
	require.Empty(t, list.Items)

	_, err = Parse("SELECT 1 FROM DUAL;\nSELEC 2")
	var parseErr *parser.ParseError
	require.ErrorAs(t, err, &parseErr)
	require.GreaterOrEqual(t, parseErr.Position, 20)
}
