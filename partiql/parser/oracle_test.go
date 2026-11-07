package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bytebase/omni/partiql/internal/testcontainer"
)

var dynamoOracle testcontainer.Container

func TestMain(m *testing.M) {
	code := m.Run()
	_ = testcontainer.TerminateContainer(dynamoOracle)
	os.Exit(code)
}

func dynamoRequest(t *testing.T, endpoint, operation string, body any) (int, []byte) {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-amz-json-1.0")
	req.Header.Set("X-Amz-Target", "DynamoDB_20120810."+operation)
	req.Header.Set("X-Amz-Date", time.Now().UTC().Format("20060102T150405Z"))
	// DynamoDB Local accepts a dummy access key and signature.
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=omni/20261107/us-east-1/dynamodb/aws4_request, SignedHeaders=host;x-amz-date;x-amz-target, Signature=local")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, result
}

// DynamoDB implements a subset of PartiQL. These probes cover that subset.
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/ql-reference.html
func TestDynamoDBPartiQLOracle(t *testing.T) {
	var err error
	dynamoOracle, err = testcontainer.Start(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	endpoint, err := dynamoOracle.PortEndpoint(context.Background(), "8000/tcp", "http")
	if err != nil {
		t.Fatal(err)
	}
	status, output := dynamoRequest(t, endpoint, "CreateTable", map[string]any{
		"TableName": "omni", "BillingMode": "PAY_PER_REQUEST",
		"AttributeDefinitions": []any{map[string]string{"AttributeName": "id", "AttributeType": "S"}},
		"KeySchema":            []any{map[string]string{"AttributeName": "id", "KeyType": "HASH"}},
	})
	if status != http.StatusOK && !bytes.Contains(output, []byte("ResourceInUseException")) {
		t.Fatalf("fixture: HTTP %d: %s", status, output)
	}
	for _, tc := range []struct {
		sql    string
		accept bool
	}{
		{"SELECT * FROM omni", true},
		{"INSERT INTO omni VALUE {'id':'one', 'v':1}", true},
		{"SELECT id FROM omni WHERE id = 'one'", true},
		{"UPDATE omni SET v = 2 WHERE id = 'one'", true},
		{"DELETE FROM omni WHERE id = 'one'", true},
		{"SELECT FROM", false},
		{"SELECT * FROM omni WHERE", false},
		{"INSERT INTO omni VALUE", false},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			_, parseErr := Parse(tc.sql)
			status, body := dynamoRequest(t, endpoint, "ExecuteStatement", map[string]string{"Statement": tc.sql})
			if status != http.StatusOK && status != http.StatusBadRequest {
				t.Fatalf("no query verdict: HTTP %d: %s", status, body)
			}
			if (parseErr == nil) != tc.accept || (status == http.StatusOK) != tc.accept {
				t.Errorf("parser error %v, DynamoDB HTTP %d: %s", parseErr, status, body)
			}
		})
	}
}
