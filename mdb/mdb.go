package mdb

import (
	"github.com/globalsign/mgo"
	"log"
	"github.com/globalsign/mgo/bson"
	"InclusION/static"
)

func InitDB() *mgo.Database {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	//defer session.Close()

	db := session.DB("inclusion")

	// make username unique in tblUsers
	indexUser := mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	}
	db.C(static.TBL_USERS).EnsureIndex(indexUser)

	return db
}


func Insert(collection string, obj interface{}) error {

	db := InitDB()

	err := db.C(collection).Insert(obj)

	if err != nil {
		return err
	}

	return nil
}


func QueryById(collection string, objectId bson.ObjectId) interface{} {

	db := InitDB()
	c := db.C(collection)

	var result interface{}
	err := c.Find(bson.M{"_id": objectId}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}


func QueryAll(collection string) interface{} {
	db := InitDB()
	c := db.C(collection)

	var results interface{}
	err := c.Find(nil).All(&results)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

func Delete(collection string, objectId bson.ObjectId) error {
	db := InitDB()
	c := db.C(collection)

	err := c.RemoveId(objectId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateById(collection string, objectId bson.ObjectId, updateObject interface{}) error {
	db := InitDB()

	c := db.C(collection)
	//id := bson.M{"_id": objectId}
	//log.Println(id)

	err := c.UpdateId(objectId, updateObject)
	if err != nil {
		return err
	}

	return nil
}


func UpdateByKey(collection string, key string, value string, updateObject interface{}) error {
	db := InitDB()

	c := db.C(collection)
	where := bson.M{key : value}

	err := c.Update(where, updateObject)

	if err != nil {
		return err
	}

	return nil
}

