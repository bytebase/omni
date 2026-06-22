package parser

import "testing"

// TestUUIDInetTypeAccept covers the MariaDB-only scalar types UUID, INET4 and
// INET6 (BYT-9135). They are non-reserved — UUID is also a function, INET4/INET6
// are usable as identifiers — so they are recognised only in type position.
func TestUUIDInetTypeAccept(t *testing.T) {
	accept := []string{
		"CREATE TABLE t (a UUID)",
		"CREATE TABLE t (a uuid)", // case-insensitive
		"CREATE TABLE t (a INET4, b INET6)",
		"CREATE TABLE t (id INT PRIMARY KEY, addr INET6 NOT NULL DEFAULT '::1')",
		"CREATE TABLE t (u UUID DEFAULT UUID())",
		"ALTER TABLE t ADD COLUMN addr INET4",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestUUIDInetNonReserved pins that the words stay usable as function names and
// identifiers (the contextual recognition must not reserve them).
func TestUUIDInetNonReserved(t *testing.T) {
	accept := []string{
		"SELECT UUID()",
		"SELECT uuid, inet4, inet6 FROM t",
		"CREATE TABLE t (uuid INT, inet4 VARCHAR(10), inet6 TEXT)",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}
