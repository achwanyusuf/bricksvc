package account

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

func (a *Account) insertPSQL(ctx context.Context, data *entity.Account) error {
	tx, err := a.DB.BeginTx(ctx, nil)
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

func (a *Account) getSingleByParamPSQL(ctx context.Context, param *model.GetAccountByParam) (entity.Account, error) {
	var res entity.Account
	qr := param.GetQuery()
	account, err := entity.Accounts(qr...).One(ctx, a.DB)
	if err == sql.ErrNoRows {
		return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get accounts")
	}

	if err != nil {
		return res, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get accounts")
	}

	return *account, nil
}

func (a *Account) updatePSQL(ctx context.Context, account *entity.Account) error {
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorTransaction, err, "error begin transaction")
	}

	_, err = account.Update(ctx, tx, boil.Infer())
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			logger.Log.Warn(ctx, errormsg.WrapErr(svcerr.BrickSVCPSQLErrorRollback, err, "error rollback"))
		}
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorUpdate, err, "error update")
	}
	err = tx.Commit()
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error commit")
	}
	return nil
}

func (a *Account) deletePSQL(ctx context.Context, account *entity.Account, id int64, isHardDelete bool) error {
	tx, err := a.DB.BeginTx(ctx, nil)
	if err != nil {
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorTransaction, err, "error begin transaction")
	}

	_, err = account.Delete(ctx, tx, isHardDelete)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			logger.Log.Warn(ctx, errormsg.WrapErr(svcerr.BrickSVCPSQLErrorRollback, err, "error rollback"))
		}
		return errormsg.WrapErr(svcerr.BrickSVCPSQLErrorUpdate, err, "error delete")
	}

	if !isHardDelete {
		account.DeletedBy = null.NewInt(int(id), true)
		_, err = account.Update(ctx, tx, boil.Infer())
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				logger.Log.Warn(ctx, errormsg.WrapErr(svcerr.BrickSVCPSQLErrorRollback, err, "error rollback"))
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

func (a *Account) getByParamPSQL(ctx context.Context, param *model.GetAccountsByParam) (entity.AccountSlice, model.Pagination, error) {
	var totalPages int64 = 1
	if param.Limit == 0 {
		param.Limit = int64(a.Conf.DefaultPageLimit)
	}

	if param.Page == 0 {
		param.Page = 1
	}

	qr := param.GetQuery()
	count, err := entity.Accounts(qr...).Count(ctx, a.DB)
	if err != nil {
		return entity.AccountSlice{}, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error count data")
	}
	qr = append(qr, qm.Offset(int((param.Page-1)*param.Limit)))
	qr = append(qr, qm.Limit(int(param.Limit)))
	accounts, err := entity.Accounts(qr...).All(ctx, a.DB)
	if err == sql.ErrNoRows {
		return accounts, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get accounts")
	}
	if err != nil {
		return accounts, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get accounts")
	}
	if count > 0 {
		totalPages = (count / param.Limit) + 1
	}
	return accounts, model.Pagination{
		CurrentPage:     param.Page,
		CurrentElements: int64(len(accounts)),
		TotalElements:   count,
		TotalPages:      totalPages,
		SortBy:          param.OrderBy.String,
	}, nil
}
