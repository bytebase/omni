package driver

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/driver"
)

// TestClientConfigBuildsTheConnection covers the settings Open hands the client
// library: which ones are required, which have a default, and how a caller's
// endpoint and credential reach it. None of it needs a project.
func TestClientConfigBuildsTheConnection(t *testing.T) {
	for _, tc := range []struct {
		name         string
		cfg          driver.Config
		wantErr      string
		wantProject  string
		wantEndpoint string
		wantOptions  int
	}{
		{
			name:        "the project is required",
			cfg:         driver.Config{Database: "analytics"},
			wantErr:     "project must be set",
			wantOptions: 0,
		},
		{
			name:        "a project alone is enough, and the endpoint is the library's own",
			cfg:         driver.Config{Params: map[string]string{ParamProject: "example-project"}},
			wantProject: "example-project",
		},
		{
			name: "the project is trimmed, and a dataset is not a project",
			cfg: driver.Config{
				Database: "analytics",
				Params:   map[string]string{ParamProject: "  example-project  "},
			},
			wantProject: "example-project",
		},
		{
			name: "a host and port name the endpoint",
			cfg: driver.Config{
				Host:   "bigquery.example.com",
				Port:   "443",
				Params: map[string]string{ParamProject: "example-project"},
			},
			wantProject:  "example-project",
			wantEndpoint: "bigquery.example.com:443",
			wantOptions:  1,
		},
		{
			name: "a host that already carries its port is kept whole",
			cfg: driver.Config{
				Host:   "bigquery.example.com:8443",
				Params: map[string]string{ParamProject: "example-project"},
			},
			wantProject:  "example-project",
			wantEndpoint: "bigquery.example.com:8443",
			wantOptions:  1,
		},
		{
			name: "a credential and an endpoint are both applied",
			cfg: driver.Config{
				Host:     "bigquery.example.com",
				Password: `{"type": "service_account", "project_id": "example-project"}`,
				Params:   map[string]string{ParamProject: "example-project"},
			},
			wantProject:  "example-project",
			wantEndpoint: "bigquery.example.com",
			wantOptions:  2,
		},
		{
			name: "a credential the library cannot read is refused before a client exists",
			cfg: driver.Config{
				Password: `{"type": "not_a_credential"}`,
				Params:   map[string]string{ParamProject: "example-project"},
			},
			wantErr: `unsupported GCP credential type: "not_a_credential"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := clientConfigOf(tc.cfg)
			if tc.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.wantProject, got.project)
			require.Equal(t, tc.wantEndpoint, got.endpoint)
			require.Len(t, got.options, tc.wantOptions)
		})
	}
}

// TestCredentialOption covers the four credential documents Google Cloud states
// and the three ways one is refused, which is what Bytebase's GCPCredentialOption
// answered for the same input.
func TestCredentialOption(t *testing.T) {
	for _, tc := range []struct {
		name        string
		credentials string
		wantErr     string
	}{
		{name: "service account", credentials: `{"type": "service_account", "project_id": "test"}`},
		{name: "external account", credentials: `{"type": "external_account", "audience": "test"}`},
		{name: "impersonated service account", credentials: `{"type": "impersonated_service_account"}`},
		{name: "authorized user", credentials: `{"type": "authorized_user"}`},
		{
			name:        "not JSON",
			credentials: "not json",
			wantErr:     "failed to parse the GCP credential JSON",
		},
		{
			name:        "no type member",
			credentials: `{"project_id": "test"}`,
			wantErr:     `holds no "type" member`,
		},
		{
			name:        "a type Google Cloud does not state",
			credentials: `{"type": "unknown_type"}`,
			wantErr:     `unsupported GCP credential type: "unknown_type"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := credentialOption([]byte(tc.credentials))
			if tc.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}

// TestEndpointOf covers the three shapes a caller's endpoint takes.
func TestEndpointOf(t *testing.T) {
	for _, tc := range []struct {
		host, port, want string
	}{
		{host: "", port: "443", want: ""},
		{host: "bigquery.example.com", want: "bigquery.example.com"},
		{host: "bigquery.example.com", port: "443", want: "bigquery.example.com:443"},
		{host: "  bigquery.example.com  ", port: " 443 ", want: "bigquery.example.com:443"},
		{host: "bigquery.example.com:8443", port: "", want: "bigquery.example.com:8443"},
	} {
		require.Equal(t, tc.want, endpointOf(tc.host, tc.port), "endpointOf(%q, %q)", tc.host, tc.port)
	}
}

// TestClientConfigRefusesACredentialButKeepsTheProject pins that a refused
// credential is not reported as a missing project, which a caller reads
// differently.
func TestClientConfigRefusesACredentialButKeepsTheProject(t *testing.T) {
	_, err := clientConfigOf(driver.Config{
		Password: "{}",
		Params:   map[string]string{ParamProject: "example-project"},
	})
	require.Error(t, err)
	require.False(t, strings.Contains(err.Error(), ParamProject+" must be set"))
}
