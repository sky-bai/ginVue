package v1

import (
	"ginVue/model"
	"ginVue/utils/errmsg"
	"ginVue/utils/validator"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 查询用户是否存在
func UserExist(c *gin.Context) {

}

var code int

// 添加用户   根据前端传过来的请求添加用户
func AddUser(c *gin.Context) {
	var data model.User         // 获取用户结构体
	_ = c.ShouldBindJSON(&data) // 将json数据绑定到后端结构体上  绑定好前端传过来的数据
	// 这时就拿到前端传过来的json数据 现在通过操控后端结构体的数据就可以了

	// 1.验证前端发过来的数据是否有效
	msg, code := validator.Validate(&data)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": msg,
		})
		return // 不准往下走
	}
	code = model.CheckUser(data.Username) // 判定传过的数据是否存在
	// 如果不存在就创建新的用户
	if code == errmsg.SUCCESS {
		model.CreateUser(&data)
	}
	// 如果判定为用户名已使用 就将code赋值为错误代码
	if code == errmsg.ERROR_USERNAME_USED {
		code = errmsg.ERROR_USERNAME_USED
	}
	// 往前端返回数据 第一个是请求成功 返回请求成功的状态码 第二个返回一个结构体
	c.JSON(http.StatusOK, gin.H{
		// 本身自定义的错误信息
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})

}

// 删除用户
func DeleteUser(c *gin.Context) {
	// 根据传入的请求 获取id 返回删除正确与否的状态码
	id, _ := strconv.Atoi(c.Param("id"))
	code = model.DeleteUser(id)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

//获取单个用户
func GetUserInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   1,
		"message": errmsg.GetErrMsg(code),
	})
}

// 获取用户列表   params 查询参数
func GetUsers(c *gin.Context) {
	// 1.先从请求的url获取到需要查询到值
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	userName := c.Query("username")

	// 2. 在gorm中 参数为-1表示获取所有记录
	if pageSize == 0 {
		pageSize = -1
	}
	if pageNum == 0 {
		pageNum = -1
	}

	// 3.调用dao层 传入每页的个数 所有记录个数
	data, total := model.GetUsers(userName, pageSize, pageNum)
	code = errmsg.SUCCESS
	// 设置json的数据
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": errmsg.GetErrMsg(code),
	})
}

// 编辑用户
func EditUser(c *gin.Context) {
	// 根据传入的id 和jsondata数据进行更新
	// 1.获取id
	id, _ := strconv.Atoi(c.Param("id"))

	// 2.获取json_data数据
	var data model.User         // 获取用户结构体
	_ = c.ShouldBindJSON(&data) // 将json数据绑定到后端结构体上

	// 3.如果新设置的用户名不存在就根据id修改用户名
	code = model.CheckUpUser(id, data.Username)
	if code == errmsg.SUCCESS {
		model.EddUser(id, &data)
	}
	// 4.如果用户存在就暂停
	if code == errmsg.ERROR_USERNAME_USED {
		// 阻止调用后续函数
		c.Abort()
	}
	// 5.返回状态结果 根据前面的code返回对应的状态信息
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})

}
