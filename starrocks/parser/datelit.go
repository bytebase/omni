package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// parseDateLiteral reads dateLiteral (StarRocks.g4:2432), DATE or DATETIME
// and a string, as a cast of the string to that type. cur is the keyword.
func (p *Parser) parseDateLiteral() (ast.Node, error) {
	kw := p.advance()
	str, err := p.expectDateString()
	if err != nil {
		return nil, err
	}
	name := "DATE"
	if kw.Kind == kwDATETIME {
		name = "DATETIME"
	}
	return &ast.CastExpr{
		Expr:     &ast.Literal{Kind: ast.LitString, Value: str.Str, Loc: str.Loc},
		TypeName: &ast.TypeName{Name: name, Loc: kw.Loc},
		Loc:      kw.Loc.Merge(str.Loc),
	}, nil
}

// expectDateString consumes the string of a dateLiteral, failing as the
// pinned server does on text that is no date (see validDateLiteral).
func (p *Parser) expectDateString() (Token, error) {
	if p.cur.Kind != tokString {
		return Token{}, p.syntaxErrorAtCur()
	}
	if !validDateLiteral(p.cur.Str) {
		return Token{}, &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "invalid date literal " + p.cur.Str,
		}
	}
	return p.advance(), nil
}

// validDateLiteral reports whether StarRocks takes s, a string's text after
// its quotes, as the value of DATE s or DATETIME s. AstBuilder.visitDateLiteral
// builds a DateLiteral, which parses s with DateUtils.parseStrictDateTime
// (fe-core common/util/DateUtils.java at 3.4.10) and fails the parse with
// "Invalid date literal" when that throws. The checks below follow that
// method's branches and the java.time formatters each picks.
func validDateLiteral(s string) bool {
	if len(s) < 5 {
		return false
	}
	switch {
	case strings.HasSuffix(s, "Z") || strings.HasSuffix(s, "]"):
		return isoZonedDateTime(s)
	case strings.Contains(s, "+"):
		rest, ok := isoLocalDateTime(s)
		if !ok {
			return false
		}
		rest, ok = isoOffset(rest)
		return ok && rest == ""
	case strings.Contains(s, ":"):
		// DATETIME_FORMATTERS[isTwoDigit][withMs][withSplitT][withSec]:
		// %Y-%m-%e or %y-%m-%e, ' ' or 'T', %H:%i, then :%s, then .%f, each
		// read by a STRICT resolver. No formatter has a fraction but no
		// seconds.
		twoDigit := len(strings.SplitN(s, "-", 2)[0]) == 2
		withSec := len(javaSplit(s, ":")) > 2
		withMs := strings.Contains(s, ".")
		sep := byte(' ')
		if strings.Contains(s, "T") {
			sep = 'T'
		}
		if withMs && !withSec {
			return false
		}
		y, rest, ok := dateYear(s, twoDigit)
		if !ok {
			return false
		}
		m, d, rest, ok := monthDay(rest)
		if !ok || !validDay(y, m, d) || rest == "" || rest[0] != sep {
			return false
		}
		rest = rest[1:]
		var hh, mm, ss int
		if hh, rest, ok = fixedDigits(rest, 2); !ok {
			return false
		}
		if rest, ok = cut(rest, ':'); !ok {
			return false
		}
		if mm, rest, ok = fixedDigits(rest, 2); !ok {
			return false
		}
		if withSec {
			if rest, ok = cut(rest, ':'); !ok {
				return false
			}
			if ss, rest, ok = fixedDigits(rest, 2); !ok {
				return false
			}
		}
		if withMs {
			if rest, ok = cut(rest, '.'); !ok {
				return false
			}
			n := leadingDigits(rest, 6)
			if n == 0 {
				return false
			}
			rest = rest[n:]
		}
		return rest == "" && hh <= 23 && mm <= 59 && ss <= 59
	}
	// A date: %y-%m-%e when the text before the first '-' is two
	// characters, %Y-%m-%e when it splits into three parts at '-',
	// %Y%m%e (fixed-width month) at eight characters, all STRICT, and
	// %Y%m%d%H%i%s at fourteen, which resolves SMART.
	switch {
	case len(strings.SplitN(s, "-", 2)[0]) == 2:
		return strictDate(s, true)
	case len(javaSplit(s, "-")) == 3:
		return strictDate(s, false)
	case len(s) == 8:
		y, rest, ok := fixedDigits(s, 4)
		if !ok {
			return false
		}
		m, rest, ok := fixedDigits(rest, 2)
		if !ok {
			return false
		}
		d, rest, ok := fixedDigits(rest, 2)
		return ok && rest == "" && validDay(y, m, d)
	case len(s) == 14:
		var v [6]int
		rest := s
		for i, w := range []int{4, 2, 2, 2, 2, 2} {
			var ok bool
			if v[i], rest, ok = fixedDigits(rest, w); !ok {
				return false
			}
		}
		// SMART keeps a day of 1 to 31 for any month, and 24:00:00 as the
		// next midnight.
		return v[1] >= 1 && v[1] <= 12 && v[2] >= 1 && v[2] <= 31 &&
			(v[3] <= 23 || v[3] == 24 && v[4] == 0 && v[5] == 0) && v[4] <= 59 && v[5] <= 59
	}
	return strictDate(s, false)
}

// strictDate reads all of s as %Y-%m-%e, or %y-%m-%e when twoDigit.
func strictDate(s string, twoDigit bool) bool {
	y, rest, ok := dateYear(s, twoDigit)
	if !ok {
		return false
	}
	m, d, rest, ok := monthDay(rest)
	return ok && rest == "" && validDay(y, m, d)
}

// dateYear reads %Y, four digits, or %y, two digits read as 1970 to 2069,
// and the '-' after it.
func dateYear(s string, twoDigit bool) (int, string, bool) {
	w := 4
	if twoDigit {
		w = 2
	}
	y, rest, ok := fixedDigits(s, w)
	if !ok {
		return 0, "", false
	}
	if twoDigit {
		if y < 70 {
			y += 2000
		} else {
			y += 1900
		}
	}
	rest, ok = cut(rest, '-')
	return y, rest, ok
}

// monthDay reads %m-%e: a month and a day of one or two digits each.
func monthDay(s string) (int, int, string, bool) {
	n := leadingDigits(s, 2)
	if n == 0 {
		return 0, 0, "", false
	}
	m := atoiDigits(s[:n])
	rest, ok := cut(s[n:], '-')
	if !ok {
		return 0, 0, "", false
	}
	n = leadingDigits(rest, 2)
	if n == 0 {
		return 0, 0, "", false
	}
	return m, atoiDigits(rest[:n]), rest[n:], true
}

// isoZonedDateTime reads DateTimeFormatter.ISO_ZONED_DATE_TIME: a local date
// and time, an offset, and an optional bracketed zone ID (see zoneID).
func isoZonedDateTime(s string) bool {
	rest, ok := isoLocalDateTime(s)
	if !ok {
		return false
	}
	if rest, ok = isoOffset(rest); !ok {
		return false
	}
	if rest == "" {
		return true
	}
	if len(rest) < 3 || rest[0] != '[' || rest[len(rest)-1] != ']' {
		return false
	}
	return zoneID(rest[1 : len(rest)-1])
}

// zoneID reports whether ISO_ZONED_DATE_TIME's bracket takes id, as Java
// 11's ZoneIdPrinterParser reads it (case-sensitively): GMT0, or UT, UTC or
// GMT alone or with a strict offset; then a strict offset alone; then a
// region the pinned server's JRE knows (javaZoneIDs); then Z.
func zoneID(id string) bool {
	if id == "GMT0" {
		return true
	}
	for _, prefix := range []string{"UTC", "GMT", "UT"} {
		if rest, ok := strings.CutPrefix(id, prefix); ok {
			return rest == "" || strictOffset(rest)
		}
	}
	if id[0] == '+' || id[0] == '-' {
		return strictOffset(id)
	}
	return javaZoneIDs[id] || id == "Z"
}

// strictOffset reports whether s is all one offset as the pattern
// +HH:MM:ss reads it strictly: a sign, HH:MM and an optional :SS, no more
// than 18 hours.
func strictOffset(s string) bool {
	if len(s) != 6 && len(s) != 9 || s[0] != '+' && s[0] != '-' || strings.Count(s, ":") != len(s)/3-1 {
		return false
	}
	rest, ok := isoOffset(s)
	return ok && rest == ""
}

// isoLocalDateTime reads ISO_LOCAL_DATE_TIME, uuuu-MM-dd'T'HH:mm[:ss[.n]],
// and returns what follows it. The year here takes four digits, or a sign
// and four to nine.
func isoLocalDateTime(s string) (string, bool) {
	var y int
	var ok bool
	rest := s
	if rest != "" && (rest[0] == '+' || rest[0] == '-') {
		// SignStyle.EXCEEDS_PAD: a minus with four digits or more, a plus
		// only with more than four, and YEAR's range of ±999,999,999.
		neg := rest[0] == '-'
		n := leadingDigits(rest[1:], 10)
		if n < 4 || !neg && n == 4 || n > 9 {
			return "", false
		}
		y = atoiDigits(rest[1 : 1+n])
		if neg {
			y = -y
		}
		rest = rest[1+n:]
	} else {
		n := leadingDigits(rest, 4)
		if n != 4 || leadingDigits(rest, 5) != 4 {
			return "", false
		}
		y, rest = atoiDigits(rest[:4]), rest[4:]
	}
	var m, d, hh, mm int
	if rest, ok = cut(rest, '-'); !ok {
		return "", false
	}
	if m, rest, ok = fixedDigits(rest, 2); !ok {
		return "", false
	}
	if rest, ok = cut(rest, '-'); !ok {
		return "", false
	}
	if d, rest, ok = fixedDigits(rest, 2); !ok || !validDay(y, m, d) {
		return "", false
	}
	if rest == "" || rest[0] != 'T' && rest[0] != 't' {
		return "", false
	}
	if hh, rest, ok = fixedDigits(rest[1:], 2); !ok || hh > 23 {
		return "", false
	}
	if rest, ok = cut(rest, ':'); !ok {
		return "", false
	}
	if mm, rest, ok = fixedDigits(rest, 2); !ok || mm > 59 {
		return "", false
	}
	if r, ok := cut(rest, ':'); ok {
		ss, r, ok := fixedDigits(r, 2)
		if !ok || ss > 59 {
			return "", false
		}
		rest = r
		if r, ok := cut(rest, '.'); ok {
			rest = r[leadingDigits(r, 9):]
		}
	}
	return rest, true
}

// isoOffset reads an offset ID: Z or z, or a sign and HH, HH:mm or HH:mm:ss, no
// more than 18 hours.
func isoOffset(s string) (string, bool) {
	if s == "" {
		return "", false
	}
	// ISO_OFFSET_DATE_TIME reads its offset case-insensitively.
	if s[0] == 'Z' || s[0] == 'z' {
		return s[1:], true
	}
	if s[0] != '+' && s[0] != '-' {
		return "", false
	}
	hh, rest, ok := fixedDigits(s[1:], 2)
	if !ok {
		return "", false
	}
	secs := hh * 3600
	for _, unit := range []int{60, 1} {
		r, ok := cut(rest, ':')
		if !ok {
			break
		}
		v, r, ok := fixedDigits(r, 2)
		if !ok || v > 59 {
			return "", false
		}
		secs += v * unit
		rest = r
	}
	return rest, secs <= 18*3600
}

// validDay reports whether d is a day of month m of the proleptic Gregorian
// year y.
func validDay(y, m, d int) bool {
	if m < 1 || m > 12 || d < 1 {
		return false
	}
	days := [...]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}[m-1]
	if m == 2 && y%4 == 0 && (y%100 != 0 || y%400 == 0) {
		days = 29
	}
	return d <= days
}

// fixedDigits reads exactly w ASCII digits.
func fixedDigits(s string, w int) (int, string, bool) {
	if leadingDigits(s, w) != w {
		return 0, "", false
	}
	return atoiDigits(s[:w]), s[w:], true
}

// leadingDigits counts the ASCII digits s starts with, up to max.
func leadingDigits(s string, max int) int {
	n := 0
	for n < len(s) && n < max && s[n] >= '0' && s[n] <= '9' {
		n++
	}
	return n
}

func atoiDigits(s string) int {
	v := 0
	for i := 0; i < len(s); i++ {
		v = v*10 + int(s[i]-'0')
	}
	return v
}

// cut consumes the byte c that s starts with.
func cut(s string, c byte) (string, bool) {
	if s == "" || s[0] != c {
		return "", false
	}
	return s[1:], true
}

// javaSplit is Java's String.split on a literal separator: trailing empty
// parts are dropped.
func javaSplit(s, sep string) []string {
	parts := strings.Split(s, sep)
	for len(parts) > 0 && parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}
