package settings

import (
	"log"

	env "github.com/Netflix/go-env"
	_ "github.com/joho/godotenv/autoload"
)

var Env struct {
	RedisUrl            string `env:"REDIS_URL"`
	DatabaseUrl         string `env:"DATABASE_URL"`
	AuthorizeServiceUrl string `env:"AUTHORIZE_URL"`
	BrokerStreamUrl     string `env:"BROKER_STEAM"`
}

func init() {
	if _, err := env.UnmarshalFromEnviron(&Env); err != nil {
		log.Fatal(err)
	}
}
