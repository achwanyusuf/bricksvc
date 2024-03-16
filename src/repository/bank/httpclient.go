package bank

import (
	"errors"
	"fmt"

	"github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/domain/svcerr"
	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func (b *Bank) getBankAccountHTTPClient(v model.GetBankAccount) (clientresponse.BankAccount, error) {
	var bankAccount []clientresponse.BankAccount
	client := fiber.Get(b.Conf.GetBankAccountURL)
	client.QueryString(fmt.Sprintf("bank_id=%d", v.BankID))
	client.QueryString("account_name=" + v.AccountName)
	client.QueryString("account_number=" + v.AccountNumber)
	statusCode, body, err := client.Bytes()
	if statusCode == fiber.StatusNotFound {
		return clientresponse.BankAccount{}, errormsg.WrapErr(svcerr.BrickSVCNotFound, errors.Join(err...), "status not found")
	}

	if statusCode == fiber.StatusUnauthorized {
		return clientresponse.BankAccount{}, errormsg.WrapErr(svcerr.BrickSVCNotAuthorized, errors.Join(err...), "status unauthorized")
	}

	if err != nil {
		return clientresponse.BankAccount{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, errors.Join(err...))
	}

	if err := jsoniter.Unmarshal(body, &bankAccount); err != nil {
		return clientresponse.BankAccount{}, errormsg.WrapErr(svcerr.BrickSVCBadRequest, err)
	}

	return bankAccount[0], nil
}
