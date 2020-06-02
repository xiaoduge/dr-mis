/**********************************************************
 * @Author       dcj
 * @Date         2020-05-22 10:22:15
 * @Description  初始化数据库接口
 * @Version      V0.0.1
 **********************************************************/

package data

import (
	"database/sql"
	"log"
)

var Db *sql.DB

func init() {
	var err error
	// Db, err = sql.Open("postgres", "host=111.229.167.91 port=5432 dbname=dcjmis user=dcj password=dcj sslmode=disable")
	Db, err = sql.Open("postgres", "host=47.102.123.193 port=5432 dbname=safekids user=safekids password=safekids sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return
}
