package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/landerrors"
	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/pkg/httpclient"
	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error, requestInfo *domain.RequestInfo) {
	status, errs := mapToDIGITErrors(err)
	c.JSON(status, domain.NewErrorResponse(requestInfo, errs))
}

func mapToDIGITErrors(err error) (int, []domain.Error) {
	if err == nil {
		return http.StatusInternalServerError, []domain.Error{{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "unknown error",
		}}
	}

	var custom *landerrors.CustomException
	if errors.As(err, &custom) {
		return http.StatusBadRequest, errorsFromCustomException(custom)
	}

	var httpErr *httpclient.HTTPError
	if errors.As(err, &httpErr) {
		return http.StatusBadGateway, []domain.Error{{
			Code:        "EXTERNAL_SERVICE_EXCEPTION",
			Message:     httpErr.Error(),
			Description: httpErr.Body,
		}}
	}

	if isBadRequest(err) {
		return http.StatusBadRequest, []domain.Error{{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		}}
	}

	return http.StatusInternalServerError, []domain.Error{{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: err.Error(),
	}}
}

func isBadRequest(err error) bool {
	var syntax *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &syntax) || errors.As(err, &typeErr) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "json") || strings.Contains(msg, "binding") || strings.Contains(msg, "unmarshal")
}

func errorsFromCustomException(ce *landerrors.CustomException) []domain.Error {
	if len(ce.Errors) > 0 {
		out := make([]domain.Error, 0, len(ce.Errors))
		for code, message := range ce.Errors {
			out = append(out, domain.Error{Code: code, Message: message})
		}
		return out
	}
	return []domain.Error{{Code: ce.Code, Message: ce.Message}}
}

// ErrorMiddleware converts panics and unhandled errors; handlers call writeError directly.
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		var requestInfo *domain.RequestInfo
		if ri, ok := c.Get("requestInfo"); ok {
			if r, ok := ri.(*domain.RequestInfo); ok {
				requestInfo = r
			}
		}
		writeError(c, err, requestInfo)
	}
}
