package service

import (
	"WorkloadQuery/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type Response struct {
	Code    int
	Message string
}

var res Response

func ChangeProductInfoService(c *gin.Context) {
	// 查询条件
	var where []string
	// 入参
	var req controller.RequestInfo
	_ = c.ShouldBindBodyWith(&req.C, binding.JSON)
	// 将入参多条数据Code整合为一个where条件
	var seen = make(map[string]bool)
	for _, v := range req.C {
		if !seen[v.Code] {
			seen[v.Code] = true
			where = append(where, v.Code)
		} else {
			res.Code = 1
			res.Message = fmt.Sprintf("%s入参存在多条;", v.Code)
			c.JSON(http.StatusCreated, res)
			return
		}
	}
	// 获取怡道系统产品基本信息
	_, msg := req.GetProductInfo(where)
	res.Message = msg
	c.JSON(http.StatusCreated, res)
}
