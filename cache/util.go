package cache

import (
	"net/url"
	"sync"
)

func redacted(u *url.URL) string {
	if u == nil {
		return ""
	}

	ru := *u
	if _, has := ru.User.Password(); has {
		ru.User = url.UserPassword(ru.User.Username(), "xxxxx")
	}
	return ru.String()
}

type keyMutex struct {
	m *sync.Map
}

func (km keyMutex) Lock(key interface{}) {
	m := sync.Mutex{}
	act, _ := km.m.LoadOrStore(key, &m)
	mm := act.(*sync.Mutex)
	mm.Lock()
	if mm != &m {
		mm.Unlock()
		km.Lock(key)
	}
}

func (km keyMutex) Unlock(key interface{}) {
	act, exist := km.m.Load(key)
	if !exist {
		panic("unlock of unlocked mutex")
	}
	m := act.(*sync.Mutex)
	km.m.Delete(key)
	m.Unlock()
}
