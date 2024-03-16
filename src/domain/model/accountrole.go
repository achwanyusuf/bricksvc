package model

import (
	"strings"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	GetSingleByParamAccountRoleKey string = "gspAccountRole:%s"
	GetByParamAccountRoleKey       string = "gpAccountRole:%s"
	GetByParamAccountRolePgKey     string = "gppgAccountRole:%s"
)

type GetAccountRoleByParam struct {
	ID        null.Int64 `schema:"id" json:"id" query:"id"`
	AccountID null.Int64 `schema:"account_id" json:"account_id" query:"account_id"`
	RoleID    null.Int64 `schema:"role_id" json:"role_id" query:"role_id"`
}

func (g *GetAccountRoleByParam) GetQuery() []qm.QueryMod {
	var res []qm.QueryMod
	if g.AccountID.Valid {
		res = append(res, qm.Where("account_id=?", g.AccountID.Int64))
	}

	if g.RoleID.Valid {
		res = append(res, qm.Where("role_id=?", g.RoleID.Int64))
	}
	return res
}

type GetAccountRolesByParam struct {
	GetAccountRoleByParam
	OrderBy null.String `schema:"order_by" json:"order_by" query:"order_by"`
	Limit   int64       `schema:"limit" json:"limit" query:"limit"`
	Page    int64       `schema:"page" json:"page" query:"page"`
}

func (g *GetAccountRolesByParam) GetQuery() []qm.QueryMod {
	var res []qm.QueryMod
	if g.AccountID.Valid {
		res = append(res, qm.Where("account_id=?", g.AccountID.Int64))
	}

	if g.RoleID.Valid {
		res = append(res, qm.Where("role_id=?", g.RoleID.Int64))
	}

	if g.OrderBy.Valid {
		order := strings.Split(g.OrderBy.String, ",")
		for _, o := range order {
			res = append(res, qm.OrderBy(o))
		}
	}

	return res
}

type CreateAccountRole struct {
	AccountID int64 `json:"account_id"`
	RoleID    int64 `json:"role_id"`
	CreatedBy int64 `json:"-"`
}

func (v *CreateAccountRole) Validate() error {
	if v.AccountID == 0 {
		return errormsg.WrapErr(svcerr.BrickSVCInvalidScope, nil, "invalid account id")
	}

	if v.RoleID == 0 {
		return errormsg.WrapErr(svcerr.BrickSVCInvalidClientIDClientSecret, nil, "invalid role id")
	}
	return nil
}

type UpdateAccountRole struct {
	AccountID null.Int64 `json:"account_id"`
	RoleID    null.Int64 `json:"role_id"`
	UpdatedBy int64      `json:"-"`
}

func (v *UpdateAccountRole) FillEntity(accountRole *entity.AccountRole) {
	if v.AccountID.Valid {
		accountRole.AccountID = int(v.AccountID.Int64)
	}

	if v.RoleID.Valid {
		accountRole.RoleID = int(v.RoleID.Int64)
	}
}

type AccountRole struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	RoleID    int64 `json:"role_id"`
	BaseInformation
}

func TransformPSQLSingleAccountRole(accountRole *entity.AccountRole) AccountRole {
	creationInfo := BaseInformation{
		CreatedBy: int64(accountRole.CreatedBy),
		CreatedAt: accountRole.CreatedAt,
		UpdatedBy: int64(accountRole.UpdatedBy),
		UpdatedAt: accountRole.UpdatedAt,
		DeletedBy: int64(accountRole.DeletedBy.Int),
		DeletedAt: accountRole.DeletedAt.Time,
	}

	return AccountRole{
		ID:              int64(accountRole.ID),
		AccountID:       int64(accountRole.AccountID),
		RoleID:          int64(accountRole.RoleID),
		BaseInformation: creationInfo,
	}
}

func TransformPSQLAccountRole(role *entity.AccountRoleSlice) []AccountRole {
	var res []AccountRole
	for _, v := range *role {
		creationInfo := BaseInformation{
			CreatedBy: int64(v.CreatedBy),
			CreatedAt: v.CreatedAt,
			UpdatedBy: int64(v.UpdatedBy),
			UpdatedAt: v.UpdatedAt,
			DeletedBy: int64(v.DeletedBy.Int),
			DeletedAt: v.DeletedAt.Time,
		}

		res = append(res, AccountRole{
			ID:              int64(v.ID),
			AccountID:       int64(v.AccountID),
			RoleID:          int64(v.RoleID),
			BaseInformation: creationInfo,
		})
	}

	return res
}
