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

func NewMutexService() mutexService {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)
	mu := mutexService{client: rs}
	return mu

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	// mutex := rs.NewMutex(mutexname, redsync.WithExpiry(1*time.Second))
	// mutexname := "my-global-mutex"

}

/* func (m mutex) NewMutex(mutexname string, ttl time.Duration) *redsync.Mutex {
	return m.client.NewMutex(mutexname, redsync.WithExpiry(ttl))
}

func (m mutex) Lock(ctx context.Context, mu *redsync.Mutex) error {
	return mu.Lock()
}

func (m mutex) Unlock(ctx context.Context, mutexname string) error {
	mu := m.client.NewMutex(mutexname, redi)
} */

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
