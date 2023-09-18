package service

import (
	"WorkloadQuery/controller"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func GetNoDeliveredPurchaseSummary(c *gin.Context) {
	// 时间区间
	TimeInterval := controller.PurchaseSummaryTimeInterval{}
	// 请求时间绑定
	_ = c.ShouldBindBodyWith(&TimeInterval, binding.JSON)
	TimeInterval.StartTime += " 00:00:00.000"
	TimeInterval.EndTime += " 23:59:59.000"
	//
	purchase := TimeInterval.NotDeliveredPurchaseSummary()
	if purchase == nil || len(*purchase) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"msg":  "无数据",
			"Data": "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "成功",
			"Data": purchase,
		})
	}
}
