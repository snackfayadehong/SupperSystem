package clientDb

import (
	"WorkloadQuery/encry"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"os"
)

var (
	DB      *gorm.DB
	configs *config
)

type config struct {
	DBClient struct {
		IP       string `json:"ip"`
		Username string `json:"username"`
		Password string `json:"password"`
		DbName   string `json:"db_name"`
	} `json:"DBClient"`
	IsEc int `json:"isEc"`
}

func InitDb() error {
	var DbPwd string
	file, err := os.OpenFile("../config.json", os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	v, _ := io.ReadAll(file)
	err = json.Unmarshal(v, &configs)
	if configs.IsEc == 0 {
		file.Close()
		err = writeEncryptionPwd()
		if err != nil {
			return err
		}
	}
	DbPwd, err = decryptionPwd()
	if err != nil {
		return err
	}
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&loc=%s&encrypt=disable", userName, password, ipAddr, port, dbName, loc)
	// dsn := "sqlserver://sa:密码@127.0.0.1:1433?database=dbStatus&encrypt=disable"
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=disable", configs.DBClient.Username, DbPwd, configs.DBClient.IP, configs.DBClient.DbName)
	DB, _ = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "TB_",
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	return nil
}

// 读取配置文件密码加密后重新写入配置文件
func writeEncryptionPwd() error {
	// 生成公钥密钥文件
	encry.GenerateRSAKey(2048)
	// 打开public.pem公钥文件
	file, err := os.Open("../encry/public.pem")
	defer file.Close()
	// 读取公钥
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	// pem解码
	block, _ := pem.Decode(buf)
	// x509
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	// 类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	configFile, err := os.OpenFile("../config.json", os.O_RDWR, 0666)
	defer configFile.Close()
	if err != nil {
		return err
	}
	v, _ := io.ReadAll(configFile)
	err = json.Unmarshal(v, &configs)
	if err != nil {
		return err
	}
	// 密码加密
	// 对明文进行加密
	encPwd, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(configs.DBClient.Password))
	if err != nil {
		return err
	}
	// 转base64
	base64Pwd := base64.StdEncoding.EncodeToString(encPwd)
	configs.DBClient.Password = base64Pwd
	// 标记加密
	configs.IsEc = 1
	newConfig, err := json.Marshal(configs)
	_, err = configFile.WriteAt(newConfig, 0)
	if err != nil {
		return err
	}
	return err
}

// 解密
func decryptionPwd() (pwd string, err error) {
	// 打开私钥文件
	file, err := os.Open("../encry/private.pem")
	if err != nil {
		return
	}
	defer file.Close()
	// 读取私钥文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	// pem解码
	block, _ := pem.Decode(buf)
	// x509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	pwdByte, err := base64.StdEncoding.DecodeString(configs.DBClient.Password)
	if err != nil {
		return
	}
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, pwdByte)
	if err != nil {
		return
	}
	return string(decryptedBytes), err
}
