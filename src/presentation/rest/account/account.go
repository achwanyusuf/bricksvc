package account

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/response"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/usecase/account"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/httpserver"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type Account struct {
	log     logger.LoggerInterface
	account account.AccountInterface
	conf    Conf
}

type Conf struct {
	TokenSecret  string `mapstructure:"token_secret"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type AccountInterface interface {
	Oauth2(ctx *fiber.Ctx) error
	CurrentAccount(ctx *fiber.Ctx) error
	UpdateCurrentAccount(ctx *fiber.Ctx) error
	UpdatePasswordAccount(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Read(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	UpdateByID(ctx *fiber.Ctx) error
	DeleteByID(ctx *fiber.Ctx) error
}

func New(conf Conf, log *logger.LoggerInterface, acc account.AccountInterface) AccountInterface {
	return &Account{
		conf:    conf,
		log:     *log,
		account: acc,
	}
}

// Oauth2 godoc
// @Summary OAUTH2 Authorization
// @Description OAUTH2 Authorization Code flow will show generated token to access apps
// @Tags account
// @Accept x-www-form-urlencoded
// @Produce json
// @Param client_id header string true "Client ID"
// @Param client_secret header string true "Client Secret"
// @Param username formData string true "Account Email"
// @Param password formData string true "Account Password"
// @Success 200 {object} response.LoginResponse
// @Success 400 {object} response.LoginResponse
// @Success 401 {object} response.LoginResponse
// @Success 500 {object} response.LoginResponse
// @Router /oauth2 [post]
func (a *Account) Oauth2(ctx *fiber.Ctx) error {
	var (
		response               response.LoginResponse
		header                 model.Header
		clientID, clientSecret string
	)

	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}

	if header.Authorization == "" {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, nil, "error unmarshal body"))
	}
	clientID, clientSecret = a.decodeClient(ctx, header.Authorization)

	loginData := model.Login{
		Email:        ctx.FormValue("username"),
		Password:     ctx.FormValue("password"),
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	if loginData.Email == "" {
		if err := ctx.BodyParser(&loginData); err != nil {
			return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
		}
	}

	auth, err := a.account.Oauth2(ctx.Context(), loginData)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Auth = auth

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

func (a *Account) decodeClient(ctx *fiber.Ctx, auth string) (string, string) {
	client := strings.Split(auth, " ")
	if len(client) < 2 {
		return "", ""
	}
	rawDecodedText, err := base64.StdEncoding.DecodeString(client[1])
	if err != nil {
		a.log.WarnWithContext(ctx, err)
	}
	secret := strings.Split(string(rawDecodedText), ":")
	if len(secret) < 2 {
		return "", ""
	}
	return secret[0], secret[1]
}

// Current Account godoc
// @Summary Get current account data
// @Description Get current account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Success 200 {object} response.SingleAccountResponse
// @Success 400 {object} response.SingleAccountResponse
// @Success 401 {object} response.SingleAccountResponse
// @Success 500 {object} response.SingleAccountResponse
// @Router /me [get]
func (a *Account) CurrentAccount(ctx *fiber.Ctx) error {
	var (
		header   model.Header
		response response.SingleAccountResponse
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	userData := httpserver.GetUserData(ctx)
	result, err := a.account.GetByID(ctx.Context(), header.CacheControl, userData.ID)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Update Current Account godoc
// @Summary Update current account data
// @Description Update current account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param data body model.UpdateAccountData true "Update Account Data"
// @Success 200 {object} response.SingleAccountResponse
// @Success 400 {object} response.SingleAccountResponse
// @Success 401 {object} response.SingleAccountResponse
// @Success 500 {object} response.SingleAccountResponse
// @Router /me [put]
func (a *Account) UpdateCurrentAccount(ctx *fiber.Ctx) error {
	var (
		updateData model.UpdateAccountData
		response   response.SingleAccountResponse
	)

	if err := ctx.BodyParser(&updateData); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
	}
	userData := httpserver.GetUserData(ctx)
	updateData.UpdateBy = userData.ID
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "invalid id"))
	}
	result, err := a.account.UpdateByID(ctx.Context(), id, updateData)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Update Password Account godoc
// @Summary Update password account data
// @Description Update password account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param data body model.UpdatePasswordData true "Update Account Data"
// @Success 200 {object} response.SingleAccountResponse
// @Success 400 {object} response.SingleAccountResponse
// @Success 401 {object} response.SingleAccountResponse
// @Success 500 {object} response.SingleAccountResponse
// @Router /me/password [put]
func (a *Account) UpdatePasswordAccount(ctx *fiber.Ctx) error {
	var (
		updateData model.UpdatePasswordData
		response   response.SingleAccountResponse
	)

	if err := ctx.BodyParser(&updateData); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
	}
	userData := httpserver.GetUserData(ctx)
	updateData.UpdateBy = userData.ID
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "invalid id"))
	}
	result, err := a.account.UpdatePasswordByID(ctx.Context(), id, updateData)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Register godoc
// @Summary Register account
// @Description Register to create access from guest
// @Tags account
// @Accept json
// @Produce json
// @Param data body model.Register true "Account Data"
// @Success 200 {object} response.RegisterResponse
// @Success 400 {object} response.RegisterResponse
// @Success 500 {object} response.RegisterResponse
// @Router /register [post]
func (a *Account) Register(ctx *fiber.Ctx) error {
	var (
		registerData model.Register
		result       model.Account
		response     response.RegisterResponse
	)
	if err := ctx.BodyParser(&registerData); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
	}

	result, err := a.account.Create(ctx.Context(), registerData)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusCreated, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusCreated, nil)
}

// Create Account godoc
// @Summary Create account
// @Description Create account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param data body model.Register true "Account Data"
// @Success 200 {object} response.RegisterResponse
// @Success 400 {object} response.RegisterResponse
// @Success 500 {object} response.RegisterResponse
// @Router /account [post]
func (a *Account) Create(ctx *fiber.Ctx) error {
	var (
		registerData model.Register
		result       model.Account
		response     response.RegisterResponse
	)

	if err := ctx.BodyParser(&registerData); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
	}
	userData := httpserver.GetUserData(ctx)
	registerData.CreatedBy = userData.ID
	result, err := a.account.Create(ctx.Context(), registerData)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusCreated, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusCreated, nil)
}

// Get Accounts Data godoc
// @Summary Get accounts data
// @Description Get accounts data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id query string false "search by id"
// @Param name query string false "search by name"
// @Param email query string false "search by email"
// @Param sort_by query string false "sort result by attributes"
// @Param page query int false " "
// @Param limit query int false " "
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} response.AccountsResponse
// @Success 400 {object} response.AccountsResponse
// @Success 500 {object} response.AccountsResponse
// @Router /account [get]
func (a *Account) Read(ctx *fiber.Ctx) error {
	var (
		param    model.GetAccountsByParam
		response response.AccountsResponse
		header   model.Header
	)

	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	if err := ctx.QueryParser(&param); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal query"))
	}
	accounts, pagination, err := a.account.GetByParam(ctx.Context(), header.CacheControl, param)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = accounts
	response.Pagination = pagination

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Get Accounts Data godoc
// @Summary Get accounts data
// @Description Get accounts data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "get by id"
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} response.SingleAccountResponse
// @Success 400 {object} response.SingleAccountResponse
// @Success 500 {object} response.SingleAccountResponse
// @Router /account/{id} [get]
func (a *Account) GetByID(ctx *fiber.Ctx) error {
	var (
		header   model.Header
		response response.SingleAccountResponse
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
	}
	result, err := a.account.GetByID(ctx.Context(), header.CacheControl, id)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Update Account Data godoc
// @Summary Update account data
// @Description Update account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "update by id"
// @Param data body model.UpdateAccountData true "Account Data"
// @Success 200 {object} response.SingleAccountResponse
// @Success 400 {object} response.SingleAccountResponse
// @Success 500 {object} response.SingleAccountResponse
// @Router /account/{id} [put]
func (a *Account) UpdateByID(ctx *fiber.Ctx) error {
	var (
		updateData model.UpdateAccountData
		response   response.SingleAccountResponse
	)

	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(errormsg.Error400, err, "error get id"))
	}

	if err := ctx.BodyParser(&updateData); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
	}
	userData := httpserver.GetUserData(ctx)
	updateData.UpdateBy = userData.ID
	scope := userData.Scope
	if scope != model.SuperAdminScope {
		id = userData.ID
	}
	result, err := a.account.UpdateByID(ctx.Context(), id, updateData)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = result

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}

// Delete Account Data godoc
// @Summary Delete account data
// @Description Delete account data
// @Tags account
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id path string true "delete by id"
// @Success 200 {object} response.EmptyResponse
// @Success 400 {object} response.EmptyResponse
// @Success 500 {object} response.EmptyResponse
// @Router /account/{id} [delete]
func (a *Account) DeleteByID(ctx *fiber.Ctx) error {
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
	err = a.account.DeleteByID(ctx.Context(), userData.ID, false, id)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	return response.Transform(ctx, a.log, http.StatusOK, nil)
}
