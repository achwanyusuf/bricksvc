package transfer

import (
	"context"
	"database/sql"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (t *Transfer) insertPSQL(ctx context.Context, data *entity.TransferJob) error {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorTransaction, err, "error begin transaction")
	}

	err = data.Insert(ctx, tx, boil.Infer())
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			logger.Log.Warn(errormsg.WrapErr(svcerr.BrickSVCPSQLErrorRollback, err, "error rollback"))
		}
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorInsert, err, "error insert")
	}
	err = tx.Commit()
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error commit")
	}
	return nil
}

func (t *Transfer) getSingleByParamPSQL(ctx context.Context, param *model.GetTransferJobByParam) (entity.TransferJob, error) {
	var res entity.TransferJob
	qr := param.GetQuery()
	transferJob, err := entity.TransferJobs(qr...).One(ctx, t.DB)
	if err == sql.ErrNoRows {
		return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get transferJobs")
	}

	if err != nil {
		return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get transferJobs")
	}

	return *transferJob, nil
}

func (t *Transfer) updatePSQL(ctx context.Context, transferJob *entity.TransferJob) error {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorTransaction, err, "error begin transaction")
	}

	_, err = transferJob.Update(ctx, tx, boil.Infer())
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			logger.Log.Warn(errormsg.WrapErr(svcerr.BrickSVCPSQLErrorRollback, err, "error rollback"))
		}
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorUpdate, err, "error update")
	}
	err = tx.Commit()
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error commit")
	}
	return nil
}

func (t *Transfer) deletePSQL(ctx context.Context, transferJob *entity.TransferJob, id int64, isHardDelete bool) error {
	tx, err := t.DB.BeginTx(ctx, nil)
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorTransaction, err, "error begin transaction")
	}

	_, err = transferJob.Delete(ctx, tx, isHardDelete)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			logger.Log.Warn(errormsg.WrapErr(svcerr.BrickSVCPSQLErrorRollback, err, "error rollback"))
		}
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorUpdate, err, "error delete")
	}

	if !isHardDelete {
		transferJob.DeletedBy = null.NewInt(int(id), true)
		_, err = transferJob.Update(ctx, tx, boil.Infer())
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				logger.Log.Warn(errormsg.WrapErr(svcerr.BrickSVCPSQLErrorRollback, err, "error rollback"))
			}
			return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorUpdate, err, "error update")
		}
	}
	err = tx.Commit()
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error commit")
	}
	return nil
}

func (t *Transfer) getByParamPSQL(ctx context.Context, param *model.GetTransferJobsByParam) (entity.TransferJobSlice, model.Pagination, error) {
	var totalPages int64 = 1
	if param.Limit == 0 {
		param.Limit = int64(t.Conf.DefaultPageLimit)
	}

	if param.Page == 0 {
		param.Page = 1
	}

	qr := param.GetQuery()
	count, err := entity.TransferJobs(qr...).Count(ctx, t.DB)
	if err != nil {
		return entity.TransferJobSlice{}, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error count data")
	}
	qr = append(qr, qm.Offset(int((param.Page-1)*param.Limit)))
	qr = append(qr, qm.Limit(int(param.Limit)))
	transferJobs, err := entity.TransferJobs(qr...).All(ctx, t.DB)
	if err == sql.ErrNoRows {
		return transferJobs, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get transferJobs")
	}
	if err != nil {
		return transferJobs, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get transferJobs")
	}
	if count > 0 {
		totalPages = (count / param.Limit) + 1
	}
	return transferJobs, model.Pagination{
		CurrentPage:     param.Page,
		CurrentElements: int64(len(transferJobs)),
		TotalElements:   count,
		TotalPages:      totalPages,
		SortBy:          param.OrderBy.String,
	}, nil
}
