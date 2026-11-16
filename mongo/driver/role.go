package driver

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/bytebase/omni/driver"
)

const (
	// bytebaseDefaultDatabase is the database the driver reads its build and
	// user information from. MongoDB's commands run against a database, and
	// this one is Bytebase's own name for an instance-level read.
	bytebaseDefaultDatabase = "bytebase"
)

// UsersInfo is the subset of the mongodb command result of "usersInfo".
type UsersInfo struct {
	Users []User `bson:"users"`
}

// User is the subset of the `users` field in the `User`.
type User struct {
	ID       string `json:"_id" bson:"_id"`
	UserName string `json:"user" bson:"user"`
	DB       string `json:"db" bson:"db"`
	Roles    []Role `json:"roles" bson:"roles"`
}

// Role is the subset of the `roles` field in the `User`.
type Role struct {
	RoleName string `json:"role" bson:"role"`
	DB       string `json:"db" bson:"db"`
}

// getInstanceRoles returns the list of users.
func (d *Driver) getInstanceRoles(ctx context.Context) ([]*driver.Role, error) {
	commandResult, err := d.client.RunCommand(ctx, bytebaseDefaultDatabase, bson.D{{
		Key: "usersInfo",
		Value: bson.D{{
			Key:   "forAllDBs",
			Value: true,
		}},
	}})
	if err != nil {
		if isAtlasUnauthorizedError(err) {
			slog.Info("Skip getting user list because the user is not authorized to run the command 'usersInfo' in atlas cluster M0/M2/M5")
			return nil, nil
		}
		return nil, fmt.Errorf("cannot run usersInfo command: %w", err)
	}
	users, err := decodeUsers(commandResult)
	if err != nil {
		return nil, err
	}
	var instanceRoles []*driver.Role
	for _, user := range users {
		instanceRoles = append(instanceRoles, &driver.Role{
			Name: user.ID,
		})
	}
	return instanceRoles, nil
}

// decodeUsers reads the usersInfo answer. Bytebase decoded the command result
// straight into UsersInfo; the client hands back a document, so the same
// decode runs here.
func decodeUsers(commandResult bson.M) ([]User, error) {
	raw, err := bson.Marshal(commandResult)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal usersInfo command result: %w", err)
	}
	var usersInfo UsersInfo
	if err := bson.Unmarshal(raw, &usersInfo); err != nil {
		return nil, fmt.Errorf("cannot decode usersInfo command result: %w", err)
	}
	return usersInfo.Users, nil
}
