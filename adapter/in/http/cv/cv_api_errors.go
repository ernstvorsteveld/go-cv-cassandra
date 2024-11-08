package cv

import (
	"log/slog"
	"net/http"

	m "github.com/ernstvorsteveld/go-cv-cassandra/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewError(ctx *gin.Context, err error, code string, status int) {
	cId := m.GetCorrelationIdHeader(ctx)
	text := m.ExperienceErrors[code]
	slog.Debug("cv.NewError", "content", text, "correlationId", cId, "error-code", code, "error", err.Error())
	ctx.AbortWithStatusJSON(status, newError(ctx, code))
}

func NewExperienceBindError(ctx *gin.Context, err error) {
	code := "EXP0000003"
	NewError(ctx, err, code, http.StatusBadRequest)
}

func NewCreateExperienceError(ctx *gin.Context, err error) {
	code := "EXP0000004"
	NewError(ctx, err, code, http.StatusInternalServerError)
}

func NewListExperienceError(ctx *gin.Context, err error) {
	code := "EXP0000001"
	NewError(ctx, err, code, http.StatusInternalServerError)
}

func NewListExperienceMarshalError(ctx *gin.Context, err error) {
	code := "EXP0000002"
	NewError(ctx, err, code, http.StatusInternalServerError)
}

func NewGetExperienceByIdMarshalError(ctx *gin.Context, err error) {
	code := "EXP0000005"
	NewError(ctx, err, code, http.StatusInternalServerError)
}

func NewGetExperienceByIdNotFoundError(ctx *gin.Context, err error) {
	if isNotFoundError(err) {
		code := "EXP0000006"
		NewError(ctx, err, code, http.StatusNotFound)
	} else {
		code := "EXP0000007"
		NewError(ctx, err, code, http.StatusInternalServerError)
	}
}

func NewListTagsError(ctx *gin.Context, err error) {
	if isNotFoundError(err) {
		code := "TAG0000004"
		NewError(ctx, err, code, http.StatusNotFound)
	} else {
		code := "TAG0000001"
		NewError(ctx, err, code, http.StatusInternalServerError)
	}
}

func NewCreateTagMarshalError(ctx *gin.Context, err error) {
	code := "TAG0000002"
	NewError(ctx, err, code, http.StatusBadRequest)
}

func NewCreateTagError(ctx *gin.Context, err error) {
	code := "TAG0000004"
	NewError(ctx, err, code, http.StatusInternalServerError)
}

func isNotFoundError(err error) bool {
	return err != nil && err.Error() == "not found"
}

func NewGetTagByIdMarshalError(ctx *gin.Context, err error) {
	code := "TAG0000005"
	NewError(ctx, err, code, http.StatusBadRequest)
}

func NewGetTagByIdNotFoundError(ctx *gin.Context, err error) {
	if isNotFoundError(err) {
		code := "TAG0000004"
		NewError(ctx, err, code, http.StatusNotFound)
	} else {
		code := "TAG0000001"
		NewError(ctx, err, code, http.StatusInternalServerError)
	}
}

func newError(ctx *gin.Context, code string) *Error {
	return &Error{
		Code:      code,
		Message:   m.ExperienceErrors[code],
		RequestId: uuid.MustParse(m.GetCorrelationIdHeader(ctx)),
	}
}
