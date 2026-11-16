package testcontainer

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest"
)

func TestPin(t *testing.T) {
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		t.Fatal(err)
	}
	if pin.Image != "clickhouse/clickhouse-server:"+pin.Tag {
		t.Errorf("image = %q, tag = %q", pin.Image, pin.Tag)
	}
	if !strings.HasPrefix(pin.Digest, "sha256:") {
		t.Errorf("digest = %q", pin.Digest)
	}
}

func TestDatabase(t *testing.T) {
	for _, tc := range []struct {
		test, want string
	}{
		{"TestSyncDatabase", "omni_testsyncdatabase"},
		{"TestExecuteAutoCommit", "omni_testexecuteautocommit"},
	} {
		got := Database(fakeTB{name: tc.test})
		if got != tc.want {
			t.Errorf("Database(%q) = %q, want %q", tc.test, got, tc.want)
		}
	}
}

type fakeTB struct {
	testing.TB
	name string
}

func (f fakeTB) Name() string { return f.name }
