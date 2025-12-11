package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/AkhmadeevRus/sublease-app/pkg/di"
	emailsmtp "github.com/AkhmadeevRus/sublease-app/pkg/email_smtp"
	"github.com/AkhmadeevRus/sublease-app/pkg/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initialazing configs: %s", err.Error())
	}

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error while load dotenv: %s", err.Error())
	}

	db, err := server.NewPostgresDB(server.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	dbNum, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		logrus.Fatalf("err while create connection to db: %s", err.Error())
	}

	redisOption := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNum,
	}

	cacheDb, err := server.NewRedisDb(&redisOption)
	if err != nil {
		logrus.Fatalf("err while create connection to db: %s", err.Error())
	}

	codeExp, err := time.ParseDuration(os.Getenv("CODE_EXP"))
	if err != nil {
		logrus.Fatalf("Error while create connection to db: %s", err.Error())
	}

	codeLength, err := strconv.Atoi(os.Getenv("CODE_LENGTH"))
	if err != nil {
		logrus.Fatalf("Error while create connection to db: %s", err.Error())
	}

	emailCfg := emailsmtp.NewEmailCfg(
		os.Getenv("OWNER_EMAIL"),
		os.Getenv("OWNER_PASSWORD"),
		os.Getenv("SMTP_ADDR"),
		codeLength,
		codeExp,
	)

	repos := di.NewRepository(db, cacheDb, emailCfg)
	emailSmtpService := emailsmtp.NewEmailSmtpService(
		repos.EmailSmtpRepository,
		repos.EmailSmtpCacheRepository,
	)
	services := di.NewService(repos, emailSmtpService)
	handlers := di.NewHandler(services)
	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
