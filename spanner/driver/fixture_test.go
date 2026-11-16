package driver

import (
	"testing"

	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

// fixtureDDL is the schema the sync and query tests read. It covers every
// Spanner type a result cell can hold, a column that is NOT NULL, one that is
// nullable, one with a default, a table's primary key, a unique index, an index
// whose key descends, a foreign key, and a view.
var fixtureDDL = []string{
	`CREATE TABLE Singers (
		Id INT64 NOT NULL,
		FirstName STRING(1024),
		LastName STRING(1024) NOT NULL,
		Active BOOL,
		Rating FLOAT64,
		Photo BYTES(MAX),
		Born DATE,
		UpdatedAt TIMESTAMP,
		Meta JSON,
		Scores ARRAY<INT64>,
	) PRIMARY KEY (Id)`,
	`CREATE TABLE Albums (
		Id INT64 NOT NULL,
		SingerId INT64 NOT NULL,
		Title STRING(1024) NOT NULL DEFAULT ('untitled'),
		CONSTRAINT FK_Singer FOREIGN KEY (SingerId) REFERENCES Singers (Id),
	) PRIMARY KEY (Id)`,
	`CREATE UNIQUE INDEX SingersByFirstName ON Singers (FirstName)`,
	`CREATE INDEX SingersByLastName ON Singers (LastName DESC)`,
	`CREATE VIEW SingerNames SQL SECURITY INVOKER AS SELECT s.Id AS Id, s.LastName AS LastName FROM Singers s`,
}

// fixtureRows is the fixture's data: one singer carrying every type, one with
// most of them left out, one whose name sorts between theirs, and an album per
// singer with a name.
var fixtureRows = []string{
	`INSERT INTO Singers (Id, FirstName, LastName, Active, Rating, Photo, Born, UpdatedAt, Meta, Scores) VALUES
		(1, 'Ann', 'Smith', TRUE, 4.5, b'\x01\x02\x03', DATE '2024-01-15', TIMESTAMP '2024-01-15T10:30:45.123456Z', JSON '{"a":1,"b":[2,3]}', [1, 2, 3])`,
	`INSERT INTO Singers (Id, FirstName, LastName) VALUES (2, 'Bob', 'Jones')`,
	`INSERT INTO Singers (Id, FirstName, LastName, Active, Rating) VALUES (3, 'Cid', 'Adams', FALSE, 3.25)`,
	`INSERT INTO Albums (Id, SingerId, Title) VALUES (10, 1, 'First')`,
	`INSERT INTO Albums (Id, SingerId) VALUES (11, 2)`,
}

// syncFixture returns the instance and database the sync tests read.
func syncFixture(t testing.TB) (string, string) {
	t.Helper()
	return createFixture(t, "sync")
}

// queryFixture returns the instance and database the query tests read. It is a
// database of its own, so a query test that writes cannot change what a sync
// test recorded.
func queryFixture(t testing.TB) (string, string) {
	t.Helper()
	return createFixture(t, "query")
}

func createFixture(t testing.TB, id string) (string, string) {
	t.Helper()
	testcontainer.Create(t, id, id, fixtureDDL, fixtureRows)
	return id, id
}

// openFixture returns a driver over a fixture database named id.
func openFixture(t testing.TB, id string) *Driver {
	t.Helper()
	testcontainer.Create(t, id, id, fixtureDDL, fixtureRows)
	return openDriver(t, id, id)
}
