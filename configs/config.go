package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

var cfg *conf

type conf struct {
	DBDriver        string `mapstructure:"DB_DRIVER"`
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBUser          string `mapstructure:"DB_USER"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	DBName          string `mapstructure:"DB_NAME"`
	WebServerPort   string `mapstructure:"WEBSERVER_PORT"`
	JWTSecret       string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn    int    `mapstructure:"JWT_EXPIRESIN"`
	URL_CEP1        string `mapstructure:"URL_CEP_1"`
	URL_CEP2        string `mapstructure:"URL_CEP_2"`
	LIMITE_CONTAGEM string `mapstructure:"LIMITE_CONTAGEM"`
	TokenAuth       *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	println("url1: " + cfg.URL_CEP1)

	return cfg, err
}
