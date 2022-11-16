package db

import (
	"Excel-Props/pkg/utils"
	"gorm.io/gorm"
	"time"
)

//用户注册表
type Register struct {
	Account    string    `gorm:"column:account;primary_key;type:char(64)" json:"account"`
	Password   string    `gorm:"column:password;type:varchar(255)" json:"password"`
	UserName   string    `gorm:"column:user_name;type:varchar(64)" json:"userName"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
	ManagerLevel int     `gorm:"column:manager_level" json:"managerLevel"`
	Ex         string    `gorm:"column:ex;type:varchar(1024)"  json:"ex"`
	DB         *gorm.DB  `gorm:"-" json:"-"`
}

func NewRegister(DB *gorm.DB) *Register {
	return &Register{DB: DB}
}
func (r *Register) GetAccountInfo(account string) (*Register, error) {
	user := Register{}
	err := DB.MysqlDB.db.Model(&user).Where("account = ? ", account).Take(&user).Error
	return &user, utils.Wrap(err, "")
}
