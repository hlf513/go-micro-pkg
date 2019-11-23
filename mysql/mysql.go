package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/util/log"

	"github.com/hlf513/go-micro-pkg/config"
)

// dbs 这里是只读
var dbs = make(map[string]*gorm.DB)

// Connect 连接 MySQL
func Connect(configs map[string]config.DB) {
	for name, dbCfg := range configs {
		db, err := gorm.Open(
			"mysql",
			fmt.Sprintf(
				"%s:%s@tcp(%s)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local",
				dbCfg.Username,
				dbCfg.Password,
				dbCfg.Host,
				dbCfg.DBName))
		if err != nil {
			log.Fatal("mysql.connect:", err.Error())
			return
		}

		// 可使用单数表格
		db.SingularTable(true)
		db.DB().SetMaxIdleConns(dbCfg.MaxIdleConn)
		db.DB().SetMaxOpenConns(dbCfg.MaxOpenConn)
		db.DB().SetConnMaxLifetime(dbCfg.MaxLifeTime * time.Second)
		AddGormCallbacks(db)
		dbs[name] = db

		if dbCfg.Debug {
			db.LogMode(true)
		}
		log.Info("初始化 MySQL 连接：" + name)
	}
}

// Close 关闭数据库连接
func Close() {
	for _, db := range dbs {
		if db != nil {
			_ = db.Close()
		}
	}
}

// GetDB 获取 DB 实例
func GetDB(ctx context.Context, name ...string) *gorm.DB {
	var dbName string
	if len(name) == 0 {
		dbName = "default"
	} else {
		dbName = name[0]
	}

	if db, ok := dbs[dbName]; ok {
		return SetSpanToGorm(ctx, db)
	} else {
		log.Fatal("未找到 DB[" + dbName + "] 实例")
	}
	return nil
}
