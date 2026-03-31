// Package analysis extracts field-level information from parsed CosmosDB queries.
package analysis

import "strings"

// Selector represents one step in a document property path.
type Selector struct {
	Name       string // property name
	ArrayIndex int    // -1 for item access (.name), >= 0 for array index ([n])
}

// ItemSelector creates a Selector for property access.
func ItemSelector(name string) Selector {
	return Selector{Name: name, ArrayIndex: -1}
}

// ArraySelector creates a Selector for array index access.
func ArraySelector(name string, index int) Selector {
	return Selector{Name: name, ArrayIndex: index}
}

// IsArray returns true if this selector represents an array index access.
func (s Selector) IsArray() bool {
	return s.ArrayIndex >= 0
}

// FieldPath represents a path through a JSON document.
type FieldPath []Selector

// String returns a human-readable representation like "container.addresses[1].country".
func (fp FieldPath) String() string {
	var sb strings.Builder
	for i, s := range fp {
		if s.IsArray() {
			sb.WriteString("[" + itoa(s.ArrayIndex) + "]")
		} else {
			if i > 0 {
				sb.WriteByte('.')
			}
			sb.WriteString(s.Name)
		}
	}
	return sb.String()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
