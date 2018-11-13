package db


type MongoObject interface {
	ToMap() map[string]string
}