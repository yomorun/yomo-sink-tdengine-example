package main

import (
	"database/sql"
	"fmt"
	_ "github.com/taosdata/driver-go/taosSql"
	"os"
)

var url = "root:taosdata@/tcp(172.16.0.191:6030)/yomo"

func main() {
	db, err := sql.Open("taosSql", url)
	if err != nil {
		fmt.Printf("Open database error: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Ensure 'noise' table exists
	sql := "CREATE TABLE IF NOT EXISTS noise (ts TIMESTAMP, v FLOAT)"
	res, err := db.Exec(sql)
	if err != nil {
		fmt.Printf("db.Exec error: %s\n", err)
	}
	fmt.Printf("res=%v\n", res)

	// Insert data
	var val float32 = 5.83
	sql = "INSERT INTO noise VALUES (NOW, " + fmt.Sprintf("%f", val) + ")"
	res, err = db.Exec(sql)
	if err != nil {
		fmt.Printf("Insert error: %s\n", err)
	}
	fmt.Printf("res=%v\n", res)

}
