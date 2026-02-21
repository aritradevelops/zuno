package mongodb

import (
	"context"
	"fmt"
	"goserve/internal/db"
	"net/url"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ db.Database = &MongoDB{}

type MongoDB struct {
	connectionUrl string
	client        *mongo.Client
	dbName        string
	db            *mongo.Database
}

func New(connectionUrl string) (*MongoDB, error) {
	raw, err := url.Parse(connectionUrl)
	if err != nil {
		return nil, (fmt.Errorf("mongodb connection url is invalid: %w", err))
	} else if raw.Path == "/" {
		fmt.Printf("raw connection: %+v", raw)
		return nil, (fmt.Errorf("mongodb connection url should contain a default database name: %w", err))
	}
	// trim the leading '/'
	dbName := raw.Path[1:]
	return &MongoDB{connectionUrl: connectionUrl, dbName: dbName}, nil
}

// Connect implements db.Database.
func (m *MongoDB) Connect(ctx context.Context) error {
	c, err := mongo.Connect(options.Client().ApplyURI(m.connectionUrl))
	if err != nil {
		return err
	}
	m.client = c
	m.db = c.Database(m.dbName)
	return nil
}

// Disconnect implements db.Database.
func (m *MongoDB) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Ping implements db.Database.
func (m *MongoDB) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, nil)
}

func (m *MongoDB) Client() *mongo.Client {
	return m.client
}

func (m *MongoDB) Database() *mongo.Database {
	return m.db
}
