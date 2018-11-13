package db

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
)


func innit() *mongo.Client {

	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil { log.Println(err) }

	err = client.Connect(context.TODO())
	if err != nil { log.Println(err) }

	//colUsers = client.Database("inclusion").Collection("tblUser")
	return client
}


func Insert(collectionName string, object MongoObject ) {

	client := innit()

	//user := NewUser("huynhbinh", "12344566")
	//var data = make(map[string]string)
	//data["name"] = "huynhbinh"
	//data["age"] = "12"
	//data["email"] = "huynhbinh@gmail.com"
	//user := User{username:"huynbinh", password: "000000", email: "000000@gmail.co"}

	collection := client.Database("inclusion").Collection(collectionName)

	res, err := collection.InsertOne(context.Background(), object.ToMap())
	if err != nil { log.Fatal(err) }
	id := res.InsertedID

	log.Println(id)
}

func Update(collectionName string, oldObj MongoObject, newObj MongoObject) {

	client := innit()

	collection := client.Database("inclusion").Collection(collectionName)

	res, err := collection.UpdateOne(context.Background(), oldObj, newObj);
	if err != nil { log.Fatal(err) }

	id := res.UpsertedID
	log.Println(id)
}

func delete() {

}

func LoadAll(collectionName string)  {

	client := innit()

	collection := client.Database("inclusion").Collection(collectionName)

	cur, err := collection.Find(context.Background(), nil)
	if err != nil { log.Fatal(err) }

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		elem := bson.NewDocument()
		err := cur.Decode(elem)
		if err != nil { log.Fatal(err) }
		// do something with elem....

		log.Println(elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

}


func Load(collectionName string) {
	client := innit()

	collection := client.Database("inclusion").Collection(collectionName)


	filter := make(map[string]string)
	filter["username"] = "hnb"

	cur, err := collection.Find(context.Background(), bson.NewDocument(bson.EC.String("email", "22222@gmail.io" )))


	if err != nil { log.Fatal(err) }

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		elem := bson.NewDocument()
		err := cur.Decode(elem)
		if err != nil { log.Fatal(err) }
		// do something with elem....

		log.Println(elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
}
