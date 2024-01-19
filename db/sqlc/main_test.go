package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/18941325888/sb/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

//全局变量，包含一个DBTX，可以数据库连接和处理事务

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	//创建数据库连接，传入数据库驱动程序和数据库源字符串
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
	//m.Run()运行单元测试，此函数将返回一个退出代码，它告诉我们测试是通过还是失败。然后我们通过os.Exit命令将它报告回测试运行器。
}
