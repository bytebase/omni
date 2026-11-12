package parser

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of TiDB's parser.y and its near-misses, with omni's verdict and
// TiDB's. It reads the file and starts no server. OMNI_RECORD=1 judges the
// rows that have no verdict yet on the pinned TiDB.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: filepath.Join("..", "testdata", "parser.y"),
		Tokens:  filepath.Join("..", "testdata", "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// compatJudge asks the package's TiDB server, through PREPARE, whether it
// accepts each statement's syntax. It fails t unless the server is the pinned
// image and release.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	if pin.Image != tidbImage || pin.Tag != tidbCorpusVersion {
		t.Fatalf("pins.tsv pins %s at %s, but this package starts %s at %s", pin.Image, pin.Tag, tidbImage, tidbCorpusVersion)
	}
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	assertImageDigest(t, tc, pin.Digest)
	previous := "<start of compat.tsv>"
	return func(sql string) (string, error) {
		accepts, _, err := tc.syntaxVerdict(t, sql, previous)
		previous = sql
		if err != nil {
			return "", err
		}
		if accepts {
			return compat.Accept, nil
		}
		return compat.Reject, nil
	}
}

// assertImageDigest fails t unless the container runs the image with the
// pinned repository digest. A tag can be pushed again; a digest cannot.
func assertImageDigest(t testing.TB, tc *tidbContainer, digest string) {
	t.Helper()
	info, err := tc.container.Inspect(tc.ctx)
	if err != nil {
		t.Fatalf("inspect the TiDB container: %v", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(tc.ctx)
	if err != nil {
		t.Fatalf("docker client: %v", err)
	}
	defer cli.Close()
	img, err := cli.ImageInspect(tc.ctx, info.Image)
	if err != nil {
		t.Fatalf("inspect the TiDB image: %v", err)
	}
	for _, d := range img.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return
		}
	}
	t.Fatalf("the TiDB container runs an image with digests %v, but pins.tsv pins %s", img.RepoDigests, digest)
}
