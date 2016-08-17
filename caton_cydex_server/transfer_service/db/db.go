package db
//it is related to database
import (

	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Engine *xorm.Engine
)

func CreateEngine(driver, url string, show_sql bool) error {
	var err error
	if Engine, err = xorm.NewEngine(driver, url); err != nil {
		return err
	}
	Engine.ShowSQL(show_sql)
	return nil
}

func SetupEngine(e *xorm.Engine) error {
	if Engine != nil {
		return errors.New("Engine has been setup")
	}
	Engine = e
	return nil
}

func SyncTables(tables []interface{}) error {
	for _, t := range tables {
		if err := Engine.Sync2(t); err != nil {
			return err
		}
	}
	return nil
}
