package driver

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bytebase/omni/driver"
)

// TestPing reads one database name as a connectivity check.
func TestPing(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()
	a.NoError(d.Ping(context.Background()))
}

// TestClose closes the clients and their sessions.
func TestClose(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	a.NoError(d.Close())
}

// TestOpenRequiresTheResourceName refuses a config that cannot name a database:
// the project and the instance are the two parts the contract has no field for.
func TestOpenRequiresTheResourceName(t *testing.T) {
	a := require.New(t)
	_, err := Open(context.Background(), driver.Config{Params: map[string]string{ParamInstance: "i"}})
	a.EqualError(err, "project ID cannot be empty")
	_, err = Open(context.Background(), driver.Config{Params: map[string]string{ParamProject: "p"}})
	a.EqualError(err, "instance ID cannot be empty")
}

// TestClassifyStatement holds the routing each statement's first keyword gets.
func TestClassifyStatement(t *testing.T) {
	tests := []struct {
		text string
		kind statementKind
		ok   bool
	}{
		{text: "SELECT 1", kind: query, ok: true},
		{text: "select 1", kind: query, ok: true},
		{text: "WITH n AS (SELECT 1 AS x) SELECT x FROM n", kind: query, ok: true},
		{text: "CREATE TABLE t (Id INT64 NOT NULL) PRIMARY KEY (Id)", kind: ddl, ok: true},
		{text: "ALTER TABLE t ADD COLUMN a INT64", kind: ddl, ok: true},
		{text: "DROP TABLE t", kind: ddl, ok: true},
		{text: "INSERT INTO t (Id) VALUES (1)", kind: dml, ok: true},
		{text: "UPDATE t SET a = 1", kind: dml, ok: true},
		{text: "DELETE FROM t", kind: dml, ok: true},
		// Comments before the keyword are the lexer's to skip, as Bytebase
		// stripped them before its prefix test.
		{text: "-- a note\nSELECT 1", kind: query, ok: true},
		{text: "/* a note */ CREATE TABLE t (Id INT64 NOT NULL) PRIMARY KEY (Id)", kind: ddl, ok: true},
		// Spanner's GRANT and REVOKE are DDL, which the admin API takes and the
		// DML path refuses.
		{text: "GRANT SELECT ON TABLE t TO ROLE r", kind: ddl, ok: true},
		{text: "revoke select on table t from role r", kind: ddl, ok: true},
		// A segment that holds only a comment holds no statement.
		{text: "-- only a note", ok: false},
		{text: "/* only a note */", ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.text, func(t *testing.T) {
			kind, ok := classifyStatement(tc.text)
			require.Equal(t, tc.ok, ok)
			if ok {
				require.Equal(t, tc.kind, kind)
			}
		})
	}
}

// TestStatements splits a script into the statements the driver runs, skipping
// the segments that ask Spanner nothing.
// TestFirstTwoKeywords reads the keywords a creation script's check needs.
func TestFirstTwoKeywords(t *testing.T) {
	tests := map[string][2]string{
		"CREATE DATABASE `d`":                       {"CREATE", "DATABASE"},
		"create database `d`":                       {"CREATE", "DATABASE"},
		"-- a note\nCREATE DATABASE `d`":            {"CREATE", "DATABASE"},
		"/* a note */\nCREATE /* x */ DATABASE `d`": {"CREATE", "DATABASE"},
		"CREATE":         {"CREATE", ""},
		"-- only a note": {"", ""},
	}
	for text, want := range tests {
		require.Equal(t, want, firstTwoKeywords(text), text)
	}
}

func TestStatements(t *testing.T) {
	a := require.New(t)
	stmts := statements(`SELECT 1;
		-- a note
		UPDATE t SET a = 1;
		;
		/* done */`)
	a.Len(stmts, 2)
	a.Equal("SELECT 1", stmts[0].Text)
	a.Equal(query, stmts[0].Kind)
	// The comment before the second statement is part of its text: the moved
	// splitter keeps the text the user wrote, and the first keyword decides the
	// routing.
	a.Equal("-- a note\n\t\tUPDATE t SET a = 1", stmts[1].Text)
	a.Equal(dml, stmts[1].Kind)
}

func TestGetDatabaseFromDSN(t *testing.T) {
	tests := []struct {
		dsn   string
		match bool
		want  string
	}{
		{dsn: "projects/p/instances/i/databases/d", match: true, want: "d"},
		{dsn: "projects/p/instances/i/databases/", match: false},
	}
	a := require.New(t)
	for _, test := range tests {
		got, err := getDatabaseFromDSN(test.dsn)
		if !test.match {
			a.Error(err)
			continue
		}
		a.NoError(err)
		a.Equal(test.want, got)
	}
}

func TestGetColumnTypeName(t *testing.T) {
	tests := []struct {
		spannerType *sppb.Type
		want        string
		wantErr     string
	}{
		{spannerType: &sppb.Type{Code: sppb.TypeCode_JSON}, want: "JSON"},
		{spannerType: &sppb.Type{Code: sppb.TypeCode_DATE}, want: "DATE"},
		{spannerType: &sppb.Type{Code: sppb.TypeCode_TIMESTAMP}, want: "TIMESTAMP"},
		{
			spannerType: &sppb.Type{Code: sppb.TypeCode_ARRAY, ArrayElementType: &sppb.Type{Code: sppb.TypeCode_BYTES}},
			want:        "[]BYTES",
		},
		{spannerType: &sppb.Type{Code: sppb.TypeCode_STRUCT}, wantErr: "spanner STRUCT type is not supported"},
		{
			spannerType: &sppb.Type{Code: sppb.TypeCode_ARRAY, ArrayElementType: &sppb.Type{Code: sppb.TypeCode_STRUCT}},
			wantErr:     "spanner STRUCT type is not supported",
		},
	}
	a := require.New(t)
	for _, test := range tests {
		got, err := getColumnTypeName(test.spannerType)
		if test.wantErr != "" {
			a.EqualError(err, test.wantErr)
			continue
		}
		a.NoError(err)
		a.Equal(test.want, got)
	}
}

// timestampCell is the cell a TIMESTAMP reads as, for the conversion tests.
func timestampCell(t *testing.T, text string, precision int32) driver.Value {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339Nano, text)
	require.NoError(t, err)
	require.True(t, timestamppb.New(parsed).IsValid())
	return driver.Value{Kind: driver.Timestamp, Time: parsed, Precision: precision}
}

// TestConvertSpannerValue holds every type the client decodes to the contract's
// cell, with the cases Bytebase's driver covered before the move.
func TestConvertSpannerValue(t *testing.T) {
	tests := []struct {
		name     string
		colType  *sppb.Type
		value    *structpb.Value
		expected driver.Value
	}{
		{name: "nil value", expected: driver.Value{Kind: driver.Null}},
		{
			name:     "null value",
			colType:  &sppb.Type{Code: sppb.TypeCode_STRING},
			value:    structpb.NewNullValue(),
			expected: driver.Value{Kind: driver.Null},
		},
		{
			name:     "string value",
			colType:  &sppb.Type{Code: sppb.TypeCode_STRING},
			value:    structpb.NewStringValue("hello"),
			expected: driver.Value{Kind: driver.String, String: "hello"},
		},
		{
			name:     "bool value",
			colType:  &sppb.Type{Code: sppb.TypeCode_BOOL},
			value:    structpb.NewBoolValue(true),
			expected: driver.Value{Kind: driver.Bool, Bool: true},
		},
		{
			name:     "float64 value",
			colType:  &sppb.Type{Code: sppb.TypeCode_FLOAT64},
			value:    structpb.NewNumberValue(3.14159),
			expected: driver.Value{Kind: driver.Float64, Float64: 3.14159},
		},
		{
			name:     "int64 as a number",
			colType:  &sppb.Type{Code: sppb.TypeCode_INT64},
			value:    structpb.NewNumberValue(42),
			expected: driver.Value{Kind: driver.Int64, Int64: 42},
		},
		{
			name:     "number without type metadata defaults to double",
			value:    structpb.NewNumberValue(42),
			expected: driver.Value{Kind: driver.Float64, Float64: 42},
		},
		{
			name:     "int64 as a string, the way Spanner writes it",
			colType:  &sppb.Type{Code: sppb.TypeCode_INT64},
			value:    structpb.NewStringValue("9223372036854775807"),
			expected: driver.Value{Kind: driver.Int64, Int64: math.MaxInt64},
		},
		{
			name:     "int64 as a string, smallest",
			colType:  &sppb.Type{Code: sppb.TypeCode_INT64},
			value:    structpb.NewStringValue("-9223372036854775808"),
			expected: driver.Value{Kind: driver.Int64, Int64: math.MinInt64},
		},
		{
			name:     "int64 as a string, unparsable, stays text",
			colType:  &sppb.Type{Code: sppb.TypeCode_INT64},
			value:    structpb.NewStringValue("not-a-number"),
			expected: driver.Value{Kind: driver.String, String: "not-a-number"},
		},
		{
			name:     "array value",
			colType:  &sppb.Type{Code: sppb.TypeCode_ARRAY},
			value:    structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{structpb.NewNumberValue(1), structpb.NewNumberValue(2), structpb.NewNumberValue(3)}}),
			expected: driver.Value{Kind: driver.String, String: "[1,2,3]"},
		},
		{
			name:    "struct value",
			colType: &sppb.Type{Code: sppb.TypeCode_STRUCT},
			value: structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
				"name": structpb.NewStringValue("Alice"),
				"age":  structpb.NewNumberValue(30),
			}}),
			expected: driver.Value{Kind: driver.String, String: `{"age":30,"name":"Alice"}`},
		},
		{
			name:    "nested array",
			colType: &sppb.Type{Code: sppb.TypeCode_ARRAY},
			value: structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
				structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{structpb.NewNumberValue(1), structpb.NewNumberValue(2)}}),
				structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{structpb.NewNumberValue(3), structpb.NewNumberValue(4)}}),
			}}),
			expected: driver.Value{Kind: driver.String, String: "[[1,2],[3,4]]"},
		},
		{
			name:    "array of structs",
			colType: &sppb.Type{Code: sppb.TypeCode_ARRAY},
			value: structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{"id": structpb.NewNumberValue(1)}}),
				structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{"id": structpb.NewNumberValue(2)}}),
			}}),
			expected: driver.Value{Kind: driver.String, String: `[{"id":1},{"id":2}]`},
		},
		{
			name:     "timestamp with microseconds",
			colType:  &sppb.Type{Code: sppb.TypeCode_TIMESTAMP},
			value:    structpb.NewStringValue("2024-01-15T10:30:45.123456Z"),
			expected: timestampCell(t, "2024-01-15T10:30:45.123456Z", 6),
		},
		{
			name:     "timestamp with nanoseconds",
			colType:  &sppb.Type{Code: sppb.TypeCode_TIMESTAMP},
			value:    structpb.NewStringValue("2024-01-15T10:30:45.123456789Z"),
			expected: timestampCell(t, "2024-01-15T10:30:45.123456789Z", 9),
		},
		{
			name:     "timestamp without a fraction keeps microseconds",
			colType:  &sppb.Type{Code: sppb.TypeCode_TIMESTAMP},
			value:    structpb.NewStringValue("2024-01-15T10:30:45Z"),
			expected: timestampCell(t, "2024-01-15T10:30:45Z", 6),
		},
		{
			name:     "timestamp with more than nine digits is capped",
			colType:  &sppb.Type{Code: sppb.TypeCode_TIMESTAMP},
			value:    structpb.NewStringValue("2024-01-15T10:30:45.123456789012Z"),
			expected: timestampCell(t, "2024-01-15T10:30:45.123456789012Z", 9),
		},
		{
			name:    "timestamp unparsable stays text",
			colType: &sppb.Type{Code: sppb.TypeCode_TIMESTAMP},
			value:   structpb.NewStringValue("not-a-timestamp"),
			expected: driver.Value{
				Kind:   driver.String,
				String: "not-a-timestamp",
			},
		},
		{
			name:     "bytes value",
			colType:  &sppb.Type{Code: sppb.TypeCode_BYTES},
			value:    structpb.NewStringValue("SGVsbG8gV29ybGQ="),
			expected: driver.Value{Kind: driver.Bytes, Bytes: []byte("Hello World")},
		},
		{
			name:     "bytes that are not base64 stay text",
			colType:  &sppb.Type{Code: sppb.TypeCode_BYTES},
			value:    structpb.NewStringValue("not-valid-base64!"),
			expected: driver.Value{Kind: driver.String, String: "not-valid-base64!"},
		},
		{
			name:     "date value",
			colType:  &sppb.Type{Code: sppb.TypeCode_DATE},
			value:    structpb.NewStringValue("2024-01-15"),
			expected: driver.Value{Kind: driver.String, String: "2024-01-15"},
		},
		{
			name:     "json value",
			colType:  &sppb.Type{Code: sppb.TypeCode_JSON},
			value:    structpb.NewStringValue(`{"name":"Alice","age":30}`),
			expected: driver.Value{Kind: driver.String, String: `{"name":"Alice","age":30}`},
		},
		{
			name:     "numeric value stays text",
			colType:  &sppb.Type{Code: sppb.TypeCode_NUMERIC},
			value:    structpb.NewStringValue("123.456"),
			expected: driver.Value{Kind: driver.String, String: "123.456"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, convertSpannerValue(tc.colType, tc.value))
		})
	}
}

// TestConvertQueryPlanToJSON renders a plan as the document Bytebase showed: a
// plan of two nodes, one of them with child links, a short representation and
// metadata.
func TestConvertQueryPlanToJSON(t *testing.T) {
	a := require.New(t)
	got, err := convertQueryPlanToJSON(&sppb.QueryPlan{PlanNodes: []*sppb.PlanNode{
		{
			Index:       0,
			Kind:        sppb.PlanNode_RELATIONAL,
			DisplayName: "Project",
			ChildLinks:  []*sppb.PlanNode_ChildLink{{ChildIndex: 1, Type: "Scalar", Variable: "v"}},
			ShortRepresentation: &sppb.PlanNode_ShortRepresentation{
				Description: "1, 2",
				Subqueries:  map[string]int32{"sub": 1},
			},
			Metadata: mustStruct(t, map[string]any{"rows": 3}),
		},
		{
			Index:       1,
			Kind:        sppb.PlanNode_SCALAR,
			DisplayName: "Scan",
		},
	}})
	a.NoError(err)
	a.Equal(`{"planNodes":[{"index":0,"kind":"RELATIONAL","displayName":"Project","childLinks":[{"childIndex":1,"type":"Scalar","variable":"v"}],"shortRepresentation":{"description":"1, 2","subqueries":{"sub":1}},"metadata":{"rows":3}},{"index":1,"kind":"SCALAR","displayName":"Scan"}]}`, got)

	empty, err := convertQueryPlanToJSON(nil)
	a.NoError(err)
	a.Equal("{}", empty)
}

func mustStruct(t *testing.T, fields map[string]any) *structpb.Struct {
	t.Helper()
	value, err := structpb.NewStruct(fields)
	require.NoError(t, err)
	return value
}

// TestCredentialOption dispatches the credential document the caller resolved,
// and refuses one the authentication library has no kind for.
func TestCredentialOption(t *testing.T) {
	a := require.New(t)
	option, err := credentialOption("")
	a.NoError(err)
	a.Nil(option)

	option, err = credentialOption(`{"type":"service_account","project_id":"p"}`)
	a.NoError(err)
	a.NotNil(option)

	_, err = credentialOption(`{"type":"unknown"}`)
	a.EqualError(err, `unsupported GCP credential type: "unknown"`)
	_, err = credentialOption(`{"type":""}`)
	a.EqualError(err, `GCP credential JSON missing "type" field`)
	_, err = credentialOption(`not json`)
	a.ErrorContains(err, "failed to parse GCP credential JSON")
}

// TestEndpoint is the address the client is pointed at: host:port when the caller
// named a host, and nothing to leave the client's own.
func TestEndpoint(t *testing.T) {
	tests := []struct {
		host, port string
		want       string
		wantHost   string
	}{
		{want: ""},
		{host: "spanner.example.com", want: "spanner.example.com", wantHost: "spanner.example.com"},
		{host: "spanner.example.com", port: "443", want: "spanner.example.com:443", wantHost: "spanner.example.com"},
		{host: "10.0.0.1", port: "8443", want: "10.0.0.1:8443", wantHost: "10.0.0.1"},
	}
	for _, test := range tests {
		cfg := driver.Config{Host: test.host, Port: test.port}
		endpoint := endpoint(cfg)
		require.Equal(t, test.want, endpoint)
		require.Equal(t, test.wantHost, endpointHost(endpoint))
	}
}

// TestDataCallWithoutADatabase answers an error for a driver opened without a
// database, which Bytebase's driver dereferenced as a nil client.
func TestDataCallWithoutADatabase(t *testing.T) {
	a := require.New(t)
	client := &liveClient{}
	_, err := client.Read(context.Background(), "SELECT 1")
	a.EqualError(err, errNoDatabase.Error())
	_, err = client.Snapshot(context.Background())
	a.EqualError(err, errNoDatabase.Error())
	_, err = client.Update(context.Background(), "UPDATE t SET a = 1")
	a.EqualError(err, errNoDatabase.Error())
	_, err = client.BatchUpdate(context.Background(), []string{"UPDATE t SET a = 1"})
	a.EqualError(err, errNoDatabase.Error())
	_, err = client.PartitionedUpdate(context.Background(), "UPDATE t SET a = 1")
	a.EqualError(err, errNoDatabase.Error())
	_, err = client.AnalyzeQuery(context.Background(), "SELECT 1", false)
	a.EqualError(err, errNoDatabase.Error())
	a.True(errors.Is(err, errNoDatabase))
}

// TestScan decodes the cells a sync's queries read.
func TestScan(t *testing.T) {
	a := require.New(t)
	var (
		name     string
		nullable bool
		position int64
		count    nullString
		columns  []string
	)
	cells := []Cell{
		{Type: &sppb.Type{Code: sppb.TypeCode_STRING}, Value: structpb.NewStringValue("Singers")},
		{Type: &sppb.Type{Code: sppb.TypeCode_BOOL}, Value: structpb.NewBoolValue(true)},
		{Type: &sppb.Type{Code: sppb.TypeCode_INT64}, Value: structpb.NewStringValue("7")},
		{Type: &sppb.Type{Code: sppb.TypeCode_STRING}, Value: structpb.NewNullValue()},
		{Type: &sppb.Type{Code: sppb.TypeCode_ARRAY}, Value: structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
			structpb.NewStringValue("Id"), structpb.NewStringValue("Name"),
		}})},
	}
	a.NoError(scan(cells, &name, &nullable, &position, &count, &columns))
	a.Equal("Singers", name)
	a.True(nullable)
	a.Equal(int64(7), position)
	a.False(count.Valid)
	a.Equal([]string{"Id", "Name"}, columns)

	// A NULL into a destination that cannot hold it, and a row of the wrong
	// width, are errors rather than a silently wrong reading.
	a.Error(scan([]Cell{{Value: structpb.NewNullValue()}}, &name))
	a.Error(scan(cells[:2], &name))
	// A destination the sync does not read has no decoding of its own.
	var unsupported []byte
	a.Error(scan([]Cell{{Value: structpb.NewStringValue("x")}}, &unsupported))
}

// TestTimestampAccuracy reads the fractional digits the server wrote.
func TestTimestampAccuracy(t *testing.T) {
	tests := map[string]int32{
		"2024-01-15T10:30:45Z":                6,
		"2024-01-15T10:30:45.1Z":              1,
		"2024-01-15T10:30:45.123456Z":         6,
		"2024-01-15T10:30:45.123456789Z":      9,
		"2024-01-15T10:30:45.123456789012Z":   9,
		"2024-01-15T10:30:45.123456+08:00":    6,
		"2024-01-15T10:30:45.123456789-07:00": 9,
		"2024-01-15T10:30:45.123456789012+08": 9,
	}
	for timestamp, want := range tests {
		require.Equal(t, want, timestampAccuracy(timestamp), timestamp)
	}
}

// TestStatementPaths are the resource names the driver builds and reads.
func TestStatementPaths(t *testing.T) {
	a := require.New(t)
	a.Equal("projects/p/instances/i", InstancePath("p", "i"))
	a.Equal("projects/p/instances/i/databases/d", DatabasePath("projects/p/instances/i", "d"))
}

// TestLogTransactionControlWithoutALogger logs nothing, so a caller that wants no
// log need not pass one.
func TestLogTransactionControlWithoutALogger(t *testing.T) {
	logTransactionControl(nil, driver.Begin, nil)
}

// TestSessionIsTheDriver states that a session holds no state of its own: closing
// it closes nothing.
func TestSessionIsTheDriver(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()
	session, err := d.Session(context.Background())
	a.NoError(err)
	a.NoError(session.Close())
}

// TestStatementTextKeepsComments states what changed with the move's splitter:
// the statement the driver sends and reports is the text the user wrote, where
// Bytebase removed its comments before sending it.
func TestStatementTextKeepsComments(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "query")
	defer d.Close()

	results, err := sessionOf(t, d).Query(context.Background(), "SELECT 1 AS a -- a note", driver.QueryOptions{MaxResultBytes: maxResultBytes})
	a.NoError(err)
	a.Equal("SELECT 1 AS a -- a note", results[0].Statement)
	a.NoError(results[0].Err)
	a.Equal(int64(1), results[0].Rows[0][0].Int64)
}
