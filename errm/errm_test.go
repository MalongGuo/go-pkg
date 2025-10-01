package errm

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TestNewContainsMessageAndCause 验证 New 返回的错误文本包含自定义消息与底层错误文本
func TestNewContainsMessageAndCause(t *testing.T) {
	root := errors.New("root cause")
	e := New(root, "additional context")

	if e == nil {
		t.Fatalf("expected non-nil error")
	}

	if !strings.Contains(e.Error(), "additional context") {
		t.Fatalf("error text should contain custom message, got: %q", e.Error())
	}

	if !strings.Contains(e.Error(), root.Error()) {
		t.Fatalf("error text should contain root cause message, got: %q", e.Error())
	}

	if got := e.Cause(); got == nil {
		t.Fatalf("Cause() should not be nil")
	} else if got.Error() != e.Error() {
		t.Fatalf("Cause().Error() should equal wrapper Error(); got=%q want=%q", got.Error(), e.Error())
	}
}

// TestWithMessagefFormats 验证 WithMessagef 的格式化输出
func TestWithMessagefFormats(t *testing.T) {
	e := WithMessagef(errors.New("E"), "code=%d", 42)
	if !strings.Contains(e.Error(), "code=42") {
		t.Fatalf("formatted message not found, got: %q", e.Error())
	}
}

// TestWithMessageNilCause 验证当 cause 为 nil 时也能返回包含消息的错误
func TestWithMessageNilCause(t *testing.T) {
	e := WithMessage(errors.New("E"), "hello")
	fmt.Println(e.Cause())
	if e == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(e.Error(), "hello") {
		t.Fatalf("missing custom message, got: %q", e.Error())
	}
	if e.Cause() == nil {
		t.Fatalf("Cause() should still be non-nil error wrapping location+message")
	}
}

// TestImplementsCause 接口实现性检查
func TestImplementsCause(t *testing.T) {
	var c Cause = New(errors.New("x"), "y")
	if c == nil {
		t.Fatalf("expected Cause to be implemented")
	}
}
