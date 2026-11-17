package analysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnoseUsesParseErrorPosition(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		line      int
	}{
		{
			name:      "single line",
			statement: "SELECT * FROM",
			line:      0,
		},
		{
			name:      "second sql statement",
			statement: "SELECT 1;\n\nSELECT * FROM",
			line:      2,
		},
		{
			name:      "line after non-ascii",
			statement: "SELECT '你好';\nSELECT * FROM",
			line:      1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diagnostic := Diagnose(tc.statement)
			require.NotNil(t, diagnostic)
			require.NotEmpty(t, diagnostic.Message)
			require.Equal(t, tc.line, strings.Count(tc.statement[:diagnostic.Offset], "\n"))
		})
	}
}

func TestDiagnoseAcceptsRedshiftSyntax(t *testing.T) {
	tests := []string{
		"SELECT * FROM sales QUALIFY row_number() OVER () = 1;",
		"COPY sales FROM 's3://bucket/path' IAM_ROLE 'arn:aws:iam::123456789012:role/redshift';",
		"UNLOAD ('SELECT * FROM sales') TO 's3://bucket/out' IAM_ROLE 'arn:aws:iam::123456789012:role/redshift';",
	}

	for _, statement := range tests {
		t.Run(statement, func(t *testing.T) {
			require.Nil(t, Diagnose(statement))
		})
	}
}
