package driver

import (
	"context"
	"io"

	metadatapb "github.com/bytebase/omni/metadata"
)

// Dump writes nothing and returns nil. The package doc says why Redis has no
// schema for a dump to render.
func (*Driver) Dump(_ context.Context, _ io.Writer, _ *metadatapb.DatabaseSchemaMetadata) error {
	return nil
}
