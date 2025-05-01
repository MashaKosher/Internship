package setup

// import (
// 	"coreservice/internal/di"

// 	"github.com/hibiken/asynq"
// )

// const ClientPoolSize = 10

// func mustDelay(cfg di.ConfigType, logger di.LoggerType) di.DelayTaskType {

// 	return di.DelayTaskType{
// 		SeasonProducer: createDelayProducer(cfg, logger),
// 	}
// }

// func createDelayProducer(cfg di.ConfigType, logger di.LoggerType) di.DelayProducer {
// 	producer := asynq.NewClient(asynq.RedisClientOpt{
// 		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
// 		PoolSize: ClientPoolSize,
// 	})
// 	logger.Info("Asynq Client connected successfully")

// 	return producer
// }

// func createDelayConsumer(cfg di.ConfigType, logger di.LoggerType) di.DelayProducer {
// 	producer := asynq.NewClient(asynq.RedisClientOpt{
// 		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
// 		PoolSize: ClientPoolSize,
// 	})
// 	logger.Info("Asynq Client connected successfully")

// 	return producer
// }

// func deferDelayTask(delayTask di.DelayTaskType) {
// 	// bus.AuthConsumer.Close()
// 	// bus.DailyTaskConsumer.Close()
// 	// bus.SeasonConsumer.Close()
// 	// bus.AuthProducer.Close()
// 	delayTask.SeasonProducer.Close()
// }
