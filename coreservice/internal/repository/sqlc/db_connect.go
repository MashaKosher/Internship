package sqlc

import (
	"context"
	"coreservice/internal/config"
	"coreservice/internal/logger"
	db "coreservice/internal/repository/sqlc/generated"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var Query *db.Queries
var Ctx context.Context

func DBConnect() (*pgx.Conn, context.Context) {
	dan, err := pgx.ParseConfig(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.AppConfig.DB.Host, config.AppConfig.DB.Port, config.AppConfig.DB.User, config.AppConfig.DB.Password, config.AppConfig.DB.Name, config.AppConfig.DB.SSLMode))
	if err != nil {
		logger.Logger.Error("Unable to parse connection string: " + err.Error())
		panic(err)
	}

	Ctx = context.Background()
	conn, err := pgx.ConnectConfig(Ctx, dan)

	if err != nil {
		defer conn.Close(Ctx)
		logger.Logger.Error(err.Error())
		panic(err)
	}

	Query = db.New(conn)
	return conn, Ctx
}
