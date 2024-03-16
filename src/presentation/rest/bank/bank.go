package bank

import (
	"net/http"

	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/response"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/src/usecase/bank"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type Bank struct {
	log  logger.LoggerInterface
	bank bank.BankInterface
	conf Conf
}

type Conf struct{}

type BankInterface interface {
	GetBankAccount(ctx *fiber.Ctx) error
}

func New(conf Conf, log *logger.LoggerInterface, b bank.BankInterface) BankInterface {
	return &Bank{
		conf: conf,
		log:  *log,
		bank: b,
	}
}

// Get Bank Data godoc
// @Summary Get Bank data
// @Description Get banks data
// @Tags bank
// @Accept json
// @Produce json
// @Security APIKey
// @Param bank_id query string true "search by bank id"
// @Param account_name query string true "search by account name"
// @Param account_number query string true "search by account number"
// @Param Cache-Control header string false "Request Cache Control" Enums(must-revalidate, none)
// @Success 200 {object} response.GetBankAccountResponse
// @Success 400 {object} response.GetBankAccountResponse
// @Success 500 {object} response.GetBankAccountResponse
// @Router /bank [get]
func (a *Bank) GetBankAccount(ctx *fiber.Ctx) error {
	var (
		param    model.GetBankAccount
		header   model.Header
		response response.GetBankAccountResponse
	)
	if err := ctx.ReqHeaderParser(&header); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal header"))
	}
	if err := ctx.QueryParser(&param); err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err, "error unmarshal query"))
	}
	bank, err := a.bank.GetBankAccount(ctx.Context(), header.CacheControl, param, header.APIKey)
	if err != nil {
		return response.Transform(ctx, a.log, http.StatusOK, err)
	}

	response.Data = bank
	return response.Transform(ctx, a.log, http.StatusOK, nil)
}
