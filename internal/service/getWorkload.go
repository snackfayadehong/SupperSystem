package service

import (
	"SupperSystem/internal/controller"
	"SupperSystem/internal/model"
	http2 "SupperSystem/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

var WorkloadServiceInstance = &WorkloadService{
	workloadCtrl: controller.SupperSystemController{},
}

type WorkloadService struct {
	workloadCtrl controller.SupperSystemController
}

func (s *WorkloadService) HandleWorkloadRequest(c *gin.Context) {
	res := http2.NewBaseResponse()
	// 从中间件 context 获取时间
	starTime := c.GetString("startTime")
	endTime := c.GetString("endTime")
	data, err := s.GetWorkloadReport(starTime, endTime)
	if err != nil {
		res.Code = 1
		res.Message = "查询失败" + err.Error()
		res.Data = []model.WorkloadGroup{}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	res.Data = data
	res.Message = "查询成功"
	c.JSON(http.StatusOK, res)
}

// mapDeptToCategory 统一维护映射逻辑
func (s *WorkloadService) mapDeptToCategory(code string, fallback string) string {
	switch code {
	case "200346", "200416":
		return "高值扫码材料"
	case "200650", "200632":
		return "高值扫码材料(经开分院)"
	case "200418":
		return "低值扫码材料"
	case "200420":
		return "普通卫生材料"
	case "200438", "200426":
		return "试剂与化验材料库"
	default:
		if fallback != "" {
			return fallback
		}
		return code
	}
}

// GetWorkloadReport 聚合业务逻辑：从 Controller 获取数据并重新组装结构
func (s *WorkloadService) GetWorkloadReport(startTime, endTime string) ([]model.WorkloadGroup, error) {
	raws, err := s.workloadCtrl.GetWorkloadRawData(startTime, endTime)
	if err != nil {
		return nil, err
	}

	groupMap := make(map[string]*model.WorkloadGroup)
	for _, r := range raws {
		if _, ok := groupMap[r.OperatorName]; !ok {
			groupMap[r.OperatorName] = &model.WorkloadGroup{
				Operator:        r.OperatorName,
				Inbound:         []model.WorkloadDetail{},
				Outbound:        []model.WorkloadDetail{},
				Return:          []model.WorkloadDetail{},
				InReg:           []model.WorkloadDetail{},
				Purchase:        []model.WorkloadDetail{},
				Push:            []model.WorkloadDetail{},
				SecondaryRefund: []model.WorkloadDetail{},
			}
		}
		group := groupMap[r.OperatorName]
		detail := model.WorkloadDetail{
			// 调用 Go 内部的转换函数
			Category:    s.mapDeptToCategory(r.StorehouseCode, r.FallbackName),
			SpecCount:   r.SpecCount,
			BillCount:   r.BillCount,
			TotalAmount: r.TotalAmount,
		}

		switch r.OperationType {
		case 0:
			group.Inbound = append(group.Inbound, detail)
		case 1:
			group.Outbound = append(group.Outbound, detail)
		case 2:
			group.Return = append(group.Return, detail)
		case 3:
			group.InReg = append(group.InReg, detail)
		case 4:
			detail.Category = "采购单发送"
			group.Purchase = append(group.Purchase, detail)
		case 5:
			detail.Category = "催货业务"
			group.Push = append(group.Push, detail)
		case 6:
			group.SecondaryRefund = append(group.SecondaryRefund, detail)
		}
	}

	var results []model.WorkloadGroup
	for _, v := range groupMap {
		results = append(results, *v)
	}
	return results, nil
}
