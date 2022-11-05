package conf

import (
	"os"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
)

var AppConfig Config

type Config struct {
	JWTConfig   *JWTConfig
	ApilayerKey string
}

type JWTConfig struct {
	JWTSignatureKey   string
	JWTExpirationTime time.Duration
}

func init() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		panic(err)
	}

	beego.BConfig.RunMode = os.Getenv("beego_runmode")

	AppConfig.ApilayerKey = os.Getenv("apilayer_api_key")

	AppConfig.JWTConfig = &JWTConfig{}

	AppConfig.JWTConfig.JWTSignatureKey = os.Getenv("jwt_signature_key")
	if AppConfig.JWTConfig.JWTSignatureKey == "" {
		panic("jwt_signature_key not set")
	}

	jwtExpirationTimeStr := os.Getenv("jwt_expiration_time")
	jwtExpirationTime, _ := strconv.Atoi(jwtExpirationTimeStr)
	if jwtExpirationTime == 0 {
		jwtExpirationTime = 24 * 60 * 60
	}
	AppConfig.JWTConfig.JWTExpirationTime = time.Duration(jwtExpirationTime) * time.Second

}
