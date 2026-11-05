package parser

import "strings"

// charsetIntroducerNames is the set of character set names MariaDB reads as an
// introducer when one follows an underscore (_utf8mb4'...'): the character sets of SHOW
// CHARACTER SET on mariadb:11.8.8, plus utf8, which names utf8mb3, and filename, which
// SHOW leaves out. MariaDB has no gb18030, and a collation name such as utf8mb4_bin is no
// introducer. Matching is case-insensitive. TestIntroducerWordsMatchMariaDB holds the set
// to the server.
//
// Ref: MY_LEX_IDENT (get_charset_by_csname with MY_CS_PRIMARY) in sql_lex.cc
// (mariadb-11.8.8).
var charsetIntroducerNames = map[string]bool{
	"armscii8": true, "ascii": true, "big5": true, "binary": true,
	"cp1250": true, "cp1251": true, "cp1256": true, "cp1257": true,
	"cp850": true, "cp852": true, "cp866": true, "cp932": true,
	"dec8": true, "eucjpms": true, "euckr": true, "filename": true,
	"gb2312": true, "gbk": true, "geostd8": true, "greek": true,
	"hebrew": true, "hp8": true, "keybcs2": true, "koi8r": true,
	"koi8u": true, "latin1": true, "latin2": true, "latin5": true,
	"latin7": true, "macce": true, "macroman": true, "sjis": true,
	"swe7": true, "tis620": true, "ucs2": true, "ujis": true,
	"utf16": true, "utf16le": true, "utf32": true, "utf8": true,
	"utf8mb3": true, "utf8mb4": true,
}

// isCharsetIntroducer reports whether ident is an underscore and a character set name
// MariaDB knows, which its lexer reads as an introducer and not as a name.
func isCharsetIntroducer(ident string) bool {
	if len(ident) < 2 || ident[0] != '_' {
		return false
	}
	return charsetIntroducerNames[strings.ToLower(ident[1:])]
}
