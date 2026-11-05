package parser

import "strings"

// charsetIntroducerNames is the set of character set names OceanBase reads as an
// introducer when one follows an underscore (_utf8mb4'...'). OceanBase's lexer has a
// token for each, a fixed list that is not MySQL's: it has hkscs, hkscs31 and
// gb18030_2022, and no ucs2, utf32, cp1250 or koi8r, which are names there. Matching is
// case-insensitive (_UTF8MB4 is valid). Collation names are not introducers.
// TestIntroducerWordsMatchOceanBase holds the set to the server.
//
// Ref: the _UTF8 ... _SWE7 rules of sql_parser_mysql_mode.l (oceanbase v4.3.5_CE_BP6).
var charsetIntroducerNames = map[string]bool{
	"ascii": true, "big5": true, "binary": true, "cp850": true,
	"cp932": true, "dec8": true, "eucjpms": true, "euckr": true,
	"gb18030": true, "gb18030_2022": true, "gb2312": true, "gbk": true,
	"hkscs": true, "hkscs31": true, "hp8": true, "latin1": true,
	"macroman": true, "sjis": true, "swe7": true, "tis620": true,
	"ujis": true, "utf16": true, "utf16le": true, "utf8": true,
	"utf8mb3": true, "utf8mb4": true,
}

// isCharsetIntroducer reports whether ident is a charset introducer: an
// underscore followed by a character set name OceanBase's lexer has a token
// for. Any other _ident before a string literal is an ordinary identifier.
func isCharsetIntroducer(ident string) bool {
	if len(ident) < 2 || ident[0] != '_' {
		return false
	}
	return charsetIntroducerNames[strings.ToLower(ident[1:])]
}
