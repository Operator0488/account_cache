package main

import (
	"fmt"
	"log"
)

func main() {

	var bufferSize int64 = 100000

	account := Init(bufferSize)

	user, err := account.CreateUser("user")
	log.Println(err)
	target, err := user.CreateTarget("target")
	log.Println(err)
	_, err = target.CreateSubTarget("subtarget", 10)
	log.Println(err)

	account.User("user").Transaction(1000)
	account.User("user").Target("target").SubTarget("subtarget").Transaction(100)
	account.User("user").Target("target").Transaction(1000)

	w, err := account.User("user").SyncBalance()
	log.Println(err)
	w2, err := account.User("user").Target("target").CheckLimit()
	log.Println(err)
	w3, err := account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	account.User("user").Transaction(-50)
	account.User("user").Target("target").Transaction(-20)
	_, err = account.User("user").Target("target").SubTarget("subtarget").Transaction(-100)
	log.Println(err)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)
	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	account.User("user").Target("target").SubTarget("subtarget").Transaction(-5)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)
	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	account.User("user").Target("target").SubTarget("subtarget").Transaction(-4)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)
	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	account.User("user").Target("target").SubTarget("subtarget").Transaction(-5)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)
	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	account.User("user").SetCurrentBalance(1000)
	account.User("user").Target("target").SubTarget("subtarget").SetCurrentBalance(5)
	account.User("user").Target("target").SetCurrentBalance(1000)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)
	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	account.User("user").Target("target").SubTarget("subtarget").Transaction(-5)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)
	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)
}

func Case2() {
	var bufferSize int64 = 100000

	account := Init(bufferSize)

	user, err := account.CreateUser("user")
	log.Println(err)
	log.Println(account.User("user"))

	target, err := user.CreateTarget("target")
	log.Println(err)
	log.Println(account.User("user").Target("target"))
	_, _ = target.CreateSubTarget("subtarget", 10)
	w3, err := account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(err, w3)

	account.User("user").Transaction(100)
	account.User("user").Target("target").Transaction(100)
	_, err = account.User("user").Target("target").SubTarget("subtarget").Transaction(10)
	log.Println(err)
	_, err = account.User("user").Target("target").SubTarget("subtarget").Transaction(-1)
	log.Println(err)

	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(err, w3)
}

func Case1() {
	var bufferSize int64 = 100000

	account := Init(bufferSize)
	user, err := account.CreateUser("user")
	log.Println(err)
	log.Println(account.User("user"))

	target, err := user.CreateTarget("target")
	log.Println(err)
	log.Println(account.User("user").Target("target"))
	_, _ = target.CreateSubTarget("subtarget", 100)

	account.User("user").Transaction(100)
	account.User("user").Target("target").SubTarget("subtarget").Transaction(100)
	account.User("user").Target("target").Transaction(100)

	w, err := account.User("user").SyncBalance()
	log.Println(err)
	w2, err := account.User("user").Target("target").CheckLimit()
	log.Println(err)
	w3, err := account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	account.User("user").Transaction(-50)
	account.User("user").Target("target").Transaction(-20)
	account.User("user").Target("target").SubTarget("subtarget").Transaction(-5)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)
	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	txid, err := account.User("user").Target("target").SubTarget("subtarget").Transaction(-25)
	log.Println(err)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)

	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	err = account.Rollback(txid)
	log.Println(err)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)

	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	txid, err = account.User("user").Target("target").Transaction(-25)
	log.Println(err)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)

	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	err = account.Rollback(txid)
	log.Println(err)

	w, err = account.User("user").SyncBalance()
	log.Println(err)
	w2, err = account.User("user").Target("target").CheckLimit()
	log.Println(err)

	w3, err = account.User("user").Target("target").SubTarget("subtarget").CheckLimit()
	log.Println(err)

	fmt.Println(w, w2, w3)

	_, err = account.User("user1").Transaction(100)
	log.Println(err)
	_, err = account.User("user").Target("target1").SubTarget("subtarget").Transaction(100)
	log.Println(err)
	_, err = account.User("user1").Target("target").Transaction(100)
	log.Println(err)
}

func Chernovik() {
	//var bufferSize int64 = 100000

	//off1 := SubTargetInput{
	//	Meta:        "finSubTarget-1",
	//	MaxBidRate:  120,
	//	DailyBudget: 1000,
	//}
	//
	//off2 := SubTargetInput{
	//	Meta:        "finSubTarget-2",
	//	MaxBidRate:  120,
	//	DailyBudget: 1000,
	//}
	//
	//camp1 := CampaignInput{
	//	Meta:   "campaign-1",
	//	Budget: 4000,
	//	SubTargets: []SubTargetInput{
	//		off1,
	//		off2,
	//	},
	//}
	//
	//wal := CreateWalletInput{
	//	Meta:   "user-1",
	//	Budget: 10000,
	//	Campaigns: []CampaignInput{
	//		camp1,
	//	},
	//}

	//account := Init(bufferSize) // buffer queue size.
	//
	//account.CreateWallet("user")
	//
	//if user != nil {
	//	tr, err := user.Transaction(100)
	//}
	//
	//if err := account.CreateUser(wal); err == nil {
	//
	//	str, err := account.Transaction("user-1", "campaign-1", "finSubTarget-1", -100) // add first default balance.
	//
	//	log.Println(str)
	//	log.Println(err)
	//
	//	user, err := account.CreateWallet(wal)
	//
	//	trx, err := account.User("score int").Transaction(100)
	//	account.User("string").Target("string").SubTarget("string").Transaction("score int")
	//	account.User("string").Target("string").Transaction("score int")
	//}
	//
	//w, err := account.SyncBalanceWallet("user-1")
	//log.Println(w)
	//log.Println(err)
	//
	//w2, err := account.SyncBalanceCampaign("user-1", "campaign-1")
	//log.Println(w2)
	//log.Println(err)
	//
	//
	//
	//account.User().Transaction("score int")
	//account.User("string").Campaign("string").Offer("string").Transaction("score int")
	//account.User("string").Campaign("string").Transaction("score int")
}
