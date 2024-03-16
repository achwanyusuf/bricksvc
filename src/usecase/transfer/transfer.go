package transfer

import (
	"context"
	"time"

	"github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/repository/account"
	"github.com/achwanyusuf/bricksvc/src/repository/transfer"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/volatiletech/null/v8"
)

type Transfer struct {
	log      logger.LoggerInterface
	conf     Conf
	account  account.AccountInterface
	transfer transfer.TransferInterface
}

type Conf struct {
	JobActiveDuration time.Duration `mapstructure:"job_active_duration"`
}

type TransferInterface interface {
	Create(ctx context.Context, key string, v model.CreateTransfer) error
	Transfer(ctx context.Context, v model.CreateTransfer, apikey string) (model.TransferJob, error)
	GetByParam(ctx context.Context, cacheControl string, v model.GetTransferJobsByParam) ([]model.TransferJob, model.Pagination, error)
	GetByJobID(ctx context.Context, cacheControl string, id string) (model.TransferJob, error)
	ProccessGetCallback(ctx context.Context, param *model.GetTransferJobsByParam)
}

func New(conf Conf, logger *logger.LoggerInterface, account account.AccountInterface, transfer transfer.TransferInterface) TransferInterface {
	return &Transfer{
		conf:     conf,
		log:      *logger,
		account:  account,
		transfer: transfer,
	}
}

func (t *Transfer) Create(ctx context.Context, key string, v model.CreateTransfer) error {
	cParam := v.ToClientRes()
	resClient, err := t.transfer.InsertTransfer(ctx, cParam)
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}

	job, err := t.transfer.GetSingleByParam(ctx, model.MustRevalidate, &model.GetTransferJobByParam{
		JobID: null.StringFrom(key),
	})
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}
	payload, err := jsoniter.Marshal(resClient)
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}
	job.Payload = payload
	return t.transfer.Update(ctx, &job)
}

func (t *Transfer) Transfer(ctx context.Context, v model.CreateTransfer, apikey string) (model.TransferJob, error) {
	acc, err := t.account.GetSingleByParam(ctx, "", &model.GetAccountByParam{
		APIKey: null.StringFrom(apikey),
	})
	if err != nil {
		return model.TransferJob{}, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, err)
	}

	data, err := v.ToEntity(acc.APIKey.String)
	if err != nil {
		return model.TransferJob{}, err
	}

	err = t.transfer.Insert(ctx, &data)
	if err != nil {
		return model.TransferJob{}, err
	}
	res, err := model.TransformTransferJob(data)
	if err != nil {
		return model.TransferJob{}, err
	}

	return res, nil
}

func (t *Transfer) GetByParam(ctx context.Context, cacheControl string, v model.GetTransferJobsByParam) ([]model.TransferJob, model.Pagination, error) {
	transferJobSlice, pagination, err := t.transfer.GetByParam(ctx, cacheControl, &v)
	if err != nil {
		return []model.TransferJob{}, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get by param")
	}
	res, err := model.TransformPSQLTransferJob(&transferJobSlice)
	if err != nil {
		return []model.TransferJob{}, pagination, err
	}
	return res, pagination, nil
}

func (t *Transfer) GetByJobID(ctx context.Context, cacheControl string, id string) (model.TransferJob, error) {
	transferJob, err := t.transfer.GetSingleByParam(ctx, cacheControl, &model.GetTransferJobByParam{
		JobID: null.StringFrom(id),
	})
	if err != nil {
		return model.TransferJob{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "data not found")
	}
	res, err := model.TransformPSQLSingleTransferJob(&transferJob)
	if err != nil {
		return model.TransferJob{}, err
	}
	return res, nil
}

func (t *Transfer) ProccessGetCallback(ctx context.Context, param *model.GetTransferJobsByParam) {
	tNow := time.Now().UTC().Add(-t.conf.JobActiveDuration)
	transferJobSlice, _, err := t.transfer.GetByParam(ctx, model.MustRevalidate, param)
	if err != nil {
		logger.Log.Error(errormsg.WriteErr(errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get by param")))
	}

	for _, v := range transferJobSlice {
		var payload clientresponse.Transfer
		if err := v.Payload.Unmarshal(&payload); err != nil {
			logger.Log.Error(errormsg.WriteErr(errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshall")))
			continue
		}

		if payload.ID == "" {
			v.Status = entity.TransferstatusFailed
			if err := t.transfer.Update(ctx, v); err != nil {
				logger.Log.Error(errormsg.WriteErr(errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error delete data")))
			}
			continue
		}

		cb, err := t.transfer.GetTransferByID(payload.ID)
		if err != nil {
			logger.Log.Error(errormsg.WriteErr(errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get data")))
			continue
		}

		if cb.Status == entity.TransferstatusSuccess.String() {
			v.Status = entity.TransferstatusSuccess
			if err := t.transfer.Update(ctx, v); err != nil {
				logger.Log.Error(errormsg.WriteErr(errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error delete data")))
			}
			continue
		}

		if cb.Status == entity.TransferstatusFailed.String() {
			v.Status = entity.TransferstatusFailed
			if err := t.transfer.Update(ctx, v); err != nil {
				logger.Log.Error(errormsg.WriteErr(errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error delete data")))
			}
			continue
		}

		if v.CreatedAt.Before(tNow) && v.Status == entity.TransferstatusPending {
			v.Status = entity.TransferstatusFailed
			if err := t.transfer.Update(ctx, v); err != nil {
				logger.Log.Error(errormsg.WriteErr(errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error delete data")))
			}
		}

	}
}
