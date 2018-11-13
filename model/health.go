package model

import (
	"github.com/globalsign/mgo/bson"
	"time"
	"InclusION/mdb"
	"log"
	"InclusION/static"
)

type Health struct {
	MongoID bson.ObjectId `bson:"_id,omitempty"`
	Username string `json:"username"`
	Token string

	HeartRate float32
	Temperature float32
	UV float32
	GPS string
	BloodPressure float32
	IsMeltdown bool

	Timestamp time.Time
	DeletedAt time.Time
}


func (h *Health) Insert() error {

	h.Timestamp = time.Now().UTC()
	err := mdb.Insert(static.TBL_HEALTHS, h)
	if err  != nil {
		return err
	}

	return nil
}

func (h *Health) QueryAll() []Health {
	result := mdb.QueryAll(static.TBL_HEALTHS)
	healths, ok := result.([]Health)

	if ok == false {
		log.Println("No results")
		return nil
	} else {
		return healths
	}
}

func (h *Health) QueryById() {

	result := mdb.QueryById(static.TBL_HEALTHS, h.MongoID)

	health, ok := result.(Health)

	if ok == false {
		log.Println("No results")
	} else {
		h = &health
	}
}

func (h *Health) QueryByUser(user User) (error, []Health) {
	db := mdb.InitDB()

	c := db.C(static.TBL_HEALTHS)

	var result []Health
	err := c.Find(bson.M{"username": user.Username}).All(&result)
	if err != nil {
		return err, nil
	}

	return nil, result
}

func (h *Health) UpdateById() {

	updateObject := bson.M{"$set":
		bson.M{"heartrate": h.HeartRate,
			"temperature": h.Temperature,
			"uv": h.UV,
			"gps": h.GPS,
			"bloodpressure": h.BloodPressure,
			"ismeltdown": h.IsMeltdown,
			"timestamp": time.Now()}}

	err := mdb.UpdateById(static.TBL_HEALTHS, h.MongoID, updateObject)
	if err != nil {
		log.Println(err)
	}

	log.Println("Updated")
}

func (h *Health) HardDelete() error {

	err := mdb.Delete(static.TBL_HEALTHS, h.MongoID)
	if err != nil {
		return err
	}

	log.Println("Hard Deleted")
	return nil
}

func (h *Health) SoftDelete() error {

	h.DeletedAt = time.Now().UTC()
	updateObject := bson.M{"$set": bson.M{"deletedat" : h.DeletedAt}}

	err := mdb.UpdateById(static.TBL_HEALTHS, h.MongoID, updateObject)
	if err != nil {
		return err
	}

	log.Println("Soft Deleted")
	return nil

}
