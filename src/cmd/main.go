package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/achwanyusuf/bricksvc/adapter/psqlclient"
	"github.com/achwanyusuf/bricksvc/adapter/redisclient"
	"github.com/achwanyusuf/bricksvc/conf"
	"github.com/achwanyusuf/bricksvc/docs"
	"github.com/achwanyusuf/bricksvc/src/presentation/consumer"
	"github.com/achwanyusuf/bricksvc/src/presentation/rest"
	"github.com/achwanyusuf/bricksvc/src/presentation/scheduler"
	"github.com/achwanyusuf/bricksvc/src/repository"
	"github.com/achwanyusuf/bricksvc/src/usecase"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/httpserver"
	"github.com/achwanyusuf/bricksvc/utils/kafkalib"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/achwanyusuf/bricksvc/utils/migration"
	"github.com/achwanyusuf/bricksvc/utils/schedulerengine"
)

var (
	staticConfPath, Namespace, BuildTime, Version string
	migrateup, migratedown                        bool
)

// @contact.name   Brick Support
// @contact.url		https://www.brick.com/support
// @contact.email 	support@brick.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl %s
// @securitydefinitions.apikey APIKey
// @in header
// @name API-Key
// @description Type "APIKey" followed by api key
func main() {
	flag.StringVar(&staticConfPath, "staticConfPath", "./conf/conf.yaml", "config path")
	flag.BoolVar(&migrateup, "migrateup", false, "run migration up")
	flag.BoolVar(&migratedown, "migratedown", false, "run migration up")
	flag.Parse()
	cfg, err := conf.New(staticConfPath)
	if err != nil {
		panic(err)
	}
	if Version == "" {
		Version = "v1.0.0"
	}
	log := logger.New(&logger.Config{
		IsFile: false,
		Level:  logger.LevelDebug,
		CustomFields: map[string]interface{}{
			"namespace":  Namespace,
			"version":    Version,
			"build_time": BuildTime,
			"pid":        os.Getpid(),
		},
	})

	psql, err := psqlclient.PsqlConnect(cfg.App.PSQL)
	if err != nil {
		logger.Log.Panic(errormsg.WriteErr(err))
		panic(err)
	}

	if migrateup {
		migrate, err := migration.New(migration.Conf{
			DB:   psql,
			Path: cfg.App.PSQL.MigrationPath,
		})
		if err != nil {
			logger.Log.Error(err)
		}
		if err == nil {
			if err := migrate.Up(); err != nil {
				logger.Log.Error(err)
			}
		}
	}

	if migratedown {
		migrate, err := migration.New(migration.Conf{
			DB:   psql,
			Path: cfg.App.PSQL.MigrationPath,
		})
		if err == nil {
			if err := migrate.Down(); err != nil {
				logger.Log.Error(err)
			}
		}
	}

	redis, err := redisclient.RedisConnect(cfg.App.Redis)
	if err != nil {
		logger.Log.Panic(errormsg.WriteErr(err))
		panic(err)
	}

	httpServer := httpserver.New(&cfg.App.HTTPServer, log)
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%v", cfg.App.HTTPServer.Host, cfg.App.HTTPServer.Port)
	docs.SwaggerInfo.Title = Namespace
	docs.SwaggerInfo.Description = cfg.App.Description
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Version = Version
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.SwaggerTemplate = fmt.Sprintf(docs.SwaggerInfo.SwaggerTemplate, cfg.App.OAuth2PasswordTokenUrl)
	httpServer.Setup()

	kafkaConsumer, err := kafkalib.NewConsumer(&cfg.App.Kafka)
	if err != nil {
		logger.Log.Panic(errormsg.WriteErr(err))
		panic(err)
	}

	kafkaProducer, err := kafkalib.NewProducer(&cfg.App.Kafka)
	if err != nil {
		logger.Log.Panic(errormsg.WriteErr(err))
		panic(err)
	}

	rp := repository.New(&repository.Repository{
		Conf:  cfg.Repository,
		DB:    psql,
		Redis: redis,
		Kafka: kafkaProducer,
	})

	uc := usecase.New(&usecase.Usecase{
		Conf:       cfg.Usecase,
		Log:        &log,
		Repository: rp,
	})
	re := &rest.Rest{
		Conf:       cfg.Rest,
		Log:        &log,
		Usecase:    uc,
		HTTPServer: httpServer.Get(),
	}
	rinf := rest.New(re)
	re.Serve(rinf)

	sched := schedulerengine.New()
	schedulerConf := &scheduler.Scheduler{
		Conf:      cfg.Scheduler,
		Scheduler: sched,
		Usecase:   *uc,
	}

	readSignal := make(chan os.Signal, 1)

	signal.Notify(
		readSignal,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	// running application
	go func() {
		if err := httpServer.Run(); err != nil {
			logger.Log.Panic(errormsg.WriteErr(err))
			panic(err)
		}
	}()
	if cfg.App.Kafka.Enabled {
		go func() {
			err := consumer.New(&consumer.Consumer{
				Conf:     cfg.Consumer,
				Consumer: kafkaConsumer,
				Usecase:  *uc,
			})
			if err != nil {
				logger.Log.Panic(errormsg.WriteErr(err))
				panic(err)
			}
		}()
	}

	go func() {
		scheduler.New(schedulerConf)
	}()
	<-readSignal

	logger.Log.Warn("closing gracefully . . . ")
	st := time.Now()
	if cfg.App.Kafka.Enabled {
		kafkaConsumer.Stop()
	}
	httpServer.Close()
	sched.Close()

	// close all connection here before shutdown
	logger.Log.Error("service shutdown!", time.Since(st).Seconds(), "sec")
}
