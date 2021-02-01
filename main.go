package main

import (
	"database/sql"
	"fmt"
	"os"
	"log"
    "context"

	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/quic"
	_ "github.com/taosdata/driver-go/taosSql"
)

var url = "root:taosdata@/tcp(172.16.0.191:6030)/yomo"

func main() {
	log.Print("Starting YoMo Sink server: -> TDEngine")
	srv := quic.NewServer(&srvHandler{})
	err := srv.ListenAndServe(context.Background(), "0.0.0.0:9333")
	if err != nil {
		log.Printf("YoMo Sink server start failed: %s\n", err.Error())
	}
	select {}
}

type srvHandler struct {}

func (s *srvHandler) Listen() error {
	return nil
}

func (s *srvHandler) Read(qs quic.Stream) error {
	ch := y3.FromStream(qs).
		Subscribe(0x10).
		OnObserve(decode)

	go func() {
		for item := range ch {
			err := store(item)
			if err != nil {
				log.Printf("write to TDEngine error : %s", err.Error())
			} else {
				log.Printf("saved: %v", item)
			}
		}
	}()

	return nil
}

func decode(v []byte) (interface{}, error) {
	data, err := y3.ToFloat32(v)
	if err != nil {
		log.Printf("err: %s\n", err.Error())
	}
	return data, err
}

func store(v interface{}) error {
	db, err := sql.Open("taosSql", url)
	if err != nil {
		fmt.Printf("Open database error: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Ensure 'noise' table exists
	sql := "CREATE TABLE IF NOT EXISTS noise (ts TIMESTAMP, v FLOAT)"
	_, err = db.Exec(sql)
	if err != nil {
		fmt.Printf("db.Exec error: %s\n", err)
	}

	// Insert data
	var val float32 = v.(float32)
	sql = "INSERT INTO noise VALUES (NOW, " + fmt.Sprintf("%f", val) + ")"
	_, err = db.Exec(sql)
	if err != nil {
		fmt.Printf("Insert error: %s\n", err)
	}

	return err
}
