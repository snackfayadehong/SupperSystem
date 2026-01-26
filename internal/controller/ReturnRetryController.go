package controller

import (
	"SupperSystem/internal/model"
	clientDb "SupperSystem/pkg/db"
	"SupperSystem/pkg/integration"
	"SupperSystem/pkg/logger"
	"encoding/json"
	"fmt"
	"strings"
)

type ReturnRequestInfo struct {
	Count *int64
	Rn    *[]model.ReturnNo
}

type ReturnResponseInfo struct {
	integration.KLBRBaseResponse
	Data integration.DeliveryData `json:"data"`
}

func (r *ReturnRequestInfo) processSingleReturn(raw model.ReturnNo) error {
	full := &model.ReturnFullSerializer{ReturnNo: &raw} // 构造用于请求 HIS 的接口请求参数
	raw2 := full.ReturnSerializer()
	data, err := json.Marshal(raw2)
	if err != nil {
		return fmt.Errorf("序列化请求数据失败: %w", err)
	}

	// 构建请求
	k := integration.KLBRRequest{
		Headers: integration.NewReqHeaders("herp-clckgl"),
		Url:     integration.BaseUrl + "herp-clckgl/1.0",
		ReqData: data,
	}
	// 发送 HTTP请求
	res, err := k.KLBRHttpPost()
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %w", err)
	}

	// 解析响应
	var HisRes DeliveryResponseInfo
	if err = json.Unmarshal(*res, &HisRes); err != nil {
		logMsg := fmt.Sprintf("\r\n事件:接口请求跟踪\r\n出参：%v\r\n%s\r\n", res, logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 检查响应状态
	if HisRes.AckCode != "200.1" {
		return fmt.Errorf("接口返回错误1: %s", HisRes.AckMessage)
	}

	if len(HisRes.Data.Fhxx) == 0 {
		return fmt.Errorf("响应数据中缺少Fhxx信息")
	}

	fhxx := HisRes.Data.Fhxx[0]
	if strings.TrimSpace(fhxx.Ckdh) == "" && strings.TrimSpace(fhxx.Sczt) != "0" {
		return fmt.Errorf("接口返回错误2: %s", fhxx.Scsm)
	}

	// 记录成功日志
	logMsg := fmt.Sprintf("\r\n事件:接口请求跟踪\r\n出参:%+v\r\n%s\r\n", HisRes.Data, logger.LoggerEndStr)
	logger.AsyncLog(logMsg)

	// 更新数据库
	tx := clientDb.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	hisCkdh := fhxx.Ckdh
	if err = tx.Table("TB_ReturnPurchase").
		Where("ReturnCode = ?", raw.Ckdh).
		Updates(map[string]interface{}{
			"DeliveryNoteCode": hisCkdh,
			"SendStatus":       "1",
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新数据库失败:%w", err)
	}
	if err := tx.Commit(); err != nil {
		return err.Error
	}
	return nil
}

func (r *ReturnRequestInfo) ReturnNoRetryToHis() (err error) {
	var successCount int64
	for _, raw := range *r.Rn {
		if err = r.processSingleReturn(raw); err != nil {
			// 记录错误但继续处理其他项目
			logMsg := fmt.Sprintf("\r\n事件:处理退货单失败\r\n退货单: %s\r\n错误: %v\r\n%s\r\n",
				raw.Ckdh, err, logger.LoggerEndStr)
			logger.AsyncLog(logMsg)
			continue
		}
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("所有退货单处理失败")
	}

	logMsg := fmt.Sprintf("\r\n事件:退货单处理完成\r\n成功数量: %d\r\n总数量: %d\r\n%s\r\n",
		successCount, len(*r.Rn), logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	return nil
}

func (r *ReturnRequestInfo) GetReturnNo(startDate, endDate string) (err error) {
	db := clientDb.DB.Table("TB_ReturnPurchase").
		Select("'03' as ckfs,ReturnCode as ckdh, StoreHouseName,SupplierName").
		Where("Source = ?", 0).
		Where("IsBuy = ?", 0).
		Where("Status = ?", 21).
		Where("SendStatus <> ?", 1).
		Where("AuditorDate >= ?", startDate).
		Where("AuditorDate <= ?", endDate).Find(&r.Rn)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		*r.Count = 0
		logMsg := fmt.Sprintf("\r\n事件:查询产品退货失败业务数据\r\n查询结果:无数据记录\r\n%s\r\n", logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
	} else {
		*r.Count = db.RowsAffected
	}
	return nil
}
