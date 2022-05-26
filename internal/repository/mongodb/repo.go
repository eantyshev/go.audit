package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.audit/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DB = "dbevent"
const COLL = "events"

type Repo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (r *Repo) InsertEvent(ctx context.Context, event entity.Event) (entity.ID, error) {
	// assign UTC timestamp
	event.CreatedAt = time.Now().UTC()

	doc, err := bson.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("failed to serialize to BSON: %w", err)
	}
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert: %w", err)
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	if _id, ok := result.InsertedID.(primitive.ObjectID); !ok {
		return "", errors.New("failed to get _id")
	} else {
		return entity.ID(_id.String()), nil
	}
}

func MakeQueryBSON(params entity.QueryParams) bson.M {
	query := bson.M{}
	if params.Consumer != nil {
		query["consumer"] = *params.Consumer
	}
	if params.Type != nil {
		query["type"] = *params.Type
	}
	if params.CreatedFrom != nil || params.CreatedTo != nil {
		query_created := bson.A{}
		if params.CreatedFrom != nil {
			query_created = append(query_created,
				bson.D{
					bson.E{Key: "created_at",
						Value: bson.D{bson.E{Key: "$gt", Value: *params.CreatedFrom}}}})
		}
		if params.CreatedTo != nil {
			query_created = append(query_created,
				bson.D{
					bson.E{Key: "created_at",
						Value: bson.D{bson.E{Key: "$lt", Value: *params.CreatedTo}}}})
		}
		query["$and"] = query_created
	}
	return query
}

func (r *Repo) FindEvents(ctx context.Context, params entity.QueryParams) ([]entity.Event, error) {
	lst := make([]entity.Event, 0)
	queryBson := MakeQueryBSON(params)
	cursor, err := r.collection.Find(ctx, queryBson)
	if err != nil {
		return lst, err
	}
	if err = cursor.All(ctx, &lst); err != nil {
		return lst, err
	}
	return lst, nil
}

func MakeMongoRepo(mongoUri string) *Repo {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	coll := client.Database(DB).Collection(COLL)
	return &Repo{client, coll}
}

func (r *Repo) Close() {
	if err := r.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
