package db

import (
	"Excel-Props/pkg/log"
	"Excel-Props/pkg/utils"
	"gorm.io/gorm"
)

//料件模板表
type TemplateMaterial struct {
	IndexNumber        int32    `gorm:"column:index_number" json:"indexNumber"`
	MaterialKey        string   `gorm:"column:material_key;primary_key;type:char(64)" json:"materialKey"`
	MaterialStandard   string   `gorm:"column:material_standard;type:varchar(64)" json:"materialStandard"`
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

func NewTemplate2(DB *gorm.DB) *TemplateMaterial {
	return &TemplateMaterial{DB: DB}
}

func (t *TemplateMaterial) ImportDataToTemplateMaterial(data []*TemplateMaterial) error {
	for _, v := range data {
		t := TemplateMaterial{}
		if e := DB.MysqlDB.db.Model(&TemplateMaterial{}).Where("material_key = ? And material_standard = ?", v.MaterialKey, v.MaterialStandard).Take(&t).Error; e != nil {
			log.Debug("new Material find : ", v, e.Error())
			if err := DB.MysqlDB.db.Model(v).Create(v).Error; err != nil {
				log.Error("import sheet  db error  : ", v, e.Error())
			}
		}
	}
	return nil
}
func (t *TemplateMaterial) GetMaterialInfo(materialKey string) (*TemplateMaterial, error) {
	temp := TemplateMaterial{}
	err := DB.MysqlDB.db.Model(&temp).Where("material_key = ?", materialKey).Take(&temp).Error
	return &temp, utils.Wrap(err, "")
}
func (t *TemplateMaterial) DeleteAllTemplateMaterial() error {
	err := DB.MysqlDB.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&TemplateMaterial{}).Error
	return utils.Wrap(err, "")
}
func (t *TemplateMaterial) GetAllMaterialTemplates() ([]*TemplateMaterial, error) {
	var templateList []TemplateMaterial
	err := utils.Wrap(DB.MysqlDB.db.Debug().Order("index_number asc").Find(&templateList).Error,
		"GetAllMaterialTemplates failed")
	var transfer []*TemplateMaterial
	for _, v := range templateList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}
