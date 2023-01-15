package db

import (
	"context"
	"github.com/el27egs/twittor-tilotta/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func CreateTweet(t models.UserTweet) (string, bool, error) {

	timeoutCtx, cancelHandler := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelHandler()

	db := MongoConnection.Database("twittor")
	collection := db.Collection("tweet")

	tweet := bson.M{
		"userid":  t.UserID,
		"message": t.Message,
		"date":    t.Date,
	}

	result, err := collection.InsertOne(timeoutCtx, tweet)
	if err != nil {
		return "", false, err
	}
	objID, _ := result.InsertedID.(primitive.ObjectID)

	return objID.Hex(), true, nil
}

func GetTweetsWithPager(userId string, page int64) ([]*models.TweetResponse, bool) {

	timeoutCtx, cancelHandler := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelHandler()

	db := MongoConnection.Database("twittor")
	collection := db.Collection("tweet")

	var tweets []*models.TweetResponse

	filter := bson.M{"userid": userId}

	opts := options.Find()
	opts.SetLimit(20)
	//La diferencia entre bson.M y bson.D es que la ultima conserve al orden
	// en el que se declaran los parametros
	opts.SetSort(bson.D{{Key: "date", Value: -1}})
	opts.SetSkip((page - 1) * 20)

	rs, err := collection.Find(timeoutCtx, filter, opts)
	if err != nil {
		log.Fatal(err.Error())
		return []*models.TweetResponse{}, false
	}

	for rs.Next(context.TODO()) {
		var tweet models.TweetResponse
		err := rs.Decode(&tweet)
		if err != nil {
			log.Fatal(err.Error())
			return tweets, false
		}
		tweets = append(tweets, &tweet)
	}
	return tweets, true

}
