package errormsg

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var (
	Error500 Message = Message{
		Code:       50000,
		StatusCode: 500,
		Message:    "Terjadi kendala dalam sistem! Silakan hubungi admin!",
		Translation: Translation{
			EN: "Oops! There is something wrong. Please contact us!",
		},
	}

	Error400 Message = Message{
		Code:       40000,
		StatusCode: 400,
		Message:    "Kesalahan input. Silakan cek kembali masukan anda!",
		Translation: Translation{
			EN: "Invalid input. Please validate your input!",
		},
	}

	Error404 Message = Message{
		Code:       40400,
		StatusCode: 404,
		Message:    "Data tidak ditemukan!",
		Translation: Translation{
			EN: "Data not found!",
		},
	}

	ErrorPage404 Message = Message{
		Code:       40400,
		StatusCode: 404,
		Message:    "Halaman tidak ditemukan!",
		Translation: Translation{
			EN: "Page not found!",
		},
	}

	Error401 Message = Message{
		Code:       40100,
		StatusCode: 401,
		Message:    "Akses tidak diijinkan! Silakan login kembali!",
		Translation: Translation{
			EN: "Access not authorized! Please login again!",
		},
	}
	defaultMessage = "unknown error"
)

type ErrorMsg struct {
	Code           int64   `json:"code"`
	DebugError     error   `json:"debug_error"`
	WrappedMessage Message `json:"wrapped_message"`
	Line           int     `json:"int"`
	FilePath       string  `json:"file_path"`
	Func           string  `json:"func"`
	PC             int     `json:"pc"`
}

func (w *ErrorMsg) Error() string {
	return fmt.Sprint(&w)
}

type Message struct {
	Code        int64       `json:"code"`
	StatusCode  int64       `json:"status_code"`
	Message     string      `json:"message"`
	Translation Translation `json:"translation"`
}

type Translation struct {
	EN string `json:"en"`
}

func WrapErr(msg Message, err error, customMessage ...string) error {
	if err == nil {
		msg := strings.Join(customMessage, " ")
		if msg == "" {
			msg = defaultMessage
		}
		err = errors.New(msg)
	}

	if msg.Code == 0 {
		return generateError(Error500, err)
	}

	if err != nil {
		if errData, ok := err.(*ErrorMsg); ok {
			err = errData.DebugError
		}
	}

	return generateError(msg, err)
}

func generateError(msg Message, err error) *ErrorMsg {
	newErr := &ErrorMsg{
		Code:           msg.Code,
		DebugError:     err,
		WrappedMessage: msg,
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return newErr
	}
	newErr.Line = line
	newErr.FilePath = file

	f := runtime.FuncForPC(pc)
	if f == nil {
		return newErr
	}
	newErr.Func = f.Name()

	return newErr
}

func GetErrorData(err error) ErrorMsg {
	// Handle when error is nil
	if err == nil {
		return ErrorMsg{}
	}

	// Get error data
	if errData, ok := err.(*ErrorMsg); ok {
		return *errData
	}

	// Get wrapped error when error is not wrapped yet
	e := generateError(Error400, err)
	return *e
}

func WriteErr(err error) string {
	// Write empty error
	if err == nil {
		return ""
	}

	// Write when error is wrapped
	if errData, ok := err.(*ErrorMsg); ok {
		return fmt.Sprintf("%s%s:%v %v %s", errData.FilePath, errData.Func, errData.Line, errData.PC, errData.DebugError)
	}

	// Write unwrapped error
	return err.Error()
}
