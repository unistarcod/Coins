package APP

import (
	"Coins/api"
	"Coins/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/ygcool/go-hdwallet"
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
			//w := api.NewWallet()
			//fmt.Println(w.GetAddress())
			//fmt.Println("比特币地址:", w.GetAddress())
			//fmt.Println("比特币公钥", base58.Encode(w.PublicKey))
			//fmt.Println("比特币私钥", base58.Encode(w.PrivateKey.X.Bytes()))
			//fmt.Printf("比特币地址是否有效:%v\n：", api.ValidateAddress(w.GetAddress()))
			//entropy, _ := bip39.NewEntropy(128)
			//mnemonic, _ := bip39.NewMnemonic(entropy)
			//
			//// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
			//seed := bip39.NewSeed(mnemonic, "Secret Passphrase")
			//
			//masterKey, _ := bip32.NewMasterKey(seed)
			//publicKey := masterKey.PublicKey()
			//
			//seed1 := bip39.NewSeed(mnemonic, "Secret Passphrase")
			//fmt.Println(hex.EncodeToString(seed1))
			//
			//fmt.Println(seed)
			//fmt.Println(seed1)
			//
			//// Display mnemonic and keys
			//fmt.Println("Mnemonic: ", mnemonic)
			//fmt.Println("Master private key: ", masterKey)
			//fmt.Println("Master public key: ", publicKey)
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

func setupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}
