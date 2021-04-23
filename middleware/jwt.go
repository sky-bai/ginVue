package middleware

import (
	"ginVue/utils/errmsg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var JwtKey = []byte("key")

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 1.生成token
func SetToken(username string) (string, int) {
	expiredTime := time.Now().Add(10 * time.Hour)
	setClaims := MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
			Issuer:    "wangxinyue",
		},
	}
	// 1.编码生成前两段
	requestClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, setClaims)
	//  2.签名算法生成后一段
	token, err := requestClaim.SignedString(JwtKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	return token, errmsg.SUCCESS
}

// 2.验证token
func CheckToken(token string) (*MyClaims, int) {
	// 前端传入过来的未验证的tokenstring 后面回调函数接受已解析但未验证的token  根据传过来的tokenstring和数据库中的自己声明的结构体解析出密钥key
	// keyFunc将接收已解析的令牌，并返回用于验证的密钥
	setToken, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if setToken != nil {
		// 这里判断解析后的token是否与数据库里面的token一样 和 token是否过期
		if claims, _ := setToken.Claims.(*MyClaims); setToken.Valid {
			return claims, errmsg.SUCCESS
		}
		return nil, errmsg.ERROR
	}

	return nil, errmsg.ERROR
}

var code int

// jwt中间件  这里处理token是否正确
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		checkToken := strings.SplitN(tokenHeader, " ", 2)
		if len(checkToken) != 2 || checkToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		key, tCode := CheckToken(checkToken[1])
		if tCode == errmsg.ERROR {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		if time.Now().Unix() > key.ExpiresAt {
			code = errmsg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			return
		}
		// 这里设置值
		c.Set("username", key.Username)
		c.Next()
	}
}
