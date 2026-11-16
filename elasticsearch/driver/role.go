package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/bytebase/omni/driver"
)

// securityRoute is the path that lists the instance's users: Elasticsearch's
// security API, or OpenSearch's, which nests it under its own plugin.
func (d *Driver) securityRoute() string {
	if d.openSearch {
		return "/_plugins/_security/api/internalusers"
	}
	return "/_security/user"
}

// getInstanceRoles reads the users/roles the endpoint reports, for display
// only. It fails open on any error, since most deployments do not use internal
// users. Users are ordered by name: the answer is an object keyed by them, and
// the metadata and its golden must hold one order.
func (d *Driver) getInstanceRoles(ctx context.Context) ([]*driver.Role, error) {
	resp, err := d.client.Do(ctx, http.MethodGet, d.securityRoute(), nil)
	if err != nil {
		slog.Debug("skip the instance's users: the endpoint answered no user list", "error", err)
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, nil
	}
	defer resp.Body.Close()

	var results map[string]struct {
		Roles []string `json:"roles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, nil
	}

	names := make([]string, 0, len(results))
	for name := range results {
		names = append(names, name)
	}
	slices.Sort(names)

	var roles []*driver.Role
	for _, name := range names {
		attribute := fmt.Sprintf("[%s]", strings.Join(results[name].Roles, ", "))
		roles = append(roles, &driver.Role{Name: name, Attribute: &attribute})
	}
	return roles, nil
}
