// Package testcontainer is how cosmosdb's driver tests reach Cosmos DB: a
// pinned Cosmos DB emulator, started on first use in a recording
// (OMNI_RECORD=1) and removed by Terminate, which the package's TestMain
// calls. A replay starts nothing and needs no Docker.
//
// Cosmos DB for NoSQL has no server CI can start more than one of: the
// emulator publishes the address it serves in its own account document and
// takes most of a minute to hand out its partitions, so it binds one fixed
// port and one runs per machine.
//
// The package also holds the fixtures the driver's tests read and the tape
// server their recordings are written to. The tape itself sits in the driver's
// tests, at the driver's client seam (cosmosdb/driver.Client): a request is
// one call the driver makes of the client, and an answer is the names it
// listed, the items a query returned, or the client's error. Those tests own
// the exchange's spelling because they own the seam's types.
//
// Two fixtures are seeded, once per emulator, before the first request of a
// recording:
//
//	cities  the aggregate, DISTINCT, ORDER BY, TOP, OFFSET/LIMIT and GROUP BY
//	        corpus, whose 8 partition key values each hold 5 cities
//	paged   2,130 rows over 8 partition key values, two of them past the
//	        gateway's per-page cap, for the query classes a range can only
//	        deliver over several pages
//
// Both are written through the client, not through a taped request: a replay
// answers as if they were there, and the seeding is not part of what the driver
// asked.
package testcontainer

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Pins is cosmosdb/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Emulator is the pin role of the one Cosmos DB emulator the driver's tests
// ask. The parser's answers come from the grammar pin in the same file, which
// is role default and starts nothing.
const Emulator = "emulator"

// EmulatorKey is the well-known master key the emulator reads from every
// request. https://learn.microsoft.com/en-us/azure/cosmos-db/emulator
const EmulatorKey = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="

// Port is the port the emulator serves and advertises in its account document.
// It is published unchanged, so only one emulator runs on a machine.
const Port = "8081"

// Server is the tape server the driver's recordings are written to. Its pins
// name the emulator, and its codec holds no error spelling: an answer carries
// the client's own error, so nothing else is ever answered by an error.
var Server = &tape.Server{Pins: Pins, Codec: codec{}}

type codec struct{}

func (codec) EncodeError(error) (string, bool) { return "", false }

func (codec) DecodeError(text string) error { return errors.New(text) }

// InsecureTLS accepts the certificate the emulator signs itself, at the oldest
// TLS version it serves. It is what a driver's Config.TLS returns for the
// emulator, and what the fixtures are written over.
func InsecureTLS(string) (*tls.Config, error) {
	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true, //nolint:gosec // the emulator signs its certificate itself
	}, nil
}

var shared = struct {
	sync.Mutex
	container testcontainers.Container
	endpoint  string
	err       error
	seeded    bool
}{}

// Endpoint starts the pinned emulator on first use and returns the account
// endpoint, with the fixtures in place. It reaches the emulator, so only a
// recording calls it.
func Endpoint(t testing.TB) (string, error) {
	t.Helper()
	pin, err := enginetest.LoadPin(Pins, Emulator)
	if err != nil {
		return "", err
	}
	shared.Lock()
	defer shared.Unlock()
	if shared.endpoint == "" && shared.err == nil {
		shared.endpoint, shared.err = startEmulator(pin)
	}
	if shared.err != nil {
		return "", shared.err
	}
	if !shared.seeded {
		if err := seed(shared.endpoint); err != nil {
			return "", err
		}
		shared.seeded = true
	}
	return shared.endpoint, nil
}

// startEmulator runs the pinned emulator and waits for its gateway.
func startEmulator(pin enginetest.Pin) (endpoint string, err error) {
	// testcontainers panics when Docker is missing; fail the test instead.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	emulatorPort := nat.Port(Port + "/tcp")
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.Image,
			ExposedPorts: []string{string(emulatorPort)},
			Env: map[string]string{
				"AZURE_COSMOS_EMULATOR_PARTITION_COUNT":         "8",
				"AZURE_COSMOS_EMULATOR_ENABLE_DATA_PERSISTENCE": "false",
				"AZURE_COSMOS_EMULATOR_IP_ADDRESS_OVERRIDE":     "127.0.0.1",
			},
			// Readiness is an answer from the gateway rather than a log line: the
			// emulator logs "Started N/M partitions" per partition before the bare
			// "Started" that means it is up, and terminates its lines with CRLF.
			//
			// Bringing up the partitions takes about half a minute on an idle host
			// and was measured at 216 s on a loaded one, so the timeout clears the
			// longest boot seen rather than the typical one.
			WaitingFor: wait.ForHTTP("/_explorer/emulator.pem").
				WithPort(emulatorPort).
				WithTLS(true, &tls.Config{
					MinVersion:         tls.VersionTLS12,
					InsecureSkipVerify: true, //nolint:gosec // the emulator signs its certificate itself
				}).
				WithStartupTimeout(10 * time.Minute),
			// The emulator names its own endpoint in the database account document
			// and the client follows that rather than the endpoint it was configured
			// with. The address override makes it advertise 127.0.0.1, which only
			// reaches the emulator while the host publishes the same port it serves
			// — hence the fixed binding, and one emulator per machine. A second
			// concurrent start fails on that binding instead of sharing the first
			// one's databases.
			HostConfigModifier: func(hc *dockercontainer.HostConfig) {
				hc.PortBindings = nat.PortMap{emulatorPort: []nat.PortBinding{{HostPort: Port}}}
			},
		},
		Started: true,
	})
	if c != nil {
		shared.container = c
	}
	if err != nil {
		return "", fmt.Errorf("start the Cosmos DB emulator: %w", err)
	}
	if err := checkDigest(ctx, c, pin.Digest); err != nil {
		return "", err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("the emulator's host: %w", err)
	}
	return fmt.Sprintf("https://%s:%s", host, Port), nil
}

// checkDigest holds the container to the pin. An image reference is a tag,
// which a registry may move; the digest is what the answers were recorded from.
func checkDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the Cosmos DB emulator container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	image, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the Cosmos DB emulator image: %w", err)
	}
	for _, d := range image.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("the Cosmos DB emulator runs an image with digests %v, but pins.tsv pins %s", image.RepoDigests, digest)
}

// Terminate removes the emulator this run started, if any. A package's
// TestMain calls it.
func Terminate() {
	if shared.container != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		_ = shared.container.Terminate(ctx, testcontainers.StopTimeout(0))
		shared.container = nil
	}
}

// The fixture the driver's tests read. A recording seeds it; a replay answers
// as if it were there.
const (
	// Database holds both fixture containers.
	Database = "cosmosfixture"
	// SpareDatabase is a database the fixture puts no container in, so an
	// instance sync reports a database whose container list is empty.
	SpareDatabase = "cosmosspare"
	// CitiesContainer holds CitiesPerCountry cities for each of Cities, one
	// partition each, so every aggregate over it has one correct answer.
	CitiesContainer = "cities"
	// PagedContainer holds PagedRows rows over eight partition key values, two
	// of them past the gateway's per-page cap.
	PagedContainer = "paged"
	// CitiesThroughput is the least throughput that spreads eight partition key
	// values over several physical partitions: the emulator hands out one
	// partition per 10k RU, and the gateway serves a single-partition query
	// itself, so the merge would never run below this.
	CitiesThroughput = 30000
	// PagedThroughput leaves the cities fixture the rest of the emulator's
	// eight-partition budget.
	PagedThroughput = 12000
	// PageSize is the gateway's own per-page item cap, measured both on the
	// emulator and on a real Cosmos account.
	PageSize = 1000
	// PagedBigRows puts one partition key value past the cap, so the range
	// holding it must be read over several pages.
	PagedBigRows = PageSize + 50
	// PagedSmallRows keeps the other partition key values short, so a merge
	// has ranges that drain while the big ones still have pages to come.
	PagedSmallRows = 5
	// CitiesPerCountry is the rows each Cities partition key value holds.
	CitiesPerCountry = 5
	// PopulationGap is the step between the cities of one country, so every
	// aggregate over the fixture has one correct answer.
	PopulationGap = 100000
)

// Cities are the fixture's partition key values, each with CitiesPerCountry
// cities whose populations step by PopulationGap.
var Cities = []string{"AE", "BR", "CN", "DE", "FR", "IN", "JP", "US"}

// CitiesItems is the fixture's row count.
func CitiesItems() int { return len(Cities) * CitiesPerCountry }

// CitiesPopulationSum is the total SUM must return over the cities fixture.
func CitiesPopulationSum() int {
	perCountry := 0
	for i := 1; i <= CitiesPerCountry; i++ {
		perCountry += i * PopulationGap
	}
	return perCountry * len(Cities)
}

// PagedBigBuckets are the partition key values whose rows exceed the page cap.
var PagedBigBuckets = []string{"big-1", "big-2"}

// PagedSmallBuckets are the partition key values that drain in one page.
var PagedSmallBuckets = []string{"small-1", "small-2", "small-3", "small-4", "small-5", "small-6"}

// PagedRows is the fixture's row count. Every row carries a distinct seq, from
// 0 to PagedRows-1, so an ordered result has exactly one correct answer and
// DISTINCT and GROUP BY over seq have exactly PagedRows outputs.
func PagedRows() int {
	return len(PagedBigBuckets)*PagedBigRows + len(PagedSmallBuckets)*PagedSmallRows
}

// SortedPagedBuckets is the fixture's partition key values in ascending order.
func SortedPagedBuckets() []string {
	buckets := append(append([]string(nil), PagedBigBuckets...), PagedSmallBuckets...)
	slices.Sort(buckets)
	return buckets
}

// PagedBucketRows is how many rows one partition key value holds.
func PagedBucketRows(bucket string) int {
	if slices.Contains(PagedBigBuckets, bucket) {
		return PagedBigRows
	}
	return PagedSmallRows
}

// Shape is one row's `v`, one per shape a document can hold.
type Shape struct {
	Seq   int
	Value string
}

// PagedShapes gives eleven of the fixture's rows a `v`, chosen so that the
// values span the container's ranges: two rows in each big bucket and one in
// each small one. Every other row leaves `v` undefined, so a query grouping or
// ordering on it meets both.
//
// Ordering them is what the service does by type — undefined < null < bool <
// number < string < array < object — and the two seqs inside each of the array
// and object classes are interchangeable, since the service orders those by
// its index's own encoding.
var PagedShapes = []Shape{
	{Seq: 0, Value: `null`},
	{Seq: 1, Value: `false`},
	{Seq: PagedBigRows, Value: `true`},
	{Seq: PagedBigRows + 1, Value: `1`},
	{Seq: 2 * PagedBigRows, Value: `20`},
	{Seq: 2*PagedBigRows + 1, Value: `"a"`},
	{Seq: 2*PagedBigRows + PagedSmallRows, Value: `"b"`},
	{Seq: 2*PagedBigRows + 2*PagedSmallRows, Value: `[1,2]`},
	{Seq: 2*PagedBigRows + 3*PagedSmallRows, Value: `[2]`},
	{Seq: 2*PagedBigRows + 4*PagedSmallRows, Value: `{"a":1}`},
	{Seq: 2*PagedBigRows + 5*PagedSmallRows, Value: `{"b":1}`},
}

// ShapeOf returns the shape attached to a row, or "" when the row leaves `v`
// undefined.
func ShapeOf(seq int) string {
	for _, shape := range PagedShapes {
		if shape.Seq == seq {
			return shape.Value
		}
	}
	return ""
}

// seed writes the fixture through the client, once per emulator.
func seed(endpoint string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	credential, err := azcosmos.NewKeyCredential(EmulatorKey)
	if err != nil {
		return fmt.Errorf("emulator key credential: %w", err)
	}
	insecure, err := InsecureTLS("")
	if err != nil {
		return err
	}
	c, err := azcosmos.NewClientWithKey(endpoint, credential, &azcosmos.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: &http.Client{Transport: &http.Transport{TLSClientConfig: insecure}},
		},
	})
	if err != nil {
		return fmt.Errorf("emulator client: %w", err)
	}

	if _, err := c.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: Database}, nil); err != nil && !alreadyExists(err) {
		return fmt.Errorf("create the fixture database: %w", err)
	}
	if _, err := c.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: SpareDatabase}, nil); err != nil && !alreadyExists(err) {
		return fmt.Errorf("create the spare database: %w", err)
	}
	database, err := c.NewDatabase(Database)
	if err != nil {
		return fmt.Errorf("open the fixture database: %w", err)
	}
	if err := createContainer(ctx, database, CitiesContainer, "/country", CitiesThroughput); err != nil {
		return err
	}
	if err := createContainer(ctx, database, PagedContainer, "/bucket", PagedThroughput); err != nil {
		return err
	}
	cities, err := database.NewContainer(CitiesContainer)
	if err != nil {
		return fmt.Errorf("open the cities container: %w", err)
	}
	paged, err := database.NewContainer(PagedContainer)
	if err != nil {
		return fmt.Errorf("open the paged container: %w", err)
	}
	if err := seedCities(ctx, cities); err != nil {
		return err
	}
	return seedPaged(ctx, paged)
}

func createContainer(ctx context.Context, database *azcosmos.DatabaseClient, id, path string, throughput int32) error {
	properties := azcosmos.NewManualThroughputProperties(throughput)
	if _, err := database.CreateContainer(ctx, azcosmos.ContainerProperties{
		ID:                     id,
		PartitionKeyDefinition: azcosmos.PartitionKeyDefinition{Paths: []string{path}},
	}, &azcosmos.CreateContainerOptions{ThroughputProperties: &properties}); err != nil && !alreadyExists(err) {
		return fmt.Errorf("create the %s container: %w", id, err)
	}
	return nil
}

func seedCities(ctx context.Context, container *azcosmos.ContainerClient) error {
	for _, country := range Cities {
		for i := 0; i < CitiesPerCountry; i++ {
			item, err := json.Marshal(map[string]any{
				"id":         fmt.Sprintf("%s-%d", country, i),
				"country":    country,
				"name":       fmt.Sprintf("city-%s-%d", country, i),
				"population": (i + 1) * PopulationGap,
			})
			if err != nil {
				return fmt.Errorf("marshal a fixture city: %w", err)
			}
			if err := createItem(ctx, container, country, item); err != nil {
				return fmt.Errorf("seed %s-%d: %w", country, i, err)
			}
		}
	}
	return nil
}

// seedPaged writes the multi-page fixture. A few thousand writes are most of
// this harness's runtime, so they go out concurrently; the worker count is
// what the emulator answered without timing out under a loaded machine.
func seedPaged(ctx context.Context, container *azcosmos.ContainerClient) error {
	type row struct {
		bucket string
		seq    int
	}
	rows := make(chan row, PagedRows())
	seq := 0
	for _, bucket := range PagedBigBuckets {
		for i := 0; i < PagedBigRows; i++ {
			rows <- row{bucket: bucket, seq: seq}
			seq++
		}
	}
	for _, bucket := range PagedSmallBuckets {
		for i := 0; i < PagedSmallRows; i++ {
			rows <- row{bucket: bucket, seq: seq}
			seq++
		}
	}
	close(rows)

	var wg sync.WaitGroup
	errs := make(chan error, 1)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range rows {
				doc := map[string]any{"id": fmt.Sprintf("row-%d", r.seq), "bucket": r.bucket, "seq": r.seq}
				if shape := ShapeOf(r.seq); shape != "" {
					doc["v"] = json.RawMessage(shape)
				}
				item, err := json.Marshal(doc)
				if err == nil {
					err = createItem(ctx, container, r.bucket, item)
				}
				if err != nil {
					select {
					case errs <- fmt.Errorf("seed row %d: %w", r.seq, err):
					default:
					}
					return
				}
			}
		}()
	}
	wg.Wait()
	select {
	case err := <-errs:
		return err
	default:
	}
	return nil
}

// createItem writes one fixture item, retrying what a loaded emulator does to a
// write it cannot serve in time: a 408, a 429, a 503, or a dropped connection
// that surfaces as an EOF. An error that is not transient costs the attempts
// and then surfaces unchanged.
func createItem(ctx context.Context, container *azcosmos.ContainerClient, partitionKey string, item []byte) error {
	var err error
	for attempt := 0; attempt < 6; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(attempt) * 500 * time.Millisecond):
			}
		}
		if _, err = container.CreateItem(ctx, azcosmos.NewPartitionKeyString(partitionKey), item, nil); err == nil || alreadyExists(err) {
			return nil
		}
	}
	return err
}

func alreadyExists(err error) bool {
	var responseErr *azcore.ResponseError
	return errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusConflict
}
