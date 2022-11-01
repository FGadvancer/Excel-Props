package db

var DB DataBases

type DataBases struct {
	MysqlDB mysqlDB
	redis   *Redis
}

func init() {
	//mysql init
	initMysqlDB()
	DB.redis = initRedis()

}
