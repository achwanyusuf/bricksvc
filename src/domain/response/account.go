package response

import (
	"time"

	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type RegisterResponse struct {
	Response
	Data model.Account `json:"data"`
}

func (r *RegisterResponse) Transform(ctx *fiber.Ctx, log logger.LoggerInterface, code int, err error) error {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request().URI().String(),
			RequestMethod: ctx.Method(),
			RequestID:     ctx.Locals("requestid").(string),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.ErrorWithContext(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}
	return ctx.Status(int(r.Response.Code)).JSON(r)
}

type LoginResponse struct {
	Response
	model.Auth
}

func (r *LoginResponse) Transform(ctx *fiber.Ctx, log logger.LoggerInterface, code int, err error) error {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request().URI().String(),
			RequestMethod: ctx.Method(),
			RequestID:     ctx.Locals("requestid").(string),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.ErrorWithContext(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	return ctx.Status(int(r.Response.Code)).JSON(r)
}

type SingleAccountResponse struct {
	Response
	Data model.Account `json:"data"`
}

func (r *SingleAccountResponse) Transform(ctx *fiber.Ctx, log logger.LoggerInterface, code int, err error) error {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request().URI().String(),
			RequestMethod: ctx.Method(),
			RequestID:     ctx.Locals("requestid").(string),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.ErrorWithContext(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	return ctx.Status(int(r.Response.Code)).JSON(r)
}

type EmptyResponse struct {
	Response
}

func (r *EmptyResponse) Transform(ctx *fiber.Ctx, log logger.LoggerInterface, code int, err error) error {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request().URI().String(),
			RequestMethod: ctx.Method(),
			RequestID:     ctx.Locals("requestid").(string),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.ErrorWithContext(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	return ctx.Status(int(r.Response.Code)).JSON(r)
}

type AccountsResponse struct {
	Response
	Data       []model.Account  `json:"data"`
	Pagination model.Pagination `json:"pagination"`
}

func (r *AccountsResponse) Transform(ctx *fiber.Ctx, log logger.LoggerInterface, code int, err error) error {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request().URI().String(),
			RequestMethod: ctx.Method(),
			RequestID:     ctx.Locals("requestid").(string),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.ErrorWithContext(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	if len(r.Data) == 0 {
		r.Data = []model.Account{}
	}

	return ctx.Status(int(r.Response.Code)).JSON(r)
}
