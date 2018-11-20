package mongodb

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	base_model "base/model"
)

type Config struct {
	Hosts          string `json:"hosts"`
	Username       string `json:"user"`
	Password       string `json:"password"`
	Database       string `json:"database"`
	ReplicaSetName string `json:"replicaSet"`
	Poolsize       int    `json:"poolsize"`
	Source         string `json:"source"`
	Mode           int32  `json:"mode"`
}

type MongoRepository struct {
	dataSource *DataSource
}

// Constructor
func NewRepository(dataSource *DataSource) *MongoRepository {
	self := &MongoRepository{
		dataSource: dataSource,
	}
	return self
}

func (self *MongoRepository) GetDataSource() *DataSource {
	return self.DataSource
}

func (self *MongoRepository) All(m Model, result interface{}) error {
	session := self.DataSource.GetSession()
	db := session.DB(self.DataSource.Database)
	defer session.Close()
	return Where(db, m, nil).All(result)
}

func (self *MongoRepository) Count(m Model) (int64, error) {
	session := self.DataSource.GetSession()
	db := session.DB(self.DataSource.Database)
	defer session.Close()
	r, err := Where(db, m, nil).Count()
	return int64(r), err
}

func (self *MongoRepository) Find(m Model) *mgo.Query {
	session := self.DataSource.GetSession()
	db := session.DB(self.DataSource.Database)

	return Where(db, m, m.Unique()).Limit(1)
}

func (self *MongoRepository) Where(db *mgo.Database, m Model, q interface{}) *mgo.Query {
	return db.C(m.Collection()).Find(q)
}

func (self *MongoRepository) Update(m Model) (*mgo.ChangeInfo, error) {
	session := self.DataSource.GetSession()
	db := session.DB(self.DataSource.Database)
	defer session.Close()

	return Find(db, m).Apply(mgo.Change{
		ReturnNew: true,
		Update: bson.M{
			"$set": m,
		},
	}, m)
}

func (self *MongoRepository) UpdateSelective(m Model, updateData bson.M) error {
	session := self.DataSource.GetSession()
	db := session.DB(self.DataSource.Database)
	defer session.Close()

	return db.C(m.Collection()).Update(m.Unique(), bson.M{"$set": updateData})
}

func (self *MongoRepository) Insert(m Model) error {
	session := self.DataSource.GetSession()
	db := session.DB(self.DataSource.Database)
	defer session.Close()

	return db.C(m.Collection()).Insert(m)
}

func (self *MongoRepository) Upsert(m Model) (*mgo.ChangeInfo, error) {
	session := self.DataSource.GetSession()
	db := session.DB(self.DataSource.Database)
	defer session.Close()
	changeInfo, err := db.C(m.Collection()).Upsert(m.Unique(), bson.M{"$set": m})
	return changeInfo, err
}

func (self *MongoRepository) Delete(m Model) error {
	session := self.DataSource.GetSession()
	db := session.DB(self.DataSource.Database)
	defer session.Close()

	return db.C(m.Collection()).Remove(m.Unique())
}

func (self *MongoRepository) Page(pageQuery *base_model.PageQuery, m Model, list interface{}) (int64, int64, int32, error) {
	session := self.DataSource.GetSession()
	defer session.Close()

	filters, pageNo, pageSize, sorts := ParsePageQuery(m, pageQuery)

	c := session.DB(self.DataSource.Database).C(m.Collection())
	ms := bson.M(filters)
	total, _ := c.Find(ms).Count()

	offset := int((int32(pageNo) - 1) * pageSize)
	limit := int(pageSize)

	err := c.Find(ms).Skip(offset).Limit(limit).Sort(sorts...).All(list)

	return int64(total), pageNo, pageSize, err
}

func (self *MongoRepository) ensureIndexes(m Indexed) {
	session := self.DataSource.GetSession()
	db := session.DB(self.DataSource.Database)
	defer session.Close()

	coll := db.C(m.Collection())
	for _, i := range m.Indexes() {
		coll.EnsureIndex(i)
	}
}
