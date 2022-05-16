package main

import (
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/config"
	userRepository "github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/repository"
	"github.com/ZeeeUs/BMSTU-Diploma-project/internal/user/usecase"
	"github.com/sirupsen/logrus"
)

func main() {
	// logger
	log := logrus.New()

	// configs
	cfg := config.NewConfig()

	// repositories
	userRepo, err := userRepository.NewUserRepository(cfg, log)
	if err != nil {
		log.Fatal(err)
	}

	// usecases
	userUseCase := usecase.NewUserUsecase(userRepo, cfg.Timeouts.ContextTimeout)

}
