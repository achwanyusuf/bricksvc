package rest

import (
	"github.com/achwanyusuf/bricksvc/src/domain/model"
	"github.com/achwanyusuf/bricksvc/src/presentation/rest/account"
	"github.com/achwanyusuf/bricksvc/src/presentation/rest/accountrole"
	"github.com/achwanyusuf/bricksvc/src/presentation/rest/bank"
	"github.com/achwanyusuf/bricksvc/src/presentation/rest/role"
	"github.com/achwanyusuf/bricksvc/src/presentation/rest/transfer"
	"github.com/achwanyusuf/bricksvc/src/usecase"
	"github.com/achwanyusuf/bricksvc/utils/httpserver"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	"github.com/gofiber/fiber/v2"
)

type Rest struct {
	Conf       Config
	Log        *logger.LoggerInterface
	Usecase    *usecase.UsecaseInterface
	HTTPServer *fiber.App
}

type Config struct {
	Account     account.Conf     `mapstructure:"account"`
	Role        role.Conf        `mapstructure:"role"`
	AccountRole accountrole.Conf `mapstructure:"account_role"`
	Bank        bank.Conf        `mapstructure:"bank"`
	Transfer    transfer.Conf    `mapstructure:"transfer"`
	TokenSecret string           `mapstructure:"token_secret"`
}

type RestInterface struct {
	Account     account.AccountInterface
	Role        role.RoleInterface
	AccountRole accountrole.AccountRoleInterface
	Bank        bank.BankInterface
	Transfer    transfer.TransferInterface
}

func New(r *Rest) *RestInterface {
	return &RestInterface{
		account.New(r.Conf.Account, r.Log, r.Usecase.Account),
		role.New(r.Conf.Role, r.Log, r.Usecase.Role),
		accountrole.New(r.Conf.AccountRole, r.Log, r.Usecase.AccountRole),
		bank.New(r.Conf.Bank, r.Log, r.Usecase.Bank),
		transfer.New(r.Conf.Transfer, r.Log, r.Usecase.Transfer),
	}
}

func (r *Rest) Serve(handler *RestInterface) {
	api := r.HTTPServer.Group("/api/v1")
	api.Post("/oauth2", handler.Account.Oauth2)
	api.Post("/register", handler.Account.Register)

	api.Get("/me", httpserver.Protected(r.Conf.TokenSecret), handler.Account.CurrentAccount)
	api.Put("/me", httpserver.Protected(r.Conf.TokenSecret), handler.Account.UpdateCurrentAccount)
	api.Put("/me/password", httpserver.Protected(r.Conf.TokenSecret), handler.Account.UpdatePasswordAccount)
	api.Post("/account", httpserver.Protected(r.Conf.TokenSecret), handler.Account.Create)
	api.Get("/account", httpserver.Protected(r.Conf.TokenSecret), httpserver.ValidateScope([]string{model.SuperAdminScope}), handler.Account.Read)
	api.Get("/account/:id", httpserver.Protected(r.Conf.TokenSecret), handler.Account.GetByID)
	api.Put("/account/:id", httpserver.Protected(r.Conf.TokenSecret), handler.Account.UpdateByID)
	api.Delete("/account/:id", httpserver.Protected(r.Conf.TokenSecret), handler.Account.DeleteByID)

	api.Post("/role", httpserver.Protected(r.Conf.TokenSecret), httpserver.ValidateScope([]string{model.SuperAdminScope}), handler.Role.Create)
	api.Get("/role", httpserver.Protected(r.Conf.TokenSecret), httpserver.ValidateScope([]string{model.SuperAdminScope}), handler.Role.Read)
	api.Get("/role/:id", httpserver.Protected(r.Conf.TokenSecret), httpserver.ValidateScope([]string{model.SuperAdminScope}), handler.Role.GetByID)
	api.Put("/role/:id", httpserver.Protected(r.Conf.TokenSecret), httpserver.ValidateScope([]string{model.SuperAdminScope}), handler.Role.UpdateByID)
	api.Delete("/role/:id", httpserver.Protected(r.Conf.TokenSecret), httpserver.ValidateScope([]string{model.SuperAdminScope}), handler.Role.DeleteByID)

	api.Post("/account-role", httpserver.Protected(r.Conf.TokenSecret), handler.AccountRole.Create)
	api.Get("/account-role", httpserver.Protected(r.Conf.TokenSecret), httpserver.ValidateScope([]string{model.SuperAdminScope}), handler.AccountRole.Read)
	api.Get("/account-role/:id", httpserver.Protected(r.Conf.TokenSecret), httpserver.ValidateScope([]string{model.SuperAdminScope}), handler.AccountRole.GetByID)
	api.Delete("/account-role/:id", httpserver.Protected(r.Conf.TokenSecret), httpserver.ValidateScope([]string{model.SuperAdminScope}), handler.AccountRole.DeleteByID)

	api.Get("/bank", handler.Bank.GetBankAccount)
	api.Post("/transfer", handler.Transfer.Transfer)
	api.Get("/transfer", handler.Transfer.Read)
	api.Get("/transfer/:job_id", handler.Transfer.GetByID)
}
