package db

import (
	"Excel-Props/pkg/utils"
	"errors"
	"gorm.io/gorm"
	"time"
)

//生成模号表
type Sheet struct {
	SheetID            string    `gorm:"column:sheet_id;primary_key;type:char(64)" json:"sheetID"`
	CommodityName      string    `gorm:"column:commodity_name;type:char(64)" json:"commodityName"`
	Version            int32     `gorm:"column:version" json:"version"`
	Code               string    `gorm:"column:code;type:varchar(64)" json:"code"`
	CreatorUserID      string    `gorm:"column:creator_user_id;size:64"`
	CreateTime         time.Time `gorm:"column:create_time" json:"createTime"`
	LastModifierUserID string    `gorm:"column:last_modifier_userID;size:64" json:"lastModifierUserID"`
	LastModifierIP     string    `gorm:"column:last_modifier_ip;size:64" json:"LastModifierIP"`
	LastModifyTime     time.Time `gorm:"column:last_modify_time" json:"lastModifyTime"`
	Ex                 string    `gorm:"column:ex;type:varchar(1024)"  json:"ex,omitempty"`
	DB                 *gorm.DB  `gorm:"-" json:"-"`
}

func NewSheet(DB *gorm.DB) *Sheet {
	return &Sheet{DB: DB}
}

func (s *Sheet) GetSheetInfo(sheetID string) (*Sheet, error) {
	temp := Sheet{}
	err := DB.MysqlDB.db.Model(&temp).Where("sheet_id = ? ", sheetID).Take(&temp).Error
	return &temp, utils.Wrap(err, "")
}

func (s *Sheet) InsertSheet(temp *Sheet) error {
	return utils.Wrap(DB.MysqlDB.db.Create(temp).Error, "InsertSheet failed")
}
func (s *Sheet) UpdateSheet(sheet *Sheet) error {
	t := DB.MysqlDB.db.Updates(sheet)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateSheet failed")
}

func (s *Sheet) GetAllSheetsInfo() ([]*Sheet, error) {
	var sheetList []Sheet
	err := DB.MysqlDB.db.Find(&sheetList).Error
	var transfer []*Sheet
	for _, v := range sheetList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetAllSheetsInfo failed ")
}
