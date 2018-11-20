package mongodb

import (
	"github.com/mofancloud/xmicro/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Config struct {
	Hosts          string `json:"hosts"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Database       string `json:"database"`
	ReplicaSetName string `json:"replicaSetName"`
	Poolsize       int    `json:"poolsize"`
	Source         string `json:"source"`
	Mode           int    `json:"mode"`
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
	return self.dataSource
}

func (self *MongoRepository) All(m Model, result interface{}) error {
	err := Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		return Where(c, nil).All(result)
	})
	return err
}

func (self *MongoRepository) Count(m Model) (count int64, err error) {
	Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		c1, err := Where(c, nil).Count()
		count = int64(c1)
		return err
	})
	return
}

func (self *MongoRepository) Update(m Model) (info *mgo.ChangeInfo, err error) {
	Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		info, err = c.Find(m.Unique()).Apply(mgo.Change{
			ReturnNew: true,
			Update: bson.M{
				"$set": m,
			},
		}, m)
		return err
	})

	return
}

func (self *MongoRepository) UpdateSelective(m Model, updateData bson.M) error {
	err := Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		return c.Update(m.Unique(), bson.M{"$set": updateData})
	})
	return err
}

func (self *MongoRepository) Insert(m Model) error {
	err := Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		return c.Insert(m)
	})
	return err
}

func (self *MongoRepository) Upsert(m Model) (changeInfo *mgo.ChangeInfo, err error) {
	Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		changeInfo, err = c.Upsert(m.Unique(), bson.M{"$set": m})
		return err
	})

	return
}

func (self *MongoRepository) FindOne(m Model) error {
	return Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		err := c.Find(m.Unique()).One(m)
		return err
	})
}

func (self *MongoRepository) Delete(m Model) error {
	return Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		return c.Remove(m.Unique())
	})
}

func (self *MongoRepository) Page(pageQuery *data.PageQuery, m Model, list interface{}) (total int64, pageNo int64, pageSize int32, err error) {
	filters, pageNo, pageSize, _ := ParsePageQuery(m, pageQuery)

	Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		t, err := c.Find(filters).Count()
		total = int64(t)
		if err != nil {
			return err
		}

		return Page(c, pageQuery, m, list)
	})

	return
}

func (self *MongoRepository) EnsureIndexes(m Indexed) {
	Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		for _, i := range m.Indexes() {
			c.EnsureIndex(i)
		}

		return nil
	})
}
