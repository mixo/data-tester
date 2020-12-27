package tester

import (
	"github.com/mixo/gosql"
)

type DBConfig struct {
    Driver string
    Host string
    Port string
    User string
    Password string
    Database string
}

var (
    dbConfigs = []DBConfig{
        DBConfig{"mysql", "127.0.0.1", "3306", "test", "test", "test"},
//         DBConfig{"postgres", "127.0.0.1", "5432", "test", "test", "test"},
    }
)

type FunctionalTest struct {}

func (this FunctionalTest) GetDatabases() (dbs map[string]gosql.DB) {
    dbs = make(map[string]gosql.DB)
    for _, db := range dbConfigs {
        dbs[db.Driver] = gosql.DB{db.Driver, db.Host, db.Port, db.User, db.Password, db.Database}
    }

    return
}
