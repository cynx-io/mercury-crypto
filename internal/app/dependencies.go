package app

import (
	"log"
	"mercury/internal/dependencies"
	"mercury/internal/dependencies/coingecko"
	"mercury/internal/dependencies/gopluslabs"
	"mercury/internal/pkg"
	"mercury/internal/pkg/logger"
)

type Dependencies struct {
	Config *dependencies.Config

	DatabaseClient *dependencies.DatabaseClient

	CoinGecko  *coingecko.Client
	GoPlusLabs *gopluslabs.Client
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
	coinGeckoClient := coingecko.NewClient(&config.CoinGecko)

	logger.Info("Creating GoPlusLabs Client")
	goPlusLabsClient := gopluslabs.NewClient(&config.GoPlusLabs)

	logger.Info("Dependencies initialized")
	return &Dependencies{
		Config:         config,
		DatabaseClient: databaseClient,

		CoinGecko:  coinGeckoClient,
		GoPlusLabs: goPlusLabsClient,
	}
}
