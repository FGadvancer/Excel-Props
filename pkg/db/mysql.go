package db

import (
	"Excel-Props/pkg/config"
	"Excel-Props/pkg/utils"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mysqlDB struct {
	sync.RWMutex
	db *gorm.DB
	*TemplateSheet
	*TemplateMaterial
	*Register
	*Sheet
	*SheetAndMaterial
	*VersionUploadRecord
	*SheetSub
}

func (m mysqlDB) Db() *gorm.DB {
	return m.db
}

type Writer struct{}

func (w Writer) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func initMysqlDB() {
	fmt.Println("init mysqlDB start")
	//When there is no open IM database, connect to the mysql built-in database to create openIM database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], "mysql")
	var db *gorm.DB
	var err1 error
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		fmt.Println("Open failed ", err.Error(), dsn)
	}
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		db, err1 = gorm.Open(mysql.Open(dsn), nil)
		if err1 != nil {
			fmt.Println("Open failed ", err1.Error(), dsn)
			panic(err1.Error())
		}
	}

	//Check the database and table during initialization
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci;", config.Config.Mysql.DBDatabaseName)
	fmt.Println("exec sql: ", sql, " begin")
	err = db.Exec(sql).Error
	if err != nil {
		fmt.Println("Exec failed ", err.Error(), sql)
		panic(err.Error())
	}
	fmt.Println("exec sql: ", sql, " end")
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.DBUserName, config.Config.Mysql.DBPassword, config.Config.Mysql.DBAddress[0], config.Config.Mysql.DBDatabaseName)

	newLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold:             time.Duration(config.Config.Mysql.SlowThreshold) * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.LogLevel(config.Config.Mysql.LogLevel),                       // Log level
			IgnoreRecordNotFoundError: true,                                                                // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                                                // Disable color
		},
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("Open failed ", err.Error(), dsn)
		panic(err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}

	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Config.Mysql.DBMaxLifeTime))
	sqlDB.SetMaxOpenConns(config.Config.Mysql.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Config.Mysql.DBMaxIdleConns)

	fmt.Println("open mysql ok ", dsn)
	db.AutoMigrate(
		&Register{}, &Sheet{}, &SheetAndMaterial{}, &TemplateSheet{}, &TemplateMaterial{}, &VersionUploadRecord{})
	db.Set("gorm:table_options", "CHARSET=utf8")
	db.Set("gorm:table_options", "collation=utf8_unicode_ci")

	if !db.Migrator().HasTable(&Register{}) {
		fmt.Println("CreateTable Register")
		db.Migrator().CreateTable(&Register{})
	}
	if !db.Migrator().HasTable(&Sheet{}) {
		fmt.Println("CreateTable Sheet")
		db.Migrator().CreateTable(&Sheet{})
	}
	if !db.Migrator().HasTable(&SheetAndMaterial{}) {

		err := db.Migrator().CreateTable(&SheetAndMaterial{})
		if err != nil {
			fmt.Println("CreateTable SheetAndMaterial err:", err.Error())
		}
	}
	if !db.Migrator().HasTable(&TemplateSheet{}) {
		fmt.Println("CreateTable TemplateSheet")
		db.Migrator().CreateTable(&TemplateSheet{})
	}
	if !db.Migrator().HasTable(&TemplateMaterial{}) {
		fmt.Println("CreateTable TemplateMaterial")
		db.Migrator().CreateTable(&TemplateMaterial{})
	}
	if !db.Migrator().HasTable(&VersionUploadRecord{}) {
		fmt.Println("CreateTable VersionUploadRecord")
		db.Migrator().CreateTable(&VersionUploadRecord{})
	}
	if !db.Migrator().HasTable(&SheetSub{}) {
		fmt.Println("CreateTable SheetSub")
		db.Migrator().CreateTable(&SheetSub{})
	}
	DB.MysqlDB.db = db
	DB.MysqlDB.Register = NewRegister(db)
	DB.MysqlDB.TemplateSheet = NewTemplate1(db)
	DB.MysqlDB.TemplateMaterial = NewTemplate2(db)
	DB.MysqlDB.Sheet = NewSheet(db)
	DB.MysqlDB.SheetAndMaterial = NewSheetAndMaterial(db)
	DB.MysqlDB.VersionUploadRecord = NewVersionUpLoadRecord(db)
	DB.MysqlDB.SheetSub = NewSheetSub(db)
	for i, v := range config.Config.Manager.AppManagerUid {
		user := Register{}
		if e := DB.MysqlDB.db.Model(&Register{}).Where("account = ? ", v).Take(&user).Error; e != nil {
			fmt.Println("admin find : ", v, e.Error())
			user.Account = config.Config.Manager.AppManagerUid[i]
			user.Password = utils.Md5(config.Config.Manager.Secrets[i])
			user.UserName = config.Config.Manager.AppManagerUid[i]
			if v == "admin1" {
				user.ManagerLevel = 1
			}
			user.CreateTime = time.Now()
			if err := DB.MysqlDB.db.Model(&user).Create(&user).Error; err != nil {
				fmt.Println("init db error  : ", v, e.Error())
			}
		}
	}

	return
}
