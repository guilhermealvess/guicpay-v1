package service

import (
	"context"
	"time"
)

type MutexService interface {
	NewMutex(mutexname string, ttl time.Duration) Mutex
}

type Mutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}
