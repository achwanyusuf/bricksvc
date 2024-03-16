package usecase

import (
	"github.com/achwanyusuf/bricksvc/src/repository"
	"github.com/achwanyusuf/bricksvc/src/usecase/account"
	"github.com/achwanyusuf/bricksvc/src/usecase/accountrole"
	"github.com/achwanyusuf/bricksvc/src/usecase/bank"
	"github.com/achwanyusuf/bricksvc/src/usecase/role"
	"github.com/achwanyusuf/bricksvc/src/usecase/transfer"
	"github.com/achwanyusuf/bricksvc/utils/logger"
)

type Usecase struct {
	Conf       Config
	Log        *logger.LoggerInterface
	Repository *repository.RepositoryInterface
}

type Config struct {
	Account     account.Conf     `mapstructure:"account"`
	Role        role.Conf        `mapstructure:"role"`
	AccountRole accountrole.Conf `mapstructure:"account_role"`
	Bank        bank.Conf        `mapstructure:"bank"`
	Transfer    transfer.Conf    `mapstructure:"transfer"`
}

type UsecaseInterface struct {
	Account     account.AccountInterface
	Role        role.RoleInterface
	AccountRole accountrole.AccountRoleInterface
	Bank        bank.BankInterface
	Transfer    transfer.TransferInterface
}

func New(u *Usecase) *UsecaseInterface {
	return &UsecaseInterface{
		account.New(u.Conf.Account, u.Log, u.Repository.Account, u.Repository.Role, u.Repository.AccountRole),
		role.New(u.Conf.Role, u.Log, u.Repository.Role),
		accountrole.New(u.Conf.AccountRole, u.Log, u.Repository.AccountRole),
		bank.New(u.Conf.Bank, u.Log, u.Repository.Bank, u.Repository.Account),
		transfer.New(u.Conf.Transfer, u.Log, u.Repository.Account, u.Repository.Transfer),
	}
}
