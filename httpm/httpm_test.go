package httpm

import "testing"

func TestNewResultOk_DefaultMsg(t *testing.T) {
	data := map[string]any{"k": "v"}
	r := NewResultOk(data)
	if r == nil {
		t.Fatalf("NewResultOk returned nil")
	}
	if r.Code != 0 {
		t.Fatalf("Code want 0, got %d", r.Code)
	}
	if r.Msg != "success" {
		t.Fatalf("Msg want 'success', got %q", r.Msg)
	}
	if r.Data == nil {
		t.Fatalf("Data should not be nil")
	}
}

func TestNewResultOk_CustomMsg(t *testing.T) {
	data := []int{1, 2, 3}
	r := NewResultOk(data, "ok")
	if r.Msg != "ok" {
		t.Fatalf("Msg want 'ok', got %q", r.Msg)
	}
	if r.Code != 0 {
		t.Fatalf("Code want 0, got %d", r.Code)
	}
	if _, ok := r.Data.([]int); !ok {
		t.Fatalf("Data type want []int")
	}
}

func TestNewResultErr_DefaultCode(t *testing.T) {
	r := NewResultErr("err")
	if r == nil {
		t.Fatalf("NewResultErr returned nil")
	}
	if r.Code != 1 {
		t.Fatalf("Code want 1, got %d", r.Code)
	}
	if r.Msg != "err" {
		t.Fatalf("Msg want 'err', got %q", r.Msg)
	}
}

func TestNewResultErr_CustomCode(t *testing.T) {
	r := NewResultErr("bad request", 4001)
	if r.Code != 4001 {
		t.Fatalf("Code want 4001, got %d", r.Code)
	}
	if r.Msg != "bad request" {
		t.Fatalf("Msg want 'bad request', got %q", r.Msg)
	}
}
