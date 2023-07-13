package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GetConfig()
}

var (
	c *Config
)

type Config struct {
	Redis Redis
	Mongo Mongo
}

type Redis struct {
	Network  string
	Address  string
	Username string
	Password string
}

type Mongo struct {
	Address string
}

func GetConfig() *Config {
	if c == nil {
		//REDIS
		network := os.Getenv("REDIS_NETWORK")
		if network == "" {
			panic("REDIS_NETWORK is not set")
		}

		address := os.Getenv("REDIS_ADDRESS")
		if address == "" {
			panic("REDIS_ADDRESS is not set")
		}

		username := os.Getenv("REDIS_USERNAME")
		if username == "" {
			panic("REDIS_USERNAME is not set")
		}

		password := os.Getenv("REDIS_PASSWORD")
		if password == "" {
			panic("REDIS_PASSWORD is not set")
		}

		//MONGO
		mongoAddress := os.Getenv("MONGO_ADDRESS")
		if address == "" {
			panic("MONGO_ADDRESS is not set")
		}

		c = &Config{
			Redis: Redis{
				Network:  network,
				Address:  address,
				Username: username,
				Password: password,
			},
			Mongo: Mongo{
				Address: mongoAddress,
			},
		}

		return c
	}

	return c
}
