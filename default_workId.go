package idGenerator

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
)

type DBConfig struct {
	ConnectInfo string `json:"connect_info"`
}

func GetWorkId(configPath string) uint64 {
	file, e := ioutil.ReadFile(configPath)
	if e != nil {
		log.Panicf("open file err %v", e)
	}

	config := new(DBConfig)
	e = json.Unmarshal(file, config)
	if e != nil {
		log.Panicf("unable to marshal %v", e)
	}

	db, e := sql.Open("mysql", config.ConnectInfo)
	if e != nil {
		log.Panicf("unable to connect db %v", e)
	}

	stmt, e := db.Prepare("insert into worker_node (host_name, port) values (?, ?)")
	if e != nil {
		log.Panicf("prepare stmt err, %v", stmt)
	}

	result, e := stmt.Exec(GetIp(), rand.Int31())
	if e != nil {
		log.Panicf("execute sql err, %v", e)
	}
	id, e := result.LastInsertId()
	if e != nil {
		log.Panicf("get last insert id err, %v", e)
	}
	return uint64(id)
}

func GetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Panicf("unable to get addrs, %v", err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}
