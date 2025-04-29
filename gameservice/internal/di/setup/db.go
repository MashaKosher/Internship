package setup

import (
	"context"
	"gameservice/internal/di"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

const maxConns = 10    // Максимальное количество открытых соединений
const maxdeadConns = 5 // Максимальное количество бездействующих соединений
const sessionLife = 0  // Максимальное время жизни соединения

func mustDB(cfg di.ConfigType, logger di.LoggerType) di.DBType {

	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{cfg.Clickhouse.Host + ":" + cfg.Clickhouse.Port},
		Auth: clickhouse.Auth{
			Database: cfg.Clickhouse.Name,
			Username: cfg.Clickhouse.User,
			Password: cfg.Clickhouse.Password,
		},
		DialTimeout: time.Second * 10,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})

	conn.SetMaxOpenConns(maxConns)
	conn.SetMaxIdleConns(maxdeadConns)
	conn.SetConnMaxLifetime(sessionLife)

	ctx := context.Background()
	if err := conn.PingContext(ctx); err != nil {
		logger.Fatal("Ping failed:" + err.Error())
	}
	logger.Info("ClickHouse connected successfully")

	return conn
}
