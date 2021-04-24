package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

var db *gorm.DB
var err error

func InitDb() {
	db, err = gorm.Open("mysql", "root:12345678@tcp(localhost:3306)/ginblog?charset=utf8mb4&parseTime=True&loc=Local",
		//utils.DbPassword,
		//utils.DbHost,
		//utils.DbPort,
		//utils.DbName,

	)
	if err != nil {
		fmt.Printf("连接数据库错误，请查看参数", err)
	}

	db.SingularTable(true)
	db.AutoMigrate(&Article{}, &Category{}, &User{})

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	db.DB().SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	db.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	db.DB().SetConnMaxLifetime(10 * time.Second)
}
