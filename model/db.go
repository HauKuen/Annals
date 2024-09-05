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
		//gorm日志模式：silent
		Logger: logger.Default.LogMode(logger.Silent),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 使用单数表名
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})

	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}
	err = db.AutoMigrate(&User{}, &Category{})
	if err != nil {
		fmt.Println("auto migrate failed, err:" + err.Error())
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

}
