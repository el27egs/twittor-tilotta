package db

import (
	"context"
	"github.com/el27egs/twittor-tilotta/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func SearchUserByEmail(email string) (models.User, bool, string) {
	timeoutCtx, cancelHandler := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelHandler()

	db := MongoConnection.Database("twittor")
	collection := db.Collection("usuarios")

	predicate := bson.M{"email": email}

	var userFound models.User

	err := collection.FindOne(timeoutCtx, predicate).Decode(&userFound)
	ID := userFound.ID.Hex()
	if err != nil {
		return userFound, false, ID
	}

	return userFound, true, ID
}

func SaveUser(user models.User) (string, bool, error) {

	timeoutCtx, cancelHandler := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelHandler()

	db := MongoConnection.Database("twittor")
	collection := db.Collection("usuarios")

	user.Password, _ = EncriptPasword(user.Password)

	result, err := collection.InsertOne(timeoutCtx, user)
	if err != nil {
		return "", false, err
	}
	ObjID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjID.String(), true, nil
}
