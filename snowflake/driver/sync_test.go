package driver

import (
	"testing"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestExcludedSchemaPredicate pins the WHERE clause that leaves the system
// schemas out of a read of INFORMATION_SCHEMA.
func TestExcludedSchemaPredicate(t *testing.T) {
	const want = "LOWER(TABLE_SCHEMA) NOT IN ('information_schema')"
	if got := excludedSchemaPredicate("TABLE_SCHEMA"); got != want {
		t.Errorf("excludedSchemaPredicate = %q, want %q", got, want)
	}
}

// TestConvertYesNo pins the nullability answer INFORMATION_SCHEMA.COLUMNS gives.
func TestConvertYesNo(t *testing.T) {
	for in, want := range map[string]bool{"YES": true, "Y": true, "1": true, "NO": false, "N": false, "0": false} {
		got, err := convertYesNo(in)
		if err != nil {
			t.Errorf("convertYesNo(%q): %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("convertYesNo(%q) = %v, want %v", in, got, want)
		}
	}
	if _, err := convertYesNo("maybe"); err == nil {
		t.Error("convertYesNo accepted an answer that is neither")
	}
}

// TestStreamTypeAndMode pins the two vocabularies SHOW STREAMS states into the
// metadata's enums, and that an unknown word stays unspecified rather than
// becoming a value.
func TestStreamTypeAndMode(t *testing.T) {
	if got := streamTypeOf("DELTA"); got != metadatapb.StreamMetadata_TYPE_DELTA {
		t.Errorf("streamTypeOf(DELTA) = %v", got)
	}
	for _, word := range []string{"", "delta", "APPEND_ONLY"} {
		if got := streamTypeOf(word); got != metadatapb.StreamMetadata_TYPE_UNSPECIFIED {
			t.Errorf("streamTypeOf(%q) = %v, want unspecified", word, got)
		}
	}
	for word, want := range map[string]metadatapb.StreamMetadata_Mode{
		"DEFAULT":     metadatapb.StreamMetadata_MODE_DEFAULT,
		"APPEND_ONLY": metadatapb.StreamMetadata_MODE_APPEND_ONLY,
		"INSERT_ONLY": metadatapb.StreamMetadata_MODE_INSERT_ONLY,
	} {
		if got := streamModeOf(word); got != want {
			t.Errorf("streamModeOf(%q) = %v, want %v", word, got, want)
		}
	}
	for _, word := range []string{"", "default", "TRIGGERLESS"} {
		if got := streamModeOf(word); got != metadatapb.StreamMetadata_MODE_UNSPECIFIED {
			t.Errorf("streamModeOf(%q) = %v, want unspecified", word, got)
		}
	}
}

// TestTaskState pins the states SHOW TASKS reports in the metadata's
// vocabulary.
func TestTaskState(t *testing.T) {
	if got := taskStateOf("started"); got != metadatapb.TaskMetadata_STATE_STARTED {
		t.Errorf("taskStateOf(started) = %v", got)
	}
	if got := taskStateOf("suspended"); got != metadatapb.TaskMetadata_STATE_SUSPENDED {
		t.Errorf("taskStateOf(suspended) = %v", got)
	}
	for _, state := range []string{"", "STARTED", "complete"} {
		if got := taskStateOf(state); got != metadatapb.TaskMetadata_STATE_UNSPECIFIED {
			t.Errorf("taskStateOf(%q) = %v, want unspecified", state, got)
		}
	}
}

// TestArrayStringScansAPredecessorList reads what SHOW TASKS states for a task's
// predecessors, which is an array of names written as a JSON document.
func TestArrayStringScansAPredecessorList(t *testing.T) {
	for _, tc := range []struct {
		name    string
		src     any
		want    []string
		wantErr bool
	}{
		{name: "text", src: `["A","B"]`, want: []string{"A", "B"}},
		{name: "bytes", src: []byte(`["A"]`), want: []string{"A"}},
		{name: "empty", src: `[]`, want: []string{}},
		// A task that is no other task's predecessor, or a server that left the
		// column NULL, has none.
		{name: "NULL", src: nil, want: nil},
		{name: "not a document", src: `A,B`, wantErr: true},
		{name: "not text", src: int64(1), wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got ArrayString
			err := got.Scan(tc.src)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("Scan(%v) = %v, want an error", tc.src, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("Scan(%v): %v", tc.src, err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("Scan(%v) = %v, want %v", tc.src, got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("predecessor %d = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

// TestSystemSchemas pins the schemas a sync leaves out of a database.
func TestSystemSchemas(t *testing.T) {
	if !systemSchemas["information_schema"] {
		t.Error("information_schema is not excluded")
	}
	for _, name := range []string{"PUBLIC", "APP", "public"} {
		if systemSchemas[name] {
			t.Errorf("%q is excluded", name)
		}
	}
}

// TestInformationSchema pins the schema path a sync reads: the database is one
// identifier, quoted once, and the schema's own name follows it.
func TestInformationSchema(t *testing.T) {
	for _, tc := range []struct{ database, want string }{
		{database: "DEMO_DB", want: `"DEMO_DB".INFORMATION_SCHEMA`},
		{database: "analytics", want: `"analytics".INFORMATION_SCHEMA`},
		{database: `a"b`, want: `"a""b".INFORMATION_SCHEMA`},
	} {
		if got := informationSchema(tc.database); got != tc.want {
			t.Errorf("informationSchema(%q) = %q, want %q", tc.database, got, tc.want)
		}
	}
}
