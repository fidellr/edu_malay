package http

import (
	"context"
	"net/http"
	"strconv"

	model "github.com/fidellr/edu_malay/model/teacher"
	"github.com/fidellr/edu_malay/teacher"
	"github.com/fidellr/edu_malay/utils"
	"github.com/labstack/echo"
)

type ResponseError struct {
	Message string `json:"message"`
}

type Handler struct {
	service teacher.ProfileUsecase
}

func NewTeacherProfileHandler(e *echo.Echo, service teacher.ProfileUsecase) {
	handler := &Handler{
		service,
	}

	e.GET("/teachers", handler.FindAll)
	e.GET("/teacher/:id", handler.GetByID)
	e.POST("/teacher", handler.Create)
	e.PUT("/teacher/:id", handler.Update)
}

func (h *Handler) FindAll(c echo.Context) error {
	var num int
	var err error

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	if c.QueryParam("num") != "" {
		num, err = strconv.Atoi(c.QueryParam("num"))
		if err != nil {
			return utils.ConstraintErrorf("%s", err.Error())
		}
	}

	filter := &model.Filter{
		Num:    num,
		Cursor: c.QueryParam("cursor"),
	}

	teachers, nextCursor, err := h.service.FindAll(ctx, filter)
	if err != nil {
		return utils.ConstraintErrorf("%s", err.Error())
	}

	c.Response().Header().Set("X-Cursor", nextCursor)

	return c.JSON(http.StatusOK, teachers)
}

func (h *Handler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	t := new(model.ProfileEntity)
	if err := c.Bind(t); err != nil {
		return c.JSON(utils.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	err := h.service.Create(ctx, t)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	t, err := h.service.GetByID(ctx, c.Param("id"))
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(utils.GetStatusCode(err), t)
}

func (h *Handler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	t := new(model.ProfileEntity)
	if err := c.Bind(t); err != nil {
		return c.JSON(utils.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	err := h.service.Update(ctx, c.Param("id"), t)
	if err != nil {
		return c.JSON(utils.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}
