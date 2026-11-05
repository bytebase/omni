package parser

import "testing"

// SET TRANSACTION as TiDB v8.5.5 reads it, every verdict measured with PREPARE on
// pingcap/tidb:v8.5.5 on 2026-11-05: a 1064 is the parser's refusal, and any other
// answer means the text parsed.
//
// The rules are SetStmt, parser.y:10758-10770, and TransactionChars and TransactionChar,
// parser.y:10840 and 10859. GLOBAL and SESSION are the only scopes it names. The list
// may not be empty, each characteristic has to be whole, and one may repeat. READ ONLY
// alone takes AsOfClause, parser.y:6559, which START TRANSACTION also reads.
//
// Not in this table, because omni still disagrees and the change is not about it:
// TRANSACTION and ISOLATION are unreserved upstream and so variable names, as in `SET
// TRANSACTION = 1`, `SET LOCAL TRANSACTION = 1`, `SET transaction.a = 1`, `SET a = 1,
// transaction = 2` and `SET isolation = 1`. omni refuses them because
// parseLvalueIdent takes neither word.

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
	// The clause on SET TRANSACTION.
	{"set transaction read only as of timestamp now()", true},
	{"set transaction read only as of timestamp '2021-04-21 00:42:12'", true},
	{"SET TRANSACTION READ ONLY AS OF TIMESTAMP NOW()", true},
	{"set session transaction read only as of timestamp now()", true},
	{"set global transaction read only as of timestamp now()", true},
	{"set local transaction read only as of timestamp now()", false},
	{"set persist transaction read only as of timestamp now()", false},
	{"set persist_only transaction read only as of timestamp now()", false},
	{"set @@session transaction read only as of timestamp now()", false},
	{"set transaction read only as of timestamp now() - interval 3 second", true},
	{"set transaction read only as of timestamp tidb_bounded_staleness(now(), now())", true},
	{"set transaction read only as of timestamp (select now())", true},
	{"set transaction read only as of timestamp a", true},
	{"set transaction read only as of timestamp @a", true},
	{"set transaction read only as of timestamp 1 or 2", true},
	{"set transaction read only as /* c */ of timestamp now()", true},
	{"set transaction read only as\nof\ntimestamp now()", true},
	{"set transaction read only as of", false},
	{"set transaction read only as of timestamp", false},
	{"set transaction read only as of now()", false},
	{"set transaction read only as of time now()", false},
	{"set transaction read only as now()", false},
	{"set transaction read only as timestamp now()", false},
	{"set transaction read only of timestamp now()", false},
	{"set transaction read only timestamp now()", false},
	{"set transaction read only as `of` timestamp now()", false},
	{"set transaction read only as of `timestamp` now()", false},
	{"set transaction read only as of timestamp now() a", false},
	{"set transaction read only as of timestamp now() as of timestamp now()", false},
	{"set transaction read only as of timestamp now() read write", false},
	{"set transaction read only as of timestamp now() for update", false},
	{"set transaction read write as of timestamp now()", false},
	{"set transaction as of timestamp now()", false},
	{"set transaction isolation level serializable as of timestamp now()", false},
	{"set transaction read only, as of timestamp now()", false},

	// The clause in a list of characteristics.
	{"set transaction read only as of timestamp now(), isolation level serializable", true},
	{"set transaction isolation level serializable, read only as of timestamp now()", true},
	{"set transaction read only as of timestamp now(), read write", true},
	{"set transaction read write, read only as of timestamp now()", true},
	{"set transaction read only as of timestamp 1, read only as of timestamp 2", true},
	{"set transaction read only as of timestamp now(),", false},
	{"set transaction read only as of timestamp 1, 2", false},
	{"set session transaction isolation level read committed, read only as of timestamp now()", true},
	{"set global transaction read only as of timestamp now(), isolation level repeatable read", true},

	// The rest of the characteristics.
	{"set transaction", false},
	{"set session transaction", false},
	{"set global transaction", false},
	{"set transaction read", false},
	{"set transaction read only", true},
	{"set transaction read write", true},
	{"set transaction read committed", false},
	{"set transaction only", false},
	{"set transaction write", false},
	{"set transaction isolation", false},
	{"set transaction isolation level", false},
	{"set transaction level serializable", false},
	{"set transaction isolation serializable", false},
	{"set transaction isolation level serializable", true},
	{"set transaction isolation level repeatable read", true},
	{"set transaction isolation level repeatable", false},
	{"set transaction isolation level read committed", true},
	{"set transaction isolation level read uncommitted", true},
	{"set transaction isolation level read", false},
	{"set transaction isolation level read only", false},
	{"set transaction isolation level committed", false},
	{"set transaction isolation level snapshot", false},
	{"set transaction isolation level `serializable`", false},
	{"set transaction isolation level serializable, read only", true},
	{"set transaction read only, isolation level serializable", true},
	{"set transaction read only, read write", true},
	{"set transaction isolation level serializable, isolation level read committed", true},
	{"set transaction read only,", false},
	{"set transaction, read only", false},
	{"set transaction read only read write", false},
	{"set transaction isolation level serializable read only", false},
	{"set transaction read only foo", false},
	{"set local transaction read only", false},
	{"set persist transaction read only", false},
	{"set persist_only transaction read only", false},
	{"set session transaction isolation level serializable", true},
	{"set global transaction isolation level read committed", true},
	{"set local transaction isolation level serializable", false},
	{"set @@global transaction read only", false},
	{"set @@transaction = 1", true},
	{"set transaction_isolation = 'READ-COMMITTED'", true},

	// TRANSACTION as a variable name.
	{"set persist transaction = 1", false},
	{"set a = 1, transaction read only", false},
	{"set transaction read only, a = 1", false},
	{"set transaction read only, @a = 1", false},
	{"set transaction read only, transaction read write", false},
	{"set transaction isolation level serializable = 1", false},
	{"set transaction read only as of timestamp a = 1", true},

	// START TRANSACTION, which shares the rule.
	{"start transaction read only as of timestamp now()", true},
	{"start transaction read only as of", false},
	{"start transaction read only as of timestamp", false},
	{"start transaction read only as now()", false},
	{"start transaction read only as timestamp now()", false},
	{"start transaction read only as `of` timestamp now()", false},
	{"start transaction read only of timestamp now()", false},
	{"start transaction read only timestamp now()", false},
	{"start transaction read only as of timestamp now(), read write", false},
}
