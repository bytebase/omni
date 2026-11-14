package parser

// checkJMXObjectName is what Parser.g's canonicalizeObjectName (Parser.g:162)
// refuses in the name of MBEAN 'name' or MBEANS 'name' while it parses: an
// empty name, the *:* wildcard, a name javax.management.ObjectName does not
// accept, and a name that is no pattern after MBEANS. It returns the
// refusal, or "" for a name Cassandra takes.
func checkJMXObjectName(s string, plural bool) string {
	if s == "" {
		return "Empty JMX object name supplied"
	}
	if s == "*:*" {
		return "Please use ALL MBEANS instead of wildcard pattern"
	}
	pattern, ok := parseJMXObjectName(s)
	if !ok {
		return s + " is not a valid JMX object name"
	}
	if plural && !pattern {
		return "Plural form used, but non-pattern JMX object name specified (" + s + ")"
	}
	return ""
}

// parseJMXObjectName follows javax.management.ObjectName.construct(String)
// (OpenJDK 17): domain ':' key '=' value (',' key '=' value)*, with '*' as
// a property-list wildcard and '*' or '?' making the domain or a value a
// pattern. It reports whether the name is a pattern and whether it is well
// formed.
func parseJMXObjectName(name string) (pattern, ok bool) {
	chars := []rune(name)
	n := len(chars)
	i := 0

	// The domain runs to the first ':'.
domain:
	for i < n {
		switch chars[i] {
		case ':':
			i++
			break domain
		case '=':
			// ObjectName looks for a ':' after an '=' in the domain.
			i++
			j := i
			for j < n && chars[j] != ':' {
				j++
			}
			if j == n {
				return false, false // Domain part must be specified
			}
		case '\n':
			return false, false
		case '*', '?':
			pattern = true
			i++
		default:
			i++
		}
	}
	if i == n {
		return false, false // Key properties cannot be empty (or no ':' at all)
	}

	keys := map[string]bool{}
	listPattern := false
	for i < n {
		if chars[i] == '*' {
			if listPattern {
				return false, false // several '*' in the property list
			}
			listPattern = true
			pattern = true
			i++
			if i < n && chars[i] != ',' {
				return false, false // end of name or ',' expected after '*'
			}
			if i == n {
				break
			}
			i++
			continue
		}

		keyStart := i
		if chars[i] == '=' {
			return false, false // Invalid key (empty)
		}
		for i < n {
			c := chars[i]
			i++
			if c == '=' {
				break
			}
			switch c {
			case '*', '?', ',', ':', '\n':
				return false, false // invalid character in key
			}
		}
		if chars[i-1] != '=' {
			return false, false // Unterminated key property part
		}
		key := string(chars[keyStart : i-1])

		if i < n && chars[i] == '"' {
			i++
			for i < n && chars[i] != '"' {
				switch chars[i] {
				case '\\':
					i++
					if i == n {
						return false, false // Unterminated quoted value
					}
					switch chars[i] {
					case '\\', '"', '?', '*', 'n':
					default:
						return false, false // Invalid escape sequence
					}
				case '\n':
					return false, false // Newline in quoted value
				case '?', '*':
					pattern = true
				}
				i++
			}
			if i == n {
				return false, false // Unterminated quoted value
			}
			i++ // closing quote
		} else {
			for i < n && chars[i] != ',' {
				switch chars[i] {
				case '*', '?':
					pattern = true
				case '=', ':', '"', '\n':
					return false, false // invalid character in value
				}
				i++
			}
		}
		if i == n-1 {
			return false, false // Invalid ending character or comma
		}
		i++ // the ',' after the value

		if keys[key] {
			return false, false // key already defined
		}
		keys[key] = true
	}
	return pattern, true
}
