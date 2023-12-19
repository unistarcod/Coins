package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// user的结构
type User struct {
	Id       string `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Age      string `json:"age" form:"age"`
	Sex      string `json:"sex" form:"sex"`
}

func UserHandler(c *gin.Context) {
	var user User
	// 定义一个切片数据，存放user数据
	userList := make([]User, 0)
	rows, err := Db.Query("select * from articles where username=?", "张三")
	defer rows.Close() // 即时关闭
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	for rows.Next() { //循环显示所有的数据
		rows.Scan(&user.Id, &user.Username, &user.Age, &user.Sex)
		userList = append(userList, user)
		i++
	}
	// 接口返回给前端的数据
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "请求成功",
		"data": &userList,
	})
}
