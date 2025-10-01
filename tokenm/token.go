package tokenm

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenM[T any] struct {
	jwt.RegisteredClaims
	Data      T
	secretKey string
}

// NewTokenM 创建一个基于对称密钥的 Token 管理器实例。
//
// 参数:
//   - secretKey: HMAC 对称密钥, 用于签名与解析
//
// 返回:
//   - *TokenM[T]: 泛型 Token 管理器, 通过方法 Sign 与 Parse 使用
func NewTokenM[T any](secretKey string) *TokenM[T] {
	return &TokenM[T]{
		secretKey: secretKey,
	}
}

// Sign 生成并签名包含指定载荷数据的 JWT 字符串。
//
// 参数:
//   - data:      业务载荷数据, 将写入 TokenM[T].Data
//   - expiresAt: 过期时间
//
// 返回:
//   - string: 已签名的 JWT 字符串
//   - error:  生成或签名失败时返回错误
func (t *TokenM[T]) Sign(data T, expiresAt time.Time) (string, error) {
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
		Data:      data,
		secretKey: t.secretKey,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.secretKey))
}

// Parse 解析 JWT 字符串并返回其中的业务载荷数据。
//
// 参数:
//   - tokenStr: 需要解析的 JWT 字符串
//
// 返回:
//   - T:     Token 中的业务载荷数据
//   - error: 解析或验证失败时返回错误
func (t *TokenM[T]) Parse(tokenStr string) (T, error) {
	var zero T
	token, err := jwt.ParseWithClaims(tokenStr, &TokenM[T]{}, func(token *jwt.Token) (any, error) {
		return []byte(t.secretKey), nil
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

// ParseFromContext 从 http 请求读取 Authorization 并解析业务载荷数据。
// 读取的 Header 名为 "Authorization"。
//
// 参数:
//   - r: http 请求
//
// 返回:
//   - T:     Token 中的业务载荷数据
//   - error: 解析或验证失败时返回错误
func (t *TokenM[T]) ParseFromContext(r *http.Request) (T, error) {
	var zero T
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		return zero, fmt.Errorf("authorization header is empty")
	}
	parts := strings.Split(tokenStr, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return zero, fmt.Errorf("invalid authorization header")
	}

	return t.Parse(parts[1])
}
