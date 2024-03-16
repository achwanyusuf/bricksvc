package account

import (
	"context"
	"time"

	"github.com/achwanyusuf/bricksvc/src/domain/entity"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/repository/account"
	"github.com/achwanyusuf/bricksvc/src/repository/accountrole"
	"github.com/achwanyusuf/bricksvc/src/repository/role"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/hash"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/golang-jwt/jwt"
	"github.com/volatiletech/null/v8"
)

type Account struct {
	log         logger.LoggerInterface
	conf        Conf
	account     account.AccountInterface
	role        role.RoleInterface
	accountRole accountrole.AccountRoleInterface
}

type Conf struct {
	TokenTimeout time.Duration `mapstructure:"token_timeout"`
	TokenSecret  string        `mapstructure:"token_secret"`
	AESSecret    string        `mapstructure:"aes_secret"`
}

type AccountInterface interface {
	Oauth2(ctx context.Context, v model.Login) (model.Auth, error)
	Create(ctx context.Context, v model.Register) (model.Account, error)
	GetByParam(ctx context.Context, cacheControl string, v model.GetAccountsByParam) ([]model.Account, model.Pagination, error)
	GetByID(ctx context.Context, cacheControl string, id int64) (model.Account, error)
	UpdateByID(ctx context.Context, id int64, v model.UpdateAccountData) (model.Account, error)
	UpdatePasswordByID(ctx context.Context, id int64, v model.UpdatePasswordData) (model.Account, error)
	DeleteByID(ctx context.Context, id int64, isHardDelete bool, vid int64) error
}

func New(conf Conf, logger *logger.LoggerInterface, account account.AccountInterface, role role.RoleInterface, accountRole accountrole.AccountRoleInterface) AccountInterface {
	return &Account{
		conf:        conf,
		log:         *logger,
		account:     account,
		role:        role,
		accountRole: accountRole,
	}
}

func (a *Account) Oauth2(ctx context.Context, v model.Login) (model.Auth, error) {
	var auth model.Auth
	err := v.Validate()
	if err != nil {
		return auth, err
	}
	role, err := a.role.GetSingleByParam(ctx, "", &model.GetRoleByParam{
		Cid: null.NewString(v.ClientID, true),
	})
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, err, "role not found")
	}

	if match := hash.CompareAES(role.Sec, a.conf.AESSecret, v.ClientSecret); !match {
		return auth, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, err, "invalid client id/client secret")
	}

	account, err := a.account.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountByParam{
		Email: null.NewString(v.Email, true),
	})
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, err, "account not found")
	}

	_, err = a.accountRole.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountRoleByParam{
		AccountID: null.NewInt64(int64(account.ID), true),
		RoleID:    null.NewInt64(int64(role.ID), true),
	})
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, err, "invalid client id/client secret")
	}

	err = hash.Compare(account.Password, v.Password)
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.BrickSVCInvalidPasswordNotMatch, err, "password not match")
	}

	token := jwt.New(jwt.SigningMethodHS512)
	expired := time.Now().Add(a.conf.TokenTimeout)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = account.ID
	claims["username"] = account.Email
	claims["exp"] = expired.Unix()
	claims["scope"] = role.Scope
	t, err := token.SignedString([]byte(a.conf.TokenSecret))
	if err != nil {
		return auth, errormsg.WrapErr(svcerr.BrickSVCInvalidPasswordNotMatch, err, "invalid token")
	}

	auth = model.Auth{
		AccessToken: t,
		Exp:         &expired,
		TokenType:   model.TokenTypeBearer,
		Scope:       role.Scope,
	}

	return auth, nil
}

func (a *Account) Create(ctx context.Context, v model.Register) (model.Account, error) {
	var result model.Account
	err := v.Validate()
	if err != nil {
		return result, err
	}

	pwd, err := hash.Hash(v.Password)
	if err != nil {
		return result, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error hash password")
	}

	apiKey := hash.SHA(v.Email + time.Now().Format(time.RFC3339Nano))

	account := &entity.Account{
		Name:      v.Name,
		Email:     v.Email,
		Password:  pwd,
		APIKey:    null.StringFrom(apiKey),
		CreatedBy: int(v.CreatedBy),
		UpdatedBy: int(v.CreatedBy),
	}

	err = a.account.Insert(ctx, account)
	if err != nil {
		return result, err
	}

	return model.TransformPSQLSingleAccount(account), nil
}

func (a *Account) GetByParam(ctx context.Context, cacheControl string, v model.GetAccountsByParam) ([]model.Account, model.Pagination, error) {
	accountSlice, pagination, err := a.account.GetByParam(ctx, cacheControl, &v)
	if err != nil {
		return []model.Account{}, model.Pagination{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error get by param")
	}
	return model.TransformPSQLAccount(&accountSlice), pagination, nil
}

func (a *Account) GetByID(ctx context.Context, cacheControl string, id int64) (model.Account, error) {
	account, err := a.account.GetSingleByParam(ctx, cacheControl, &model.GetAccountByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Account{}, errormsg.WrapErr(svcerr.BrickSVCNotFound, err, "data not found")
	}
	return model.TransformPSQLSingleAccount(&account), nil
}

func (a *Account) UpdateByID(ctx context.Context, id int64, v model.UpdateAccountData) (model.Account, error) {
	if err := v.Validate(); err != nil {
		return model.Account{}, err
	}

	account, err := a.account.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Account{}, err
	}

	if v.Name == account.Name {
		return model.TransformPSQLSingleAccount(&account), nil
	}

	account.Name = v.Name
	account.UpdatedBy = int(v.UpdateBy)

	err = a.account.Update(ctx, &account)
	if err != nil {
		return model.Account{}, err
	}

	return model.TransformPSQLSingleAccount(&account), nil
}

func (a *Account) UpdatePasswordByID(ctx context.Context, id int64, v model.UpdatePasswordData) (model.Account, error) {
	if err := v.IsValid(); err != nil {
		return model.Account{}, err
	}

	account, err := a.account.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountByParam{
		ID: null.NewInt64(id, true),
	})
	if err != nil {
		return model.Account{}, err
	}
	pwd, err := hash.Hash(v.Password)
	if err != nil {
		return model.Account{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error hash password")
	}

	account.Password = pwd
	account.UpdatedBy = int(v.UpdateBy)

	err = a.account.Update(ctx, &account)
	if err != nil {
		return model.Account{}, err
	}
	return model.TransformPSQLSingleAccount(&account), nil
}

func (a *Account) DeleteByID(ctx context.Context, id int64, isHardDelete bool, vid int64) error {
	account, err := a.account.GetSingleByParam(ctx, model.MustRevalidate, &model.GetAccountByParam{
		ID: null.NewInt64(vid, true),
	})
	if err != nil {
		return err
	}
	return a.account.Delete(ctx, &account, id, isHardDelete)
}
