package accountrole

import (
	"context"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/repository/accountrole"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/volatiletech/null/v8"
)

type AccountRole struct {
	log         logger.LoggerInterface
	conf        Conf
	accountRole accountrole.AccountRoleInterface
}

type Conf struct{}

type AccountRoleInterface interface {
	Create(ctx context.Context, v model.CreateAccountRole) (model.AccountRole, error)
	GetByParam(ctx context.Context, cacheControl string, v model.GetAccountRolesByParam) ([]model.AccountRole, model.Pagination, error)
	GetByID(ctx context.Context, cacheControl string, id int64) (model.AccountRole, error)
	DeleteByID(ctx context.Context, id int64, isHardDelete bool, vid int64) error
}

func New(conf Conf, logger *logger.LoggerInterface, accountRole accountrole.AccountRoleInterface) AccountRoleInterface {
	return &AccountRole{
		conf:        conf,
		log:         *logger,
		accountRole: accountRole,
	}
}

func (a *AccountRole) Create(ctx context.Context, v model.CreateAccountRole) (model.AccountRole, error) {
	var result model.AccountRole
	err := v.Validate()
	if err != nil {
		return result, err
	}

	role := &entity.AccountRole{
		AccountID: int(v.AccountID),
		RoleID:    int(v.RoleID),
		CreatedBy: int(v.CreatedBy),
		UpdatedBy: int(v.CreatedBy),
	}

	err = a.accountRole.Insert(ctx, role)
	if err != nil {
		return result, err
	}

	return model.TransformPSQLSingleAccountRole(role), nil
}

func (a *AccountRole) GetByParam(ctx context.Context, cacheControl string, v model.GetAccountRolesByParam) ([]model.AccountRole, model.Pagination, error) {
	accountRoleSlice, pagination, err := a.accountRole.GetByParam(ctx, cacheControl, &v)
	if err != nil {
		return []model.AccountRole{}, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get by param")
	}
	return model.TransformPSQLAccountRole(&accountRoleSlice), pagination, nil
}

func (a *AccountRole) GetByID(ctx context.Context, cacheControl string, id int64) (model.AccountRole, error) {
	accountRole, err := a.accountRole.GetSingleByParam(ctx, cacheControl, &model.GetAccountRoleByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.AccountRole{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "data not found")
	}
	return model.TransformPSQLSingleAccountRole(&accountRole), nil
}

func (a *AccountRole) DeleteByID(ctx context.Context, id int64, isHardDelete bool, vid int64) error {
	accountRole, err := a.accountRole.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountRoleByParam{
		ID: null.NewInt64(vid, true),
	})
	if err != nil {
		return err
	}
	return a.accountRole.Delete(ctx, &accountRole, id, isHardDelete)
}
