package domain

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// custom claims
// https://mgleon08.github.io/blog/2018/07/16/jwt/
// iss (Issuer) - jwt簽發者
// sub (Subject) - jwt所面向的用戶
// aud (Audience) - 接收jwt的一方
// exp (Expiration Time) - jwt的過期時間，這個過期時間必須要大於簽發時間
// nbf (Not Before) - 定義在什麼時間之前，該jwt都是不可用的
// iat (Issued At) - jwt的簽發時間
// jti (JWT ID) - jwt的唯一身份標識，主要用來作為一次性token,從而迴避重放攻擊
type Claims struct {
	Account         string   `json:"account,omitempty"`
	Roles           []string `json:"roles"`
	Uid             string   `json:"uid,omitempty"`
	Telephone       string   `json:"telephone,omitempty"`
	TelephoneRegion string   `json:"telephoneRegion,omitempty"`
	jwt.StandardClaims
}

type TokenClaims interface{}

// parse and validate token for six things:
// validationErrorMalformed => token is malformed
// validationErrorUnverifiable => token could not be verified because of signing problems
// validationErrorSignatureInvalid => signature validation failed
// validationErrorExpired => exp validation failed
// validationErrorNotValidYet => nbf validation failed
// validationErrorIssuedAt => iat validation failed

type JwtService interface {
	// 簽發jwt
	JwtSign(expiresTime time.Duration, user *UserData, jwtId *string) (expiresAt int64, token string, err error)
	// 驗證jwt
	JwtVerify(token string) (*Claims, error)
	// 驗證過期的jwt
	JwtVerifyExpired(token string) (*Claims, error)
	// 不驗證jwt 解析jwt
	JwtDecode(token string) (TokenClaims, error)
}
