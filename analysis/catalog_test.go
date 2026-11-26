package analysis

import (
	"context"
	"errors"
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestCatalogMetadataLoader(t *testing.T) {
	if (Catalog{}).MetadataLoader() != nil {
		t.Fatal("a catalog that cannot look databases up has no loader")
	}
	failure := errors.New("unreachable")
	load := Catalog{GetDatabase: func(_ context.Context, name string) (*Database, error) {
		switch name {
		case "db":
			return &Database{Name: "db", Metadata: &metadata.DatabaseSchemaMetadata{Name: "db"}}, nil
		case "broken":
			return nil, failure
		default:
			return nil, nil
		}
	}}.MetadataLoader()
	if m, err := load(context.Background(), "db"); err != nil || m.GetName() != "db" {
		t.Fatalf("db = %v, %v", m, err)
	}
	if m, err := load(context.Background(), "missing"); err != nil || m != nil {
		t.Fatalf("missing = %v, %v", m, err)
	}
	if _, err := load(context.Background(), "broken"); !errors.Is(err, failure) {
		t.Fatalf("broken err = %v", err)
	}
}
