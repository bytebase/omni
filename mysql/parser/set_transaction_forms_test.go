package parser

import "testing"

// SET TRANSACTION as MySQL reads it, every verdict measured with PREPARE on mysql:8.0
// (8.0.46) on 2026-11-05: a 1064 is the parser's refusal, and any other answer here is a
// statement that parsed (1295 is PREPARE refusing to carry a CREATE PROCEDURE it read).
//
// The rules are start_option_value_list and transaction_characteristics, sql_yacc.yy:15760
// and 16035. The list holds one characteristic or one of each kind, so unlike TiDB it
// refuses `SET TRANSACTION READ ONLY, READ WRITE`. option_type, sql_yacc.yy:15918, has every
// scope, so unlike TiDB it reads LOCAL, PERSIST and PERSIST_ONLY.
//
// Not in this table, because 1193 is no syntax verdict here: the parser raises it in an
// action for an unknown variable and stops there, as `SET TRANSACTION = 1 2` shows, so the
// answer says nothing of the text after it. TRANSACTION is unreserved upstream, so `SET
// TRANSACTION = 1`, `SET SESSION TRANSACTION = 1` and `SET transaction.a = 1` all reach
// that action, and omni refuses them; each is an error on the server either way.

// TestSetTransactionForms holds omni to the recorded verdicts without a server, so a run
// with no Docker still reports a widened or narrowed grammar.
func TestSetTransactionForms(t *testing.T) {
	for _, c := range setTransactionCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

var setTransactionCases = []asOfCase{
	// A list with no characteristic, or one cut short.
	{"set transaction", false},
	{"set session transaction", false},
	{"set global transaction", false},
	{"set local transaction", false},
	{"set persist transaction", false},
	{"set persist_only transaction", false},
	{"set transaction read", false},
	{"set transaction read only,", false},
	{"set transaction read write,", false},
	{"set transaction, read only", false},
	{"set transaction, isolation level serializable", false},
	{"set transaction only", false},
	{"set transaction write", false},
	{"set transaction read committed", false},
	{"set transaction read uncommitted", false},
	{"set transaction repeatable read", false},
	{"set transaction serializable", false},
	{"set transaction isolation", false},
	{"set transaction level serializable", false},
	{"set transaction isolation serializable", false},
	{"set transaction isolation level", false},
	{"set transaction isolation level repeatable", false},
	{"set transaction isolation level repeatable write", false},
	{"set transaction isolation level read", false},
	{"set transaction isolation level read only", false},
	{"set transaction isolation level read write", false},
	{"set transaction isolation level committed", false},
	{"set transaction isolation level uncommitted", false},
	{"set transaction isolation level snapshot", false},
	{"set transaction isolation level default", false},
	{"set transaction isolation level `serializable`", false},
	{"set transaction isolation level 'serializable'", false},
	{"set transaction `read` only", false},
	{"set transaction read `only`", false},
	{"set transaction `isolation` level serializable", false},
	{"set transaction isolation `level` serializable", false},

	// One characteristic.
	{"set transaction read only", true},
	{"set transaction read write", true},
	{"set transaction read only;", true},
	{"SET TRANSACTION READ ONLY", true},
	{"SET TRANSACTION ISOLATION LEVEL READ COMMITTED", true},
	{"set transaction /* c */ read /* c */ only", true},
	{"set transaction isolation -- c\nlevel serializable", true},
	{"set transaction isolation level serializable", true},
	{"set transaction isolation level repeatable read", true},
	{"set transaction isolation level read committed", true},
	{"set transaction isolation level read uncommitted", true},

	// Two characteristics of different kinds, in either order.
	{"set transaction read only, isolation level serializable", true},
	{"set transaction read write, isolation level read committed", true},
	{"set transaction isolation level serializable, read only", true},
	{"set transaction isolation level repeatable read, read write", true},
	{"set transaction isolation level read uncommitted, read write", true},
	{"set transaction isolation level read committed, read only", true},

	// A kind named twice, or a third characteristic.
	{"set transaction read only, read write", false},
	{"set transaction read only, read only", false},
	{"set transaction read write, read only", false},
	{"set transaction read write, read write", false},
	{"set transaction isolation level serializable, isolation level read committed", false},
	{"set transaction isolation level serializable, isolation level serializable", false},
	{"set transaction read only, isolation level serializable, read write", false},
	{"set transaction read only, isolation level serializable, isolation level read committed", false},
	{"set transaction isolation level serializable, read only, isolation level read committed", false},
	{"set transaction isolation level serializable, read only, read write", false},

	// A list cut short after its comma, or with no comma.
	{"set transaction read only, isolation level serializable,", false},
	{"set transaction isolation level serializable, read only,", false},
	{"set transaction isolation level serializable,", false},
	{"set transaction read only, isolation", false},
	{"set transaction read only, isolation level", false},
	{"set transaction read only, isolation level read", false},
	{"set transaction read only, isolation level repeatable", false},
	{"set transaction read only, read", false},
	{"set transaction read only, serializable", false},
	{"set transaction isolation level serializable, read", false},
	{"set transaction isolation level serializable, only", false},
	{"set transaction isolation level serializable, isolation", false},
	{"set transaction read only read write", false},
	{"set transaction read only isolation level serializable", false},
	{"set transaction isolation level serializable read only", false},
	{"set transaction isolation level serializable isolation level read committed", false},

	// Text after the characteristics.
	{"set transaction read only foo", false},
	{"set transaction read write foo", false},
	{"set transaction isolation level serializable foo", false},
	{"set transaction read only, foo", false},
	{"set transaction read only with consistent snapshot", false},
	{"set transaction with consistent snapshot", false},
	{"set transaction read only for update", false},

	// Each scope. All five read the same list.
	{"set session transaction read only", true},
	{"set session transaction read write", true},
	{"set session transaction isolation level serializable", true},
	{"set session transaction isolation level read committed, read only", true},
	{"set session transaction read only, read write", false},
	{"set session transaction read", false},
	{"set global transaction read only", true},
	{"set global transaction read write", true},
	{"set global transaction isolation level serializable", true},
	{"set global transaction read only, isolation level serializable", true},
	{"set global transaction read only, read write", false},
	{"set global transaction isolation", false},
	{"set local transaction read only", true},
	{"set local transaction read write", true},
	{"set local transaction isolation level serializable", true},
	{"set local transaction isolation level read committed, read only", true},
	{"set local transaction read only, read write", false},
	{"set local transaction read", false},
	{"set persist transaction read only", true},
	{"set persist transaction read write", true},
	{"set persist transaction isolation level serializable", true},
	{"set persist transaction read only, isolation level serializable", true},
	{"set persist transaction read only, read write", false},
	{"set persist transaction read", false},
	{"set persist_only transaction read only", true},
	{"set persist_only transaction read write", true},
	{"set persist_only transaction isolation level serializable", true},
	{"set persist_only transaction read only, isolation level serializable", true},
	{"set persist_only transaction read only, read write", false},
	{"set persist_only transaction read", false},
	{"set @@session transaction read only", false},
	{"set @@global transaction read only", false},
	{"set @@local transaction read only", false},
	{"set @@persist transaction read only", false},
	{"set @@persist_only transaction read only", false},
	{"set session session transaction read only", false},
	{"set global session transaction read only", false},
	{"set session global transaction read only", false},
	{"set local local transaction read only", false},
	{"set session.transaction read only", false},
	{"set `session` transaction read only", false},
	{"set `global` transaction read only", false},
	{"set `transaction` read only", false},

	// TRANSACTION beside an assignment, which no list takes.
	{"set @a = 1, transaction read only", false},
	{"set sql_mode = '', transaction read only", false},
	{"set session sql_mode = '', transaction read only", false},
	{"set sql_mode = '', session transaction read only", false},
	{"set transaction read only, @a = 1", false},
	{"set transaction read only, sql_mode = ''", false},
	{"set transaction read only, session sql_mode = ''", false},
	{"set transaction isolation level serializable, @a = 1", false},
	{"set transaction read only, transaction read write", false},
	{"set transaction read only, session transaction read write", false},
	{"set transaction isolation level serializable = 1", false},
	{"set transaction read only = 1", false},
	{"set transaction read = 1", false},

	// The names that start like it.
	{"set @transaction = 1", true},
	{"set transaction_isolation = 'READ-COMMITTED'", true},
	{"set session transaction_isolation = 'READ-COMMITTED'", true},
	{"set global transaction_isolation = 'READ-COMMITTED'", true},
	{"set transaction_read_only = 1", true},
	{"set session transaction_read_only = 1", true},

	// In a stored program, which reads the same rule.
	{"create procedure p() set transaction read only", true},
	{"create procedure p() begin set transaction read only; end", true},
	{"create procedure p() begin set transaction isolation level serializable, read only; end", true},
	{"create procedure p() begin set transaction; end", false},
	{"create procedure p() begin set transaction read; end", false},
	{"create procedure p() begin set transaction read only, read write; end", false},
	{"create procedure p() begin set local transaction read only; end", true},
	{"create procedure p() begin set persist transaction read only; end", true},
}
