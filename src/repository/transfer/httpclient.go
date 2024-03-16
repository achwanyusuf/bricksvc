package transfer

import (
	"context"
	"errors"

	"github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func (t *Transfer) insertTransferHTTPClient(ctx context.Context, data clientresponse.Transfer) (clientresponse.Transfer, error) {
	var transfer clientresponse.Transfer
	client := fiber.Post(t.Conf.CreateTransactionURL)
	body, err := jsoniter.Marshal(data)
	if err != nil {
		return transfer, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}
	client.ContentType("application/json")
	client.Body(body)
	statusCode, body, errs := client.Bytes()
	if statusCode == fiber.StatusNotFound {
		return clientresponse.Transfer{}, errormsg.WrapErr(svcerr.BrickSVCNotFound, errors.Join(errs...), "status not found")
	}

	if statusCode == fiber.StatusUnauthorized {
		return clientresponse.Transfer{}, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, errors.Join(errs...), "status unauthorized")
	}

	if err != nil {
		return clientresponse.Transfer{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, errors.Join(errs...))
	}

	if err := jsoniter.Unmarshal(body, &transfer); err != nil {
		return clientresponse.Transfer{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}

	return transfer, nil
}

func (t *Transfer) getTransferHTTPClient(id string) (clientresponse.Transfer, error) {
	var transfer clientresponse.Transfer
	client := fiber.Get(t.Conf.GetTransactionURL + id)
	statusCode, body, err := client.Bytes()
	if statusCode == fiber.StatusNotFound {
		return clientresponse.Transfer{}, errormsg.WrapErr(svcerr.BrickSVCNotFound, errors.Join(err...), "status not found")
	}

	if statusCode == fiber.StatusUnauthorized {
		return clientresponse.Transfer{}, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, errors.Join(err...), "status unauthorized")
	}

	if err != nil {
		return clientresponse.Transfer{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, errors.Join(err...))
	}

	if err := jsoniter.Unmarshal(body, &transfer); err != nil {
		return clientresponse.Transfer{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}

	return transfer, nil
}
