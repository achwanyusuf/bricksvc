package account

import (
	"context"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	jsoniter "github.com/json-iterator/go"
)

func (a *Account) getSingleByParamRedis(ctx context.Context, key string) (entity.Account, error) {
	var res entity.Account
	data, err := a.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (a *Account) setRedis(ctx context.Context, key string, data string) error {
	expTime := a.Conf.RedisExpirationTime
	if a.Conf.RedisExpirationTime == 0 {
		expTime = model.DefaultRedisExpiration
	}
	_, err := a.Redis.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	_, err = a.Redis.Set(ctx, key, data, expTime).Result()
	return err
}

func (a *Account) getByParamRedis(ctx context.Context, key string) (entity.AccountSlice, error) {
	var res entity.AccountSlice
	data, err := a.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (a *Account) getByParamPaginationRedis(ctx context.Context, key string) (model.Pagination, error) {
	var res model.Pagination
	data, err := a.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
