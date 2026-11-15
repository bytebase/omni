package catalog

import (
	"strings"
	"testing"
)

func TestParseVersion(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want Version
		ok   bool
	}{
		{"10.6.12-MariaDB-1:10.6.12+maria~ubu2004", MariaDB107, false},
		{"10.7.8-MariaDB-1:10.7.8+maria~ubu2004", MariaDB107, true},
		{"10.9.2-MariaDB", MariaDB107, true},
		{"10.11.19-MariaDB-ubu2204", MariaDB1011, true},
		{"11.2.4-MariaDB", MariaDB1011, true},
		{"11.4.13-MariaDB-ubu2404", MariaDB114, true},
		{"11.8.8-MariaDB-ubu2404", MariaDB118, true},
		{"12.1.1-MariaDB", MariaDB118, true},
		{"garbage", MariaDB118, false},
	} {
		got, ok := ParseVersion(tc.in)
		if got != tc.want || ok != tc.ok {
			t.Errorf("ParseVersion(%q) = %v, %v; want %v, %v", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestVersionDefaults(t *testing.T) {
	const schema = "CREATE DATABASE d; USE d; CREATE TABLE t (a VARCHAR(10), b VARCHAR(10) CHARACTER SET utf8mb3, ts TIMESTAMP, k INT, KEY ik (k DESC)) DEFAULT CHARSET=utf8mb4;"
	for _, tc := range []struct {
		v    Version
		want []string
	}{
		{MariaDB107, []string{"COLLATE=utf8mb4_general_ci", "CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci", "`ts` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()", "KEY `ik` (`k`)\n"}},
		{MariaDB1011, []string{"COLLATE=utf8mb4_general_ci", "CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci", "`ts` timestamp NULL DEFAULT NULL", "KEY `ik` (`k` DESC)"}},
		{MariaDB114, []string{"COLLATE=utf8mb4_uca1400_ai_ci", "CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci", "KEY `ik` (`k` DESC)"}},
		{MariaDB118, []string{"COLLATE=utf8mb4_uca1400_ai_ci", "CHARACTER SET utf8mb3 COLLATE utf8mb3_uca1400_ai_ci", "KEY `ik` (`k` DESC)"}},
	} {
		c, err := LoadSDLWithVersion(schema, tc.v)
		if err != nil {
			t.Fatal(err)
		}
		got := c.ShowCreateTable("d", "t")
		for _, w := range tc.want {
			if !strings.Contains(got, w) {
				t.Errorf("version %v: want %q in\n%s", tc.v, w, got)
			}
		}
	}
}
