package transfer

import (
	"net/http"

	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/response"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/usecase/transfer"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type Transfer struct {
	log      logger.LoggerInterface
	transfer transfer.TransferInterface
	conf     Conf
}

type Conf struct {
	TokenSecret  string `mapstructure:"token_secret"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type TransferInterface interface {
	Transfer(ctx *fiber.Ctx) error
	Read(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
}

func New(conf Conf, log *logger.LoggerInterface, t transfer.TransferInterface) TransferInterface {
	return &Transfer{
		conf:     conf,
		log:      *log,
		transfer: t,
	}
}

// Create Transfer Data godoc
// @Summary Create Transfer data
// @Description create transfer bank
// @Tags transfer
// @Accept json
// @Produce json
// @Security APIKey
// @Param data body model.CreateTransfer true "Create Transfer Data"
// @Success 200 {object} response.SingleTransferJobResponse
// @Success 400 {object} response.SingleTransferJobResponse
// @Success 401 {object} response.SingleTransferJobResponse
// @Success 500 {object} response.SingleTransferJobResponse
// @Router /transfer [post]
func (t *Transfer) Transfer(ctx *fiber.Ctx) error {
	var (
		createTransfer model.CreateTransfer
		header         model.Header
		result         model.TransferJob
		response       response.SingleTransferJobResponse
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, t.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}

	if err := ctx.BodyParser(&createTransfer); err != nil {
		return response.Transform(ctx, t.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal body"))
	}

	result, err := t.transfer.Transfer(ctx.Context(), createTransfer, header.APIKey)
	if err != nil {
		return response.Transform(ctx, t.log, http.StatusCreated, err)
	}

	response.Data = result

	return response.Transform(ctx, t.log, http.StatusCreated, nil)
}

// Get Transfer Data godoc
// @Summary get Transfer data
// @Description get transfer bank
// @Tags transfer
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param id query int false "search by id"
// @Param job_id query string false "search by job id"
// @Param api_key query string false "search by api key"
// @Param status query string false "search by status"
// @Param sort_by query string false "sort result by attributes"
// @Param page query int false " "
// @Param limit query int false " "
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} response.TransferJobsResponse
// @Success 400 {object} response.TransferJobsResponse
// @Success 500 {object} response.TransferJobsResponse
// @Router /transfer [get]
func (t *Transfer) Read(ctx *fiber.Ctx) error {
	var (
		param    model.GetTransferJobsByParam
		header   model.Header
		response response.TransferJobsResponse
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, t.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	if err := ctx.QueryParser(&param); err != nil {
		return response.Transform(ctx, t.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal query"))
	}
	roles, pagination, err := t.transfer.GetByParam(ctx.Context(), header.CacheControl, param)
	if err != nil {
		return response.Transform(ctx, t.log, http.StatusOK, err)
	}

	response.Data = roles
	response.Pagination = pagination

	return response.Transform(ctx, t.log, http.StatusOK, nil)
}

// Get Transfer Data By Job ID godoc
// @Summary Get Transfer Data By Job ID
// @Description get transfer data by job id
// @Tags transfer
// @Accept json
// @Produce json
// @Security OAuth2Password
// @Param job_id path string true "get by job id"
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} response.SingleTransferJobResponse
// @Success 400 {object} response.SingleTransferJobResponse
// @Success 500 {object} response.SingleTransferJobResponse
// @Router /transfer/{job_id} [get]
func (t *Transfer) GetByID(ctx *fiber.Ctx) error {
	var (
		header   model.Header
		response response.SingleTransferJobResponse
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, t.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	id := ctx.Params("job_id")
	result, err := t.transfer.GetByJobID(ctx.Context(), header.CacheControl, id)
	if err != nil {
		return response.Transform(ctx, t.log, http.StatusOK, err)
	}

	response.Data = result

	return response.Transform(ctx, t.log, http.StatusOK, nil)
}
