package db

import (
	"gorm.io/gorm"
	"time"
)

//模号对应料件表
type SheetAndMaterial struct {
	SheetID            string    `gorm:"column:sheet_id;primary_key;type:char(32)" json:"sheetID"`
	MaterialKey        string    `gorm:"column:material_key;primary_key;type:char(32)" json:"materialKey"`
	MaterialCategory   string    `gorm:"column:material_category;type:varchar(64)" json:"materialCategory"`
	MaterialName       string    `gorm:"column:material_name;type:varchar(64)" json:"materialName"`
	MaterialSubstance  string    `gorm:"column:material_substance;type:varchar(64)" json:"materialSubstance"`
	MaterialStandard   string    `gorm:"column:material_standard;type:varchar(64)" json:"materialStandard"`
	Quantity           int32     `gorm:"column:quantity" json:"quantity"`
	MaterialUnit       string    `gorm:"column:material_unit;type:varchar(64)" json:"materialUnit"`
	ProcessingCategory string    `gorm:"column:processing_category;type:varchar(64)" json:"processingCategory"`
	RemarkOne          string    `gorm:"column:remark_one;type:varchar(64)" json:"remarkOne"`
	RemarkTwo          string    `gorm:"column:remark_two;type:varchar(64)" json:"remarkTwo"`
	IsPurchase         string    `gorm:"column:is_purchase;type:varchar(64)" json:"isPurchase"`
	StandardCraft      string    `gorm:"column:standard_craft;type:varchar(64)" json:"standardCraft"`
	SubMaterialKey     string    `gorm:"column:sub_material_key;type:varchar(1024)" json:"subMaterialKey"`
	LastModifyTime     time.Time `gorm:"column:last_modify_time" json:"lastModifyTime"`
	LastModifierUserID string    `gorm:"column:last_modifier_userID;char(64)" json:"lastModifierUserID"`
	LastModifyCount    int32     `gorm:"column:last_modify_count" json:"lastModifyCount"`
	Ex                 string    `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
	DB                 *gorm.DB  `gorm:"-" json:"-"`
}

func NewSheetAndMaterial(DB *gorm.DB) *SheetAndMaterial {
	return &SheetAndMaterial{DB: DB}
}
