package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRewriteResultLimit(t *testing.T) {
	testCases := []struct {
		stmt  string
		count int
		want  string
	}{
		{
			stmt:  "SELECT * FROM t;",
			count: 10,
			want:  "SELECT * FROM `t` LIMIT 10",
		},
		{
			stmt:  "WITH a AS (SELECT * FROM t) SELECT * FROM a;",
			count: 10,
			want:  "WITH `a` AS (SELECT * FROM `t`) SELECT * FROM `a` LIMIT 10",
		},
		{
			stmt:  "SELECT * FROM t1 UNION SELECT * FROM t2;",
			count: 10,
			want:  "SELECT * FROM `t1` UNION SELECT * FROM `t2` LIMIT 10",
		},
		{
			stmt:  "SELECT * FROM t1 INTERSECT SELECT * FROM t2;",
			count: 10,
			want:  "SELECT * FROM `t1` INTERSECT SELECT * FROM `t2` LIMIT 10",
		},
		{
			stmt:  "SELECT * FROM t LIMIT 5;",
			count: 10,
			want:  "SELECT * FROM `t` LIMIT 5",
		},
		{
			stmt:  "SELECT * FROM t LIMIT 123;",
			count: 10,
			want:  "SELECT * FROM `t` LIMIT 10",
		},
		{
			stmt:  "SELECT * FROM t WHERE nickname = 'bb'",
			count: 10,
			want:  "SELECT * FROM `t` WHERE `nickname`='bb' LIMIT 10",
		},
		{
			stmt:  "SELECT 'Hello world'",
			count: 10,
			want:  "SELECT 'Hello world' LIMIT 10",
		},
	}

	for _, tc := range testCases {
		got, err := RewriteResultLimit(tc.stmt, tc.count)
		require.NoError(t, err, tc.stmt)
		require.Equal(t, tc.want, got, tc.stmt)
	}
}

func TestRewriteResultLimitErrors(t *testing.T) {
	for _, stmt := range []string{"SELEC 1", "SELECT 1; SELECT 2", ""} {
		_, err := RewriteResultLimit(stmt, 10)
		require.Error(t, err, stmt)
	}
	got, err := RewriteResultLimit("UPDATE t SET a = 1", 10)
	require.NoError(t, err)
	require.Equal(t, "UPDATE t SET a = 1", got)
}
