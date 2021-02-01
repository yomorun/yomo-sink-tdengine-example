![yomo-integrate-with-TDengine](taos-yomo.jpg)

# Integrate TDengine to YoMo

TDengine 🙌 YoMo. Demonstrates how to integrate TDengine to YoMo and store data to TDengine after stream processing.

## About TDengine

TDengine is an open-sourced big data platform under GNU AGPL v3.0, designed and optimized for the Internet of Things (IoT), Connected Cars, Industrial IoT, and IT Infrastructure and Application Monitoring. Besides the 10x faster time-series database, it provides caching, stream computing, message queuing and other functionalities to reduce the complexity and cost of development and operation.

For more information, please visit [TDengine homepage](https://www.taosdata.com)

## 1: Installing TDengine

```bash
$ sudo apt-get install -y gcc cmake build-essential git
$ git clone --depth 1 https://github.com/taosdata/TDengine.git
$ cd TDengine
$ git submodule update --init --recursive
$ mkdir debug && cd $_
$ cmake ..
$ cmake --build .
$ make install
```

[TDengine officai installation page](https://github.com/taosdata/TDengine#installing)

## 2: Create database and table

```bash
$ taos
taos> create database yomo;
Query OK, 0 row(s) affected (0.004701s)

taos> use yomo;
Database changed.

taos> create table in not exists noise (ts timestamp, v float);
Query OK, 0 row(s) affected (0.011501s)

taos> insert into noise values ('2021-01-01 00:00:00', 41.1);
Query OK, 1 row(s) affected (0.002012s)

taos> select * from t;
           ts            |           v          |
=================================================
 2021-01-01 00:00:00.000 |             41.10000 |
Query OK, 1 row(s) in set (0.004414s)

taos>
```

## 3: Integrate TDengine with YoMo

### Start YoMo-Zipper

Configure [YoMo-Zipper](https://yomo.run/zipper):

```yaml
name: YoMoZipper 
host: localhost
port: 9000
sinks:
  - name: TDEngine
    host: localhost
    port: 9333
```

Start this zipper will listen on `9000` port, send data streams directly to `9333`:

```
cd ./zipper && yomo wf run
```

### Store data to TDengine

```go
var url = "root:taosdata@/tcp(localhost:6030)/yomo"

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
```

Start this [YoMo-Sink](https://yomo.run/sink), will save data to TDengine wherever data arrives

```bash
go run main.go
```

### Emulate a data source for test

```bash
cd source && go run main.go
```

This will start a [YoMo-Source](https://yomo.run/source), demonstrates a random float every 100ms to YoMo-Zipper

## 4. Verify TDengine

```bash
taos> use yomo;
Database changed.

taos> select * from noise;
           ts            |          v           |
=================================================
 2021-02-01 02:11:54.581 |              5.83000 |
 2021-02-01 02:14:19.372 |              5.83000 |
 2021-02-01 04:35:12.875 |             44.58845 |
 2021-02-01 04:35:12.963 |            157.36317 |
 2021-02-01 04:35:13.062 |             16.95439 |
 2021-02-01 04:35:13.163 |            180.45207 |
 2021-02-01 04:35:13.263 |             96.63864 |
 2021-02-01 04:35:13.364 |            134.08540 |
 2021-02-01 04:35:13.464 |             59.86330 |
 2021-02-01 04:35:13.565 |            197.74881 |
 2021-02-01 04:35:13.666 |            171.70944 |
 2021-02-01 04:35:13.765 |             36.40285 |
 ```
