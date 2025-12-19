package controller

import (
	clientDb "WorkloadQuery/db"
	"WorkloadQuery/model"
)

type WorkloadQueryController struct{}

// GetWorkloadRawData 获取扁平的原始记录
func (c *WorkloadQueryController) GetWorkloadRawData(startTime, endTime string) ([]model.RawWorkloadRow, error) {
	var rows []model.RawWorkloadRow
	err := clientDb.DB.Raw(clientDb.FullWorkloadSQL,
		startTime, endTime, startTime, endTime, startTime, endTime, // 入库
		startTime, endTime, startTime, endTime, // 出库
		startTime, endTime, startTime, endTime, startTime, endTime, // 退货
	).Scan(&rows).Error
	return rows, err
}
