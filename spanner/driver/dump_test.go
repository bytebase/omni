package driver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDump writes the database's DDL as the server states it and holds it to
// testdata/dump/fixture.sql.
func TestDump(t *testing.T) {
	a := require.New(t)
	d := openFixture(t, "dump")
	defer d.Close()

	var dumped strings.Builder
	a.NoError(d.Dump(context.Background(), &dumped, nil))

	golden := filepath.Join("testdata", "dump", "fixture.sql")
	if *update {
		a.NoError(os.WriteFile(golden, []byte(dumped.String()), 0o644))
		return
	}
	want, err := os.ReadFile(golden)
	a.NoError(err)
	a.Equal(string(want), dumped.String())
}

// TestDumpMissing holds a database the server does not have to the server's own
// error: a dump is a read of the schema, not a sync that reports a missing
// database.
func TestDumpMissing(t *testing.T) {
	a := require.New(t)
	instance, _ := syncFixture(t)
	d := openDriver(t, instance, "missing")
	defer d.Close()

	var dumped strings.Builder
	err := d.Dump(context.Background(), &dumped, nil)
	a.Error(err)
	a.Contains(err.Error(), "NotFound")
}
