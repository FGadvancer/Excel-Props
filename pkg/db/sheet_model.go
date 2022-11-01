package db

import (
	"gorm.io/gorm"
	"time"
)

//生成模号表
type Sheet struct {
	SheetID            string    `gorm:"column:sheet_id;primary_key;type:char(64)" json:"sheetID"`
	CommodityName      string    `gorm:"column:commodity_name;type:char(64)" json:"commodityName"`
	Version            int32     `gorm:"column:version" json:"version"`
	Code               int32     `gorm:"column:code" json:"code"`
	CreatorUserID      string    `gorm:"column:creator_user_id;size:64"`
	CreateTime         time.Time `gorm:"column:create_time" json:"createTime"`
	LastModifierUserID string    `gorm:"column:last_modifier_userID;size:64" json:"lastModifierUserID"`
	LastModifierIP     string    `gorm:"column:last_modifier_ip;size:64" json:"LastModifierIP"`
	LastModifyTime     time.Time `gorm:"column:last_modify_time" json:"lastModifyTime"`
	Ex                 string    `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
	DB                 *gorm.DB  `gorm:"-" json:"-"`
}

func NewSheet(DB *gorm.DB) *Sheet {
	return &Sheet{DB: DB}
}
