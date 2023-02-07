package mongo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/first_ddd_go/aggregates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db        *mongo.Database
	customers *mongo.Collection
}

func New(ctx context.Context, connection string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	if err != nil {
		return nil, err
	}
	db := client.Database("ddd")
	customers := db.Collection("customers")
	return &MongoRepository{
		db:        db,
		customers: customers,
	}, nil
}

func (mr *MongoRepository) Get(id uuid.UUID) (aggregates.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := mr.customers.FindOne(ctx, bson.M{"id": id})
	var c mongoCustomer
	err := result.Decode(&c)
	if err != nil {
	}
	return c.ToAggregate(), nil
}

func (mr *MongoRepository) Add(c aggregates.Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := mr.customers.InsertOne(ctx, NewFromCustomer(c))
	if err != nil {

	}
	return nil
}
func (mr *MongoRepository) Update(c aggregates.Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := mr.customers.UpdateOne(ctx, c.GetID(), NewFromCustomer(c))
	if err != nil {

	}
	return nil
}

func NewFromCustomer(c aggregates.Customer) mongoCustomer {
	return mongoCustomer{
		ID:   c.GetID(),
		Name: c.GetName(),
	}
}

type mongoCustomer struct {
	ID   uuid.UUID `bson:"id"`
	Name string    `bson:"name"`
}

func (m mongoCustomer) ToAggregate() aggregates.Customer {
	c := aggregates.Customer{}
	c.SetID(m.ID)
	c.SetName(m.Name)
	return c
}
