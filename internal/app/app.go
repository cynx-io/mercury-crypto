package app

import (
	"log"
	"mercury/internal/pkg/logger"
)

type App struct {
	Dependencies *Dependencies
	Repos        *Repos
	Services     *Services
}

func NewApp(configPath string) (*App, error) {

	log.Println("Initializing Dependencies")
	dependencies := NewDependencies(configPath)

	if dependencies.Config.Database.AutoMigrate {
		logger.Info("Running database migrations")
		err := dependencies.DatabaseClient.RunMigrations()
		if err != nil {
			logger.Fatal("Failed to run migrations: ", err)
		}
	}

	logger.Info("Initializing Repositories")
	repos := NewRepos(dependencies)

	logger.Info("Initializing Services")
	services := NewServices(repos, dependencies)

	logger.Info("App initialized")
	return &App{
		Dependencies: dependencies,
		Repos:        repos,
		Services:     services,
	}, nil
}
