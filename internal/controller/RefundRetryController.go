package controller

import (
	"SupperSystem/internal/model"
	clientDb "SupperSystem/pkg/db"
	"SupperSystem/pkg/integration"
	"SupperSystem/pkg/logger"
	"fmt"
	"strings"
)

type RefundRequestInfo struct {
	Count *int64
	Re    *[]model.RefundNo
}

func (r *RefundRequestInfo) processSingleRefund(raw model.RefundNo) error {
	full := &model.RefundFullSerializer{RefundNo: &raw}
	raw2 := full.RefundSerialize()
	fhxx, err := integration.SendToHis(raw2, "herp-clrkgl/1.0", "herp-clrkgl")
	if err != nil {
		return err
	}
	// 校验单据号
	// 逻辑：如果状态不是 "0"(未处理)，则必须有回传的单据号(Rkdh)
	// strings.TrimSpace 用于防止全是空格的情况
	hasId := strings.TrimSpace(fhxx.Rkdh) != ""
	if !hasId && strings.TrimSpace(fhxx.Sczt) != "0" {
		return fmt.Errorf("接口业务处理异常: %s", fhxx.Scsm)
	}
	// 记录成功的业务日志
	logger.AsyncLog(fmt.Sprintf("\r\n事件:接口调用成功\r\n单号:%s\r\n说明:%s\r\n%s\r\n",
		fhxx.Rkdh, fhxx.Scsm, logger.LoggerEndStr))
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
	if err := tx.Table("TB_Refund").Where("RetWarhouCode = ?", raw.Yddh).Update("SendStatus", 1).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新数据失败: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return err.Error
	}
	return nil
}

func (r *RefundRequestInfo) RetryRefundToHis() (err error) {
	var successCount int64
	for _, raw := range *r.Re {
		if err = r.processSingleRefund(raw); err != nil {
			// 记录错误但继续处理其他项目
			logMsg := fmt.Sprintf("\r\n事件:处理科室退库单失败\r\n退库单号: %s\r\n错误: %v\r\n%s\r\n",
				raw.Yddh, err, logger.LoggerEndStr)
			logger.AsyncLog(logMsg)
			continue
		}
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("所有科室退库单处理失败")
	}

	logMsg := fmt.Sprintf("\r\n事件:科室退库单处理完成\r\n成功数量: %d\r\n总数量: %d\r\n%s\r\n",
		successCount, len(*r.Re), logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	return nil
}
func (r *RefundRequestInfo) GetRefundNo(startDate, endDate string) (err error) {
	//db := clientDb.DB.Raw(clientDb.QueryRefundBillno, startDate, endDate).Find(&r.Re)
	db := clientDb.DB.Table("TB_Refund a").
		Select("a.RetWarhouCode as yddh,'02' as rkfs,store.DepartmentName as storeHouseName,dept.DepartmentName as leaderDepartName").
		Joins("Left Join TB_Department store on store.DeptCode = a.TargetStorehouseID").
		Joins("Left Join TB_Department dept on dept.DeptCode = a.DeptCode").
		Where("ISNULL(a.SendStatus,?) = ?", "", "").
		Where("a.Status = ?", 51).
		Where("a.CreateTime >= ? And a.CreateTime < ?", startDate, endDate).
		Find(&r.Re)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		*r.Count = 0
		logMsg := fmt.Sprintf("\r\n事件:查询科室退库失败业务数据\r\n查询结果:无数据记录\r\n%s\r\n", logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
	} else {
		*r.Count = db.RowsAffected
	}
	return nil
}
