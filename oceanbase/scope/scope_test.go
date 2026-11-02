package scope

import (
	"math/rand"
	"strings"
	"testing"
)

// TestLookupsMatchScans checks FirstEqualFold and FirstWithColumn against the
// scans over AllEntries they replace, on names spelled with runes whose case
// folding is not ASCII's: the Kelvin sign folds with k, the long s with s,
// and invalid bytes compare as U+FFFD.
func TestLookupsMatchScans(t *testing.T) {
	alphabet := []string{"a", "A", "k", "K", "K", "s", "S", "ſ", "ß", "ẞ", "σ", "ς", "Σ", "é", "É", "\xff", "\xfe", "�", "_", "1"}
	r := rand.New(rand.NewSource(1))
	word := func() string {
		var b strings.Builder
		for i := 0; i < 1+r.Intn(3); i++ {
			b.WriteString(alphabet[r.Intn(len(alphabet))])
		}
		return b.String()
	}
	for round := 0; round < 200; round++ {
		s := New()
		for i := 0; i < 1+r.Intn(12); i++ {
			var cols []Column
			for j := 0; j < r.Intn(4); j++ {
				cols = append(cols, Column{Name: word(), Position: j + 1})
			}
			s.Add(word(), &Table{Columns: cols})
			for probe := 0; probe < 20; probe++ {
				name := word()
				wantName, wantFound := "", false
				for _, e := range s.AllEntries() {
					if strings.EqualFold(e.Name, name) {
						wantName, wantFound = e.Name, true
						break
					}
				}
				if e, found := s.FirstEqualFold(name); found != wantFound || e.Name != wantName {
					t.Fatalf("FirstEqualFold(%q) = %q, %v; scan finds %q, %v", name, e.Name, found, wantName, wantFound)
				}
				wantName, wantFound = "", false
				for _, e := range s.AllEntries() {
					if e.Table.GetColumn(name) != nil {
						wantName, wantFound = e.Name, true
						break
					}
				}
				if e, found := s.FirstWithColumn(name); found != wantFound || e.Name != wantName {
					t.Fatalf("FirstWithColumn(%q) = %q, %v; scan finds %q, %v", name, e.Name, found, wantName, wantFound)
				}
			}
		}
	}
}
