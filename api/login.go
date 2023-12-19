package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct {
	Id       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// 登录接口
func LoginHandler(c *gin.Context) {
	var loginForm, user Login
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		// 返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 确保数据库中有这张user的表
	rows, err := Db.Query("select * from user where username=?", loginForm.Username)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&user.Id, &user.Username, &user.Password)
	}
	// 判断用户名密码是否正确
	if loginForm.Username != user.Username || loginForm.Password != user.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "用户名或密码不正确",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "登录成功",
		"token": "123jkasdqwe1231a12r13",
	})
}
