package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ygcool/go-hdwallet"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
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

	mnemonic, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	fmt.Println(string(mnemonic))

	master, err := hdwallet.NewKey(
		hdwallet.Mnemonic(string(mnemonic)),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("助记词：", string(mnemonic))
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

	//w := NewWallet()
	gormResponse.Code = http.StatusOK
	/*gormResponse.Message = w.GetAddress()*/
	/*gormResponse.Data = w.PrivateKey.D.Bytes()*/
	gormResponse.Message = string(mnemonic)
	gormResponse.Data = addressP2WPKH
	c.JSON(http.StatusOK, gormResponse)

	//=============捕获异常============
	//defer func() {
	//	err := recover()
	//	if err != nil {
	//		gormResponse.Code = http.StatusBadRequest
	//		gormResponse.Message = "错误"
	//		gormResponse.Data = err
	//		c.JSON(http.StatusBadRequest, gormResponse)
	//	}
	//}()
	////============
	//p := &Product{ID: 2, Number: "23456", Category: "NOTHONG", Name: "HHHH"}
	////err := c.Bind(&p)
	////if err != nil {
	////	gormResponse.Code = http.StatusBadRequest
	////	gormResponse.Message = "参数错误"
	////	gormResponse.Data = err
	////	c.JSON(http.StatusOK, gormResponse)
	////	return
	////}
	//
	//fmt.Println(p)
	//tx := gormDB.Delete(&p)
	//if tx.RowsAffected > 0 {
	//	gormResponse.Code = http.StatusOK
	//	gormResponse.Message = "删除成功"
	//	gormResponse.Data = "OK"
	//	c.JSON(http.StatusOK, gormResponse)
	//	return
	//}
	////fmt.Printf("insert failed, err:%v\n", err)
	//gormResponse.Code = http.StatusBadRequest
	//gormResponse.Message = "写入失败"
	//gormResponse.Data = tx
	//c.JSON(http.StatusOK, gormResponse)
	//fmt.Println(tx) //打印结果
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
