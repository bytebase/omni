package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/mariadb/parser"
)

type tableResource struct {
	Database string
	Schema   string
	Table    string
}

func TestCollectAccessTables(t *testing.T) {
	tests := []struct {
		statement string
		expected  []tableResource
	}{
		{
			statement: "SELECT * FROM t1 WHERE c1 = 1;",
			expected: []tableResource{
				{
					Database: "db",
					Table:    "t1",
				},
			},
		},
		{
			statement: "SELECT * FROM db1.t1 JOIN db2.t2 ON t1.c1 = t2.c1;",
			expected: []tableResource{
				{
					Database: "db1",
					Table:    "t1",
				},
				{
					Database: "db2",
					Table:    "t2",
				},
			},
		},
		{
			statement: "SELECT a > (select max(a) from t1) FROM t2;",
			expected: []tableResource{
				{
					Database: "db",
					Table:    "t1",
				},
				{
					Database: "db",
					Table:    "t2",
				},
			},
		},
	}

	for _, test := range tests {
		parseResult, err := parser.Parse(test.statement)
		require.NoError(t, err, "failed to parse statement: %s", test.statement)
		require.Len(t, parseResult.Items, 1, "expected one parse result for statement: %s", test.statement)

		sourceColumns := collectAccessTables(parseResult.Items[0], "db")
		var resources []tableResource
		for resource := range sourceColumns {
			resources = append(resources, tableResource{
				Database: resource.Database,
				Schema:   resource.Schema,
				Table:    resource.Table,
			})
		}
		require.ElementsMatch(t, test.expected, resources, test.statement)
	}
}
