package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Model interface {
	Unique() bson.M
	Collection() string
	Database() string
}

type Indexed interface {
	Indexes() []mgo.Index
	Model
}
