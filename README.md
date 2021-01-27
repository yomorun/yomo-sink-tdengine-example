![yomo-integrate-with-TDengine]()

# yomo-source-TDengine-starter

TDengine 🙌 YoMo

## About TDengine

EMQ X broker is a fully open source, highly scalable, highly available distributed MQTT messaging broker for IoT, M2M and Mobile applications that can handle tens of millions of concurrent clients.

Starting from 3.0 release, EMQ X broker fully supports MQTT V5.0 protocol specifications and backward compatible with MQTT V3.1 and V3.1.1, as well as other communication protocols such as MQTT-SN, CoAP, LwM2M, WebSocket and STOMP. The 3.0 release of the EMQ X broker can scaled to 10+ million concurrent MQTT connections on one cluster.

For more information, please visit [EMQ X homepage](https://www.emqx.io/)

## 1/x Installing TDengine

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

## 2/x Setup Database

```bash
$ taos -c test/cfg
taos> create database noise;
Query OK, 0 row(s) affected (0.004701s)

taos> use noise;
Database changed.

taos> create table t (ts timestamp, noise float);
Query OK, 0 row(s) affected (0.011501s)

taos> insert into t values ('2021-01-01 00:00:00', 41.1);
Query OK, 1 row(s) affected (0.002012s)

taos> select * from t;
           ts            |        noise         |
=================================================
 2021-01-01 00:00:00.000 |             41.10000 |
Query OK, 1 row(s) in set (0.004414s)

taos>
```

## 3/x Integrate TDengine as YoMo-Sink
