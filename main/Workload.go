package main

import (
	clientDb "WorkloadQuery/db"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Employee struct {
	HRCode       int
	EmployeeName string
}

func main() {
	user := Employee{}
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		clientDb.DB.Select("HRCode,EmployeeName").First(&user)
		u, _ := json.Marshal(user)
		c.String(http.StatusOK, string(u))
	})
	err := r.Run("127.0.0.1:3305")
	if err != nil {
		return
	}
}

func init() {
	err := clientDb.InitDb()
	if err != nil {
		return
	}
}
