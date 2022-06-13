package main

import (
	"net/http"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/config"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/handler"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/storage"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/usecase"
	supersDelivery "github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/delivery"
	supersRRepo "github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/repository"
	supersCase "github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/usecase"
	userDelivery "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/delivery"
	userRepository "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	userCase "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/usecase"
	"github.com/ZeeeUs/BMSTU-Diploma-project/pkg/middleware"
	redisClient "github.com/ZeeeUs/BMSTU-Diploma-project/pkg/redis"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

func main() {
	// logger
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	log := logrus.New()
	log.SetFormatter(formatter)

	// configs
	cfg := config.NewConfig()

	// router
	router := mux.NewRouter()

	// Postgres
	pgConn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     cfg.DbConfig.DbHostName,
			Port:     uint16(cfg.DbConfig.DbPort),
			Database: cfg.DbConfig.DbName,
			User:     cfg.DbConfig.DbUser,
			Password: cfg.DbConfig.DbPassword,
		},
	})
	if err != nil {
		log.Fatalf("Error %s occurred during connection to database", err)
	}

	// redis
	rCClient := redis.NewClient(&redis.Options{
		Addr:       cfg.RedisConfig.Addr,
		Password:   cfg.RedisConfig.Password,
		MaxRetries: cfg.RedisConfig.MaxRetries,
	})
	pong, err := rCClient.Ping().Result()
	if err != nil {
		log.Fatalf("redis: %s", err)
	}
	log.Infof("succsessfully connetc to redis: %s", pong)

	rc := redisClient.New(rCClient)

	// repositories
	userRepo := userRepository.NewUserRepository(pgConn, log)
	sessionRepo := userRepository.NewSessionRepository(rc, log)
	supersRepo := supersRRepo.NewSupersRepository(pgConn, log)
	studentRepo := storage.NewStudentRepository(pgConn, log)

	// usecases
	userUseCase := userCase.NewUserUsecase(userRepo, cfg.Timeouts.ContextTimeout, log)
	sessionUseCase := userCase.NewSessionUsecase(sessionRepo, cfg.Timeouts.ContextTimeout, log)
	supersUseCase := supersCase.NewSupersUsecase(supersRepo, log)
	studentUseCase := usecase.NewStudentUseCase(studentRepo, log)

	m := middleware.NewMiddleware(userRepo, sessionRepo)

	userDelivery.SetUserRouting(router, log, userUseCase, sessionUseCase, m)
	supersDelivery.SetSupersRouting(router, log, supersUseCase, m)
	handler.SetStudentRouting(router, log, studentUseCase, m)

	server := &http.Server{
		Handler:      router,
		Addr:         cfg.ServerConfig.Addr,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Infof("Server start at addr %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
