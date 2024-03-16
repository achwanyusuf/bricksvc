package transfer

import (
	"context"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	jsoniter "github.com/json-iterator/go"
)

func (t *Transfer) getSingleByParamRedis(ctx context.Context, key string) (entity.TransferJob, error) {
	var res entity.TransferJob
	data, err := t.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t *Transfer) setRedis(ctx context.Context, key string, data string) error {
	expTime := t.Conf.RedisExpirationTime
	if t.Conf.RedisExpirationTime == 0 {
		expTime = model.DefaultRedisExpiration
	}
	_, err := t.Redis.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	_, err = t.Redis.Set(ctx, key, data, expTime).Result()
	return err
}

func (t *Transfer) getByParamRedis(ctx context.Context, key string) (entity.TransferJobSlice, error) {
	var res entity.TransferJobSlice
	data, err := t.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (t *Transfer) getByParamPaginationRedis(ctx context.Context, key string) (model.Pagination, error) {
	var res model.Pagination
	data, err := t.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
