package driver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
)

// TestDump renders each server's synced corpus, testdata/sync/<role>.json,
// and holds it to testdata/dump/<role>.sql.
func TestDump(t *testing.T) {
	for _, role := range versions {
		t.Run(role, func(t *testing.T) {
			m := readGolden(t, filepath.Join("testdata", "sync", role+".json"))
			before := proto.Clone(m)
			var b strings.Builder
			if err := (&Driver{}).Dump(context.Background(), &b, m); err != nil {
				t.Fatal(err)
			}
			if !proto.Equal(before, m) {
				t.Errorf("Dump changed the metadata it was given")
			}
			golden := filepath.Join("testdata", "dump", role+".sql")
			if *update {
				if err := os.WriteFile(golden, []byte(b.String()), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(string(want), b.String()); diff != "" {
				t.Errorf("Dump of %s differs from %s (-want +got):\n%s", role, golden, diff)
			}
		})
	}
}
