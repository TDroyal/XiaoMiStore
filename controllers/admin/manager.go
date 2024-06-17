package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ManagerController struct { //管理员管理
	BaseController
}

func (con ManagerController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"status":  0,
		"data":    nil,
	})
}

func (con ManagerController) Add(c *gin.Context) {

}

func (con ManagerController) Edit(c *gin.Context) {

}

func (con ManagerController) Delete(c *gin.Context) {

}
