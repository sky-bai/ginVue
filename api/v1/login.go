package v1

import (
	"fmt"
	"ginVue/middleware"
	"ginVue/model"
	"ginVue/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context)  {
	var data model.User
	// 1.将前端传过来的用户数据绑定到data结构体上
	err := c.ShouldBindJSON(&data)
	if err != nil {
		fmt.Print(err)
	}
	// 2.检查是否为管理员
	code := model.CheckLogin(data.Username, data.Password)
	// 3.如果是管理员，就根据用户的用户名发放token
	if code == errmsg.SUCCESS{
		token, _ := middleware.SetToken(data.Username)
		c.JSON(http.StatusOK,gin.H{
			"status":code,
			"message":errmsg.GetErrMsg(code),
			"token":token,
		})
	}else {
		c.JSON(http.StatusOK, gin.H{
			"status":    code,
			"message": errmsg.GetErrMsg(code),
		})
	}
}
