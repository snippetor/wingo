package wingo

import (
	"github.com/jinzhu/gorm"
)

var msgpackCodec Codec

func init() {
	msgpackCodec = &MsgPackCodec{}
}

type BaseSqlModel struct {
	db   *gorm.DB
	self interface{}
	Id   uint32 `gorm:"primary_key"`
}

func (m *BaseSqlModel) init(db *gorm.DB, self interface{}) {
	m.db = db
	m.self = self
}

func (m *BaseSqlModel) DB() *gorm.DB {
	return m.db
}

// 更新到数据库
func (m *BaseSqlModel) Sync(cols ...interface{}) {
	if cols != nil && len(cols) > 0 {
		CheckError(m.DB().Model(m.self).UpdateColumn(cols).Error)
	} else {
		CheckError(m.DB().Model(m.self).Updates(m.self).Error)
	}
}

// 更新到数据库
func (m *BaseSqlModel) SyncInTx(tx *gorm.DB, cols ...interface{}) {
	if cols != nil && len(cols) > 0 {
		CheckError(tx.Model(m.self).UpdateColumn(cols).Error)
	} else {
		CheckError(tx.Model(m.self).Updates(m.self).Error)
	}
}

// 从数据库移除，ID必须存在
func (m *BaseSqlModel) Del() {
	CheckError(m.DB().Delete(m.self).Error)
}

// 从数据库移除，ID必须存在
func (m *BaseSqlModel) DelInTx(tx *gorm.DB) {
	CheckError(tx.Delete(m.self).Error)
}

func (m *BaseSqlModel) FieldToString(f interface{}) string {
	bs, err := msgpackCodec.Marshal(f)
	CheckError(err)
	return Bytes2String(bs)
}

func (m *BaseSqlModel) FieldFromString(s string, f interface{}) {
	msgpackCodec.Unmarshal(String2Bytes(s), f)
}
