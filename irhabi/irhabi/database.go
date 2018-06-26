package irhabi

import (
	"fmt"
	"time"

	"github.com/alfatih/irhabi/common/log"
	"github.com/alfatih/irhabi/orm"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// DbSetup registering database connection
// by reading database config variables.
func DbSetup() error {
	orm.Debug = IsDebug()
	orm.DebugLog = log.Log
	orm.DefaultTimeLoc = time.Local
	orm.DefaultRelsDepth = 3

	ds := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", Config.DbUser, Config.DbPassword, Config.DbHost, Config.DbName, "charset=utf8&loc=Asia%2FJakarta")
	return orm.RegisterDataBase("default", Config.DbEngine, ds)
}
