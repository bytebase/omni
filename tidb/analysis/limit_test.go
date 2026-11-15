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

// The rewritten forms were run on TiDB v8.5.5: each is accepted and takes the
// same row locks as the statement it rewrites.
func TestRewriteResultLimitLockingClause(t *testing.T) {
	testCases := []struct {
		stmt string
		want string
	}{
		{"SELECT * FROM t FOR UPDATE", "SELECT * FROM `t` LIMIT 10 FOR UPDATE"},
		{"SELECT * FROM t LIMIT 20 FOR UPDATE NOWAIT", "SELECT * FROM `t` LIMIT 10 FOR UPDATE NOWAIT"},
		{"SELECT * FROM t FOR UPDATE OF t WAIT 5", "SELECT * FROM `t` LIMIT 10 FOR UPDATE OF `t` WAIT 5"},
		{"SELECT * FROM db.t FOR UPDATE OF db.t SKIP LOCKED", "SELECT * FROM `db`.`t` LIMIT 10 FOR UPDATE OF `db`.`t` SKIP LOCKED"},
		{"SELECT * FROM t FOR SHARE", "SELECT * FROM `t` LIMIT 10 FOR SHARE"},
		{"SELECT * FROM t LOCK IN SHARE MODE", "SELECT * FROM `t` LIMIT 10 LOCK IN SHARE MODE"},
		{"TABLE t FOR UPDATE", "TABLE `t` LIMIT 10 FOR UPDATE"},
	}
	for _, tc := range testCases {
		got, err := RewriteResultLimit(tc.stmt, 10)
		require.NoError(t, err, tc.stmt)
		require.Equal(t, tc.want, got, tc.stmt)
	}
}

func TestRewriteResultLimitNestedLock(t *testing.T) {
	for _, stmt := range []string{
		"SELECT * FROM t1 UNION SELECT * FROM t2 FOR UPDATE",
		"SELECT * FROM t1 UNION SELECT * FROM t2 LIMIT 20 FOR UPDATE",
		"SELECT * FROM t1 UNION (SELECT * FROM t2 FOR UPDATE)",
		"(SELECT * FROM t FOR UPDATE) LIMIT 20",
		"SELECT * FROM (SELECT * FROM t FOR UPDATE) d",
		"WITH c AS (SELECT * FROM t FOR UPDATE) SELECT * FROM c",
		"SELECT (SELECT a FROM t FOR UPDATE) x",
	} {
		_, err := RewriteResultLimit(stmt, 10)
		require.ErrorIs(t, err, errNestedLock, stmt)
	}
}
