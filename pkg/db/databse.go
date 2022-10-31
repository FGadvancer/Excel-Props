package db

var DB DataBases

type DataBases struct {
	MysqlDB mysqlDB
	Rdb redis_model.DB
}

func key(dbAddress, dbName string) string {
	return dbAddress + "_" + dbName
}

func init() {
	//mysql init
	initMysqlDB()
	initRedis(&DB.Rdb)

}
