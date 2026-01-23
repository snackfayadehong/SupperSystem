package controller

import (
	"SupperSystem/internal/model"
	clientDb "SupperSystem/pkg/db"
)

type SupperSystemController struct{}

// GetWorkloadRawData 获取扁平的原始记录
func (c *SupperSystemController) GetWorkloadRawData(startTime, endTime string) ([]model.RawWorkloadRow, error) {
	var rows []model.RawWorkloadRow
	err := clientDb.DB.Raw(clientDb.FullWorkloadSQL,
		startTime, endTime, // 入库登记
		startTime, endTime, startTime, endTime, startTime, endTime, // 入库验收
		startTime, endTime, startTime, endTime, // 出库
		startTime, endTime, startTime, endTime, startTime, endTime, // 退货
		startTime, endTime, // 采购订单发送
		startTime, endTime, // 采购订单催货
		startTime, endTime, // 二级库退库
	).Scan(&rows).Error
	return rows, err
}
