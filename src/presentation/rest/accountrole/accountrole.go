package accountrole

import (
	"net/http"
	"strconv"

	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/response"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/usecase/accountrole"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/httpserver"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type AccountRole struct {
	log         logger.LoggerInterface
	accountrole accountrole.AccountRoleInterface
	conf        Conf
}

type Conf struct {
	TokenSecret  string `mapstructure:"token_secret"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type AccountRoleInterface interface {
	Create(ctx *fiber.Ctx) error
	Read(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	DeleteByID(ctx *fiber.Ctx) error
}

func New(conf Conf, log *logger.LoggerInterface, accountrole accountrole.AccountRoleInterface) AccountRoleInterface {
	return &AccountRole{
		conf:        conf,
		log:         *log,
		accountrole: accountrole,
	}
}

// Create AccountRole godoc
// @Summary Create AccountRole
// @Description Create account role data
// @Tags account-role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param data body model.CreateAccountRole true "AccountRole Data"
// @Success 200 {object} response.SingleAccountRoleResponse
// @Success 400 {object} response.SingleAccountRoleResponse
// @Success 500 {object} response.SingleAccountRoleResponse
// @Router /account-role [post]
func (a *AccountRole) Create(ctx *fiber.Ctx) error {
	var (
		roleData model.CreateAccountRole
		result   model.AccountRole
		response response.SingleAccountRoleResponse
	)

	if err := ctx.BodyParser(&roleData); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
	}
	userData := httpserver.GetUserData(ctx)
	roleData.CreatedBy = userData.ID
	result, err := a.accountrole.Create(ctx.Context(), roleData)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusCreated, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusCreated, nil)
}

// Get AccountRoles Data godoc
// @Summary Get account roles data
// @Description Get account roles data
// @Tags account-role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id query string false "search by id"
// @Param account_id query int false "search by account id"
// @Param role_id query int false "search by role id"
// @Param sort_by query string false "sort result by attributes"
// @Param page query int false " "
// @Param limit query int false " "
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} response.AccountRolesResponse
// @Success 400 {object} response.AccountRolesResponse
// @Success 500 {object} response.AccountRolesResponse
// @Router /account-role [get]
func (a *AccountRole) Read(ctx *fiber.Ctx) error {
	var (
		param    model.GetAccountRolesByParam
		response response.AccountRolesResponse
		header   model.Header
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	if err := ctx.QueryParser(&param); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal query"))
	}
	roles, pagination, err := a.accountrole.GetByParam(ctx.Context(), header.CacheControl, param)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = roles
	response.Pagination = pagination

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Get AccountRoles Data godoc
// @Summary Get account role by id data
// @Description Get account role by id data
// @Tags account-role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "get by id"
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} response.SingleAccountRoleResponse
// @Success 400 {object} response.SingleAccountRoleResponse
// @Success 500 {object} response.SingleAccountRoleResponse
// @Router /account-role/{id} [get]
func (a *AccountRole) GetByID(ctx *fiber.Ctx) error {
	var (
		header   model.Header
		response response.SingleAccountRoleResponse
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
	}
	result, err := a.accountrole.GetByID(ctx.Context(), header.CacheControl, id)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Delete AccountRole Data godoc
// @Summary Delete account role data
// @Description Delete account role data
// @Tags account-role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "delete by id"
// @Success 200 {object} response.EmptyResponse
// @Success 400 {object} response.EmptyResponse
// @Success 500 {object} response.EmptyResponse
// @Router /account-role/{id} [delete]
func (a *AccountRole) DeleteByID(ctx *fiber.Ctx) error {
	var (
		response response.EmptyResponse
	)
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
	}
	userData := httpserver.GetUserData(ctx)
	scope := userData.Scope
	if scope != model.SuperAdminScope {
		id = userData.ID
	}
	err = a.accountrole.DeleteByID(ctx.Context(), userData.ID, false, id)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}
