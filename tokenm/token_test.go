package tokenm

import (
	"fmt"
	"testing"
)

// Test_GetTokeStr_And_Parse_String 生成并解析包含字符串数据的 Token
func Test_GetTokeStr_And_Parse_String(t *testing.T) {
	secret := "secret-123"
	data := "hello"

	tokenStr, err := GetTokeStr[string](secret, data)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	got, err := GetTokenMData[string](tokenStr, secret)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	fmt.Println(tokenStr, got)
	if got != data {
		t.Fatalf("unexpected data: got=%q want=%q", got, data)
	}
}

// Test_GetTokeStr_And_Parse_Struct 生成并解析包含结构体数据的 Token
func Test_GetTokeStr_And_Parse_Struct(t *testing.T) {
	type payload struct {
		Name string
		ID   int
	}
	secret := "secret-xyz"
	data := payload{Name: "bob", ID: 7}

	tokenStr, err := GetTokeStr[payload](secret, data)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	got, err := GetTokenMData[payload](tokenStr, secret)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if got != data {
		t.Fatalf("unexpected data: got=%+v want=%+v", got, data)
	}
}

// Test_Parse_BadSecret 使用错误密钥解析应返回错误
func Test_Parse_BadSecret(t *testing.T) {
	secret := "secret-a"
	tokenStr, err := GetTokeStr[string](secret, "x")
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	if _, err := GetTokenMData[string](tokenStr, "wrong-secret"); err == nil {
		t.Fatalf("expected error when using wrong secret")
	}
}
