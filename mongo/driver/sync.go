package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/bytebase/omni/driver"
	metadatapb "github.com/bytebase/omni/metadata"
)

// systemDatabase are the databases no instance sync reports: MongoDB's own
// three and the one this driver reads its build information from.
var systemDatabase = map[string]bool{
	"admin":    true,
	"config":   true,
	"local":    true,
	"bytebase": true,
}

// SyncInstance reads the server version, its databases without their
// collections, and its users.
func (d *Driver) SyncInstance(ctx context.Context) (*driver.Instance, error) {
	version, err := d.getVersion(ctx)
	if err != nil {
		return nil, err
	}
	roles, err := d.getInstanceRoles(ctx)
	if err != nil {
		return nil, err
	}
	databaseNames, err := d.getNonSystemDatabaseList(ctx)
	if err != nil {
		return nil, err
	}
	databases := make([]*metadatapb.DatabaseSchemaMetadata, 0, len(databaseNames))
	for _, databaseName := range databaseNames {
		databases = append(databases, &metadatapb.DatabaseSchemaMetadata{Name: databaseName})
	}
	return &driver.Instance{Version: version, Databases: databases, Roles: roles}, nil
}

// SyncDatabase reads the full schema of Config.Database: its collections with
// their sizes and indexes, and its views. A database that does not exist is
// driver.ErrNotFound.
func (d *Driver) SyncDatabase(ctx context.Context) (*metadatapb.DatabaseSchemaMetadata, error) {
	schemaMetadata := &metadatapb.SchemaMetadata{Name: ""}

	exist, err := d.isDatabaseExist(ctx, d.databaseName)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("%w: database %s does not exist", driver.ErrNotFound, d.databaseName)
	}

	collections, err := d.client.ListCollections(ctx, d.databaseName)
	if err != nil {
		return nil, fmt.Errorf("failed to list collection names: %w", err)
	}
	var collectionNames []string
	var viewNames []string
	for _, collection := range collections {
		var tp string
		if t, ok := collection["type"]; ok {
			switch s, _ := t.(string); s {
			case "collection", "timeseries":
				// A time-series collection is a user collection whose
				// documents the server stores in buckets of its own.
				tp = "collection"
			case "view":
				tp = "view"
			}
		}
		name, ok := collection["name"]
		if !ok {
			return nil, errors.New("cannot get collection name from collection info")
		}
		collectionName, ok := name.(string)
		if !ok {
			return nil, errors.New("cannot convert collection name to string")
		}
		switch tp {
		case "collection":
			collectionNames = append(collectionNames, collectionName)
		case "view":
			viewNames = append(viewNames, collectionName)
		default:
			// Other types like system collections
		}
	}
	slices.Sort(collectionNames)
	slices.Sort(viewNames)

	for _, collectionName := range collectionNames {
		if isSystemCollection(collectionName) {
			continue
		}

		count, err := d.estimatedDocumentCount(ctx, collectionName)
		if err != nil {
			return nil, fmt.Errorf("failed to get estimated document count: %w", err)
		}
		// Get collection data size and total index size in byte.
		commandResult, err := d.client.CollStats(ctx, d.databaseName, collectionName)
		if err != nil {
			return nil, fmt.Errorf("cannot run collStats command: %w", err)
		}
		dataSize, ok := commandResult["size"]
		if !ok {
			return nil, errors.New("cannot get size from collStats command result")
		}
		// Bytebase logged a value it could not convert and read the size as
		// zero; the driver contract has no logger, so the reading stays.
		dataSize64, _ := convertEmptyInterfaceToInt64(dataSize)

		totalIndexSize, ok := commandResult["totalIndexSize"]
		if !ok {
			return nil, errors.New("cannot get totalIndexSize from collStats command result")
		}
		totalIndexSize64, _ := convertEmptyInterfaceToInt64(totalIndexSize)

		indexes, err := getIndexes(ctx, d.client, d.databaseName, collectionName)
		if err != nil {
			return nil, fmt.Errorf("failed to get index schema of collection %s: %w", collectionName, err)
		}
		schemaMetadata.Tables = append(schemaMetadata.Tables, &metadatapb.TableMetadata{
			Name:      collectionName,
			RowCount:  count,
			DataSize:  dataSize64,
			IndexSize: totalIndexSize64,
			Indexes:   indexes,
		})
	}

	for _, viewName := range viewNames {
		schemaMetadata.Views = append(schemaMetadata.Views, &metadatapb.ViewMetadata{Name: viewName})
	}

	return &metadatapb.DatabaseSchemaMetadata{
		Name:    d.databaseName,
		Schemas: []*metadatapb.SchemaMetadata{schemaMetadata},
	}, nil
}

// getIndexes returns all indexes schema of a collection.
// https://www.mongodb.com/docs/manual/reference/command/listIndexes/#output
func getIndexes(ctx context.Context, c client, database, collectionName string) ([]*metadatapb.IndexMetadata, error) {
	indexInfos, err := c.ListIndexes(ctx, database, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	indexMap := make(map[string]*metadatapb.IndexMetadata)
	for _, indexInfo := range indexInfos {
		name, ok := indexInfo["name"]
		if !ok {
			return nil, errors.New("cannot get index name from index info")
		}
		indexName, ok := name.(string)
		if !ok {
			return nil, errors.New("cannot convert index name to string")
		}
		key, ok := indexInfo["key"]
		if !ok {
			return nil, errors.New("cannot get index key from index info")
		}
		expression, err := json.Marshal(key)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal index key to json: %w", err)
		}
		unique := false
		if u, ok := indexInfo["unique"]; ok {
			unique, ok = u.(bool)
			if !ok {
				return nil, errors.New("cannot convert unique to bool")
			}
		}

		if _, ok := indexMap[indexName]; !ok {
			indexMap[indexName] = &metadatapb.IndexMetadata{
				Name:   indexName,
				Unique: unique,
			}
		}
		indexMap[indexName].Expressions = append(indexMap[indexName].Expressions, string(expression))
	}

	var indexes []*metadatapb.IndexMetadata
	var indexNames []string
	for name := range indexMap {
		indexNames = append(indexNames, name)
	}
	slices.Sort(indexNames)
	for _, name := range indexNames {
		indexes = append(indexes, indexMap[name])
	}
	return indexes, nil
}

func isSystemCollection(collectionName string) bool {
	return strings.HasPrefix(collectionName, "system.")
}

// estimatedDocumentCount reads a collection's estimated count.
func (d *Driver) estimatedDocumentCount(ctx context.Context, collectionName string) (int64, error) {
	return d.client.EstimatedDocumentCount(ctx, d.databaseName, collectionName)
}

// getVersion returns the version of mongod or mongos instance.
func (d *Driver) getVersion(ctx context.Context) (string, error) {
	commandResult, err := d.client.RunCommand(ctx, bytebaseDefaultDatabase, bson.D{{Key: "buildInfo", Value: 1}})
	if err != nil {
		return "", fmt.Errorf("cannot run buildInfo command: %w", annotateUnsupportedServerVersion(err))
	}
	version, ok := commandResult["version"]
	if !ok {
		return "", errors.New("cannot get version from buildInfo command result")
	}
	v, ok := version.(string)
	if !ok {
		return "", errors.New("cannot convert version to string")
	}
	return v, nil
}

// isDatabaseExist returns true if the database exists.
func (d *Driver) isDatabaseExist(ctx context.Context, databaseName string) (bool, error) {
	// We do the filter by hand instead of using the filter option of ListDatabaseNames because we may encounter the following error:
	// Unallowed argument in listDatabases command: filter
	databaseList, err := d.client.ListDatabaseNames(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to list database names: %w", err)
	}
	for _, database := range databaseList {
		if database == databaseName {
			return true, nil
		}
	}
	return false, nil
}

// getNonSystemDatabaseList returns the list of non system databases.
func (d *Driver) getNonSystemDatabaseList(ctx context.Context) ([]string, error) {
	databaseNames, err := d.client.ListDatabaseNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list database names: %w", err)
	}
	var nonSystemDatabaseList []string
	for _, databaseName := range databaseNames {
		if _, ok := systemDatabase[databaseName]; !ok {
			nonSystemDatabaseList = append(nonSystemDatabaseList, databaseName)
		}
	}
	return nonSystemDatabaseList, nil
}

// isAtlasUnauthorizedError returns true if the error is an Atlas unauthorized error.
func isAtlasUnauthorizedError(err error) bool {
	commandError, ok := err.(mongo.CommandError)
	if !ok {
		return strings.Contains(err.Error(), "AtlasError: Unauthorized")
	}
	// Atlas M0/M2/M5 shared cluster does not support usersInfo command.
	if commandError.Name == "AtlasError" && commandError.Code == 8000 && strings.Contains(commandError.Message, "Unauthorized") {
		return true
	}
	// M10/M20/M30 returns the following error.
	if commandError.Name == "Unauthorized" && commandError.Code == 13 {
		return true
	}

	return false
}

func convertEmptyInterfaceToInt64(value any) (int64, error) {
	// NOTE: convert uint64 to int64 may cause overflow. But we don't care about it.
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %v to int64", value)
	}
}
