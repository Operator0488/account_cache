package main

import (
	"errors"
	"time"
)

// Engin - engine for caching user account balance.
type Engin struct {
	queue        chan task
	accounts     []*User
	transactions map[string]operations
	now          func() time.Time
}

// Init - Engin struct.
func Init(queueBufferSize int64) *Engin {
	engin := &Engin{
		queue:        make(chan task, queueBufferSize),
		accounts:     []*User{},
		transactions: map[string]operations{},
		now:          time.Now,
	}

	// run queue worker reader
	go engin.worker()

	return engin
}

func (engin *Engin) TransactionNegative(user any, target any, subTarget any, score float64) (string, error) {
	if subTarget == nil || user == nil || target == nil {
		return "", errors.New("fail request")
	}

	if score > 0 {
		return "", errors.New("transaction in engine can not be more than 0")
	}

	result := make(chan response)

	engin.queue <- task{
		action: "transaction",
		request: &request{
			User:      user,
			Target:    target,
			SubTarget: subTarget,
			Score:     score,
			Level:     subTargetLevel,
		},
		result: result,
	}

	res := <-result
	txID, ok := res.Meta.(string)
	if !ok {
		return "", errors.New("fail recording")
	}
	return txID, res.Error
}

func (user *User) Transaction(score float64) (string, error) {
	if user == nil {
		return "", errors.New("not found user")
	}

	result := make(chan response)

	user.engin.queue <- task{
		action: "transaction",
		request: &request{
			User:  user,
			Score: score,
			Level: userLevel,
		},
		result: result,
	}

	res := <-result
	txID, ok := res.Meta.(string)
	if !ok {
		return "", errors.New("fail recording")
	}
	return txID, res.Error
}

func (target *Target) Transaction(score float64) (string, error) {
	if target == nil {
		return "", errors.New("not found target")
	}

	result := make(chan response)

	target.user.engin.queue <- task{
		action: "transaction",
		request: &request{
			User:   target.user,
			Target: target,
			Score:  score,
			Level:  targetLevel,
		},
		result: result,
	}

	res := <-result
	txID, ok := res.Meta.(string)
	if !ok {
		return "", errors.New("fail recording")
	}
	return txID, res.Error
}

func (subTarget *SubTarget) Transaction(score float64) (string, error) {
	if subTarget == nil {
		return "", errors.New("not found subTarget")
	}
	if score > 0 {
		return "", errors.New("transaction subtarget can not be more than 0")
	}

	result := make(chan response)

	subTarget.target.user.engin.queue <- task{
		action: "transaction",
		request: &request{
			User:      subTarget.target.user,
			Target:    subTarget.target,
			SubTarget: subTarget,
			Score:     score,
			Level:     subTargetLevel,
		},
		result: result,
	}

	res := <-result
	txID, ok := res.Meta.(string)
	if !ok {
		return "", errors.New("fail recording")
	}
	return txID, res.Error
}

func (engin *Engin) CreateUser(user any) (*User, error) {
	result := make(chan response)

	engin.queue <- task{
		action: "create",
		request: &request{
			User:  user,
			Level: userLevel,
		},
		result: result,
	}

	res := <-result
	u, ok := res.Meta.(*User)
	if !ok {
		return nil, errors.New("fail create wallet")
	}

	return u, res.Error
}

func (user *User) CreateTarget(target any) (*Target, error) {
	if user == nil {
		return nil, errors.New("not found user")
	}

	result := make(chan response)

	user.engin.queue <- task{
		action: "create",
		request: &request{
			User:   user,
			Target: target,
			Level:  targetLevel,
		},
		result: result,
	}

	res := <-result
	t, ok := res.Meta.(*Target)
	if !ok {
		return nil, errors.New("fail create wallet")
	}

	return t, res.Error
}

func (target *Target) CreateSubTarget(subtarget any, limit float64) (*SubTarget, error) {
	if target == nil {
		return nil, errors.New("not found target")
	}

	result := make(chan response)

	target.user.engin.queue <- task{
		action: "create",
		request: &request{
			User:      target.user,
			Target:    target,
			SubTarget: subtarget,
			Score:     limit,
			Level:     subTargetLevel,
		},
		result: result,
	}

	res := <-result
	s, ok := res.Meta.(*SubTarget)
	if !ok {
		return nil, errors.New("fail create wallet")
	}

	return s, res.Error
}

func (engin *Engin) User(user any) *User {
	result := make(chan response)

	engin.queue <- task{
		action: "find",
		request: &request{
			User:  user,
			Level: userLevel,
		},
		result: result,
	}

	res := <-result
	u, ok := res.Meta.(*User)
	if !ok {
		return nil
	}

	return u
}

func (user *User) Target(target any) *Target {
	if user == nil {
		return nil
	}
	result := make(chan response)

	user.engin.queue <- task{
		action: "find",
		request: &request{
			User:   user,
			Target: target,
			Level:  targetLevel,
		},
		result: result,
	}

	res := <-result
	t, ok := res.Meta.(*Target)
	if !ok {
		return nil
	}

	return t
}

func (target *Target) SubTarget(subtarget any) *SubTarget {
	if target == nil {
		return nil
	}
	result := make(chan response)

	target.user.engin.queue <- task{
		action: "find",
		request: &request{
			User:      target.user,
			Target:    target,
			SubTarget: subtarget,
			Level:     subTargetLevel,
		},
		result: result,
	}

	res := <-result
	s, ok := res.Meta.(*SubTarget)
	if !ok {
		return nil
	}

	return s
}

func (engin *Engin) Rollback(transactionId string) error {
	result := make(chan response)

	engin.queue <- task{
		action: "rollback",
		request: &request{
			Transaction: transactionId,
		},
		result: result,
	}

	res := <-result
	return res.Error
}

func (user *User) SyncBalance() (float64, error) {
	if user == nil {
		return 0, errors.New("not found user")
	}

	result := make(chan response)

	user.engin.queue <- task{
		action: "balance",
		request: &request{
			Level: userLevel,
			User:  user,
		},
		result: result,
	}

	res := <-result
	balance, ok := res.Meta.(float64)
	if !ok {
		return 0, errors.New("fail sync balance")
	}

	return balance, res.Error
}

func (target *Target) CheckLimit() (float64, error) {
	if target == nil {
		return 0, errors.New("not found target")
	}

	result := make(chan response)

	target.user.engin.queue <- task{
		action: "balance",
		request: &request{
			Level:  targetLevel,
			User:   target.user,
			Target: target,
		},
		result: result,
	}

	res := <-result
	limit, ok := res.Meta.(float64)
	if !ok {
		return 0, errors.New("fail sync balance")
	}

	return limit, res.Error
}

func (subTarget *SubTarget) CheckLimit() (float64, error) {
	if subTarget == nil {
		return 0, errors.New("not found subTarget")
	}

	result := make(chan response)

	subTarget.target.user.engin.queue <- task{
		action: "balance",
		request: &request{
			Level:     subTargetLevel,
			User:      subTarget.target.user,
			Target:    subTarget.target,
			SubTarget: subTarget,
		},
		result: result,
	}

	res := <-result
	limit, ok := res.Meta.(float64)
	if !ok {
		return 0, errors.New("fail sync balance")
	}

	return limit, res.Error
}

func (subTarget *SubTarget) ChangeLimit(limit float64) error {
	if subTarget == nil {
		return errors.New("not found subTarget")
	}

	result := make(chan response)

	subTarget.target.user.engin.queue <- task{
		action: "change_limit",
		request: &request{
			User:      subTarget.target.user,
			Target:    subTarget.target,
			SubTarget: subTarget,
			Score:     limit,
			Level:     subTargetLevel,
		},
		result: result,
	}

	res := <-result

	return res.Error
}

func (user *User) SetCurrentBalance(score float64) (string, error) {
	if user == nil {
		return "", errors.New("not found user")
	}

	result := make(chan response)

	user.engin.queue <- task{
		action: "set_current_balance",
		request: &request{
			User:  user,
			Score: score,
			Level: userLevel,
		},
		result: result,
	}

	res := <-result
	txID, ok := res.Meta.(string)
	if !ok {
		return "", errors.New("fail set current balance")
	}

	return txID, res.Error
}

func (target *Target) SetCurrentBalance(score float64) (string, error) {
	if target == nil {
		return "", errors.New("not found target")
	}

	result := make(chan response)

	target.user.engin.queue <- task{
		action: "set_current_balance",
		request: &request{
			User:   target.user,
			Target: target,
			Score:  score,
			Level:  targetLevel,
		},
		result: result,
	}

	res := <-result
	txID, ok := res.Meta.(string)
	if !ok {
		return "", errors.New("fail set current balance")
	}

	return txID, res.Error
}

func (subTarget *SubTarget) ResetBalance() (string, error) {
	if subTarget == nil {
		return "", errors.New("not found subtarget")
	}

	result := make(chan response)

	subTarget.target.user.engin.queue <- task{
		action: "set_current_balance",
		request: &request{
			User:      subTarget.target.user,
			Target:    subTarget.target,
			SubTarget: subTarget,
			Score:     0,
			Level:     subTargetLevel,
		},
		result: result,
	}

	res := <-result
	txID, ok := res.Meta.(string)
	if !ok {
		return "", errors.New("fail set current balance")
	}

	return txID, res.Error
}

func (subTarget *SubTarget) SetCurrentBalance(score float64) (string, error) {
	if subTarget == nil {
		return "", errors.New("not found subtarget")
	}

	result := make(chan response)

	subTarget.target.user.engin.queue <- task{
		action: "set_current_balance",
		request: &request{
			User:      subTarget.target.user,
			Target:    subTarget.target,
			SubTarget: subTarget,
			Score:     score,
			Level:     subTargetLevel,
		},
		result: result,
	}

	res := <-result
	txID, ok := res.Meta.(string)
	if !ok {
		return "", errors.New("fail set current balance")
	}

	return txID, res.Error
}
