package driver

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"

	"github.com/bytebase/omni/driver"
)

// TestNotFoundIsMappedToTheContract covers the two shapes BigQuery says "not
// found" in: the API's 404, and a job the service failed with its notFound reason,
// which is how a statement against a missing dataset fails. Everything else stays
// the client's own error.
func TestNotFoundIsMappedToTheContract(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "the API's 404",
			err:  &googleapi.Error{Code: http.StatusNotFound, Message: "Not found: Dataset example-project:analytics"},
			want: true,
		},
		{
			name: "a job failed with the notFound reason",
			err:  &bigquery.Error{Reason: "notFound", Message: "Not found: Dataset example-project:analytics"},
			want: true,
		},
		{
			name: "an answer wrapped by the driver",
			err:  fmt.Errorf("failed to read columns: %w", &googleapi.Error{Code: http.StatusNotFound}),
			want: true,
		},
		{
			name: "the API's 403",
			err:  &googleapi.Error{Code: http.StatusForbidden, Message: "Access Denied: Dataset example-project:analytics"},
		},
		{
			name: "a job failed with another reason",
			err:  &bigquery.Error{Reason: "invalidQuery", Message: "Syntax error"},
		},
		{
			name: "an error that is not the client's",
			err:  errors.New("connection refused"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, notFound(tc.err))
			mapped := notFoundError(tc.err)
			if tc.want {
				require.ErrorIs(t, mapped, driver.ErrNotFound)
				// The service's own error stays readable behind the contract's.
				require.ErrorIs(t, mapped, tc.err)
				return
			}
			require.NotErrorIs(t, mapped, driver.ErrNotFound)
			require.Equal(t, tc.err, mapped)
		})
	}
}
