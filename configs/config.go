package conf

import (
	"SupperSystem/encry"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io"
	"os"
	"path/filepath" // 引入 filepath 包

	"go.uber.org/zap"
)

var Configs *Config

type Config struct {
	Server struct {
		IP              string `json:"Ip"`
		Port            string `json:"Port"`
		RunModel        string `json:"RunModel"`
		HealthCheckPort string `json:"HealthCheckPort"`
	} `json:"Server"`
	DBClient struct {
		IP       string `json:"ip"`
		Username string `json:"username"`
		Password string `json:"password"`
		DbName   string `json:"db_name"`
		IsEc     int    `json:"isEc"`
	} `json:"DBClient"`
	IPWhite struct {
		IPWhiteList []string `json:"IPWhiteList"`
	} `json:"IPWhite"`
	CustomTaskTime struct {
		Run       int `json:"Run"`
		StartTime int `json:"StartTime"`
		EndTime   int `json:"EndTime"`
	} `json:"CustomTaskTime"`
}

// InitSetting 初始化配置
func InitSetting(rootPath string) error {
	// 使用 filepath.Join 安全地拼接路径
	configPath := filepath.Join(rootPath, "configs", "config.json")

	file, err := os.OpenFile(configPath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close() // 确保文件关闭

	v, _ := io.ReadAll(file)
	err = json.Unmarshal(v, &Configs)
	if err != nil {
		return err
	}

	if Configs.DBClient.IsEc == 0 {
		if err = writeEncryptionPwd(rootPath); err != nil {
			return err
		}
	}
	return nil
}

// writeEncryptionPwd 读取配置文件密码加密后重新写入配置文件
func writeEncryptionPwd(rootPath string) error {

	// 秘钥统一存在 configs 目录下
	keyDir := filepath.Join(rootPath, "configs")
	// 1. 生成公钥密钥文件
	if err := encry.GenerateRSAKey(keyDir, 2048); err != nil {
		zap.L().Error("encrypt rsa failed", zap.Error(err))
		return err
	}

	// 2. 读取公钥文件
	pubKeyPath := filepath.Join(keyDir, "public.pem")
	file, err := os.Open(pubKeyPath)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, _ := io.ReadAll(file)
	block, _ := pem.Decode(buf)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	// 3. 打开配置文件准备更新
	configPath := filepath.Join(rootPath, "configs", "config.json")
	configFile, err := os.OpenFile(configPath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer configFile.Close()

	v, _ := io.ReadAll(configFile)
	if err = json.Unmarshal(v, &Configs); err != nil {
		return err
	}

	// 4. 对明文进行加密
	encPwd, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(Configs.DBClient.Password))
	if err != nil {
		return err
	}

	base64Pwd := base64.StdEncoding.EncodeToString(encPwd)
	Configs.DBClient.Password = base64Pwd
	Configs.DBClient.IsEc = 1

	// 5. 重新写入
	newConfig, err := json.MarshalIndent(Configs, "", "  ") // 使用 Indent 使 JSON 更易读
	if err != nil {
		return err
	}

	// 清空原内容并从头写入
	if err = configFile.Truncate(0); err != nil {
		return err
	}
	if _, err = configFile.WriteAt(newConfig, 0); err != nil {
		return err
	}
	return nil
}

// DecryptionPwd 解密
func DecryptionPwd(rootPath string) (pwd string, err error) {
	// 使用 filepath.Join 查找私钥
	privKeyPath := filepath.Join(rootPath, "configs", "private.pem")
	file, err := os.Open(privKeyPath)
	if err != nil {
		return
	}
	defer file.Close()

	buf, _ := io.ReadAll(file)
	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	pwdByte, err := base64.StdEncoding.DecodeString(Configs.DBClient.Password)
	if err != nil {
		return
	}

	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, pwdByte)
	if err != nil {
		return
	}
	return string(decryptedBytes), err
}
