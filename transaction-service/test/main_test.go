package test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/majidmohsenifar/heli-tech/transaction-service/config"
	"github.com/majidmohsenifar/heli-tech/transaction-service/logger"
	"github.com/spf13/viper"
)

var loggerService *slog.Logger
var mainDB *pgxpool.Pool
var viperConfig *viper.Viper

var Seeders = map[string]func(ctx context.Context, db *pgxpool.Pool){}

func TestMain(m *testing.M) {
	ctx := context.Background()
	setupDB(ctx)
	runSeeders(ctx)
	exitCode := m.Run()
	os.Exit(exitCode)

}

func setupDB(ctx context.Context) {
	db := getDB()
	content, err := os.ReadFile("./data/db.sql")
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = db.Exec(ctx, string(content))

	if err != nil {
		fmt.Println(err.Error())
	}
}

func runSeeders(ctx context.Context) {
	db := getDB()
	for name, seedFunc := range Seeders {
		fmt.Println("running seed ", name)
		seedFunc(ctx, db)
	}
}

func getLogger() *slog.Logger {
	if loggerService != nil {
		return loggerService
	}
	loggerService = logger.NewLogger()
	return loggerService
}

func getViperConfig() *viper.Viper {
	if viperConfig != nil {
		return viperConfig
	}
	viperConfig = config.NewViper("../config/")
	viperConfig.Set("jwt.privatekey", "../config/jwt/private.pem")
	viperConfig.Set("jwt.publickey", "../config/jwt/public.pem")
	return viperConfig
}

func getDB() *pgxpool.Pool {
	if mainDB != nil {
		return mainDB
	}
	viper := getViperConfig()
	dbPool, err := pgxpool.New(context.Background(), viper.GetString("db.dsn"))
	if err != nil {
		panic("can not establish connection to database")
	}
	mainDB = dbPool
	return mainDB
}

func truncateDB() error {
	ctx := context.Background()
	db := getDB()
	_, err := db.Exec(ctx, "TRUNCATE TABLE user_balances CASCADE")
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, "TRUNCATE TABLE transactions CASCADE")
	if err != nil {
		return err
	}
	return nil
}
