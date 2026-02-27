package service

import (
	"context"
	"testing"
	"time"
)

func TestCreditService_CheckCredit(t *testing.T) {
	svc := NewCreditService()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	score, err := svc.CheckCredit(ctx, "0000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if score < 1 || score > 999 {
		t.Errorf("score %d out of range [1, 999]", score)
	}
}

func TestCreditService_CheckCredit_Cancelled(t *testing.T) {
	svc := NewCreditService()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.CheckCredit(ctx, "0000000001")
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestCreditService_CheckCredit_Consistency(t *testing.T) {
	svc := NewCreditService()
	ctx := context.Background()

	for i := 0; i < 20; i++ {
		score, err := svc.CheckCredit(ctx, "0000000001")
		if err != nil {
			t.Fatalf("run %d: unexpected error: %v", i, err)
		}
		if score < 1 || score > 999 {
			t.Errorf("run %d: score %d out of range [1, 999]", i, score)
		}
	}
}
