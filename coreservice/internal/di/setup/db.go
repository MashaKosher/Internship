package setup

import (
	"context"
	"coreservice/internal/di"
	db "coreservice/internal/repository/sqlc/generated"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn
var ctx context.Context

func mustDB(cfg di.ConfigType, logger di.LoggerType) di.DBType {
	dan, err := pgx.ParseConfig(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode))
	if err != nil {
		logger.Fatal("Unable to parse connection string: " + err.Error())
	}

	ctx := context.Background()
	conn, err = pgx.ConnectConfig(ctx, dan)

	if err != nil {
		defer conn.Close(ctx)
		logger.Fatal(err.Error())
	}

	return db.New(conn)
}

func deferDB() {
	conn.Close(ctx)
}
