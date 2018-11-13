package model

import (
	"github.com/globalsign/mgo/bson"
	"time"
	"image"
	"github.com/pkg/errors"
	"github.com/InclusION/mdb"
	"github.com/InclusION/static"
)

type Blog struct {
	MongoID bson.ObjectId `bson:"_id,omitempty"`

	Tittle string
	BlogContent string
	Thumbnail image.Image

	CreatorName string
	CreatedTime time.Time
	UpdatedTime time.Time
}


func NewBlog(id bson.ObjectId) Blog {
	b := Blog{MongoID: id}
	return b
}


func (blog *Blog) QueryAllPaging(page int) (error, []Blog) {

	var result []Blog

	db := mdb.InitDB()
	c := db.C(static.TBL_BLOGS)

	// Use Skip() and Limit() to denote the page you want to send, e.g.
	q := c.Find(nil).Sort("createdtime").Skip((page - 1)*20).Limit(20)
	err := q.All(&result)
	if err != nil {
		return err, result
	}

	return nil, result
}

func (blog *Blog) QueryById() (error, Blog) {

	result := mdb.QueryById(static.TBL_BLOGS, blog.MongoID)

	blg, ok := result.(Blog)
	if ok {
		return nil, blg
	}

	return errors.New("Type missed match"), blg
}


func (blog *Blog) Insert() error {

	err := mdb.Insert(static.TBL_BLOGS, blog)
	if err != nil {
		return err
	}

	return nil
}
