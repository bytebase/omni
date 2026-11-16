package driver

import (
	"context"
	"sort"

	"github.com/bytebase/gomongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// client is what the driver asks of a MongoDB client: the commands sync reads
// and the mongosh execution query and execute run. It is the seam a test
// replaces with a tape, so the driver's own code is what a recording and a
// replay exercise; the live implementation below is the one line per call it
// takes to reach the mongo-go-driver client.
//
// A cursor's documents come back ordered by name, not in the order the server
// sent them: MongoDB lists collections in a hash order that changes between
// runs, and the metadata and its dump must not.
type client interface {
	Ping(ctx context.Context) error
	Disconnect(ctx context.Context) error
	ListDatabaseNames(ctx context.Context) ([]string, error)
	ListCollections(ctx context.Context, database string) ([]bson.M, error)
	EstimatedDocumentCount(ctx context.Context, database, collection string) (int64, error)
	CollStats(ctx context.Context, database, collection string) (bson.M, error)
	ListIndexes(ctx context.Context, database, collection string) ([]bson.M, error)
	RunCommand(ctx context.Context, database string, command bson.D) (bson.M, error)
	Execute(ctx context.Context, database, statement string, maxRows int64) (*gomongo.Result, error)
}

// liveClient is the client the driver opens against a server.
type liveClient struct {
	client *mongo.Client
}

func (c liveClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx, nil)
}

func (c liveClient) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c liveClient) ListDatabaseNames(ctx context.Context) ([]string, error) {
	return c.client.ListDatabaseNames(ctx, bson.M{})
}

func (c liveClient) ListCollections(ctx context.Context, database string) ([]bson.M, error) {
	cursor, err := c.client.Database(database).ListCollections(ctx, bson.M{})
	return readCursor(ctx, cursor, err)
}

func (c liveClient) EstimatedDocumentCount(ctx context.Context, database, collection string) (int64, error) {
	return c.client.Database(database).Collection(collection).EstimatedDocumentCount(ctx)
}

func (c liveClient) CollStats(ctx context.Context, database, collection string) (bson.M, error) {
	doc, err := c.RunCommand(ctx, database, bson.D{{Key: "collStats", Value: collection}})
	if err != nil {
		return nil, err
	}
	return collStatsFields(doc), nil
}

// collStatsFields keeps what SyncDatabase reads of a collStats answer. The
// answer also carries a wiredTiger statistics tree whose counters move between
// runs, and a collection's namespace and UUID: none of that is the driver's, so
// a tape of these answers holds bytes that do not change.
func collStatsFields(doc bson.M) bson.M {
	fields := bson.M{}
	for _, key := range []string{"size", "totalIndexSize"} {
		if value, ok := doc[key]; ok {
			fields[key] = value
		}
	}
	return fields
}

func (c liveClient) ListIndexes(ctx context.Context, database, collection string) ([]bson.M, error) {
	cursor, err := c.client.Database(database).Collection(collection).Indexes().List(ctx)
	return readCursor(ctx, cursor, err)
}

func (c liveClient) RunCommand(ctx context.Context, database string, command bson.D) (bson.M, error) {
	var result bson.M
	if err := c.client.Database(database).RunCommand(ctx, command).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c liveClient) Execute(ctx context.Context, database, statement string, maxRows int64) (*gomongo.Result, error) {
	var opts []gomongo.ExecuteOption
	if maxRows > 0 {
		opts = append(opts, gomongo.WithMaxRows(maxRows))
	}
	return gomongo.NewClient(c.client).Execute(ctx, database, statement, opts...)
}

// readCursor decodes every document a cursor holds, ordered by name.
func readCursor(ctx context.Context, cursor *mongo.Cursor, err error) ([]bson.M, error) {
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(ctx) }()
	var docs []bson.M
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	sortDocsByName(docs)
	return docs, nil
}

// sortDocsByName orders documents by their "name", the field every
// listCollections and listIndexes document carries.
func sortDocsByName(docs []bson.M) {
	sort.SliceStable(docs, func(i, j int) bool {
		left, _ := docs[i]["name"].(string)
		right, _ := docs[j]["name"].(string)
		return left < right
	})
}
