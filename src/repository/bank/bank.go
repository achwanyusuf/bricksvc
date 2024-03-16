package bank

import (
	"context"
	"fmt"
	"time"

	"github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	jsoniter "github.com/json-iterator/go"
	goredislib "github.com/redis/go-redis/v9"
)

type Bank struct {
	Redis *goredislib.Client
	Conf  Conf
}

type Conf struct {
	DefaultPageLimit    int           `mapstructure:"page_limit"`
	RedisExpirationTime time.Duration `mapstructure:"expiration_time"`
	GetBankAccountURL   string        `mapstructure:"get_bank_account_url"`
}

type BankInterface interface {
	GetBankAccount(ctx context.Context, cacheControl string, v model.GetBankAccount) (clientresponse.BankAccount, error)
}

func New(conf Conf, rds *goredislib.Client) BankInterface {
	return &Bank{
		Redis: rds,
		Conf:  conf,
	}
}

func (b *Bank) GetBankAccount(ctx context.Context, cacheControl string, v model.GetBankAccount) (clientresponse.BankAccount, error) {
	str, err := jsoniter.MarshalToString(v)
	if err != nil {
		return clientresponse.BankAccount{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetBankAccountKey, str)
	if cacheControl != model.MustRevalidate {
		res, err1 := b.getSingleByParamRedis(ctx, key)
		if err1 != nil {
			if err1 == goredislib.Nil {
				res, err := b.getBankAccountHTTPClient(v)
				if err == nil {
					dataStr, err := jsoniter.Marshal(&res)
					if err != nil {
						return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = b.setRedis(ctx, key, string(dataStr))
					if err != nil {
						return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
					}
				}
				return res, err
			}
			return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error marshal param")
		}
		return res, nil
	}

	res, err := b.getBankAccountHTTPClient(v)
	if err == nil {
		dataStr, err := jsoniter.Marshal(&res)
		if err != nil {
			return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get client")
		}
		err = b.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
	}
	return res, nil
}
