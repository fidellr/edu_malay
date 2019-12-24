package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo"

	assembler "github.com/fidellr/edu_malay/assembler/profile"
	"github.com/fidellr/edu_malay/model"
	assemblerModel "github.com/fidellr/edu_malay/model/assembler"
	"github.com/fidellr/edu_malay/utils"
)

type Handler struct {
	service assembler.ProfileAssemblerUsecase
}

func NewProfileAssemblerHandler(e *echo.Echo, service assembler.ProfileAssemblerUsecase) {
	handler := &Handler{
		service,
	}

	e.POST("/assemble-profile/:clc_id", handler.Create)
	e.GET("/assemble-profile", handler.FetchAll)
	e.POST("/assemble-profile/drop/:id", handler.Remove)
}

func (h *Handler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	payload := new(assemblerModel.ProfileAssemblerParam)
	if err := c.Bind(payload); err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	err := h.service.Create(ctx, c.Param("clc_id"), payload)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) FetchAll(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := h.service.FetchAll(ctx)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Remove(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := h.service.Remove(ctx, c.Param("id"))
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), model.ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
