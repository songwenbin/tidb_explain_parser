package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("test")
	db, err := sql.Open("mysql", "root:root@/test?charset=utf8")
	checkErr(err)
	defer db.Close()

	fmt.Println(db)

	rows, err := db.Query("SELECT * FROM t1")
	checkErr(err)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		checkErr(err)
		fmt.Println(id)
	}

}
