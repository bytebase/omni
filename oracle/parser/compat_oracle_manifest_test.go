package parser

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

var allowedOracleCompatExpect = map[string]struct{}{
	"accept":            {},
	"reject":            {},
	"unsafe":            {},
	"version_dependent": {},
}

var allowedOracleCompatReferenceClass = map[string]struct{}{
	"syntax":            {},
	"semantic_catalog":  {},
	"privilege_catalog": {},
	"fixture_sensitive": {},
	"version_dependent": {},
}

var requiredOracleCompatColumns = []string{
	"id",
	"phase",
	"family",
	"bnf_file",
	"oracle_mode",
	"expect",
	"reference_class",
	"sql",
	"setup",
	"cleanup",
	"min_version",
	"notes",
}

const (
	oracleCompatLaneStandard   = "standard"
	oracleCompatLanePrivileged = "privileged"
)

func TestOracleCompatibilityManifest(t *testing.T) {
	rows, header := readCoverageTSV(t, filepath.Join("testdata", "coverage", "compat_oracle.tsv"))
	requireOracleCompatHeader(t, header)
	if len(rows) < 200 {
		t.Fatalf("Oracle compatibility manifest rows=%d, want at least 200 expanded compatibility rows", len(rows))
	}

	seen := make(map[string]struct{}, len(rows))
	phaseCounts := make(map[string]int)
	familyCounts := make(map[string]int)
	expectCounts := make(map[string]int)
	referenceClassCounts := make(map[string]int)
	for _, row := range rows {
		id := row.Fields["id"]
		if _, ok := seen[id]; ok {
			t.Fatalf("%s: duplicate compatibility row id", id)
		}
		seen[id] = struct{}{}
		if strings.TrimSpace(row.Fields["phase"]) == "" {
			t.Fatalf("%s: phase is required", id)
		}
		if strings.TrimSpace(row.Fields["family"]) == "" {
			t.Fatalf("%s: family is required", id)
		}
		if _, ok := allowedOracleCompatExpect[row.Fields["expect"]]; !ok {
			t.Fatalf("%s: unknown compatibility expectation %q", id, row.Fields["expect"])
		}
		if _, ok := allowedOracleCompatReferenceClass[row.Fields["reference_class"]]; !ok {
			t.Fatalf("%s: unknown compatibility reference_class %q", id, row.Fields["reference_class"])
		}
		if strings.TrimSpace(row.Fields["sql"]) == "" {
			t.Fatalf("%s: sql is required", id)
		}
		if bnf := row.Fields["bnf_file"]; bnf != "" {
			if _, err := filepath.Glob(filepath.Join("bnf", bnf)); err != nil {
				t.Fatalf("%s: invalid bnf_file glob %q: %v", id, bnf, err)
			}
			if matches, _ := filepath.Glob(filepath.Join("bnf", bnf)); len(matches) == 0 {
				t.Fatalf("%s: bnf_file %q has no matching file", id, bnf)
			}
		}

		phaseCounts[row.Fields["phase"]]++
		familyCounts[row.Fields["family"]]++
		expectCounts[row.Fields["expect"]]++
		referenceClassCounts[row.Fields["reference_class"]]++
	}

	t.Logf("Oracle compatibility manifest: rows=%d phases=%s families=%s expectations=%s reference_classes=%s",
		len(rows), sortedCountString(phaseCounts), sortedCountString(familyCounts), sortedCountString(expectCounts), sortedCountString(referenceClassCounts))
}

func oracleCompatLaneIncludesRow(lane string, row coverageRow) bool {
	switch row.Fields["expect"] {
	case "unsafe", "version_dependent":
		return false
	}
	if lane == oracleCompatLanePrivileged {
		return row.Fields["reference_class"] == "privilege_catalog"
	}
	return true
}

func oracleCompatPrivilegedFixture(row coverageRow, runID string) (string, string) {
	id := row.Fields["id"]
	var setup []string
	var cleanup []string

	addUser := func() {
		setup = append(setup, "CREATE USER compat_user_"+runID+" IDENTIFIED BY pass123")
		cleanup = append(cleanup, "DROP USER compat_user_"+runID+" CASCADE")
	}
	addRole := func() {
		setup = append(setup, "CREATE ROLE compat_role_"+runID)
		cleanup = append(cleanup, "DROP ROLE compat_role_"+runID)
	}
	addProfile := func() {
		setup = append(setup, "CREATE PROFILE compat_profile_"+runID+" LIMIT SESSIONS_PER_USER 2")
		cleanup = append(cleanup, "DROP PROFILE compat_profile_"+runID+" CASCADE")
	}
	addTableT := func() {
		setup = append(setup, "CREATE TABLE t (id NUMBER, a NUMBER) ENABLE ROW MOVEMENT")
		cleanup = append(cleanup, "DROP TABLE t PURGE")
	}
	addDirectory := func() {
		setup = append(setup, "CREATE DIRECTORY compat_dir_"+runID+" AS '/tmp'")
		cleanup = append(cleanup, "DROP DIRECTORY compat_dir_"+runID)
	}
	addAuditPolicy := func() {
		setup = append(setup, "CREATE AUDIT POLICY compat_audit_pol_"+runID+" PRIVILEGES CREATE SESSION")
		cleanup = append(cleanup, "NOAUDIT POLICY compat_audit_pol_"+runID)
		cleanup = append(cleanup, "DROP AUDIT POLICY compat_audit_pol_"+runID)
	}
	addDroppedTable := func(name string) {
		setup = append(setup, "CREATE TABLE "+name+" (id NUMBER)")
		setup = append(setup, "DROP TABLE "+name)
		cleanup = append(cleanup, "PURGE TABLE "+name)
	}

	switch id {
	case "compat_admin_005", "compat_admin_006", "compat_admin_007", "compat_admin_008",
		"compat_admin_018", "compat_admin_022":
		addUser()
	case "compat_admin_012", "compat_admin_013", "compat_admin_014":
		addRole()
	case "compat_admin_016", "compat_admin_017":
		addProfile()
	case "compat_admin_019", "compat_admin_020", "compat_admin_038",
		"compat_admin_069", "compat_admin_070":
		addTableT()
		addUser()
	case "compat_admin_024":
		addUser()
		setup = append(setup, "GRANT CREATE SESSION TO compat_user_"+runID)
	case "compat_admin_027":
		addUser()
		setup = append(setup, "GRANT ALL PRIVILEGES TO compat_user_"+runID)
	case "compat_admin_025":
		addTableT()
		addUser()
		setup = append(setup, "GRANT SELECT ON t TO compat_user_"+runID)
	case "compat_admin_021", "compat_admin_026":
		addUser()
		addRole()
		if id == "compat_admin_026" {
			setup = append(setup, "GRANT compat_role_"+runID+" TO compat_user_"+runID)
		}
	case "compat_admin_023":
		addDirectory()
		addUser()
	case "compat_admin_032":
		addTableT()
	case "compat_admin_033", "compat_admin_034", "compat_admin_035", "compat_admin_036":
		addAuditPolicy()
		addUser()
	case "compat_admin_071":
		addDroppedTable("t")
		cleanup = append(cleanup, "DROP TABLE t_old PURGE")
	case "compat_admin_072":
		addDroppedTable("compat_purged_" + runID)
	case "compat_admin_073":
		addDroppedTable("compat_purged_idx_" + runID)
	}

	return strings.Join(setup, ";"), strings.Join(cleanup, ";")
}

func TestOracleCompatibilityLocalReport(t *testing.T) {
	rows, _ := readCoverageTSV(t, filepath.Join("testdata", "coverage", "compat_oracle.tsv"))

	outcomes := make(map[string]int)
	for _, row := range rows {
		switch row.Fields["expect"] {
		case "unsafe", "version_dependent":
			outcomes[row.Fields["expect"]]++
			continue
		}
		_, err := Parse(oracleCompatSQL(row.Fields["sql"], "LOCALRUN"))
		omniAccepts := err == nil
		switch {
		case row.Fields["expect"] == "accept" && omniAccepts:
			outcomes["local_accept_match"]++
		case row.Fields["expect"] == "accept" && !omniAccepts:
			outcomes["local_too_strict"]++
			t.Logf("%s: omni rejects expected accept: %v", row.Fields["id"], err)
		case row.Fields["expect"] == "reject" && !omniAccepts:
			outcomes["local_reject_match"]++
		case row.Fields["expect"] == "reject" && omniAccepts:
			outcomes["local_too_lenient"]++
			t.Logf("%s: omni accepts expected reject SQL", row.Fields["id"])
		}
	}

	t.Logf("Oracle compatibility local report: %s", sortedCountString(outcomes))
}

func TestOracleCompatibilityLaneIncludesRows(t *testing.T) {
	rows := []coverageRow{
		{Key: "syntax", Fields: map[string]string{"expect": "accept", "reference_class": "syntax"}},
		{Key: "privileged", Fields: map[string]string{"expect": "accept", "reference_class": "privilege_catalog"}},
		{Key: "unsafe", Fields: map[string]string{"expect": "unsafe", "reference_class": "syntax"}},
	}

	for _, row := range rows {
		if row.Key != "unsafe" && !oracleCompatLaneIncludesRow(oracleCompatLaneStandard, row) {
			t.Fatalf("standard lane excluded %s", row.Key)
		}
	}
	if oracleCompatLaneIncludesRow(oracleCompatLaneStandard, rows[2]) {
		t.Fatal("standard lane included unsafe row")
	}
	if !oracleCompatLaneIncludesRow(oracleCompatLanePrivileged, rows[1]) {
		t.Fatal("privileged lane excluded privilege_catalog row")
	}
	if oracleCompatLaneIncludesRow(oracleCompatLanePrivileged, rows[0]) {
		t.Fatal("privileged lane included syntax row")
	}
	if oracleCompatLaneIncludesRow(oracleCompatLanePrivileged, rows[2]) {
		t.Fatal("privileged lane included unsafe row")
	}
}

func TestOracleCompatibilityPrivilegedFixture(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupParts  []string
		cleanupPart string
	}{
		{
			name:        "alter user gets generated user",
			id:          "compat_admin_005",
			setupParts:  []string{"CREATE USER compat_user_RUN IDENTIFIED BY pass123"},
			cleanupPart: "DROP USER compat_user_RUN CASCADE",
		},
		{
			name:        "role grant gets generated user and role",
			id:          "compat_admin_021",
			setupParts:  []string{"CREATE USER compat_user_RUN IDENTIFIED BY pass123", "CREATE ROLE compat_role_RUN"},
			cleanupPart: "DROP ROLE compat_role_RUN",
		},
		{
			name:        "object grant gets generated table and user",
			id:          "compat_admin_019",
			setupParts:  []string{"CREATE TABLE t (id NUMBER, a NUMBER)", "CREATE USER compat_user_RUN IDENTIFIED BY pass123"},
			cleanupPart: "DROP TABLE t PURGE",
		},
		{
			name:        "audit action gets generated table",
			id:          "compat_admin_032",
			setupParts:  []string{"CREATE TABLE t (id NUMBER, a NUMBER)"},
			cleanupPart: "DROP TABLE t PURGE",
		},
		{
			name:        "directory grant gets generated directory and user",
			id:          "compat_admin_023",
			setupParts:  []string{"CREATE DIRECTORY compat_dir_RUN AS '/tmp'", "CREATE USER compat_user_RUN IDENTIFIED BY pass123"},
			cleanupPart: "DROP DIRECTORY compat_dir_RUN",
		},
		{
			name:        "system privilege revoke gets prior grant",
			id:          "compat_admin_024",
			setupParts:  []string{"CREATE USER compat_user_RUN IDENTIFIED BY pass123", "GRANT CREATE SESSION TO compat_user_RUN"},
			cleanupPart: "DROP USER compat_user_RUN CASCADE",
		},
		{
			name:        "all privileges revoke gets prior all privileges grant",
			id:          "compat_admin_027",
			setupParts:  []string{"CREATE USER compat_user_RUN IDENTIFIED BY pass123", "GRANT ALL PRIVILEGES TO compat_user_RUN"},
			cleanupPart: "DROP USER compat_user_RUN CASCADE",
		},
		{
			name:        "recycle bin table purge gets dropped table",
			id:          "compat_admin_072",
			setupParts:  []string{"CREATE TABLE compat_purged_RUN (id NUMBER)", "DROP TABLE compat_purged_RUN"},
			cleanupPart: "PURGE TABLE compat_purged_RUN",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			row := coverageRow{Fields: map[string]string{"id": tc.id}}
			setup, cleanup := oracleCompatPrivilegedFixture(row, "RUN")
			for _, part := range tc.setupParts {
				if !strings.Contains(setup, part) {
					t.Fatalf("setup %q missing %q", setup, part)
				}
			}
			if !strings.Contains(cleanup, tc.cleanupPart) {
				t.Fatalf("cleanup %q missing %q", cleanup, tc.cleanupPart)
			}
		})
	}
}

func requireOracleCompatHeader(t *testing.T, header []string) {
	t.Helper()
	got := make(map[string]struct{}, len(header))
	for _, name := range header {
		got[name] = struct{}{}
	}
	for _, name := range requiredOracleCompatColumns {
		if _, ok := got[name]; !ok {
			t.Fatalf("compat_oracle.tsv missing required column %q", name)
		}
	}
}

func oracleCompatSQL(sqlText, runID string) string {
	return strings.ReplaceAll(sqlText, "{run}", runID)
}

func sortedCountString(counts map[string]int) string {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+"="+itoaForTest(counts[key]))
	}
	return strings.Join(parts, ",")
}

func itoaForTest(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
