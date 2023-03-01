package token

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xssed/owlcache/config"
)

//Claim是一些实体（通常指的用户）的状态和额外的元数据
type Claims struct {
	TokenId string `json:"token_id"`
	jwt.StandardClaims
}

// 根据用户的用户名和密码产生token
func GenerateToken(token_id string) (string, error) {
	//设置token有效时间
	nowTime := time.Now()
	var expireTime time.Time
	expire_time_int, c_err := strconv.Atoi(config.OwlConfigModel.Tonken_expire_time)
	if c_err != nil {
		return "", c_err
	}
	if expire_time_int == 0 {
		expireTime = nowTime.Add(31536000 * time.Second) //一年有效期
	} else {
		expiration, _ := time.ParseDuration(config.OwlConfigModel.Tonken_expire_time + "s")
		expireTime = nowTime.Add(expiration)
	}

	claims := Claims{
		TokenId: token_id,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "owlcache",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString([]byte(config.OwlConfigModel.Tonken_jwt_secret))
	return token, err
}

// 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseToken(token string) (*Claims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.OwlConfigModel.Tonken_jwt_secret), nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
