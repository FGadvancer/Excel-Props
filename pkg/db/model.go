package db

var DB DataBases

type DataBases struct {
	MysqlDB mysqlDB
}

func key(dbAddress, dbName string) string {
	return dbAddress + "_" + dbName
}

func init() {
	//mysql init
	initMysqlDB()

}
