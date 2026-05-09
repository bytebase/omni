package parser

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type splitGoldenCase struct {
	Description string                 `yaml:"description"`
	Input       string                 `yaml:"input"`
	Parse       bool                   `yaml:"parse,omitempty"`
	Statements  []splitGoldenStatement `yaml:"statements"`
}

type splitGoldenStatement struct {
	Text string `yaml:"text"`
}

func TestSplitGoldenCorpus(t *testing.T) {
	data, err := os.ReadFile("testdata/splitter.yaml")
	if err != nil {
		t.Fatalf("read splitter corpus: %v", err)
	}

	var cases []splitGoldenCase
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatalf("unmarshal splitter corpus: %v", err)
	}
	if len(cases) < 80 {
		t.Fatalf("splitter corpus has %d cases, want at least 80", len(cases))
	}

	for _, tc := range cases {
		t.Run(tc.Description, func(t *testing.T) {
			got := Split(tc.Input)
			if len(got) != len(tc.Statements) {
				t.Fatalf("got %d segments %q, want %d", len(got), splitTexts(got), len(tc.Statements))
			}

			prevEnd := 0
			for i, seg := range got {
				if seg.ByteStart < prevEnd {
					t.Fatalf("segment[%d] starts at %d before previous end %d", i, seg.ByteStart, prevEnd)
				}
				if seg.ByteStart < 0 || seg.ByteEnd < seg.ByteStart || seg.ByteEnd > len(tc.Input) {
					t.Fatalf("segment[%d] invalid range [%d,%d] for input length %d", i, seg.ByteStart, seg.ByteEnd, len(tc.Input))
				}
				if seg.Text != tc.Input[seg.ByteStart:seg.ByteEnd] {
					t.Fatalf("segment[%d] Text does not match input range [%d,%d]", i, seg.ByteStart, seg.ByteEnd)
				}
				if seg.Empty() {
					t.Fatalf("segment[%d] is empty: %q", i, seg.Text)
				}
				if seg.Text != tc.Statements[i].Text {
					t.Fatalf("segment[%d] = %q, want %q", i, seg.Text, tc.Statements[i].Text)
				}
				if tc.Parse {
					if _, err := Parse(strings.Repeat(" ", seg.ByteStart) + seg.Text); err != nil {
						t.Fatalf("segment[%d] parse error: %v\ntext:\n%s", i, err, seg.Text)
					}
				}
				prevEnd = seg.ByteEnd
			}
		})
	}
}
