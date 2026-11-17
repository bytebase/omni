package driver

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	metadatapb "github.com/bytebase/omni/metadata"
)

// TestDumpWritesNothing pins what Bytebase's BigQuery driver dumped: nothing, with
// no error, whatever metadata the caller passes.
func TestDumpWritesNothing(t *testing.T) {
	var out bytes.Buffer
	require.NoError(t, (&Driver{}).Dump(context.Background(), &out, &metadatapb.DatabaseSchemaMetadata{Name: "analytics"}))
	require.Empty(t, out.String())
}
