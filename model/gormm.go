package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type Product struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Number   string `gorm:"unique" json:"number"`                       //商品编号（唯一）
	Category string `gorm:"type:varchar(256);not null" json:"category"` //商品类别
	Name     string `gorm:"type:varchar(20);not null" json:"name"`      //商品名称
}

// 应答体
type GormResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

var gormDB *gorm.DB
var gormResponse GormResponse

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

//	func main() {
//		r := gin.Default()
//		//数据库的CRUD--->gin的 post、get、put、delete方法
//		r.POST("gorm/insert", gormInsertData)   //添加数据
//		r.GET("gorm/get", gormGetData)          //查询数据（单条记录）
//		r.GET("gorm/mulget", gormGetMulData)    //查询数据（多条记录）
//		r.PUT("gorm/update", gormUpdateData)    //更新数据
//		r.DELETE("gorm/delete", gormDeleteData) //删除数据
//		r.Run(":9090")
//	}
func GormGetData(c *gin.Context) {
	fmt.Println(c)
	//=============捕获异常============
	defer func() {
		err := recover()
		if err != nil {
			gormResponse.Code = http.StatusBadRequest
			gormResponse.Message = "错误"
			gormResponse.Data = err
			c.JSON(http.StatusBadRequest, gormResponse)
		}
	}()
	//============
	number := c.Query("number")
	fmt.Println(number)
	product := Product{}
	tx := gormDB.Where("number=?", "123456").First(&product)
	if tx.Error != nil {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "查询错误"
		gormResponse.Data = tx.Error
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	gormResponse.Code = http.StatusOK
	gormResponse.Message = "读取成功"
	gormResponse.Data = product
	c.JSON(http.StatusOK, gormResponse)
}

func GormInsertData(c *gin.Context) {
	fmt.Println(c)
	//=============捕获异常============
	defer func() {
		err := recover()
		if err != nil {
			gormResponse.Code = http.StatusBadRequest
			gormResponse.Message = "错误"
			gormResponse.Data = err
			c.JSON(http.StatusBadRequest, gormResponse)
		}
	}()
	//============
	p := &Product{ID: 2, Number: "23456", Category: "NOTHONG", Name: "HHHH"}
	//err := c.Bind(&p)
	//if err != nil {
	//	gormResponse.Code = http.StatusBadRequest
	//	gormResponse.Message = "参数错误"
	//	gormResponse.Data = err
	//	c.JSON(http.StatusOK, gormResponse)
	//	return
	//}

	fmt.Println(p)
	tx := gormDB.Create(&p)
	if tx.RowsAffected > 0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "写入成功"
		gormResponse.Data = "OK"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	//fmt.Printf("insert failed, err:%v\n", err)
	gormResponse.Code = http.StatusBadRequest
	gormResponse.Message = "写入失败"
	gormResponse.Data = tx
	c.JSON(http.StatusOK, gormResponse)
	fmt.Println(tx) //打印结果
}
func GormDeleteData(c *gin.Context) {
	fmt.Println(c)
	//=============捕获异常============
	defer func() {
		err := recover()
		if err != nil {
			gormResponse.Code = http.StatusBadRequest
			gormResponse.Message = "错误"
			gormResponse.Data = err
			c.JSON(http.StatusBadRequest, gormResponse)
		}
	}()
	//============
	p := &Product{ID: 2, Number: "23456", Category: "NOTHONG", Name: "HHHH"}
	//err := c.Bind(&p)
	//if err != nil {
	//	gormResponse.Code = http.StatusBadRequest
	//	gormResponse.Message = "参数错误"
	//	gormResponse.Data = err
	//	c.JSON(http.StatusOK, gormResponse)
	//	return
	//}

	fmt.Println(p)
	tx := gormDB.Delete(&p)
	if tx.RowsAffected > 0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "删除成功"
		gormResponse.Data = "OK"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	//fmt.Printf("insert failed, err:%v\n", err)
	gormResponse.Code = http.StatusBadRequest
	gormResponse.Message = "写入失败"
	gormResponse.Data = tx
	c.JSON(http.StatusOK, gormResponse)
	fmt.Println(tx) //打印结果
}

func GormGetaddress(c *gin.Context) {
	fmt.Println(c)
	//=============捕获异常============
	defer func() {
		err := recover()
		if err != nil {
			gormResponse.Code = http.StatusBadRequest
			gormResponse.Message = "错误"
			gormResponse.Data = err
			c.JSON(http.StatusBadRequest, gormResponse)
		}
	}()
	//============
	number := c.Query("number")
	fmt.Println(number)
	product := Product{}
	tx := gormDB.Where("number=?", "123456").First(&product)
	if tx.Error != nil {
		gormResponse.Code = http.StatusBadRequest
		gormResponse.Message = "查询错误"
		gormResponse.Data = tx.Error
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	gormResponse.Code = http.StatusOK
	gormResponse.Message = "读取成功"
	gormResponse.Data = product
	c.JSON(http.StatusOK, gormResponse)
}
