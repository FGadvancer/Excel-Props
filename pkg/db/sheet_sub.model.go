package db

import (
	"Excel-Props/pkg/utils"
	"gorm.io/gorm"
)

//生成模号表
type SheetSub struct {
	SubSheetID string   `gorm:"column:sub_sheet_id;primary_key;type:char(64)" json:"subSheetID"`
	Ex         string   `gorm:"column:ex;type:varchar(1024)"  json:"ex,omitempty"`
	DB         *gorm.DB `gorm:"-" json:"-"`
}

func NewSheetSub(DB *gorm.DB) *SheetSub {
	return &SheetSub{DB: DB}
}
func (s *SheetSub) BatchInsertSheetSubList(sheetSubList []*SheetSub) error {
	if sheetSubList == nil {
		return nil
	}
	return utils.Wrap(DB.MysqlDB.db.Create(sheetSubList).Error, "BatchInsertSheetSubList failed")
}

//func (s *SheetSub)DeleteSheetAndMaterialInfoBySheetIDAndVersion(subSheetID string) error {
//	return DB.MysqlDB.db.Where("sheet_id=? and version=? ", sheetID, version).Delete(&SheetAndMaterial{}).Error
//}
func (s *SheetSub) GetSheetSubList() ([]string, error) {
	var temp []string
	return temp, utils.Wrap(DB.MysqlDB.db.Model(&SheetSub{}).Select("sub_sheet_id").Find(&temp).Error, "")
}

func (s *SheetSub) GetSheetSubInfo(subSheetID string) (*SheetSub, error) {
	temp := SheetSub{}
	err := DB.MysqlDB.db.Model(&temp).Where("sub_sheet_id = ? ", subSheetID).Take(&temp).Error
	return &temp, utils.Wrap(err, "")
}
