package service

import (
	"WorkloadQuery/controller"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func GetNoAccountEntry(c *gin.Context) {
	// 查询时间区间
	NoAccountEntryTimeInterval := controller.QueryAccountEntryTime{}
	// 将请求的时间绑定到NoAccountEntryTimeInterval
	_ = c.ShouldBindBodyWith(&NoAccountEntryTimeInterval, binding.JSON)
	NoAccountEntryTimeInterval.StartTime += " 00:00:00.000"
	NoAccountEntryTimeInterval.EndTime += " 23:59:59.000"
	// 调用方法返回未上账数据
	NoAccountEntryBills := NoAccountEntryTimeInterval.NoAccountEntryQuery()
	if NoAccountEntryBills == nil || len(*NoAccountEntryBills) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"msg":  "无数据",
			"Data": "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "成功",
			"Data": NoAccountEntryBills,
		})
	}
}
