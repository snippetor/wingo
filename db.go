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
	AutoMigrate(model SqlModel)

	Create(model SqlModel) bool
	Find(model SqlModel) bool
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
