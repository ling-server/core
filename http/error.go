package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	openapi "github.com/go-openapi/errors"

	"github.com/ling-server/core/errors"
	"github.com/ling-server/core/log"
)

var (
	codeMap = map[string]int{
		errors.BadRequestCode:                  http.StatusBadRequest,
		errors.DIGESTINVALID:                   http.StatusBadRequest,
		errors.MANIFESTINVALID:                 http.StatusBadRequest,
		errors.UNSUPPORTED:                     http.StatusBadRequest,
		errors.UnAuthorizedCode:                http.StatusUnauthorized,
		errors.ForbiddenCode:                   http.StatusForbidden,
		errors.MethodNotAllowedCode:            http.StatusMethodNotAllowed,
		errors.DENIED:                          http.StatusForbidden,
		errors.NotFoundCode:                    http.StatusNotFound,
		errors.ConflictCode:                    http.StatusConflict,
		errors.PreconditionCode:                http.StatusPreconditionFailed,
		errors.ViolateForeignKeyConstraintCode: http.StatusPreconditionFailed,
		errors.PROJECTPOLICYVIOLATION:          http.StatusPreconditionFailed,
		errors.GeneralCode:                     http.StatusInternalServerError,
	}
)

// Error wrap HTTP status code and message as an error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error ...
func (e *Error) Error() string {
	return fmt.Sprintf("http error: code %d, message %s", e.Code, e.Message)
}

// String wraps the error msg to the well formatted error message
func (e *Error) String() string {
	data, err := json.Marshal(&e)
	if err != nil {
		return e.Message
	}
	return string(data)
}

// SendError tries to parse the HTTP status code from the specified error, envelops it into
// an error array as the error payload and returns the code and payload to the response.
// And the error is logged as well
func SendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	statusCode, errPayload, stackTrace := apiError(err)
	// the error detail is logged only, and will not be sent to the client to avoid leaking server information
	if statusCode >= http.StatusInternalServerError {
		log.Errorf("%s %s", errPayload, stackTrace)
		err = errors.New(nil).WithCode(errors.GeneralCode).WithMessage("internal server error")
		errPayload = errors.NewErrs(err).Error()
	} else {
		// only log the error whose status code < 500 when debugging to avoid log flooding
		log.Debug(errPayload)
	}
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, errPayload)
}

// generates the HTTP status code based on the specified error,
// envelops the error into an error array as the payload and return them
func apiError(err error) (statusCode int, errPayload, stackTrace string) {
	code := 0
	var openAPIErr openapi.Error
	if errors.As(err, &openAPIErr) {
		// Before executing operation handler, go-swagger will bind a parameters object to a request and validate the request,
		// it will return directly when bind and validate failed.
		// The response format of the default ServeError implementation does not match the internal error response format.
		// So we needed to convert the format to the internal error response format.
		code = int(openAPIErr.Code())
		errCode := strings.Replace(strings.ToUpper(http.StatusText(code)), " ", "_", -1)
		err = errors.New(nil).WithCode(errCode).WithMessage(openAPIErr.Error())
	} else if legacyErr, ok := err.(*Error); ok {
		// make sure the legacy error format is align with the new one
		code = legacyErr.Code
		errCode := strings.Replace(strings.ToUpper(http.StatusText(code)), " ", "_", -1)
		err = errors.New(nil).WithCode(errCode).WithMessage(legacyErr.Message)
	} else {
		code = codeMap[errors.ErrorCode(err)]
	}
	if code == 0 {
		code = http.StatusInternalServerError
	}
	fullStack := ""
	if _, ok := err.(*errors.Error); ok {
		fullStack = err.(*errors.Error).StackTrace()
	}
	return code, errors.NewErrs(err).Error(), fullStack
}
