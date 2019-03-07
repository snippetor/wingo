package wingo

import (
	"github.com/jinzhu/gorm"
)

type SqlDBType byte

type sqlDB struct {
	tbPrefix string
	db       *gorm.DB
}

func (m *sqlDB) Dial(dbType SqlDBType, host, port, user, pwd, defaultDb, tbPrefix string, debugMode bool) {
	m.tbPrefix = tbPrefix
	// db
	var db *gorm.DB
	var err error
	switch dbType {
	case MySql:
		db, err = gorm.Open("mysql", user+":"+pwd+"@tcp("+host+":"+port+")/"+defaultDb+"?charset=utf8&parseTime=True&loc=Local")
	case MsSql:
		db, err = gorm.Open("mssql", "sqlserver://"+user+":"+pwd+"@"+host+":"+port+"?database="+defaultDb)
	case Postgres:
		db, err = gorm.Open("postgres", "host="+host+" port="+port+" user="+user+" dbname="+defaultDb+" password="+pwd)
	}
	CheckError(err)
	m.db = db
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return m.tbPrefix + "_" + defaultTableName
	}
	db.LogMode(debugMode)
}

func (m *sqlDB) DB() *gorm.DB {
	if m.db == nil {
		panic("- DB not be initialized, invoke Dial at first!")
	}
	return m.db
}

func (m *sqlDB) TableName(tbName string) string {
	if m.tbPrefix != "" {
		return m.tbPrefix + "_" + tbName
	}
	return tbName
}

func (m *sqlDB) AutoMigrate(model interface{}) {
	if mod, ok := model.(interface {
		Init(*gorm.DB, interface{})
	}); ok {
		mod.Init(m.db, model)
	}
	m.DB().AutoMigrate(model)
}

func (m *sqlDB) Create(model interface{}) bool {
	if mod, ok := model.(interface {
		Init(*gorm.DB, interface{})
	}); ok {
		mod.Init(m.db, model)
	}
	res := m.DB().Create(model)
	if res.Error != nil {
		panic(res.Error)
	}
	return true
}

func (m *sqlDB) Find(model interface{}) bool {
	if mod, ok := model.(interface {
		Init(*gorm.DB, interface{})
	}); ok {
		mod.Init(m.db, model)
	}
	res := m.DB().Where(model).First(model)
	return res.Error == nil
}

func (m *sqlDB) FindAll(models interface{}) bool {
	res := m.DB().Find(models)
	if res.Error != nil {
		panic(res.Error)
	}
	return true
}

func (m *sqlDB) FindMany(models interface{}, limit int, orderBy string, whereAndArgs ...interface{}) bool {
	db := m.DB()
	if limit > 0 {
		db = db.Limit(limit)
	}
	if orderBy != "" {
		db = db.Order(orderBy)
	}
	if len(whereAndArgs) > 0 && len(whereAndArgs)%2 == 0 {
		var args = make(map[string]interface{})
		for i := 0; i < len(whereAndArgs); i += 2 {
			args[whereAndArgs[i].(string)] = whereAndArgs[i+1]
		}
		db = db.Where(args)
	}
	db = db.Find(models)
	if db.Error != nil {
		panic(db.Error)
	}
	return true
}

func (m *sqlDB) Begin() *gorm.DB {
	return m.DB().Begin()
}

func (m *sqlDB) Rollback() {
	m.DB().Rollback()
}

func (m *sqlDB) Commit() {
	m.DB().Rollback()
}

func (m *sqlDB) Close() {
	if m.db != nil {
		m.db.Close()
	}
}
