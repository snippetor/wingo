package wingo

import (
	"github.com/jinzhu/gorm"
)

type SqlModel interface {
	init(db *gorm.DB, self interface{})
	DB() *gorm.DB
	Sync(cols ...interface{})
	SyncInTx(tx *gorm.DB, cols ...interface{})
	Del()
	DelInTx(tx *gorm.DB)
	FieldToString(f interface{}) string
	FieldFromString(s string, f interface{})
}
