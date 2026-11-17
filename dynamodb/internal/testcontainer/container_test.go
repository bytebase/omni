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

// TestCodecNamesTheRequestIDsBothEndpointsWrite holds the request ID an error's
// text carries to the stand-in. DynamoDB Local writes a UUID and the real
// service writes 52 upper-case letters and digits; a recording from either
// must replace what it wrote, or two recordings of one failure differ.
func TestCodecNamesTheRequestIDsBothEndpointsWrite(t *testing.T) {
	c := codec{}
	for _, in := range []string{
		`RequestID: 123e4567-e89b-12d3-a456-426614174000`,
		`RequestID: 64MQM2JG6L55RB5H2B2F65IA5JVV4KQNSO5AEMVJF66Q9ASUAAJG`,
		`"nextToken":"vBTmMeYwW6Y0+/="`,
	} {
		if !c.MadeUp().MatchString(in) {
			t.Errorf("MadeUp does not match %q", in)
		}
	}
	for _, in := range []string{
		`RequestID: 123`,
		`RequestID: lower-case-request-id`,
	} {
		if c.MadeUp().MatchString(in) {
			t.Errorf("MadeUp matches %q, which names no request", in)
		}
	}
	if got, want := c.StandIn(`RequestID: x`, 3), "RequestID: 00000000-0000-0000-0000-000000000003"; got != want {
		t.Errorf("StandIn = %q, want %q", got, want)
	}
}
