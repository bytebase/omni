package parser

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oracle/ast"
)

func TestOracleLocNodeCoverage(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "loc_node_coverage.tsv"))
	manifest := make(map[string]coverageRow, len(rows))
	statusCounts := make(map[string]int)
	for _, row := range rows {
		status := coverageStatus(row)
		if _, ok := allowedCoverageStatuses[status]; !ok {
			t.Fatalf("%s: unknown status %q", row.Key, status)
		}
		if status == "unknown" {
			t.Fatalf("%s: Loc node status must be classified, got unknown", row.Key)
		}
		manifest[row.Key] = row
		statusCounts[status]++
	}

	for _, nodeType := range nodeLocSwitchTypes(t) {
		if _, ok := manifest["*nodes."+nodeType]; !ok {
			t.Fatalf("Loc-bearing node *nodes.%s is missing from Loc coverage manifest", nodeType)
		}
	}

	for _, row := range rows {
		if coverageStatus(row) != "covered" {
			continue
		}
		fixture := row.Fields["fixture"]
		if fixture == "" {
			t.Fatalf("%s: covered Loc node row has empty fixture", row.Key)
		}
		tree, err := Parse(fixture)
		if err != nil {
			t.Errorf("%s: Parse(%q): %v", row.Key, fixture, err)
			continue
		}
		want := strings.TrimPrefix(row.Key, "*nodes.")
		if !astContainsType(tree, want) {
			t.Errorf("%s: fixture %q did not produce node type %s", row.Key, fixture, want)
			continue
		}
		if violations := CheckLocations(t, fixture); len(violations) > 0 {
			t.Errorf("%s: Loc violations: %v", row.Key, violations)
		}
	}

	t.Logf("Oracle Loc node coverage: rows=%d status=%v", len(rows), statusCounts)
}

func TestOracleLocNodeFixtureCoverageRequiresApproval(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "loc_node_coverage.tsv"))
	statusCounts := make(map[string]int)
	covered := 0
	for _, row := range rows {
		status := coverageStatus(row)
		statusCounts[status]++
		if status == "covered" {
			covered++
			if row.Fields["fixture"] == "" {
				t.Fatalf("%s: covered Loc row must have fixture", row.Key)
			}
			continue
		}
		if row.Fields["debt_class"] == "" {
			t.Fatalf("%s: non-covered Loc row must have debt_class", row.Key)
		}
		if row.Fields["approval"] == "" {
			t.Fatalf("%s: non-covered Loc row must have approval", row.Key)
		}
		if row.Fields["next_action"] == "" {
			t.Fatalf("%s: non-covered Loc row must have next_action", row.Key)
		}
	}
	if covered < 80 {
		t.Fatalf("direct Loc fixture coverage = %d, want at least 80", covered)
	}
	t.Logf("Oracle Loc node direct fixture coverage=%d status=%v", covered, statusCounts)
}

func nodeLocSwitchTypes(t *testing.T) []string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "ast", "loc.go"))
	if err != nil {
		t.Fatalf("read ast loc.go: %v", err)
	}
	re := regexp.MustCompile(`case \*([A-Za-z0-9_]+):`)
	matches := re.FindAllStringSubmatch(string(content), -1)
	types := make([]string, 0, len(matches))
	for _, match := range matches {
		types = append(types, match[1])
	}
	return types
}

func astContainsType(n nodes.Node, typeName string) bool {
	return valueContainsType(reflect.ValueOf(n), typeName)
}

func valueContainsType(v reflect.Value, typeName string) bool {
	if !v.IsValid() {
		return false
	}
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		if v.Type().Name() == typeName {
			return true
		}
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if field.IsExported() && valueContainsType(v.Field(i), typeName) {
				return true
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if valueContainsType(v.Index(i), typeName) {
				return true
			}
		}
	}
	return false
}
