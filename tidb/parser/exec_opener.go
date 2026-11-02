package parser

// execOpener reports whether s[i:], which starts with "/*", opens an
// executable comment the way TiDB v8.5.5 reads one, and the offset its body
// starts at. TiDB opens one at
//
//   - "/*!" and any version digits after it;
//   - "/*T!" and a feature list "[f1,f2]" whose features it supports;
//   - "/*T!" alone when no well-formed feature list follows, so the body
//     starts at whatever does: "SELECT 1 /*T![a b] */" is a 1064 at the "[".
//
// A "/*T!" whose well-formed feature list names a feature TiDB does not
// support is an ordinary comment, and so is any other "/*".
func execOpener(s string, i int) (int, bool) {
	if i+2 >= len(s) || s[i] != '/' || s[i+1] != '*' {
		return 0, false
	}
	switch {
	case s[i+2] == '!':
		j := i + 3
		for j < len(s) && s[j] >= '0' && s[j] <= '9' {
			j++
		}
		return j, true
	case s[i+2] == 'T' && i+3 < len(s) && s[i+3] == '!':
		end, features, ok := tidbFeatureList(s, i+4)
		if !ok {
			return i + 4, true
		}
		if !allTiDBFeaturesSupported(features) {
			return 0, false
		}
		return end, true
	}
	return 0, false
}

// tidbFeatureList reads a feature list "[f1,f2]" at s[i:] the way TiDB's
// scanner does: one or more names of letters, digits, '_', '$' or non-ASCII
// bytes, separated by single commas with no space, then ']'. It returns the
// offset after the ']' and the names between the brackets. The read stops at
// the first byte that does not fit, so it never scans past the list.
func tidbFeatureList(s string, i int) (end int, features string, ok bool) {
	if i >= len(s) || s[i] != '[' {
		return 0, "", false
	}
	j := i + 1
	for {
		k := j
		for k < len(s) && isFeatureNameByte(s[k]) {
			k++
		}
		if k == j || k >= len(s) {
			return 0, "", false
		}
		switch s[k] {
		case ',':
			j = k + 1
		case ']':
			return k + 1, s[i+1 : k], true
		default:
			return 0, "", false
		}
	}
}

func isFeatureNameByte(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
		c == '_' || c == '$' || c >= 0x80
}

// IsExecutableComment reports whether sql[i:] opens an executable comment,
// "/*!", "/*T!" or a supported "/*T![feature]" gate, whose text TiDB runs. Any
// other "/*", an unsupported gate included, is an ordinary comment.
func IsExecutableComment(sql string, i int) bool {
	_, ok := execOpener(sql, i)
	return ok
}

// FeatureGate reports whether sql[i:] opens a TiDB feature gate, a
// "/*T![feature,...]" comment whose features TiDB supports, and where it ends:
// body is the offset its text starts at, and closeStart the offset of the "*/"
// that closes it, found as the lexer finds it. closeStart is -1 when nothing
// closes the gate, and TiDB then drops the opener and reads the rest as
// ordinary SQL.
func FeatureGate(sql string, i int) (body, closeStart int, ok bool) {
	if i+4 >= len(sql) || sql[i+2] != 'T' || sql[i+4] != '[' {
		return 0, 0, false
	}
	body, ok = execOpener(sql, i)
	if !ok || body == i+4 {
		return 0, 0, false
	}
	return body, scanExecClose(sql, body, nil), true
}
