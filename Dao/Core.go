package Dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitDB() {

	dsn := "root:123456@tcp(127.0.0.1:3306)/gin?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//	DB, err = gorm.Open(mysql.Open(config.DBConnectString()), &gorm.Config{
	//		PrepareStmt:            true, //缓存预编译命令
	//		SkipDefaultTransaction: true, //禁用默认事务操作
	//Logger:                 logger.Default.LogMode(logger.Info), //打印sql语句
	//	})
	if err != nil {
		fmt.Println("err1")
		panic(err)
	}
	err = DB.AutoMigrate(&Comment{}, &TableVideo{}, &Follow{}, &Like{}, &TableUser{})
	if err != nil {
		panic(err)
	}

}
