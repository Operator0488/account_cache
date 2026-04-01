package main

import (
	"testing"
)

func BenchmarkGetBalance(b *testing.B) {
	account := Init(100000)

	b.SetBytes(int64(10000))
	b.ResetTimer()

	wallet, _ := account.CreateUser("test_user")
	if wallet != nil {
		wallet.Transaction(0)
	}

	for i := 0; i < b.N; i++ {
		wallet.Transaction(float64(i))
	}

	for i := 0; i < b.N; i++ {
		_, _ = wallet.SyncBalance()
	}
}
