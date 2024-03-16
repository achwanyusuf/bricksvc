package httpserver

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/achwanyusuf/bricksvc/utils/errormsg"
	"github.com/achwanyusuf/bricksvc/utils/findvariable"
	"github.com/achwanyusuf/bricksvc/utils/logger"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/lucsky/cuid"
)

type HTTPServerInterface interface {
	Setup()
	Get() *fiber.App
	Run() error
	Close() error
}

type Conf struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	PrintRoutes     bool          `mapstructure:"print_routes"`
	CaseSensitive   bool          `mapstructure:"case_sensitive"`
	Concurency      int           `mapstructure:"concurency"`
	ReadBufferSize  int           `mapstructure:"read_buffer_size"`
	WriteBufferSize int           `mapstructure:"write_buffer_size"`
}

type httpserver struct {
	host string
	port int
	app  *fiber.App
	log  logger.LoggerInterface
}

type reqData struct {
	Header string
	Method string
	Body   string
}

type respData struct {
	Latency    string
	StatusCode int
	Data       string
}

type AuthData struct {
	ID       int64
	Username string
	Expired  time.Time
	Scope    string
}

type transactionInfo struct {
	RequestURI    string    `json:"request_uri"`
	RequestMethod string    `json:"request_method"`
	RequestID     string    `json:"request_id"`
	Timestamp     time.Time `json:"timestamp"`
	ErrorCode     int64     `json:"error_code,omitempty"`
	Cause         string    `json:"cause,omitempty"`
}

type response struct {
	TransactionInfo transactionInfo `json:"transaction_info"`
	Code            int64           `json:"status_code"`
	Message         string          `json:"message,omitempty"`
	Translation     *translation    `json:"translation,omitempty"`
}

type translation struct {
	EN string `json:"en"`
}

func New(c *Conf, log logger.LoggerInterface) HTTPServerInterface {
	conf := fiber.Config{
		JSONEncoder:     jsoniter.Marshal,
		JSONDecoder:     jsoniter.Unmarshal,
		WriteTimeout:    c.WriteTimeout,
		ReadTimeout:     c.ReadTimeout,
		CaseSensitive:   c.CaseSensitive,
		Concurrency:     c.Concurency,
		ReadBufferSize:  c.ReadBufferSize,
		WriteBufferSize: c.WriteBufferSize,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var e *fiber.Error
			code := fiber.StatusInternalServerError
			errMsg := errormsg.Error500
			log.ErrorWithContext(ctx, errormsg.WriteErr(err))
			if errors.As(err, &e) {
				code = e.Code
			}

			switch code {
			case 401:
				errMsg = errormsg.Error401
			case 404:
				errMsg = errormsg.ErrorPage404
			case 405:
				errMsg = errormsg.ErrorPage404
			case 400:
				errMsg = errormsg.Error400
			default:
				errMsg = errormsg.Error500
			}

			resp := response{
				TransactionInfo: transactionInfo{
					RequestURI:    ctx.Request().URI().String(),
					RequestMethod: ctx.Method(),
					RequestID:     ctx.Locals("requestid").(string),
					Timestamp:     time.Now(),
					ErrorCode:     errMsg.Code,
				},
				Code:    errMsg.StatusCode,
				Message: errMsg.Message,
				Translation: &translation{
					EN: errMsg.Translation.EN,
				},
			}

			return ctx.Status(int(resp.Code)).JSON(resp)
		},
	}

	httpServer := fiber.New(conf)
	return &httpserver{
		host: c.Host,
		port: c.Port,
		app:  httpServer,
		log:  log,
	}
}

func (h *httpserver) Setup() {
	h.app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e interface{}) {
			h.log.ErrorWithContext(ctx, errormsg.WriteErr(e.(error)))

		},
	}))
	h.app.Use(compress.New())
	h.app.Use(requestid.New(requestid.Config{
		Header:    "X-Custom-Header",
		Generator: cuid.New,
	}))
	h.app.Use(func(ctx *fiber.Ctx) error {
		var (
			start time.Time
			next  func(c *fiber.Ctx) bool
		)
		if next != nil && next(ctx) {
			return ctx.Next()
		}

		start = time.Now()

		req := reqData{
			Header: ctx.Request().Header.String(),
			Method: ctx.Method(),
			Body:   string(ctx.Body()),
		}

		h.log.InfoWithContext(ctx, req)

		err := ctx.Next()
		if err != nil {
			return err
		}

		latency := time.Since(start)
		response := ctx.Response()
		body := bytes.NewBuffer(response.Body()).String()

		resp := respData{
			Latency:    latency.String(),
			StatusCode: response.StatusCode(),
			Data:       body,
		}

		if response.StatusCode() > 299 {
			h.log.ErrorWithContext(ctx, resp)
		} else {
			h.log.InfoWithContext(ctx, resp)
		}
		return nil
	})
	h.app.Use(pprof.New())
	h.app.Get("/swagger/*", swagger.HandlerDefault)
	h.app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]interface{}{
			"message": "health",
		})
	})
	h.app.Get("/live", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]interface{}{
			"message": "live",
		})
	})
	h.app.Get("/ready", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]interface{}{
			"message": "ready",
		})
	})
}

func GetUserData(ctx *fiber.Ctx) AuthData {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return AuthData{
		ID:       int64(claims["id"].(float64)),
		Username: claims["username"].(string),
		Expired:  time.Unix(int64(claims["exp"].(float64)), 0),
		Scope:    claims["scope"].(string),
	}

}

// Protected protect routes
func Protected(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(secret),
		},
		ErrorHandler: jwtError,
	})
}

func jwtError(ctx *fiber.Ctx, err error) error {
	resp := response{
		TransactionInfo: transactionInfo{
			RequestURI:    ctx.Request().URI().String(),
			RequestMethod: ctx.Method(),
			RequestID:     ctx.Locals("requestid").(string),
			Timestamp:     time.Now(),
			ErrorCode:     errormsg.Error401.Code,
		},
		Code:    errormsg.Error401.StatusCode,
		Message: errormsg.Error401.Message,
		Translation: &translation{
			EN: errormsg.Error401.Translation.EN,
		},
	}
	if err.Error() == "Missing or malformed JWT" {
		resp.Code = errormsg.Error400.StatusCode
		resp.Message = errormsg.Error400.Message
		resp.TransactionInfo.ErrorCode = errormsg.Error400.Code
		resp.Translation = (*translation)(&errormsg.Error400.Translation)
		return ctx.Status(int(resp.Code)).JSON(resp)
	}
	return ctx.Status(int(resp.Code)).JSON(resp)
}

func (h *httpserver) Get() *fiber.App {
	return h.app
}

func (h *httpserver) Run() error {
	return h.app.Listen(fmt.Sprintf("%s:%d", h.host, h.port))
}

func (h *httpserver) Close() error {
	return h.app.Shutdown()
}

func ValidateScope(scopes []string) fiber.Handler {
	resp := response{}
	return func(ctx *fiber.Ctx) error {
		userData := GetUserData(ctx)
		scope := userData.Scope
		if found := findvariable.FindStrInSlice(scope, scopes); !found {
			resp = response{
				TransactionInfo: transactionInfo{
					RequestURI:    ctx.Request().URI().String(),
					RequestMethod: ctx.Method(),
					RequestID:     ctx.Locals("requestid").(string),
					Timestamp:     time.Now(),
					ErrorCode:     errormsg.Error401.Code,
				},
				Code:    errormsg.Error401.StatusCode,
				Message: errormsg.Error401.Message,
				Translation: &translation{
					EN: errormsg.Error401.Translation.EN,
				},
			}
			return ctx.Status(int(resp.Code)).JSON(resp)
		}
		return ctx.Next()
	}
}
