package service

import (
	clientDb "WorkloadQuery/db"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type Rep struct {
	Start string `json:"startTime" binding:"required"`
	End   string `json:"endTime" binding:"required"`
}

func GetWorkload(c *gin.Context) {
	rep := Rep{}
	_ = c.ShouldBindBodyWith(&rep, binding.JSON)
	rep.Start += " 00:00:00.000"
	rep.End += " 23:59:59.000"
	uw := clientDb.UserWorkloadQuery(rep.Start, rep.End)
	if uw == nil || len(uw) == 0 {
		c.JSON(http.StatusNoContent, gin.H{
			"msg":  "无数据",
			"Data": "",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "成功",
			"Data": uw,
		})
	}
}
