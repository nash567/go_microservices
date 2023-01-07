package repo

import (
	"context"
	"log"
	"time"

	"github.com/loger-service/internal/service/log/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	db *mongo.Client
}

func NewRepository(db *mongo.Client) *Repo {
	return &Repo{db}
}
func (r *Repo) Insert(entry *model.LogEntry) error {
	collection := r.db.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), model.LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("err inserting into logs :%s", err)
		return err
	}
	return nil
}

func (r *Repo) All() ([]*model.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	collection := r.db.Database("logs").Collection("logs")
	options := options.Find()
	options.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, options)

	if err != nil {
		log.Println("err finding all docs:%s", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*model.LogEntry
	for cursor.Next(ctx) {
		var item model.LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("err decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}

	}

	return logs, nil

}

func (r *Repo) GetOne(id string) (*model.LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	collection := r.db.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var entry model.LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil

}

func (r *Repo) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	collection := r.db.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(l model.LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()
	collection := r.db.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}
