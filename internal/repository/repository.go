package repository

import (
	"context"
	"log"

	Openapi "github.com/karthik2304/xm-go-service/api/v1/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoClient struct {
	db *mongo.Database
}

type Repository interface {
	// company CRUDS
	Create(ctx context.Context, data Openapi.CreateCompanyJSONBody) error
	Update(ctx context.Context, filter map[string]interface{}, details Openapi.UpdateCompanyDetailsJSONRequestBody) error
	GetOne(ctx context.Context, filter map[string]interface{}) (*Openapi.CreateCompanyJSONRequestBody, error)
	Delete(ctx context.Context, filter map[string]interface{}) error
	GetAll(ctx context.Context) (*[]Openapi.SuccessResponse, error)

	// Auth
	CreateAuth(ctx context.Context, data Openapi.SignUpJSONRequestBody) error
	GetAuth(ctx context.Context, email string) (*Openapi.SignUpJSONRequestBody, error)

	// Events Upload
	Send(ctx context.Context, data Openapi.Event) error
	GetAllEvents(ctx context.Context) (*[]Openapi.Event, error)
}

// GetAllEvents implements Repository.
func (m *MongoClient) GetAllEvents(ctx context.Context) (*[]Openapi.Event, error) {
	var response []Openapi.Event
	collection := m.db.Collection("events")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Send implements Repository.
func (m *MongoClient) Send(ctx context.Context, data Openapi.Event) error {
	collection := m.db.Collection("events")
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// GetAuth implements Repository.
func (m *MongoClient) GetAuth(ctx context.Context, email string) (*Openapi.SignUpJSONRequestBody, error) {
	var response Openapi.SignUpJSONRequestBody
	collection := m.db.Collection("users")
	filter := bson.D{{Key: "username", Value: email}}
	err := collection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// CreateAuth implements Repository.
func (m *MongoClient) CreateAuth(ctx context.Context, data Openapi.SignUpJSONRequestBody) error {
	collection := m.db.Collection("users")
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// Create implements Repository.
func (m *MongoClient) Create(ctx context.Context, data Openapi.CreateCompanyJSONBody) error {
	collection := m.db.Collection("companies")
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements Repository.
func (m *MongoClient) Delete(ctx context.Context, filter map[string]interface{}) error {
	collection := m.db.Collection("companies")
	filters := bson.M(filter)
	_, err := collection.DeleteOne(ctx, filters)
	if err != nil {
		return err
	}
	return nil
}

// GetOne implements Repository.
func (m *MongoClient) GetOne(ctx context.Context, filter map[string]interface{}) (*Openapi.CreateCompanyJSONRequestBody, error) {
	var response Openapi.CreateCompanyJSONRequestBody
	collection := m.db.Collection("companies")
	filters := bson.M(filter)
	err := collection.FindOne(ctx, filters).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetAll implements Repository.
func (m *MongoClient) GetAll(ctx context.Context) (*[]Openapi.SuccessResponse, error) {
	var response []Openapi.SuccessResponse
	collection := m.db.Collection("companies")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Update implements Repository.
func (m *MongoClient) Update(ctx context.Context, filter map[string]interface{}, details Openapi.UpdateCompanyDetailsJSONRequestBody) error {
	var updateData bson.M
	data, _ := bson.Marshal(details)
	_ = bson.Unmarshal(data, &updateData)
	collection := m.db.Collection("companies")
	filters := bson.M(filter)
	update := bson.M{"$set": updateData}
	_, err := collection.UpdateOne(ctx, filters, update)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func NewDB(db *mongo.Database) Repository {
	return &MongoClient{
		db: db,
	}
}
