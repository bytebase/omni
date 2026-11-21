package analysis

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type splitCorpusCase struct {
	Description string               `yaml:"description"`
	Input       string               `yaml:"input"`
	Result      []splitCorpusSegment `yaml:"result"`
}

type splitCorpusSegment struct {
	Text  string           `yaml:"text"`
	Range splitCorpusRange `yaml:"range"`
	Empty bool             `yaml:"empty"`
}

type splitCorpusRange struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
}

func TestSplitCorpus(t *testing.T) {
	content, err := os.ReadFile("testdata/split.yaml")
	require.NoError(t, err)
	var cases []splitCorpusCase
	require.NoError(t, yaml.Unmarshal(content, &cases))
	require.NotEmpty(t, cases)

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			segments := omnipg.Split(tc.Input)
			require.Len(t, segments, len(tc.Result))
			for i, want := range tc.Result {
				got := segments[i]
				require.Equal(t, want.Text, got.Text, "segment %d text", i)
				require.Equal(t, want.Empty, got.Empty(), "segment %d empty", i)
				require.Equal(t, want.Range.Start, got.ByteStart, "segment %d start", i)
				require.Equal(t, want.Range.End, got.ByteEnd, "segment %d end", i)
			}
		})
	}
}

func TestSplitLargeInsertScriptScalesLinearly(t *testing.T) {
	const rowCount = 2000
	padding := strings.Repeat("x", 1024)
	var builder strings.Builder
	for i := range rowCount {
		fmt.Fprintf(&builder, "INSERT INTO perf_omni_pg (id, payload) VALUES (%d, '%s');\n", i, padding)
	}

	started := time.Now()
	segments := omnipg.Split(builder.String())
	elapsed := time.Since(started)

	var nonEmpty int
	for _, seg := range segments {
		if !seg.Empty() {
			nonEmpty++
		}
	}
	require.Equal(t, rowCount, nonEmpty)
	require.Less(t, elapsed, time.Second)
}
