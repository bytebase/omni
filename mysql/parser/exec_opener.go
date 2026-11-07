package parser

// execOpener reports whether s[i:], which starts with "/*", opens an
// executable comment, "/*!" and any version digits after it, and the offset
// its body starts at.
func execOpener(s string, i int) (int, bool) {
	if i+2 >= len(s) || s[i] != '/' || s[i+1] != '*' || s[i+2] != '!' {
		return 0, false
	}
	j := i + 3
	version := 0
	for j < len(s) && s[j] >= '0' && s[j] <= '9' {
		version = version*10 + int(s[j]-'0')
		j++
		if version > 80413 {
			return 0, false
		}
	}
	// MySQL executes a versioned comment only at or above its version.
	// https://dev.mysql.com/doc/refman/8.4/en/comments.html
	return j, true
}
