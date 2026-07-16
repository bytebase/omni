package splittest

// KnownBetterThanPsql lists inputs where this splitter intentionally
// diverges from psql's client-side scanner (psqlscan): the server grammar
// accepts each shape as a single statement — engine-verified via psql -c
// on PostgreSQL 17 — but psql's script scanner mis-splits it inside
// BEGIN ATOMIC bodies. The S3 server-differential harness must treat a
// psql-vs-omni mismatch on these inputs as expected (omni matches the
// server, which is the higher truth anchor); without this allowance S3
// would report them as false mismatches.
var KnownBetterThanPsql = []string{
	// Dot-qualified reserved word inside a SQL-standard function body:
	// psqlscan closes the block at ".end" and truncates the statement.
	"CREATE FUNCTION f6() RETURNS int BEGIN ATOMIC SELECT t4.end FROM t4; END;",
	// AS alias using a reserved word: same psqlscan truncation.
	"CREATE FUNCTION f7() RETURNS int BEGIN ATOMIC SELECT 1 AS end; END;",
}
