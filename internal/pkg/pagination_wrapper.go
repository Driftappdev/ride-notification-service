package pkg

import (
	gopagination "github.com/driftappdev/libpackage/resilience/pagination"
	"net/http"
)

func ParsePage(r *http.Request) (gopagination.PageParams, error) {
	return gopagination.ParsePageParams(r, gopagination.DefaultConfig())
}
func BuildPage[T any](items []T, params gopagination.PageParams, r *http.Request) gopagination.Page[T] {
	pageData, _ := gopagination.Paginate(items, params)
	return gopagination.NewPage(pageData, len(items), params, r, gopagination.DefaultConfig())
}
func RespondPage[T any](w http.ResponseWriter, page gopagination.Page[T]) error {
	return gopagination.RespondPage(w, page)
}
