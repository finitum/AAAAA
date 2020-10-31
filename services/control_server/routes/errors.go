package routes

import (
	"github.com/go-chi/render"
	"net/http"
)

type AppCode int64

const (
	AppCodeGeneric AppCode = iota
	AppCodeGitRepoUnreachable
	AppCodeAlreadyExists
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string  `json:"status"`          // user-level status message
	AppCode    AppCode `json:"code,omitempty"`  // application-specific error code
	ErrorText  string  `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error, code ...AppCode) render.Renderer {
	retcode := AppCodeGeneric
	if len(code) > 0 {
		retcode = code[0]
	}

	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid Request",
		AppCode:        retcode,
		ErrorText:      err.Error(),
	}
}

func ErrServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Server Error",
		ErrorText:      err.Error(),
	}
}

func ErrUnauthorized() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized",
	}
}

func ErrExists() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Object already exists",
		AppCode:        AppCodeAlreadyExists,
	}
}
