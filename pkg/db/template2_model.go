package db

import (
	"fmt"
	"gorm.io/gorm"
)

//料件模板表
type Template2 struct {
	MaterialKey        string   `gorm:"column:material_key;primary_key;type:char(64)" json:"materialKey"`
	MaterialStandard   string   `gorm:"column:material_standard;primary_key;type:varchar(64)" json:"materialStandard"`
	MaterialCategory   string   `gorm:"column:material_category;type:varchar(64)" json:"materialCategory"`
	MaterialName       string   `gorm:"column:material_name;type:varchar(64)" json:"materialName"`
	MaterialSubstance  string   `gorm:"column:material_substance;type:varchar(64)" json:"materialSubstance"`
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

func (t *Template2) ImportDataToTemplate2(data []*Template2) error {
	for _, v := range data {
		t := Template2{}
		if e := DB.MysqlDB.db.Model(&Template1{}).Where("material_key = ? And material_standard", v.MaterialKey, v.MaterialStandard).Take(&t).Error; e != nil {
			fmt.Println("new Material find : ", v, e.Error())
			if err := DB.MysqlDB.db.Model(v).Create(v).Error; err != nil {
				fmt.Println("import sheet  db error  : ", v, e.Error())
			}
		}
	}
	return nil
}
