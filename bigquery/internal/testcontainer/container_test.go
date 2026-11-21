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

func TestResolveDryRunConfig(t *testing.T) {
	env := func(values map[string]string) func(string) string {
		return func(name string) string { return values[name] }
	}

	_, err := ResolveDryRunConfig(env(nil))
	require.ErrorContains(t, err, CredentialVariable)

	_, err = ResolveDryRunConfig(env(map[string]string{
		CredentialVariable: CredentialDefault,
	}))
	require.ErrorContains(t, err, ProjectVariable)

	cfg, err := ResolveDryRunConfig(env(map[string]string{
		CredentialVariable: CredentialDefault,
		ProjectVariable:    " bytebase-dev ",
	}))
	require.NoError(t, err)
	require.Equal(t, LiveConfig{Project: "bytebase-dev"}, cfg)

	document := `{"type":"service_account","project_id":"key-project"}`
	cfg, err = ResolveDryRunConfig(env(map[string]string{
		CredentialVariable: document,
	}))
	require.NoError(t, err)
	require.Equal(t, LiveConfig{Project: "key-project", Password: document}, cfg)
}

func TestDryRunDigest(t *testing.T) {
	got := DryRunDigest("bytebase-dev", "US")
	require.Equal(t, "sha256:fef5b54092b50edbb5d7c363182e0cd578a2147d527595e968a26f5131429231", got)
	require.NotEqual(t, got, DryRunDigest("bytebase-dev", "EU"))
	require.NotEqual(t, got, DryRunDigest("other-project", "US"))
}

func TestCredentialOption(t *testing.T) {
	for _, typ := range []string{"service_account", "external_account", "impersonated_service_account", "authorized_user"} {
		opt, err := credentialOption([]byte(`{"type":"` + typ + `"}`))
		require.NoError(t, err, typ)
		require.NotNil(t, opt, typ)
	}
	_, err := credentialOption([]byte(`not-json`))
	require.ErrorContains(t, err, "failed to parse")
	_, err = credentialOption([]byte(`{}`))
	require.ErrorContains(t, err, `"type"`)
	_, err = credentialOption([]byte(`{"type":"unknown"}`))
	require.ErrorContains(t, err, "unsupported")
}
