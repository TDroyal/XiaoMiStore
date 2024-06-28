package admin

import (
	"XiaoMiStore/dao"
	"XiaoMiStore/logic"

	"github.com/gin-gonic/gin"
)

// 封装一些公有的方法

type MainController struct{ BaseController }

// 公共修改数据库表中status字段
// 前端传入     数据id     表名     以及status字段

type StatusInfo struct {
	ID    uint   `json:"id"`    //数据id
	Table string `json:"table"` //表名
	Field string `json:"field"` //字段
}

func (con MainController) ChangeStatus(c *gin.Context) {
	statusInfo := StatusInfo{}
	if err := c.ShouldBind(&statusInfo); err != nil {
		con.Error(c, "修改状态失败", -1, nil)
		return
	}
	//执行原生sql
	if err := dao.DB.Exec("update "+statusInfo.Table+" set "+statusInfo.Field+" = abs("+statusInfo.Field+" - 1) where id = ?", statusInfo.ID).Error; err != nil {
		con.Error(c, "修改状态失败", -1, nil)
		return
	}
	con.Success(c, "修改状态成功", 0, nil)
}

// 公共修改数据库表中的数字字段
type NumInfo struct {
	ID     uint   `json:"id"`     //数据id
	Table  string `json:"table"`  //表名
	Field  string `json:"field"`  //字段
	Number int    `json:"number"` //修改的值
}

func (con MainController) ChangeNum(c *gin.Context) {
	numInfo := NumInfo{}
	if err := c.ShouldBind(&numInfo); err != nil {
		con.Error(c, "修改数字失败", -1, nil)
		return
	}
	//执行原生sql
	if err := dao.DB.Exec("update "+numInfo.Table+" set "+numInfo.Field+" = "+logic.IntToString(numInfo.Number)+" where id = ?", numInfo.ID).Error; err != nil {
		con.Error(c, "修改数字失败", -1, nil)
		return
	}
	con.Success(c, "修改数字成功", 0, nil)
}
