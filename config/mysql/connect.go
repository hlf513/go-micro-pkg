package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/util/log"
)

// dbInstance 保存多 db 实例
var dbInstance = make(map[string]*gorm.DB)

// Connect 连接 MySQL
func Connect() {
	configs := GetConf()
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
			log.Fatal("[mysql connect] ", err.Error())
			return
		}

		// 可使用单数表格
		db.SingularTable(true)
		db.DB().SetMaxIdleConns(dbCfg.MaxIdleConn)
		db.DB().SetMaxOpenConns(dbCfg.MaxOpenConn)
		db.DB().SetConnMaxLifetime(dbCfg.MaxLifeTime * time.Second)
		// 设置 trance 钩子
		AddGormCallbacks(db)
		// 保存 db 实例
		dbInstance[name] = db

		if dbCfg.Debug {
			db.LogMode(true)
		}
		log.Info("[mysql connect] 初始化 MySQL 连接：" + name)
	}
}

// Close 关闭数据库连接
func Close() {
	for _, db := range dbInstance {
		if db != nil {
			_ = db.Close()
		}
	}
}

// GetDB 获取带有 span 的 DB 实例
func GetDB(ctx context.Context, name ...string) *gorm.DB {
	var dbName string
	if len(name) == 0 {
		dbName = "default"
	} else {
		dbName = name[0]
	}

	if db, ok := dbInstance[dbName]; ok {
		return SetSpanToGorm(ctx, db)
	} else {
		log.Fatal("[mysql GetDB] 未找到 " + dbName + " 实例")
	}
	return nil
}
