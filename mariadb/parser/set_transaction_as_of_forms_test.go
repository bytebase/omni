package parser

import "testing"

// MariaDB has no AS OF TIMESTAMP on SET TRANSACTION. Every verdict was measured with PREPARE
// on mariadb:11.8.8 on 2026-11-05: a 1064 is the parser's refusal, and every form of the clause is
// one. These rows hold only the clause; setTransactionCases holds the rest of the
// statement.

// TestSetTransactionAsOfForms holds omni to the recorded verdicts without a server.
func TestSetTransactionAsOfForms(t *testing.T) {
	for _, c := range setTransactionAsOfCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

var setTransactionAsOfCases = []asOfCase{
	// The clause on SET TRANSACTION.
	{"set transaction read only as of timestamp now()", false},
	{"set transaction read only as of timestamp '2021-04-21 00:42:12'", false},
	{"SET TRANSACTION READ ONLY AS OF TIMESTAMP NOW()", false},
	{"set session transaction read only as of timestamp now()", false},
	{"set global transaction read only as of timestamp now()", false},
	{"set local transaction read only as of timestamp now()", false},
	{"set persist transaction read only as of timestamp now()", false},
	{"set persist_only transaction read only as of timestamp now()", false},
	{"set @@session transaction read only as of timestamp now()", false},
	{"set transaction read only as of timestamp now() - interval 3 second", false},
	{"set transaction read only as of timestamp tidb_bounded_staleness(now(), now())", false},
	{"set transaction read only as of timestamp (select now())", false},
	{"set transaction read only as of timestamp a", false},
	{"set transaction read only as of timestamp @a", false},
	{"set transaction read only as of timestamp 1 or 2", false},
	{"set transaction read only as /* c */ of timestamp now()", false},
	{"set transaction read only as\nof\ntimestamp now()", false},
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
	{"set transaction read only as of timestamp now(), isolation level serializable", false},
	{"set transaction isolation level serializable, read only as of timestamp now()", false},
	{"set transaction read only as of timestamp now(), read write", false},
	{"set transaction read write, read only as of timestamp now()", false},
	{"set transaction read only as of timestamp 1, read only as of timestamp 2", false},
	{"set transaction read only as of timestamp now(),", false},
	{"set transaction read only as of timestamp 1, 2", false},
	{"set session transaction isolation level read committed, read only as of timestamp now()", false},
	{"set global transaction read only as of timestamp now(), isolation level repeatable read", false},
}
