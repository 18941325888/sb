package main

import (
	"database/sql"
	"log"

	"github.com/18941325888/sb/api"
	db "github.com/18941325888/sb/db/sqlc"
	"github.com/18941325888/sb/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("不能连接数据库：", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("不能创建服务器：", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("不能启动服务器：", err)

	}
}
