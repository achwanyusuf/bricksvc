package accountrole

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	jsoniter "github.com/json-iterator/go"
	goredislib "github.com/redis/go-redis/v9"
)

type AccountRole struct {
	DB    *sql.DB
	Redis *goredislib.Client
	Conf  Conf
}

type Conf struct {
	DefaultPageLimit    int           `mapstructure:"page_limit"`
	RedisExpirationTime time.Duration `mapstructure:"expiration_time"`
}

type AccountRoleInterface interface {
	Insert(ctx context.Context, data *entity.AccountRole) error
	GetSingleByParam(ctx context.Context, cacheControl string, param *model.GetAccountRoleByParam) (entity.AccountRole, error)
	Update(ctx context.Context, AccountRole *entity.AccountRole) error
	Delete(ctx context.Context, AccountRole *entity.AccountRole, id int64, isHardDelete bool) error
	GetByParam(ctx context.Context, cacheControl string, param *model.GetAccountRolesByParam) (entity.AccountRoleSlice, model.Pagination, error)
}

func New(conf Conf, db *sql.DB, rds *goredislib.Client) AccountRoleInterface {
	return &AccountRole{
		DB:    db,
		Redis: rds,
		Conf:  conf,
	}
}

func (a *AccountRole) Insert(ctx context.Context, data *entity.AccountRole) error {
	return a.insertPSQL(ctx, data)
}

func (a *AccountRole) GetSingleByParam(ctx context.Context, cacheControl string, param *model.GetAccountRoleByParam) (entity.AccountRole, error) {
	str, err := jsoniter.Marshal(param)
	if err != nil {
		return entity.AccountRole{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetSingleByParamAccountRoleKey, str)
	if cacheControl != model.MustRevalidate {
		res, err := a.getSingleByParamRedis(ctx, key)
		if err != nil {
			if err == goredislib.Nil {
				res, err := a.getSingleByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := jsoniter.Marshal(&res)
					if err != nil {
						return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = a.setRedis(ctx, key, string(dataStr))
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

	res, err := a.getSingleByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := jsoniter.Marshal(&res)
		if err != nil {
			return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
		}
		err = a.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
	}
	return res, nil
}

func (a *AccountRole) Update(ctx context.Context, AccountRole *entity.AccountRole) error {
	return a.updatePSQL(ctx, AccountRole)
}

func (a *AccountRole) Delete(ctx context.Context, AccountRole *entity.AccountRole, id int64, isHardDelete bool) error {
	return a.deletePSQL(ctx, AccountRole, id, isHardDelete)
}
func (a *AccountRole) GetByParam(ctx context.Context, cacheControl string, param *model.GetAccountRolesByParam) (entity.AccountRoleSlice, model.Pagination, error) {
	var pg model.Pagination
	var res entity.AccountRoleSlice

	str, err := jsoniter.Marshal(param)
	if err != nil {
		return entity.AccountRoleSlice{}, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetByParamAccountRoleKey, str)
	keyPg := fmt.Sprintf(model.GetByParamAccountRolePgKey, str)
	if cacheControl != model.MustRevalidate {
		res, err1 := a.getByParamRedis(ctx, key)
		pg, err2 := a.getByParamPaginationRedis(ctx, keyPg)
		if err1 != nil || err2 != nil {
			if err1 == goredislib.Nil || err2 == goredislib.Nil {
				res, pg, err := a.getByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := jsoniter.Marshal(&res)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = a.setRedis(ctx, key, string(dataStr))
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
					}
					dataStr, err = jsoniter.Marshal(&pg)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = a.setRedis(ctx, key, string(dataStr))
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
					}
				}
				return res, pg, err
			}
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error marshal param")
		}
		return res, pg, nil
	}

	res, pg, err = a.getByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := jsoniter.Marshal(&res)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
		}
		err = a.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
		dataStr, err = jsoniter.Marshal(&pg)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
		}
		err = a.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
	}
	return res, pg, err
}
