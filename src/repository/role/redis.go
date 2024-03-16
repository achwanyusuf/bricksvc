package role

import (
	"context"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	jsoniter "github.com/json-iterator/go"
)

func (r *Role) getSingleByParamRedis(ctx context.Context, key string) (entity.Role, error) {
	var res entity.Role
	data, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *Role) setRedis(ctx context.Context, key string, data string) error {
	expTime := r.Conf.RedisExpirationTime
	if r.Conf.RedisExpirationTime == 0 {
		expTime = model.DefaultRedisExpiration
	}
	_, err := r.Redis.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	_, err = r.Redis.Set(ctx, key, data, expTime).Result()
	return err
}

func (r *Role) getByParamRedis(ctx context.Context, key string) (entity.RoleSlice, error) {
	var res entity.RoleSlice
	data, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *Role) getByParamPaginationRedis(ctx context.Context, key string) (model.Pagination, error) {
	var res model.Pagination
	data, err := r.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
