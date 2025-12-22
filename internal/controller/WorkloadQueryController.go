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
		startTime, endTime, startTime, endTime, startTime, endTime, // 入库
		startTime, endTime, startTime, endTime, // 出库
		startTime, endTime, startTime, endTime, startTime, endTime, // 退货
	).Scan(&rows).Error
	return rows, err
}
