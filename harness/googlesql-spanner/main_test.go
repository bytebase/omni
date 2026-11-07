package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/bytebase/omni/harness/googlesql-spanner/internal/testcontainer"
)

// TestMain points the tests at $SPANNER_EMULATOR_HOST when set, and otherwise at
// an emulator testcontainer it starts and removes.
func TestMain(m *testing.M) {
	if os.Getenv("SPANNER_EMULATOR_HOST") != "" {
		os.Exit(m.Run())
	}
	c, host, err := testcontainer.Start(context.Background())
	if err != nil {
		_ = testcontainer.TerminateContainer(c)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Setenv("SPANNER_EMULATOR_HOST", host)
	code := m.Run()
	_ = testcontainer.TerminateContainer(c)
	os.Exit(code)
}
