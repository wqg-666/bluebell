package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func signingSecret() ([]byte, error) {
	secret := viper.GetString("auth.jwt_secret")
	if secret == "" {
		return nil, errors.New("auth.jwt_secret or BLUEBELL_JWT_SECRET must be configured")
	}
	return []byte(secret), nil
}

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	secret, err := signingSecret()
	if err != nil {
		return "", err
	}
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			//viper 是 Go 中一个非常常用的配置管理库，主要用来读取项目中的配置文件，例如 config.yaml、config.json 等。
			//你这段代码是在使用 Viper 从配置文件中读取 JWT 的过期时间。
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(), // 过期时间
			Issuer:    "bluebell",                                                                        // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	//jwt.SigningMethodHS256 JWT 签名算法
	//c  JWT 中携带的声明信息
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名,并获得完整的编码后的字符串token
	return token.SignedString(secret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	secret, err := signingSecret()
	if err != nil {
		return nil, err
	}
	// 解析token
	var mc = new(MyClaims) //创建一个指向MyClaims对象的指针,一定要手动申请内存，即使返回值是mc *MyClaims
	//jwt.ParseWithClaims() 主要完成三个工作：
	//1. 解析 JWT 字符串。
	//2. 将 JWT 中的 Payload 数据解析到 mc 中。也就是把 Token 中的 Claims 数据填充到 mc 结构体中
	//3. 使用密钥验证 JWT 的签名。
	//tokenString传入需要解析的 JWT 字符串。mc告诉 JWT 库：解析出来的用户信息应该存放在哪里。也就是把 Token 中的 Claims 数据填充到 mc 结构体中
	//func(token *jwt.Token) (interface{}, error) {
	//    return secret, nil
	//}
	//这是一个回调函数，也叫密钥提供函数。它的作用是：告诉 JWT 库，使用什么密钥来验证 Token 的签名。
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected JWT signing method")
		}
		return secret, nil
	})
	//如果 JWT 解析过程中出现错误，就直接返回。Token 格式不正确。签名验证失败。Token 已经过期。
	if err != nil {
		return nil, err
	}
	//判断 Token 是否有效
	if token.Valid {
		return mc, nil
	}
	//返回无效错误
	return nil, errors.New("invalid token")
}
