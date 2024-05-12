package dao

import (
	"fmt"
	"time"

	fllog "github.com/forlifeproj/msf/log"
)

const (
	AccountTable = "table_account"
	OpenidTable  = "table_openid"
	UnionidTable = "table_unionid"
)

type Account struct {
	Uid        int64  `gorm:"column:uid;primaryKey;autoIncrement"`
	OpenId     string `gorm:"column:openid"`
	UnionId    string `gorm:"column:unionid"`
	OpenType   int    `gorm:"column:open_type"`
	Status     int    `gorm:"column:status"`
	UpdateTime int64  `gorm:"column:update_time"`
	CreateTime int64  `gorm:"column:create_time"`
}

func AddAccount(record *Account) error {
	result := gormDb.Table(AccountTable).Save(record)
	return result.Error
}

func GetAccount(uid int64) (*Account, error) {
	var account Account
	result := gormDb.Table(AccountTable).First(&account).Where("uid = ?", uid)
	return &account, result.Error
}

type OpenID struct {
	OpenId     string `gorm:"column:openid;primaryKey"`
	OpenType   int    `gorm:"column:open_type;primaryKey"`
	Uid        int64  `gorm:"column:uid"`
	UpdateTime int64  `gorm:"column:update_time"`
	CreateTime int64  `gorm:"column:create_time"`
}

func AddOpenID(record *OpenID) error {
	result := gormDb.Table(OpenidTable).Save(record)
	return result.Error
}

func GetOpenID(openId string, openType int) (*OpenID, error) {
	var openID OpenID
	result := gormDb.Table(OpenidTable).First(&openID).Where("openid = ? and open_type = ?", openId, openType)
	return &openID, result.Error
}

type UnionID struct {
	UnionId    string `gorm:"column:unionid;primaryKey"`
	OpenType   int    `gorm:"column:open_type;primaryKey"`
	Uid        int64  `gorm:"column:uid"`
	UpdateTime int64  `gorm:"column:update_time"`
	CreateTime int64  `gorm:"column:create_time"`
}

func AddUnionID(record *OpenID) error {
	result := gormDb.Table(UnionidTable).Save(record)
	return result.Error
}

func GetUnionID(unionId string, openType int) (*UnionID, error) {
	var unionID UnionID
	result := gormDb.Table(UnionidTable).First(&unionID).Where("openid = ? and open_type = ?", unionId, openType)
	return &unionID, result.Error
}

func RegisterUid(openId, unionId string, openType int) (int64, error) {
	tx := gormDb.Begin()
	var err error
	nowTime := time.Now().Unix()
	addAccount := &Account{
		OpenId:     openId,
		OpenType:   openType,
		UnionId:    unionId,
		Status:     1,
		CreateTime: nowTime,
		UpdateTime: nowTime,
	}

	for {
		// insert into account
		result := tx.Table(AccountTable).Save(addAccount)
		if result.Error != nil {
			fllog.Log().Error(fmt.Sprintf("insert into table_account failed. err:%+v addAccount:%+v",
				result.Error, addAccount))
			err = result.Error
			break
		}

		// insert into openid
		addOpenID := &OpenID{
			OpenId:     openId,
			OpenType:   openType,
			Uid:        addAccount.Uid,
			CreateTime: nowTime,
			UpdateTime: nowTime,
		}
		result = tx.Table(OpenidTable).Save(addOpenID)
		if result.Error != nil {
			fllog.Log().Error(fmt.Sprintf("insert into table_openid failed. err:%+v addOpenID:%+v", result.Error, addOpenID))
			err = result.Error
			break
		}

		// insert into unionid
		if len(unionId) > 0 {
			addUnionID := &UnionID{
				UnionId:    unionId,
				OpenType:   openType,
				Uid:        addAccount.Uid,
				CreateTime: nowTime,
				UpdateTime: nowTime,
			}
			result = tx.Table(UnionidTable).Save(addUnionID)
			if result.Error != nil {
				fllog.Log().Error(fmt.Sprintf("insert into table_unionid failed. err:%+v addUnionID:%+v", result.Error, addUnionID))
				err = result.Error
				break
			}
		}
		break
	}

	if err == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}

	return addAccount.Uid, err
}
