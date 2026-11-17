package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	"google.golang.org/api/option"

	"github.com/bytebase/omni/driver"
)

// ParamProject is the Config.Params key that names the Google Cloud project:
// BigQuery addresses every object as project.dataset.table, and the client
// library signs every request for a project.
const ParamProject = "project"

// clientConfig is what Open builds out of driver.Config before a client exists.
// It is a separate value so that the settings a caller supplies can be read
// without a project, a credential or a network.
type clientConfig struct {
	project string
	// endpoint is the BigQuery API endpoint the caller named, or empty for the
	// client library's own bigquery.googleapis.com.
	endpoint string
	// options are the settings above as client options, in the order they are
	// applied.
	options []option.ClientOption
}

// clientConfigOf reads the settings Open hands the client library.
//
// Project comes from Params[ParamProject] and is required: no BigQuery request
// can be built without one, and Config.Database is a dataset, not a project.
// Host and Port optionally name the endpoint, which is how a Private Service
// Connect endpoint or a proxy is reached; an empty Host leaves the library's
// own. Password is the credential, as the JSON document of a service account, an
// external account, an impersonated service account or an authorized user; empty
// leaves the Google authentication library's default chain, which reads the
// environment, the well-known file or an attached service account.
func clientConfigOf(cfg driver.Config) (*clientConfig, error) {
	project := strings.TrimSpace(cfg.Params[ParamProject])
	if project == "" {
		return nil, fmt.Errorf("%s must be set: BigQuery addresses every object as project.dataset.table", ParamProject)
	}
	c := &clientConfig{project: project, endpoint: endpointOf(cfg.Host, cfg.Port)}
	if c.endpoint != "" {
		c.options = append(c.options, option.WithEndpoint(c.endpoint))
	}
	if cfg.Password != "" {
		credential, err := credentialOption([]byte(cfg.Password))
		if err != nil {
			return nil, err
		}
		c.options = append(c.options, credential)
	}
	return c, nil
}

// endpointOf joins the host and port a caller named into an endpoint, or returns
// an empty endpoint for the client library's own. A host the caller already
// wrote with a port keeps it, and a bare host is left for the library to merge
// with its default scheme and path.
func endpointOf(host, port string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}
	if port = strings.TrimSpace(port); port != "" {
		return net.JoinHostPort(host, port)
	}
	return host
}

// credentialOption reads the "type" member of a Google Cloud credential JSON
// document and returns the client option that resolves it, as Bytebase's
// GCPCredentialOption did.
//
// SECURITY NOTE: an external account or an impersonated service account
// configuration is not validated here; a URL it names is fetched when the
// credential is used. See
// https://cloud.google.com/docs/authentication/external/externally-sourced-credentials.
func credentialOption(credentials []byte) (option.ClientOption, error) {
	var credential struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(credentials, &credential); err != nil {
		return nil, fmt.Errorf("failed to parse the GCP credential JSON: %w", err)
	}
	switch credential.Type {
	case "service_account":
		return option.WithAuthCredentialsJSON(option.ServiceAccount, credentials), nil
	case "external_account":
		return option.WithAuthCredentialsJSON(option.ExternalAccount, credentials), nil
	case "impersonated_service_account":
		return option.WithAuthCredentialsJSON(option.ImpersonatedServiceAccount, credentials), nil
	case "authorized_user":
		return option.WithAuthCredentialsJSON(option.AuthorizedUser, credentials), nil
	case "":
		return nil, errors.New(`the GCP credential JSON holds no "type" member`)
	default:
		return nil, fmt.Errorf("unsupported GCP credential type: %q", credential.Type)
	}
}
