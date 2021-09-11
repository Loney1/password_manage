package model

import (
	utime "adp_backend/infra/time"

	"sync"
)

type UserBucket struct {
	List map[string]*UserLoginCountInfo
	Mu   sync.RWMutex
}

type UserLoginCountInfo struct {
	LoginErrCount       int
	LastLoginExpireTime int64
	LastLoginErrorTime  int64
}

func (bucket *UserBucket) Add(username string, clt *UserLoginCountInfo) {
	bucket.Mu.Lock()
	defer bucket.Mu.Unlock()
	bucket.List[username] = clt
}

func (bucket *UserBucket) Get(username string) *UserLoginCountInfo {
	bucket.Mu.Lock()
	defer bucket.Mu.Unlock()

	clt, exists := bucket.List[username]
	if exists {
		isExp := clt.LastLoginErrorTime + 60*5

		if utime.CurSecond() > isExp {
			clt.LoginErrCount = 0
			clt.LastLoginErrorTime = 0
		}
	} else {
		bucket.List = make(map[string]*UserLoginCountInfo)
		bucket.List[username] = &UserLoginCountInfo{LoginErrCount: 0, LastLoginExpireTime: 0}
	}

	return bucket.List[username]
}

func (bucket *UserBucket) SetLoginErrCount(username string, value int) *UserLoginCountInfo {
	bucket.Mu.Lock()
	defer bucket.Mu.Unlock()

	clt, exists := bucket.List[username]

	if exists {
		if value > 0 {
			clt.LoginErrCount += value
			clt.LastLoginErrorTime = utime.CurSecond()
		} else {
			clt.LoginErrCount = 0
			clt.LastLoginErrorTime = 0
		}
	}

	return clt
}

func (bucket *UserBucket) SetLastLoginExpireTime(username string, LastLoginExpireTime int64) *UserLoginCountInfo {
	bucket.Mu.Lock()
	defer bucket.Mu.Unlock()

	clt, exists := bucket.List[username]

	if exists {
		clt.LastLoginExpireTime = LastLoginExpireTime
	}

	return clt
}
