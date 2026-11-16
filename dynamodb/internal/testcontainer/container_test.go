package testcontainer

import (
	"slices"
	"testing"
)

// TestFixture holds the fixture's premises, the properties the driver's tests
// draw their conclusions from. A change to a constant here changes what those
// tests prove even where they still pass, so it fails here first.
func TestFixture(t *testing.T) {
	t.Run("the paged fixture takes more than one listing page", func(t *testing.T) {
		// ListTables answers at most 100 names, whatever its limit, so a listing
		// that reaches the second page is the only one that proves the walk.
		if PagedTables <= 100 {
			t.Fatalf("the paged fixture holds %d tables, and one listing page holds 100", PagedTables)
		}
	})

	t.Run("the paged table names are the listing's order", func(t *testing.T) {
		names := PagedTableNames()
		if len(names) != PagedTables {
			t.Fatalf("PagedTableNames holds %d names, but the fixture holds %d tables", len(names), PagedTables)
		}
		if !slices.IsSorted(names) {
			t.Fatalf("the table names are %v, which a listing does not answer in", names)
		}
		for i := 1; i < len(names); i++ {
			if names[i] == names[i-1] {
				t.Fatalf("two tables are named %q", names[i])
			}
		}
	})

	t.Run("the paged fixture's rows take more than one read", func(t *testing.T) {
		if PagedCount <= 2 {
			t.Fatalf("the paged fixture holds %d rows, which one read of 2 items answers", PagedCount)
		}
	})

	t.Run("the people fixture's items are its rows", func(t *testing.T) {
		if PeopleCount != 3 {
			t.Fatalf("the people fixture holds %d rows, but its items are three", PeopleCount)
		}
		for _, id := range []string{PersonWithEveryType, PersonWithFewAttributes, PersonWithALateAttribute} {
			if id == "" {
				t.Fatal("an item has no id")
			}
		}
	})
}
