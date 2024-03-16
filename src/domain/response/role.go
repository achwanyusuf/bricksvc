package response

import (
	"time"

	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type SingleRoleResponse struct {
	Response
	Data model.Role `json:"data"`
}

func (r *SingleRoleResponse) Transform(ctx *fiber.Ctx, log logger.LoggerInterface, code int, err error) error {
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

type RolesResponse struct {
	Response
	Data       []model.Role     `json:"data"`
	Pagination model.Pagination `json:"pagination"`
}

func (r *RolesResponse) Transform(ctx *fiber.Ctx, log logger.LoggerInterface, code int, err error) error {
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
		r.Data = []model.Role{}
	}

	return ctx.Status(int(r.Response.Code)).JSON(r)
}

type SingleAccountRoleResponse struct {
	Response
	Data model.AccountRole `json:"data"`
}

func (r *SingleAccountRoleResponse) Transform(ctx *fiber.Ctx, log logger.LoggerInterface, code int, err error) error {
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

type AccountRolesResponse struct {
	Response
	Data       []model.AccountRole `json:"data"`
	Pagination model.Pagination    `json:"pagination"`
}

func (r *AccountRolesResponse) Transform(ctx *fiber.Ctx, log logger.LoggerInterface, code int, err error) error {
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
		r.Data = []model.AccountRole{}
	}

	return ctx.Status(int(r.Response.Code)).JSON(r)
}
