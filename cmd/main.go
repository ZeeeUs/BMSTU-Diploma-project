package main

import (
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/config"
	userDelivery "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/delivery"
	userRepository "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/usecase"
	redisClient "github.com/ZeeeUs/BMSTU-Diploma-project/pkg/redis"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	// logger
	log := logrus.New()

	// configs
	cfg := config.NewConfig()

	router := mux.NewRouter()

	// redis
	rCClient := redis.NewClient(&redis.Options{
		Addr:       cfg.RedisConfig.Addr,
		Password:   cfg.RedisConfig.Password,
		MaxRetries: cfg.RedisConfig.MaxRetries,
	})
	rc := redisClient.New(rCClient)

	// repositories
	userRepo, err := userRepository.NewUserRepository(cfg, log)
	if err != nil {
		log.Fatal(err)
	}
	sessionRepo := userRepository.NewSessionRepository(rc, log)

	// usecases
	userUseCase := usecase.NewUserUsecase(userRepo, cfg.Timeouts.ContextTimeout)

	userDelivery.SetUserRouting(router, log, userUseCase)

}
