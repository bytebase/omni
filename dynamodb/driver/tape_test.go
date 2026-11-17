package driver

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/bytebase/omni/driver"
	"github.com/bytebase/omni/dynamodb/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/tape"
)

// tapeClient is the driver's client, answered by t's tape. It is what a test
// hands the driver in place of a live client; its methods are the calls the
// driver makes, and the recordings hold what a live client answered.
type tapeClient struct {
	tape *tape.Raw
}

// OpenTape returns t's client to the pinned endpoint and its fixture: a replay
// of its tape or, in a recording, the endpoint itself with every call written
// down.
func OpenTape(t testing.TB) *tapeClient {
	t.Helper()
	return openTape(t, testcontainer.Endpoint)
}

// OpenTapePaged returns t's client to the endpoint that holds the listing that
// takes two pages.
func OpenTapePaged(t testing.TB) *tapeClient {
	t.Helper()
	return openTape(t, testcontainer.PagedEndpoint)
}

func openTape(t testing.TB, endpoint func(testing.TB) (string, error)) *tapeClient {
	t.Helper()
	raw := tape.OpenRaw(t, testcontainer.Server, testcontainer.Default, func(request string) (string, error) {
		// A recording reaches the endpoint on its first call, so a test that
		// replays starts no container.
		return live(t, endpoint, request)
	})
	return &tapeClient{tape: raw}
}

func (c *tapeClient) Ping(ctx context.Context) error {
	_, err := c.do(ctx, request{Op: opPing})
	return err
}

func (c *tapeClient) Close() error {
	_, err := c.do(context.Background(), request{Op: opClose})
	return err
}

func (c *tapeClient) ListTables(ctx context.Context, exclusiveStartTableName string, limit int32) ([]string, string, error) {
	a, err := c.do(ctx, request{Op: opListTables, ExclusiveStartTableName: exclusiveStartTableName, Limit: limit})
	if err != nil {
		return nil, "", err
	}
	return a.Tables, a.LastEvaluatedTableName, nil
}

func (c *tapeClient) DescribeTable(ctx context.Context, tableName string) (*Table, error) {
	a, err := c.do(ctx, request{Op: opDescribeTable, Table: tableName})
	if err != nil {
		return nil, err
	}
	return a.Table.table(), nil
}

func (c *tapeClient) ExecuteStatement(ctx context.Context, statement StatementRequest) (StatementAnswer, error) {
	a, err := c.do(ctx, request{
		Op:        opExecuteStatement,
		Statement: statement.Statement,
		Limit:     statement.Limit,
		NextToken: statement.NextToken,
	})
	// The items come back with the error: a read that failed after a page must
	// not be read again, in a replay as in a recording.
	items := make([]map[string]types.AttributeValue, len(a.Items))
	for i, item := range a.Items {
		items[i] = item
	}
	return StatementAnswer{Items: items, NextToken: a.NextToken}, err
}

func (c *tapeClient) ExecuteTransaction(ctx context.Context, statement string) error {
	_, err := c.do(ctx, request{Op: opExecuteTransaction, Statement: statement})
	return err
}

func (c *tapeClient) CallerAccount(ctx context.Context) (string, error) {
	a, err := c.do(ctx, request{Op: opCallerAccount})
	if err != nil {
		return "", err
	}
	return a.Account, nil
}

// do sends one call and returns its answer.
func (c *tapeClient) do(_ context.Context, r request) (answer, error) {
	text, err := r.text()
	if err != nil {
		return answer{}, err
	}
	answered, err := c.tape.Do(text)
	if err != nil {
		return answer{}, err
	}
	var a answer
	if err := unmarshalJSON([]byte(answered), &a); err != nil {
		return answer{}, fmt.Errorf("test: unreadable answer to %s: %w", text, err)
	}
	if a.Error != nil {
		return a, a.Error.err()
	}
	return a, nil
}

// The calls the driver makes of its client, as a tape names them.
const (
	opPing               = "ping"
	opClose              = "close"
	opListTables         = "listTables"
	opDescribeTable      = "describeTable"
	opExecuteStatement   = "executeStatement"
	opExecuteTransaction = "executeTransaction"
	opCallerAccount      = "callerAccount"
)

// request is one call as a tape holds it.
type request struct {
	Op                      string `json:"op"`
	Table                   string `json:"table,omitempty"`
	ExclusiveStartTableName string `json:"exclusiveStartTableName,omitempty"`
	Limit                   int32  `json:"limit,omitempty"`
	Statement               string `json:"statement,omitempty"`
	NextToken               string `json:"nextToken,omitempty"`
}

func (r request) text() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// answer is one call's answer as a tape holds it.
type answer struct {
	Tables                 []string `json:"tables,omitempty"`
	LastEvaluatedTableName string   `json:"lastEvaluatedTableName,omitempty"`
	Table                  *table   `json:"table,omitempty"`
	Items                  []item   `json:"items,omitempty"`
	NextToken              string   `json:"nextToken,omitempty"`
	Account                string   `json:"account,omitempty"`
	Error                  *failure `json:"error,omitempty"`
}

// table is one DescribeTable answer as a tape holds it: the parts SyncDatabase
// reads, and none of the endpoint's creation time, size or throughput, which
// only the endpoint that answered can name.
type table struct {
	Name             string             `json:"name"`
	Status           string             `json:"status"`
	ItemCount        int64              `json:"itemCount,omitempty"`
	KeySchema        []keySchemaElement `json:"keySchema,omitempty"`
	SecondaryIndexes []secondaryIndex   `json:"secondaryIndexes,omitempty"`
}

type secondaryIndex struct {
	Name      string             `json:"name"`
	Global    bool               `json:"global,omitempty"`
	KeySchema []keySchemaElement `json:"keySchema,omitempty"`
}

type keySchemaElement struct {
	AttributeName string `json:"attributeName"`
	KeyType       string `json:"keyType"`
}

// tableAnswer is one DescribeTable answer as the tape holds it.
func tableAnswer(described *Table) *table {
	answered := &table{
		Name:      described.Name,
		Status:    described.Status,
		ItemCount: described.ItemCount,
		KeySchema: keySchemaElements(described.KeySchema),
	}
	for _, index := range described.SecondaryIndexes {
		answered.SecondaryIndexes = append(answered.SecondaryIndexes, secondaryIndex{
			Name:      index.Name,
			Global:    index.Global,
			KeySchema: keySchemaElements(index.KeySchema),
		})
	}
	return answered
}

// table is the driver's view of the table a tape holds.
func (t *table) table() *Table {
	described := &Table{Name: t.Name, Status: t.Status, ItemCount: t.ItemCount}
	described.KeySchema = keySchemaOf(t.KeySchema)
	for _, index := range t.SecondaryIndexes {
		described.SecondaryIndexes = append(described.SecondaryIndexes, SecondaryIndex{
			Name:      index.Name,
			Global:    index.Global,
			KeySchema: keySchemaOf(index.KeySchema),
		})
	}
	return described
}

func keySchemaElements(schema []KeySchemaElement) []keySchemaElement {
	keys := make([]keySchemaElement, 0, len(schema))
	for _, key := range schema {
		keys = append(keys, keySchemaElement{AttributeName: key.AttributeName, KeyType: key.KeyType})
	}
	return keys
}

func keySchemaOf(keys []keySchemaElement) []KeySchemaElement {
	schema := make([]KeySchemaElement, 0, len(keys))
	for _, key := range keys {
		schema = append(schema, KeySchemaElement{AttributeName: key.AttributeName, KeyType: key.KeyType})
	}
	return schema
}

// failure is an error the client returned, as a tape holds it: its text. The
// driver reports an error's text and reads nothing else of it, so a replay
// answers the text a recording saw.
type failure struct {
	Text string `json:"text"`
}

func (f *failure) err() error { return fmt.Errorf("%s", f.Text) }

// item is one DynamoDB item as a tape holds it: the JSON form DynamoDB's own
// API is written in, where every value is tagged with its type, so a replay
// answers with the types the endpoint answered rather than what a decoder made
// of them.
type item map[string]types.AttributeValue

func (i item) MarshalJSON() ([]byte, error) {
	encoded := make(map[string]any, len(i))
	for name, value := range i {
		attribute, err := encodeAttributeValue(value)
		if err != nil {
			return nil, err
		}
		encoded[name] = attribute
	}
	return json.Marshal(encoded)
}

func (i *item) UnmarshalJSON(b []byte) error {
	var encoded map[string]json.RawMessage
	if err := json.Unmarshal(b, &encoded); err != nil {
		return err
	}
	decoded := make(item, len(encoded))
	for name, attribute := range encoded {
		value, err := decodeAttributeValue(attribute)
		if err != nil {
			return err
		}
		decoded[name] = value
	}
	*i = decoded
	return nil
}

func encodeAttributeValue(value types.AttributeValue) (map[string]any, error) {
	switch value := value.(type) {
	case *types.AttributeValueMemberS:
		return map[string]any{"S": value.Value}, nil
	case *types.AttributeValueMemberN:
		return map[string]any{"N": value.Value}, nil
	case *types.AttributeValueMemberB:
		return map[string]any{"B": base64.StdEncoding.EncodeToString(value.Value)}, nil
	case *types.AttributeValueMemberBOOL:
		return map[string]any{"BOOL": value.Value}, nil
	case *types.AttributeValueMemberNULL:
		return map[string]any{"NULL": value.Value}, nil
	case *types.AttributeValueMemberSS:
		return map[string]any{"SS": value.Value}, nil
	case *types.AttributeValueMemberNS:
		return map[string]any{"NS": value.Value}, nil
	case *types.AttributeValueMemberBS:
		encoded := make([]string, 0, len(value.Value))
		for _, b := range value.Value {
			encoded = append(encoded, base64.StdEncoding.EncodeToString(b))
		}
		return map[string]any{"BS": encoded}, nil
	case *types.AttributeValueMemberL:
		encoded := make([]map[string]any, 0, len(value.Value))
		for _, element := range value.Value {
			attribute, err := encodeAttributeValue(element)
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, attribute)
		}
		return map[string]any{"L": encoded}, nil
	case *types.AttributeValueMemberM:
		encoded := make(map[string]map[string]any, len(value.Value))
		for name, element := range value.Value {
			attribute, err := encodeAttributeValue(element)
			if err != nil {
				return nil, err
			}
			encoded[name] = attribute
		}
		return map[string]any{"M": encoded}, nil
	}
	return nil, fmt.Errorf("test: unsupported attribute value type %T", value)
}

func decodeAttributeValue(b []byte) (types.AttributeValue, error) {
	var tagged map[string]json.RawMessage
	if err := json.Unmarshal(b, &tagged); err != nil {
		return nil, err
	}
	if len(tagged) != 1 {
		return nil, fmt.Errorf("test: an attribute value is tagged with one type, not %d", len(tagged))
	}
	for tag, raw := range tagged {
		switch tag {
		case "S":
			var v string
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			return &types.AttributeValueMemberS{Value: v}, nil
		case "N":
			var v string
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			return &types.AttributeValueMemberN{Value: v}, nil
		case "B":
			var v string
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			decoded, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return nil, err
			}
			return &types.AttributeValueMemberB{Value: decoded}, nil
		case "BOOL":
			var v bool
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			return &types.AttributeValueMemberBOOL{Value: v}, nil
		case "NULL":
			var v bool
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			return &types.AttributeValueMemberNULL{Value: v}, nil
		case "SS":
			var v []string
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			return &types.AttributeValueMemberSS{Value: v}, nil
		case "NS":
			var v []string
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			return &types.AttributeValueMemberNS{Value: v}, nil
		case "BS":
			var v []string
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			decoded := make([][]byte, 0, len(v))
			for _, s := range v {
				b, err := base64.StdEncoding.DecodeString(s)
				if err != nil {
					return nil, err
				}
				decoded = append(decoded, b)
			}
			return &types.AttributeValueMemberBS{Value: decoded}, nil
		case "L":
			var v []json.RawMessage
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			decoded := make([]types.AttributeValue, 0, len(v))
			for _, element := range v {
				attribute, err := decodeAttributeValue(element)
				if err != nil {
					return nil, err
				}
				decoded = append(decoded, attribute)
			}
			return &types.AttributeValueMemberL{Value: decoded}, nil
		case "M":
			var v map[string]json.RawMessage
			if err := json.Unmarshal(raw, &v); err != nil {
				return nil, err
			}
			decoded := make(map[string]types.AttributeValue, len(v))
			for name, element := range v {
				attribute, err := decodeAttributeValue(element)
				if err != nil {
					return nil, err
				}
				decoded[name] = attribute
			}
			return &types.AttributeValueMemberM{Value: decoded}, nil
		}
		return nil, fmt.Errorf("test: unknown attribute value type %q", tag)
	}
	return nil, nil
}

// ordersResult matches the clause that decides a result's own order. A read
// that names none answers in whatever order the endpoint's table iterated in,
// which is not the same between boots.
var ordersResult = regexp.MustCompile(`(?i)\border\s+by\b`)

// live answers one request from the pinned endpoint, through the client the
// driver itself opens.
func live(t testing.TB, endpoint func(testing.TB) (string, error), text string) (string, error) {
	var req request
	if err := unmarshalJSON([]byte(text), &req); err != nil {
		return "", fmt.Errorf("test: unreadable request %s: %w", text, err)
	}
	// The identity call reaches STS, which DynamoDB Local is not and whose
	// real endpoint a recording must not ask with the fixture's placeholder
	// key pair. The harness answers as the fixture's account, before any
	// endpoint starts, so the tape holds the call and the replay needs none.
	if req.Op == opCallerAccount {
		b, err := json.Marshal(answer{Account: testcontainer.AccountID})
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	url, err := endpoint(t)
	if err != nil {
		return "", err
	}
	d, err := Open(context.Background(), driver.Config{
		Host:     url,
		Params:   map[string]string{ParamRegion: testcontainer.Region},
		User:     testcontainer.AccessKeyID,
		Password: testcontainer.SecretAccessKey,
	})
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	a := answer{}
	switch req.Op {
	case opPing:
		err = d.client.Ping(ctx)
	case opClose:
		err = d.client.Close()
	case opListTables:
		a.Tables, a.LastEvaluatedTableName, err = d.client.ListTables(ctx, req.ExclusiveStartTableName, req.Limit)
	case opDescribeTable:
		var described *Table
		described, err = d.client.DescribeTable(ctx, req.Table)
		if described != nil {
			a.Table = tableAnswer(described)
		}
	case opExecuteStatement:
		var answered StatementAnswer
		answered, err = d.client.ExecuteStatement(ctx, StatementRequest{
			Statement: req.Statement,
			Limit:     req.Limit,
			NextToken: req.NextToken,
		})
		for _, read := range answered.Items {
			a.Items = append(a.Items, canonicalItem(read))
		}
		a.NextToken = answered.NextToken
		if !ordersResult.MatchString(req.Statement) {
			// The order of a read that names no ORDER BY is the endpoint's to
			// choose, and its table iteration differs between boots. The tape
			// holds one page in canonical order, and a test that draws a
			// conclusion from a row's position names the clause that decides it.
			slices.SortFunc(a.Items, func(x, y item) int {
				return strings.Compare(mustJSON(x), mustJSON(y))
			})
		}
	case opExecuteTransaction:
		err = d.client.ExecuteTransaction(ctx, req.Statement)
	default:
		return "", fmt.Errorf("test: unknown request %q", req.Op)
	}
	if err != nil {
		a.Error = &failure{Text: err.Error()}
	}
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// canonicalItem holds a set's elements in one order. A DynamoDB set is
// unordered: the endpoint may answer {a, b} as {b, a}, so a recorded answer is
// held in one canonical order and a test that draws a conclusion from a set
// reads it as a set, not as a sequence.
func canonicalItem(read map[string]types.AttributeValue) item {
	for name, value := range read {
		read[name] = canonicalAttributeValue(value)
	}
	return read
}

func canonicalAttributeValue(value types.AttributeValue) types.AttributeValue {
	switch value := value.(type) {
	case *types.AttributeValueMemberSS:
		slices.Sort(value.Value)
	case *types.AttributeValueMemberNS:
		slices.Sort(value.Value)
	case *types.AttributeValueMemberBS:
		slices.SortFunc(value.Value, func(x, y []byte) int { return bytes.Compare(x, y) })
	case *types.AttributeValueMemberL:
		for i, element := range value.Value {
			value.Value[i] = canonicalAttributeValue(element)
		}
	case *types.AttributeValueMemberM:
		for name, element := range value.Value {
			value.Value[name] = canonicalAttributeValue(element)
		}
	}
	return value
}

func mustJSON(i item) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// unmarshalJSON decodes JSON strictly enough that a tape's typo fails the test
// rather than answering something else.
func unmarshalJSON(b []byte, v any) error {
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
