package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

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

// parseIntroducedByteString reads an introducer and the hex or bit literal after it, which
// make a string of that character set (UNDERSCORE_CHARSET hex_or_bin_String in literal,
// MariaDB 11.8 sql/sql_yacc.yy). The digits are held to the literal's own rules, and the
// bytes to the character set's (check_well_formed_result), while the statement parses:
// mariadb:11.8.8 answers 1300 for _utf8mb4 0xC0 and _utf8mb3 0xF09F9880. Only the UTF-8
// sets are checked here.
func (p *Parser) parseIntroducedByteString() (nodes.ExprNode, error) {
	charsetTok := p.advance()
	litTok := p.advance()
	value, err := p.decodeByteLiteral(litTok)
	if err != nil {
		return nil, err
	}
	cs := strings.ToLower(charsetTok.Str[1:])
	if (cs == "utf8mb4" || cs == "utf8mb3" || cs == "utf8") && !wellFormedUTF8(value, cs == "utf8mb4") {
		return nil, &ParseError{Message: "Invalid " + cs + " character string", Position: litTok.Loc}
	}
	return &nodes.StringLit{
		Loc:     nodes.Loc{Start: charsetTok.Loc, End: litTok.End},
		Value:   value,
		Charset: charsetTok.Str,
	}, nil
}

// wellFormedUTF8 reports whether s is UTF-8 the way MariaDB's utf8mb4 and utf8mb3 read it
// (my_charlen_utf8mb4 in strings/ctype-utf8.c): no overlong form and nothing past
// U+10FFFF, four bytes only when fourByte, and the surrogates allowed: mariadb:11.8.8 takes
// _utf8mb4 0xEDA080 and answers 1300 for 0xC080 and 0xF4908080.
func wellFormedUTF8(s string, fourByte bool) bool {
	cont := func(c byte) bool { return c&0xC0 == 0x80 }
	for i := 0; i < len(s); {
		c := s[i]
		n := 0
		switch {
		case c < 0x80:
			n = 1
		case c >= 0xC2 && c <= 0xDF:
			n = 2
		case c >= 0xE0 && c <= 0xEF:
			n = 3
		case fourByte && c >= 0xF0 && c <= 0xF4:
			n = 4
		default:
			return false
		}
		if i+n > len(s) {
			return false
		}
		for k := 1; k < n; k++ {
			if !cont(s[i+k]) {
				return false
			}
		}
		if (c == 0xE0 && s[i+1] < 0xA0) || (c == 0xF0 && s[i+1] < 0x90) || (c == 0xF4 && s[i+1] >= 0x90) {
			return false
		}
		i += n
	}
	return true
}

// isCharsetIntroducer reports whether ident is an underscore and a character set name
// MariaDB knows, which its lexer reads as an introducer and not as a name.
func isCharsetIntroducer(ident string) bool {
	if len(ident) < 2 || ident[0] != '_' {
		return false
	}
	return charsetIntroducerNames[strings.ToLower(ident[1:])]
}
