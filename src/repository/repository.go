package repository

import (
	"database/sql"

	"github.com/achwanyusuf/bricksvc/src/repository/account"
	"github.com/achwanyusuf/bricksvc/src/repository/accountrole"
	"github.com/achwanyusuf/bricksvc/src/repository/bank"
	"github.com/achwanyusuf/bricksvc/src/repository/role"
	"github.com/achwanyusuf/bricksvc/src/repository/transfer"
	"github.com/achwanyusuf/bricksvc/utils/kafkalib"
	goredislib "github.com/redis/go-redis/v9"
)

type Repository struct {
	Conf  Config
	DB    *sql.DB
	Redis *goredislib.Client
	Kafka kafkalib.ProducerInterface
}

type Config struct {
	Account     account.Conf     `mapstructure:"account"`
	Role        role.Conf        `mapstructure:"role"`
	AccountRole accountrole.Conf `mapstructure:"account_role"`
	Bank        bank.Conf        `mapstructure:"bank"`
	Transfer    transfer.Conf    `mapstructure:"transfer"`
}

type RepositoryInterface struct {
	Account     account.AccountInterface
	Role        role.RoleInterface
	AccountRole accountrole.AccountRoleInterface
	Bank        bank.BankInterface
	Transfer    transfer.TransferInterface
}

func New(d *Repository) *RepositoryInterface {
	return &RepositoryInterface{
		account.New(d.Conf.Account, d.DB, d.Redis),
		role.New(d.Conf.Role, d.DB, d.Redis),
		accountrole.New(d.Conf.AccountRole, d.DB, d.Redis),
		bank.New(d.Conf.Bank, d.Redis),
		transfer.New(d.Conf.Transfer, d.DB, d.Redis, d.Kafka),
	}
}
