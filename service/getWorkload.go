package service

import (
	"WorkloadQuery/controller"
	http2 "WorkloadQuery/http"
	"WorkloadQuery/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetWorkload(c *gin.Context) {
	res := http2.NewBaseResponse()
	workloadTimeInterval := controller.QueryDate{}
	_ = c.ShouldBindBodyWith(&workloadTimeInterval, binding.JSON)
	workloads, _ := workloadTimeInterval.GetWorkloadStatistics()
	if workloads == nil || len(*workloads) == 0 {
		res.Code = 1
		res.Message = "无数据记录"
		res.Data = []model.WorkloadData{}
		c.JSON(http.StatusOK, res)
		return
	}
	res.Data = *workloads
	res.Message = "查询成功"
	res.Code = http.StatusOK
	c.JSON(http.StatusOK, res)
}
