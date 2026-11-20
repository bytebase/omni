package catalog

import (
	"errors"
	"testing"
)

func TestFunctionalIndexAliasedArgValidation(t *testing.T) {
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE t (a INT, KEY k ((COALESCE(RAND() x, a))));", nil)
	if err != nil {
		t.Fatal(err)
	}
	last := results[len(results)-1]
	var catErr *Error
	if !errors.As(last.Error, &catErr) || catErr.Code != 3757 {
		t.Fatalf("got err=%v, want catalog Error 3757 for RAND() inside aliased arg", last.Error)
	}

	// Lambda parameter inside CHECK constraint is not treated as a missing table column.
	results, err = c.Exec("CREATE TABLE t2 (arr INT[], CHECK (ARRAY_LENGTH(ARRAY_FILTER(x -> x > 0, arr)) > 0));", nil)
	if err != nil {
		t.Fatal(err)
	}
	if last := results[len(results)-1]; last.Error != nil {
		t.Fatalf("unexpected error on lambda inside CHECK: %v", last.Error)
	}
}
