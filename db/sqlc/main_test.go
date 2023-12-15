package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/18941325888/sb/util"
	_ "github.com/lib/pq"
	//导入驱动，但是没有在代码中调佣任何函数。 需要引入空白标识符
)

var testQueries *Queries
var testDB *sql.DB

//全局变量，包含一个DBTX，可以数据库连接和处理事务

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	//创建数据库连接，传入数据库驱动程序和数据库源字符串
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	//使用连接创建新的testQueries对象，New定义在sqlc为我们生成的db.go文件中

	os.Exit(m.Run())
	//m.Run()运行单元测试，此函数将返回一个退出代码，它告诉我们测试是通过还是失败。然后我们通过os.Exit命令将它报告回测试运行器。
}
