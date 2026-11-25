// Package testcontainer is how spanner's tests reach the pinned Cloud Spanner
// emulator: through a tape (enginetest/tape) that replays its answers, or, in
// a recording (OMNI_RECORD=1), through the emulator itself. The emulator is one
// testcontainer per test binary, started by the first test that records and
// removed by Terminate, which each package's TestMain calls. A replay starts
// nothing and needs no Docker.
//
// The grammar's tests (spanner/parser, spanner/analysis) ask through
// spanner/spanner-oracle, a harness that sends one statement to the emulator
// and answers a JSON verdict for it (spanner/spanner-oracle/README.md); see
// harness.go. Their tape is the harness's: each request is one statement, and
// each answer the harness's JSON for it. The harness is a program of its own
// module, so the tape is the boundary the spec gives such programs: replay
// builds and runs nothing.
//
// The driver's tests (spanner/driver) ask through Spanner's own clients; see
// driver.go. Spanner's client is not a database/sql driver, so the tape cannot
// sit at that seam. The driver's own tests answer the client seam
// (spanner/driver's client) instead: a request is one call the driver makes of
// its clients — the admin API's listing, get, DDL update, creation and DDL
// read, a single-use read, a snapshot's query, an update, a batch update, a
// partitioned update, a plan ask, a ping and a close — and an answer is the
// databases it listed, the DDL it read, the rows a read returned or the
// client's error. Those tests own the exchange's spelling because they own the
// seam's types.
//
// The emulator answers the admin API, DDL, DML, queries and INFORMATION_SCHEMA,
// so sync, query, execute, dump and ping are all recorded from it. Three facts a
// real instance has and an emulator does not are never recorded, because it
// cannot answer them: IAM and the credential chain (it accepts any caller), a
// Private Service Connect endpoint and multi-region routing (it is one process
// on one address), and query plans (it has no optimizer). The driver's package
// doc states each.
//
// Create creates a GoogleSQL database from a backtick-quoted CREATE DATABASE.
// CreatePostgres creates a PostgreSQL-dialect one, which the driver skips and
// refuses; it is there so a test records the database beside a GoogleSQL one.
//
// A database is created per fixture, through the admin API, before the first
// request of a recording: a replay answers as if it were there, and the fixture
// is not part of what the driver asked. The database's id is the caller's, so
// two tests hold their own state and two recordings of one test write the same
// resource names.
package testcontainer

import (
	"path/filepath"
	"runtime"
)

// Pins is spanner/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// The pin roles the tests ask: Default is the emulator the grammar's harness
// tapes record from, and Driver the one the driver's tapes record from.
const (
	Default = "default"
	Driver  = "driver"
)

// Terminate removes the emulators this run started and the harness binary it
// built, if any. Each package's TestMain calls it.
func Terminate() {
	terminateHarness()
	terminateDriver()
}
