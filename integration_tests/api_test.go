package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ClearMongo() {
	mongoUri := os.Getenv("AUDIT_API_MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	coll := client.Database("dbevent").Collection("events")
	_, err = coll.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
}

//nolint:gomnd
func TestAuth(t *testing.T) {
	defer ClearMongo()
	client := NewAuditApiClient()

	t.Run("ok", func(t *testing.T) {
		code, _ := client.AddEvent(Event{Type: "t1", Consumer: "c1"})
		assert.Equal(t, code, 200)
		code, _ = client.ListEvents(QueryParams{})
		assert.Equal(t, code, 200)
	})
	t.Run("bad key", func(t *testing.T) {
		var old_key string
		old_key, client.apiKey = client.apiKey, "wrongsecret"
		defer func() { client.apiKey = old_key }()
		code, _ := client.AddEvent(Event{Type: "t1", Consumer: "c1"})
		assert.Equal(t, code, 403)
		code, _ = client.ListEvents(QueryParams{})
		assert.Equal(t, code, 403)
	})
}

func TestPostEvent(t *testing.T) {
	defer ClearMongo()
	client := NewAuditApiClient()

	t.Run("ok", func(t *testing.T) {
		code, _ := client.AddEvent(Event{Type: "t1", Consumer: "c1"})
		assert.Equal(t, code, 200)
	})
	t.Run("ok,payload", func(t *testing.T) {
		payload := map[string]interface{}{"key": "val", "k2": 1234}
		code, _ := client.AddEvent(Event{Type: "t1", Consumer: "c1", Payload: payload})
		assert.Equal(t, code, 200)
	})
	t.Run("bad format", func(t *testing.T) {
		code, resp := client.AddEvent(Event{Consumer: "c1"})
		assert.Equal(t, code, 400)
		assert.NotEmpty(t, resp.Error)
	})
}

func TestListEvents(t *testing.T) {
	defer ClearMongo()
	client := NewAuditApiClient()

	{
		code, _ := client.AddEvent(Event{Type: "t10", Consumer: "c10"})
		assert.Equal(t, code, 200)
		code, _ = client.AddEvent(Event{Type: "t11", Consumer: "c11"})
		assert.Equal(t, code, 200)
		payload := map[string]interface{}{"key": "val", "k2": 1234}
		code, _ = client.AddEvent(Event{Type: "t12", Consumer: "c10", Payload: payload})
		assert.Equal(t, code, 200)
	}

	consumer := "c10"
	type_ := "t12"
	created_from := time.Now().Add(-10 * time.Hour)
	created_to := time.Now().Add(10 * time.Hour)

	// by consumer
	code, resp := client.ListEvents(QueryParams{Consumer: &consumer})
	assert.Equal(t, code, 200)
	assert.Len(t, resp.Events, 2)

	// by type
	code, resp = client.ListEvents(QueryParams{Type: &type_})
	assert.Equal(t, code, 200)
	assert.Len(t, resp.Events, 1)

	// by created_from
	code, resp = client.ListEvents(QueryParams{CreatedFrom: &created_from})
	assert.Equal(t, code, 200)
	assert.Len(t, resp.Events, 3)

	code, resp = client.ListEvents(QueryParams{CreatedFrom: &created_to})
	assert.Equal(t, code, 200)
	assert.Len(t, resp.Events, 0)

	// by created_to
	code, resp = client.ListEvents(QueryParams{CreatedTo: &created_from})
	assert.Equal(t, code, 200)
	assert.Len(t, resp.Events, 0)

	code, resp = client.ListEvents(QueryParams{CreatedTo: &created_to})
	assert.Equal(t, code, 200)
	assert.Len(t, resp.Events, 3)
}
