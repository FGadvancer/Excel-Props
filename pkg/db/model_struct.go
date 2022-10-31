package db

import "time"
//用户注册表
type Register struct {
	Account    string    `gorm:"column:account;primary_key;type:char(64)" json:"account"`
	Password   string    `gorm:"column:password;type:varchar(255)" json:"password"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	Ex         string    `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
}

//生成模号表
type Sheet struct {
	SheetID            string    `gorm:"column:sheet_id;primary_key;type:char(64)" json:"sheetID"`
	CommodityName      string    `gorm:"column:commodity_name;type:char(64)" json:"commodityName"`
	Version            int32     `gorm:"column:version" json:"version"`
	Code               int32     `gorm:"column:code" json:"code"`
	CreatorUserID      string    `gorm:"column:creator_user_id;size:64"`
	CreateTime         time.Time `gorm:"column:create_time" json:"createTime"`
	LastModifierUserID string    `gorm:"column:last_modifier_userID;size:64" json:"lastModifierUserID"`
	LastModifierIP     string    `gorm:"column:last_modifier_ip;size:64" json:"LastModifierIP"`
	LastModifyTime     time.Time `gorm:"column:last_modify_time" json:"lastModifyTime"`
	Ex                 string    `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
}

//模号对应料件表
type SheetAndMaterial struct {
	SheetID            string    `gorm:"column:sheet_id;primary_key;type:char(64)" json:"sheetID"`
	MaterialKey        string `gorm:"column:material_key;primary_key;type:char(64)" json:"materialKey"`
	MaterialCategory   string `gorm:"column:material_category;type:varchar(64)" json:"materialCategory"`
	MaterialName       string `gorm:"column:material_name;type:varchar(64)" json:"materialName"`
	MaterialSubstance  string `gorm:"column:material_substance;type:varchar(64)" json:"materialSubstance"`
	MaterialStandard   string `gorm:"column:material_standard;type:varchar(64)" json:"materialStandard"`
	Quantity           int32  `gorm:"column:quantity" json:"quantity"`
	MaterialUnit       string `gorm:"column:material_unit;type:varchar(64)" json:"materialUnit"`
	ProcessingCategory string `gorm:"column:processing_category;type:varchar(64)" json:"processingCategory"`
	RemarkOne          string `gorm:"column:remark_one;type:varchar(64)" json:"remarkOne"`
	RemarkTwo          string `gorm:"column:remark_two;type:varchar(64)" json:"remarkTwo"`
	IsPurchase         string   `gorm:"column:is_purchase;type:varchar(64)" json:"isPurchase"`
	StandardCraft      string `gorm:"column:standard_craft;type:varchar(64)" json:"standardCraft"`
	SubMaterialKey     string    `gorm:"column:sub_material_key;primary_key;type:varchar(1024)" json:"subMaterialKey"`
	LastModifyTime     time.Time `gorm:"column:last_modify_time" json:"lastModifyTime"`
	LastModifierUserID string    `gorm:"column:last_modifier_userID;size:64" json:"lastModifierUserID"`
	LastModifyCount int32    `gorm:"column:last_modify_count" json:"lastModifyCount"`

	Ex string `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
}
//模号模板表
type Template1 struct {
	SheetID     string `gorm:"column:sheet_id;primary_key;type:char(64)" json:"sheetID"`
	MachineKind string `gorm:"column:machine_kind;type:varchar(64)" json:"machineKind"`
	ProductName string `gorm:"column:product_name;type:varchar(64)" json:"productName"`
	Code        string  `gorm:"column:code;type:varchar(64)" json:"code"`
	Ex          string `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
}
//料件模板表
type Template2 struct {
	MaterialKey        string `gorm:"column:material_key;primary_key;type:char(64)" json:"materialKey"`
	MaterialCategory   string `gorm:"column:material_category;type:varchar(64)" json:"materialCategory"`
	MaterialName       string `gorm:"column:material_name;type:varchar(64)" json:"materialName"`
	MaterialSubstance  string `gorm:"column:material_substance;type:varchar(64)" json:"materialSubstance"`
	MaterialStandard   string `gorm:"column:material_standard;type:varchar(64)" json:"materialStandard"`
	Quantity           int32  `gorm:"column:quantity" json:"quantity"`
	MaterialUnit       string `gorm:"column:material_unit;type:varchar(64)" json:"materialUnit"`
	ProcessingCategory string `gorm:"column:processing_category;type:varchar(64)" json:"processingCategory"`
	RemarkOne          string `gorm:"column:remark_one;type:varchar(64)" json:"remarkOne"`
	RemarkTwo          string `gorm:"column:remark_two;type:varchar(64)" json:"remarkTwo"`
	IsPurchase         string  `gorm:"column:is_purchase;type:varchar(64)" json:"isPurchase"`
	StandardCraft      string `gorm:"column:standard_craft;type:varchar(64)" json:"standardCraft"`
	Ex                 string `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
}
