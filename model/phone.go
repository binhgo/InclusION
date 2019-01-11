package model

import (
	"github.com/globalsign/mgo/bson"
	"github.com/InclusION/mdb"
	"github.com/InclusION/static"
	"log"
)

type Phone struct {
	MongoID  bson.ObjectId `bson:"_id,omitempty"`
	Username string

	FCMToken string
}


func (p *Phone) Insert() error {
	err := mdb.Insert(static.TBL_DEVICES, p)
	if err  != nil {
		return err
	}

	return nil
}

func (p *Phone) HardDelete() error {

	err := mdb.Delete(static.TBL_DEVICES, p.MongoID)
	if err != nil {
		return err
	}

	log.Println("Hard Deleted")
	return nil
}

func (p *Phone) QueryPhonesByUsername(username string) (error, []Phone) {
	db := mdb.InitDB()

	c := db.C(static.TBL_DEVICES)

	var result []Phone
	err := c.Find(bson.M{"username": username}).All(&result)
	if err != nil {
		return err, nil
	}

	return nil, result
}
