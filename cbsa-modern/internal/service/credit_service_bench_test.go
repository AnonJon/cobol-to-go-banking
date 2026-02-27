package service

import (
	"context"
	"testing"
)

// BenchmarkCreditCheck measures throughput of the credit scoring fan-out.
// This replaces 5 separate CICS RUN TRANSID + FETCH ANY + GET CONTAINER
// calls in the original COBOL (CRDTAGY1-5.cbl). The goroutine-based
// implementation eliminates all CICS container marshalling overhead.
func BenchmarkCreditCheck(b *testing.B) {
	svc := NewCreditService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := svc.CheckCredit(ctx, "0000000001")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCreditCheck_Parallel measures concurrent credit check throughput.
// In CICS, each credit check consumes a task (typically capped at 500).
// With goroutines, we can run thousands concurrently.
func BenchmarkCreditCheck_Parallel(b *testing.B) {
	svc := NewCreditService()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := svc.CheckCredit(context.Background(), "0000000001")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
