package catalog

import (
	"fmt"
	"strconv"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	oracleparser "github.com/bytebase/omni/oracle/parser"
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

type identityToken struct {
	tokType int
	str     string
	raw     string
	quoted  bool
}

func tokenizeIdentityOptions(text string) []identityToken {
	lexer := oracleparser.NewLexer(text)
	var all []identityToken
	for {
		tok := lexer.NextToken()
		if tok.Type == 0 {
			break
		}
		quoted := tok.Loc >= 0 && tok.Loc < len(text) && text[tok.Loc] == '"'
		raw := tok.Str
		if tok.Loc >= 0 && tok.End >= tok.Loc && tok.End <= len(text) {
			raw = text[tok.Loc:tok.End]
		}
		all = append(all, identityToken{
			tokType: tok.Type,
			str:     tok.Str,
			raw:     raw,
			quoted:  quoted,
		})
	}
	// IdentityClause.Loc.Start starts at the column definition, which includes
	// the column name (e.g. a quoted identifier like "CYCLE") and datatype.
	// Restrict option parsing to tokens after the unquoted IDENTITY keyword.
	start := 0
	for i, tok := range all {
		if !tok.quoted && strings.EqualFold(tok.str, "IDENTITY") {
			start = i + 1
			break
		}
	}
	return all[start:]
}

func readSignedIntToken(tokens []identityToken, idx int) (string, int, bool) {
	if idx >= len(tokens) {
		return "", idx, false
	}
	sign := ""
	if !tokens[idx].quoted && (tokens[idx].raw == "+" || tokens[idx].raw == "-") {
		sign = tokens[idx].raw
		idx++
	}
	if idx >= len(tokens) || tokens[idx].quoted {
		return "", idx, false
	}
	digits := tokens[idx].raw
	if digits == "" {
		return "", idx, false
	}
	for i := 0; i < len(digits); i++ {
		if digits[i] < '0' || digits[i] > '9' {
			return "", idx, false
		}
	}
	return sign + digits, idx + 1, true
}

func applyIdentitySequenceOptionsFromText(column *metadatapb.ColumnMetadata, text string) {
	tokens := tokenizeIdentityOptions(text)
	v := parseColumnIdentityOptions(column.GetOptions())
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		if tok.quoted {
			continue
		}
		switch strings.ToUpper(tok.str) {
		case "START":
			if i+1 < len(tokens) && !tokens[i+1].quoted && strings.EqualFold(tokens[i+1].str, "WITH") {
				if num, next, ok := readSignedIntToken(tokens, i+2); ok {
					if val, err := strconv.ParseInt(num, 10, 64); err == nil {
						column.IdentitySeed = val
					}
					i = next - 1
				}
			}
		case "INCREMENT":
			if i+1 < len(tokens) && !tokens[i+1].quoted && strings.EqualFold(tokens[i+1].str, "BY") {
				if num, next, ok := readSignedIntToken(tokens, i+2); ok {
					if val, err := strconv.ParseInt(num, 10, 64); err == nil {
						column.IdentityIncrement = val
					}
					i = next - 1
				}
			}
		case "MAXVALUE":
			if num, next, ok := readSignedIntToken(tokens, i+1); ok {
				v.maxValue = num
				i = next - 1
			}
		case "MINVALUE":
			if num, next, ok := readSignedIntToken(tokens, i+1); ok {
				v.minValue = num
				i = next - 1
			}
		case "CACHE":
			if num, next, ok := readSignedIntToken(tokens, i+1); ok {
				v.cache = num
				i = next - 1
			}
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
