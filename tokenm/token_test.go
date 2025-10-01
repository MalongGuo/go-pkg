package tokenm

import (
	"testing"
	"time"
)

// 表驱动测试覆盖字符串、结构体与错误密钥场景
func Test_TokenM_Table(t *testing.T) {
	type (
		payload struct {
			Name string
			ID   int
		}
		testCase[T comparable] struct {
			name      string
			secretKey string
			data      T
			parseKey  string
			wantErr   bool
		}
	)

	cases := []any{
		testCase[string]{
			name:      "string payload ok",
			secretKey: "secret-123",
			data:      "hello",
			parseKey:  "secret-123",
			wantErr:   false,
		},
		testCase[payload]{
			name:      "struct payload ok",
			secretKey: "secret-xyz",
			data:      payload{Name: "bob", ID: 7},
			parseKey:  "secret-xyz",
			wantErr:   false,
		},
		testCase[string]{
			name:      "wrong secret returns error",
			secretKey: "secret-a",
			data:      "x",
			parseKey:  "wrong-secret",
			wantErr:   true,
		},
	}

	for _, c := range cases {
		switch tc := c.(type) {
		case testCase[string]:
			signer := NewTokenM[string](tc.secretKey)
			tokenStr, err := signer.Sign(tc.data, time.Now().AddDate(0, 0, 1))
			if err != nil {
				t.Fatalf("%s: sign token: %v", tc.name, err)
			}
			parser := NewTokenM[string](tc.parseKey)
			got, err := parser.Parse(tokenStr)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("%s: expected error, got nil", tc.name)
				}
				continue
			}
			if err != nil {
				t.Fatalf("%s: parse token: %v", tc.name, err)
			}
			if got != tc.data {
				t.Fatalf("%s: got=%v want=%v", tc.name, got, tc.data)
			}
		case testCase[payload]:
			signer := NewTokenM[payload](tc.secretKey)
			tokenStr, err := signer.Sign(tc.data, time.Now().AddDate(0, 0, 1))
			if err != nil {
				t.Fatalf("%s: sign token: %v", tc.name, err)
			}
			parser := NewTokenM[payload](tc.parseKey)
			got, err := parser.Parse(tokenStr)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("%s: expected error, got nil", tc.name)
				}
				continue
			}
			if err != nil {
				t.Fatalf("%s: parse token: %v", tc.name, err)
			}
			if got != tc.data {
				t.Fatalf("%s: got=%+v want=%+v", tc.name, got, tc.data)
			}
		}
	}
}
