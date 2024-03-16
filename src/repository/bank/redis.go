package bank

import (
	"context"

	"github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	jsoniter "github.com/json-iterator/go"
)

func (b *Bank) getSingleByParamRedis(ctx context.Context, key string) (clientresponse.BankAccount, error) {
	var res clientresponse.BankAccount
	data, err := b.Redis.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}
	err = jsoniter.Unmarshal([]byte(data), &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (b *Bank) setRedis(ctx context.Context, key string, data string) error {
	expTime := b.Conf.RedisExpirationTime
	if b.Conf.RedisExpirationTime == 0 {
		expTime = model.DefaultRedisExpiration
	}
	_, err := b.Redis.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	_, err = b.Redis.Set(ctx, key, data, expTime).Result()
	return err
}
