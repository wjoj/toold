package toold

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//JWTCodeType type
type JWTCodeType int

//JWTCodeType type
const (
	JWTCodeTypeSuccess          JWTCodeType = 0
	JWTCodeTypeTokenFail                    = 40001
	JWTCodeTypeTokenGetInfoFail             = 40002
	JWTCodeTypeTokenTimeOutFail             = 40003
)

var jwtMessage = map[JWTCodeType]string{
	JWTCodeTypeSuccess:          "成功",
	JWTCodeTypeTokenFail:        "无验证错误",
	JWTCodeTypeTokenGetInfoFail: "获取验证信息失败",
	JWTCodeTypeTokenTimeOutFail: "验证超时",
}

//Authorization auth
const Authorization = "Authorization"

var jwtSecret = []byte("")

var expireTime time.Duration

//Claims Claims
type Claims struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	jwt.StandardClaims
}

//NewJWTConfig NewJWTConfig
func NewJWTConfig(expire time.Duration, secret string) {
	expireTime = expire
	jwtSecret = []byte(secret)
}

func getMessage(code JWTCodeType) string {
	jwtMessage := map[JWTCodeType]string{
		JWTCodeTypeSuccess:          "成功",
		JWTCodeTypeTokenFail:        "验证信息错误",
		JWTCodeTypeTokenGetInfoFail: "验证信息已失效",
		JWTCodeTypeTokenTimeOutFail: "验证信息已过期",
	}
	msg := jwtMessage[code]
	if len(msg) == 0 {
		return "验证错误"
	}
	return msg
}

/*
GenerateToken 获取token
*/
func GenerateToken(username, password string) (string, int64, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(expireTime)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.UnixNano(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, expireTime.Unix(), err
}

//ParseToken ParseToken
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
		if claims != nil {
			return claims, nil
		}
	}

	return nil, err
}

//JWT JWT
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code JWTCodeType
		var data interface{}
		code = JWTCodeTypeSuccess
		token := c.GetHeader(Authorization)
		if token == "" {
			code = JWTCodeTypeTokenFail
		} else {
			claims, err := ParseToken(token)
			if err != nil {
				code = JWTCodeTypeTokenGetInfoFail
			} else if time.Now().UnixNano() > claims.ExpiresAt {
				code = JWTCodeTypeTokenTimeOutFail
			} else {
				claims.ExpiresAt = time.Now().Add(expireTime).UnixNano()
			}
		}

		if code != JWTCodeTypeSuccess {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  getMessage(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
