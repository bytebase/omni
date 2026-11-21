package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/enginetest"
)

// TestGeneratedOutputsMatchSnapshot proves redis/testdata/RedisCommands.g4 and
// redis/parser/commands_gen.go match what commandgrammar emits from the
// committed redis/testdata/commands.tsv snapshot with zero drift.
func TestGeneratedOutputsMatchSnapshot(t *testing.T) {
	root := filepath.Join("..", "..", "..", "..")
	pin, err := enginetest.LoadPin(filepath.Join(root, defaultPins), "default")
	require.NoError(t, err)

	rawTSV, err := os.ReadFile(filepath.Join(root, defaultCommandsTSV))
	require.NoError(t, err)
	entries, err := parseCommandsTSV(string(rawTSV))
	require.NoError(t, err)

	wantGrammar, err := renderGrammar(pin, entries)
	require.NoError(t, err)
	gotGrammar, err := os.ReadFile(filepath.Join(root, defaultGrammarOut))
	require.NoError(t, err)
	require.Equal(t, wantGrammar, string(gotGrammar),
		"RedisCommands.g4 drifted from commands.tsv; regenerate with: go run ./redis/parser/cmd/commandgrammar")

	wantGo, err := renderGoRegistry(pin, entries)
	require.NoError(t, err)
	gotGo, err := os.ReadFile(filepath.Join(root, defaultGoOut))
	require.NoError(t, err)
	require.Equal(t, string(wantGo), string(gotGo),
		"redis/parser/commands_gen.go drifted from commands.tsv; regenerate with: go run ./redis/parser/cmd/commandgrammar")
}
