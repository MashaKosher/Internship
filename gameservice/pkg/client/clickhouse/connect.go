package clickhouse

import (
	"context"
	"database/sql"
	"gameservice/internal/config"
	"gameservice/pkg/logger"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func ConnectClickHouse() *sql.DB {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{config.AppConfig.Clickhouse.Host + ":" + config.AppConfig.Clickhouse.Port},
		Auth: clickhouse.Auth{
			Database: config.AppConfig.Clickhouse.Name,
			Username: config.AppConfig.Clickhouse.User,
			Password: config.AppConfig.Clickhouse.Password,
		},
		DialTimeout: time.Second * 10,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})

	conn.SetMaxOpenConns(10)   // Максимальное количество открытых соединений
	conn.SetMaxIdleConns(5)    // Максимальное количество бездействующих соединений
	conn.SetConnMaxLifetime(0) // Максимальное время жизни соединения

	// 2. Проверка подключения
	ctx := context.Background()
	if err := conn.PingContext(ctx); err != nil {
		logger.L.Fatal("Ping failed:" + err.Error())
	}
	logger.L.Info("Connected to ClickHouse!")

	// // 3. Выполнение запроса
	// rows, err := conn.QueryContext(ctx, "SELECT version()")
	// if err != nil {
	// 	logger.L.Fatal("Query failed:" + err.Error())
	// }
	// defer rows.Close()

	// var version string
	// for rows.Next() {
	// 	if err := rows.Scan(&version); err != nil {
	// 		logger.L.Fatal("Scan failed:" + err.Error())
	// 	}
	// 	logger.L.Info("ClickHouse version:" + version)
	// }
	return conn
}
