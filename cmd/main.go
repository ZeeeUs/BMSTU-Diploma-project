package main

import (
	"net/http"

	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/config"
	minioStor "github.com/ZeeeUs/BMSTU-Diploma-project/internal/minio/storage"
	minioCase "github.com/ZeeeUs/BMSTU-Diploma-project/internal/minio/usecase"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/handler"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/storage"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/student/usecase"
	supersHandler "github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/handler"
	supersStor "github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/storage"
	supersCase "github.com/ZeeeUs/BMSTU-Diploma-project/internal/supervisor/usecase"
	userHandler "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/handler"
	userStor "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/storage"
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
	userStorage := userStor.NewUserStorage(pgConn, log)
	sessionStorage := userStor.NewSessionStorage(rc, log)
	supersStorage := supersStor.NewSupersStorage(pgConn, log)
	studentStorage := storage.NewStudentStorage(pgConn, log)
	minioStorage, err := minioStor.NewMinioStorage(
		cfg.MinioConfig.Endpoint,
		cfg.MinioConfig.AccessKeyID,
		cfg.MinioConfig.SecretAccessKey,
		log,
	)
	if err != nil {
		log.Fatal(err)
	}

	// UseCases
	userUseCase := userCase.NewUserUsecase(userStorage, cfg.Timeouts.ContextTimeout, log)
	sessionUseCase := userCase.NewSessionUsecase(sessionStorage, cfg.Timeouts.ContextTimeout, log)
	supersUseCase := supersCase.NewSupersUsecase(supersStorage, log)
	studentUseCase := usecase.NewStudentUseCase(studentStorage, log)
	minioUseCase := minioCase.NewMinioUseCase(minioStorage, log)

	m := middleware.NewMiddleware(userStorage, sessionStorage)

	userHandler.SetUserRouting(router, log, userUseCase, sessionUseCase, m)
	supersHandler.SetSupersRouting(router, log, supersUseCase, m)
	handler.SetStudentRouting(router, log, studentUseCase, minioUseCase, m)

	server := &http.Server{
		Handler:      router,
		Addr:         cfg.ServerConfig.Addr,
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Infof("Server start at addr %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
