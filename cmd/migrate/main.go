package main

import (
	"github.com/rs/zerolog/log"

	"github.com/huypham67/bookmark-service/pkg/logger"
	"github.com/huypham67/bookmark-service/pkg/sqldb"
)

func main() {
	if err := logger.NewLoggerClient(""); err != nil {
		log.Error().Err(err).Msg("failed to initialize logger")
		return
	}

	dbClient, err := sqldb.NewDBClient("")
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize postgres client")
		return
	}

	if err := sqldb.MigratePostgresDB(dbClient, "migrations"); err != nil {
		log.Error().Err(err).Msg("failed to run database migrations")
		return
	}

	log.Info().Msg("database migrations completed successfully")
}
