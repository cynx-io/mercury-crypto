package app

import (
	"log"
	"mercury/internal/dependencies"
	"mercury/internal/dependencies/coingecko"
	"mercury/internal/pkg"
	"mercury/internal/pkg/logger"
)

type Dependencies struct {
	Config *dependencies.Config

	DatabaseClient *dependencies.DatabaseClient

	CoinGecko *coingecko.Client
}

func NewDependencies(configPath string) *Dependencies {

	log.Println("Loading Config")
	config, err := dependencies.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	log.Println("Initializing Logger")
	logger.InitLogger()

	logger.Info("Initializing Validator")
	pkg.InitValidator()

	logger.Info("Connecting to Database")
	databaseClient, err := dependencies.NewDatabaseClient(config)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}

	logger.Info("Creating CoinGecko Client")
	coinGeckoClient := coingecko.NewCoinGecko(&config.CoinGecko)

	logger.Info("Dependencies initialized")
	return &Dependencies{
		Config:         config,
		DatabaseClient: databaseClient,

		CoinGecko: coinGeckoClient,
	}
}
