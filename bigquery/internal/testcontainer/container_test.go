package testcontainer

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestResolveLiveConfig(t *testing.T) {
	env := func(values map[string]string) func(string) string {
		return func(name string) string { return values[name] }
	}

	_, err := ResolveLiveConfig(env(nil))
	require.ErrorContains(t, err, CredentialVariable)

	_, err = ResolveLiveConfig(env(map[string]string{CredentialVariable: "  "}))
	require.ErrorContains(t, err, CredentialVariable)

	_, err = ResolveLiveConfig(env(map[string]string{
		CredentialVariable: CredentialDefault,
		DatasetVariable:    "omni_record",
	}))
	require.ErrorContains(t, err, ProjectVariable)

	cfg, err := ResolveLiveConfig(env(map[string]string{
		CredentialVariable: CredentialDefault,
		ProjectVariable:    "  bytebase-dev ",
		DatasetVariable:    " omni_record ",
	}))
	require.NoError(t, err)
	require.Equal(t, LiveConfig{
		Project: "bytebase-dev",
		Dataset: "omni_record",
	}, cfg)

	document := `{"type":"service_account","project_id":"key-project"}`
	cfg, err = ResolveLiveConfig(env(map[string]string{
		CredentialVariable: document,
		DatasetVariable:    "omni_record",
	}))
	require.NoError(t, err)
	require.Equal(t, LiveConfig{
		Project:  "key-project",
		Dataset:  "omni_record",
		Password: document,
	}, cfg)

	keyPath := filepath.Join(t.TempDir(), "key.json")
	require.NoError(t, os.WriteFile(keyPath, []byte(document), 0o600))
	cfg, err = ResolveLiveConfig(env(map[string]string{
		CredentialVariable: keyPath,
		ProjectVariable:    "override-project",
		DatasetVariable:    "omni_record",
	}))
	require.NoError(t, err)
	require.Equal(t, LiveConfig{
		Project:  "override-project",
		Dataset:  "omni_record",
		Password: document,
	}, cfg)

	_, err = ResolveLiveConfig(env(map[string]string{
		CredentialVariable: filepath.Join(t.TempDir(), "missing.json"),
		ProjectVariable:    "bytebase-dev",
		DatasetVariable:    "omni_record",
	}))
	require.ErrorContains(t, err, "cannot be read")

	_, err = ResolveLiveConfig(env(map[string]string{
		CredentialVariable: CredentialDefault,
		ProjectVariable:    "bytebase-dev",
	}))
	require.ErrorContains(t, err, DatasetVariable)
}

func TestDatasetDigest(t *testing.T) {
	created := time.UnixMilli(1700000000123)
	got := DatasetDigest("bytebase-dev:omni_record", "US", created)
	require.True(t, len(got) == len("sha256:")+64)
	require.Equal(t, got, DatasetDigest("bytebase-dev:omni_record", "US", created))
	require.NotEqual(t, got, DatasetDigest("bytebase-dev:omni_record", "US", time.UnixMilli(1700000000124)))
	require.Equal(t, "gcp:bigquery:bytebase-dev:US", DatasetImage("bytebase-dev", "US"))
}
