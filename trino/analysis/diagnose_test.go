package analysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/trino/parser"
)

// TestDiagnose pins the editor's syntax diagnostics, including the stubs that
// are valid Trino and must not be reported. The corpus is the adapter's, which
// recorded each case against the parser.
func TestDiagnose(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		expectErr bool
	}{
		{name: "Valid SELECT statement", sql: "SELECT * FROM users"},
		{name: "Valid CREATE TABLE statement", sql: "CREATE TABLE users (id INT, name VARCHAR)"},
		{name: "Missing FROM clause", sql: "SELECT * FFROM users", expectErr: true},
		{name: "Invalid SQL statement", sql: "CREATE TABLEE users (id INT)", expectErr: true},
		{name: "Syntax error", sql: "SELECT FROM users", expectErr: true},
		// Empty input is zero statements, not a syntax error.
		{name: "Empty statement", sql: ""},
		{name: "DESCRIBE INPUT", sql: "DESCRIBE INPUT x"},
		{name: "DESCRIBE OUTPUT", sql: "DESCRIBE OUTPUT x"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagnostics := Diagnose(tt.sql)
			if tt.expectErr {
				assert.NotEmpty(t, diagnostics, "Expected diagnostics for invalid SQL")
			} else {
				assert.Empty(t, diagnostics, "Expected no diagnostics for valid SQL")
			}
		})
	}
}

// TestIsValidButUnimplementedStub pins which parser stubs are valid Trino. A
// CREATE/DROP/ALTER stub fires on an unrecognized object keyword, which is a
// genuine syntax error and stays.
func TestIsValidButUnimplementedStub(t *testing.T) {
	require.True(t, isValidButUnimplementedStub("DESCRIBE INPUT statement parsing is not yet supported"))
	require.True(t, isValidButUnimplementedStub("DESCRIBE OUTPUT statement parsing is not yet supported"))
	require.False(t, isValidButUnimplementedStub("CREATE statement parsing is not yet supported"))
	require.False(t, isValidButUnimplementedStub("DESCRIBE table"))
}

// TestDiagnoseKeepsGenuineErrors proves the filter does not swallow a real
// parser error: an unrecognized CREATE object keyword remains a diagnostic.
func TestDiagnoseKeepsGenuineErrors(t *testing.T) {
	for _, statement := range []string{"CREATE TABLEE users (id INT)", "DROP TABL users"} {
		diags := parser.Diagnose(statement)
		require.NotEmpty(t, diags, "%q", statement)
		require.Equal(t, diags, Diagnose(statement), "%q", statement)
	}
}
