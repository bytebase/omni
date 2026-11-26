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
		t, rest, ok := isoLocalDateTime(s)
		if !ok {
			return false
		}
		rest, off, ok := isoOffset(rest)
		return ok && rest == "" && t.inRangeAfter(-off)
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
	t, rest, ok := isoLocalDateTime(s)
	if !ok {
		return false
	}
	rest, off, ok := isoOffset(rest)
	if !ok || !t.inRangeAfter(-off) {
		return false
	}
	if rest == "" {
		return true
	}
	if len(rest) < 3 || rest[0] != '[' || rest[len(rest)-1] != ']' {
		return false
	}
	zoneOff, fixed, ok := zoneID(rest[1 : len(rest)-1])
	// ZonedDateTime resolves the instant into the bracket's zone first. A
	// region's offset at year ±999,999,999 is its tz rules', which omni
	// does not hold, so only a fixed offset is checked there.
	return ok && (!fixed || t.inRangeAfter(zoneOff-off))
}

// zoneID reports whether ISO_ZONED_DATE_TIME's bracket takes id, as Java
// 11's ZoneIdPrinterParser reads it (case-sensitively): GMT0, or UT, UTC or
// GMT alone or with a strict offset; then a strict offset alone; then a
// region the pinned server's JRE knows (javaZoneIDs); then Z. For all but a
// region it returns the fixed offset in seconds.
func zoneID(id string) (off int, fixed, ok bool) {
	if id == "GMT0" || id == "Z" {
		return 0, true, true
	}
	for _, prefix := range []string{"UTC", "GMT", "UT"} {
		if rest, found := strings.CutPrefix(id, prefix); found {
			if rest == "" {
				return 0, true, true
			}
			off, ok = strictOffset(rest)
			return off, true, ok
		}
	}
	if id[0] == '+' || id[0] == '-' {
		off, ok = strictOffset(id)
		return off, true, ok
	}
	return 0, false, javaZoneIDs[id]
}

// strictOffset reads s as one offset as the pattern +HH:MM:ss reads it
// strictly: a sign, HH:MM and an optional :SS, no more than 18 hours. It
// returns the offset in seconds.
func strictOffset(s string) (int, bool) {
	if len(s) != 6 && len(s) != 9 || s[0] != '+' && s[0] != '-' || strings.Count(s, ":") != len(s)/3-1 {
		return 0, false
	}
	rest, off, ok := isoOffset(s)
	return off, ok && rest == ""
}

// An isoTime is the date and the whole seconds of the day an ISO local
// date-time names.
type isoTime struct{ y, m, d, secs int }

// inRangeAfter reports whether t moved by shift seconds, as a change of
// zone moves it, is still a LocalDateTime: java.time's years end at
// ±999,999,999, and a shift of at most 36 hours crosses that only from the
// last or first day. DateUtils converts to the server's zone, UTC on the
// pinned image, which is a shift of minus the parsed offset.
func (t isoTime) inRangeAfter(shift int) bool {
	switch {
	case t.y == 999999999 && t.m == 12 && t.d == 31:
		return t.secs+shift < 86400
	case t.y == 999999999 && t.m == 12 && t.d == 30:
		return t.secs+shift < 2*86400
	case t.y == -999999999 && t.m == 1 && t.d == 1:
		return t.secs+shift >= 0
	case t.y == -999999999 && t.m == 1 && t.d == 2:
		return t.secs+shift >= -86400
	}
	return true
}

// isoLocalDateTime reads ISO_LOCAL_DATE_TIME, uuuu-MM-dd'T'HH:mm[:ss[.n]],
// and returns what follows it. The year here takes four digits, or a sign
// and four to ten.
func isoLocalDateTime(s string) (isoTime, string, bool) {
	var t isoTime
	var y int
	var ok bool
	rest := s
	if rest != "" && (rest[0] == '+' || rest[0] == '-') {
		// appendValue(YEAR, 4, 10, SignStyle.EXCEEDS_PAD): four to ten
		// digits, a plus only with more than four, and YEAR's range of
		// ±999,999,999.
		neg := rest[0] == '-'
		n := leadingDigits(rest[1:], 10)
		if n < 4 || !neg && n == 4 {
			return t, "", false
		}
		// A strict parse refuses a minus on zero.
		if y = atoiDigits(rest[1 : 1+n]); y > 999999999 || neg && y == 0 {
			return t, "", false
		}
		if neg {
			y = -y
		}
		rest = rest[1+n:]
	} else {
		n := leadingDigits(rest, 4)
		if n != 4 || leadingDigits(rest, 5) != 4 {
			return t, "", false
		}
		y, rest = atoiDigits(rest[:4]), rest[4:]
	}
	var m, d, hh, mm, ss int
	if rest, ok = cut(rest, '-'); !ok {
		return t, "", false
	}
	if m, rest, ok = fixedDigits(rest, 2); !ok {
		return t, "", false
	}
	if rest, ok = cut(rest, '-'); !ok {
		return t, "", false
	}
	if d, rest, ok = fixedDigits(rest, 2); !ok || !validDay(y, m, d) {
		return t, "", false
	}
	if rest == "" || rest[0] != 'T' && rest[0] != 't' {
		return t, "", false
	}
	if hh, rest, ok = fixedDigits(rest[1:], 2); !ok || hh > 23 {
		return t, "", false
	}
	if rest, ok = cut(rest, ':'); !ok {
		return t, "", false
	}
	if mm, rest, ok = fixedDigits(rest, 2); !ok || mm > 59 {
		return t, "", false
	}
	if r, ok := cut(rest, ':'); ok {
		if ss, r, ok = fixedDigits(r, 2); !ok || ss > 59 {
			return t, "", false
		}
		rest = r
		if r, ok := cut(rest, '.'); ok {
			rest = r[leadingDigits(r, 9):]
		}
	}
	return isoTime{y, m, d, hh*3600 + mm*60 + ss}, rest, true
}

// isoOffset reads an offset ID: Z or z, or a sign and HH, HH:mm or
// HH:mm:ss, no more than 18 hours. It returns what follows and the offset
// in seconds.
func isoOffset(s string) (string, int, bool) {
	if s == "" {
		return "", 0, false
	}
	// ISO_OFFSET_DATE_TIME reads its offset case-insensitively.
	if s[0] == 'Z' || s[0] == 'z' {
		return s[1:], 0, true
	}
	if s[0] != '+' && s[0] != '-' {
		return "", 0, false
	}
	hh, rest, ok := fixedDigits(s[1:], 2)
	if !ok {
		return "", 0, false
	}
	secs := hh * 3600
	for _, unit := range []int{60, 1} {
		r, ok := cut(rest, ':')
		if !ok {
			break
		}
		v, r, ok := fixedDigits(r, 2)
		if !ok || v > 59 {
			return "", 0, false
		}
		secs += v * unit
		rest = r
	}
	if s[0] == '-' {
		secs = -secs
	}
	return rest, secs, secs <= 18*3600 && secs >= -18*3600
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
