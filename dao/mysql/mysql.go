package mysql

import (
	settings "bluebell/setting"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init(cfg *settings.MySQLConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	// 也可以使用MustConnect连接不成功就panic
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	// 获取底层 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	// 配置连接池
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	return
}

func Close() {
	//获取底层 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	// 配置连接池
	_ = sqlDB.Close()
}
