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

type MongoRepository interface {
	All(m Model, result interface{}) error
	Count(m Model) (count int64, err error)
	Update(m Model) (updated int, err error)
	UpdateSelective(m Model, updateData map[string]interface{}) error
	Insert(m Model) error
	Upsert(m Model) (upserted int, err error)
	FindOne(in Model, out Model) error
	Delete(m Model) error
	Page(pageQuery *data.PageQuery, m Model, list interface{}) (total int64, pageNo int64, pageSize int32, err error)
	Execute(m Model, fn DBFunc) error
	EnsureIndexes(m Indexed)
}

type mongoRepositoryImpl struct {
	dataSource *DataSource
}

// Constructor
func NewMongoRepository(dataSource *DataSource) MongoRepository {
	return &mongoRepositoryImpl{
		dataSource: dataSource,
	}
}

func (self *mongoRepositoryImpl) GetDataSource() *DataSource {
	return self.dataSource
}

func (self *mongoRepositoryImpl) All(m Model, result interface{}) error {
	err := self.Execute(m, func(c *mgo.Collection) error {
		return Where(c, nil).All(result)
	})
	return err
}

func (self *mongoRepositoryImpl) Count(m Model) (count int64, err error) {
	self.Execute(m, func(c *mgo.Collection) error {
		c1, err := Where(c, nil).Count()
		count = int64(c1)
		return err
	})
	return
}

func (self *mongoRepositoryImpl) Update(m Model) (updated int, err error) {
	self.Execute(m, func(c *mgo.Collection) error {
		info, err := c.Find(m.Unique()).Apply(mgo.Change{
			ReturnNew: true,
			Update: bson.M{
				"$set": m,
			},
		}, m)

		if err != nil {
			return err
		}

		updated = info.Updated
		return nil
	})

	return
}

func (self *mongoRepositoryImpl) UpdateSelective(m Model, updateData map[string]interface{}) error {
	err := self.Execute(m, func(c *mgo.Collection) error {
		return c.Update(m, bson.M{"$set": updateData})
	})
	return err
}

func (self *mongoRepositoryImpl) Insert(m Model) error {
	err := self.Execute(m, func(c *mgo.Collection) error {
		return c.Insert(m)
	})
	return err
}

func (self *mongoRepositoryImpl) Upsert(m Model) (upserted int, err error) {
	self.Execute(m, func(c *mgo.Collection) error {
		changeInfo, err := c.Upsert(m.Unique(), bson.M{"$set": m})
		if err != nil {
			return err
		}
		upserted = changeInfo.Updated
		return nil
	})

	return
}

func (self *mongoRepositoryImpl) FindOne(in Model, out Model) error {
	return self.Execute(in, func(c *mgo.Collection) error {
		err := c.Find(in.Unique()).One(out)
		return err
	})
}

func (self *mongoRepositoryImpl) Delete(m Model) error {
	return self.Execute(m, func(c *mgo.Collection) error {
		return c.Remove(m.Unique())
	})
}

func (self *mongoRepositoryImpl) Page(pageQuery *data.PageQuery, m Model, list interface{}) (total int64, pageNo int64, pageSize int32, err error) {
	filters, pageNo, pageSize, _ := ParsePageQuery(m, pageQuery)

	self.Execute(m, func(c *mgo.Collection) error {
		t, err := c.Find(filters).Count()
		total = int64(t)
		if err != nil {
			return err
		}

		return Page(c, pageQuery, m, list)
	})

	return
}

func (self *mongoRepositoryImpl) Execute(m Model, fn DBFunc) error {
	return Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), fn)
}

func (self *mongoRepositoryImpl) EnsureIndexes(m Indexed) {
	Execute(self.dataSource.GetSession(), m.Database(), m.Collection(), func(c *mgo.Collection) error {
		for _, i := range m.Indexes() {
			c.EnsureIndex(i)
		}

		return nil
	})
}
