package review

import (
	"strings"

	"github.com/bytebase/omni/review/walk"
)

// The walk does not know the database's collation, and on SQL Server two
// ASCII names a collation equates need not differ only in case: Hungarian
// equates ssz with szsz, Welsh_100 f with g, and the accent-insensitive
// Estonian, Finnish_Swedish, and Frisian_100 collations v with w and i
// with y (docs/specs/2026-10-26-walk-through-mssql.md §2.1, verified on
// every collation). So every name has a class that any name some
// collation may equate to it shares, and the ledger keeps two keys per
// name: its class, and its exact spelling within the class.
//
//   - The class key is Present when the class certainly holds an object,
//     with the kind and spelling when it holds exactly one, and Absent
//     when it certainly holds none.
//   - The exact key is Present when an object of exactly that spelling
//     exists, and Absent when no object the collation equates to it does.
//
// A name is claimed to exist only for its own spelling, and to be missing
// only when its class or its exact key says so.

// markedPrefix begins the class of a marked name: one with a byte other
// than an ASCII letter, digit, or underscore. Such a name may equal names
// of any class (é and e under an accent-insensitive collation, an
// ignorable hyphen in Maori_100), so it is a class of its own, and
// creating or dropping one may add or remove any name of its space.
const markedPrefix = "\x01"

// exactSep joins a class to a spelling in an exact key. A name holding
// it, or the ledger's own separator, cannot be keyed.
const exactSep = "\x02"

// name is one name as a statement writes it.
type name struct {
	spelling string
	class    string
	plain    bool
}

// keyed returns a name's class, false for a name the walk cannot key.
//
// A plain name's class is built in three steps: lowercase it and map w to
// v, y to i, and g to f; keep the first occurrence of each byte; append
// the last byte. Doubling a Hungarian digraph adds no first occurrence and
// keeps the last byte, so ssz and szsz share sz|z.
func keyed(spelling string) (name, bool) {
	if spelling == "" || strings.ContainsAny(spelling, "\x00"+exactSep) {
		return name{}, false
	}
	for i := 0; i < len(spelling); i++ {
		if !plainByte(spelling[i]) {
			return name{spelling: spelling, class: markedPrefix + spelling}, true
		}
	}
	return name{spelling: spelling, class: classPrefix(spelling) + "|" + string(fold(spelling[len(spelling)-1])), plain: true}, true
}

func plainByte(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_'
}

// fold maps a plain byte to the byte every collation may equate it to.
func fold(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		c += 'a' - 'A'
	}
	switch c {
	case 'w':
		return 'v'
	case 'y':
		return 'i'
	case 'g':
		return 'f'
	}
	return c
}

// classPrefix returns what the class of every plain name that begins with
// the plain text p begins with: p's bytes folded, each kept once.
func classPrefix(p string) string {
	var seen [128]bool
	var b strings.Builder
	for i := 0; i < len(p); i++ {
		c := fold(p[i])
		if !seen[c] {
			seen[c] = true
			b.WriteByte(c)
		}
	}
	return b.String()
}

func (n name) exact() string { return n.class + exactSep + n.spelling }

// classKey and exactKey key a name in a space under a parent.
func (n name) classKey(space, parent string) walk.Key {
	return walk.Key{Space: space, Parent: parent, Name: n.class}
}

func (n name) exactKey(space, parent string) walk.Key {
	return walk.Key{Space: space, Parent: parent, Name: n.exact()}
}

// isClassKey reports whether a key is a class key, not an exact one.
func isClassKey(k walk.Key) bool { return !strings.Contains(k.Name, exactSep) }

// A class key's Present fact holds the kind and spelling of the one
// object the class holds, or "" when it may hold several.
func classKind(kind, spelling string) string { return kind + "\x00" + spelling }

// single returns the kind and spelling of a class key's one object.
func single(fact string) (kind, spelling string, ok bool) {
	kind, spelling, ok = strings.Cut(fact, "\x00")
	return kind, spelling, ok
}
