package model

import (
	"fmt"
	"github.com/HauKuen/Annals/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var db *gorm.DB

func InitDb() {
	var err error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.User, utils.Password, utils.Host, utils.Port, utils.Dbname)
	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
	})

	if err != nil {
		utils.Log.Fatal("Failed to connect database:", err)
	}

	err = db.AutoMigrate(&User{}, &Category{})
	if err != nil {
		utils.Log.Error("Auto migrate failed:", err)
	}

	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	utils.Log.Info("Database initialized successfully")
}
