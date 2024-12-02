package cmd

import (
	"flag"
	"hacko-app/db/seeds"
	"hacko-app/internal/adapter"

	"github.com/rs/zerolog/log"
)

func RunSeed(cmd *flag.FlagSet, args []string) {
	var (
		table = cmd.String("table", "", "seed to run")
		total = cmd.Int("total", 1, "total of records to seed")
	)

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	adapter.Adapters.Sync(
		adapter.WithHackoPostgres(),
	)
	defer func() {
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
	}()

	seeds.Execute(adapter.Adapters.HackoPostgres, *table, *total)
}
