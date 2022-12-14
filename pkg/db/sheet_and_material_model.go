package db

import (
	"Excel-Props/pkg/utils"
	"errors"
	"gorm.io/gorm"
	"time"
)

//模号对应料件表
type SheetAndMaterial struct {
	SheetID           string `gorm:"column:sheet_id;primary_key;type:char(64)" json:"sheetID"`
	MaterialKey       string `gorm:"column:material_key;primary_key;type:char(64)" json:"materialKey"`
	MaterialStandard  string `gorm:"column:material_standard;primary_key;type:char(64)" json:"materialStandard"`
	Version           int32  `gorm:"column:version;primary_key;" json:"version"`
	IndexNumber       int32  `gorm:"column:index_number" json:"indexNumber"`
	MaterialCategory  string `gorm:"column:material_category;type:varchar(64)" json:"materialCategory"`
	MaterialName      string `gorm:"column:material_name;type:varchar(64)" json:"materialName"`
	MaterialSubstance string `gorm:"column:material_substance;type:varchar(64)" json:"materialSubstance"`

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

func NewSheetAndMaterial(DB *gorm.DB) *SheetAndMaterial {
	return &SheetAndMaterial{DB: DB}
}

func (s *SheetAndMaterial) GetSheetAndMaterialInfo(sheetID string, materialKey string, materialStandard string, version int32) (*SheetAndMaterial, error) {
	temp := SheetAndMaterial{}
	err := DB.MysqlDB.db.Model(&temp).Where("sheet_id = ? And material_key = ? And material_standard = ? And version = ?", sheetID, materialKey, materialStandard, version).Take(&temp).Error
	return &temp, utils.Wrap(err, "")
}
func (s *SheetAndMaterial) BatchInsertSheetAndMaterialList(materialList []*SheetAndMaterial) error {
	if materialList == nil {
		return nil
	}
	return utils.Wrap(DB.MysqlDB.db.Create(materialList).Error, "BatchInsertSheetAndMaterialList failed")
}

func (s *SheetAndMaterial) InsertSheetAndMaterial(temp *SheetAndMaterial) error {
	return utils.Wrap(DB.MysqlDB.db.Create(temp).Error, "InsertSheetAndMaterial failed")
}
func (s *SheetAndMaterial) UpdateSheetAndMaterial(material *SheetAndMaterial) error {
	t := DB.MysqlDB.db.Updates(material)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateSheetAndMaterial failed")
}

func (s *SheetAndMaterial) GetSheetAndMaterialInfoBySheetID(sheetID string) ([]*SheetAndMaterial, error) {
	var temp []SheetAndMaterial
	err := DB.MysqlDB.db.Debug().Model(&temp).Where("sheet_id = ?", sheetID).Order("last_modify_time DESC").Find(&temp).Error

	var transfer []*SheetAndMaterial
	for _, v := range temp {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetSheetAndMaterialInfoBySheetID failed")
}
func (s *SheetAndMaterial) DeleteSheetAndMaterialInfoBySheetIDAndVersion(sheetID string, version int32) error {
	return DB.MysqlDB.db.Where("sheet_id=? and version=? ", sheetID, version).Delete(&SheetAndMaterial{}).Error
}

func (s *SheetAndMaterial) DecrMaterialQuantity(sheetID string, version int32, materialKey string, materialStandard string, quantity int32) error {
	c := SheetAndMaterial{SheetID: sheetID, Version: version, MaterialKey: materialKey, MaterialStandard: materialStandard}
	return DB.MysqlDB.db.Model(&c).Update("quantity", gorm.Expr("quantity-?", quantity)).Error
}
