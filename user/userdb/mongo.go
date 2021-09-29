package userdb

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Age      int32              `bson:"age,omitempty"`
	Greeting string             `bson:"greeting,omitempty"`
	Salary   int32              `bson:"salary,omitempty"`
	Power    string             `bson:"power,omitempty"`
}

var _ = loadLocalEnv()

var (
	db   = GetEnv("MONGO_INITDB_DATABASE")
	user = GetEnv("MONGO_INITDB_USER")
	pwd  = GetEnv("MONGO_INITDB_PWD")
	coll = GetEnv("MONGO_COLLECTION")
	addr = GetEnv("MONGO_CONN")
)

func NewClient(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(addr).
			SetAuth(options.Credential{
				AuthSource: db,
				Username:   user,
				Password:   pwd,
			}))

	if err != nil {
		return nil, errors.New("invalid mongodb options")
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.New("cannot connect to mongodb instance")
	}

	return client, nil
}

func UpsertOne(client *mongo.Client, ctx context.Context, user *User) error {
	collection := client.Database(db).Collection(coll)

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"id": user.ID}

	update := bson.M{"$set": bson.M{"age": user.Age, "name": user.Name,
		"salary": user.Salary, "greeting": user.Greeting, "power": user.Power}}

	_, err := collection.UpdateOne(ctx, filter, update, opts)

	return err
}

func FindOne(client mongo.Client, ctx context.Context, id primitive.ObjectID) (*User, error) {

	collection := client.Database(db).Collection(coll)

	var user User
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func Find(client *mongo.Client, ctx context.Context) (*[]User, error) {
	collection := client.Database(db).Collection(coll)

	var users []User

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}

	return &users, nil
}

func loadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("MONGO_CONN"); !runningInContainer {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func GetEnv(k string) string {
	value, ok := os.LookupEnv(k)

	if !ok {
		log.Fatal("Enviroment variable not found: ", k)
	}

	return value
}
