package configs

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Settings Config

type Config struct {
	APP_KAFKATOPIC       string `envconfig:"APP_KAFKATOPIC"`
	APP_PORT             int    `envconfig:"APP_PORT"`
	MONGO_ADDR           string `envconfig:"APP_MONGODB_ADDR"`
	APP_KAFKA_ADDR       string `envconfig:"APP_KAFKA_ADDR"`
	APP_KAFKA_GROUPID    string `envconfig:"APP_KAFKA_GROUPID"`
	APP_SECURE_SECRETKEY string `envconfig:"APP_SECRET"`
}

func LoadConfig() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables")
	}

	if err := envconfig.Process("", &Settings); err != nil {
		log.Fatal(err)
	}
	log.Println("Configuration Done")

}
