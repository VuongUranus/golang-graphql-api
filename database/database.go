package database

import (
	"context"
	"fmt"
	"golang-graphql-api/graph/model"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
	dbName		string
}

func Connect(dbUrl string) *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Error(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error(err)
	}

	return &DB{
		client: client,
		dbName: "graphql-mongodb-api-db",
	}
}

func (db *DB) InsertMovieById(movie model.NewMovie) *model.Movie {
	movieColl := db.client.Database(db.dbName).Collection("movie")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := movieColl.InsertOne(ctx, bson.D{{Key: "name", Value: movie.Name}})
	if err != nil {
		log.Error(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnMovie := model.Movie{ID: insertedID, Name: movie.Name}

	return &returnMovie
}

func (db *DB) FindMovieById(id string) *model.Movie {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
	}
	movieColl := db.client.Database(db.dbName).Collection("movie")
	ctx, cancel := context.WithTimeout(context.Background(), 30* time.Second)
	defer cancel()
	res := movieColl.FindOne(ctx, bson.M{"_id": ObjectID})

	movie := model.Movie{ID: id}

	res.Decode(&movie)

	return &movie
}

func (db *DB) All() []*model.Movie {
	movieColl := db.client.Database(db.dbName).Collection("movie")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := movieColl.Find(ctx, bson.D{})
	if err != nil {
		 log.Error(err)
	}

	var movies []*model.Movie
	for cur.Next(ctx){
		sus, err := cur.Current.Elements()
		fmt.Println(sus)
		if err != nil {
			log.Error(err)
		}

		movie := model.Movie{ID: (sus[0].String()), Name: (sus[1]).String()}

		movies = append(movies, &movie)
	}
	return movies
}
