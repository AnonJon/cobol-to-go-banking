package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkJSON measures JSON response serialization throughput.
// This replaces the Liberty JAX-RS JSON serialization layer and the
// z/OS Connect JSON-to-COBOL transformation.
func BenchmarkJSON(b *testing.B) {
	payload := map[string]interface{}{
		"accountNumber":    "00000001",
		"sortCode":         "987654",
		"customerNumber":   "0000000001",
		"accountType":      "CURRENT",
		"interestRate":     "2.25",
		"availableBalance": "1234.56",
		"actualBalance":    "1234.56",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		JSON(w, http.StatusOK, payload)
	}
}

// BenchmarkJSON_Parallel measures concurrent JSON serialization.
func BenchmarkJSON_Parallel(b *testing.B) {
	payload := map[string]interface{}{
		"accountNumber":  "00000001",
		"availableBalance": "1234.56",
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			JSON(w, http.StatusOK, payload)
		}
	})
}
