package tokenm

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenM[T any] struct {
	jwt.RegisteredClaims
	Data T
}

// GetTokeStr 使用给定的密钥与载荷数据生成签名后的 JWT 字符串。
//
// 参数:
//   - secretKey: HMAC 对称密钥字符串
//   - data:      任意类型的业务载荷数据, 将存入 TokenM[T].Data
//
// 返回:
//   - string: 已签名的 JWT 字符串
//   - error:  生成或签名失败时返回错误
func GetTokeStr[T any](secretKey string, data T, expiresAt time.Time) (string, error) {
	claims := &TokenM[T]{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  "token",
			Issuer:   "token",
			Audience: jwt.ClaimStrings{"token"},

			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		Data: data,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// GetTokenMData 解析 JWT 字符串并返回其中的业务载荷数据。
//
// 参数:
//   - tokenStr: 需要解析的 JWT 字符串
//   - secretKey: 用于验证签名的 HMAC 对称密钥
//
// 返回:
//   - T:     Token 中的业务载荷数据
//   - error: 解析或验证失败时返回错误
func GetTokenMData[T any](tokenStr string, secretKey string) (T, error) {
	var zero T
	token, err := jwt.ParseWithClaims(tokenStr, &TokenM[T]{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return zero, err
	}
	claims, ok := token.Claims.(*TokenM[T])
	if !ok {
		return zero, fmt.Errorf("invalid claims type")
	}
	return claims.Data, nil
}
