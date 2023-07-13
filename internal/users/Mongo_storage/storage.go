package Mongo_storage

import (
	"context"
	"github.com/sillamilla/user_microservice/internal/users/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage interface {
	SignUp(ctx context.Context, user model.User) error
	SignIn(ctx context.Context, input model.Input) (model.User, error)

	EditProfile(ctx context.Context, id string, input model.UpdateUser) error
	EditPassword(ctx context.Context, id string, password string) error

	GetByID(ctx context.Context, id string) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)

	UpsertSession(ctx context.Context, id string, session string) error
	GetBySession(ctx context.Context, session string) (model.User, error)
	SearchByUsername(ctx context.Context, username string) (model.UserInfo, error)
}

type mongoDB struct {
	mo *mongo.Client
}

func New(mo *mongo.Client) Storage {
	return &mongoDB{
		mo: mo,
	}
}

func (db *mongoDB) SignUp(ctx context.Context, user model.User) error {
	_, err := db.mo.Database("users_microservice").Collection("users").InsertOne(ctx, bson.M{"id": user.ID, "username": user.Username, "password": user.Password, "session": user.Session, "createdAt": user.CreatedAt})
	if err != nil {
		return err
	}

	return nil
}

func (db *mongoDB) SignIn(ctx context.Context, input model.Input) (model.User, error) {
	var user model.User

	filter := bson.M{"username": input.Username, "password": input.Password}
	err := db.mo.Database("users_microservice").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (db *mongoDB) EditProfile(ctx context.Context, id string, input model.UpdateUser) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"username": input.Username, "email": input.Email, "bio": input.Bio, "icon": input.Icon}}
	options := options.Update().SetUpsert(true)

	_, err := db.mo.Database("users_microservice").Collection("users").UpdateOne(ctx, filter, update, options)
	if err != nil {
		return err
	}

	return nil
}

func (db *mongoDB) EditPassword(ctx context.Context, id string, password string) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"password": password}}
	options := options.Update().SetUpsert(true)

	_, err := db.mo.Database("users_microservice").Collection("users").UpdateOne(ctx, filter, update, options)
	if err != nil {
		return err
	}

	return nil
}

func (db *mongoDB) GetByID(ctx context.Context, id string) (model.User, error) {
	var user model.User

	filter := bson.M{"id": id}
	err := db.mo.Database("users_microservice").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (db *mongoDB) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User

	filter := bson.M{"username": username}
	err := db.mo.Database("users_microservice").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (db *mongoDB) SearchByUsername(ctx context.Context, username string) (model.UserInfo, error) {
	var user model.UserInfo

	filter := bson.M{"username": username}
	err := db.mo.Database("users_microservice").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return model.UserInfo{}, err
	}

	return user, nil
}

func (db *mongoDB) GetBySession(ctx context.Context, session string) (model.User, error) {
	var user model.User

	filter := bson.M{"session": session}
	err := db.mo.Database("users_microservice").Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (db *mongoDB) UpsertSession(ctx context.Context, id string, session string) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"session": session}}
	options := options.Update().SetUpsert(true)
	_, err := db.mo.Database("users_microservice").Collection("users").UpdateOne(ctx, filter, update, options)
	if err != nil {
		return err
	}

	return nil
}
