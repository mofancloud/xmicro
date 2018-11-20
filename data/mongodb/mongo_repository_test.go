package mongodb

import (
	"fmt"
	"testing"
	"time"
)

type User struct {
	Id       string    `bson:"_id" json:"id"`
	Nick     string    `bson:"nick" json:"nick"`
	Age      int       `bson:"age" json:"age"`
	Ctime    time.Time `bson:"ctime" json:"ctime"`
	TenantId string    `bson:"tenantId" json:"tenantId"`
}

func (self *User) Database() string {
	return fmt.Sprintf("%s_user_db", self.TenantId)
}

func (self *User) Collection() string {
	return fmt.Sprintf("%s_users", self.TenantId)
}

func (self *User) Unique() string {
	return self.Id
}

func TestMongoRepository(t *testing.T) {
	Config := Config{
		Hosts:    "localhost:27017",
		Username: "",
		Password: "",
		Database: "admin",
	}

	userRepository := NewMongoRepository(Config)

	user := User{
		Id:       bson.NewObjectId(),
		Nick:     "Marry张",
		Age:      23,
		Ctime:    time.Unix(),
		TenantId: "t1",
	}

	err := userRepository.Insert(&user)
	if err != nil {
		t.Errorf("ParsePageQueryFromRequest err: %v", err)
	}
}
