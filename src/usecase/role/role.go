package role

import (
	"context"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/repository/role"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/hash"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/volatiletech/null/v8"
)

type Role struct {
	log  logger.LoggerInterface
	conf Conf
	role role.RoleInterface
}

type Conf struct {
	SecretKey string `mapstructure:"secret_key"`
}

type RoleInterface interface {
	Create(ctx context.Context, v model.CreateRole) (model.Role, error)
	GetByParam(ctx context.Context, cacheControl string, v model.GetRolesByParam) ([]model.Role, model.Pagination, error)
	GetByID(ctx context.Context, cacheControl string, id int64) (model.Role, error)
	UpdateByID(ctx context.Context, id int64, v model.UpdateRole) (model.Role, error)
	DeleteByID(ctx context.Context, id int64, isHardDelete bool, vid int64) error
}

func New(conf Conf, logger *logger.LoggerInterface, role role.RoleInterface) RoleInterface {
	return &Role{
		conf: conf,
		log:  *logger,
		role: role,
	}
}

func (r *Role) Create(ctx context.Context, v model.CreateRole) (model.Role, error) {
	var result model.Role
	err := v.Validate()
	if err != nil {
		return result, err
	}

	secret, err := hash.EncAES(v.Sec, r.conf.SecretKey)
	if err != nil {
		return result, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error hash secret key")
	}

	role := &entity.Role{
		Scope:     v.Scope,
		Cid:       v.Cid,
		Sec:       secret,
		CreatedBy: int(v.CreatedBy),
		UpdatedBy: int(v.CreatedBy),
	}

	err = r.role.Insert(ctx, role)
	if err != nil {
		return result, err
	}

	return model.TransformPSQLSingleRole(role), nil
}

func (r *Role) GetByParam(ctx context.Context, cacheControl string, v model.GetRolesByParam) ([]model.Role, model.Pagination, error) {
	roleSlice, pagination, err := r.role.GetByParam(ctx, cacheControl, &v)
	if err != nil {
		return []model.Role{}, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get by param")
	}
	return model.TransformPSQLRole(&roleSlice), pagination, nil
}

func (r *Role) GetByID(ctx context.Context, cacheControl string, id int64) (model.Role, error) {
	role, err := r.role.GetSingleByParam(ctx, cacheControl, &model.GetRoleByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Role{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "data not found")
	}
	return model.TransformPSQLSingleRole(&role), nil
}

func (r *Role) UpdateByID(ctx context.Context, id int64, v model.UpdateRole) (model.Role, error) {
	role, err := r.role.GetSingleByParam(ctx, model.MustRevalidate, &model.GetRoleByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Role{}, err
	}

	if !v.Scope.Valid && !v.Cid.Valid && !v.Sec.Valid {
		return model.TransformPSQLSingleRole(&role), nil
	}

	if v.Scope.Valid {
		role.Scope = v.Scope.String
	}

	if v.Sec.Valid {
		v.Sec.String, err = hash.EncAES(v.Sec.String, r.conf.SecretKey)
		if err != nil {
			return model.Role{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error hash secret key")
		}
	}

	v.FillEntity(&role)
	role.UpdatedBy = int(v.UpdatedBy)

	err = r.role.Update(ctx, &role)
	if err != nil {
		return model.Role{}, err
	}

	return model.TransformPSQLSingleRole(&role), nil
}

func (r *Role) DeleteByID(ctx context.Context, id int64, isHardDelete bool, vid int64) error {
	role, err := r.role.GetSingleByParam(ctx, model.MustRevalidate, &model.GetRoleByParam{
		ID: null.NewInt64(vid, true),
	})
	if err != nil {
		return err
	}
	return r.role.Delete(ctx, &role, id, isHardDelete)
}
