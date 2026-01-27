package controller

import (
	"SupperSystem/internal/model"
	clientDb "SupperSystem/pkg/db"
	"SupperSystem/pkg/integration"
	"SupperSystem/pkg/logger"
	"fmt"
	"strings"
)

// DeliveryRequestInfo 接口入参
type DeliveryRequestInfo struct {
	Count *int64
	De    *[]model.DeliveryNo
}

func (d *DeliveryRequestInfo) processSingleDelivery(raw model.DeliveryNo) error {
	full := &model.DeliveryFullSerializer{DeliveryNo: &raw} // 构造用于请求 HIS 的接口请求参数
	raw2 := full.DeliverySerialize()
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
	hisCkdh := fhxx.Ckdh

	if err := tx.Table("TB_DeliveryApplyDetailRecord").
		Where("DeliveryId = ? AND DetailSort = ?", raw.Ckdh, raw.DetailSort).
		Update("OutNumber", hisCkdh).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新数据库失败:%w", err)
	}
	if err := tx.Commit(); err != nil {
		return err.Error
	}
	return nil
}

func (d *DeliveryRequestInfo) DeliveryNoRetryToHis() (err error) {
	var successCount int64
	for _, raw := range *d.De {
		if err = d.processSingleDelivery(raw); err != nil {
			// 记录错误但继续处理其他项目
			logMsg := fmt.Sprintf("\r\n事件:处理配送单失败\r\n配送单: %s\r\n错误: %v\r\n%s\r\n",
				raw.Ckdh, err, logger.LoggerEndStr)
			logger.AsyncLog(logMsg)
			continue
		}
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("所有配送单处理失败")
	}

	logMsg := fmt.Sprintf("\r\n事件:配送单处理完成\r\n成功数量: %d\r\n总数量: %d\r\n%s\r\n",
		successCount, len(*d.De), logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	return nil
}

func (d *DeliveryRequestInfo) GetDeliveryNo(startDate, endDate string) (err error) {
	db := clientDb.DB.Table("TB_DeliveryApplyDetailRecord As dr").
		Select("'01' AS ckfs, d.DeliveryID AS ckdh, dr.DetailSort AS detailSort,"+
			"d.DeliveryCode,DE.DepartmentName as leaderDepartName,DEPT.DepartmentName as storeHouseName").
		Joins("JOIN TB_DeliveryApply d on dr.DeliveryID = d.DeliveryID").
		Joins("Left Join TB_Department dept on dept.DeptCode = d.StorehouseID").
		Joins("Left Join TB_Department de on de.DeptCode = d.DeptCode").
		Where("dr.IsVoid = ?", 0).
		Where("d.Source = ?", "1").
		Where("d.IsStockGoods = ?", "0").
		Where("d.Type = ?", "1").
		Where("d.Status IN ?", []int{61, 71, 41, 81, 22, 91, 19, 29, 99}).
		Where("(d.IsStockGoods <> '1' OR d.IsStockGoods IS NULL)").
		Where("ISNULL(dr.OutNumber, '') = ?", "").
		Where("dr.SendStatus <>?", 1).
		Where("dr.UpdateTime >= ? AND dr.UpdateTime <= ?", startDate, endDate).
		Group("dr.DetailSort, d.DeliveryID,d.DeliveryCode,DEPT.DepartmentName,de.DepartmentName").
		Find(&d.De)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		*d.Count = 0
		logMsg := fmt.Sprintf("\r\n事件:查询领用出库失败业务数据\r\n查询结果:无数据记录\r\n%s\r\n", logger.LoggerEndStr)
		logger.AsyncLog(logMsg)
	} else {
		*d.Count = db.RowsAffected
	}
	return nil
}
