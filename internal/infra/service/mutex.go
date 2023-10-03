package service

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/guilhermealvess/guicpay/internal/domain/service"
	goredislib "github.com/redis/go-redis/v9"
)

type mutexService struct {
	client *redsync.Redsync
}

func NewMutexService(address string) mutexService {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: address,
	})
	pool := goredis.NewPool(client)

	rs := redsync.New(pool)
	mu := mutexService{client: rs}
	return mu
}

func (m mutexService) NewMutex(mutexname string, ttl time.Duration) service.Mutex {
	mu := m.client.NewMutex(mutexname, redsync.WithExpiry(ttl))
	return mutex{mu: mu}
}

type mutex struct {
	mu *redsync.Mutex
}

func (m mutex) Lock(ctx context.Context) error {
	return m.mu.Lock()
}

func (m mutex) Unlock(ctx context.Context) error {
	_, err := m.mu.Unlock()
	return err
}
