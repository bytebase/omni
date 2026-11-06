package catalog

import (
	"fmt"
	"strings"
)

// TiDB's character sets and collations, as pingcap/tidb:v8.5.5 lists them in
// SHOW CHARACTER SET and information_schema.COLLATIONS (read on 2026-11-06).
// The server runs with the new collation framework, so this is the whole list:
// a charset outside it answers 1115 and a collation outside it 1273, even one
// MySQL has (`latin1_swedish_ci`, `utf8mb4_0900_as_cs`).
//
// TiDB names MySQL's utf8mb3 `utf8`. It reads `utf8mb3` and `utf8mb3_*` as
// aliases and stores and prints the `utf8` spelling: `CHARACTER SET utf8mb3`
// comes back from SHOW CREATE TABLE as `CHARACTER SET utf8 COLLATE utf8_bin`.

// tidbCollationCharset maps every collation TiDB has to its charset.
var tidbCollationCharset = map[string]string{
	"ascii_bin":          "ascii",
	"binary":             "binary",
	"gbk_bin":            "gbk",
	"gbk_chinese_ci":     "gbk",
	"latin1_bin":         "latin1",
	"utf8_bin":           "utf8",
	"utf8_general_ci":    "utf8",
	"utf8_unicode_ci":    "utf8",
	"utf8mb4_0900_ai_ci": "utf8mb4",
	"utf8mb4_0900_bin":   "utf8mb4",
	"utf8mb4_bin":        "utf8mb4",
	"utf8mb4_general_ci": "utf8mb4",
	"utf8mb4_unicode_ci": "utf8mb4",
}

// normalizeCharsetName lowers a charset name and spells utf8mb3 as TiDB does.
func normalizeCharsetName(cs string) string {
	cs = strings.ToLower(cs)
	if cs == "utf8mb3" {
		return "utf8"
	}
	return cs
}

// normalizeCollationName lowers a collation name and spells utf8mb3_* as TiDB
// does.
func normalizeCollationName(coll string) string {
	coll = strings.ToLower(coll)
	if strings.HasPrefix(coll, "utf8mb3_") {
		return "utf8_" + strings.TrimPrefix(coll, "utf8mb3_")
	}
	return coll
}

func errUnknownCharset(cs string) error {
	return &Error{Code: ErrUnknownCharacterSet, SQLState: sqlState(ErrUnknownCharacterSet),
		Message: fmt.Sprintf("Unknown character set: '%s'", cs)}
}

func errUnknownCollation(coll string) error {
	if knownCollations[strings.ToLower(coll)] {
		return &Error{Code: ErrUnknownCollation, SQLState: sqlState(ErrUnknownCollation),
			Message: fmt.Sprintf("Unsupported collation when new collation is enabled: '%s'", coll)}
	}
	return &Error{Code: ErrUnknownCollation, SQLState: sqlState(ErrUnknownCollation),
		Message: fmt.Sprintf("Unknown collation: '%s'", coll)}
}

func errCollationCharsetMismatch(coll, cs string) error {
	return &Error{Code: ErrCollationCharsetMismatch, SQLState: sqlState(ErrCollationCharsetMismatch),
		Message: fmt.Sprintf("COLLATION '%s' is not valid for CHARACTER SET '%s'", coll, cs)}
}

// checkCharset returns the TiDB spelling of a charset the user named, or 1115.
func checkCharset(cs string) (string, error) {
	n := normalizeCharsetName(cs)
	if _, ok := defaultCollationForCharset[n]; !ok {
		return "", errUnknownCharset(cs)
	}
	return n, nil
}

// checkCollation returns the TiDB spelling of a collation the user named, or
// 1273.
func checkCollation(coll string) (string, error) {
	n := normalizeCollationName(coll)
	if _, ok := tidbCollationCharset[n]; !ok {
		return "", errUnknownCollation(coll)
	}
	return n, nil
}

// resolveCharsetCollation checks a charset and collation a statement names,
// either of which may be empty, and fills in whichever is missing the way TiDB
// does: a charset's default collation, a collation's charset. An empty pair
// stays empty, for the caller to inherit. A collation of another charset
// answers 1253.
func resolveCharsetCollation(cs, coll string) (string, string, error) {
	var err error
	if cs != "" {
		if cs, err = checkCharset(cs); err != nil {
			return "", "", err
		}
	}
	if coll != "" {
		if coll, err = checkCollation(coll); err != nil {
			return "", "", err
		}
	}
	switch {
	case cs != "" && coll != "":
		if tidbCollationCharset[coll] != cs {
			return "", "", errCollationCharsetMismatch(coll, cs)
		}
	case cs != "":
		coll = defaultCollationForCharset[cs]
	case coll != "":
		cs = tidbCollationCharset[coll]
	}
	return cs, coll, nil
}

// applyColumnCollate applies a COLLATE written after a column's type, as in
// `a VARCHAR(5) NOT NULL COLLATE utf8mb4_general_ci`. The collation decides the
// charset unless the type named one, which it must then belong to.
func applyColumnCollate(col *Column, coll string, charsetExplicit bool) error {
	n, err := checkCollation(coll)
	if err != nil {
		return err
	}
	cs := tidbCollationCharset[n]
	if charsetExplicit && col.Charset != cs {
		return errCollationCharsetMismatch(n, col.Charset)
	}
	col.Charset, col.Collation = cs, n
	return nil
}
