package api //bitcoin account
import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	_ "github.com/btcsuite/btcutil/base58"
	"github.com/gin-gonic/gin"
	"github.com/ygcool/go-hdwallet"
	"golang.org/x/crypto/ripemd160"
	"io/ioutil"
	"net/http"
)

type GormResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

var gormResponse GormResponse

// 用于生成地址的版本
const Version = byte(0x00)

// 用于生成地址的校验和位数
const AddressChecksumLen = 4

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pubKey
}

func (w Wallet) GetAddress() string {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte{Version}, pubKeyHash...)
	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)
	return address
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	RIPEMD160Hasher.Write(publicSHA256[:])
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])
	return secondSHA[:AddressChecksumLen]
}

// /  验证地址是否合法
func ValidateAddress(address string) bool {
	pubKeyHash := base58.Decode(address)
	actualChecksum := pubKeyHash[len(pubKeyHash)-AddressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-AddressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// /  获取账户地址
func GetAddress(c *gin.Context) {
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

	//w := NewWallet()
	gormResponse.Code = http.StatusOK
	/*gormResponse.Message = w.GetAddress()*/
	/*gormResponse.Data = w.PrivateKey.D.Bytes()*/
	gormResponse.Message = mnemonic
	gormResponse.Data = addressP2WPKH
	c.JSON(http.StatusOK, gormResponse)
}

func Importaddress(c *gin.Context) {
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
		gormResponse.Code = 100
		gormResponse.Message = string(mnemonic)
		gormResponse.Data = "助记词错误"
		c.JSON(http.StatusOK, gormResponse)
		panic(err)
	}
	fmt.Println("助记词：", string(mnemonic))
	wallet, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+44), hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(0))
	address, _ := wallet.GetAddress()                               // 1AwEPfoojHnKrhgt1vfuZAhrvPrmz7Rh44
	addressP2WPKH, _ := wallet.GetKey().AddressP2WPKH()             // bc1qdnavt2xqvmc58ktff7rhvtn9s62zylp5lh5sry
	addressP2WPKHInP2SH, _ := wallet.GetKey().AddressP2WPKHInP2SH() // 39vtu9kWfGigXTKMMyc8tds7q36JBCTEHg

	btcwif, err := wallet.GetKey().PrivateWIF(true)
	if err != nil {
		panic(err)
	}
	fmt.Println("BTC私钥：", btcwif)
	fmt.Println("BTC: ", address, addressP2WPKH, addressP2WPKHInP2SH)
	gormResponse.Code = http.StatusOK
	gormResponse.Message = string(mnemonic)
	gormResponse.Data = addressP2WPKH
	c.JSON(http.StatusOK, gormResponse)

}
