package driver

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// scriptedClient answers each Page call with the next answer of a script, and
// holds the requests it was asked.
type scriptedClient struct {
	answers  []PageAnswer
	errs     []error
	requests []PageRequest
}

func (*scriptedClient) Ping(context.Context) error { return nil }

func (*scriptedClient) Close() error { return nil }

func (*scriptedClient) ListDatabases(context.Context) ([]string, error) { return nil, nil }

func (*scriptedClient) ListContainers(context.Context, string) ([]string, error) {
	return nil, nil
}

func (c *scriptedClient) Page(_ context.Context, req PageRequest) (PageAnswer, error) {
	c.requests = append(c.requests, req)
	i := len(c.requests) - 1
	if i >= len(c.answers) {
		panic("scriptedClient: more Pages than answers")
	}
	return c.answers[i], c.errs[i]
}

// gatewayRefusal is the answer the gateway gives a cross-partition query it
// cannot merge itself.
func gatewayRefusal() error {
	body := `{"code":"BadRequest","message":"The provided cross partition query can not be directly served by the gateway."}`
	return azruntime.NewResponseErrorWithErrorCode(&http.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 BadRequest",
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(body)),
	}, "BadRequest")
}

// notFound is any other failure: a 404 is a container that is not there, which
// no engine can serve.
func notFound() error {
	return azruntime.NewResponseErrorWithErrorCode(&http.Response{
		StatusCode: http.StatusNotFound,
		Status:     "404 Not Found",
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(`{"code":"NotFound"}`)),
	}, "NotFound")
}

// TestQueryEveryPartition pins when a query runs again through the client's
// query engine: only when the gateway refused it outright. A query that already
// read a page cannot resume from the engine, and its error is its answer.
func TestQueryEveryPartition(t *testing.T) {
	items := [][]byte{[]byte(`{"seq":1}`)}

	t.Run("a query the gateway serves", func(t *testing.T) {
		c := &scriptedClient{answers: []PageAnswer{{Items: items, Pages: 2}}, errs: []error{nil}}
		answer, err := queryEveryPartition(context.Background(), c, PageRequest{Statement: "SELECT * FROM c"})
		require.NoError(t, err)
		require.Equal(t, PageAnswer{Items: items, Pages: 2}, answer)
		require.Len(t, c.requests, 1)
		require.False(t, c.requests[0].Engine)
	})

	t.Run("a refusal before any page runs again through the engine", func(t *testing.T) {
		c := &scriptedClient{
			answers: []PageAnswer{{}, {Items: items, Pages: 1}},
			errs:    []error{gatewayRefusal(), nil},
		}
		answer, err := queryEveryPartition(context.Background(), c, PageRequest{Statement: "SELECT * FROM c"})
		require.NoError(t, err)
		require.Equal(t, PageAnswer{Items: items, Pages: 1}, answer)
		require.Len(t, c.requests, 2)
		require.False(t, c.requests[0].Engine)
		require.True(t, c.requests[1].Engine)
	})

	t.Run("a refusal after a page is the answer", func(t *testing.T) {
		c := &scriptedClient{
			answers: []PageAnswer{{Items: items, Pages: 1}},
			errs:    []error{gatewayRefusal()},
		}
		answer, err := queryEveryPartition(context.Background(), c, PageRequest{Statement: "SELECT * FROM c"})
		require.Error(t, err)
		require.Equal(t, PageAnswer{Items: items, Pages: 1}, answer, "the pages read are kept")
		require.Len(t, c.requests, 1, "a query that read a page cannot resume through the engine")
	})

	t.Run("another failure before any page is the answer", func(t *testing.T) {
		c := &scriptedClient{answers: []PageAnswer{{}}, errs: []error{notFound()}}
		_, err := queryEveryPartition(context.Background(), c, PageRequest{Statement: "SELECT * FROM c"})
		require.Error(t, err)
		require.Len(t, c.requests, 1, "only the gateway's cross-partition refusal runs again")
	})

	t.Run("an error that is not the gateway's is the answer", func(t *testing.T) {
		c := &scriptedClient{answers: []PageAnswer{{}}, errs: []error{errors.New("connection reset")}}
		_, err := queryEveryPartition(context.Background(), c, PageRequest{Statement: "SELECT * FROM c"})
		require.Error(t, err)
		require.Len(t, c.requests, 1)
	})
}
