package database

import (
	"fmt"
	"os"

	"github.com/mike504110403/goutils/dbconn"

	mlog "github.com/mike504110403/goutils/log"
)

// 控管所有連線，其他模組僅組裝連線字串與進行連線等作業

type Config struct {
	DefaultPostgresConnConfig dbconn.Config
}

func Init() {
	dbNew([]dbconn.DBName{POSRECORDS})
}

func dbNew(dbnames []dbconn.DBName) {
	for _, dbname := range dbnames {
		if err := setDsn(dbname); err != nil {
			mlog.Fatal(err.Error())
		}
		dbconn.New(map[dbconn.DBName]dbconn.DBConfig{
			dbname: {
				DBDriver:  dbconn.DBDriverMySQL,
				DSNSource: os.Getenv(Envkey),
			},
		})
	}
}

func setDsn(dbName dbconn.DBName) error {
	// 從環境變數中讀取相關設定
	username := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	setting := "parseTime=true&loc=Asia%2FTaipei"
	// 構建 CRYPTO_MYSQL_URL
	cryptoMySQLURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", username, password, host, DB_Name_Map[dbName], setting)
	// 寫入環境變數
	err := os.Setenv(string(Envkey), cryptoMySQLURL)
	if err != nil {
		return err
	}
	return nil
}
