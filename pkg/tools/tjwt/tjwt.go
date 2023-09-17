package tjwt

import (
	"github.com/golang-jwt/jwt/v4"
	"rt-chat/pkg/tools/tencrypt"
	"time"
)

// JSON Web Token
/*
在RFC标准中，JWT由三部分组成：
Header头部
Payload载荷 主要包含声明(claims)部分
	registered(iss签发者)、(exp过期时间)、(aud受众)
	public由使用JWT的人随意定义
	private claims 自定义，用于在服务双方共享一些信息
Signature签名
*/

const (
	TokenMaxExpireHour      = 2  // token 最长有效期
	TokenMaxRemainingMinute = 15 // token 还有多久过期就返回新token
)

type MyClaims struct {
	jwt.RegisteredClaims
	User string `json:"user"`
}

type Token struct {
	AccessToken  string `json:"accessToken" form:"accessToken" comment:"token 1小时"`
	RefreshToken string `json:"refreshToken" form:"refreshToken" comment:"刷续期的token(7天 只可用一次)"`
}

// GenerateToken 生成token，返回string型token和err
func GenerateToken(uid string) (string, error) {
	//获取rsa的私钥
	privateKey, err := tencrypt.GetPrivateKey()
	if err != nil {
		return "", err
	}
	//------------
	// claims
	// 签发者：Server、过期时间：两小时、（在签发之前无效、签发时间）
	claims := MyClaims{
		User: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Server",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * TokenMaxExpireHour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// 私钥加密
	signedString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return signedString, nil
}

// CheckToken 校验token，返回自定义claims、err
func CheckToken(signedString string) (*MyClaims, error) {

	// 获取rsa的公钥
	publicKey, err := tencrypt.GetPublicKey()
	if err != nil {
		return nil, err
	}
	// 公钥解密
	token, err := jwt.ParseWithClaims(signedString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// RefreshToken 更新token，token过期时间 - 最大截止刷新时间 > 当前时间 ？；如果true，则返回token。添加headers['new-token']
func RefreshToken(claims *MyClaims) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	// 即将超过过期期限
	if t := claims.ExpiresAt.Time.Add(-time.Minute * TokenMaxRemainingMinute); t.Before(time.Now()) {
		return GenerateToken(claims.User)
	}
	return "", nil
}
