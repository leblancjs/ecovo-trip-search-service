package search

import (
	"context"
	"fmt"

	"azure.com/ecovo/trip-search-service/pkg/entity"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// A MongoRepository is a repository that performs CRUD operations on searches
// in a MongoDB collection.
type MongoRepository struct {
	collection *mongo.Collection
}

type document struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

func newDocumentFromEntity(t *entity.Search) (*document, error) {
	if t == nil {
		return nil, fmt.Errorf("trop.MongoRepository: entity is nil")
	}

	var id primitive.ObjectID
	if t.ID.IsZero() {
		id = primitive.NilObjectID
	} else {
		objectID, err := primitive.ObjectIDFromHex(t.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("search.MongoRepository: failed to create object")
		}

		id = objectID
	}

	return &document{
		id,
	}, nil
}

func (d document) Entity() *entity.Search {
	return &entity.Search{
		entity.NewIDFromHex(d.ID.Hex()),
	}
}

// NewMongoRepository creates a search repository for a MongoDB collection.
func NewMongoRepository(collection *mongo.Collection) (Repository, error) {
	if collection == nil {
		return nil, fmt.Errorf("search.MongoRepository: collection is nil")
	}

	return &MongoRepository{collection}, nil
}

// FindByID retrieves the search with the given ID, if it exists.
func (r *MongoRepository) FindByID(ID entity.ID) (*entity.Search, error) {
	objectID, err := primitive.ObjectIDFromHex(string(ID))
	if err != nil {
		return nil, fmt.Errorf("search.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	var d document
	err = r.collection.FindOne(context.TODO(), filter).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("search.MongoRepository: no search found with ID \"%s\" (%s)", ID, err)
	}
	return d.Entity(), nil
}

// Create stores the new search in the database and returns the unique
// identifier that was generated for it.
func (r *MongoRepository) Create(t *entity.Search) (entity.ID, error) {
	if t == nil {
		return entity.NilID, fmt.Errorf("search.MongoRepository: failed to create search (search is nil)")
	}

	d, err := newDocumentFromEntity(t)
	if err != nil {
		return entity.NilID, fmt.Errorf("search.MongoRepository: failed to create search document from entity (%s)", err)
	}

	res, err := r.collection.InsertOne(context.TODO(), d)
	if err != nil {
		return entity.NilID, fmt.Errorf("search.MongoRepository: failed to create search (%s)", err)
	}

	ID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return entity.NilID, fmt.Errorf("search.MongoRepository: failed to get ID of created search (%s)", err)
	}

	return entity.ID(ID.Hex()), nil
}

// Delete removes the search with the given ID from the database.
func (r *MongoRepository) Delete(ID entity.ID) error {
	objectID, err := primitive.ObjectIDFromHex(string(ID))
	if err != nil {
		return fmt.Errorf("search.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("search.MongoRepository: failed to delete search with ID \"%s\" (%s)", ID, err)
	}

	return nil
}
