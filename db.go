package wingo

import "github.com/jinzhu/gorm"

const (
	MySql SqlDBType = iota
	MsSql
	Postgres
)

type SqlDB interface {
	Dial(dbType SqlDBType, host, port, user, pwd, defaultDb, tbPrefix string, debugMode bool)
	DB() *gorm.DB
	TableName(tbName string) string
	AutoMigrate(model interface{})

	Create(model interface{}) bool
	Find(model interface{}) bool
	FindAll(models interface{}) bool
	FindMany(models interface{}, limit int, orderBy string, whereAndArgs ...interface{}) bool
	Begin() *gorm.DB
	Rollback()
	Commit()
	Close()
}

type NoSqlDB interface {
}

var (
	Sql   SqlDB
	NoSql NoSqlDB
)

func init() {
	Sql = &sqlDB{}
}
