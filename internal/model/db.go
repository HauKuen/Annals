package model

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/HauKuen/Annals/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func InitDb() error {
	var err error
	maxRetries := 5
	retryDelay := time.Second * 3

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.User, utils.Password, utils.Host, utils.Port, utils.Dbname)

	config := getGormConfig()

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dns), config)
		if err == nil {
			break
		}

		utils.Log.Warnf("数据库连接失败，%d秒后重试 (尝试 %d/%d): %v",
			retryDelay/time.Second, i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(utils.DbMaxIdleConns)
	sqlDB.SetMaxOpenConns(utils.DbMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(utils.DbConnMaxLifetime) * time.Minute)

	if err := autoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	utils.Log.Info("数据库初始化成功")
	return nil
}

func autoMigrate() error {
	if !needMigration() {
		return nil
	}

	if err := db.AutoMigrate(&User{}, &Category{}, &Article{}); err != nil {
		return err
	}

	return nil
}

func needMigration() bool {
	return true
}

func getLogLevel() logger.LogLevel {
	if utils.AppMode == "debug" {
		return logger.Silent
	}
	return logger.Silent
}

func getGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  getLogLevel(),
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

func CheckDatabaseHealth() error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	return nil
}
