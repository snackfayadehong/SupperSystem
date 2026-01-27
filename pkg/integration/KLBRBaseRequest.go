package integration

import (
	"SupperSystem/pkg/logger"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const BaseUrl = `http://172.21.1.248:17980/api/bsp-api-engine-others/`

// KLBRReqHeaders 柯林布瑞业务中台请求头
type KLBRReqHeaders struct {
	appId       string
	timestamp   string
	messageId   string
	signature   string
	contentType string
}

// KLBRBaseResponse 柯林布瑞业务中台接口返回
type KLBRBaseResponse struct {
	AckCode      string `json:"ackCode"`
	AckMessage   string `json:"ackMessage"`
	AckMessageID string `json:"ackMessageId"`
}

// FhxxData 统一明细数据结构
// 包含出库、退库、入库、产品信息变更等可能用到的字段,向下兼容
type FhxxData struct {
	Ckdh   string `json:"ckdh"` // 出库单号
	Rkdh   string `json:"rkdh"` // 入库单号
	Sczt   string `json:"sczt"`
	Scsm   string `json:"scsm"`   //说明信息
	Ydcldm string `json:"ydcldm"` // 怡道材料代码
	Ypdm   string `json:"ypdm"`   // 材料代码
	Scdm   string `json:"scdm"`   // 0成功
}

// GenericResponse 通用响应包装结构
type GenericResponse struct {
	KLBRBaseResponse
	Data struct {
		Fhxx []FhxxData `json:"fhxx"`
	} `json:"data"`
}

// KLBRRequest 柯林布瑞接口请求参数
type KLBRRequest struct {
	Headers *KLBRReqHeaders
	Url     string
	ReqData []byte
}

// SendToHis 通用发送与解析方法
// requestData: 待序列化的请求结构体 如；model.DeliveryFullSerializer
// apiSuffix: 接口的后缀的Url 如: "herp-clckgl/1.0"
// headerType: 签名用的ServiceCode 如: "herp-clckgl"
func SendToHis(requestData interface{}, apiSuffix, headerType string) (*FhxxData, error) {
	// 1. 序列化请求Data
	data, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("请求数据序列化失败: %w", err)
	}
	// 2. 构建基础请求对象
	k := KLBRRequest{
		Headers: NewReqHeaders(headerType),
		Url:     BaseUrl + apiSuffix,
		ReqData: data,
	}
	// 3. Http Post
	resBytes, err := k.KLBRHttpPost()
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	// 4. 解析到通用返回结构
	var genericRes GenericResponse
	if err = json.Unmarshal(*resBytes, &genericRes); err != nil {
		logMsg := fmt.Sprintf("\r\n事件:响应解析异常\r\n原始响应:%s\r\n%s\r\n", string(*resBytes), logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
		return nil, fmt.Errorf("响应JSON解析失败: %w", err)
	}
	// 5. 业务校验
	if genericRes.AckCode != "200.1" {
		return nil, fmt.Errorf("接口返回失败状态(%s): %s", genericRes.AckCode, genericRes.AckMessage)
	}
	// 6. Fhxx检查
	if len(genericRes.Data.Fhxx) == 0 {
		return nil, fmt.Errorf("响应数据中缺少Fhxx明细")
	}
	fhxx := genericRes.Data.Fhxx[0]
	return &fhxx, nil
}

// NewReqHeaders  KLBRReqHeaders请求头构造函数，根据入参信息生成请求Headers
func NewReqHeaders(serviceCode string) *KLBRReqHeaders {
	reqHeaders := new(KLBRReqHeaders)
	reqHeaders.appId = "HERP"
	reqHeaders.timestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)
	uuid, _ := uuid.NewUUID()
	reqHeaders.messageId = uuid.String()
	reqHeaders.contentType = "json"
	var signStr = fmt.Sprintf("appId=%s&serviceCode=%s&version=%s&timestamp=%v",
		reqHeaders.appId, serviceCode, "1.0", reqHeaders.timestamp)
	reqHeaders.signature = HMACSHA1(signStr)
	return reqHeaders
}

// HMACSHA1 加密转base64
func HMACSHA1(str string) string {
	keyStr := "123456"
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(str))
	//进行base64编码
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return res
}

// KLBRHttpPost 柯林布瑞业务中台接口统一封装方法
func (k *KLBRRequest) KLBRHttpPost() (*[]byte, error) {
	reqData := bytes.NewBuffer(k.ReqData)
	reqBody, err := http.NewRequest("POST", k.Url, reqData)
	if err != nil {
		return nil, err
	}
	logMsg := fmt.Sprintf("\r\n事件:接口请求跟踪\r\n接口地址:%s\r\n入参:%s\r\n%s\r\n", k.Url, string(k.ReqData), logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	defer reqBody.Body.Close()
	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("appId", k.Headers.appId)
	reqBody.Header.Set("timestamp", k.Headers.timestamp)
	reqBody.Header.Set("messageId", k.Headers.messageId)
	reqBody.Header.Set("signature", k.Headers.signature)
	reqBody.Header.Set("contentType", k.Headers.contentType)
	rep, err := http.DefaultClient.Do(reqBody)
	if err != nil {
		return nil, err
	}
	repBytes, _ := io.ReadAll(rep.Body)
	return &repBytes, nil
}
