package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// databaseName is the one database an Elasticsearch cluster is reported as.
// Bytebase named it "node" and left the naming to a later question.
const databaseName = "node"

// SyncInstance reads the server version, the cluster's indices and its users.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, err := d.getVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch version from Elasticsearch server: %w", err)
	}

	dbMetadata, err := d.SyncDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch indices from Elasticsearch server: %w", err)
	}

	instanceRoles, err := d.getInstanceRoles(ctx)
	if err != nil {
		return nil, err
	}

	return &driver.Instance{
		Version:   version,
		Databases: []*metadatapb.DatabaseSchemaMetadata{dbMetadata},
		Roles:     instanceRoles,
	}, nil
}

// SyncDatabase reads the cluster's indices, with the document count and data
// size each reports, and their aliases. Indices are ordered by name: the cat
// API's order is the server's, and the metadata and its golden must hold one
// order.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	indices, err := d.getIndices(ctx)
	if err != nil {
		return nil, err
	}
	aliases, err := d.getAliases(ctx, indices)
	if err != nil {
		return nil, err
	}
	return &metadatapb.DatabaseSchemaMetadata{
		Name: databaseName,
		Schemas: []*metadatapb.SchemaMetadata{{
			Tables:   indices,
			Synonyms: aliases,
		}},
	}, nil
}

// aliasesRoute lists every index's aliases.
const aliasesRoute = "/_alias"

// aliasesResult is the part of one index's _alias entry the sync reads.
type aliasesResult struct {
	Aliases map[string]json.RawMessage `json:"aliases"`
}

// getAliases reads the aliases of the synced indices. An alias is another
// name a request may read an index by, which the engine resolves when it is
// used, so each is a synonym: one per index the alias names, ordered by
// alias and then index. An alias of an index the sync leaves out (a hidden
// one, a data stream's) is left out with it.
func (d *Driver) getAliases(ctx context.Context, indices []*metadatapb.TableMetadata) ([]*metadatapb.SynonymMetadata, error) {
	resp, err := d.client.Do(ctx, http.MethodGet, aliasesRoute, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list aliases: %w", err)
	}
	body, err := readBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list aliases: %w", err)
	}
	var results map[string]aliasesResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse Elasticsearch aliases response: %s: %w", truncate(body, 500), err)
	}
	var synonyms []*metadatapb.SynonymMetadata
	for _, index := range indices {
		for alias := range results[index.Name].Aliases {
			synonyms = append(synonyms, &metadatapb.SynonymMetadata{Name: alias, TargetName: index.Name})
		}
	}
	slices.SortFunc(synonyms, func(a, b *metadatapb.SynonymMetadata) int {
		if c := strings.Compare(a.Name, b.Name); c != 0 {
			return c
		}
		return strings.Compare(a.TargetName, b.TargetName)
	})
	return synonyms, nil
}

// versionResult is the part of the root document getVersion reads.
type versionResult struct {
	Version struct {
		Number string `json:"number"`
	} `json:"version"`
}

// getVersion reads the endpoint's version from its root document.
func (d *Driver) getVersion(ctx context.Context) (string, error) {
	resp, err := d.client.Do(ctx, http.MethodGet, "/", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result versionResult
	if err := json.Unmarshal(body, &result); err != nil {
		// Include response body to help debug parsing issues
		return "", fmt.Errorf("failed to parse version response: %s: %w", truncate(body, 500), err)
	}
	return result.Version.Number, nil
}

// indicesResult is the part of one _cat/indices entry the sync reads.
type indicesResult struct {
	IndexSize string `json:"store.size"`
	DocsCount string `json:"docs.count"`
	Index     string `json:"index"`
}

// indexSettingsResult is the part of one index's settings the sync reads.
type indexSettingsResult struct {
	Settings struct {
		Index struct {
			Hidden string `json:"hidden"`
		} `json:"index"`
	} `json:"settings"`
}

// hiddenSettingsRoute asks every index's settings for its hidden flag. "all"
// reaches hidden indices from Elasticsearch 7.7 on, and is a value 6.x knows
// too, where "hidden" is refused (no valid expand wildcard value [hidden]) and
// no index is hidden.
//
// Ref: https://www.elastic.co/guide/en/elasticsearch/reference/7.17/multi-index.html (expand_wildcards)
const hiddenSettingsRoute = "/*/_settings?filter_path=**.hidden&expand_wildcards=all"

// getHiddenIndices returns the index names that are marked as hidden.
func (d *Driver) getHiddenIndices(ctx context.Context) (map[string]bool, error) {
	hiddenIndices := map[string]bool{}
	resp, err := d.client.Do(ctx, http.MethodGet, hiddenSettingsRoute, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get index settings: %w", err)
	}
	body, err := readBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get index settings: %w", err)
	}
	var results map[string]indexSettingsResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse index settings response: %w", err)
	}
	for indexName, indexData := range results {
		if indexData.Settings.Index.Hidden == "true" {
			hiddenIndices[indexName] = true
		}
	}
	return hiddenIndices, nil
}

// catIndicesRoute lists every index the cat API reports.
const catIndicesRoute = "/_cat/indices?format=json&pretty"

// getIndices reads every index's name, document count and data size, hidden
// ones left out.
func (d *Driver) getIndices(ctx context.Context) ([]*metadatapb.TableMetadata, error) {
	// Get hidden indices to filter them out.
	hiddenIndices, err := d.getHiddenIndices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get hidden indices: %w", err)
	}

	resp, err := d.client.Do(ctx, http.MethodGet, catIndicesRoute, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list indices: %w", err)
	}
	body, err := readBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list indices: %w", err)
	}
	var results []indicesResult
	if err := json.Unmarshal(body, &results); err != nil {
		// Include response body to help debug parsing issues
		return nil, fmt.Errorf("failed to parse Elasticsearch indices response: %s: %w", truncate(body, 500), err)
	}
	slices.SortFunc(results, func(a, b indicesResult) int {
		return strings.Compare(a.Index, b.Index)
	})

	var indicesMetadata []*metadatapb.TableMetadata
	for _, m := range results {
		if hiddenIndices[m.Index] {
			continue
		}
		dataSize, err := unitConversion(m.IndexSize)
		if err != nil {
			return nil, err
		}
		docCount, err := parseDocCount(m.DocsCount)
		if err != nil {
			return nil, err
		}
		indicesMetadata = append(indicesMetadata, &metadatapb.TableMetadata{
			Name:     m.Index,
			DataSize: dataSize,
			RowCount: docCount,
		})
	}
	return indicesMetadata, nil
}

// parseDocCount parses document count string, returning 0 for empty values
// (e.g., closed indices).
func parseDocCount(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	count, err := strconv.Atoi(s)
	return int64(count), err
}

var sizeRe = regexp.MustCompile(`([0-9.]+)([kmgtpe]?b)`)

func unitConversion(sizeWithUnit string) (int64, error) {
	// Empty size is expected for closed indices (Elasticsearch answers null,
	// which unmarshals to "").
	if sizeWithUnit == "" {
		return 0, nil
	}

	sizeWithUnit = strings.ToLower(sizeWithUnit)
	match := sizeRe.FindSubmatch([]byte(sizeWithUnit))

	// Non-empty string that doesn't match expected format is unexpected
	if len(match) < 3 {
		return 0, fmt.Errorf("unexpected size format: %q", sizeWithUnit)
	}

	unit := string(match[2])

	size, err := strconv.ParseFloat(string(match[1]), 64)
	if err != nil {
		return 0, err
	}

	switch unit {
	case "kb":
		size *= 1024
	case "mb":
		size *= 1024 * 1024
	case "gb":
		size *= 1024 * 1024 * 1024
	case "tb":
		size *= 1024 * 1024 * 1024 * 1024
	case "pb":
		size *= 1024 * 1024 * 1024 * 1024 * 1024
	case "eb":
		size *= 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		// For "b" (bytes) or any other unit, keep the size as-is
	}

	return int64(size), nil
}

// readBody reads and closes a response, failing on an error status with the
// body the endpoint sent.
func readBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// truncate cuts text to limit code points and says so, so an unreadable
// answer's error carries a usable preview of it.
func truncate(text []byte, limit int) string {
	s := string(text)
	count := 0
	for i := range s {
		if count >= limit {
			return s[:i] + "..."
		}
		count++
	}
	return s
}
