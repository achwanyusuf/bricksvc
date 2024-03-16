package role

import (
	"net/http"
	"strconv"

	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/response"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/usecase/role"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/httpserver"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type RoleDep struct {
	log  logger.LoggerInterface
	role role.RoleInterface
	conf Conf
}

type Conf struct {
	TokenSecret  string `mapstructure:"token_secret"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type RoleInterface interface {
	Create(ctx *fiber.Ctx) error
	Read(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	UpdateByID(ctx *fiber.Ctx) error
	DeleteByID(ctx *fiber.Ctx) error
}

func New(conf Conf, log *logger.LoggerInterface, role role.RoleInterface) RoleInterface {
	return &RoleDep{
		conf: conf,
		log:  *log,
		role: role,
	}
}

// Create Role godoc
// @Summary Create Role
// @Description Create role data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param data body model.CreateRole true "Role Data"
// @Success 200 {object} response.SingleRoleResponse
// @Success 400 {object} response.SingleRoleResponse
// @Success 500 {object} response.SingleRoleResponse
// @Router /role [post]
func (a *RoleDep) Create(ctx *fiber.Ctx) error {
	var (
		roleData model.CreateRole
		result   model.Role
		response response.SingleRoleResponse
	)

	if err := ctx.BodyParser(&roleData); err != nil {
		return response.Transform(ctx, a.log, http.StatusCreated, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
	}
	userData := httpserver.GetUserData(ctx)

	roleData.CreatedBy = userData.ID
	result, err := a.role.Create(ctx.Context(), roleData)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusCreated, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusCreated, nil)
}

// Get Roles Data godoc
// @Summary Get roles data
// @Description Get roles data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id query string false "search by id"
// @Param scope query string false "search by scope"
// @Param cid query string false "search by client_id"
// @Param sort_by query string false "sort result by attributes"
// @Param page query int false " "
// @Param limit query int false " "
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} response.RolesResponse
// @Success 400 {object} response.RolesResponse
// @Success 500 {object} response.RolesResponse
// @Router /role [get]
func (a *RoleDep) Read(ctx *fiber.Ctx) error {
	var (
		param    model.GetRolesByParam
		header   model.Header
		response response.RolesResponse
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	if err := ctx.QueryParser(&param); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal query"))
	}
	roles, pagination, err := a.role.GetByParam(ctx.Context(), header.CacheControl, param)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = roles
	response.Pagination = pagination

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Get Roles Data godoc
// @Summary Get roles data
// @Description Get roles data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "get by id"
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} response.SingleRoleResponse
// @Success 400 {object} response.SingleRoleResponse
// @Success 500 {object} response.SingleRoleResponse
// @Router /role/{id} [get]
func (a *RoleDep) GetByID(ctx *fiber.Ctx) error {
	var (
		header   model.Header
		response response.SingleRoleResponse
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
	}
	result, err := a.role.GetByID(ctx.Context(), header.CacheControl, id)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Update Role Data godoc
// @Summary Update role data
// @Description Update role data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "update by id"
// @Param data body model.UpdateRole true "Role Data"
// @Success 200 {object} response.SingleRoleResponse
// @Success 400 {object} response.SingleRoleResponse
// @Success 500 {object} response.SingleRoleResponse
// @Router /role/{id} [put]
func (a *RoleDep) UpdateByID(ctx *fiber.Ctx) error {
	var (
		updateData model.UpdateRole
		response   response.SingleRoleResponse
	)

	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
	}

	if err = ctx.BodyParser(&updateData); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
	}

	userData := httpserver.GetUserData(ctx)
	updateData.UpdatedBy = userData.ID
	scope := userData.Scope
	if scope != model.SuperAdminScope {
		id = userData.ID
	}
	result, err := a.role.UpdateByID(ctx.Context(), id, updateData)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Delete Role Data godoc
// @Summary Delete role data
// @Description Delete role data
// @Tags role
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "delete by id"
// @Success 200 {object} response.EmptyResponse
// @Success 400 {object} response.EmptyResponse
// @Success 500 {object} response.EmptyResponse
// @Router /role/{id} [delete]
func (a *RoleDep) DeleteByID(ctx *fiber.Ctx) error {
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
	err = a.role.DeleteByID(ctx.Context(), userData.ID, false, id)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}
