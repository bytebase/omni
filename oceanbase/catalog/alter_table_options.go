package catalog

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// checkAlterTableOptions refuses an ALTER TABLE whose table options OceanBase
// refuses on tbl. OceanBase checks every table option before any other command
// of the statement, in the order the statement writes them, so the first bad
// option is the answer even after ADD b where b exists (1060) or DROP COLUMN zz
// (1091). Measured on the pinned 4.3.5.6 against a plain table: the
// simulate.tsv rows of rule table_option and TestOracle_OceanBaseAlterTableOptions.
func checkAlterTableOptions(tbl *Table, cmds []*nodes.AlterTableCmd) error {
	for _, cmd := range cmds {
		if cmd.Type != nodes.ATTableOption || cmd.Option == nil {
			continue
		}
		if err := checkAlterTableOption(tbl, cmd.Option); err != nil {
			return err
		}
	}
	return nil
}

func checkAlterTableOption(tbl *Table, opt *nodes.TableOption) error {
	name := strings.ToUpper(opt.Name)
	value := unquoteOption(opt.Value)
	switch {
	case name == "TABLE_MODE":
		switch strings.ToLower(value) {
		case "normal", "queuing", "moderate", "super", "extreme":
			return nil
		}
		shown := value
		if shown == "" {
			shown = "(null)"
		}
		return errNotSupported(fmt.Sprintf("Table mode %s is", shown))
	case name == "DUPLICATE_SCOPE":
		return oneOfOption(value, "duplicate_scope", "none", "cluster")
	case name == "DUPLICATE_READ_CONSISTENCY":
		return oneOfOption(value, "duplicate_read_consistency", "strong", "weak")
	case name == "SEMISTRUCT_ENCODING_TYPE":
		switch strings.ToLower(value) {
		case "":
			return nil
		case "encoding":
			// Taken on a ROW_FORMAT = COMPRESSED table and refused on the default
			// DYNAMIC and on CONDENSED; other row formats are not measured.
			if strings.EqualFold(tbl.RowFormat, "COMPRESSED") {
				return nil
			}
			return errNotSupported("semistruct_encoding if cs encoding is not set is")
		}
		return wrongArguments("semistruct encoding type, should be none or encoding")
	case name == "AUTO_INCREMENT_MODE":
		switch strings.ToLower(value) {
		case "order", "noorder":
			return nil
		}
		return errOptionSyntax()
	case name == "KV_ATTRIBUTES":
		// Text that is not JSON is 1149. JSON of the wrong shape ('{}') is 1235,
		// and what shape is right is not measured, so JSON is taken.
		if value != "" && !json.Valid([]byte(value)) {
			return errOptionSyntax()
		}
		return nil
	case name == "BLOCK_SIZE":
		if n, err := strconv.ParseInt(value, 10, 64); err == nil && n >= 1024 && n <= 1048576 {
			return nil
		}
		return &Error{Code: 5258, SQLState: "HY000", Message: "Invalid block size, block size should between 1024 and 1048576"}
	case name == "STORAGE_FORMAT_VERSION":
		// A table is created at the release's version, 4, and never goes back.
		n, err := strconv.ParseInt(value, 10, 64)
		switch {
		case err != nil || n < 1 || n > 4:
			return wrongArguments("invalid storage format version")
		case n < 4:
			return &Error{Code: 4179, SQLState: "HY000", Message: "target storage format version cannot be smaller than current version not allowed"}
		}
		return nil
	case name == "TABLET_SIZE":
		// 0, or a whole number of 2 MiB: 1024, 1048576 and 3145728 are refused.
		if n, err := strconv.ParseInt(value, 10, 64); err == nil && n%(2<<20) == 0 {
			return nil
		}
		return &Error{Code: 4147, SQLState: "HY000", Message: "Invalid config"}
	case strings.HasPrefix(name, "ORGANIZATION"):
		return errNotSupported("organization in the alter table statement")
	case name == "MERGE_ENGINE":
		return errNotSupported("alter merge engine table option")
	case name == "AUTO_REFRESH":
		return errNotSupported("Using ALTER AUTO REFRESH as a ALTER TABLE option")
	case name == "STORAGE_CACHE_POLICY":
		// The server runs in shared-nothing mode, as every release image does.
		return errNotSupported("storage cache policy is not supported in shared storage mode")
	case name == "DYNAMIC_PARTITION_POLICY":
		return checkDynamicPartitionPolicy(tbl, opt.Value)
	case name == "JSON":
		// JSON(col) STORE AS (CHUNK n) is refused whatever the column, a JSON one too.
		return &Error{Code: ErrWrongArguments, SQLState: sqlState(ErrWrongArguments), Message: "Invalid argument"}
	case name == "FORMAT", name == "LOCATION", name == "PATTERN", name == "PROPERTIES", name == "PARTITION_TYPE":
		// An external table's options. The catalog holds no external table, and
		// OceanBase answers an internal error for them on any other.
		return &Error{Code: 4016, SQLState: "HY000", Message: "Internal error"}
	}
	return nil
}

var (
	// dynamicPartitionTime is PRECREATE_TIME's and EXPIRE_TIME's value: -1, 0, or
	// a count of a unit no shorter than an hour. '1', '2minute' and '1 hours' are 1210.
	dynamicPartitionTime = regexp.MustCompile(`(?i)^\s*(-1|0|\+?\d+\s*(hour|day|week|month|year))\s*$`)
	// dynamicPartitionZone is TIME_ZONE's offset form; DEFAULT is the other one.
	dynamicPartitionZone = regexp.MustCompile(`^([+-])(\d{1,2}):(\d{2})$`)
)

// checkDynamicPartitionPolicy checks DYNAMIC_PARTITION_POLICY(k = v, ...). Each
// entry is checked in order, then the table: only a RANGE-partitioned table takes
// a policy. What a RANGE table requires further (ENABLE = TRUE alone is 1210, a
// precision must suit the partitioning column) is not modeled.
func checkDynamicPartitionPolicy(tbl *Table, text string) error {
	invalid := &Error{Code: ErrWrongArguments, SQLState: sqlState(ErrWrongArguments), Message: "Invalid argument"}
	entries, ok := optionEntries(text)
	if !ok {
		// Two entries with no comma between them, as in (ENABLE = TRUE ENABLE = TRUE).
		return invalid
	}
	for _, entry := range entries {
		key, value := entry[0], entry[1]
		switch strings.ToUpper(key) {
		case "TIME_UNIT":
			switch strings.ToLower(value) {
			case "hour", "day", "week", "month", "year":
			default:
				return invalid
			}
		case "PRECREATE_TIME", "EXPIRE_TIME":
			if !dynamicPartitionTime.MatchString(value) {
				return invalid
			}
		case "TIME_ZONE":
			if !strings.EqualFold(value, "default") && !validOffsetZone(value) {
				return &Error{Code: 1298, SQLState: "HY000", Message: "Unknown or incorrect time zone"}
			}
		case "BIGINT_PRECISION":
			switch strings.ToLower(value) {
			case "none", "us", "ms", "s":
			default:
				return invalid
			}
		}
	}
	switch {
	case tbl.Partitioning == nil:
		return errNotSupported("specify dynamic partition on non partitioned table is")
	case !strings.HasPrefix(strings.ToUpper(tbl.Partitioning.Type), "RANGE"):
		return errNotSupported("specify dynamic partition on table not partitioned by range is")
	}
	return nil
}

// validOffsetZone reports whether text is an offset OceanBase takes as a time
// zone: a sign, hours and two digits of minutes, from -12:59 to +13:00.
func validOffsetZone(text string) bool {
	m := dynamicPartitionZone.FindStringSubmatch(text)
	if m == nil {
		return false
	}
	h, _ := strconv.Atoi(m[2])
	mm, _ := strconv.Atoi(m[3])
	if mm > 59 {
		return false
	}
	minutes := h*60 + mm
	if m[1] == "-" {
		return minutes <= 12*60+59
	}
	return minutes <= 13*60
}

// An optionToken is one token of an option's value as the parser keeps it, the
// written text: a string (adjacent strings folded into one, escapes decoded), one
// of the punctuation marks ( ) , =, or a run of other characters. Comments are
// dropped, so a comma or = inside one separates nothing.
type optionToken struct {
	text  string
	str   bool // a string literal
	punct byte // '(', ')', ',' or '=' for punctuation, else 0
}

func lexOption(text string) []optionToken {
	var toks []optionToken
	for i := 0; i < len(text); {
		c := text[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case c == '#' || (c == '-' && strings.HasPrefix(text[i:], "-- ")):
			for i < len(text) && text[i] != '\n' {
				i++
			}
		case strings.HasPrefix(text[i:], "/*"):
			if j := strings.Index(text[i+2:], "*/"); j >= 0 {
				i += j + 4
			} else {
				i = len(text)
			}
		case c == '\'' || c == '"':
			value, n := decodeOptionString(text[i:])
			i += n
			if k := len(toks) - 1; k >= 0 && toks[k].str {
				toks[k].text += value // 'DA' 'Y' is 'DAY'
			} else {
				toks = append(toks, optionToken{text: value, str: true})
			}
		case c == '(' || c == ')' || c == ',' || c == '=':
			toks = append(toks, optionToken{text: string(c), punct: c})
			i++
		default:
			j := i
			for j < len(text) && !strings.ContainsRune(" \t\n\r#'\"(),=", rune(text[j])) && !strings.HasPrefix(text[j:], "/*") && !strings.HasPrefix(text[j:], "-- ") {
				j++
			}
			toks = append(toks, optionToken{text: text[i:j]})
			i = j
		}
	}
	return toks
}

// decodeOptionString decodes the string literal text starts with and returns it and
// the bytes it took: a doubled quote is one, and a backslash escapes the next byte.
func decodeOptionString(text string) (string, int) {
	quote := text[0]
	var b strings.Builder
	i := 1
	for i < len(text) {
		c := text[i]
		switch {
		case c == '\\' && i+1 < len(text):
			switch e := text[i+1]; e {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case '0':
				b.WriteByte(0)
			default:
				b.WriteByte(e)
			}
			i += 2
		case c == quote && i+1 < len(text) && text[i+1] == quote:
			b.WriteByte(quote)
			i += 2
		case c == quote:
			return b.String(), i + 1
		default:
			b.WriteByte(c)
			i++
		}
	}
	return b.String(), i
}

// unquoteOption is an option's value as OceanBase reads it: a string literal decoded,
// or the written text of anything else.
func unquoteOption(value string) string {
	toks := lexOption(value)
	if len(toks) == 1 {
		return toks[0].text
	}
	return strings.TrimSpace(value)
}

// optionEntries reads a (key = value, ...) list. ok is false when an entry is not
// one key, one = and one value, as when a comma is missing between two entries.
func optionEntries(text string) (entries [][2]string, ok bool) {
	toks := lexOption(text)
	if len(toks) >= 2 && toks[0].punct == '(' && toks[len(toks)-1].punct == ')' {
		toks = toks[1 : len(toks)-1]
	}
	var entry []optionToken
	flush := func() bool {
		if len(entry) != 3 || entry[1].punct != '=' || entry[0].punct != 0 || entry[2].punct != 0 {
			return false
		}
		entries = append(entries, [2]string{entry[0].text, entry[2].text})
		entry = nil
		return true
	}
	for _, t := range toks {
		if t.punct == ',' {
			if !flush() {
				return nil, false
			}
			continue
		}
		entry = append(entry, t)
	}
	return entries, flush()
}

func oneOfOption(value, option string, allowed ...string) error {
	for _, a := range allowed {
		if strings.EqualFold(value, a) {
			return nil
		}
	}
	return wrongArguments(option)
}

func wrongArguments(what string) error {
	return &Error{Code: ErrWrongArguments, SQLState: sqlState(ErrWrongArguments), Message: "Incorrect arguments to " + what}
}

// errOptionSyntax is the 1149 OceanBase answers for an option value it cannot read.
func errOptionSyntax() error {
	return &Error{Code: 1149, SQLState: "42000", Message: "You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use"}
}

// errNotSupported is OceanBase's 1235, whose text is what it names and "not supported".
func errNotSupported(what string) error {
	return &Error{Code: 1235, SQLState: "0A000", Message: what + " not supported"}
}
