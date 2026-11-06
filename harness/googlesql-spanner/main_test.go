package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// emulatorImage is the emulator build the recorded verdicts came from; keep it
// equal to googlesql/internal/spannertest.Image.
const emulatorImage = "gcr.io/cloud-spanner-emulator/emulator:1.5.54"

// TestMain points the tests at $SPANNER_EMULATOR_HOST when set, and otherwise at
// an emulator testcontainer it starts and removes.
func TestMain(m *testing.M) {
	if os.Getenv("SPANNER_EMULATOR_HOST") != "" {
		os.Exit(m.Run())
	}
	c, host, err := startEmulator(context.Background())
	if err != nil {
		_ = testcontainers.TerminateContainer(c)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Setenv("SPANNER_EMULATOR_HOST", host)
	code := m.Run()
	_ = testcontainers.TerminateContainer(c)
	os.Exit(code)
}

func startEmulator(ctx context.Context) (testcontainers.Container, string, error) {
	c, err := testcontainers.Run(ctx, emulatorImage,
		testcontainers.WithExposedPorts("9010/tcp"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("9010/tcp").WithStartupTimeout(2*time.Minute)),
	)
	if err != nil {
		return c, "", fmt.Errorf("start Spanner emulator: %w", err)
	}
	host, err := c.Host(ctx)
	if err != nil {
		return c, "", err
	}
	port, err := c.MappedPort(ctx, "9010/tcp")
	if err != nil {
		return c, "", err
	}
	return c, host + ":" + port.Port(), nil
}
