package db

import (
	"context"
	"fmt"
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

func SearchUserByID(ID string) (models.User, error) {
	timeoutCtx, cancelHandler := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelHandler()

	db := MongoConnection.Database("twittor")
	collection := db.Collection("usuarios")

	var useProfile = models.User{}
	objID, _ := primitive.ObjectIDFromHex(ID)
	predicate := bson.M{
		"_id": objID,
	}
	err := collection.FindOne(timeoutCtx, predicate).Decode(&useProfile)
	useProfile.Password = ""
	if err != nil {
		fmt.Println("Usuario no encontrado " + err.Error())
		return models.User{}, err
	}
	return useProfile, nil
}
