package role

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

type Role struct {
	DB    *sql.DB
	Redis *goredislib.Client
	Conf  Conf
}

type Conf struct {
	DefaultPageLimit    int           `mapstructure:"page_limit"`
	RedisExpirationTime time.Duration `mapstructure:"expiration_time"`
}

type RoleInterface interface {
	Insert(ctx context.Context, data *entity.Role) error
	GetSingleByParam(ctx context.Context, cacheControl string, param *model.GetRoleByParam) (entity.Role, error)
	Update(ctx context.Context, v *entity.Role) error
	Delete(ctx context.Context, v *entity.Role, id int64, isHardDelete bool) error
	GetByParam(ctx context.Context, cacheControl string, param *model.GetRolesByParam) (entity.RoleSlice, model.Pagination, error)
}

func New(conf Conf, db *sql.DB, rds *goredislib.Client) RoleInterface {
	return &Role{
		DB:    db,
		Redis: rds,
		Conf:  conf,
	}
}

func (r *Role) Insert(ctx context.Context, data *entity.Role) error {
	return r.insertPSQL(ctx, data)
}

func (r *Role) GetSingleByParam(ctx context.Context, cacheControl string, param *model.GetRoleByParam) (entity.Role, error) {
	str, err := jsoniter.Marshal(param)
	if err != nil {
		return entity.Role{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetSingleByParamRoleKey, str)
	if cacheControl != model.MustRevalidate {
		res, err := r.getSingleByParamRedis(ctx, key)
		if err != nil {
			if err == goredislib.Nil {
				res, err := r.getSingleByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := jsoniter.Marshal(&res)
					if err != nil {
						return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = r.setRedis(ctx, key, string(dataStr))
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

	res, err := r.getSingleByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := jsoniter.Marshal(&res)
		if err != nil {
			return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
		}
		err = r.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
	}
	return res, nil
}

func (r *Role) Update(ctx context.Context, v *entity.Role) error {
	return r.updatePSQL(ctx, v)
}

func (r *Role) Delete(ctx context.Context, v *entity.Role, id int64, isHardDelete bool) error {
	return r.deletePSQL(ctx, v, id, isHardDelete)
}
func (r *Role) GetByParam(ctx context.Context, cacheControl string, param *model.GetRolesByParam) (entity.RoleSlice, model.Pagination, error) {
	var pg model.Pagination
	var res entity.RoleSlice

	str, err := jsoniter.Marshal(param)
	if err != nil {
		return entity.RoleSlice{}, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetByParamRoleKey, str)
	keyPg := fmt.Sprintf(model.GetByParamRolePgKey, str)
	if cacheControl != model.MustRevalidate {
		res, err1 := r.getByParamRedis(ctx, key)
		pg, err2 := r.getByParamPaginationRedis(ctx, keyPg)
		if err1 != nil || err2 != nil {
			if err1 == goredislib.Nil || err2 == goredislib.Nil {
				res, pg, err := r.getByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := jsoniter.Marshal(&res)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = r.setRedis(ctx, key, string(dataStr))
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
					}
					dataStr, err = jsoniter.Marshal(&pg)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = r.setRedis(ctx, key, string(dataStr))
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

	res, pg, err = r.getByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := jsoniter.Marshal(&res)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
		}
		err = r.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
		dataStr, err = jsoniter.Marshal(&pg)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
		}
		err = r.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
	}
	return res, pg, err
}
