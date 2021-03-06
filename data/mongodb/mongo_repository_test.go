package mongodb

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id       bson.ObjectId `bson:"_id" json:"id"`
	Nick     string        `bson:"nick" json:"nick"`
	Age      int           `bson:"age" json:"age"`
	Type     int           `bson:"type" json:"type"`
	Ctime    time.Time     `bson:"ctime" json:"ctime"`
	TenantId string        `bson:"tenantId" json:"tenantId"`
}

func (self *User) Database() string {
	return fmt.Sprintf("%s_user_db", self.TenantId)
}

func (self *User) Collection() string {
	return fmt.Sprintf("%s_users", self.TenantId)
}

func (self *User) Unique() bson.M {
	return bson.M{"_id": self.Id}
}

func (self *User) Indexes() []mgo.Index {
	return []mgo.Index{
		mgo.Index{Key: []string{"nick"}, Unique: true},
		mgo.Index{Key: []string{"type"}},
	}
}

type UserRepository struct {
	MongoRepository
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		NewMongoRepository(),
	}
}

func TestMongoRepository(t *testing.T) {
	config := &Config{
		Addrs:          "localhost:27017",
		Username:       "admin",
		Password:       "admin",
		Database:       "admin",
		ReplicaSetName: "",
		Poolsize:       200,
		Source:         "admin",
		Mode:           2,
	}

	err := RegisterDataSource("default", config)
	if err != nil {
		t.Errorf("register dataSource err: %v\n", err)
		return
	}

	userRepository := NewUserRepository()

	user := User{
		Id:       bson.NewObjectId(),
		Nick:     "Marry张",
		Age:      23,
		Ctime:    time.Now(),
		TenantId: "t1",
	}

	userRepository.EnsureIndexes(&user)

	err = userRepository.Insert(&user)
	if err != nil {
		t.Errorf("insert err: %v", err)
	}

	user.Nick = "哈哈"
	updated, err := userRepository.Update(&user)
	if err != nil {
		t.Errorf("update err: %v", err)
	}
	t.Logf("updated: %d, user: %v\n", updated, user)

	m := User{Id: user.Id, TenantId: "t1"}
	err = userRepository.FindOne(&m)
	if err != nil {
		t.Errorf("find err: %v", err)
	}
	t.Logf("find: %v\n", m)

	err = userRepository.Delete(&m)
	if err != nil {
		t.Errorf("delete err: %v", err)
	}
}
