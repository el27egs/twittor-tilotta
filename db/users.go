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

	filter := bson.M{"email": email}

	var userFound models.User

	err := collection.FindOne(timeoutCtx, filter).Decode(&userFound)
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

func UpdateUser(u models.User, ID string) (bool, error) {
	timeoutCtx, cancelHandler := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelHandler()

	db := MongoConnection.Database("twittor")
	collection := db.Collection("usuarios")

	data := make(map[string]interface{})

	if len(u.Name) > 0 {
		data["name"] = u.Name
	}
	if len(u.LastName) > 0 {
		data["lastName"] = u.LastName
	}
	data["birthday"] = u.Birthday
	if len(u.Avatar) > 0 {
		data["avatar"] = u.Avatar
	}
	if len(u.Banner) > 0 {
		data["banner"] = u.Banner
	}
	if len(u.Biography) > 0 {
		data["biography"] = u.Biography
	}
	if len(u.Location) > 0 {
		data["location"] = u.Location
	}

	if len(u.Web) > 0 {
		data["web"] = u.Web
	}

	update := bson.M{"$set": data}

	objID, _ := primitive.ObjectIDFromHex(ID)

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	updateResult, err := collection.UpdateOne(timeoutCtx, filter, update)

	if err != nil {
		return false, err
	}
	if updateResult.MatchedCount == 0 || updateResult.ModifiedCount == 0 {
		return false, err
	}
	return true, nil
}
