package testcontainer

import (
	"errors"
	"reflect"
	"testing"

	mssql "github.com/microsoft/go-mssqldb"
)

// TestCodecRoundTrips holds Codec to giving back the error it wrote: every
// error the batch reported, first to last, with the last one's fields on
// the error itself, as go-mssqldb gives them.
func TestCodecRoundTrips(t *testing.T) {
	first := mssql.Error{Number: 105, State: 1, Class: 15, LineNo: 1, Message: "Unclosed quotation mark after the character string 'abc'."}
	last := mssql.Error{Number: 2627, State: 1, Class: 14, LineNo: 3, ProcName: "dbo.my proc", Message: "Violation of PRIMARY KEY constraint.\nThe duplicate key value is (1) \\ 2."}
	for _, in := range []mssql.Error{
		{Number: 102, State: 1, Class: 15, LineNo: 1, Message: "Incorrect syntax near 'FROM'.", All: []mssql.Error{{Number: 102, State: 1, Class: 15, LineNo: 1, Message: "Incorrect syntax near 'FROM'."}}},
		{Number: last.Number, State: last.State, Class: last.Class, LineNo: last.LineNo, ProcName: last.ProcName, Message: last.Message, All: []mssql.Error{first, last}},
		{Number: 0, Message: ""},
	} {
		text, ok := Codec{}.EncodeError(in)
		if !ok {
			t.Fatalf("EncodeError(%v) refused a server's error", in)
		}
		var out mssql.Error
		if !errors.As(Codec{}.DecodeError(text), &out) {
			t.Fatalf("DecodeError(%q) is no mssql.Error", text)
		}
		want := in
		if len(want.All) == 0 {
			want.All = []mssql.Error{{Number: in.Number, State: in.State, Class: in.Class, LineNo: in.LineNo, ProcName: in.ProcName, Message: in.Message}}
		}
		if !reflect.DeepEqual(out, want) {
			t.Errorf("DecodeError(EncodeError(%+v)) = %+v\ntext:\n%s", in, out, text)
		}
	}
	if _, ok := (Codec{}).EncodeError(errors.New("connection reset")); ok {
		t.Errorf("EncodeError took an error the server did not send")
	}
}

// TestMadeUpNames pins what Codec gives a stand-in: an object id read as
// text, and a constraint name SQL Server made up, and the stand-ins
// themselves, so a test that sends one back is read again.
func TestMadeUpNames(t *testing.T) {
	for _, name := range []string{
		"oid:581577110",
		"oid:-1366940038",
		"PK__b__3213E83F2F1A1B2C",
		"UQ__t__357D4CF9A1B2C3D4",
		"DF__t__a__5AEE82B9",
		"CK__nokey__a__4F7CD00D",
		"FK__x__a__6383C8BA",
		"UQ__a b__3BD019E4E299938C",
	} {
		if got := madeUp.FindString("x " + name + " y"); got != name {
			t.Errorf("madeUp finds %q in %q", got, name)
		}
		in := Codec{}.StandIn(name, 7)
		if madeUp.FindString(in) != in {
			t.Errorf("the stand-in %q of %q is not one madeUp finds", in, name)
		}
	}
	if got := (Codec{}).StandIn("DF__t__a__5AEE82B9", 3); got != "DF__t__a__standin3" {
		t.Errorf("StandIn = %q", got)
	}
	if got := (Codec{}).StandIn("oid:581577110", 3); got != "oid:03" {
		t.Errorf("StandIn = %q", got)
	}
	for _, name := range []string{"pk_t", "PK__b", "oid:", "ix_a", "DF__t__a__5aee82b9"} {
		if madeUp.MatchString(name) {
			t.Errorf("madeUp matches %q, a name the user gave", name)
		}
	}
}
