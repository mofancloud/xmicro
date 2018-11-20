package mongodb

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/mofancloud/xmicro/data"
)

func All(db *mgo.Database, m Model) *mgo.Query {
	return Where(db, m, nil)
}

func Find(db *mgo.Database, m Model) *mgo.Query {
	return Where(db, m, m.Unique()).Limit(1)
}

func Where(db *mgo.Database, m Model, q interface{}) *mgo.Query {
	return db.C(m.Collection()).Find(q)
}

func Update(db *mgo.Database, m Model) (*mgo.ChangeInfo, error) {
	return Find(db, m).Apply(mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": m,
		},
	}, m)
}

func Insert(db *mgo.Database, m Model) error {
	return db.C(m.Collection()).Insert(m)
}

func Delete(db *mgo.Database, m Model) error {
	return db.C(m.Collection()).Remove(m.Unique())
}

func Page(db *mgo.Database, pageQuery *data.PageQuery, m Model, list interface{}) (int64, error) {
	filters, pageNo, pageSize, sorts := ParsePageQuery(m, pageQuery)

	c := db.C(m.Collection())
	total, _ := c.Find(bson.M(filters)).Count()

	offset := int((int32(pageNo) - 1) * pageSize)
	limit := int(pageSize)

	err := c.Find(bson.M(filters)).Skip(offset).Limit(limit).Sort(sorts...).All(&list)
	return int64(total), err
}

func ensureIndexes(db *mgo.Database, m Indexed) {
	coll := db.C(m.Collection())
	for _, i := range m.Indexes() {
		coll.EnsureIndex(i)
	}
}

type DBFunc func(*mgo.Collection) error

func Execute(mongoSession *mgo.Session, databaseName string, collectionName string, fn DBFunc) error {
	db := session.DB(databaseName)
	defer session.Close()

	collection := db.C(collectionName)
	if collection == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		return err
	}
	err := fn(collection)
	if err != nil {
		return err
	}

	return nil
}
