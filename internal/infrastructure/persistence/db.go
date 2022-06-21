package persistence

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{}) {

}

func NewDB(dbUsername, dbPassword, dbIP, dbName string, dbPort int) (*sql.DB, *gorm.DB) {
	loggerDefine := New(zap.L())
	loggerDefine.SetAsDefault()
	connArgs := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbUsername, dbPassword, dbIP, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(connArgs), &gorm.Config{
		Logger: loggerDefine,
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return sqlDB, db
}
