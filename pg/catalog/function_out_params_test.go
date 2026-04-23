package catalog

import "testing"

func TestCreateFunctionInfersReturnTypeFromSingleOutParam(t *testing.T) {
	c := New()

	execSQL(t, c, `
CREATE FUNCTION f_out_scalar(IN input_value integer, OUT output_value text)
LANGUAGE plpgsql
AS 'BEGIN output_value := input_value::text; END';
`)

	procs := c.LookupProcByName("f_out_scalar")
	if len(procs) != 1 {
		t.Fatalf("procs: got %d, want 1", len(procs))
	}
	up := c.userProcs[procs[0].OID]
	if up == nil {
		t.Fatal("user proc not found")
	}
	if got, want := string(up.ArgModes), "io"; got != want {
		t.Fatalf("arg modes: got %q, want %q", got, want)
	}
	if got, want := len(up.AllArgTypes), 2; got != want {
		t.Fatalf("all arg types count: got %d, want %d", got, want)
	}
	if procs[0].RetType != TEXTOID {
		t.Fatalf("ret type: got %d, want %d", procs[0].RetType, TEXTOID)
	}
	if procs[0].NArgs != 1 {
		t.Fatalf("nargs: got %d, want 1", procs[0].NArgs)
	}
}

func TestCreateFunctionInfersRecordReturnTypeFromMultipleOutParams(t *testing.T) {
	c := New()

	execSQL(t, c, `
CREATE FUNCTION f_out_record(IN input_value integer, OUT doubled integer, OUT label text)
LANGUAGE plpgsql
AS 'BEGIN doubled := input_value * 2; label := input_value::text; END';
`)

	procs := c.LookupProcByName("f_out_record")
	if len(procs) != 1 {
		t.Fatalf("procs: got %d, want 1", len(procs))
	}
	if procs[0].RetType != RECORDOID {
		t.Fatalf("ret type: got %d, want %d", procs[0].RetType, RECORDOID)
	}
	if procs[0].NArgs != 1 {
		t.Fatalf("nargs: got %d, want 1", procs[0].NArgs)
	}
}

func TestCreateFunctionWithoutReturnsAndWithoutOutParamsIsRejected(t *testing.T) {
	c := New()

	err := execSQLErr(c, `
CREATE FUNCTION f_missing_returns(input_value integer)
LANGUAGE sql
AS 'SELECT input_value';
`)
	assertCode(t, err, CodeInvalidFunctionDefinition)
}
