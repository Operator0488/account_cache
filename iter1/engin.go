package main

import (
	"errors"
	"github.com/lithammer/shortuuid"
)

func (engin *Engin) chageLimit(req *request) response {
	_, _, subtarget, err := engin.resolve(req.User, req.Target, req.SubTarget)
	if err != nil {
		return response{Error: err}
	}
	subtarget.dailyBudget = req.Score
	return response{Error: nil}
}

func (engin *Engin) rollback(req *request) response {
	transactionID := req.Transaction

	operation, ok := engin.transactions[transactionID]
	if !ok {
		return response{Error: errors.New("transaction not found")}
	}

	if operation.userLink != nil {
		operation.userLink.Delete(transactionID)
	}
	if operation.targetLink != nil {
		operation.targetLink.Delete(transactionID)
	}
	if operation.subTargetLink != nil {
		operation.subTargetLink.Delete(transactionID)
	}

	delete(engin.transactions, transactionID)

	return response{Error: nil}
}

func (engin *Engin) balance(req *request) (float64, error) {

	switch req.Level {
	case userLevel:
		user, _, _, err := engin.resolve(req.User, nil, nil)
		if err != nil {
			return 0, err
		}

		return accounting(user.link), nil

	case targetLevel:
		_, target, _, err := engin.resolve(req.User, req.Target, nil)
		if err != nil {
			return 0, err
		}

		return accounting(target.link), nil

	case subTargetLevel:
		_, _, subtarget, err := engin.resolve(req.User, req.Target, req.SubTarget)
		if err != nil {
			return 0, err
		}

		return subTargetDailyBalance(subtarget.link, subtarget.dailyBudget, engin.now().Format("2006-01-02")), nil

	}

	return 0, errors.New("unknown balance level")
}

func (engin *Engin) setCurrentUserBalance(reqUser any, score float64) (string, error) {
	user, _, _, err := engin.resolve(reqUser, nil, nil)
	if err != nil {
		return "", err
	}

	user.link.Range(func(key, _ any) bool {
		user.link.Delete(key)
		return true
	})

	txID := shortuuid.New()
	user.link.Store(txID, transaction{
		amount: score,
		dayKey: engin.now().Format("2006-01-02"),
	})

	engin.transactions[txID] = operations{
		userLink: user.link,
	}

	return txID, nil
}

func (engin *Engin) setCurrentTargetBalance(reqUser, reqTarget any, score float64) (string, error) {
	_, target, _, err := engin.resolve(reqUser, reqTarget, nil)
	if err != nil {
		return "", err
	}

	target.link.Range(func(key, _ any) bool {
		target.link.Delete(key)
		return true
	})

	txID := shortuuid.New()
	target.link.Store(txID, transaction{
		amount: score,
		dayKey: engin.now().Format("2006-01-02"),
	})

	engin.transactions[txID] = operations{
		targetLink: target.link,
	}

	return txID, nil
}

func (engin *Engin) setCurrentSubTargetBalance(reqUser, reqTarget, reqSubTarget any, score float64) (string, error) {
	_, _, subTarget, err := engin.resolve(reqUser, reqTarget, reqSubTarget)
	if err != nil {
		return "", err
	}

	subTarget.link.Range(func(key, _ any) bool {
		subTarget.link.Delete(key)
		return true
	})

	txID := shortuuid.New()
	subTarget.link.Store(txID, transaction{
		amount: score,
		dayKey: engin.now().Format("2006-01-02"),
	})

	engin.transactions[txID] = operations{
		subTargetLink: subTarget.link,
	}

	return txID, nil
}
