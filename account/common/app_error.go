package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AppError struct {
	StatusCode int      `json:"status_code"`
	RootErr    error    `json:"-"`
	Message    string   `json:"message"`
	Log        string   `json:"log"`
	Key        string   `json:"error_key"`
	Trace      []string `json:"trace"`
	SpanID     string   `json:"span_id"`
	TraceID    string   `json:"trace_id"`
}

// type stackTracer interface {
// 	StackTrace() StackTrace
// }

func NewErrorResponse(root error, msg, log, key string) *AppError {
	stack := callers()
	traces := []string{}
	// fmt.Printf("NewErrorResponse stack %v \n", stack.StackTrace())
	for _, f := range stack.StackTrace() {
		trace := fmt.Sprintf("%+s:%d\n", f, f)
		traces = append(traces, trace)
	}
	return &AppError{

		StatusCode: http.StatusBadRequest,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
		Trace:      traces,
	}
}
func NewFullErrorResponse(statusCode int, root error, msg, log, key string) *AppError {
	stack := callers()
	traces := []string{}
	// fmt.Printf("NewErrorResponse stack %v \n", stack.StackTrace())
	for _, f := range stack.StackTrace() {
		trace := fmt.Sprintf("%+s:%d\n", f, f)
		traces = append(traces, trace)
	}
	return &AppError{

		StatusCode: statusCode,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
		Trace:      traces,
	}
}
func NewUnauthorized(root error, msg, log, key string) *AppError {
	stack := callers()
	traces := []string{}
	// fmt.Printf("NewErrorResponse stack %v \n", stack.StackTrace())
	for _, f := range stack.StackTrace() {
		trace := fmt.Sprintf("%+s:%d\n", f, f)
		traces = append(traces, trace)
	}
	return &AppError{

		StatusCode: http.StatusUnauthorized,
		RootErr:    root,
		Message:    msg,
		Key:        key,
		Log:        log,
		Trace:      traces,
	}
}
func NewCusUnauthorizedError(root error, msg string, key string) *AppError {
	if root != nil {
		return NewUnauthorized(root, msg, root.Error(), key)
	}
	return NewUnauthorized(errors.New(msg), msg, msg, key)
}
func NewCustomError(root error, msg string, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, msg, root.Error(), key)
	}
	return NewErrorResponse(errors.New(msg), msg, msg, key)
}
func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}
	return e.RootErr
}
func (e *AppError) Error() string {
	return e.RootError().Error()
}
func ErrDB(err error) *AppError {
	return NewErrorResponse(err, "something went wrong with DB", err.Error(), "DB_ERROR")
}
func ErrInvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "invalid request", err.Error(), "ErrInvalidRequest")
}
func ErrInternal(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong in the server", err.Error(), "ErrInternal")
}
func ErrCannotListEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot list %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotList%s", entity),
	)
}
func ErrCannotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot delete %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotDelete%s", entity),
	)
}
func ErrCannotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot update %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotUpdate%s", entity),
	)
}
func ErrCannotGetEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot get %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotGet%s", entity),
	)
}
func ErrEntityDeleted(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("%s deleted", strings.ToLower(entity)),
		fmt.Sprintf("Err%sDeleted", entity),
	)
}
func ErrEntityExisted(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("%s already exists", strings.ToLower(entity)),
		fmt.Sprintf("Err%sAlreadyExists", entity),
	)
}
func ErrEntityNotFound(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("%s not found", strings.ToLower(entity)),
		fmt.Sprintf("Err%sNotFound", entity),
	)
}
func ErrCannotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Cannot Create %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCanNotCreate%s", entity),
	)
}
func ErrNoPermission(err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("You have no permission"),
		fmt.Sprintf("ErrNoPermission"),
	)
}
func ErrLoginNotValid(err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintf("Email and Password not valid"),
		fmt.Sprintf("ErrLoginNotValid"),
	)
}

var ErrRecordNotFound = errors.New("record not found")
