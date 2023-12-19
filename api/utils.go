package api

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func DbConnect() *sql.DB {
	// root:password@(localhost:3306)/blog
	// 用户名:密码@(数据库地址及端口号)/数据库名
	db, err := sql.Open("mysql", "root:password@(localhost:3306)/blog")
	db.SetMaxOpenConns(10) // 最大连接数
	db.SetMaxIdleConns(5)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		fmt.Println("连接失败")
		panic(err.Error())
	}
	fmt.Println("连接成功")
	return db // 连接成功后返回db以便在路由中调用
}
