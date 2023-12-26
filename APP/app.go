package APP

import (
	"Coins/api"
	"Coins/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/ygcool/go-hdwallet"
	"log"
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
			mnemonic, _ := hdwallet.NewMnemonic(12, "")

			master, err := hdwallet.NewKey(
				hdwallet.Mnemonic(mnemonic),
			)
			if err != nil {
				panic(err)
			}
			fmt.Println("助记词：", mnemonic)
			wallet, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+44), hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(0))
			address, _ := wallet.GetAddress()                               // 1AwEPfoojHnKrhgt1vfuZAhrvPrmz7Rh44
			addressP2WPKH, _ := wallet.GetKey().AddressP2WPKH()             // bc1qdnavt2xqvmc58ktff7rhvtn9s62zylp5lh5sry
			addressP2WPKHInP2SH, _ := wallet.GetKey().AddressP2WPKHInP2SH() // 39vtu9kWfGigXTKMMyc8tds7q36JBCTEHg

			// addressP2WPKHInP2SH的特别说明:这个隔离见证的地址，是属于当前wif私钥的（默认bip44）。
			// 假设你是用生成的助记词导入到imtoken中，对应的隔离见证地址不是这个。
			// 若想和imtoken一致，请在 master.GetWallet 时传入 hdwallet.ZeroQuote+49 （即bip49）得到的隔离见证地址和对应私钥即可
			btcwif, err := wallet.GetKey().PrivateWIF(true)
			if err != nil {
				panic(err)
			}
			fmt.Println("BTC私钥：", btcwif)
			fmt.Println("BTC: ", address, addressP2WPKH, addressP2WPKHInP2SH)
			return nil
		},
	}
)

func run() error {
	a := &api.Api{Version: "1.0"}
	r := setupRouter()

	r.GET("/hello", a.Getid)
	r.POST("gorm/insert", model.GormInsertData) //添加数据
	r.GET("gorm/get", model.GormGetData)        //查询数据（单条记录）
	r.POST("gorm/delete", model.GormDeleteData) //添加数据
	r.GET("gorm/getaddress", api.GetAddress)
	r.POST("gorm/importaddress", api.Importaddress) //添加数据
	if err := r.Run("0.0.0.0:8100"); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
		//r.POST("sql/insert", model.InsertData)   //添加数据
		//r.GET("sql/get", model.GetData)          //查询数据（单条记录）
		//r.GET("sql/mulget", model.GetMulData)    //查询数据（多条记录）
		//r.PUT("sql/update", model.UpdateData)    //更新数据
		//r.DELETE("sql/delete", model.DeleteData) //删除数据
		////数据库的CRUD--->gin的 post、get、put、delete方法
	}
	return nil
}

//解决跨域请求

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())
	return r
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}
