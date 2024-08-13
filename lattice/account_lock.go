package lattice

import (
	"fmt"
	"sync"
)

func NewAccountLock() AccountLock {
	return &accountLock{}
}

type accountLock struct {
	locks sync.Map
}

type AccountLock interface {
	// Obtain 获取账户锁
	//
	// Parameters:
	//   -chainId, address string
	Obtain(chainId, address string)

	// Unlock 释放账户锁
	//
	// Parameters:
	//   - chainId, address string
	Unlock(chainId, address string)
}

func (l *accountLock) Obtain(chainId, address string) {
	v, _ := l.locks.LoadOrStore(fmt.Sprintf("%s_%s", chainId, address), &sync.Mutex{})
	mutex := v.(*sync.Mutex)
	mutex.Lock()
}

func (l *accountLock) Unlock(chainId, address string) {
	if v, ok := l.locks.Load(fmt.Sprintf("%s_%s", chainId, address)); ok {
		mutex := v.(*sync.Mutex)
		mutex.Unlock()
		// fixme whether to delete the value
		// _i.locks.Delete(mutexName)
	}
}
