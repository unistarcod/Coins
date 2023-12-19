package APP

import (
	"Coins/api"
	"Coins/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"net/http"
)

var (
	StartCmd = &cobra.Command{
		Use:          "server",
		Short:        "server",
		Example:      "server",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)
var (
	Account = &cobra.Command{
		Use:          "address",
		Short:        "address",
		Example:      "address",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			w := api.NewWallet()
			fmt.Println(w.GetAddress())
			fmt.Println("比特币地址:", w.GetAddress())
			fmt.Printf("比特币地址是否有效:%v\n：", api.ValidateAddress(w.GetAddress()))
			return nil
		},
	}
)

func run() error {
	a := &api.Api{Version: "1.0"}
	r := setupRouter()
	r.GET("/hello", a.Getid)
	//r.POST("sql/insert", model.InsertData)   //添加数据
	//r.GET("sql/get", model.GetData)          //查询数据（单条记录）
	//r.GET("sql/mulget", model.GetMulData)    //查询数据（多条记录）
	//r.PUT("sql/update", model.UpdateData)    //更新数据
	//r.DELETE("sql/delete", model.DeleteData) //删除数据
	////数据库的CRUD--->gin的 post、get、put、delete方法
	r.POST("gorm/insert", model.GormInsertData) //添加数据
	r.GET("gorm/get", model.GormGetData)        //查询数据（单条记录）
	r.POST("gorm/delete", model.GormDeleteData) //添加数据
	if err := r.Run(); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
	return nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}
