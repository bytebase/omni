package parser

// execOpener reports whether s[i:], which starts with "/*", opens an
// executable comment, "/*!" and any version digits after it, and the offset
// its body starts at.
func execOpener(s string, i int) (int, bool) {
	if i+2 >= len(s) || s[i] != '/' || s[i+1] != '*' || s[i+2] != '!' {
		return 0, false
	}
	j := i + 3
	for j < len(s) && s[j] >= '0' && s[j] <= '9' {
		j++
	}
	return j, true
}
