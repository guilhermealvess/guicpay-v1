package settings

import (
	"log"
	"time"

	env "github.com/Netflix/go-env"
	_ "github.com/joho/godotenv/autoload"
)

var Env struct {
	RedisUrl            string        `env:"REDIS_URL,default=localhost:6379"`
	DatabaseUrl         string        `env:"DATABASE_URL,default="`
	AuthorizeServiceUrl string        `env:"AUTHORIZE_URL"`
	BrokerStreamUrl     string        `env:"BROKER_STREAM"`
	Salt                string        `env:"SALT"`
	TransactionTimeout  time.Duration `env:"TRANSACTION_TIMEOUT,default=7s"`
}

func init() {
	if _, err := env.UnmarshalFromEnviron(&Env); err != nil {
		log.Fatal(err)
	}
}
