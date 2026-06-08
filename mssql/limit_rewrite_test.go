package mssql_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/mssql"
)

func TestStatementWithResultLimit(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		limit     int
		want      string
	}{
		{
			name:      "adds top after select",
			statement: "SELECT id FROM t1;",
			limit:     10,
			want:      "SELECT TOP 10 id FROM t1;",
		},
		{
			name:      "adds top after distinct",
			statement: "SELECT DISTINCT S.*, DS.DistrictId FROM Streets S LEFT JOIN DistrictStreets DS ON S.StreetId = DS.StreetId;",
			limit:     10,
			want:      "SELECT DISTINCT TOP 10 S.*, DS.DistrictId FROM Streets S LEFT JOIN DistrictStreets DS ON S.StreetId = DS.StreetId;",
		},
		{
			name:      "keeps smaller top",
			statement: "SELECT TOP 3 id FROM t1;",
			limit:     10,
			want:      "SELECT TOP 3 id FROM t1;",
		},
		{
			name:      "lowers larger top",
			statement: "SELECT TOP 123 id FROM t1;",
			limit:     10,
			want:      "SELECT TOP 10 id FROM t1;",
		},
		{
			name:      "lowers top percent even when numeric value is smaller",
			statement: "SELECT TOP (3) PERCENT id FROM t1;",
			limit:     10,
			want:      "SELECT TOP 10 id FROM t1;",
		},
		{
			name:      "lowers top with ties even when numeric value is smaller",
			statement: "SELECT TOP (3) WITH TIES id FROM t1 ORDER BY id;",
			limit:     10,
			want:      "SELECT TOP 10 id FROM t1 ORDER BY id;",
		},
		{
			name:      "adds offset fetch after order by",
			statement: "SELECT id FROM t1 ORDER BY price;",
			limit:     10,
			want:      "SELECT id FROM t1 ORDER BY price OFFSET 0 ROWS FETCH NEXT 10 ROWS ONLY;",
		},
		{
			name:      "adds fetch after existing offset",
			statement: "SELECT id FROM t1 ORDER BY price OFFSET 123 ROWS;",
			limit:     10,
			want:      "SELECT id FROM t1 ORDER BY price OFFSET 123 ROWS FETCH NEXT 10 ROWS ONLY;",
		},
		{
			name:      "lowers existing fetch",
			statement: "SELECT name, price FROM toy ORDER BY price OFFSET 0 ROWS FETCH FIRST 123 ROWS ONLY;",
			limit:     10,
			want:      "SELECT name, price FROM toy ORDER BY price OFFSET 0 ROWS FETCH FIRST 10 ROWS ONLY;",
		},
		{
			name: "adds top to both union branches",
			statement: `SELECT
    first_name,
    last_name
FROM
    sales.staffs
UNION ALL
SELECT
    first_name,
    last_name
FROM
    sales.customers;`,
			limit: 10,
			want: `SELECT TOP 10
    first_name,
    last_name
FROM
    sales.staffs
UNION ALL
SELECT TOP 10
    first_name,
    last_name
FROM
    sales.customers;`,
		},
		{
			name: "adds fetch after union order by",
			statement: `SELECT
    first_name,
    last_name
FROM
    sales.staffs
UNION ALL
SELECT
    first_name,
    last_name
FROM
    sales.customers
ORDER BY
    first_name,
    last_name;`,
			limit: 10,
			want: `SELECT
    first_name,
    last_name
FROM
    sales.staffs
UNION ALL
SELECT
    first_name,
    last_name
FROM
    sales.customers
ORDER BY
    first_name,
    last_name OFFSET 0 ROWS FETCH NEXT 10 ROWS ONLY;`,
		},
		{
			name: "adds top to cte outer select",
			statement: `WITH cte_org AS (
SELECT
    staff_id,
    first_name,
    manager_id
FROM
    sales.staffs
)
SELECT * FROM cte_org;`,
			limit: 10,
			want: `WITH cte_org AS (
SELECT
    staff_id,
    first_name,
    manager_id
FROM
    sales.staffs
)
SELECT TOP 10 * FROM cte_org;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mssql.StatementWithResultLimit(tt.statement, tt.limit)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestStatementWithResultLimitErrors(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		limit     int
	}{
		{
			name:      "empty statement",
			statement: "",
			limit:     10,
		},
		{
			name:      "multi statement",
			statement: "SELECT 1; SELECT 2;",
			limit:     10,
		},
		{
			name:      "invalid limit",
			statement: "SELECT 1",
			limit:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mssql.StatementWithResultLimit(tt.statement, tt.limit)
			require.Error(t, err)
		})
	}
}
