package migrate

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stuckinforloop/fabrik/db"
)

type DSNConfig struct {
	MasterDSN string `mapstructure:"master_dsn"`
	ReaderDSN string `mapstructure:"reader_dsn"`
}

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "run db migrations",
	RunE: func(cmd *cobra.Command, _ []string) error {
		configFile, _ := cmd.Flags().GetString("config")
		if configFile != "" {
			viper.SetConfigFile(configFile)
		}

		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}

		var config DSNConfig
		if err := viper.UnmarshalKey("database", &config); err != nil {
			return fmt.Errorf("error unmarshaling config: %w", err)
		}

		dsn := config.MasterDSN
		if dsn == "" {
			dsn = os.Getenv("DATABASE_URL")
			if dsn == "" {
				return fmt.Errorf("database DSN not found in config or environment")
			}
		}

		db, err := db.New(context.Background(), db.Config{
			MasterDSN: dsn,
			ReaderDSN: dsn,
		})
		if err != nil {
			return fmt.Errorf("failed to create db: %w", err)
		}

		if err := db.Migrate(); err != nil {
			return fmt.Errorf("failed to migrate: %w", err)
		}

		return nil
	},
}

func init() {
	MigrateCmd.Flags().StringP("config", "c", "", "config file")
}
