package main

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/guilhermealvess/guicpay/internal/infra/service"
	"github.com/guilhermealvess/guicpay/internal/settings"
	goredislib "github.com/redis/go-redis/v9"
)

func lock() {
	start := time.Now()
	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.
	mutexname := "my-global-mutex"
	mutex := rs.NewMutex(mutexname, redsync.WithExpiry(1*time.Second))

	// Obtain a lock for our given mutex. After this is successful, no one else
	// can obtain the same lock (the same mutex name) until we unlock it.
	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)
	if err := mutex.Lock(); err != nil {
		println(time.Since(start).String())
		panic("TESTE TODO")
	}
	// Do your work that requires the lock.

	// Release the lock so other processes or threads can obtain a lock.
	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
}

func auth() {
	authService := service.NewAuthorizeService(settings.Env.AuthorizeServiceUrl)
	ctx := context.Background()

	if err := authService.RegisterUser(ctx, "id_random", "test@gmail.com", "Pass123"); err != nil {
		panic(err)
	}

	println("Success")

	if err := authService.Authorize(ctx, "test@gmail.com", "Pass123"); err != nil {
		panic(err)
	}

	println("Success")
}

func main() {
	auth()
}
