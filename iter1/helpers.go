package main

import (
	"sync"
)

func accounting(link *sync.Map) float64 {
	var balance float64

	link.Range(func(_, value any) bool {
		balance += value.(transaction).amount
		return true
	})

	return balance
}

func subTargetDailyBalance(link *sync.Map, dailyBudget float64, currentDay string) float64 {
	balance := dailyBudget

	link.Range(func(_, value any) bool {
		tx := value.(transaction)
		if tx.dayKey == currentDay {
			if tx.amount < 0 {
				balance += tx.amount
			}
		}
		return true
	})

	return balance
}
