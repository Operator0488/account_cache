package main

import (
	"sync"
)

type User struct {
	engin   *Engin
	meta    any
	budget  float64
	targets []*Target
	link    *sync.Map
}

type Target struct {
	user       *User
	meta       any
	budget     float64
	subTargets []*SubTarget
	link       *sync.Map
}

type SubTarget struct {
	target      *Target
	meta        any
	dailyBudget float64
	link        *sync.Map
}

type transaction struct {
	amount float64
	dayKey string
}

type task struct {
	action  string
	request *request
	result  chan response
}

type request struct {
	User        any
	Target      any
	SubTarget   any
	Score       float64
	Transaction string
	Level       level
}

type operations struct {
	userLink      *sync.Map
	targetLink    *sync.Map
	subTargetLink *sync.Map
}

type response struct {
	Meta  any
	Error error
}

type level string

const (
	userLevel      level = "user"
	targetLevel    level = "target"
	subTargetLevel level = "subTarget"
)

func (user *User) HistoryLength() (length int64) {
	user.link.Range(func(key, value interface{}) bool {
		length += 1

		return true
	})

	return
}

func (target *Target) HistoryLength() (length int64) {
	target.link.Range(func(key, value interface{}) bool {
		length += 1

		return true
	})

	return
}

func (subTarget *SubTarget) HistoryLength() (length int64) {
	subTarget.link.Range(func(key, value interface{}) bool {
		length += 1

		return true
	})

	return
}
