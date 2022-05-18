package main

import (
	"net/http"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/config"
	userDelivery "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/delivery"
	userRepository "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/usecase"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/middleware"
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
	userUseCase := usecase.NewUserUsecase(userRepo, cfg.Timeouts.ContextTimeout, log)
	sessionUseCase := usecase.NewSessionUsecase(sessionRepo, cfg.Timeouts.ContextTimeout, log)

	m := middleware.NewMiddleware(userRepo, sessionRepo)

	userDelivery.SetUserRouting(router, log, userUseCase, sessionUseCase, m)

	server := &http.Server{
		Handler:      router,
		Addr:         cfg.ServerConfig.Addr,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Infof("Server start at addt %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
