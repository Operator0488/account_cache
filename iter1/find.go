package main

import (
	"errors"
	"reflect"
)

func (engin *Engin) resolve(user, target, subtarget any) (*User, *Target, *SubTarget, error) {
	if user != nil {
		u, ok := engin.findUser(user)
		if !ok {
			return nil, nil, nil, errors.New("user not found")
		}

		if target != nil {
			c, ok := u.findTarget(target)
			if !ok {
				return nil, nil, nil, errors.New("target not found")
			}
			if subtarget != nil {
				o, ok := c.finSubTarget(subtarget)
				if !ok {
					return nil, nil, nil, errors.New("subTarget not found")
				}
				return u, c, o, nil
			}
			return u, c, nil, nil
		}
		return u, nil, nil, nil
	}

	return nil, nil, nil, errors.New("error in request all nil")
}

func (engin *Engin) findUser(meta any) (*User, bool) {
	if u, ok := meta.(*User); ok {
		meta = u.meta
	}

	for _, user := range engin.accounts {
		if reflect.DeepEqual(user.meta, meta) {
			return user, true
		}
	}
	return nil, false
}

func (user *User) findTarget(meta any) (*Target, bool) {
	if t, ok := meta.(*Target); ok {
		meta = t.meta
	}

	for _, target := range user.targets {
		if reflect.DeepEqual(target.meta, meta) {
			return target, true
		}
	}
	return nil, false
}

func (target *Target) finSubTarget(meta any) (*SubTarget, bool) {
	if s, ok := meta.(*SubTarget); ok {
		meta = s.meta
	}

	for _, subTarget := range target.subTargets {
		if reflect.DeepEqual(subTarget.meta, meta) {
			return subTarget, true
		}
	}
	return nil, false
}
