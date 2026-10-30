package parser

import "strings"

// charsetIntroducerNames is the set of character set names MySQL accepts as a
// string-literal charset introducer (_utf8mb4'...'). Taken from SHOW CHARACTER
// SET on live MySQL 8.0.32 and 5.7.25 (union), plus the utf8/utf8mb3 aliases —
// both engines accept _utf8 and _utf8mb3 even where SHOW lists only one of
// them. Collation names (utf8mb4_0900_ai_ci, ...) are NOT introducers and are
// deliberately absent. Matching is case-insensitive (_UTF8MB4 is valid).
var charsetIntroducerNames = map[string]bool{
	"armscii8": true, "ascii": true, "big5": true, "binary": true,
	"cp1250": true, "cp1251": true, "cp1256": true, "cp1257": true,
	"cp850": true, "cp852": true, "cp866": true, "cp932": true,
	"dec8": true, "eucjpms": true, "euckr": true, "gb18030": true,
	"gb2312": true, "gbk": true, "geostd8": true, "greek": true,
	"hebrew": true, "hp8": true, "keybcs2": true, "koi8r": true,
	"koi8u": true, "latin1": true, "latin2": true, "latin5": true,
	"latin7": true, "macce": true, "macroman": true, "sjis": true,
	"swe7": true, "tis620": true, "ucs2": true, "ujis": true,
	"utf16": true, "utf16le": true, "utf32": true, "utf8": true,
	"utf8mb3": true, "utf8mb4": true,
}

// isCharsetIntroducer reports whether ident is a charset introducer: an
// underscore followed by a known MySQL character set name. MySQL's lexer only
// forms an introducer token for real charset names — any other _ident before a
// string literal is an ordinary identifier.
func isCharsetIntroducer(ident string) bool {
	if len(ident) < 2 || ident[0] != '_' {
		return false
	}
	return charsetIntroducerNames[strings.ToLower(ident[1:])]
}
