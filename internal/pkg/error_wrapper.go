package pkg

import goerror "github.com/driftappdev/libpackage/goerror"

type BaseError = goerror.BaseError
type Category = goerror.Category

func categoryFromHTTP(status int) goerror.Category {
	switch status {
	case 400, 422:
		return goerror.CategoryValidation
	case 401:
		return goerror.CategoryUnauthorized
	case 403:
		return goerror.CategoryForbidden
	case 404:
		return goerror.CategoryNotFound
	case 409:
		return goerror.CategoryConflict
	case 429:
		return goerror.CategoryRateLimit
	case 503:
		return goerror.CategoryUnavailable
	case 504:
		return goerror.CategoryTimeout
	default:
		return goerror.CategoryInternal
	}
}

func NewBaseError(code string, status int, message string) *goerror.BaseError {
	return goerror.New(code, categoryFromHTTP(status), message)
}

func HTTPStatus(err error) int                               { return goerror.HTTPStatusFrom(err) }
func APIResponse(err *goerror.BaseError) goerror.APIResponse { return goerror.ToAPIResponse(err) }
