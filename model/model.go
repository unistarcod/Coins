package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	var err error
	sqlStr := "root:123456@tcp(127.0.0.1:3306)/DATA?charset=utf8mb4&parseTime=true&loc=Local"
	gormDB, err = gorm.Open(mysql.Open(sqlStr), &gorm.Config{}) //配置项中预设了连接池 ConnPool
	if err != nil {
		fmt.Println("数据库连接出现了问题：", err)
		return
	}
	fmt.Println("链接成功")
}
