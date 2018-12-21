package model

import (
	"time"
	"github.com/globalsign/mgo/bson"
	"log"
	"github.com/InclusION/static"
	"github.com/InclusION/mdb"
	"github.com/mitchellh/mapstructure"
)

type User struct {
	MongoID bson.ObjectId `bson:"_id,omitempty"`
	Username string
	Password string
	LoginNonce uint64

	Email string
	FirstName string
	LastName string
	Weight float32
	Height float32
	Class int

	FatherName string
	FatherPhone string

	MotherName string
	MotherPhone string

	CoachName string
	CoachPhone string

	AdditionalInfo string
	Timestamp time.Time

	Token string
	TokenExpiryTime time.Time

	DeletedAt time.Time


	// chat
	ClientId string
}


func NewUser(username string) User {
	u := User{Username: username}
	return u
}


func (u *User) ToMap() map[string]string {
	var uMap = make(map[string]string)
	uMap["username"] = u.Username
	uMap["password"] = u.Password
	uMap["email"] = u.Email

	return uMap
}


func (u *User) Insert() error {

	u.Timestamp = time.Now().UTC()
	err := mdb.Insert(static.TBL_USERS, u)
	if err  != nil {
		return err
	}

	log.Println("Inserted")
	return nil
}

func (u *User) QueryAll() []User {
	results := mdb.QueryAll(static.TBL_USERS)

	//log.Printf("Results %s", results)

	var users []User

	for _, u := range results {

		var user User
		err := mapstructure.Decode(u, &user)
		if err  != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	return users
}


func (u *User) QueryByUsernameAndPassword() (error, User) {

	db := mdb.InitDB()
	c := db.C(static.TBL_USERS)

	var result User
	err := c.Find(bson.M{"username": u.Username, "password": u.Password}).One(&result)
	if err != nil {
		log.Println(err)
		return err, result
	}

	return nil, result
}

func (u *User) QueryByUsername() (error, User) {

	db := mdb.InitDB()
	c := db.C(static.TBL_USERS)

	var result User
	err := c.Find(bson.M{"username": u.Username}).One(&result)
	if err != nil {
		return err, result
	}

	return nil, result
}

func (u *User) UpdateById() error  {

	u.Timestamp = time.Now().UTC()

	updateObject := bson.M{"$set":
		bson.M{
		"password" : u.Password, "email": u.Email, "firstname" : u.FirstName, "lastname" : u.LastName,
		"weight" : u.Weight, "height" : u.Height, "class" : u.Class,
		"fathername" : u.FatherName, "fatherphone" : u.FatherPhone ,
		"mothername" : u.MotherName, "motherphone" : u.MotherPhone,
		"coachname" : u.CoachName, "coachphone" : u.CoachPhone,
		"additionalinfo" : u.AdditionalInfo,"timestamp": u.Timestamp,
		"loginnonce": u.LoginNonce,
		"token" : u.Token, "tokenexpirytime" : u.TokenExpiryTime }}

	err := mdb.UpdateById(static.TBL_USERS, u.MongoID, updateObject)
	if err != nil {
		return err
	}

	log.Println("Updated")
	return nil
}

func (u *User) UpdateByUsername() error  {

	u.Timestamp = time.Now().UTC()

	updateObject := bson.M{"$set":
	bson.M{
		"password" : u.Password, "email": u.Email, "firstname" : u.FirstName, "lastname" : u.LastName,
		"weight" : u.Weight, "height" : u.Height, "class" : u.Class,
		"fathername" : u.FatherName, "fatherphone" : u.FatherPhone ,
		"mothername" : u.MotherName, "motherphone" : u.MotherPhone,
		"coachname" : u.CoachName, "coachphone" : u.CoachPhone,
		"additionalinfo" : u.AdditionalInfo,"timestamp": u.Timestamp,
		"loginnonce": u.LoginNonce,
		"token" : u.Token, "tokenexpirytime" : u.TokenExpiryTime}}


	err := mdb.UpdateByKey(static.TBL_USERS,"username", u.Username, updateObject)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Updated")
	return nil
}

func (u *User) HardDelete() error {

	err := mdb.Delete(static.TBL_USERS, u.MongoID)
	if err != nil {
		return err
	}

	log.Println("Hard Deleted")
	return nil
}

func (u *User) SoftDelete() error {

	u.DeletedAt = time.Now().UTC()
	updateObject := bson.M{"$set": bson.M{"deletedat" : u.DeletedAt}}

	err := mdb.UpdateByKey(static.TBL_USERS,"username", u.Username, updateObject)
	if err != nil {
		return err
	}

	log.Println("Soft Deleted")
	return nil
}


