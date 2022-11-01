package db

import "gorm.io/gorm"

//料件模板表
type Template2 struct {
	MaterialKey        string   `gorm:"column:material_key;primary_key;type:char(64)" json:"materialKey"`
	MaterialCategory   string   `gorm:"column:material_category;type:varchar(64)" json:"materialCategory"`
	MaterialName       string   `gorm:"column:material_name;type:varchar(64)" json:"materialName"`
	MaterialSubstance  string   `gorm:"column:material_substance;type:varchar(64)" json:"materialSubstance"`
	MaterialStandard   string   `gorm:"column:material_standard;type:varchar(64)" json:"materialStandard"`
	Quantity           int32    `gorm:"column:quantity" json:"quantity"`
	MaterialUnit       string   `gorm:"column:material_unit;type:varchar(64)" json:"materialUnit"`
	ProcessingCategory string   `gorm:"column:processing_category;type:varchar(64)" json:"processingCategory"`
	RemarkOne          string   `gorm:"column:remark_one;type:varchar(64)" json:"remarkOne"`
	RemarkTwo          string   `gorm:"column:remark_two;type:varchar(64)" json:"remarkTwo"`
	IsPurchase         string   `gorm:"column:is_purchase;type:varchar(64)" json:"isPurchase"`
	StandardCraft      string   `gorm:"column:standard_craft;type:varchar(64)" json:"standardCraft"`
	Ex                 string   `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
	DB                 *gorm.DB `gorm:"-" json:"-"`
}

func NewTemplate2(DB *gorm.DB) *Template2 {
	return &Template2{DB: DB}
}
