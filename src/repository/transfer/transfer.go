package transfer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/kafkalib"
	jsoniter "github.com/json-iterator/go"
	goredislib "github.com/redis/go-redis/v9"
)

type Transfer struct {
	DB    *sql.DB
	Redis *goredislib.Client
	Conf  Conf
	Kafka kafkalib.ProducerInterface
}

type Conf struct {
	DefaultPageLimit     int           `mapstructure:"page_limit"`
	RedisExpirationTime  time.Duration `mapstructure:"expiration_time"`
	TransferTopic        string        `mapstructure:"transfer_topic"`
	CreateTransactionURL string        `mapstructure:"create_transaction_url"`
	GetTransactionURL    string        `mapstructure:"get_transaction_url"`
}

type TransferInterface interface {
	Insert(ctx context.Context, data *entity.TransferJob) error
	GetSingleByParam(ctx context.Context, cacheControl string, param *model.GetTransferJobByParam) (entity.TransferJob, error)
	Update(ctx context.Context, v *entity.TransferJob) error
	Delete(ctx context.Context, v *entity.TransferJob, id int64, isHardDelete bool) error
	GetByParam(ctx context.Context, cacheControl string, param *model.GetTransferJobsByParam) (entity.TransferJobSlice, model.Pagination, error)
	InsertTransfer(ctx context.Context, data clientresponse.Transfer) (clientresponse.Transfer, error)
	GetTransferByID(id string) (clientresponse.Transfer, error)
}

func New(conf Conf, db *sql.DB, rds *goredislib.Client, kafka kafkalib.ProducerInterface) TransferInterface {
	return &Transfer{
		DB:    db,
		Redis: rds,
		Conf:  conf,
		Kafka: kafka,
	}
}

func (t *Transfer) Insert(ctx context.Context, data *entity.TransferJob) error {
	if err := t.insertKafka(ctx, data); err != nil {
		return err
	}
	return t.insertPSQL(ctx, data)
}

func (t *Transfer) InsertTransfer(ctx context.Context, data clientresponse.Transfer) (clientresponse.Transfer, error) {
	return t.insertTransferHTTPClient(ctx, data)
}

func (t *Transfer) GetSingleByParam(ctx context.Context, cacheControl string, param *model.GetTransferJobByParam) (entity.TransferJob, error) {
	str, err := jsoniter.Marshal(param)
	if err != nil {
		return entity.TransferJob{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetSingleByParamTransferKey, str)
	if cacheControl != model.MustRevalidate {
		res, err := t.getSingleByParamRedis(ctx, key)
		if err != nil {
			if err == goredislib.Nil {
				res, err := t.getSingleByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := jsoniter.Marshal(&res)
					if err != nil {
						return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = t.setRedis(ctx, key, string(dataStr))
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

	res, err := t.getSingleByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := jsoniter.Marshal(&res)
		if err != nil {
			return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
		}
		err = t.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
	}
	return res, nil
}

func (t *Transfer) Update(ctx context.Context, transferJob *entity.TransferJob) error {
	return t.updatePSQL(ctx, transferJob)
}

func (t *Transfer) Delete(ctx context.Context, transferJob *entity.TransferJob, id int64, isHardDelete bool) error {
	return t.deletePSQL(ctx, transferJob, id, isHardDelete)
}
func (t *Transfer) GetByParam(ctx context.Context, cacheControl string, param *model.GetTransferJobsByParam) (entity.TransferJobSlice, model.Pagination, error) {
	var pg model.Pagination
	var res entity.TransferJobSlice

	str, err := jsoniter.Marshal(param)
	if err != nil {
		return entity.TransferJobSlice{}, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error marshal param")
	}

	key := fmt.Sprintf(model.GetByParamTransferKey, str)
	keyPg := fmt.Sprintf(model.GetByParamTransferPgKey, str)
	if cacheControl != model.MustRevalidate {
		res, err1 := t.getByParamRedis(ctx, key)
		pg, err2 := t.getByParamPaginationRedis(ctx, keyPg)
		if err1 != nil || err2 != nil {
			if err1 == goredislib.Nil || err2 == goredislib.Nil {
				res, pg, err := t.getByParamPSQL(ctx, param)
				if err == nil {
					dataStr, err := jsoniter.Marshal(&res)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = t.setRedis(ctx, key, string(dataStr))
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
					}
					dataStr, err = jsoniter.Marshal(&pg)
					if err != nil {
						return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
					}
					err = t.setRedis(ctx, key, string(dataStr))
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

	res, pg, err = t.getByParamPSQL(ctx, param)
	if err == nil {
		dataStr, err := jsoniter.Marshal(&res)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
		}
		err = t.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
		dataStr, err = jsoniter.Marshal(&pg)
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get psql")
		}
		err = t.setRedis(ctx, key, string(dataStr))
		if err != nil {
			return res, pg, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error set redis")
		}
	}
	return res, pg, err
}

func (t *Transfer) GetTransferByID(id string) (clientresponse.Transfer, error) {
	return t.getTransferHTTPClient(id)
}
