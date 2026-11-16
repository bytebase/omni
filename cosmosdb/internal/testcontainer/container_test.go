package testcontainer

import (
	"slices"
	"testing"
)

// TestFixture holds the fixture's premises, the properties the driver's tests
// draw their conclusions from. A change to a constant here changes what those
// tests prove even where they still pass, so it fails here first.
func TestFixture(t *testing.T) {
	t.Run("two partition key values pass the page cap and six drain in one page", func(t *testing.T) {
		if len(PagedBigBuckets) == 0 {
			t.Fatal("no partition key value passes the page cap, so no range has to page")
		}
		for _, bucket := range PagedBigBuckets {
			if PagedBucketRows(bucket) <= PageSize {
				t.Fatalf("%q holds %d rows, not more than the %d-row page cap", bucket, PagedBucketRows(bucket), PageSize)
			}
		}
		for _, bucket := range PagedSmallBuckets {
			if PagedBucketRows(bucket) > PageSize {
				t.Fatalf("%q holds %d rows, more than the %d-row page cap", bucket, PagedBucketRows(bucket), PageSize)
			}
		}
	})

	t.Run("the paged fixture's rows and shapes line up", func(t *testing.T) {
		want := len(PagedBigBuckets)*PagedBigRows + len(PagedSmallBuckets)*PagedSmallRows
		if got := PagedRows(); got != want {
			t.Fatalf("PagedRows is %d, but the buckets hold %d rows", got, want)
		}
		seqs := make([]int, 0, len(PagedShapes))
		for _, shape := range PagedShapes {
			if shape.Seq < 0 || shape.Seq >= PagedRows() {
				t.Fatalf("shape %d is not a row of the fixture", shape.Seq)
			}
			seqs = append(seqs, shape.Seq)
		}
		if len(slices.Compact(slices.Sorted(slices.Values(seqs)))) != len(seqs) {
			t.Fatalf("two rows share a shape: %v", seqs)
		}
	})

	t.Run("the shapes reach every item-order class", func(t *testing.T) {
		// One row per class the service orders by: undefined (no row carries
		// it), null, bool, number, string, array, object.
		for _, want := range []string{`null`, `true`, `1`, `"a"`, `[1,2]`, `{"a":1}`} {
			if !slices.ContainsFunc(PagedShapes, func(shape Shape) bool { return shape.Value == want }) {
				t.Fatalf("no row carries v = %s", want)
			}
		}
	})

	t.Run("the cities fixture adds up as the tests compute it", func(t *testing.T) {
		perCountry := 0
		for i := 1; i <= CitiesPerCountry; i++ {
			perCountry += i * PopulationGap
		}
		if got, want := CitiesPopulationSum(), perCountry*len(Cities); got != want {
			t.Fatalf("CitiesPopulationSum is %d, want %d", got, want)
		}
		if got, want := CitiesItems(), len(Cities)*CitiesPerCountry; got != want {
			t.Fatalf("CitiesItems is %d, want %d", got, want)
		}
	})
}
