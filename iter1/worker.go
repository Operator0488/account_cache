package main

import "errors"

// worker - queue worker.
func (engin *Engin) worker() {
	for t := range engin.queue {
		switch t.action {
		case "create":
			t.result <- engin.create(t.request)

		case "rollback":
			t.result <- engin.rollback(t.request)

		case "transaction":
			t.result <- engin.transaction(t.request)

		case "find":
			t.result <- engin.find(t.request)

		case "balance":
			value, err := engin.balance(t.request)
			t.result <- response{
				Meta:  value,
				Error: err,
			}

		case "change_limit":
			t.result <- engin.chageLimit(t.request)

		case "set_current_balance":
			t.result <- engin.setCurrentBalance(t.request)

		}
	}
}

func (engin *Engin) setCurrentBalance(req *request) response {
	switch req.Level {
	case userLevel:
		txID, err := engin.setCurrentUserBalance(req.User, req.Score)
		return response{Meta: txID, Error: err}

	case targetLevel:
		txID, err := engin.setCurrentTargetBalance(req.User, req.Target, req.Score)
		return response{Meta: txID, Error: err}

	case subTargetLevel:
		txID, err := engin.setCurrentSubTargetBalance(req.User, req.Target, req.SubTarget, req.Score)
		return response{Meta: txID, Error: err}

	default:
		return response{Meta: nil, Error: errors.New("unknown level")}
	}
}

func (engin *Engin) find(req *request) response {
	switch req.Level {
	case userLevel:
		u, _, _, err := engin.resolve(req.User, req.Target, req.SubTarget)
		return response{Meta: u, Error: err}
	case targetLevel:
		_, t, _, err := engin.resolve(req.User, req.Target, req.SubTarget)
		return response{Meta: t, Error: err}
	case subTargetLevel:
		_, _, s, err := engin.resolve(req.User, req.Target, req.SubTarget)
		return response{Meta: s, Error: err}
	default:
		return response{Meta: 0, Error: errors.New("unknown level")}
	}

}

func (engin *Engin) transaction(req *request) response {
	switch req.Level {
	case userLevel:
		tr, err := engin.transactionUser(req.User, req.Score)
		return response{Meta: tr, Error: err}
	case targetLevel:
		tr, err := engin.transactionTarget(req.User, req.Target, req.Score)
		return response{Meta: tr, Error: err}
	case subTargetLevel:
		tr, err := engin.transactionSubTarget(req.User, req.Target, req.SubTarget, req.Score)
		return response{Meta: tr, Error: err}
	default:
		return response{Meta: 0, Error: errors.New("unknown level")}
	}
}

func (engin *Engin) create(req *request) response {
	switch req.Level {
	case userLevel:
		u, err := engin.createUser(req.User)
		return response{Meta: u, Error: err}
	case targetLevel:
		t, err := engin.createTarget(req.User, req.Target)
		return response{Meta: t, Error: err}
	case subTargetLevel:
		s, err := engin.createSubTarget(req.User, req.Target, req.SubTarget, req.Score)
		return response{Meta: s, Error: err}
	default:
		return response{Meta: nil, Error: errors.New("unknown level")}
	}
}
