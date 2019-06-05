package bumble

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// A Database is an abstract dating profile database.
type Database interface {
	AddUser(u *User) error
	GetUser(userID string) (*User, error)

	PhotoExists(id string) (bool, error)
	AddPhoto(id string, data []byte) error
	GetPhoto(id string) ([]byte, error)
}

type mongoDatabase struct {
	config   *Config
	client   *mongo.Client
	db       *mongo.Database
	photos   *mongo.Collection
	profiles *mongo.Collection
}

func OpenDatabase(c *Config) (Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.DatabaseURI))
	if err != nil {
		return nil, err
	}
	db := client.Database("bumble")
	return &mongoDatabase{
		config:   c,
		client:   client,
		db:       db,
		photos:   db.Collection("photos"),
		profiles: db.Collection("profiles"),
	}, nil
}

func (m *mongoDatabase) AddUser(u *User) error {
	// TODO: this.
	return errors.New("not yet implemented")
}

func (m *mongoDatabase) GetUser(userID string) (*User, error) {
	// TODO: this.
	return nil, errors.New("not yet implemented")
}

func (m *mongoDatabase) PhotoExists(id string) (bool, error) {
	_, err := m.photos.FindOne(context.Background(), bson.D{{Key: "id", Value: id}}).DecodeBytes()
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err == nil {
		return true, nil
	}
	return false, errors.Wrap(err, "check photo exists")
}

func (m *mongoDatabase) AddPhoto(id string, data []byte) error {
	// TODO: this.
	return nil
}

func (m *mongoDatabase) GetPhoto(id string) ([]byte, error) {
	// TODO: this.
	return nil, nil
}
