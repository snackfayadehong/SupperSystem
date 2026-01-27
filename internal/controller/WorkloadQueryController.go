package controller

import (
	"SupperSystem/internal/model"
	clientDb "SupperSystem/pkg/db"
)

type SupperSystemController struct{}

// GetWorkloadRawData 获取扁平的原始记录
func (c *SupperSystemController) GetWorkloadRawData(startTime, endTime string) ([]model.RawWorkloadRow, error) {
	var rows []model.RawWorkloadRow
	params := map[string]interface{}{
		"startTime": startTime,
		"endTime":   endTime,
	}
	err := clientDb.DB.Raw(clientDb.FullWorkloadSQL, params).Scan(&rows).Error
	return rows, err
}
