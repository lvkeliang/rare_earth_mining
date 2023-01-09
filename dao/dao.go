package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

func InitDB() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/rare_earth_mining?charset=utf8mb4&locLocal&parseTime=True")
	//db, err := sql.Open("mysql", "LvKeliang:lkl12190711@tcp(localhost:3306)/rare_earth_mining?charset=utf8mb4&locLocal&parseTime=True")
	fmt.Println("执行InitDB")
	if err != nil {
		fmt.Println("InitDB出错了")
		log.Fatalf("connect mysql error : %v", err)
	}

	DB = db

	fmt.Println(db.Ping())
}
