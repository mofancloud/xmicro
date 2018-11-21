package mongodb

import (
	"fmt"
	"sync"
)

type alias struct {
	Name       string
	DataSource DataSource
}

type _dsCache struct {
	mux   sync.RWMutex
	cache map[string]*alias
}

var (
	dataSourceCache = &_dsCache{cache: make(map[string]*alias)}
)

// add database alias with original name.
func (ac *_dsCache) add(name string, al *alias) (added bool) {
	ac.mux.Lock()
	defer ac.mux.Unlock()
	if _, ok := ac.cache[name]; !ok {
		ac.cache[name] = al
		added = true
	}
	return
}

// get database alias if cached.
func (ac *_dsCache) get(name string) (al *alias, ok bool) {
	ac.mux.RLock()
	defer ac.mux.RUnlock()
	al, ok = ac.cache[name]
	return
}

// get default alias.
func (ac *_dsCache) getDefault() (al *alias) {
	al, _ = ac.get("default")
	return
}

func RegisterDataSource(aliasName string, config *Config) error {
	ds := NewDataSource(config)
	err := ds.Connect()
	if err != nil {
		return err
	}

	al := new(alias)
	al.Name = aliasName
	if !dataSourceCache.add(aliasName, al) {
		return fmt.Errorf("DataBase alias name `%s` already registered, cannot reuse", aliasName)
	}

	al.DataSource = ds

	return nil
}

func GetDataSource(aliasNames ...string) (DataSource, error) {
	var name string
	if len(aliasNames) > 0 {
		name = aliasNames[0]
	} else {
		name = "default"
	}
	al, ok := dataSourceCache.get(name)
	if ok {
		return al.DataSource, nil
	}
	return nil, fmt.Errorf("DataSource of alias name `%s` not found", name)
}
