package catalog_test

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/mongo/catalog"
)

func TestNewCatalogEmpty(t *testing.T) {
	cat := catalog.New()
	if got := cat.Collections(); len(got) != 0 {
		t.Errorf("new catalog Collections() = %v, want empty", got)
	}
}

func TestAddAndListCollections(t *testing.T) {
	cat := catalog.New()
	cat.AddCollection("users")
	cat.AddCollection("orders")

	got := cat.Collections()
	want := []string{"orders", "users"}
	if !slices.Equal(got, want) {
		t.Errorf("Collections() = %v, want %v", got, want)
	}
}

func TestAddCollectionDedup(t *testing.T) {
	cat := catalog.New()
	cat.AddCollection("users")
	cat.AddCollection("users")
	cat.AddCollection("users")

	got := cat.Collections()
	if len(got) != 1 {
		t.Errorf("Collections() = %v, want 1 entry", got)
	}
}

func TestCollectionsSortOrder(t *testing.T) {
	cat := catalog.New()
	cat.AddCollection("zebra")
	cat.AddCollection("alpha")
	cat.AddCollection("middle")

	got := cat.Collections()
	want := []string{"alpha", "middle", "zebra"}
	if !slices.Equal(got, want) {
		t.Errorf("Collections() = %v, want %v", got, want)
	}
}
