package main

import (
	"errors"
	"github.com/lithammer/shortuuid"
	"sync"
)

func (engin *Engin) createUser(reqUser any) (*User, error) {
	if _, ok := engin.findUser(reqUser); ok {
		return nil, errors.New("user already exists")
	}

	user := &User{
		engin:   engin,
		meta:    reqUser,
		budget:  0,
		targets: make([]*Target, 0),
		link:    &sync.Map{},
	}

	// юзер
	user.link.Store(shortuuid.New(), transaction{
		amount: 0,
		dayKey: engin.now().Format("2006-01-02"),
	})

	engin.accounts = append(engin.accounts, user)

	return user, nil
}

func (engin *Engin) createTarget(reqUser any, reqTarget any) (*Target, error) {
	user, _, _, err := engin.resolve(reqUser, nil, nil)
	if err != nil {
		return nil, err
	}

	if _, ok := user.findTarget(reqTarget); ok {
		return nil, errors.New("target already exists")
	}

	target := &Target{
		user:       user,
		meta:       reqTarget,
		budget:     0,
		subTargets: make([]*SubTarget, 0),
		link:       &sync.Map{},
	}

	target.link.Store(shortuuid.New(), transaction{
		amount: 0,
		dayKey: engin.now().Format("2006-01-02"),
	})

	user.targets = append(user.targets, target)

	return target, nil
}

func (engin *Engin) createSubTarget(reqUser any, reqTarget any, reqSubTarget any, reqLimit float64) (*SubTarget, error) {
	_, target, _, err := engin.resolve(reqUser, reqTarget, nil)
	if err != nil {
		return nil, err
	}

	if _, ok := target.finSubTarget(reqSubTarget); ok {
		return nil, errors.New("target already exists")
	}

	subTarget := &SubTarget{
		target:      target,
		meta:        reqSubTarget,
		dailyBudget: reqLimit,
		link:        &sync.Map{},
	}

	subTarget.link.Store(shortuuid.New(), transaction{
		amount: 0,
		dayKey: engin.now().Format("2006-01-02"),
	})

	target.subTargets = append(target.subTargets, subTarget)

	return subTarget, nil
}
