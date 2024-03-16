package conf

import (
	"github.com/achwanyusuf/bricksvc/adapter/psqlclient"
	"github.com/achwanyusuf/bricksvc/adapter/redisclient"
	"github.com/achwanyusuf/bricksvc/src/presentation/consumer"
	"github.com/achwanyusuf/bricksvc/src/presentation/rest"
	"github.com/achwanyusuf/bricksvc/src/presentation/scheduler"
	"github.com/achwanyusuf/bricksvc/src/repository"
	"github.com/achwanyusuf/bricksvc/src/usecase"
	"github.com/achwanyusuf/bricksvc/utils/httpserver"
	"github.com/achwanyusuf/bricksvc/utils/kafkalib"
)

type Conf struct {
	App        App               `mapstructure:"app"`
	Consumer   consumer.Conf     `mapstructure:"consumer"`
	Rest       rest.Config       `mapstructure:"rest"`
	Scheduler  scheduler.Conf    `mapstructure:"scheduler"`
	Usecase    usecase.Config    `mapstructure:"usecase"`
	Repository repository.Config `mapstructure:"repository"`
}

type App struct {
	Env                    string            `mapstructure:"env"`
	Description            string            `mapstructure:"description"`
	OAuth2PasswordTokenUrl string            `mapstructure:"oauth2_password_token_url"`
	Kafka                  kafkalib.Conf     `mapstructure:"kafka"`
	HTTPServer             httpserver.Conf   `mapstructure:"http_server"`
	PSQL                   psqlclient.PSQL   `mapstructure:"psql"`
	Redis                  redisclient.Redis `mapstructure:"redis"`
}
