package catalog

import (
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// OceanBase gives every fulltext index a parser and that parser's properties, and SHOW
// CREATE TABLE prints both whether or not the DDL named them: FULLTEXT KEY fk (c)
// reads back as FULLTEXT KEY `fk` (`c`) WITH PARSER space
// PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384 on 4.3.5.6.
// So the catalog keeps what the DDL wrote, renders that, and compares the effective
// form, with the engine's defaults filled in (TestOracle_OceanBaseFulltextParserEcho).
//
// The defaults are those of ObFTParserJsonProps::rebuild_props_for_ddl and
// OB_DEFAULT_FULLTEXT_PARSER_NAME (src/storage/fts/ob_fts_parser_property.cpp,
// src/share/ob_fts_index_builder_util.cpp, v4.3.5_CE_BP6), and each was read back from
// the engine.

// defaultFulltextParser is the parser OceanBase gives a fulltext index whose DDL names
// none.
const defaultFulltextParser = "space"

// fulltextParserDefaults are the properties each built-in parser gets when the DDL
// leaves them out. ik also gets its dictionary tables, which SHOW CREATE TABLE does not
// print.
var fulltextParserDefaults = map[string]map[string]string{
	"space":  {"min_token_size": "3", "max_token_size": "84"},
	"beng":   {"min_token_size": "3", "max_token_size": "84"},
	"ngram":  {"ngram_token_size": "2"},
	"ngram2": {"min_ngram_size": "2", "max_ngram_size": "3"},
	"ik":     {"ik_mode": "smart"},
}

// shownParserProperties are the properties SHOW CREATE TABLE prints, in its order
// (ObFTParserJsonProps::show_parser_properties). It never prints stopword_table,
// dict_table or quantifier_table, so those cannot take part in a comparison against a
// readback.
var shownParserProperties = []string{
	"min_token_size", "max_token_size", "ik_mode", "ngram_token_size", "min_ngram_size", "max_ngram_size",
}

// parserPropertyOrder is the order the catalog renders properties in: the shown ones,
// then the tables.
var parserPropertyOrder = append(append([]string{}, shownParserProperties...), "stopword_table", "dict_table", "quantifier_table")

// stringParserProperties are the properties whose value is a string literal.
var stringParserProperties = map[string]bool{
	"ik_mode": true, "stopword_table": true, "dict_table": true, "quantifier_table": true,
}

// parserPropertiesFromOption reads a PARSER_PROPERTIES option into a map. A key given
// twice keeps its last value, and ik_mode is lower-cased, as 4.3.5.6 stores them:
// (min_token_size=3, min_token_size=6) reads back min_token_size=6, and
// ik_mode='MAX_WORD' reads back ik_mode="max_word".
func parserPropertiesFromOption(opt *nodes.IndexOption) map[string]string {
	props := make(map[string]string, len(opt.Properties))
	for _, prop := range opt.Properties {
		key := strings.ToLower(prop.Name)
		switch v := prop.Value.(type) {
		case *nodes.IntLit:
			props[key] = strconv.FormatInt(v.Value, 10)
		case *nodes.StringLit:
			props[key] = v.Value
			if key == "ik_mode" {
				props[key] = strings.ToLower(v.Value)
			}
		}
	}
	return props
}

// copyParserProperties returns a copy of props, nil for nil, so a copied index does not
// share its properties with the one it was copied from.
func copyParserProperties(props map[string]string) map[string]string {
	if props == nil {
		return nil
	}
	out := make(map[string]string, len(props))
	for k, v := range props {
		out[k] = v
	}
	return out
}

// canonicalFulltextParser returns a fulltext index's parser and the properties SHOW
// CREATE TABLE prints for it, with OceanBase's defaults filled in, as a key fragment;
// "" for any other index. A readback and the DDL it came from give the same fragment,
// and a change of parser or of a shown property gives a different one.
func canonicalFulltextParser(idx *Index) string {
	if !idx.Fulltext {
		return ""
	}
	parser := idx.Parser
	if parser == "" {
		parser = defaultFulltextParser
	}
	effective := make(map[string]string)
	for k, v := range fulltextParserDefaults[parser] {
		effective[k] = v
	}
	for k, v := range idx.ParserProperties {
		effective[k] = v
	}
	var b strings.Builder
	b.WriteString(parser)
	b.WriteString("(")
	first := true
	for _, key := range shownParserProperties {
		v, ok := effective[key]
		if !ok {
			continue
		}
		if !first {
			b.WriteString(",")
		}
		first = false
		b.WriteString(key + "=" + v)
	}
	b.WriteString(")")
	return b.String()
}

// formatFulltextParser renders a fulltext index's WITH PARSER and PARSER_PROPERTIES as
// the DDL wrote them, with a leading space, or "" when it wrote neither. OceanBase
// prints them in this order, right after the columns.
func formatFulltextParser(idx *Index) string {
	if !idx.Fulltext {
		return ""
	}
	var b strings.Builder
	if idx.Parser != "" {
		b.WriteString(" WITH PARSER " + idx.Parser)
	}
	if len(idx.ParserProperties) > 0 {
		parts := make([]string, 0, len(idx.ParserProperties))
		for _, key := range parserPropertyOrder {
			v, ok := idx.ParserProperties[key]
			if !ok {
				continue
			}
			if stringParserProperties[key] {
				v = quoteStringLiteral(v)
			}
			parts = append(parts, key+"="+v)
		}
		b.WriteString(" PARSER_PROPERTIES=(" + strings.Join(parts, ",") + ")")
	}
	return b.String()
}
