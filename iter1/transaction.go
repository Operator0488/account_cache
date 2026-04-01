package main

import (
	"errors"
	"github.com/lithammer/shortuuid"
)

func (engin *Engin) transactionUser(reqUser any, reqScore float64) (string, error) {
	user, _, _, err := engin.resolve(reqUser, nil, nil)
	if err != nil {
		return "", err
	}

	bal := accounting(user.link)

	if reqScore < 0 && bal < -reqScore {
		return "", errors.New("not enough balance on wallet")
	}

	txID := shortuuid.New()
	tx := transaction{
		amount: reqScore,
		dayKey: user.engin.now().Format("2006-01-02"),
	}

	user.link.Store(txID, tx)

	engin.transactions[txID] = operations{
		userLink: user.link,
	}

	return txID, nil
}

func (engin *Engin) transactionTarget(reqUser any, reqTarget any, reqScore float64) (string, error) {
	user, target, _, err := engin.resolve(reqUser, reqTarget, nil)
	if err != nil {
		return "", err
	}

	balUser := accounting(user.link)
	if reqScore < 0 && balUser < -reqScore {
		return "", errors.New("not enough balance on wallet")
	}

	balTar := accounting(target.link)
	if reqScore < 0 && balTar < -reqScore {
		return "", errors.New("not enough balance on target")
	}

	txID := shortuuid.New()
	tx := transaction{
		amount: reqScore,
		dayKey: user.engin.now().Format("2006-01-02"),
	}

	if reqScore <= 0 {
		user.link.Store(txID, tx)
		target.link.Store(txID, tx)
		engin.transactions[txID] = operations{
			userLink:   user.link,
			targetLink: target.link,
		}
	} else {
		target.link.Store(txID, tx)
		engin.transactions[txID] = operations{
			targetLink: target.link,
		}
	}

	return txID, nil
}

func (engin *Engin) transactionSubTarget(reqUser any, reqTarget any, reqSubTarget any, reqScore float64) (string, error) {
	user, target, subTarget, err := engin.resolve(reqUser, reqTarget, reqSubTarget)
	if err != nil {
		return "", err
	}

	balUser := accounting(user.link)
	if reqScore < 0 && balUser < -reqScore {
		return "", errors.New("not enough balance on wallet")
	}

	balTar := accounting(target.link)
	if reqScore < 0 && balTar < -reqScore {
		return "", errors.New("not enough balance on target")
	}

	balSub := subTargetDailyBalance(subTarget.link, subTarget.dailyBudget, engin.now().Format("2006-01-02"))
	if reqScore < 0 && balSub < -reqScore {
		return "", errors.New("not enough balance on subtarget today")
	}

	txID := shortuuid.New()
	tx := transaction{
		amount: reqScore,
		dayKey: engin.now().Format("2006-01-02"),
	}

	if reqScore <= 0 {
		user.link.Store(txID, tx)
		target.link.Store(txID, tx)
		subTarget.link.Store(txID, tx)
		engin.transactions[txID] = operations{
			userLink:      user.link,
			targetLink:    target.link,
			subTargetLink: subTarget.link,
		}
	} else {
		subTarget.link.Store(txID, tx)
		engin.transactions[txID] = operations{
			targetLink: target.link,
		}
	}

	return txID, nil
}
