package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-filestore-server/config"
	"go-filestore-server/util"
	"log"
	"os"
)

var db *sql.DB

func InitMysql() {
	source := util.GetMysqlSource(config.DefaultConfig.MysqlUser,
		config.DefaultConfig.MysqlPwd,
		config.DefaultConfig.MysqlHost,
		config.DefaultConfig.MysqlPort,
		config.DefaultConfig.MysqlDb,
		config.DefaultConfig.MysqlCharset)
	fmt.Println("mysql source:\t", source)
	db, _ = sql.Open("mysql", source)
	db.SetMaxOpenConns(config.DefaultConfig.MysqlMaxConn)
	err := db.Ping()
	if err != nil {
		fmt.Println("failed to connect to mysql, err:\t", err.Error())
		os.Exit(1)
	}
}

// DBConn 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		// 将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
