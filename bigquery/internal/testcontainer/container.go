// Package testcontainer is how bigquery's driver tests reach BigQuery: through
// a tape (enginetest/tape) that replays a pinned dataset's answers, or, in a
// recording (OMNI_RECORD=1), through the pinned Google Cloud project and dataset
// itself. A replay starts nothing and needs no credentials.
//
// BigQuery is a Google Cloud service with no container image and no local
// emulator. The tape sits at the driver's client seam (bigquery/driver's
// client): a request is one call the driver makes of its client — Ping,
// ListDatasets, DatasetMetadata, ListTables, TableMetadata, Read, Run, DryRun,
// and a job's Wait, Read and Status — and an answer is the narrowed metadata,
// typed rows, job statistics or error the service returned. The driver's own
// tests own the exchange's spelling because they own the seam's types.
//
// A recording reads BIGQUERY_CREDENTIALS, BIGQUERY_PROJECT and BIGQUERY_DATASET
// from the environment and fails if they are missing rather than skipping, so
// OMNI_RECORD=1 can never write an empty tape and call it a pass. It verifies
// the live dataset's identity — "<id> <location> <creationTime>" SHA-256 and
// "gcp:bigquery:<project>:<location>" — against bigquery/testdata/pins.tsv so
// a recording cannot silently write answers under the wrong pin. Credentials
// are never written to a tape or a pin.
package testcontainer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is bigquery/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Driver is the pin role of the Google Cloud project and dataset the driver's
// tapes were recorded against.
const Driver = "driver"

// Project and Dataset are the pinned project and scratch dataset the driver's
// tapes were recorded against. Replay uses them directly so offline tests need
// no environment variables.
const (
	Project = "bytebase-dev"
	Dataset = "omni_record"
)

// Environment variables read during recording (OMNI_RECORD=1).
const (
	CredentialVariable = "BIGQUERY_CREDENTIALS"
	ProjectVariable    = "BIGQUERY_PROJECT"
	DatasetVariable    = "BIGQUERY_DATASET"

	// CredentialDefault selects the Google authentication library's default
	// credential chain (environment, well-known file, or attached service
	// account). Because it names no key file, BIGQUERY_PROJECT must be set.
	CredentialDefault = "default"
)

// Server is the tape server the driver's recordings are written to. Answers
// carry their own structured errors, so the codec holds no error spelling.
var Server = &tape.Server{Pins: Pins, Codec: codec{}}

type codec struct{}

func (codec) EncodeError(error) (string, bool) { return "", false }

func (codec) DecodeError(text string) error { return fmt.Errorf("%s", text) } //nolint:revive // an answer's text, not a cause to wrap

// LiveConfig holds the resolved live connection settings for a recording.
type LiveConfig struct {
	Project  string
	Dataset  string
	Password string // resolved credential JSON, or empty for the default chain
}

// Live resolves the recording environment variables and fails t if any required
// setting is missing or unreadable. It never skips: a recording without a
// project must fail rather than write an empty tape.
func Live(t testing.TB) LiveConfig {
	t.Helper()
	cfg, err := ResolveLiveConfig(os.Getenv)
	if err != nil {
		t.Fatalf("testcontainer: %v", err)
	}
	return cfg
}

// ResolveLiveConfig resolves BIGQUERY_CREDENTIALS, BIGQUERY_PROJECT and
// BIGQUERY_DATASET using lookup.
func ResolveLiveConfig(lookup func(string) string) (LiveConfig, error) {
	credentials := strings.TrimSpace(lookup(CredentialVariable))
	if credentials == "" {
		return LiveConfig{}, fmt.Errorf("recording needs %s (a key path, key JSON, or %q)", CredentialVariable, CredentialDefault)
	}
	project := strings.TrimSpace(lookup(ProjectVariable))
	var document string
	if credentials == CredentialDefault {
		if project == "" {
			return LiveConfig{}, fmt.Errorf("%s is %q and names no key to read a project from, so %s must be set", CredentialVariable, CredentialDefault, ProjectVariable)
		}
	} else {
		document = credentials
		if !strings.HasPrefix(credentials, "{") {
			key, err := os.ReadFile(credentials)
			if err != nil {
				return LiveConfig{}, fmt.Errorf("%s names a file that cannot be read: %w", CredentialVariable, err)
			}
			document = string(key)
		}
		if project == "" {
			var key struct {
				ProjectID string `json:"project_id"`
			}
			if err := json.Unmarshal([]byte(document), &key); err != nil {
				return LiveConfig{}, fmt.Errorf("%s does not parse as a credential: %w", CredentialVariable, err)
			}
			project = strings.TrimSpace(key.ProjectID)
			if project == "" {
				return LiveConfig{}, fmt.Errorf("%s is not set and the key names no project_id", ProjectVariable)
			}
		}
	}
	dataset := strings.TrimSpace(lookup(DatasetVariable))
	if dataset == "" {
		return LiveConfig{}, fmt.Errorf("%s must be set: a recording runs against a scratch dataset of its own", DatasetVariable)
	}
	return LiveConfig{
		Project:  project,
		Dataset:  dataset,
		Password: document,
	}, nil
}

// DatasetDigest returns the pin digest bigquery/testdata/pins.tsv holds for a
// dataset: the SHA-256 of "<id> <location> <creationTime>" as datasets.get
// states it, with creationTime in Unix milliseconds.
func DatasetDigest(fullID, location string, creationTime time.Time) string {
	payload := fmt.Sprintf("%s %s %d", fullID, location, creationTime.UnixMilli())
	sum := sha256.Sum256([]byte(payload))
	return "sha256:" + hex.EncodeToString(sum[:])
}

// DatasetImage returns the pin image column for project and location.
func DatasetImage(project, location string) string {
	return fmt.Sprintf("gcp:bigquery:%s:%s", project, location)
}

var verifiedPins = struct {
	sync.Mutex
	ok map[string]bool
}{ok: map[string]bool{}}

// VerifyPin checks that the live dataset matches pin's image and SHA-256 digest
// in bigquery/testdata/pins.tsv. Successful verifications are cached per process
// so a multi-test recording asks datasets.get for the pin check once.
func VerifyPin(ctx context.Context, client *bigquery.Client, pin enginetest.Pin, project, dataset string) error {
	key := pin.String() + "|" + project + "|" + dataset
	verifiedPins.Lock()
	already := verifiedPins.ok[key]
	verifiedPins.Unlock()
	if already {
		return nil
	}

	meta, err := client.Dataset(dataset).Metadata(ctx)
	if err != nil {
		return fmt.Errorf("read dataset %s:%s metadata to verify pin: %w", project, dataset, err)
	}
	gotImage := DatasetImage(project, meta.Location)
	if gotImage != pin.Image {
		return fmt.Errorf("live dataset is %q, but %s pins %q", gotImage, Pins, pin.Image)
	}
	gotDigest := DatasetDigest(meta.FullID, meta.Location, meta.CreationTime)
	if gotDigest != pin.Digest {
		return fmt.Errorf("live dataset %s (%s, created %d) has digest %s, but %s pins %s",
			meta.FullID, meta.Location, meta.CreationTime.UnixMilli(), gotDigest, Pins, pin.Digest)
	}

	verifiedPins.Lock()
	verifiedPins.ok[key] = true
	verifiedPins.Unlock()
	return nil
}
