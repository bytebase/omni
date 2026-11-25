package catalog

import (
	"errors"
	"testing"
)

// TestSequenceRefusals pins the code and SQLSTATE of each kind of sequence
// refusal as OceanBase 4.3.5.6 answers it (testdata/tapes/
// TestSequenceSimulation.tape), which TestSequenceSimulation compares by code
// alone.
func TestSequenceRefusals(t *testing.T) {
	for _, tc := range []struct {
		sql   string
		code  int
		state string
	}{
		{"CREATE SEQUENCE a", 4297, "HY000"},
		{"DROP SEQUENCE b", 4298, "HY000"},
		{"ALTER SEQUENCE a", 4314, "HY000"},
		{"CREATE SEQUENCE b CYCLE CYCLE", 4303, "HY000"},
		{"CREATE SEQUENCE b CYCLE NOCYCLE", 4308, "HY000"},
		{"CREATE SEQUENCE b RESTART", 1149, "42000"},
		{"CREATE SEQUENCE b MAXVALUE 1e5", 5114, "42000"},
		{"CREATE SEQUENCE b MAXVALUE 1.5", 4317, "HY000"},
		{"CREATE SEQUENCE b CYCLE", 4330, "HY000"},
		{"ALTER SEQUENCE a START WITH 5", 4311, "HY000"},
		{"CREATE SEQUENCE nope.b", 1049, "42000"},
		{"CREATE SEQUENCE bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", 1059, "42000"},
	} {
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE SEQUENCE a;")
		if err != nil {
			t.Fatal(err)
		}
		results, err := c.Exec(tc.sql, nil)
		if err != nil {
			t.Fatalf("%s: %v", tc.sql, err)
		}
		var refusal *Error
		if !errors.As(results[0].Error, &refusal) || refusal.Code != tc.code || refusal.SQLState != tc.state {
			t.Errorf("%s: %v, want %d (%s)", tc.sql, results[0].Error, tc.code, tc.state)
		}
	}
}

// TestSequenceStored pins what ALTER SEQUENCE stores where OceanBase keeps a
// bound past 64 bits (NOMAXVALUE) and where RESTART moves the start.
func TestSequenceStored(t *testing.T) {
	c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE SEQUENCE a START WITH 10 INCREMENT BY -2 MINVALUE -100 MAXVALUE 100 CYCLE CACHE 5 ORDER; ALTER SEQUENCE a NOMAXVALUE; ALTER SEQUENCE a RESTART;")
	if err != nil {
		t.Fatal(err)
	}
	seq := c.GetDatabase("d").Sequences["a"]
	if seq == nil || seq.Max.String() != "9223372036854775808" || seq.Start.String() != "9223372036854775808" {
		t.Fatalf("got %+v, want MAXVALUE and START WITH 9223372036854775808", seq)
	}
}
