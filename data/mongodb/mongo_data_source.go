package mongodb

import (
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type DataSource struct {
	config  *Config
	session *mgo.Session
}

// Constructor
func NewDataSource(config *Config) *DataSource {
	self := &DataSource{
		config: config,
	}
	return self
}

func (self *DataSource) Connect() error {
	var info = &mgo.DialInfo{
		Addrs:     strings.Split(self.config.Addrs, ";"),
		Username:  self.config.Username,
		Password:  self.config.Password,
		Database:  self.config.Database,
		PoolLimit: self.config.Poolsize,
		Source:    self.config.Source,
		Timeout:   time.Second * 5,
	}
	if len(self.config.ReplicaSetName) > 0 {
		info.ReplicaSetName = self.config.ReplicaSetName
	}
	var GlobalMgoSession, err = mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal(err)
		return err
	}
	GlobalMgoSession.SetSocketTimeout(time.Second * 10)
	GlobalMgoSession.SetSyncTimeout(time.Second * 10)

	mode := self.config.Mode
	// 没有配置就取primary
	if mode <= 0 {
		mode = int(mgo.Primary)
	}
	GlobalMgoSession.SetMode(mgo.Mode(mode), true)

	self.session = GlobalMgoSession

	return nil
}

func (s *DataSource) GetSession() *mgo.Session {
	return s.session
}
