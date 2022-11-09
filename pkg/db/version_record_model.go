package db

import (
	"Excel-Props/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type VersionUpLoadRecord struct {
	SheetID            string    `gorm:"column:sheet_id;primary_key;type:char(64)" json:"sheetID"`
	Version            int32     `gorm:"column:version;primary_key;" json:"version"`
	SubVersion         int32     `gorm:"column:sub_version;primary_key;" json:"subVersion"`
	MaterialKey        string    `gorm:"column:material_key;primary_key;type:char(64)" json:"materialKey"`
	MaterialStandard   string    `gorm:"column:material_standard;primary_key;type:char(64)" json:"materialStandard"`
	MaterialCategory   string    `gorm:"column:material_category;type:varchar(64)" json:"materialCategory"`
	MaterialName       string    `gorm:"column:material_name;type:varchar(64)" json:"materialName"`
	MaterialSubstance  string    `gorm:"column:material_substance;type:varchar(64)" json:"materialSubstance"`
	Quantity           int32     `gorm:"column:quantity" json:"quantity"`
	MaterialUnit       string    `gorm:"column:material_unit;type:varchar(64)" json:"materialUnit"`
	ProcessingCategory string    `gorm:"column:processing_category;type:varchar(64)" json:"processingCategory"`
	RemarkOne          string    `gorm:"column:remark_one;type:varchar(64)" json:"remarkOne"`
	RemarkTwo          string    `gorm:"column:remark_two;type:varchar(64)" json:"remarkTwo"`
	IsPurchase         string    `gorm:"column:is_purchase;type:varchar(64)" json:"isPurchase"`
	StandardCraft      string    `gorm:"column:standard_craft;type:varchar(64)" json:"standardCraft"`
	SubMaterialKey     string    `gorm:"column:sub_material_key;type:varchar(1024)" json:"subMaterialKey"`
	LastModifyTime     time.Time `gorm:"column:last_modify_time;index:index_last_modify_time;" json:"lastModifyTime"`
	LastModifierUserID string    `gorm:"column:last_modifier_userID;char(64)" json:"lastModifierUserID"`
	LastModifierName   string    `gorm:"column:last_modifier_name;type:varchar(64)" json:"lastModifierName"`
	LastModifyCount    int32     `gorm:"column:last_modify_count" json:"lastModifyCount"`
	Ex                 string    `gorm:"column:ex;type:varchar(1024)"  json:"ex,omitempty"`
	DB                 *gorm.DB  `gorm:"-" json:"-"`
}

func NewVersionUpLoadRecord(DB *gorm.DB) *VersionUpLoadRecord {
	return &VersionUpLoadRecord{DB: DB}
}

func (s *VersionUpLoadRecord) BatchInsertVersionUpLoadRecordList(recordList []*VersionUpLoadRecord) error {
	if recordList == nil {
		return nil
	}
	return utils.Wrap(DB.MysqlDB.db.Create(recordList).Error, "BatchInsertVersionUpLoadRecordList failed")
}
