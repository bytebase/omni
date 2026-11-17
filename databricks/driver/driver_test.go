package driver

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/databricks/databricks-sdk-go/apierr"

	"github.com/bytebase/omni/driver"
)

// TestWorkspaceConfig maps the contract's settings onto the SDK's. None of it
// needs a workspace: Open builds the client from these values, and Ping is what
// reaches the service.
func TestWorkspaceConfig(t *testing.T) {
	t.Run("the workspace, token and warehouse are kept whole", func(t *testing.T) {
		cfg, warehouseID, err := workspaceConfig(driver.Config{
			Host:     "https://dbc-00000000-0000.cloud.databricks.com/",
			Password: "dapi-token",
			Params:   map[string]string{ParamWarehouseID: "warehouse-1"},
		})
		if err != nil {
			t.Fatalf("workspaceConfig: %v", err)
		}
		if cfg.Host != "https://dbc-00000000-0000.cloud.databricks.com" {
			t.Errorf("Host = %q", cfg.Host)
		}
		if cfg.Token != "dapi-token" {
			t.Errorf("Token = %q", cfg.Token)
		}
		if warehouseID != "warehouse-1" || cfg.WarehouseID != "warehouse-1" {
			t.Errorf("warehouse = %q / %q, want warehouse-1", warehouseID, cfg.WarehouseID)
		}
		if cfg.HTTPTransport != nil {
			t.Error("HTTPTransport is set without a dialer or TLS")
		}
		if cfg.HostMetadataResolver == nil {
			t.Fatal("HostMetadataResolver is not installed, so Open would fetch host metadata")
		}
		if meta, err := cfg.HostMetadataResolver(context.Background(), cfg.Host); err != nil || meta != nil {
			t.Errorf("HostMetadataResolver = %+v, %v; want no metadata and no error", meta, err)
		}
	})

	t.Run("warehouse_id is accepted as an alias", func(t *testing.T) {
		_, warehouseID, err := workspaceConfig(driver.Config{
			Host:     "dbc.example.com",
			Password: "token",
			Params:   map[string]string{"WAREHOUSE_ID": "warehouse-2"},
		})
		if err != nil {
			t.Fatalf("workspaceConfig: %v", err)
		}
		if warehouseID != "warehouse-2" {
			t.Errorf("warehouse = %q, want warehouse-2", warehouseID)
		}
	})

	t.Run("TLS is asked for the workspace host", func(t *testing.T) {
		asked := ""
		cfg, _, err := workspaceConfig(driver.Config{
			Host:     "https://dbc.example.com/",
			Password: "token",
			Params:   map[string]string{ParamWarehouseID: "warehouse-1"},
			TLS: func(host string) (*tls.Config, error) {
				asked = host
				return &tls.Config{ServerName: host}, nil
			},
		})
		if err != nil {
			t.Fatalf("workspaceConfig: %v", err)
		}
		if asked != "dbc.example.com" {
			t.Errorf("TLS was asked for %q, want dbc.example.com", asked)
		}
		transport, ok := cfg.HTTPTransport.(*http.Transport)
		if !ok {
			t.Fatalf("HTTPTransport = %T, want *http.Transport", cfg.HTTPTransport)
		}
		if transport.TLSClientConfig == nil || transport.TLSClientConfig.ServerName != "dbc.example.com" {
			t.Errorf("TLSClientConfig = %+v", transport.TLSClientConfig)
		}
	})

	t.Run("the caller's dialer reaches the SDK", func(t *testing.T) {
		cfg, _, err := workspaceConfig(driver.Config{
			Host:     "dbc.example.com",
			Password: "token",
			Params:   map[string]string{ParamWarehouseID: "warehouse-1"},
			Dial: func(context.Context, string, string) (net.Conn, error) {
				return nil, errors.New("dialed")
			},
		})
		if err != nil {
			t.Fatalf("workspaceConfig: %v", err)
		}
		transport, ok := cfg.HTTPTransport.(*http.Transport)
		if !ok {
			t.Fatalf("HTTPTransport = %T, want *http.Transport", cfg.HTTPTransport)
		}
		if transport.DialContext == nil {
			t.Fatal("DialContext is not the caller's dialer")
		}
		if _, err := transport.DialContext(context.Background(), "tcp", "dbc.example.com:443"); err == nil || !strings.Contains(err.Error(), "dialed") {
			t.Errorf("DialContext: %v, want the caller's error", err)
		}
	})
}

// TestWorkspaceConfigRefusesWhatCannotNameAWorkspace covers the values Open
// refuses rather than handing the SDK a client that cannot authenticate.
func TestWorkspaceConfigRefusesWhatCannotNameAWorkspace(t *testing.T) {
	for _, tc := range []struct {
		name    string
		cfg     driver.Config
		wantErr string
	}{
		{
			name:    "no host",
			cfg:     driver.Config{Password: "token", Params: map[string]string{ParamWarehouseID: "w"}},
			wantErr: "host must be set",
		},
		{
			name:    "no token",
			cfg:     driver.Config{Host: "dbc.example.com", Params: map[string]string{ParamWarehouseID: "w"}},
			wantErr: "token must be set",
		},
		{
			name:    "no warehouse",
			cfg:     driver.Config{Host: "dbc.example.com", Password: "token"},
			wantErr: "warehouse ID must be set",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := workspaceConfig(tc.cfg)
			if err == nil {
				t.Fatal("workspaceConfig accepted the configuration")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("error = %q, want it to hold %q", err, tc.wantErr)
			}
		})
	}
}

func TestOpenRefusesAnIncompleteConfiguration(t *testing.T) {
	if _, err := Open(context.Background(), driver.Config{}); err == nil {
		t.Fatal("Open accepted an empty configuration")
	}
}

// TestWrapError keeps the SDK's own error reachable and joins the contract's
// not-found marker when the workspace answered 404.
func TestWrapError(t *testing.T) {
	apiErr := &apierr.APIError{StatusCode: 404, ErrorCode: "RESOURCE_DOES_NOT_EXIST", Message: "no such catalog"}
	err := wrapError("failed to read a catalog", apiErr)
	if !errors.Is(err, driver.ErrNotFound) {
		t.Errorf("wrapError(404) = %v, want driver.ErrNotFound", err)
	}
	var got *apierr.APIError
	if !errors.As(err, &got) {
		t.Errorf("wrapError(404) = %v, want the SDK's error to survive", err)
	}

	plain := errors.New("boom")
	wrapped := wrapError("failed to ping", plain)
	if !errors.Is(wrapped, plain) {
		t.Errorf("wrapError = %v, want the original error", wrapped)
	}
	if !strings.Contains(wrapped.Error(), "failed to ping") {
		t.Errorf("wrapError = %v, want the operation", wrapped)
	}
}

func TestHostnameOf(t *testing.T) {
	for _, tc := range []struct {
		host string
		want string
	}{
		{host: "https://dbc.example.com/", want: "dbc.example.com"},
		{host: "dbc.example.com", want: "dbc.example.com"},
		{host: "https://dbc.example.com:443/path", want: "dbc.example.com"},
	} {
		got, err := hostnameOf(tc.host)
		if err != nil {
			t.Errorf("hostnameOf(%q): %v", tc.host, err)
			continue
		}
		if got != tc.want {
			t.Errorf("hostnameOf(%q) = %q, want %q", tc.host, got, tc.want)
		}
	}
	if _, err := hostnameOf("https:///path"); err == nil {
		t.Error("hostnameOf accepted a URL with no hostname")
	}
}
