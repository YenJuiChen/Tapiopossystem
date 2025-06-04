package database

import (
	"github.com/mike504110403/goutils/dbconn"
)

// 以下為連線設定範例 需調整
// 連線字符定義在.env
const Envkey string = "MYSQL_URL"

// mysql使用的參數設定
const POSRECORDS dbconn.DBName = "PosRecords"

// 組裝用字串
const PosRecord_dsn string = "PosRecords"

var DB_Name_Map = map[dbconn.DBName]string{
	POSRECORDS: PosRecord_dsn,
}
