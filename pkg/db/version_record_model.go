package db

import (
	"Excel-Props/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type VersionUploadRecord struct {
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
	CommitTime         time.Time `gorm:"column:commit_time;index:index_commit_time;" json:"commitTime"`
	ModifierUserID     string    `gorm:"column:modifier_userID;char(64)" json:"modifierUserID"`
	ModifierName       string    `gorm:"column:modifier_name;type:varchar(64)" json:"modifierName"`
	Ex                 string    `gorm:"column:ex;type:varchar(1024)"  json:"ex,omitempty"`
	DB                 *gorm.DB  `gorm:"-" json:"-"`
}

func NewVersionUpLoadRecord(DB *gorm.DB) *VersionUploadRecord {
	return &VersionUploadRecord{DB: DB}
}

func (s *VersionUploadRecord) BatchInsertVersionUpLoadRecordList(recordList []*VersionUploadRecord) error {
	if recordList == nil {
		return nil
	}
	return utils.Wrap(DB.MysqlDB.db.Create(recordList).Error, "BatchInsertVersionUpLoadRecordList failed")
}
func (s *VersionUploadRecord) GetVersionRecordList(sheetID string, version int32) ([]*VersionUploadRecord, error) {
	var temp []VersionUploadRecord
	err := DB.MysqlDB.db.Debug().Model(&temp).Where("sheet_id = ? And version = ?", sheetID, version).Order("commit_time DESC").Find(&temp).Error

	var transfer []*VersionUploadRecord
	for _, v := range temp {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetVersionRecordList failed")
}
func (s *VersionUploadRecord) GetSubVersionRecordList(sheetID string, version int32, subVersion int32) ([]*VersionUploadRecord, error) {
	var temp []VersionUploadRecord
	err := DB.MysqlDB.db.Debug().Model(&temp).Where("sheet_id = ? And version = ? And sub_version=?", sheetID, version, subVersion).Find(&temp).Error

	var transfer []*VersionUploadRecord
	for _, v := range temp {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetSubVersionRecordList failed")
}
func (s *VersionUploadRecord) DeleteVersionRecordListBySheetIDAndVersion(sheetID string, version int32) error {
	return DB.MysqlDB.db.Where("sheet_id=? and version=? ", sheetID, version).Delete(&VersionUploadRecord{}).Error
}
func (s *VersionUploadRecord) DeleteSubVersionRecordListBySheetIDAndVersion(sheetID string, version int32, subVersion int32) error {
	return DB.MysqlDB.db.Where("sheet_id=? and version=? and sub_version=?", sheetID, version, subVersion).Delete(&VersionUploadRecord{}).Error
}
