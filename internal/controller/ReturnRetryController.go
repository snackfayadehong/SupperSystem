package controller

import (
	"SupperSystem/internal/model"
	clientDb "SupperSystem/pkg/db"
	"SupperSystem/pkg/integration"
	"SupperSystem/pkg/logger"
	"fmt"
	"strings"
)

type ReturnRequestInfo struct {
	Count *int64
	Rn    *[]model.ReturnNo
}

func (r *ReturnRequestInfo) processSingleReturn(raw model.ReturnNo) error {
	full := &model.ReturnFullSerializer{ReturnNo: &raw} // 构造用于请求 HIS 的接口请求参数
	raw2 := full.ReturnSerializer()
	fhxx, err := integration.SendToHis(raw2, "herp-clckgl/1.0", "herp-clckgl")
	if err != nil {
		return err
	}
	// 校验单据号
	// 逻辑：如果状态不是 "0"(未处理)，则必须有回传的单据号(Ckdh)
	// strings.TrimSpace 用于防止全是空格的情况
	hasId := strings.TrimSpace(fhxx.Ckdh) != ""
	if !hasId && strings.TrimSpace(fhxx.Sczt) != "0" {
		return fmt.Errorf("接口业务处理异常: %s", fhxx.Scsm)
	}
	// 记录成功的业务日志
	logger.AsyncLog(fmt.Sprintf("\r\n事件:接口调用成功\r\n单号:%s\r\n说明:%s\r\n%s\r\n",
		fhxx.Ckdh, fhxx.Scsm, logger.LoggerEndStr))

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
	if err = tx.Table("TB_ReturnPurchase").
		Where("ReturnCode = ?", raw.Ckdh).
		Updates(map[string]interface{}{
			"DeliveryNoteCode": fhxx.Ckdh,
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
