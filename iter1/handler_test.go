package main

import (
	"sync"
	"testing"
	"time"
)

func newTestEngine(t *testing.T) *Engin {
	t.Helper()

	engin := Init(32)
	engin.now = func() time.Time {
		return time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	}

	t.Cleanup(func() {
		close(engin.queue)
	})

	return engin
}

func createWallets(t *testing.T, engin *Engin, subLimit float64) (*User, *Target, *SubTarget) {
	t.Helper()

	user, err := engin.CreateUser("user-1")
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	target, err := user.CreateTarget("target-1")
	if err != nil {
		t.Fatalf("CreateTarget() error = %v", err)
	}

	subTarget, err := target.CreateSubTarget("sub-1", subLimit)
	if err != nil {
		t.Fatalf("CreateSubTarget() error = %v", err)
	}

	return user, target, subTarget
}

func userBalance(t *testing.T, user *User) float64 {
	t.Helper()

	got, err := user.SyncBalance()
	if err != nil {
		t.Fatalf("SyncBalance() error = %v", err)
	}

	return got
}

func targetBalance(t *testing.T, target *Target) float64 {
	t.Helper()

	got, err := target.CheckLimit()
	if err != nil {
		t.Fatalf("Target.CheckLimit() error = %v", err)
	}

	return got
}

func subTargetLimit(t *testing.T, subTarget *SubTarget) float64 {
	t.Helper()

	got, err := subTarget.CheckLimit()
	if err != nil {
		t.Fatalf("SubTarget.CheckLimit() error = %v", err)
	}

	return got
}

// Create Wallets and
func Test_Create(t *testing.T) {
	engin := newTestEngine(t)
	user, target, subTarget := createWallets(t, engin, 100)

	if engin.User("user-1") == nil {
		t.Fatal("engin.User(user-1) = nil, want user")
	}

	foundTarget := engin.User("user-1").Target("target-1")
	if foundTarget == nil {
		t.Fatal("user.Target(target-1) = nil, want target")
	}

	foundSubTarget := foundTarget.SubTarget("sub-1")
	if foundSubTarget == nil {
		t.Fatal("target.SubTarget(sub-1) = nil, want subTarget")
	}

	got, err := user.SyncBalance()
	if err != nil {
		t.Fatalf("SyncBalance() error = %v", err)
	}

	if got != 0 {
		t.Fatalf("user balance = %v, want 0", got)
	}

	got, err = target.CheckLimit()
	if err != nil {
		t.Fatalf("Target.CheckLimit() error = %v", err)
	}

	if got != 0 {
		t.Fatalf("target balance = %v, want 0", got)
	}

	got, err = subTarget.CheckLimit()
	if err != nil {
		t.Fatalf("SubTarget.CheckLimit() error = %v", err)
	}

	if got != 100 {
		t.Fatalf("subTarget limit = %v, want 100", got)
	}
}

func Test_CreateDuplicates(t *testing.T) {
	engin := newTestEngine(t)
	user, target, _ := createWallets(t, engin, 100)

	if _, err := engin.CreateUser("user-1"); err == nil {
		t.Fatal("CreateUser duplicate: expected error, got nil")
	}

	if _, err := user.CreateTarget("target-1"); err == nil {
		t.Fatal("CreateTarget duplicate: expected error, got nil")
	}

	if _, err := target.CreateSubTarget("sub-1", 100); err == nil {
		t.Fatal("CreateSubTarget duplicate: expected error, got nil")
	}
}

func TestUser_SetCurrentBalanceTransactionRollback(t *testing.T) {
	engin := newTestEngine(t)
	user, _, _ := createWallets(t, engin, 100)

	if _, err := user.SetCurrentBalance(100); err != nil {
		t.Fatalf("SetCurrentBalance() error = %v", err)
	}

	got, err := user.SyncBalance()
	if err != nil {
		t.Fatalf("SyncBalance() error = %v", err)
	}

	if got != 100 {
		t.Fatalf("user balance = %v, want 100", got)
	}

	txID, err := user.Transaction(-30)
	if err != nil {
		t.Fatalf("user.Transaction() error = %v", err)
	}

	got, err = user.SyncBalance()
	if err != nil {
		t.Fatalf("SyncBalance() error = %v", err)
	}

	if got != 70 {
		t.Fatalf("user balance = %v, want 70", got)
	}

	txID, err = user.Transaction(20)
	if err != nil {
		t.Fatalf("user.Transaction() error = %v", err)
	}

	got, err = user.SyncBalance()
	if err != nil {
		t.Fatalf("SyncBalance() error = %v", err)
	}

	if got != 90 {
		t.Fatalf("user balance = %v, want 70", got)
	}

	if err := engin.Rollback(txID); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}

	got, err = user.SyncBalance()
	if err != nil {
		t.Fatalf("SyncBalance() error = %v", err)
	}

	if got != 70 {
		t.Fatalf("user balance = %v, want 70", got)
	}

}

func TestTarget_TransactionDebetCredit(t *testing.T) {
	engin := newTestEngine(t)
	user, target, _ := createWallets(t, engin, 100)

	if _, err := user.SetCurrentBalance(100); err != nil {
		t.Fatalf("user.SetCurrentBalance() error = %v", err)
	}

	if _, err := target.SetCurrentBalance(60); err != nil {
		t.Fatalf("target.SetCurrentBalance() error = %v", err)
	}

	if got := userBalance(t, user); got != 100 {
		t.Fatalf("initial user balance = %v, want 100", got)
	}
	if got := targetBalance(t, target); got != 60 {
		t.Fatalf("initial target balance = %v, want 60", got)
	}

	txID, err := target.Transaction(-20)
	if err != nil {
		t.Fatalf("target.Transaction(-20) error = %v", err)
	}

	if got := userBalance(t, user); got != 80 {
		t.Fatalf("user balance after target credit = %v, want 80", got)
	}
	if got := targetBalance(t, target); got != 40 {
		t.Fatalf("target balance after credit = %v, want 40", got)
	}

	if err := engin.Rollback(txID); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}

	if got := userBalance(t, user); got != 100 {
		t.Fatalf("user balance after rollback = %v, want 100", got)
	}
	if got := targetBalance(t, target); got != 60 {
		t.Fatalf("target balance after rollback = %v, want 60", got)
	}

	txID, err = target.Transaction(15)
	if err != nil {
		t.Fatalf("target.Transaction(15) error = %v", err)
	}

	if got := userBalance(t, user); got != 100 {
		t.Fatalf("user balance after target debet = %v, want 100", got)
	}
	if got := targetBalance(t, target); got != 75 {
		t.Fatalf("target balance after debet = %v, want 75", got)
	}

	if err := engin.Rollback(txID); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}

	if got := userBalance(t, user); got != 100 {
		t.Fatalf("user balance after debet rollback = %v, want 100", got)
	}

	if got := targetBalance(t, target); got != 60 {
		t.Fatalf("target balance after debet rollback = %v, want 60", got)
	}
}

func TestSubTarget_CreditLimitRollbackResetChangeLimit(t *testing.T) {
	engin := newTestEngine(t)
	user, target, subTarget := createWallets(t, engin, 50)

	if _, err := user.SetCurrentBalance(100); err != nil {
		t.Fatalf("user.SetCurrentBalance() error = %v", err)
	}
	if _, err := target.SetCurrentBalance(80); err != nil {
		t.Fatalf("target.SetCurrentBalance() error = %v", err)
	}

	if got := subTargetLimit(t, subTarget); got != 50 {
		t.Fatalf("initial subTarget limit = %v, want 50", got)
	}

	txID, err := subTarget.Transaction(-20)
	if err != nil {
		t.Fatalf("subTarget.Transaction(-20) error = %v", err)
	}

	if got := userBalance(t, user); got != 80 {
		t.Fatalf("user balance after subTarget credit = %v, want 80", got)
	}
	if got := targetBalance(t, target); got != 60 {
		t.Fatalf("target balance after subTarget credit = %v, want 60", got)
	}
	if got := subTargetLimit(t, subTarget); got != 30 {
		t.Fatalf("subTarget limit after credit = %v, want 30", got)
	}

	if err := engin.Rollback(txID); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}

	if got := userBalance(t, user); got != 100 {
		t.Fatalf("user balance after rollback = %v, want 100", got)
	}
	if got := targetBalance(t, target); got != 80 {
		t.Fatalf("target balance after rollback = %v, want 80", got)
	}
	if got := subTargetLimit(t, subTarget); got != 50 {
		t.Fatalf("subTarget limit after rollback = %v, want 50", got)
	}

	if err := subTarget.ChangeLimit(70); err != nil {
		t.Fatalf("ChangeLimit() error = %v", err)
	}

	if got := subTargetLimit(t, subTarget); got != 70 {
		t.Fatalf("subTarget limit after ChangeLimit = %v, want 70", got)
	}

	if _, err := subTarget.Transaction(-10); err != nil {
		t.Fatalf("subTarget.Transaction() error = %v", err)
	}

	if got := subTargetLimit(t, subTarget); got != 60 {
		t.Fatalf("subTarget limit after second credit = %v, want 60", got)
	}

	if _, err := subTarget.ResetBalance(); err != nil {
		t.Fatalf("ResetBalance() error = %v", err)
	}

	if got := subTargetLimit(t, subTarget); got != 70 {
		t.Fatalf("subTarget limit after reset = %v, want 70", got)
	}
}

func TestSubTarget_TransactionPositiv(t *testing.T) {
	engin := newTestEngine(t)
	_, _, subTarget := createWallets(t, engin, 50)

	if _, err := subTarget.Transaction(10); err == nil {
		t.Fatal("subTarget.Transaction(): expected error, got nil")
	}
}

func TestEngine_Transaction(t *testing.T) {
	engin := newTestEngine(t)
	user, target, subTarget := createWallets(t, engin, 40)

	if _, err := user.SetCurrentBalance(100); err != nil {
		t.Fatalf("user.SetCurrentBalance() error = %v", err)
	}
	if _, err := target.SetCurrentBalance(70); err != nil {
		t.Fatalf("target.SetCurrentBalance() error = %v", err)
	}

	txID, err := engin.TransactionNegative("user-1", "target-1", "sub-1", -15)
	if err != nil {
		t.Fatalf("engin.Transaction() error = %v", err)
	}

	if got := userBalance(t, user); got != 85 {
		t.Fatalf("user balance after engin.Transaction = %v, want 85", got)
	}
	if got := targetBalance(t, target); got != 55 {
		t.Fatalf("target balance after engin.Transaction = %v, want 55", got)
	}
	if got := subTargetLimit(t, subTarget); got != 25 {
		t.Fatalf("subTarget limit after engin.Transaction = %v, want 25", got)
	}

	if err := engin.Rollback(txID); err != nil {
		t.Fatalf("Rollback() error = %v", err)
	}

	if got := userBalance(t, user); got != 100 {
		t.Fatalf("user balance after rollback = %v, want 100", got)
	}
	if got := targetBalance(t, target); got != 70 {
		t.Fatalf("target balance after rollback = %v, want 70", got)
	}
	if got := subTargetLimit(t, subTarget); got != 40 {
		t.Fatalf("subTarget limit after rollback = %v, want 40", got)
	}

	if _, err = engin.TransactionNegative("user-1", "target-1", "sub-1", 15); err == nil {
		t.Fatalf("expected error, but error = nil")
	}
}

func Test_BalanceLimit(t *testing.T) {
	engin := newTestEngine(t)
	user, target, subTarget := createWallets(t, engin, 10)

	if _, err := user.SetCurrentBalance(5); err != nil {
		t.Fatalf("user.SetCurrentBalance() error = %v", err)
	}
	if _, err := target.SetCurrentBalance(4); err != nil {
		t.Fatalf("target.SetCurrentBalance() error = %v", err)
	}

	if _, err := user.Transaction(-10); err == nil {
		t.Fatal("user.Transaction(): expected error, but  nil")
	}

	if _, err := target.Transaction(-10); err == nil {
		t.Fatal("target.Transaction(): expected error, but  nil")
	}

	if _, err := subTarget.Transaction(-11); err == nil {
		t.Fatal("subTarget.Transaction(): expected error, but nil")
	}
}

func TestSubTarget_ConcurrentTransactions(t *testing.T) {
	engin := Init(1000)
	engin.now = func() time.Time {
		return time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	}

	t.Cleanup(func() {
		close(engin.queue)
	})

	user, err := engin.CreateUser("user-1")
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	target, err := user.CreateTarget("target-1")
	if err != nil {
		t.Fatalf("CreateTarget() error = %v", err)
	}

	subTarget, err := target.CreateSubTarget("sub-1", 10_000)
	if err != nil {
		t.Fatalf("CreateSubTarget() error = %v", err)
	}

	if _, err := user.SetCurrentBalance(10_000); err != nil {
		t.Fatalf("user.SetCurrentBalance() error = %v", err)
	}

	if _, err := target.SetCurrentBalance(10_000); err != nil {
		t.Fatalf("target.SetCurrentBalance() error = %v", err)
	}

	const (
		goroutines        = 50
		transactions      = 100
		transactionAmount = -1.0
	)

	start := make(chan struct{})
	errCh := make(chan error, goroutines*transactions)

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			<-start

			for j := 0; j < transactions; j++ {
				if _, err := subTarget.Transaction(transactionAmount); err != nil {
					errCh <- err
					return
				}
			}

		}()
	}

	close(start)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent transaction error = %v", err)
		}
	}

	userBalance, err := user.SyncBalance()
	if err != nil {
		t.Fatalf("user.SyncBalance() error = %v", err)
	}

	targetBalance, err := target.CheckLimit()
	if err != nil {
		t.Fatalf("target.CheckLimit() error = %v", err)
	}

	subTargetBalance, err := subTarget.CheckLimit()
	if err != nil {
		t.Fatalf("subTarget.CheckLimit() error = %v", err)
	}

	totalSpent := float64(goroutines * transactions)

	wantUser := 10_000 - totalSpent
	wantTarget := 10_000 - totalSpent
	wantSubTarget := 10_000 - totalSpent

	if userBalance != wantUser {
		t.Fatalf("user balance = %v, want %v", userBalance, wantUser)
	}

	if targetBalance != wantTarget {
		t.Fatalf("target balance = %v, want %v", targetBalance, wantTarget)
	}

	if subTargetBalance != wantSubTarget {
		t.Fatalf("subTarget balance = %v, want %v", subTargetBalance, wantSubTarget)
	}
}
