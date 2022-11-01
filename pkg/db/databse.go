package db

var DB DataBases

type DataBases struct {
	MysqlDB mysqlDB
	Redis   *Redis
}

func init() {
	//mysql init
	initMysqlDB()
	DB.Redis = initRedis()

}
