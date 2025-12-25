package service

import (
	"SupperSystem/internal/controller"
	"SupperSystem/internal/model"
	http2 "SupperSystem/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type RetryListRequest struct {
	QueryType string `json:"queryType" binding:"required"`
	StartTime string `json:"startTime" binding:"required"`
	EndTime   string `json:"endTime" binding:"required"`
}

// HandleRetryList 获取待重试列表
func HandleRetryList(c *gin.Context) {
	var req RetryListRequest
	res := http2.NewBaseResponse()
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.QueryType == "delivery" {
		d := &controller.DeliveryRequestInfo{
			Count: new(int64),
			De:    &[]model.DeliveryNo{},
		}
		if err := d.GetDeliveryNo(req.StartTime, req.EndTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		res.Message = "查询成功"
		res.Data = d.De
		c.JSON(http.StatusOK, res)
	} else {
		r := &controller.RefundRequestInfo{
			Count: new(int64),
			Re:    &[]model.RefundNo{},
		}
		if err := r.GetRefundNo(req.StartTime, req.EndTime); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		res.Message = "查询成功"
		res.Data = r.Re
		c.JSON(http.StatusOK, r)
	}
}
