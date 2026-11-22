package catalog

import (
	"fmt"
	"regexp"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
)

const (
	oracleDefaultSeqMaxAsc  = "9999999999999999999999999999"
	oracleDefaultSeqMinAsc  = "1"
	oracleDefaultSeqMaxDesc = "-1"
	oracleDefaultSeqMinDesc = "-999999999999999999999999999"
	oracleDefaultSeqCache   = "20"
)

// identityOptionValues holds the sequence options beyond START WITH and
// INCREMENT BY that an Oracle identity column carries in ColumnMetadata.Options.
type identityOptionValues struct {
	maxValue string
	minValue string
	cache    string
	cycle    bool
	order    bool
	keep     bool
	scale    bool
	extend   bool
	session  bool
	shard    bool
}

func parseColumnIdentityOptions(options []string) identityOptionValues {
	var v identityOptionValues
	for _, entry := range options {
		key, val, _ := strings.Cut(entry, "=")
		switch strings.ToUpper(strings.TrimSpace(key)) {
		case "MAXVALUE":
			v.maxValue = strings.TrimSpace(val)
		case "MINVALUE":
			v.minValue = strings.TrimSpace(val)
		case "CACHE":
			v.cache = strings.TrimSpace(val)
		case "NOCACHE":
			v.cache = "0"
		case "CYCLE":
			v.cycle = val == "" || strings.EqualFold(val, "true") || strings.EqualFold(val, "Y")
		case "ORDER":
			v.order = val == "" || strings.EqualFold(val, "true") || strings.EqualFold(val, "Y")
		case "KEEP":
			v.keep = val == "" || strings.EqualFold(val, "true") || strings.EqualFold(val, "Y")
		case "SCALE":
			v.scale = val == "" || strings.EqualFold(val, "true") || strings.EqualFold(val, "Y")
		case "EXTEND":
			v.extend = val == "" || strings.EqualFold(val, "true") || strings.EqualFold(val, "Y")
		case "SESSION":
			v.session = val == "" || strings.EqualFold(val, "true") || strings.EqualFold(val, "Y")
		case "SHARD":
			v.shard = val == "" || strings.EqualFold(val, "true") || strings.EqualFold(val, "Y")
		}
	}
	return v
}

// canonicalizeIdentityOptions encodes non-default Oracle identity sequence
// options into canonical key=value entries for ColumnMetadata.Options.
func canonicalizeIdentityOptions(v identityOptionValues, increment int64) []string {
	descending := increment < 0
	defaultMax := oracleDefaultSeqMaxAsc
	defaultMin := oracleDefaultSeqMinAsc
	if descending {
		defaultMax = oracleDefaultSeqMaxDesc
		defaultMin = oracleDefaultSeqMinDesc
	}
	var out []string
	if v.maxValue != "" && v.maxValue != defaultMax {
		out = append(out, "MAXVALUE="+v.maxValue)
	}
	if v.minValue != "" && v.minValue != defaultMin {
		out = append(out, "MINVALUE="+v.minValue)
	}
	if v.cache != "" && v.cache != oracleDefaultSeqCache {
		out = append(out, "CACHE="+v.cache)
	}
	if v.cycle {
		out = append(out, "CYCLE=true")
	}
	if v.order {
		out = append(out, "ORDER=true")
	}
	if v.keep {
		out = append(out, "KEEP=true")
	}
	if v.scale {
		out = append(out, "SCALE=true")
		if v.extend {
			out = append(out, "EXTEND=true")
		}
	}
	if v.session {
		out = append(out, "SESSION=true")
	}
	if v.shard {
		out = append(out, "SHARD=true")
	}
	return out
}

func (v identityOptionValues) sqlClauses() []string {
	var options []string
	if v.maxValue != "" {
		options = append(options, "MAXVALUE "+v.maxValue)
	}
	if v.minValue != "" {
		options = append(options, "MINVALUE "+v.minValue)
	}
	if v.cache == "0" {
		options = append(options, "NOCACHE")
	} else if v.cache != "" {
		options = append(options, "CACHE "+v.cache)
	}
	if v.cycle {
		options = append(options, "CYCLE")
	}
	if v.order {
		options = append(options, "ORDER")
	}
	if v.keep {
		options = append(options, "KEEP")
	}
	if v.scale {
		if v.extend {
			options = append(options, "SCALE EXTEND")
		} else {
			options = append(options, "SCALE")
		}
	}
	if v.session {
		options = append(options, "SESSION")
	}
	if v.shard {
		options = append(options, "SHARD")
	}
	return options
}

var (
	identityMaxValuePattern = regexp.MustCompile(`(?i)\bMAXVALUE\s+([-+]?[0-9]+)`)
	identityMinValuePattern = regexp.MustCompile(`(?i)\bMINVALUE\s+([-+]?[0-9]+)`)
	identityCachePattern    = regexp.MustCompile(`(?i)\bCACHE\s+([0-9]+)`)
	identityWordPattern     = regexp.MustCompile(`(?i)\b(NOMAXVALUE|NOMINVALUE|NOCACHE|CYCLE|NOCYCLE|ORDER|NOORDER|KEEP|NOKEEP|SCALE|NOSCALE|EXTEND|NOEXTEND|SESSION|SHARD|NOSHARD)\b`)
)

func applyIdentitySequenceOptionsFromText(column *metadatapb.ColumnMetadata, text string) {
	v := parseColumnIdentityOptions(column.GetOptions())
	if match := identityMaxValuePattern.FindStringSubmatch(text); len(match) == 2 {
		v.maxValue = match[1]
	}
	if match := identityMinValuePattern.FindStringSubmatch(text); len(match) == 2 {
		v.minValue = match[1]
	}
	if match := identityCachePattern.FindStringSubmatch(text); len(match) == 2 {
		v.cache = match[1]
	}
	for _, match := range identityWordPattern.FindAllStringSubmatch(text, -1) {
		switch strings.ToUpper(match[1]) {
		case "NOMAXVALUE":
			v.maxValue = ""
		case "NOMINVALUE":
			v.minValue = ""
		case "NOCACHE":
			v.cache = "0"
		case "CYCLE":
			v.cycle = true
		case "NOCYCLE":
			v.cycle = false
		case "ORDER":
			v.order = true
		case "NOORDER":
			v.order = false
		case "KEEP":
			v.keep = true
		case "NOKEEP":
			v.keep = false
		case "SCALE":
			v.scale = true
		case "NOSCALE":
			v.scale = false
			v.extend = false
		case "EXTEND":
			v.extend = true
		case "NOEXTEND":
			v.extend = false
		case "SESSION":
			v.session = true
		case "SHARD":
			v.shard = true
		case "NOSHARD":
			v.shard = false
		}
	}
	column.Options = canonicalizeIdentityOptions(v, column.GetIdentityIncrement())
}

// columnIdentityClause renders the GENERATED ... AS IDENTITY clause a synced
// identity column carries.
func columnIdentityClause(column *metadatapb.ColumnMetadata) string {
	if !column.GetIsIdentity() {
		return ""
	}
	mode := " GENERATED ALWAYS AS IDENTITY"
	if column.GetIdentityGeneration() == metadatapb.ColumnMetadata_BY_DEFAULT {
		mode = " GENERATED BY DEFAULT"
		if column.GetDefaultOnNull() {
			mode += " ON NULL"
		}
		mode += " AS IDENTITY"
	}
	var options []string
	if seed := column.GetIdentitySeed(); seed != 0 {
		options = append(options, fmt.Sprintf("START WITH %d", seed))
	}
	if increment := column.GetIdentityIncrement(); increment != 0 {
		options = append(options, fmt.Sprintf("INCREMENT BY %d", increment))
	}
	options = append(options, parseColumnIdentityOptions(column.GetOptions()).sqlClauses()...)
	if len(options) > 0 {
		mode += " (" + strings.Join(options, " ") + ")"
	}
	return mode
}
