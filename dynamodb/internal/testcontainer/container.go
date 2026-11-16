// Package testcontainer is how dynamodb's tests reach DynamoDB: through a tape
// (enginetest/tape) that replays a pinned endpoint's answers, or, in a
// recording (OMNI_RECORD=1), through the pinned endpoint itself. The endpoint
// is a DynamoDB Local container, started by the first test that records and
// removed by Terminate, which the package's TestMain calls. A replay starts
// nothing and needs no Docker.
//
// DynamoDB Local has no database/sql driver, so the tape cannot sit at that
// seam. The driver's own tests answer the client seam (dynamodb/driver's
// client) instead: a request is one call the driver makes of its client —
// listtables, describetable, executestatement, executetransaction, ping — and
// an answer is the names it listed, the items a statement read, or the client's
// error. Those tests own the exchange's spelling because they own the seam's
// types.
//
// Two fixtures are seeded, once per container, before the first request of a
// recording:
//
//	people  three items over a partition key and a sort key, whose attributes
//	        cover every DynamoDB type, an item that leaves most of them out, and
//	        an item that names one the others do not
//	paged   five items under one partition key, ordered by their sort key, for
//	        the reads that take several execute-statement calls
//
// Two more tables are created: orders, whose key has no sort key, and writes,
// which Execute's tests write to so that the fixture above stays as recorded.
//
// Tables are written through the AWS SDK client, not through a taped request: a
// replay answers as if they were there, and the seeding is not part of what the
// driver asked.
//
// A second, separate container holds PagedTables tables, because ListTables
// answers at most 100 names per page: the walk over a listing's pages cannot be
// taped on the fixture above, and giving the fixture 101 tables would bury its
// answer.
//
// DynamoDB Local reads any credential and needs one to sign a request, so the
// endpoint is reached with the placeholder key pair below. It also needs a
// region, which names what the signed request is for. The container runs with
// -sharedDb so that one endpoint presents one set of tables whatever key pair
// signs the request; without it DynamoDB Local keeps a separate namespace per
// key pair, which is neither what the service does nor something a test could
// see if the driver and a fixture writer held different keys.
//
// Two things an answer carries cannot be held as answered, and the server's
// codec gives each a stand-in: the token that resumes a PartiQL read, which
// names the moment the table was created, and the request ID an error's text
// carries, which names one request. A stand-in is what a recording compares and
// what a replay answers, so two recordings write the same bytes.
//
// Two orders the endpoint chooses are held in one canonical order instead: the
// items of a read that names no ORDER BY, whose table iteration differs between
// boots, and the elements of a set, which DynamoDB does not order at all. The
// driver hands both on as the endpoint answered them; only the tape is
// canonical, so a test that draws a conclusion from a position names the clause
// that decides it or reads a set as a set.
package testcontainer

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awssdk "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is dynamodb/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Default is the pin role of the one DynamoDB Local the tests ask.
const Default = "default"

// Port is the port DynamoDB Local serves.
const Port = "8000"

// Region is the region every request is signed for. DynamoDB Local reads the
// signature but not the region, which the AWS SDK requires all the same.
const Region = "us-east-1"

// The placeholder key pair every request is signed with. DynamoDB Local reads
// any credential and rejects a request that carries none, so these are values a
// test can hold and no account has. The key ID is letters and digits only
// because DynamoDB Local answers an access key ID holding a dash with
// UnrecognizedClientException.
const (
	AccessKeyID     = "omniLocalDynamoDB"
	SecretAccessKey = "omniLocalDynamoDBSecret"
)

// Server is the tape server the driver's recordings are written to.
var Server = &tape.Server{Pins: Pins, Codec: codec{}}

type codec struct{}

func (codec) EncodeError(error) (string, bool) { return "", false }

func (codec) DecodeError(text string) error { return fmt.Errorf("%s", text) } //nolint:revive // an answer's text, not a cause to wrap

// madeUp matches the two values an endpoint makes up that a tape holds as
// stand-ins: a page token and the request ID of an error's text.
var madeUp = regexp.MustCompile(`"nextToken":"[A-Za-z0-9+/=]+"|RequestID: [0-9a-f-]{36}`)

func (codec) MadeUp() *regexp.Regexp { return madeUp }

// StandIn is a made-up value's stand-in, counting n from 1. A token becomes one
// that says it is a stand-in, base64-encoded as the endpoint's own tokens are;
// a request ID becomes the all-zero one, since a UUID a test reads means
// nothing.
func (codec) StandIn(name string, n int) string {
	if strings.HasPrefix(name, `"nextToken"`) {
		token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("omni-token-%d", n)))
		return fmt.Sprintf(`"nextToken":"%s"`, token)
	}
	return fmt.Sprintf("RequestID: 00000000-0000-0000-0000-%012x", n)
}

type server struct {
	sync.Mutex
	container testcontainers.Container
	endpoint  string
	err       error
	seeded    bool
}

var (
	shared = &server{}
	paged  = &server{}
)

// Endpoint starts the pinned endpoint on first use and returns its URL, with
// the fixture in place. It reaches the endpoint, so only a recording calls it.
func Endpoint(t testing.TB) (string, error) {
	t.Helper()
	return shared.start(t, "the fixture endpoint", seed)
}

// PagedEndpoint starts the endpoint that holds PagedTables tables on first use
// and returns its URL. It reaches the endpoint, so only a recording calls it.
func PagedEndpoint(t testing.TB) (string, error) {
	t.Helper()
	return paged.start(t, "the paged endpoint", seedPaged)
}

// start runs the pinned image once for s, waits for its endpoint and writes the
// fixture.
func (s *server) start(t testing.TB, what string, write func(*awssdk.Client) error) (string, error) {
	t.Helper()
	pin, err := enginetest.LoadPin(Pins, Default)
	if err != nil {
		return "", err
	}
	s.Lock()
	defer s.Unlock()
	if s.endpoint == "" && s.err == nil {
		s.endpoint, s.container, s.err = startContainer(pin, what)
	}
	if s.err != nil {
		return "", s.err
	}
	if !s.seeded {
		if err := write(clientFor(s.endpoint)); err != nil {
			return "", err
		}
		s.seeded = true
	}
	return s.endpoint, nil
}

// clientFor returns a client to endpoint with the placeholder credential, which
// is what a fixture is written over.
func clientFor(endpoint string) *awssdk.Client {
	return awssdk.NewFromConfig(aws.Config{
		Region:      Region,
		Credentials: credentials.NewStaticCredentialsProvider(AccessKeyID, SecretAccessKey, ""),
	}, func(o *awssdk.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
}

// startContainer runs the pinned image and waits for the endpoint to answer.
func startContainer(pin enginetest.Pin, what string) (endpoint string, container testcontainers.Container, err error) {
	// testcontainers panics when Docker is missing; fail the test instead.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("docker provider panic: %v", r)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// -sharedDb keeps one set of tables for the endpoint rather than one per
	// access key, which is what the service does: a table belongs to the
	// caller's account, not to the key pair that signed the request.
	containerPort := nat.Port(Port + "/tcp")
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        pin.Image,
			ExposedPorts: []string{string(containerPort)},
			Cmd:          []string{"-jar", "DynamoDBLocal.jar", "-inMemory", "-sharedDb"},
			WaitingFor:   wait.ForListeningPort(containerPort).WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	})
	if err != nil {
		if c != nil {
			_ = c.Terminate(ctx, testcontainers.StopTimeout(0))
		}
		return "", nil, fmt.Errorf("start %s: %w", what, err)
	}
	if err := checkDigest(ctx, c, pin.Digest); err != nil {
		_ = c.Terminate(ctx, testcontainers.StopTimeout(0))
		return "", nil, err
	}
	host, err := c.Host(ctx)
	if err != nil {
		return "", c, fmt.Errorf("%s's host: %w", what, err)
	}
	port, err := c.MappedPort(ctx, containerPort)
	if err != nil {
		return "", c, fmt.Errorf("%s's port: %w", what, err)
	}
	return fmt.Sprintf("http://%s:%s", host, port.Port()), c, nil
}

// checkDigest holds the container to the pin. An image reference is a tag,
// which a registry may move; the digest is what the answers were recorded from.
func checkDigest(ctx context.Context, c testcontainers.Container, digest string) error {
	info, err := c.Inspect(ctx)
	if err != nil {
		return fmt.Errorf("inspect the DynamoDB Local container: %w", err)
	}
	cli, err := testcontainers.NewDockerClientWithOpts(ctx)
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()
	image, err := cli.ImageInspect(ctx, info.Image)
	if err != nil {
		return fmt.Errorf("inspect the DynamoDB Local image: %w", err)
	}
	for _, d := range image.RepoDigests {
		if strings.HasSuffix(d, "@"+digest) {
			return nil
		}
	}
	return fmt.Errorf("DynamoDB Local runs an image with digests %v, but pins.tsv pins %s", image.RepoDigests, digest)
}

// Terminate removes the endpoints this run started, if any. The package's
// TestMain calls it.
func Terminate() {
	for _, s := range []*server{shared, paged} {
		s.Lock()
		container := s.container
		s.container = nil
		s.Unlock()
		if container == nil {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		_ = container.Terminate(ctx, testcontainers.StopTimeout(0))
		cancel()
	}
}

// The tables the fixture holds.
const (
	// People is the table whose items carry every DynamoDB type.
	People = "people"
	// Orders is a table with a partition key and no sort key.
	Orders = "orders"
	// Paged is the table whose items one partition key holds, ordered by their
	// sort key.
	Paged = "paged"
	// Writes is the table Execute's tests write to, so that the fixture above
	// stays as it was recorded.
	Writes = "writes"
	// PagedTables is the table count of the second container: one more than the
	// 100 names one ListTables answer holds.
	PagedTables = 101
	// GlobalIndex is Paged's global secondary index, over the label attribute.
	GlobalIndex = "by-label"
	// LocalIndex is Paged's local secondary index, over the label attribute.
	LocalIndex = "by-owner-label"
)

// The fixture's items, as the tests read them.
const (
	// PersonWithEveryType carries one attribute of every DynamoDB type,
	// including a NULL.
	PersonWithEveryType = "u1"
	// PersonWithFewAttributes leaves most attributes out, so a read of the
	// table answers a row with nils in it.
	PersonWithFewAttributes = "u2"
	// PersonWithALateAttribute names an attribute the items before it do not.
	PersonWithALateAttribute = "u3"
	// PeopleCount is how many items People holds.
	PeopleCount = 3
	// OrdersCount is how many items Orders holds.
	OrdersCount = 2
	// PagedCount is how many items Paged holds under its one partition key.
	PagedCount = 5
	// PagedOwner is the partition key value Paged's items share.
	PagedOwner = "o1"
)

// seed writes the fixture: the four tables above and their items.
func seed(c *awssdk.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := createTable(ctx, c, People, []types.KeySchemaElement{
		{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		{AttributeName: aws.String("email"), KeyType: types.KeyTypeRange},
	}, []types.AttributeDefinition{
		{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
		{AttributeName: aws.String("email"), AttributeType: types.ScalarAttributeTypeS},
	}); err != nil {
		return err
	}
	if err := createTable(ctx, c, Orders, []types.KeySchemaElement{
		{AttributeName: aws.String("pk"), KeyType: types.KeyTypeHash},
	}, []types.AttributeDefinition{
		{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
	}); err != nil {
		return err
	}
	if _, err := c.CreateTable(ctx, &awssdk.CreateTableInput{
		TableName: aws.String(Paged),
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("owner"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("seq"), KeyType: types.KeyTypeRange},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("owner"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("seq"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("label"), AttributeType: types.ScalarAttributeTypeS},
		},
		BillingMode: types.BillingModePayPerRequest,
		// One index of each kind, over an attribute that is neither key, so that
		// a sync reports both and the column set grows past the table's keys.
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(GlobalIndex),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("label"), KeyType: types.KeyTypeHash},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			},
		},
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName: aws.String(LocalIndex),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("owner"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("label"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			},
		},
	}); err != nil {
		return fmt.Errorf("create table %s: %w", Paged, err)
	}
	if err := createTable(ctx, c, Writes, []types.KeySchemaElement{
		{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
	}, []types.AttributeDefinition{
		{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
	}); err != nil {
		return err
	}

	for _, item := range []map[string]types.AttributeValue{
		{
			"id":      &types.AttributeValueMemberS{Value: PersonWithEveryType},
			"email":   &types.AttributeValueMemberS{Value: "u1@example.com"},
			"name":    &types.AttributeValueMemberS{Value: "Ann"},
			"age":     &types.AttributeValueMemberN{Value: "33"},
			"active":  &types.AttributeValueMemberBOOL{Value: true},
			"tags":    &types.AttributeValueMemberSS{Value: []string{"a", "b"}},
			"scores":  &types.AttributeValueMemberNS{Value: []string{"1", "2"}},
			"payload": &types.AttributeValueMemberB{Value: []byte{0xff, 0x00, 0xfe}},
			"flags":   &types.AttributeValueMemberBS{Value: [][]byte{{0xff}, {0xfe}}},
			"meta":    &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{"k": &types.AttributeValueMemberS{Value: "v"}, "raw": &types.AttributeValueMemberB{Value: []byte{0xff}}}},
			"list":    &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberN{Value: "7"}, &types.AttributeValueMemberS{Value: "s"}, &types.AttributeValueMemberB{Value: []byte{0xfe}}}},
			"missing": &types.AttributeValueMemberNULL{Value: true},
		},
		{
			"id":    &types.AttributeValueMemberS{Value: PersonWithFewAttributes},
			"email": &types.AttributeValueMemberS{Value: "u2@example.com"},
			"name":  &types.AttributeValueMemberS{Value: "Bob"},
			"age":   &types.AttributeValueMemberN{Value: "44"},
		},
		{
			"id":    &types.AttributeValueMemberS{Value: PersonWithALateAttribute},
			"email": &types.AttributeValueMemberS{Value: "u3@example.com"},
			"late":  &types.AttributeValueMemberS{Value: "last"},
		},
	} {
		if err := putItem(ctx, c, People, item); err != nil {
			return err
		}
	}

	for _, item := range []map[string]types.AttributeValue{
		{
			"pk":    &types.AttributeValueMemberS{Value: "o1"},
			"total": &types.AttributeValueMemberN{Value: "10"},
		},
		{
			"pk":    &types.AttributeValueMemberS{Value: "o2"},
			"total": &types.AttributeValueMemberN{Value: "20"},
		},
	} {
		if err := putItem(ctx, c, Orders, item); err != nil {
			return err
		}
	}

	for i := range PagedCount {
		if err := putItem(ctx, c, Paged, map[string]types.AttributeValue{
			"owner": &types.AttributeValueMemberS{Value: PagedOwner},
			"seq":   &types.AttributeValueMemberS{Value: fmt.Sprintf("s%d", i)},
			"label": &types.AttributeValueMemberS{Value: fmt.Sprintf("row-%d", i)},
		}); err != nil {
			return err
		}
	}
	return nil
}

// seedPaged writes the container whose listing takes two ListTables answers.
func seedPaged(c *awssdk.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	for i := range PagedTables {
		if err := createTable(ctx, c, PagedTableName(i), []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		}, []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
		}); err != nil {
			return err
		}
	}
	return nil
}

// PagedTableName is the name of one of the paged container's tables.
func PagedTableName(i int) string {
	return fmt.Sprintf("t%03d", i)
}

// PagedTableNames is every name of the paged container's tables, in the order
// ListTables answers them.
func PagedTableNames() []string {
	names := make([]string, 0, PagedTables)
	for i := range PagedTables {
		names = append(names, PagedTableName(i))
	}
	return names
}

func createTable(ctx context.Context, c *awssdk.Client, name string, keySchema []types.KeySchemaElement, attributeDefinitions []types.AttributeDefinition) error {
	if _, err := c.CreateTable(ctx, &awssdk.CreateTableInput{
		TableName:            aws.String(name),
		KeySchema:            keySchema,
		AttributeDefinitions: attributeDefinitions,
		BillingMode:          types.BillingModePayPerRequest,
	}); err != nil {
		return fmt.Errorf("create table %s: %w", name, err)
	}
	return nil
}

func putItem(ctx context.Context, c *awssdk.Client, table string, item map[string]types.AttributeValue) error {
	if _, err := c.PutItem(ctx, &awssdk.PutItemInput{
		TableName: aws.String(table),
		Item:      item,
	}); err != nil {
		return fmt.Errorf("write an item of %s: %w", table, err)
	}
	return nil
}
