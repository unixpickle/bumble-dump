package bumble

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
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
	AddPhoto(id string, photo *Photo, data []byte) error
	GetPhoto(id string) (*Photo, []byte, error)
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
	_, err := m.profiles.InsertOne(context.Background(), u)
	if err != nil {
		return errors.Wrap(err, "add user")
	}
	return nil
}

func (m *mongoDatabase) GetUser(userID string) (*User, error) {
	var user User
	res := m.photos.FindOne(context.Background(), bson.D{{Key: "id", Value: userID}})
	if err := res.Decode(&user); err != nil {
		return nil, errors.Wrap(err, "get user")
	}
	return &user, nil
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

func (m *mongoDatabase) AddPhoto(id string, photo *Photo, data []byte) error {
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		return errors.Wrap(err, "add photo")
	}

	_, err = tmpFile.Write(data)
	tmpFile.Close()
	if err != nil {
		os.Remove(tmpFile.Name())
		return errors.Wrap(err, "add photo")
	}

	newPath := filepath.Join(m.config.PhotosPath, id+".jpg")
	if err := os.Rename(tmpFile.Name(), newPath); err != nil {
		os.Remove(tmpFile.Name())
		return errors.Wrap(err, "add photo")
	}

	_, err = m.photos.InsertOne(context.Background(), photo)
	if err != nil {
		os.Remove(newPath)
		return errors.Wrap(err, "add photo")
	}

	return nil
}

func (m *mongoDatabase) GetPhoto(id string) (*Photo, []byte, error) {
	var photo Photo
	res := m.photos.FindOne(context.Background(), bson.D{{Key: "id", Value: id}})
	if err := res.Decode(&photo); err != nil {
		return nil, nil, errors.Wrap(err, "get photo")
	}
	data, err := ioutil.ReadFile(filepath.Join(m.config.PhotosPath, id+".jpg"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "get photo")
	}
	return &photo, data, nil
}
