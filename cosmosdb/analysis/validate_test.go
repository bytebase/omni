package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestValidateQuery pins the Cosmos DB editor gate: the query language has no
// DML or DDL, so every statement the parser accepts is read-only and returns
// data, and a syntax error is the only refusal. The corpus includes the shapes
// the package's parse and query-span tests already carry.
func TestValidateQuery(t *testing.T) {
	tests := []struct {
		name        string
		statement   string
		readOnly    bool
		returnsData bool
		wantErr     bool
	}{
		{name: "single select", statement: "SELECT * FROM c", readOnly: true, returnsData: true},
		{name: "select with where", statement: "SELECT c.id, c.name FROM users c WHERE c.age > 18", readOnly: true, returnsData: true},
		{name: "select distinct", statement: "SELECT DISTINCT c.category FROM products c", readOnly: true, returnsData: true},
		{name: "multi-line select", statement: "SELECT\n  * FROM c", readOnly: true, returnsData: true},
		{name: "empty statement", statement: "", readOnly: true, returnsData: true},
		{name: "whitespace only", statement: "   \n  \t  ", readOnly: true, returnsData: true},
		{name: "invalid syntax", statement: "SELCT * FRM c", wantErr: true},
		{name: "unclosed string", statement: "SELECT * FROM c WHERE c.name = 'test", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			readOnly, returnsData, err := ValidateQuery(tc.statement)
			if tc.wantErr {
				require.Error(t, err)
				require.False(t, readOnly)
				require.False(t, returnsData)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.readOnly, readOnly)
			require.Equal(t, tc.returnsData, returnsData)
		})
	}
}
