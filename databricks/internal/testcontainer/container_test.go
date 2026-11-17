package testcontainer

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
)

// TestRequestAndAnswerRoundTrip keeps the tape's own encoding readable by the
// transport after a recording.
func TestRequestAndAnswerRoundTrip(t *testing.T) {
	text, err := request{Method: "POST", Route: "/api/2.0/sql/statements", Body: `{"statement":"SELECT 1"}`}.text()
	if err != nil {
		t.Fatalf("request text: %v", err)
	}
	var got request
	if err := json.Unmarshal([]byte(text), &got); err != nil {
		t.Fatalf("unreadable request: %v", err)
	}
	if got.Method != "POST" || got.Route != "/api/2.0/sql/statements" || !strings.Contains(got.Body, "SELECT 1") {
		t.Errorf("request round trip = %+v", got)
	}

	a, err := json.Marshal(answer{Status: 200, ContentType: "application/json", Body: `{"statement_id":"x"}`})
	if err != nil {
		t.Fatalf("answer text: %v", err)
	}
	var decoded answer
	if err := json.Unmarshal(a, &decoded); err != nil {
		t.Fatalf("unreadable answer: %v", err)
	}
	if decoded.Status != 200 || decoded.ContentType != "application/json" || decoded.Body != `{"statement_id":"x"}` {
		t.Errorf("answer round trip = %+v", decoded)
	}
}

// TestNarrowStatementKeepsWhatTheDriverReads drops the service's timings,
// chunks and links, and keeps the state, columns and rows.
func TestNarrowStatementKeepsWhatTheDriverReads(t *testing.T) {
	body := []byte(`{
		"statement_id":"01f1e285-b8ba-1656-acc4-631d2033d907",
		"status":{"state":"SUCCEEDED"},
		"manifest":{"format":"JSON_ARRAY","schema":{"column_count":1,"columns":[{"name":"n","type_text":"INT","type_name":"INT","position":0}]},"total_chunk_count":1,"total_row_count":1,"truncated":false,"chunks":[{"chunk_index":0}]},
		"result":{"chunk_index":0,"row_offset":0,"row_count":1,"data_array":[["1"]],"next_chunk_index":1}
	}`)
	narrowed, ok := narrowStatement(body)
	if !ok {
		t.Fatal("narrowStatement refused the answer")
	}
	text := string(narrowed)
	for _, want := range []string{`"state":"SUCCEEDED"`, `"name":"n"`, `"data_array":[["1"]]`, `"truncated":false`} {
		if !strings.Contains(text, want) {
			t.Errorf("narrowed answer %s does not hold %s", text, want)
		}
	}
	for _, unwanted := range []string{"total_chunk_count", "next_chunk_index", "chunks", "row_count"} {
		if strings.Contains(text, unwanted) {
			t.Errorf("narrowed answer %s still holds %s", text, unwanted)
		}
	}
}

// TestNarrowTablesSortsByNameAndPosition makes a listing's order the tape's.
func TestNarrowTablesSortsByNameAndPosition(t *testing.T) {
	body := []byte(`{"tables":[
		{"name":"b","columns":[{"name":"z","position":1},{"name":"a","position":0}]},
		{"name":"a","columns":[]}
	]}`)
	narrowed, ok := narrowTables(body)
	if !ok {
		t.Fatal("narrowTables refused the answer")
	}
	var got struct {
		Tables []table `json:"tables"`
	}
	if err := json.Unmarshal(narrowed, &got); err != nil {
		t.Fatalf("unreadable narrowed answer: %v", err)
	}
	if len(got.Tables) != 2 || got.Tables[0].Name != "a" || got.Tables[1].Name != "b" {
		t.Fatalf("tables = %+v, want a then b", got.Tables)
	}
	if got.Tables[1].Columns[0].Name != "a" || got.Tables[1].Columns[1].Name != "z" {
		t.Errorf("columns = %+v, want a then z", got.Tables[1].Columns)
	}
}

// TestMadeUpStatementIDHasAStandIn keeps a poll's path consistent with the
// answer that named it, without recording the service's own ID.
func TestMadeUpStatementIDHasAStandIn(t *testing.T) {
	c := codec{}
	real := "01f1e285-b8ba-1656-acc4-631d2033d907"
	standIn := c.StandIn(real, 1)
	if !c.MadeUp().MatchString(standIn) {
		t.Fatalf("stand-in %q does not match the made-up pattern", standIn)
	}
	if standIn == real {
		t.Fatal("stand-in is the service's own ID")
	}
	if got := c.MadeUp().ReplaceAllString(`"statement_id":"`+real+`"`, standIn); !strings.Contains(got, standIn) {
		t.Errorf("replacement = %q, want the stand-in", got)
	}
}

// TestPinNamesTheRecordedWorkspace keeps the constants the tests hand the
// driver from drifting from the pin row the tapes name.
func TestPinNamesTheRecordedWorkspace(t *testing.T) {
	pin, err := enginetest.LoadPin(Pins, Driver)
	if err != nil {
		t.Fatalf("LoadPin: %v", err)
	}
	if !strings.Contains(pin.Tag, "warehouse_id="+WarehouseID) {
		t.Errorf("pin source tag %q does not name warehouse %q", pin.Tag, WarehouseID)
	}
	if !strings.Contains(pin.Tag, "catalog="+Catalog) {
		t.Errorf("pin source tag %q does not name catalog %q", pin.Tag, Catalog)
	}
}
