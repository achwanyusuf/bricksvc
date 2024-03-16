package bank

import (
	"context"

	"github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/repository/account"
	"github.com/achwanyusuf/bricksvc/src/repository/bank"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/volatiletech/null/v8"
)

type Bank struct {
	log     logger.LoggerInterface
	conf    Conf
	bank    bank.BankInterface
	account account.AccountInterface
}

type Conf struct{}

type BankInterface interface {
	GetBankAccount(ctx context.Context, cacheControl string, v model.GetBankAccount, apikey string) (clientresponse.BankAccount, error)
}

func New(conf Conf, logger *logger.LoggerInterface, bank bank.BankInterface, account account.AccountInterface) BankInterface {
	return &Bank{
		conf:    conf,
		log:     *logger,
		bank:    bank,
		account: account,
	}
}

func (b *Bank) GetBankAccount(ctx context.Context, cacheControl string, v model.GetBankAccount, apikey string) (clientresponse.BankAccount, error) {
	_, err := b.account.GetSingleByParam(ctx, "", &model.GetAccountByParam{
		APIKey: null.StringFrom(apikey),
	})
	if err != nil {
		return clientresponse.BankAccount{}, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, err)
	}

	return b.bank.GetBankAccount(ctx, cacheControl, v)
}
