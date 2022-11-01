package db

import (
	"Excel-Props/pkg/utils"
	"fmt"
	"gorm.io/gorm"
)

//模号模板表
type Template1 struct {
	SheetID     string   `gorm:"column:sheet_id;primary_key;type:char(64)" json:"sheetID"`
	MachineKind string   `gorm:"column:machine_kind;type:varchar(64)" json:"machineKind"`
	ProductName string   `gorm:"column:product_name;type:varchar(64)" json:"productName"`
	Code        string   `gorm:"column:code;type:varchar(64)" json:"code"`
	Ex          string   `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
	DB          *gorm.DB `gorm:"-" json:"-"`
}

func NewTemplate1(DB *gorm.DB) *Template1 {
	return &Template1{DB: DB}
}

func (t *Template1) ImportDataToTemplate1(data []*Template1) error {
	for _, v := range data {
		t := Template1{}
		if e := DB.MysqlDB.db.Model(&Template1{}).Where("sheet_id = ? ", v.SheetID).Take(&t).Error; e != nil {
			fmt.Println("new sheetID find : ", v, e.Error())
			if err := DB.MysqlDB.db.Model(v).Create(v).Error; err != nil {
				fmt.Println("import sheet  db error  : ", v, e.Error())
			}
		}
	}
	return nil
}

func (t *Template1) GetTemplateInfo(sheetID string) (*Template1, error) {
	temp := Template1{}
	err := DB.MysqlDB.db.Model(&temp).Where("sheet_id = ? ", sheetID).Take(&temp).Error
	return &temp, utils.Wrap(err, "")
}

func (t *Template1) GetAllTemplates() ([]*Template1, error) {
	var templateList []Template1
	err := utils.Wrap(DB.MysqlDB.db.Find(&templateList).Error,
		"GetAllTemplates failed")
	var transfer []*Template1
	for _, v := range templateList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, err
}
