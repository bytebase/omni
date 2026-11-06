package catalog

import "testing"

// TestSetForeignKeyChecksUnaryPlus: the parser keeps a unary plus, so SET
// reads its value through it. On TiDB v8.5.5, SET foreign_key_checks = +0
// turns the checks off and +1 turns them back on.
func TestSetForeignKeyChecksUnaryPlus(t *testing.T) {
	c := New()
	results, err := c.Exec(`CREATE DATABASE d; USE d;
SET foreign_key_checks = +0;
CREATE TABLE c1 (a INT, FOREIGN KEY (a) REFERENCES nosuch (id));
SET foreign_key_checks = +1;
CREATE TABLE c2 (a INT, FOREIGN KEY (a) REFERENCES nosuch (id));`, &ExecOptions{ContinueOnError: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 6 {
		t.Fatalf("got %d results, want 6", len(results))
	}
	if err := results[3].Error; err != nil {
		t.Errorf("CREATE TABLE with checks off (+0): %v", err)
	}
	if results[5].Error == nil {
		t.Error("CREATE TABLE with checks back on (+1): want an error for the missing table")
	}
}
