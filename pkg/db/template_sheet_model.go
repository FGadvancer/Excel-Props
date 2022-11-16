package db

import (
	"Excel-Props/pkg/log"
	"Excel-Props/pkg/utils"
	"gorm.io/gorm"
)

//模号模板表
type TemplateSheet struct {
	SheetID     string   `gorm:"column:sheet_id;primary_key;type:char(64)" json:"sheetID"`
	MachineKind string   `gorm:"column:machine_kind;type:varchar(64)" json:"machineKind"`
	ProductName string   `gorm:"column:product_name;type:varchar(64)" json:"productName"`
	Code        string   `gorm:"column:code;type:varchar(64)" json:"code"`
	Ex          string   `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
	DB          *gorm.DB `gorm:"-" json:"-"`
}

func NewTemplate1(DB *gorm.DB) *TemplateSheet {
	return &TemplateSheet{DB: DB}
}

func (t *TemplateSheet) ImportDataToTemplateSheet(data []*TemplateSheet) error {
	for _, v := range data {
		t := TemplateSheet{}
		if e := DB.MysqlDB.db.Model(&TemplateSheet{}).Where("sheet_id = ? ", v.SheetID).Take(&t).Error; e != nil {
			log.Error("new sheetID find : ", v, e.Error())
			if err := DB.MysqlDB.db.Model(v).Create(v).Error; err != nil {
				log.Error("import sheet  db error  : ", v, e.Error())
			}
		}
	}
	return nil
}
func (t *TemplateSheet) DeleteAllTemplateSheet()  error{
	err := DB.MysqlDB.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&TemplateSheet{}).Error
	return utils.Wrap(err, "")
}
func (t *TemplateSheet) GetTemplateSheetInfo(sheetID string) (*TemplateSheet, error) {
	temp := TemplateSheet{}
	err := DB.MysqlDB.db.Model(&temp).Where("sheet_id = ? ", sheetID).Take(&temp).Error
	return &temp, utils.Wrap(err, "")
}

func (t *TemplateSheet) GetAllSheetTemplates() ([]*TemplateSheet, error) {
	var templateList []TemplateSheet
	err := utils.Wrap(DB.MysqlDB.db.Find(&templateList).Error,
		"GetAllSheetTemplates failed")
	var transfer []*TemplateSheet
	for _, v := range templateList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}
